package payout

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

// PayoutInfo 付款记录 + 资金账户 / 凭证展示信息。
type PayoutInfo struct {
	Payout          *defaultEntity.ProjectPayout
	VendorName      string
	FundAccountName string
	VoucherNo       string
}

// DetailResult 付款记录详情：记录 + 资金账户 + 记账信息。
type DetailResult struct {
	Payout          *defaultEntity.ProjectPayout
	VendorName      string
	FundAccountName string
	VoucherNo       string
	Entries         []*booking.EntryView
}

// Form 付款记录新增 / 编辑入参。
type Form struct {
	Amount        float64
	PayoutType    string
	PayoutDate    string
	VendorName    string
	FundAccountId int64
	Remark        string
}

type Payout struct{}

func New() *Payout { return &Payout{} }

// List 付款记录分页列表。
func (p *Payout) List(
	ctx context.Context,
	keyword string,
	payoutType string,
	current int,
	size int,
) (list []*PayoutInfo, total int, err error) {
	cols := defaultDao.ProjectPayout.Columns()
	model := defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL")

	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		model = model.WhereLike(cols.PayoutType, like).
			WhereOrLike(cols.VendorName, like).
			WhereOrLike(cols.Remark, like)
	}
	if payoutType = strings.TrimSpace(payoutType); payoutType != "" {
		model = model.Where(cols.PayoutType, payoutType)
	}

	var records []*defaultEntity.ProjectPayout
	err = model.Page(current, size).
		OrderDesc(cols.PayoutDate).
		OrderDesc(cols.Id).
		ScanAndCount(&records, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询付款记录分页列表失败: %v", err)
		return nil, 0, gerror.New("查询付款记录失败，请稍后重试")
	}

	return p.attach(ctx, records), total, nil
}

