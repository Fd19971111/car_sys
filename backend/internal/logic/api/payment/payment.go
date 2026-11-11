package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"engineServ/internal/consts"
	"engineServ/internal/logic/fin/booking"
	"engineServ/internal/logic/oplog"
	defaultDao "engineServ/model/default/dao"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// PaymentInfo 收款记录本体 + 所属项目 / 客户 / 资金账户 / 凭证展示信息
type PaymentInfo struct {
	Payment         *defaultEntity.ProjectPayment
	ProjectName     string
	CustomerName    string
	FundAccountName string
	VoucherNo       string
}

// DetailResult 收款记录详情：记录 + 所属项目 + 客户 + 该项目累计已收款 + 记账信息
type DetailResult struct {
	Payment         *defaultEntity.ProjectPayment
	Project         *defaultEntity.Project
	Customer        *defaultEntity.Customer
	PaidAmount      float64
	FundAccountName string
	VoucherNo       string
	Entries         []*booking.EntryView
}

// Form 收款记录新增 / 编辑入参
type Form struct {
	ProjectId     int64
	Amount        float64
	PaymentType   string
	PaymentDate   string
	FundAccountId int64
	Remark        string
}

type Payment struct{}

func New() *Payment {
	return &Payment{}
}

// List 收款记录分页列表：关键词（项目名称 / 客户姓名 / 款项类型 / 备注）+ 项目 + 款项类型过滤。
func (p *Payment) List(
	ctx context.Context,
	keyword string,
	projectId *int64,
	paymentType string,
	current int,
	size int,
) (list []*PaymentInfo, total int, err error) {
	cols := defaultDao.ProjectPayment.Columns()
	model := defaultDao.ProjectPayment.Ctx(ctx).Where("company_id", currentCompanyId(ctx))

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		matchedIds, matchErr := p.matchPaymentIds(ctx, keyword)
		if matchErr != nil {
			return nil, 0, matchErr
		}

		if len(matchedIds) == 0 {
			return []*PaymentInfo{}, 0, nil
		}

		model = model.WhereIn(cols.Id, matchedIds)
	}

	if projectId != nil && *projectId > 0 {
		model = model.Where(cols.ProjectId, *projectId)
	}

	paymentType = strings.TrimSpace(paymentType)
	if paymentType != "" {
		model = model.Where(cols.PaymentType, paymentType)
	}

	var records []*defaultEntity.ProjectPayment
	err = model.
		Page(current, size).
		OrderDesc(cols.PaymentDate).
		OrderDesc(cols.Id).
		ScanAndCount(&records, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询收款记录分页列表失败: %v", err)
		return nil, 0, gerror.New("查询收款记录失败，请稍后重试")
	}

	return p.attach(ctx, records), total, nil
}

