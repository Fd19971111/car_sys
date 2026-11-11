package fin_voucher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"engineServ/internal/logic/fin/voucher"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ListFilter 凭证列表筛选条件。
type ListFilter struct {
	PeriodCode string
	StatusCode string
	Keyword    string
	Current    int
	Size       int
}

// VoucherRow 凭证列表行。
type VoucherRow struct {
	Id              int64
	VoucherNo       string
	VoucherDate     string
	PeriodCode      string
	VoucherTypeCode string
	Summary         string
	StatusCode      string
	SourceType      string
	SourceId        int64
	DebitTotal      float64
	CreditTotal     float64
	CreatedAt       string
}

// VoucherEntryRow 凭证分录。
type VoucherEntryRow struct {
	LineNo        int
	AccountCode   string
	AccountName   string
	Summary       string
	DirectionCode string
	Amount        float64
}

// VoucherDetail 凭证详情。
type VoucherDetail struct {
	VoucherRow
	Entries     []*VoucherEntryRow
	DebitTotal  float64
	CreditTotal float64
}

// PeriodRow 会计期间。
type PeriodRow struct {
	PeriodCode   string
	PeriodYear   int
	PeriodMonth  int
	StatusCode   string
	VoucherCount int
	DebitTotal   float64
	CreditTotal  float64
	CanClose     bool
	CanReopen    bool
	BeforeStart  bool
}

// PeriodList 会计期间列表。
type PeriodList struct {
	Initialized   bool
	CurrentPeriod string
	StartPeriod   string
	Records       []*PeriodRow
}

// ExportFilter 凭证导出筛选条件。
type ExportFilter struct {
	PeriodCode string
	StatusCode string
	Keyword    string
}

// ExportRow 凭证导出行（一条分录一行）。
type ExportRow struct {
	VoucherNo       string
	VoucherDate     string
	PeriodCode      string
	VoucherTypeCode string
	VoucherSummary  string
	StatusCode      string
	SourceType      string
	SourceId        int64
	LineNo          int
	AccountCode     string
	AccountName     string
	EntrySummary    string
	DirectionCode   string
	Amount          float64
}

// Service 记账凭证服务。
type Service struct{}

// New 创建记账凭证服务。
func New() *Service { return &Service{} }

