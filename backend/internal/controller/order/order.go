package order

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	orderApi "engineServ/api/order"
	v1 "engineServ/api/order/v1"
	"engineServ/internal/consts"
	"engineServ/internal/logic/fin/booking"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// Controller 订单中心控制器。
type Controller struct{}

// NewV1 创建订单中心控制器。
func NewV1() orderApi.IOrderV1 { return &Controller{} }

const orderSelectFields = `o.id,o.order_number,o.customer_id,
	COALESCE(NULLIF(o.customer_name,''),c.name,'') AS customer_name,
	COALESCE(NULLIF(o.customer_phone,''),c.phone,'') AS customer_phone,
	o.status,o.passenger_count,o.vehicle_category,o.pickup_address,o.destination_address,
	o.estimated_distance_km,o.estimated_duration_min,o.quoted_fare,o.final_fare,o.driver_fare,o.customer_response,
	DATE_FORMAT(o.scheduled_at,'%Y-%m-%d %H:%i:%s') AS scheduled_at,
	o.driver_id,COALESCE(d.name,'') AS driver_name,COALESCE(d.phone,'') AS driver_phone,
	o.vehicle_id,COALESCE(v.plate_number,'') AS vehicle_plate,
	DATE_FORMAT(o.completed_at,'%Y-%m-%d %H:%i:%s') AS completed_at,
	o.internal_notes,
	DATE_FORMAT(o.created_at,'%Y-%m-%d %H:%i:%s') AS created_at,
	DATE_FORMAT(o.updated_at,'%Y-%m-%d %H:%i:%s') AS updated_at`

const orderFromJoin = `FROM ride_order o
	LEFT JOIN customer c ON c.id=o.customer_id AND c.deleted_at IS NULL
	LEFT JOIN driver d ON d.id=o.driver_id AND d.deleted_at IS NULL
	LEFT JOIN vehicle v ON v.id=o.vehicle_id AND v.deleted_at IS NULL`

// List 订单分页列表。
func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	current, size := normalizePage(req.Current, req.Size)
	companyId := currentCompanyId(ctx)

	where := []string{"o.company_id=?", "o.deleted_at IS NULL"}
	args := []any{companyId}

	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		where = append(where, `(o.order_number LIKE ? OR o.customer_name LIKE ? OR o.customer_phone LIKE ?
			OR o.pickup_address LIKE ? OR o.destination_address LIKE ? OR c.name LIKE ? OR d.name LIKE ? OR v.plate_number LIKE ?)`)
		args = append(args, like, like, like, like, like, like, like, like)
	}
	if status := strings.TrimSpace(req.Status); status != "" {
		where = append(where, "o.status=?")
		args = append(args, status)
	}
	if response := strings.TrimSpace(req.CustomerResponse); response != "" {
		where = append(where, "o.customer_response=?")
		args = append(args, response)
	}
	if dateFrom := strings.TrimSpace(req.DateFrom); dateFrom != "" {
		where = append(where, "DATE(o.scheduled_at)>=?")
		args = append(args, dateFrom)
	}
	if dateTo := strings.TrimSpace(req.DateTo); dateTo != "" {
		where = append(where, "DATE(o.scheduled_at)<=?")
		args = append(args, dateTo)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) "+orderFromJoin+" WHERE "+whereSQL, args...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单总数失败")
	}

	rows, err := g.DB().GetAll(ctx,
		"SELECT "+orderSelectFields+" "+orderFromJoin+" WHERE "+whereSQL+
			" ORDER BY o.scheduled_at DESC,o.id DESC LIMIT ?,?",
		append(args, pageOffset(current, size), size)...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单列表失败")
	}

	records := make([]*v1.OrderItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, orderItemFromRow(row))
	}
	return &v1.ListRes{Current: current, Size: size, Total: total.Int(), Records: records}, nil
}

// Detail 订单详情。
func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (*v1.DetailRes, error) {
	row, err := g.DB().GetOne(ctx,
		"SELECT "+orderSelectFields+" "+orderFromJoin+
			" WHERE o.id=? AND o.company_id=? AND o.deleted_at IS NULL",
		req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单详情失败")
	}
	if row.IsEmpty() {
		return nil, gerror.New("订单不存在或已删除")
	}
	return orderItemFromRow(row), nil
}

