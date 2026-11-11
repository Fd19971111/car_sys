package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
)

func (c *ControllerV1) SelectOptions(ctx context.Context, req *v1.SelectOptionsReq) (res *v1.SelectOptionsRes, err error) {
	projects, err := c.project.SelectOptions(ctx)
	if err != nil {
		return nil, err
	}

	res = &v1.SelectOptionsRes{Projects: make([]*v1.ProjectSelectItem, 0, len(projects))}

	for _, project := range projects {
		res.Projects = append(res.Projects, &v1.ProjectSelectItem{
			Value:  project.Id,
			Label:  project.Label,
			Amount: project.Amount,
			Status: project.Status,
		})
	}

	return res, nil
}
