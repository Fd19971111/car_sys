package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 用户列表 =====================

type ListReq struct {
	g.Meta  `path:"/user/list" method:"get" sm:"用户分页列表" tags:"User"`
	Current int    `json:"current" d:"1"`
	Size    int    `json:"size" d:"10"`
	Keyword string `json:"keyword"`
	Status  *int   `json:"status"`
	Role    *int   `json:"role"`
}

type UserItem struct {
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	Phone          string `json:"phone"`
	Role           int    `json:"role"`
	RoleName       string `json:"roleName"`
	Status         int    `json:"status"`
	DepartmentId   uint   `json:"departmentId"`
	DepartmentName string `json:"departmentName"`
	LastLoginAt    string `json:"lastLoginAt"`
	CreatedAt      string `json:"createdAt"`
}

type ListRes struct {
	Current int         `json:"current"`
	Size    int         `json:"size"`
	Total   int         `json:"total"`
	Records []*UserItem `json:"records"`
}

// ===================== 重置密码 =====================

type ResetPwdReq struct {
	g.Meta `path:"/user/resetPwd" method:"post" sm:"重置密码" tags:"User"`
	Id     uint `json:"id" v:"required"`
}

type ResetPwdRes struct {
	Password string `json:"password"`
}

// ===================== 启停账号 =====================

type ToggleReq struct {
	g.Meta `path:"/user/toggle" method:"post" sm:"启停账号" tags:"User"`
	Id     uint `json:"id" v:"required"`
	Status int  `json:"status" v:"required|in:0,1"`
}

type ToggleRes struct{}

// ===================== 创建账号（关联员工） =====================

type CreateAccountReq struct {
	g.Meta   `path:"/user/create" method:"post" sm:"创建登录账号" tags:"User"`
	StaffId  uint   `json:"staffId" v:"required"`
	Username string `json:"username" v:"required|length:4,20"`
	Password string `json:"password" v:"required|length:6,20"`
}

type CreateAccountRes struct {
	Id uint `json:"id"`
}

type AvailableStaffReq struct {
	g.Meta `path:"/user/availableStaff" method:"get" sm:"待开通账号员工列表" tags:"User"`
}

type StaffOption struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type AvailableStaffRes struct {
	Items []StaffOption `json:"items"`
}

type UpdateAccountReq struct {
	g.Meta   `path:"/user/update" method:"post" sm:"编辑登录账号" tags:"User"`
	Id       uint   `json:"id" v:"required"`
	Username string `json:"username" v:"required|length:4,20"`
	Role     int    `json:"role" v:"required|in:1,2"`
}

type UpdateAccountRes struct{}

// ===================== 修改密码 =====================

type ChangePwdReq struct {
	g.Meta   `path:"/user/changePwd" method:"post" sm:"修改密码" tags:"User"`
	Id       uint   `json:"id" v:"required"`
	Password string `json:"password" v:"required|length:6,20"`
}

type ChangePwdRes struct{}
