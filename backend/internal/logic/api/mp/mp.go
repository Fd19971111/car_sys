package mp

import (
	"context"
	"database/sql"
	"errors"
	"math"

	progressLogic "engineServ/internal/logic/api/progress"
	defaultDao "engineServ/model/default/dao"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectCard 项目卡片：项目本体 + 首页展示所需的收款 / 负责人 / 最近进度
type ProjectCard struct {
	Project          *defaultEntity.Project
	PaidAmount       float64
	ManagerName      string
	LatestProgress   string
	LatestProgressAt *gtime.Time
	ProgressPercent  int
}

// LogItem 进度记录 + 记录人
type LogItem struct {
	Log      *defaultEntity.ProjectLog
	UserName string
}

// ProjectSnapshot 项目详情聚合：项目 + 业主 + 负责人 + 报价 + 收款 + 进度
type ProjectSnapshot struct {
	Project         *defaultEntity.Project
	Owner           *defaultEntity.Customer
	ManagerName     string
	ManagerPhone    string
	QuoteAmount     float64
	PaidAmount      float64
	ProgressPercent int
	Logs            []*LogItem
	Quotes          []*defaultEntity.ProjectQuote
	Payments        []*defaultEntity.ProjectPayment
}

type Mp struct {
	progress *progressLogic.Progress
}

func New() *Mp {
	return &Mp{progress: progressLogic.New()}
}

// Overview 业主首页聚合：业主资料 + 名下所有项目的资金与进度。
//
// customerId 由 controller 从令牌上下文取出后传入，不接受前端直接指定。
func (m *Mp) Overview(ctx context.Context, customerId uint) (*defaultEntity.Customer, []*ProjectCard, error) {
	owner, err := m.getCustomer(ctx, customerId)
	if err != nil {
		return nil, nil, err
	}

	projectCols := defaultDao.Project.Columns()

	var projects []*defaultEntity.Project
	err = defaultDao.Project.Ctx(ctx).
		Where(projectCols.CustomerId, customerId).
		OrderDesc(projectCols.CreatedAt).
		OrderDesc(projectCols.Id).
		Scan(&projects)
	if err != nil {
		glog.Errorf(ctx, "查询业主项目失败, customerId=%d: %v", customerId, err)
		return nil, nil, gerror.New("获取项目列表失败，请稍后重试")
	}

	if len(projects) == 0 {
		return owner, []*ProjectCard{}, nil
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

	paidMap := m.sumPaymentByProject(ctx, projectIds)
	latestLogs := m.latestLogByProject(ctx, projectIds)
	users := m.loadUsers(ctx, managerIds)

	cards := make([]*ProjectCard, 0, len(projects))

	for _, project := range projects {
		card := &ProjectCard{
			Project:         project,
			PaidAmount:      roundAmount(paidMap[project.Id]),
			ManagerName:     displayName(users[project.ManagerId]),
			ProgressPercent: progressPercent(project.Status, project.StartDate, project.EndDate),
		}

		if logRow, ok := latestLogs[project.Id]; ok {
			card.LatestProgress = logRow.Content
			// 展示业务日期（进度发生日），历史补录的进度不会错显示为当天
			card.LatestProgressAt = logRow.LogDate
			if card.LatestProgressAt == nil {
				card.LatestProgressAt = logRow.CreatedAt
			}
		}

		cards = append(cards, card)
	}

	return owner, cards, nil
}

// Project 项目详情聚合：一次返回进度 / 报价 / 收款 / 负责人，减少小程序端请求次数。
//
// 关键：必须校验项目归属，否则业主 A 只要把 id 改成 5 就能读到业主 B 的项目。
func (m *Mp) Project(ctx context.Context, projectId uint, customerId uint) (*ProjectSnapshot, error) {
	project, err := m.getOwnedProject(ctx, projectId, customerId)
	if err != nil {
		return nil, err
	}

	snapshot := &ProjectSnapshot{
		Project:         project,
		ProgressPercent: progressPercent(project.Status, project.StartDate, project.EndDate),
	}

	// 业主
	if project.CustomerId > 0 {
		if owner, ownerErr := m.getCustomer(ctx, project.CustomerId); ownerErr == nil {
			snapshot.Owner = owner
		} else {
			// 业主信息仅用于展示，查询失败不阻断详情
			glog.Errorf(ctx, "查询项目业主失败, customerId=%d: %v", project.CustomerId, ownerErr)
		}
	}

	// 负责人（含联系方式，业主端可直接拨号）
	if project.ManagerId > 0 {
		manager := m.loadUsers(ctx, []uint{project.ManagerId})[project.ManagerId]
		snapshot.ManagerName = displayName(manager)
		if manager != nil {
			snapshot.ManagerPhone = manager.Phone
		}
	}

	quoteCols := defaultDao.ProjectQuote.Columns()

	// 报价明细 + 报价合计
	var quotes []*defaultEntity.ProjectQuote
	err = defaultDao.ProjectQuote.Ctx(ctx).
		Where(quoteCols.ProjectId, projectId).
		OrderAsc(quoteCols.Id).
		Scan(&quotes)
	if err != nil {
		glog.Errorf(ctx, "查询项目报价失败, projectId=%d: %v", projectId, err)
		return nil, gerror.New("获取项目报价失败，请稍后重试")
	}

	snapshot.Quotes = quotes

	quoteAmount, err := defaultDao.ProjectQuote.Ctx(ctx).
		Where(quoteCols.ProjectId, projectId).
		Sum(quoteCols.Total)
	if err != nil {
		glog.Errorf(ctx, "统计项目报价合计失败, projectId=%d: %v", projectId, err)
	}
	snapshot.QuoteAmount = roundAmount(quoteAmount)

	paymentCols := defaultDao.ProjectPayment.Columns()

	// 收款记录（按收款日期倒序）
	var payments []*defaultEntity.ProjectPayment
	err = defaultDao.ProjectPayment.Ctx(ctx).
		Where(paymentCols.ProjectId, projectId).
		OrderDesc(paymentCols.PaymentDate).
		OrderDesc(paymentCols.Id).
		Scan(&payments)
	if err != nil {
		glog.Errorf(ctx, "查询项目收款失败, projectId=%d: %v", projectId, err)
		return nil, gerror.New("获取项目收款失败，请稍后重试")
	}

	snapshot.Payments = payments

	paidAmount, err := defaultDao.ProjectPayment.Ctx(ctx).
		Where(paymentCols.ProjectId, projectId).
		Sum(paymentCols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计项目已收款失败, projectId=%d: %v", projectId, err)
	}
	snapshot.PaidAmount = roundAmount(paidAmount)

	// 施工进度（按时间倒序）
	logCols := defaultDao.ProjectLog.Columns()

	var logs []*defaultEntity.ProjectLog
	err = defaultDao.ProjectLog.Ctx(ctx).
		Where(logCols.ProjectId, projectId).
		OrderDesc(logCols.LogDate).
		OrderDesc(logCols.Id).
		Scan(&logs)
	if err != nil {
		glog.Errorf(ctx, "查询项目进度失败, projectId=%d: %v", projectId, err)
		return nil, gerror.New("获取施工进度失败，请稍后重试")
	}

	snapshot.Logs = m.attachLogUsers(ctx, logs)

	return snapshot, nil
}

// CreateProgress 小程序端新增施工进度。
//
// 记录人取「项目负责人」；同时校验项目归属，防止往别人的项目里写数据。
func (m *Mp) CreateProgress(
	ctx context.Context,
	projectId uint,
	customerId uint,
	content string,
	images []string,
) (int64, error) {
	project, err := m.getOwnedProject(ctx, projectId, customerId)
	if err != nil {
		return 0, err
	}

	return m.progress.Create(ctx, &progressLogic.Form{
		ProjectId: projectId,
		UserId:    project.ManagerId,
		Content:   content,
		Images:    images,
	})
}

// getOwnedProject 查询项目并校验归属当前业主。
//
// 归属不符时统一返回「项目不存在」，避免通过错误信息探测出别人的项目ID。
func (m *Mp) getOwnedProject(ctx context.Context, projectId uint, customerId uint) (*defaultEntity.Project, error) {
	if customerId == 0 {
		return nil, gerror.New("登录已失效，请重新登录")
	}

	project, err := m.getProject(ctx, projectId)
	if err != nil {
		return nil, err
	}

	if project.CustomerId != customerId {
		glog.Warningf(ctx, "越权访问项目: customerId=%d 试图访问 projectId=%d(属于 customerId=%d)",
			customerId, projectId, project.CustomerId)

		return nil, gerror.New("项目不存在或已被删除")
	}

	return project, nil
}

// getProject 按主键查询未删除项目。
func (m *Mp) getProject(ctx context.Context, projectId uint) (*defaultEntity.Project, error) {
	var project defaultEntity.Project
	err := defaultDao.Project.Ctx(ctx).
		Where(defaultDao.Project.Columns().Id, projectId).
		Scan(&project)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询项目失败, id=%d: %v", projectId, err)
		return nil, gerror.New("获取项目信息失败，请稍后重试")
	}

	if project.Id == 0 {
		return nil, gerror.New("项目不存在或已被删除")
	}

	return &project, nil
}

// getCustomer 按主键查询未删除业主。
func (m *Mp) getCustomer(ctx context.Context, id uint) (*defaultEntity.Customer, error) {
	var customer defaultEntity.Customer
	err := defaultDao.Customer.Ctx(ctx).
		Where(defaultDao.Customer.Columns().Id, id).
		Scan(&customer)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询业主失败, id=%d: %v", id, err)
		return nil, gerror.New("获取业主信息失败，请稍后重试")
	}

	if customer.Id == 0 {
		return nil, gerror.New("业主不存在或已被删除")
	}

	return &customer, nil
}

