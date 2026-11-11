package quote

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"engineServ/internal/consts"
	defaultDao "engineServ/model/default/dao"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// QuoteInfo 报价明细本体 + 所属项目 / 客户展示信息
type QuoteInfo struct {
	Quote        *defaultEntity.ProjectQuote
	ProjectName  string
	CustomerName string
}

// DetailResult 报价明细详情：明细 + 所属项目 + 客户 + 该项目报价合计
type DetailResult struct {
	Quote       *defaultEntity.ProjectQuote
	Project     *defaultEntity.Project
	Customer    *defaultEntity.Customer
	QuoteAmount float64
}

// Form 报价明细新增 / 编辑入参（与 project_quote 表既有列一一对应）
type Form struct {
	ProjectId int64
	ItemName  string
	Unit      string
	Price     float64
	Quantity  float64
	Remark    string
}

type Quote struct{}

func New() *Quote {
	return &Quote{}
}

// List 报价明细分页列表：关键词（报价项名称 / 项目名称 / 客户姓名）+ 项目过滤，按创建时间倒序。
func (q *Quote) List(
	ctx context.Context,
	keyword string,
	projectId *int64,
	current int,
	size int,
) (list []*QuoteInfo, total int, err error) {
	cols := defaultDao.ProjectQuote.Columns()
	model := defaultDao.ProjectQuote.Ctx(ctx).Where("company_id", currentCompanyId(ctx))

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		matchedIds, matchErr := q.matchQuoteIds(ctx, keyword)
		if matchErr != nil {
			return nil, 0, matchErr
		}

		if len(matchedIds) == 0 {
			return []*QuoteInfo{}, 0, nil
		}

		model = model.WhereIn(cols.Id, matchedIds)
	}

	if projectId != nil && *projectId > 0 {
		model = model.Where(cols.ProjectId, *projectId)
	}

	var records []*defaultEntity.ProjectQuote
	// 报价明细是「项目报价单」的展开：先按项目分组（最新的项目在前），
	// 项目内部按录入顺序展示，这样一屏内就是一份完整可读的报价单。
	err = model.
		Page(current, size).
		OrderDesc(cols.ProjectId).
		OrderAsc(cols.Id).
		ScanAndCount(&records, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询报价明细分页列表失败: %v", err)
		return nil, 0, gerror.New("查询报价明细失败，请稍后重试")
	}

	return q.attachProjects(ctx, records), total, nil
}

// Stats 报价概览统计：条目数 + 报价总额 + 涉及项目数 + 平均单条金额。
func (q *Quote) Stats(ctx context.Context) (totalCount int, totalAmount float64, projectCount int, avgAmount float64, err error) {
	cols := defaultDao.ProjectQuote.Columns()

	totalCount, err = defaultDao.ProjectQuote.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Count()
	if err != nil {
		glog.Errorf(ctx, "统计报价条目数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取报价统计失败，请稍后重试")
	}

	totalAmount, err = defaultDao.ProjectQuote.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Sum(cols.Total)
	if err != nil {
		glog.Errorf(ctx, "统计报价总额失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取报价统计失败，请稍后重试")
	}

	projectCountValue, err := defaultDao.ProjectQuote.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(fmt.Sprintf("COUNT(DISTINCT %s) AS projectCount", cols.ProjectId)).
		Value()
	if err != nil {
		glog.Errorf(ctx, "统计报价涉及项目数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取报价统计失败，请稍后重试")
	}

	projectCount = projectCountValue.Int()

	if totalCount > 0 {
		avgAmount = roundAmount(totalAmount / float64(totalCount))
	}

	return totalCount, roundAmount(totalAmount), projectCount, avgAmount, nil
}

// Detail 报价明细详情：明细 + 所属项目 + 客户 + 该项目报价合计。
func (q *Quote) Detail(ctx context.Context, id int64) (*DetailResult, error) {
	record, err := q.getQuote(ctx, id)
	if err != nil {
		return nil, err
	}

	project, customer := q.loadProject(ctx, record.ProjectId)

	quoteAmount, sumErr := defaultDao.ProjectQuote.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectQuote.Columns().ProjectId, record.ProjectId).
		Sum(defaultDao.ProjectQuote.Columns().Total)
	if sumErr != nil {
		// 合计仅用于展示，失败不阻断详情
		glog.Errorf(ctx, "统计项目报价合计失败, projectId=%d: %v", record.ProjectId, sumErr)
	}

	return &DetailResult{
		Quote:       record,
		Project:     project,
		Customer:    customer,
		QuoteAmount: roundAmount(quoteAmount),
	}, nil
}

