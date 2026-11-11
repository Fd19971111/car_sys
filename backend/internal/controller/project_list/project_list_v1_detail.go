package project_list

import (
	"context"

	v1 "engineServ/api/project_list/v1"
	uploadLogic "engineServ/internal/logic/upload"
)

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	project, customer, logs, err := c.project.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res = &v1.DetailRes{
		Id:         project.Id,
		CustomerId: project.CustomerId,
		Name:       project.Name,
		Address:    project.Address,
		ManagerId:  project.ManagerId,
		Status:     project.Status,
		Amount:     project.Amount,
		StartDate:  formatDate(project.StartDate),
		EndDate:    formatDate(project.EndDate),
		Remark:     project.Remark,
		CreatedAt:  formatDateTime(project.CreatedAt),
		UpdatedAt:  formatDateTime(project.UpdatedAt),
		Logs:       make([]*v1.ProjectLogItem, 0, len(logs)),
	}

	if customer != nil {
		res.CustomerName = customer.Name
		res.CustomerPhone = customer.Phone
		res.CustomerAddress = customer.Address
	}

	for _, item := range logs {
		logDate := formatDate(item.Log.LogDate)
		res.Logs = append(res.Logs, &v1.ProjectLogItem{
			Id:        item.Log.Id,
			Content:   item.Log.Content,
			Images:    uploadLogic.AbsoluteURLs(ctx, parseImages(item.Log.Images)),
			UserName:  item.UserName,
			LogDate:   logDate,
			CreatedAt: formatDateTime(item.Log.CreatedAt),
		})
	}

	// 最近一条进度即时间轴首条（按业务日期倒序），展示用业务日期
	if len(res.Logs) > 0 {
		res.LatestProgress = res.Logs[0].Content
		res.LatestProgressAt = res.Logs[0].LogDate
		if res.LatestProgressAt == "" {
			res.LatestProgressAt = res.Logs[0].CreatedAt
		}
	}

	// 负责人姓名：详情单独查一次，避免复用列表的批量逻辑
	res.ManagerName = c.project.ManagerName(ctx, project.ManagerId)

	return res, nil
}
