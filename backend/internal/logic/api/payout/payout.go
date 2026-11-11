package payout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"engineServ/internal/consts"
	projectDataLogic "engineServ/internal/logic/api/project_data"
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

// PayoutInfo 付款记录本体 + 所属项目 / 客户 / 资金账户 / 凭证展示信息
type PayoutInfo struct {
	Payout          *defaultEntity.ProjectPayout
	ProjectName     string
	CustomerName    string
	VendorName      string
	FundAccountName string
	VoucherNo       string
}

// DetailResult 付款记录详情：记录 + 所属项目 + 客户 + 该项目累计已付款 + 记账信息
type DetailResult struct {
	Payout          *defaultEntity.ProjectPayout
	Project         *defaultEntity.Project
	Customer        *defaultEntity.Customer
	VendorName      string
	PaidAmount      float64
	FundAccountName string
	VoucherNo       string
	Entries         []*booking.EntryView
}

// Form 付款记录新增 / 编辑入参
type Form struct {
	ProjectId       int64
	Amount          float64
	PayoutType      string
	PayoutDate      string
	VendorId        int64
	PurchaseOrderId int64
	VendorName      string
	FundAccountId   int64
	Remark          string
}

type Payout struct{}

func New() *Payout {
	return &Payout{}
}

// List 付款记录分页列表：关键词（项目名称 / 客户姓名 / 收款方 / 费用类型 / 备注）+ 项目 + 费用类型过滤。
func (p *Payout) List(
	ctx context.Context,
	keyword string,
	projectId *int64,
	payoutType string,
	source string,
	current int,
	size int,
) (list []*PayoutInfo, total int, err error) {
	cols := defaultDao.ProjectPayout.Columns()
	model := defaultDao.ProjectPayout.Ctx(ctx).Where("company_id", currentCompanyId(ctx))

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		matchedIds, matchErr := p.matchPayoutIds(ctx, keyword)
		if matchErr != nil {
			return nil, 0, matchErr
		}

		if len(matchedIds) == 0 {
			return []*PayoutInfo{}, 0, nil
		}

		model = model.WhereIn(cols.Id, matchedIds)
	}

	if projectId != nil && *projectId > 0 {
		model = model.Where(cols.ProjectId, *projectId)
	}

	payoutType = strings.TrimSpace(payoutType)
	if payoutType != "" {
		model = model.Where(cols.PayoutType, payoutType)
	}

	source = strings.TrimSpace(source)
	if source == "purchase" {
		model = model.Where("EXISTS (SELECT 1 FROM project_purchase_order po WHERE po.payout_id=project_payout.id AND po.company_id=? AND po.deleted_at IS NULL)", currentCompanyId(ctx))
	} else if source == "manual" {
		model = model.Where("NOT EXISTS (SELECT 1 FROM project_purchase_order po WHERE po.payout_id=project_payout.id AND po.company_id=? AND po.deleted_at IS NULL)", currentCompanyId(ctx))
	}

	var records []*defaultEntity.ProjectPayout
	err = model.
		Page(current, size).
		OrderDesc(cols.PayoutDate).
		OrderDesc(cols.Id).
		ScanAndCount(&records, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询付款记录分页列表失败: %v", err)
		return nil, 0, gerror.New("查询付款记录失败，请稍后重试")
	}

	return p.attach(ctx, records), total, nil
}

