package mp

import (
	"context"
	"strings"

	mpApi "engineServ/api/mp"
	v1 "engineServ/api/mp/v1"
	"engineServ/internal/consts"
	customerAuthLogic "engineServ/internal/logic/api/customer_auth"
	mpLogic "engineServ/internal/logic/api/mp"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type ControllerV1 struct {
	mp   *mpLogic.Mp
	auth *customerAuthLogic.CustomerAuth
}

func NewV1() mpApi.IMpV1 {
	return &ControllerV1{
		mp:   mpLogic.New(),
		auth: customerAuthLogic.New(),
	}
}

// currentCustomerId 当前登录业主ID。
//
// 只从 CustomerAuth 中间件写入的上下文里取，**绝不再从请求参数接收**，
// 这是数据隔离成立的前提。
func currentCustomerId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCustomerIdKey))
}

// toOwnerProfile 客户实体 → 业主资料
func toOwnerProfile(customer *defaultEntity.Customer) *v1.OwnerProfile {
	if customer == nil {
		return nil
	}

	return &v1.OwnerProfile{
		Id:      customer.Id,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
	}
}

// percent 计算百分比并夹取到 0~100
func percent(part, total float64) int {
	if total <= 0 {
		return 0
	}

	value := int(float64(part)/total*100 + 0.5)

	if value < 0 {
		return 0
	}

	if value > 100 {
		return 100
	}

	return value
}

// subAmount 差额，负数归零（避免超额收款时显示负数）
func subAmount(total, part float64) float64 {
	diff := total - part
	if diff < 0 {
		return 0
	}

	return diff
}

// formatDate 格式化为 YYYY-MM-DD，空值返回空串。
func formatDate(t *gtime.Time) string {
	if t == nil {
		return ""
	}

	return t.Format("Y-m-d")
}

// formatDateTime 格式化为 YYYY-MM-DD HH:mm，空值返回空串。
func formatDateTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}

	return t.Format("Y-m-d H:i")
}

// parseImages 把逗号分隔的图片相对路径拆成数组
func parseImages(images string) []string {
	result := make([]string, 0)

	for _, item := range strings.Split(images, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
