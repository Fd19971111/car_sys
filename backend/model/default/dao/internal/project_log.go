// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProjectLogDao is the data access object for the table project_log.
type ProjectLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ProjectLogColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ProjectLogColumns defines and stores column names for the table project_log.
type ProjectLogColumns struct {
	Id        string //
	ProjectId string // 项目ID
	UserId    string // 记录人ID（site_user.id）
	Content   string // 进度内容
	LogDate   string // 进度业务日期（YYYY-MM-DD，进度实际发生日期）
	Images    string // 图片URL，多个用逗号分隔
	CreatedAt string // 创建时间
}

// projectLogColumns holds the columns for the table project_log.
var projectLogColumns = ProjectLogColumns{
	Id:        "id",
	ProjectId: "project_id",
	UserId:    "user_id",
	Content:   "content",
	LogDate:   "log_date",
	Images:    "images",
	CreatedAt: "created_at",
}

// NewProjectLogDao creates and returns a new DAO object for table data access.
func NewProjectLogDao(handlers ...gdb.ModelHandler) *ProjectLogDao {
	return &ProjectLogDao{
		group:    "default",
		table:    "project_log",
		columns:  projectLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProjectLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProjectLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProjectLogDao) Columns() ProjectLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProjectLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProjectLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProjectLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
