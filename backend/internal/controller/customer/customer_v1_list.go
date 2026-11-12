package customer

import (
	"context"

	v1 "engineServ/api/customer/v1"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.customer.List(ctx, req.Keyword, req.Source, req.Status, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.CustomerListItem, 0, len(list))

	for _, item := range list {
		records = append(records, &v1.CustomerListItem{
			Id:           item.Id,
			Name:         item.Name,
			Phone:        item.Phone,
			Address:      item.Address,
			Source:       item.Source,
			Status:       item.Status,
			Remark:       item.Remark,
			CreatedAt:    formatDateTime(item.CreatedAt),
			UpdatedAt:    formatDateTime(item.UpdatedAt),
		})
	}

	return &v1.ListRes{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		Records: records,
	}, nil
}
