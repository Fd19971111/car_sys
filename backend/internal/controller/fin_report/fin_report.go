package fin_report

import (
	"context"

	finReportApi "engineServ/api/fin_report"
	v1 "engineServ/api/fin_report/v1"
	finReportLogic "engineServ/internal/logic/api/fin_report"
)

// ControllerV1 财务报表控制器。
type ControllerV1 struct{}

// NewV1 创建财务报表控制器。
func NewV1() finReportApi.IFinReportV1 { return &ControllerV1{} }

// Overview 财务报表总览。
func (c *ControllerV1) Overview(ctx context.Context, req *v1.ReportOverviewReq) (*v1.ReportOverviewRes, error) {
	overview, err := finReportLogic.New().Overview(ctx, req.Year)
	if err != nil {
		return nil, err
	}

	monthly := make([]*v1.ReportMonthlyItem, 0, len(overview.Monthly))
	for _, item := range overview.Monthly {
		monthly = append(monthly, &v1.ReportMonthlyItem{
			Month:   item.Month,
			Label:   item.Label,
			Income:  item.Income,
			Expense: item.Expense,
			Net:     item.Net,
		})
	}

	funds := make([]*v1.ReportFundBalance, 0, len(overview.FundAccounts))
	for _, item := range overview.FundAccounts {
		funds = append(funds, &v1.ReportFundBalance{
			Id:           item.Id,
			Type:         item.Type,
			Name:         item.Name,
			BankName:     item.BankName,
			AccountLast4: item.AccountLast4,
			AccountCode:  item.AccountCode,
			AccountName:  item.AccountName,
			Balance:      item.Balance,
		})
	}

	balances := make([]*v1.ReportAccountBalance, 0, len(overview.AccountBalances))
	for _, item := range overview.AccountBalances {
		balances = append(balances, &v1.ReportAccountBalance{
			Code:         item.Code,
			Name:         item.Name,
			CategoryCode: item.CategoryCode,
			Balance:      item.Balance,
		})
	}

	recentFlows := make([]*v1.ReportFlowItem, 0, len(overview.RecentFlows))
	for _, item := range overview.RecentFlows {
		recentFlows = append(recentFlows, &v1.ReportFlowItem{
			Id:              item.Id,
			OccurredAt:      item.OccurredAt,
			DirectionCode:   item.DirectionCode,
			Amount:          item.Amount,
			Remark:          item.Remark,
			AccountName:     item.AccountName,
			FundAccountName: item.FundAccountName,
		})
	}

	return &v1.ReportOverviewRes{
		Initialized:   overview.Initialized,
		CurrentPeriod: overview.CurrentPeriod,
		PeriodStatus:  overview.PeriodStatus,
		Year:          overview.Year,
		Summary: v1.ReportSummary{
			Income:         overview.Summary.Income,
			Expense:        overview.Summary.Expense,
			Net:            overview.Summary.Net,
			NetRate:        overview.Summary.NetRate,
			TotalBalance:   overview.Summary.TotalBalance,
			CashBalance:    overview.Summary.CashBalance,
			BankBalance:    overview.Summary.BankBalance,
			Receivable:     overview.Summary.Receivable,
			Payable:        overview.Summary.Payable,
			ContractAmount: overview.Summary.ContractAmount,
			ReceivedAmount: overview.Summary.ReceivedAmount,
			YearReceived:   overview.Summary.YearReceived,
			UnpaidAmount:   overview.Summary.UnpaidAmount,
			ReceivedRate:   overview.Summary.ReceivedRate,
		},
		Monthly:         monthly,
		IncomeItems:     mapAccountAmounts(overview.IncomeItems),
		ExpenseItems:    mapAccountAmounts(overview.ExpenseItems),
		FundAccounts:    funds,
		AccountBalances: balances,
		RecentFlows:     recentFlows,
	}, nil
}

// Flows 资金流水明细。
func (c *ControllerV1) Flows(ctx context.Context, req *v1.ReportFlowListReq) (*v1.ReportFlowListRes, error) {
	list, err := finReportLogic.New().Flows(ctx, &finReportLogic.FlowFilter{
		Year:          req.Year,
		Month:         req.Month,
		DirectionCode: req.DirectionCode,
		FundAccountId: req.FundAccountId,
		Keyword:       req.Keyword,
		Current:       req.Current,
		Size:          req.Size,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*v1.ReportFlowItem, 0, len(list.Records))
	for _, item := range list.Records {
		records = append(records, &v1.ReportFlowItem{
			Id:              item.Id,
			OccurredAt:      item.OccurredAt,
			DirectionCode:   item.DirectionCode,
			Amount:          item.Amount,
			Remark:          item.Remark,
			AccountName:     item.AccountName,
			FundAccountName: item.FundAccountName,
		})
	}

	return &v1.ReportFlowListRes{
		Records: records,
		Total:   list.Total,
		Summary: v1.ReportFlowSummary{
			Income:  list.Summary.Income,
			Expense: list.Summary.Expense,
			Net:     list.Summary.Net,
		},
	}, nil
}

// mapAccountAmounts 转换收入 / 支出构成列表。
func mapAccountAmounts(items []*finReportLogic.AccountAmount) []*v1.ReportAccountAmount {
	records := make([]*v1.ReportAccountAmount, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.ReportAccountAmount{
			Code:    item.Code,
			Name:    item.Name,
			Amount:  item.Amount,
			Percent: item.Percent,
		})
	}
	return records
}

