package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 项目分页列表 =====================

// ListReq 项目分页列表请求
type ListReq struct {
	g.Meta `path:"/project/list" method:"get" sm:"项目分页列表" tags:"Project"`

	Current    int    `json:"current"    d:"1"  dc:"当前页码，默认 1"`
	Size       int    `json:"size"       d:"10" dc:"每页条数，默认 10"`
	Keyword    string `json:"keyword"           dc:"关键词：匹配项目名称 / 施工地址 / 客户姓名"`
	Status     *int   `json:"status"            dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	CustomerId *uint  `json:"customerId"        dc:"所属客户ID"`
}

// ProjectListItem 项目列表行
type ProjectListItem struct {
	Id               uint    `json:"id"               dc:"项目ID"`
	CustomerId       uint    `json:"customerId"       dc:"所属客户ID"`
	CustomerName     string  `json:"customerName"     dc:"所属客户姓名"`
	Name             string  `json:"name"             dc:"项目名称"`
	Address          string  `json:"address"          dc:"施工地址"`
	ManagerId        uint    `json:"managerId"        dc:"负责人ID"`
	ManagerName      string  `json:"managerName"      dc:"负责人姓名"`
	Status           int     `json:"status"           dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	Amount           float64 `json:"amount"           dc:"合同金额"`
	StartDate        string  `json:"startDate"        dc:"开工日期"`
	EndDate          string  `json:"endDate"          dc:"预计完工日期"`
	Remark           string  `json:"remark"           dc:"备注"`
	LatestProgress   string  `json:"latestProgress"   dc:"最近一条进度记录"`
	LatestProgressAt string  `json:"latestProgressAt" dc:"最近进度记录时间"`
	CreatedAt        string  `json:"createdAt"        dc:"创建时间"`
	UpdatedAt        string  `json:"updatedAt"        dc:"更新时间"`
}

// ListRes 项目分页列表响应
type ListRes struct {
	Current int                `json:"current" dc:"当前页码"`
	Size    int                `json:"size"    dc:"每页条数"`
	Total   int                `json:"total"   dc:"总条数"`
	Records []*ProjectListItem `json:"records" dc:"项目列表"`
}

// ===================== 项目概览统计 =====================

// StatsReq 项目概览统计请求
type StatsReq struct {
	g.Meta `path:"/project/stats" method:"get" sm:"项目概览统计" tags:"Project"`
}

// StatsRes 项目概览统计响应
type StatsRes struct {
	Total       int     `json:"total"       dc:"项目总数"`
	Pending     int     `json:"pending"     dc:"待开工数量"`
	Building    int     `json:"building"    dc:"施工中数量"`
	Finished    int     `json:"finished"    dc:"已完工数量"`
	Settled     int     `json:"settled"     dc:"已结算数量"`
	TotalAmount float64 `json:"totalAmount" dc:"合同总金额"`
}

// ===================== 表单下拉选项 =====================

// OptionsReq 表单下拉选项请求
type OptionsReq struct {
	g.Meta `path:"/project/options" method:"get" sm:"项目表单下拉选项" tags:"Project"`
}

// OptionItem 下拉选项
type OptionItem struct {
	Value uint   `json:"value" dc:"选项值"`
	Label string `json:"label" dc:"选项文案"`
}

// OptionsRes 表单下拉选项响应
type OptionsRes struct {
	Customers []*OptionItem `json:"customers" dc:"客户下拉选项"`
	Managers  []*OptionItem `json:"managers"  dc:"负责人下拉选项"`
}

// ===================== 项目下拉选项（供报价 / 收款模块选择项目） =====================

// SelectOptionsReq 项目下拉选项请求
type SelectOptionsReq struct {
	g.Meta `path:"/project/selectOptions" method:"get" sm:"项目下拉选项" tags:"Project"`
}

