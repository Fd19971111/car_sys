package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 付款记录分页列表 =====================

// ListReq 付款记录分页列表请求
type ListReq struct {
	g.Meta `path:"/payout/list" method:"get" sm:"付款记录分页列表" tags:"Payout"`

	Current    int    `json:"current"    d:"1"  dc:"当前页码，默认 1"`
	Size       int    `json:"size"       d:"10" dc:"每页条数，默认 10"`
	Keyword    string `json:"keyword"          dc:"关键词：收款方 / 费用类型 / 备注"`
	PayoutType string `json:"payoutType"       dc:"费用类型：材料款/人工费/施工费/设计费/管理费/销售费/其他"`
}

// PayoutListItem 付款记录列表行
type PayoutListItem struct {
	Id              int64   `json:"id"              dc:"付款记录ID"`
	Amount          float64 `json:"amount"          dc:"付款金额"`
	PayoutType      string  `json:"payoutType"      dc:"费用类型"`
	PayoutDate      string  `json:"payoutDate"      dc:"付款日期"`
	VendorName      string  `json:"vendorName"      dc:"收款方名称"`
	FundAccountId   int64   `json:"fundAccountId"   dc:"付款账户ID"`
	FundAccountName string  `json:"fundAccountName" dc:"付款账户（现金 / 银行）"`
	VoucherNo       string  `json:"voucherNo"       dc:"自动生成的凭证号，空字符串表示未记账"`
	Remark          string  `json:"remark"          dc:"备注"`
	CreatedAt       string  `json:"createdAt"       dc:"创建时间"`
}

// ListRes 付款记录分页列表响应
type ListRes struct {
	Current int               `json:"current" dc:"当前页码"`
	Size    int               `json:"size"    dc:"每页条数"`
	Total   int               `json:"total"   dc:"总条数"`
	Records []*PayoutListItem `json:"records" dc:"付款记录列表"`
}

// ===================== 付款概览统计 =====================

// StatsReq 付款概览统计请求
type StatsReq struct {
	g.Meta `path:"/payout/stats" method:"get" sm:"付款概览统计" tags:"Payout"`
}

// StatsRes 付款概览统计响应
type StatsRes struct {
	TotalCount  int     `json:"totalCount"  dc:"付款笔数"`
	TotalAmount float64 `json:"totalAmount" dc:"付款总额"`
	MonthAmount float64 `json:"monthAmount" dc:"本月付款金额"`
	MonthCount  int     `json:"monthCount"  dc:"本月付款笔数"`
}

// ===================== 付款记录详情 =====================

// DetailReq 付款记录详情请求
type DetailReq struct {
	g.Meta `path:"/payout/detail" method:"get" sm:"付款记录详情" tags:"Payout"`

	Id int64 `json:"id" v:"required#付款记录ID不能为空" dc:"付款记录ID"`
}

// LedgerEntryItem 凭证分录，用于详情页展示「借什么、贷什么」。
type LedgerEntryItem struct {
	Direction   string  `json:"direction"   dc:"借贷方向：DEBIT 借 / CREDIT 贷"`
	AccountName string  `json:"accountName" dc:"会计科目名称"`
	Amount      float64 `json:"amount"      dc:"分录金额"`
	Summary     string  `json:"summary"     dc:"分录摘要"`
}

// DetailRes 付款记录详情响应
type DetailRes struct {
	Id              int64              `json:"id"              dc:"付款记录ID"`
	Amount          float64            `json:"amount"          dc:"本次付款金额"`
	PayoutType      string             `json:"payoutType"      dc:"费用类型"`
	PayoutDate      string             `json:"payoutDate"      dc:"付款日期"`
	VendorName      string             `json:"vendorName"      dc:"收款方名称"`
	FundAccountId   int64              `json:"fundAccountId"   dc:"付款账户ID"`
	FundAccountName string             `json:"fundAccountName" dc:"付款账户（现金 / 银行）"`
	VoucherId       int64              `json:"voucherId"       dc:"自动生成的凭证ID，0 表示未记账"`
	VoucherNo       string             `json:"voucherNo"       dc:"自动生成的凭证号"`
	LedgerEntries   []*LedgerEntryItem `json:"ledgerEntries"   dc:"凭证借贷分录"`
	Remark          string             `json:"remark"          dc:"备注"`
	CreatedAt       string             `json:"createdAt"       dc:"创建时间"`
}

// ===================== 新增 / 编辑 / 删除 =====================

// PayoutFormFields 付款记录表单公共字段
type PayoutFormFields struct {
	Amount        float64 `json:"amount"        v:"required|min:0.01#请输入付款金额|付款金额必须大于 0" dc:"付款金额（元）"`
	PayoutType    string  `json:"payoutType"    v:"required|max-length:50#请选择费用类型|费用类型最多 50 个字符" dc:"费用类型：材料款/人工费/施工费/设计费/管理费/销售费/其他"`
	PayoutDate    string  `json:"payoutDate"    v:"required#请选择付款日期"                                    dc:"付款日期（YYYY-MM-DD）"`
	VendorName    string  `json:"vendorName"    v:"required|max-length:100#请填写收款方|收款方最多 100 个字符" dc:"收款方名称：供应商 / 工人 / 公司"`
	FundAccountId int64   `json:"fundAccountId" dc:"付款账户ID（财务已初始化时必填，决定借记哪个资金科目）"`
	Remark        string  `json:"remark"        v:"max-length:255#备注最多 255 个字符"                        dc:"备注"`
}

// CreateReq 新增付款记录请求
type CreateReq struct {
	g.Meta `path:"/payout/create" method:"post" sm:"新增付款记录" tags:"Payout"`

	PayoutFormFields
}

// CreateRes 新增付款记录响应
type CreateRes struct {
	Id        int64  `json:"id"        dc:"新建付款记录ID"`
	VoucherNo string `json:"voucherNo" dc:"自动生成的凭证号，空字符串表示未记账"`
}

// UpdateReq 更新付款记录请求
type UpdateReq struct {
	g.Meta `path:"/payout/update" method:"post" sm:"更新付款记录" tags:"Payout"`

	Id int64 `json:"id" v:"required#付款记录ID不能为空" dc:"付款记录ID"`
	PayoutFormFields
}

// UpdateRes 更新付款记录响应
type UpdateRes struct {
	Id        int64  `json:"id"        dc:"付款记录ID"`
	VoucherNo string `json:"voucherNo" dc:"自动生成的凭证号，空字符串表示未记账"`
}

// DeleteReq 删除付款记录请求（后端软删除，同时作废自动生成的凭证）
type DeleteReq struct {
	g.Meta `path:"/payout/delete" method:"post" sm:"删除付款记录" tags:"Payout"`

	Id int64 `json:"id" v:"required#付款记录ID不能为空" dc:"付款记录ID"`
}

// DeleteRes 删除付款记录响应
type DeleteRes struct{}

// BatchDeleteReq 批量删除付款记录请求
type BatchDeleteReq struct {
	g.Meta `path:"/payout/batchDelete" method:"post" sm:"批量删除付款记录" tags:"Payout"`

	Ids []int64 `json:"ids" v:"required#请选择要删除的付款记录" dc:"付款记录ID集合"`
}

// BatchDeleteRes 批量删除付款记录响应
type BatchDeleteRes struct{}
