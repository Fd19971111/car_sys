package inventory

import (
	"context"
	"strings"

	v1 "engineServ/api/inventory/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// WarehouseList 仓库分页列表。
func (c *ControllerV1) WarehouseList(ctx context.Context, req *v1.WarehouseListReq) (*v1.WarehouseListRes, error) {
	current, size := req.Current, req.Size
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	model := g.DB().Model("warehouse").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL")
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		model = model.Where("(name LIKE ? OR code LIKE ? OR address LIKE ?)", like, like, like)
	}
	var records []*v1.WarehouseItem
	var total int
	if err := model.Page(current, size).
		Order("is_default DESC, sort ASC, id ASC").
		ScanAndCount(&records, &total, false); err != nil {
		return nil, err
	}
	return &v1.WarehouseListRes{Current: current, Size: size, Total: total, Records: records}, nil
}

// WarehouseCreate 新增仓库；同公司下仓库名称不可重复，设为默认时清除其它默认标记。
func (c *ControllerV1) WarehouseCreate(ctx context.Context, req *v1.WarehouseCreateReq) (*v1.WarehouseCreateRes, error) {
	companyId := currentCompanyId(ctx)
	if err := checkWarehouseNameUnique(ctx, companyId, req.Name, 0); err != nil {
		return nil, err
	}
	if req.IsDefault {
		if _, err := g.DB().Model("warehouse").Ctx(ctx).
			Where("company_id", companyId).
			Where("deleted_at IS NULL").
			Data(g.Map{"is_default": 0}).
			Update(); err != nil {
			return nil, err
		}
	}
	result, err := g.DB().Model("warehouse").Ctx(ctx).Data(g.Map{
		"company_id": companyId,
		"name":       req.Name,
		"code":       req.Code,
		"address":    req.Address,
		"remark":     req.Remark,
		"is_default": req.IsDefault,
		"sort":       req.Sort,
	}).Insert()
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &v1.WarehouseCreateRes{Id: uint(id)}, nil
}

// WarehouseUpdate 更新仓库；设为默认时清除其它默认标记。
func (c *ControllerV1) WarehouseUpdate(ctx context.Context, req *v1.WarehouseUpdateReq) (*v1.WarehouseUpdateRes, error) {
	companyId := currentCompanyId(ctx)
	if err := checkWarehouseNameUnique(ctx, companyId, req.Name, req.Id); err != nil {
		return nil, err
	}
	if err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if req.IsDefault {
			if _, err := tx.Model("warehouse").Ctx(ctx).
				Where("company_id", companyId).
				Where("deleted_at IS NULL").
				WhereNot("id", req.Id).
				Data(g.Map{"is_default": 0}).
				Update(); err != nil {
				return err
			}
		}
		_, err := tx.Model("warehouse").Ctx(ctx).
			Where("company_id", companyId).
			Where("id", req.Id).
			Where("deleted_at IS NULL").
			Data(g.Map{
				"name":       req.Name,
				"code":       req.Code,
				"address":    req.Address,
				"remark":     req.Remark,
				"is_default": req.IsDefault,
				"sort":       req.Sort,
			}).Update()
		return err
	}); err != nil {
		return nil, err
	}
	return &v1.WarehouseUpdateRes{Id: req.Id}, nil
}

// WarehouseDelete 删除仓库（软删除）；默认仓库不允许删除。
func (c *ControllerV1) WarehouseDelete(ctx context.Context, req *v1.WarehouseDeleteReq) (*v1.WarehouseDeleteRes, error) {
	if err := checkWarehouseDeletable(ctx, currentCompanyId(ctx), []uint{req.Id}); err != nil {
		return nil, err
	}
	_, err := g.DB().Model("warehouse").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn("id", []uint{req.Id}).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now(), "is_default": 0}).
		Update()
	if err != nil {
		return nil, err
	}
	return &v1.WarehouseDeleteRes{}, nil
}

// WarehouseBatchDelete 批量删除仓库（软删除）；包含默认仓库时整体拒绝。
func (c *ControllerV1) WarehouseBatchDelete(ctx context.Context, req *v1.WarehouseBatchDeleteReq) (*v1.WarehouseBatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return &v1.WarehouseBatchDeleteRes{}, nil
	}
	if err := checkWarehouseDeletable(ctx, currentCompanyId(ctx), req.Ids); err != nil {
		return nil, err
	}
	_, err := g.DB().Model("warehouse").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn("id", req.Ids).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now(), "is_default": 0}).
		Update()
	if err != nil {
		return nil, err
	}
	return &v1.WarehouseBatchDeleteRes{}, nil
}

// checkWarehouseNameUnique 校验同公司下仓库名称唯一（排除自身与已删除记录）。
func checkWarehouseNameUnique(ctx context.Context, companyId uint, name string, excludeId uint) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return gerror.New("请输入仓库名称")
	}
	count, err := g.DB().Model("warehouse").Ctx(ctx).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		WhereNot("id", excludeId).
		Where("name", name).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("仓库名称已存在，请更换名称")
	}
	return nil
}

// checkWarehouseDeletable 校验待删除仓库存在且不含默认仓库。
func checkWarehouseDeletable(ctx context.Context, companyId uint, ids []uint) error {
	var records []*v1.WarehouseItem
	if err := g.DB().Model("warehouse").Ctx(ctx).
		Where("company_id", companyId).
		WhereIn("id", ids).
		Where("deleted_at IS NULL").
		Scan(&records); err != nil {
		return err
	}
	if len(records) != len(ids) {
		return gerror.New("仓库不存在或已删除")
	}
	for _, row := range records {
		if row.IsDefault {
			return gerror.New("默认仓库不允许删除，请先将其它仓库设为默认")
		}
	}
	return nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}
