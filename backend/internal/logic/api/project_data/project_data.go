package project_data

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// maxImages 单条记录最多关联的图片数（与 project_log 保持一致）
	maxImages = 9
	// maxImagesFieldLength images 字段长度上限
	maxImagesFieldLength = 1000
	// defaultPageSize 默认每页条数
	defaultPageSize = 10
	// maxPageSize 每页最大条数
	maxPageSize = 100
)

// Material 材料清单
type Material struct {
	Id             int64
	CatalogId      int64
	ProjectId      int64
	WarehouseId    int64
	WarehouseName  string
	Name           string
	Spec           string
	Unit           string
	Quantity       float64
	UnitPrice      float64
	TaxMode        string
	TaxRate        float64
	Amount         float64
	UsedQuantity   float64
	RemainQuantity float64
	Supplier       string
	Remark         string
	CreatedAt      string
}

// MaterialInput 材料清单入参
type MaterialInput struct {
	ProjectId   int64
	CatalogId   int64
	WarehouseId int64
	Name        string
	Spec        string
	Unit        string
	Quantity    float64
	UnitPrice   float64
	TaxMode     string
	TaxRate     float64
	Supplier    string
	Remark      string
}

// Usage 材料领用记录
type Usage struct {
	Id            int64
	ProjectId     int64
	MaterialId    int64
	MaterialName  string
	WarehouseId   int64
	WarehouseName string
	Quantity      float64
	UsedAt        string
	Receiver      string
	Remark        string
	StockedAt     string
	CreatedBy     int64
	CreatedAt     string
}

// UsageInput 领用记录入参
type UsageInput struct {
	ProjectId  int64
	MaterialId int64
	Quantity   float64
	UsedAt     string
	Receiver   string
	Remark     string
}

// SiteLog 施工日志
type SiteLog struct {
	Id        int64
	ProjectId int64
	Content   string
	Images    []string
	LogDate   string
	UserName  string
	CreatedAt string
}

// Change 变更记录
type Change struct {
	Id         int64
	ProjectId  int64
	ChangeType string
	Title      string
	Amount     float64
	ChangeDate string
	Reason     string
	Images     []string
	Remark     string
	CreatedAt  string
}

// ChangeInput 变更记录入参
type ChangeInput struct {
	ProjectId  int64
	ChangeType string
	Title      string
	Amount     float64
	ChangeDate string
	Reason     string
	Images     []string
	Remark     string
}

// ChangeListResult 变更记录列表结果
type ChangeListResult struct {
	Records     []*Change
	Total       int
	AddTotal    float64
	ReduceTotal float64
}

type Contract struct {
	Id          int64
	ProjectId   int64
	Name        string
	ContractNo  string
	Amount      float64
	SignedAt    string
	Remark      string
	Attachments []*ContractAttachment
	CreatedAt   string
}

type ContractAttachment struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type ContractInput struct {
	ProjectId   int64
	Name        string
	ContractNo  string
	Amount      float64
	SignedAt    string
	Remark      string
	Attachments []*ContractAttachment
}

// Service 项目资料服务
type Service struct{}

// New 创建项目资料服务
func New() *Service { return &Service{} }

// ===================== 材料清单 =====================

