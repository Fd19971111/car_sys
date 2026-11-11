package voucher

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
)

// VoidWithTx 在调用方事务中作废一张已记账凭证，并同步撤销它产生的余额与资金流水。
//
// 使用场景：业务单据（收款 / 付款）被修改或删除时，先作废旧凭证，再按新数据
// 重新记账。凭证行保留（状态改为 VOID）以便追溯，但不再出现在任何报表统计中：
//   - 报表与余额查询都只统计 status_code='POSTED' 的凭证；
//   - 科目期间余额统一在作废后整表重算，撤销该凭证带来的发生额；
//   - 凭证生成的资金流水按来源（source_type + source_id）删除，随后重新生成。
//
// 已结账期间不允许作废，避免回溯修改历史账目。
func (s *Service) VoidWithTx(ctx context.Context, tx gdb.TX, voucherId int64) error {
	if voucherId <= 0 {
		return nil
	}
	companyId := currentCompanyId(ctx)

	header, err := tx.GetOne(`SELECT id,voucher_no,period_code,status_code,source_type,source_id
		FROM fin_voucher WHERE company_id=? AND id=? FOR UPDATE`, companyId, voucherId)
	if err != nil {
		return gerror.Wrap(err, "读取待作废凭证失败")
	}
	if header == nil || header.IsEmpty() {
		// 凭证已被清理：业务侧数据仍可继续保存，不阻断用户操作
		return nil
	}

	// 已作废的凭证无需重复处理（例如并发重复提交的删除请求）
	if header["status_code"].String() != voucherStatusPosted {
		return nil
	}

	periodCode := header["period_code"].String()
	row, err := tx.GetOne("SELECT status_code FROM fin_accounting_period WHERE company_id=? AND period_code=? FOR UPDATE", companyId, periodCode)
	if err != nil {
		return gerror.Wrap(err, "读取会计期间失败")
	}
	if row != nil && !row.IsEmpty() && row["status_code"].String() == "CLOSED" {
		return gerror.Newf("会计期间 %s 已结账，不能修改或删除该笔业务", periodCode)
	}

	// 资金流水以凭证为唯一事实来源：作废凭证时同步删除它生成的流水，
	// 业务侧重新记账后会按新数据重新生成，避免报表重复统计。
	if _, err = tx.Exec("DELETE FROM fin_transaction WHERE company_id=? AND source_type=? AND source_id=?",
		companyId, header["source_type"].String(), header["source_id"].Int64()); err != nil {
		return gerror.Wrap(err, "撤销凭证资金流水失败")
	}

	if _, err = tx.Exec("UPDATE fin_voucher SET status_code='VOID' WHERE company_id=? AND id=?", companyId, voucherId); err != nil {
		return gerror.Wrap(err, "作废凭证失败")
	}

	// 状态改为 VOID 后再重算余额，撤销该凭证的发生额
	if err = rebuildBalances(tx, companyId); err != nil {
		return err
	}

	return nil
}
