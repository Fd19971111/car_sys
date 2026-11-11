package payout

import (
	payoutApi "engineServ/api/payout"
	payoutLogic "engineServ/internal/logic/api/payout"

	"github.com/gogf/gf/v2/os/gtime"
)

type ControllerV1 struct {
	payout *payoutLogic.Payout
}

func NewV1() payoutApi.IPayoutV1 {
	return &ControllerV1{payout: payoutLogic.New()}
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
