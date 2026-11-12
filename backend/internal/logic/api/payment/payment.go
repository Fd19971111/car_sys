package payment

import (
	"context"
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

// PaymentInfo 收款记录 + 资金账户 / 凭证展示信息。
type PaymentInfo struct {
	Payment         *defaultEntity.ProjectPayment
	FundAccountName string
	VoucherNo       string
}

// DetailResult 收款记录详情：记录 + 资金账户 + 记账信息。
type DetailResult struct {
	Payment         *defaultEntity.ProjectPayment
	FundAccountName string
	VoucherNo       string
	Entries         []*booking.EntryView
}

// Form 收款记录新增 / 编辑入参。
type Form struct {
	RideOrderId   int64
	Amount        float64
	PaymentType   string
	PaymentDate   string
	FundAccountId int64
	Remark        string
}

type Payment struct{}

func New() *Payment { return &Payment{} }

// List 收款记录分页列表。
func (p *Payment) List(
	ctx context.Context,
	keyword string,
	paymentType string,
	current int,
	size int,
) (list []*PaymentInfo, total int, err error) {
	cols := defaultDao.ProjectPayment.Columns()
	model := defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL")

	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		model = model.WhereLike(cols.PaymentType, like).WhereOrLike(cols.Remark, like)
	}
	if paymentType = strings.TrimSpace(paymentType); paymentType != "" {
		model = model.Where(cols.PaymentType, paymentType)
	}

	var records []*defaultEntity.ProjectPayment
	err = model.Page(current, size).
		OrderDesc(cols.PaymentDate).
		OrderDesc(cols.Id).
		ScanAndCount(&records, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询收款记录分页列表失败: %v", err)
		return nil, 0, gerror.New("查询收款记录失败，请稍后重试")
	}

	return p.attach(ctx, records), total, nil
}

