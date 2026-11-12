package cmd

import (
	"context"

	"engineServ/internal/controller/auth"
	"engineServ/internal/controller/customer"
	"engineServ/internal/controller/department"
	driverFinanceController "engineServ/internal/controller/driver_finance"
	finBasicController "engineServ/internal/controller/fin_basic"
	finInitController "engineServ/internal/controller/fin_init"
	finReportController "engineServ/internal/controller/fin_report"
	finVoucherController "engineServ/internal/controller/fin_voucher"
	menuController "engineServ/internal/controller/menu"
	oplogController "engineServ/internal/controller/oplog"
	orderController "engineServ/internal/controller/order"
	"engineServ/internal/controller/payment"
	"engineServ/internal/controller/payout"
	pricingRuleController "engineServ/internal/controller/pricing_rule"
	printController "engineServ/internal/controller/print"
	"engineServ/internal/controller/staff"
	transportController "engineServ/internal/controller/transport"
	uploadController "engineServ/internal/controller/upload"
	"engineServ/internal/controller/user"
	"engineServ/internal/logic/middleware"
	"engineServ/internal/logic/upload"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "start http server",
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		s := g.Server()

		// 本地静态资源：附件（图片 / 合同）落在 app.uploadDir（默认 storage/file），
		// 通过 app.uploadUrlPath（默认 /uploads/**）访问。
		// 正式环境建议由 Nginx 托管，或切到 OSS / CDN（改 app.uploadUrlPrefix 即可）。
		s.AddStaticPath(upload.URLPathPrefix(ctx), upload.StorageDir(ctx))

		s.Group("/", func(group *ghttp.RouterGroup) {
			// 同 8006 项目的同一套中间件/路由分层：
			// RequestLogging 在最外层包裹所有中间件，然后 CORS，然后统一响应包装
			group.Middleware(middleware.RequestLogging)
			group.Middleware(ghttp.MiddlewareCORS)
			group.Middleware(ghttp.MiddlewareHandlerResponse)

			// /apiv1 分组：按「客户端」分成三条独立的鉴权链路
			group.Group("/apiv1", func(group *ghttp.RouterGroup) {
				// ---------- 1. 后台：登录 / 刷新令牌免鉴权 ----------
				group.Bind(auth.NewV1())

				// ---------- 2. 后台：需要 site_user_token ----------
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					group.Bind(auth.NewV1Plus())
					group.Bind(customer.NewV1())
					group.Bind(department.NewV1())
					group.Bind(driverFinanceController.NewV1())
					group.Bind(finInitController.NewV1())
					group.Bind(finBasicController.NewV1())
					group.Bind(finReportController.NewV1())
					group.Bind(finVoucherController.NewV1())
					group.Bind(menuController.NewV1())
					group.Bind(orderController.NewV1())
					group.Bind(oplogController.NewV1())
					group.Bind(transportController.NewV1())
					group.Bind(staff.NewV1())
					group.Bind(user.NewV1())
					group.Bind(payment.NewV1())
					group.Bind(payout.NewV1())
					group.Bind(printController.NewV1())
					group.Bind(pricingRuleController.NewV1())
					group.Bind(uploadController.NewAdminV1())
				})
			})
		})

		s.Run()
		return nil
	},
}