// Stats 收款概览统计：笔数 + 收款总额 + 本月收款 + 涉及项目数。
func (p *Payment) Stats(ctx context.Context) (totalCount int, totalAmount float64, monthAmount float64, projectCount int, err error) {
	cols := defaultDao.ProjectPayment.Columns()

	totalCount, err = defaultDao.ProjectPayment.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Count()
	if err != nil {
		glog.Errorf(ctx, "统计收款笔数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	totalAmount, err = defaultDao.ProjectPayment.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计收款总额失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	// 本月收款：payment_date ∈ [本月1号, 下月1号)
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	monthAmount, err = defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereGTE(cols.PaymentDate, monthStart.Format("2006-01-02")).
		WhereLT(cols.PaymentDate, nextMonthStart.Format("2006-01-02")).
		Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计本月收款失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	projectCountValue, err := defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(fmt.Sprintf("COUNT(DISTINCT %s) AS projectCount", cols.ProjectId)).
		Value()
	if err != nil {
		glog.Errorf(ctx, "统计收款涉及项目数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	return totalCount, roundAmount(totalAmount), roundAmount(monthAmount), projectCountValue.Int(), nil
}

// Detail 收款记录详情：记录 + 所属项目 + 客户 + 该项目累计已收款 + 凭证分录。
func (p *Payment) Detail(ctx context.Context, id int64) (*DetailResult, error) {
	record, err := p.getPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	project, customer := p.loadProject(ctx, record.ProjectId)

	paidAmount, sumErr := defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectPayment.Columns().ProjectId, record.ProjectId).
		Sum(defaultDao.ProjectPayment.Columns().Amount)
	if sumErr != nil {
		// 累计已收仅用于展示，失败不阻断详情
		glog.Errorf(ctx, "统计项目累计收款失败, projectId=%d: %v", record.ProjectId, sumErr)
	}

	result := &DetailResult{
		Payment:    record,
		Project:    project,
		Customer:   customer,
		PaidAmount: roundAmount(paidAmount),
	}

	result.FundAccountName = p.fundAccountName(ctx, record.FundAccountId)
	result.VoucherNo = p.voucherNo(ctx, record.VoucherId)
	result.Entries = p.ledgerEntries(ctx, record.VoucherId)

	return result, nil
}

// Create 新增收款记录并自动生成收款凭证，返回新记录ID与凭证号。
//
// 业务数据与凭证在同一个事务里写入：凭证记账失败（例如会计期间已结账、
// 资金账户不可用）时，收款记录也不会落库，避免出现「有业务没记账」的账。
func (p *Payment) Create(ctx context.Context, form *Form, userId uint) (id int64, voucherNo string, err error) {
	if err = form.validate(); err != nil {
		return 0, "", err
	}

	if err = p.ensureProjectExists(ctx, form.ProjectId); err != nil {
		return 0, "", err
	}

	data, err := form.toData()
	if err != nil {
		return 0, "", err
	}

	err = defaultDao.ProjectPayment.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 超收校验：累计收款不能超过合同金额
		if checkErr := p.ensureContractNotExceeded(ctx, tx, form.ProjectId, form.Amount, 0); checkErr != nil {
			return checkErr
		}

		// 防重复提交：短时间内出现完全相同的收款记录直接拒绝，避免双击 / 网络重试重复入账
		if dupErr := p.ensureNotDuplicate(ctx, tx, form.ProjectId, data[defaultDao.ProjectPayment.Columns().PaymentDate], form.Amount, form.PaymentType); dupErr != nil {
			return dupErr
		}

		data["company_id"] = currentCompanyId(ctx)
		result, insertErr := defaultDao.ProjectPayment.Ctx(ctx).TX(tx).Data(data).Insert()
		if insertErr != nil {
			glog.Errorf(ctx, "新增收款记录失败: %v", insertErr)
			return gerror.New("新增收款记录失败，请稍后重试")
		}

		id, insertErr = result.LastInsertId()
		if insertErr != nil {
			glog.Errorf(ctx, "读取新增收款记录ID失败: %v", insertErr)
			return gerror.New("新增收款记录失败，请稍后重试")
		}

		voucherNo, insertErr = p.book(ctx, tx, id, form, userId)
		return insertErr
	})
	if err != nil {
		return 0, "", err
	}

	p.logReceipt(ctx, oplog.ActionCreate, id, form, voucherNo)

	return id, voucherNo, nil
}

// Update 更新收款记录，并同步更新自动生成的凭证。
//
// 只有原本已经记账的记录才会重新记账：历史数据（voucher_id=0）保持未记账，
// 这样补录 / 修正很久以前的收款时不会被已结账期间挡住。
func (p *Payment) Update(ctx context.Context, id int64, form *Form, userId uint) (voucherNo string, err error) {
	if err = form.validate(); err != nil {
		return "", err
	}

	if err = p.ensureProjectExists(ctx, form.ProjectId); err != nil {
		return "", err
	}

	data, err := form.toData()
	if err != nil {
		return "", err
	}
	data["company_id"] = currentCompanyId(ctx)

	err = defaultDao.ProjectPayment.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		oldVoucherId, lockErr := p.lockVoucherId(ctx, tx, id, true)
		if lockErr != nil {
			return lockErr
		}

		// 超收校验：修改后的累计收款同样不能超过合同金额（排除本条记录）
		if checkErr := p.ensureContractNotExceeded(ctx, tx, form.ProjectId, form.Amount, id); checkErr != nil {
			return checkErr
		}

		// 先作废旧凭证，撤销它的余额与流水，再按新数据重新记账
		if oldVoucherId > 0 {
			if lockErr = booking.New().Void(ctx, tx, oldVoucherId); lockErr != nil {
				return lockErr
			}
			if _, lockErr = defaultDao.ProjectPayment.Ctx(ctx).TX(tx).
				Where("company_id", currentCompanyId(ctx)).
				Where(defaultDao.ProjectPayment.Columns().Id, id).
				Data(g.Map{defaultDao.ProjectPayment.Columns().VoucherId: 0}).
				Update(); lockErr != nil {
				glog.Errorf(ctx, "重置收款记录凭证ID失败, id=%d: %v", id, lockErr)
				return gerror.New("更新收款记录失败，请稍后重试")
			}
		}

		if _, lockErr = defaultDao.ProjectPayment.Ctx(ctx).TX(tx).
			Where("company_id", currentCompanyId(ctx)).
			Where(defaultDao.ProjectPayment.Columns().Id, id).
			Data(data).
			Update(); lockErr != nil {
			glog.Errorf(ctx, "更新收款记录失败, id=%d: %v", id, lockErr)
			return gerror.New("更新收款记录失败，请稍后重试")
		}

		if oldVoucherId == 0 {
			return nil
		}

		voucherNo, lockErr = p.book(ctx, tx, id, form, userId)
		return lockErr
	})
	if err != nil {
		return "", err
	}

	p.logReceipt(ctx, oplog.ActionUpdate, id, form, voucherNo)

	return voucherNo, nil
}

