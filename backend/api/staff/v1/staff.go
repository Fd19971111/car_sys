package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 人员列表 =====================

type ListReq struct {
	g.Meta       `path:"/staff/list" method:"get" sm:"人员分页列表" tags:"Staff"`
	Current      int    `json:"current" d:"1"`
	Size         int    `json:"size" d:"10"`
	Keyword      string `json:"keyword"`
	DepartmentId *int   `json:"departmentId"`
	Status       *int   `json:"status"`
}

type StaffItem struct {
	Id             uint   `json:"id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Gender         string `json:"gender"`
	Position       string `json:"position"`
	DepartmentId   uint   `json:"departmentId"`
	DepartmentName string `json:"departmentName"`
	JoinDate       string `json:"joinDate"`
	Status         int    `json:"status"`
	Remark         string `json:"remark"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type ListRes struct {
	Current int          `json:"current"`
	Size    int          `json:"size"`
	Total   int          `json:"total"`
	Records []*StaffItem `json:"records"`
}

// ===================== 人员详情 =====================

type DetailReq struct {
	g.Meta `path:"/staff/detail" method:"get" sm:"人员详情" tags:"Staff"`
	Id     uint `json:"id" v:"required"`
}

type DetailRes = StaffItem

// ===================== 新增 / 编辑 / 删除 =====================

type StaffForm struct {
	Name         string `json:"name" v:"required|length:2,20"`
	Phone        string `json:"phone" v:"required|regex:^1[3-9][0-9]{9}$"`
	Gender       string `json:"gender"`
	Position     string `json:"position"`
	DepartmentId uint   `json:"departmentId"`
	JoinDate     string `json:"joinDate"`
	Status       int    `json:"status" v:"required|in:1,2"`
	Remark       string `json:"remark"`
}

type CreateReq struct {
	g.Meta `path:"/staff/create" method:"post" sm:"新增人员" tags:"Staff"`
	StaffForm
}

type CreateRes struct {
	Id uint `json:"id"`
}

type UpdateReq struct {
	g.Meta `path:"/staff/update" method:"post" sm:"更新人员" tags:"Staff"`
	Id     uint `json:"id" v:"required"`
	StaffForm
}

type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/staff/delete" method:"post" sm:"删除人员" tags:"Staff"`
	Id     uint `json:"id" v:"required"`
}

type DeleteRes struct{}

type BatchDeleteReq struct {
	g.Meta `path:"/staff/batchDelete" method:"post" sm:"批量删除人员" tags:"Staff"`
	Ids    []uint `json:"ids" v:"required"`
}

type BatchDeleteRes struct{}

// ===================== 部门下拉选项 =====================

type DepartmentOptionReq struct {
	g.Meta `path:"/department/options" method:"get" sm:"部门下拉选项" tags:"Staff"`
}

type DepartmentOption struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

type DepartmentOptionRes struct {
	Items []DepartmentOption `json:"items"`
}