// ProjectProfitList 项目成本 / 毛利列表。
func (c *ControllerV1) ProjectProfitList(ctx context.Context, req *v1.ProjectProfitListReq) (*v1.ProjectProfitListRes, error) {
	list, err := finReportLogic.New().ProjectProfits(ctx, req.Keyword, req.Status, req.OnlyLoss)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.ProjectProfitItem, 0, len(list.Records))
	for _, item := range list.Records {
		records = append(records, mapProjectProfit(item))
	}

	return &v1.ProjectProfitListRes{
		Records: records,
		Summary: &v1.ProjectProfitSummary{
			ProjectCount:      list.Summary.ProjectCount,
			SettlementTotal:   list.Summary.SettlementTotal,
			ReceivedTotal:     list.Summary.ReceivedTotal,
			UnpaidTotal:       list.Summary.UnpaidTotal,
			CostTotal:         list.Summary.CostTotal,
			GrossProfit:       list.Summary.GrossProfit,
			GrossRate:         list.Summary.GrossRate,
			LossCount:         list.Summary.LossCount,
			CompanyCost:       list.Summary.CompanyCost,
			MaterialUsedTotal: list.Summary.MaterialUsedTotal,
		},
	}, nil
}

// ProjectProfitDetail 单个项目的成本归集与毛利明细。
func (c *ControllerV1) ProjectProfitDetail(ctx context.Context, req *v1.ProjectProfitDetailReq) (*v1.ProjectProfitDetailRes, error) {
	detail, err := finReportLogic.New().ProjectProfitDetail(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	res := &v1.ProjectProfitDetailRes{
		ProjectProfitItem: *mapProjectProfit(&detail.ProjectProfit),
		Receipts:          make([]*v1.ProjectProfitReceipt, 0, len(detail.Receipts)),
		Payouts:           make([]*v1.ProjectProfitPayout, 0, len(detail.Payouts)),
		Materials:         make([]*v1.ProjectProfitMaterial, 0, len(detail.Materials)),
		Changes:           make([]*v1.ProjectProfitChange, 0, len(detail.Changes)),
	}

	for _, item := range detail.Receipts {
		res.Receipts = append(res.Receipts, &v1.ProjectProfitReceipt{
			Id:          item.Id,
			PaymentType: item.PaymentType,
			Amount:      item.Amount,
			PaymentDate: item.PaymentDate,
			Remark:      item.Remark,
		})
	}
	for _, item := range detail.Payouts {
		res.Payouts = append(res.Payouts, &v1.ProjectProfitPayout{
			Id:         item.Id,
			PayoutType: item.PayoutType,
			Amount:     item.Amount,
			PayoutDate: item.PayoutDate,
			VendorName: item.VendorName,
			Remark:     item.Remark,
		})
	}
	for _, item := range detail.Materials {
		res.Materials = append(res.Materials, &v1.ProjectProfitMaterial{
			Id:           item.Id,
			MaterialName: item.MaterialName,
			Unit:         item.Unit,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			Amount:       item.Amount,
			UsedAt:       item.UsedAt,
			Receiver:     item.Receiver,
		})
	}
	for _, item := range detail.Changes {
		res.Changes = append(res.Changes, &v1.ProjectProfitChange{
			Id:         item.Id,
			ChangeType: item.ChangeType,
			Title:      item.Title,
			Amount:     item.Amount,
			ChangeDate: item.ChangeDate,
			Reason:     item.Reason,
		})
	}

	return res, nil
}

// mapProjectProfit 转换项目毛利行。
func mapProjectProfit(item *finReportLogic.ProjectProfit) *v1.ProjectProfitItem {
	costItems := make([]*v1.ProjectProfitCostItem, 0, len(item.CostItems))
	for _, cost := range item.CostItems {
		costItems = append(costItems, &v1.ProjectProfitCostItem{
			Type:    cost.Type,
			Amount:  cost.Amount,
			Percent: cost.Percent,
		})
	}

	return &v1.ProjectProfitItem{
		ProjectId:        item.ProjectId,
		ProjectName:      item.ProjectName,
		CustomerName:     item.CustomerName,
		Status:           item.Status,
		StartDate:        item.StartDate,
		EndDate:          item.EndDate,
		ContractAmount:   item.ContractAmount,
		ChangeNet:        item.ChangeNet,
		SettlementAmount: item.SettlementAmount,
		ReceivedAmount:   item.ReceivedAmount,
		UnpaidAmount:     item.UnpaidAmount,
		ReceivedRate:     item.ReceivedRate,
		CostTotal:        item.CostTotal,
		CostItems:        costItems,
		MaterialPaid:     item.MaterialPaid,
		MaterialUsed:     item.MaterialUsed,
		GrossProfit:      item.GrossProfit,
		GrossRate:        item.GrossRate,
		CashProfit:       item.CashProfit,
	}
}
