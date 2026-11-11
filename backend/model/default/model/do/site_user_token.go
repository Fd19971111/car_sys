// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteUserToken is the golang structure of table site_user_token for DAO operations like Where/Data.
type SiteUserToken struct {
	g.Meta    `orm:"table:site_user_token, do:true"`
	Id        any         //
	UserId    any         // 用户ID，对应 site_user.id
	Token     any         // Token（JWT的jti或随机串）
	Client    any         // 客户端类型：web/app
	Ip        any         // 登录IP
	ExpireAt  *gtime.Time // 过期时间
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