// Stats 付款概览统计：笔数 + 付款总额 + 本月付款 + 涉及项目数。
func (p *Payout) Stats(ctx context.Context) (totalCount int, totalAmount float64, monthAmount float64, projectCount int, err error) {
	cols := defaultDao.ProjectPayout.Columns()

	totalCount, err = defaultDao.ProjectPayout.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Count()
	if err != nil {
		glog.Errorf(ctx, "统计付款笔数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	totalAmount, err = defaultDao.ProjectPayout.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计付款总额失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	// 本月付款：payout_date ∈ [本月1号, 下月1号)
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	monthAmount, err = defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		WhereGTE(cols.PayoutDate, monthStart.Format("2006-01-02")).
		WhereLT(cols.PayoutDate, nextMonthStart.Format("2006-01-02")).
		Sum(cols.Amount)
	if err != nil {
		glog.Errorf(ctx, "统计本月付款失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	// 公司公共支出（project_id=0）不计入「涉及项目数」
	projectCountValue, err := defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields("COUNT(DISTINCT CASE WHEN project_id>0 THEN project_id END) AS projectCount").
		Value()
	if err != nil {
		glog.Errorf(ctx, "统计付款涉及项目数失败: %v", err)
		return 0, 0, 0, 0, gerror.New("获取付款统计失败，请稍后重试")
	}

	return totalCount, roundAmount(totalAmount), roundAmount(monthAmount), projectCountValue.Int(), nil
}

// Detail 付款记录详情：记录 + 所属项目 + 客户 + 该项目累计已付款 + 凭证分录。
func (p *Payout) Detail(ctx context.Context, id int64) (*DetailResult, error) {
	record, err := p.getPayout(ctx, id)
	if err != nil {
		return nil, err
	}

	project, customer := p.loadProject(ctx, record.ProjectId)

	paidAmount := 0.0
	if record.ProjectId > 0 {
		value, sumErr := defaultDao.ProjectPayout.Ctx(ctx).
			Where("company_id", currentCompanyId(ctx)).
			Where(defaultDao.ProjectPayout.Columns().ProjectId, record.ProjectId).
			Sum(defaultDao.ProjectPayout.Columns().Amount)
		if sumErr != nil {
			// 累计已付仅用于展示，失败不阻断详情
			glog.Errorf(ctx, "统计项目累计付款失败, projectId=%d: %v", record.ProjectId, sumErr)
		}
		paidAmount = value
	}

	result := &DetailResult{
		Payout:     record,
		Project:    project,
		Customer:   customer,
		PaidAmount: roundAmount(paidAmount),
	}
	result.VendorName = p.vendorName(ctx, record.VendorId)

	result.FundAccountName = p.fundAccountName(ctx, record.FundAccountId)
	result.VoucherNo = p.voucherNo(ctx, record.VoucherId)
	result.Entries = p.ledgerEntries(ctx, record.VoucherId)

	return result, nil
}

// Create 新增付款记录并自动生成付款凭证，返回新记录ID与凭证号。
//
// 业务数据与凭证在同一个事务里写入：凭证记账失败（例如会计期间已结账、
// 资金账户不可用）时，付款记录也不会落库，避免出现「有业务没记账」的账。
func (p *Payout) Create(ctx context.Context, form *Form, userId uint) (id int64, voucherNo string, err error) {
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

	err = defaultDao.ProjectPayout.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		vendorName, vendorErr := p.vendorNameTx(ctx, tx, form.VendorId)
		if vendorErr != nil {
			return vendorErr
		}
		form.VendorName = vendorName

		// 超付校验：项目累计付款不能超过合同金额
		if checkErr := p.ensureContractNotExceeded(ctx, tx, form.ProjectId, form.Amount, 0); checkErr != nil {
			return checkErr
		}
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
		if form.PurchaseOrderId > 0 {
			order, orderErr := tx.GetOne("SELECT project_id,material_id,quantity,amount,status,stocked_at FROM project_purchase_order WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", form.PurchaseOrderId, currentCompanyId(ctx))
			if orderErr != nil {
				glog.Errorf(ctx, "读取采购单失败, id=%d: %v", form.PurchaseOrderId, orderErr)
				return gerror.New("关联采购单失败，请稍后重试")
			}
			if order == nil || order.IsEmpty() || order["status"].String() != "CONFIRMED" {
				return gerror.New("采购单不存在或已处理")
			}
			if order["project_id"].Int64() != form.ProjectId || roundAmount(order["amount"].Float64()) != roundAmount(form.Amount) {
				return gerror.New("付款项目或金额与采购单不一致")
			}
			// 采购单已在列表中补做过入库（stocked_at 非空）时，付款不再重复加库存
			if order["stocked_at"].IsEmpty() {
				if orderErr = projectDataLogic.New().AddMaterialStockTx(ctx, tx, form.ProjectId, order["material_id"].Int64(), order["quantity"].Float64()); orderErr != nil {
					return orderErr
				}
			}
			if _, orderErr = tx.Exec("UPDATE project_purchase_order SET status='PAID',payout_id=?,paid_at=NOW(),stocked_at=IFNULL(stocked_at,NOW()) WHERE id=? AND company_id=? AND status='CONFIRMED'", id, form.PurchaseOrderId, currentCompanyId(ctx)); orderErr != nil {
				glog.Errorf(ctx, "更新采购单状态失败, id=%d: %v", form.PurchaseOrderId, orderErr)
				return gerror.New("更新采购单状态失败，请稍后重试")
			}
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

// Update 更新付款记录，并同步更新自动生成的凭证。
//
// 只有原本已经记账的记录才会重新记账：未记账的历史数据保持未记账，
// 这样补录 / 修正很久以前的付款时不会被已结账期间挡住。
func (p *Payout) Update(ctx context.Context, id int64, form *Form, userId uint) (voucherNo string, err error) {
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

	err = defaultDao.ProjectPayout.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		vendorName, vendorErr := p.vendorNameTx(ctx, tx, form.VendorId)
		if vendorErr != nil {
			return vendorErr
		}
		form.VendorName = vendorName

		oldVoucherId, lockErr := p.lockVoucherId(ctx, tx, id, true)
		if lockErr != nil {
			return lockErr
		}

		// 超付校验：修改后的累计付款同样不能超过合同金额（排除本条记录）
		if checkErr := p.ensureContractNotExceeded(ctx, tx, form.ProjectId, form.Amount, id); checkErr != nil {
			return checkErr
		}
		// 先作废旧凭证，撤销它的余额与流水，再按新数据重新记账
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
	// 日志需要原始金额与收款方，删除前先留快照
	snapshot := p.loadPayoutSnapshot(ctx, id)

	err := defaultDao.ProjectPayout.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if lockErr := p.reversePurchaseInbound(ctx, tx, id); lockErr != nil {
			return lockErr
		}
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
		if lockErr = p.resetPurchaseOrderStatus(ctx, tx, id); lockErr != nil {
			return lockErr
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

// BatchDelete 批量删除付款记录：逐条作废凭证后软删除。
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
			if lockErr := p.reversePurchaseInbound(ctx, tx, id); lockErr != nil {
				return lockErr
			}
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
		for _, id := range ids {
			if lockErr := p.resetPurchaseOrderStatus(ctx, tx, id); lockErr != nil {
				return lockErr
			}
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

// resetPurchaseOrderStatus 付款撤销后恢复采购单状态。
// 若库存来自列表中的「采购单入库」补做，则保留 stocked_at，避免单据回到未入库但流水仍存在而无法再次入库。
func (p *Payout) resetPurchaseOrderStatus(ctx context.Context, tx gdb.TX, payoutId int64) error {
	keepStocked, err := p.hasBackfillInbound(ctx, tx, payoutId)
	if err != nil {
		return err
	}

	sqlText := "UPDATE project_purchase_order SET status='CONFIRMED',payout_id=0,paid_at=NULL,stocked_at=NULL WHERE payout_id=? AND company_id=?"
	if keepStocked {
		sqlText = "UPDATE project_purchase_order SET status='CONFIRMED',payout_id=0,paid_at=NULL WHERE payout_id=? AND company_id=?"
	}
	if _, err = tx.Exec(sqlText, payoutId, currentCompanyId(ctx)); err != nil {
		glog.Errorf(ctx, "恢复采购单状态失败, payoutId=%d: %v", payoutId, err)
		return gerror.New("恢复采购单状态失败，请稍后重试")
	}

	return nil
}

// hasBackfillInbound 判断采购单库存是否来自列表补做的入库流水。
func (p *Payout) hasBackfillInbound(ctx context.Context, tx gdb.TX, payoutId int64) (bool, error) {
	count, err := tx.GetValue(`SELECT COUNT(*) FROM project_purchase_order po
		INNER JOIN inventory_record ir ON ir.source_type='PURCHASE' AND ir.source_id=po.id AND ir.deleted_at IS NULL
		WHERE po.payout_id=? AND po.company_id=?`, payoutId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "查询采购单补入库流水失败, payoutId=%d: %v", payoutId, err)
		return false, gerror.New("恢复采购单状态失败，请稍后重试")
	}

	return count.Int() > 0, nil
}

func (p *Payout) reversePurchaseInbound(ctx context.Context, tx gdb.TX, payoutId int64) error {
	order, err := tx.GetOne(`SELECT id,project_id,material_id,quantity,stocked_at
		FROM project_purchase_order WHERE payout_id=? AND company_id=? AND status='PAID' AND deleted_at IS NULL FOR UPDATE`, payoutId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "读取付款关联采购单失败, payoutId=%d: %v", payoutId, err)
		return gerror.New("读取关联采购单失败，请稍后重试")
	}
	if order == nil || order.IsEmpty() {
		return nil
	}
	if order["stocked_at"].IsEmpty() {
		return nil
	}
	// 该采购单库存来自列表中的「采购单入库」补做，付款本身没加库存，撤销付款无需回退
	count, err := tx.GetValue(`SELECT COUNT(*) FROM inventory_record
		WHERE company_id=? AND source_type='PURCHASE' AND source_id=? AND deleted_at IS NULL`, currentCompanyId(ctx), order["id"].Int64())
	if err != nil {
		glog.Errorf(ctx, "读取采购单入库流水失败, orderId=%d: %v", order["id"].Int64(), err)
		return gerror.New("撤销采购入库失败，请稍后重试")
	}
	if count.Int() > 0 {
		return nil
	}
	return projectDataLogic.New().RemoveMaterialStockTx(ctx, tx, order["project_id"].Int64(), order["material_id"].Int64(), order["quantity"].Float64())
}

// book 为付款记录生成借贷凭证，并把凭证ID写回业务表。
//
// 财务尚未初始化时直接跳过：旧数据与未启用财务的公司仍然可以正常登记付款，
// 只是该笔记录保持「未记账」状态。
func (p *Payout) book(ctx context.Context, tx gdb.TX, id int64, form *Form, userId uint) (string, error) {
	initialized, err := booking.New().IsInitialized(ctx, tx)
	if err != nil {
		return "", err
	}
	if !initialized {
		return "", nil
	}

	projectName, err := p.projectName(ctx, tx, form.ProjectId)
	if err != nil {
		return "", err
	}

	result, err := booking.New().PostPayout(ctx, tx, &booking.Payout{
		SourceId:      id,
		ProjectName:   projectName,
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

// validate 业务校验：金额必须大于 0，费用类型必填。
func (f *Form) validate() error {
	if f.Amount <= 0 {
		return gerror.New("付款金额必须大于 0")
	}

	if strings.TrimSpace(f.PayoutType) == "" {
		return gerror.New("请选择费用类型")
	}
	if f.VendorId <= 0 {
		return gerror.New("请选择供应商")
	}

	return nil
}

// toData 转换为可落库的列值映射；付款日期为必填，解析失败返回业务错误。
func (f *Form) toData() (g.Map, error) {
	cols := defaultDao.ProjectPayout.Columns()

	payoutDate, err := parseDate(f.PayoutDate)
	if err != nil {
		return nil, err
	}

	return g.Map{
		cols.ProjectId:     f.ProjectId,
		cols.Amount:        f.Amount,
		cols.PayoutType:    f.PayoutType,
		cols.PayoutDate:    payoutDate,
		cols.VendorId:      f.VendorId,
		cols.FundAccountId: f.FundAccountId,
		cols.Remark:        f.Remark,
	}, nil
}

// matchPayoutIds 关键词命中：费用类型 / 收款方 / 备注直接匹配，或所属项目名称 / 客户姓名匹配。
func (p *Payout) matchPayoutIds(ctx context.Context, keyword string) ([]int64, error) {
	cols := defaultDao.ProjectPayout.Columns()
	like := "%" + keyword + "%"

	projectIds, err := p.matchProjectIds(ctx, like)
	if err != nil {
		return nil, err
	}

	vendorRows, err := g.DB().GetAll(ctx, "SELECT id FROM vendor WHERE name LIKE ? AND company_id=? AND deleted_at IS NULL", like, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "按供应商名称查询付款记录失败: %v", err)
		return nil, gerror.New("查询付款记录失败，请稍后重试")
	}
	vendorIds := make([]int64, 0, len(vendorRows))
	for _, row := range vendorRows {
		vendorIds = append(vendorIds, row["id"].Int64())
	}

	matchedModel := defaultDao.ProjectPayout.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Fields(cols.Id).
		WhereLike(cols.PayoutType, like).
		WhereOrLike(cols.Remark, like)
	if len(vendorIds) > 0 {
		matchedModel = matchedModel.WhereOrIn(cols.VendorId, vendorIds)
	}

	if len(projectIds) > 0 {
		matchedModel = matchedModel.WhereOrIn(cols.ProjectId, projectIds)
	}

	rows, err := matchedModel.All()
	if err != nil {
		glog.Errorf(ctx, "查询关键词命中的付款记录失败: %v", err)
		return nil, gerror.New("查询付款记录失败，请稍后重试")
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row[cols.Id].Int64())
	}

	return ids, nil
}

// matchProjectIds 关键词命中：项目名称直接匹配，或所属客户姓名匹配。
func (p *Payout) matchProjectIds(ctx context.Context, like string) ([]uint, error) {
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
		return nil, gerror.New("查询付款记录失败，请稍后重试")
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
		return nil, gerror.New("查询付款记录失败，请稍后重试")
	}

	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row[projectCols.Id].Uint())
	}

	return ids, nil
}

// attach 为付款记录批量补充项目名称、客户姓名、资金账户与凭证号。
func (p *Payout) attach(ctx context.Context, records []*defaultEntity.ProjectPayout) []*PayoutInfo {
	result := make([]*PayoutInfo, 0, len(records))
	if len(records) == 0 {
		return result
	}

	projectIds := make([]int64, 0, len(records))
	projectIdSet := make(map[int64]struct{})
	fundAccountIds := make([]int64, 0, len(records))
	fundAccountSet := make(map[int64]struct{})
	voucherIds := make([]int64, 0, len(records))
	vendorIds := make([]int64, 0, len(records))
	vendorIdSet := make(map[int64]struct{})

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
		if record.VendorId > 0 {
			if _, ok := vendorIdSet[record.VendorId]; !ok {
				vendorIdSet[record.VendorId] = struct{}{}
				vendorIds = append(vendorIds, record.VendorId)
			}
		}
	}

	projectNames, projectCustomers, customerNames := p.loadProjectMaps(ctx, projectIds)
	fundAccountNames := p.fundAccountNames(ctx, fundAccountIds)
	vendorNames := p.vendorNames(ctx, vendorIds)

	voucherNos := make(map[int64]string, len(voucherIds))
	if len(voucherIds) > 0 {
		if values, err := booking.New().VoucherNos(ctx, voucherIds); err != nil {
			// 凭证号仅用于展示，查询失败不阻断列表
			glog.Errorf(ctx, "查询付款凭证号失败, ids=%v: %v", voucherIds, err)
		} else {
			voucherNos = values
		}
	}

	for _, record := range records {
		result = append(result, &PayoutInfo{
			Payout:          record,
			ProjectName:     projectNames[record.ProjectId],
			CustomerName:    customerNames[projectCustomers[record.ProjectId]],
			VendorName:      vendorNames[record.VendorId],
			FundAccountName: fundAccountNames[record.FundAccountId],
			VoucherNo:       voucherNos[record.VoucherId],
		})
	}

	return result
}

func (p *Payout) vendorNames(ctx context.Context, ids []int64) map[int64]string {
	names := make(map[int64]string, len(ids))
	if len(ids) == 0 {
		return names
	}

	rows, err := g.DB().Model("vendor").Ctx(ctx).Fields("id", "name").Where("company_id", currentCompanyId(ctx)).WhereIn("id", ids).All()
	if err != nil {
		glog.Errorf(ctx, "查询付款记录供应商名称失败, vendorIds=%v: %v", ids, err)
		return names
	}
	for _, row := range rows {
		names[row["id"].Int64()] = row["name"].String()
	}
	return names
}

func (p *Payout) vendorName(ctx context.Context, id int64) string {
	if id <= 0 {
		return ""
	}
	value, err := g.DB().Model("vendor").Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where("id", id).Value("name")
	if err != nil {
		glog.Errorf(ctx, "查询供应商名称失败, vendorId=%d: %v", id, err)
		return ""
	}
	return value.String()
}

func (p *Payout) vendorNameTx(ctx context.Context, tx gdb.TX, id int64) (string, error) {
	if id <= 0 {
		return "", gerror.New("请选择供应商")
	}
	value, err := tx.GetValue("SELECT name FROM vendor WHERE id=? AND company_id=? AND deleted_at IS NULL", id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "校验付款供应商失败, vendorId=%d: %v", id, err)
		return "", gerror.New("校验供应商失败，请稍后重试")
	}
	if value.IsEmpty() {
		return "", gerror.New("供应商不存在或已删除")
	}
	return value.String(), nil
}

// loadProjectMaps 批量加载项目名称、项目所属客户、客户姓名。
func (p *Payout) loadProjectMaps(ctx context.Context, projectIds []int64) (map[int64]string, map[int64]int64, map[int64]string) {
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
		glog.Errorf(ctx, "查询付款记录所属项目失败, projectIds=%v: %v", projectIds, err)
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
		glog.Errorf(ctx, "查询付款记录所属客户失败, customerIds=%v: %v", customerIds, err)
		return projectNames, projectCustomers, customerNames
	}

	for _, customer := range customers {
		customerNames[int64(customer.Id)] = customer.Name
	}

	return projectNames, projectCustomers, customerNames
}

// fundAccountNames 批量查询资金账户名称。
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

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(
		"SELECT id,account_name FROM fin_fund_account WHERE id IN (%s) AND company_id=?", placeholders), append(args, currentCompanyId(ctx))...)
	if err != nil {
		glog.Errorf(ctx, "查询付款资金账户失败, ids=%v: %v", ids, err)
		return names
	}

	for _, record := range records {
		names[record["id"].Int64()] = record["account_name"].String()
	}

	return names
}

