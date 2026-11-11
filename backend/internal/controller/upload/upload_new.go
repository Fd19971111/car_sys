package upload

import (
	uploadApi "engineServ/api/upload"
)

// ControllerV1 小程序 / 业主端图片上传。
type ControllerV1 struct{}

// AdminControllerV1 后台图片上传（单独一个类型，避免接口绑定重复注册路由）。
type AdminControllerV1 struct{}

func NewV1() uploadApi.IUploadV1 {
	return &ControllerV1{}
}

func NewAdminV1() uploadApi.IAdminUploadV1 {
	return &AdminControllerV1{}
}
