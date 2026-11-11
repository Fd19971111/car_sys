package fin_basic

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	accountCodePattern = regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z._-]{0,31}$`)
	last4Pattern       = regexp.MustCompile(`^\d{4}$`)
)

// AccountFilter 会计科目查询条件。
type AccountFilter struct {
	Keyword      string
	CategoryCode string
	Enabled      *bool
}

// Account 会计科目。
type Account struct {
	Id            int64
	Code          string
	Name          string
	CategoryCode  string
	DirectionCode string
	Enabled       bool
	IsSystem      bool
	SortOrder     int
	Used          bool
}

// AccountSummary 会计科目统计。
type AccountSummary struct {
	Total            int
	Enabled          int
	Disabled         int
	System           int
	FundAccountCount int
}

// AccountForm 会计科目保存参数。
type AccountForm struct {
	Id           int64
	Code         string
	Name         string
	CategoryCode string
	Enabled      bool
}

// FundAccountFilter 资金账户查询条件。
type FundAccountFilter struct {
	Keyword string
	Type    string
}

// FundAccount 资金账户。
type FundAccount struct {
	Id                int64
	Type              string
	Name              string
	BankName          string
	AccountLast4      string
	LedgerAccountId   int64
	LedgerAccountCode string
	LedgerAccountName string
	Enabled           bool
}

// FundAccountForm 资金账户保存参数。
type FundAccountForm struct {
	Id                int64
	Type              string
	Name              string
	BankName          string
	AccountLast4      string
	LedgerAccountCode string
	Enabled           bool
}

// Service 财务基础资料服务。
type Service struct{}

// New 创建财务基础资料服务。
func New() *Service { return &Service{} }

// Accounts 查询会计科目及统计信息。
func (s *Service) Accounts(ctx context.Context, filter *AccountFilter) ([]*Account, *AccountSummary, error) {
	where := []string{"a.company_id=?"}
	args := []interface{}{currentCompanyId(ctx)}
	if filter != nil {
		if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
			like := "%" + keyword + "%"
			where = append(where, "(a.account_code LIKE ? OR a.account_name LIKE ?)")
			args = append(args, like, like)
		}
		if category := strings.TrimSpace(filter.CategoryCode); category != "" {
			where = append(where, "a.category_code=?")
			args = append(args, category)
		}
		if filter.Enabled != nil {
			where = append(where, "a.enabled=?")
			args = append(args, *filter.Enabled)
		}
	}

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			a.id,a.account_code,a.account_name,a.category_code,a.direction_code,
			a.enabled,a.is_system,a.sort_order,
			(EXISTS(SELECT 1 FROM fin_voucher_entry e WHERE e.account_id=a.id AND e.company_id=a.company_id)
				OR EXISTS(SELECT 1 FROM fin_account_balance b WHERE b.account_id=a.id AND b.company_id=a.company_id)
				OR EXISTS(SELECT 1 FROM fin_fund_account f WHERE f.ledger_account_id=a.id AND f.company_id=a.company_id)) AS used
		FROM fin_account a
		WHERE %s
		ORDER BY a.sort_order,a.account_code`, strings.Join(where, " AND ")), args...)
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询会计科目失败，请确认已执行 fin_init.sql")
	}

	items := make([]*Account, 0, len(records))
	for _, row := range records {
		items = append(items, &Account{
			Id:            row["id"].Int64(),
			Code:          row["account_code"].String(),
			Name:          row["account_name"].String(),
			CategoryCode:  row["category_code"].String(),
			DirectionCode: row["direction_code"].String(),
			Enabled:       row["enabled"].Bool(),
			IsSystem:      row["is_system"].Bool(),
			SortOrder:     row["sort_order"].Int(),
			Used:          row["used"].Bool(),
		})
	}

	summary := &AccountSummary{}
	row, err := g.DB().GetOne(ctx, `SELECT
		(SELECT COUNT(*) FROM fin_account WHERE company_id=?) AS total,
		(SELECT COUNT(*) FROM fin_account WHERE company_id=? AND enabled=1) AS enabled,
		(SELECT COUNT(*) FROM fin_account WHERE company_id=? AND enabled=0) AS disabled,
		(SELECT COUNT(*) FROM fin_account WHERE company_id=? AND is_system=1) AS system_count,
		(SELECT COUNT(*) FROM fin_fund_account WHERE company_id=?) AS fund_account_count`,
		currentCompanyId(ctx), currentCompanyId(ctx), currentCompanyId(ctx), currentCompanyId(ctx), currentCompanyId(ctx))
	if err != nil {
		return nil, nil, gerror.Wrap(err, "统计会计科目失败")
	}
	if row != nil && !row.IsEmpty() {
		summary.Total = row["total"].Int()
		summary.Enabled = row["enabled"].Int()
		summary.Disabled = row["disabled"].Int()
		summary.System = row["system_count"].Int()
		summary.FundAccountCount = row["fund_account_count"].Int()
	}

	return items, summary, nil
}

