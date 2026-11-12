package payment

import (
	"context"

	v1 "engineServ/api/payment/v1"
	"engineServ/internal/logic/fin/booking"
)

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	detail, err := c.payment.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res = &v1.DetailRes{
		Id:              detail.Payment.Id,
		Amount:          detail.Payment.Amount,
		PaymentType:     detail.Payment.PaymentType,
		PaymentDate:     formatDate(detail.Payment.PaymentDate),
		FundAccountId:   detail.Payment.FundAccountId,
		FundAccountName: detail.FundAccountName,
		VoucherId:       detail.Payment.VoucherId,
		VoucherNo:       detail.VoucherNo,
		Remark:          detail.Payment.Remark,
		CreatedAt:       formatDateTime(detail.Payment.CreatedAt),
		LedgerEntries:   mapLedgerEntries(detail.Entries),
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
