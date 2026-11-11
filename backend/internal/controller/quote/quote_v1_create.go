package quote

import (
	"context"

	v1 "engineServ/api/quote/v1"
	quoteLogic "engineServ/internal/logic/api/quote"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, warning, err := c.quote.Create(ctx, &quoteLogic.Form{
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

	return &v1.CreateRes{Id: id, Warning: warning}, nil
}
