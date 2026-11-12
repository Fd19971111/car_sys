package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 车辆列表 =====================

// VehicleListReq 车辆分页列表请求。
type VehicleListReq struct {
	g.Meta  `path:"/transport/vehicle/list" method:"get" tags:"运力管理-车辆" summary:"车辆分页列表"`
	Current int    `json:"current"  d:"1"  dc:"当前页码"`
	Size    int    `json:"size"     d:"10" dc:"每页条数"`
	Keyword string `json:"keyword"  dc:"车牌号 / 品牌 / 司机姓名"`
	Status  string `json:"status"   dc:"车辆状态：active maintenance inactive"`
}

// VehicleItem 车辆列表行。
type VehicleItem struct {
	Id          uint   `json:"id"          dc:"车辆ID"`
	PlateNumber string `json:"plateNumber" dc:"车牌号"`
	Brand       string `json:"brand"       dc:"品牌车型"`
	Category    string `json:"category"    dc:"服务车型"`
	SeatCount   uint   `json:"seatCount"   dc:"可乘人数"`
	Status      string `json:"status"      dc:"车辆状态"`
	DriverId    uint   `json:"driverId"    dc:"当前司机ID"`
	DriverName  string `json:"driverName"  dc:"当前司机姓名"`
	Remark      string `json:"remark"      dc:"备注"`
	CreatedAt   string `json:"createdAt"   dc:"创建时间"`
	UpdatedAt   string `json:"updatedAt"   dc:"更新时间"`
}

// VehicleListRes 车辆分页列表响应。
type VehicleListRes struct {
	Current int            `json:"current" dc:"当前页码"`
	Size    int            `json:"size"    dc:"每页条数"`
	Total   int            `json:"total"   dc:"总条数"`
	Records []*VehicleItem `json:"records" dc:"车辆列表"`
}

// ===================== 车辆详情 =====================

// VehicleDetailReq 车辆详情请求。
type VehicleDetailReq struct {
	g.Meta `path:"/transport/vehicle/detail" method:"get" tags:"运力管理-车辆" summary:"车辆详情"`
	Id     uint `json:"id" v:"required#车辆ID不能为空" dc:"车辆ID"`
}

// VehicleDetailRes 车辆详情响应。
type VehicleDetailRes = VehicleItem

// ===================== 新增 / 编辑 / 删除 =====================

// VehicleForm 车辆表单字段。
type VehicleForm struct {
	PlateNumber string `json:"plateNumber" v:"required|max-length:20#请填写车牌号|车牌号最多20个字符" dc:"车牌号"`
	Brand       string `json:"brand"       v:"max-length:100#品牌车型最多100个字符" dc:"品牌车型"`
	Category    string `json:"category"    v:"required|in:economy,comfort,business,six-seat#请选择服务车型|服务车型无效" dc:"服务车型"`
	SeatCount   uint   `json:"seatCount"   v:"required|min:1|max:8#请填写可乘人数|可乘人数至少1人|可乘人数最多8人" dc:"可乘人数"`
	Status      string `json:"status"      v:"required|in:active,maintenance,inactive#请选择车辆状态|车辆状态无效" dc:"车辆状态"`
	DriverId    uint   `json:"driverId"    dc:"当前司机ID，0 表示未绑定"`
	Remark      string `json:"remark"      v:"max-length:255#备注最多255个字符" dc:"备注"`
}

// VehicleCreateReq 新增车辆请求。
type VehicleCreateReq struct {
	g.Meta `path:"/transport/vehicle/create" method:"post" tags:"运力管理-车辆" summary:"新增车辆"`
	VehicleForm
}

// VehicleCreateRes 新增车辆响应。
type VehicleCreateRes struct {
	Id uint `json:"id" dc:"车辆ID"`
}

// VehicleUpdateReq 更新车辆请求。
type VehicleUpdateReq struct {
	g.Meta `path:"/transport/vehicle/update" method:"post" tags:"运力管理-车辆" summary:"更新车辆"`
	Id     uint `json:"id" v:"required#车辆ID不能为空" dc:"车辆ID"`
	VehicleForm
}

// VehicleUpdateRes 更新车辆响应。
type VehicleUpdateRes struct{}

// VehicleDeleteReq 删除车辆请求。
type VehicleDeleteReq struct {
	g.Meta `path:"/transport/vehicle/delete" method:"post" tags:"运力管理-车辆" summary:"删除车辆"`
	Id     uint `json:"id" v:"required#车辆ID不能为空" dc:"车辆ID"`
}

// VehicleDeleteRes 删除车辆响应。
type VehicleDeleteRes struct{}

// VehicleBatchDeleteReq 批量删除车辆请求。
type VehicleBatchDeleteReq struct {
	g.Meta `path:"/transport/vehicle/batchDelete" method:"post" tags:"运力管理-车辆" summary:"批量删除车辆"`
	Ids    []uint `json:"ids" v:"required#请选择要删除的车辆" dc:"车辆ID列表"`
}

// VehicleBatchDeleteRes 批量删除车辆响应。
type VehicleBatchDeleteRes struct{}
