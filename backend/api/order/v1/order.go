package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 全部订单 =====================

// ListReq 订单分页列表请求。
type ListReq struct {
	g.Meta           `path:"/order/list" method:"get" tags:"订单中心" summary:"订单分页列表"`
	Current          int    `json:"current"          d:"1"  dc:"当前页码"`
	Size             int    `json:"size"             d:"10" dc:"每页条数"`
	Keyword          string `json:"keyword"          dc:"订单号 / 客户 / 手机号 / 上下车地点 / 司机 / 车牌"`
	Status           string `json:"status"           dc:"订单状态"`
	CustomerResponse string `json:"customerResponse" dc:"客户反馈"`
	DateFrom         string `json:"dateFrom"         dc:"预约开始日期 YYYY-MM-DD"`
	DateTo           string `json:"dateTo"           dc:"预约结束日期 YYYY-MM-DD"`
}

// OrderItem 订单列表行。
type OrderItem struct {
	Id                   int64   `json:"id"                 dc:"订单ID"`
	OrderNumber          string  `json:"orderNumber"        dc:"订单编号"`
	CustomerId           int64   `json:"customerId"         dc:"客户ID"`
	CustomerName         string  `json:"customerName"       dc:"客户姓名"`
	CustomerPhone        string  `json:"customerPhone"      dc:"客户手机号"`
	Status               string  `json:"status"             dc:"订单状态"`
	PassengerCount       uint    `json:"passengerCount"     dc:"乘车人数"`
	VehicleCategory      string  `json:"vehicleCategory"    dc:"车型"`
	PickupAddress        string  `json:"pickupAddress"      dc:"上车地点"`
	DestinationAddress   string  `json:"destinationAddress" dc:"目的地点"`
	EstimatedDistanceKm  float64 `json:"estimatedDistanceKm" dc:"预计里程（公里）"`
	EstimatedDurationMin uint    `json:"estimatedDurationMin" dc:"预计时长（分钟）"`
	QuotedFare           float64 `json:"quotedFare"         dc:"客户报价"`
	FinalFare            float64 `json:"finalFare"          dc:"最终费用"`
	DriverFare           float64 `json:"driverFare"         dc:"本单司机应得车费"`
	CustomerResponse     string  `json:"customerResponse"   dc:"客户反馈"`
	ScheduledAt          string  `json:"scheduledAt"        dc:"预约用车时间"`
	DriverId             int64   `json:"driverId"           dc:"司机ID"`
	DriverName           string  `json:"driverName"         dc:"司机姓名"`
	DriverPhone          string  `json:"driverPhone"        dc:"司机手机号"`
	VehicleId            int64   `json:"vehicleId"          dc:"车辆ID"`
	VehiclePlate         string  `json:"vehiclePlate"       dc:"车牌号"`
	CompletedAt          string  `json:"completedAt"        dc:"行程结束时间"`
	InternalNotes        string  `json:"internalNotes"      dc:"内部备注"`
	CreatedAt            string  `json:"createdAt"          dc:"创建时间"`
	UpdatedAt            string  `json:"updatedAt"          dc:"更新时间"`
}

// ListRes 订单分页列表响应。
type ListRes struct {
	Current int          `json:"current" dc:"当前页码"`
	Size    int          `json:"size"    dc:"每页条数"`
	Total   int          `json:"total"   dc:"总条数"`
	Records []*OrderItem `json:"records" dc:"订单列表"`
}

// ===================== 订单详情 =====================

// DetailReq 订单详情请求。
type DetailReq struct {
	g.Meta `path:"/order/detail" method:"get" tags:"订单中心" summary:"订单详情"`
	Id     int64 `json:"id" v:"required#订单ID不能为空" dc:"订单ID"`
}

// DetailRes 订单详情响应。
type DetailRes = OrderItem

// ===================== 新增 / 编辑 / 删除 =====================

