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
			Income:        overview.Summary.Income,
			Expense:       overview.Summary.Expense,
			Net:           overview.Summary.Net,
			NetRate:       overview.Summary.NetRate,
			TotalBalance:  overview.Summary.TotalBalance,
			CashBalance:   overview.Summary.CashBalance,
			BankBalance:   overview.Summary.BankBalance,
			Receivable:    overview.Summary.Receivable,
			Payable:       overview.Summary.Payable,
			RideRevenue:   overview.Summary.RideRevenue,
			PassengerPaid: overview.Summary.PassengerPaid,
			YearReceived:  overview.Summary.YearReceived,
			UnpaidFare:    overview.Summary.UnpaidFare,
			DriverPayable: overview.Summary.DriverPayable,
			DriverSettled: overview.Summary.DriverSettled,
			DriverExpense: overview.Summary.DriverExpense,
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
