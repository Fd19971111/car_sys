package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.project.List(ctx, req.Keyword, req.Status, req.CustomerId, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.ProjectListItem, 0, len(list))

	for _, item := range list {
		records = append(records, &v1.ProjectListItem{
			Id:               item.Project.Id,
			CustomerId:       item.Project.CustomerId,
			CustomerName:     item.CustomerName,
			Name:             item.Project.Name,
			Address:          item.Project.Address,
			ManagerId:        item.Project.ManagerId,
			ManagerName:      item.ManagerName,
			Status:           item.Project.Status,
			Amount:           item.Project.Amount,
			StartDate:        formatDate(item.Project.StartDate),
			EndDate:          formatDate(item.Project.EndDate),
			Remark:           item.Project.Remark,
			LatestProgress:   item.LatestProgress,
			LatestProgressAt: formatDate(item.LatestProgressAt),
			CreatedAt:        formatDateTime(item.Project.CreatedAt),
			UpdatedAt:        formatDateTime(item.Project.UpdatedAt),
		})
	}

	return &v1.ListRes{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		Records: records,
	}, nil
}
