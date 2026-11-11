package mp

import (
	"context"
	"strings"

	v1 "engineServ/api/mp/v1"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	// 本接口在 CustomerAuth 之后执行，令牌必然是有效的；直接删掉即可
	token := strings.TrimSpace(strings.TrimPrefix(g.RequestFromCtx(ctx).Header.Get("Authorization"), "Bearer "))

	if err = c.auth.Logout(ctx, token); err != nil {
		return nil, err
	}

	return &v1.LogoutRes{}, nil
}
