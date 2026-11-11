// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerToken is the golang structure of table customer_token for DAO operations like Where/Data.
type CustomerToken struct {
	g.Meta     `orm:"table:customer_token, do:true"`
	Id         any         // 主键ID
	CustomerId any         // 业主ID，对应 customer.id
	Token      any         // 访问令牌
	Client     any         // 客户端类型：mp=小程序
	Ip         any         // 登录IP
	ExpireAt   *gtime.Time // 过期时间
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}
