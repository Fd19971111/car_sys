package auth

import (
	"context"

	v1 "engineServ/api/auth/v1"
)

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	token, err := c.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &v1.RefreshTokenRes{
		Token:        token,
		RefreshToken: token,
	}, nil
}
