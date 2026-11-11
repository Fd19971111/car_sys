package mp

import (
	"context"

	v1 "engineServ/api/mp/v1"
)

// IMpV1 小程序（业主端）需登录的接口。
//
// 这些接口全部挂在 CustomerAuth 中间件下：业主身份从令牌解析并写入上下文，
// 业务层不再从请求参数接收 customerId。
type IMpV1 interface {
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
	Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error)
	Project(ctx context.Context, req *v1.ProjectReq) (res *v1.ProjectRes, err error)
	ProgressCreate(ctx context.Context, req *v1.ProgressCreateReq) (res *v1.ProgressCreateRes, err error)
}

// IMpAuthV1 小程序（业主端）免登录接口。
//
// 单独拆一份接口，是因为 GoFrame 的 group.Bind 会绑定对象的全部方法，
// 登录接口必须放在没有 CustomerAuth 中间件的分组里。
type IMpAuthV1 interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
}
