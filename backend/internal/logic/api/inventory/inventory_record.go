package inventory

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	payoutLogic "engineServ/internal/logic/api/payout"
	"engineServ/internal/logic/fin/booking"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// 出入库流水类型 / 业务类型
const (
	recordTypeIn  = "IN"
	recordTypeOut = "OUT"

	bizTypeNormal = "NORMAL" // 业务单据：可选择联动财务
	bizTypeAdjust = "ADJUST" // 补录：只调库存，不生成财务单据
)

// MovementInput 出入库 / 补录通用入参。
//
// 库存口径与既有规则保持一致：
//   - 入库：project_material.quantity 累加（与采购付款入库同源）；
//   - 出库（项目领用）：写入 project_material_usage 领用记录，结存 = quantity - 已领用，不生成财务单据；
//   - 入库联动付款单在同一个事务里完成，任一步失败整体回滚。
type MovementInput struct {
	ProjectId       int64
	MaterialId      int64
	Quantity        float64
	UnitPrice       float64
	GenerateFinance bool // 仅入库使用：是否同步生成付款单
	VendorName      string
	FundAccountId   int64
	RecordType      string // IN / OUT（仅补录使用）
	Receiver        string
	BizDate         string
	Remark          string
}

// RecordItem 出入库流水行
type RecordItem struct {
	Id               int64
	RecordType       string
	BizType          string
	ProjectName      string
	MaterialName     string
	Spec             string
	Unit             string
	WarehouseName    string
	Quantity         float64
	UnitPrice        float64
	Amount           float64
	VendorName       string
	FundAccountName  string
	PayoutId         int64
	PaymentId        int64
	PayoutVoucherNo  string
	PaymentVoucherNo string
	BizDate          string
	Remark           string
	CreatedAt        string
}

// materialRow 事务内锁定的材料行
type materialRow struct {
	Id          int64
	ProjectId   int64
	WarehouseId int64
	Name        string
	Spec        string
	Unit        string
	Quantity    float64
	UnitPrice   float64
}

// StockIn 材料入库：库存数量累加，可选生成付款单（材料款）。
//
// 库存变更、流水记录、付款单与凭证在同一个事务里写入：
// 付款单复用 payout 模块的创建逻辑（含凭证记账、超付校验），
// GoFrame 嵌套事务保证任一步失败时整体回滚。
func (s *Service) StockIn(ctx context.Context, in *MovementInput, userId uint) (recordId, payoutId int64, voucherNo string, err error) {
	if err = validateFinanceOption(ctx, in); err != nil {
		return 0, 0, "", err
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		mat, matErr := lockMaterial(ctx, tx, in.MaterialId)
		if matErr != nil {
			return matErr
		}

		price := in.UnitPrice
		if price <= 0 {
			price = mat.UnitPrice
		}
		amount := round2(in.Quantity * price)

		if recordId, matErr = insertRecord(ctx, tx, mat, recordTypeIn, bizTypeNormal, in.Quantity, price, amount,
			in.VendorName, in.FundAccountId, in.BizDate, in.Remark, userId); matErr != nil {
			return matErr
		}

		if _, matErr = tx.Exec("UPDATE project_material SET quantity=quantity+?, amount=amount+? WHERE id=?",
			in.Quantity, amount, mat.Id); matErr != nil {
			glog.Errorf(ctx, "入库更新库存失败, materialId=%d: %v", mat.Id, matErr)
			return gerror.New("材料入库失败，请稍后重试")
		}

		if in.GenerateFinance {
			payoutId, voucherNo, matErr = payoutLogic.New().Create(ctx, &payoutLogic.Form{
				ProjectId:     mat.ProjectId,
				Amount:        amount,
				PayoutType:    "材料款",
				PayoutDate:    in.BizDate,
				VendorName:    in.VendorName,
				FundAccountId: in.FundAccountId,
				Remark:        movementRemark("材料入库", mat, in.Quantity, in.Remark),
			}, userId)
			if matErr != nil {
				return matErr
			}
			if _, matErr = tx.Exec("UPDATE inventory_record SET payout_id=? WHERE id=?", payoutId, recordId); matErr != nil {
				glog.Errorf(ctx, "写回入库流水付款单ID失败, recordId=%d: %v", recordId, matErr)
				return gerror.New("材料入库失败，请稍后重试")
			}
		}

		return nil
	})
	if err != nil {
		return 0, 0, "", err
	}

	return recordId, payoutId, voucherNo, nil
}