// Create 新增报价明细。小计由后端按 单价 × 数量 统一计算，累计报价不得超过合同金额。
func (q *Quote) Create(ctx context.Context, form *Form) (id int64, warning string, err error) {
	if err = form.validate(); err != nil {
		return 0, "", err
	}

	if err = q.ensureProjectExists(ctx, form.ProjectId); err != nil {
		return 0, "", err
	}
	if err = q.ensureQuoteWithinContract(ctx, form.ProjectId, calcTotal(form.Price, form.Quantity), 0); err != nil {
		return 0, "", err
	}

	// 注意：gdb.Model.Data() 是覆盖语义，链式多次调用只保留最后一次的数据，
	// 因此 company_id 必须与表单数据合并到同一个 map 里一次性写入。
	data := form.toData()
	data["company_id"] = currentCompanyId(ctx)

	result, err := defaultDao.ProjectQuote.Ctx(ctx).Data(data).Insert()
	if err != nil {
		glog.Errorf(ctx, "新增报价明细失败: %v", err)
		return 0, "", gerror.New("新增报价明细失败，请稍后重试")
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		glog.Errorf(ctx, "读取新增报价明细ID失败: %v", err)
		return 0, "", gerror.New("新增报价明细失败，请稍后重试")
	}

	return lastInsertId, "", nil
}

// Update 更新报价明细（仅允许修改 project_quote 表既有列），累计报价不得超过合同金额。
func (q *Quote) Update(ctx context.Context, id int64, form *Form) (warning string, err error) {
	if _, err := q.getQuote(ctx, id); err != nil {
		return "", err
	}

	if err := form.validate(); err != nil {
		return "", err
	}

	if err := q.ensureProjectExists(ctx, form.ProjectId); err != nil {
		return "", err
	}
	if err := q.ensureQuoteWithinContract(ctx, form.ProjectId, calcTotal(form.Price, form.Quantity), id); err != nil {
		return "", err
	}

	_, err = defaultDao.ProjectQuote.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectQuote.Columns().Id, id).
		Data(form.toData()).
		Update()
	if err != nil {
		glog.Errorf(ctx, "更新报价明细失败, id=%d: %v", id, err)
		return "", gerror.New("更新报价明细失败，请稍后重试")
	}

	return "", nil
}

// ensureQuoteWithinContract 校验项目累计报价不得超过合同金额；未维护合同金额时不设上限。
func (q *Quote) ensureQuoteWithinContract(ctx context.Context, projectId int64, newTotal float64, excludeQuoteId int64) error {
	companyId := currentCompanyId(ctx)

	row, err := g.DB().GetOne(ctx,
		"SELECT amount FROM project WHERE id=? AND company_id=? AND deleted_at IS NULL", projectId, companyId)
	if err != nil || row == nil || row.IsEmpty() {
		if err != nil {
			glog.Errorf(ctx, "查询项目合同金额失败, projectId=%d: %v", projectId, err)
			return gerror.New("校验项目合同金额失败，请稍后重试")
		}
		return gerror.New("所选项目不存在或已被删除")
	}
	contract := roundAmount(row["amount"].Float64())
	if contract <= 0 {
		return nil
	}

	query := "SELECT COALESCE(SUM(total),0) FROM project_quote WHERE project_id=? AND company_id=? AND deleted_at IS NULL"
	args := []any{projectId, companyId}
	if excludeQuoteId > 0 {
		query += " AND id<>?"
		args = append(args, excludeQuoteId)
	}
	quoteValue, err := g.DB().GetValue(ctx, query, args...)
	if err != nil {
		glog.Errorf(ctx, "统计项目报价合计失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验项目报价合计失败，请稍后重试")
	}

	quoteTotal := roundAmount(quoteValue.Float64() + newTotal)
	if quoteTotal <= contract {
		return nil
	}

	return gerror.New(fmt.Sprintf("项目报价合计 ¥%.2f，不能超过合同金额 ¥%.2f", quoteTotal, contract))
}

// Delete 删除单条报价明细（ORM 软删除，写入 deleted_at）。
func (q *Quote) Delete(ctx context.Context, id int64) error {
	if _, err := q.getQuote(ctx, id); err != nil {
		return err
	}

	_, err := defaultDao.ProjectQuote.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.ProjectQuote.Columns().Id, id).Delete()
	if err != nil {
		glog.Errorf(ctx, "删除报价明细失败, id=%d: %v", id, err)
		return gerror.New("删除报价明细失败，请稍后重试")
	}

	return nil
}

