package payout

import (
	"context"

	v1 "engineServ/api/payout/v1"
	"engineServ/internal/consts"
	payoutLogic "engineServ/internal/logic/api/payout"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	id, voucherNo, err := c.payout.Create(ctx, &payoutLogic.Form{
		ProjectId:       req.ProjectId,
		Amount:          req.Amount,
		PayoutType:      req.PayoutType,
		PayoutDate:      req.PayoutDate,
		VendorId:        req.VendorId,
		PurchaseOrderId: req.PurchaseOrderId,
		FundAccountId:   req.FundAccountId,
		Remark:          req.Remark,
	}, userId)
	if err != nil {
		return nil, err
	}

	return &v1.CreateRes{Id: id, VoucherNo: voucherNo}, nil
}
