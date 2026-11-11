package vendor

import (
	"context"
	"strings"

	vendorApi "engineServ/api/supplier"
	v1 "engineServ/api/supplier/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() vendorApi.ISupplierV1 { return &Controller{} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	current, size := req.Current, req.Size
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	model := g.DB().Model("vendor").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where("deleted_at IS NULL")
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		model = model.Where("(name LIKE ? OR contact_name LIKE ? OR phone LIKE ? OR address LIKE ?)", like, like, like, like)
	}
	var records []*v1.VendorItem
	var total int
	if err := model.Page(current, size).OrderDesc("created_at").OrderDesc("id").ScanAndCount(&records, &total, false); err != nil {
		return nil, err
	}
	return &v1.ListRes{Current: current, Size: size, Total: total, Records: records}, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (*v1.SaveRes, error) {
	result, err := g.DB().Model("vendor").Ctx(ctx).Data(g.Map{
		"company_id": currentCompanyId(ctx), "name": req.Name, "contact_name": req.ContactName, "phone": req.Phone, "address": req.Address, "remark": req.Remark,
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
	_, err := g.DB().Model("vendor").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where("id", req.Id).Where("deleted_at IS NULL").Data(g.Map{
		"name": req.Name, "contact_name": req.ContactName, "phone": req.Phone, "address": req.Address, "remark": req.Remark,
	}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.SaveRes{Id: req.Id}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	_, err := g.DB().Model("vendor").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where("id", req.Id).Where("deleted_at IS NULL").Data(g.Map{"deleted_at": gtime.Now()}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (*v1.BatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return &v1.BatchDeleteRes{}, nil
	}
	_, err := g.DB().Model("vendor").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).WhereIn("id", req.Ids).Where("deleted_at IS NULL").Data(g.Map{"deleted_at": gtime.Now()}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.BatchDeleteRes{}, nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}