// Stats 付款概览统计：笔数 + 付款总额 + 本月付款 + 本月笔数。
func (p *Payout) Stats(ctx context.Context) (totalCount int, totalAmount float64, monthAmount float64, monthCount int, err error) {
	cols := defaultDao.ProjectPayout.Columns()

	totalCount, err = defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		glog.Errorf(ctx, "统计付款笔数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	totalAmount, err = defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计付款总额失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	monthAmount, err = defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		WhereGTE(cols.PayoutDate, monthStart.Format("2006-01-02")).
		WhereLT(cols.PayoutDate, nextMonthStart.Format("2006-01-02")).
		Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计本月付款失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	monthCount, err = defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where("deleted_at IS NULL").
		WhereGTE(cols.PayoutDate, monthStart.Format("2006-01-02")).
		WhereLT(cols.PayoutDate, nextMonthStart.Format("2006-01-02")).
		Count()
	if err != nil {
		glog.Errorf(ctx, "统计本月付款笔数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	return totalCount, roundAmount(totalAmount), roundAmount(monthAmount), monthCount, nil
}

// Detail 付款记录详情。
func (p *Payout) Detail(ctx context.Context, id int64) (*DetailResult, error) {
	record, err := p.getPayout(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DetailResult{
		Payout:          record,
		VendorName:      record.VendorName,
		FundAccountName: p.fundAccountName(ctx, record.FundAccountId),
		VoucherNo:       p.voucherNo(ctx, record.VoucherId),
		Entries:         p.ledgerEntries(ctx, record.VoucherId),
	}, nil
}

// Create 新增付款记录并自动生成付款凭证。
func (p *Payout) Create(ctx context.Context, form *Form, userId uint) (id int64, voucherNo string, err error) {
	if err = form.validate(); err != nil {
		return 0, "", err
	}

	data, err := form.toData()
	if err != nil {
		return 0, "", err
	}

	err = defaultDao.ProjectPayout.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		data["company_id"] = currentCompanyId(ctx)
		result, insertErr := defaultDao.ProjectPayout.Ctx(ctx).TX(tx).Data(data).Insert()
		if insertErr != nil {
			glog.Errorf(ctx, "新增付款记录失败: %v", insertErr)
			return gerror.New("新增付款记录失败，请稍后重试")
		}

		id, insertErr = result.LastInsertId()
		if insertErr != nil {
			glog.Errorf(ctx, "读取新增付款记录ID失败: %v", insertErr)
			return gerror.New("新增付款记录失败，请稍后重试")
		}

		voucherNo, insertErr = p.book(ctx, tx, id, form, userId)
		return insertErr
	})
	if err != nil {
		return 0, "", err
	}

	p.logPayout(ctx, oplog.ActionCreate, id, form, voucherNo)
	return id, voucherNo, nil
}

// Update 更新付款记录，并同步重记凭证。
func (p *Payout) Update(ctx context.Context, id int64, form *Form, userId uint) (voucherNo string, err error) {
	if err = form.validate(); err != nil {
		return "", err
	}

	data, err := form.toData()
	if err != nil {
		return "", err
	}
	data["company_id"] = currentCompanyId(ctx)

	err = defaultDao.ProjectPayout.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		oldVoucherId, lockErr := p.lockVoucherId(ctx, tx, id, true)
		if lockErr != nil {
			return lockErr
		}

		if oldVoucherId > 0 {
			if lockErr = booking.New().Void(ctx, tx, oldVoucherId); lockErr != nil {
				return lockErr
			}
			if _, lockErr = defaultDao.ProjectPayout.Ctx(ctx).TX(tx).
				Where("company_id", currentCompanyId(ctx)).
				Where(defaultDao.ProjectPayout.Columns().Id, id).
				Data(g.Map{defaultDao.ProjectPayout.Columns().VoucherId: 0}).
				Update(); lockErr != nil {
				glog.Errorf(ctx, "重置付款记录凭证ID失败, id=%d: %v", id, lockErr)
				return gerror.New("更新付款记录失败，请稍后重试")
			}
		}

		if _, lockErr = defaultDao.ProjectPayout.Ctx(ctx).TX(tx).
			Where("company_id", currentCompanyId(ctx)).
			Where(defaultDao.ProjectPayout.Columns().Id, id).
			Data(data).
			Update(); lockErr != nil {
			glog.Errorf(ctx, "更新付款记录失败, id=%d: %v", id, lockErr)
			return gerror.New("更新付款记录失败，请稍后重试")
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

	p.logPayout(ctx, oplog.ActionUpdate, id, form, voucherNo)
	return voucherNo, nil
}

// Delete 删除单条付款记录：先作废凭证，再软删除业务数据。
func (p *Payout) Delete(ctx context.Context, id int64) error {
	snapshot := p.loadPayoutSnapshot(ctx, id)

	err := defaultDao.ProjectPayout.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		voucherId, lockErr := p.lockVoucherId(ctx, tx, id, true)
		if lockErr != nil {
			return lockErr
		}
		if lockErr = booking.New().Void(ctx, tx, voucherId); lockErr != nil {
			return lockErr
		}
		if _, lockErr = defaultDao.ProjectPayout.Ctx(ctx).TX(tx).
			Where("company_id", currentCompanyId(ctx)).
			Where(defaultDao.ProjectPayout.Columns().Id, id).
			Delete(); lockErr != nil {
			glog.Errorf(ctx, "删除付款记录失败, id=%d: %v", id, lockErr)
			return gerror.New("删除付款记录失败，请稍后重试")
		}
		return nil
	})
	if err != nil {
		return err
	}

	if snapshot != nil {
		p.logPayout(ctx, oplog.ActionDelete, id, snapshot.Form, snapshot.VoucherNo)
	}
	return nil
}

// BatchDelete 批量删除付款记录。
func (p *Payout) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return gerror.New("请选择要删除的付款记录")
	}

	snapshots := make(map[int64]*payoutSnapshot, len(ids))
	for _, id := range ids {
		if snapshot := p.loadPayoutSnapshot(ctx, id); snapshot != nil {
			snapshots[id] = snapshot
		}
	}

	err := defaultDao.ProjectPayout.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, id := range ids {
			voucherId, lockErr := p.lockVoucherId(ctx, tx, id, false)
			if lockErr != nil {
				return lockErr
			}
			if lockErr = booking.New().Void(ctx, tx, voucherId); lockErr != nil {
				return lockErr
			}
		}
		if _, lockErr := defaultDao.ProjectPayout.Ctx(ctx).TX(tx).
			Where("company_id", currentCompanyId(ctx)).
			WhereIn(defaultDao.ProjectPayout.Columns().Id, ids).
			Delete(); lockErr != nil {
			glog.Errorf(ctx, "批量删除付款记录失败, ids=%v: %v", ids, lockErr)
			return gerror.New("删除付款记录失败，请稍后重试")
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, id := range ids {
		if snapshot, ok := snapshots[id]; ok {
			p.logPayout(ctx, oplog.ActionDelete, id, snapshot.Form, snapshot.VoucherNo)
		}
	}
	return nil
}

