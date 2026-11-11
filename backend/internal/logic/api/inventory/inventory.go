package inventory

import (
	"context"
	"fmt"
	"math"
	"strings"

	"engineServ/internal/consts"
	projectDataLogic "engineServ/internal/logic/api/project_data"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// Service 库存管理服务。
//
// 库存的唯一数据源是 project_material.quantity（累计入库，采购付款时增加），
// 出库量实时汇总自 project_material_usage（领用记录），
// 因此本模块只做展示与台账维护，不单独记录出入库流水。
type Service struct{}

// New 创建库存管理服务。
func New() *Service { return &Service{} }

// StockItem 库存行
type StockItem struct {
	Id             int64
	ProjectId      int64
	ProjectName    string
	WarehouseId    int64
	WarehouseName  string
	Name           string
	Spec           string
	Unit           string
	Quantity       float64
	UnitPrice      float64
	Amount         float64
	UsedQuantity   float64
	RemainQuantity float64
	Supplier       string
	Remark         string
	UpdatedAt      string
}

// StockSummary 库存汇总
type StockSummary struct {
	Kinds    int
	Amount   float64
	Shortage int
}

// ListStock 库存分页列表（跨项目），附带汇总统计。warehouseId 为 0 时表示全部仓库。
func (s *Service) ListStock(ctx context.Context, projectId, warehouseId int64, keyword string, current, size int) ([]*StockItem, *StockSummary, int, error) {
	current, size = normalizePage(current, size)

	where := []string{"m.deleted_at IS NULL", "m.company_id=?"}
	args := []interface{}{currentCompanyId(ctx)}
	if projectId > 0 {
		where = append(where, "m.project_id=?")
		args = append(args, projectId)
	}
	if warehouseId > 0 {
		where = append(where, "m.warehouse_id=?")
		args = append(args, warehouseId)
	}
	if text := strings.TrimSpace(keyword); text != "" {
		like := "%" + text + "%"
		where = append(where, "(m.name LIKE ? OR m.spec LIKE ? OR m.supplier LIKE ?)")
		args = append(args, like, like, like)
	}
	whereSQL := strings.Join(where, " AND ")

	joinSQL := `LEFT JOIN (
			SELECT company_id,material_id,SUM(quantity) AS used_quantity
			FROM project_material_usage
			WHERE deleted_at IS NULL
			GROUP BY company_id,material_id
		) u ON u.material_id=m.id AND u.company_id=m.company_id`

	summaryRow, err := g.DB().GetOne(ctx, fmt.Sprintf(`SELECT
			COUNT(*) AS kinds,
			COALESCE(SUM(m.amount),0) AS amount,
			COALESCE(SUM(CASE WHEN m.quantity-COALESCE(u.used_quantity,0)<=0 THEN 1 ELSE 0 END),0) AS shortage
		FROM project_material m %s WHERE %s`, joinSQL, whereSQL), args...)
	if err != nil {
		glog.Errorf(ctx, "统计库存汇总失败: %v", err)
		return nil, nil, 0, gerror.New("查询库存失败，请稍后重试")
	}
	summary := &StockSummary{
		Kinds:    summaryRow["kinds"].Int(),
		Amount:   round2(summaryRow["amount"].Float64()),
		Shortage: summaryRow["shortage"].Int(),
	}
	if summary.Kinds == 0 {
		return []*StockItem{}, summary, 0, nil
	}

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_material m WHERE "+whereSQL, args...)
	if err != nil {
		glog.Errorf(ctx, "统计库存总数失败: %v", err)
		return nil, nil, 0, gerror.New("查询库存失败，请稍后重试")
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			m.id,m.project_id,m.warehouse_id,m.name,m.spec,m.unit,m.quantity,m.unit_price,m.amount,m.supplier,m.remark,
			DATE_FORMAT(m.updated_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS updated_at,
			COALESCE(u.used_quantity,0) AS used_quantity,
			COALESCE(p.name,'') AS project_name,
			COALESCE(w.name,'') AS warehouse_name
		FROM project_material m
		%s
		LEFT JOIN project p ON p.id=m.project_id
		LEFT JOIN warehouse w ON w.id=m.warehouse_id AND w.deleted_at IS NULL
		WHERE %s
		ORDER BY m.updated_at DESC,m.id DESC
		LIMIT ?,?`, joinSQL, whereSQL), queryArgs...)
	if err != nil {
		glog.Errorf(ctx, "查询库存列表失败: %v", err)
		return nil, nil, 0, gerror.New("查询库存失败，请稍后重试")
	}

	items := make([]*StockItem, 0, len(records))
	for _, row := range records {
		used := round2(row["used_quantity"].Float64())
		items = append(items, &StockItem{
			Id:             row["id"].Int64(),
			ProjectId:      row["project_id"].Int64(),
			ProjectName:    row["project_name"].String(),
			WarehouseId:    row["warehouse_id"].Int64(),
			WarehouseName:  row["warehouse_name"].String(),
			Name:           row["name"].String(),
			Spec:           row["spec"].String(),
			Unit:           row["unit"].String(),
			Quantity:       round2(row["quantity"].Float64()),
			UnitPrice:      round2(row["unit_price"].Float64()),
			Amount:         round2(row["amount"].Float64()),
			UsedQuantity:   used,
			RemainQuantity: round2(row["quantity"].Float64() - used),
			Supplier:       row["supplier"].String(),
			Remark:         row["remark"].String(),
			UpdatedAt:      row["updated_at"].String(),
		})
	}

	return items, summary, total.Int(), nil
}

// CreateStock 新增库存（复用项目资料的材料写入，保持公司 / 项目归属校验一致）。
func (s *Service) CreateStock(ctx context.Context, in *projectDataLogic.MaterialInput) (int64, error) {
	return projectDataLogic.New().CreateMaterial(ctx, in)
}

// UpdateStock 更新库存。
func (s *Service) UpdateStock(ctx context.Context, id int64, in *projectDataLogic.MaterialInput) error {
	return projectDataLogic.New().UpdateMaterial(ctx, id, in)
}

// DeleteStock 删除库存（软删除，领用记录保留）。
func (s *Service) DeleteStock(ctx context.Context, id int64) error {
	return projectDataLogic.New().DeleteMaterial(ctx, id)
}

func normalizePage(current, size int) (int, int) {
	if current < 1 {
		current = 1
	}
	if size < 1 || size > maxPageSize {
		size = defaultPageSize
	}

	return current, size
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