// ListMaterials 材料清单分页列表，附带已领用 / 剩余数量。
func (s *Service) ListMaterials(ctx context.Context, projectId int64, keyword string, current, size int) ([]*Material, int, error) {
	if err := s.ensureProject(ctx, projectId); err != nil {
		return nil, 0, err
	}

	current, size = normalizePage(current, size)

	where := []string{"m.deleted_at IS NULL", "m.company_id=?", "m.project_id=?"}
	args := []interface{}{currentCompanyId(ctx), projectId}
	if text := strings.TrimSpace(keyword); text != "" {
		like := "%" + text + "%"
		where = append(where, "(m.name LIKE ? OR m.spec LIKE ? OR m.supplier LIKE ?)")
		args = append(args, like, like, like)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_material m WHERE "+whereSQL, args...)
	if err != nil {
		glog.Errorf(ctx, "统计材料清单失败, projectId=%d: %v", projectId, err)
		return nil, 0, gerror.New("查询材料清单失败，请稍后重试")
	}
	if total.Int() == 0 {
		return []*Material{}, 0, nil
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			m.id,m.catalog_id,m.project_id,m.warehouse_id,m.name,m.spec,m.unit,m.quantity,m.unit_price,m.tax_mode,m.tax_rate,m.amount,m.supplier,m.remark,
			DATE_FORMAT(m.created_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS created_at,
			COALESCE(u.used_quantity,0) AS used_quantity,
			COALESCE(w.name,'') AS warehouse_name
		FROM project_material m
		LEFT JOIN (
			SELECT company_id,material_id,SUM(quantity) AS used_quantity
			FROM project_material_usage
			WHERE deleted_at IS NULL
			GROUP BY company_id,material_id
		) u ON u.material_id=m.id AND u.company_id=m.company_id
		LEFT JOIN warehouse w ON w.id=m.warehouse_id AND w.deleted_at IS NULL
		WHERE %s
		ORDER BY m.id DESC
		LIMIT ?,?`, whereSQL), queryArgs...)
	if err != nil {
		glog.Errorf(ctx, "查询材料清单失败, projectId=%d: %v", projectId, err)
		return nil, 0, gerror.New("查询材料清单失败，请稍后重试")
	}

	items := make([]*Material, 0, len(records))
	for _, row := range records {
		used := round2(row["used_quantity"].Float64())
		items = append(items, &Material{
			Id:             row["id"].Int64(),
			CatalogId:      row["catalog_id"].Int64(),
			ProjectId:      row["project_id"].Int64(),
			WarehouseId:    row["warehouse_id"].Int64(),
			WarehouseName:  row["warehouse_name"].String(),
			Name:           row["name"].String(),
			Spec:           row["spec"].String(),
			Unit:           row["unit"].String(),
			Quantity:       round2(row["quantity"].Float64()),
			UnitPrice:      round2(row["unit_price"].Float64()),
			TaxMode:        row["tax_mode"].String(),
			TaxRate:        round2(row["tax_rate"].Float64()),
			Amount:         round2(row["amount"].Float64()),
			UsedQuantity:   used,
			RemainQuantity: round2(row["quantity"].Float64() - used),
			Supplier:       row["supplier"].String(),
			Remark:         row["remark"].String(),
			CreatedAt:      row["created_at"].String(),
		})
	}

	return items, total.Int(), nil
}

func (s *Service) applyMaterialCatalog(ctx context.Context, in *MaterialInput) error {
	row, err := g.DB().GetOne(ctx, `SELECT name,spec,unit,unit_price,tax_mode,tax_rate
		FROM material_catalog WHERE id=? AND company_id=? AND deleted_at IS NULL`, in.CatalogId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "读取材料档案失败, catalogId=%d: %v", in.CatalogId, err)
		return gerror.New("读取材料档案失败，请稍后重试")
	}
	if row == nil || row.IsEmpty() {
		return gerror.New("材料档案不存在或已删除")
	}
	in.Name = row["name"].String()
	in.Spec = row["spec"].String()
	in.Unit = row["unit"].String()
	in.UnitPrice = row["unit_price"].Float64()
	in.TaxMode = row["tax_mode"].String()
	in.TaxRate = row["tax_rate"].Float64()
	in.Supplier = ""
	return nil
}

// CreateMaterial 新增材料。
func (s *Service) CreateMaterial(ctx context.Context, in *MaterialInput) (int64, error) {
	if err := s.ensureProject(ctx, in.ProjectId); err != nil {
		return 0, err
	}

	name := strings.TrimSpace(in.Name)
	if name == "" {
		return 0, gerror.New("请填写材料名称")
	}
	if in.CatalogId > 0 {
		if err := s.applyMaterialCatalog(ctx, in); err != nil {
			return 0, err
		}
	}
	if err := ValidateTaxSelection(in.TaxMode, in.TaxRate); err != nil {
		return 0, err
	}
	warehouseId, err := s.resolveWarehouseId(ctx, in.WarehouseId)
	if err != nil {
		return 0, err
	}

	result, err := g.DB().Exec(ctx, `INSERT INTO project_material
		(company_id,project_id,catalog_id,warehouse_id,name,spec,unit,quantity,unit_price,tax_mode,tax_rate,amount,supplier,remark)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		currentCompanyId(ctx), in.ProjectId, in.CatalogId, warehouseId, name, strings.TrimSpace(in.Spec), strings.TrimSpace(in.Unit),
		round2(in.Quantity), round2(in.UnitPrice), in.TaxMode, round2(in.TaxRate), round2(in.Quantity*in.UnitPrice),
		strings.TrimSpace(in.Supplier), strings.TrimSpace(in.Remark))
	if err != nil {
		glog.Errorf(ctx, "新增材料失败, projectId=%d: %v", in.ProjectId, err)
		return 0, gerror.New("新增材料失败，请稍后重试")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.New("新增材料失败，请稍后重试")
	}

	return id, nil
}

// UpdateMaterial 更新材料。
func (s *Service) UpdateMaterial(ctx context.Context, id int64, in *MaterialInput) error {
	if _, err := s.ensureMaterial(ctx, id); err != nil {
		return err
	}
	if err := s.ensureProject(ctx, in.ProjectId); err != nil {
		return err
	}

	name := strings.TrimSpace(in.Name)
	if name == "" {
		return gerror.New("请填写材料名称")
	}
	if in.CatalogId > 0 {
		if err := s.applyMaterialCatalog(ctx, in); err != nil {
			return err
		}
	}
	if err := ValidateTaxSelection(in.TaxMode, in.TaxRate); err != nil {
		return err
	}
	warehouseId, err := s.resolveWarehouseId(ctx, in.WarehouseId)
	if err != nil {
		return err
	}

	_, err = g.DB().Exec(ctx, `UPDATE project_material SET
		project_id=?,catalog_id=?,warehouse_id=?,name=?,spec=?,unit=?,quantity=?,unit_price=?,tax_mode=?,tax_rate=?,amount=?,supplier=?,remark=?
		WHERE id=? AND company_id=? AND deleted_at IS NULL`,
		in.ProjectId, in.CatalogId, warehouseId, name, strings.TrimSpace(in.Spec), strings.TrimSpace(in.Unit),
		round2(in.Quantity), round2(in.UnitPrice), in.TaxMode, round2(in.TaxRate), round2(in.Quantity*in.UnitPrice),
		strings.TrimSpace(in.Supplier), strings.TrimSpace(in.Remark), id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "更新材料失败, id=%d: %v", id, err)
		return gerror.New("更新材料失败，请稍后重试")
	}

	return nil
}

// DeleteMaterial 删除材料（软删除，领用记录保留）。
func (s *Service) DeleteMaterial(ctx context.Context, id int64) error {
	if _, err := s.ensureMaterial(ctx, id); err != nil {
		return err
	}

	if _, err := g.DB().Exec(ctx, "UPDATE project_material SET deleted_at=NOW() WHERE id=? AND company_id=?", id, currentCompanyId(ctx)); err != nil {
		glog.Errorf(ctx, "删除材料失败, id=%d: %v", id, err)
		return gerror.New("删除材料失败，请稍后重试")
	}

	return nil
}

// resolveWarehouseId 解析材料所属仓库：指定了仓库则校验存在性，否则回退到公司默认仓库。
func (s *Service) resolveWarehouseId(ctx context.Context, warehouseId int64) (int64, error) {
	companyId := currentCompanyId(ctx)
	if warehouseId > 0 {
		count, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM warehouse WHERE id=? AND company_id=? AND deleted_at IS NULL",
			warehouseId, companyId)
		if err != nil {
			glog.Errorf(ctx, "校验仓库失败, warehouseId=%d: %v", warehouseId, err)
			return 0, gerror.New("校验仓库失败，请稍后重试")
		}
		if count.Int() == 0 {
			return 0, gerror.New("所选仓库不存在或已删除")
		}
		return warehouseId, nil
	}

	id, err := g.DB().GetValue(ctx,
		"SELECT id FROM warehouse WHERE company_id=? AND deleted_at IS NULL ORDER BY is_default DESC, sort ASC, id ASC LIMIT 1",
		companyId)
	if err != nil {
		glog.Errorf(ctx, "查询默认仓库失败: %v", err)
		return 0, gerror.New("查询仓库失败，请稍后重试")
	}
	if id.IsNil() || id.Int64() <= 0 {
		return 0, gerror.New("尚未创建仓库，请先在「库存查询-仓库管理」中添加")
	}
	return id.Int64(), nil
}

