package project_data

import (
	"context"

	v1 "engineServ/api/project_data/v1"
	"engineServ/internal/consts"
	uploadLogic "engineServ/internal/logic/upload"

	"github.com/gogf/gf/v2/util/gconv"
)

// SiteLogList 施工日志列表。
func (c *ControllerV1) SiteLogList(ctx context.Context, req *v1.SiteLogListReq) (*v1.SiteLogListRes, error) {
	items, total, err := c.data.ListSiteLogs(ctx, req.ProjectId, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.SiteLogItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.SiteLogItem{
			Id:        item.Id,
			ProjectId: item.ProjectId,
			Content:   item.Content,
			Images:    item.Images,
			ImageUrls: uploadLogic.AbsoluteURLs(ctx, item.Images),
			LogDate:   item.LogDate,
			UserName:  item.UserName,
			CreatedAt: item.CreatedAt,
		})
	}

	return &v1.SiteLogListRes{Current: req.Current, Size: req.Size, Total: total, Records: records}, nil
}

// SiteLogCreate 新增施工日志。
func (c *ControllerV1) SiteLogCreate(ctx context.Context, req *v1.SiteLogCreateReq) (*v1.SiteLogCreateRes, error) {
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	id, err := c.data.CreateSiteLog(ctx, req.ProjectId, userId, req.Content, req.Images, req.LogDate)
	if err != nil {
		return nil, err
	}

	return &v1.SiteLogCreateRes{Id: id}, nil
}

// SiteLogUpdate 更新施工日志。
func (c *ControllerV1) SiteLogUpdate(ctx context.Context, req *v1.SiteLogUpdateReq) (*v1.SiteLogUpdateRes, error) {
	if err := c.data.UpdateSiteLog(ctx, req.Id, req.Content, req.Images); err != nil {
		return nil, err
	}

	return &v1.SiteLogUpdateRes{Id: req.Id}, nil
}

// SiteLogDelete 删除施工日志。
func (c *ControllerV1) SiteLogDelete(ctx context.Context, req *v1.SiteLogDeleteReq) (*v1.SiteLogDeleteRes, error) {
	if err := c.data.DeleteSiteLog(ctx, req.Id); err != nil {
		return nil, err
	}

	return &v1.SiteLogDeleteRes{}, nil
}
