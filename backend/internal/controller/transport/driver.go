package transport

import (
	"context"
	"strings"

	v1 "engineServ/api/transport/v1"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

const driverSelectFields = `id,name,phone,license_number,license_type,service_status,rating,completed_trips,
	DATE_FORMAT(joined_at,'%Y-%m-%d') AS joined_at,remark,
	DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s') AS created_at,
	DATE_FORMAT(updated_at,'%Y-%m-%d %H:%i:%s') AS updated_at`

// DriverList 司机分页列表。
func (c *Controller) DriverList(ctx context.Context, req *v1.DriverListReq) (*v1.DriverListRes, error) {
	current, size := normalizePage(req.Current, req.Size)
	companyId := currentCompanyId(ctx)

	where := []string{"company_id=?", "deleted_at IS NULL"}
	args := []any{companyId}

	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		where = append(where, "(name LIKE ? OR phone LIKE ? OR license_number LIKE ?)")
		args = append(args, like, like, like)
	}
	if status := strings.TrimSpace(req.ServiceStatus); status != "" {
		where = append(where, "service_status=?")
		args = append(args, status)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM driver WHERE "+whereSQL, args...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询司机总数失败")
	}

	rows, err := g.DB().GetAll(ctx,
		"SELECT "+driverSelectFields+" FROM driver WHERE "+whereSQL+" ORDER BY id DESC LIMIT ?,?",
		append(args, pageOffset(current, size), size)...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询司机列表失败")
	}

	records := make([]*v1.DriverItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, driverItemFromRow(row))
	}
	return &v1.DriverListRes{Current: current, Size: size, Total: total.Int(), Records: records}, nil
}

// DriverDetail 司机详情。
func (c *Controller) DriverDetail(ctx context.Context, req *v1.DriverDetailReq) (*v1.DriverDetailRes, error) {
	row, err := g.DB().GetOne(ctx,
		"SELECT "+driverSelectFields+" FROM driver WHERE id=? AND company_id=? AND deleted_at IS NULL",
		req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询司机详情失败")
	}
	if row.IsEmpty() {
		return nil, gerror.New("司机不存在或已删除")
	}
	return driverItemFromRow(row), nil
}

// DriverCreate 新增司机。
func (c *Controller) DriverCreate(ctx context.Context, req *v1.DriverCreateReq) (*v1.DriverCreateRes, error) {
	companyId := currentCompanyId(ctx)
	if err := ensureDriverPhoneUnique(ctx, companyId, req.Phone, 0); err != nil {
		return nil, err
	}

	data := g.Map{
		"company_id":      companyId,
		"name":            strings.TrimSpace(req.Name),
		"phone":           strings.TrimSpace(req.Phone),
		"license_number":  strings.TrimSpace(req.LicenseNumber),
		"license_type":    strings.TrimSpace(req.LicenseType),
		"service_status":  req.ServiceStatus,
		"rating":          req.Rating,
		"completed_trips": req.CompletedTrips,
		"joined_at":       optionalDate(req.JoinedAt),
		"remark":          strings.TrimSpace(req.Remark),
	}
	result, err := g.DB().Model("driver").Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "新增司机失败")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "新增司机失败")
	}
	return &v1.DriverCreateRes{Id: uint(id)}, nil
}