// Create 新增订单。
func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	companyId := currentCompanyId(ctx)
	customerName, customerPhone, err := customerSnapshot(ctx, companyId, req.CustomerId)
	if err != nil {
		return nil, err
	}
	if err = ensureDriverAndVehicle(ctx, companyId, req.DriverId, req.VehicleId); err != nil {
		return nil, err
	}
	orderNumber, err := generateOrderNumber(ctx, companyId)
	if err != nil {
		return nil, err
	}

	data := g.Map{
		"company_id":             companyId,
		"order_number":           orderNumber,
		"customer_id":            req.CustomerId,
		"customer_name":          customerName,
		"customer_phone":         customerPhone,
		"status":                 req.Status,
		"passenger_count":        req.PassengerCount,
		"vehicle_category":       req.VehicleCategory,
		"pickup_address":         strings.TrimSpace(req.PickupAddress),
		"destination_address":    strings.TrimSpace(req.DestinationAddress),
		"estimated_distance_km":  req.EstimatedDistanceKm,
		"estimated_duration_min": req.EstimatedDurationMin,
		"quoted_fare":            req.QuotedFare,
		"final_fare":             req.FinalFare,
		"driver_fare":            req.DriverFare,
		"customer_response":      req.CustomerResponse,
		"scheduled_at":           parseDateTime(req.ScheduledAt),
		"driver_id":              req.DriverId,
		"vehicle_id":             req.VehicleId,
		"completed_at":           completedAtValue(req.Status, req.CompletedAt),
		"internal_notes":         strings.TrimSpace(req.InternalNotes),
		"created_by":             currentUserId(ctx),
	}
	var id int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, insertErr := tx.Model("ride_order").Ctx(ctx).Data(data).Insert()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "新增订单失败")
		}
		id, insertErr = result.LastInsertId()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "新增订单失败")
		}
		return c.postCompletedOrder(ctx, tx, id, companyId, currentUserId(ctx), req.Status, req.DriverId, orderNumber, req.FinalFare, req.DriverFare, req.CompletedAt)
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id, OrderNumber: orderNumber}, nil
}

