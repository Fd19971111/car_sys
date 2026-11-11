package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
)

// IProjectListV1 项目模块接口（均需登录鉴权，挂在 Auth 中间件分组下）
type IProjectListV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error)
	Options(ctx context.Context, req *v1.OptionsReq) (res *v1.OptionsRes, err error)
	SelectOptions(ctx context.Context, req *v1.SelectOptionsReq) (res *v1.SelectOptionsRes, err error)
	Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
	BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (res *v1.BatchDeleteRes, err error)
	ProgressCreate(ctx context.Context, req *v1.ProgressCreateReq) (res *v1.ProgressCreateRes, err error)
}
