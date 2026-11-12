// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectPayment is the golang structure for table payment.
type ProjectPayment struct {
	Id            int64       `json:"id"            orm:"id"              description:"主键"`                     // 主键
	Amount        float64     `json:"amount"        orm:"amount"          description:"收款金额"`                     // 收款金额
	PaymentType   string      `json:"paymentType"   orm:"payment_type"    description:"款项类型：定金/首期款/中期款/尾款"`      // 款项类型：定金/首期款/中期款/尾款
	PaymentDate   *gtime.Time `json:"paymentDate"   orm:"payment_date"    description:"收款日期"`                     // 收款日期
	FundAccountId int64       `json:"fundAccountId" orm:"fund_account_id" description:"收款资金账户ID，对应 fin_fund_account.id"` // 收款资金账户ID，对应 fin_fund_account.id
	VoucherId     int64       `json:"voucherId"     orm:"voucher_id"      description:"自动生成的记账凭证ID，0 表示未记账"`         // 自动生成的记账凭证ID，0 表示未记账
	Remark        string      `json:"remark"        orm:"remark"          description:"备注，如：微信转账"`                // 备注，如：微信转账
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:"创建时间"`                     // 创建时间
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"      description:"软删除时间，NULL表示未删除"`          // 软删除时间，NULL表示未删除
}
