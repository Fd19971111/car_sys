package transport

import (
	"context"

	transportApi "engineServ/api/transport"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/util/gconv"
)

// Controller 运力管理控制器（司机 / 车辆）。
type Controller struct{}

// NewV1 创建运力管理控制器。
func NewV1() transportApi.ITransportV1 { return &Controller{} }

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func normalizePage(current, size int) (int, int) {
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 200 {
		size = 10
	}
	return current, size
}

func pageOffset(current, size int) int {
	return (current - 1) * size
}
