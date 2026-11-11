package fin_voucher

import (
	"context"

	v1 "engineServ/api/fin_voucher/v1"
)

// IFinVoucherV1 记账凭证与月结接口（登录后可操作）。
type IFinVoucherV1 interface {
	List(ctx context.Context, req *v1.VoucherListReq) (res *v1.VoucherListRes, err error)
	Detail(ctx context.Context, req *v1.VoucherDetailReq) (res *v1.VoucherDetailRes, err error)
	Periods(ctx context.Context, req *v1.PeriodListReq) (res *v1.PeriodListRes, err error)
	Close(ctx context.Context, req *v1.ClosePeriodReq) (res *v1.ClosePeriodRes, err error)
	Reopen(ctx context.Context, req *v1.ReopenPeriodReq) (res *v1.ReopenPeriodRes, err error)
	Export(ctx context.Context, req *v1.VoucherExportReq) (res *v1.VoucherExportRes, err error)
}
