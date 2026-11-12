package v1

import "github.com/gogf/gf/v2/frame/g"

type ListReq struct {
	g.Meta  `path:"/pricing-rule/list" method:"get" sm:"计价配置分页列表" tags:"PricingRule"`
	Current int    `json:"current" d:"1"`
	Size    int    `json:"size" d:"10"`
	Keyword string `json:"keyword" dc:"规则名称或适用车型"`
}

type PricingRuleItem struct {
	Id         uint    `json:"id"`
	Name       string  `json:"name"`
	CarType    string  `json:"carType"`
	BasePrice  float64 `json:"basePrice"`
	PerKmPrice float64 `json:"perKmPrice"`
	Enabled    bool    `json:"enabled"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

type ListRes struct {
	Current int                `json:"current"`
	Size    int                `json:"size"`
	Total   int                `json:"total"`
	Records []*PricingRuleItem `json:"records"`
}

type PricingRuleFormFields struct {
	Name       string  `json:"name"       v:"required|max-length:100#请输入规则名称|规则名称最多100个字符"`
	CarType    string  `json:"carType"    v:"required|max-length:50#请选择适用车型|适用车型最多50个字符"`
	BasePrice  float64 `json:"basePrice"  v:"min:0#起步价不能为负数"`
	PerKmPrice float64 `json:"perKmPrice" v:"min:0#超里程单价不能为负数"`
	Enabled    bool    `json:"enabled"`
}

type CreateReq struct {
	g.Meta `path:"/pricing-rule/create" method:"post" sm:"新增计价规则" tags:"PricingRule"`
	PricingRuleFormFields
}

type UpdateReq struct {
	g.Meta `path:"/pricing-rule/update" method:"post" sm:"更新计价规则" tags:"PricingRule"`
	Id     uint `json:"id" v:"required#计价规则ID不能为空"`
	PricingRuleFormFields
}

type SaveRes struct {
	Id uint `json:"id"`
}

type DeleteReq struct {
	g.Meta `path:"/pricing-rule/delete" method:"post" sm:"删除计价规则" tags:"PricingRule"`
	Id     uint `json:"id" v:"required#计价规则ID不能为空"`
}

type DeleteRes struct{}

type BatchDeleteReq struct {
	g.Meta `path:"/pricing-rule/batchDelete" method:"post" sm:"批量删除计价规则" tags:"PricingRule"`
	Ids    []uint `json:"ids" v:"required#请选择要删除的计价规则"`
}

type BatchDeleteRes struct{}
