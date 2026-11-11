package customer

import (
	"context"

	v1 "engineServ/api/customer/v1"
)

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	info, projects, err := c.customer.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res = &v1.DetailRes{
		Id:        info.Id,
		Name:      info.Name,
		Phone:     info.Phone,
		Address:   info.Address,
		Source:    info.Source,
		Status:    info.Status,
		Remark:    info.Remark,
		CreatedAt: formatDateTime(info.CreatedAt),
		UpdatedAt: formatDateTime(info.UpdatedAt),
		Projects:  make([]*v1.CustomerProject, 0, len(projects)),
	}

	for _, project := range projects {
		res.Projects = append(res.Projects, &v1.CustomerProject{
			Id:               project.Project.Id,
			Name:             project.Project.Name,
			Status:           project.Project.Status,
			Amount:           project.Project.Amount,
			ManagerName:      project.ManagerName,
			StartDate:        formatDate(project.Project.StartDate),
			EndDate:          formatDate(project.Project.EndDate),
			LatestProgress:   project.LatestProgress,
			LatestProgressAt: formatDate(project.LatestProgressAt),
		})
	}

	return res, nil
}
