package driver_finance

import (
	"context"
	"fmt"
	"strings"
	"time"

	driverFinanceApi "engineServ/api/driver_finance"
	v1 "engineServ/api/driver_finance/v1"
	"engineServ/internal/consts"
	"engineServ/internal/logic/fin/booking"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() driverFinanceApi.IDriverFinanceV1 { return &Controller{} }

func (c *Controller) Drivers(ctx context.Context, _ *v1.DriversReq) (*v1.DriversRes, error) {
	companyId := currentCompanyId(ctx)
	rows, err := g.DB().GetAll(ctx, `SELECT d.id,d.name,d.phone,
		COUNT(e.id) AS trip_count,COALESCE(SUM(e.driver_amount),0) AS earned_total,
		COALESCE(SUM(COALESCE(st.settled,0)),0) AS settled_total
		FROM driver d
		LEFT JOIN ride_driver_earning e ON e.driver_id=d.id AND e.company_id=d.company_id AND e.status='POSTED'
		LEFT JOIN (
			SELECT i.earning_id,SUM(i.amount) AS settled FROM ride_driver_settlement_item i
			INNER JOIN ride_driver_settlement s ON s.id=i.settlement_id AND s.company_id=i.company_id AND s.status='POSTED' AND s.deleted_at IS NULL
			GROUP BY i.earning_id
		) st ON st.earning_id=e.id
		WHERE d.company_id=? AND d.deleted_at IS NULL
		GROUP BY d.id,d.name,d.phone ORDER BY earned_total-settled_total DESC,d.name`, companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "查询司机待结余额失败")
	}
	items := make([]*v1.DriverBalance, 0, len(rows))
	for _, row := range rows {
		earned, settled := row["earned_total"].Float64(), row["settled_total"].Float64()
		items = append(items, &v1.DriverBalance{
			DriverId: row["id"].Int64(), DriverName: row["name"].String(), Phone: row["phone"].String(),
			TripCount: row["trip_count"].Int(), EarnedTotal: round2(earned), SettledTotal: round2(settled), Payable: round2(earned - settled),
		})
	}
	return &v1.DriversRes{Records: items}, nil
}

func (c *Controller) Settlements(ctx context.Context, req *v1.SettlementsReq) (*v1.SettlementsRes, error) {
	current, size := req.Current, req.Size
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	companyId := currentCompanyId(ctx)
	where := "s.company_id=? AND s.deleted_at IS NULL AND s.status='POSTED'"
	args := []any{companyId}
	if req.DriverId > 0 {
		where += " AND s.driver_id=?"
		args = append(args, req.DriverId)
	}
	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM ride_driver_settlement s WHERE "+where, args...)
	if err != nil {
		return nil, gerror.Wrap(err, "统计司机结算单失败")
	}
	queryArgs := append(append([]any{}, args...), (current-1)*size, size)
	rows, err := g.DB().GetAll(ctx, `SELECT s.id,s.settlement_no,s.driver_id,COALESCE(d.name,'') AS driver_name,
		DATE_FORMAT(s.settlement_date,'%Y-%m-%d') AS settlement_date,s.amount,s.remark,
		COALESCE(f.account_name,'') AS fund_account_name,COALESCE(v.voucher_no,'') AS voucher_no
		FROM ride_driver_settlement s
		LEFT JOIN driver d ON d.id=s.driver_id AND d.company_id=s.company_id
		LEFT JOIN fin_fund_account f ON f.id=s.fund_account_id AND f.company_id=s.company_id
		LEFT JOIN fin_voucher v ON v.id=s.voucher_id AND v.company_id=s.company_id
		WHERE `+where+` ORDER BY s.settlement_date DESC,s.id DESC LIMIT ?,?`, queryArgs...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询司机结算记录失败")
	}
	items := make([]*v1.SettlementItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &v1.SettlementItem{Id: row["id"].Int64(), SettlementNo: row["settlement_no"].String(), DriverId: row["driver_id"].Int64(), DriverName: row["driver_name"].String(), SettlementDate: row["settlement_date"].String(), Amount: row["amount"].Float64(), FundAccountName: row["fund_account_name"].String(), VoucherNo: row["voucher_no"].String(), Remark: row["remark"].String()})
	}
	return &v1.SettlementsRes{Records: items, Total: total.Int()}, nil
}

