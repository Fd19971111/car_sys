package quote

import (
	"context"

	v1 "engineServ/api/quote/v1"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error) {
	totalCount, totalAmount, projectCount, avgAmount, err := c.quote.Stats(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.StatsRes{
		TotalCount:   totalCount,
		TotalAmount:  totalAmount,
		ProjectCount: projectCount,
		AvgAmount:    avgAmount,
	}, nil
}