// Update 更新订单。
func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (*v1.UpdateRes, error) {
	companyId := currentCompanyId(ctx)
	customerName, customerPhone, err := customerSnapshot(ctx, companyId, req.CustomerId)
	if err != nil {
		return nil, err
	}
	if err = ensureDriverAndVehicle(ctx, companyId, req.DriverId, req.VehicleId); err != nil {
		return nil, err
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		old, lockErr := tx.GetOne("SELECT status,final_fare,driver_fare,driver_id FROM ride_order WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", req.Id, companyId)
		if lockErr != nil {
			return gerror.Wrap(lockErr, "校验订单失败")
		}
		if old.IsEmpty() {
			return gerror.New("订单不存在或已删除")
		}
		earningCount, lockErr := tx.GetValue("SELECT COUNT(*) FROM ride_driver_earning WHERE company_id=? AND order_id=? AND status='POSTED'", companyId, req.Id)
		if lockErr != nil {
			return gerror.Wrap(lockErr, "校验司机应得记录失败")
		}
		if earningCount.Int() > 0 && (req.Status != "completed" || old["final_fare"].Float64() != req.FinalFare || old["driver_fare"].Float64() != req.DriverFare || old["driver_id"].Int64() != req.DriverId) {
			return gerror.New("该订单已生成司机应得记录，不能修改完单状态、车费或司机；请先处理财务冲销")
		}
		_, lockErr = tx.Model("ride_order").Ctx(ctx).
			Where("id", req.Id).
			Where("company_id", companyId).
			Where("deleted_at IS NULL").
			Data(g.Map{
				"customer_id":            req.CustomerId,
				"customer_name":          customerName,
				"customer_phone":         customerPhone,
				"status":                 req.Status,
				"passenger_count":        req.PassengerCount,
				"vehicle_category":       req.VehicleCategory,
				"pickup_address":         strings.TrimSpace(req.PickupAddress),
				"destination_address":    strings.TrimSpace(req.DestinationAddress),
				"estimated_distance_km":  req.EstimatedDistanceKm,
				"estimated_duration_min": req.EstimatedDurationMin,
				"quoted_fare":            req.QuotedFare,
				"final_fare":             req.FinalFare,
				"driver_fare":            req.DriverFare,
				"customer_response":      req.CustomerResponse,
				"scheduled_at":           parseDateTime(req.ScheduledAt),
				"driver_id":              req.DriverId,
				"vehicle_id":             req.VehicleId,
				"completed_at":           completedAtValue(req.Status, req.CompletedAt),
				"internal_notes":         strings.TrimSpace(req.InternalNotes),
			}).Update()
		if lockErr != nil {
			return gerror.Wrap(lockErr, "更新订单失败")
		}
		if req.Status == "completed" && earningCount.Int() == 0 {
			row, rowErr := tx.GetOne("SELECT order_number FROM ride_order WHERE id=? AND company_id=?", req.Id, companyId)
			if rowErr != nil {
				return gerror.Wrap(rowErr, "读取订单编号失败")
			}
			return c.postCompletedOrder(ctx, tx, req.Id, companyId, currentUserId(ctx), req.Status, req.DriverId, row["order_number"].String(), req.FinalFare, req.DriverFare, req.CompletedAt)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

func (c *Controller) postCompletedOrder(ctx context.Context, tx gdb.TX, orderId int64, companyId uint, userId uint, status string, driverId int64, orderNumber string, fare float64, driverFare float64, completedAt string) error {
	if status != "completed" {
		return nil
	}
	if driverId <= 0 {
		return gerror.New("完单前请先指派司机")
	}
	if fare <= 0 || driverFare < 0 || driverFare > fare {
		return gerror.New("完单车费必须大于0，司机应得不能超过完单车费")
	}
	date := time.Now().Format("2006-01-02")
	if parsed := completedAtValue(status, completedAt); parsed != nil {
		date = gconv.String(parsed)
		if len(date) > 10 {
			date = date[:10]
		}
	}
	result, err := tx.Exec(`INSERT INTO ride_driver_earning
		(company_id,order_id,driver_id,fare_amount,driver_amount,earning_date,status)
		VALUES(?,?,?,?,?,?,'POSTED')`, companyId, orderId, driverId, fare, driverFare, date)
	if err != nil {
		return gerror.Wrap(err, "保存司机应得车费失败")
	}
	earningId, err := result.LastInsertId()
	if err != nil {
		return gerror.Wrap(err, "读取司机应得记录失败")
	}
	initialized, err := booking.New().IsInitialized(ctx, tx)
	if err != nil {
		return err
	}
	if !initialized {
		return nil
	}
	voucherResult, err := booking.New().PostRideCompletion(ctx, tx, orderId, orderNumber, date, fare, driverFare, userId)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE ride_driver_earning SET voucher_id=? WHERE id=? AND company_id=?", voucherResult.VoucherId, earningId, companyId)
	if err != nil {
		return gerror.Wrap(err, "保存司机应得凭证关联失败")
	}
	return nil
}

// Delete 删除订单（软删除）。
func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	protected, err := g.DB().GetValue(ctx, `SELECT
		(SELECT COUNT(*) FROM ride_driver_earning WHERE company_id=? AND order_id=? AND status='POSTED')
		+ (SELECT COUNT(*) FROM payment WHERE company_id=? AND ride_order_id=? AND deleted_at IS NULL)`,
		currentCompanyId(ctx), req.Id, currentCompanyId(ctx), req.Id)
	if err != nil {
		return nil, gerror.Wrap(err, "检查订单财务记录失败")
	}
	if protected.Int() > 0 {
		return nil, gerror.New("订单已有司机应得或收款记录，不能删除")
	}
	result, err := g.DB().Model("ride_order").Ctx(ctx).
		Where("id", req.Id).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now()}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "删除订单失败")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.New("订单不存在或已删除")
	}
	return &v1.DeleteRes{}, nil
}

