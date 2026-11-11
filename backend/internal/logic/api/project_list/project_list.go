package project_list

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"engineServ/internal/consts"
	progressLogic "engineServ/internal/logic/api/progress"
	defaultDao "engineServ/model/default/dao"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// 项目状态：与 project.status 注释保持一致（project 表结构固定，不新增/修改列）
const (
	statusPending  = 1 // 待开工
	statusBuilding = 2 // 施工中
	statusFinished = 3 // 已完工
	statusSettled  = 4 // 已结算
)

// optionLimit 表单下拉选项的最大返回条数，避免一次性拉全表
const optionLimit = 500

// ProjectOption 项目下拉选项（供报价 / 收款模块使用）
type ProjectOption struct {
	Id     uint
	Label  string
	Amount float64
	Status int
}

// ProjectInfo 项目本体 + 列表 / 详情展示所需的关联信息
type ProjectInfo struct {
	Project          *defaultEntity.Project
	CustomerName     string
	ManagerName      string
	LatestProgress   string
	LatestProgressAt *gtime.Time
}

// LogInfo 施工进度记录 + 记录人姓名
type LogInfo struct {
	Log      *defaultEntity.ProjectLog
	UserName string
}

// Form 项目新增 / 编辑入参（与 project 表既有列一一对应）
type Form struct {
	CustomerId uint
	Name       string
	Address    string
	ManagerId  uint
	Status     int
	Amount     float64
	StartDate  string
	EndDate    string
	Remark     string
}

type Project struct {
	progress *progressLogic.Progress
}

func New() *Project {
	return &Project{progress: progressLogic.New()}
}

// List 项目分页列表：关键词（项目名称 / 施工地址 / 客户姓名）+ 状态 + 所属客户过滤，按创建时间倒序。
// 返回结果已补充客户姓名、负责人姓名与最近一条施工进度。
func (p *Project) List(
	ctx context.Context,
	keyword string,
	status *int,
	customerId *uint,
	current int,
	size int,
) (list []*ProjectInfo, total int, err error) {
	cols := defaultDao.Project.Columns()
	model := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	if !isSuperUser(ctx) {
		userId := currentUserId(ctx)
		model = model.Where("(manager_id = ? OR customer_id IN (SELECT id FROM customer WHERE owner_user_id = ?))", userId, userId)
	}

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"

		// 关键词命中的客户（姓名匹配）→ 这些客户名下的项目同样算命中
		var matchedCustomerIds []uint
		customerRows, customerErr := defaultDao.Customer.Ctx(ctx).
			Where("company_id", currentCompanyId(ctx)).
			Fields(defaultDao.Customer.Columns().Id).
			WhereLike(defaultDao.Customer.Columns().Name, like).
			All()
		if customerErr != nil {
			glog.Errorf(ctx, "查询关键词命中的客户失败: %v", customerErr)
			return nil, 0, gerror.New("查询项目列表失败，请稍后重试")
		}

		for _, row := range customerRows {
			matchedCustomerIds = append(matchedCustomerIds, row[defaultDao.Customer.Columns().Id].Uint())
		}

		// 命中的项目：项目名称 / 施工地址 直接匹配，或所属客户姓名匹配
		matchedModel := defaultDao.Project.Ctx(ctx).
			Where("company_id", currentCompanyId(ctx)).
			Fields(cols.Id).
			WhereLike(cols.Name, like).
			WhereOrLike(cols.Address, like)

		if len(matchedCustomerIds) > 0 {
			matchedModel = matchedModel.WhereOrIn(cols.CustomerId, matchedCustomerIds)
		}

		matchedRows, matchedErr := matchedModel.All()
		if matchedErr != nil {
			glog.Errorf(ctx, "查询关键词命中的项目失败: %v", matchedErr)
			return nil, 0, gerror.New("查询项目列表失败，请稍后重试")
		}

		matchedProjectIds := make([]uint, 0, len(matchedRows))
		for _, row := range matchedRows {
			matchedProjectIds = append(matchedProjectIds, row[cols.Id].Uint())
		}

		if len(matchedProjectIds) == 0 {
			return []*ProjectInfo{}, 0, nil
		}

		model = model.WhereIn(cols.Id, matchedProjectIds)
	}

	// 前端未选择状态时会带 `status=` 空串请求，GoFrame 会把它绑定成「指向 0 的非空指针」，
	// 因此这里用 > 0 作为「有筛选」的判据（项目状态只有 1/2/3/4）。
	if status != nil && *status > 0 {
		model = model.Where(cols.Status, *status)
	}

	if customerId != nil && *customerId > 0 {
		model = model.Where(cols.CustomerId, *customerId)
	}

	var projects []*defaultEntity.Project
	err = model.
		Page(current, size).
		OrderDesc(cols.CreatedAt).
		OrderDesc(cols.Id).
		ScanAndCount(&projects, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询项目分页列表失败: %v", err)
		return nil, 0, gerror.New("查询项目列表失败，请稍后重试")
	}

	return p.attachExtras(ctx, projects), total, nil
}

