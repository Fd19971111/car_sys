package payment

import (
	"context"

	v1 "engineServ/api/payment/v1"
)

// IPaymentV1 收款记录模块接口（均需登录鉴权，挂在 Auth 中间件分组下）
type IPaymentV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error)
	Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
	BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (res *v1.BatchDeleteRes, err error)
}
