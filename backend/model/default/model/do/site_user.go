// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteUser is the golang structure of table site_user for DAO operations like Where/Data.
type SiteUser struct {
	g.Meta      `orm:"table:site_user, do:true"`
	Id          any         // 主键ID
	Username    any         // 登录账号
	Password    any         // 密码（bcrypt加密）
	Nickname    any         // 昵称/姓名
	Phone       any         // 手机号
	Avatar      any         // 头像URL
	Role        any         // 角色：1=管理员 2=员工 3=商家
	Status      any         // 状态：0=禁用 1=正常
	MerchantId  any         // 所属商家ID（商家自身或员工归属）
	Remark      any         // 备注
	LastLoginAt *gtime.Time // 最后登录时间
	LastLoginIp any         // 最后登录IP
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 软删除时间
}