// Stats 收款概览统计：笔数 + 收款总额 + 本月收款 + 本月笔数。
func (p *Payment) Stats(ctx context.Context) (totalCount int, totalAmount float64, monthAmount float64, monthCount int, err error) {
	cols := defaultDao.ProjectPayment.Columns()

	totalCount, err = defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		glog.Errorf(ctx, "统计收款笔数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	totalAmount, err = defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计收款总额失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	monthAmount, err = defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		WhereGTE(cols.PaymentDate, monthStart.Format("2006-01-02")).
		WhereLT(cols.PaymentDate, nextMonthStart.Format("2006-01-02")).
		Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计本月收款失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	monthCount, err = defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		WhereGTE(cols.PaymentDate, monthStart.Format("2006-01-02")).
		WhereLT(cols.PaymentDate, nextMonthStart.Format("2006-01-02")).
		Count()
	if err != nil {
		glog.Errorf(ctx, "统计本月收款笔数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取收款统计失败，请稍后重试")
	}

	return totalCount, roundAmount(totalAmount), roundAmount(monthAmount), monthCount, nil
}

// Detail 收款记录详情。
func (p *Payment) Detail(ctx context.Context, id int64) (*DetailResult, error) {
	record, err := p.getPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DetailResult{
		Payment:         record,
		FundAccountName: p.fundAccountName(ctx, record.FundAccountId),
		VoucherNo:       p.voucherNo(ctx, record.VoucherId),
		Entries:         p.ledgerEntries(ctx, record.VoucherId),
	}, nil
}

// Create 新增收款记录并自动生成收款凭证。
func (p *Payment) Create(ctx context.Context, form *Form, userId uint) (id int64, voucherNo string, err error) {
	if err = form.validate(); err != nil {
		return 0, "", err
	}

	data, err := form.toData()
	if err != nil {
		return 0, "", err
	}

	err = defaultDao.ProjectPayment.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if orderErr := p.ensureRideOrderReceivable(ctx, tx, form, 0); orderErr != nil {
			return orderErr
		}
		if dupErr := p.ensureNotDuplicate(ctx, tx, data[defaultDao.ProjectPayment.Columns().PaymentDate], form.Amount, form.PaymentType); dupErr != nil {
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

// Update 更新收款记录，并同步重记凭证。
func (p *Payment) Update(ctx context.Context, id int64, form *Form, userId uint) (voucherNo string, err error) {
	if err = form.validate(); err != nil {
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
		if lockErr = p.ensureRideOrderReceivable(ctx, tx, form, id); lockErr != nil {
			return lockErr
		}

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

func (p *Payment) ensureRideOrderReceivable(ctx context.Context, tx gdb.TX, form *Form, excludePaymentId int64) error {
	if form.RideOrderId < 0 {
		return gerror.New("行程订单ID不能小于0")
	}
	if form.RideOrderId == 0 {
		return nil
	}
	order, err := tx.GetOne("SELECT final_fare,status FROM ride_order WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", form.RideOrderId, currentCompanyId(ctx))
	if err != nil {
		return gerror.Wrap(err, "校验行程订单失败")
	}
	if order.IsEmpty() || order["status"].String() != "completed" {
		return gerror.New("收款只能关联当前公司的已完成行程")
	}
	paid, err := tx.GetValue("SELECT COALESCE(SUM(amount),0) FROM payment WHERE company_id=? AND ride_order_id=? AND deleted_at IS NULL AND id<>?", currentCompanyId(ctx), form.RideOrderId, excludePaymentId)
	if err != nil {
		return gerror.Wrap(err, "校验行程已收金额失败")
	}
	if paid.Float64()+form.Amount > order["final_fare"].Float64()+0.005 {
		return gerror.New("关联收款金额超过该行程剩余应收车费")
	}
	return nil
}

// Delete 删除单条收款记录：先作废凭证，再软删除业务数据。
func (p *Payment) Delete(ctx context.Context, id int64) error {
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

// BatchDelete 批量删除收款记录。
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

// book 为收款记录生成借贷凭证，并把凭证 ID 写回业务表。
func (p *Payment) book(ctx context.Context, tx gdb.TX, id int64, form *Form, userId uint) (string, error) {
	initialized, err := booking.New().IsInitialized(ctx, tx)
	if err != nil {
		return "", err
	}
	if !initialized {
		return "", nil
	}
	if form.FundAccountId <= 0 {
		return "", gerror.New("请选择资金账户，系统需要据此记录这笔钱进入了哪个账户")
	}

	result, err := booking.New().PostReceipt(ctx, tx, &booking.Receipt{
		SourceId:      id,
		RideOrderId:   form.RideOrderId,
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

// validate 业务校验：金额必须大于 0，款项类型必填。
func (f *Form) validate() error {
	f.Amount = roundAmount(f.Amount)
	if f.Amount <= 0 {
		return gerror.New("收款金额必须大于 0")
	}
	if strings.TrimSpace(f.PaymentType) == "" {
		return gerror.New("请选择款项类型")
	}
	return nil
}

// toData 转换为可落库的列值映射。
func (f *Form) toData() (g.Map, error) {
	cols := defaultDao.ProjectPayment.Columns()
	paymentDate, err := parseDate(f.PaymentDate)
	if err != nil {
		return nil, err
	}
	return g.Map{
		cols.RideOrderId:   f.RideOrderId,
		cols.Amount:        f.Amount,
		cols.PaymentType:   f.PaymentType,
		cols.PaymentDate:   paymentDate,
		cols.FundAccountId: f.FundAccountId,
		cols.Remark:        f.Remark,
	}, nil
}

// attach 为收款记录批量补充资金账户名与凭证号。
func (p *Payment) attach(ctx context.Context, records []*defaultEntity.ProjectPayment) []*PaymentInfo {
	result := make([]*PaymentInfo, 0, len(records))
	if len(records) == 0 {
		return result
	}

	fundAccountIds := make([]int64, 0, len(records))
	fundAccountSet := make(map[int64]struct{})
	voucherIds := make([]int64, 0, len(records))

	for _, record := range records {
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

	fundAccountNames := p.fundAccountNames(ctx, fundAccountIds)
	voucherNos := make(map[int64]string, len(voucherIds))
	if len(voucherIds) > 0 {
		if values, err := booking.New().VoucherNos(ctx, voucherIds); err != nil {
			glog.Errorf(ctx, "查询收款凭证号失败, ids=%v: %v", voucherIds, err)
		} else {
			voucherNos = values
		}
	}

	for _, record := range records {
		result = append(result, &PaymentInfo{
			Payment:         record,
			FundAccountName: fundAccountNames[record.FundAccountId],
			VoucherNo:       voucherNos[record.VoucherId],
		})
	}
	return result
}

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

func (p *Payment) fundAccountName(ctx context.Context, id int64) string {
	if id <= 0 {
		return ""
	}
	value, err := g.DB().GetValue(ctx,
		"SELECT account_name FROM fin_fund_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "查询收款资金账户失败, id=%d: %v", id, err)
		return ""
	}
	return value.String()
}

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

func (p *Payment) ledgerEntries(ctx context.Context, voucherId int64) []*booking.EntryView {
	entries, err := booking.New().Entries(ctx, voucherId)
	if err != nil {
		glog.Errorf(ctx, "查询收款凭证分录失败, voucherId=%d: %v", voucherId, err)
		return []*booking.EntryView{}
	}
	return entries
}

func (p *Payment) getPayment(ctx context.Context, id int64) (*defaultEntity.ProjectPayment, error) {
	var record defaultEntity.ProjectPayment
	err := defaultDao.ProjectPayment.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectPayment.Columns().Id, id).
		Where("deleted_at IS NULL").
		Scan(&record)
	if err != nil {
		glog.Errorf(ctx, "查询收款记录失败, id=%d: %v", id, err)
		return nil, gerror.New("获取收款记录失败，请稍后重试")
	}
	if record.Id == 0 {
		return nil, gerror.New("收款记录不存在或已被删除")
	}
	return &record, nil
}

func (p *Payment) lockVoucherId(ctx context.Context, tx gdb.TX, id int64, mustExist bool) (int64, error) {
	row, err := tx.GetOne(
		"SELECT id,voucher_id FROM payment WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE",
		id, currentCompanyId(ctx))
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

func (p *Payment) ensureNotDuplicate(ctx context.Context, tx gdb.TX, dateValue interface{}, amount float64, paymentType string) error {
	date, ok := dateValue.(*time.Time)
	if !ok || date == nil {
		return nil
	}
	count, err := tx.GetValue(`SELECT COUNT(*) FROM payment
		WHERE company_id=? AND payment_date=? AND amount=? AND payment_type=?
		AND deleted_at IS NULL AND created_at >= DATE_SUB(NOW(), INTERVAL 2 MINUTE)`,
		currentCompanyId(ctx), date.Format("2006-01-02"), amount, strings.TrimSpace(paymentType))
	if err != nil {
		glog.Errorf(ctx, "校验重复收款失败: %v", err)
		return gerror.New("校验收款记录失败，请稍后重试")
	}
	if count.Int() > 0 {
		return gerror.New("检测到 2 分钟内已保存过完全相同的收款记录，请刷新列表确认是否提交成功，避免重复登记")
	}
	return nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

// ===================== 操作日志 =====================

type receiptSnapshot struct {
	Form      *Form
	VoucherNo string
}

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
			Amount:        roundAmount(record["amount"].Float64()),
			PaymentType:   record["payment_type"].String(),
			PaymentDate:   record["payment_date"].String(),
			FundAccountId: record["fund_account_id"].Int64(),
			Remark:        record["remark"].String(),
		},
		VoucherNo: p.voucherNoText(ctx, record["voucher_id"].Int64()),
	}
}

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

func (p *Payment) logReceipt(ctx context.Context, action string, id int64, form *Form, voucherNo string) {
	if form == nil {
		return
	}
	label := oplog.ActionLabels[action]
	if label == "" {
		label = action
	}
	oplog.New().Record(ctx, &oplog.Entry{
		Module:     oplog.ModuleReceipt,
		Action:     action,
		TargetType: "payment",
		TargetId:   id,
		Title:      fmt.Sprintf("%s收款：%s ¥%.2f", label, form.PaymentType, form.Amount),
		Amount:     form.Amount,
		Detail: g.Map{
			"paymentId":   id,
			"amount":      roundAmount(form.Amount),
			"paymentType": form.PaymentType,
			"paymentDate": form.PaymentDate,
			"voucherNo":   voucherNo,
			"remark":      form.Remark,
		},
	})
}

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

func roundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
