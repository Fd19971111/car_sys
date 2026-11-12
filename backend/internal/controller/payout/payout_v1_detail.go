package payout

import (
	"context"

	v1 "engineServ/api/payout/v1"
	"engineServ/internal/logic/fin/booking"
)

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	detail, err := c.payout.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res = &v1.DetailRes{
		Id:              detail.Payout.Id,
		ProjectId:       detail.Payout.ProjectId,
		Amount:          detail.Payout.Amount,
		PayoutType:      detail.Payout.PayoutType,
		PayoutDate:      formatDate(detail.Payout.PayoutDate),
		VendorName:      detail.VendorName,
		FundAccountId:   detail.Payout.FundAccountId,
		FundAccountName: detail.FundAccountName,
		VoucherId:       detail.Payout.VoucherId,
		VoucherNo:       detail.VoucherNo,
		Remark:          detail.Payout.Remark,
		CreatedAt:       formatDateTime(detail.Payout.CreatedAt),
		ProjectPaid:     detail.PaidAmount,
		LedgerEntries:   mapLedgerEntries(detail.Entries),
	}

	if detail.Project != nil {
		res.ProjectName = detail.Project.Name
		res.ProjectAmount = detail.Project.Amount
	}

	if detail.Customer != nil {
		res.CustomerName = detail.Customer.Name
	}

	return res, nil
}

// mapLedgerEntries 转换凭证分录展示结构。
func mapLedgerEntries(entries []*booking.EntryView) []*v1.LedgerEntryItem {
	items := make([]*v1.LedgerEntryItem, 0, len(entries))
	for _, entry := range entries {
		items = append(items, &v1.LedgerEntryItem{
			Direction:   entry.Direction,
			AccountName: entry.AccountName,
			Amount:      entry.Amount,
			Summary:     entry.Summary,
		})
	}

	return items
}
