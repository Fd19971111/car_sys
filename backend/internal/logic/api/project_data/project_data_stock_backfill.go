package project_data

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// 出入库流水来源单据类型
const (
	stockSourcePurchase = "PURCHASE" // 采购单
	stockSourceUsage    = "USAGE"    // 领用记录
)

// StockInPurchaseOrder 采购单补入库：选择仓库后直接增加库存并生成入库流水。
//
// 已入库（stocked_at 非空）或已生成过流水的单据不允许重复操作；
// 只有创建人本人可以补做，created_by=0 的历史数据不限，用于修正历史错误。
func (s *Service) StockInPurchaseOrder(ctx context.Context, orderId, warehouseId int64, userId uint, remark string) (recordId int64, err error) {
	if orderId <= 0 {
		return 0, gerror.New("采购单ID不能为空")
	}
	if warehouseId <= 0 {
		return 0, gerror.New("请选择仓库")
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		order, orderErr := tx.GetOne(`SELECT id,project_id,material_id,supplier,order_no,quantity,unit_price,
				stocked_at,created_by,DATE_FORMAT(purchase_date,'%Y-%m-%d') AS purchase_date
			FROM project_purchase_order
			WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE`, orderId, currentCompanyId(ctx))
		if orderErr != nil {
			glog.Errorf(ctx, "采购单入库读取单据失败, id=%d: %v", orderId, orderErr)
			return gerror.New("读取采购单失败，请稍后重试")
		}
		if order == nil || order.IsEmpty() {
			return gerror.New("采购单不存在或已被删除")
		}
		if !order["stocked_at"].IsEmpty() {
			return gerror.New("该采购单已完成入库，请勿重复出入库")
		}
		if orderErr = ensureStockOperator(order["created_by"].Int64(), userId); orderErr != nil {
			return orderErr
		}
		if orderErr = ensureWarehouse(ctx, tx, warehouseId); orderErr != nil {
			return orderErr
		}
		if orderErr = ensureNoStockRecord(ctx, tx, stockSourcePurchase, orderId); orderErr != nil {
			return orderErr
		}

		materialId := order["material_id"].Int64()
		material, materialErr := tx.GetOne(`SELECT id,name,spec,unit,unit_price FROM project_material
			WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE`, materialId, currentCompanyId(ctx))
		if materialErr != nil {
			glog.Errorf(ctx, "采购单入库读取材料失败, materialId=%d: %v", materialId, materialErr)
			return gerror.New("读取采购材料失败，请稍后重试")
		}
		if material == nil || material.IsEmpty() {
			return gerror.New("采购单关联的材料不存在或已被删除")
		}

		quantity := round2(order["quantity"].Float64())
		amount := round2(order["amount"].Float64())
		price := 0.0
		if quantity > 0 {
			price = round2(amount / quantity)
		}

		// 入库：材料库存累加（与付款入库口径一致，只加数量），并按所选仓库修正材料归属仓库
		if _, materialErr = tx.Exec(`UPDATE project_material
			SET quantity=quantity+?, warehouse_id=?, updated_at=NOW()
			WHERE id=? AND company_id=? AND deleted_at IS NULL`,
			quantity, warehouseId, materialId, currentCompanyId(ctx)); materialErr != nil {
			glog.Errorf(ctx, "采购单入库更新库存失败, materialId=%d: %v", materialId, materialErr)
			return gerror.New("采购单入库失败，请稍后重试")
		}

		if recordId, materialErr = insertSourceRecord(ctx, tx, &sourceRecordInput{
			ProjectId:   order["project_id"].Int64(),
			MaterialId:  materialId,
			WarehouseId: warehouseId,
			RecordType:  "IN",
			Quantity:    quantity,
			UnitPrice:   price,
			Amount:      amount,
			VendorName: order["supplier"].String(),
			SourceType:  stockSourcePurchase,
			SourceId:    orderId,
			BizDate:     order["purchase_date"].String(),
			Remark:      stockRemark("采购单入库", "采购单 "+order["order_no"].String(), material, quantity, remark),
			CreatedBy:   userId,
		}); materialErr != nil {
			return materialErr
		}

		if _, materialErr = tx.Exec(`UPDATE project_purchase_order SET warehouse_id=?, stocked_at=NOW()
			WHERE id=? AND company_id=?`, warehouseId, orderId, currentCompanyId(ctx)); materialErr != nil {
			glog.Errorf(ctx, "采购单入库标记失败, id=%d: %v", orderId, materialErr)
			return gerror.New("采购单入库失败，请稍后重试")
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return recordId, nil
}

// StockOutMaterialUsage 领用记录补出库：选择仓库后生成出库流水。
//
// 领用记录在登记时已按领用量扣减结存，这里不再重复扣库存，只补齐仓库出库流水，
// 否则结存会被扣两次。已出库或已生成过流水的记录不允许重复操作。
func (s *Service) StockOutMaterialUsage(ctx context.Context, usageId, warehouseId int64, userId uint, remark string) (recordId int64, err error) {
	if usageId <= 0 {
		return 0, gerror.New("领用记录ID不能为空")
	}
	if warehouseId <= 0 {
		return 0, gerror.New("请选择仓库")
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		usage, usageErr := tx.GetOne(`SELECT id,project_id,material_id,quantity,stocked_at,created_by,
				DATE_FORMAT(used_at,'%Y-%m-%d') AS used_at
			FROM project_material_usage
			WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE`, usageId, currentCompanyId(ctx))
		if usageErr != nil {
			glog.Errorf(ctx, "领用出库读取记录失败, id=%d: %v", usageId, usageErr)
			return gerror.New("读取领用记录失败，请稍后重试")
		}
		if usage == nil || usage.IsEmpty() {
			return gerror.New("领用记录不存在或已被删除")
		}
		if !usage["stocked_at"].IsEmpty() {
			return gerror.New("该领用记录已完成出库，请勿重复出入库")
		}
		if usageErr = ensureStockOperator(usage["created_by"].Int64(), userId); usageErr != nil {
			return usageErr
		}
		if usageErr = ensureWarehouse(ctx, tx, warehouseId); usageErr != nil {
			return usageErr
		}
		if usageErr = ensureNoStockRecord(ctx, tx, stockSourceUsage, usageId); usageErr != nil {
			return usageErr
		}

		materialId := usage["material_id"].Int64()
		material, materialErr := tx.GetOne(`SELECT id,name,spec,unit,unit_price FROM project_material
			WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE`, materialId, currentCompanyId(ctx))
		if materialErr != nil {
			glog.Errorf(ctx, "领用出库读取材料失败, materialId=%d: %v", materialId, materialErr)
			return gerror.New("读取领用材料失败，请稍后重试")
		}
		if material == nil || material.IsEmpty() {
			return gerror.New("领用记录关联的材料不存在或已被删除")
		}

		quantity := round2(usage["quantity"].Float64())
		price := round2(material["unit_price"].Float64())
		amount := round2(quantity * price)

		if recordId, materialErr = insertSourceRecord(ctx, tx, &sourceRecordInput{
			ProjectId:   usage["project_id"].Int64(),
			MaterialId:  materialId,
			WarehouseId: warehouseId,
			RecordType:  "OUT",
			Quantity:    quantity,
			UnitPrice:   price,
			Amount:      amount,
			SourceType:  stockSourceUsage,
			SourceId:    usageId,
			BizDate:     usage["used_at"].String(),
			Remark:      stockRemark("领用出库", fmt.Sprintf("领用单 #%d", usageId), material, quantity, remark),
			CreatedBy:   userId,
		}); materialErr != nil {
			return materialErr
		}

		// 出库不重复扣库存（结存已按领用记录扣减），只同步仓库归属与出库标记
		if _, materialErr = tx.Exec(`UPDATE project_material SET warehouse_id=?, updated_at=NOW()
			WHERE id=? AND company_id=? AND deleted_at IS NULL`, warehouseId, materialId, currentCompanyId(ctx)); materialErr != nil {
			glog.Errorf(ctx, "领用出库同步材料仓库失败, materialId=%d: %v", materialId, materialErr)
			return gerror.New("领用出库失败，请稍后重试")
		}
		if _, materialErr = tx.Exec(`UPDATE project_material_usage SET warehouse_id=?, stocked_at=NOW()
			WHERE id=? AND company_id=?`, warehouseId, usageId, currentCompanyId(ctx)); materialErr != nil {
			glog.Errorf(ctx, "领用出库标记失败, id=%d: %v", usageId, materialErr)
			return gerror.New("领用出库失败，请稍后重试")
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return recordId, nil
}

// sourceRecordInput 补做出入库流水的入参
type sourceRecordInput struct {
	ProjectId   int64
	MaterialId  int64
	WarehouseId int64
	RecordType  string
	Quantity    float64
	UnitPrice   float64
	Amount      float64
	VendorName  string
	SourceType  string
	SourceId    int64
	BizDate     string
	Remark      string
	CreatedBy   uint
}

// ensureStockOperator 校验操作人：只有创建人本人可以出入库，created_by=0 的历史数据不限。
func ensureStockOperator(createdBy int64, userId uint) error {
	if createdBy > 0 && createdBy != int64(userId) {
		return gerror.New("只能对本人创建的单据做出入库")
	}

	return nil
}

// ensureWarehouse 校验仓库存在且属于当前公司。
func ensureWarehouse(ctx context.Context, tx gdb.TX, warehouseId int64) error {
	count, err := tx.GetValue("SELECT COUNT(*) FROM warehouse WHERE id=? AND company_id=? AND deleted_at IS NULL",
		warehouseId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "校验仓库失败, warehouseId=%d: %v", warehouseId, err)
		return gerror.New("校验仓库失败，请稍后重试")
	}
	if count.Int() == 0 {
		return gerror.New("仓库不存在或已停用")
	}

	return nil
}

// ensureNoStockRecord 校验该单据是否已生成出入库流水，防止重复出入库。
func ensureNoStockRecord(ctx context.Context, tx gdb.TX, sourceType string, sourceId int64) error {
	count, err := tx.GetValue(`SELECT COUNT(*) FROM inventory_record
		WHERE company_id=? AND source_type=? AND source_id=? AND deleted_at IS NULL`,
		currentCompanyId(ctx), sourceType, sourceId)
	if err != nil {
		glog.Errorf(ctx, "校验出入库流水失败, %s#%d: %v", sourceType, sourceId, err)
		return gerror.New("校验出入库流水失败，请稍后重试")
	}
	if count.Int() > 0 {
		return gerror.New("该单据已生成出入库流水，请勿重复出入库")
	}

	return nil
}

// insertSourceRecord 写入出入库流水并返回流水ID。
func insertSourceRecord(ctx context.Context, tx gdb.TX, in *sourceRecordInput) (int64, error) {
	bizDate, err := parseDate(in.BizDate, "请选择业务日期")
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(`INSERT INTO inventory_record
		(company_id,project_id,material_id,warehouse_id,record_type,biz_type,quantity,unit_price,amount,
		vendor_name,fund_account_id,payout_id,payment_id,source_type,source_id,biz_date,remark,created_by)
		VALUES(?,?,?,?,?,?,?,?,?,?,0,0,0,?,?,?,?,?)`,
		currentCompanyId(ctx), in.ProjectId, in.MaterialId, in.WarehouseId, in.RecordType, "ADJUST",
		in.Quantity, in.UnitPrice, in.Amount, in.VendorName,
		in.SourceType, in.SourceId, bizDate, strings.TrimSpace(in.Remark), in.CreatedBy)
	if err != nil {
		glog.Errorf(ctx, "写入出入库流水失败, source=%s#%d: %v", in.SourceType, in.SourceId, err)
		return 0, gerror.New("保存出入库流水失败，请稍后重试")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.New("保存出入库流水失败，请稍后重试")
	}

	return id, nil
}

// stockRemark 组装流水备注：采购单入库：采购单 POxxx 多层实木板（18mm）× 20 张；备注
func stockRemark(kind, source string, material gdb.Record, quantity float64, remark string) string {
	spec := ""
	if text := strings.TrimSpace(material["spec"].String()); text != "" {
		spec = "（" + text + "）"
	}
	text := fmt.Sprintf("%s：%s %s%s × %g %s", kind, source, material["name"].String(), spec, quantity, material["unit"].String())
	if extra := strings.TrimSpace(remark); extra != "" {
		text += "；" + extra
	}

	return text
}
