package fin_report

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Summary 报表核心指标，金额单位：元。
type Summary struct {
	Income         float64
	Expense        float64
	Net            float64
	NetRate        float64
	TotalBalance   float64
	CashBalance    float64
	BankBalance    float64
	Receivable     float64
	Payable        float64
	ContractAmount float64
	ReceivedAmount float64
	YearReceived   float64
	UnpaidAmount   float64
	ReceivedRate   float64
}

// Monthly 月度收支。
type Monthly struct {
	Month   int
	Label   string
	Income  float64
	Expense float64
	Net     float64
}

// AccountAmount 收入 / 支出按科目汇总。
type AccountAmount struct {
	Code    string
	Name    string
	Amount  float64
	Percent float64
}

// FundBalance 资金账户余额。
type FundBalance struct {
	Id           int64
	Type         string
	Name         string
	BankName     string
	AccountLast4 string
	AccountCode  string
	AccountName  string
	Balance      float64
}

// AccountBalance 科目余额。
type AccountBalance struct {
	Code         string
	Name         string
	CategoryCode string
	Balance      float64
}

// Flow 资金流水。
type Flow struct {
	Id              int64
	OccurredAt      string
	DirectionCode   string
	Amount          float64
	Remark          string
	AccountName     string
	FundAccountName string
}

// Overview 报表总览。
type Overview struct {
	Initialized     bool
	CurrentPeriod   string
	PeriodStatus    string
	Year            int
	Summary         Summary
	Monthly         []*Monthly
	IncomeItems     []*AccountAmount
	ExpenseItems    []*AccountAmount
	FundAccounts    []*FundBalance
	AccountBalances []*AccountBalance
	RecentFlows     []*Flow
}

// FlowFilter 资金流水筛选条件。
type FlowFilter struct {
	Year          int
	Month         int
	DirectionCode string
	FundAccountId int64
	Keyword       string
	Current       int
	Size          int
}

// FlowSummary 筛选范围内的收支合计。
type FlowSummary struct {
	Income  float64
	Expense float64
	Net     float64
}

// FlowList 资金流水分页结果。
type FlowList struct {
	Records []*Flow
	Total   int
	Summary FlowSummary
}

// Service 财务报表服务。
type Service struct{}

// New 创建财务报表服务。
func New() *Service { return &Service{} }

// Overview 汇总指定年份的经营报表数据。
func (s *Service) Overview(ctx context.Context, year int) (*Overview, error) {
	now := time.Now()
	if year < 2000 || year > 2100 {
		year = now.Year()
	}

	overview := &Overview{Year: year}

	if value, err := g.DB().GetValue(ctx, "SELECT is_initialized FROM fin_setup_config WHERE id=1"); err == nil {
		overview.Initialized = value.Int() == 1
	}

	overview.CurrentPeriod = now.Format("2006-01")
	if value, err := g.DB().GetValue(ctx, "SELECT status_code FROM fin_accounting_period WHERE period_code=?", overview.CurrentPeriod); err == nil {
		overview.PeriodStatus = value.String()
	}

	if err := s.fillMonthly(ctx, overview); err != nil {
		return nil, err
	}
	if err := s.fillBreakdown(ctx, overview); err != nil {
		return nil, err
	}
	if err := s.fillFundAccounts(ctx, overview); err != nil {
		return nil, err
	}
	if err := s.fillAccountBalances(ctx, overview); err != nil {
		return nil, err
	}
	if err := s.fillReceivable(ctx, overview); err != nil {
		return nil, err
	}
	if err := s.fillContract(ctx, overview); err != nil {
		return nil, err
	}
	if err := s.fillRecentFlows(ctx, overview); err != nil {
		return nil, err
	}

	s.balanceSummary(overview)

	return overview, nil
}

