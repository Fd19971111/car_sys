package mp

import (
	"context"

	v1 "engineServ/api/mp/v1"
)

func (c *ControllerV1) ProgressCreate(ctx context.Context, req *v1.ProgressCreateReq) (res *v1.ProgressCreateRes, err error) {
	id, err := c.mp.CreateProgress(ctx, req.ProjectId, currentCustomerId(ctx), req.Content, req.Images)
	if err != nil {
		return nil, err
	}

	return &v1.ProgressCreateRes{Id: id}, nil
}