// Stats 项目概览统计：总数 + 各状态数量 + 合同总金额。
func (p *Project) Stats(ctx context.Context) (total, pending, building, finished, settled int, totalAmount float64, err error) {
	cols := defaultDao.Project.Columns()

	model := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	if !isSuperUser(ctx) {
		userId := currentUserId(ctx)
		model = model.Where("(manager_id = ? OR customer_id IN (SELECT id FROM customer WHERE owner_user_id = ?))", userId, userId)
	}
	rows, queryErr := model.
		Fields(cols.Status, "COUNT(*) AS total").
		Group(cols.Status).
		All()
	if queryErr != nil {
		glog.Errorf(ctx, "统计项目状态分布失败: %v", queryErr)
		return 0, 0, 0, 0, 0, 0, gerror.New("获取项目统计失败，请稍后重试")
	}

	for _, row := range rows {
		count := row["total"].Int()
		total += count

		switch row[cols.Status].Int() {
		case statusPending:
			pending = count
		case statusBuilding:
			building = count
		case statusFinished:
			finished = count
		case statusSettled:
			settled = count
		}
	}

	totalAmount, queryErr = model.Sum(cols.Amount)
	if queryErr != nil {
		glog.Errorf(ctx, "统计项目合同总金额失败: %v", queryErr)
		return 0, 0, 0, 0, 0, 0, gerror.New("获取项目统计失败，请稍后重试")
	}

	return total, pending, building, finished, settled, totalAmount, nil
}

// Options 项目表单下拉选项：客户列表 + 可指派负责人列表。
func (p *Project) Options(ctx context.Context) (customers []*defaultEntity.Customer, managers []*defaultEntity.SiteUser, err error) {
	customerCols := defaultDao.Customer.Columns()
	userCols := defaultDao.SiteUser.Columns()
	customerModel := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	userModel := defaultDao.SiteUser.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	if !isSuperUser(ctx) {
		userId := currentUserId(ctx)
		customerModel = customerModel.Where("owner_user_id", userId)
		userModel = userModel.Where(userCols.Id, userId)
	}
	err = customerModel.
		Fields(customerCols.Id, customerCols.Name, customerCols.Phone, customerCols.Status).
		OrderDesc(customerCols.Id).
		Limit(optionLimit).
		Scan(&customers)
	if err != nil {
		glog.Errorf(ctx, "查询客户下拉选项失败: %v", err)
		return nil, nil, gerror.New("获取项目表单选项失败，请稍后重试")
	}

	err = userModel.
		Fields(userCols.Id, userCols.Username, userCols.Nickname).
		OrderAsc(userCols.Id).
		Limit(optionLimit).
		Scan(&managers)
	if err != nil {
		glog.Errorf(ctx, "查询负责人下拉选项失败: %v", err)
		return nil, nil, gerror.New("获取项目表单选项失败，请稍后重试")
	}

	return customers, managers, nil
}

// Detail 项目详情：项目本体 + 所属客户 + 施工进度记录（按时间倒序）。
func (p *Project) Detail(ctx context.Context, id uint) (
	project *defaultEntity.Project,
	customer *defaultEntity.Customer,
	logs []*LogInfo,
	err error,
) {
	project, err = p.getProject(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	if project.CustomerId > 0 {
		var customerRow defaultEntity.Customer
		scanErr := defaultDao.Customer.Ctx(ctx).
			Where("company_id", currentCompanyId(ctx)).
			Where(defaultDao.Customer.Columns().Id, project.CustomerId).
			Scan(&customerRow)
		if scanErr != nil {
			// 客户信息仅用于展示，查询失败不阻断详情
			glog.Errorf(ctx, "查询项目所属客户失败, customerId=%d: %v", project.CustomerId, scanErr)
		} else if customerRow.Id > 0 {
			customer = &customerRow
		}
	}

	var logRows []*defaultEntity.ProjectLog
	err = defaultDao.ProjectLog.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectLog.Columns().ProjectId, id).
		OrderDesc(defaultDao.ProjectLog.Columns().LogDate).
		OrderDesc(defaultDao.ProjectLog.Columns().Id).
		Scan(&logRows)
	if err != nil {
		glog.Errorf(ctx, "查询项目进度记录失败, projectId=%d: %v", id, err)
		return nil, nil, nil, gerror.New("获取项目详情失败，请稍后重试")
	}

	return project, customer, p.attachLogUsers(ctx, logRows), nil
}

