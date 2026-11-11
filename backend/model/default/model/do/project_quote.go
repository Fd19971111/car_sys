// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectQuote is the golang structure of table project_quote for DAO operations like Where/Data.
type ProjectQuote struct {
	g.Meta    `orm:"table:project_quote, do:true"`
	Id        any         // 主键
	ProjectId any         // 关联项目ID
	ItemName  any         // 项目名称，如：拆墙、水电改造
	Unit      any         // 单位，如：平米、米、项
	Price     any         // 单价
	Quantity  any         // 数量
	Total     any         // 小计 = 单价 * 数量
	Remark    any         // 备注，如：不含瓷砖
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 软删除时间，NULL表示未删除
}