// Delete 删除单条收款记录：先作废凭证，再软删除业务数据。
func (p *Payment) Delete(ctx context.Context, id int64) error {
	// 日志需要原始金额与项目，删除前先留快照
	snapshot := p.loadReceiptSnapshot(ctx, id)

	err := defaultDao.ProjectPayment.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		voucherId, lockErr := p.lockVoucherId(ctx, tx, id, true)
		if lockErr != nil {
			return lockErr
		}

		if lockErr = booking.New().Void(ctx, tx, voucherId); lockErr != nil {
			return lockErr
		}

		if _, lockErr = defaultDao.ProjectPayment.Ctx(ctx).TX(tx).
			Where("company_id", currentCompanyId(ctx)).
			Where(defaultDao.ProjectPayment.Columns().Id, id).
			Delete(); lockErr != nil {
			glog.Errorf(ctx, "删除收款记录失败, id=%d: %v", id, lockErr)
			return gerror.New("删除收款记录失败，请稍后重试")
		}

		return nil
	})
	if err != nil {
		return err
	}

	if snapshot != nil {
		p.logReceipt(ctx, oplog.ActionDelete, id, snapshot.Form, snapshot.VoucherNo)
	}

	return nil
}

// BatchDelete 批量删除收款记录：逐条作废凭证后软删除。
func (p *Payment) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return gerror.New("请选择要删除的收款记录")
	}

	snapshots := make(map[int64]*receiptSnapshot, len(ids))
	for _, id := range ids {
		if snapshot := p.loadReceiptSnapshot(ctx, id); snapshot != nil {
			snapshots[id] = snapshot
		}
	}

	err := defaultDao.ProjectPayment.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, id := range ids {
			voucherId, lockErr := p.lockVoucherId(ctx, tx, id, false)
			if lockErr != nil {
				return lockErr
			}
			if lockErr = booking.New().Void(ctx, tx, voucherId); lockErr != nil {
				return lockErr
			}
		}

		if _, lockErr := defaultDao.ProjectPayment.Ctx(ctx).TX(tx).
			Where("company_id", currentCompanyId(ctx)).
			WhereIn(defaultDao.ProjectPayment.Columns().Id, ids).
			Delete(); lockErr != nil {
			glog.Errorf(ctx, "批量删除收款记录失败, ids=%v: %v", ids, lockErr)
			return gerror.New("删除收款记录失败，请稍后重试")
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, id := range ids {
		if snapshot, ok := snapshots[id]; ok {
			p.logReceipt(ctx, oplog.ActionDelete, id, snapshot.Form, snapshot.VoucherNo)
		}
	}

	return nil
}

