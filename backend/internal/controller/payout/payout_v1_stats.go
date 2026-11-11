package payout

import (
	"context"

	v1 "engineServ/api/payout/v1"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error) {
	totalCount, totalAmount, monthAmount, projectCount, err := c.payout.Stats(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.StatsRes{
		TotalCount:   totalCount,
		TotalAmount:  totalAmount,
		MonthAmount:  monthAmount,
		ProjectCount: projectCount,
	}, nil
}
