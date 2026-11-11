// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerToken is the golang structure for table customer_token.
type CustomerToken struct {
	Id         uint64      `json:"id"         orm:"id"          description:"主键ID"`                // 主键ID
	CustomerId uint        `json:"customerId" orm:"customer_id" description:"业主ID，对应 customer.id"` // 业主ID，对应 customer.id
	Token      string      `json:"token"      orm:"token"       description:"访问令牌"`                // 访问令牌
	Client     string      `json:"client"     orm:"client"      description:"客户端类型：mp=小程序"`        // 客户端类型：mp=小程序
	Ip         string      `json:"ip"         orm:"ip"          description:"登录IP"`                // 登录IP
	ExpireAt   *gtime.Time `json:"expireAt"   orm:"expire_at"   description:"过期时间"`                // 过期时间
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`                // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`                // 更新时间
}
