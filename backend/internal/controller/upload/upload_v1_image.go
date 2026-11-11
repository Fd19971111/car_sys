package upload

import (
	"context"

	v1 "engineServ/api/upload/v1"
	uploadLogic "engineServ/internal/logic/upload"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Image(ctx context.Context, req *v1.ImageReq) (res *v1.ImageRes, err error) {
	request := g.RequestFromCtx(ctx)

	// 主字段名 file；兼容前端传 files / file[] 的情况
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

	res = &v1.ImageRes{Files: make([]*v1.UploadedFile, 0, len(files))}

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