// ProjectSelectItem 项目下拉选项
type ProjectSelectItem struct {
	Value uint   `json:"value" dc:"项目ID"`
	Label string `json:"label" dc:"项目名称 · 客户姓名"`
	// 以下字段供表单联动展示，不参与选择
	Amount float64 `json:"amount" dc:"合同金额"`
	Status int     `json:"status" dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
}

// SelectOptionsRes 项目下拉选项响应
type SelectOptionsRes struct {
	Projects []*ProjectSelectItem `json:"projects" dc:"项目下拉选项"`
}

// ===================== 项目详情 =====================

// DetailReq 项目详情请求
type DetailReq struct {
	g.Meta `path:"/project/detail" method:"get" sm:"项目详情" tags:"Project"`

	Id uint `json:"id" v:"required#项目ID不能为空" dc:"项目ID"`
}

// ProjectLogItem 施工进度记录
type ProjectLogItem struct {
	Id        uint     `json:"id"        dc:"进度记录ID"`
	Content   string   `json:"content"   dc:"进度内容"`
	Images    []string `json:"images"    dc:"进度图片URL集合"`
	UserName  string   `json:"userName"  dc:"记录人姓名"`
	LogDate   string   `json:"logDate"   dc:"进度业务日期（YYYY-MM-DD，进度实际发生日期）"`
	CreatedAt string   `json:"createdAt" dc:"记录时间"`
}

// DetailRes 项目详情响应
type DetailRes struct {
	Id               uint              `json:"id"               dc:"项目ID"`
	CustomerId       uint              `json:"customerId"       dc:"所属客户ID"`
	CustomerName     string            `json:"customerName"     dc:"所属客户姓名"`
	CustomerPhone    string            `json:"customerPhone"    dc:"客户手机号"`
	CustomerAddress  string            `json:"customerAddress"  dc:"客户房屋地址"`
	Name             string            `json:"name"             dc:"项目名称"`
	Address          string            `json:"address"          dc:"施工地址"`
	ManagerId        uint              `json:"managerId"        dc:"负责人ID"`
	ManagerName      string            `json:"managerName"      dc:"负责人姓名"`
	Status           int               `json:"status"           dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	Amount           float64           `json:"amount"           dc:"合同金额"`
	StartDate        string            `json:"startDate"        dc:"开工日期"`
	EndDate          string            `json:"endDate"          dc:"预计完工日期"`
	Remark           string            `json:"remark"           dc:"备注"`
	LatestProgress   string            `json:"latestProgress"   dc:"最近一条进度记录"`
	LatestProgressAt string            `json:"latestProgressAt" dc:"最近进度记录时间"`
	CreatedAt        string            `json:"createdAt"        dc:"创建时间"`
	UpdatedAt        string            `json:"updatedAt"        dc:"更新时间"`
	Logs             []*ProjectLogItem `json:"logs"             dc:"施工进度记录（按时间倒序）"`
}

// ===================== 新增 / 编辑 / 删除 =====================

// ProjectFormFields 项目表单公共字段（project 表结构固定，不新增业务列）
type ProjectFormFields struct {
	CustomerId uint    `json:"customerId" v:"required#请选择所属客户"                                    dc:"所属客户ID"`
	Name       string  `json:"name"       v:"required|length:2,50#请输入项目名称|项目名称长度为 2-50 个字符" dc:"项目名称"`
	Address    string  `json:"address"    v:"max-length:60#施工地址最多 60 个字符"                       dc:"施工地址"`
	ManagerId  uint    `json:"managerId"  dc:"负责人ID（site_user.id）"`
	Status     int     `json:"status"     v:"required|in:1,2,3,4#请选择项目状态|项目状态不正确"             dc:"项目状态：1=待开工 2=施工中 3=已完工 4=已结算"`
	Amount     float64 `json:"amount"     v:"min:0#合同金额不能为负数"                                   dc:"合同金额"`
	StartDate  string  `json:"startDate"  dc:"开工日期（YYYY-MM-DD）"`
	EndDate    string  `json:"endDate"    dc:"预计完工日期（YYYY-MM-DD）"`
	Remark     string  `json:"remark"     v:"max-length:200#备注最多 200 个字符"                        dc:"备注"`
}

// CreateReq 新增项目请求
type CreateReq struct {
	g.Meta `path:"/project/create" method:"post" sm:"新增项目" tags:"Project"`

	ProjectFormFields
}

// CreateRes 新增项目响应
type CreateRes struct {
	Id uint `json:"id" dc:"新建项目ID"`
}

// UpdateReq 更新项目请求
type UpdateReq struct {
	g.Meta `path:"/project/update" method:"post" sm:"更新项目" tags:"Project"`

	Id uint `json:"id" v:"required#项目ID不能为空" dc:"项目ID"`
	ProjectFormFields
}

// UpdateRes 更新项目响应
type UpdateRes struct {
	Id uint `json:"id" dc:"项目ID"`
}

// DeleteReq 删除项目请求（后端软删除）
type DeleteReq struct {
	g.Meta `path:"/project/delete" method:"post" sm:"删除项目" tags:"Project"`

	Id uint `json:"id" v:"required#项目ID不能为空" dc:"项目ID"`
}

// DeleteRes 删除项目响应
type DeleteRes struct{}

// BatchDeleteReq 批量删除项目请求
type BatchDeleteReq struct {
	g.Meta `path:"/project/batchDelete" method:"post" sm:"批量删除项目" tags:"Project"`

	Ids []uint `json:"ids" v:"required#请选择要删除的项目" dc:"项目ID集合"`
}

// BatchDeleteRes 批量删除项目响应
type BatchDeleteRes struct{}

// ===================== 新增施工进度（project_log 写入） =====================

// ProgressCreateReq 新增施工进度请求（后台端，记录人取当前登录用户）
type ProgressCreateReq struct {
	g.Meta `path:"/project/progressCreate" method:"post" sm:"新增施工进度" tags:"Project"`

	ProjectId uint     `json:"projectId" v:"required#请选择项目" dc:"项目ID"`
	Content   string   `json:"content"   v:"required|max-length:500#请填写进度描述|进度描述最多 500 个字符" dc:"进度描述"`
	Images    []string `json:"images" dc:"图片相对路径列表（上传接口返回的 path），最多 9 张"`
	LogDate   string   `json:"logDate" dc:"进度业务日期（YYYY-MM-DD），为空默认当天；补录历史进度时传入真实日期"`
}

// ProgressCreateRes 新增施工进度响应
type ProgressCreateRes struct {
	Id int64 `json:"id" dc:"新建的进度记录ID"`
}