// StockOut 材料出库（项目领用）：只扣减库存并写入领用记录，不生成任何财务单据。
//
// 「一进一出管库存」：领用属于项目内部耗用，库存减少的同时计入项目材料成本，
// 既不产生收款单也不产生付款单，更不生成记账凭证。
func (s *Service) StockOut(ctx context.Context, in *MovementInput, userId uint) (recordId int64, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		mat, matErr := lockMaterial(ctx, tx, in.MaterialId)
		if matErr != nil {
			return matErr
		}

		remain, matErr := materialRemain(ctx, tx, mat)
		if matErr != nil {
			return matErr
		}
		if in.Quantity > remain+1e-9 {
			return gerror.Newf("出库数量超过当前结存 %.2f %s", remain, mat.Unit)
		}

		// 领用成本单价：优先使用出库时手动填写的单价，未填则按材料清单单价
		price := in.UnitPrice
		if price <= 0 {
			price = mat.UnitPrice
		}
		amount := round2(in.Quantity * price)

		if recordId, matErr = insertRecord(ctx, tx, mat, recordTypeOut, bizTypeNormal, in.Quantity, price, amount,
			"", 0, in.BizDate, in.Remark, userId); matErr != nil {
			return matErr
		}

		return insertUsage(ctx, tx, mat, in, "材料出库")
	})
	if err != nil {
		return 0, err
	}

	return recordId, nil
}

// StockAdjust 库存补录：漏补历史数据时直接调整库存，不生成任何财务单据。
// IN 方向累加库存；OUT 方向写入领用记录并扣减结存。
func (s *Service) StockAdjust(ctx context.Context, in *MovementInput, userId uint) (recordId int64, err error) {
	recordType := strings.ToUpper(strings.TrimSpace(in.RecordType))
	if recordType != recordTypeIn && recordType != recordTypeOut {
		return 0, gerror.New("请选择补录方向")
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		mat, matErr := lockMaterial(ctx, tx, in.MaterialId)
		if matErr != nil {
			return matErr
		}

		kind := "补录入库"
		if recordType == recordTypeOut {
			kind = "补录出库"
		}

		if recordType == recordTypeIn {
			price := in.UnitPrice
			if price <= 0 {
				price = mat.UnitPrice
			}
			amount := round2(in.Quantity * price)

			if recordId, matErr = insertRecord(ctx, tx, mat, recordTypeIn, bizTypeAdjust, in.Quantity, price, amount,
				"", 0, in.BizDate, in.Remark, userId); matErr != nil {
				return matErr
			}
			if _, matErr = tx.Exec("UPDATE project_material SET quantity=quantity+?, amount=amount+? WHERE id=?",
				in.Quantity, amount, mat.Id); matErr != nil {
				glog.Errorf(ctx, "补录入库更新库存失败, materialId=%d: %v", mat.Id, matErr)
				return gerror.New("库存补录失败，请稍后重试")
			}

			return nil
		}

		remain, matErr := materialRemain(ctx, tx, mat)
		if matErr != nil {
			return matErr
		}
		if in.Quantity > remain+1e-9 {
			return gerror.Newf("补录出库数量超过当前结存 %.2f %s", remain, mat.Unit)
		}

		amount := round2(in.Quantity * mat.UnitPrice)
		if recordId, matErr = insertRecord(ctx, tx, mat, recordTypeOut, bizTypeAdjust, in.Quantity, mat.UnitPrice, amount,
			"", 0, in.BizDate, in.Remark, userId); matErr != nil {
			return matErr
		}

		return insertUsage(ctx, tx, mat, in, kind)
	})
	if err != nil {
		return 0, err
	}

	return recordId, nil
}

