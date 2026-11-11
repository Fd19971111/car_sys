package quote

import (
	"context"

	v1 "engineServ/api/quote/v1"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.quote.List(ctx, req.Keyword, req.ProjectId, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.QuoteListItem, 0, len(list))

	for _, item := range list {
		records = append(records, &v1.QuoteListItem{
			Id:           item.Quote.Id,
			ProjectId:    item.Quote.ProjectId,
			ProjectName:  item.ProjectName,
			CustomerName: item.CustomerName,
			ItemName:     item.Quote.ItemName,
			Unit:         item.Quote.Unit,
			Price:        item.Quote.Price,
			Quantity:     item.Quote.Quantity,
			Total:        item.Quote.Total,
			Remark:       item.Quote.Remark,
			CreatedAt:    formatDateTime(item.Quote.CreatedAt),
			UpdatedAt:    formatDateTime(item.Quote.UpdatedAt),
		})
	}

	return &v1.ListRes{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		Records: records,
	}, nil
}
