package middleware

import (
	"strings"

	"engineServ/internal/consts"
	authLogic "engineServ/internal/logic/api/auth"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth 登录鉴权中间件：以请求头 Bearer 令牌查询 site_user_token，
// 校验有效期及 site_user 用户状态，通过后把用户信息注入请求上下文。
//
// 业务错误码（经 MiddlewareHandlerResponse 统一包装为 HTTP 200 响应体）：
//   - 8888 令牌缺失/无效/用户被禁用：前端清除登录态并跳转登录页
//   - 9999 令牌已过期：前端使用 refreshToken 静默续期后重试
func Auth(r *ghttp.Request) {
	user, err := authLogic.New().Authenticate(r.Context())
	if err != nil {
		r.SetError(err)
		r.Exit()
		return
	}

	r.SetCtxVar(consts.CtxUserIdKey, user.Id)
	r.SetCtxVar(consts.CtxUsernameKey, user.Username)
	r.SetCtxVar(consts.CtxRoleKey, user.Role)
	r.SetCtxVar(consts.CtxMerchantIdKey, user.MerchantId)
	r.SetCtxVar(consts.CtxCompanyIdKey, user.CompanyId)
	if err := checkModulePermission(r, user.Id, user.Role, user.CompanyId); err != nil {
		r.SetError(err)
		r.Exit()
		return
	}

	r.Middleware.Next()
}

func checkModulePermission(r *ghttp.Request, userId, role, companyId uint) error {
	if role == 1 {
		return nil
	}

	// 前缀 → 允许的菜单路由名。一个前缀可能对应多个菜单：
	// 例如「全部订单」和「行程监控」共用 /apiv1/order/ 前缀，授予任一即可访问。
	moduleRules := []struct {
		prefix string
		routes []string
	}{
		{"/apiv1/customer/", []string{"settings_customer"}},
		{"/apiv1/order/", []string{"order_all", "order_monitor"}},
		{"/apiv1/payment/", []string{"finance_payment"}},
		{"/apiv1/payout/", []string{"finance_payout"}},
		{"/apiv1/fin-voucher/", []string{"finance_voucher"}},
		{"/apiv1/menu/", []string{"settings_menu"}},
		{"/apiv1/pricing-rule/", []string{"settings_pricing"}},
		{"/apiv1/transport/", []string{"transport_driver", "transport_vehicle"}},
	}

	for _, rule := range moduleRules {
		if !strings.HasPrefix(r.URL.Path, rule.prefix) {
			continue
		}

		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(rule.routes)), ",")
		args := make([]interface{}, 0, len(rule.routes)+2)
		args = append(args, userId, companyId)
		for _, route := range rule.routes {
			args = append(args, route)
		}

		count, err := g.DB().GetValue(r.Context(), `SELECT COUNT(*) FROM site_user u
			JOIN sys_department_menu dm ON dm.department_id=u.department_id
			JOIN sys_menu granted ON granted.id=dm.menu_id
			JOIN sys_menu target ON (target.id=granted.id OR target.parent_id=granted.id)
			WHERE u.id=? AND dm.company_id=? AND target.enabled=1 AND target.route_name IN (`+placeholders+`)`, args...)
		if err != nil {
			return gerror.Wrap(err, "校验部门权限失败")
		}
		if count.Int() == 0 {
			return gerror.New("当前部门无权访问该功能")
		}
		return nil
	}

	return nil
}
