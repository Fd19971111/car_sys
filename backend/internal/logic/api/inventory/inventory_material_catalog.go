package inventory

import (
	"context"
	"strings"

	v1 "engineServ/api/inventory/v1"
	projectDataLogic "engineServ/internal/logic/api/project_data"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ListMaterialCatalog 查询公司级材料档案。
func (s *Service) ListMaterialCatalog(ctx context.Context, keyword string, current, size int) ([]*v1.MaterialCatalogItem, int, error) {
	current, size = normalizePage(current, size)
	model := g.DB().Model("material_catalog").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL")
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		model = model.Where("(name LIKE ? OR spec LIKE ?)", like, like)
	}
	var records []*v1.MaterialCatalogItem
	var total int
	if err := model.Page(current, size).Order("updated_at DESC,id DESC").ScanAndCount(&records, &total, false); err != nil {
		return nil, 0, gerror.New("查询材料档案失败，请稍后重试")
	}
	return records, total, nil
}

// CreateMaterialCatalog 新增公司级材料档案。
func (s *Service) CreateMaterialCatalog(ctx context.Context, form *v1.MaterialCatalogForm) (int64, error) {
	if err := validateMaterialCatalog(form); err != nil {
		return 0, err
	}
	companyId := currentCompanyId(ctx)
	if err := ensureMaterialCatalogUnique(ctx, companyId, form.Name, form.Spec, form.Unit, 0); err != nil {
		return 0, err
	}
	result, err := g.DB().Model("material_catalog").Ctx(ctx).Data(g.Map{
		"company_id": companyId,
		"name":       strings.TrimSpace(form.Name),
		"spec":       strings.TrimSpace(form.Spec),
		"unit":       strings.TrimSpace(form.Unit),
		"unit_price": round2(form.UnitPrice),
		"tax_mode":   form.TaxMode,
		"tax_rate":   round2(form.TaxRate),
		"remark":     strings.TrimSpace(form.Remark),
	}).Insert()
	if err != nil {
		return 0, gerror.New("新增材料档案失败，请稍后重试")
	}
	return result.LastInsertId()
}

// UpdateMaterialCatalog 更新公司级材料档案。
func (s *Service) UpdateMaterialCatalog(ctx context.Context, id int64, form *v1.MaterialCatalogForm) error {
	if err := validateMaterialCatalog(form); err != nil {
		return err
	}
	companyId := currentCompanyId(ctx)
	if err := ensureMaterialCatalogExists(ctx, companyId, id); err != nil {
		return err
	}
	if err := ensureMaterialCatalogUnique(ctx, companyId, form.Name, form.Spec, form.Unit, id); err != nil {
		return err
	}
	_, err := g.DB().Model("material_catalog").Ctx(ctx).
		Where("id", id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		Data(g.Map{
			"name":       strings.TrimSpace(form.Name),
			"spec":       strings.TrimSpace(form.Spec),
			"unit":       strings.TrimSpace(form.Unit),
			"unit_price": round2(form.UnitPrice),
			"tax_mode":   form.TaxMode,
			"tax_rate":   round2(form.TaxRate),
			"remark":     strings.TrimSpace(form.Remark),
		}).Update()
	return err
}

// DeleteMaterialCatalog 软删除未被项目材料引用的档案。
func (s *Service) DeleteMaterialCatalog(ctx context.Context, id int64) error {
	companyId := currentCompanyId(ctx)
	if err := ensureMaterialCatalogExists(ctx, companyId, id); err != nil {
		return err
	}
	count, err := g.DB().Model("project_material").Ctx(ctx).
		Where("catalog_id", id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").Count()
	if err != nil {
		return gerror.New("校验材料档案引用失败，请稍后重试")
	}
	if count > 0 {
		return gerror.New("该材料档案已用于项目材料，不能删除")
	}
	_, err = g.DB().Model("material_catalog").Ctx(ctx).
		Where("id", id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now()}).Update()
	return err
}

func validateMaterialCatalog(form *v1.MaterialCatalogForm) error {
	if strings.TrimSpace(form.Name) == "" {
		return gerror.New("请填写材料名称")
	}
	if form.UnitPrice < 0 {
		return gerror.New("单价不能为负数")
	}
	if err := projectDataLogic.ValidateTaxSelection(form.TaxMode, form.TaxRate); err != nil {
		return err
	}
	return nil
}

func ensureMaterialCatalogExists(ctx context.Context, companyId uint, id int64) error {
	count, err := g.DB().Model("material_catalog").Ctx(ctx).
		Where("id", id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").Count()
	if err != nil {
		return gerror.New("查询材料档案失败，请稍后重试")
	}
	if count == 0 {
		return gerror.New("材料档案不存在或已删除")
	}
	return nil
}

func ensureMaterialCatalogUnique(ctx context.Context, companyId uint, name, spec, unit string, excludeId int64) error {
	count, err := g.DB().Model("material_catalog").Ctx(ctx).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		Where("name", strings.TrimSpace(name)).
		Where("spec", strings.TrimSpace(spec)).
		Where("unit", strings.TrimSpace(unit)).
		WhereNot("id", excludeId).Count()
	if err != nil {
		return gerror.New("校验材料名称规格失败，请稍后重试")
	}
	if count > 0 {
		return gerror.New("相同名称和规格的材料档案已存在")
	}
	return nil
}
