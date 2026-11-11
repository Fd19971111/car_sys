package customer

import (
	"context"

	v1 "engineServ/api/customer/v1"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error) {
	total, following, deal, lost, dealAmount, err := c.customer.Stats(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.StatsRes{
		Total:      total,
		Following:  following,
		Deal:       deal,
		Lost:       lost,
		DealAmount: dealAmount,
	}, nil
}