// List 凭证分页列表。
func (s *Service) List(ctx context.Context, filter *ListFilter) ([]*VoucherRow, int, error) {
	current := filter.Current
	if current < 1 {
		current = 1
	}
	size := filter.Size
	if size < 1 || size > 100 {
		size = 10
	}

	where, args := voucherWhere(filter.PeriodCode, filter.StatusCode, filter.Keyword)

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM fin_voucher v WHERE "+where, args...)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "统计凭证数量失败")
	}
	if total.Int() == 0 {
		return []*VoucherRow{}, 0, nil
	}

	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)
	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			v.id,v.voucher_no,
			DATE_FORMAT(v.voucher_date,'%%Y-%%m-%%d') AS voucher_date,
			v.period_code,v.voucher_type_code,v.summary,v.status_code,v.source_type,v.source_id,
			DATE_FORMAT(v.created_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS created_at,
			COALESCE(SUM(CASE WHEN e.direction_code='DEBIT' THEN e.amount ELSE 0 END),0) AS debit_total,
			COALESCE(SUM(CASE WHEN e.direction_code='CREDIT' THEN e.amount ELSE 0 END),0) AS credit_total
		FROM fin_voucher v
		LEFT JOIN fin_voucher_entry e ON e.voucher_id=v.id
		WHERE %s
		GROUP BY v.id,v.voucher_no,v.voucher_date,v.period_code,v.voucher_type_code,v.summary,
			v.status_code,v.source_type,v.source_id,v.created_at
		ORDER BY v.voucher_date DESC,v.id DESC
		LIMIT ?,?`, where), queryArgs...)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询凭证列表失败")
	}

	items := make([]*VoucherRow, 0, len(records))
	for _, record := range records {
		items = append(items, &VoucherRow{
			Id:              record["id"].Int64(),
			VoucherNo:       record["voucher_no"].String(),
			VoucherDate:     record["voucher_date"].String(),
			PeriodCode:      record["period_code"].String(),
			VoucherTypeCode: record["voucher_type_code"].String(),
			Summary:         record["summary"].String(),
			StatusCode:      record["status_code"].String(),
			SourceType:      record["source_type"].String(),
			SourceId:        record["source_id"].Int64(),
			DebitTotal:      round2(record["debit_total"].Float64()),
			CreditTotal:     round2(record["credit_total"].Float64()),
			CreatedAt:       record["created_at"].String(),
		})
	}

	return items, total.Int(), nil
}

// Detail 凭证详情（含分录）。
func (s *Service) Detail(ctx context.Context, id int64) (*VoucherDetail, error) {
	record, err := g.DB().GetOne(ctx, `SELECT
			v.id,v.voucher_no,
			DATE_FORMAT(v.voucher_date,'%Y-%m-%d') AS voucher_date,
			v.period_code,v.voucher_type_code,v.summary,v.status_code,v.source_type,v.source_id,
			DATE_FORMAT(v.created_at,'%Y-%m-%d %H:%i:%s') AS created_at
		FROM fin_voucher v WHERE v.id=?`, id)
	if err != nil {
		return nil, gerror.Wrap(err, "查询凭证失败")
	}
	if record == nil || record.IsEmpty() {
		return nil, gerror.New("凭证不存在")
	}

	detail := &VoucherDetail{VoucherRow: VoucherRow{
		Id:              record["id"].Int64(),
		VoucherNo:       record["voucher_no"].String(),
		VoucherDate:     record["voucher_date"].String(),
		PeriodCode:      record["period_code"].String(),
		VoucherTypeCode: record["voucher_type_code"].String(),
		Summary:         record["summary"].String(),
		StatusCode:      record["status_code"].String(),
		SourceType:      record["source_type"].String(),
		SourceId:        record["source_id"].Int64(),
		CreatedAt:       record["created_at"].String(),
	}}

	entries, err := g.DB().GetAll(ctx, `SELECT line_no,account_code,account_name,summary,direction_code,amount
		FROM fin_voucher_entry WHERE voucher_id=? ORDER BY line_no`, id)
	if err != nil {
		return nil, gerror.Wrap(err, "查询凭证分录失败")
	}

	detail.Entries = make([]*VoucherEntryRow, 0, len(entries))
	for _, entry := range entries {
		amount := round2(entry["amount"].Float64())
		detail.Entries = append(detail.Entries, &VoucherEntryRow{
			LineNo:        entry["line_no"].Int(),
			AccountCode:   entry["account_code"].String(),
			AccountName:   entry["account_name"].String(),
			Summary:       entry["summary"].String(),
			DirectionCode: entry["direction_code"].String(),
			Amount:        amount,
		})
		if entry["direction_code"].String() == voucher.DirectionDebit {
			detail.DebitTotal += amount
		} else {
			detail.CreditTotal += amount
		}
	}
	detail.DebitTotal = round2(detail.DebitTotal)
	detail.CreditTotal = round2(detail.CreditTotal)

	return detail, nil
}

// Periods 会计期间列表，并计算每个期间是否可结账 / 可反结账。
func (s *Service) Periods(ctx context.Context) (*PeriodList, error) {
	result := &PeriodList{
		CurrentPeriod: time.Now().Format("2006-01"),
		Records:       make([]*PeriodRow, 0),
	}

	initialized, err := g.DB().GetValue(ctx, "SELECT is_initialized FROM fin_setup_config WHERE id=1")
	if err != nil {
		return nil, gerror.Wrap(err, "读取财务初始化状态失败")
	}
	result.Initialized = initialized.Int() == 1
	if !result.Initialized {
		return result, nil
	}

	startPeriod, err := setupStartPeriod(ctx)
	if err != nil {
		return nil, err
	}
	result.StartPeriod = startPeriod

	records, err := g.DB().GetAll(ctx, `SELECT
			p.period_code,p.period_year,p.period_month,p.status_code,
			COUNT(DISTINCT CASE WHEN v.status_code='POSTED' THEN v.id END) AS voucher_count,
			COALESCE(SUM(CASE WHEN v.status_code='POSTED' AND e.direction_code='DEBIT' THEN e.amount ELSE 0 END),0) AS debit_total,
			COALESCE(SUM(CASE WHEN v.status_code='POSTED' AND e.direction_code='CREDIT' THEN e.amount ELSE 0 END),0) AS credit_total
		FROM fin_accounting_period p
		LEFT JOIN fin_voucher v ON v.period_code=p.period_code
		LEFT JOIN fin_voucher_entry e ON e.voucher_id=v.id
		GROUP BY p.id,p.period_code,p.period_year,p.period_month,p.status_code
		ORDER BY p.period_code`)
	if err != nil {
		return nil, gerror.Wrap(err, "查询会计期间失败")
	}

	for _, record := range records {
		result.Records = append(result.Records, &PeriodRow{
			PeriodCode:   record["period_code"].String(),
			PeriodYear:   record["period_year"].Int(),
			PeriodMonth:  record["period_month"].Int(),
			StatusCode:   record["status_code"].String(),
			VoucherCount: record["voucher_count"].Int(),
			DebitTotal:   round2(record["debit_total"].Float64()),
			CreditTotal:  round2(record["credit_total"].Float64()),
		})
	}

	// 结账：只要求「更早的、已有凭证的期间」都已结账；空期间会在结账时自动关闭
	for index, item := range result.Records {
		item.BeforeStart = item.PeriodCode < startPeriod
		if item.StatusCode == "CLOSED" {
			continue
		}
		blocked := false
		for before := 0; before < index; before++ {
			prev := result.Records[before]
			if prev.StatusCode != "CLOSED" && prev.VoucherCount > 0 {
				blocked = true
				break
			}
		}
		item.CanClose = !item.BeforeStart && !blocked && item.VoucherCount > 0
	}

	// 反结账：只要求「更晚的、已有凭证的期间」都已反结账
	for index, item := range result.Records {
		if item.StatusCode != "CLOSED" {
			continue
		}
		blocked := false
		for after := index + 1; after < len(result.Records); after++ {
			next := result.Records[after]
			if next.StatusCode == "CLOSED" && next.VoucherCount > 0 {
				blocked = true
				break
			}
		}
		item.CanReopen = !item.BeforeStart && !blocked
	}

	return result, nil
}

// Close 会计期间结账：锁定期间、自动关闭更早的空期间、置为 CLOSED 并重算余额。
func (s *Service) Close(ctx context.Context, periodCode string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		startPeriod, err := setupStartPeriodTx(ctx, tx)
		if err != nil {
			return err
		}
		if periodCode < startPeriod {
			return gerror.Newf("会计期间 %s 早于财务启用期间 %s，不能结账", periodCode, startPeriod)
		}

		row, err := tx.GetOne("SELECT status_code FROM fin_accounting_period WHERE period_code=? FOR UPDATE", periodCode)
		if err != nil {
			return gerror.Wrap(err, "读取会计期间失败")
		}
		if row == nil || row.IsEmpty() {
			return gerror.Newf("会计期间 %s 不存在", periodCode)
		}
		if row["status_code"].String() == "CLOSED" {
			return gerror.Newf("会计期间 %s 已结账", periodCode)
		}
		if row["status_code"].String() == "FUTURE" {
			return gerror.Newf("会计期间 %s 还没有记账，无需结账", periodCode)
		}

		blocker, err := earlierBlockingPeriod(tx, startPeriod, periodCode)
		if err != nil {
			return err
		}
		if blocker != "" {
			return gerror.Newf("请先结账更早的期间 %s，再结账 %s", blocker, periodCode)
		}

		// 更早的期间如果没有任何已记账凭证，顺手关闭，避免空期间卡住月结
		if _, err = tx.Exec(`UPDATE fin_accounting_period p
			SET p.status_code='CLOSED'
			WHERE p.period_code<? AND p.period_code>=? AND p.status_code<>'CLOSED'
			AND NOT EXISTS(SELECT 1 FROM fin_voucher v WHERE v.period_code=p.period_code AND v.status_code='POSTED')`,
			periodCode, startPeriod); err != nil {
			return gerror.Wrap(err, "关闭空期间失败")
		}

		if _, err = tx.Exec("UPDATE fin_accounting_period SET status_code='CLOSED' WHERE period_code=?", periodCode); err != nil {
			return gerror.Wrap(err, "结账失败")
		}

		return voucher.New().RebuildWithTx(ctx, tx)
	})
}

// Reopen 会计期间反结账：只允许反结账最近一个有凭证的已结账期间。
func (s *Service) Reopen(ctx context.Context, periodCode string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		startPeriod, err := setupStartPeriodTx(ctx, tx)
		if err != nil {
			return err
		}
		if periodCode < startPeriod {
			return gerror.Newf("会计期间 %s 早于财务启用期间 %s，不能反结账", periodCode, startPeriod)
		}

		row, err := tx.GetOne("SELECT status_code FROM fin_accounting_period WHERE period_code=? FOR UPDATE", periodCode)
		if err != nil {
			return gerror.Wrap(err, "读取会计期间失败")
		}
		if row == nil || row.IsEmpty() {
			return gerror.Newf("会计期间 %s 不存在", periodCode)
		}
		if row["status_code"].String() != "CLOSED" {
			return gerror.Newf("会计期间 %s 当前不是已结账状态", periodCode)
		}

		blocker, err := laterBlockingPeriod(tx, periodCode)
		if err != nil {
			return err
		}
		if blocker != "" {
			return gerror.Newf("请先反结账更晚的期间 %s，再反结账 %s", blocker, periodCode)
		}

		if _, err = tx.Exec("UPDATE fin_accounting_period SET status_code='OPEN' WHERE period_code=?", periodCode); err != nil {
			return gerror.Wrap(err, "反结账失败")
		}

		return voucher.New().RebuildWithTx(ctx, tx)
	})
}

// Export 导出凭证分录明细（最多 5000 行）。
func (s *Service) Export(ctx context.Context, filter *ExportFilter) ([]*ExportRow, error) {
	where, args := voucherWhere(filter.PeriodCode, filter.StatusCode, filter.Keyword)

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			v.voucher_no,
			DATE_FORMAT(v.voucher_date,'%%Y-%%m-%%d') AS voucher_date,
			v.period_code,v.voucher_type_code,v.summary AS voucher_summary,v.status_code,v.source_type,v.source_id,
			e.line_no,e.account_code,e.account_name,e.summary AS entry_summary,e.direction_code,e.amount
		FROM fin_voucher v
		INNER JOIN fin_voucher_entry e ON e.voucher_id=v.id
		WHERE %s
		ORDER BY v.voucher_date,v.voucher_no,e.line_no
		LIMIT 5000`, where), args...)
	if err != nil {
		return nil, gerror.Wrap(err, "导出凭证失败")
	}

	rows := make([]*ExportRow, 0, len(records))
	for _, record := range records {
		rows = append(rows, &ExportRow{
			VoucherNo:       record["voucher_no"].String(),
			VoucherDate:     record["voucher_date"].String(),
			PeriodCode:      record["period_code"].String(),
			VoucherTypeCode: record["voucher_type_code"].String(),
			VoucherSummary:  record["voucher_summary"].String(),
			StatusCode:      record["status_code"].String(),
			SourceType:      record["source_type"].String(),
			SourceId:        record["source_id"].Int64(),
			LineNo:          record["line_no"].Int(),
			AccountCode:     record["account_code"].String(),
			AccountName:     record["account_name"].String(),
			EntrySummary:    record["entry_summary"].String(),
			DirectionCode:   record["direction_code"].String(),
			Amount:          round2(record["amount"].Float64()),
		})
	}

	return rows, nil
}