// DriverUpdate 更新司机。
func (c *Controller) DriverUpdate(ctx context.Context, req *v1.DriverUpdateReq) (*v1.DriverUpdateRes, error) {
	companyId := currentCompanyId(ctx)
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM driver WHERE id=? AND company_id=? AND deleted_at IS NULL", req.Id, companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "校验司机失败")
	}
	if count.Int() == 0 {
		return nil, gerror.New("司机不存在或已删除")
	}
	if err := ensureDriverPhoneUnique(ctx, companyId, req.Phone, req.Id); err != nil {
		return nil, err
	}

	_, err = g.DB().Model("driver").Ctx(ctx).
		Where("id", req.Id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		Data(g.Map{
			"name":            strings.TrimSpace(req.Name),
			"phone":           strings.TrimSpace(req.Phone),
			"license_number":  strings.TrimSpace(req.LicenseNumber),
			"license_type":    strings.TrimSpace(req.LicenseType),
			"service_status":  req.ServiceStatus,
			"rating":          req.Rating,
			"completed_trips": req.CompletedTrips,
			"joined_at":       optionalDate(req.JoinedAt),
			"remark":          strings.TrimSpace(req.Remark),
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新司机失败")
	}
	return &v1.DriverUpdateRes{}, nil
}

// DriverDelete 删除司机（软删除，并解除车辆绑定）。
func (c *Controller) DriverDelete(ctx context.Context, req *v1.DriverDeleteReq) (*v1.DriverDeleteRes, error) {
	companyId := currentCompanyId(ctx)
	if _, err := g.DB().Model("vehicle").Ctx(ctx).
		Where("driver_id", req.Id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		Data(g.Map{"driver_id": 0}).Update(); err != nil {
		return nil, gerror.Wrap(err, "解除车辆司机绑定失败")
	}
	result, err := g.DB().Model("driver").Ctx(ctx).
		Where("id", req.Id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now()}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "删除司机失败")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.New("司机不存在或已删除")
	}
	return &v1.DriverDeleteRes{}, nil
}

// DriverBatchDelete 批量删除司机。
func (c *Controller) DriverBatchDelete(ctx context.Context, req *v1.DriverBatchDeleteReq) (*v1.DriverBatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return &v1.DriverBatchDeleteRes{}, nil
	}
	companyId := currentCompanyId(ctx)
	if _, err := g.DB().Model("vehicle").Ctx(ctx).
		Where("company_id", companyId).
		WhereIn("driver_id", req.Ids).
		Where("deleted_at IS NULL").
		Data(g.Map{"driver_id": 0}).Update(); err != nil {
		return nil, gerror.Wrap(err, "解除车辆司机绑定失败")
	}
	if _, err := g.DB().Model("driver").Ctx(ctx).
		Where("company_id", companyId).
		WhereIn("id", req.Ids).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now()}).Update(); err != nil {
		return nil, gerror.Wrap(err, "批量删除司机失败")
	}
	return &v1.DriverBatchDeleteRes{}, nil
}

// DriverOptions 司机下拉选项。
func (c *Controller) DriverOptions(ctx context.Context, _ *v1.DriverOptionReq) (*v1.DriverOptionRes, error) {
	rows, err := g.DB().GetAll(ctx,
		`SELECT id,name,phone,service_status FROM driver
		 WHERE company_id=? AND deleted_at IS NULL
		 ORDER BY FIELD(service_status,'available','busy','offline','suspended'), id`,
		currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询司机选项失败")
	}
	items := make([]*v1.DriverOption, 0, len(rows))
	for _, row := range rows {
		items = append(items, &v1.DriverOption{
			Id:            row["id"].Uint(),
			Name:          row["name"].String(),
			Phone:         row["phone"].String(),
			ServiceStatus: row["service_status"].String(),
		})
	}
	return &v1.DriverOptionRes{Items: items}, nil
}

func ensureDriverPhoneUnique(ctx context.Context, companyId uint, phone string, excludeId uint) error {
	count, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM driver WHERE company_id=? AND phone=? AND deleted_at IS NULL AND id<>?",
		companyId, strings.TrimSpace(phone), excludeId)
	if err != nil {
		return gerror.Wrap(err, "校验手机号失败")
	}
	if count.Int() > 0 {
		return gerror.New("该手机号已被其他司机使用")
	}
	return nil
}

func driverItemFromRow(row gdb.Record) *v1.DriverItem {
	return &v1.DriverItem{
		Id:             row["id"].Uint(),
		Name:           row["name"].String(),
		Phone:          row["phone"].String(),
		LicenseNumber:  row["license_number"].String(),
		LicenseType:    row["license_type"].String(),
		ServiceStatus:  row["service_status"].String(),
		Rating:         row["rating"].Float64(),
		CompletedTrips: row["completed_trips"].Uint(),
		JoinedAt:       row["joined_at"].String(),
		Remark:         row["remark"].String(),
		CreatedAt:      row["created_at"].String(),
		UpdatedAt:      row["updated_at"].String(),
	}
}

func optionalDate(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}
