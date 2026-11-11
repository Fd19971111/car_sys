// Package voucher 提供财务借贷记账的统一入口。
//
// 其他业务（收款、付款、进度、报表等）在写自己的业务数据时，只需要组装
// CreateInput 结构体并调用 Create / CreateWithTx，即可生成一张借贷平衡的
// 会计凭证，同时更新科目期间余额、生成资金流水，无需关心凭证编号、期间校验、
// 余额计算。
//
// 使用示例（与业务入库放在同一个事务里，保证业务数据与凭证同时成功）：
//
//	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
//		result, err := defaultDao.ProjectPayment.Ctx(ctx).Data(data).Insert()
//		if err != nil {
//			return err
//		}
//		sourceId, err := result.LastInsertId()
//		if err != nil {
//			return err
//		}
//		_, err = voucher.New().CreateWithTx(ctx, tx, &voucher.CreateInput{
//			VoucherDate: "2026-10-28",
//			VoucherType: "RECEIPT",
//			Summary:     "收到项目回款",
//			SourceType:  "PROJECT_PAYMENT",
//			SourceId:    sourceId,
//			CreatedBy:   userId,
//			Entries: []voucher.Entry{
//				{AccountCode: "1002", Direction: voucher.DirectionDebit, Amount: "10000.00", Summary: "收到项目回款"},
//				{AccountCode: "1122", Direction: voucher.DirectionCredit, Amount: "10000.00", Summary: "冲减应收账款"},
//			},
//		})
//		return err
//	})
package voucher

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// 借贷方向代码，与 fin_dictionary 中 ACCOUNT_DIRECTION 保持一致。
const (
	DirectionDebit  = "DEBIT"
	DirectionCredit = "CREDIT"
)

// 凭证状态代码，与 fin_dictionary 中 VOUCHER_STATUS 保持一致。
// 业务自动记账直接记入已记账状态，余额同步更新。
const voucherStatusPosted = "POSTED"

var amountPattern = regexp.MustCompile(`^\d{1,16}(\.\d{1,2})?$`)

// Entry 一条借贷分录。
type Entry struct {
	AccountCode string // 必填，会计科目编码，如 1002
	Direction   string // 必填，DEBIT（借）/ CREDIT（贷）
	Amount      string // 必填，正数金额，最多两位小数，如 "100.00"
	Summary     string // 可选，为空时取凭证摘要
}

// CreateInput 生成凭证所需参数。
type CreateInput struct {
	VoucherDate string  // 凭证日期 YYYY-MM-DD，为空取当前日期
	VoucherType string  // 凭证类型 GENERAL/RECEIPT/PAYMENT/TRANSFER，为空取 GENERAL
	Summary     string  // 凭证摘要，为空取第一条分录摘要
	SourceType  string  // 来源类型，如 PROJECT_PAYMENT，为空取 MANUAL
	SourceId    int64   // 来源业务记录 ID，手工录入为 0
	CreatedBy   uint    // 制单人 site_user.id
	Entries     []Entry // 必填，至少两条，且借贷合计必须相等
}

// Result 记账结果。
type Result struct {
	VoucherId  int64  // 凭证 ID
	VoucherNo  string // 凭证编号
	PeriodCode string // 所属会计期间 YYYY-MM
}

// Service 借贷记账服务。
type Service struct{}

// New 创建借贷记账服务。
func New() *Service { return &Service{} }

