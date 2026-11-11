package fin_init

import (
	"context"

	finInitApi "engineServ/api/fin_init"
	v1 "engineServ/api/fin_init/v1"
	"engineServ/internal/consts"
	finInitLogic "engineServ/internal/logic/api/fin_init"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func NewV1() finInitApi.IFinInitV1 { return &Controller{} }

func (c *Controller) Save(ctx context.Context, req *v1.SaveReq) (*v1.SaveRes, error) {
	role := gconv.Uint(ctx.Value(consts.CtxRoleKey))
	if role != 1 {
		return nil, gerror.New("仅公司超级管理员可初始化或重新配置财务")
	}
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))
	if err := finInitLogic.New().Save(ctx, userId, req); err != nil {
		return nil, err
	}
	return &v1.SaveRes{Initialized: true}, nil
}
