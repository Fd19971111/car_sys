package voucher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
)

// 会计期间辅助函数。
//
// 初始化时只会生成「启用年度」的 12 个期间，如果业务跨年（例如 2027 年继续记账），
// 旧实现会因为找不到期间而直接报错。现在改为按年自动补齐：记账时发现目标月份
// 不存在，就补齐该年度的 12 个期间（默认 FUTURE，首次记账时自动开启）。

// startPeriodCode 返回财务启用期间（YYYY-MM）。
func startPeriodCode(tx gdb.TX, companyId uint) (string, error) {
	row, err := tx.GetOne("SELECT period_year,start_month FROM fin_setup_config WHERE company_id=?", companyId)
	if err != nil {
		return "", gerror.Wrap(err, "读取财务启用期间失败")
	}
	if row == nil || row.IsEmpty() || row["period_year"].Int() < 2000 {
		return "", gerror.New("财务尚未初始化，请先完成财务基础配置")
	}

	return fmt.Sprintf("%04d-%02d", row["period_year"].Int(), row["start_month"].Int()), nil
}

// ensurePeriodYear 补齐指定期间所在年度的 12 个会计期间。
//
// 已存在的期间不会被覆盖（ON DUPLICATE KEY 只做无意义赋值），
// 因此已结账 / 已启用状态都能保留。
func ensurePeriodYear(tx gdb.TX, periodCode string, companyId uint) error {
	if len(periodCode) < 7 {
		return gerror.Newf("会计期间 %s 格式不正确", periodCode)
	}
	year, err := strconv.Atoi(periodCode[:4])
	if err != nil {
		return gerror.Newf("会计期间 %s 格式不正确", periodCode)
	}

	var builder strings.Builder
	builder.WriteString("INSERT INTO fin_accounting_period (company_id,period_code,period_year,period_month,status_code) VALUES ")

	args := make([]interface{}, 0, 48)
	for month := 1; month <= 12; month++ {
		if month > 1 {
			builder.WriteString(",")
		}
		builder.WriteString("(?,?,?,?,'FUTURE')")
		args = append(args, companyId, fmt.Sprintf("%04d-%02d", year, month), year, month)
	}
	builder.WriteString(" ON DUPLICATE KEY UPDATE period_code=VALUES(period_code)")

	if _, err = tx.Exec(builder.String(), args...); err != nil {
		return gerror.Wrapf(err, "初始化 %d 年度会计期间失败", year)
	}
	return nil
}