// Create 独立事务记账：内部自行开启事务，适合没有外层业务事务的场景。
func (s *Service) Create(ctx context.Context, in *CreateInput) (result *Result, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err = s.CreateWithTx(ctx, tx, in)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateWithTx 在调用方已有的事务中记账。
//
// 业务数据与凭证必须同时成功或同时回滚时，使用本方法并传入同一个 tx。
func (s *Service) CreateWithTx(ctx context.Context, tx gdb.TX, in *CreateInput) (*Result, error) {
	if in == nil {
		return nil, gerror.New("借贷参数不能为空")
	}

	companyId := currentCompanyId(ctx)
	if err := checkInitialized(tx, companyId); err != nil {
		return nil, err
	}

	voucherDate, err := parseVoucherDate(in.VoucherDate)
	if err != nil {
		return nil, err
	}
	periodCode := voucherDate.Format("2006-01")

	voucherType := strings.TrimSpace(in.VoucherType)
	if voucherType == "" {
		voucherType = "GENERAL"
	}
	if len(voucherType) > 20 {
		return nil, gerror.New("凭证类型代码长度不能超过20个字符")
	}

	sourceType := strings.TrimSpace(in.SourceType)
	if sourceType == "" {
		sourceType = "MANUAL"
	}
	if len(sourceType) > 20 {
		return nil, gerror.New("来源类型代码长度不能超过20个字符")
	}
	if in.SourceId < 0 {
		return nil, gerror.New("来源业务记录ID不能为负数")
	}

	if err = checkPeriod(tx, periodCode, companyId); err != nil {
		return nil, err
	}

	entries, err := prepareEntries(tx, in.Entries, strings.TrimSpace(in.Summary), companyId)
	if err != nil {
		return nil, err
	}

	summary := strings.TrimSpace(in.Summary)
	if summary == "" {
		summary = entries[0].Summary
		if summary == "" {
			summary = "系统自动生成"
		}
	}

	header := voucher{
		VoucherDate: voucherDate.Format("2006-01-02"),
		PeriodCode:  periodCode,
		VoucherType: voucherType,
		Summary:     summary,
		SourceType:  sourceType,
		SourceId:    in.SourceId,
		Status:      voucherStatusPosted,
		CreatedBy:   in.CreatedBy,
	}

	voucherId, voucherNo, err := insertVoucher(tx, header, entries, companyId)
	if err != nil {
		return nil, err
	}

	if err = insertFundTransactions(tx, header, entries, companyId); err != nil {
		return nil, err
	}

	// 统一从启用期间开始重算科目余额，保证逐月结转、与凭证完全一致
	if err = rebuildBalances(tx, companyId); err != nil {
		return nil, err
	}

	return &Result{VoucherId: voucherId, VoucherNo: voucherNo, PeriodCode: periodCode}, nil
}

// checkInitialized 校验财务基础配置已完成，未初始化时不允许记账。
func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func checkInitialized(tx gdb.TX, companyId uint) error {
	initialized, err := tx.GetValue("SELECT is_initialized FROM fin_setup_config WHERE company_id=?", companyId)
	if err != nil {
		return gerror.Wrap(err, "读取财务初始化状态失败，请确认已执行 fin_init.sql")
	}
	if initialized.Int() != 1 {
		return gerror.New("财务尚未初始化，请先完成财务基础配置")
	}
	return nil
}

// checkPeriod 校验会计期间：已结账期间禁止记账，未启用的期间在首次记账时自动启用。
//
// 期间按年度自动补齐，跨年继续记账不会再因为「期间不存在」而失败。
// 财务启用期间之前的月份永远不能记账，用来拦住系统上线前的历史数据。
func checkPeriod(tx gdb.TX, periodCode string, companyId uint) error {
	startPeriod, err := startPeriodCode(tx, companyId)
	if err != nil {
		return err
	}
	if periodCode < startPeriod {
		return gerror.Newf("会计期间 %s 早于财务启用期间 %s，不能记账", periodCode, startPeriod)
	}

	if err = ensurePeriodYear(tx, periodCode, companyId); err != nil {
		return err
	}

	periodStatus, err := tx.GetValue("SELECT status_code FROM fin_accounting_period WHERE company_id=? AND period_code=? FOR UPDATE", companyId, periodCode)
	if err != nil {
		return gerror.Wrap(err, "读取会计期间失败，请确认已执行 fin_init.sql")
	}
	current := periodStatus.String()

	switch current {
	case "":
		return gerror.Newf("会计期间 %s 不存在", periodCode)
	case "CLOSED":
		return gerror.Newf("会计期间 %s 当前为%s，不能记账", periodCode, periodStatusLabel(current))
	case "FUTURE":
		// 首次在未启用的期间记账时自动开启该期间
		if _, err = tx.Exec("UPDATE fin_accounting_period SET status_code='OPEN' WHERE company_id=? AND period_code=?", companyId, periodCode); err != nil {
			return gerror.Wrapf(err, "启用会计期间 %s 失败", periodCode)
		}
	}

	return nil
}

// preparedEntry 校验完成、并补齐了科目ID与名称的分录。
type preparedEntry struct {
	Entry
	AccountId   int64
	AccountName string
}

// prepareEntries 校验分录并批量加载分录用到的会计科目。
func prepareEntries(tx gdb.TX, input []Entry, defaultSummary string, companyId uint) ([]preparedEntry, error) {
	if len(input) < 2 {
		return nil, gerror.New("借贷分录至少需要两条")
	}

	entries := make([]preparedEntry, 0, len(input))
	codes := make([]string, 0, len(input))
	codeSet := make(map[string]struct{}, len(input))
	debitTotal := new(big.Int)
	creditTotal := new(big.Int)

	for index, item := range input {
		line := index + 1
		item.AccountCode = strings.TrimSpace(item.AccountCode)
		item.Direction = strings.ToUpper(strings.TrimSpace(item.Direction))
		item.Amount = strings.TrimSpace(item.Amount)
		item.Summary = strings.TrimSpace(item.Summary)
		if item.AccountCode == "" {
			return nil, gerror.Newf("第 %d 条分录入账科目不能为空", line)
		}
		if item.Direction != DirectionDebit && item.Direction != DirectionCredit {
			return nil, gerror.Newf("第 %d 条分录借贷方向不正确，只能是 DEBIT 或 CREDIT", line)
		}
		amountCents, normalized, ok := parseAmount(item.Amount)
		if !ok {
			return nil, gerror.Newf("第 %d 条分录金额格式不正确，最多两位小数", line)
		}
		if amountCents.Sign() <= 0 {
			return nil, gerror.Newf("第 %d 条分录金额必须大于 0", line)
		}
		if len(item.Summary) > 255 {
			return nil, gerror.Newf("第 %d 条分录摘要不能超过255个字符", line)
		}
		if item.Summary == "" {
			item.Summary = defaultSummary
		}
		item.Amount = normalized
		if item.Direction == DirectionDebit {
			debitTotal.Add(debitTotal, amountCents)
		} else {
			creditTotal.Add(creditTotal, amountCents)
		}
		if _, exists := codeSet[item.AccountCode]; !exists {
			codeSet[item.AccountCode] = struct{}{}
			codes = append(codes, item.AccountCode)
		}
		entries = append(entries, preparedEntry{Entry: item})
	}

	if debitTotal.Cmp(creditTotal) != 0 {
		return nil, gerror.Newf("借贷不平衡：借方合计 %s 元，贷方合计 %s 元", formatAmount(debitTotal), formatAmount(creditTotal))
	}

	accounts, err := loadAccounts(tx, codes, companyId)
	if err != nil {
		return nil, err
	}
	for index := range entries {
		acc := accounts[entries[index].AccountCode]
		entries[index].AccountId = acc.Id
		entries[index].AccountName = acc.Name
	}

	return entries, nil
}

// account 记账时需要用到的科目信息。
type account struct {
	Id   int64
	Code string
	Name string
}

// loadAccounts 按编码批量查询已启用的会计科目；有科目缺失时直接报错。
func loadAccounts(tx gdb.TX, codes []string, companyId uint) (map[string]account, error) {
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(codes)), ",")
	args := make([]interface{}, 0, len(codes)+1)
	args = append(args, companyId)
	for _, code := range codes {
		args = append(args, code)
	}

	rows, err := tx.GetAll(
		fmt.Sprintf("SELECT id,account_code,account_name FROM fin_account WHERE company_id=? AND enabled=1 AND account_code IN (%s)", placeholders),
		args...,
	)
	if err != nil {
		return nil, gerror.Wrap(err, "查询会计科目失败")
	}

	accounts := make(map[string]account, len(rows))
	for _, row := range rows {
		code := row["account_code"].String()
		accounts[code] = account{
			Id:   row["id"].Int64(),
			Code: code,
			Name: row["account_name"].String(),
		}
	}

	missing := make([]string, 0)
	for _, code := range codes {
		if _, exists := accounts[code]; !exists {
			missing = append(missing, code)
		}
	}
	if len(missing) > 0 {
		return nil, gerror.Newf("会计科目 %s 不存在或未启用，请先在财务基础配置中启用", strings.Join(missing, "、"))
	}

	return accounts, nil
}

