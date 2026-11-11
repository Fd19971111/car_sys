package v1

import "github.com/gogf/gf/v2/frame/g"

// ReportSummary 报表核心指标（金额单位：元）。
type ReportSummary struct {
	Income         float64 `json:"income"         dc:"收入合计"`
	Expense        float64 `json:"expense"        dc:"支出合计"`
	Net            float64 `json:"net"            dc:"结余（收入-支出）"`
	NetRate        float64 `json:"netRate"        dc:"结余率（%）"`
	TotalBalance   float64 `json:"totalBalance"   dc:"资金余额合计"`
	CashBalance    float64 `json:"cashBalance"    dc:"现金余额"`
	BankBalance    float64 `json:"bankBalance"    dc:"银行余额"`
	Receivable     float64 `json:"receivable"     dc:"应收账款余额"`
	Payable        float64 `json:"payable"        dc:"应付账款余额"`
	ContractAmount float64 `json:"contractAmount" dc:"合同总额"`
	ReceivedAmount float64 `json:"receivedAmount" dc:"累计回款"`
	YearReceived   float64 `json:"yearReceived"   dc:"所选年份回款"`
	UnpaidAmount   float64 `json:"unpaidAmount"   dc:"待回款"`
	ReceivedRate   float64 `json:"receivedRate"   dc:"回款率（%）"`
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

// ===================== 项目成本 / 毛利 =====================

// ProjectProfitCostItem 成本构成项。
type ProjectProfitCostItem struct {
	Type    string  `json:"type"    dc:"成本类型"`
	Amount  float64 `json:"amount"  dc:"金额"`
	Percent float64 `json:"percent" dc:"占成本合计百分比"`
}

// ProjectProfitItem 项目毛利行。
type ProjectProfitItem struct {
	ProjectId        int64                    `json:"projectId"        dc:"项目ID"`
	ProjectName      string                   `json:"projectName"      dc:"项目名称"`
	CustomerName     string                   `json:"customerName"     dc:"客户姓名"`
	Status           int                      `json:"status"           dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	StartDate        string                   `json:"startDate"        dc:"开工日期"`
	EndDate          string                   `json:"endDate"          dc:"预计完工日期"`
	ContractAmount   float64                  `json:"contractAmount"   dc:"合同金额"`
	ChangeNet        float64                  `json:"changeNet"        dc:"变更净额（增项-减项）"`
	SettlementAmount float64                  `json:"settlementAmount" dc:"结算金额 = 合同 + 变更净额"`
	ReceivedAmount   float64                  `json:"receivedAmount"   dc:"已收款"`
	UnpaidAmount     float64                  `json:"unpaidAmount"     dc:"待收款"`
	ReceivedRate     float64                  `json:"receivedRate"     dc:"回款率（%）"`
	CostTotal        float64                  `json:"costTotal"        dc:"成本合计（按付款归集）"`
	CostItems        []*ProjectProfitCostItem `json:"costItems"        dc:"成本构成"`
	MaterialPaid     float64                  `json:"materialPaid"     dc:"材料款付款"`
	MaterialUsed     float64                  `json:"materialUsed"     dc:"材料耗用（领用金额，参考）"`
	GrossProfit      float64                  `json:"grossProfit"      dc:"预计毛利 = 结算金额 - 成本合计"`
	GrossRate        float64                  `json:"grossRate"        dc:"毛利率（%）"`
	CashProfit       float64                  `json:"cashProfit"       dc:"现金结余 = 已收款 - 成本合计"`
}

// ProjectProfitSummary 项目毛利汇总。
type ProjectProfitSummary struct {
	ProjectCount      int     `json:"projectCount"      dc:"项目数"`
	SettlementTotal   float64 `json:"settlementTotal"   dc:"结算总额"`
	ReceivedTotal     float64 `json:"receivedTotal"     dc:"已收款合计"`
	UnpaidTotal       float64 `json:"unpaidTotal"       dc:"待收款合计"`
	CostTotal         float64 `json:"costTotal"         dc:"成本合计"`
	GrossProfit       float64 `json:"grossProfit"       dc:"预计毛利合计"`
	GrossRate         float64 `json:"grossRate"         dc:"整体毛利率（%）"`
	LossCount         int     `json:"lossCount"         dc:"亏损项目数"`
	CompanyCost       float64 `json:"companyCost"       dc:"公司公共支出（未关联项目，未摊入成本）"`
	MaterialUsedTotal float64 `json:"materialUsedTotal" dc:"材料耗用合计（参考）"`
}

// ProjectProfitListReq 项目毛利列表请求。
type ProjectProfitListReq struct {
	g.Meta `path:"/fin-report/project-profit/list" method:"get" tags:"财务报表" summary:"项目毛利列表"`

	Keyword  string `json:"keyword"  dc:"项目名称 / 客户姓名关键字"`
	Status   *int   `json:"status"   dc:"项目状态，不传为全部"`
	OnlyLoss bool   `json:"onlyLoss" dc:"只看亏损项目"`
}

// ProjectProfitListRes 项目毛利列表响应。
type ProjectProfitListRes struct {
	Records []*ProjectProfitItem  `json:"records" dc:"项目毛利"`
	Summary *ProjectProfitSummary `json:"summary" dc:"汇总"`
}

// ProjectProfitDetailReq 项目毛利明细请求。
type ProjectProfitDetailReq struct {
	g.Meta `path:"/fin-report/project-profit/detail" method:"get" tags:"财务报表" summary:"项目毛利明细"`

	ProjectId int64 `json:"projectId" v:"required#项目ID不能为空" dc:"项目ID"`
}

// ProjectProfitReceipt 项目收款明细。
type ProjectProfitReceipt struct {
	Id          int64   `json:"id"          dc:"收款ID"`
	PaymentType string  `json:"paymentType" dc:"款项类型"`
	Amount      float64 `json:"amount"      dc:"金额"`
	PaymentDate string  `json:"paymentDate" dc:"收款日期"`
	Remark      string  `json:"remark"      dc:"备注"`
}

// ProjectProfitPayout 项目付款明细（成本）。
type ProjectProfitPayout struct {
	Id         int64   `json:"id"         dc:"付款ID"`
	PayoutType string  `json:"payoutType" dc:"费用类型"`
	Amount     float64 `json:"amount"     dc:"金额"`
	PayoutDate string  `json:"payoutDate" dc:"付款日期"`
	VendorName string  `json:"vendorName" dc:"供应商名称"`
	Remark     string  `json:"remark"     dc:"备注"`
}

// ProjectProfitMaterial 项目材料领用明细（参考，不计入成本）。
type ProjectProfitMaterial struct {
	Id           int64   `json:"id"           dc:"领用ID"`
	MaterialName string  `json:"materialName" dc:"材料名称"`
	Unit         string  `json:"unit"         dc:"单位"`
	Quantity     float64 `json:"quantity"     dc:"领用数量"`
	UnitPrice    float64 `json:"unitPrice"    dc:"单价"`
	Amount       float64 `json:"amount"       dc:"金额"`
	UsedAt       string  `json:"usedAt"       dc:"领用日期"`
	Receiver     string  `json:"receiver"     dc:"领用人"`
}

// ProjectProfitChange 项目变更明细。
type ProjectProfitChange struct {
	Id         int64   `json:"id"         dc:"变更ID"`
	ChangeType string  `json:"changeType" dc:"ADD 增项 / REDUCE 减项"`
	Title      string  `json:"title"      dc:"变更事项"`
	Amount     float64 `json:"amount"     dc:"金额"`
	ChangeDate string  `json:"changeDate" dc:"变更日期"`
	Reason     string  `json:"reason"     dc:"变更原因"`
}

// ProjectProfitDetailRes 项目毛利明细响应。
type ProjectProfitDetailRes struct {
	ProjectProfitItem
	Receipts  []*ProjectProfitReceipt  `json:"receipts"  dc:"收款明细"`
	Payouts   []*ProjectProfitPayout   `json:"payouts"   dc:"付款明细（成本）"`
	Materials []*ProjectProfitMaterial `json:"materials" dc:"材料领用（参考）"`
	Changes   []*ProjectProfitChange   `json:"changes"   dc:"变更明细"`
}