// Create 新增项目，返回新项目ID。
func (p *Project) Create(ctx context.Context, form *Form) (id uint, err error) {
	if err = p.ensureCustomerExists(ctx, form.CustomerId); err != nil {
		return 0, err
	}
	if !isSuperUser(ctx) {
		form.ManagerId = currentUserId(ctx)
	}

	result, err := defaultDao.Project.Ctx(ctx).Data(form.toData(ctx)).Insert()
	if err != nil {
		glog.Errorf(ctx, "新增项目失败: %v", err)
		return 0, gerror.New("新增项目失败，请稍后重试")
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		glog.Errorf(ctx, "读取新增项目ID失败: %v", err)
		return 0, gerror.New("新增项目失败，请稍后重试")
	}

	return uint(lastInsertId), nil
}

// Update 更新项目（仅允许修改 project 表既有列）。
//
// 这里用 Map 而不是 DO 提交：DO 会自动丢弃零值字段，导致「清空负责人 / 清空工期」无法落库。
func (p *Project) Update(ctx context.Context, id uint, form *Form) error {
	if _, err := p.getProject(ctx, id); err != nil {
		return err
	}

	if err := p.ensureCustomerExists(ctx, form.CustomerId); err != nil {
		return err
	}
	if !isSuperUser(ctx) {
		form.ManagerId = currentUserId(ctx)
	}

	_, err := defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.Project.Columns().Id, id).
		Data(form.toData(ctx)).
		Update()
	if err != nil {
		glog.Errorf(ctx, "更新项目失败, id=%d: %v", id, err)
		return gerror.New("更新项目失败，请稍后重试")
	}

	return nil
}

// Delete 删除单个项目（ORM 软删除，写入 deleted_at）。
func (p *Project) Delete(ctx context.Context, id uint) error {
	if _, err := p.getProject(ctx, id); err != nil {
		return err
	}

	_, err := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Project.Columns().Id, id).Delete()
	if err != nil {
		glog.Errorf(ctx, "删除项目失败, id=%d: %v", id, err)
		return gerror.New("删除项目失败，请稍后重试")
	}

	return nil
}

// BatchDelete 批量删除项目（ORM 软删除）。
func (p *Project) BatchDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return gerror.New("请选择要删除的项目")
	}

	_, err := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).WhereIn(defaultDao.Project.Columns().Id, ids).Delete()
	if err != nil {
		glog.Errorf(ctx, "批量删除项目失败, ids=%v: %v", ids, err)
		return gerror.New("删除项目失败，请稍后重试")
	}

	return nil
}

// toData 转换为可落库的列值映射；空日期写入 NULL，零值金额 / 负责人照常写入。
func (f *Form) toData(ctx context.Context) g.Map {
	cols := defaultDao.Project.Columns()

	return g.Map{
		"company_id":    currentCompanyId(ctx),
		cols.CustomerId: f.CustomerId,
		cols.Name:       f.Name,
		cols.Address:    f.Address,
		cols.ManagerId:  f.ManagerId,
		cols.Status:     f.Status,
		cols.Amount:     f.Amount,
		cols.StartDate:  parseDate(f.StartDate),
		cols.EndDate:    parseDate(f.EndDate),
		cols.Remark:     f.Remark,
	}
}

// getProject 按主键查询未删除项目，不存在时返回业务错误。
func (p *Project) getProject(ctx context.Context, id uint) (*defaultEntity.Project, error) {
	var project defaultEntity.Project
	model := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Project.Columns().Id, id)
	if !isSuperUser(ctx) {
		userId := currentUserId(ctx)
		model = model.Where("(manager_id = ? OR customer_id IN (SELECT id FROM customer WHERE owner_user_id = ?))", userId, userId)
	}
	err := model.Scan(&project)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询项目失败, id=%d: %v", id, err)
		return nil, gerror.New("获取项目信息失败，请稍后重试")
	}

	if project.Id == 0 {
		return nil, gerror.New("项目不存在或已被删除")
	}

	return &project, nil
}

