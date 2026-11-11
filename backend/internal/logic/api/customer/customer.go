package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"engineServ/internal/consts"
	defaultDao "engineServ/model/default/dao"
	defaultDo "engineServ/model/default/model/do"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// 客户状态：与 customer.status 注释保持一致（customer 表结构固定，不新增/修改列）
const (
	statusFollowing = 1 // 跟进中
	statusDeal      = 2 // 已成交
	statusLost      = 3 // 已流失
)

// ProjectInfo 客户详情中的装修项目信息（project 子表 + site_user 负责人 + project_log 最近进度）
type ProjectInfo struct {
	Project          *defaultEntity.Project
	ManagerName      string
	LatestProgress   string
	LatestProgressAt *gtime.Time
}

type Customer struct{}

func New() *Customer {
	return &Customer{}
}

// List 客户分页列表：关键词（姓名/手机号/房屋地址）+ 状态 + 来源过滤，按创建时间倒序。
// 同时返回每个客户名下未删除项目的数量映射。
func (c *Customer) List(
	ctx context.Context,
	keyword string,
	source string,
	status *int,
	current int,
	size int,
) (list []*defaultEntity.Customer, projectCount map[uint]int, total int, err error) {
	cols := defaultDao.Customer.Columns()
	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		// 用括号包裹 OR 条件，避免与后面的 AND 过滤出现优先级问题
		kwCondition := fmt.Sprintf("(%s LIKE ? OR %s LIKE ? OR %s LIKE ?)", cols.Name, cols.Phone, cols.Address)
		model = model.Where(kwCondition, like, like, like)
	}

	// 注意：前端未选择状态时会带 `status=` 空串请求，GoFrame 会把它绑定成「指向 0 的非空指针」。
	// 若只判断 status != nil，SQL 会变成 WHERE status = 0，而客户状态只有 1/2/3，结果必然为空。
	// 因此这里统一用 > 0 作为「有筛选」的判据。
	if status != nil && *status > 0 {
		model = model.Where(cols.Status, *status)
	}

	source = strings.TrimSpace(source)
	if source != "" {
		model = model.Where(cols.Source, source)
	}

	err = model.
		Page(current, size).
		OrderDesc(cols.CreatedAt).
		OrderDesc(cols.Id).
		ScanAndCount(&list, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询客户分页列表失败: %v", err)
		return nil, nil, 0, gerror.New("查询客户列表失败，请稍后重试")
	}

	projectCount = c.batchProjectCount(ctx, list)

	return list, projectCount, total, nil
}

// Stats 客户概览统计：各状态客户数量 + 已成交客户名下项目合同总金额。
func (c *Customer) Stats(ctx context.Context) (total, following, deal, lost int, dealAmount float64, err error) {
	cols := defaultDao.Customer.Columns()

	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	total, err = model.Count()
	if err != nil {
		glog.Errorf(ctx, "统计客户总数失败: %v", err)
		return 0, 0, 0, 0, 0, gerror.New("获取客户统计失败，请稍后重试")
	}

	countByStatus := func(status int) int {
		statusModel := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(cols.Status, status)
		if !isSuperUser(ctx) {
			statusModel = statusModel.Where("owner_user_id", currentUserId(ctx))
		}
		count, countErr := statusModel.Count()
		if countErr != nil {
			glog.Errorf(ctx, "统计客户状态(%d)失败: %v", status, countErr)
			err = gerror.New("获取客户统计失败，请稍后重试")
		}

		return count
	}

	following = countByStatus(statusFollowing)
	if err != nil {
		return
	}

	deal = countByStatus(statusDeal)
	if err != nil {
		return
	}

	lost = countByStatus(statusLost)
	if err != nil {
		return
	}

	// 已成交客户（customer.status=2）名下未删除项目的合同金额合计
	dealCustomerSubQuery := defaultDao.Customer.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(cols.Id).
		Where(cols.Status, statusDeal)
	if !isSuperUser(ctx) {
		dealCustomerSubQuery = dealCustomerSubQuery.Where("owner_user_id", currentUserId(ctx))
	}

	dealAmount, err = defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn(defaultDao.Project.Columns().CustomerId, dealCustomerSubQuery).
		Sum(defaultDao.Project.Columns().Amount)
	if err != nil {
		glog.Errorf(ctx, "统计成交合同金额失败: %v", err)
		return 0, 0, 0, 0, 0, gerror.New("获取客户统计失败，请稍后重试")
	}

	return total, following, deal, lost, dealAmount, nil
}