// sumPaymentByProject 批量统计各项目已收款金额。
func (m *Mp) sumPaymentByProject(ctx context.Context, projectIds []uint) map[uint]float64 {
	result := make(map[uint]float64)

	rows, err := defaultDao.ProjectPayment.Ctx(ctx).
		Fields("project_id, SUM(amount) AS paid").
		WhereIn(defaultDao.ProjectPayment.Columns().ProjectId, projectIds).
		Group(defaultDao.ProjectPayment.Columns().ProjectId).
		All()
	if err != nil {
		// 收款金额仅用于展示，统计失败不阻断列表
		glog.Errorf(ctx, "统计项目已收款失败, projectIds=%v: %v", projectIds, err)
		return result
	}

	for _, row := range rows {
		result[row["project_id"].Uint()] = row["paid"].Float64()
	}

	return result
}

// latestLogByProject 批量取每个项目最近一条进度记录。
func (m *Mp) latestLogByProject(ctx context.Context, projectIds []uint) map[uint]*defaultEntity.ProjectLog {
	result := make(map[uint]*defaultEntity.ProjectLog)

	logs, err := defaultDao.ProjectLog.Ctx(ctx).
		WhereIn(defaultDao.ProjectLog.Columns().ProjectId, projectIds).
		OrderDesc(defaultDao.ProjectLog.Columns().LogDate).
		OrderDesc(defaultDao.ProjectLog.Columns().Id).
		All()
	if err != nil {
		// 进度仅用于展示，查询失败不阻断列表
		glog.Errorf(ctx, "查询项目最近进度失败, projectIds=%v: %v", projectIds, err)
		return result
	}

	for _, row := range logs {
		var logRow defaultEntity.ProjectLog
		if structErr := row.Struct(&logRow); structErr != nil {
			continue
		}

		if _, ok := result[logRow.ProjectId]; !ok {
			result[logRow.ProjectId] = &logRow
		}
	}

	return result
}

