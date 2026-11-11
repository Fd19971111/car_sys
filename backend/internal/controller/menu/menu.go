package menu

import (
	"context"

	menuApi "engineServ/api/menu"
	v1 "engineServ/api/menu/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() menuApi.IMenuV1 { return &Controller{} }

func (c *Controller) List(ctx context.Context, _ *v1.ListReq) (*v1.ListRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	rows, err := g.DB().GetAll(ctx, "SELECT id,parent_id,title,route_name,route_path,component,icon,sort_order,enabled,hide_in_menu FROM sys_menu ORDER BY sort_order,id")
	if err != nil {
		return nil, gerror.Wrap(err, "查询菜单失败")
	}
	items := make([]*v1.MenuItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, menuFromRow(row))
	}
	return &v1.ListRes{Items: items}, nil
}

func (c *Controller) Departments(ctx context.Context, _ *v1.DepartmentsReq) (*v1.DepartmentsRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	rows, err := g.DB().GetAll(ctx, "SELECT id,name FROM sys_department WHERE status=1 AND company_id=? ORDER BY sort_order,id", currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询部门失败")
	}
	items := make([]*v1.Department, 0, len(rows))
	for _, row := range rows {
		items = append(items, &v1.Department{Id: row["id"].Uint(), Name: row["name"].String()})
	}
	return &v1.DepartmentsRes{Items: items}, nil
}

func (c *Controller) CreateDepartment(ctx context.Context, req *v1.CreateDepartmentReq) (*v1.CreateDepartmentRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	result, err := g.DB().Exec(ctx, "INSERT INTO sys_department(name,code,company_id) VALUES(?,?,?)", req.Name, req.Code, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "新增部门失败，部门编码可能已存在")
	}
	id, err := result.LastInsertId()
	return &v1.CreateDepartmentRes{Id: uint(id)}, err
}

func (c *Controller) DepartmentMenus(ctx context.Context, req *v1.DepartmentMenusReq) (*v1.DepartmentMenusRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	rows, err := g.DB().GetAll(ctx, "SELECT menu_id FROM sys_department_menu WHERE department_id=? AND company_id=?", req.DepartmentId, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "查询部门菜单权限失败")
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row["menu_id"].Uint())
	}
	return &v1.DepartmentMenusRes{MenuIds: ids}, nil
}

