package project_data

import (
	"context"

	v1 "engineServ/api/project_data/v1"
	projectDataLogic "engineServ/internal/logic/api/project_data"
	uploadLogic "engineServ/internal/logic/upload"
)

// ChangeList 变更记录列表。
func (c *ControllerV1) ChangeList(ctx context.Context, req *v1.ChangeListReq) (*v1.ChangeListRes, error) {
	result, err := c.data.ListChanges(ctx, req.ProjectId, req.ChangeType, req.Keyword, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.ChangeItem, 0, len(result.Records))
	for _, item := range result.Records {
		records = append(records, &v1.ChangeItem{
			Id:         item.Id,
			ProjectId:  item.ProjectId,
			ChangeType: item.ChangeType,
			Title:      item.Title,
			Amount:     item.Amount,
			ChangeDate: item.ChangeDate,
			Reason:     item.Reason,
			Images:     item.Images,
			ImageUrls:  uploadLogic.AbsoluteURLs(ctx, item.Images),
			Remark:     item.Remark,
			CreatedAt:  item.CreatedAt,
		})
	}

	return &v1.ChangeListRes{
		Current:     req.Current,
		Size:        req.Size,
		Total:       result.Total,
		AddTotal:    result.AddTotal,
		ReduceTotal: result.ReduceTotal,
		Records:     records,
	}, nil
}

// ChangeCreate 新增变更记录。
func (c *ControllerV1) ChangeCreate(ctx context.Context, req *v1.ChangeCreateReq) (*v1.ChangeCreateRes, error) {
	id, err := c.data.CreateChange(ctx, &projectDataLogic.ChangeInput{
		ProjectId:  req.ProjectId,
		ChangeType: req.ChangeType,
		Title:      req.Title,
		Amount:     req.Amount,
		ChangeDate: req.ChangeDate,
		Reason:     req.Reason,
		Images:     req.Images,
		Remark:     req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ChangeCreateRes{Id: id}, nil
}

// ChangeUpdate 更新变更记录。
func (c *ControllerV1) ChangeUpdate(ctx context.Context, req *v1.ChangeUpdateReq) (*v1.ChangeUpdateRes, error) {
	err := c.data.UpdateChange(ctx, req.Id, &projectDataLogic.ChangeInput{
		ProjectId:  req.ProjectId,
		ChangeType: req.ChangeType,
		Title:      req.Title,
		Amount:     req.Amount,
		ChangeDate: req.ChangeDate,
		Reason:     req.Reason,
		Images:     req.Images,
		Remark:     req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ChangeUpdateRes{Id: req.Id}, nil
}

// ChangeDelete 删除变更记录。
func (c *ControllerV1) ChangeDelete(ctx context.Context, req *v1.ChangeDeleteReq) (*v1.ChangeDeleteRes, error) {
	if err := c.data.DeleteChange(ctx, req.Id); err != nil {
		return nil, err
	}

	return &v1.ChangeDeleteRes{}, nil
}
