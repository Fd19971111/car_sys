package v1

import "github.com/gogf/gf/v2/frame/g"

type MenuItem struct {
	Id         uint   `json:"id"`
	ParentId   uint   `json:"parentId"`
	Title      string `json:"title"`
	RouteName  string `json:"routeName"`
	RoutePath  string `json:"routePath"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	SortOrder  int    `json:"sortOrder"`
	Enabled    bool   `json:"enabled"`
	HideInMenu bool   `json:"hideInMenu"`
}

type ListReq struct {
	g.Meta `path:"/menu/list" method:"get" sm:"菜单列表" tags:"Menu"`
}

type ListRes struct {
	Items []*MenuItem `json:"items"`
}

type DepartmentsReq struct {
	g.Meta `path:"/menu/departments" method:"get" sm:"部门列表" tags:"Menu"`
}

type Department struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type DepartmentsRes struct {
	Items []*Department `json:"items"`
}

type CreateDepartmentReq struct {
	g.Meta `path:"/menu/createDepartment" method:"post" sm:"新增部门" tags:"Menu"`
	Name   string `json:"name" v:"required|length:2,40#请输入部门名称|部门名称长度为2-40个字符"`
	Code   string `json:"code" v:"required|regex:^[a-zA-Z0-9_-]+$#请输入部门编码|部门编码只能包含字母、数字、下划线和短横线"`
}

type CreateDepartmentRes struct {
	Id uint `json:"id"`
}

type UsersReq struct {
	g.Meta `path:"/menu/users" method:"get" sm:"后台人员列表" tags:"Menu"`
}

type DepartmentUser struct {
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Role         uint   `json:"role"`
	DepartmentId uint   `json:"departmentId"`
}

type UsersRes struct {
	Items []*DepartmentUser `json:"items"`
}

type SetUserDepartmentReq struct {
	g.Meta       `path:"/menu/setUserDepartment" method:"post" sm:"设置人员所属部门" tags:"Menu"`
	UserId       uint `json:"userId" v:"required#请选择人员"`
	DepartmentId uint `json:"departmentId"`
}

type SetUserDepartmentRes struct{}

type DepartmentMenusReq struct {
	g.Meta       `path:"/menu/departmentMenus" method:"get" sm:"部门菜单权限" tags:"Menu"`
	DepartmentId uint `json:"departmentId" v:"required#请选择部门"`
}

type DepartmentMenusRes struct {
	MenuIds []uint `json:"menuIds"`
}

type SaveReq struct {
	g.Meta     `path:"/menu/save" method:"post" sm:"保存菜单" tags:"Menu"`
	Id         uint   `json:"id"`
	ParentId   uint   `json:"parentId"`
	Title      string `json:"title" v:"required#请输入菜单名称"`
	RouteName  string `json:"routeName" v:"required#请输入路由名称"`
	RoutePath  string `json:"routePath" v:"required#请输入路由路径"`
	Component  string `json:"component" v:"required#请选择页面组件"`
	Icon       string `json:"icon"`
	SortOrder  int    `json:"sortOrder"`
	Enabled    bool   `json:"enabled"`
	HideInMenu bool   `json:"hideInMenu"`
}

type SaveRes struct {
	Id uint `json:"id"`
}

type DeleteReq struct {
	g.Meta `path:"/menu/delete" method:"post" sm:"删除菜单" tags:"Menu"`
	Id     uint `json:"id" v:"required#菜单ID不能为空"`
}

type DeleteRes struct{}

type SetDepartmentMenusReq struct {
	g.Meta       `path:"/menu/setDepartmentMenus" method:"post" sm:"设置部门菜单权限" tags:"Menu"`
	DepartmentId uint   `json:"departmentId" v:"required#请选择部门"`
	MenuIds      []uint `json:"menuIds"`
}

type SetDepartmentMenusRes struct{}

type GetUserRoutesReq struct {
	g.Meta `path:"/route/getUserRoutes" method:"get" sm:"获取用户菜单路由" tags:"Menu"`
}

type RouteMeta struct {
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Order      int    `json:"order"`
	HideInMenu bool   `json:"hideInMenu"`
}

type MenuRoute struct {
	Name      string       `json:"name"`
	Path      string       `json:"path"`
	Component string       `json:"component"`
	Meta      RouteMeta    `json:"meta"`
	Children  []*MenuRoute `json:"children,omitempty"`
}

type UserRoute struct {
	Routes []*MenuRoute `json:"routes"`
	Home   string       `json:"home"`
}

type GetUserRoutesRes struct {
	UserRoute
}
