package project_data

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// IncreaseMaterialStock 独立增加材料库存，供后续单独入库流程复用。
func (s *Service) IncreaseMaterialStock(ctx context.Context, projectId, materialId int64, quantity float64) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return s.AddMaterialStockTx(ctx, tx, projectId, materialId, quantity)
	})
}

// AddMaterialStockTx 采购入库：在调用方事务内增加材料库存数量。
func (s *Service) AddMaterialStockTx(ctx context.Context, tx gdb.TX, projectId, materialId int64, quantity float64) error {
	if quantity <= 0 {
		return gerror.New("入库数量必须大于 0")
	}
	material, err := tx.GetOne("SELECT project_id FROM project_material WHERE id=? AND deleted_at IS NULL FOR UPDATE", materialId)
	if err != nil {
		glog.Errorf(ctx, "读取入库材料失败, materialId=%d: %v", materialId, err)
		return gerror.New("读取材料失败，请稍后重试")
	}
	if material == nil || material.IsEmpty() {
		return gerror.New("入库材料不存在或已删除")
	}
	if material["project_id"].Int64() != projectId {
		return gerror.New("采购材料与项目不一致")
	}
	if _, err = tx.Exec("UPDATE project_material SET quantity=quantity+? WHERE id=? AND deleted_at IS NULL", round2(quantity), materialId); err != nil {
		glog.Errorf(ctx, "增加材料库存失败, materialId=%d quantity=%v: %v", materialId, quantity, err)
		return gerror.New("材料入库失败，请稍后重试")
	}
	return nil
}

// RemoveMaterialStockTx 撤销采购入库；库存已被领取时，不允许回退到负库存。
func (s *Service) RemoveMaterialStockTx(ctx context.Context, tx gdb.TX, projectId, materialId int64, quantity float64) error {
	if quantity <= 0 {
		return gerror.New("撤销入库数量必须大于 0")
	}
	material, err := tx.GetOne(`SELECT m.project_id,m.quantity,COALESCE(u.used_quantity,0) AS used_quantity
		FROM project_material m
		LEFT JOIN (
			SELECT material_id,SUM(quantity) AS used_quantity
			FROM project_material_usage WHERE deleted_at IS NULL GROUP BY material_id
		) u ON u.material_id=m.id
		WHERE m.id=? AND m.deleted_at IS NULL FOR UPDATE`, materialId)
	if err != nil {
		glog.Errorf(ctx, "读取撤销入库材料失败, materialId=%d: %v", materialId, err)
		return gerror.New("读取材料库存失败，请稍后重试")
	}
	if material == nil || material.IsEmpty() {
		return gerror.New("撤销入库材料不存在或已删除")
	}
	if material["project_id"].Int64() != projectId {
		return gerror.New("采购材料与项目不一致")
	}
	if round2(material["quantity"].Float64()-quantity) < round2(material["used_quantity"].Float64()) {
		return gerror.New("该采购材料已有领取记录，无法删除对应付款")
	}
	if _, err = tx.Exec("UPDATE project_material SET quantity=quantity-? WHERE id=? AND deleted_at IS NULL", round2(quantity), materialId); err != nil {
		glog.Errorf(ctx, "撤销材料入库失败, materialId=%d quantity=%v: %v", materialId, quantity, err)
		return gerror.New("撤销材料入库失败，请稍后重试")
	}
	return nil
}
