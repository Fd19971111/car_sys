package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 收款记录分页列表 =====================

// ListReq 收款记录分页列表请求
type ListReq struct {
	g.Meta `path:"/payment/list" method:"get" sm:"收款记录分页列表" tags:"Payment"`

	Current     int    `json:"current"     d:"1"  dc:"当前页码，默认 1"`
	Size        int    `json:"size"        d:"10" dc:"每页条数，默认 10"`
	Keyword     string `json:"keyword"           dc:"关键词：款项类型 / 备注"`
	PaymentType string `json:"paymentType"       dc:"款项类型：定金/首期款/中期款/尾款"`
}

// PaymentListItem 收款记录列表行
type PaymentListItem struct {
	Id              int64   `json:"id"              dc:"收款记录ID"`
	RideOrderId     int64   `json:"rideOrderId"     dc:"关联行程订单ID，0表示未关联"`
	Amount          float64 `json:"amount"          dc:"收款金额"`
	PaymentType     string  `json:"paymentType"     dc:"款项类型"`
	PaymentDate     string  `json:"paymentDate"     dc:"收款日期"`
	FundAccountId   int64   `json:"fundAccountId"   dc:"收款账户ID"`
	FundAccountName string  `json:"fundAccountName" dc:"收款账户（现金 / 银行）"`
	VoucherNo       string  `json:"voucherNo"       dc:"自动生成的凭证号，空字符串表示未记账"`
	Remark          string  `json:"remark"          dc:"备注"`
	CreatedAt       string  `json:"createdAt"       dc:"创建时间"`
}

// ListRes 收款记录分页列表响应
type ListRes struct {
	Current int                `json:"current" dc:"当前页码"`
	Size    int                `json:"size"    dc:"每页条数"`
	Total   int                `json:"total"   dc:"总条数"`
	Records []*PaymentListItem `json:"records" dc:"收款记录列表"`
}

// ===================== 收款概览统计 =====================

// StatsReq 收款概览统计请求
type StatsReq struct {
	g.Meta `path:"/payment/stats" method:"get" sm:"收款概览统计" tags:"Payment"`
}

// StatsRes 收款概览统计响应
type StatsRes struct {
	TotalCount  int     `json:"totalCount"  dc:"收款笔数"`
	TotalAmount float64 `json:"totalAmount" dc:"收款总额"`
	MonthAmount float64 `json:"monthAmount" dc:"本月收款金额"`
	MonthCount  int     `json:"monthCount"  dc:"本月收款笔数"`
}

// ===================== 收款记录详情 =====================

// DetailReq 收款记录详情请求
type DetailReq struct {
	g.Meta `path:"/payment/detail" method:"get" sm:"收款记录详情" tags:"Payment"`

	Id int64 `json:"id" v:"required#收款记录ID不能为空" dc:"收款记录ID"`
}

// LedgerEntryItem 凭证分录，用于详情页展示「借什么、贷什么」。
type LedgerEntryItem struct {
	Direction   string  `json:"direction"   dc:"借贷方向：DEBIT 借 / CREDIT 贷"`
	AccountName string  `json:"accountName" dc:"会计科目名称"`
	Amount      float64 `json:"amount"      dc:"分录金额"`
	Summary     string  `json:"summary"     dc:"分录摘要"`
}

// DetailRes 收款记录详情响应
type DetailRes struct {
	Id              int64              `json:"id"              dc:"收款记录ID"`
	RideOrderId     int64              `json:"rideOrderId"     dc:"关联行程订单ID，0表示未关联"`
	Amount          float64            `json:"amount"          dc:"本次收款金额"`
	PaymentType     string             `json:"paymentType"     dc:"款项类型"`
	PaymentDate     string             `json:"paymentDate"     dc:"收款日期"`
	FundAccountId   int64              `json:"fundAccountId"   dc:"收款账户ID"`
	FundAccountName string             `json:"fundAccountName" dc:"收款账户（现金 / 银行）"`
	VoucherId       int64              `json:"voucherId"       dc:"自动生成的凭证ID，0 表示未记账"`
	VoucherNo       string             `json:"voucherNo"       dc:"自动生成的凭证号"`
	LedgerEntries   []*LedgerEntryItem `json:"ledgerEntries"   dc:"凭证借贷分录"`
	Remark          string             `json:"remark"          dc:"备注"`
	CreatedAt       string             `json:"createdAt"       dc:"创建时间"`
}

// ===================== 新增 / 编辑 / 删除 =====================

// PaymentFormFields 收款记录表单公共字段
type PaymentFormFields struct {
	RideOrderId   int64   `json:"rideOrderId"   v:"min:0#订单ID不能小于0" dc:"关联行程订单ID，0表示通用收款"`
	Amount        float64 `json:"amount"        v:"required|min:0.01#请输入收款金额|收款金额必须大于 0" dc:"收款金额（元）"`
	PaymentType   string  `json:"paymentType"   v:"required|max-length:50#请选择款项类型|款项类型最多 50 个字符" dc:"款项类型：定金/首期款/中期款/尾款"`
	PaymentDate   string  `json:"paymentDate"   v:"required#请选择收款日期"                                    dc:"收款日期（YYYY-MM-DD）"`
	FundAccountId int64   `json:"fundAccountId" dc:"收款账户ID（财务已初始化时必填，决定贷记哪个资金科目）"`
	Remark        string  `json:"remark"        v:"max-length:255#备注最多 255 个字符"                        dc:"备注"`
}

// CreateReq 新增收款记录请求
type CreateReq struct {
	g.Meta `path:"/payment/create" method:"post" sm:"新增收款记录" tags:"Payment"`

	PaymentFormFields
}

// CreateRes 新增收款记录响应
type CreateRes struct {
	Id        int64  `json:"id"        dc:"新建收款记录ID"`
	VoucherNo string `json:"voucherNo" dc:"自动生成的凭证号，空字符串表示未记账"`
}

// UpdateReq 更新收款记录请求
type UpdateReq struct {
	g.Meta `path:"/payment/update" method:"post" sm:"更新收款记录" tags:"Payment"`

	Id int64 `json:"id" v:"required#收款记录ID不能为空" dc:"收款记录ID"`
	PaymentFormFields
}

// UpdateRes 更新收款记录响应
type UpdateRes struct {
	Id        int64  `json:"id"        dc:"收款记录ID"`
	VoucherNo string `json:"voucherNo" dc:"自动生成的凭证号，空字符串表示未记账"`
}

// DeleteReq 删除收款记录请求（后端软删除，同时作废自动生成的凭证）
type DeleteReq struct {
	g.Meta `path:"/payment/delete" method:"post" sm:"删除收款记录" tags:"Payment"`

	Id int64 `json:"id" v:"required#收款记录ID不能为空" dc:"收款记录ID"`
}

// DeleteRes 删除收款记录响应
type DeleteRes struct{}

// BatchDeleteReq 批量删除收款记录请求
type BatchDeleteReq struct {
	g.Meta `path:"/payment/batchDelete" method:"post" sm:"批量删除收款记录" tags:"Payment"`

	Ids []int64 `json:"ids" v:"required#请选择要删除的收款记录" dc:"收款记录ID集合"`
}

// BatchDeleteRes 批量删除收款记录响应
type BatchDeleteRes struct{}
