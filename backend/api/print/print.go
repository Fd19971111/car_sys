// Package print 打印设置模块接口定义。
//
// 本文件由 GoFrame 开发工具自动生成，也可手动维护。
package print

import (
	"context"

	v1 "engineServ/api/print/v1"
)

type IPrintV1 interface {
	SettingGet(ctx context.Context, req *v1.SettingGetReq) (res *v1.SettingGetRes, err error)
	SettingSave(ctx context.Context, req *v1.SettingSaveReq) (res *v1.SettingSaveRes, err error)
}