// book 为收款记录生成借贷凭证，并把凭证ID写回业务表。
//
// 财务尚未初始化时直接跳过：旧数据与未启用财务的公司仍然可以正常登记收款，
// 只是该笔记录保持「未记账」状态。
func (p *Payment) book(ctx context.Context, tx gdb.TX, id int64, form *Form, userId uint) (string, error) {
	initialized, err := booking.New().IsInitialized(ctx, tx)
	if err != nil {
		return "", err
	}
	if !initialized {
		return "", nil
	}
	// 财务已初始化时必须落到具体资金账户，否则会产生「有收款、无凭证」的账
	if form.FundAccountId <= 0 {
		return "", gerror.New("请选择资金账户，系统需要据此记录这笔钱进入了哪个账户")
	}

	projectName, err := p.projectName(ctx, tx, form.ProjectId)
	if err != nil {
		return "", err
	}

	result, err := booking.New().PostReceipt(ctx, tx, &booking.Receipt{
		SourceId:      id,
		ProjectName:   projectName,
		PaymentType:   form.PaymentType,
		Date:          form.PaymentDate,
		Amount:        form.Amount,
		FundAccountId: form.FundAccountId,
		CreatedBy:     userId,
	})
	if err != nil {
		return "", err
	}

	if _, err = defaultDao.ProjectPayment.Ctx(ctx).TX(tx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectPayment.Columns().Id, id).
		Data(g.Map{defaultDao.ProjectPayment.Columns().VoucherId: result.VoucherId}).
		Update(); err != nil {
		glog.Errorf(ctx, "写回收款记录凭证ID失败, id=%d: %v", id, err)
		return "", gerror.New("保存收款记录失败，请稍后重试")
	}

	return result.VoucherNo, nil
}

// validate 业务校验：项目 / 款项类型 / 收款日期必填，金额必须大于 0。
//
// 说明：接口层已有 gvalid 规则，这里再兜一次，保证金额这种关键字段不会因
// 参数绑定差异（0 值、字段缺失）而写入脏数据。
func (f *Form) validate() error {
	if f.ProjectId <= 0 {
		return gerror.New("请选择所属项目")
	}

	// 金额统一按两位小数落库，避免业务表与凭证金额出现分位差异
	f.Amount = roundAmount(f.Amount)
	if f.Amount <= 0 {
		return gerror.New("收款金额必须大于 0")
	}

	if strings.TrimSpace(f.PaymentType) == "" {
		return gerror.New("请选择款项类型")
	}

	return nil
}

// toData 转换为可落库的列值映射；收款日期为必填，解析失败返回业务错误。
func (f *Form) toData() (g.Map, error) {
	cols := defaultDao.ProjectPayment.Columns()

	paymentDate, err := parseDate(f.PaymentDate)
	if err != nil {
		return nil, err
	}

	return g.Map{
		cols.ProjectId:     f.ProjectId,
		cols.Amount:        f.Amount,
		cols.PaymentType:   f.PaymentType,
		cols.PaymentDate:   paymentDate,
		cols.FundAccountId: f.FundAccountId,
		cols.Remark:        f.Remark,
	}, nil
}

// matchPaymentIds 关键词命中：款项类型 / 备注直接匹配，或所属项目名称 / 客户姓名匹配。
func (p *Payment) matchPaymentIds(ctx context.Context, keyword string) ([]int64, error) {
	cols := defaultDao.ProjectPayment.Columns()
	like := "%" + keyword + "%"

	projectIds, err := p.matchProjectIds(ctx, like)
	if err != nil {
		return nil, err
	}

	matchedModel := defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(cols.Id).
		WhereLike(cols.PaymentType, like).
		WhereOrLike(cols.Remark, like)

	if len(projectIds) > 0 {
		matchedModel = matchedModel.WhereOrIn(cols.ProjectId, projectIds)
	}

	rows, err := matchedModel.All()
	if err != nil {
		glog.Errorf(ctx, "查询关键词命中的收款记录失败: %v", err)
		return nil, gerror.New("查询收款记录失败，请稍后重试")
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row[cols.Id].Int64())
	}

	return ids, nil
}

