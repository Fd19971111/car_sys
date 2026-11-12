package customer

import (
	"context"

	v1 "engineServ/api/customer/v1"
)

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	info, err := c.customer.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DetailRes{
		Id:        info.Id,
		Name:      info.Name,
		Phone:     info.Phone,
		Address:   info.Address,
		Source:    info.Source,
		Status:    info.Status,
		Remark:    info.Remark,
		CreatedAt: formatDateTime(info.CreatedAt),
		UpdatedAt: formatDateTime(info.UpdatedAt),
	}, nil
}