// fundAccountName 查询单个资金账户名称。
func (p *Payout) fundAccountName(ctx context.Context, id int64) string {
	if id <= 0 {
		return ""
	}

	value, err := g.DB().GetValue(ctx, "SELECT account_name FROM fin_fund_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "查询付款资金账户失败, id=%d: %v", id, err)
		return ""
	}

	return value.String()
}

// voucherNo 查询凭证号（仅已记账凭证有效）。
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

// ledgerEntries 查询凭证分录，用于详情页展示借贷方向。
func (p *Payout) ledgerEntries(ctx context.Context, voucherId int64) []*booking.EntryView {
	entries, err := booking.New().Entries(ctx, voucherId)
	if err != nil {
		glog.Errorf(ctx, "查询付款凭证分录失败, voucherId=%d: %v", voucherId, err)
		return []*booking.EntryView{}
	}

	return entries
}

// projectName 查询项目名称（事务内使用）；公司公共支出返回空串。
func (p *Payout) projectName(ctx context.Context, tx gdb.TX, projectId int64) (string, error) {
	if projectId <= 0 {
		return "", nil
	}

	value, err := tx.GetValue("SELECT name FROM project WHERE id=? AND company_id=? AND deleted_at IS NULL", projectId, currentCompanyId(ctx))
	if err != nil {
		return "", gerror.Wrap(err, "查询所属项目失败")
	}

	return value.String(), nil
}

