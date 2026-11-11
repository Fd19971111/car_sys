package auth

import (
	"context"

	v1 "engineServ/api/auth/v1"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	userId, token, err := c.auth.Login(ctx, req.UserName, req.Password, req.Client)
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{
		UserId:       userId,
		Token:        token,
		RefreshToken: token,
	}, nil
}
