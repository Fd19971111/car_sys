package v1

import "github.com/gogf/gf/v2/frame/g"

// OpLogItem 操作日志行。
type OpLogItem struct {
	Id         int64   `json:"id"         dc:"日志ID"`
	CreatedAt  string  `json:"createdAt"  dc:"操作时间 YYYY-MM-DD HH:mm:ss"`
	UserId     int64   `json:"userId"     dc:"操作人ID"`
	UserName   string  `json:"userName"   dc:"操作人"`
	Module     string  `json:"module"     dc:"业务模块"`
	Action     string  `json:"action"     dc:"动作 CREATE/UPDATE/DELETE"`
	ActionName string  `json:"actionName" dc:"动作中文名"`
	TargetType string  `json:"targetType" dc:"对象类型"`
	TargetId   int64   `json:"targetId"   dc:"对象ID"`
	Title      string  `json:"title"      dc:"操作摘要"`
	Amount     float64 `json:"amount"     dc:"涉及金额"`
	Detail     string  `json:"detail"     dc:"变更详情 JSON"`
	Ip         string  `json:"ip"         dc:"操作来源IP"`
}

// OpLogListReq 操作日志列表请求。
type OpLogListReq struct {
	g.Meta `path:"/operation-log/list" method:"get" tags:"操作日志" summary:"操作日志列表"`

	Module    string `json:"module"    dc:"业务模块，空为全部"`
	Action    string `json:"action"    dc:"动作 CREATE/UPDATE/DELETE，空为全部"`
	Keyword   string `json:"keyword"   dc:"摘要 / 操作人 / 详情关键字"`
	StartDate string `json:"startDate" dc:"开始日期 YYYY-MM-DD"`
	EndDate   string `json:"endDate"   dc:"结束日期 YYYY-MM-DD"`
	Current   int    `json:"current"   dc:"页码，从1开始"`
	Size      int    `json:"size"      dc:"每页条数"`
}

// OpLogListRes 操作日志列表响应。
type OpLogListRes struct {
	Records []*OpLogItem `json:"records" dc:"操作日志"`
	Total   int          `json:"total"   dc:"总条数"`
	Modules []string     `json:"modules" dc:"已出现过的业务模块，供筛选下拉"`
}