// ===================== 材料领用记录 =====================

// ListUsages 材料领用记录分页列表。
func (s *Service) ListUsages(ctx context.Context, projectId, materialId int64, current, size int) ([]*Usage, int, error) {
	if err := s.ensureProject(ctx, projectId); err != nil {
		return nil, 0, err
	}

	current, size = normalizePage(current, size)

	where := []string{"u.deleted_at IS NULL", "u.company_id=?", "u.project_id=?"}
	args := []interface{}{currentCompanyId(ctx), projectId}
	if materialId > 0 {
		where = append(where, "u.material_id=?")
		args = append(args, materialId)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_material_usage u WHERE "+whereSQL, args...)
	if err != nil {
		glog.Errorf(ctx, "统计领用记录失败, projectId=%d: %v", projectId, err)
		return nil, 0, gerror.New("查询领用记录失败，请稍后重试")
	}
	if total.Int() == 0 {
		return []*Usage{}, 0, nil
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			u.id,u.project_id,u.material_id,u.quantity,u.receiver,u.remark,u.created_by,
			DATE_FORMAT(u.stocked_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS stocked_at,
			DATE_FORMAT(u.used_at,'%%Y-%%m-%%d') AS used_at,
			DATE_FORMAT(u.created_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS created_at,
			COALESCE(m.name,'') AS material_name,
			COALESCE(NULLIF(u.warehouse_id,0),m.warehouse_id,0) AS warehouse_id,
			COALESCE(w.name,'') AS warehouse_name
		FROM project_material_usage u
		LEFT JOIN project_material m ON m.id=u.material_id
		LEFT JOIN warehouse w ON w.id=COALESCE(NULLIF(u.warehouse_id,0),m.warehouse_id) AND w.deleted_at IS NULL
		WHERE %s
		ORDER BY u.used_at DESC,u.id DESC
		LIMIT ?,?`, whereSQL), queryArgs...)
	if err != nil {
		glog.Errorf(ctx, "查询领用记录失败, projectId=%d: %v", projectId, err)
		return nil, 0, gerror.New("查询领用记录失败，请稍后重试")
	}

	items := make([]*Usage, 0, len(records))
	for _, row := range records {
		items = append(items, &Usage{
			Id:            row["id"].Int64(),
			ProjectId:     row["project_id"].Int64(),
			MaterialId:    row["material_id"].Int64(),
			MaterialName:  row["material_name"].String(),
			WarehouseId:   row["warehouse_id"].Int64(),
			WarehouseName: row["warehouse_name"].String(),
			Quantity:      round2(row["quantity"].Float64()),
			UsedAt:        row["used_at"].String(),
			Receiver:      row["receiver"].String(),
			Remark:        row["remark"].String(),
			StockedAt:     row["stocked_at"].String(),
			CreatedBy:     row["created_by"].Int64(),
			CreatedAt:     row["created_at"].String(),
		})
	}

	return items, total.Int(), nil
}

// CreateUsage 新增材料领用记录。createdBy 记录创建人，用于限定只能由本人补做出库。
func (s *Service) CreateUsage(ctx context.Context, in *UsageInput, createdBy uint) (int64, error) {
	if err := s.ensureProject(ctx, in.ProjectId); err != nil {
		return 0, err
	}

	materialProjectId, err := s.ensureMaterial(ctx, in.MaterialId)
	if err != nil {
		return 0, err
	}
	if materialProjectId != in.ProjectId {
		return 0, gerror.New("所选材料不属于当前项目")
	}

	if in.Quantity <= 0 {
		return 0, gerror.New("领用数量必须大于 0")
	}

	usedAt, err := parseDate(in.UsedAt, "请选择领用日期")
	if err != nil {
		return 0, err
	}

	result, err := g.DB().Exec(ctx, `INSERT INTO project_material_usage
		(company_id,project_id,material_id,warehouse_id,quantity,used_at,receiver,remark,created_by)
		VALUES(?,?,?,0,?,?,?,?,?)`,
		currentCompanyId(ctx), in.ProjectId, in.MaterialId, round2(in.Quantity), usedAt,
		strings.TrimSpace(in.Receiver), strings.TrimSpace(in.Remark), createdBy)
	if err != nil {
		glog.Errorf(ctx, "新增领用记录失败, materialId=%d: %v", in.MaterialId, err)
		return 0, gerror.New("新增领用记录失败，请稍后重试")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.New("新增领用记录失败，请稍后重试")
	}

	return id, nil
}

// DeleteUsage 删除领用记录（软删除），同时清理该记录补做的出库流水。
func (s *Service) DeleteUsage(ctx context.Context, id int64) error {
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_material_usage WHERE id=? AND company_id=? AND deleted_at IS NULL", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.New("查询领用记录失败，请稍后重试")
	}
	if count.Int() == 0 {
		return gerror.New("领用记录不存在或已被删除")
	}

	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Exec("UPDATE project_material_usage SET deleted_at=NOW() WHERE id=? AND company_id=?", id, currentCompanyId(ctx)); err != nil {
			glog.Errorf(ctx, "删除领用记录失败, id=%d: %v", id, err)
			return gerror.New("删除领用记录失败，请稍后重试")
		}
		if _, err := tx.Exec("UPDATE inventory_record SET deleted_at=NOW() WHERE company_id=? AND source_type=? AND source_id=? AND deleted_at IS NULL",
			currentCompanyId(ctx), stockSourceUsage, id); err != nil {
			glog.Errorf(ctx, "删除领用记录关联流水失败, usageId=%d: %v", id, err)
			return gerror.New("删除领用记录失败，请稍后重试")
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// ===================== 施工日志 / 现场照片 =====================

// ListSiteLogs 施工日志分页列表（含记录人）。
func (s *Service) ListSiteLogs(ctx context.Context, projectId int64, current, size int) ([]*SiteLog, int, error) {
	if err := s.ensureProject(ctx, projectId); err != nil {
		return nil, 0, err
	}

	current, size = normalizePage(current, size)

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_log WHERE project_id=? AND company_id=?", projectId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "统计施工日志失败, projectId=%d: %v", projectId, err)
		return nil, 0, gerror.New("查询施工日志失败，请稍后重试")
	}
	if total.Int() == 0 {
		return []*SiteLog{}, 0, nil
	}

	records, err := g.DB().GetAll(ctx, `SELECT
			l.id,l.project_id,l.content,l.images,
			DATE_FORMAT(COALESCE(l.log_date,DATE(l.created_at)),'%Y-%m-%d') AS log_date,
			DATE_FORMAT(l.created_at,'%Y-%m-%d %H:%i:%s') AS created_at,
			COALESCE(NULLIF(u.nickname,''),u.username,'') AS user_name
		FROM project_log l
		LEFT JOIN site_user u ON u.id=l.user_id
		WHERE l.project_id=? AND l.company_id=?
		ORDER BY COALESCE(l.log_date,DATE(l.created_at)) DESC, l.id DESC
		LIMIT ?,?`, projectId, currentCompanyId(ctx), (current-1)*size, size)
	if err != nil {
		glog.Errorf(ctx, "查询施工日志失败, projectId=%d: %v", projectId, err)
		return nil, 0, gerror.New("查询施工日志失败，请稍后重试")
	}

	items := make([]*SiteLog, 0, len(records))
	for _, row := range records {
		items = append(items, &SiteLog{
			Id:        row["id"].Int64(),
			ProjectId: row["project_id"].Int64(),
			Content:   row["content"].String(),
			Images:    splitImages(row["images"].String()),
			LogDate:   row["log_date"].String(),
			UserName:  row["user_name"].String(),
			CreatedAt: row["created_at"].String(),
		})
	}

	return items, total.Int(), nil
}

// CreateSiteLog 新增施工日志。logDate 为空取当天（YYYY-MM-DD），补录历史日志时传真实日期。
func (s *Service) CreateSiteLog(ctx context.Context, projectId int64, userId uint, content string, images []string, logDate string) (int64, error) {
	if err := s.ensureProject(ctx, projectId); err != nil {
		return 0, err
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return 0, gerror.New("请填写日志内容")
	}

	logDate = strings.TrimSpace(logDate)
	if logDate == "" {
		logDate = gtime.Now().Format("Y-m-d")
	} else if gtime.New(logDate).Format("Y-m-d") != logDate {
		return 0, gerror.New("日志日期格式不正确，应为 YYYY-MM-DD")
	}

	imagesStr, err := joinImages(images)
	if err != nil {
		return 0, err
	}

	result, err := g.DB().Exec(ctx, `INSERT INTO project_log(company_id,project_id,user_id,content,log_date,images) VALUES(?,?,?,?,?,?)`,
		currentCompanyId(ctx), projectId, userId, content, logDate, imagesStr)
	if err != nil {
		glog.Errorf(ctx, "新增施工日志失败, projectId=%d: %v", projectId, err)
		return 0, gerror.New("新增施工日志失败，请稍后重试")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.New("新增施工日志失败，请稍后重试")
	}

	return id, nil
}

// UpdateSiteLog 更新施工日志。
func (s *Service) UpdateSiteLog(ctx context.Context, id int64, content string, images []string) error {
	projectId, err := s.getLogProjectId(ctx, id)
	if err != nil {
		return err
	}
	if err = s.ensureProject(ctx, projectId); err != nil {
		return err
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return gerror.New("请填写日志内容")
	}

	imagesStr, err := joinImages(images)
	if err != nil {
		return err
	}

	if _, err = g.DB().Exec(ctx, "UPDATE project_log SET content=?,images=? WHERE id=? AND company_id=?", content, imagesStr, id, currentCompanyId(ctx)); err != nil {
		glog.Errorf(ctx, "更新施工日志失败, id=%d: %v", id, err)
		return gerror.New("更新施工日志失败，请稍后重试")
	}

	return nil
}

// DeleteSiteLog 删除施工日志（硬删除，业主端同步不可见）。
func (s *Service) DeleteSiteLog(ctx context.Context, id int64) error {
	projectId, err := s.getLogProjectId(ctx, id)
	if err != nil {
		return err
	}
	if err = s.ensureProject(ctx, projectId); err != nil {
		return err
	}

	if _, err = g.DB().Exec(ctx, "DELETE FROM project_log WHERE id=? AND company_id=?", id, currentCompanyId(ctx)); err != nil {
		glog.Errorf(ctx, "删除施工日志失败, id=%d: %v", id, err)
		return gerror.New("删除施工日志失败，请稍后重试")
	}

	return nil
}

// ===================== 变更记录 =====================

// ListChanges 变更记录分页列表，并返回增项 / 减项合计。
func (s *Service) ListChanges(ctx context.Context, projectId int64, changeType, keyword string, current, size int) (*ChangeListResult, error) {
	if err := s.ensureProject(ctx, projectId); err != nil {
		return nil, err
	}

	current, size = normalizePage(current, size)

	where := []string{"deleted_at IS NULL", "company_id=?", "project_id=?"}
	args := []interface{}{currentCompanyId(ctx), projectId}
	if kind := strings.ToUpper(strings.TrimSpace(changeType)); kind == "ADD" || kind == "REDUCE" {
		where = append(where, "change_type=?")
		args = append(args, kind)
	}
	if text := strings.TrimSpace(keyword); text != "" {
		like := "%" + text + "%"
		where = append(where, "(title LIKE ? OR reason LIKE ?)")
		args = append(args, like, like)
	}
	whereSQL := strings.Join(where, " AND ")

	stat, err := g.DB().GetOne(ctx, fmt.Sprintf(`SELECT
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN change_type='ADD' THEN amount ELSE 0 END),0) AS add_total,
			COALESCE(SUM(CASE WHEN change_type='REDUCE' THEN amount ELSE 0 END),0) AS reduce_total
		FROM project_change WHERE %s`, whereSQL), args...)
	if err != nil {
		glog.Errorf(ctx, "统计变更记录失败, projectId=%d: %v", projectId, err)
		return nil, gerror.New("查询变更记录失败，请稍后重试")
	}

	result := &ChangeListResult{Records: make([]*Change, 0)}
	if stat != nil && !stat.IsEmpty() {
		result.Total = stat["total"].Int()
		result.AddTotal = round2(stat["add_total"].Float64())
		result.ReduceTotal = round2(stat["reduce_total"].Float64())
	}
	if result.Total == 0 {
		return result, nil
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			id,project_id,change_type,title,amount,reason,images,remark,
			DATE_FORMAT(change_date,'%%Y-%%m-%%d') AS change_date,
			DATE_FORMAT(created_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS created_at
		FROM project_change
		WHERE %s
		ORDER BY change_date DESC,id DESC
		LIMIT ?,?`, whereSQL), queryArgs...)
	if err != nil {
		glog.Errorf(ctx, "查询变更记录失败, projectId=%d: %v", projectId, err)
		return nil, gerror.New("查询变更记录失败，请稍后重试")
	}

	for _, row := range records {
		result.Records = append(result.Records, &Change{
			Id:         row["id"].Int64(),
			ProjectId:  row["project_id"].Int64(),
			ChangeType: row["change_type"].String(),
			Title:      row["title"].String(),
			Amount:     round2(row["amount"].Float64()),
			ChangeDate: row["change_date"].String(),
			Reason:     row["reason"].String(),
			Images:     splitImages(row["images"].String()),
			Remark:     row["remark"].String(),
			CreatedAt:  row["created_at"].String(),
		})
	}

	return result, nil
}