// fillMonthly 按月份统计收入与支出。
func (s *Service) fillMonthly(ctx context.Context, overview *Overview) error {
	records, err := g.DB().GetAll(ctx, `SELECT
			MONTH(v.voucher_date) AS month_no,
			COALESCE(SUM(CASE WHEN a.category_code='REVENUE' AND e.direction_code='CREDIT' THEN e.amount ELSE 0 END),0) AS revenue_credit,
			COALESCE(SUM(CASE WHEN a.category_code='REVENUE' AND e.direction_code='DEBIT' THEN e.amount ELSE 0 END),0) AS revenue_debit,
			COALESCE(SUM(CASE WHEN a.category_code='EXPENSE' AND e.direction_code='DEBIT' THEN e.amount ELSE 0 END),0) AS expense_debit,
			COALESCE(SUM(CASE WHEN a.category_code='EXPENSE' AND e.direction_code='CREDIT' THEN e.amount ELSE 0 END),0) AS expense_credit
		FROM fin_voucher_entry e
		INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.status_code='POSTED'
		INNER JOIN fin_account a ON a.id=e.account_id
		WHERE YEAR(v.voucher_date)=? AND a.category_code IN ('REVENUE','EXPENSE')
		GROUP BY MONTH(v.voucher_date)`, overview.Year)
	if err != nil {
		return gerror.Wrap(err, "统计月度收支失败")
	}

	incomeByMonth := make(map[int]float64, 12)
	expenseByMonth := make(map[int]float64, 12)
	for _, row := range records {
		month := row["month_no"].Int()
		incomeByMonth[month] = round2(row["revenue_credit"].Float64() - row["revenue_debit"].Float64())
		expenseByMonth[month] = round2(row["expense_debit"].Float64() - row["expense_credit"].Float64())
	}

	items := make([]*Monthly, 0, 12)
	for month := 1; month <= 12; month++ {
		income := incomeByMonth[month]
		expense := expenseByMonth[month]
		items = append(items, &Monthly{
			Month:   month,
			Label:   fmt.Sprintf("%d月", month),
			Income:  income,
			Expense: expense,
			Net:     round2(income - expense),
		})
	}
	overview.Monthly = items
	return nil
}

// fillBreakdown 按会计科目汇总收入与支出构成。
func (s *Service) fillBreakdown(ctx context.Context, overview *Overview) error {
	records, err := g.DB().GetAll(ctx, `SELECT
			a.account_code,a.account_name,a.category_code,
			SUM(CASE WHEN a.category_code='REVENUE'
				THEN (CASE WHEN e.direction_code='CREDIT' THEN e.amount ELSE -e.amount END)
				ELSE (CASE WHEN e.direction_code='DEBIT' THEN e.amount ELSE -e.amount END) END) AS net_amount
		FROM fin_voucher_entry e
		INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.status_code='POSTED'
		INNER JOIN fin_account a ON a.id=e.account_id
		WHERE YEAR(v.voucher_date)=? AND a.category_code IN ('REVENUE','EXPENSE')
		GROUP BY a.id,a.account_code,a.account_name,a.category_code
		HAVING net_amount <> 0
		ORDER BY net_amount DESC`, overview.Year)
	if err != nil {
		return gerror.Wrap(err, "统计收支构成失败")
	}

	incomes := make([]*AccountAmount, 0, len(records))
	expenses := make([]*AccountAmount, 0, len(records))
	incomeTotal, expenseTotal := 0.0, 0.0
	for _, row := range records {
		amount := round2(row["net_amount"].Float64())
		if amount == 0 {
			continue
		}
		item := &AccountAmount{
			Code:   row["account_code"].String(),
			Name:   row["account_name"].String(),
			Amount: amount,
		}
		if row["category_code"].String() == "REVENUE" {
			incomes = append(incomes, item)
			incomeTotal += amount
		} else {
			expenses = append(expenses, item)
			expenseTotal += amount
		}
	}

	overview.IncomeItems = withPercent(incomes, incomeTotal)
	overview.ExpenseItems = withPercent(expenses, expenseTotal)
	return nil
}

// fillFundAccounts 统计各资金账户当前余额（期初余额 + 全部已记账发生额）。
func (s *Service) fillFundAccounts(ctx context.Context, overview *Overview) error {
	startPeriod, err := s.startPeriod(ctx)
	if err != nil {
		return err
	}

	records, err := g.DB().GetAll(ctx, `SELECT
			f.id,f.account_type,f.account_name,f.bank_name,f.account_last4,
			a.account_code,a.account_name AS ledger_name,
			COALESCE(ob.opening_debit-ob.opening_credit,0) + COALESCE(mv.movement,0) AS balance
		FROM fin_fund_account f
		INNER JOIN fin_account a ON a.id=f.ledger_account_id
		LEFT JOIN fin_account_balance ob ON ob.account_id=a.id AND ob.period_code=?
		LEFT JOIN (
			SELECT e.account_id,
				SUM(CASE WHEN e.direction_code='DEBIT' THEN e.amount ELSE -e.amount END) AS movement
			FROM fin_voucher_entry e
			INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.status_code='POSTED'
			GROUP BY e.account_id
		) mv ON mv.account_id=a.id
		WHERE f.enabled=1
		ORDER BY f.account_type,f.id`, startPeriod)
	if err != nil {
		return gerror.Wrap(err, "统计资金账户余额失败")
	}

	items := make([]*FundBalance, 0, len(records))
	for _, row := range records {
		items = append(items, &FundBalance{
			Id:           row["id"].Int64(),
			Type:         row["account_type"].String(),
			Name:         row["account_name"].String(),
			BankName:     row["bank_name"].String(),
			AccountLast4: row["account_last4"].String(),
			AccountCode:  row["account_code"].String(),
			AccountName:  row["ledger_name"].String(),
			Balance:      round2(row["balance"].Float64()),
		})
	}
	overview.FundAccounts = items
	return nil
}

