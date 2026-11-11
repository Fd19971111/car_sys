package upload

import (
	"context"

	v1 "engineServ/api/upload/v1"
	uploadLogic "engineServ/internal/logic/upload"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminImage 后台图片上传（site_user 令牌）。
//
// 与小程序端 Image 的落库方式完全一致，只是鉴权链路不同：
// 后台登记施工日志、变更附件时需要这个入口。
func (c *AdminControllerV1) AdminImage(ctx context.Context, req *v1.AdminImageReq) (res *v1.AdminImageRes, err error) {
	request := g.RequestFromCtx(ctx)

	files := request.GetUploadFiles("file")
	if len(files) == 0 {
		files = request.GetUploadFiles("files")
	}

	if len(files) == 0 {
		return nil, gerror.New("请选择要上传的图片")
	}

	if len(files) > uploadLogic.MaxFilesPerRequest {
		return nil, gerror.Newf("单次最多上传 %d 张图片", uploadLogic.MaxFilesPerRequest)
	}

	res = &v1.AdminImageRes{Files: make([]*v1.UploadedFile, 0, len(files))}

	for _, file := range files {
		saved, saveErr := uploadLogic.SaveImage(ctx, file)
		if saveErr != nil {
			return nil, saveErr
		}

		res.Files = append(res.Files, &v1.UploadedFile{
			Path: saved.Path,
			Url:  saved.Url,
			Name: saved.Name,
			Size: saved.Size,
		})
	}

	return res, nil
}

func (c *AdminControllerV1) AdminContractFile(ctx context.Context, _ *v1.AdminContractFileReq) (res *v1.AdminContractFileRes, err error) {
	request := g.RequestFromCtx(ctx)
	files := request.GetUploadFiles("file")
	if len(files) == 0 {
		files = request.GetUploadFiles("files")
	}
	if len(files) == 0 {
		return nil, gerror.New("请选择要上传的合同附件")
	}
	if len(files) > uploadLogic.MaxContractFilesPerRequest {
		return nil, gerror.Newf("单次最多上传 %d 个文件", uploadLogic.MaxContractFilesPerRequest)
	}

	res = &v1.AdminContractFileRes{Files: make([]*v1.UploadedFile, 0, len(files))}
	for _, file := range files {
		saved, saveErr := uploadLogic.SaveContractFile(ctx, file)
		if saveErr != nil {
			return nil, saveErr
		}
		res.Files = append(res.Files, &v1.UploadedFile{Path: saved.Path, Url: saved.Url, Name: saved.Name, Size: saved.Size})
	}
	return res, nil
}
