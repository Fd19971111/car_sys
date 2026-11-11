package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 凭证分页列表 =====================

// VoucherListReq 凭证分页列表请求
type VoucherListReq struct {
	g.Meta `path:"/fin-voucher/list" method:"get" tags:"财务凭证" summary:"凭证分页列表"`

	Current    int    `json:"current"    d:"1"  dc:"当前页码"`
	Size       int    `json:"size"       d:"10" dc:"每页条数"`
	PeriodCode string `json:"periodCode" dc:"会计期间 YYYY-MM，空为全部"`
	StatusCode string `json:"statusCode" dc:"凭证状态 POSTED/VOID，空为全部"`
	Keyword    string `json:"keyword"    dc:"凭证号或摘要关键字"`
}

// VoucherListItem 凭证列表行
type VoucherListItem struct {
	Id              int64   `json:"id"              dc:"凭证ID"`
	VoucherNo       string  `json:"voucherNo"       dc:"凭证编号"`
	VoucherDate     string  `json:"voucherDate"     dc:"凭证日期"`
	PeriodCode      string  `json:"periodCode"      dc:"会计期间"`
	VoucherTypeCode string  `json:"voucherTypeCode" dc:"凭证类型：GENERAL/RECEIPT/PAYMENT/TRANSFER"`
	Summary         string  `json:"summary"         dc:"凭证摘要"`
	StatusCode      string  `json:"statusCode"      dc:"凭证状态：DRAFT/POSTED/VOID"`
	SourceType      string  `json:"sourceType"      dc:"来源类型：PROJECT_PAYMENT/PROJECT_PAYOUT/MANUAL"`
	SourceId        int64   `json:"sourceId"        dc:"来源业务记录ID"`
	DebitTotal      float64 `json:"debitTotal"      dc:"借方合计"`
	CreditTotal     float64 `json:"creditTotal"     dc:"贷方合计"`
	CreatedAt       string  `json:"createdAt"       dc:"创建时间"`
}

// VoucherListRes 凭证分页列表响应
type VoucherListRes struct {
	Current int                `json:"current" dc:"当前页码"`
	Size    int                `json:"size"    dc:"每页条数"`
	Total   int                `json:"total"   dc:"总条数"`
	Records []*VoucherListItem `json:"records" dc:"凭证列表"`
}

// ===================== 凭证详情 =====================

// VoucherDetailReq 凭证详情请求
type VoucherDetailReq struct {
	g.Meta `path:"/fin-voucher/detail" method:"get" tags:"财务凭证" summary:"凭证详情"`

	Id int64 `json:"id" v:"required#凭证ID不能为空" dc:"凭证ID"`
}

// VoucherEntryItem 凭证分录
type VoucherEntryItem struct {
	LineNo        int     `json:"lineNo"        dc:"行号"`
	AccountCode   string  `json:"accountCode"   dc:"科目编码"`
	AccountName   string  `json:"accountName"   dc:"科目名称"`
	Summary       string  `json:"summary"       dc:"分录摘要"`
	DirectionCode string  `json:"directionCode" dc:"借贷方向 DEBIT/CREDIT"`
	Amount        float64 `json:"amount"        dc:"分录金额"`
}

// VoucherDetailRes 凭证详情响应
type VoucherDetailRes struct {
	Id              int64               `json:"id"              dc:"凭证ID"`
	VoucherNo       string              `json:"voucherNo"       dc:"凭证编号"`
	VoucherDate     string              `json:"voucherDate"     dc:"凭证日期"`
	PeriodCode      string              `json:"periodCode"      dc:"会计期间"`
	VoucherTypeCode string              `json:"voucherTypeCode" dc:"凭证类型"`
	Summary         string              `json:"summary"         dc:"凭证摘要"`
	StatusCode      string              `json:"statusCode"      dc:"凭证状态"`
	SourceType      string              `json:"sourceType"      dc:"来源类型"`
	SourceId        int64               `json:"sourceId"        dc:"来源业务记录ID"`
	CreatedAt       string              `json:"createdAt"       dc:"创建时间"`
	DebitTotal      float64             `json:"debitTotal"      dc:"借方合计"`
	CreditTotal     float64             `json:"creditTotal"     dc:"贷方合计"`
	Entries         []*VoucherEntryItem `json:"entries"         dc:"凭证分录"`
}

// ===================== 会计期间与月结 =====================

// PeriodListReq 会计期间列表请求
type PeriodListReq struct {
	g.Meta `path:"/fin-voucher/periods" method:"get" tags:"财务凭证" summary:"会计期间列表"`
}

