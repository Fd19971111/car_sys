package mp

import (
	"context"

	mpApi "engineServ/api/mp"
	v1 "engineServ/api/mp/v1"
	customerAuthLogic "engineServ/internal/logic/api/customer_auth"
)

// AuthControllerV1 免登录的小程序接口（仅登录）
type AuthControllerV1 struct {
	auth *customerAuthLogic.CustomerAuth
}

func NewAuthV1() mpApi.IMpAuthV1 {
	return &AuthControllerV1{auth: customerAuthLogic.New()}
}

func (c *AuthControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	token, expireAt, customer, err := c.auth.Login(ctx, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}

	return &v1.LoginRes{
		Token:    token,
		ExpireAt: formatDateTime(expireAt),
		Owner:    toOwnerProfile(customer),
	}, nil
}
