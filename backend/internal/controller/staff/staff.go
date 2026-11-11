package staff

import (
	"context"
	"strings"

	staffApi "engineServ/api/staff"
	v1 "engineServ/api/staff/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() staffApi.IStaffV1 { return &Controller{} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	base := `FROM site_user u LEFT JOIN sys_department d ON d.id=u.department_id WHERE u.deleted_at IS NULL AND u.company_id=?`
	args := []any{currentCompanyId(ctx)}

	if req.Keyword != "" {
		base += ` AND (u.nickname LIKE ? OR u.phone LIKE ?)`
		kw := "%" + req.Keyword + "%"
		args = append(args, kw, kw)
	}
	if req.DepartmentId != nil && *req.DepartmentId > 0 {
		base += ` AND u.department_id=?`
		args = append(args, *req.DepartmentId)
	}
	if req.Status != nil {
		base += ` AND u.status=?`
		args = append(args, *req.Status)
	} else {
		base += ` AND u.status=1`
	}

	total, _ := g.DB().GetValue(ctx, "SELECT COUNT(*)"+base, args...)
	offset := (req.Current - 1) * req.Size
	rows, err := g.DB().GetAll(ctx, "SELECT u.id,u.username,u.nickname,u.phone,u.role,u.status,u.department_id,u.remark,u.created_at,u.updated_at,d.name department_name "+base+" ORDER BY u.role,u.id LIMIT ?,?", append(args, offset, req.Size)...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询人员列表失败")
	}

	records := make([]*v1.StaffItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, rowToItem(row))
	}
	return &v1.ListRes{Current: req.Current, Size: req.Size, Total: total.Int(), Records: records}, nil
}

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (*v1.DetailRes, error) {
	row, _ := g.DB().GetOne(ctx, "SELECT u.id,u.username,u.nickname,u.phone,u.role,u.status,u.department_id,u.remark,u.created_at,u.updated_at,d.name department_name FROM site_user u LEFT JOIN sys_department d ON d.id=u.department_id WHERE u.id=? AND u.company_id=? AND u.deleted_at IS NULL", req.Id, currentCompanyId(ctx))
	if row.IsEmpty() {
		return nil, gerror.New("人员不存在")
	}
	return rowToItem(row), nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	count, _ := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE phone=? AND company_id=? AND deleted_at IS NULL", req.Phone, currentCompanyId(ctx))
	if count.Int() > 0 {
		return nil, gerror.New("该手机号已存在")
	}
	pwd := req.Phone
	if len(pwd) > 6 {
		pwd = pwd[len(pwd)-6:]
	}
	result, err := g.DB().Exec(ctx, "INSERT INTO site_user(username,nickname,phone,password,role,status,department_id,company_id,account_created,remark) VALUES(?,?,?,?,?,?,?,?,?,?)",
		req.Phone, req.Name, req.Phone, pwd, 2, req.Status, req.DepartmentId, currentCompanyId(ctx), 0, req.Remark)
	if err != nil {
		return nil, gerror.Wrap(err, "新增人员失败")
	}
	id, _ := result.LastInsertId()
	return &v1.CreateRes{Id: uint(id)}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (*v1.UpdateRes, error) {
	count, _ := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE phone=? AND id!=? AND company_id=? AND deleted_at IS NULL", req.Phone, req.Id, currentCompanyId(ctx))
	if count.Int() > 0 {
		return nil, gerror.New("该手机号已被其他人员使用")
	}
	_, err := g.DB().Exec(ctx, "UPDATE site_user SET nickname=?,phone=?,username=?,status=?,department_id=?,remark=? WHERE id=? AND company_id=?",
		req.Name, req.Phone, req.Phone, req.Status, req.DepartmentId, req.Remark, req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "更新人员失败")
	}
	return &v1.UpdateRes{}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	_, err := g.DB().Exec(ctx, "UPDATE site_user SET deleted_at=NOW() WHERE id=? AND company_id=? AND deleted_at IS NULL", req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "删除人员失败")
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (*v1.BatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return nil, gerror.New("请选择要删除的人员")
	}
	placeholders := strings.Repeat("?,", len(req.Ids))
	placeholders = placeholders[:len(placeholders)-1]
	args := make([]any, len(req.Ids))
	for i, id := range req.Ids {
		args[i] = id
	}
	args = append(args, currentCompanyId(ctx))
	_, err := g.DB().Exec(ctx, "UPDATE site_user SET deleted_at=NOW() WHERE id IN("+placeholders+") AND company_id=? AND deleted_at IS NULL", args...)
	if err != nil {
		return nil, gerror.Wrap(err, "批量删除人员失败")
	}
	return &v1.BatchDeleteRes{}, nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func rowToItem(row gdb.Record) *v1.StaffItem {
	role := row["role"].Uint()
	position := ""
	if role == 1 {
		position = "管理员"
	} else if role == 2 {
		position = "员工"
	}
	return &v1.StaffItem{
		Id: row["id"].Uint(), Name: row["nickname"].String(), Phone: row["phone"].String(),
		Gender: "", Position: position,
		DepartmentId: row["department_id"].Uint(), DepartmentName: row["department_name"].String(),
		JoinDate: safeFormatTime(row["created_at"], "2006-01-02"),
		Status:   row["status"].Int(), Remark: row["remark"].String(),
		CreatedAt: safeFormatTime(row["created_at"], "2006-01-02 15:04:05"),
		UpdatedAt: safeFormatTime(row["updated_at"], "2006-01-02 15:04:05"),
	}
}

func safeFormatTime(v gdb.Value, layout string) string {
	if v.IsNil() || v.IsEmpty() {
		return ""
	}
	t := v.Time()
	if t.IsZero() {
		return ""
	}
	return t.Format(layout)
}
