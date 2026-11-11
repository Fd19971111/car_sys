// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteUserDao is the data access object for the table site_user.
type SiteUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SiteUserColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SiteUserColumns defines and stores column names for the table site_user.
type SiteUserColumns struct {
	Id          string // 主键ID
	Username    string // 登录账号
	Password    string // 密码（bcrypt加密）
	Nickname    string // 昵称/姓名
	Phone       string // 手机号
	Avatar      string // 头像URL
	Role        string // 角色：1=管理员 2=员工 3=商家
	Status      string // 状态：0=禁用 1=正常
	MerchantId  string // 所属商家ID（商家自身或员工归属）
	Remark      string // 备注
	LastLoginAt string // 最后登录时间
	LastLoginIp string // 最后登录IP
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 软删除时间
}

// siteUserColumns holds the columns for the table site_user.
var siteUserColumns = SiteUserColumns{
	Id:          "id",
	Username:    "username",
	Password:    "password",
	Nickname:    "nickname",
	Phone:       "phone",
	Avatar:      "avatar",
	Role:        "role",
	Status:      "status",
	MerchantId:  "merchant_id",
	Remark:      "remark",
	LastLoginAt: "last_login_at",
	LastLoginIp: "last_login_ip",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewSiteUserDao creates and returns a new DAO object for table data access.
func NewSiteUserDao(handlers ...gdb.ModelHandler) *SiteUserDao {
	return &SiteUserDao{
		group:    "default",
		table:    "site_user",
		columns:  siteUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteUserDao) Columns() SiteUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
