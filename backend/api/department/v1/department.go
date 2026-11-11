package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 部门列表 =====================

type ListReq struct {
	g.Meta  `path:"/department/list" method:"get" sm:"部门分页列表" tags:"Department"`
	Current int    `json:"current" d:"1"`
	Size    int    `json:"size" d:"10"`
	Keyword string `json:"keyword"`
	Status  *int   `json:"status"`
}

type DepartmentItem struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	LeaderId    *uint  `json:"leaderId"`
	LeaderName  string `json:"leaderName"`
	MemberCount int    `json:"memberCount"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
	SortOrder   int    `json:"sortOrder"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ListRes struct {
	Current int               `json:"current"`
	Size    int               `json:"size"`
	Total   int               `json:"total"`
	Records []*DepartmentItem `json:"records"`
}

// ===================== 部门详情 =====================

type DetailReq struct {
	g.Meta `path:"/department/detail" method:"get" sm:"部门详情" tags:"Department"`
	Id     uint `json:"id" v:"required"`
}

type DetailRes = DepartmentItem

// ===================== 新增 / 编辑 / 删除 =====================

type DepartmentForm struct {
	Name      string `json:"name" v:"required|length:2,30"`
	Code      string `json:"code" v:"required|length:2,20"`
	LeaderId  *uint  `json:"leaderId"`
	Status    int    `json:"status" v:"required|in:1,2"`
	SortOrder int    `json:"sortOrder"`
	Remark    string `json:"remark"`
}

type CreateReq struct {
	g.Meta `path:"/department/create" method:"post" sm:"新增部门" tags:"Department"`
	DepartmentForm
}

type CreateRes struct {
	Id uint `json:"id"`
}

type UpdateReq struct {
	g.Meta `path:"/department/update" method:"post" sm:"更新部门" tags:"Department"`
	Id     uint `json:"id" v:"required"`
	DepartmentForm
}

type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/department/delete" method:"post" sm:"删除部门" tags:"Department"`
	Id     uint `json:"id" v:"required"`
}

type DeleteRes struct{}

type BatchDeleteReq struct {
	g.Meta `path:"/department/batchDelete" method:"post" sm:"批量删除部门" tags:"Department"`
	Ids    []uint `json:"ids" v:"required"`
}

type BatchDeleteRes struct{}

// ===================== 下拉选项 =====================

type OptionsReq struct {
	g.Meta `path:"/department/options" method:"get" sm:"部门下拉选项" tags:"Department"`
}

type OptionsRes struct {
	Items []Option `json:"items"`
}

type Option struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}