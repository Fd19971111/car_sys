// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteUserTokenDao is the data access object for the table site_user_token.
type SiteUserTokenDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  SiteUserTokenColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// SiteUserTokenColumns defines and stores column names for the table site_user_token.
type SiteUserTokenColumns struct {
	Id        string //
	UserId    string // 用户ID，对应 site_user.id
	Token     string // Token（JWT的jti或随机串）
	Client    string // 客户端类型：web/app
	Ip        string // 登录IP
	ExpireAt  string // 过期时间
	CreatedAt string //
	UpdatedAt string //
}

// siteUserTokenColumns holds the columns for the table site_user_token.
var siteUserTokenColumns = SiteUserTokenColumns{
	Id:        "id",
	UserId:    "user_id",
	Token:     "token",
	Client:    "client",
	Ip:        "ip",
	ExpireAt:  "expire_at",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewSiteUserTokenDao creates and returns a new DAO object for table data access.
func NewSiteUserTokenDao(handlers ...gdb.ModelHandler) *SiteUserTokenDao {
	return &SiteUserTokenDao{
		group:    "default",
		table:    "site_user_token",
		columns:  siteUserTokenColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteUserTokenDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteUserTokenDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteUserTokenDao) Columns() SiteUserTokenColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteUserTokenDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteUserTokenDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteUserTokenDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