// matchProjectIds 关键词命中：项目名称直接匹配，或所属客户姓名匹配。
func (p *Payment) matchProjectIds(ctx context.Context, like string) ([]uint, error) {
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
		return nil, gerror.New("查询收款记录失败，请稍后重试")
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
		return nil, gerror.New("查询收款记录失败，请稍后重试")
	}

	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row[projectCols.Id].Uint())
	}

	return ids, nil
}

// attach 为收款记录批量补充项目名称、客户姓名、资金账户与凭证号。
func (p *Payment) attach(ctx context.Context, records []*defaultEntity.ProjectPayment) []*PaymentInfo {
	result := make([]*PaymentInfo, 0, len(records))
	if len(records) == 0 {
		return result
	}

	projectIds := make([]int64, 0, len(records))
	projectIdSet := make(map[int64]struct{})
	fundAccountIds := make([]int64, 0, len(records))
	fundAccountSet := make(map[int64]struct{})
	voucherIds := make([]int64, 0, len(records))

	for _, record := range records {
		if record.ProjectId > 0 {
			if _, ok := projectIdSet[record.ProjectId]; !ok {
				projectIdSet[record.ProjectId] = struct{}{}
				projectIds = append(projectIds, record.ProjectId)
			}
		}
		if record.FundAccountId > 0 {
			if _, ok := fundAccountSet[record.FundAccountId]; !ok {
				fundAccountSet[record.FundAccountId] = struct{}{}
				fundAccountIds = append(fundAccountIds, record.FundAccountId)
			}
		}
		if record.VoucherId > 0 {
			voucherIds = append(voucherIds, record.VoucherId)
		}
	}

	projectNames, projectCustomers, customerNames := p.loadProjectMaps(ctx, projectIds)
	fundAccountNames := p.fundAccountNames(ctx, fundAccountIds)

	voucherNos := make(map[int64]string, len(voucherIds))
	if len(voucherIds) > 0 {
		if values, err := booking.New().VoucherNos(ctx, voucherIds); err != nil {
			// 凭证号仅用于展示，查询失败不阻断列表
			glog.Errorf(ctx, "查询收款凭证号失败, ids=%v: %v", voucherIds, err)
		} else {
			voucherNos = values
		}
	}

	for _, record := range records {
		result = append(result, &PaymentInfo{
			Payment:         record,
			ProjectName:     projectNames[record.ProjectId],
			CustomerName:    customerNames[projectCustomers[record.ProjectId]],
			FundAccountName: fundAccountNames[record.FundAccountId],
			VoucherNo:       voucherNos[record.VoucherId],
		})
	}

	return result
}

// loadProjectMaps 批量加载项目名称、项目所属客户、客户姓名。
func (p *Payment) loadProjectMaps(ctx context.Context, projectIds []int64) (map[int64]string, map[int64]int64, map[int64]string) {
	projectNames := make(map[int64]string)
	projectCustomers := make(map[int64]int64)
	customerNames := make(map[int64]string)

	if len(projectIds) == 0 {
		return projectNames, projectCustomers, customerNames
	}

	var projects []*defaultEntity.Project
	if err := defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn(defaultDao.Project.Columns().Id, projectIds).
		Scan(&projects); err != nil {
		glog.Errorf(ctx, "查询收款记录所属项目失败, projectIds=%v: %v", projectIds, err)
		return projectNames, projectCustomers, customerNames
	}

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

	if len(customerIds) == 0 {
		return projectNames, projectCustomers, customerNames
	}

	var customers []*defaultEntity.Customer
	if err := defaultDao.Customer.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereIn(defaultDao.Customer.Columns().Id, customerIds).
		Scan(&customers); err != nil {
		glog.Errorf(ctx, "查询收款记录所属客户失败, customerIds=%v: %v", customerIds, err)
		return projectNames, projectCustomers, customerNames
	}

	for _, customer := range customers {
		customerNames[int64(customer.Id)] = customer.Name
	}

	return projectNames, projectCustomers, customerNames
}