// ensureCustomerExists 校验所属客户存在，避免挂到已删除 / 不存在的客户上。
func (p *Project) ensureCustomerExists(ctx context.Context, customerId uint) error {
	if customerId == 0 {
		return gerror.New("请选择所属客户")
	}

	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Customer.Columns().Id, customerId)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	count, err := model.Count()
	if err != nil {
		glog.Errorf(ctx, "校验项目所属客户失败, customerId=%d: %v", customerId, err)
		return gerror.New("保存项目失败，请稍后重试")
	}

	if count == 0 {
		return gerror.New("所选客户不存在或已被删除")
	}

	return nil
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

// SelectOptions 项目下拉选项：当前用户可见的未删除项目，label 由「项目名称 · 客户姓名」拼装。
func (p *Project) SelectOptions(ctx context.Context) ([]*ProjectOption, error) {
	cols := defaultDao.Project.Columns()

	model := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	// 非超管仅可看到自己负责的项目，或所属客户归属自己的项目（与 List 的可见范围保持一致）
	if !isSuperUser(ctx) {
		userId := currentUserId(ctx)
		model = model.Where("(manager_id = ? OR customer_id IN (SELECT id FROM customer WHERE owner_user_id = ?))", userId, userId)
	}

	var projects []*defaultEntity.Project
	err := model.
		Fields(cols.Id, cols.Name, cols.CustomerId, cols.Status, cols.Amount).
		OrderDesc(cols.Id).
		Limit(optionLimit).
		Scan(&projects)
	if err != nil {
		glog.Errorf(ctx, "查询项目下拉选项失败: %v", err)
		return nil, gerror.New("获取项目下拉选项失败，请稍后重试")
	}

	if len(projects) == 0 {
		return []*ProjectOption{}, nil
	}

	customerIds := make([]uint, 0, len(projects))
	customerIdSet := make(map[uint]struct{})

	for _, project := range projects {
		if project.CustomerId > 0 {
			if _, ok := customerIdSet[project.CustomerId]; !ok {
				customerIdSet[project.CustomerId] = struct{}{}
				customerIds = append(customerIds, project.CustomerId)
			}
		}
	}

	customerNames := make(map[uint]string)
	if len(customerIds) > 0 {
		var customers []*defaultEntity.Customer
		if scanErr := defaultDao.Customer.Ctx(ctx).
			Where("company_id", currentCompanyId(ctx)).
			WhereIn(defaultDao.Customer.Columns().Id, customerIds).
			Scan(&customers); scanErr != nil {
			// 客户姓名仅用于展示，查询失败不阻断下拉
			glog.Errorf(ctx, "查询项目下拉选项的客户失败, customerIds=%v: %v", customerIds, scanErr)
		} else {
			for _, customer := range customers {
				customerNames[customer.Id] = customer.Name
			}
		}
	}

	result := make([]*ProjectOption, 0, len(projects))

	for _, project := range projects {
		label := project.Name
		if customerName := customerNames[project.CustomerId]; customerName != "" {
			label = label + " · " + customerName
		}

		result = append(result, &ProjectOption{
			Id:     project.Id,
			Label:  label,
			Amount: project.Amount,
			Status: project.Status,
		})
	}

	return result, nil
}

// CreateProgress 后台端新增施工进度，记录人由 controller 从登录上下文取。
func (p *Project) CreateProgress(
	ctx context.Context,
	projectId uint,
	userId uint,
	content string,
	images []string,
	logDate string,
) (int64, error) {
	if _, err := p.getProject(ctx, projectId); err != nil {
		return 0, err
	}

	return p.progress.Create(ctx, &progressLogic.Form{
		ProjectId: projectId,
		UserId:    userId,
		Content:   content,
		Images:    images,
		LogDate:   logDate,
	})
}

// ManagerName 查询单个负责人的展示名（项目详情单条查询用）。
func (p *Project) ManagerName(ctx context.Context, managerId uint) string {
	if managerId == 0 {
		return ""
	}

	return p.loadUserNames(ctx, []uint{managerId})[managerId]
}