// CreateChange 新增变更记录。
func (s *Service) CreateChange(ctx context.Context, in *ChangeInput) (int64, error) {
	if err := s.ensureProject(ctx, in.ProjectId); err != nil {
		return 0, err
	}

	kind, err := normalizeChangeType(in.ChangeType)
	if err != nil {
		return 0, err
	}

	title := strings.TrimSpace(in.Title)
	if title == "" {
		return 0, gerror.New("请填写变更事项")
	}
	if in.Amount < 0 {
		return 0, gerror.New("变更金额不能为负数")
	}

	changeDate, err := parseDate(in.ChangeDate, "请选择变更日期")
	if err != nil {
		return 0, err
	}

	imagesStr, err := joinImages(in.Images)
	if err != nil {
		return 0, err
	}

	result, err := g.DB().Exec(ctx, `INSERT INTO project_change
		(company_id,project_id,change_type,title,amount,change_date,reason,images,remark)
		VALUES(?,?,?,?,?,?,?,?,?)`,
		currentCompanyId(ctx), in.ProjectId, kind, title, round2(in.Amount), changeDate,
		strings.TrimSpace(in.Reason), imagesStr, strings.TrimSpace(in.Remark))
	if err != nil {
		glog.Errorf(ctx, "新增变更记录失败, projectId=%d: %v", in.ProjectId, err)
		return 0, gerror.New("新增变更记录失败，请稍后重试")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.New("新增变更记录失败，请稍后重试")
	}

	return id, nil
}