// SaveAccount 新增或修改会计科目。
func (s *Service) SaveAccount(ctx context.Context, form *AccountForm) (int64, error) {
	code := strings.TrimSpace(form.Code)
	name := strings.TrimSpace(form.Name)
	category := strings.ToUpper(strings.TrimSpace(form.CategoryCode))

	if !accountCodePattern.MatchString(code) {
		return 0, gerror.New("科目编码只能由字母、数字、点、中划线或下划线组成，最长32位")
	}
	if name == "" || len([]rune(name)) > 100 {
		return 0, gerror.New("科目名称不能为空且不能超过100个字符")
	}
	direction, ok := directionByCategory(category)
	if !ok {
		return 0, gerror.New("科目类别不正确")
	}

	if form.Id > 0 {
		row, err := g.DB().GetOne(ctx, "SELECT id,account_code,category_code,is_system FROM fin_account WHERE id=? AND company_id=?", form.Id, currentCompanyId(ctx))
		if err != nil {
			return 0, gerror.Wrap(err, "读取会计科目失败")
		}
		if row == nil || row.IsEmpty() {
			return 0, gerror.New("会计科目不存在")
		}

		oldCode := row["account_code"].String()
		oldCategory := row["category_code"].String()
		if row["is_system"].Bool() && (code != oldCode || category != oldCategory) {
			return 0, gerror.New("系统预置科目不允许修改编码和类别")
		}
		if code != oldCode || category != oldCategory {
			used, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM fin_voucher_entry WHERE account_id=? AND company_id=?", form.Id, currentCompanyId(ctx))
			if err != nil {
				return 0, gerror.Wrap(err, "校验科目使用情况失败")
			}
			if used.Int() > 0 {
				return 0, gerror.New("该科目已有凭证分录，不能修改科目编码或类别")
			}
		}
		if !form.Enabled {
			if err = s.ensureAccountNotReferencedByFund(ctx, form.Id); err != nil {
				return 0, err
			}
		}

		_, err = g.DB().Exec(ctx, `UPDATE fin_account
			SET account_code=?,account_name=?,category_code=?,direction_code=?,enabled=?
			WHERE id=? AND company_id=?`, code, name, category, direction, form.Enabled, form.Id, currentCompanyId(ctx))
		if err != nil {
			if isDuplicateError(err) {
				return 0, gerror.Newf("科目编码 %s 已存在", code)
			}
			return 0, gerror.Wrap(err, "更新会计科目失败")
		}
		return form.Id, nil
	}

	exists, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM fin_account WHERE account_code=? AND company_id=?", code, currentCompanyId(ctx))
	if err != nil {
		return 0, gerror.Wrap(err, "校验科目编码失败")
	}
	if exists.Int() > 0 {
		return 0, gerror.Newf("科目编码 %s 已存在", code)
	}

	sortOrder, err := g.DB().GetValue(ctx, "SELECT COALESCE(MAX(sort_order),0)+10 FROM fin_account WHERE company_id=?", currentCompanyId(ctx))
	if err != nil {
		return 0, gerror.Wrap(err, "生成科目排序号失败")
	}

	result, err := g.DB().Exec(ctx, `INSERT INTO fin_account
		(company_id,account_code,account_name,category_code,direction_code,enabled,is_system,sort_order)
		VALUES(?,?,?,?,?,?,0,?)`, currentCompanyId(ctx), code, name, category, direction, form.Enabled, sortOrder.Int())
	if err != nil {
		if isDuplicateError(err) {
			return 0, gerror.Newf("科目编码 %s 已存在", code)
		}
		return 0, gerror.Wrap(err, "新增会计科目失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取科目ID失败")
	}
	return id, nil
}

