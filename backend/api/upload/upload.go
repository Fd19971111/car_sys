package upload

import (
	"context"

	v1 "engineServ/api/upload/v1"
)

// IUploadV1 文件上传模块接口（小程序 / 业主端）
type IUploadV1 interface {
	Image(ctx context.Context, req *v1.ImageReq) (res *v1.ImageRes, err error)
}

// IAdminUploadV1 文件上传模块接口（后台 / site_user 端）
type IAdminUploadV1 interface {
	AdminImage(ctx context.Context, req *v1.AdminImageReq) (res *v1.AdminImageRes, err error)
	AdminContractFile(ctx context.Context, req *v1.AdminContractFileReq) (res *v1.AdminContractFileRes, err error)
}