// UpdateChange 更新变更记录。
func (s *Service) UpdateChange(ctx context.Context, id int64, in *ChangeInput) error {
	if err := s.ensureChange(ctx, id); err != nil {
		return err
	}
	if err := s.ensureProject(ctx, in.ProjectId); err != nil {
		return err
	}

	kind, err := normalizeChangeType(in.ChangeType)
	if err != nil {
		return err
	}

	title := strings.TrimSpace(in.Title)
	if title == "" {
		return gerror.New("请填写变更事项")
	}
	if in.Amount < 0 {
		return gerror.New("变更金额不能为负数")
	}

	changeDate, err := parseDate(in.ChangeDate, "请选择变更日期")
	if err != nil {
		return err
	}

	imagesStr, err := joinImages(in.Images)
	if err != nil {
		return err
	}

	if _, err = g.DB().Exec(ctx, `UPDATE project_change SET
		project_id=?,change_type=?,title=?,amount=?,change_date=?,reason=?,images=?,remark=?
		WHERE id=? AND company_id=? AND deleted_at IS NULL`,
		in.ProjectId, kind, title, round2(in.Amount), changeDate,
		strings.TrimSpace(in.Reason), imagesStr, strings.TrimSpace(in.Remark), id, currentCompanyId(ctx)); err != nil {
		glog.Errorf(ctx, "更新变更记录失败, id=%d: %v", id, err)
		return gerror.New("更新变更记录失败，请稍后重试")
	}

	return nil
}

