package payout

import (
	"context"

	v1 "engineServ/api/payout/v1"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.payout.List(ctx, req.Keyword, req.ProjectId, req.PayoutType, req.Source, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.PayoutListItem, 0, len(list))

	for _, item := range list {
		records = append(records, &v1.PayoutListItem{
			Id:              item.Payout.Id,
			ProjectId:       item.Payout.ProjectId,
			ProjectName:     item.ProjectName,
			CustomerName:    item.CustomerName,
			Amount:          item.Payout.Amount,
			PayoutType:      item.Payout.PayoutType,
			PayoutDate:      formatDate(item.Payout.PayoutDate),
			VendorName:      item.VendorName,
			FundAccountId:   item.Payout.FundAccountId,
			FundAccountName: item.FundAccountName,
			VoucherNo:       item.VoucherNo,
			Remark:          item.Payout.Remark,
			CreatedAt:       formatDateTime(item.Payout.CreatedAt),
		})
	}

	return &v1.ListRes{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		Records: records,
	}, nil
}
