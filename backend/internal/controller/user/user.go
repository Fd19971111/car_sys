package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	_userApi "engineServ/api/user"
	v1 "engineServ/api/user/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() _userApi.IUserV1 { return &Controller{} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	base := `FROM site_user u LEFT JOIN sys_department d ON d.id=u.department_id WHERE u.deleted_at IS NULL AND u.company_id=?`
	args := []any{currentCompanyId(ctx)}
	if req.Keyword != "" {
		base += ` AND (u.nickname LIKE ? OR u.username LIKE ? OR u.phone LIKE ?)`
		kw := "%" + req.Keyword + "%"
		args = append(args, kw, kw, kw)
	}
	if req.Status != nil {
		base += ` AND u.status=?`
		args = append(args, *req.Status)
	}
	if req.Role != nil {
		base += ` AND u.role=?`
		args = append(args, *req.Role)
	}

	total, _ := g.DB().GetValue(ctx, "SELECT COUNT(*)"+base, args...)
	offset := (req.Current - 1) * req.Size
	rows, err := g.DB().GetAll(ctx, "SELECT u.id,u.username,u.nickname,u.phone,u.role,u.status,u.department_id,u.last_login_at,u.created_at,d.name department_name "+base+" ORDER BY u.role,u.id LIMIT ?,?", append(args, offset, req.Size)...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户列表失败")
	}

	records := make([]*v1.UserItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, rowToItem(row))
	}
	return &v1.ListRes{Current: req.Current, Size: req.Size, Total: total.Int(), Records: records}, nil
}

func (c *Controller) ResetPwd(ctx context.Context, req *v1.ResetPwdReq) (*v1.ResetPwdRes, error) {
	phone, err := g.DB().GetValue(ctx, "SELECT phone FROM site_user WHERE id=? AND company_id=? AND deleted_at IS NULL", req.Id, currentCompanyId(ctx))
	if err != nil || phone.IsEmpty() {
		return nil, gerror.New("用户不存在")
	}
	pwd := phone.String()
	if len(pwd) > 6 {
		pwd = pwd[len(pwd)-6:]
	}
	pwdHash := fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))
	_, err = g.DB().Exec(ctx, "UPDATE site_user SET password=? WHERE id=? AND company_id=?", pwdHash, req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "重置密码失败")
	}
	return &v1.ResetPwdRes{Password: pwd}, nil
}

func (c *Controller) Toggle(ctx context.Context, req *v1.ToggleReq) (*v1.ToggleRes, error) {
	_, err := g.DB().Exec(ctx, "UPDATE site_user SET status=? WHERE id=? AND company_id=?", req.Status, req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "操作失败")
	}
	return &v1.ToggleRes{}, nil
}

func (c *Controller) AvailableStaff(ctx context.Context, _ *v1.AvailableStaffReq) (*v1.AvailableStaffRes, error) {
	rows, err := g.DB().GetAll(ctx, `SELECT id,nickname,phone FROM site_user
		WHERE company_id=? AND role=2 AND status=1 AND deleted_at IS NULL AND account_created=0
		ORDER BY id`, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询待开通账号员工失败")
	}
	items := make([]v1.StaffOption, 0, len(rows))
	for _, row := range rows {
		items = append(items, v1.StaffOption{Id: row["id"].Uint(), Name: row["nickname"].String(), Phone: row["phone"].String()})
	}
	return &v1.AvailableStaffRes{Items: items}, nil
}

func (c *Controller) CreateAccount(ctx context.Context, req *v1.CreateAccountReq) (*v1.CreateAccountRes, error) {
	row, err := g.DB().GetOne(ctx, "SELECT id,nickname,phone,role,account_created FROM site_user WHERE id=? AND company_id=? AND deleted_at IS NULL", req.StaffId, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询员工失败")
	}
	if row.IsEmpty() {
		return nil, gerror.New("员工不存在")
	}
	if row["role"].Int() != 2 || row["account_created"].Int() != 0 {
		return nil, gerror.New("该员工已创建登录账号")
	}
	// 检查用户名是否已被占用
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE username=? AND id!=? AND company_id=? AND deleted_at IS NULL", req.Username, req.StaffId, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "检查登录用户名失败")
	}
	if count.Int() > 0 {
		return nil, gerror.New("用户名已被占用")
	}
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.Wrap(err, "密码加密失败")
	}
	_, err = g.DB().Exec(ctx, "UPDATE site_user SET username=?,password=?,account_created=1 WHERE id=? AND company_id=?", req.Username, pwdHash, req.StaffId, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "创建账号失败")
	}
	return &v1.CreateAccountRes{Id: req.StaffId}, nil
}

func (c *Controller) UpdateAccount(ctx context.Context, req *v1.UpdateAccountReq) (*v1.UpdateAccountRes, error) {
	exists, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE id=? AND company_id=? AND deleted_at IS NULL", req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	if exists.Int() == 0 {
		return nil, gerror.New("用户不存在")
	}
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE username=? AND id!=? AND company_id=? AND deleted_at IS NULL", req.Username, req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "检查登录用户名失败")
	}
	if count.Int() > 0 {
		return nil, gerror.New("用户名已被占用")
	}
	_, err = g.DB().Exec(ctx, "UPDATE site_user SET username=?,role=? WHERE id=? AND company_id=? AND deleted_at IS NULL", req.Username, req.Role, req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "更新账号失败")
	}
	return &v1.UpdateAccountRes{}, nil
}

func (c *Controller) ChangePwd(ctx context.Context, req *v1.ChangePwdReq) (*v1.ChangePwdRes, error) {
	exists, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE id=? AND company_id=? AND deleted_at IS NULL", req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	if exists.Int() == 0 {
		return nil, gerror.New("用户不存在")
	}
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.Wrap(err, "密码加密失败")
	}
	_, err = g.DB().Exec(ctx, "UPDATE site_user SET password=? WHERE id=? AND company_id=?", string(pwdHash), req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "修改密码失败")
	}
	return &v1.ChangePwdRes{}, nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func rowToItem(row gdb.Record) *v1.UserItem {
	role := row["role"].Int()
	roleName := "员工"
	if role == 1 {
		roleName = "管理员"
	}
	lastLogin := ""
	if t := row["last_login_at"].Time(); !t.IsZero() {
		lastLogin = t.Format("2006-01-02 15:04")
	}
	return &v1.UserItem{
		Id: row["id"].Uint(), Username: row["username"].String(), Nickname: row["nickname"].String(),
		Phone: row["phone"].String(), Role: role, RoleName: roleName, Status: row["status"].Int(),
		DepartmentId: row["department_id"].Uint(), DepartmentName: row["department_name"].String(),
		LastLoginAt: lastLogin,
		CreatedAt:   row["created_at"].Time().Format("2006-01-02 15:04"),
	}
}
