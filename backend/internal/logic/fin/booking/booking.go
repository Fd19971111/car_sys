// Package booking 是业务单据与复式记账之间的适配层。
//
// 业务模块（收款 / 付款）只负责保存自己的单据，记账规则集中在这里：
//   - 收款：借「资金账户对应的资产科目（库存现金/银行存款）」，贷「主营业务收入」
//   - 付款：借「成本 / 费用科目」，贷「资金账户对应的资产科目」
//
// 这里采用收付实现制：收到钱确认收入，付出钱确认成本或费用。这样「报表」里的
// 收入、支出、结余和资金余额都能由业务单据自动汇总出来，不需要财务人员再手工
// 做账，也能保证业务单据与凭证永远在同一事务里成功或回滚。
//
// 记账科目不是写死的：优先使用下面约定的标准科目编码，若该科目未启用，则自动
// 回退到同类别中排序最靠前的启用科目，避免因为科目停用导致业务无法保存。
package booking

import (
	"context"
	"fmt"
	"math"
	"strings"

	"engineServ/internal/consts"
	"engineServ/internal/logic/fin/voucher"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	SourceTypePayment          = "RIDE_PAYMENT"
	SourceTypePayout           = "OPERATING_PAYOUT"
	SourceTypeRide             = "RIDE_ORDER"
	SourceTypeDriverSettlement = "DRIVER_SETTLEMENT"

	// 收款 / 付款凭证类型，对应 fin_dictionary.VOUCHER_TYPE。
	voucherTypeReceipt = "RECEIPT"
	voucherTypePayment = "PAYMENT"

	// 默认收入科目：主营业务收入。
	incomeAccountCode        = "5001"
	receivableAccountCode    = "1122"
	driverPayableAccountCode = "2202"
	// 默认成本科目：主营业务成本。
	costAccountCode = "5401"
	// 默认管理费用科目。
	adminExpenseAccountCode = "5602"
	// 默认销售费用科目。
	sellingExpenseAccountCode = "5601"

	// summaryMaxRunes 摘要最大字符数，避免超出凭证表的 255 字符限制。
	summaryMaxRunes = 60
)

// payoutAccountCodes 付款费用类型 → 默认借方科目编码。
//
// 材料 / 人工 / 施工 / 设计属于项目履约成本，默认计入主营业务成本；
// 管理费和销售费按性质分别计入对应的期间费用；其他支出按成本处理。
var payoutAccountCodes = map[string]string{
	"燃油费":    costAccountCode,
	"维修保养":   costAccountCode,
	"保险费":    costAccountCode,
	"路桥停车":   costAccountCode,
	"车辆租赁":   costAccountCode,
	"办公管理费":  adminExpenseAccountCode,
	"其他运营支出": costAccountCode,
	"材料款":    costAccountCode,
	"人工费":    costAccountCode,
	"施工费":    costAccountCode,
	"设计费":    costAccountCode,
	"管理费":    adminExpenseAccountCode,
	"销售费":    sellingExpenseAccountCode,
	"其他":     costAccountCode,
}

// Receipt 收款记账参数。
type Receipt struct {
	SourceId      int64  // 收款记录ID
	RideOrderId   int64  // 行程订单ID，大于0时冲减应收账款
	PaymentType   string // 款项类型，用于摘要
	Date          string // 收款日期 YYYY-MM-DD
	Amount        float64
	FundAccountId int64 // 收款资金账户（决定借记现金还是银行存款）
	CreatedBy     uint  // 制单人
}

// Payout 付款记账参数。
type Payout struct {
	SourceId      int64
	PayoutType    string // 费用类型
	Payee         string // 收款方
	Date          string
	Amount        float64
	FundAccountId int64
	CreatedBy     uint
}

// Result 记账结果。
type Result struct {
	VoucherId  int64
	VoucherNo  string
	PeriodCode string
}

// EntryView 凭证分录展示结构，详情页用它解释这张业务单据是怎么记账的。
type EntryView struct {
	Direction   string // DEBIT 借 / CREDIT 贷
	AccountName string
	Amount      float64
	Summary     string
}

// account 记账时用到的会计科目。
type account struct {
	Code string
	Name string
}

