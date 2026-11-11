package v1

import "github.com/gogf/gf/v2/frame/g"

// UploadedFile 已保存的文件
type UploadedFile struct {
	Path string `json:"path" dc:"相对路径（写库用），如 /uploads/20260930/xxx.jpg"`
	Url  string `json:"url"  dc:"绝对地址（前端预览用）"`
	Name string `json:"name" dc:"原始文件名"`
	Size int64  `json:"size" dc:"文件大小（字节）"`
}

// ImageReq 图片上传请求（小程序端，业主令牌）
//
// multipart/form-data，文件字段名统一用 file（同名字段可传多个，单次最多 9 个）。
type ImageReq struct {
	g.Meta `path:"/upload/image" method:"post" mime:"multipart/form-data" sm:"上传图片" tags:"Upload"`

	// 说明：文件本身通过 r.GetUploadFiles("file") 读取，
	// 这里不声明字段是为了同时兼容单文件与多文件上传。
}

// ImageRes 图片上传响应
type ImageRes struct {
	Files []*UploadedFile `json:"files" dc:"已上传的文件列表"`
}

// AdminImageReq 图片上传请求（后台端，site_user 令牌）
//
// 与小程序端的区别只是鉴权链路不同；后台登记施工日志 / 变更附件都用这个接口。
type AdminImageReq struct {
	g.Meta `path:"/upload/adminImage" method:"post" mime:"multipart/form-data" sm:"后台上传图片" tags:"Upload"`
}

// AdminImageRes 后台上传图片响应
type AdminImageRes struct {
	Files []*UploadedFile `json:"files" dc:"已上传的文件列表"`
}

type AdminContractFileReq struct {
	g.Meta `path:"/upload/adminContractFile" method:"post" mime:"multipart/form-data" sm:"上传合同附件" tags:"Upload"`
}

type AdminContractFileRes struct {
	Files []*UploadedFile `json:"files" dc:"已上传的合同附件"`
}
