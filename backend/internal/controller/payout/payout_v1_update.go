package payout

import (
	"context"

	v1 "engineServ/api/payout/v1"
	"engineServ/internal/consts"
	payoutLogic "engineServ/internal/logic/api/payout"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	voucherNo, err := c.payout.Update(ctx, req.Id, &payoutLogic.Form{
		ProjectId:     req.ProjectId,
		Amount:        req.Amount,
		PayoutType:    req.PayoutType,
		PayoutDate:    req.PayoutDate,
		VendorId:      req.VendorId,
		FundAccountId: req.FundAccountId,
		Remark:        req.Remark,
	}, userId)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateRes{Id: req.Id, VoucherNo: voucherNo}, nil
}
