package consts

import "time"

const (
	// CtxUserIdKey 登录用户 site_user.id 在请求上下文中的键，由 Auth 中间件注入。
	CtxUserIdKey = "ctx:userId"

	// CtxUsernameKey 登录用户名在请求上下文中的键，由 Auth 中间件注入。
	CtxUsernameKey = "ctx:username"

	// CtxRoleKey 登录用户角色（1=管理员 2=员工 3=商家）在请求上下文中的键。
	CtxRoleKey = "ctx:role"

	// CtxMerchantIdKey 登录用户所属商家ID在请求上下文中的键。
	CtxMerchantIdKey = "ctx:merchantId"

	// CtxCompanyIdKey 当前租户公司ID，由 Auth 中间件从已认证用户记录注入。
	CtxCompanyIdKey = "ctx:companyId"

	// CtxCustomerIdKey 登录业主 customer.id 在请求上下文中的键，由 CustomerAuth 中间件注入。
	//
	// 小程序端的数据隔离全靠它：所有 /mp 接口都只认这个值，
	// 绝不允许再从请求参数里接收 customerId（否则等于没隔离）。
	CtxCustomerIdKey = "ctx:customerId"

	// DefaultTokenExpire 未配置 app.tokenExpire 时的令牌默认有效期（7天）。
	DefaultTokenExpire = 7 * 24 * time.Hour

	// DefaultMpTokenExpire 未配置 app.mpTokenExpire 时业主令牌的默认有效期（30天）。
	//
	// 业主是小程序用户，不应该频繁登录，所以比后台令牌长；
	// 并且每次请求会做滑动续期，持续使用的业主不会被登出。
	DefaultMpTokenExpire = 30 * 24 * time.Hour

	// CodeTokenExpired 令牌已过期：前端收到后使用 refreshToken 静默续期。
	// 与前端 .env VITE_SERVICE_EXPIRED_TOKEN_CODES 保持一致。
	CodeTokenExpired = 9999

	// CodeAuthFailed 登录态缺失/无效/被禁用：前端收到后清除登录态并跳转登录页。
	// 与前端 .env VITE_SERVICE_LOGOUT_CODES 保持一致。
	CodeAuthFailed = 8888
)
