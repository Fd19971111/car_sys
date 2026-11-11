package quote

import (
	"context"

	v1 "engineServ/api/quote/v1"
)

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	detail, err := c.quote.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res = &v1.DetailRes{
		Id:          detail.Quote.Id,
		ProjectId:   detail.Quote.ProjectId,
		ItemName:    detail.Quote.ItemName,
		Unit:        detail.Quote.Unit,
		Price:       detail.Quote.Price,
		Quantity:    detail.Quote.Quantity,
		Total:       detail.Quote.Total,
		Remark:      detail.Quote.Remark,
		CreatedAt:   formatDateTime(detail.Quote.CreatedAt),
		UpdatedAt:   formatDateTime(detail.Quote.UpdatedAt),
		QuoteAmount: detail.QuoteAmount,
	}

	if detail.Project != nil {
		res.ProjectName = detail.Project.Name
		res.ProjectAmount = detail.Project.Amount
	}

	if detail.Customer != nil {
		res.CustomerName = detail.Customer.Name
	}

	return res, nil
}
