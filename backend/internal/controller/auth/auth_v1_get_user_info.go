package auth

import (
	"context"

	v1 "engineServ/api/auth/v1"
	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1Plus) GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (res *v1.GetUserInfoRes, err error) {
	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	user, roles, buttons, err := c.auth.GetUserInfo(ctx, userId)
	if err != nil {
		return nil, err
	}
	initialized, err := g.DB().GetValue(ctx, "SELECT is_initialized FROM fin_setup_config WHERE company_id=?", user.CompanyId)
	if err != nil {
		return nil, err
	}

	return &v1.GetUserInfoRes{
		UserId:         user.Id,
		UserName:       user.Username,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Phone:          user.Phone,
		Role:           user.Role,
		MerchantId:     user.MerchantId,
		CompanyId:      user.CompanyId,
		FinInitialized: initialized.Bool(),
		Roles:          roles,
		Buttons:        buttons,
	}, nil
}