// loadProject 查询付款记录所属项目与客户，缺失时返回 nil。
func (p *Payout) loadProject(ctx context.Context, projectId int64) (*defaultEntity.Project, *defaultEntity.Customer) {
	if projectId <= 0 {
		return nil, nil
	}

	var project defaultEntity.Project
	if err := defaultDao.Project.Ctx(ctx).
		Where("company_id", currentCompanyId(ctx)).
		Where(defaultDao.Project.Columns().Id, projectId).
		Scan(&project); err != nil {
		glog.Errorf(ctx, "查询付款记录所属项目失败, projectId=%d: %v", projectId, err)
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
		glog.Errorf(ctx, "查询付款记录所属客户失败, customerId=%d: %v", project.CustomerId, err)
		return &project, nil
	}

	return &project, &customer
}

// getPayout 按主键查询未删除付款记录，不存在时返回业务错误。
func (p *Payout) getPayout(ctx context.Context, id int64) (*defaultEntity.ProjectPayout, error) {
	var record defaultEntity.ProjectPayout
	err := defaultDao.ProjectPayout.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.ProjectPayout.Columns().Id, id).Scan(&record)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询付款记录失败, id=%d: %v", id, err)
		return nil, gerror.New("获取付款记录失败，请稍后重试")
	}

	if record.Id == 0 {
		return nil, gerror.New("付款记录不存在或已被删除")
	}

	return &record, nil
}

