package customer

import (
	"context"

	v1 "engineServ/api/customer/v1"
	defaultDo "engineServ/model/default/model/do"
)

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.customer.Update(ctx, req.Id, &defaultDo.Customer{
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

	return &v1.UpdateRes{Id: req.Id}, nil
}