// Detail 客户详情：客户本体 + 名下装修项目（含负责人姓名、最近一条施工进度）。
func (c *Customer) Detail(ctx context.Context, id uint) (customer *defaultEntity.Customer, projects []*ProjectInfo, err error) {
	customer, err = c.getCustomer(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	var projectList []*defaultEntity.Project
	err = defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.Project.Columns().CustomerId, id).
		OrderAsc(defaultDao.Project.Columns().Id).
		Scan(&projectList)
	if err != nil {
		glog.Errorf(ctx, "查询客户项目失败, customerId=%d: %v", id, err)
		return nil, nil, gerror.New("获取客户详情失败，请稍后重试")
	}

	projects = c.attachProjectExtras(ctx, projectList)

	return customer, projects, nil
}

// Create 新增客户，返回新客户ID。
func (c *Customer) Create(ctx context.Context, data *defaultDo.Customer) (id uint, err error) {
	insertData := g.Map{
		"company_id": currentCompanyId(ctx), "name": data.Name, "phone": data.Phone, "address": data.Address,
		"source": data.Source, "status": data.Status, "remark": data.Remark,
	}
	if !isSuperUser(ctx) {
		insertData["owner_user_id"] = currentUserId(ctx)
	}
	result, err := defaultDao.Customer.Ctx(ctx).Data(insertData).Insert()
	if err != nil {
		glog.Errorf(ctx, "新增客户失败: %v", err)
		return 0, gerror.New("新增客户失败，请稍后重试")
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		glog.Errorf(ctx, "读取新增客户ID失败: %v", err)
		return 0, gerror.New("新增客户失败，请稍后重试")
	}

	return uint(lastInsertId), nil
}

// Update 更新客户基础信息（仅允许修改 customer 表既有列）。
func (c *Customer) Update(ctx context.Context, id uint, data *defaultDo.Customer) error {
	if _, err := c.getCustomer(ctx, id); err != nil {
		return err
	}

	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Customer.Columns().Id, id)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	_, err := model.
		Data(data).
		Update()
	if err != nil {
		glog.Errorf(ctx, "更新客户失败, id=%v: %v", data.Id, err)
		return gerror.New("更新客户失败，请稍后重试")
	}

	return nil
}

// Delete 删除单个客户（ORM 软删除，写入 deleted_at）。
func (c *Customer) Delete(ctx context.Context, id uint) error {
	if _, err := c.getCustomer(ctx, id); err != nil {
		return err
	}

	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Customer.Columns().Id, id)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	_, err := model.Delete()
	if err != nil {
		glog.Errorf(ctx, "删除客户失败, id=%d: %v", id, err)
		return gerror.New("删除客户失败，请稍后重试")
	}

	return nil
}

// BatchDelete 批量删除客户（ORM 软删除）。
func (c *Customer) BatchDelete(ctx context.Context, ids []uint) error {
	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).WhereIn(defaultDao.Customer.Columns().Id, ids)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	_, err := model.Delete()
	if err != nil {
		glog.Errorf(ctx, "批量删除客户失败, ids=%v: %v", ids, err)
		return gerror.New("删除客户失败，请稍后重试")
	}

	return nil
}

// getCustomer 按主键查询未删除客户，不存在时返回业务错误。
func (c *Customer) getCustomer(ctx context.Context, id uint) (*defaultEntity.Customer, error) {
	var customer defaultEntity.Customer
	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Customer.Columns().Id, id)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	err := model.Scan(&customer)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询客户失败, id=%d: %v", id, err)
		return nil, gerror.New("获取客户信息失败，请稍后重试")
	}

	if customer.Id == 0 {
		return nil, gerror.New("客户不存在或已被删除")
	}

	return &customer, nil
}

