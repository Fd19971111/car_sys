package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 客户列表 =====================

// ListReq 客户分页列表请求
type ListReq struct {
	g.Meta `path:"/customer/list" method:"get" sm:"客户分页列表" tags:"Customer"`

	Current int    `json:"current" d:"1"  dc:"当前页码，默认 1"`
	Size    int    `json:"size"    d:"10" dc:"每页条数，默认 10"`
	Keyword string `json:"keyword"          dc:"关键词：匹配客户姓名 / 手机号 / 房屋地址"`
	Status  *int   `json:"status"           dc:"客户状态：1=跟进中 2=已成交 3=已流失"`
	Source  string `json:"source"           dc:"客户来源"`
}

// CustomerListItem 客户列表行
type CustomerListItem struct {
	Id           uint   `json:"id"           dc:"客户ID"`
	Name         string `json:"name"         dc:"客户姓名"`
	Phone        string `json:"phone"        dc:"手机号"`
	Address      string `json:"address"      dc:"房屋地址"`
	Source       string `json:"source"       dc:"客户来源"`
	Status       int    `json:"status"       dc:"客户状态：1=跟进中 2=已成交 3=已流失"`
	Remark       string `json:"remark"       dc:"备注"`
	ProjectCount int    `json:"projectCount" dc:"关联装修项目数量"`
	CreatedAt    string `json:"createdAt"    dc:"创建时间"`
	UpdatedAt    string `json:"updatedAt"    dc:"更新时间"`
}

// ListRes 客户分页列表响应
type ListRes struct {
	Current int                 `json:"current" dc:"当前页码"`
	Size    int                 `json:"size"    dc:"每页条数"`
	Total   int                 `json:"total"   dc:"总条数"`
	Records []*CustomerListItem `json:"records" dc:"客户列表"`
}

// ===================== 客户概览统计 =====================

// StatsReq 客户概览统计请求
type StatsReq struct {
	g.Meta `path:"/customer/stats" method:"get" sm:"客户概览统计" tags:"Customer"`
}

// StatsRes 客户概览统计响应
type StatsRes struct {
	Total      int     `json:"total"      dc:"客户总数"`
	Following  int     `json:"following"  dc:"跟进中数量"`
	Deal       int     `json:"deal"       dc:"已成交数量"`
	Lost       int     `json:"lost"       dc:"已流失数量"`
	DealAmount float64 `json:"dealAmount" dc:"已成交客户的合同总金额"`
}

// ===================== 客户详情 =====================

// DetailReq 客户详情请求
type DetailReq struct {
	g.Meta `path:"/customer/detail" method:"get" sm:"客户详情" tags:"Customer"`

	Id uint `json:"id" v:"required#客户ID不能为空" dc:"客户ID"`
}

// CustomerProject 客户名下的装修项目（含负责人与最近一条进度）
type CustomerProject struct {
	Id               uint    `json:"id"               dc:"项目ID"`
	Name             string  `json:"name"             dc:"项目名称"`
	Status           int     `json:"status"           dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	Amount           float64 `json:"amount"          dc:"合同金额"`
	ManagerName      string  `json:"managerName"      dc:"项目负责人姓名"`
	StartDate        string  `json:"startDate"        dc:"开工日期"`
	EndDate          string  `json:"endDate"          dc:"预计完工日期"`
	LatestProgress   string  `json:"latestProgress"   dc:"最近一条进度记录"`
	LatestProgressAt string  `json:"latestProgressAt" dc:"最近进度记录时间"`
}

// DetailRes 客户详情响应
type DetailRes struct {
	Id        uint               `json:"id"        dc:"客户ID"`
	Name      string             `json:"name"      dc:"客户姓名"`
	Phone     string             `json:"phone"     dc:"手机号"`
	Address   string             `json:"address"   dc:"房屋地址"`
	Source    string             `json:"source"    dc:"客户来源"`
	Status    int                `json:"status"    dc:"客户状态：1=跟进中 2=已成交 3=已流失"`
	Remark    string             `json:"remark"    dc:"备注"`
	CreatedAt string             `json:"createdAt" dc:"创建时间"`
	UpdatedAt string             `json:"updatedAt" dc:"更新时间"`
	Projects  []*CustomerProject `json:"projects"  dc:"关联装修项目"`
}

// ===================== 新增 / 编辑 / 删除 =====================

// CustomerFormFields 客户表单公共字段（customer 表结构固定，不新增业务列）
type CustomerFormFields struct {
	Name string `json:"name"    v:"required|length:2,20#请输入客户姓名|客户姓名长度为 2-20 个字符" dc:"客户姓名"`
	// 注意：不要写 length:11 —— GoFrame 的 length 规则在只给一个参数时 max 会取 0，
	// 变成「长度必须 >=11 且 <=0」，任何值都会被拒。11 位手机号由下面的 regex 保证。
	Phone   string `json:"phone"   v:"required|regex:^1[3-9][0-9]{9}$#请输入手机号|手机号格式不正确" dc:"手机号"`
	Address string `json:"address" v:"max-length:60#房屋地址最多 60 个字符" dc:"房屋地址"`
	Source  string `json:"source"  v:"max-length:20#客户来源最多 20 个字符" dc:"客户来源"`
	Status  int    `json:"status"  v:"required|in:1,2,3#请选择客户状态|客户状态不正确" dc:"客户状态：1=跟进中 2=已成交 3=已流失"`
	Remark  string `json:"remark"  v:"max-length:200#备注最多 200 个字符" dc:"备注"`
}

// CreateReq 新增客户请求
type CreateReq struct {
	g.Meta `path:"/customer/create" method:"post" sm:"新增客户" tags:"Customer"`

	CustomerFormFields
}

// CreateRes 新增客户响应
type CreateRes struct {
	Id uint `json:"id" dc:"新建客户ID"`
}

// UpdateReq 更新客户请求
type UpdateReq struct {
	g.Meta `path:"/customer/update" method:"post" sm:"更新客户" tags:"Customer"`

	Id uint `json:"id" v:"required#客户ID不能为空" dc:"客户ID"`
	CustomerFormFields
}

// UpdateRes 更新客户响应
type UpdateRes struct {
	Id uint `json:"id" dc:"客户ID"`
}

// DeleteReq 删除客户请求（后端软删除）
type DeleteReq struct {
	g.Meta `path:"/customer/delete" method:"post" sm:"删除客户" tags:"Customer"`

	Id uint `json:"id" v:"required#客户ID不能为空" dc:"客户ID"`
}

// DeleteRes 删除客户响应
type DeleteRes struct{}

// BatchDeleteReq 批量删除客户请求
type BatchDeleteReq struct {
	g.Meta `path:"/customer/batchDelete" method:"post" sm:"批量删除客户" tags:"Customer"`

	Ids []uint `json:"ids" v:"required#请选择要删除的客户" dc:"客户ID集合"`
}

// BatchDeleteRes 批量删除客户响应
type BatchDeleteRes struct{}
