package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 业主端（小程序）只读接口 =====================
//
// 说明：本期为了先看效果，ownerId 由前端显式传入，接口不挂 Auth 中间件。
// 后续接入微信登录后，只需把 customerId 的来源从「前端传参」改为「从 token 解析」，
// 接口路径与响应结构都不用变，前端无需改动。

// OwnerProfile 业主资料
type OwnerProfile struct {
	Id      uint   `json:"id"      dc:"业主ID（customer.id）"`
	Name    string `json:"name"    dc:"姓名"`
	Phone   string `json:"phone"   dc:"手机号"`
	Address string `json:"address" dc:"房屋地址"`
}

// ProjectCard 项目卡片（首页列表用）
type ProjectCard struct {
	Id               uint    `json:"id"               dc:"项目ID"`
	Name             string  `json:"name"             dc:"项目名称"`
	Address          string  `json:"address"          dc:"施工地址"`
	Status           int     `json:"status"           dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	Amount           float64 `json:"amount"           dc:"合同金额"`
	PaidAmount       float64 `json:"paidAmount"       dc:"已收款金额"`
	ManagerName      string  `json:"managerName"      dc:"项目负责人"`
	StartDate        string  `json:"startDate"        dc:"开工日期"`
	EndDate          string  `json:"endDate"          dc:"预计完工日期"`
	LatestProgress   string  `json:"latestProgress"   dc:"最近一条进度"`
	LatestProgressAt string  `json:"latestProgressAt" dc:"最近进度时间"`
	ProgressPercent  int     `json:"progressPercent"  dc:"施工进度百分比"`
}

// OverviewSummary 首页资金 / 进度汇总
type OverviewSummary struct {
	ProjectCount    int     `json:"projectCount"    dc:"项目数量"`
	ContractAmount  float64 `json:"contractAmount"  dc:"合同总金额"`
	PaidAmount      float64 `json:"paidAmount"      dc:"已收款总额"`
	UnpaidAmount    float64 `json:"unpaidAmount"    dc:"待收款金额"`
	ProgressPercent int     `json:"progressPercent" dc:"整体施工进度百分比"`
}

// ===================== 业主登录 =====================

// LoginReq 业主登录请求
//
// 身份凭据是「手机号 + 验证码」，手机号需与 customer.phone 一致。
// 验证码当前由配置 app.mpLoginCode 提供（开发阶段固定值），接入短信后改为动态下发。
type LoginReq struct {
	g.Meta `path:"/mp/auth/login" method:"post" sm:"业主登录" tags:"MiniProgram"`

	Phone string `json:"phone" v:"required|regex:^1[3-9][0-9]{9}$#请输入手机号|手机号格式不正确" dc:"手机号"`
	Code  string `json:"code"  v:"required#请输入验证码" dc:"验证码"`
}

// LoginRes 业主登录响应
type LoginRes struct {
	Token    string        `json:"token"     dc:"访问令牌，后续请求放在 Authorization: Bearer <token>"`
	ExpireAt string        `json:"expireAt"  dc:"令牌过期时间"`
	Owner    *OwnerProfile `json:"owner"     dc:"业主资料"`
}

// LogoutReq 业主退出登录请求
type LogoutReq struct {
	g.Meta `path:"/mp/auth/logout" method:"post" sm:"业主退出登录" tags:"MiniProgram"`
}

// LogoutRes 业主退出登录响应
type LogoutRes struct{}

// ===================== 业主首页聚合 =====================

// OverviewReq 业主首页聚合请求
//
// 注意：这里**不再接收 customerId** —— 业主身份一律从令牌解析，
// 否则前端可以随便传别人的ID，数据隔离形同虚设。
type OverviewReq struct {
	g.Meta `path:"/mp/overview" method:"get" sm:"业主首页聚合" tags:"MiniProgram"`
}

// OverviewRes 业主首页聚合响应
type OverviewRes struct {
	Owner    *OwnerProfile    `json:"owner"    dc:"业主资料"`
	Summary  *OverviewSummary `json:"summary"  dc:"资金与进度汇总"`
	Projects []*ProjectCard   `json:"projects" dc:"项目列表"`
}

// ===================== 项目详情聚合 =====================

// ProjectReq 项目详情聚合请求
type ProjectReq struct {
	g.Meta `path:"/mp/project" method:"get" sm:"项目详情聚合" tags:"MiniProgram"`

	Id uint `json:"id" v:"required#项目ID不能为空" dc:"项目ID"`
}

// ProjectInfo 项目基本信息
type ProjectInfo struct {
	Id           uint   `json:"id"           dc:"项目ID"`
	Name         string `json:"name"         dc:"项目名称"`
	Address      string `json:"address"      dc:"施工地址"`
	Status       int    `json:"status"       dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	ManagerName  string `json:"managerName"  dc:"项目负责人"`
	ManagerPhone string `json:"managerPhone" dc:"项目负责人手机号"`
	StartDate    string `json:"startDate"    dc:"开工日期"`
	EndDate      string `json:"endDate"      dc:"预计完工日期"`
	Remark       string `json:"remark"       dc:"备注"`
	CreatedAt    string `json:"createdAt"    dc:"签约时间"`
}

