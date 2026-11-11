package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
	projectListLogic "engineServ/internal/logic/api/project_list"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.project.Create(ctx, &projectListLogic.Form{
		CustomerId: req.CustomerId,
		Name:       req.Name,
		Address:    req.Address,
		ManagerId:  req.ManagerId,
		Status:     req.Status,
		Amount:     req.Amount,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Remark:     req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateRes{Id: id}, nil
}