// lockVoucherId 在事务内锁定付款记录，返回它关联的凭证ID。
//
// mustExist=false 用于批量删除：记录已被删除时跳过，保证操作幂等。
func (p *Payout) lockVoucherId(ctx context.Context, tx gdb.TX, id int64, mustExist bool) (int64, error) {
	row, err := tx.GetOne("SELECT id,voucher_id FROM project_payout WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", id, currentCompanyId(ctx))
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

// ensureContractNotExceeded 校验「本次付款后累计付款」不超过项目合同金额。
//
// 公司公共支出（projectId=0）与未填写合同金额的项目不做限制。
// excludeId 用于修改场景排除当前记录自身。
func (p *Payout) ensureContractNotExceeded(ctx context.Context, tx gdb.TX, projectId int64, amount float64, excludeId int64) error {
	if projectId <= 0 {
		return nil
	}

	row, err := tx.GetOne("SELECT amount FROM project WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", projectId, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "读取项目合同金额失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验合同金额失败，请稍后重试")
	}
	if row == nil || row.IsEmpty() {
		return gerror.New("所选项目不存在或已被删除")
	}

	contract := roundAmount(row["amount"].Float64())
	if contract <= 0 {
		return nil
	}

	value, err := tx.GetValue(`SELECT COALESCE(SUM(amount),0) FROM project_payout
		WHERE project_id=? AND company_id=? AND deleted_at IS NULL AND id<>?`, projectId, currentCompanyId(ctx), excludeId)
	if err != nil {
		glog.Errorf(ctx, "统计项目累计付款失败, projectId=%d: %v", projectId, err)
		return gerror.New("校验累计付款失败，请稍后重试")
	}

	paid := roundAmount(value.Float64())
	total := roundAmount(paid + amount)
	if total > contract {
		return gerror.Newf("本次付款后该项目累计付款 ¥%.2f，将超过合同金额 ¥%.2f（当前已付 ¥%.2f），请核对金额",
			total, contract, paid)
	}

	return nil
}