// voucherWhere 拼接凭证筛选条件。
func voucherWhere(periodCode, statusCode, keyword string) (string, []interface{}) {
	where := []string{"1=1"}
	args := make([]interface{}, 0, 4)

	if period := strings.TrimSpace(periodCode); period != "" {
		where = append(where, "v.period_code=?")
		args = append(args, period)
	}
	if status := strings.TrimSpace(statusCode); status != "" {
		where = append(where, "v.status_code=?")
		args = append(args, status)
	}
	if text := strings.TrimSpace(keyword); text != "" {
		like := "%" + text + "%"
		where = append(where, "(v.voucher_no LIKE ? OR v.summary LIKE ?)")
		args = append(args, like, like)
	}

	return strings.Join(where, " AND "), args
}

// earlierBlockingPeriod 返回更早的、有凭证但未结账的期间。
func earlierBlockingPeriod(tx gdb.TX, startPeriod, periodCode string) (string, error) {
	value, err := tx.GetValue(`SELECT p.period_code FROM fin_accounting_period p
		WHERE p.period_code<? AND p.period_code>=? AND p.status_code<>'CLOSED'
		AND EXISTS(SELECT 1 FROM fin_voucher v WHERE v.period_code=p.period_code AND v.status_code='POSTED')
		ORDER BY p.period_code LIMIT 1`, periodCode, startPeriod)
	if err != nil {
		return "", gerror.Wrap(err, "检查更早期间状态失败")
	}
	return value.String(), nil
}