// ProjectSummary 项目资金与进度汇总
type ProjectSummary struct {
	ContractAmount  float64 `json:"contractAmount"  dc:"合同金额"`
	QuoteAmount     float64 `json:"quoteAmount"     dc:"报价合计"`
	PaidAmount      float64 `json:"paidAmount"      dc:"已收款"`
	UnpaidAmount    float64 `json:"unpaidAmount"    dc:"待收款"`
	PaidPercent     int     `json:"paidPercent"     dc:"收款进度百分比"`
	ProgressPercent int     `json:"progressPercent" dc:"施工进度百分比"`
}

// ProgressLog 施工进度记录
type ProgressLog struct {
	Id        uint     `json:"id"        dc:"进度记录ID"`
	Content   string   `json:"content"   dc:"进度内容"`
	Images    []string `json:"images"    dc:"现场照片"`
	UserName  string   `json:"userName"  dc:"记录人"`
	CreatedAt string   `json:"createdAt" dc:"记录时间"`
}

// QuoteItem 报价明细项
type QuoteItem struct {
	Id       int64   `json:"id"       dc:"明细ID"`
	ItemName string  `json:"itemName" dc:"报价项名称"`
	Unit     string  `json:"unit"     dc:"单位"`
	Price    float64 `json:"price"    dc:"单价"`
	Quantity float64 `json:"quantity" dc:"数量"`
	Total    float64 `json:"total"    dc:"小计"`
	Remark   string  `json:"remark"   dc:"备注"`
}

// PaymentItem 收款记录项
type PaymentItem struct {
	Id          int64   `json:"id"          dc:"收款记录ID"`
	PaymentType string  `json:"paymentType" dc:"款项类型"`
	Amount      float64 `json:"amount"      dc:"收款金额"`
	PaymentDate string  `json:"paymentDate" dc:"收款日期"`
	Remark      string  `json:"remark"      dc:"备注"`
}

// ProjectRes 项目详情聚合响应
type ProjectRes struct {
	Project  *ProjectInfo    `json:"project"  dc:"项目信息"`
	Owner    *OwnerProfile   `json:"owner"    dc:"业主资料"`
	Summary  *ProjectSummary `json:"summary"  dc:"资金与进度汇总"`
	Logs     []*ProgressLog  `json:"logs"     dc:"施工进度（按时间倒序）"`
	Quotes   []*QuoteItem    `json:"quotes"   dc:"报价明细"`
	Payments []*PaymentItem  `json:"payments" dc:"收款记录（按时间倒序）"`
}

// ===================== 新增施工进度 =====================

// ProgressCreateReq 新增施工进度请求（小程序端）
type ProgressCreateReq struct {
	g.Meta `path:"/mp/progress" method:"post" sm:"新增施工进度（小程序端）" tags:"MiniProgram"`

	ProjectId uint     `json:"projectId" v:"required#请选择项目" dc:"项目ID"`
	Content   string   `json:"content"   v:"required|max-length:500#请填写进度描述|进度描述最多 500 个字符" dc:"进度描述"`
	Images    []string `json:"images" dc:"图片相对路径列表（上传接口返回的 path），最多 9 张"`
}

// ProgressCreateRes 新增施工进度响应
type ProgressCreateRes struct {
	Id int64 `json:"id" dc:"新建的进度记录ID"`
}
