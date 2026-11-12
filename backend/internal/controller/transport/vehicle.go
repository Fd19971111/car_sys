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

const vehicleSelectFields = `v.id,v.plate_number,v.brand,v.category,v.seat_count,v.status,v.driver_id,
	COALESCE(d.name,'') AS driver_name,v.remark,
	DATE_FORMAT(v.created_at,'%Y-%m-%d %H:%i:%s') AS created_at,
	DATE_FORMAT(v.updated_at,'%Y-%m-%d %H:%i:%s') AS updated_at`

const vehicleFromJoin = `FROM vehicle v
	LEFT JOIN driver d ON d.id=v.driver_id AND d.deleted_at IS NULL`

// VehicleList 车辆分页列表。
func (c *Controller) VehicleList(ctx context.Context, req *v1.VehicleListReq) (*v1.VehicleListRes, error) {
	current, size := normalizePage(req.Current, req.Size)
	companyId := currentCompanyId(ctx)

	where := []string{"v.company_id=?", "v.deleted_at IS NULL"}
	args := []any{companyId}

	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		where = append(where, "(v.plate_number LIKE ? OR v.brand LIKE ? OR d.name LIKE ?)")
		args = append(args, like, like, like)
	}
	if status := strings.TrimSpace(req.Status); status != "" {
		where = append(where, "v.status=?")
		args = append(args, status)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) "+vehicleFromJoin+" WHERE "+whereSQL, args...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询车辆总数失败")
	}

	rows, err := g.DB().GetAll(ctx,
		"SELECT "+vehicleSelectFields+" "+vehicleFromJoin+" WHERE "+whereSQL+" ORDER BY v.id DESC LIMIT ?,?",
		append(args, pageOffset(current, size), size)...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询车辆列表失败")
	}

	records := make([]*v1.VehicleItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, vehicleItemFromRow(row))
	}
	return &v1.VehicleListRes{Current: current, Size: size, Total: total.Int(), Records: records}, nil
}

// VehicleDetail 车辆详情。
func (c *Controller) VehicleDetail(ctx context.Context, req *v1.VehicleDetailReq) (*v1.VehicleDetailRes, error) {
	row, err := g.DB().GetOne(ctx,
		"SELECT "+vehicleSelectFields+" "+vehicleFromJoin+" WHERE v.id=? AND v.company_id=? AND v.deleted_at IS NULL",
		req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询车辆详情失败")
	}
	if row.IsEmpty() {
		return nil, gerror.New("车辆不存在或已删除")
	}
	return vehicleItemFromRow(row), nil
}

// VehicleCreate 新增车辆。
func (c *Controller) VehicleCreate(ctx context.Context, req *v1.VehicleCreateReq) (*v1.VehicleCreateRes, error) {
	companyId := currentCompanyId(ctx)
	plateNumber := strings.TrimSpace(req.PlateNumber)
	if err := ensureVehiclePlateUnique(ctx, companyId, plateNumber, 0); err != nil {
		return nil, err
	}
	if err := ensureDriverAssignable(ctx, companyId, req.DriverId); err != nil {
		return nil, err
	}

	result, err := g.DB().Model("vehicle").Ctx(ctx).Data(g.Map{
		"company_id":   companyId,
		"plate_number": plateNumber,
		"brand":        strings.TrimSpace(req.Brand),
		"category":     req.Category,
		"seat_count":   req.SeatCount,
		"status":       req.Status,
		"driver_id":    req.DriverId,
		"remark":       strings.TrimSpace(req.Remark),
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "新增车辆失败")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "新增车辆失败")
	}
	return &v1.VehicleCreateRes{Id: uint(id)}, nil
}

// VehicleUpdate 更新车辆。
func (c *Controller) VehicleUpdate(ctx context.Context, req *v1.VehicleUpdateReq) (*v1.VehicleUpdateRes, error) {
	companyId := currentCompanyId(ctx)
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM vehicle WHERE id=? AND company_id=? AND deleted_at IS NULL", req.Id, companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "校验车辆失败")
	}
	if count.Int() == 0 {
		return nil, gerror.New("车辆不存在或已删除")
	}

	plateNumber := strings.TrimSpace(req.PlateNumber)
	if err := ensureVehiclePlateUnique(ctx, companyId, plateNumber, req.Id); err != nil {
		return nil, err
	}
	if err := ensureDriverAssignable(ctx, companyId, req.DriverId); err != nil {
		return nil, err
	}

	_, err = g.DB().Model("vehicle").Ctx(ctx).
		Where("id", req.Id).
		Where("company_id", companyId).
		Where("deleted_at IS NULL").
		Data(g.Map{
			"plate_number": plateNumber,
			"brand":        strings.TrimSpace(req.Brand),
			"category":     req.Category,
			"seat_count":   req.SeatCount,
			"status":       req.Status,
			"driver_id":    req.DriverId,
			"remark":       strings.TrimSpace(req.Remark),
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新车辆失败")
	}
	return &v1.VehicleUpdateRes{}, nil
}

// VehicleDelete 删除车辆（软删除）。
func (c *Controller) VehicleDelete(ctx context.Context, req *v1.VehicleDeleteReq) (*v1.VehicleDeleteRes, error) {
	result, err := g.DB().Model("vehicle").Ctx(ctx).
		Where("id", req.Id).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now()}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "删除车辆失败")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.New("车辆不存在或已删除")
	}
	return &v1.VehicleDeleteRes{}, nil
}

// VehicleBatchDelete 批量删除车辆。
func (c *Controller) VehicleBatchDelete(ctx context.Context, req *v1.VehicleBatchDeleteReq) (*v1.VehicleBatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return &v1.VehicleBatchDeleteRes{}, nil
	}
	if _, err := g.DB().Model("vehicle").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn("id", req.Ids).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now()}).Update(); err != nil {
		return nil, gerror.Wrap(err, "批量删除车辆失败")
	}
	return &v1.VehicleBatchDeleteRes{}, nil
}

func ensureVehiclePlateUnique(ctx context.Context, companyId uint, plateNumber string, excludeId uint) error {
	count, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM vehicle WHERE company_id=? AND plate_number=? AND deleted_at IS NULL AND id<>?",
		companyId, plateNumber, excludeId)
	if err != nil {
		return gerror.Wrap(err, "校验车牌号失败")
	}
	if count.Int() > 0 {
		return gerror.New("该车牌号已存在")
	}
	return nil
}

func ensureDriverAssignable(ctx context.Context, companyId, driverId uint) error {
	if driverId == 0 {
		return nil
	}
	count, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM driver WHERE id=? AND company_id=? AND deleted_at IS NULL",
		driverId, companyId)
	if err != nil {
		return gerror.Wrap(err, "校验司机失败")
	}
	if count.Int() == 0 {
		return gerror.New("所选司机不存在或已删除")
	}
	return nil
}

func vehicleItemFromRow(row gdb.Record) *v1.VehicleItem {
	return &v1.VehicleItem{
		Id:          row["id"].Uint(),
		PlateNumber: row["plate_number"].String(),
		Brand:       row["brand"].String(),
		Category:    row["category"].String(),
		SeatCount:   row["seat_count"].Uint(),
		Status:      row["status"].String(),
		DriverId:    row["driver_id"].Uint(),
		DriverName:  row["driver_name"].String(),
		Remark:      row["remark"].String(),
		CreatedAt:   row["created_at"].String(),
		UpdatedAt:   row["updated_at"].String(),
	}
}
