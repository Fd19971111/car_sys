package fin_init

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	v1 "engineServ/api/fin_init/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type FinanceInit struct{}

func New() *FinanceInit { return &FinanceInit{} }

func (s *FinanceInit) Save(ctx context.Context, userId uint, req *v1.SaveReq) error {
	if req.Year < 2000 || req.Year > 2100 || req.StartMonth < 1 || req.StartMonth > 12 {
		return gerror.New("会计年度或启用月份不正确")
	}
	if len(req.Accounts) == 0 {
		return gerror.New("请至少启用一个会计科目")
	}

	accounts := make(map[string]v1.Account, len(req.Accounts))
	for _, account := range req.Accounts {
		account.Code = strings.TrimSpace(account.Code)
		account.Name = strings.TrimSpace(account.Name)
		if account.Code == "" || account.Name == "" || len(account.Code) > 32 || len(account.Name) > 100 {
			return gerror.New("科目编码和名称不能为空且长度需符合要求")
		}
		if _, exists := accounts[account.Code]; exists {
			return gerror.New("科目编码不能重复")
		}
		if _, ok := directionByCategory(account.CategoryCode); !ok {
			return gerror.New("科目类别不正确")
		}
		accounts[account.Code] = account
	}
	fundAccounts := make([]v1.FundAccount, 0, len(req.FundAccounts))
	fundAccountNames := make(map[string]struct{}, len(req.FundAccounts))
	for _, account := range req.FundAccounts {
		account.Type = strings.TrimSpace(account.Type)
		account.Name = strings.TrimSpace(account.Name)
		account.BankName = strings.TrimSpace(account.BankName)
		account.AccountLast4 = strings.TrimSpace(account.AccountLast4)
		account.LedgerAccountCode = strings.TrimSpace(account.LedgerAccountCode)
		if account.Type != "CASH" && account.Type != "BANK" {
			return gerror.New("资金账户类型不正确")
		}
		if account.Name == "" || len(account.Name) > 100 {
			return gerror.New("请输入有效的资金账户名称")
		}
		if account.Type == "BANK" && account.BankName == "" {
			return gerror.New("银行账户需要填写开户银行")
		}
		if len(account.BankName) > 100 || len(account.AccountLast4) > 4 || (account.AccountLast4 != "" && !regexp.MustCompile(`^\d{4}$`).MatchString(account.AccountLast4)) {
			return gerror.New("银行卡只允许填写4位数字末尾号")
		}
		ledgerAccount, ok := accounts[account.LedgerAccountCode]
		if !ok || !ledgerAccount.Enabled || ledgerAccount.CategoryCode != "ASSET" {
			return gerror.New("资金账户必须关联一个已启用的资产科目")
		}
		uniqueName := strings.ToUpper(account.Type + ":" + account.Name)
		if _, exists := fundAccountNames[uniqueName]; exists {
			return gerror.New("同类型资金账户名称不能重复")
		}
		fundAccountNames[uniqueName] = struct{}{}
		fundAccounts = append(fundAccounts, account)
	}
	openingAmounts := make(map[string]string, len(req.OpeningBalances))
	debitTotal := new(big.Int)
	creditTotal := new(big.Int)
	for _, balance := range req.OpeningBalances {
		balance.AccountCode = strings.TrimSpace(balance.AccountCode)
		balance.Amount = strings.TrimSpace(balance.Amount)
		account, ok := accounts[balance.AccountCode]
		if !ok || !account.Enabled {
			return gerror.New("期初余额只能填写在已启用的科目上")
		}
		if _, exists := openingAmounts[balance.AccountCode]; exists {
			return gerror.New("同一科目的期初余额不能重复")
		}
		amountCents, normalizedAmount, ok := parseOpeningAmount(balance.Amount)
		if !ok {
			return gerror.Newf("科目 %s 的期初余额格式不正确，最多输入两位小数", balance.AccountCode)
		}
		openingAmounts[balance.AccountCode] = normalizedAmount
		direction, _ := directionByCategory(account.CategoryCode)
		if direction == "DEBIT" {
			debitTotal.Add(debitTotal, amountCents)
		} else {
			creditTotal.Add(creditTotal, amountCents)
		}
	}
	if debitTotal.Cmp(creditTotal) != 0 {
		return gerror.Newf("期初借贷不平衡：借方合计 %s 元，贷方合计 %s 元", formatOpeningAmount(debitTotal), formatOpeningAmount(creditTotal))
	}

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Exec(`INSERT INTO fin_setup_config
			(id,company_id,is_initialized,period_year,start_month,initialized_by,initialized_at)
			VALUES(?,?,0,0,1,0,NULL) ON DUPLICATE KEY UPDATE company_id=VALUES(company_id)`, currentCompanyId(ctx), currentCompanyId(ctx)); err != nil {
			return gerror.Wrap(err, "创建财务初始化配置失败")
		}
		initialized, err := tx.GetValue("SELECT is_initialized FROM fin_setup_config WHERE company_id=? FOR UPDATE", currentCompanyId(ctx))
		if err != nil {
			return gerror.Wrap(err, "读取财务初始化状态失败，请确认已执行 fin_init.sql")
		}
		if initialized.Int() == 1 && !req.Reconfigure {
			return gerror.New("财务已完成初始化")
		}
		if req.Reconfigure {
			voucherCount, countErr := tx.GetValue("SELECT COUNT(*) FROM fin_voucher WHERE company_id=?", currentCompanyId(ctx))
			if countErr != nil {
				return gerror.Wrap(countErr, "检查现有凭证失败")
			}
			transactionCount, countErr := tx.GetValue("SELECT COUNT(*) FROM fin_transaction WHERE company_id=?", currentCompanyId(ctx))
			if countErr != nil {
				return gerror.Wrap(countErr, "检查现有资金流水失败")
			}
			if voucherCount.Int() > 0 || transactionCount.Int() > 0 {
				return gerror.New("已有凭证或资金流水数据，不能重新配置财务基础资料")
			}
			if _, err = tx.Exec("DELETE FROM fin_account_balance WHERE company_id=?", currentCompanyId(ctx)); err != nil {
				return gerror.Wrap(err, "清理原有期初余额失败")
			}
		}

		for _, account := range accounts {
			direction, _ := directionByCategory(account.CategoryCode)
			_, err = tx.Exec(`INSERT INTO fin_account
				(company_id,account_code,account_name,category_code,direction_code,enabled,is_system)
				VALUES(?,?,?,?,?,?,0)
				ON DUPLICATE KEY UPDATE account_name=VALUES(account_name),category_code=VALUES(category_code),direction_code=VALUES(direction_code),enabled=VALUES(enabled)`,
				currentCompanyId(ctx), account.Code, account.Name, account.CategoryCode, direction, account.Enabled)
			if err != nil {
				return gerror.Wrapf(err, "保存科目 %s 失败", account.Code)
			}
		}

		if _, err = tx.Exec("UPDATE fin_account SET enabled=0 WHERE company_id=?", currentCompanyId(ctx)); err != nil {
			return gerror.Wrap(err, "重置科目启用状态失败")
		}
		for _, account := range accounts {
			if account.Enabled {
				if _, err = tx.Exec("UPDATE fin_account SET enabled=1 WHERE company_id=? AND account_code=?", currentCompanyId(ctx), account.Code); err != nil {
					return gerror.Wrap(err, "启用科目失败")
				}
			}
		}

		if _, err = tx.Exec("DELETE FROM fin_fund_account WHERE company_id=?", currentCompanyId(ctx)); err != nil {
			return gerror.Wrap(err, "更新现金及银行账户失败")
		}
		for _, account := range fundAccounts {
			_, err = tx.Exec(`INSERT INTO fin_fund_account
				(company_id,account_type,account_name,bank_name,account_last4,ledger_account_id,enabled)
				SELECT ?,?,?,?, ?,id,1 FROM fin_account WHERE company_id=? AND account_code=? AND enabled=1`,
				currentCompanyId(ctx), account.Type, account.Name, account.BankName, account.AccountLast4, currentCompanyId(ctx), account.LedgerAccountCode)
			if err != nil {
				return gerror.Wrapf(err, "保存资金账户 %s 失败", account.Name)
			}
		}

		for month := 1; month <= 12; month++ {
			status := "CLOSED"
			if month > req.StartMonth {
				status = "FUTURE"
			} else if month == req.StartMonth {
				status = "OPEN"
			}
			periodCode := fmt.Sprintf("%04d-%02d", req.Year, month)
			if _, err = tx.Exec(`INSERT INTO fin_accounting_period(company_id,period_code,period_year,period_month,status_code)
				VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE status_code=VALUES(status_code)`, currentCompanyId(ctx), periodCode, req.Year, month, status); err != nil {
				return gerror.Wrapf(err, "保存会计期间 %s 失败", periodCode)
			}
		}

		periodCode := fmt.Sprintf("%04d-%02d", req.Year, req.StartMonth)
		for _, account := range accounts {
			amount := openingAmounts[account.Code]
			if amount == "" {
				amount = "0.00"
			}
			debitOpening, creditOpening := "0.00", "0.00"
			if direction, _ := directionByCategory(account.CategoryCode); direction == "DEBIT" {
				debitOpening = amount
			} else {
				creditOpening = amount
			}
			_, err = tx.Exec(`INSERT INTO fin_account_balance
				(company_id,account_id,period_code,opening_debit,opening_credit,current_debit,current_credit,ending_debit,ending_credit)
				SELECT ?,id,?,?,?,?,?,?,? FROM fin_account WHERE company_id=? AND account_code=? AND enabled=1
				ON DUPLICATE KEY UPDATE opening_debit=VALUES(opening_debit),opening_credit=VALUES(opening_credit),
				current_debit=VALUES(current_debit),current_credit=VALUES(current_credit),
				ending_debit=VALUES(ending_debit),ending_credit=VALUES(ending_credit)`,
				currentCompanyId(ctx), periodCode, debitOpening, creditOpening, "0.00", "0.00", debitOpening, creditOpening, currentCompanyId(ctx), account.Code)
			if err != nil {
				return gerror.Wrapf(err, "保存科目 %s 的期初余额失败", account.Code)
			}
		}

		_, err = tx.Exec(`UPDATE fin_setup_config SET is_initialized=1,period_year=?,start_month=?,initialized_by=?,initialized_at=NOW() WHERE company_id=?`, req.Year, req.StartMonth, userId, currentCompanyId(ctx))
		if err != nil {
			return gerror.Wrap(err, "保存财务初始化状态失败")
		}
		return nil
	})
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func parseOpeningAmount(amount string) (*big.Int, string, bool) {
	if !regexp.MustCompile(`^\d{1,16}(\.\d{1,2})?$`).MatchString(amount) {
		return nil, "", false
	}
	parts := strings.SplitN(amount, ".", 2)
	fraction := "00"
	if len(parts) == 2 {
		fraction = parts[1] + strings.Repeat("0", 2-len(parts[1]))
	}
	cents, ok := new(big.Int).SetString(parts[0]+fraction, 10)
	if !ok {
		return nil, "", false
	}
	return cents, parts[0] + "." + fraction, true
}

func formatOpeningAmount(cents *big.Int) string {
	whole, fraction := new(big.Int), new(big.Int)
	whole.QuoRem(cents, big.NewInt(100), fraction)
	return fmt.Sprintf("%s.%02d", whole.String(), fraction.Int64())
}

func directionByCategory(category string) (string, bool) {
	switch category {
	case "ASSET", "EXPENSE":
		return "DEBIT", true
	case "LIABILITY", "EQUITY", "REVENUE":
		return "CREDIT", true
	default:
		return "", false
	}
}