// attachExtras 为项目列表补充客户姓名、负责人姓名与最近一条进度（project_log）。
func (p *Project) attachExtras(ctx context.Context, projects []*defaultEntity.Project) []*ProjectInfo {
	result := make([]*ProjectInfo, 0, len(projects))
	if len(projects) == 0 {
		return result
	}

	projectIds := make([]uint, 0, len(projects))
	customerIds := make([]uint, 0, len(projects))
	managerIds := make([]uint, 0, len(projects))
	customerIdSet := make(map[uint]struct{})
	managerIdSet := make(map[uint]struct{})

	for _, project := range projects {
		projectIds = append(projectIds, project.Id)

		if project.CustomerId > 0 {
			if _, ok := customerIdSet[project.CustomerId]; !ok {
				customerIdSet[project.CustomerId] = struct{}{}
				customerIds = append(customerIds, project.CustomerId)
			}
		}

		if project.ManagerId > 0 {
			if _, ok := managerIdSet[project.ManagerId]; !ok {
				managerIdSet[project.ManagerId] = struct{}{}
				managerIds = append(managerIds, project.ManagerId)
			}
		}
	}

	// 客户姓名映射
	customerNames := make(map[uint]string)
	if len(customerIds) > 0 {
		var customers []*defaultEntity.Customer
		if err := defaultDao.Customer.Ctx(ctx).
			Where("company_id", currentCompanyId(ctx)).
			WhereIn(defaultDao.Customer.Columns().Id, customerIds).
			Scan(&customers); err != nil {
			glog.Errorf(ctx, "查询项目所属客户失败, customerIds=%v: %v", customerIds, err)
		} else {
			for _, customer := range customers {
				customerNames[customer.Id] = customer.Name
			}
		}
	}

	managerNames := p.loadUserNames(ctx, managerIds)

	// 最近进度：按业务日期倒序批量取出，遍历时每个项目只保留第一条
	latestLogs := make(map[uint]*defaultEntity.ProjectLog)
	logs, err := defaultDao.ProjectLog.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn(defaultDao.ProjectLog.Columns().ProjectId, projectIds).
		OrderDesc(defaultDao.ProjectLog.Columns().LogDate).
		OrderDesc(defaultDao.ProjectLog.Columns().Id).
		All()
	if err != nil {
		// 进度仅用于展示，查询失败不阻断列表
		glog.Errorf(ctx, "查询项目最近进度失败, projectIds=%v: %v", projectIds, err)
	} else {
		for _, row := range logs {
			var logRow defaultEntity.ProjectLog
			if structErr := row.Struct(&logRow); structErr != nil {
				continue
			}

			if _, ok := latestLogs[logRow.ProjectId]; !ok {
				latestLogs[logRow.ProjectId] = &logRow
			}
		}
	}

	for _, project := range projects {
		info := &ProjectInfo{
			Project:      project,
			CustomerName: customerNames[project.CustomerId],
			ManagerName:  managerNames[project.ManagerId],
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

// attachLogUsers 为进度记录补充记录人姓名。
func (p *Project) attachLogUsers(ctx context.Context, logs []*defaultEntity.ProjectLog) []*LogInfo {
	result := make([]*LogInfo, 0, len(logs))
	if len(logs) == 0 {
		return result
	}

	userIds := make([]uint, 0, len(logs))
	userIdSet := make(map[uint]struct{})
	for _, logRow := range logs {
		if logRow.UserId > 0 {
			if _, ok := userIdSet[logRow.UserId]; !ok {
				userIdSet[logRow.UserId] = struct{}{}
				userIds = append(userIds, logRow.UserId)
			}
		}
	}

	userNames := p.loadUserNames(ctx, userIds)

	for _, logRow := range logs {
		result = append(result, &LogInfo{Log: logRow, UserName: userNames[logRow.UserId]})
	}

	return result
}

// loadUserNames 批量查询 site_user，返回 id → 展示名 的映射。
func (p *Project) loadUserNames(ctx context.Context, userIds []uint) map[uint]string {
	result := make(map[uint]string)
	if len(userIds) == 0 {
		return result
	}

	var users []*defaultEntity.SiteUser
	if err := defaultDao.SiteUser.Ctx(ctx).
		WhereIn(defaultDao.SiteUser.Columns().Id, userIds).
		Scan(&users); err != nil {
		glog.Errorf(ctx, "查询站点用户失败, userIds=%v: %v", userIds, err)
		return result
	}

	for _, user := range users {
		result[user.Id] = firstNonEmpty(user.Nickname, user.Username)
	}

	return result
}

// parseDate 解析 YYYY-MM-DD 日期字符串，空值 / 非法值返回 nil（落库为 NULL）。
func parseDate(value string) *gtime.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	date, err := gtime.StrToTime(value)
	if err != nil || date == nil {
		return nil
	}

	return date
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}
