package voucher

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
)

// 科目期间余额重算说明
//
// fin_account_balance 是「按月结转」的科目余额表：
//   - 本期发生额 = 该期间内已记账（POSTED）凭证分录的借贷合计；
//   - 本期期初   = 上一期间的期末（首期取初始化时录入的期初余额）；
//   - 期末借方/贷方 = 期初 + 本期净发生额，按方向拆分为借/贷两列。
//
// 以前这张表是「增量维护」的，跨月不会结转，和报表口径也不一致。现在统一改为：
// 每次凭证新增 / 作废 / 月结 / 反结账后，都从启用期间开始整表重算一次。
// 数据量很小（科目数 × 期间数），一次批量 upsert 即可，换来的是任何时刻
// 余额表、报表和凭证三者都严格一致。

// balanceKey 余额表的唯一键：科目 + 期间。
type balanceKey struct {
	AccountId  int64
	PeriodCode string
}

// balanceRow 待写入的一行余额。
type balanceRow struct {
	Key                         balanceKey
	OpeningDebit, OpeningCredit int64
	CurrentDebit, CurrentCredit int64
	EndingDebit, EndingCredit   int64
}

// RebuildWithTx 在调用方事务中重算全部科目期间余额。
//
// 月结 / 反结账会改变期间状态，进而影响结转结果，因此对外暴露本方法。
func (s *Service) RebuildWithTx(ctx context.Context, tx gdb.TX) error {
	return rebuildBalances(tx, currentCompanyId(ctx))
}

// rebuildBalances 从财务启用期间开始，按月重算所有科目的余额。
func rebuildBalances(tx gdb.TX, companyId uint) error {
	initialized, err := tx.GetValue("SELECT is_initialized FROM fin_setup_config WHERE company_id=?", companyId)
	if err != nil {
		return gerror.Wrap(err, "读取财务初始化状态失败")
	}
	if initialized.Int() != 1 {
		return nil
	}

	startPeriod, err := startPeriodCode(tx, companyId)
	if err != nil {
		return err
	}

	periods, err := loadPeriodCodes(tx, startPeriod, companyId)
	if err != nil {
		return err
	}
	if len(periods) == 0 {
		return nil
	}

	accounts, err := loadAccountIds(tx, companyId)
	if err != nil {
		return err
	}

	initialOpening, err := loadInitialOpening(tx, startPeriod, companyId)
	if err != nil {
		return err
	}

	movements, err := loadMovements(tx, startPeriod, companyId)
	if err != nil {
		return err
	}

	existing, err := loadExistingBalances(tx, startPeriod, companyId)
	if err != nil {
		return err
	}

	rows := make([]balanceRow, 0, len(accounts)*len(periods))
	staleKeys := make([]balanceKey, 0)

	for _, accountId := range accounts {
		openingDebit, openingCredit := initialOpening[accountId].Debit, initialOpening[accountId].Credit
		lastEndingDebit, lastEndingCredit := openingDebit, openingCredit

		for _, periodCode := range periods {
			key := balanceKey{AccountId: accountId, PeriodCode: periodCode}

			// 首期用初始化录入的期初，之后逐月沿用上期期末
			if periodCode != startPeriod {
				openingDebit, openingCredit = lastEndingDebit, lastEndingCredit
			}

			move := movements[key]
			net := (openingDebit - openingCredit) + (move.Debit - move.Credit)

			endingDebit, endingCredit := int64(0), int64(0)
			if net >= 0 {
				endingDebit = net
			} else {
				endingCredit = -net
			}

			empty := openingDebit == 0 && openingCredit == 0 &&
				move.Debit == 0 && move.Credit == 0 &&
				endingDebit == 0 && endingCredit == 0

			if empty {
				// 删掉历史上存在、现在已归零的余额行，保持表干净
				if existing[key] {
					staleKeys = append(staleKeys, key)
				}
			} else {
				rows = append(rows, balanceRow{
					Key:           key,
					OpeningDebit:  openingDebit,
					OpeningCredit: openingCredit,
					CurrentDebit:  move.Debit,
					CurrentCredit: move.Credit,
					EndingDebit:   endingDebit,
					EndingCredit:  endingCredit,
				})
			}

			lastEndingDebit, lastEndingCredit = endingDebit, endingCredit
		}
	}

	if err = upsertBalances(tx, rows, companyId); err != nil {
		return err
	}

	return deleteStaleBalances(tx, staleKeys, companyId)
}

