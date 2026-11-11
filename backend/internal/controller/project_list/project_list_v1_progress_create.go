package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) ProgressCreate(ctx context.Context, req *v1.ProgressCreateReq) (res *v1.ProgressCreateRes, err error) {
	// 记录人 = 当前登录的后台用户（Auth 中间件已写入上下文）
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	id, err := c.project.CreateProgress(ctx, req.ProjectId, userId, req.Content, req.Images, req.LogDate)
	if err != nil {
		return nil, err
	}

	return &v1.ProgressCreateRes{Id: id}, nil
}
