package middleware

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	_accessLog     *glog.Logger
	_accessLogOnce sync.Once
)

// AccessLog 返回访问日志记录器，所有 HTTP 请求/响应日志写入 request_{Ymd}.log
func AccessLog() *glog.Logger {
	_accessLogOnce.Do(func() {
		_accessLog = glog.New()
		_accessLog.SetConfig(g.Log().GetConfig())
		_accessLog.SetFile("request_{Ymd}.log")
	})
	return _accessLog
}

// RequestLogging 记录所有 HTTP 请求与响应的中间件。
//
// 记录的字段：
//   - method, url, path, clientIp         请求基本信息
//   - query                               请求 Query 参数
//   - body                                请求体 (POST/PUT/PATCH)
//   - statusCode                          响应 HTTP 状态码
//   - response                            最终返回的业务数据 (handler 返回值)
//   - error                               业务异常信息
//   - cost_ms                             请求耗时(毫秒)
//
// 中间件顺序：该中间件应放在 MiddlewareHandlerResponse 之后注册，
// 以最外层身份包裹所有下游中间件，确保即使是 Auth 等中间件的提前退出也能被记录。
func RequestLogging(r *ghttp.Request) {
	start := time.Now()

	// ---- 1. 记录请求参数 ----
	reqInfo := g.Map{
		"method":   r.Method,
		"url":      r.URL.String(),
		"path":     r.URL.Path,
		"clientIp": r.GetClientIp(),
	}

	// Query 参数
	if queryMap := r.GetQueryMap(); len(queryMap) > 0 {
		reqInfo["query"] = queryMap
	}

	// Body 参数 (仅 POST/PUT/PATCH 等有体的方法)
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		if bodyMap := r.GetMap(); len(bodyMap) > 0 {
			reqInfo["body"] = bodyMap
		}
	}

	// ---- 2. 执行后续中间件链 ----
	r.Middleware.Next()

	// ---- 3. 收集响应信息 ----
	cost := time.Since(start).Milliseconds()

	resInfo := g.Map{
		"request":    reqInfo,
		"statusCode": r.Response.Status,
		"cost_ms":    cost,
	}

	// 重建 MiddlewareHandlerResponse 写出的完整响应（与客户端看到的一致）
	var (
		handlerErr  = r.GetError()
		handlerRes  = r.GetHandlerResponse()
		handlerCode = gerror.Code(handlerErr)
		handlerMsg  string
	)
	if handlerErr != nil {
		if handlerCode == gcode.CodeNil {
			handlerCode = gcode.CodeInternalError
		}
		handlerMsg = handlerErr.Error()
	} else {
		handlerCode = gcode.CodeOK
		handlerMsg = handlerCode.Message()
	}
	resInfo["response"] = g.Map{
		"code":    handlerCode.Code(),
		"message": handlerMsg,
		"data":    handlerRes,
	}

	// ---- 4. 输出日志 ----
	// 异常请求用 Warning 级别, 正常请求用 Info 级别
	logFunc := AccessLog().Info
	if r.GetError() != nil || r.Response.Status >= 400 {
		logFunc = AccessLog().Warning
	}
	logFunc(r.Context(), gconv.String(resInfo))
}