func (c *Controller) Save(ctx context.Context, req *v1.SaveReq) (*v1.SaveRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	if !allowedFrontendRoute(req.RouteName, req.RoutePath, req.Component) {
		return nil, gerror.New("菜单路由必须对应已编译的前端页面")
	}
	if req.Id > 0 && req.ParentId == req.Id {
		return nil, gerror.New("菜单不能设置自身为父级")
	}
	if req.ParentId > 0 {
		parentCount, checkErr := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM sys_menu WHERE id=? AND component='layout.base' AND parent_id=0", req.ParentId)
		if checkErr != nil {
			return nil, gerror.Wrap(checkErr, "校验父级菜单失败")
		}
		if parentCount.Int() == 0 {
			return nil, gerror.New("父级菜单必须是有效的一级分组菜单")
		}
	}
	if req.Id == 0 {
		result, err := g.DB().Exec(ctx, "INSERT INTO sys_menu(parent_id,title,route_name,route_path,component,icon,sort_order,enabled,hide_in_menu) VALUES(?,?,?,?,?,?,?,?,?)",
			req.ParentId, req.Title, req.RouteName, req.RoutePath, req.Component, req.Icon, req.SortOrder, req.Enabled, req.HideInMenu)
		if err != nil {
			return nil, gerror.Wrap(err, "新增菜单失败")
		}
		id, err := result.LastInsertId()
		return &v1.SaveRes{Id: uint(id)}, err
	}
	_, err := g.DB().Exec(ctx, "UPDATE sys_menu SET parent_id=?,title=?,route_name=?,route_path=?,component=?,icon=?,sort_order=?,enabled=?,hide_in_menu=? WHERE id=?",
		req.ParentId, req.Title, req.RouteName, req.RoutePath, req.Component, req.Icon, req.SortOrder, req.Enabled, req.HideInMenu, req.Id)
	if err != nil {
		return nil, gerror.Wrap(err, "更新菜单失败")
	}
	return &v1.SaveRes{Id: req.Id}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM sys_menu WHERE parent_id=?", req.Id)
	if err != nil {
		return nil, gerror.Wrap(err, "检查菜单子项失败")
	}
	if count.Int() > 0 {
		return nil, gerror.New("请先删除该菜单下的子菜单")
	}
	if _, err := g.DB().Exec(ctx, "DELETE FROM sys_department_menu WHERE menu_id=?", req.Id); err != nil {
		return nil, gerror.Wrap(err, "清理部门菜单权限失败")
	}
	if _, err := g.DB().Exec(ctx, "DELETE FROM sys_menu WHERE id=?", req.Id); err != nil {
		return nil, gerror.Wrap(err, "删除菜单失败")
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) SetDepartmentMenus(ctx context.Context, req *v1.SetDepartmentMenusReq) (*v1.SetDepartmentMenusRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	departmentCount, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM sys_department WHERE id=? AND company_id=?", req.DepartmentId, currentCompanyId(ctx))
	if err != nil {
		return nil, gerror.Wrap(err, "校验部门失败")
	}
	if departmentCount.Int() == 0 {
		return nil, gerror.New("部门不存在")
	}
	if _, err := g.DB().Exec(ctx, "DELETE FROM sys_department_menu WHERE department_id=? AND company_id=?", req.DepartmentId, currentCompanyId(ctx)); err != nil {
		return nil, gerror.Wrap(err, "清除部门菜单权限失败")
	}
	for _, menuId := range req.MenuIds {
		if _, err := g.DB().Exec(ctx, "INSERT IGNORE INTO sys_department_menu(department_id,company_id,menu_id) VALUES(?,?,?)", req.DepartmentId, currentCompanyId(ctx), menuId); err != nil {
			return nil, gerror.Wrap(err, "保存部门菜单权限失败")
		}
	}
	return &v1.SetDepartmentMenusRes{}, nil
}

func (c *Controller) GetUserRoutes(ctx context.Context, _ *v1.GetUserRoutesReq) (*v1.GetUserRoutesRes, error) {
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))
	role := gconv.Uint(ctx.Value(consts.CtxRoleKey))
	var rows gdb.Result
	var err error
	if role == 1 {
		rows, err = g.DB().GetAll(ctx, "SELECT m.id,m.parent_id,m.title,m.route_name,m.route_path,m.component,m.icon,m.sort_order,m.hide_in_menu FROM sys_menu m WHERE m.enabled=1 ORDER BY m.sort_order,m.id")
	} else {
		rows, err = g.DB().GetAll(ctx, "SELECT DISTINCT m.id,m.parent_id,m.title,m.route_name,m.route_path,m.component,m.icon,m.sort_order,m.hide_in_menu FROM sys_menu m JOIN site_user u ON u.id=? LEFT JOIN sys_department_menu dm ON dm.department_id=u.department_id AND dm.company_id=u.company_id LEFT JOIN sys_menu granted ON granted.id=dm.menu_id WHERE u.company_id=? AND m.enabled=1 AND (m.id=granted.id OR m.id=granted.parent_id OR m.parent_id=granted.id) ORDER BY m.sort_order,m.id", userId, currentCompanyId(ctx))
	}
	if err != nil {
		return nil, gerror.Wrap(err, "获取用户菜单失败")
	}
	byId := make(map[uint]*v1.MenuRoute, len(rows))
	parents := make(map[uint]uint, len(rows))
	for _, row := range rows {
		id := row["id"].Uint()
		byId[id] = &v1.MenuRoute{
			Name: row["route_name"].String(), Path: row["route_path"].String(), Component: row["component"].String(),
			Meta: v1.RouteMeta{Title: row["title"].String(), Icon: row["icon"].String(), Order: row["sort_order"].Int(), HideInMenu: row["hide_in_menu"].Bool()},
		}
		parents[id] = row["parent_id"].Uint()
	}
	roots := make([]*v1.MenuRoute, 0)
	for _, row := range rows {
		id, parentId := row["id"].Uint(), parents[row["id"].Uint()]
		if parentId == 0 || byId[parentId] == nil {
			roots = append(roots, byId[id])
		} else {
			parent := byId[parentId]
			parent.Children = append(parent.Children, byId[id])
		}
	}
	if role == 1 {
		setupRoute := &v1.MenuRoute{
			Name:      "finance_setup",
			Path:      "/finance/setup",
			Component: "view.finance_setup",
			Meta:      v1.RouteMeta{Title: "finance_setup", HideInMenu: true},
		}
		roots = append(roots, setupRoute)
	}
	return &v1.GetUserRoutesRes{UserRoute: v1.UserRoute{Routes: roots, Home: "home"}}, nil
}