// voucher 凭证主表落库参数。
type voucher struct {
	VoucherDate string
	PeriodCode  string
	VoucherType string
	Summary     string
	SourceType  string
	SourceId    int64
	Status      string
	CreatedBy   uint
}

// insertVoucher 生成凭证编号并写入凭证主表与分录明细。
//
// 凭证编号按月递增（如 记-2026-10-0001）。并发下依赖 fin_voucher 的
// voucher_no 唯一索引，撞号时重新取号重试。
func insertVoucher(tx gdb.TX, header voucher, entries []preparedEntry, companyId uint) (int64, string, error) {
	prefix := voucherPrefix(header.VoucherType)
	periodNumber := strings.ReplaceAll(header.PeriodCode, "-", "")

	var (
		voucherId int64
		voucherNo string
	)
	for attempt := 0; attempt < 5; attempt++ {
		seq, err := nextSequence(tx, header.PeriodCode, companyId)
		if err != nil {
			return 0, "", err
		}
		voucherNo = fmt.Sprintf("%s-%s-%04d", prefix, periodNumber, seq)

		result, err := tx.Exec(`INSERT INTO fin_voucher
			(company_id,voucher_no,voucher_date,period_code,voucher_type_code,summary,source_type,source_id,status_code,created_by)
			VALUES(?,?,?,?,?,?,?,?,?,?)`,
			companyId, voucherNo, header.VoucherDate, header.PeriodCode, header.VoucherType, header.Summary,
			header.SourceType, header.SourceId, header.Status, header.CreatedBy)
		if err != nil {
			if isDuplicateError(err) && attempt < 4 {
				continue
			}
			return 0, "", gerror.Wrap(err, "保存会计凭证失败")
		}

		voucherId, err = result.LastInsertId()
		if err != nil {
			return 0, "", gerror.Wrap(err, "读取会计凭证ID失败")
		}
		break
	}
	if voucherId == 0 {
		return 0, "", gerror.New("生成凭证编号失败，请稍后重试")
	}

	for index, entry := range entries {
		if _, err := tx.Exec(`INSERT INTO fin_voucher_entry
			(company_id,voucher_id,line_no,account_id,account_code,account_name,summary,direction_code,amount)
			VALUES(?,?,?,?,?,?,?,?,?)`,
			companyId, voucherId, index+1, entry.AccountId, entry.AccountCode, entry.AccountName,
			entry.Summary, entry.Direction, entry.Amount,
		); err != nil {
			return 0, "", gerror.Wrapf(err, "保存第 %d 条凭证分录失败", index+1)
		}
	}

	return voucherId, voucherNo, nil
}