// ensureProjectExists 校验所属项目存在；projectId=0 表示公司公共支出，跳过校验。
func (p *Payout) ensureProjectExists(ctx context.Context, projectId int64) error {
	if projectId <= 0 {
		return nil
	}

	count, err := defaultDao.Project.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Project.Columns().Id, projectId).Count()
	if err != nil {
		glog.Errorf(ctx, "校验付款记录所属项目失败, projectId=%d: %v", projectId, err)
		return gerror.New("保存付款记录失败，请稍后重试")
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

// payoutSnapshot 付款单快照，删除时用于还原日志所需信息。
type payoutSnapshot struct {
	Form      *Form
	VoucherNo string
}

// loadPayoutSnapshot 删除前读取付款单快照，读不到时返回 nil（不影响删除）。
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
			ProjectId:     record["project_id"].Int64(),
			Amount:        roundAmount(record["amount"].Float64()),
			PayoutType:    record["payout_type"].String(),
			PayoutDate:    record["payout_date"].String(),
			VendorId:      record["vendor_id"].Int64(),
			FundAccountId: record["fund_account_id"].Int64(),
			Remark:        record["remark"].String(),
		},
		VoucherNo: p.voucherNoText(ctx, record["voucher_id"].Int64()),
	}
}

// voucherNoText 凭证号，未记账或已作废时返回空串。
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