func (c *Controller) CreateSettlement(ctx context.Context, req *v1.CreateSettlementReq) (*v1.CreateSettlementRes, error) {
	date, err := time.Parse("2006-01-02", strings.TrimSpace(req.SettlementDate))
	if err != nil {
		return nil, gerror.New("结算日期格式不正确")
	}
	amount := round2(req.Amount)
	if amount <= 0 {
		return nil, gerror.New("结算金额必须大于0")
	}
	companyId, userId := currentCompanyId(ctx), currentUserId(ctx)
	var response *v1.CreateSettlementRes
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		driver, txErr := tx.GetOne("SELECT name FROM driver WHERE id=? AND company_id=? AND deleted_at IS NULL FOR UPDATE", req.DriverId, companyId)
		if txErr != nil {
			return gerror.Wrap(txErr, "读取司机信息失败")
		}
		if driver.IsEmpty() {
			return gerror.New("司机不存在或已停用")
		}
		initialized, txErr := booking.New().IsInitialized(ctx, tx)
		if txErr != nil {
			return txErr
		}
		if !initialized {
			return gerror.New("请先完成财务初始化，再进行司机结算")
		}
		rows, txErr := tx.GetAll(`SELECT e.id,e.driver_amount-COALESCE(SUM(CASE WHEN s.id IS NOT NULL THEN i.amount ELSE 0 END),0) AS outstanding
			FROM ride_driver_earning e
			LEFT JOIN ride_driver_settlement_item i ON i.earning_id=e.id AND i.company_id=e.company_id
			LEFT JOIN ride_driver_settlement s ON s.id=i.settlement_id AND s.company_id=i.company_id AND s.status='POSTED' AND s.deleted_at IS NULL
			WHERE e.company_id=? AND e.driver_id=? AND e.status='POSTED'
			GROUP BY e.id,e.driver_amount,e.earning_date
			HAVING outstanding>0 ORDER BY e.earning_date,e.id FOR UPDATE`, companyId, req.DriverId)
		if txErr != nil {
			return gerror.Wrap(txErr, "读取司机未结车费失败")
		}
		available := 0.0
		for _, row := range rows {
			available += row["outstanding"].Float64()
		}
		if amount > available+0.005 {
			return gerror.Newf("司机待结余额仅 ¥%.2f，不能结算 ¥%.2f", round2(available), amount)
		}
		result, txErr := tx.Exec(`INSERT INTO ride_driver_settlement
			(company_id,driver_id,settlement_no,settlement_date,amount,fund_account_id,created_by,remark)
			VALUES(?,?, '',?,?,?, ?,?)`, companyId, req.DriverId, date.Format("2006-01-02"), amount, req.FundAccountId, userId, strings.TrimSpace(req.Remark))
		if txErr != nil {
			return gerror.Wrap(txErr, "创建司机结算单失败")
		}
		settlementId, txErr := result.LastInsertId()
		if txErr != nil {
			return gerror.Wrap(txErr, "读取司机结算单ID失败")
		}
		settlementNo := fmt.Sprintf("DS%s%06d", date.Format("20060102"), settlementId)
		if _, txErr = tx.Exec("UPDATE ride_driver_settlement SET settlement_no=? WHERE id=? AND company_id=?", settlementNo, settlementId, companyId); txErr != nil {
			return gerror.Wrap(txErr, "生成司机结算单号失败")
		}
		remaining := amount
		for _, row := range rows {
			if remaining <= 0.005 {
				break
			}
			outstanding := round2(row["outstanding"].Float64())
			allocated := outstanding
			if allocated > remaining {
				allocated = remaining
			}
			if _, txErr = tx.Exec("INSERT INTO ride_driver_settlement_item (company_id,settlement_id,earning_id,amount) VALUES(?,?,?,?)", companyId, settlementId, row["id"].Int64(), allocated); txErr != nil {
				return gerror.Wrap(txErr, "分摊司机结算明细失败")
			}
			remaining = round2(remaining - allocated)
		}
		voucherResult, txErr := booking.New().PostDriverSettlement(ctx, tx, settlementId, req.DriverId, driver["name"].String(), date.Format("2006-01-02"), amount, req.FundAccountId, userId)
		if txErr != nil {
			return txErr
		}
		if _, txErr = tx.Exec("UPDATE ride_driver_settlement SET voucher_id=? WHERE id=? AND company_id=?", voucherResult.VoucherId, settlementId, companyId); txErr != nil {
			return gerror.Wrap(txErr, "保存司机结算凭证关联失败")
		}
		response = &v1.CreateSettlementRes{Id: settlementId, SettlementNo: settlementNo, VoucherNo: voucherResult.VoucherNo}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func currentCompanyId(ctx context.Context) uint { return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey)) }
func currentUserId(ctx context.Context) uint    { return gconv.Uint(ctx.Value(consts.CtxUserIdKey)) }
func round2(value float64) float64              { return float64(int64(value*100+0.5)) / 100 }