func currentUserId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxUserIdKey))
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func isSuperUser(ctx context.Context) bool {
	role := gconv.Uint(ctx.Value(consts.CtxRoleKey))
	return role == 1 || role == 3
}

// batchProjectCount 批量统计客户名下项目数量，无项目的客户不出现于映射中。
func (c *Customer) batchProjectCount(ctx context.Context, list []*defaultEntity.Customer) map[uint]int {
	result := make(map[uint]int)

	if len(list) == 0 {
		return result
	}

	customerIds := make([]uint, 0, len(list))
	for _, item := range list {
		customerIds = append(customerIds, item.Id)
	}

	rows, err := defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields("customer_id, COUNT(*) AS total").
		WhereIn(defaultDao.Project.Columns().CustomerId, customerIds).
		Group(defaultDao.Project.Columns().CustomerId).
		All()
	if err != nil {
		// 项目数量仅用于展示，统计失败不阻断列表
		glog.Errorf(ctx, "统计客户项目数量失败: %v", err)
		return result
	}

	for _, row := range rows {
		result[row["customer_id"].Uint()] = row["total"].Int()
	}

	return result
}

// attachProjectExtras 为项目列表补充负责人姓名（site_user）与最近一条进度（project_log）。
func (c *Customer) attachProjectExtras(ctx context.Context, projects []*defaultEntity.Project) []*ProjectInfo {
	result := make([]*ProjectInfo, 0, len(projects))
	if len(projects) == 0 {
		return result
	}

	projectIds := make([]uint, 0, len(projects))
	managerIds := make([]uint, 0, len(projects))
	managerIdSet := make(map[uint]struct{})

	for _, project := range projects {
		projectIds = append(projectIds, project.Id)

		if project.ManagerId > 0 {
			if _, ok := managerIdSet[project.ManagerId]; !ok {
				managerIdSet[project.ManagerId] = struct{}{}
				managerIds = append(managerIds, project.ManagerId)
			}
		}
	}

	// 负责人姓名映射
	managerNames := make(map[uint]string)
	if len(managerIds) > 0 {
		var users []*defaultEntity.SiteUser
		if err := defaultDao.SiteUser.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).WhereIn(defaultDao.SiteUser.Columns().Id, managerIds).Scan(&users); err != nil {
			glog.Errorf(ctx, "查询项目负责人失败, managerIds=%v: %v", managerIds, err)
		} else {
			for _, user := range users {
				managerNames[user.Id] = firstNonEmpty(user.Nickname, user.Username)
			}
		}
	}

	// 最近进度：按业务日期倒序批量取出，遍历时每个项目只保留第一条
	latestLogs := make(map[uint]*defaultEntity.ProjectLog)
	logs, err := defaultDao.ProjectLog.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn(defaultDao.ProjectLog.Columns().ProjectId, projectIds).
		OrderDesc(defaultDao.ProjectLog.Columns().LogDate).
		OrderDesc(defaultDao.ProjectLog.Columns().Id).
		All()
	if err != nil {
		// 进度仅用于展示，查询失败不阻断详情
		glog.Errorf(ctx, "查询项目最近进度失败, projectIds=%v: %v", projectIds, err)
	} else {
		for _, row := range logs {
			var logRow defaultEntity.ProjectLog
			if err = row.Struct(&logRow); err != nil {
				continue
			}

			if _, ok := latestLogs[logRow.ProjectId]; !ok {
				latestLogs[logRow.ProjectId] = &logRow
			}
		}
	}

	for _, project := range projects {
		info := &ProjectInfo{
			Project:     project,
			ManagerName: managerNames[project.ManagerId],
		}

		if logRow, ok := latestLogs[project.Id]; ok {
			info.LatestProgress = logRow.Content
			// 展示业务日期（进度发生日），历史补录的进度不会错显示为当天
			info.LatestProgressAt = logRow.LogDate
			if info.LatestProgressAt == nil {
				info.LatestProgressAt = logRow.CreatedAt
			}
		}

		result = append(result, info)
	}

	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}
