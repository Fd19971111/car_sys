package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error) {
	total, pending, building, finished, settled, totalAmount, err := c.project.Stats(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.StatsRes{
		Total:       total,
		Pending:     pending,
		Building:    building,
		Finished:    finished,
		Settled:     settled,
		TotalAmount: totalAmount,
	}, nil
}
