package inventory

import (
	"context"

	v1 "engineServ/api/inventory/v1"
)

// IInventoryV1 库存管理接口（跨项目库存总览 + 库存台账 CRUD + 仓库管理 CRUD）。
type IInventoryV1 interface {
	MaterialCatalogList(ctx context.Context, req *v1.MaterialCatalogListReq) (res *v1.MaterialCatalogListRes, err error)
	MaterialCatalogCreate(ctx context.Context, req *v1.MaterialCatalogCreateReq) (res *v1.MaterialCatalogCreateRes, err error)
	MaterialCatalogUpdate(ctx context.Context, req *v1.MaterialCatalogUpdateReq) (res *v1.MaterialCatalogUpdateRes, err error)
	MaterialCatalogDelete(ctx context.Context, req *v1.MaterialCatalogDeleteReq) (res *v1.MaterialCatalogDeleteRes, err error)

	StockList(ctx context.Context, req *v1.StockListReq) (res *v1.StockListRes, err error)
	StockCreate(ctx context.Context, req *v1.StockCreateReq) (res *v1.StockCreateRes, err error)
	StockUpdate(ctx context.Context, req *v1.StockUpdateReq) (res *v1.StockUpdateRes, err error)
	StockDelete(ctx context.Context, req *v1.StockDeleteReq) (res *v1.StockDeleteRes, err error)
	StockIn(ctx context.Context, req *v1.StockInReq) (res *v1.StockInRes, err error)
	StockOut(ctx context.Context, req *v1.StockOutReq) (res *v1.StockOutRes, err error)
	StockAdjust(ctx context.Context, req *v1.StockAdjustReq) (res *v1.StockAdjustRes, err error)
	StockRecordList(ctx context.Context, req *v1.StockRecordListReq) (res *v1.StockRecordListRes, err error)

	WarehouseList(ctx context.Context, req *v1.WarehouseListReq) (res *v1.WarehouseListRes, err error)
	WarehouseCreate(ctx context.Context, req *v1.WarehouseCreateReq) (res *v1.WarehouseCreateRes, err error)
	WarehouseUpdate(ctx context.Context, req *v1.WarehouseUpdateReq) (res *v1.WarehouseUpdateRes, err error)
	WarehouseDelete(ctx context.Context, req *v1.WarehouseDeleteReq) (res *v1.WarehouseDeleteRes, err error)
	WarehouseBatchDelete(ctx context.Context, req *v1.WarehouseBatchDeleteReq) (res *v1.WarehouseBatchDeleteRes, err error)
}
