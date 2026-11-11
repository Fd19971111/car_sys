package inventory

import (
	"context"

	v1 "engineServ/api/inventory/v1"
	"engineServ/internal/consts"
	inventoryLogic "engineServ/internal/logic/api/inventory"
	projectDataLogic "engineServ/internal/logic/api/project_data"

	"github.com/gogf/gf/v2/util/gconv"
)

// StockList 库存列表（跨项目）。
func (c *ControllerV1) StockList(ctx context.Context, req *v1.StockListReq) (*v1.StockListRes, error) {
	items, summary, total, err := c.stock.ListStock(ctx, req.ProjectId, req.WarehouseId, req.Keyword, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.StockItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.StockItem{
			Id:             item.Id,
			ProjectId:      item.ProjectId,
			ProjectName:    item.ProjectName,
			WarehouseId:    item.WarehouseId,
			WarehouseName:  item.WarehouseName,
			Name:           item.Name,
			Spec:           item.Spec,
			Unit:           item.Unit,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			Amount:         item.Amount,
			UsedQuantity:   item.UsedQuantity,
			RemainQuantity: item.RemainQuantity,
			Supplier:       item.Supplier,
			Remark:         item.Remark,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	return &v1.StockListRes{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		Summary: &v1.StockSummary{
			Kinds:    summary.Kinds,
			Amount:   summary.Amount,
			Shortage: summary.Shortage,
		},
		Records: records,
	}, nil
}

// StockCreate 新增库存。
func (c *ControllerV1) StockCreate(ctx context.Context, req *v1.StockCreateReq) (*v1.StockCreateRes, error) {
	id, err := c.stock.CreateStock(ctx, &projectDataLogic.MaterialInput{
		ProjectId:   req.ProjectId,
		WarehouseId: req.WarehouseId,
		Name:        req.Name,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		Supplier:    req.Supplier,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.StockCreateRes{Id: id}, nil
}

// StockUpdate 更新库存。
func (c *ControllerV1) StockUpdate(ctx context.Context, req *v1.StockUpdateReq) (*v1.StockUpdateRes, error) {
	err := c.stock.UpdateStock(ctx, req.Id, &projectDataLogic.MaterialInput{
		ProjectId:   req.ProjectId,
		WarehouseId: req.WarehouseId,
		Name:        req.Name,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		Supplier:    req.Supplier,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.StockUpdateRes{Id: req.Id}, nil
}

// StockDelete 删除库存。
func (c *ControllerV1) StockDelete(ctx context.Context, req *v1.StockDeleteReq) (*v1.StockDeleteRes, error) {
	if err := c.stock.DeleteStock(ctx, req.Id); err != nil {
		return nil, err
	}

	return &v1.StockDeleteRes{}, nil
}

// StockIn 材料入库：库存增加，可选生成付款单。
func (c *ControllerV1) StockIn(ctx context.Context, req *v1.StockInReq) (*v1.StockInRes, error) {
	recordId, payoutId, voucherNo, err := c.stock.StockIn(ctx, &inventoryLogic.MovementInput{
		ProjectId:       req.ProjectId,
		MaterialId:      req.MaterialId,
		Quantity:        req.Quantity,
		UnitPrice:       req.UnitPrice,
		GenerateFinance: req.GenerateFinance,
		VendorId:        req.VendorId,
		FundAccountId:   req.FundAccountId,
		BizDate:         req.BizDate,
		Remark:          req.Remark,
	}, currentUserId(ctx))
	if err != nil {
		return nil, err
	}

	return &v1.StockInRes{Id: recordId, PayoutId: payoutId, VoucherNo: voucherNo}, nil
}

// StockOut 材料出库（项目领用）：写入领用记录并扣减结存，不生成任何财务单据。
func (c *ControllerV1) StockOut(ctx context.Context, req *v1.StockOutReq) (*v1.StockOutRes, error) {
	recordId, err := c.stock.StockOut(ctx, &inventoryLogic.MovementInput{
		ProjectId:  req.ProjectId,
		MaterialId: req.MaterialId,
		Quantity:   req.Quantity,
		UnitPrice:  req.UnitPrice,
		Receiver:   req.Receiver,
		BizDate:    req.BizDate,
		Remark:     req.Remark,
	}, currentUserId(ctx))
	if err != nil {
		return nil, err
	}

	return &v1.StockOutRes{Id: recordId}, nil
}

// StockAdjust 库存补录：直接调整库存，不生成财务单据。
func (c *ControllerV1) StockAdjust(ctx context.Context, req *v1.StockAdjustReq) (*v1.StockAdjustRes, error) {
	recordId, err := c.stock.StockAdjust(ctx, &inventoryLogic.MovementInput{
		RecordType: req.RecordType,
		ProjectId:  req.ProjectId,
		MaterialId: req.MaterialId,
		Quantity:   req.Quantity,
		UnitPrice:  req.UnitPrice,
		Receiver:   req.Receiver,
		BizDate:    req.BizDate,
		Remark:     req.Remark,
	}, currentUserId(ctx))
	if err != nil {
		return nil, err
	}

	return &v1.StockAdjustRes{Id: recordId}, nil
}

// StockRecordList 出入库流水分页列表。
func (c *ControllerV1) StockRecordList(ctx context.Context, req *v1.StockRecordListReq) (*v1.StockRecordListRes, error) {
	items, total, err := c.stock.ListRecords(ctx, req.ProjectId, req.MaterialId, req.RecordType, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.StockRecordItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.StockRecordItem{
			Id:               item.Id,
			RecordType:       item.RecordType,
			BizType:          item.BizType,
			ProjectName:      item.ProjectName,
			MaterialName:     item.MaterialName,
			Spec:             item.Spec,
			Unit:             item.Unit,
			WarehouseName:    item.WarehouseName,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			Amount:           item.Amount,
			VendorName:       item.VendorName,
			FundAccountName:  item.FundAccountName,
			PayoutId:         item.PayoutId,
			PaymentId:        item.PaymentId,
			PayoutVoucherNo:  item.PayoutVoucherNo,
			PaymentVoucherNo: item.PaymentVoucherNo,
			BizDate:          item.BizDate,
			Remark:           item.Remark,
			CreatedAt:        item.CreatedAt,
		})
	}

	return &v1.StockRecordListRes{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		Records: records,
	}, nil
}

func currentUserId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxUserIdKey))
}
