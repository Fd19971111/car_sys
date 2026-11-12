package payment

import (
	"context"

	v1 "engineServ/api/payment/v1"
	"engineServ/internal/consts"
	paymentLogic "engineServ/internal/logic/api/payment"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	voucherNo, err := c.payment.Update(ctx, req.Id, &paymentLogic.Form{
		RideOrderId:   req.RideOrderId,
		Amount:        req.Amount,
		PaymentType:   req.PaymentType,
		PaymentDate:   req.PaymentDate,
		FundAccountId: req.FundAccountId,
		Remark:        req.Remark,
	}, userId)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateRes{Id: req.Id, VoucherNo: voucherNo}, nil
}
