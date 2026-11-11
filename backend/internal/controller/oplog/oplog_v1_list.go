package oplog

import (
	"context"

	v1 "engineServ/api/oplog/v1"
	oplogLogic "engineServ/internal/logic/oplog"
)

// List 操作日志分页列表。
func (c *ControllerV1) List(ctx context.Context, req *v1.OpLogListReq) (res *v1.OpLogListRes, err error) {
	items, total, modules, err := c.oplog.List(ctx, &oplogLogic.Filter{
		Module:    req.Module,
		Action:    req.Action,
		Keyword:   req.Keyword,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Current:   req.Current,
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*v1.OpLogItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.OpLogItem{
			Id:         item.Id,
			CreatedAt:  item.CreatedAt,
			UserId:     item.UserId,
			UserName:   item.UserName,
			Module:     item.Module,
			Action:     item.Action,
			ActionName: item.ActionName,
			TargetType: item.TargetType,
			TargetId:   item.TargetId,
			Title:      item.Title,
			Amount:     item.Amount,
			Detail:     item.Detail,
			Ip:         item.Ip,
		})
	}

	return &v1.OpLogListRes{Records: records, Total: total, Modules: modules}, nil
}