// DeleteAccount 删除会计科目：系统预置或已被引用的科目不允许删除。
func (s *Service) DeleteAccount(ctx context.Context, id int64) error {
	row, err := g.DB().GetOne(ctx, "SELECT is_system FROM fin_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.Wrap(err, "读取会计科目失败")
	}
	if row == nil || row.IsEmpty() {
		return gerror.New("会计科目不存在")
	}
	if row["is_system"].Bool() {
		return gerror.New("系统预置科目不能删除，可以改为停用")
	}

	references, err := g.DB().GetValue(ctx, `SELECT
		(SELECT COUNT(*) FROM fin_voucher_entry WHERE account_id=? AND company_id=?)
		+ (SELECT COUNT(*) FROM fin_account_balance WHERE account_id=? AND company_id=?)
		+ (SELECT COUNT(*) FROM fin_fund_account WHERE ledger_account_id=? AND company_id=?)`,
		id, currentCompanyId(ctx), id, currentCompanyId(ctx), id, currentCompanyId(ctx))
	if err != nil {
		return gerror.Wrap(err, "校验科目引用情况失败")
	}
	if references.Int() > 0 {
		return gerror.New("该科目已被凭证、余额或资金账户引用，不能删除，请改为停用")
	}

	if _, err = g.DB().Exec(ctx, "DELETE FROM fin_account WHERE id=? AND company_id=? AND is_system=0", id, currentCompanyId(ctx)); err != nil {
		return gerror.Wrap(err, "删除会计科目失败")
	}
	return nil
}

// ToggleAccount 启用或停用会计科目。
func (s *Service) ToggleAccount(ctx context.Context, id int64, enabled bool) error {
	exists, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM fin_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.Wrap(err, "读取会计科目失败")
	}
	if exists.Int() == 0 {
		return gerror.New("会计科目不存在")
	}
	if !enabled {
		if err = s.ensureAccountNotReferencedByFund(ctx, id); err != nil {
			return err
		}
	}

	if _, err = g.DB().Exec(ctx, "UPDATE fin_account SET enabled=? WHERE id=? AND company_id=?", enabled, id, currentCompanyId(ctx)); err != nil {
		return gerror.Wrap(err, "更新科目状态失败")
	}
	return nil
}

// ensureAccountNotReferencedByFund 校验科目没有被启用中的资金账户引用。
//
// 历史凭证与余额仍可保留，停用科目时只拦截资金账户引用。
func (s *Service) ensureAccountNotReferencedByFund(ctx context.Context, id int64) error {
	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM fin_fund_account WHERE ledger_account_id=? AND company_id=? AND enabled=1", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.Wrap(err, "校验资金账户引用失败")
	}
	if count.Int() > 0 {
		return gerror.New("该科目正被启用中的资金账户引用，不能停用或删除")
	}
	return nil
}

