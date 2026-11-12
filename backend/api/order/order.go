package order

import (
	"context"

	v1 "engineServ/api/order/v1"
)

// IOrderV1 订单中心接口（全部订单 / 行程监控）。
type IOrderV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
	BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (res *v1.BatchDeleteRes, err error)
	Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error)
	Monitor(ctx context.Context, req *v1.MonitorReq) (res *v1.MonitorRes, err error)
}
