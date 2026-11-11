package customer

import (
	customerApi "engineServ/api/customer"
	customerLogic "engineServ/internal/logic/api/customer"

	"github.com/gogf/gf/v2/os/gtime"
)

type ControllerV1 struct {
	customer *customerLogic.Customer
}

func NewV1() customerApi.ICustomerV1 {
	return &ControllerV1{customer: customerLogic.New()}
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