// DeleteChange 删除变更记录（软删除）。
func (s *Service) DeleteChange(ctx context.Context, id int64) error {
	if err := s.ensureChange(ctx, id); err != nil {
		return err
	}

	if _, err := g.DB().Exec(ctx, "UPDATE project_change SET deleted_at=NOW() WHERE id=? AND company_id=?", id, currentCompanyId(ctx)); err != nil {
		glog.Errorf(ctx, "删除变更记录失败, id=%d: %v", id, err)
		return gerror.New("删除变更记录失败，请稍后重试")
	}

	return nil
}

// ListContracts 查询项目合同。
func (s *Service) ListContracts(ctx context.Context, projectId int64, current, size int) ([]*Contract, int, error) {
	if err := s.ensureProject(ctx, projectId); err != nil {
		return nil, 0, err
	}
	current, size = normalizePage(current, size)
	companyId := currentCompanyId(ctx)
	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_contract WHERE project_id=? AND company_id=? AND deleted_at IS NULL", projectId, companyId)
	if err != nil {
		return nil, 0, gerror.New("查询合同列表失败，请稍后重试")
	}
	if total.Int() == 0 {
		return []*Contract{}, 0, nil
	}
	rows, err := g.DB().GetAll(ctx, `SELECT id,project_id,name,contract_no,amount,
		DATE_FORMAT(signed_at,'%Y-%m-%d') AS signed_at,remark,attachments,
		DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s') AS created_at
		FROM project_contract WHERE project_id=? AND company_id=? AND deleted_at IS NULL
		ORDER BY signed_at DESC,id DESC LIMIT ?,?`, projectId, companyId, (current-1)*size, size)
	if err != nil {
		return nil, 0, gerror.New("查询合同列表失败，请稍后重试")
	}
	items := make([]*Contract, 0, len(rows))
	for _, row := range rows {
		attachments := make([]*ContractAttachment, 0)
		if value := row["attachments"].String(); value != "" {
			if err := json.Unmarshal([]byte(value), &attachments); err != nil {
				glog.Errorf(ctx, "解析合同附件失败, contractId=%d: %v", row["id"].Int64(), err)
			}
		}
		items = append(items, &Contract{
			Id: row["id"].Int64(), ProjectId: row["project_id"].Int64(), Name: row["name"].String(),
			ContractNo: row["contract_no"].String(), Amount: round2(row["amount"].Float64()),
			SignedAt: row["signed_at"].String(), Remark: row["remark"].String(),
			Attachments: attachments, CreatedAt: row["created_at"].String(),
		})
	}
	return items, total.Int(), nil
}

