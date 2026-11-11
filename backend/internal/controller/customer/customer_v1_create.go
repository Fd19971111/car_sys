package customer

import (
	"context"

	v1 "engineServ/api/customer/v1"
	defaultDo "engineServ/model/default/model/do"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.customer.Create(ctx, &defaultDo.Customer{
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
		Source:  req.Source,
		Status:  req.Status,
		Remark:  req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateRes{Id: id}, nil
}
