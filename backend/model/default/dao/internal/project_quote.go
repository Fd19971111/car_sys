// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProjectQuoteDao is the data access object for the table project_quote.
type ProjectQuoteDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ProjectQuoteColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ProjectQuoteColumns defines and stores column names for the table project_quote.
type ProjectQuoteColumns struct {
	Id        string // 主键
	ProjectId string // 关联项目ID
	ItemName  string // 项目名称，如：拆墙、水电改造
	Unit      string // 单位，如：平米、米、项
	Price     string // 单价
	Quantity  string // 数量
	Total     string // 小计 = 单价 * 数量
	Remark    string // 备注，如：不含瓷砖
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
	DeletedAt string // 软删除时间，NULL表示未删除
}

// projectQuoteColumns holds the columns for the table project_quote.
var projectQuoteColumns = ProjectQuoteColumns{
	Id:        "id",
	ProjectId: "project_id",
	ItemName:  "item_name",
	Unit:      "unit",
	Price:     "price",
	Quantity:  "quantity",
	Total:     "total",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewProjectQuoteDao creates and returns a new DAO object for table data access.
func NewProjectQuoteDao(handlers ...gdb.ModelHandler) *ProjectQuoteDao {
	return &ProjectQuoteDao{
		group:    "default",
		table:    "project_quote",
		columns:  projectQuoteColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProjectQuoteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProjectQuoteDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProjectQuoteDao) Columns() ProjectQuoteColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProjectQuoteDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProjectQuoteDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProjectQuoteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