// BatchDelete 批量删除报价明细（ORM 软删除）。
func (q *Quote) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return gerror.New("请选择要删除的报价明细")
	}

	_, err := defaultDao.ProjectQuote.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).WhereIn(defaultDao.ProjectQuote.Columns().Id, ids).Delete()
	if err != nil {
		glog.Errorf(ctx, "批量删除报价明细失败, ids=%v: %v", ids, err)
		return gerror.New("删除报价明细失败，请稍后重试")
	}

	return nil
}

// validate 业务校验：项目 / 报价项名称必填，单价与数量不能为负数。
//
// 说明：接口层已有 gvalid 规则，这里再兜一次，避免参数绑定差异写入脏数据。
func (f *Form) validate() error {
	if f.ProjectId <= 0 {
		return gerror.New("请选择所属项目")
	}

	if strings.TrimSpace(f.ItemName) == "" {
		return gerror.New("请输入报价项名称")
	}

	if f.Price < 0 {
		return gerror.New("单价不能为负数")
	}

	if f.Quantity < 0 {
		return gerror.New("数量不能为负数")
	}

	return nil
}

// toData 转换为可落库的列值映射；小计由后端按 单价 × 数量 计算。
func (f *Form) toData() g.Map {
	cols := defaultDao.ProjectQuote.Columns()

	return g.Map{
		cols.ProjectId: f.ProjectId,
		cols.ItemName:  f.ItemName,
		cols.Unit:      f.Unit,
		cols.Price:     f.Price,
		cols.Quantity:  f.Quantity,
		cols.Total:     calcTotal(f.Price, f.Quantity),
		cols.Remark:    f.Remark,
	}
}

// matchQuoteIds 关键词命中：报价项名称直接匹配，或所属项目名称 / 客户姓名匹配。
func (q *Quote) matchQuoteIds(ctx context.Context, keyword string) ([]int64, error) {
	cols := defaultDao.ProjectQuote.Columns()
	like := "%" + keyword + "%"

	projectIds, err := q.matchProjectIds(ctx, like)
	if err != nil {
		return nil, err
	}

	matchedModel := defaultDao.ProjectQuote.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(cols.Id).
		WhereLike(cols.ItemName, like)

	if len(projectIds) > 0 {
		matchedModel = matchedModel.WhereOrIn(cols.ProjectId, projectIds)
	}

	rows, err := matchedModel.All()
	if err != nil {
		glog.Errorf(ctx, "查询关键词命中的报价明细失败: %v", err)
		return nil, gerror.New("查询报价明细失败，请稍后重试")
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row[cols.Id].Int64())
	}

	return ids, nil
}

// matchProjectIds 关键词命中：项目名称直接匹配，或所属客户姓名匹配。
func (q *Quote) matchProjectIds(ctx context.Context, like string) ([]uint, error) {
	projectCols := defaultDao.Project.Columns()
	customerCols := defaultDao.Customer.Columns()

	var matchedCustomerIds []uint
	customerRows, customerErr := defaultDao.Customer.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(customerCols.Id).
		WhereLike(customerCols.Name, like).
		All()
	if customerErr != nil {
		glog.Errorf(ctx, "查询关键词命中的客户失败: %v", customerErr)
		return nil, gerror.New("查询报价明细失败，请稍后重试")
	}

	for _, row := range customerRows {
		matchedCustomerIds = append(matchedCustomerIds, row[customerCols.Id].Uint())
	}

	matchedModel := defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(projectCols.Id).
		WhereLike(projectCols.Name, like)

	if len(matchedCustomerIds) > 0 {
		matchedModel = matchedModel.WhereOrIn(projectCols.CustomerId, matchedCustomerIds)
	}

	rows, err := matchedModel.All()
	if err != nil {
		glog.Errorf(ctx, "查询关键词命中的项目失败: %v", err)
		return nil, gerror.New("查询报价明细失败，请稍后重试")
	}

	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row[projectCols.Id].Uint())
	}

	return ids, nil
}