// laterBlockingPeriod 返回更晚的、有凭证且已结账的期间。
func laterBlockingPeriod(tx gdb.TX, periodCode string) (string, error) {
	value, err := tx.GetValue(`SELECT p.period_code FROM fin_accounting_period p
		WHERE p.period_code>? AND p.status_code='CLOSED'
		AND EXISTS(SELECT 1 FROM fin_voucher v WHERE v.period_code=p.period_code AND v.status_code='POSTED')
		ORDER BY p.period_code LIMIT 1`, periodCode)
	if err != nil {
		return "", gerror.Wrap(err, "检查更晚期间状态失败")
	}
	return value.String(), nil
}

// setupStartPeriod 查询财务启用期间。
func setupStartPeriod(ctx context.Context) (string, error) {
	row, err := g.DB().GetOne(ctx, "SELECT is_initialized,period_year,start_month FROM fin_setup_config WHERE id=1")
	if err != nil {
		return "", gerror.Wrap(err, "读取财务初始化状态失败")
	}
	return startPeriodFromRow(row)
}

// setupStartPeriodTx 在事务内查询财务启用期间。
func setupStartPeriodTx(ctx context.Context, tx gdb.TX) (string, error) {
	row, err := tx.GetOne("SELECT is_initialized,period_year,start_month FROM fin_setup_config WHERE id=1")
	if err != nil {
		return "", gerror.Wrap(err, "读取财务初始化状态失败")
	}
	return startPeriodFromRow(row)
}

// startPeriodFromRow 从初始化配置中解析启用期间。
func startPeriodFromRow(row gdb.Record) (string, error) {
	if row == nil || row.IsEmpty() || row["is_initialized"].Int() != 1 || row["period_year"].Int() < 2000 {
		return "", gerror.New("财务尚未初始化，请先完成财务基础配置")
	}
	return fmt.Sprintf("%04d-%02d", row["period_year"].Int(), row["start_month"].Int()), nil
}

// round2 保留两位小数。
func round2(value float64) float64 {
	return float64(int64(value*100+0.5)) / 100
}
