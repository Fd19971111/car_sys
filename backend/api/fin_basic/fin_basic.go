package fin_basic

import (
	"context"

	v1 "engineServ/api/fin_basic/v1"
)

// IFinBasicV1 财务基础资料接口（登录后可读，写操作仅超级管理员）。
type IFinBasicV1 interface {
	Accounts(ctx context.Context, req *v1.AccountListReq) (res *v1.AccountListRes, err error)
	SaveAccount(ctx context.Context, req *v1.AccountSaveReq) (res *v1.AccountSaveRes, err error)
	DeleteAccount(ctx context.Context, req *v1.AccountDeleteReq) (res *v1.AccountDeleteRes, err error)
	ToggleAccount(ctx context.Context, req *v1.AccountToggleReq) (res *v1.AccountToggleRes, err error)

	FundAccounts(ctx context.Context, req *v1.FundAccountListReq) (res *v1.FundAccountListRes, err error)
	SaveFundAccount(ctx context.Context, req *v1.FundAccountSaveReq) (res *v1.FundAccountSaveRes, err error)
	DeleteFundAccount(ctx context.Context, req *v1.FundAccountDeleteReq) (res *v1.FundAccountDeleteRes, err error)
	ToggleFundAccount(ctx context.Context, req *v1.FundAccountToggleReq) (res *v1.FundAccountToggleRes, err error)
}
