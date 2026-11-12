// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectPayout is the golang structure for table payout.
type ProjectPayout struct {
	Id            int64       `json:"id"            orm:"id"              description:"主键"`                              // 主键
	Amount        float64     `json:"amount"        orm:"amount"          description:"付款金额"`                            // 付款金额
	PayoutType    string      `json:"payoutType"    orm:"payout_type"     description:"费用类型：材料款/人工费/施工费/设计费/管理费/销售费/其他"` // 费用类型：材料款/人工费/施工费/设计费/管理费/销售费/其他
	PayoutDate    *gtime.Time `json:"payoutDate"    orm:"payout_date"     description:"付款日期"`                            // 付款日期
	VendorName    string      `json:"vendorName"    orm:"vendor_name"     description:"收款方名称：供应商 / 工人 / 公司"`            // 收款方名称：供应商 / 工人 / 公司
	FundAccountId int64       `json:"fundAccountId" orm:"fund_account_id" description:"付款资金账户ID，对应 fin_fund_account.id"` // 付款资金账户ID，对应 fin_fund_account.id
	VoucherId     int64       `json:"voucherId"     orm:"voucher_id"      description:"自动生成的记账凭证ID，0 表示未记账"`             // 自动生成的记账凭证ID，0 表示未记账
	Remark        string      `json:"remark"        orm:"remark"          description:"备注，如：转账、结算尾款"`                    // 备注，如：转账、结算尾款
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:"创建时间"`                            // 创建时间
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"      description:"软删除时间，NULL表示未删除"`                 // 软删除时间，NULL表示未删除
}
