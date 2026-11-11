package project_data

import (
	"context"

	v1 "engineServ/api/project_data/v1"
	"engineServ/internal/consts"
	projectDataLogic "engineServ/internal/logic/api/project_data"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) PurchaseOrderList(ctx context.Context, req *v1.PurchaseOrderListReq) (*v1.PurchaseOrderListRes, error) {
	items, total, err := c.data.ListPurchaseOrders(ctx, req.ProjectId, req.Current, req.Size, false)
	if err != nil {
		return nil, err
	}
	return &v1.PurchaseOrderListRes{Current: req.Current, Size: req.Size, Total: total, Records: purchaseOrderItems(items)}, nil
}

func (c *ControllerV1) PurchaseOrderPendingList(ctx context.Context, _ *v1.PurchaseOrderPendingListReq) (*v1.PurchaseOrderPendingListRes, error) {
	items, err := c.data.ListPendingPurchaseOrders(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.PurchaseOrderPendingListRes{Records: purchaseOrderItems(items)}, nil
}

func (c *ControllerV1) PurchaseOrderCreate(ctx context.Context, req *v1.PurchaseOrderCreateReq) (*v1.PurchaseOrderCreateRes, error) {
	id, err := c.data.CreatePurchaseOrder(ctx, req.ProjectId, req.MaterialId, req.VendorId, req.PurchaseDate,
		req.Quantity, req.UnitPrice, req.TaxMode, req.TaxRate, req.Remark, gconv.Uint(ctx.Value(consts.CtxUserIdKey)))
	if err != nil {
		return nil, err
	}
	return &v1.PurchaseOrderCreateRes{Id: id}, nil
}

// PurchaseOrderStock 采购单入库：选择仓库后直接完成出入库。
func (c *ControllerV1) PurchaseOrderStock(ctx context.Context, req *v1.PurchaseOrderStockReq) (*v1.PurchaseOrderStockRes, error) {
	recordId, err := c.data.StockInPurchaseOrder(ctx, req.Id, req.WarehouseId, gconv.Uint(ctx.Value(consts.CtxUserIdKey)), req.Remark)
	if err != nil {
		return nil, err
	}

	return &v1.PurchaseOrderStockRes{Id: req.Id, RecordId: recordId}, nil
}

func purchaseOrderItems(items []*projectDataLogic.PurchaseOrder) []*v1.PurchaseOrderItem {
	result := make([]*v1.PurchaseOrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, &v1.PurchaseOrderItem{
			Id: item.Id, OrderNo: item.OrderNo, ProjectId: item.ProjectId, ProjectName: item.ProjectName,
			MaterialId: item.MaterialId, VendorId: item.VendorId, PurchaseDate: item.PurchaseDate,
			MaterialName: item.MaterialName, Spec: item.Spec, Unit: item.Unit,
			Quantity: item.Quantity, UnitPrice: item.UnitPrice, TaxMode: item.TaxMode, TaxRate: item.TaxRate,
			NetAmount: item.NetAmount, TaxAmount: item.TaxAmount, Amount: item.Amount, Supplier: item.Supplier,
			Remark: item.Remark, Status: item.Status, PayoutId: item.PayoutId,
			WarehouseId: item.WarehouseId, WarehouseName: item.WarehouseName,
			StockedAt: item.StockedAt, CreatedBy: item.CreatedBy, CreatedAt: item.CreatedAt,
		})
	}
	return result
}
