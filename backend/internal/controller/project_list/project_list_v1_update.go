package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
	projectListLogic "engineServ/internal/logic/api/project_list"
)

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.project.Update(ctx, req.Id, &projectListLogic.Form{
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

	return &v1.UpdateRes{Id: req.Id}, nil
}