// attachProjects 为报价明细批量补充项目名称与客户姓名。
func (q *Quote) attachProjects(ctx context.Context, records []*defaultEntity.ProjectQuote) []*QuoteInfo {
	result := make([]*QuoteInfo, 0, len(records))
	if len(records) == 0 {
		return result
	}

	projectIds := make([]int64, 0, len(records))
	projectIdSet := make(map[int64]struct{})

	for _, record := range records {
		if record.ProjectId > 0 {
			if _, ok := projectIdSet[record.ProjectId]; !ok {
				projectIdSet[record.ProjectId] = struct{}{}
				projectIds = append(projectIds, record.ProjectId)
			}
		}
	}

	projectNames := make(map[int64]string)
	projectCustomers := make(map[int64]int64)
	customerNames := make(map[int64]string)

	if len(projectIds) > 0 {
		var projects []*defaultEntity.Project
		if err := defaultDao.Project.Ctx(ctx).
			Where("company_id", currentCompanyId(ctx)).
			WhereIn(defaultDao.Project.Columns().Id, projectIds).
			Scan(&projects); err != nil {
			glog.Errorf(ctx, "查询报价明细所属项目失败, projectIds=%v: %v", projectIds, err)
		} else {
			customerIds := make([]uint, 0, len(projects))
			customerIdSet := make(map[uint]struct{})

			for _, project := range projects {
				projectNames[int64(project.Id)] = project.Name
				projectCustomers[int64(project.Id)] = int64(project.CustomerId)

				if project.CustomerId > 0 {
					if _, ok := customerIdSet[project.CustomerId]; !ok {
						customerIdSet[project.CustomerId] = struct{}{}
						customerIds = append(customerIds, project.CustomerId)
					}
				}
			}

			if len(customerIds) > 0 {
				var customers []*defaultEntity.Customer
				if customerErr := defaultDao.Customer.Ctx(ctx).
					Where("company_id", currentCompanyId(ctx)).
					WhereIn(defaultDao.Customer.Columns().Id, customerIds).
					Scan(&customers); customerErr != nil {
					glog.Errorf(ctx, "查询报价明细所属客户失败, customerIds=%v: %v", customerIds, customerErr)
				} else {
					for _, customer := range customers {
						customerNames[int64(customer.Id)] = customer.Name
					}
				}
			}
		}
	}

	for _, record := range records {
		result = append(result, &QuoteInfo{
			Quote:        record,
			ProjectName:  projectNames[record.ProjectId],
			CustomerName: customerNames[projectCustomers[record.ProjectId]],
		})
	}

	return result
}

// loadProject 查询报价明细所属项目与客户，缺失时返回 nil。
func (q *Quote) loadProject(ctx context.Context, projectId int64) (*defaultEntity.Project, *defaultEntity.Customer) {
	if projectId <= 0 {
		return nil, nil
	}

	var project defaultEntity.Project
	if err := defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.Project.Columns().Id, projectId).
		Scan(&project); err != nil {
		glog.Errorf(ctx, "查询报价明细所属项目失败, projectId=%d: %v", projectId, err)
		return nil, nil
	}

	if project.Id == 0 {
		return nil, nil
	}

	if project.CustomerId == 0 {
		return &project, nil
	}

	var customer defaultEntity.Customer
	if err := defaultDao.Customer.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.Customer.Columns().Id, project.CustomerId).
		Scan(&customer); err != nil {
		glog.Errorf(ctx, "查询报价明细所属客户失败, customerId=%d: %v", project.CustomerId, err)
		return &project, nil
	}

	return &project, &customer
}

// getQuote 按主键查询未删除报价明细，不存在时返回业务错误。
func (q *Quote) getQuote(ctx context.Context, id int64) (*defaultEntity.ProjectQuote, error) {
	var record defaultEntity.ProjectQuote
	err := defaultDao.ProjectQuote.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.ProjectQuote.Columns().Id, id).Scan(&record)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询报价明细失败, id=%d: %v", id, err)
		return nil, gerror.New("获取报价明细失败，请稍后重试")
	}

	if record.Id == 0 {
		return nil, gerror.New("报价明细不存在或已被删除")
	}

	return &record, nil
}

// ensureProjectExists 校验所属项目存在，避免挂到已删除 / 不存在的项目上。
func (q *Quote) ensureProjectExists(ctx context.Context, projectId int64) error {
	if projectId <= 0 {
		return gerror.New("请选择所属项目")
	}

	count, err := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Project.Columns().Id, projectId).Count()
	if err != nil {
		glog.Errorf(ctx, "校验报价明细所属项目失败, projectId=%d: %v", projectId, err)
		return gerror.New("保存报价明细失败，请稍后重试")
	}

	if count == 0 {
		return gerror.New("所选项目不存在或已被删除")
	}

	return nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

// calcTotal 小计 = 单价 × 数量，保留两位小数。
func calcTotal(price, quantity float64) float64 {
	return math.Round(price*quantity*100) / 100
}

// roundAmount 金额保留两位小数。
func roundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
