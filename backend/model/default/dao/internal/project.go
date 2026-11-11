// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProjectDao is the data access object for the table project.
type ProjectDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ProjectColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ProjectColumns defines and stores column names for the table project.
type ProjectColumns struct {
	Id         string //
	CustomerId string // 客户ID
	Name       string // 项目名称
	Address    string // 施工地址
	ManagerId  string // 负责人ID（site_user.id）
	Status     string // 状态：1=待开工 2=施工中 3=已完工 4=已结算
	Amount     string // 合同金额
	StartDate  string // 开工日期
	EndDate    string // 预计完工日期
	Remark     string // 备注
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 软删除时间
}

// projectColumns holds the columns for the table project.
var projectColumns = ProjectColumns{
	Id:         "id",
	CustomerId: "customer_id",
	Name:       "name",
	Address:    "address",
	ManagerId:  "manager_id",
	Status:     "status",
	Amount:     "amount",
	StartDate:  "start_date",
	EndDate:    "end_date",
	Remark:     "remark",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewProjectDao creates and returns a new DAO object for table data access.
func NewProjectDao(handlers ...gdb.ModelHandler) *ProjectDao {
	return &ProjectDao{
		group:    "default",
		table:    "project",
		columns:  projectColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProjectDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProjectDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProjectDao) Columns() ProjectColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProjectDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProjectDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProjectDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