// PeriodItem 会计期间
type PeriodItem struct {
	PeriodCode   string  `json:"periodCode"    dc:"会计期间 YYYY-MM"`
	PeriodYear   int     `json:"periodYear"    dc:"会计年度"`
	PeriodMonth  int     `json:"periodMonth"   dc:"会计月份"`
	StatusCode   string  `json:"statusCode"    dc:"期间状态 OPEN/CLOSED/FUTURE"`
	VoucherCount int     `json:"voucherCount"  dc:"已记账凭证张数"`
	DebitTotal   float64 `json:"debitTotal"    dc:"借方发生额"`
	CreditTotal  float64 `json:"creditTotal"   dc:"贷方发生额"`
	CanClose     bool    `json:"canClose"      dc:"是否可结账"`
	CanReopen    bool    `json:"canReopen"     dc:"是否可反结账"`
	BeforeStart  bool    `json:"beforeStart"   dc:"是否早于财务启用期间"`
}

// PeriodListRes 会计期间列表响应
type PeriodListRes struct {
	Initialized   bool          `json:"initialized"   dc:"财务是否已初始化"`
	CurrentPeriod string        `json:"currentPeriod" dc:"当前自然期间 YYYY-MM"`
	StartPeriod   string        `json:"startPeriod"   dc:"财务启用期间 YYYY-MM"`
	Records       []*PeriodItem `json:"records"       dc:"会计期间列表"`
}

// ClosePeriodReq 月结请求
type ClosePeriodReq struct {
	g.Meta `path:"/fin-voucher/close" method:"post" tags:"财务凭证" summary:"会计期间结账"`

	PeriodCode string `json:"periodCode" v:"required#请选择要结账的期间" dc:"会计期间 YYYY-MM"`
}

// ClosePeriodRes 月结响应
type ClosePeriodRes struct {
	PeriodCode string `json:"periodCode" dc:"会计期间"`
	StatusCode string `json:"statusCode" dc:"结账后状态"`
}

// ReopenPeriodReq 反结账请求
type ReopenPeriodReq struct {
	g.Meta `path:"/fin-voucher/reopen" method:"post" tags:"财务凭证" summary:"会计期间反结账"`

	PeriodCode string `json:"periodCode" v:"required#请选择要反结账的期间" dc:"会计期间 YYYY-MM"`
}

// ReopenPeriodRes 反结账响应
type ReopenPeriodRes struct {
	PeriodCode string `json:"periodCode" dc:"会计期间"`
	StatusCode string `json:"statusCode" dc:"反结账后状态"`
}

// ===================== 凭证导出 =====================

// VoucherExportReq 凭证导出请求（按分录展开，便于会计做账）
type VoucherExportReq struct {
	g.Meta `path:"/fin-voucher/export" method:"get" tags:"财务凭证" summary:"凭证导出"`

	PeriodCode string `json:"periodCode" dc:"会计期间 YYYY-MM，空为全部"`
	StatusCode string `json:"statusCode" dc:"凭证状态 POSTED/VOID，空为全部"`
	Keyword    string `json:"keyword"    dc:"凭证号或摘要关键字"`
}

// VoucherExportRow 导出行（一张凭证的每一条分录一行）
type VoucherExportRow struct {
	VoucherNo       string  `json:"voucherNo"       dc:"凭证编号"`
	VoucherDate     string  `json:"voucherDate"     dc:"凭证日期"`
	PeriodCode      string  `json:"periodCode"      dc:"会计期间"`
	VoucherTypeCode string  `json:"voucherTypeCode" dc:"凭证类型"`
	VoucherSummary  string  `json:"voucherSummary"  dc:"凭证摘要"`
	StatusCode      string  `json:"statusCode"      dc:"凭证状态"`
	SourceType      string  `json:"sourceType"      dc:"来源类型"`
	SourceId        int64   `json:"sourceId"        dc:"来源业务记录ID"`
	LineNo          int     `json:"lineNo"          dc:"分录行号"`
	AccountCode     string  `json:"accountCode"     dc:"科目编码"`
	AccountName     string  `json:"accountName"     dc:"科目名称"`
	EntrySummary    string  `json:"entrySummary"    dc:"分录摘要"`
	DirectionCode   string  `json:"directionCode"   dc:"借贷方向"`
	Amount          float64 `json:"amount"          dc:"金额"`
}

// VoucherExportRes 凭证导出响应
type VoucherExportRes struct {
	Records []*VoucherExportRow `json:"records" dc:"凭证分录明细"`
}
