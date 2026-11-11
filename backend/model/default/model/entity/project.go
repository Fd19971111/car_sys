// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Project is the golang structure for table project.
type Project struct {
	Id         uint        `json:"id"         orm:"id"          description:""`                           //
	CompanyId  uint        `json:"companyId"  orm:"company_id"  description:"所属公司ID"`                     // 所属公司ID
	CustomerId uint        `json:"customerId" orm:"customer_id" description:"客户ID"`                       // 客户ID
	Name       string      `json:"name"       orm:"name"        description:"项目名称"`                       // 项目名称
	Address    string      `json:"address"    orm:"address"     description:"施工地址"`                       // 施工地址
	ManagerId  uint        `json:"managerId"  orm:"manager_id"  description:"负责人ID（site_user.id）"`        // 负责人ID（site_user.id）
	Status     int         `json:"status"     orm:"status"      description:"状态：1=待开工 2=施工中 3=已完工 4=已结算"` // 状态：1=待开工 2=施工中 3=已完工 4=已结算
	Amount     float64     `json:"amount"     orm:"amount"      description:"合同金额"`                       // 合同金额
	StartDate  *gtime.Time `json:"startDate"  orm:"start_date"  description:"开工日期"`                       // 开工日期
	EndDate    *gtime.Time `json:"endDate"    orm:"end_date"    description:"预计完工日期"`                     // 预计完工日期
	Remark     string      `json:"remark"     orm:"remark"      description:"备注"`                         // 备注
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`                       // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`                       // 更新时间
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:"软删除时间"`                      // 软删除时间
}
