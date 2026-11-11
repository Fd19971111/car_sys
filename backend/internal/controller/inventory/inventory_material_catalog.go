package inventory

import (
	"context"

	v1 "engineServ/api/inventory/v1"
)

// MaterialCatalogList 查询公司级材料档案。
func (c *ControllerV1) MaterialCatalogList(ctx context.Context, req *v1.MaterialCatalogListReq) (*v1.MaterialCatalogListRes, error) {
	records, total, err := c.stock.ListMaterialCatalog(ctx, req.Keyword, req.Current, req.Size)
	if err != nil {
		return nil, err
	}
	current, size := req.Current, req.Size
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return &v1.MaterialCatalogListRes{Current: current, Size: size, Total: total, Records: records}, nil
}

// MaterialCatalogCreate 新增公司级材料档案。
func (c *ControllerV1) MaterialCatalogCreate(ctx context.Context, req *v1.MaterialCatalogCreateReq) (*v1.MaterialCatalogCreateRes, error) {
	id, err := c.stock.CreateMaterialCatalog(ctx, &req.MaterialCatalogForm)
	if err != nil {
		return nil, err
	}
	return &v1.MaterialCatalogCreateRes{Id: id}, nil
}

// MaterialCatalogUpdate 更新公司级材料档案。
func (c *ControllerV1) MaterialCatalogUpdate(ctx context.Context, req *v1.MaterialCatalogUpdateReq) (*v1.MaterialCatalogUpdateRes, error) {
	if err := c.stock.UpdateMaterialCatalog(ctx, req.Id, &req.MaterialCatalogForm); err != nil {
		return nil, err
	}
	return &v1.MaterialCatalogUpdateRes{Id: req.Id}, nil
}

// MaterialCatalogDelete 删除未被项目材料引用的档案。
func (c *ControllerV1) MaterialCatalogDelete(ctx context.Context, req *v1.MaterialCatalogDeleteReq) (*v1.MaterialCatalogDeleteRes, error) {
	if err := c.stock.DeleteMaterialCatalog(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.MaterialCatalogDeleteRes{}, nil
}
