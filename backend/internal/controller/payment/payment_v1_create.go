package payment

import (
	"context"

	v1 "engineServ/api/payment/v1"
	"engineServ/internal/consts"
	paymentLogic "engineServ/internal/logic/api/payment"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	id, voucherNo, err := c.payment.Create(ctx, &paymentLogic.Form{
		ProjectId:     req.ProjectId,
		Amount:        req.Amount,
		PaymentType:   req.PaymentType,
		PaymentDate:   req.PaymentDate,
		FundAccountId: req.FundAccountId,
		Remark:        req.Remark,
	}, userId)
	if err != nil {
		return nil, err
	}

	return &v1.CreateRes{Id: id, VoucherNo: voucherNo}, nil
}