// fillAccountBalances 统计所有有余额的会计科目，供「科目余额」页面使用。
//
// 余额 = 启用期间的期初余额 + 全部已记账发生额，正数表示该科目的正常方向余额。
func (s *Service) fillAccountBalances(ctx context.Context, overview *Overview) error {
	startPeriod, err := s.startPeriod(ctx)
	if err != nil {
		return err
	}

	records, err := g.DB().GetAll(ctx, `SELECT
			a.account_code,a.account_name,a.category_code,
			COALESCE(CASE WHEN a.direction_code='DEBIT'
				THEN ob.opening_debit-ob.opening_credit
				ELSE ob.opening_credit-ob.opening_debit END,0) + COALESCE(mv.movement,0) AS balance
		FROM fin_account a
		LEFT JOIN fin_account_balance ob ON ob.account_id=a.id AND ob.period_code=?
		LEFT JOIN (
			SELECT e.account_id,
				SUM(CASE WHEN a2.direction_code='DEBIT' AND e.direction_code='DEBIT' THEN e.amount
					WHEN a2.direction_code='DEBIT' THEN -e.amount
					WHEN e.direction_code='CREDIT' THEN e.amount
					ELSE -e.amount END) AS movement
			FROM fin_voucher_entry e
			INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.status_code='POSTED'
			INNER JOIN fin_account a2 ON a2.id=e.account_id
			GROUP BY e.account_id
		) mv ON mv.account_id=a.id
		WHERE a.category_code IN ('ASSET','LIABILITY','EQUITY','REVENUE','EXPENSE')
		GROUP BY a.id,a.account_code,a.account_name,a.category_code,a.direction_code,a.sort_order,ob.opening_debit,ob.opening_credit,mv.movement
		HAVING balance <> 0
		ORDER BY FIELD(a.category_code,'ASSET','LIABILITY','EQUITY','REVENUE','EXPENSE'),a.sort_order,a.account_code`, startPeriod)
	if err != nil {
		return gerror.Wrap(err, "统计科目余额失败")
	}

	items := make([]*AccountBalance, 0, len(records))
	for _, row := range records {
		items = append(items, &AccountBalance{
			Code:         row["account_code"].String(),
			Name:         row["account_name"].String(),
			CategoryCode: row["category_code"].String(),
			Balance:      round2(row["balance"].Float64()),
		})
	}
	overview.AccountBalances = items
	return nil
}

// fillReceivable 统计应收 / 应付账款余额（含期初）。
func (s *Service) fillReceivable(ctx context.Context, overview *Overview) error {
	startPeriod, err := s.startPeriod(ctx)
	if err != nil {
		return err
	}

	records, err := g.DB().GetAll(ctx, `SELECT
			a.account_code,
			COALESCE(CASE WHEN a.direction_code='DEBIT'
				THEN ob.opening_debit-ob.opening_credit
				ELSE ob.opening_credit-ob.opening_debit END,0) + COALESCE(mv.movement,0) AS balance
		FROM fin_account a
		LEFT JOIN fin_account_balance ob ON ob.account_id=a.id AND ob.period_code=?
		LEFT JOIN (
			SELECT e.account_id,
				SUM(CASE WHEN a2.direction_code='DEBIT' AND e.direction_code='DEBIT' THEN e.amount
					WHEN a2.direction_code='DEBIT' THEN -e.amount
					WHEN e.direction_code='CREDIT' THEN e.amount
					ELSE -e.amount END) AS movement
			FROM fin_voucher_entry e
			INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.status_code='POSTED'
			INNER JOIN fin_account a2 ON a2.id=e.account_id
			GROUP BY e.account_id
		) mv ON mv.account_id=a.id
		WHERE a.account_code IN ('1122','2202')
		GROUP BY a.id,a.account_code,a.direction_code,ob.opening_debit,ob.opening_credit,mv.movement`, startPeriod)
	if err != nil {
		return gerror.Wrap(err, "统计应收应付失败")
	}

	for _, row := range records {
		balance := round2(row["balance"].Float64())
		if row["account_code"].String() == "1122" {
			overview.Summary.Receivable = balance
		} else {
			overview.Summary.Payable = balance
		}
	}
	return nil
}

