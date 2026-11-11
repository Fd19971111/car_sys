package menu

import (
	"context"

	v1 "engineServ/api/menu/v1"
)

type IMenuV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Departments(ctx context.Context, req *v1.DepartmentsReq) (res *v1.DepartmentsRes, err error)
	CreateDepartment(ctx context.Context, req *v1.CreateDepartmentReq) (res *v1.CreateDepartmentRes, err error)
	Users(ctx context.Context, req *v1.UsersReq) (res *v1.UsersRes, err error)
	SetUserDepartment(ctx context.Context, req *v1.SetUserDepartmentReq) (res *v1.SetUserDepartmentRes, err error)
	DepartmentMenus(ctx context.Context, req *v1.DepartmentMenusReq) (res *v1.DepartmentMenusRes, err error)
	Save(ctx context.Context, req *v1.SaveReq) (res *v1.SaveRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
	SetDepartmentMenus(ctx context.Context, req *v1.SetDepartmentMenusReq) (res *v1.SetDepartmentMenusRes, err error)
	GetUserRoutes(ctx context.Context, req *v1.GetUserRoutesReq) (res *v1.GetUserRoutesRes, err error)
}
