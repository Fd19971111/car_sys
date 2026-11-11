package project_data

import (
	"context"

	v1 "engineServ/api/project_data/v1"
)

// IProjectDataV1 项目管理 - 项目资料接口（材料清单 / 领用记录 / 施工日志 / 变更记录）。
type IProjectDataV1 interface {
	MaterialList(ctx context.Context, req *v1.MaterialListReq) (res *v1.MaterialListRes, err error)
	MaterialCreate(ctx context.Context, req *v1.MaterialCreateReq) (res *v1.MaterialCreateRes, err error)
	MaterialUpdate(ctx context.Context, req *v1.MaterialUpdateReq) (res *v1.MaterialUpdateRes, err error)
	MaterialDelete(ctx context.Context, req *v1.MaterialDeleteReq) (res *v1.MaterialDeleteRes, err error)
	MaterialUsageList(ctx context.Context, req *v1.MaterialUsageListReq) (res *v1.MaterialUsageListRes, err error)
	MaterialUsageCreate(ctx context.Context, req *v1.MaterialUsageCreateReq) (res *v1.MaterialUsageCreateRes, err error)
	MaterialUsageDelete(ctx context.Context, req *v1.MaterialUsageDeleteReq) (res *v1.MaterialUsageDeleteRes, err error)
	PurchaseOrderList(ctx context.Context, req *v1.PurchaseOrderListReq) (res *v1.PurchaseOrderListRes, err error)
	PurchaseOrderPendingList(ctx context.Context, req *v1.PurchaseOrderPendingListReq) (res *v1.PurchaseOrderPendingListRes, err error)
	PurchaseOrderCreate(ctx context.Context, req *v1.PurchaseOrderCreateReq) (res *v1.PurchaseOrderCreateRes, err error)

	SiteLogList(ctx context.Context, req *v1.SiteLogListReq) (res *v1.SiteLogListRes, err error)
	SiteLogCreate(ctx context.Context, req *v1.SiteLogCreateReq) (res *v1.SiteLogCreateRes, err error)
	SiteLogUpdate(ctx context.Context, req *v1.SiteLogUpdateReq) (res *v1.SiteLogUpdateRes, err error)
	SiteLogDelete(ctx context.Context, req *v1.SiteLogDeleteReq) (res *v1.SiteLogDeleteRes, err error)

	ChangeList(ctx context.Context, req *v1.ChangeListReq) (res *v1.ChangeListRes, err error)
	ChangeCreate(ctx context.Context, req *v1.ChangeCreateReq) (res *v1.ChangeCreateRes, err error)
	ChangeUpdate(ctx context.Context, req *v1.ChangeUpdateReq) (res *v1.ChangeUpdateRes, err error)
	ChangeDelete(ctx context.Context, req *v1.ChangeDeleteReq) (res *v1.ChangeDeleteRes, err error)

	ContractList(ctx context.Context, req *v1.ContractListReq) (res *v1.ContractListRes, err error)
	ContractCreate(ctx context.Context, req *v1.ContractCreateReq) (res *v1.ContractCreateRes, err error)
	ContractUpdate(ctx context.Context, req *v1.ContractUpdateReq) (res *v1.ContractUpdateRes, err error)
	ContractDelete(ctx context.Context, req *v1.ContractDeleteReq) (res *v1.ContractDeleteRes, err error)
}