// fundAccountNames 批量查询资金账户名称。
func (p *Payment) fundAccountNames(ctx context.Context, ids []int64) map[int64]string {
	names := make(map[int64]string, len(ids))
	if len(ids) == 0 {
		return names
	}

	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")

	args = append(args, currentCompanyId(ctx))
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(
		"SELECT id,account_name FROM fin_fund_account WHERE id IN (%s) AND company_id=?", placeholders), args...)
	if err != nil {
		glog.Errorf(ctx, "查询收款资金账户失败, ids=%v: %v", ids, err)
		return names
	}

	for _, record := range records {
		names[record["id"].Int64()] = record["account_name"].String()
	}

	return names
}

// fundAccountName 查询单个资金账户名称。
func (p *Payment) fundAccountName(ctx context.Context, id int64) string {
	if id <= 0 {
		return ""
	}

	value, err := g.DB().GetValue(ctx, "SELECT account_name FROM fin_fund_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "查询收款资金账户失败, id=%d: %v", id, err)
		return ""
	}

	return value.String()
}

// voucherNo 查询凭证号（仅已记账凭证有效）。
func (p *Payment) voucherNo(ctx context.Context, voucherId int64) string {
	if voucherId <= 0 {
		return ""
	}

	values, err := booking.New().VoucherNos(ctx, []int64{voucherId})
	if err != nil {
		glog.Errorf(ctx, "查询收款凭证号失败, voucherId=%d: %v", voucherId, err)
		return ""
	}

	return values[voucherId]
}

// ledgerEntries 查询凭证分录，用于详情页展示借贷方向。
func (p *Payment) ledgerEntries(ctx context.Context, voucherId int64) []*booking.EntryView {
	entries, err := booking.New().Entries(ctx, voucherId)
	if err != nil {
		glog.Errorf(ctx, "查询收款凭证分录失败, voucherId=%d: %v", voucherId, err)
		return []*booking.EntryView{}
	}

	return entries
}

// projectName 查询项目名称（事务内使用）。
func (p *Payment) projectName(ctx context.Context, tx gdb.TX, projectId int64) (string, error) {
	value, err := tx.GetValue("SELECT name FROM project WHERE id=? AND company_id=? AND deleted_at IS NULL", projectId, currentCompanyId(ctx))
	if err != nil {
		return "", gerror.Wrap(err, "查询所属项目失败")
	}

	return value.String(), nil
}

// loadProject 查询收款记录所属项目与客户，缺失时返回 nil。
func (p *Payment) loadProject(ctx context.Context, projectId int64) (*defaultEntity.Project, *defaultEntity.Customer) {
	if projectId <= 0 {
		return nil, nil
	}

	var project defaultEntity.Project
	if err := defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.Project.Columns().Id, projectId).
		Scan(&project); err != nil {
		glog.Errorf(ctx, "查询收款记录所属项目失败, projectId=%d: %v", projectId, err)
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
		glog.Errorf(ctx, "查询收款记录所属客户失败, customerId=%d: %v", project.CustomerId, err)
		return &project, nil
	}

	return &project, &customer
}

// getPayment 按主键查询未删除收款记录，不存在时返回业务错误。
func (p *Payment) getPayment(ctx context.Context, id int64) (*defaultEntity.ProjectPayment, error) {
	var record defaultEntity.ProjectPayment
	err := defaultDao.ProjectPayment.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.ProjectPayment.Columns().Id, id).Scan(&record)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询收款记录失败, id=%d: %v", id, err)
		return nil, gerror.New("获取收款记录失败，请稍后重试")
	}

	if record.Id == 0 {
		return nil, gerror.New("收款记录不存在或已被删除")
	}

	return &record, nil
}

