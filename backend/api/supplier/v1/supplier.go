package v1

import "github.com/gogf/gf/v2/frame/g"

type ListReq struct {
	g.Meta  `path:"/vendor/list" method:"get" sm:"供应商分页列表" tags:"Vendor"`
	Current int    `json:"current" d:"1"`
	Size    int    `json:"size" d:"10"`
	Keyword string `json:"keyword" dc:"供应商名称、联系人、电话或地址"`
}

type VendorItem struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	ContactName string `json:"contactName"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Remark      string `json:"remark"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ListRes struct {
	Current int           `json:"current"`
	Size    int           `json:"size"`
	Total   int           `json:"total"`
	Records []*VendorItem `json:"records"`
}

type VendorFormFields struct {
	Name        string `json:"name" v:"required|max-length:100#请输入供应商名称|供应商名称最多100个字符"`
	ContactName string `json:"contactName" v:"max-length:50#联系人最多50个字符"`
	Phone       string `json:"phone" v:"max-length:30#联系电话最多30个字符"`
	Address     string `json:"address" v:"max-length:255#地址最多255个字符"`
	Remark      string `json:"remark" v:"max-length:255#备注最多255个字符"`
}

type CreateReq struct {
	g.Meta `path:"/vendor/create" method:"post" sm:"新增供应商" tags:"Vendor"`
	VendorFormFields
}

type UpdateReq struct {
	g.Meta `path:"/vendor/update" method:"post" sm:"更新供应商" tags:"Vendor"`
	Id     uint `json:"id" v:"required#供应商ID不能为空"`
	VendorFormFields
}

type SaveRes struct {
	Id uint `json:"id"`
}

type DeleteReq struct {
	g.Meta `path:"/vendor/delete" method:"post" sm:"删除供应商" tags:"Vendor"`
	Id     uint `json:"id" v:"required#供应商ID不能为空"`
}

type DeleteRes struct{}

type BatchDeleteReq struct {
	g.Meta `path:"/vendor/batchDelete" method:"post" sm:"批量删除供应商" tags:"Vendor"`
	Ids    []uint `json:"ids" v:"required#请选择要删除的供应商"`
}

type BatchDeleteRes struct{}