// fundAccount 资金账户及其关联的资产科目。
type fundAccount struct {
	Name       string
	LedgerId   int64
	LedgerCode string
	LedgerName string
}

// Service 业务记账服务。
type Service struct{}

// New 创建业务记账服务。
func New() *Service { return &Service{} }

// IsInitialized 判断财务基础配置是否已经完成。
//
// 未初始化时业务模块仍然可以保存单据，只是不生成凭证（voucher_id=0），
// 避免旧数据 / 未启用财务的公司无法登记收付款。
func (s *Service) IsInitialized(ctx context.Context, tx gdb.TX) (bool, error) {
	value, err := tx.GetValue("SELECT is_initialized FROM fin_setup_config WHERE company_id=?", currentCompanyId(ctx))
	if err != nil {
		return false, gerror.Wrap(err, "读取财务初始化状态失败，请确认已执行 fin_init.sql")
	}
	return value.Int() == 1, nil
}

// PostReceipt 生成收款凭证：借 资金账户科目 / 贷 收入科目。
func (s *Service) PostReceipt(ctx context.Context, tx gdb.TX, in *Receipt) (*Result, error) {
	fund, err := s.fundAccount(ctx, tx, in.FundAccountId)
	if err != nil {
		return nil, err
	}

	creditCategory, creditCode := "REVENUE", incomeAccountCode
	if in.RideOrderId > 0 {
		creditCategory, creditCode = "ASSET", receivableAccountCode
	}
	income, err := s.resolveAccount(ctx, tx, creditCategory, creditCode)
	if err != nil {
		return nil, err
	}

	amount := formatAmount(in.Amount)
	summary := receiptSummary(in)

	result, err := voucher.New().CreateWithTx(ctx, tx, &voucher.CreateInput{
		VoucherDate: in.Date,
		VoucherType: voucherTypeReceipt,
		Summary:     summary,
		SourceType:  SourceTypePayment,
		SourceId:    in.SourceId,
		CreatedBy:   in.CreatedBy,
		Entries: []voucher.Entry{
			{
				AccountCode: fund.LedgerCode,
				Direction:   voucher.DirectionDebit,
				Amount:      amount,
				Summary:     summary,
			},
			{
				AccountCode: income.Code,
				Direction:   voucher.DirectionCredit,
				Amount:      amount,
				Summary:     receiptCreditSummary(in),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &Result{VoucherId: result.VoucherId, VoucherNo: result.VoucherNo, PeriodCode: result.PeriodCode}, nil
}

// PostRideCompletion 在完单时确认车费收入、应收车费和司机应付成本。
func (s *Service) PostRideCompletion(ctx context.Context, tx gdb.TX, orderId int64, orderNumber, date string, fare, driverFare float64, createdBy uint) (*Result, error) {
	receivable, err := s.resolveAccount(ctx, tx, "ASSET", receivableAccountCode)
	if err != nil {
		return nil, err
	}
	income, err := s.resolveAccount(ctx, tx, "REVENUE", incomeAccountCode)
	if err != nil {
		return nil, err
	}
	cost, err := s.resolveAccount(ctx, tx, "EXPENSE", costAccountCode)
	if err != nil {
		return nil, err
	}
	payable, err := s.resolveAccount(ctx, tx, "LIABILITY", driverPayableAccountCode)
	if err != nil {
		return nil, err
	}

	fareAmount := formatAmount(fare)
	driverAmount := formatAmount(driverFare)
	summary := "行程完单 " + strings.TrimSpace(orderNumber)
	entries := []voucher.Entry{
		{AccountCode: receivable.Code, Direction: voucher.DirectionDebit, Amount: fareAmount, Summary: "确认行程应收车费"},
		{AccountCode: income.Code, Direction: voucher.DirectionCredit, Amount: fareAmount, Summary: "确认行程营业收入"},
	}
	if driverFare > 0 {
		entries = append(entries,
			voucher.Entry{AccountCode: cost.Code, Direction: voucher.DirectionDebit, Amount: driverAmount, Summary: "确认司机车费成本"},
			voucher.Entry{AccountCode: payable.Code, Direction: voucher.DirectionCredit, Amount: driverAmount, Summary: "确认应付司机车费"},
		)
	}
	result, err := voucher.New().CreateWithTx(ctx, tx, &voucher.CreateInput{
		VoucherDate: date,
		VoucherType: "GENERAL",
		Summary:     summary,
		SourceType:  SourceTypeRide,
		SourceId:    orderId,
		CreatedBy:   createdBy,
		Entries:     entries,
	})
	if err != nil {
		return nil, err
	}
	return &Result{VoucherId: result.VoucherId, VoucherNo: result.VoucherNo, PeriodCode: result.PeriodCode}, nil
}

// PostDriverSettlement 确认司机结算：借记应付司机，贷记资金账户。
func (s *Service) PostDriverSettlement(ctx context.Context, tx gdb.TX, sourceId, driverId int64, driverName, date string, amount float64, fundAccountId int64, createdBy uint) (*Result, error) {
	fund, err := s.fundAccount(ctx, tx, fundAccountId)
	if err != nil {
		return nil, err
	}
	if err = s.ensureFundBalance(ctx, tx, fund, amount); err != nil {
		return nil, err
	}
	payable, err := s.resolveAccount(ctx, tx, "LIABILITY", driverPayableAccountCode)
	if err != nil {
		return nil, err
	}
	amountText := formatAmount(amount)
	summary := "结算司机车费：" + strings.TrimSpace(driverName)
	result, err := voucher.New().CreateWithTx(ctx, tx, &voucher.CreateInput{
		VoucherDate: date,
		VoucherType: voucherTypePayment,
		Summary:     summary,
		SourceType:  SourceTypeDriverSettlement,
		SourceId:    sourceId,
		CreatedBy:   createdBy,
		Entries: []voucher.Entry{
			{AccountCode: payable.Code, Direction: voucher.DirectionDebit, Amount: amountText, Summary: "冲减应付司机车费"},
			{AccountCode: fund.LedgerCode, Direction: voucher.DirectionCredit, Amount: amountText, Summary: "支付司机车费"},
		},
	})
	if err != nil {
		return nil, err
	}
	return &Result{VoucherId: result.VoucherId, VoucherNo: result.VoucherNo, PeriodCode: result.PeriodCode}, nil
}

// PostPayout 生成付款凭证：借 成本/费用科目 / 贷 资金账户科目。
func (s *Service) PostPayout(ctx context.Context, tx gdb.TX, in *Payout) (*Result, error) {
	fund, err := s.fundAccount(ctx, tx, in.FundAccountId)
	if err != nil {
		return nil, err
	}

	// 钱不够不让付，避免把现金 / 银行付成负数
	if err = s.ensureFundBalance(ctx, tx, fund, in.Amount); err != nil {
		return nil, err
	}

	preferCode := payoutAccountCodes[strings.TrimSpace(in.PayoutType)]
	if preferCode == "" {
		preferCode = costAccountCode
	}
	expense, err := s.resolveAccount(ctx, tx, "EXPENSE", preferCode)
	if err != nil {
		return nil, err
	}

	amount := formatAmount(in.Amount)
	summary := payoutSummary(in)

	result, err := voucher.New().CreateWithTx(ctx, tx, &voucher.CreateInput{
		VoucherDate: in.Date,
		VoucherType: voucherTypePayment,
		Summary:     summary,
		SourceType:  SourceTypePayout,
		SourceId:    in.SourceId,
		CreatedBy:   in.CreatedBy,
		Entries: []voucher.Entry{
			{
				AccountCode: expense.Code,
				Direction:   voucher.DirectionDebit,
				Amount:      amount,
				Summary:     summary,
			},
			{
				AccountCode: fund.LedgerCode,
				Direction:   voucher.DirectionCredit,
				Amount:      amount,
				Summary:     "资金流出：" + fund.Name,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &Result{VoucherId: result.VoucherId, VoucherNo: result.VoucherNo, PeriodCode: result.PeriodCode}, nil
}

// Void 作废业务单据对应的凭证，并撤销余额与资金流水。
//
// voucherId 为 0（未记账）时直接返回，方便业务侧无脑调用。
func (s *Service) Void(ctx context.Context, tx gdb.TX, voucherId int64) error {
	return voucher.New().VoidWithTx(ctx, tx, voucherId)
}

// Entries 读取一张凭证的分录，用于详情页展示「借什么、贷什么」。
func (s *Service) Entries(ctx context.Context, voucherId int64) ([]*EntryView, error) {
	items := make([]*EntryView, 0)
	if voucherId <= 0 {
		return items, nil
	}

	records, err := g.DB().GetAll(ctx, `SELECT e.direction_code,e.account_name,e.amount,e.summary
		FROM fin_voucher_entry e
		INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.status_code='POSTED'
		WHERE e.voucher_id=? ORDER BY e.line_no`, voucherId)
	if err != nil {
		return nil, gerror.Wrap(err, "读取记账凭证分录失败")
	}

	for _, row := range records {
		items = append(items, &EntryView{
			Direction:   row["direction_code"].String(),
			AccountName: row["account_name"].String(),
			Amount:      roundAmount(row["amount"].Float64()),
			Summary:     row["summary"].String(),
		})
	}

	return items, nil
}

// VoucherNos 批量读取仍然有效的凭证号（已作废的凭证不返回）。
func (s *Service) VoucherNos(ctx context.Context, voucherIds []int64) (map[int64]string, error) {
	result := make(map[int64]string, len(voucherIds))
	if len(voucherIds) == 0 {
		return result, nil
	}

	args := make([]interface{}, 0, len(voucherIds))
	for _, id := range voucherIds {
		args = append(args, id)
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(voucherIds)), ",")

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(
		"SELECT id,voucher_no FROM fin_voucher WHERE status_code='POSTED' AND id IN (%s)", placeholders), args...)
	if err != nil {
		return nil, gerror.Wrap(err, "读取记账凭证编号失败")
	}

	for _, row := range records {
		result[row["id"].Int64()] = row["voucher_no"].String()
	}

	return result, nil
}

// ensureFundBalance 校验付款账户余额是否足够。
//
// 余额 = 财务启用期间的期初余额 + 全部已记账凭证的净发生额。
// 修改付款记录时，旧凭证已先被作废，因此这里算到的是「不含本笔」的余额。
func (s *Service) ensureFundBalance(ctx context.Context, tx gdb.TX, fund *fundAccount, amount float64) error {
	startPeriod, err := s.startPeriod(ctx, tx)
	if err != nil {
		return err
	}
	if startPeriod == "" {
		return nil
	}

	value, err := tx.GetValue(`SELECT
			COALESCE((SELECT opening_debit-opening_credit FROM fin_account_balance WHERE account_id=? AND period_code=?),0)
			+ COALESCE((SELECT SUM(CASE WHEN e.direction_code='DEBIT' THEN e.amount ELSE -e.amount END)
				FROM fin_voucher_entry e
				INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.status_code='POSTED'
				WHERE e.account_id=?),0)`, fund.LedgerId, startPeriod, fund.LedgerId)
	if err != nil {
		return gerror.Wrap(err, "查询付款账户余额失败")
	}

	balance := roundAmount(value.Float64())
	if balance+0.005 < amount {
		return gerror.Newf("「%s」当前余额 ¥%.2f，不足以支付 ¥%.2f，请先确认账户余额或选择其他付款账户",
			fund.Name, balance, amount)
	}

	return nil
}

// startPeriod 返回财务启用期间（YYYY-MM），未初始化时返回空字符串。
func (s *Service) startPeriod(ctx context.Context, tx gdb.TX) (string, error) {
	row, err := tx.GetOne("SELECT is_initialized,period_year,start_month FROM fin_setup_config WHERE id=1")
	if err != nil {
		return "", gerror.Wrap(err, "读取财务初始化状态失败")
	}
	if row == nil || row.IsEmpty() || row["is_initialized"].Int() != 1 {
		return "", nil
	}

	return fmt.Sprintf("%04d-%02d", row["period_year"].Int(), row["start_month"].Int()), nil
}

// fundAccount 校验资金账户，并返回它关联的资产科目编码。
func (s *Service) fundAccount(ctx context.Context, tx gdb.TX, id int64) (*fundAccount, error) {
	if id <= 0 {
		return nil, gerror.New("请选择资金账户，系统需要据此判断钱进了/出了哪个账户")
	}

	row, err := tx.GetOne(`SELECT f.account_name,f.enabled AS fund_enabled,
			a.id AS ledger_id,a.account_code,a.account_name AS ledger_name,a.enabled AS ledger_enabled
		FROM fin_fund_account f
		INNER JOIN fin_account a ON a.id=f.ledger_account_id
		WHERE f.id=? AND f.company_id=?`, id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询资金账户失败，请确认已执行 fin_init.sql")
	}
	if row == nil || row.IsEmpty() {
		return nil, gerror.New("所选资金账户不存在，请先在「财务基础资料」中配置现金 / 银行账户")
	}
	if !row["fund_enabled"].Bool() {
		return nil, gerror.Newf("资金账户「%s」已停用，请选择其他账户", row["account_name"].String())
	}
	if !row["ledger_enabled"].Bool() {
		return nil, gerror.Newf("资金账户「%s」关联的会计科目已停用，请先在「财务基础资料」中启用",
			row["account_name"].String())
	}

	return &fundAccount{
		Name:       row["account_name"].String(),
		LedgerId:   row["ledger_id"].Int64(),
		LedgerCode: row["account_code"].String(),
		LedgerName: row["ledger_name"].String(),
	}, nil
}

// currentCompanyId 读取当前请求所属公司，确保记账只使用本公司的资金账户与科目。
func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

// resolveAccount 按「优先编码 → 同类别首个启用科目」的顺序取记账科目。
func (s *Service) resolveAccount(ctx context.Context, tx gdb.TX, category, preferCode string) (*account, error) {
	if preferCode != "" {
		row, err := tx.GetOne(
			"SELECT account_code,account_name FROM fin_account WHERE account_code=? AND category_code=? AND company_id=? AND enabled=1",
			preferCode, category, currentCompanyId(ctx))
		if err != nil {
			return nil, gerror.Wrap(err, "查询会计科目失败")
		}
		if row != nil && !row.IsEmpty() {
			return &account{Code: row["account_code"].String(), Name: row["account_name"].String()}, nil
		}
	}

	row, err := tx.GetOne(`SELECT account_code,account_name FROM fin_account
		WHERE category_code=? AND company_id=? AND enabled=1 ORDER BY sort_order,account_code LIMIT 1`,
		category, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询会计科目失败")
	}
	if row == nil || row.IsEmpty() {
		return nil, gerror.Newf("请先在「财务基础资料」中启用一个%s类科目，用于自动记账", categoryLabel(category))
	}

	return &account{Code: row["account_code"].String(), Name: row["account_name"].String()}, nil
}

// receiptSummary 生成收款凭证摘要，例如「收到阳光美地整装定金」。
func receiptSummary(in *Receipt) string {
	paymentType := strings.TrimSpace(in.PaymentType)

	switch {
	case paymentType != "":
		return truncate("收到"+paymentType, summaryMaxRunes)
	default:
		return "收到款项"
	}
}

func receiptCreditSummary(in *Receipt) string {
	if in.RideOrderId > 0 {
		return "冲减行程应收车费"
	}
	return "确认营业收入"
}

// payoutSummary 生成付款凭证摘要，例如「支付恒信建材材料款」。
func payoutSummary(in *Payout) string {
	payee := strings.TrimSpace(in.Payee)
	payoutType := strings.TrimSpace(in.PayoutType)

	target := payee

	switch {
	case target != "" && payoutType != "":
		return truncate("支付"+target+payoutType, summaryMaxRunes)
	case target != "":
		return truncate("支付"+target+"款项", summaryMaxRunes)
	case payoutType != "":
		return truncate("支付"+payoutType, summaryMaxRunes)
	default:
		return "支付款项"
	}
}

// categoryLabel 科目类别的中文名称。
func categoryLabel(category string) string {
	switch category {
	case "ASSET":
		return "资产"
	case "LIABILITY":
		return "负债"
	case "EQUITY":
		return "权益"
	case "REVENUE":
		return "收入"
	case "EXPENSE":
		return "费用"
	default:
		return category
	}
}

// formatAmount 金额保留两位小数，并转成凭证要求的正数字符串。
func formatAmount(amount float64) string {
	cents := int64(math.Round(amount * 100))
	if cents < 0 {
		cents = -cents
	}
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

// roundAmount 金额保留两位小数（浮点展示用）。
func roundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}

// truncate 按字符（而非字节）截断摘要，避免中文超长。
func truncate(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}
