package middleware

import (
	"engineServ/internal/consts"
	customerAuthLogic "engineServ/internal/logic/api/customer_auth"

	"github.com/gogf/gf/v2/net/ghttp"
)

// CustomerAuth 业主端（小程序）鉴权中间件。
//
// 与后台的 Auth 中间件完全独立：
//   - 校验的是 customer_token（业主手机号登录签发）
//   - 通过后把 customerId 注入上下文，/mp 下所有接口只认这一个来源
//
// 业务错误码（经 MiddlewareHandlerResponse 统一包装）：
//   - 8888 令牌缺失/无效/业主不存在：前端清除登录态并跳登录页
//   - 9999 令牌已过期：前端同样跳登录页（业主端暂不做静默续期）
func CustomerAuth(r *ghttp.Request) {
	customerId, err := customerAuthLogic.New().Authenticate(r.Context())
	if err != nil {
		r.SetError(err)
		r.Exit()
		return
	}

	r.SetCtxVar(consts.CtxCustomerIdKey, customerId)

	r.Middleware.Next()
}
