package payment

import (
	"context"

	v1 "engineServ/api/payment/v1"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.payment.List(ctx, req.Keyword, req.ProjectId, req.PaymentType, req.Current, req.Size)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.PaymentListItem, 0, len(list))

	for _, item := range list {
		records = append(records, &v1.PaymentListItem{
			Id:              item.Payment.Id,
			ProjectId:       item.Payment.ProjectId,
			ProjectName:     item.ProjectName,
			CustomerName:    item.CustomerName,
			Amount:          item.Payment.Amount,
			PaymentType:     item.Payment.PaymentType,
			PaymentDate:     formatDate(item.Payment.PaymentDate),
			FundAccountId:   item.Payment.FundAccountId,
			FundAccountName: item.FundAccountName,
			VoucherNo:       item.VoucherNo,
			Remark:          item.Payment.Remark,
			CreatedAt:       formatDateTime(item.Payment.CreatedAt),
		})
	}

	return &v1.ListRes{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		Records: records,
	}, nil
}
