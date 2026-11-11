package user

import (
	"context"

	v1 "engineServ/api/user/v1"
)

type IUserV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	ResetPwd(ctx context.Context, req *v1.ResetPwdReq) (res *v1.ResetPwdRes, err error)
	ChangePwd(ctx context.Context, req *v1.ChangePwdReq) (res *v1.ChangePwdRes, err error)
	Toggle(ctx context.Context, req *v1.ToggleReq) (res *v1.ToggleRes, err error)
	CreateAccount(ctx context.Context, req *v1.CreateAccountReq) (res *v1.CreateAccountRes, err error)
	AvailableStaff(ctx context.Context, req *v1.AvailableStaffReq) (res *v1.AvailableStaffRes, err error)
	UpdateAccount(ctx context.Context, req *v1.UpdateAccountReq) (res *v1.UpdateAccountRes, err error)
}
