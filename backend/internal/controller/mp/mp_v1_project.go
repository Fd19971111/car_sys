package mp

import (
	"context"

	v1 "engineServ/api/mp/v1"
	uploadLogic "engineServ/internal/logic/upload"
)

func (c *ControllerV1) Project(ctx context.Context, req *v1.ProjectReq) (res *v1.ProjectRes, err error) {
	snapshot, err := c.mp.Project(ctx, req.Id, currentCustomerId(ctx))
	if err != nil {
		return nil, err
	}

	project := snapshot.Project

	res = &v1.ProjectRes{
		Project: &v1.ProjectInfo{
			Id:           project.Id,
			Name:         project.Name,
			Address:      project.Address,
			Status:       project.Status,
			ManagerName:  snapshot.ManagerName,
			ManagerPhone: snapshot.ManagerPhone,
			StartDate:    formatDate(project.StartDate),
			EndDate:      formatDate(project.EndDate),
			Remark:       project.Remark,
			CreatedAt:    formatDate(project.CreatedAt),
		},
		Owner: toOwnerProfile(snapshot.Owner),
		Summary: &v1.ProjectSummary{
			ContractAmount:  round2(project.Amount),
			QuoteAmount:     round2(snapshot.QuoteAmount),
			PaidAmount:      round2(snapshot.PaidAmount),
			UnpaidAmount:    round2(subAmount(project.Amount, snapshot.PaidAmount)),
			PaidPercent:     percent(snapshot.PaidAmount, project.Amount),
			ProgressPercent: snapshot.ProgressPercent,
		},
		Logs:     make([]*v1.ProgressLog, 0, len(snapshot.Logs)),
		Quotes:   make([]*v1.QuoteItem, 0, len(snapshot.Quotes)),
		Payments: make([]*v1.PaymentItem, 0, len(snapshot.Payments)),
	}

	for _, item := range snapshot.Logs {
		res.Logs = append(res.Logs, &v1.ProgressLog{
			Id:        item.Log.Id,
			Content:   item.Log.Content,
			Images:    uploadLogic.AbsoluteURLs(ctx, parseImages(item.Log.Images)),
			UserName:  item.UserName,
			CreatedAt: formatDateTime(item.Log.CreatedAt),
		})
	}

	for _, item := range snapshot.Quotes {
		res.Quotes = append(res.Quotes, &v1.QuoteItem{
			Id:       item.Id,
			ItemName: item.ItemName,
			Unit:     item.Unit,
			Price:    item.Price,
			Quantity: item.Quantity,
			Total:    item.Total,
			Remark:   item.Remark,
		})
	}

	for _, item := range snapshot.Payments {
		res.Payments = append(res.Payments, &v1.PaymentItem{
			Id:          item.Id,
			PaymentType: item.PaymentType,
			Amount:      item.Amount,
			PaymentDate: formatDate(item.PaymentDate),
			Remark:      item.Remark,
		})
	}

	return res, nil
}
