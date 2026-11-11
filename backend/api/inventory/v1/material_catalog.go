package v1

import "github.com/gogf/gf/v2/frame/g"

// MaterialCatalogListReq 公司级材料档案分页列表请求。
type MaterialCatalogListReq struct {
	g.Meta `path:"/inventory/materialCatalog/list" method:"get" tags:"库存管理" summary:"材料档案列表"`

	Keyword string `json:"keyword" dc:"材料名称 / 规格 / 供应商关键字"`
	Current int    `json:"current" d:"1" dc:"当前页码"`
	Size    int    `json:"size"    d:"10" dc:"每页条数"`
}

// MaterialCatalogItem 公司级材料档案。
type MaterialCatalogItem struct {
	Id        int64   `json:"id"        dc:"材料档案ID"`
	Name      string  `json:"name"      dc:"材料名称"`
	Spec      string  `json:"spec"      dc:"规格 / 型号"`
	Unit      string  `json:"unit"      dc:"单位"`
	UnitPrice float64 `json:"unitPrice" dc:"默认单价"`
	TaxMode   string  `json:"taxMode"   dc:"默认税务方式"`
	TaxRate   float64 `json:"taxRate"   dc:"默认税率（百分比）"`
	Remark    string  `json:"remark"    dc:"备注"`
	CreatedAt string  `json:"createdAt" dc:"创建时间"`
	UpdatedAt string  `json:"updatedAt" dc:"更新时间"`
}

// MaterialCatalogListRes 公司级材料档案分页列表响应。
type MaterialCatalogListRes struct {
	Current int                    `json:"current" dc:"当前页码"`
	Size    int                    `json:"size"    dc:"每页条数"`
	Total   int                    `json:"total"   dc:"总条数"`
	Records []*MaterialCatalogItem `json:"records" dc:"材料档案列表"`
}

// MaterialCatalogForm 公司级材料档案表单字段。
type MaterialCatalogForm struct {
	Name      string  `json:"name"      v:"required|max-length:100#请填写材料名称|材料名称最多 100 个字符" dc:"材料名称"`
	Spec      string  `json:"spec"      v:"max-length:100#规格最多 100 个字符"                       dc:"规格 / 型号"`
	Unit      string  `json:"unit"      v:"max-length:20#单位最多 20 个字符"                         dc:"单位"`
	UnitPrice float64 `json:"unitPrice" v:"min:0#单价不能为负数"                                    dc:"默认单价"`
	TaxMode   string  `json:"taxMode"   v:"required|in:NO_INVOICE,TAX_INCLUDED,TAX_EXCLUDED#请选择税务方式|税务方式无效" dc:"默认税务方式"`
	TaxRate   float64 `json:"taxRate"   v:"min:0|max:100#税率不能超过 100%"                           dc:"默认税率（百分比）"`
	Remark    string  `json:"remark"    v:"max-length:255#备注最多 255 个字符"                       dc:"备注"`
}

// MaterialCatalogCreateReq 新增公司级材料档案。
type MaterialCatalogCreateReq struct {
	g.Meta `path:"/inventory/materialCatalog/create" method:"post" tags:"库存管理" summary:"新增材料档案"`
	MaterialCatalogForm
}

// MaterialCatalogCreateRes 新增公司级材料档案响应。
type MaterialCatalogCreateRes struct {
	Id int64 `json:"id" dc:"材料档案ID"`
}

// MaterialCatalogUpdateReq 更新公司级材料档案。
type MaterialCatalogUpdateReq struct {
	g.Meta `path:"/inventory/materialCatalog/update" method:"post" tags:"库存管理" summary:"更新材料档案"`
	Id     int64 `json:"id" v:"required#材料档案ID不能为空" dc:"材料档案ID"`
	MaterialCatalogForm
}

// MaterialCatalogUpdateRes 更新公司级材料档案响应。
type MaterialCatalogUpdateRes struct {
	Id int64 `json:"id" dc:"材料档案ID"`
}

// MaterialCatalogDeleteReq 删除公司级材料档案。
type MaterialCatalogDeleteReq struct {
	g.Meta `path:"/inventory/materialCatalog/delete" method:"post" tags:"库存管理" summary:"删除材料档案"`
	Id     int64 `json:"id" v:"required#材料档案ID不能为空" dc:"材料档案ID"`
}

// MaterialCatalogDeleteRes 删除公司级材料档案响应。
type MaterialCatalogDeleteRes struct{}
