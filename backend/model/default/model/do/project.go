// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Project is the golang structure of table project for DAO operations like Where/Data.
type Project struct {
	g.Meta     `orm:"table:project, do:true"`
	Id         any         //
	CustomerId any         // 客户ID
	Name       any         // 项目名称
	Address    any         // 施工地址
	ManagerId  any         // 负责人ID（site_user.id）
	Status     any         // 状态：1=待开工 2=施工中 3=已完工 4=已结算
	Amount     any         // 合同金额
	StartDate  *gtime.Time // 开工日期
	EndDate    *gtime.Time // 预计完工日期
	Remark     any         // 备注
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 软删除时间
}