func (s *Service) CreateContract(ctx context.Context, in *ContractInput) (int64, error) {
	if err := validateContractInput(ctx, s, in); err != nil {
		return 0, err
	}
	attachments, err := marshalContractAttachments(in.Attachments)
	if err != nil {
		return 0, err
	}
	result, err := g.DB().Exec(ctx, `INSERT INTO project_contract
		(company_id,project_id,name,contract_no,amount,signed_at,remark,attachments) VALUES(?,?,?,?,?,?,?,?)`,
		currentCompanyId(ctx), in.ProjectId, strings.TrimSpace(in.Name), strings.TrimSpace(in.ContractNo),
		round2(in.Amount), in.SignedAt, strings.TrimSpace(in.Remark), attachments)
	if err != nil {
		glog.Errorf(ctx, "新增合同失败, projectId=%d: %v", in.ProjectId, err)
		return 0, gerror.New("新增合同失败，请稍后重试")
	}
	return result.LastInsertId()
}

func (s *Service) UpdateContract(ctx context.Context, id int64, in *ContractInput) error {
	if err := s.ensureContract(ctx, id); err != nil {
		return err
	}
	if err := validateContractInput(ctx, s, in); err != nil {
		return err
	}
	attachments, err := marshalContractAttachments(in.Attachments)
	if err != nil {
		return err
	}
	_, err = g.DB().Exec(ctx, `UPDATE project_contract SET project_id=?,name=?,contract_no=?,amount=?,signed_at=?,remark=?,attachments=?
		WHERE id=? AND company_id=? AND deleted_at IS NULL`, in.ProjectId, strings.TrimSpace(in.Name),
		strings.TrimSpace(in.ContractNo), round2(in.Amount), in.SignedAt, strings.TrimSpace(in.Remark),
		attachments, id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "更新合同失败, id=%d: %v", id, err)
		return gerror.New("更新合同失败，请稍后重试")
	}
	return nil
}

func (s *Service) DeleteContract(ctx context.Context, id int64) error {
	if err := s.ensureContract(ctx, id); err != nil {
		return err
	}
	if _, err := g.DB().Exec(ctx, "UPDATE project_contract SET deleted_at=NOW() WHERE id=? AND company_id=?", id, currentCompanyId(ctx)); err != nil {
		return gerror.New("删除合同失败，请稍后重试")
	}
	return nil
}

func validateContractInput(ctx context.Context, service *Service, in *ContractInput) error {
	if err := service.ensureProject(ctx, in.ProjectId); err != nil {
		return err
	}
	if strings.TrimSpace(in.Name) == "" {
		return gerror.New("请填写合同名称")
	}
	if in.Amount < 0 {
		return gerror.New("合同金额不能为负数")
	}
	signedAt, err := parseDate(in.SignedAt, "请选择签订日期")
	if err != nil {
		return err
	}
	in.SignedAt = signedAt
	if len(in.Attachments) > 9 {
		return gerror.New("合同附件最多 9 个")
	}
	for _, attachment := range in.Attachments {
		if attachment == nil || strings.TrimSpace(attachment.Path) == "" || strings.TrimSpace(attachment.Name) == "" {
			return gerror.New("合同附件信息不完整")
		}
	}
	return nil
}

