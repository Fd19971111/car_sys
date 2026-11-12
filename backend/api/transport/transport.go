package transport

import (
	"context"

	v1 "engineServ/api/transport/v1"
)

// ITransportV1 运力管理接口（司机档案 / 车辆档案）。
type ITransportV1 interface {
	DriverList(ctx context.Context, req *v1.DriverListReq) (res *v1.DriverListRes, err error)
	DriverDetail(ctx context.Context, req *v1.DriverDetailReq) (res *v1.DriverDetailRes, err error)
	DriverCreate(ctx context.Context, req *v1.DriverCreateReq) (res *v1.DriverCreateRes, err error)
	DriverUpdate(ctx context.Context, req *v1.DriverUpdateReq) (res *v1.DriverUpdateRes, err error)
	DriverDelete(ctx context.Context, req *v1.DriverDeleteReq) (res *v1.DriverDeleteRes, err error)
	DriverBatchDelete(ctx context.Context, req *v1.DriverBatchDeleteReq) (res *v1.DriverBatchDeleteRes, err error)
	DriverOptions(ctx context.Context, req *v1.DriverOptionReq) (res *v1.DriverOptionRes, err error)

	VehicleList(ctx context.Context, req *v1.VehicleListReq) (res *v1.VehicleListRes, err error)
	VehicleDetail(ctx context.Context, req *v1.VehicleDetailReq) (res *v1.VehicleDetailRes, err error)
	VehicleCreate(ctx context.Context, req *v1.VehicleCreateReq) (res *v1.VehicleCreateRes, err error)
	VehicleUpdate(ctx context.Context, req *v1.VehicleUpdateReq) (res *v1.VehicleUpdateRes, err error)
	VehicleDelete(ctx context.Context, req *v1.VehicleDeleteReq) (res *v1.VehicleDeleteRes, err error)
	VehicleBatchDelete(ctx context.Context, req *v1.VehicleBatchDeleteReq) (res *v1.VehicleBatchDeleteRes, err error)
}
