// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProjectPaymentDao is the data access object for the table payment.
type ProjectPaymentDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  ProjectPaymentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// ProjectPaymentColumns defines and stores column names for the table payment.
type ProjectPaymentColumns struct {
	Id            string // 主键
	Amount        string // 收款金额
	PaymentType   string // 款项类型：定金/首期款/中期款/尾款
	PaymentDate   string // 收款日期
	FundAccountId string // 收款资金账户ID，对应 fin_fund_account.id
	VoucherId     string // 自动生成的记账凭证ID，0 表示未记账
	Remark        string // 备注，如：微信转账
	CreatedAt     string // 创建时间
	DeletedAt     string // 软删除时间，NULL表示未删除
}

// projectPaymentColumns holds the columns for the table payment.
var projectPaymentColumns = ProjectPaymentColumns{
	Id:            "id",
	Amount:        "amount",
	PaymentType:   "payment_type",
	PaymentDate:   "payment_date",
	FundAccountId: "fund_account_id",
	VoucherId:     "voucher_id",
	Remark:        "remark",
	CreatedAt:     "created_at",
	DeletedAt:     "deleted_at",
}

// NewProjectPaymentDao creates and returns a new DAO object for table data access.
func NewProjectPaymentDao(handlers ...gdb.ModelHandler) *ProjectPaymentDao {
	return &ProjectPaymentDao{
		group:    "default",
		table:    "payment",
		columns:  projectPaymentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProjectPaymentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProjectPaymentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProjectPaymentDao) Columns() ProjectPaymentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProjectPaymentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProjectPaymentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProjectPaymentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
