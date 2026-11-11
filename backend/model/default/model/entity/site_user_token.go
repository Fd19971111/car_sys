// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteUserToken is the golang structure for table site_user_token.
type SiteUserToken struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                     //
	UserId    uint        `json:"userId"    orm:"user_id"    description:"用户ID，对应 site_user.id"` // 用户ID，对应 site_user.id
	Token     string      `json:"token"     orm:"token"      description:"Token（JWT的jti或随机串）"`   // Token（JWT的jti或随机串）
	Client    string      `json:"client"    orm:"client"     description:"客户端类型：web/app"`        // 客户端类型：web/app
	Ip        string      `json:"ip"        orm:"ip"         description:"登录IP"`                 // 登录IP
	ExpireAt  *gtime.Time `json:"expireAt"  orm:"expire_at"  description:"过期时间"`                 // 过期时间
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`                     //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`                     //
}