// nextSequence 取指定期间已有凭证号的最大序号，下一个编号为最大值加一。
func nextSequence(tx gdb.TX, periodCode string, companyId uint) (int64, error) {
	value, err := tx.GetValue(`SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(voucher_no,'-',-1) AS UNSIGNED)),0)
		FROM fin_voucher WHERE company_id=? AND period_code=?`, companyId, periodCode)
	if err != nil {
		return 0, gerror.Wrap(err, "生成凭证编号失败")
	}
	return value.Int64() + 1, nil
}

// parseVoucherDate 解析会计日期，支持 YYYY-MM-DD 与 YYYY/MM/DD。
func parseVoucherDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local), nil
	}
	for _, layout := range []string{"2006-01-02", "2006/01/02"} {
		if date, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return date, nil
		}
	}
	return time.Time{}, gerror.New("凭证日期格式不正确，应为 YYYY-MM-DD")
}

// parseAmount 把金额字符串换算为分，并返回两位小数的标准格式。
func parseAmount(value string) (*big.Int, string, bool) {
	if !amountPattern.MatchString(value) {
		return nil, "", false
	}
	parts := strings.SplitN(value, ".", 2)
	fraction := "00"
	if len(parts) == 2 {
		fraction = parts[1] + strings.Repeat("0", 2-len(parts[1]))
	}
	cents, ok := new(big.Int).SetString(parts[0]+fraction, 10)
	if !ok {
		return nil, "", false
	}
	return cents, parts[0] + "." + fraction, true
}

// formatAmount 把分转换为两位小数的金额文本。
func formatAmount(cents *big.Int) string {
	whole, fraction := new(big.Int), new(big.Int)
	whole.QuoRem(cents, big.NewInt(100), fraction)
	return fmt.Sprintf("%s.%02d", whole.String(), fraction.Int64())
}

// voucherPrefix 凭证编号前缀：记 / 收 / 付 / 转。
func voucherPrefix(voucherType string) string {
	switch voucherType {
	case "RECEIPT":
		return "收"
	case "PAYMENT":
		return "付"
	case "TRANSFER":
		return "转"
	default:
		return "记"
	}
}

// periodStatusLabel 会计期间状态的中文描述。
func periodStatusLabel(status string) string {
	switch status {
	case "OPEN":
		return "启用中"
	case "CLOSED":
		return "已结账"
	default:
		return "未启用"
	}
}

// isDuplicateError 判断是否命中唯一索引冲突。
func isDuplicateError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate entry")
}