// attachLogUsers 为进度记录补充记录人姓名。
func (m *Mp) attachLogUsers(ctx context.Context, logs []*defaultEntity.ProjectLog) []*LogItem {
	result := make([]*LogItem, 0, len(logs))
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

	userNames := m.loadUsers(ctx, userIds)

	for _, logRow := range logs {
		result = append(result, &LogItem{Log: logRow, UserName: displayName(userNames[logRow.UserId])})
	}

	return result
}

// loadUsers 批量查询 site_user，返回 id → 用户 的映射。
func (m *Mp) loadUsers(ctx context.Context, userIds []uint) map[uint]*defaultEntity.SiteUser {
	result := make(map[uint]*defaultEntity.SiteUser)
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
		result[user.Id] = user
	}

	return result
}

// displayName 站点用户展示名：优先昵称，其次登录账号。
func displayName(user *defaultEntity.SiteUser) string {
	if user == nil {
		return ""
	}

	return firstNonEmpty(user.Nickname, user.Username)
}

// progressPercent 施工进度百分比：待开工 0%，施工中按工期推进比例（10%~90%），已完工 / 已结算 100%。
//
// 与后台项目管理页保持同一套口径，避免两端显示不一致。
func progressPercent(status int, startDate, endDate *gtime.Time) int {
	if status <= 1 {
		return 0
	}

	if status >= 3 {
		return 100
	}

	if startDate == nil || endDate == nil || endDate.Before(startDate) {
		return 50
	}

	total := int(endDate.Sub(startDate).Hours()/24) + 1
	passed := int(gtime.Now().Sub(startDate).Hours()/24) + 1

	if total <= 0 {
		return 50
	}

	percent := int(math.Round(float64(passed) / float64(total) * 100))

	if percent < 10 {
		return 10
	}

	if percent > 90 {
		return 90
	}

	return percent
}

func roundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