// ListRecords 出入库流水分页列表。
func (s *Service) ListRecords(ctx context.Context, projectId, materialId int64, recordType string, current, size int) ([]*RecordItem, int, error) {
	current, size = normalizePage(current, size)

	where := []string{"r.deleted_at IS NULL", "r.company_id=?"}
	args := []interface{}{currentCompanyId(ctx)}
	if projectId > 0 {
		where = append(where, "r.project_id=?")
		args = append(args, projectId)
	}
	if materialId > 0 {
		where = append(where, "r.material_id=?")
		args = append(args, materialId)
	}
	switch strings.ToUpper(strings.TrimSpace(recordType)) {
	case recordTypeIn, recordTypeOut:
		where = append(where, "r.record_type=?")
		args = append(args, strings.ToUpper(recordType))
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM inventory_record r WHERE "+whereSQL, args...)
	if err != nil {
		glog.Errorf(ctx, "统计出入库流水失败: %v", err)
		return nil, 0, gerror.New("查询出入库流水失败，请稍后重试")
	}
	if total.Int() == 0 {
		return []*RecordItem{}, 0, nil
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			r.id,r.record_type,r.biz_type,r.quantity,r.unit_price,r.amount,r.payout_id,r.payment_id,r.remark,
			DATE_FORMAT(r.biz_date,'%%Y-%%m-%%d') AS biz_date,
			DATE_FORMAT(r.created_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS created_at,
			COALESCE(p.name,'') AS project_name,
			COALESCE(m.name,'') AS material_name,
			COALESCE(m.spec,'') AS spec,
			COALESCE(m.unit,'') AS unit,
			COALESCE(w.name,'') AS warehouse_name,
			COALESCE(r.vendor_name,'') AS vendor_name,
			COALESCE(fa.account_name,'') AS fund_account_name,
			COALESCE(po.voucher_id,0) AS payout_voucher_id,
			COALESCE(pm.voucher_id,0) AS payment_voucher_id
		FROM inventory_record r
		LEFT JOIN project p ON p.id=r.project_id
		LEFT JOIN project_material m ON m.id=r.material_id
		LEFT JOIN warehouse w ON w.id=r.warehouse_id AND w.deleted_at IS NULL
		LEFT JOIN fin_fund_account fa ON fa.id=r.fund_account_id
		LEFT JOIN project_payout po ON po.id=r.payout_id AND po.deleted_at IS NULL
		LEFT JOIN project_payment pm ON pm.id=r.payment_id AND pm.deleted_at IS NULL
		WHERE %s
		ORDER BY r.id DESC
		LIMIT ?,?`, whereSQL), queryArgs...)
	if err != nil {
		glog.Errorf(ctx, "查询出入库流水失败: %v", err)
		return nil, 0, gerror.New("查询出入库流水失败，请稍后重试")
	}

	// 汇总流水涉及的凭证号，一次性批量查询
	voucherIds := make([]int64, 0, len(records)*2)
	for _, row := range records {
		if id := row["payout_voucher_id"].Int64(); id > 0 {
			voucherIds = append(voucherIds, id)
		}
		if id := row["payment_voucher_id"].Int64(); id > 0 {
			voucherIds = append(voucherIds, id)
		}
	}
	voucherNos := map[int64]string{}
	if len(voucherIds) > 0 {
		if values, err := booking.New().VoucherNos(ctx, voucherIds); err != nil {
			// 凭证号仅用于展示，查询失败不阻断列表
			glog.Errorf(ctx, "查询出入库流水凭证号失败: %v", err)
		} else {
			voucherNos = values
		}
	}

	items := make([]*RecordItem, 0, len(records))
	for _, row := range records {
		payoutVoucherId := row["payout_voucher_id"].Int64()
		paymentVoucherId := row["payment_voucher_id"].Int64()
		items = append(items, &RecordItem{
			Id:               row["id"].Int64(),
			RecordType:       row["record_type"].String(),
			BizType:          row["biz_type"].String(),
			ProjectName:      row["project_name"].String(),
			MaterialName:     row["material_name"].String(),
			Spec:             row["spec"].String(),
			Unit:             row["unit"].String(),
			WarehouseName:    row["warehouse_name"].String(),
			Quantity:         round2(row["quantity"].Float64()),
			UnitPrice:        round2(row["unit_price"].Float64()),
			Amount:           round2(row["amount"].Float64()),
			VendorName:       row["vendor_name"].String(),
			FundAccountName:  row["fund_account_name"].String(),
			PayoutId:         row["payout_id"].Int64(),
			PaymentId:        row["payment_id"].Int64(),
			PayoutVoucherNo:  voucherNos[payoutVoucherId],
			PaymentVoucherNo: voucherNos[paymentVoucherId],
			BizDate:          row["biz_date"].String(),
			Remark:           row["remark"].String(),
			CreatedAt:        row["created_at"].String(),
		})
	}

	return items, total.Int(), nil
}

// ===================== 内部辅助 =====================

// validateFinanceOption 校验入库生成付款单时的必填条件。
func validateFinanceOption(ctx context.Context, in *MovementInput) error {
	if !in.GenerateFinance {
		return nil
	}

	if strings.TrimSpace(in.VendorName) == "" {
		return gerror.New("生成付款单时请填写供应商")
	}
	if in.FundAccountId <= 0 {
		return gerror.New("生成付款单时请选择资金账户")
	}

	return nil
}

// lockMaterial 事务内锁定并返回材料行。
func lockMaterial(ctx context.Context, tx gdb.TX, materialId int64) (*materialRow, error) {
	row, err := tx.GetOne(`SELECT id,project_id,warehouse_id,name,spec,unit,quantity,unit_price
		FROM project_material WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE`,
		materialId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "读取材料失败, materialId=%d: %v", materialId, err)
		return nil, gerror.New("读取材料信息失败，请稍后重试")
	}
	if row == nil || row.IsEmpty() {
		return nil, gerror.New("材料不存在或已被删除")
	}

	return &materialRow{
		Id:          row["id"].Int64(),
		ProjectId:   row["project_id"].Int64(),
		WarehouseId: row["warehouse_id"].Int64(),
		Name:        row["name"].String(),
		Spec:        row["spec"].String(),
		Unit:        row["unit"].String(),
		Quantity:    row["quantity"].Float64(),
		UnitPrice:   row["unit_price"].Float64(),
	}, nil
}

// materialRemain 计算材料当前结存 = quantity - 已领用。
func materialRemain(ctx context.Context, tx gdb.TX, mat *materialRow) (float64, error) {
	value, err := tx.GetValue(`SELECT COALESCE(SUM(quantity),0) FROM project_material_usage
		WHERE material_id=? AND company_id=? AND deleted_at IS NULL`, mat.Id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "统计材料已领用失败, materialId=%d: %v", mat.Id, err)
		return 0, gerror.New("校验库存结存失败，请稍后重试")
	}

	return round2(mat.Quantity - value.Float64()), nil
}

// insertUsage 写入领用记录（出库与补录出库共用，保持结存统计同源）。
func insertUsage(ctx context.Context, tx gdb.TX, mat *materialRow, in *MovementInput, kind string) error {
	receiver := strings.TrimSpace(in.Receiver)
	if receiver == "" {
		receiver = "补录"
	}

	usedAt, err := parseBizDate(in.BizDate)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO project_material_usage
		(company_id,project_id,material_id,quantity,used_at,receiver,remark)
		VALUES(?,?,?,?,?,?,?)`,
		currentCompanyId(ctx), mat.ProjectId, mat.Id, round2(in.Quantity), usedAt.Format("2006-01-02"),
		receiver, movementRemark(kind, mat, in.Quantity, in.Remark)); err != nil {
		glog.Errorf(ctx, "写入领用记录失败, materialId=%d: %v", mat.Id, err)
		return gerror.New("保存领用记录失败，请稍后重试")
	}

	return nil
}

// insertRecord 写入出入库流水。
func insertRecord(ctx context.Context, tx gdb.TX, mat *materialRow, recordType, bizType string,
	quantity, unitPrice, amount float64, vendorName string, fundAccountId int64, bizDate, remark string, userId uint,
) (int64, error) {
	date, err := parseBizDate(bizDate)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(`INSERT INTO inventory_record
		(company_id,project_id,material_id,warehouse_id,record_type,biz_type,quantity,unit_price,amount,
		vendor_name,fund_account_id,biz_date,remark,created_by)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		currentCompanyId(ctx), mat.ProjectId, mat.Id, mat.WarehouseId, recordType, bizType,
		round2(quantity), round2(unitPrice), round2(amount), strings.TrimSpace(vendorName), fundAccountId,
		date.Format("2006-01-02"), strings.TrimSpace(remark), userId)
	if err != nil {
		glog.Errorf(ctx, "写入出入库流水失败, materialId=%d: %v", mat.Id, err)
		return 0, gerror.New("保存出入库流水失败，请稍后重试")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.New("保存出入库流水失败，请稍后重试")
	}

	return id, nil
}

// movementRemark 组装财务单据 / 领用记录备注：材料入库：多层实木板（18mm）× 20 张；备注
func movementRemark(kind string, mat *materialRow, quantity float64, remark string) string {
	spec := ""
	if strings.TrimSpace(mat.Spec) != "" {
		spec = "（" + mat.Spec + "）"
	}
	text := fmt.Sprintf("%s：%s%s × %s %s", kind, mat.Name, spec, formatQuantity(quantity), mat.Unit)
	if text2 := strings.TrimSpace(remark); text2 != "" {
		text += "；" + text2
	}
	return text
}

// formatQuantity 数量展示：整数不带小数，小数保留有效位。
func formatQuantity(value float64) string {
	if value == math.Trunc(value) {
		return fmt.Sprintf("%d", int64(value))
	}
	return fmt.Sprintf("%g", value)
}

// parseBizDate 解析 YYYY-MM-DD 业务日期。
func parseBizDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, gerror.New("请选择业务日期")
	}

	date, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		if date, err = time.ParseInLocation("2006/01/02", value, time.Local); err != nil {
			if date, err = time.ParseInLocation("2006-01-02 15:04:05", value, time.Local); err != nil {
				return nil, gerror.New("业务日期格式不正确")
			}
		}
	}

	return &date, nil
}
