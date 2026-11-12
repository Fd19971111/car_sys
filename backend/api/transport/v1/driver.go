package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 司机列表 =====================

// DriverListReq 司机分页列表请求。
type DriverListReq struct {
	g.Meta        `path:"/transport/driver/list" method:"get" tags:"运力管理-司机" summary:"司机分页列表"`
	Current       int    `json:"current"       d:"1"  dc:"当前页码"`
	Size          int    `json:"size"          d:"10" dc:"每页条数"`
	Keyword       string `json:"keyword"       dc:"姓名 / 手机号 / 驾驶证号"`
	ServiceStatus string `json:"serviceStatus" dc:"接单状态：available busy offline suspended"`
}

// DriverItem 司机列表行。
type DriverItem struct {
	Id             uint    `json:"id"             dc:"司机ID"`
	Name           string  `json:"name"           dc:"司机姓名"`
	Phone          string  `json:"phone"          dc:"手机号"`
	LicenseNumber  string  `json:"licenseNumber"  dc:"驾驶证号"`
	LicenseType    string  `json:"licenseType"    dc:"准驾车型"`
	ServiceStatus  string  `json:"serviceStatus"  dc:"接单状态"`
	Rating         float64 `json:"rating"         dc:"服务评分"`
	CompletedTrips uint    `json:"completedTrips" dc:"完成订单数"`
	JoinedAt       string  `json:"joinedAt"       dc:"加入日期"`
	Remark         string  `json:"remark"         dc:"备注"`
	CreatedAt      string  `json:"createdAt"      dc:"创建时间"`
	UpdatedAt      string  `json:"updatedAt"      dc:"更新时间"`
}

// DriverListRes 司机分页列表响应。
type DriverListRes struct {
	Current int           `json:"current" dc:"当前页码"`
	Size    int           `json:"size"    dc:"每页条数"`
	Total   int           `json:"total"   dc:"总条数"`
	Records []*DriverItem `json:"records" dc:"司机列表"`
}

// ===================== 司机详情 =====================

// DriverDetailReq 司机详情请求。
type DriverDetailReq struct {
	g.Meta `path:"/transport/driver/detail" method:"get" tags:"运力管理-司机" summary:"司机详情"`
	Id     uint `json:"id" v:"required#司机ID不能为空" dc:"司机ID"`
}

// DriverDetailRes 司机详情响应。
type DriverDetailRes = DriverItem

// ===================== 新增 / 编辑 / 删除 =====================

// DriverForm 司机表单字段。
type DriverForm struct {
	Name           string  `json:"name"           v:"required|length:1,50#请填写司机姓名|司机姓名最多50个字符" dc:"司机姓名"`
	Phone          string  `json:"phone"          v:"required|regex:^1[3-9][0-9]{9}$#请填写手机号|手机号格式不正确" dc:"手机号"`
	LicenseNumber  string  `json:"licenseNumber"  v:"max-length:50#驾驶证号最多50个字符"  dc:"驾驶证号"`
	LicenseType    string  `json:"licenseType"    v:"max-length:20#准驾车型最多20个字符"  dc:"准驾车型"`
	ServiceStatus  string  `json:"serviceStatus"  v:"required|in:available,busy,offline,suspended#请选择接单状态|接单状态无效" dc:"接单状态"`
	Rating         float64 `json:"rating"         v:"min:0|max:5#服务评分不能小于0|服务评分不能大于5" dc:"服务评分"`
	CompletedTrips uint    `json:"completedTrips" dc:"完成订单数"`
	JoinedAt       string  `json:"joinedAt"       dc:"加入日期，格式 YYYY-MM-DD"`
	Remark         string  `json:"remark"         v:"max-length:255#备注最多255个字符" dc:"备注"`
}

// DriverCreateReq 新增司机请求。
type DriverCreateReq struct {
	g.Meta `path:"/transport/driver/create" method:"post" tags:"运力管理-司机" summary:"新增司机"`
	DriverForm
}

// DriverCreateRes 新增司机响应。
type DriverCreateRes struct {
	Id uint `json:"id" dc:"司机ID"`
}

// DriverUpdateReq 更新司机请求。
type DriverUpdateReq struct {
	g.Meta `path:"/transport/driver/update" method:"post" tags:"运力管理-司机" summary:"更新司机"`
	Id     uint `json:"id" v:"required#司机ID不能为空" dc:"司机ID"`
	DriverForm
}

// DriverUpdateRes 更新司机响应。
type DriverUpdateRes struct{}

// DriverDeleteReq 删除司机请求。
type DriverDeleteReq struct {
	g.Meta `path:"/transport/driver/delete" method:"post" tags:"运力管理-司机" summary:"删除司机"`
	Id     uint `json:"id" v:"required#司机ID不能为空" dc:"司机ID"`
}

// DriverDeleteRes 删除司机响应。
type DriverDeleteRes struct{}

// DriverBatchDeleteReq 批量删除司机请求。
type DriverBatchDeleteReq struct {
	g.Meta `path:"/transport/driver/batchDelete" method:"post" tags:"运力管理-司机" summary:"批量删除司机"`
	Ids    []uint `json:"ids" v:"required#请选择要删除的司机" dc:"司机ID列表"`
}

// DriverBatchDeleteRes 批量删除司机响应。
type DriverBatchDeleteRes struct{}

// ===================== 司机下拉选项 =====================

// DriverOptionReq 司机下拉选项请求。
type DriverOptionReq struct {
	g.Meta `path:"/transport/driver/options" method:"get" tags:"运力管理-司机" summary:"司机下拉选项"`
}

// DriverOption 司机下拉选项。
type DriverOption struct {
	Id            uint   `json:"id"            dc:"司机ID"`
	Name          string `json:"name"          dc:"司机姓名"`
	Phone         string `json:"phone"         dc:"手机号"`
	ServiceStatus string `json:"serviceStatus" dc:"接单状态"`
}

// DriverOptionRes 司机下拉选项响应。
type DriverOptionRes struct {
	Items []*DriverOption `json:"items" dc:"司机选项"`
}