// book 为付款记录生成借贷凭证，并把凭证 ID 写回业务表。
func (p *Payout) book(ctx context.Context, tx gdb.TX, id int64, form *Form, userId uint) (string, error) {
	initialized, err := booking.New().IsInitialized(ctx, tx)
	if err != nil {
		return "", err
	}
	if !initialized {
		return "", nil
	}
	if form.FundAccountId <= 0 {
		return "", gerror.New("请选择资金账户，系统需要据此记录这笔钱从哪个账户支出")
	}

	result, err := booking.New().PostPayout(ctx, tx, &booking.Payout{
		SourceId:      id,
		PayoutType:    form.PayoutType,
		Payee:         form.VendorName,
		Date:          form.PayoutDate,
		Amount:        form.Amount,
		FundAccountId: form.FundAccountId,
		CreatedBy:     userId,
	})
	if err != nil {
		return "", err
	}

	if _, err = defaultDao.ProjectPayout.Ctx(ctx).TX(tx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectPayout.Columns().Id, id).
		Data(g.Map{defaultDao.ProjectPayout.Columns().VoucherId: result.VoucherId}).
		Update(); err != nil {
		glog.Errorf(ctx, "写回付款记录凭证ID失败, id=%d: %v", id, err)
		return "", gerror.New("保存付款记录失败，请稍后重试")
	}

	return result.VoucherNo, nil
}

// validate 业务校验：金额必须大于 0，费用类型与收款方必填。
func (f *Form) validate() error {
	f.Amount = roundAmount(f.Amount)
	if f.Amount <= 0 {
		return gerror.New("付款金额必须大于 0")
	}
	if strings.TrimSpace(f.PayoutType) == "" {
		return gerror.New("请选择费用类型")
	}
	if strings.TrimSpace(f.VendorName) == "" {
		return gerror.New("请填写收款方")
	}
	return nil
}

// toData 转换为可落库的列值映射。
func (f *Form) toData() (g.Map, error) {
	cols := defaultDao.ProjectPayout.Columns()
	payoutDate, err := parseDate(f.PayoutDate)
	if err != nil {
		return nil, err
	}
	return g.Map{
		cols.Amount:        f.Amount,
		cols.PayoutType:    f.PayoutType,
		cols.PayoutDate:    payoutDate,
		cols.VendorName:    strings.TrimSpace(f.VendorName),
		cols.FundAccountId: f.FundAccountId,
		cols.Remark:        f.Remark,
	}, nil
}

// attach 为付款记录批量补充资金账户名与凭证号。
func (p *Payout) attach(ctx context.Context, records []*defaultEntity.ProjectPayout) []*PayoutInfo {
	result := make([]*PayoutInfo, 0, len(records))
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
			glog.Errorf(ctx, "查询付款凭证号失败, ids=%v: %v", voucherIds, err)
		} else {
			voucherNos = values
		}
	}

	for _, record := range records {
		result = append(result, &PayoutInfo{
			Payout:          record,
			VendorName:      record.VendorName,
			FundAccountName: fundAccountNames[record.FundAccountId],
			VoucherNo:       voucherNos[record.VoucherId],
		})
	}
	return result
}

func (p *Payout) fundAccountNames(ctx context.Context, ids []int64) map[int64]string {
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
		glog.Errorf(ctx, "查询付款资金账户失败, ids=%v: %v", ids, err)
		return names
	}
	for _, record := range records {
		names[record["id"].Int64()] = record["account_name"].String()
	}
	return names
}

func (p *Payout) fundAccountName(ctx context.Context, id int64) string {
	if id <= 0 {
		return ""
	}
	value, err := g.DB().GetValue(ctx,
		"SELECT account_name FROM fin_fund_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "查询付款资金账户失败, id=%d: %v", id, err)
		return ""
	}
	return value.String()
}

