// Package print 打印设置控制器。
package print

import (
	"context"

	printApi "engineServ/api/print"
	v1 "engineServ/api/print/v1"
	printLogic "engineServ/internal/logic/api/print"
)

// ControllerV1 打印设置控制器。
type ControllerV1 struct {
	setting *printLogic.Service
}

// NewV1 创建打印设置控制器。
func NewV1() printApi.IPrintV1 {
	return &ControllerV1{setting: printLogic.New()}
}

// SettingGet 获取打印设置。
func (c *ControllerV1) SettingGet(ctx context.Context, req *v1.SettingGetReq) (*v1.SettingGetRes, error) {
	setting, err := c.setting.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.SettingGetRes{Setting: setting}, nil
}

// SettingSave 保存打印设置。
func (c *ControllerV1) SettingSave(ctx context.Context, req *v1.SettingSaveReq) (*v1.SettingSaveRes, error) {
	if err := c.setting.Save(ctx, req); err != nil {
		return nil, err
	}

	return &v1.SettingSaveRes{}, nil
}
