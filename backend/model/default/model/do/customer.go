// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Customer is the golang structure of table customer for DAO operations like Where/Data.
type Customer struct {
	g.Meta    `orm:"table:customer, do:true"`
	Id        any         //
	Name      any         // 客户姓名
	Phone     any         // 手机号
	Address   any         // 房屋地址
	Source    any         // 客户来源：转介绍/线上/门店等
	Status    any         // 状态：1=跟进中 2=已成交 3=已流失
	Remark    any         // 备注
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 软删除时间
}