// logPayout 记录付款操作日志（封装：新增 / 修改 / 删除各调用一次即可）。
func (p *Payout) logPayout(ctx context.Context, action string, id int64, form *Form, voucherNo string) {
	if form == nil {
		return
	}

	projectName := p.projectNameText(ctx, form.ProjectId)
	vendorName := p.vendorNameText(ctx, form.VendorId)
	label := oplog.ActionLabels[action]
	if label == "" {
		label = action
	}

	oplog.New().Record(ctx, &oplog.Entry{
		Module:     oplog.ModulePayout,
		Action:     action,
		TargetType: "payout",
		TargetId:   id,
		Title:      fmt.Sprintf("%s付款：%s（%s 付 %s ￥%.2f）", label, projectName, form.PayoutType, vendorName, form.Amount),
		Amount:     form.Amount,
		Detail: g.Map{
			"payoutId":    id,
			"projectId":   form.ProjectId,
			"projectName": projectName,
			"amount":      roundAmount(form.Amount),
			"payoutType":  form.PayoutType,
			"payoutDate":  form.PayoutDate,
			"vendorId":    form.VendorId,
			"vendorName":  vendorName,
			"voucherNo":   voucherNo,
			"remark":      form.Remark,
		},
	})
}

// projectNameText 项目名称，取不到时给出占位，避免日志出现空括号。
func (p *Payout) projectNameText(ctx context.Context, projectId int64) string {
	if projectId <= 0 {
		return "公司公共支出"
	}

	value, err := g.DB().GetValue(ctx, "SELECT name FROM project WHERE id=? AND deleted_at IS NULL", projectId)
	if err != nil || value.IsEmpty() {
		return fmt.Sprintf("项目 #%d", projectId)
	}

	return value.String()
}

// vendorNameText 收款方名称，取不到时给出占位。
func (p *Payout) vendorNameText(ctx context.Context, vendorId int64) string {
	if vendorId <= 0 {
		return "未指定收款方"
	}

	value, err := g.DB().GetValue(ctx, "SELECT name FROM vendor WHERE id=? AND deleted_at IS NULL", vendorId)
	if err != nil || value.IsEmpty() {
		return fmt.Sprintf("供应商 #%d", vendorId)
	}

	return value.String()
}

// parseDate 解析 YYYY-MM-DD 日期字符串，兼容常见的斜杠 / 带时间格式。
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

// roundAmount 金额保留两位小数。
func roundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