// lockVoucherId 在事务内锁定收款记录，返回它关联的凭证ID。
//
// mustExist=false 用于批量删除：记录已被删除时跳过，保证操作幂等。
func (p *Payment) lockVoucherId(ctx context.Context, tx gdb.TX, id int64, mustExist bool) (int64, error) {
	row, err := tx.GetOne("SELECT id,voucher_id FROM project_payment WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "锁定收款记录失败, id=%d: %v", id, err)
		return 0, gerror.New("操作收款记录失败，请稍后重试")
	}

	if row == nil || row.IsEmpty() {
		if mustExist {
			return 0, gerror.New("收款记录不存在或已被删除")
		}
		return 0, nil
	}

	return row["voucher_id"].Int64(), nil
}

// ensureNotDuplicate 拦截短时间内的重复提交。
//
// 同一项目、同一天、同一金额、同一款项类型的记录在 2 分钟内重复出现时，
// 基本可以判定为双击或网络重试，此时拒绝写入，避免重复收款与重复凭证。
func (p *Payment) ensureNotDuplicate(ctx context.Context, tx gdb.TX, projectId int64, dateValue interface{}, amount float64, paymentType string) error {
	date, ok := dateValue.(*time.Time)
	if !ok || date == nil {
		return nil
	}

	count, err := tx.GetValue(`SELECT COUNT(*) FROM project_payment
		WHERE company_id=? AND project_id=? AND payment_date=? AND amount=? AND payment_type=?
		AND deleted_at IS NULL AND created_at >= DATE_SUB(NOW(), INTERVAL 2 MINUTE)`,
		currentCompanyId(ctx), projectId, date.Format("2006-01-02"), amount, strings.TrimSpace(paymentType))
	if err != nil {
		glog.Errorf(ctx, "校验重复收款失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验收款记录失败，请稍后重试")
	}

	if count.Int() > 0 {
		return gerror.New("检测到 2 分钟内已保存过完全相同的收款记录，请刷新列表确认是否提交成功，避免重复登记")
	}

	return nil
}

// ensureContractNotExceeded 校验「本次收款后累计收款」不超过项目结算金额。
//
// 结算金额 = 合同金额 + 变更净额（加项 - 减项），与「项目毛利」报表口径保持一致；
// 结算金额为 0（未填写合同且无变更）时不做限制，避免影响还没维护金额的项目。
// excludeId 用于修改场景排除当前记录自身。
func (p *Payment) ensureContractNotExceeded(ctx context.Context, tx gdb.TX, projectId int64, amount float64, excludeId int64) error {
	row, err := tx.GetOne("SELECT amount FROM project WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", projectId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "读取项目合同金额失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验合同金额失败，请稍后重试")
	}
	if row == nil || row.IsEmpty() {
		return gerror.New("所选项目不存在或已被删除")
	}

	contract := roundAmount(row["amount"].Float64())

	changeValue, err := tx.GetValue(`SELECT COALESCE(SUM(CASE WHEN change_type='ADD' THEN amount ELSE -amount END),0)
		FROM project_change WHERE project_id=? AND company_id=? AND deleted_at IS NULL`,
		projectId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "统计项目变更净额失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验项目变更金额失败，请稍后重试")
	}
	changeNet := roundAmount(changeValue.Float64())

	settlement := roundAmount(contract + changeNet)
	if settlement <= 0 {
		return nil
	}

	value, err := tx.GetValue(`SELECT COALESCE(SUM(amount),0) FROM project_payment
		WHERE project_id=? AND company_id=? AND deleted_at IS NULL AND id<>?`, projectId, currentCompanyId(ctx), excludeId)
	if err != nil {
		glog.Errorf(ctx, "统计项目累计收款失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验累计收款失败，请稍后重试")
	}

	received := roundAmount(value.Float64())
	total := roundAmount(received + amount)
	if total > settlement {
		return gerror.Newf("本次收款后该项目累计收款 ¥%.2f，将超过结算金额 ¥%.2f（合同 ¥%.2f + 变更 ¥%.2f，当前已收 ¥%.2f），请核对金额",
			total, settlement, contract, changeNet, received)
	}

	return nil
}