// fillContract 统计项目合同与回款情况，让不懂财务的人也能看懂回款进度。
func (s *Service) fillContract(ctx context.Context, overview *Overview) error {
	contract, err := g.DB().GetValue(ctx, "SELECT COALESCE(SUM(amount),0) FROM project WHERE deleted_at IS NULL")
	if err != nil {
		return gerror.Wrap(err, "统计合同总额失败")
	}
	received, err := g.DB().GetValue(ctx, "SELECT COALESCE(SUM(amount),0) FROM project_payment WHERE deleted_at IS NULL")
	if err != nil {
		return gerror.Wrap(err, "统计收款总额失败")
	}
	yearReceived, err := g.DB().GetValue(ctx, "SELECT COALESCE(SUM(amount),0) FROM project_payment WHERE deleted_at IS NULL AND YEAR(payment_date)=?", overview.Year)
	if err != nil {
		return gerror.Wrap(err, "统计本年收款失败")
	}

	contractAmount := round2(contract.Float64())
	receivedAmount := round2(received.Float64())
	unpaid := round2(contractAmount - receivedAmount)
	if unpaid < 0 {
		unpaid = 0
	}

	overview.Summary.ContractAmount = contractAmount
	overview.Summary.ReceivedAmount = receivedAmount
	overview.Summary.YearReceived = round2(yearReceived.Float64())
	overview.Summary.UnpaidAmount = unpaid
	if contractAmount > 0 {
		overview.Summary.ReceivedRate = round1(receivedAmount / contractAmount * 100)
	}
	return nil
}

// fillRecentFlows 读取最近几条资金流水，用于总览页快速浏览。
func (s *Service) fillRecentFlows(ctx context.Context, overview *Overview) error {
	records, err := g.DB().GetAll(ctx, `SELECT
			t.id,t.direction_code,t.amount,t.remark,
			DATE_FORMAT(t.occurred_at,'%Y-%m-%d %H:%i:%s') AS occurred_at,
			COALESCE(a.account_name,'') AS account_name,
			COALESCE((SELECT f.account_name FROM fin_fund_account f WHERE f.ledger_account_id=t.account_id ORDER BY f.enabled DESC,f.id LIMIT 1),'') AS fund_account_name
		FROM fin_transaction t
		LEFT JOIN fin_account a ON a.id=t.account_id
		ORDER BY t.occurred_at DESC,t.id DESC
		LIMIT 6`)
	if err != nil {
		return gerror.Wrap(err, "查询最近资金流水失败")
	}
	overview.RecentFlows = mapFlows(records)
	return nil
}

// balanceSummary 汇总资金余额与收支指标。
func (s *Service) balanceSummary(overview *Overview) {
	summary := &overview.Summary
	for _, item := range overview.FundAccounts {
		if item.Type == "CASH" {
			summary.CashBalance += item.Balance
		} else {
			summary.BankBalance += item.Balance
		}
	}
	summary.CashBalance = round2(summary.CashBalance)
	summary.BankBalance = round2(summary.BankBalance)
	summary.TotalBalance = round2(summary.CashBalance + summary.BankBalance)

	for _, item := range overview.Monthly {
		summary.Income += item.Income
		summary.Expense += item.Expense
	}
	summary.Income = round2(summary.Income)
	summary.Expense = round2(summary.Expense)
	summary.Net = round2(summary.Income - summary.Expense)
	if summary.Income > 0 {
		summary.NetRate = round1(summary.Net / summary.Income * 100)
	}
}

