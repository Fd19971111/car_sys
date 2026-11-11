// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectQuote is the golang structure for table project_quote.
type ProjectQuote struct {
	Id        int64       `json:"id"        orm:"id"         description:"主键"`              // 主键
	ProjectId int64       `json:"projectId" orm:"project_id" description:"关联项目ID"`          // 关联项目ID
	ItemName  string      `json:"itemName"  orm:"item_name"  description:"项目名称，如：拆墙、水电改造"`  // 项目名称，如：拆墙、水电改造
	Unit      string      `json:"unit"      orm:"unit"       description:"单位，如：平米、米、项"`     // 单位，如：平米、米、项
	Price     float64     `json:"price"     orm:"price"      description:"单价"`              // 单价
	Quantity  float64     `json:"quantity"  orm:"quantity"   description:"数量"`              // 数量
	Total     float64     `json:"total"     orm:"total"      description:"小计 = 单价 * 数量"`    // 小计 = 单价 * 数量
	Remark    string      `json:"remark"    orm:"remark"     description:"备注，如：不含瓷砖"`       // 备注，如：不含瓷砖
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`            // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`            // 更新时间
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间，NULL表示未删除"` // 软删除时间，NULL表示未删除
}
