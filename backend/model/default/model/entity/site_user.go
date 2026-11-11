// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteUser is the golang structure for table site_user.
type SiteUser struct {
	Id          uint        `json:"id"          orm:"id"            description:"主键ID"`               // 主键ID
	Username    string      `json:"username"    orm:"username"      description:"登录账号"`               // 登录账号
	Password    string      `json:"password"    orm:"password"      description:"密码（bcrypt加密）"`       // 密码（bcrypt加密）
	Nickname    string      `json:"nickname"    orm:"nickname"      description:"昵称/姓名"`              // 昵称/姓名
	Phone       string      `json:"phone"       orm:"phone"         description:"手机号"`                // 手机号
	Avatar      string      `json:"avatar"      orm:"avatar"        description:"头像URL"`              // 头像URL
	Role        uint        `json:"role"        orm:"role"          description:"角色：1=管理员 2=员工 3=商家"` // 角色：1=管理员 2=员工 3=商家
	Status      uint        `json:"status"      orm:"status"        description:"状态：0=禁用 1=正常"`       // 状态：0=禁用 1=正常
	MerchantId  uint        `json:"merchantId"  orm:"merchant_id"   description:"所属商家ID（商家自身或员工归属）"`  // 所属商家ID（商家自身或员工归属）
	CompanyId   uint        `json:"companyId"   orm:"company_id"    description:"所属公司ID"`             // 所属公司ID
	Remark      string      `json:"remark"      orm:"remark"        description:"备注"`                 // 备注
	LastLoginAt *gtime.Time `json:"lastLoginAt" orm:"last_login_at" description:"最后登录时间"`             // 最后登录时间
	LastLoginIp string      `json:"lastLoginIp" orm:"last_login_ip" description:"最后登录IP"`             // 最后登录IP
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:"创建时间"`               // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:"更新时间"`               // 更新时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"    description:"软删除时间"`              // 软删除时间
}