// Flows 分页查询资金流水明细。
func (s *Service) Flows(ctx context.Context, filter *FlowFilter) (*FlowList, error) {
	now := time.Now()
	if filter == nil {
		filter = &FlowFilter{}
	}

	year := filter.Year
	if year < 2000 || year > 2100 {
		year = now.Year()
	}
	current := filter.Current
	if current < 1 {
		current = 1
	}
	size := filter.Size
	if size < 1 || size > 100 {
		size = 10
	}

	where := []string{"YEAR(t.occurred_at)=?"}
	args := []interface{}{year}
	if filter.Month >= 1 && filter.Month <= 12 {
		where = append(where, "MONTH(t.occurred_at)=?")
		args = append(args, filter.Month)
	}
	if direction := strings.ToUpper(strings.TrimSpace(filter.DirectionCode)); direction == "INCOME" || direction == "EXPENSE" {
		where = append(where, "t.direction_code=?")
		args = append(args, direction)
	}
	if filter.FundAccountId > 0 {
		where = append(where, "t.account_id=(SELECT ledger_account_id FROM fin_fund_account WHERE id=?)")
		args = append(args, filter.FundAccountId)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		where = append(where, "(t.remark LIKE ? OR a.account_name LIKE ? OR EXISTS(SELECT 1 FROM fin_fund_account f2 WHERE f2.ledger_account_id=t.account_id AND f2.account_name LIKE ?))")
		args = append(args, like, like, like)
	}
	whereSQL := strings.Join(where, " AND ")

	stat, err := g.DB().GetOne(ctx, fmt.Sprintf(`SELECT
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN t.direction_code='INCOME' THEN t.amount ELSE 0 END),0) AS income,
			COALESCE(SUM(CASE WHEN t.direction_code='EXPENSE' THEN t.amount ELSE 0 END),0) AS expense
		FROM fin_transaction t
		LEFT JOIN fin_account a ON a.id=t.account_id
		WHERE %s`, whereSQL), args...)
	if err != nil {
		return nil, gerror.Wrap(err, "统计资金流水失败")
	}

	result := &FlowList{Records: make([]*Flow, 0), Summary: FlowSummary{}}
	if stat != nil && !stat.IsEmpty() {
		result.Total = stat["total"].Int()
		result.Summary.Income = round2(stat["income"].Float64())
		result.Summary.Expense = round2(stat["expense"].Float64())
		result.Summary.Net = round2(result.Summary.Income - result.Summary.Expense)
	}

	if result.Total == 0 {
		return result, nil
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			t.id,t.direction_code,t.amount,t.remark,
			DATE_FORMAT(t.occurred_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS occurred_at,
			COALESCE(a.account_name,'') AS account_name,
			COALESCE((SELECT f.account_name FROM fin_fund_account f WHERE f.ledger_account_id=t.account_id ORDER BY f.enabled DESC,f.id LIMIT 1),'') AS fund_account_name
		FROM fin_transaction t
		LEFT JOIN fin_account a ON a.id=t.account_id
		WHERE %s
		ORDER BY t.occurred_at DESC,t.id DESC
		LIMIT ?,?`, whereSQL), queryArgs...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资金流水失败")
	}

	result.Records = mapFlows(records)
	return result, nil
}

// mapFlows 把查询结果转换为流水结构。
func mapFlows(records gdb.Result) []*Flow {
	items := make([]*Flow, 0, len(records))
	for _, row := range records {
		items = append(items, &Flow{
			Id:              row["id"].Int64(),
			OccurredAt:      row["occurred_at"].String(),
			DirectionCode:   row["direction_code"].String(),
			Amount:          round2(row["amount"].Float64()),
			Remark:          row["remark"].String(),
			AccountName:     row["account_name"].String(),
			FundAccountName: row["fund_account_name"].String(),
		})
	}
	return items
}

// startPeriod 返回财务启用期间（YYYY-MM），未初始化时返回空字符串。
func (s *Service) startPeriod(ctx context.Context) (string, error) {
	value, err := g.DB().GetValue(ctx, `SELECT CONCAT(LPAD(period_year,4,'0'),'-',LPAD(start_month,2,'0'))
		FROM fin_setup_config WHERE id=1 AND is_initialized=1`)
	if err != nil {
		return "", gerror.Wrap(err, "读取财务启用期间失败")
	}
	return value.String(), nil
}

// withPercent 计算每个科目占同类合计的百分比。
func withPercent(items []*AccountAmount, total float64) []*AccountAmount {
	if len(items) == 0 {
		return items
	}
	for _, item := range items {
		if total > 0 {
			item.Percent = round1(item.Amount / total * 100)
		}
	}
	return items
}

// round2 保留两位小数。
func round2(value float64) float64 { return math.Round(value*100) / 100 }

// round1 保留一位小数。
func round1(value float64) float64 { return math.Round(value*10) / 10 }
