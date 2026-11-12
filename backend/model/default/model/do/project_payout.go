// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectPayout is the golang structure of table project_payout for DAO operations like Where/Data.
type ProjectPayout struct {
	g.Meta        `orm:"table:project_payout, do:true"`
	Id            any         // 主键
	ProjectId     any         // 关联项目ID，0 表示公司公共支出
	Amount        any         // 付款金额
	PayoutType    any         // 费用类型：材料款/人工费/施工费/设计费/管理费/销售费/其他
	PayoutDate    *gtime.Time // 付款日期
	VendorName    any         // 收款方名称：供应商 / 工人 / 公司
	FundAccountId any         // 付款资金账户ID，对应 fin_fund_account.id
	VoucherId     any         // 自动生成的记账凭证ID，0 表示未记账
	Remark        any         // 备注，如：转账、结算尾款
	CreatedAt     *gtime.Time // 创建时间
	DeletedAt     *gtime.Time // 软删除时间，NULL表示未删除
}
