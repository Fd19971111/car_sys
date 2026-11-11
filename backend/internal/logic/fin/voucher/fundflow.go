package voucher

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
)

// insertFundTransactions 为涉及现金 / 银行科目的分录生成资金流水。
//
// 资金流水以会计凭证为唯一事实来源，避免业务侧重复记流水：
//   - 现金 / 银行科目借方发生额 => INCOME（资金流入）
//   - 现金 / 银行科目贷方发生额 => EXPENSE（资金流出）
//
// 是否属于资金科目以 fin_fund_account.ledger_account_id 为准，普通科目不产生流水。
// 流水与凭证在同一个事务内写入，凭证回滚时流水一起回滚。
func insertFundTransactions(tx gdb.TX, header voucher, entries []preparedEntry, companyId uint) error {
	accountIds := make([]int64, 0, len(entries))
	seen := make(map[int64]struct{}, len(entries))
	for _, entry := range entries {
		if _, exists := seen[entry.AccountId]; exists {
			continue
		}
		seen[entry.AccountId] = struct{}{}
		accountIds = append(accountIds, entry.AccountId)
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(accountIds)), ",")
	args := make([]interface{}, 0, len(accountIds))
	for _, id := range accountIds {
		args = append(args, id)
	}

	rows, err := tx.GetAll(
		fmt.Sprintf("SELECT ledger_account_id FROM fin_fund_account WHERE company_id=? AND enabled=1 AND ledger_account_id IN (%s)", placeholders),
		append([]interface{}{companyId}, args...)...,
	)
	if err != nil {
		return gerror.Wrap(err, "查询资金账户失败，请确认已执行 fin_init.sql")
	}

	fundAccountIds := make(map[int64]struct{}, len(rows))
	for _, row := range rows {
		fundAccountIds[row["ledger_account_id"].Int64()] = struct{}{}
	}
	if len(fundAccountIds) == 0 {
		return nil
	}

	occurredAt := header.VoucherDate + " 00:00:00"
	for _, entry := range entries {
		if _, isFund := fundAccountIds[entry.AccountId]; !isFund {
			continue
		}

		direction := "EXPENSE"
		if entry.Direction == DirectionDebit {
			direction = "INCOME"
		}

		if _, err = tx.Exec(`INSERT INTO fin_transaction
			(company_id,source_type,source_id,account_id,occurred_at,amount,direction_code,remark)
			VALUES(?,?,?,?,?,?,?,?)`,
			companyId, header.SourceType, header.SourceId, entry.AccountId, occurredAt,
			entry.Amount, direction, entry.Summary,
		); err != nil {
			return gerror.Wrapf(err, "记录科目 %s 的资金流水失败", entry.AccountCode)
		}
	}

	return nil
}