func marshalContractAttachments(attachments []*ContractAttachment) (string, error) {
	value, err := json.Marshal(attachments)
	if err != nil {
		return "", gerror.New("合同附件保存失败")
	}
	return string(value), nil
}

func (s *Service) ensureContract(ctx context.Context, id int64) error {
	projectId, err := g.DB().GetValue(ctx, "SELECT project_id FROM project_contract WHERE id=? AND company_id=? AND deleted_at IS NULL", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.New("查询合同失败，请稍后重试")
	}
	if projectId.Int64() == 0 {
		return gerror.New("合同不存在或已被删除")
	}
	return s.ensureProject(ctx, projectId.Int64())
}

// ===================== 内部辅助 =====================

// ensureProject 校验项目存在，并按非超管的可见范围校验权限。
func (s *Service) ensureProject(ctx context.Context, projectId int64) error {
	if projectId <= 0 {
		return gerror.New("请选择项目")
	}

	query := "SELECT COUNT(*) FROM project WHERE id=? AND company_id=? AND deleted_at IS NULL"
	args := []interface{}{projectId, currentCompanyId(ctx)}

	if role := gconv.Uint(ctx.Value(consts.CtxRoleKey)); role != 1 {
		query += " AND (manager_id=? OR customer_id IN (SELECT id FROM customer WHERE owner_user_id=? AND company_id=?))"
		userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))
		args = append(args, userId, userId, currentCompanyId(ctx))
	}

	count, err := g.DB().GetValue(ctx, query, args...)
	if err != nil {
		glog.Errorf(ctx, "校验项目失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验项目失败，请稍后重试")
	}
	if count.Int() == 0 {
		return gerror.New("项目不存在、已删除或无权访问")
	}

	return nil
}

// ensureMaterial 校验材料存在，返回其所属项目ID。
func (s *Service) ensureMaterial(ctx context.Context, id int64) (int64, error) {
	value, err := g.DB().GetValue(ctx, "SELECT project_id FROM project_material WHERE id=? AND company_id=? AND deleted_at IS NULL", id, currentCompanyId(ctx))
	if err != nil {
		return 0, gerror.New("查询材料失败，请稍后重试")
	}
	if value.Int64() == 0 {
		return 0, gerror.New("材料不存在或已被删除")
	}

	return value.Int64(), nil
}

// ensureChange 校验变更记录存在。
func (s *Service) ensureChange(ctx context.Context, id int64) error {
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM project_change WHERE id=? AND company_id=? AND deleted_at IS NULL", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.New("查询变更记录失败，请稍后重试")
	}
	if count.Int() == 0 {
		return gerror.New("变更记录不存在或已被删除")
	}

	return nil
}

// getLogProjectId 查询施工日志所属项目。
func (s *Service) getLogProjectId(ctx context.Context, id int64) (int64, error) {
	value, err := g.DB().GetValue(ctx, "SELECT project_id FROM project_log WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		return 0, gerror.New("查询施工日志失败，请稍后重试")
	}
	if value.Int64() == 0 {
		return 0, gerror.New("施工日志不存在")
	}

	return value.Int64(), nil
}

// normalizePage 统一分页参数。
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

// parseDate 解析 YYYY-MM-DD 日期并返回统一格式。
func parseDate(value, message string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", gerror.New(message)
	}

	for _, layout := range []string{"2006-01-02", "2006/01/02", "2006-01-02 15:04:05"} {
		if date, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return date.Format("2006-01-02"), nil
		}
	}

	return "", gerror.New("日期格式不正确，应为 YYYY-MM-DD")
}

// normalizeChangeType 规范化变更类型。
func normalizeChangeType(value string) (string, error) {
	kind := strings.ToUpper(strings.TrimSpace(value))
	if kind != "ADD" && kind != "REDUCE" {
		return "", gerror.New("请选择变更类型（增项 / 减项）")
	}

	return kind, nil
}

// splitImages 逗号分隔的图片路径 → 数组。
func splitImages(value string) []string {
	result := make([]string, 0)
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}

// joinImages 图片数组 → 逗号分隔字符串（去空、去重、限 9 张）。
func joinImages(images []string) (string, error) {
	result := make([]string, 0, len(images))
	seen := make(map[string]struct{}, len(images))

	for _, item := range images {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}

	if len(result) > maxImages {
		return "", gerror.Newf("单条记录最多上传 %d 张图片", maxImages)
	}

	joined := strings.Join(result, ",")
	if len(joined) > maxImagesFieldLength {
		return "", gerror.New("图片地址过长，请减少图片数量后重试")
	}

	return joined, nil
}

// round2 金额 / 数量保留两位小数。
func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