// ensureProjectExists 校验所属项目存在，避免挂到已删除 / 不存在的项目上。
func (p *Payment) ensureProjectExists(ctx context.Context, projectId int64) error {
	if projectId <= 0 {
		return gerror.New("请选择所属项目")
	}

	count, err := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Project.Columns().Id, projectId).Count()
	if err != nil {
		glog.Errorf(ctx, "校验收款记录所属项目失败, projectId=%d: %v", projectId, err)
		return gerror.New("保存收款记录失败，请稍后重试")
	}

	if count == 0 {
		return gerror.New("所选项目不存在或已被删除")
	}

	return nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

// ===================== 操作日志 =====================

// receiptSnapshot 收款单快照，删除时用于还原日志所需信息。
type receiptSnapshot struct {
	Form      *Form
	VoucherNo string
}

// loadReceiptSnapshot 删除前读取收款单快照，读不到时返回 nil（不影响删除）。
func (p *Payment) loadReceiptSnapshot(ctx context.Context, id int64) *receiptSnapshot {
	record, err := defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectPayment.Columns().Id, id).
		One()
	if err != nil || record.IsEmpty() {
		return nil
	}

	return &receiptSnapshot{
		Form: &Form{
			ProjectId:     record["project_id"].Int64(),
			Amount:        roundAmount(record["amount"].Float64()),
			PaymentType:   record["payment_type"].String(),
			PaymentDate:   record["payment_date"].String(),
			FundAccountId: record["fund_account_id"].Int64(),
			Remark:        record["remark"].String(),
		},
		VoucherNo: p.voucherNoText(ctx, record["voucher_id"].Int64()),
	}
}

// voucherNoText 凭证号，未记账或已作废时返回空串。
func (p *Payment) voucherNoText(ctx context.Context, voucherId int64) string {
	if voucherId <= 0 {
		return ""
	}

	values, err := booking.New().VoucherNos(ctx, []int64{voucherId})
	if err != nil {
		return ""
	}

	return values[voucherId]
}

// logReceipt 记录收款操作日志（封装：新增 / 修改 / 删除各调用一次即可）。
func (p *Payment) logReceipt(ctx context.Context, action string, id int64, form *Form, voucherNo string) {
	if form == nil {
		return
	}

	projectName := p.projectNameText(ctx, form.ProjectId)
	label := oplog.ActionLabels[action]
	if label == "" {
		label = action
	}

	oplog.New().Record(ctx, &oplog.Entry{
		Module:     oplog.ModuleReceipt,
		Action:     action,
		TargetType: "payment",
		TargetId:   id,
		Title:      fmt.Sprintf("%s收款：%s（%s ￥%.2f）", label, projectName, form.PaymentType, form.Amount),
		Amount:     form.Amount,
		Detail: g.Map{
			"paymentId":   id,
			"projectId":   form.ProjectId,
			"projectName": projectName,
			"amount":      roundAmount(form.Amount),
			"paymentType": form.PaymentType,
			"paymentDate": form.PaymentDate,
			"voucherNo":   voucherNo,
			"remark":      form.Remark,
		},
	})
}

// projectNameText 项目名称，取不到时给出占位，避免日志出现空括号。
func (p *Payment) projectNameText(ctx context.Context, projectId int64) string {
	if projectId <= 0 {
		return "公司公共支出"
	}

	value, err := g.DB().GetValue(ctx, "SELECT name FROM project WHERE id=? AND deleted_at IS NULL", projectId)
	if err != nil || value.IsEmpty() {
		return fmt.Sprintf("项目 #%d", projectId)
	}

	return value.String()
}

// parseDate 解析 YYYY-MM-DD 日期字符串，兼容常见的斜杠 / 带时间格式。
func parseDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, gerror.New("请选择收款日期")
	}

	layouts := []string{"2006-01-02", "2006/01/02", "2006-01-02 15:04:05"}

	for _, layout := range layouts {
		if date, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &date, nil
		}
	}

	return nil, gerror.New("收款日期格式不正确")
}

// roundAmount 金额保留两位小数。
func roundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
