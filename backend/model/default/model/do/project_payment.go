// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectPayment is the golang structure of table payment for DAO operations like Where/Data.
type ProjectPayment struct {
	g.Meta        `orm:"table:payment, do:true"`
	Id            any         // 主键
	RideOrderId   any         // 关联行程订单ID
	Amount        any         // 收款金额
	PaymentType   any         // 款项类型：定金/首期款/中期款/尾款
	PaymentDate   *gtime.Time // 收款日期
	FundAccountId any         // 收款资金账户ID，对应 fin_fund_account.id
	VoucherId     any         // 自动生成的记账凭证ID，0 表示未记账
	Remark        any         // 备注，如：微信转账
	CreatedAt     *gtime.Time // 创建时间
	DeletedAt     *gtime.Time // 软删除时间，NULL表示未删除
}
