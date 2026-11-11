// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Customer is the golang structure for table customer.
type Customer struct {
	Id        uint        `json:"id"        orm:"id"         description:""`                     //
	CompanyId uint        `json:"companyId" orm:"company_id" description:"所属公司ID"`               // 所属公司ID
	Name      string      `json:"name"      orm:"name"       description:"客户姓名"`                 // 客户姓名
	Phone     string      `json:"phone"     orm:"phone"      description:"手机号"`                  // 手机号
	Address   string      `json:"address"   orm:"address"    description:"房屋地址"`                 // 房屋地址
	Source    string      `json:"source"    orm:"source"     description:"客户来源：转介绍/线上/门店等"`      // 客户来源：转介绍/线上/门店等
	Status    int         `json:"status"    orm:"status"     description:"状态：1=跟进中 2=已成交 3=已流失"` // 状态：1=跟进中 2=已成交 3=已流失
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`                   // 备注
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                 // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`                 // 更新时间
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间"`                // 软删除时间
}
