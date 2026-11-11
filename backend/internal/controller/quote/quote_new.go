package quote

import (
	quoteApi "engineServ/api/quote"
	quoteLogic "engineServ/internal/logic/api/quote"

	"github.com/gogf/gf/v2/os/gtime"
)

type ControllerV1 struct {
	quote *quoteLogic.Quote
}

func NewV1() quoteApi.IQuoteV1 {
	return &ControllerV1{quote: quoteLogic.New()}
}

// formatDateTime 格式化为 YYYY-MM-DD HH:mm:ss，空值返回空串。
func formatDateTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}

	return t.String()
}

// formatDate 格式化为 YYYY-MM-DD，空值返回空串。
func formatDate(t *gtime.Time) string {
	if t == nil {
		return ""
	}

	return t.Format("Y-m-d")
}