// amountCents 金额字符串（DECIMAL(18,2)）转分，避免浮点误差。
func amountCents(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	negative := false
	if strings.HasPrefix(value, "-") {
		negative = true
		value = strings.TrimPrefix(value, "-")
	}

	parts := strings.SplitN(value, ".", 2)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0
	}

	fraction := int64(0)
	if len(parts) == 2 {
		text := parts[1]
		if len(text) > 2 {
			text = text[:2]
		}
		text += strings.Repeat("0", 2-len(text))
		fraction, err = strconv.ParseInt(text, 10, 64)
		if err != nil {
			return 0
		}
	}

	cents := whole*100 + fraction
	if negative {
		return -cents
	}
	return cents
}

// centsText 分转两位小数的金额字符串。
func centsText(cents int64) string {
	negative := cents < 0
	if negative {
		cents = -cents
	}

	text := fmt.Sprintf("%d.%02d", cents/100, cents%100)
	if negative {
		return "-" + text
	}
	return text
}

func loadPeriodCodes(tx gdb.TX, startPeriod string, companyId uint) ([]string, error) {
	records, err := tx.GetAll("SELECT period_code FROM fin_accounting_period WHERE company_id=? AND period_code>=? ORDER BY period_code", companyId, startPeriod)
	if err != nil {
		return nil, gerror.Wrap(err, "读取会计期间失败")
	}

	periods := make([]string, 0, len(records))
	for _, record := range records {
		periods = append(periods, record["period_code"].String())
	}
	return periods, nil
}

func loadAccountIds(tx gdb.TX, companyId uint) ([]int64, error) {
	records, err := tx.GetAll("SELECT id FROM fin_account WHERE company_id=? ORDER BY id", companyId)
	if err != nil {
		return nil, gerror.Wrap(err, "读取会计科目失败")
	}

	ids := make([]int64, 0, len(records))
	for _, record := range records {
		ids = append(ids, record["id"].Int64())
	}
	return ids, nil
}

// openingAmount 期初余额（借/贷分列转成分）。
type openingAmount struct {
	Debit  int64
	Credit int64
}

func loadInitialOpening(tx gdb.TX, startPeriod string, companyId uint) (map[int64]openingAmount, error) {
	records, err := tx.GetAll(
		"SELECT account_id,opening_debit,opening_credit FROM fin_account_balance WHERE company_id=? AND period_code=?", companyId, startPeriod)
	if err != nil {
		return nil, gerror.Wrap(err, "读取期初余额失败")
	}

	result := make(map[int64]openingAmount, len(records))
	for _, record := range records {
		result[record["account_id"].Int64()] = openingAmount{
			Debit:  amountCents(record["opening_debit"].String()),
			Credit: amountCents(record["opening_credit"].String()),
		}
	}
	return result, nil
}

// movement 某科目某期间的发生额（分）。
type movement struct {
	Debit  int64
	Credit int64
}

func loadMovements(tx gdb.TX, startPeriod string, companyId uint) (map[balanceKey]movement, error) {
	records, err := tx.GetAll(`SELECT e.account_id,v.period_code,
			COALESCE(SUM(CASE WHEN e.direction_code='DEBIT' THEN e.amount ELSE 0 END),0) AS debit,
			COALESCE(SUM(CASE WHEN e.direction_code='CREDIT' THEN e.amount ELSE 0 END),0) AS credit
		FROM fin_voucher_entry e
		INNER JOIN fin_voucher v ON v.id=e.voucher_id AND v.company_id=e.company_id AND v.status_code='POSTED'
		WHERE e.company_id=? AND v.company_id=? AND v.period_code>=?
		GROUP BY e.account_id,v.period_code`, companyId, companyId, startPeriod)
	if err != nil {
		return nil, gerror.Wrap(err, "统计凭证发生额失败")
	}

	result := make(map[balanceKey]movement, len(records))
	for _, record := range records {
		result[balanceKey{
			AccountId:  record["account_id"].Int64(),
			PeriodCode: record["period_code"].String(),
		}] = movement{
			Debit:  amountCents(record["debit"].String()),
			Credit: amountCents(record["credit"].String()),
		}
	}
	return result, nil
}

