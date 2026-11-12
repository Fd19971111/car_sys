package v1

import "github.com/gogf/gf/v2/frame/g"

// ReportSummary 报表核心指标（金额单位：元）。
type ReportSummary struct {
	Income        float64 `json:"income"         dc:"收入合计"`
	Expense       float64 `json:"expense"        dc:"支出合计"`
	Net           float64 `json:"net"            dc:"结余（收入-支出）"`
	NetRate       float64 `json:"netRate"        dc:"结余率（%）"`
	TotalBalance  float64 `json:"totalBalance"   dc:"资金余额合计"`
	CashBalance   float64 `json:"cashBalance"    dc:"现金余额"`
	BankBalance   float64 `json:"bankBalance"    dc:"银行余额"`
	Receivable    float64 `json:"receivable"     dc:"应收账款余额"`
	Payable       float64 `json:"payable"        dc:"应付账款余额"`
	RideRevenue   float64 `json:"rideRevenue"    dc:"所选年份完单车费"`
	PassengerPaid float64 `json:"passengerPaid"  dc:"所选年份关联行程收款"`
	YearReceived  float64 `json:"yearReceived"   dc:"所选年份所有收款"`
	UnpaidFare    float64 `json:"unpaidFare"     dc:"完单待收车费余额"`
	DriverPayable float64 `json:"driverPayable"  dc:"司机待结余额"`
	DriverSettled float64 `json:"driverSettled"  dc:"所选年份司机已结"`
	DriverExpense float64 `json:"driverExpense"  dc:"所选年份司机应得"`
}

// ReportMonthlyItem 月度收支。
type ReportMonthlyItem struct {
	Month   int     `json:"month"   dc:"月份"`
	Label   string  `json:"label"   dc:"月份标签"`
	Income  float64 `json:"income"  dc:"收入"`
	Expense float64 `json:"expense" dc:"支出"`
	Net     float64 `json:"net"     dc:"结余"`
}

// ReportAccountAmount 收入 / 支出按科目汇总。
type ReportAccountAmount struct {
	Code    string  `json:"code"    dc:"科目编码"`
	Name    string  `json:"name"    dc:"科目名称"`
	Amount  float64 `json:"amount"  dc:"金额"`
	Percent float64 `json:"percent" dc:"占同类合计百分比"`
}

// ReportFundBalance 资金账户余额。
type ReportFundBalance struct {
	Id           int64   `json:"id"           dc:"资金账户ID"`
	Type         string  `json:"type"         dc:"账户类型 CASH/BANK"`
	Name         string  `json:"name"         dc:"账户名称"`
	BankName     string  `json:"bankName"     dc:"开户银行"`
	AccountLast4 string  `json:"accountLast4" dc:"卡号末四位"`
	AccountCode  string  `json:"accountCode"  dc:"关联科目编码"`
	AccountName  string  `json:"accountName"  dc:"关联科目名称"`
	Balance      float64 `json:"balance"      dc:"当前余额"`
}

// ReportAccountBalance 科目余额。
type ReportAccountBalance struct {
	Code         string  `json:"code"         dc:"科目编码"`
	Name         string  `json:"name"         dc:"科目名称"`
	CategoryCode string  `json:"categoryCode" dc:"科目类别"`
	Balance      float64 `json:"balance"      dc:"余额"`
}

// ReportFlowItem 资金流水。
type ReportFlowItem struct {
	Id              int64   `json:"id"              dc:"流水ID"`
	OccurredAt      string  `json:"occurredAt"      dc:"发生时间 YYYY-MM-DD HH:mm:ss"`
	DirectionCode   string  `json:"directionCode"   dc:"INCOME 收入 / EXPENSE 支出"`
	Amount          float64 `json:"amount"          dc:"金额"`
	Remark          string  `json:"remark"          dc:"备注"`
	AccountName     string  `json:"accountName"     dc:"对方会计科目名称"`
	FundAccountName string  `json:"fundAccountName" dc:"资金账户名称"`
}

// ReportOverviewReq 查询财务报表总览。
type ReportOverviewReq struct {
	g.Meta `path:"/fin-report/overview" method:"get" tags:"财务报表" summary:"财务报表总览"`
	Year   int `json:"year" dc:"统计年份，默认当前年份"`
}

// ReportOverviewRes 财务报表总览。
type ReportOverviewRes struct {
	Initialized     bool                    `json:"initialized"     dc:"财务是否已初始化"`
	CurrentPeriod   string                  `json:"currentPeriod"   dc:"当前会计期间 YYYY-MM"`
	PeriodStatus    string                  `json:"periodStatus"    dc:"当前期间状态 OPEN/CLOSED/FUTURE"`
	Year            int                     `json:"year"`
	Summary         ReportSummary           `json:"summary"`
	Monthly         []*ReportMonthlyItem    `json:"monthly"`
	IncomeItems     []*ReportAccountAmount  `json:"incomeItems"`
	ExpenseItems    []*ReportAccountAmount  `json:"expenseItems"`
	FundAccounts    []*ReportFundBalance    `json:"fundAccounts"`
	AccountBalances []*ReportAccountBalance `json:"accountBalances"`
	RecentFlows     []*ReportFlowItem       `json:"recentFlows"`
}

// ReportFlowListReq 查询资金流水明细。
type ReportFlowListReq struct {
	g.Meta        `path:"/fin-report/flows" method:"get" tags:"财务报表" summary:"资金流水明细"`
	Year          int    `json:"year"          dc:"年份，默认当前年份"`
	Month         int    `json:"month"         dc:"月份 1-12，0 表示全年"`
	DirectionCode string `json:"directionCode" dc:"INCOME / EXPENSE，空为全部"`
	FundAccountId int64  `json:"fundAccountId" dc:"资金账户ID，0 为全部"`
	Keyword       string `json:"keyword"       dc:"备注 / 科目 / 账户关键字"`
	Current       int    `json:"current"       dc:"页码，从1开始"`
	Size          int    `json:"size"          dc:"每页条数"`
}

// ReportFlowSummary 当前筛选条件下的收支合计。
type ReportFlowSummary struct {
	Income  float64 `json:"income"  dc:"收入合计"`
	Expense float64 `json:"expense" dc:"支出合计"`
	Net     float64 `json:"net"     dc:"净额"`
}

// ReportFlowListRes 资金流水列表。
type ReportFlowListRes struct {
	Records []*ReportFlowItem `json:"records" dc:"资金流水"`
	Total   int               `json:"total"   dc:"总条数"`
	Summary ReportFlowSummary `json:"summary" dc:"收支合计"`
}
