// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProjectPayoutDao is the data access object for the table project_payout.
type ProjectPayoutDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  ProjectPayoutColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// ProjectPayoutColumns defines and stores column names for the table project_payout.
type ProjectPayoutColumns struct {
	Id            string // 主键
	ProjectId     string // 关联项目ID，0 表示公司公共支出
	Amount        string // 付款金额
	PayoutType    string // 费用类型：材料款/人工费/施工费/设计费/管理费/销售费/其他
	PayoutDate    string // 付款日期
	VendorId      string // 供应商ID
	FundAccountId string // 付款资金账户ID，对应 fin_fund_account.id
	VoucherId     string // 自动生成的记账凭证ID，0 表示未记账
	Remark        string // 备注，如：转账、结算尾款
	CreatedAt     string // 创建时间
	DeletedAt     string // 软删除时间，NULL表示未删除
}

// projectPayoutColumns holds the columns for the table project_payout.
var projectPayoutColumns = ProjectPayoutColumns{
	Id:            "id",
	ProjectId:     "project_id",
	Amount:        "amount",
	PayoutType:    "payout_type",
	PayoutDate:    "payout_date",
	VendorId:      "vendor_id",
	FundAccountId: "fund_account_id",
	VoucherId:     "voucher_id",
	Remark:        "remark",
	CreatedAt:     "created_at",
	DeletedAt:     "deleted_at",
}

// NewProjectPayoutDao creates and returns a new DAO object for table data access.
func NewProjectPayoutDao(handlers ...gdb.ModelHandler) *ProjectPayoutDao {
	return &ProjectPayoutDao{
		group:    "default",
		table:    "project_payout",
		columns:  projectPayoutColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProjectPayoutDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProjectPayoutDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProjectPayoutDao) Columns() ProjectPayoutColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProjectPayoutDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProjectPayoutDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *ProjectPayoutDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