// BatchDelete 批量删除订单。
func (c *Controller) BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (*v1.BatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return &v1.BatchDeleteRes{}, nil
	}
	protected, err := g.DB().Model("ride_driver_earning").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("status", "POSTED").
		WhereIn("order_id", req.Ids).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "检查订单财务记录失败")
	}
	if protected > 0 {
		return nil, gerror.New("选中订单中包含已产生司机应得或收款记录的订单，不能批量删除")
	}
	paid, err := g.DB().Model("payment").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		WhereIn("ride_order_id", req.Ids).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "检查订单收款记录失败")
	}
	if paid > 0 {
		return nil, gerror.New("选中订单中包含已收款订单，不能批量删除")
	}
	if _, err := g.DB().Model("ride_order").Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn("id", req.Ids).
		Where("deleted_at IS NULL").
		Data(g.Map{"deleted_at": gtime.Now()}).Update(); err != nil {
		return nil, gerror.Wrap(err, "批量删除订单失败")
	}
	return &v1.BatchDeleteRes{}, nil
}

// Stats 订单统计。
func (c *Controller) Stats(ctx context.Context, _ *v1.StatsReq) (*v1.StatsRes, error) {
	companyId := currentCompanyId(ctx)
	whereSQL := "company_id=? AND deleted_at IS NULL"

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM ride_order WHERE "+whereSQL, companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "统计订单总数失败")
	}
	today, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM ride_order WHERE "+whereSQL+" AND DATE(scheduled_at)=CURDATE()", companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "统计今日订单失败")
	}
	active, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM ride_order WHERE "+whereSQL+
			" AND status IN('dispatching','driver-en-route','in-progress')", companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "统计进行中订单失败")
	}
	completed, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM ride_order WHERE "+whereSQL+" AND status='completed'", companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "统计已完成订单失败")
	}
	cancelled, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM ride_order WHERE "+whereSQL+" AND status='cancelled'", companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "统计已取消订单失败")
	}
	revenue, err := g.DB().GetValue(ctx,
		"SELECT COALESCE(SUM(CASE WHEN status='completed' THEN final_fare ELSE 0 END),0) FROM ride_order WHERE "+whereSQL, companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "统计订单收入失败")
	}

	rows, err := g.DB().GetAll(ctx,
		"SELECT status,COUNT(*) AS count FROM ride_order WHERE "+whereSQL+" GROUP BY status", companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "统计订单状态失败")
	}
	statusCounts := make([]*v1.StatusCount, 0, len(rows))
	for _, row := range rows {
		statusCounts = append(statusCounts, &v1.StatusCount{Status: row["status"].String(), Count: row["count"].Int()})
	}

	return &v1.StatsRes{
		Total: total.Int(), Today: today.Int(), Active: active.Int(), Completed: completed.Int(),
		Cancelled: cancelled.Int(), Revenue: revenue.Float64(), StatusCounts: statusCounts,
	}, nil
}

// Monitor 行程监控。
func (c *Controller) Monitor(ctx context.Context, _ *v1.MonitorReq) (*v1.MonitorRes, error) {
	companyId := currentCompanyId(ctx)
	rows, err := g.DB().GetAll(ctx,
		"SELECT "+orderSelectFields+" "+orderFromJoin+
			" WHERE o.company_id=? AND o.deleted_at IS NULL AND o.status IN('dispatching','driver-en-route','in-progress')"+
			" ORDER BY FIELD(o.status,'in-progress','driver-en-route','dispatching'),o.scheduled_at ASC,o.id DESC LIMIT 50",
		companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "查询行程监控失败")
	}

	records := make([]*v1.OrderItem, 0, len(rows))
	dispatching, enRoute, inProgress := 0, 0, 0
	for _, row := range rows {
		item := orderItemFromRow(row)
		records = append(records, item)
		switch item.Status {
		case "dispatching":
			dispatching++
		case "driver-en-route":
			enRoute++
		case "in-progress":
			inProgress++
		}
	}

	return &v1.MonitorRes{
		ActiveCount: len(records), Dispatching: dispatching, DriverEnRoute: enRoute, InProgress: inProgress,
		Records: records, UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func currentUserId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxUserIdKey))
}

func normalizePage(current, size int) (int, int) {
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 200 {
		size = 10
	}
	return current, size
}

func pageOffset(current, size int) int {
	return (current - 1) * size
}

