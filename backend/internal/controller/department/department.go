package department

import (
	"context"
	"strings"

	departmentApi "engineServ/api/department"
	v1 "engineServ/api/department/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() departmentApi.IDepartmentV1 { return &Controller{} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	base := `FROM sys_department d WHERE d.company_id=?`
	args := []any{currentCompanyId(ctx)}
	var whereClause string
	if req.Keyword != "" {
		whereClause = ` AND (d.name LIKE ? OR d.code LIKE ?)`
		kw := "%" + req.Keyword + "%"
		args = append(args, kw, kw)
	}
	if req.Status != nil {
		if whereClause == "" {
			whereClause = ` AND d.status=?`
		} else {
			whereClause += ` AND d.status=?`
		}
		args = append(args, *req.Status)
	}

	total, _ := g.DB().GetValue(ctx, "SELECT COUNT(*)"+base+whereClause, args...)
	offset := (req.Current - 1) * req.Size
	rows, err := g.DB().GetAll(ctx, "SELECT d.id,d.name,d.code,d.status,d.sort_order,(SELECT COUNT(*) FROM site_user su WHERE su.department_id=d.id AND su.company_id=d.company_id AND su.deleted_at IS NULL AND su.status=1) member_count "+base+whereClause+" ORDER BY d.sort_order,d.id LIMIT ?,?", append(args, offset, req.Size)...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询部门列表失败")
	}

	records := make([]*v1.DepartmentItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, rowToItem(row))
	}
	return &v1.ListRes{Current: req.Current, Size: req.Size, Total: total.Int(), Records: records}, nil
}

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (*v1.DetailRes, error) {
	row, _ := g.DB().GetOne(ctx, "SELECT d.id,d.name,d.code,d.status,d.sort_order,(SELECT COUNT(*) FROM site_user su WHERE su.department_id=d.id AND su.company_id=d.company_id AND su.deleted_at IS NULL AND su.status=1) member_count FROM sys_department d WHERE d.id=? AND d.company_id=?", req.Id, currentCompanyId(ctx))
	if row.IsEmpty() {
		return nil, gerror.New("部门不存在")
	}
	return rowToItem(row), nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	result, err := g.DB().Exec(ctx, "INSERT INTO sys_department(name,code,status,sort_order,company_id) VALUES(?,?,?,?,?)",
		req.Name, req.Code, req.Status, req.SortOrder, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "新增部门失败（编码可能已存在）")
	}
	id, _ := result.LastInsertId()
	return &v1.CreateRes{Id: uint(id)}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (*v1.UpdateRes, error) {
	_, err := g.DB().Exec(ctx, "UPDATE sys_department SET name=?,code=?,status=?,sort_order=? WHERE id=? AND company_id=?",
		req.Name, req.Code, req.Status, req.SortOrder, req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "更新部门失败")
	}
	return &v1.UpdateRes{}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	count, _ := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE department_id=? AND company_id=? AND deleted_at IS NULL", req.Id, currentCompanyId(ctx))
	if count.Int() > 0 {
		return nil, gerror.New("该部门下还有人员，请先将人员移出")
	}
	_, err := g.DB().Exec(ctx, "DELETE FROM sys_department WHERE id=? AND company_id=?", req.Id, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "删除部门失败")
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (*v1.BatchDeleteRes, error) {
	if len(req.Ids) == 0 {
		return nil, gerror.New("请选择要删除的部门")
	}
	pl := strings.Repeat("?,", len(req.Ids))
	pl = pl[:len(pl)-1]
	args := make([]any, len(req.Ids))
	for i, id := range req.Ids {
		args[i] = id
	}
	args = append(args, currentCompanyId(ctx))
	count, _ := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM site_user WHERE department_id IN("+pl+") AND company_id=? AND deleted_at IS NULL", args...)
	if count.Int() > 0 {
		return nil, gerror.New("部分部门下还有人员，请先将人员移出")
	}
	_, err := g.DB().Exec(ctx, "DELETE FROM sys_department WHERE id IN("+pl+") AND company_id=?", args...)
	if err != nil {
		return nil, gerror.Wrap(err, "批量删除部门失败")
	}
	return &v1.BatchDeleteRes{}, nil
}

func (c *Controller) Options(ctx context.Context, _ *v1.OptionsReq) (*v1.OptionsRes, error) {
	rows, _ := g.DB().GetAll(ctx, "SELECT id,name FROM sys_department WHERE status=1 AND company_id=? ORDER BY sort_order,id", currentCompanyId(ctx))
	items := make([]v1.Option, 0, len(rows))
	for _, row := range rows {
		items = append(items, v1.Option{Value: row["id"].Uint(), Label: row["name"].String()})
	}
	return &v1.OptionsRes{Items: items}, nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func rowToItem(row gdb.Record) *v1.DepartmentItem {
	return &v1.DepartmentItem{
		Id:          row["id"].Uint(),
		Name:        row["name"].String(),
		Code:        row["code"].String(),
		LeaderId:    nil,
		LeaderName:  "",
		MemberCount: row["member_count"].Int(),
		Status:      row["status"].Int(),
		Remark:      "",
		SortOrder:   row["sort_order"].Int(),
		CreatedAt:   "",
		UpdatedAt:   "",
	}
}