// OrderForm 订单表单字段。
type OrderForm struct {
	CustomerId           int64   `json:"customerId"           v:"required#请选择客户" dc:"客户ID"`
	Status               string  `json:"status"               v:"required|in:quoting,awaiting-confirmation,dispatching,driver-en-route,in-progress,completed,cancelled#请选择订单状态|订单状态无效" dc:"订单状态"`
	PassengerCount       uint    `json:"passengerCount"       v:"required|min:1|max:8#请填写乘车人数|乘车人数至少1人|乘车人数最多8人" dc:"乘车人数"`
	VehicleCategory      string  `json:"vehicleCategory"      v:"required|in:economy,comfort,business,six-seat#请选择车型|车型无效" dc:"车型"`
	PickupAddress        string  `json:"pickupAddress"        v:"required|max-length:255#请填写上车地点|上车地点最多255个字符" dc:"上车地点"`
	DestinationAddress   string  `json:"destinationAddress"   v:"required|max-length:255#请填写目的地点|目的地点最多255个字符" dc:"目的地点"`
	EstimatedDistanceKm  float64 `json:"estimatedDistanceKm"  v:"min:0#预计里程不能小于0" dc:"预计里程（公里）"`
	EstimatedDurationMin uint    `json:"estimatedDurationMin" dc:"预计时长（分钟）"`
	QuotedFare           float64 `json:"quotedFare"           v:"min:0#客户报价不能小于0" dc:"客户报价"`
	FinalFare            float64 `json:"finalFare"            v:"min:0#最终费用不能小于0" dc:"最终费用"`
	DriverFare           float64 `json:"driverFare"           v:"min:0#司机车费不能小于0" dc:"本单司机应得车费"`
	CustomerResponse     string  `json:"customerResponse"     v:"required|in:pending,accepted,rejected#请选择客户反馈|客户反馈无效" dc:"客户反馈"`
	ScheduledAt          string  `json:"scheduledAt"          v:"required#请选择预约用车时间" dc:"预约用车时间 YYYY-MM-DD HH:mm:ss"`
	DriverId             int64   `json:"driverId"             dc:"司机ID，0 表示未指派"`
	VehicleId            int64   `json:"vehicleId"            dc:"车辆ID，0 表示未指派"`
	CompletedAt          string  `json:"completedAt"          dc:"行程结束时间 YYYY-MM-DD HH:mm:ss"`
	InternalNotes        string  `json:"internalNotes"        v:"max-length:500#内部备注最多500个字符" dc:"内部备注"`
}

// CreateReq 新增订单请求。
type CreateReq struct {
	g.Meta `path:"/order/create" method:"post" tags:"订单中心" summary:"新增订单"`
	OrderForm
}

// CreateRes 新增订单响应。
type CreateRes struct {
	Id          int64  `json:"id"          dc:"订单ID"`
	OrderNumber string `json:"orderNumber" dc:"订单编号"`
}

// UpdateReq 更新订单请求。
type UpdateReq struct {
	g.Meta `path:"/order/update" method:"post" tags:"订单中心" summary:"更新订单"`
	Id     int64 `json:"id" v:"required#订单ID不能为空" dc:"订单ID"`
	OrderForm
}

// UpdateRes 更新订单响应。
type UpdateRes struct{}

// DeleteReq 删除订单请求。
type DeleteReq struct {
	g.Meta `path:"/order/delete" method:"post" tags:"订单中心" summary:"删除订单"`
	Id     int64 `json:"id" v:"required#订单ID不能为空" dc:"订单ID"`
}

// DeleteRes 删除订单响应。
type DeleteRes struct{}

// BatchDeleteReq 批量删除订单请求。
type BatchDeleteReq struct {
	g.Meta `path:"/order/batchDelete" method:"post" tags:"订单中心" summary:"批量删除订单"`
	Ids    []int64 `json:"ids" v:"required#请选择要删除的订单" dc:"订单ID列表"`
}

// BatchDeleteRes 批量删除订单响应。
type BatchDeleteRes struct{}

// ===================== 统计 / 行程监控 =====================

// StatusCount 状态数量。
type StatusCount struct {
	Status string `json:"status" dc:"状态"`
	Count  int    `json:"count"  dc:"数量"`
}

// StatsReq 订单统计请求。
type StatsReq struct {
	g.Meta `path:"/order/stats" method:"get" tags:"订单中心" summary:"订单统计"`
}

// StatsRes 订单统计响应。
type StatsRes struct {
	Total        int            `json:"total"         dc:"订单总数"`
	Today        int            `json:"today"         dc:"今日预约订单"`
	Active       int            `json:"active"        dc:"进行中订单"`
	Completed    int            `json:"completed"     dc:"已完成订单"`
	Cancelled    int            `json:"cancelled"     dc:"已取消订单"`
	Revenue      float64        `json:"revenue"       dc:"已完成订单收入"`
	StatusCounts []*StatusCount `json:"statusCounts"  dc:"各状态数量"`
}

// MonitorReq 行程监控请求。
type MonitorReq struct {
	g.Meta `path:"/order/monitor" method:"get" tags:"订单中心" summary:"行程监控"`
}

// MonitorRes 行程监控响应。
type MonitorRes struct {
	ActiveCount   int          `json:"activeCount"   dc:"进行中订单数"`
	Dispatching   int          `json:"dispatching"   dc:"待派单"`
	DriverEnRoute int          `json:"driverEnRoute" dc:"司机前往中"`
	InProgress    int          `json:"inProgress"    dc:"行程中"`
	Records       []*OrderItem `json:"records"       dc:"进行中订单"`
	UpdatedAt     string       `json:"updatedAt"     dc:"刷新时间"`
}