func orderItemFromRow(row gdb.Record) *v1.OrderItem {
	return &v1.OrderItem{
		Id:                   row["id"].Int64(),
		OrderNumber:          row["order_number"].String(),
		CustomerId:           row["customer_id"].Int64(),
		CustomerName:         row["customer_name"].String(),
		CustomerPhone:        row["customer_phone"].String(),
		Status:               row["status"].String(),
		PassengerCount:       row["passenger_count"].Uint(),
		VehicleCategory:      row["vehicle_category"].String(),
		PickupAddress:        row["pickup_address"].String(),
		DestinationAddress:   row["destination_address"].String(),
		EstimatedDistanceKm:  row["estimated_distance_km"].Float64(),
		EstimatedDurationMin: row["estimated_duration_min"].Uint(),
		QuotedFare:           row["quoted_fare"].Float64(),
		FinalFare:            row["final_fare"].Float64(),
		DriverFare:           row["driver_fare"].Float64(),
		CustomerResponse:     row["customer_response"].String(),
		ScheduledAt:          row["scheduled_at"].String(),
		DriverId:             row["driver_id"].Int64(),
		DriverName:           row["driver_name"].String(),
		DriverPhone:          row["driver_phone"].String(),
		VehicleId:            row["vehicle_id"].Int64(),
		VehiclePlate:         row["vehicle_plate"].String(),
		CompletedAt:          row["completed_at"].String(),
		InternalNotes:        row["internal_notes"].String(),
		CreatedAt:            row["created_at"].String(),
		UpdatedAt:            row["updated_at"].String(),
	}
}

func customerSnapshot(ctx context.Context, companyId uint, customerId int64) (string, string, error) {
	row, err := g.DB().GetOne(ctx,
		"SELECT name,phone FROM customer WHERE id=? AND company_id=? AND deleted_at IS NULL", customerId, companyId)
	if err != nil {
		return "", "", gerror.Wrap(err, "校验客户失败")
	}
	if row.IsEmpty() {
		return "", "", gerror.New("所选客户不存在或已删除")
	}
	return row["name"].String(), row["phone"].String(), nil
}

func ensureDriverAndVehicle(ctx context.Context, companyId uint, driverId, vehicleId int64) error {
	if driverId > 0 {
		count, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM driver WHERE id=? AND company_id=? AND deleted_at IS NULL", driverId, companyId)
		if err != nil {
			return gerror.Wrap(err, "校验司机失败")
		}
		if count.Int() == 0 {
			return gerror.New("所选司机不存在或已删除")
		}
	}
	if vehicleId > 0 {
		count, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM vehicle WHERE id=? AND company_id=? AND deleted_at IS NULL", vehicleId, companyId)
		if err != nil {
			return gerror.Wrap(err, "校验车辆失败")
		}
		if count.Int() == 0 {
			return gerror.New("所选车辆不存在或已删除")
		}
	}
	return nil
}

func generateOrderNumber(ctx context.Context, companyId uint) (string, error) {
	prefix := "RO" + time.Now().Format("20060102")
	pattern := prefix + "%"
	value, err := g.DB().GetValue(ctx,
		"SELECT order_number FROM ride_order WHERE company_id=? AND order_number LIKE ? ORDER BY order_number DESC LIMIT 1",
		companyId, pattern)
	if err != nil {
		return "", gerror.Wrap(err, "生成订单编号失败")
	}
	seq := 1
	if last := value.String(); last != "" && len(last) > len(prefix) {
		if n, parseErr := strconv.Atoi(last[len(prefix):]); parseErr == nil {
			seq = n + 1
		}
	}
	for i := 0; i < 100; i++ {
		candidate := fmt.Sprintf("%s%04d", prefix, seq)
		count, countErr := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM ride_order WHERE company_id=? AND order_number=?", companyId, candidate)
		if countErr != nil {
			return "", gerror.Wrap(countErr, "生成订单编号失败")
		}
		if count.Int() == 0 {
			return candidate, nil
		}
		seq++
	}
	return "", gerror.New("生成订单编号失败，请稍后重试")
}

func parseDateTime(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02 15:04", "2006-01-02T15:04"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed.Format("2006-01-02 15:04:05")
		}
	}
	return value
}

func completedAtValue(status, completedAt string) any {
	if strings.TrimSpace(completedAt) != "" {
		return parseDateTime(completedAt)
	}
	if status == "completed" {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	return nil
}