func requireSuper(ctx context.Context) error {
	if gconv.Uint(ctx.Value(consts.CtxRoleKey)) != 1 {
		return gerror.New("仅公司超级管理员可管理菜单和部门权限")
	}
	return nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func allowedFrontendRoute(name, path, component string) bool {
	allowed := map[string][2]string{
		"home":                  {"/home", "layout.base$view.home"},
		"project":               {"/project", "layout.base"},
		"project_project-list":  {"/project/project-list", "view.project_project-list"},
		"project_data":          {"/project/data", "view.project_data"},
		"project_quote":         {"/project/quote", "view.project_quote"},
		"finance":               {"/finance", "layout.base"},
		"finance_payment":       {"/finance/payment", "view.finance_payment"},
		"finance_payout":        {"/finance/payout", "view.finance_payout"},
		"finance_voucher":       {"/finance/voucher", "view.finance_voucher"},
		"finance_basic":         {"/finance/basic", "view.finance_basic"},
		"report":                {"/report", "layout.base"},
		"report_overview":       {"/report/overview", "view.report_overview"},
		"report_flows":          {"/report/flows", "view.report_flows"},
		"report_balance":        {"/report/balance", "view.report_balance"},
		"report_profit":         {"/report/profit", "view.report_profit"},
		"inventory":             {"/inventory", "layout.base"},
		"inventory_materials":   {"/inventory/materials", "view.inventory_materials"},
		"inventory_usages":      {"/inventory/usages", "view.inventory_usages"},
		"inventory_procurement": {"/inventory/procurement", "view.inventory_procurement"},
		"inventory_stock":       {"/inventory/stock", "view.inventory_stock"},
		"inventory_warehouse":   {"/inventory/warehouse", "view.inventory_warehouse"},
		"enterprise":            {"/enterprise", "layout.base"},
		"enterprise_staff":      {"/enterprise/staff", "view.enterprise_staff"},
		"enterprise_department": {"/enterprise/department", "view.enterprise_department"},
		"enterprise_user":       {"/enterprise/user", "view.enterprise_user"},
		"settings":              {"/settings", "layout.base"},
		"settings_customer":     {"/settings/customer", "view.settings_customer"},
		"settings_supplier":     {"/settings/supplier", "view.settings_supplier"},
		"settings_menu":         {"/settings/menu", "view.settings_menu"},
		"settings_print":        {"/settings/print", "view.settings_print"},
	}
	match, ok := allowed[name]
	return ok && match == [2]string{path, component}
}

func menuFromRow(row gdb.Record) *v1.MenuItem {
	return &v1.MenuItem{
		Id: row["id"].Uint(), ParentId: row["parent_id"].Uint(), Title: row["title"].String(),
		RouteName: row["route_name"].String(), RoutePath: row["route_path"].String(), Component: row["component"].String(),
		Icon: row["icon"].String(), SortOrder: row["sort_order"].Int(), Enabled: row["enabled"].Bool(), HideInMenu: row["hide_in_menu"].Bool(),
	}
}

func (c *Controller) Users(ctx context.Context, _ *v1.UsersReq) (*v1.UsersRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	rows, err := g.DB().GetAll(ctx, "SELECT id,username,nickname,role,department_id FROM site_user WHERE status=1 ORDER BY role,id")
	if err != nil {
		return nil, gerror.Wrap(err, "查询后台人员失败")
	}
	items := make([]*v1.DepartmentUser, 0, len(rows))
	for _, row := range rows {
		items = append(items, &v1.DepartmentUser{Id: row["id"].Uint(), Username: row["username"].String(), Nickname: row["nickname"].String(), Role: row["role"].Uint(), DepartmentId: row["department_id"].Uint()})
	}
	return &v1.UsersRes{Items: items}, nil
}

func (c *Controller) SetUserDepartment(ctx context.Context, req *v1.SetUserDepartmentReq) (*v1.SetUserDepartmentRes, error) {
	if err := requireSuper(ctx); err != nil {
		return nil, err
	}
	roleValue, err := g.DB().GetValue(ctx, "SELECT role FROM site_user WHERE id=? AND status=1", req.UserId)
	if err != nil {
		return nil, gerror.Wrap(err, "查询人员失败")
	}
	role := roleValue.Uint()
	if role == 1 && req.DepartmentId != 0 {
		return nil, gerror.New("公司超级管理员不能挂靠到部门")
	}
	if role != 1 && req.DepartmentId > 0 {
		count, checkErr := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM sys_department WHERE id=? AND status=1", req.DepartmentId)
		if checkErr != nil {
			return nil, gerror.Wrap(checkErr, "校验部门失败")
		}
		if count.Int() == 0 {
			return nil, gerror.New("部门不存在或已停用")
		}
	}
	if _, err := g.DB().Exec(ctx, "UPDATE site_user SET department_id=? WHERE id=?", req.DepartmentId, req.UserId); err != nil {
		return nil, gerror.Wrap(err, "设置人员所属部门失败")
	}
	return &v1.SetUserDepartmentRes{}, nil
}
