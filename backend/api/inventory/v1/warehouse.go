package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 仓库管理 =====================

// WarehouseListReq 仓库分页列表请求
type WarehouseListReq struct {
	g.Meta `path:"/inventory/warehouse/list" method:"get" tags:"仓库管理" summary:"仓库列表"`

	Keyword string `json:"keyword" dc:"仓库名称 / 编号 / 地址关键字"`
	Current int    `json:"current" d:"1"  dc:"当前页码"`
	Size    int    `json:"size"    d:"10" dc:"每页条数"`
}

// WarehouseItem 仓库行
type WarehouseItem struct {
	Id        uint   `json:"id"        dc:"仓库ID"`
	Name      string `json:"name"      dc:"仓库名称"`
	Code      string `json:"code"      dc:"仓库编号"`
	Address   string `json:"address"   dc:"仓库地址"`
	Remark    string `json:"remark"    dc:"备注"`
	IsDefault bool   `json:"isDefault" dc:"是否默认仓库"`
	Sort      int    `json:"sort"      dc:"排序（越小越靠前）"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
	UpdatedAt string `json:"updatedAt" dc:"更新时间"`
}

// WarehouseListRes 仓库分页列表响应
type WarehouseListRes struct {
	Current int              `json:"current" dc:"当前页码"`
	Size    int              `json:"size"    dc:"每页条数"`
	Total   int              `json:"total"   dc:"总条数"`
	Records []*WarehouseItem `json:"records" dc:"仓库列表"`
}

// WarehouseFormFields 仓库表单字段
type WarehouseFormFields struct {
	Name      string `json:"name"      v:"required|max-length:100#请输入仓库名称|仓库名称最多100个字符" dc:"仓库名称"`
	Code      string `json:"code"      v:"max-length:50#仓库编号最多50个字符"                        dc:"仓库编号"`
	Address   string `json:"address"   v:"max-length:255#仓库地址最多255个字符"                      dc:"仓库地址"`
	Remark    string `json:"remark"    v:"max-length:255#备注最多255个字符"                          dc:"备注"`
	IsDefault bool   `json:"isDefault"                                                          dc:"是否设为默认仓库"`
	Sort      int    `json:"sort"      d:"0"                                                    dc:"排序（越小越靠前）"`
}

// WarehouseCreateReq 新增仓库请求
type WarehouseCreateReq struct {
	g.Meta `path:"/inventory/warehouse/create" method:"post" tags:"仓库管理" summary:"新增仓库"`
	WarehouseFormFields
}

// WarehouseCreateRes 新增仓库响应
type WarehouseCreateRes struct {
	Id uint `json:"id" dc:"仓库ID"`
}

// WarehouseUpdateReq 更新仓库请求
type WarehouseUpdateReq struct {
	g.Meta `path:"/inventory/warehouse/update" method:"post" tags:"仓库管理" summary:"更新仓库"`
	Id     uint `json:"id" v:"required#仓库ID不能为空" dc:"仓库ID"`
	WarehouseFormFields
}

// WarehouseUpdateRes 更新仓库响应
type WarehouseUpdateRes struct {
	Id uint `json:"id" dc:"仓库ID"`
}

// WarehouseDeleteReq 删除仓库请求
type WarehouseDeleteReq struct {
	g.Meta `path:"/inventory/warehouse/delete" method:"post" tags:"仓库管理" summary:"删除仓库"`
	Id     uint `json:"id" v:"required#仓库ID不能为空" dc:"仓库ID"`
}

// WarehouseDeleteRes 删除仓库响应
type WarehouseDeleteRes struct{}

// WarehouseBatchDeleteReq 批量删除仓库请求
type WarehouseBatchDeleteReq struct {
	g.Meta `path:"/inventory/warehouse/batchDelete" method:"post" tags:"仓库管理" summary:"批量删除仓库"`
	Ids    []uint `json:"ids" v:"required#请选择要删除的仓库" dc:"仓库ID列表"`
}

// WarehouseBatchDeleteRes 批量删除仓库响应
type WarehouseBatchDeleteRes struct{}
