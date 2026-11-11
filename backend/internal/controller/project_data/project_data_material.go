package project_data

import (
	"context"

	v1 "engineServ/api/project_data/v1"
	"engineServ/internal/consts"
	projectDataLogic "engineServ/internal/logic/api/project_data"

	"github.com/gogf/gf/v2/util/gconv"
)

// MaterialList 材料清单列表。
func (c *ControllerV1) MaterialList(ctx context.Context, req *v1.MaterialListReq) (*v1.MaterialListRes, error) {
	items, total, err := c.data.ListMaterials(ctx, req.ProjectId, req.Keyword, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.MaterialItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.MaterialItem{
			Id:             item.Id,
			CatalogId:      item.CatalogId,
			ProjectId:      item.ProjectId,
			WarehouseId:    item.WarehouseId,
			WarehouseName:  item.WarehouseName,
			Name:           item.Name,
			Spec:           item.Spec,
			Unit:           item.Unit,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			TaxMode:        item.TaxMode,
			TaxRate:        item.TaxRate,
			Amount:         item.Amount,
			UsedQuantity:   item.UsedQuantity,
			RemainQuantity: item.RemainQuantity,
			Supplier:       item.Supplier,
			Remark:         item.Remark,
			CreatedAt:      item.CreatedAt,
		})
	}

	return &v1.MaterialListRes{Current: req.Current, Size: req.Size, Total: total, Records: records}, nil
}

// MaterialCreate 新增材料。
func (c *ControllerV1) MaterialCreate(ctx context.Context, req *v1.MaterialCreateReq) (*v1.MaterialCreateRes, error) {
	id, err := c.data.CreateMaterial(ctx, &projectDataLogic.MaterialInput{
		ProjectId:   req.ProjectId,
		CatalogId:   req.CatalogId,
		WarehouseId: req.WarehouseId,
		Name:        req.Name,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		TaxMode:     req.TaxMode,
		TaxRate:     req.TaxRate,
		Supplier:    req.Supplier,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.MaterialCreateRes{Id: id}, nil
}

// MaterialUpdate 更新材料。
func (c *ControllerV1) MaterialUpdate(ctx context.Context, req *v1.MaterialUpdateReq) (*v1.MaterialUpdateRes, error) {
	err := c.data.UpdateMaterial(ctx, req.Id, &projectDataLogic.MaterialInput{
		ProjectId:   req.ProjectId,
		CatalogId:   req.CatalogId,
		WarehouseId: req.WarehouseId,
		Name:        req.Name,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		TaxMode:     req.TaxMode,
		TaxRate:     req.TaxRate,
		Supplier:    req.Supplier,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.MaterialUpdateRes{Id: req.Id}, nil
}

// MaterialDelete 删除材料。
func (c *ControllerV1) MaterialDelete(ctx context.Context, req *v1.MaterialDeleteReq) (*v1.MaterialDeleteRes, error) {
	if err := c.data.DeleteMaterial(ctx, req.Id); err != nil {
		return nil, err
	}

	return &v1.MaterialDeleteRes{}, nil
}

// MaterialUsageList 领用记录列表。
func (c *ControllerV1) MaterialUsageList(ctx context.Context, req *v1.MaterialUsageListReq) (*v1.MaterialUsageListRes, error) {
	items, total, err := c.data.ListUsages(ctx, req.ProjectId, req.MaterialId, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.MaterialUsageItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.MaterialUsageItem{
			Id:            item.Id,
			ProjectId:     item.ProjectId,
			MaterialId:    item.MaterialId,
			MaterialName:  item.MaterialName,
			WarehouseId:   item.WarehouseId,
			WarehouseName: item.WarehouseName,
			Quantity:      item.Quantity,
			UsedAt:        item.UsedAt,
			Receiver:      item.Receiver,
			Remark:        item.Remark,
			StockedAt:     item.StockedAt,
			CreatedBy:     item.CreatedBy,
			CreatedAt:     item.CreatedAt,
		})
	}

	return &v1.MaterialUsageListRes{Current: req.Current, Size: req.Size, Total: total, Records: records}, nil
}

// MaterialUsageCreate 新增领用记录。
func (c *ControllerV1) MaterialUsageCreate(ctx context.Context, req *v1.MaterialUsageCreateReq) (*v1.MaterialUsageCreateRes, error) {
	id, err := c.data.CreateUsage(ctx, &projectDataLogic.UsageInput{
		ProjectId:  req.ProjectId,
		MaterialId: req.MaterialId,
		Quantity:   req.Quantity,
		UsedAt:     req.UsedAt,
		Receiver:   req.Receiver,
		Remark:     req.Remark,
	}, gconv.Uint(ctx.Value(consts.CtxUserIdKey)))
	if err != nil {
		return nil, err
	}

	return &v1.MaterialUsageCreateRes{Id: id}, nil
}

// MaterialUsageStock 领用单出库：选择仓库后直接完成出库。
func (c *ControllerV1) MaterialUsageStock(ctx context.Context, req *v1.MaterialUsageStockReq) (*v1.MaterialUsageStockRes, error) {
	recordId, err := c.data.StockOutMaterialUsage(ctx, req.Id, req.WarehouseId, gconv.Uint(ctx.Value(consts.CtxUserIdKey)), req.Remark)
	if err != nil {
		return nil, err
	}

	return &v1.MaterialUsageStockRes{Id: req.Id, RecordId: recordId}, nil
}

// MaterialUsageDelete 删除领用记录。
func (c *ControllerV1) MaterialUsageDelete(ctx context.Context, req *v1.MaterialUsageDeleteReq) (*v1.MaterialUsageDeleteRes, error) {
	if err := c.data.DeleteUsage(ctx, req.Id); err != nil {
		return nil, err
	}

	return &v1.MaterialUsageDeleteRes{}, nil
}
