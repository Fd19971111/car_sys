package payment

import (
	"context"

	v1 "engineServ/api/payment/v1"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error) {
	totalCount, totalAmount, monthAmount, monthCount, err := c.payment.Stats(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.StatsRes{
		TotalCount:  totalCount,
		TotalAmount: totalAmount,
		MonthAmount: monthAmount,
		MonthCount:  monthCount,
	}, nil
}
