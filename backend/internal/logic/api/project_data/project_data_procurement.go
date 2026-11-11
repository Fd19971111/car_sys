package project_data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// ListPurchaseOrders 查询采购单。
func (s *Service) ListPurchaseOrders(ctx context.Context, projectId int64, current, size int, pendingOnly bool) ([]*PurchaseOrder, int, error) {
	current, size = normalizePage(current, size)
	where := []string{"po.deleted_at IS NULL"}
	args := make([]interface{}, 0, 3)
	if projectId > 0 {
		where = append(where, "po.project_id=?")
		args = append(args, projectId)
	}
	if pendingOnly {
		where = append(where, "po.status='CONFIRMED'")
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_purchase_order po WHERE "+whereSQL, args...)
	if err != nil {
		glog.Errorf(ctx, "统计采购单失败: %v", err)
		return nil, 0, gerror.New("查询采购单失败，请稍后重试")
	}
	if total.Int() == 0 {
		return []*PurchaseOrder{}, 0, nil
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	rows, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			po.id,po.order_no,po.project_id,COALESCE(p.name,'') AS project_name,po.material_id,po.vendor_id,
			DATE_FORMAT(po.purchase_date,'%%Y-%%m-%%d') AS purchase_date,po.material_name,
			po.spec,po.unit,po.quantity,po.unit_price,po.tax_mode,po.tax_rate,po.net_amount,po.tax_amount,po.amount,po.supplier,po.remark,po.status,po.payout_id,po.created_by,
			DATE_FORMAT(po.stocked_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS stocked_at,
			COALESCE(NULLIF(po.warehouse_id,0),m.warehouse_id,0) AS warehouse_id,
			COALESCE(w.name,'') AS warehouse_name,
			DATE_FORMAT(po.created_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS created_at
		FROM project_purchase_order po
		LEFT JOIN project p ON p.id=po.project_id
		LEFT JOIN project_material m ON m.id=po.material_id
		LEFT JOIN warehouse w ON w.id=COALESCE(NULLIF(po.warehouse_id,0),m.warehouse_id) AND w.deleted_at IS NULL
		WHERE %s
		ORDER BY po.id DESC
		LIMIT ?,?`, whereSQL), queryArgs...)
	if err != nil {
		glog.Errorf(ctx, "查询采购单失败: %v", err)
		return nil, 0, gerror.New("查询采购单失败，请稍后重试")
	}

	return purchaseOrdersFromRows(rows), total.Int(), nil
}

// ListPendingPurchaseOrders 查询付款管理可关联的待付款采购单。
func (s *Service) ListPendingPurchaseOrders(ctx context.Context) ([]*PurchaseOrder, error) {
	rows, err := g.DB().GetAll(ctx, `SELECT
			po.id,po.order_no,po.project_id,COALESCE(p.name,'') AS project_name,po.material_id,po.vendor_id,
			DATE_FORMAT(po.purchase_date,'%Y-%m-%d') AS purchase_date,po.material_name,
			po.spec,po.unit,po.quantity,po.unit_price,po.tax_mode,po.tax_rate,po.net_amount,po.tax_amount,po.amount,po.supplier,po.remark,po.status,po.payout_id,po.created_by,
			DATE_FORMAT(po.stocked_at,'%Y-%m-%d %H:%i:%s') AS stocked_at,
			COALESCE(NULLIF(po.warehouse_id,0),m.warehouse_id,0) AS warehouse_id,
			COALESCE(w.name,'') AS warehouse_name,
			DATE_FORMAT(po.created_at,'%Y-%m-%d %H:%i:%s') AS created_at
		FROM project_purchase_order po
		LEFT JOIN project p ON p.id=po.project_id
		LEFT JOIN project_material m ON m.id=po.material_id
		LEFT JOIN warehouse w ON w.id=COALESCE(NULLIF(po.warehouse_id,0),m.warehouse_id) AND w.deleted_at IS NULL
		WHERE po.deleted_at IS NULL AND po.status='CONFIRMED'
		ORDER BY po.id DESC
		LIMIT 500`)
	if err != nil {
		glog.Errorf(ctx, "查询待付款采购单失败: %v", err)
		return nil, gerror.New("查询待付款采购单失败，请稍后重试")
	}
	return purchaseOrdersFromRows(rows), nil
}

// CreatePurchaseOrder 从项目材料创建待付款采购单。createdBy 记录创建人，用于限定只能由本人补做出入库。
func (s *Service) CreatePurchaseOrder(ctx context.Context, projectId, materialId, vendorId int64, purchaseDate string, quantity, unitPrice float64, taxMode string, taxRate float64, remark string, createdBy uint) (int64, error) {
	if err := s.ensureProject(ctx, projectId); err != nil {
		return 0, err
	}
	materialProjectId, err := s.ensureMaterial(ctx, materialId)
	if err != nil {
		return 0, err
	}
	if materialProjectId != projectId {
		return 0, gerror.New("所选材料不属于当前项目")
	}
	if quantity <= 0 {
		return 0, gerror.New("采购数量必须大于 0")
	}
	if unitPrice < 0 {
		return 0, gerror.New("采购单价不能为负数")
	}
	if err := ValidateTaxSelection(taxMode, taxRate); err != nil {
		return 0, err
	}
	if vendorId <= 0 {
		return 0, gerror.New("请选择供应商")
	}
	purchaseAt, err := parseDate(purchaseDate, "请选择采购日期")
	if err != nil {
		return 0, err
	}
	vendor, err := g.DB().GetOne(ctx, "SELECT name FROM vendor WHERE id=? AND deleted_at IS NULL", vendorId)
	if err != nil {
		glog.Errorf(ctx, "读取采购供应商失败, vendorId=%d: %v", vendorId, err)
		return 0, gerror.New("读取供应商信息失败，请稍后重试")
	}
	if vendor.IsEmpty() {
		return 0, gerror.New("供应商不存在或已停用")
	}

	material, err := g.DB().GetOne(ctx, `SELECT name,spec,unit,supplier FROM project_material WHERE id=? AND deleted_at IS NULL`, materialId)
	if err != nil {
		glog.Errorf(ctx, "读取采购材料失败, materialId=%d: %v", materialId, err)
		return 0, gerror.New("读取材料信息失败，请稍后重试")
	}
	if material.IsEmpty() {
		return 0, gerror.New("材料不存在或已被删除")
	}
	amounts, err := calculatePurchaseTax(quantity, unitPrice, taxMode, taxRate)
	if err != nil {
		return 0, err
	}

	orderNo := "PO" + time.Now().Format("20060102150405.000000000")
	result, err := g.DB().Exec(ctx, `INSERT INTO project_purchase_order
		(order_no,project_id,material_id,vendor_id,warehouse_id,purchase_date,material_name,spec,unit,quantity,unit_price,tax_mode,tax_rate,net_amount,tax_amount,amount,supplier,remark,status,confirmed_at,created_by)
		VALUES(?,?,?,?,0,?,?,?,?,?,?,?,?,?,?,?,?,?,'CONFIRMED',NOW(),?)`,
		orderNo, projectId, materialId, vendorId, purchaseAt, material["name"].String(), material["spec"].String(), material["unit"].String(),
		quantity, unitPrice, taxMode, round2(taxRate), amounts.NetAmount, amounts.TaxAmount, amounts.TotalAmount, vendor["name"].String(), strings.TrimSpace(remark), createdBy)
	if err != nil {
		glog.Errorf(ctx, "创建采购单失败, projectId=%d materialId=%d: %v", projectId, materialId, err)
		return 0, gerror.New("创建采购单失败，请稍后重试")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.New("创建采购单失败，请稍后重试")
	}
	return id, nil
}

// PurchaseOrder 采购单。
type PurchaseOrder struct {
	Id            int64
	OrderNo       string
	ProjectId     int64
	ProjectName   string
	MaterialId    int64
	VendorId      int64
	PurchaseDate  string
	MaterialName  string
	Spec          string
	Unit          string
	Quantity      float64
	UnitPrice     float64
	TaxMode       string
	TaxRate       float64
	NetAmount     float64
	TaxAmount     float64
	Amount        float64
	Supplier      string
	Remark        string
	Status        string
	PayoutId      int64
	WarehouseId   int64
	WarehouseName string
	StockedAt     string
	CreatedBy     int64
	CreatedAt     string
}

func purchaseOrdersFromRows(rows gdb.Result) []*PurchaseOrder {
	items := make([]*PurchaseOrder, 0, len(rows))
	for _, row := range rows {
		items = append(items, &PurchaseOrder{
			Id: row["id"].Int64(), OrderNo: row["order_no"].String(), ProjectId: row["project_id"].Int64(), ProjectName: row["project_name"].String(),
			VendorId: row["vendor_id"].Int64(), PurchaseDate: row["purchase_date"].String(),
			MaterialId: row["material_id"].Int64(), MaterialName: row["material_name"].String(), Spec: row["spec"].String(),
			Unit: row["unit"].String(), Quantity: row["quantity"].Float64(), UnitPrice: row["unit_price"].Float64(),
			TaxMode: row["tax_mode"].String(), TaxRate: row["tax_rate"].Float64(), NetAmount: row["net_amount"].Float64(), TaxAmount: row["tax_amount"].Float64(),
			Amount: row["amount"].Float64(), Supplier: row["supplier"].String(), Remark: row["remark"].String(),
			Status: row["status"].String(), PayoutId: row["payout_id"].Int64(),
			WarehouseId: row["warehouse_id"].Int64(), WarehouseName: row["warehouse_name"].String(),
			StockedAt: row["stocked_at"].String(), CreatedBy: row["created_by"].Int64(),
			CreatedAt: row["created_at"].String(),
		})
	}
	return items
}