// FundAccounts 查询资金账户。
func (s *Service) FundAccounts(ctx context.Context, filter *FundAccountFilter) ([]*FundAccount, error) {
	where := []string{"f.company_id=?", "a.company_id=f.company_id"}
	args := make([]interface{}, 0, 3)
	args = append(args, currentCompanyId(ctx))
	if filter != nil {
		if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
			like := "%" + keyword + "%"
			where = append(where, "(f.account_name LIKE ? OR f.bank_name LIKE ? OR f.account_last4 LIKE ? OR a.account_code LIKE ? OR a.account_name LIKE ?)")
			args = append(args, like, like, like, like, like)
		}
		if accountType := strings.ToUpper(strings.TrimSpace(filter.Type)); accountType != "" {
			where = append(where, "f.account_type=?")
			args = append(args, accountType)
		}
	}

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			f.id,f.account_type,f.account_name,f.bank_name,f.account_last4,
			f.ledger_account_id,f.enabled,a.account_code,a.account_name AS ledger_account_name
		FROM fin_fund_account f
		INNER JOIN fin_account a ON a.id=f.ledger_account_id
		WHERE %s
		ORDER BY f.account_type,f.id`, strings.Join(where, " AND ")), args...)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资金账户失败，请确认已执行 fin_init.sql")
	}

	items := make([]*FundAccount, 0, len(records))
	for _, row := range records {
		items = append(items, &FundAccount{
			Id:                row["id"].Int64(),
			Type:              row["account_type"].String(),
			Name:              row["account_name"].String(),
			BankName:          row["bank_name"].String(),
			AccountLast4:      row["account_last4"].String(),
			LedgerAccountId:   row["ledger_account_id"].Int64(),
			LedgerAccountCode: row["account_code"].String(),
			LedgerAccountName: row["ledger_account_name"].String(),
			Enabled:           row["enabled"].Bool(),
		})
	}
	return items, nil
}

// SaveFundAccount 新增或修改资金账户。
func (s *Service) SaveFundAccount(ctx context.Context, form *FundAccountForm) (int64, error) {
	accountType := strings.ToUpper(strings.TrimSpace(form.Type))
	name := strings.TrimSpace(form.Name)
	bankName := strings.TrimSpace(form.BankName)
	last4 := strings.TrimSpace(form.AccountLast4)
	ledgerCode := strings.TrimSpace(form.LedgerAccountCode)

	if accountType != "CASH" && accountType != "BANK" {
		return 0, gerror.New("资金账户类型不正确")
	}
	if name == "" || len([]rune(name)) > 100 {
		return 0, gerror.New("账户名称不能为空且不能超过100个字符")
	}
	if accountType == "BANK" && bankName == "" {
		return 0, gerror.New("银行账户必须填写开户银行")
	}
	if last4 != "" && !last4Pattern.MatchString(last4) {
		return 0, gerror.New("卡号末四位必须是4位数字")
	}
	if accountType == "CASH" {
		bankName, last4 = "", ""
	}
	if len([]rune(bankName)) > 100 {
		return 0, gerror.New("开户银行不能超过100个字符")
	}

	ledger, err := g.DB().GetOne(ctx, `SELECT id FROM fin_account
		WHERE account_code=? AND company_id=? AND enabled=1 AND category_code='ASSET'`, ledgerCode, currentCompanyId(ctx))
	if err != nil {
		return 0, gerror.Wrap(err, "校验关联科目失败")
	}
	if ledger == nil || ledger.IsEmpty() {
		return 0, gerror.Newf("关联科目 %s 不存在、未启用或不是资产类科目", ledgerCode)
	}

	exists, err := g.DB().GetValue(ctx, `SELECT COUNT(*) FROM fin_fund_account
		WHERE account_type=? AND account_name=? AND company_id=? AND id<>?`, accountType, name, currentCompanyId(ctx), form.Id)
	if err != nil {
		return 0, gerror.Wrap(err, "校验资金账户名称失败")
	}
	if exists.Int() > 0 {
		return 0, gerror.Newf("已存在同类型且名称为 %s 的资金账户", name)
	}

	if form.Id > 0 {
		row, err := g.DB().GetOne(ctx, "SELECT id FROM fin_fund_account WHERE id=? AND company_id=?", form.Id, currentCompanyId(ctx))
		if err != nil {
			return 0, gerror.Wrap(err, "读取资金账户失败")
		}
		if row == nil || row.IsEmpty() {
			return 0, gerror.New("资金账户不存在")
		}

		_, err = g.DB().Exec(ctx, `UPDATE fin_fund_account
			SET account_type=?,account_name=?,bank_name=?,account_last4=?,ledger_account_id=?,enabled=?
			WHERE id=? AND company_id=?`, accountType, name, bankName, last4, ledger["id"].Int64(), form.Enabled, form.Id, currentCompanyId(ctx))
		if err != nil {
			return 0, gerror.Wrap(err, "更新资金账户失败")
		}
		return form.Id, nil
	}

	result, err := g.DB().Exec(ctx, `INSERT INTO fin_fund_account
		(company_id,account_type,account_name,bank_name,account_last4,ledger_account_id,enabled)
		VALUES(?,?,?,?,?,?,?)`, currentCompanyId(ctx), accountType, name, bankName, last4, ledger["id"].Int64(), form.Enabled)
	if err != nil {
		return 0, gerror.Wrap(err, "新增资金账户失败")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取资金账户ID失败")
	}
	return id, nil
}

// DeleteFundAccount 删除资金账户。
func (s *Service) DeleteFundAccount(ctx context.Context, id int64) error {
	exists, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM fin_fund_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.Wrap(err, "读取资金账户失败")
	}
	if exists.Int() == 0 {
		return gerror.New("资金账户不存在")
	}
	if _, err = g.DB().Exec(ctx, "DELETE FROM fin_fund_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx)); err != nil {
		return gerror.Wrap(err, "删除资金账户失败")
	}
	return nil
}

// ToggleFundAccount 启用或停用资金账户。
func (s *Service) ToggleFundAccount(ctx context.Context, id int64, enabled bool) error {
	exists, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM fin_fund_account WHERE id=? AND company_id=?", id, currentCompanyId(ctx))
	if err != nil {
		return gerror.Wrap(err, "读取资金账户失败")
	}
	if exists.Int() == 0 {
		return gerror.New("资金账户不存在")
	}
	if _, err = g.DB().Exec(ctx, "UPDATE fin_fund_account SET enabled=? WHERE id=? AND company_id=?", enabled, id, currentCompanyId(ctx)); err != nil {
		return gerror.Wrap(err, "更新资金账户状态失败")
	}
	return nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

// directionByCategory 根据科目类别确定余额方向。
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

// isDuplicateError 判断是否命中唯一索引冲突。
func isDuplicateError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate entry")
}
