package quote

import (
	"context"

	v1 "engineServ/api/quote/v1"
	quoteLogic "engineServ/internal/logic/api/quote"
)

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	warning, err := c.quote.Update(ctx, req.Id, &quoteLogic.Form{
		ProjectId: req.ProjectId,
		ItemName:  req.ItemName,
		Unit:      req.Unit,
		Price:     req.Price,
		Quantity:  req.Quantity,
		Remark:    req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateRes{Id: req.Id, Warning: warning}, nil
}
