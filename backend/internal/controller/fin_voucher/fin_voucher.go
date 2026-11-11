package fin_voucher

import (
	"context"

	finVoucherApi "engineServ/api/fin_voucher"
	v1 "engineServ/api/fin_voucher/v1"
	finVoucherLogic "engineServ/internal/logic/api/fin_voucher"
)

// ControllerV1 记账凭证控制器。
type ControllerV1 struct{}

// NewV1 创建记账凭证控制器。
func NewV1() finVoucherApi.IFinVoucherV1 { return &ControllerV1{} }

// List 凭证分页列表。
func (c *ControllerV1) List(ctx context.Context, req *v1.VoucherListReq) (*v1.VoucherListRes, error) {
	rows, total, err := finVoucherLogic.New().List(ctx, &finVoucherLogic.ListFilter{
		PeriodCode: req.PeriodCode,
		StatusCode: req.StatusCode,
		Keyword:    req.Keyword,
		Current:    req.Current,
		Size:       req.Size,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*v1.VoucherListItem, 0, len(rows))
	for _, row := range rows {
		records = append(records, &v1.VoucherListItem{
			Id:              row.Id,
			VoucherNo:       row.VoucherNo,
			VoucherDate:     row.VoucherDate,
			PeriodCode:      row.PeriodCode,
			VoucherTypeCode: row.VoucherTypeCode,
			Summary:         row.Summary,
			StatusCode:      row.StatusCode,
			SourceType:      row.SourceType,
			SourceId:        row.SourceId,
			DebitTotal:      row.DebitTotal,
			CreditTotal:     row.CreditTotal,
			CreatedAt:       row.CreatedAt,
		})
	}

	return &v1.VoucherListRes{Current: req.Current, Size: req.Size, Total: total, Records: records}, nil
}

// Detail 凭证详情。
func (c *ControllerV1) Detail(ctx context.Context, req *v1.VoucherDetailReq) (*v1.VoucherDetailRes, error) {
	detail, err := finVoucherLogic.New().Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	entries := make([]*v1.VoucherEntryItem, 0, len(detail.Entries))
	for _, entry := range detail.Entries {
		entries = append(entries, &v1.VoucherEntryItem{
			LineNo:        entry.LineNo,
			AccountCode:   entry.AccountCode,
			AccountName:   entry.AccountName,
			Summary:       entry.Summary,
			DirectionCode: entry.DirectionCode,
			Amount:        entry.Amount,
		})
	}

	return &v1.VoucherDetailRes{
		Id:              detail.Id,
		VoucherNo:       detail.VoucherNo,
		VoucherDate:     detail.VoucherDate,
		PeriodCode:      detail.PeriodCode,
		VoucherTypeCode: detail.VoucherTypeCode,
		Summary:         detail.Summary,
		StatusCode:      detail.StatusCode,
		SourceType:      detail.SourceType,
		SourceId:        detail.SourceId,
		CreatedAt:       detail.CreatedAt,
		DebitTotal:      detail.DebitTotal,
		CreditTotal:     detail.CreditTotal,
		Entries:         entries,
	}, nil
}

// Periods 会计期间列表。
func (c *ControllerV1) Periods(ctx context.Context, req *v1.PeriodListReq) (*v1.PeriodListRes, error) {
	list, err := finVoucherLogic.New().Periods(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.PeriodItem, 0, len(list.Records))
	for _, item := range list.Records {
		records = append(records, &v1.PeriodItem{
			PeriodCode:   item.PeriodCode,
			PeriodYear:   item.PeriodYear,
			PeriodMonth:  item.PeriodMonth,
			StatusCode:   item.StatusCode,
			VoucherCount: item.VoucherCount,
			DebitTotal:   item.DebitTotal,
			CreditTotal:  item.CreditTotal,
			CanClose:     item.CanClose,
			CanReopen:    item.CanReopen,
			BeforeStart:  item.BeforeStart,
		})
	}

	return &v1.PeriodListRes{
		Initialized:   list.Initialized,
		CurrentPeriod: list.CurrentPeriod,
		StartPeriod:   list.StartPeriod,
		Records:       records,
	}, nil
}

// Close 会计期间结账。
func (c *ControllerV1) Close(ctx context.Context, req *v1.ClosePeriodReq) (*v1.ClosePeriodRes, error) {
	if err := finVoucherLogic.New().Close(ctx, req.PeriodCode); err != nil {
		return nil, err
	}

	return &v1.ClosePeriodRes{PeriodCode: req.PeriodCode, StatusCode: "CLOSED"}, nil
}

// Reopen 会计期间反结账。
func (c *ControllerV1) Reopen(ctx context.Context, req *v1.ReopenPeriodReq) (*v1.ReopenPeriodRes, error) {
	if err := finVoucherLogic.New().Reopen(ctx, req.PeriodCode); err != nil {
		return nil, err
	}

	return &v1.ReopenPeriodRes{PeriodCode: req.PeriodCode, StatusCode: "OPEN"}, nil
}

// Export 导出凭证分录明细。
func (c *ControllerV1) Export(ctx context.Context, req *v1.VoucherExportReq) (*v1.VoucherExportRes, error) {
	rows, err := finVoucherLogic.New().Export(ctx, &finVoucherLogic.ExportFilter{
		PeriodCode: req.PeriodCode,
		StatusCode: req.StatusCode,
		Keyword:    req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*v1.VoucherExportRow, 0, len(rows))
	for _, row := range rows {
		records = append(records, &v1.VoucherExportRow{
			VoucherNo:       row.VoucherNo,
			VoucherDate:     row.VoucherDate,
			PeriodCode:      row.PeriodCode,
			VoucherTypeCode: row.VoucherTypeCode,
			VoucherSummary:  row.VoucherSummary,
			StatusCode:      row.StatusCode,
			SourceType:      row.SourceType,
			SourceId:        row.SourceId,
			LineNo:          row.LineNo,
			AccountCode:     row.AccountCode,
			AccountName:     row.AccountName,
			EntrySummary:    row.EntrySummary,
			DirectionCode:   row.DirectionCode,
			Amount:          row.Amount,
		})
	}

	return &v1.VoucherExportRes{Records: records}, nil
}
