package pricing_rule

import (
	"context"
	"strings"

	pricingRuleApi "engineServ/api/pricing_rule"
	v1 "engineServ/api/pricing_rule/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() pricingRuleApi.IPricingRuleV1 { return &Controller{} }

// pricingRuleRow 数据库行：enabled 用整型承接 tinyint，避免 bool 直接扫描出现误差。
type pricingRuleRow struct {
	Id         uint    `json:"id"`
	Name       string  `json:"name"`
	CarType    string  `json:"carType"`
	BasePrice  float64 `json:"basePrice"`
	PerKmPrice float64 `json:"perKmPrice"`
	Enabled    int     `json:"enabled"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	current, size := req.Current, req.Size
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	model := g.DB().Model("pricing_rule").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where("deleted_at IS NULL")
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		model = model.Where("(name LIKE ? OR car_type LIKE ?)", like, like)
	}
	var rows []*pricingRuleRow
	var total int
	if err := model.Page(current, size).OrderAsc("id").ScanAndCount(&rows, &total, false); err != nil {
		return nil, err
	}
	records := make([]*v1.PricingRuleItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, &v1.PricingRuleItem{
			Id:         row.Id,
			Name:       row.Name,
			CarType:    row.CarType,
			BasePrice:  row.BasePrice,
			PerKmPrice: row.PerKmPrice,
			Enabled:    row.Enabled == 1,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
		})
	}
	return &v1.ListRes{Current: current, Size: size, Total: total, Records: records}, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (*v1.SaveRes, error) {
	result, err := g.DB().Model("pricing_rule").Ctx(ctx).Data(g.Map{
		"company_id":   currentCompanyId(ctx),
		"name":         req.Name,
		"car_type":     req.CarType,
		"base_price":   req.BasePrice,
		"per_km_price": req.PerKmPrice,
		"enabled":      boolToInt(req.Enabled),
	}).Insert()
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &v1.SaveRes{Id: uint(id)}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (*v1.SaveRes, error) {
	_, err := g.DB().Model("pricing_rule").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where("id", req.Id).Where("deleted_at IS NULL").Data(g.Map{
		"name":         req.Name,
		"car_type":     req.CarType,
		"base_price":   req.BasePrice,
		"per_km_price": req.PerKmPrice,
		"enabled":      boolToInt(req.Enabled),
	}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.SaveRes{Id: req.Id}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	_, err := g.DB().Model("pricing_rule").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where("id", req.Id).Where("deleted_at IS NULL").Data(g.Map{"deleted_at": gtime.Now()}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (*v1.BatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return &v1.BatchDeleteRes{}, nil
	}
	_, err := g.DB().Model("pricing_rule").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).WhereIn("id", req.Ids).Where("deleted_at IS NULL").Data(g.Map{"deleted_at": gtime.Now()}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.BatchDeleteRes{}, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}
