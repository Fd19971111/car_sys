package mp

import (
	"context"
	"math"

	v1 "engineServ/api/mp/v1"
)

func (c *ControllerV1) Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error) {
	// 业主身份只从令牌上下文取，前端不再传 customerId
	owner, cards, err := c.mp.Overview(ctx, currentCustomerId(ctx))
	if err != nil {
		return nil, err
	}

	summary := &v1.OverviewSummary{ProjectCount: len(cards)}
	projects := make([]*v1.ProjectCard, 0, len(cards))

	progressSum := 0

	for _, card := range cards {
		summary.ContractAmount += card.Project.Amount
		summary.PaidAmount += card.PaidAmount
		progressSum += card.ProgressPercent

		projects = append(projects, &v1.ProjectCard{
			Id:               card.Project.Id,
			Name:             card.Project.Name,
			Address:          card.Project.Address,
			Status:           card.Project.Status,
			Amount:           card.Project.Amount,
			PaidAmount:       card.PaidAmount,
			ManagerName:      card.ManagerName,
			StartDate:        formatDate(card.Project.StartDate),
			EndDate:          formatDate(card.Project.EndDate),
			LatestProgress:   card.LatestProgress,
			LatestProgressAt: formatDateTime(card.LatestProgressAt),
			ProgressPercent:  card.ProgressPercent,
		})
	}

	summary.ContractAmount = round2(summary.ContractAmount)
	summary.PaidAmount = round2(summary.PaidAmount)
	summary.UnpaidAmount = round2(subAmount(summary.ContractAmount, summary.PaidAmount))

	if len(cards) > 0 {
		summary.ProgressPercent = int(math.Round(float64(progressSum) / float64(len(cards))))
	}

	return &v1.OverviewRes{
		Owner:    toOwnerProfile(owner),
		Summary:  summary,
		Projects: projects,
	}, nil
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
