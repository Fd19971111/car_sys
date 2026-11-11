package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 报价明细分页列表 =====================

// ListReq 报价明细分页列表请求
type ListReq struct {
	g.Meta `path:"/quote/list" method:"get" sm:"报价明细分页列表" tags:"Quote"`

	Current   int    `json:"current"   d:"1"  dc:"当前页码，默认 1"`
	Size      int    `json:"size"      d:"10" dc:"每页条数，默认 10"`
	Keyword   string `json:"keyword"          dc:"关键词：匹配报价项名称 / 项目名称 / 客户姓名"`
	ProjectId *int64 `json:"projectId"        dc:"所属项目ID"`
}

// QuoteListItem 报价明细列表行
type QuoteListItem struct {
	Id           int64   `json:"id"           dc:"报价明细ID"`
	ProjectId    int64   `json:"projectId"    dc:"所属项目ID"`
	ProjectName  string  `json:"projectName"  dc:"项目名称"`
	CustomerName string  `json:"customerName" dc:"客户姓名"`
	ItemName     string  `json:"itemName"     dc:"报价项名称"`
	Unit         string  `json:"unit"         dc:"单位"`
	Price        float64 `json:"price"        dc:"单价"`
	Quantity     float64 `json:"quantity"     dc:"数量"`
	Total        float64 `json:"total"        dc:"小计 = 单价 × 数量"`
	Remark       string  `json:"remark"       dc:"备注"`
	CreatedAt    string  `json:"createdAt"    dc:"创建时间"`
	UpdatedAt    string  `json:"updatedAt"    dc:"更新时间"`
}

// ListRes 报价明细分页列表响应
type ListRes struct {
	Current int              `json:"current" dc:"当前页码"`
	Size    int              `json:"size"    dc:"每页条数"`
	Total   int              `json:"total"   dc:"总条数"`
	Records []*QuoteListItem `json:"records" dc:"报价明细列表"`
}

// ===================== 报价概览统计 =====================

// StatsReq 报价概览统计请求
type StatsReq struct {
	g.Meta `path:"/quote/stats" method:"get" sm:"报价概览统计" tags:"Quote"`
}

// StatsRes 报价概览统计响应
type StatsRes struct {
	TotalCount   int     `json:"totalCount"   dc:"报价条目数"`
	TotalAmount  float64 `json:"totalAmount"  dc:"报价总额"`
	ProjectCount int     `json:"projectCount" dc:"涉及项目数"`
	AvgAmount    float64 `json:"avgAmount"    dc:"平均单条金额"`
}

// ===================== 报价明细详情 =====================

// DetailReq 报价明细详情请求
type DetailReq struct {
	g.Meta `path:"/quote/detail" method:"get" sm:"报价明细详情" tags:"Quote"`

	Id int64 `json:"id" v:"required#报价明细ID不能为空" dc:"报价明细ID"`
}

// DetailRes 报价明细详情响应
type DetailRes struct {
	Id            int64   `json:"id"            dc:"报价明细ID"`
	ProjectId     int64   `json:"projectId"     dc:"所属项目ID"`
	ProjectName   string  `json:"projectName"   dc:"项目名称"`
	CustomerName  string  `json:"customerName"  dc:"客户姓名"`
	ProjectAmount float64 `json:"projectAmount" dc:"所属项目合同金额"`
	QuoteAmount   float64 `json:"quoteAmount"   dc:"所属项目报价合计"`
	ItemName      string  `json:"itemName"      dc:"报价项名称"`
	Unit          string  `json:"unit"          dc:"单位"`
	Price         float64 `json:"price"         dc:"单价"`
	Quantity      float64 `json:"quantity"      dc:"数量"`
	Total         float64 `json:"total"         dc:"小计 = 单价 × 数量"`
	Remark        string  `json:"remark"        dc:"备注"`
	CreatedAt     string  `json:"createdAt"     dc:"创建时间"`
	UpdatedAt     string  `json:"updatedAt"     dc:"更新时间"`
}

// ===================== 新增 / 编辑 / 删除 =====================

// QuoteFormFields 报价明细表单公共字段（project_quote 表结构固定，不新增业务列）
type QuoteFormFields struct {
	ProjectId int64   `json:"projectId" v:"required#请选择所属项目"                                    dc:"所属项目ID"`
	ItemName  string  `json:"itemName"  v:"required|length:1,100#请输入报价项名称|报价项名称最多 100 个字符" dc:"报价项名称"`
	Unit      string  `json:"unit"      v:"max-length:20#单位最多 20 个字符"                          dc:"单位"`
	Price     float64 `json:"price"     v:"min:0#单价不能为负数"                                      dc:"单价（元）"`
	Quantity  float64 `json:"quantity"  v:"min:0#数量不能为负数"                                      dc:"数量"`
	Remark    string  `json:"remark"    v:"max-length:255#备注最多 255 个字符"                        dc:"备注"`
}

// CreateReq 新增报价明细请求
type CreateReq struct {
	g.Meta `path:"/quote/create" method:"post" sm:"新增报价明细" tags:"Quote"`

	QuoteFormFields
}

// CreateRes 新增报价明细响应
type CreateRes struct {
	Id      int64  `json:"id"      dc:"新建报价明细ID"`
	Warning string `json:"warning" dc:"报价合计超过结算金额时的提示，空值表示无异常"`
}

// UpdateReq 更新报价明细请求
type UpdateReq struct {
	g.Meta `path:"/quote/update" method:"post" sm:"更新报价明细" tags:"Quote"`

	Id int64 `json:"id" v:"required#报价明细ID不能为空" dc:"报价明细ID"`
	QuoteFormFields
}

// UpdateRes 更新报价明细响应
type UpdateRes struct {
	Id      int64  `json:"id"      dc:"报价明细ID"`
	Warning string `json:"warning" dc:"报价合计超过结算金额时的提示，空值表示无异常"`
}

// DeleteReq 删除报价明细请求（后端软删除）
type DeleteReq struct {
	g.Meta `path:"/quote/delete" method:"post" sm:"删除报价明细" tags:"Quote"`

	Id int64 `json:"id" v:"required#报价明细ID不能为空" dc:"报价明细ID"`
}

// DeleteRes 删除报价明细响应
type DeleteRes struct{}

// BatchDeleteReq 批量删除报价明细请求
type BatchDeleteReq struct {
	g.Meta `path:"/quote/batchDelete" method:"post" sm:"批量删除报价明细" tags:"Quote"`

	Ids []int64 `json:"ids" v:"required#请选择要删除的报价明细" dc:"报价明细ID集合"`
}

// BatchDeleteRes 批量删除报价明细响应
type BatchDeleteRes struct{}