func (p *Payout) voucherNo(ctx context.Context, voucherId int64) string {
	if voucherId <= 0 {
		return ""
	}
	values, err := booking.New().VoucherNos(ctx, []int64{voucherId})
	if err != nil {
		glog.Errorf(ctx, "查询付款凭证号失败, voucherId=%d: %v", voucherId, err)
		return ""
	}
	return values[voucherId]
}

func (p *Payout) ledgerEntries(ctx context.Context, voucherId int64) []*booking.EntryView {
	entries, err := booking.New().Entries(ctx, voucherId)
	if err != nil {
		glog.Errorf(ctx, "查询付款凭证分录失败, voucherId=%d: %v", voucherId, err)
		return []*booking.EntryView{}
	}
	return entries
}

func (p *Payout) getPayout(ctx context.Context, id int64) (*defaultEntity.ProjectPayout, error) {
	var record defaultEntity.ProjectPayout
	err := defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectPayout.Columns().Id, id).
		Where("deleted_at IS NULL").
		Scan(&record)
	if err != nil {
		glog.Errorf(ctx, "查询付款记录失败, id=%d: %v", id, err)
		return nil, gerror.New("获取付款记录失败，请稍后重试")
	}
	if record.Id == 0 {
		return nil, gerror.New("付款记录不存在或已被删除")
	}
	return &record, nil
}

func (p *Payout) lockVoucherId(ctx context.Context, tx gdb.TX, id int64, mustExist bool) (int64, error) {
	row, err := tx.GetOne(
		"SELECT id,voucher_id FROM payout WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE",
		id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "锁定付款记录失败, id=%d: %v", id, err)
		return 0, gerror.New("操作付款记录失败，请稍后重试")
	}
	if row == nil || row.IsEmpty() {
		if mustExist {
			return 0, gerror.New("付款记录不存在或已被删除")
		}
		return 0, nil
	}
	return row["voucher_id"].Int64(), nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

// ===================== 操作日志 =====================

type payoutSnapshot struct {
	Form      *Form
	VoucherNo string
}

func (p *Payout) loadPayoutSnapshot(ctx context.Context, id int64) *payoutSnapshot {
	record, err := defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.ProjectPayout.Columns().Id, id).
		One()
	if err != nil || record.IsEmpty() {
		return nil
	}
	return &payoutSnapshot{
		Form: &Form{
			Amount:        roundAmount(record["amount"].Float64()),
			PayoutType:    record["payout_type"].String(),
			PayoutDate:    record["payout_date"].String(),
			VendorName:    record["vendor_name"].String(),
			FundAccountId: record["fund_account_id"].Int64(),
			Remark:        record["remark"].String(),
		},
		VoucherNo: p.voucherNoText(ctx, record["voucher_id"].Int64()),
	}
}

func (p *Payout) voucherNoText(ctx context.Context, voucherId int64) string {
	if voucherId <= 0 {
		return ""
	}
	values, err := booking.New().VoucherNos(ctx, []int64{voucherId})
	if err != nil {
		return ""
	}
	return values[voucherId]
}

func (p *Payout) logPayout(ctx context.Context, action string, id int64, form *Form, voucherNo string) {
	if form == nil {
		return
	}
	label := oplog.ActionLabels[action]
	if label == "" {
		label = action
	}
	oplog.New().Record(ctx, &oplog.Entry{
		Module:     oplog.ModulePayout,
		Action:     action,
		TargetType: "payout",
		TargetId:   id,
		Title:      fmt.Sprintf("%s付款：%s ¥%.2f", label, form.PayoutType, form.Amount),
		Amount:     form.Amount,
		Detail: g.Map{
			"payoutId":   id,
			"amount":     roundAmount(form.Amount),
			"payoutType": form.PayoutType,
			"payoutDate": form.PayoutDate,
			"vendorName": form.VendorName,
			"voucherNo":  voucherNo,
			"remark":     form.Remark,
		},
	})
}

func parseDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, gerror.New("请选择付款日期")
	}
	layouts := []string{"2006-01-02", "2006/01/02", "2006-01-02 15:04:05"}
	for _, layout := range layouts {
		if date, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &date, nil
		}
	}
	return nil, gerror.New("付款日期格式不正确")
}

func roundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