func loadExistingBalances(tx gdb.TX, startPeriod string, companyId uint) (map[balanceKey]bool, error) {
	records, err := tx.GetAll(
		"SELECT account_id,period_code FROM fin_account_balance WHERE company_id=? AND period_code>=?", companyId, startPeriod)
	if err != nil {
		return nil, gerror.Wrap(err, "读取科目余额失败")
	}

	result := make(map[balanceKey]bool, len(records))
	for _, record := range records {
		result[balanceKey{
			AccountId:  record["account_id"].Int64(),
			PeriodCode: record["period_code"].String(),
		}] = true
	}
	return result, nil
}

// upsertBalances 批量写入余额，按 100 行一批，避免单条 SQL 过长。
func upsertBalances(tx gdb.TX, rows []balanceRow, companyId uint) error {
	const batchSize = 100

	for start := 0; start < len(rows); start += batchSize {
		end := start + batchSize
		if end > len(rows) {
			end = len(rows)
		}

		var builder strings.Builder
		builder.WriteString(`INSERT INTO fin_account_balance
			(company_id,account_id,period_code,opening_debit,opening_credit,current_debit,current_credit,ending_debit,ending_credit)
			VALUES `)

		args := make([]interface{}, 0, (end-start)*9)
		for index := start; index < end; index++ {
			if index > start {
				builder.WriteString(",")
			}
			builder.WriteString("(?,?,?,?,?,?,?,?,?)")

			row := rows[index]
			args = append(args,
				companyId, row.Key.AccountId, row.Key.PeriodCode,
				centsText(row.OpeningDebit), centsText(row.OpeningCredit),
				centsText(row.CurrentDebit), centsText(row.CurrentCredit),
				centsText(row.EndingDebit), centsText(row.EndingCredit),
			)
		}

		builder.WriteString(` ON DUPLICATE KEY UPDATE
			opening_debit=VALUES(opening_debit),opening_credit=VALUES(opening_credit),
			current_debit=VALUES(current_debit),current_credit=VALUES(current_credit),
			ending_debit=VALUES(ending_debit),ending_credit=VALUES(ending_credit)`)

		if _, err := tx.Exec(builder.String(), args...); err != nil {
			return gerror.Wrap(err, "重算科目余额失败")
		}
	}

	return nil
}

// deleteStaleBalances 删除重算后已经归零的余额行。
func deleteStaleBalances(tx gdb.TX, keys []balanceKey, companyId uint) error {
	const batchSize = 200

	for start := 0; start < len(keys); start += batchSize {
		end := start + batchSize
		if end > len(keys) {
			end = len(keys)
		}

		var builder strings.Builder
		builder.WriteString("DELETE FROM fin_account_balance WHERE company_id=? AND (account_id,period_code) IN (")

		// 注意：这里必须是长度为 0 的切片，否则会多出一个 nil 占位参数，
		// 导致占位符与参数错位（sql: expected N arguments, got N+1）。
		args := make([]interface{}, 0, (end-start)*2+1)
		args = append(args, companyId)
		for index := start; index < end; index++ {
			if index > start {
				builder.WriteString(",")
			}
			builder.WriteString("(?,?)")
			args = append(args, keys[index].AccountId, keys[index].PeriodCode)
		}
		builder.WriteString(")")

		if _, err := tx.Exec(builder.String(), args...); err != nil {
			return gerror.Wrap(err, "清理归零余额失败")
		}
	}

	return nil
}
