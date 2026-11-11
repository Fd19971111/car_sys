// Package upload 本地文件上传：图片与合同附件统一落到 storage/file 目录，
// 通过静态资源映射对外提供访问；后续切 OSS 时只需替换存储实现与 access key 配置。
package upload

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	// DefaultUploadDir 附件存放根目录（相对 backend 运行目录），生产改 config.yaml 的 app.uploadDir
	DefaultUploadDir = "./storage/file"
	// DefaultURLPath 附件对外访问路径前缀，生产改 config.yaml 的 app.uploadUrlPath
	DefaultURLPath = "/uploads"
	// MaxFileSize 单文件大小上限 10MB
	MaxFileSize = 10 * 1024 * 1024
	// MaxFilesPerRequest 单次请求最多上传的文件数
	MaxFilesPerRequest         = 9
	MaxContractFilesPerRequest = 9
	MaxContractFileSize        = 20 * 1024 * 1024
)

// allowedExts 允许上传的图片格式（存库前先按扩展名校验）
var allowedExts = map[string]struct{}{
	".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".webp": {}, ".bmp": {},
}

var allowedContractExts = map[string]struct{}{
	".pdf": {}, ".doc": {}, ".docx": {}, ".xls": {}, ".xlsx": {},
	".jpg": {}, ".jpeg": {}, ".png": {},
}

// SavedFile 上传结果
type SavedFile struct {
	// Path 相对路径（写库存这个）：/uploads/20260930/xxxx.jpg
	Path string `json:"path"`
	// Url 绝对地址（前端预览用）
	Url string `json:"url"`
	// Name 原始文件名
	Name string `json:"name"`
	// Size 文件大小（字节）
	Size int64 `json:"size"`
}

// StorageDir 附件存放的物理目录（config.yaml 的 app.uploadDir）
func StorageDir(ctx context.Context) string {
	dir := strings.TrimSpace(g.Cfg().MustGet(ctx, "app.uploadDir", DefaultUploadDir).String())
	if dir == "" {
		return DefaultUploadDir
	}

	return dir
}

// URLPathPrefix 附件对外访问路径前缀（config.yaml 的 app.uploadUrlPath）
func URLPathPrefix(ctx context.Context) string {
	raw := strings.Trim(strings.TrimSpace(g.Cfg().MustGet(ctx, "app.uploadUrlPath", DefaultURLPath).String()), "/")
	if raw == "" {
		return ""
	}

	return "/" + raw
}

// SaveImage 校验并保存单张图片，返回相对路径与绝对地址
func SaveImage(ctx context.Context, file *ghttp.UploadFile) (*SavedFile, error) {
	if file == nil {
		return nil, gerror.New("请选择要上传的图片")
	}

	ext := strings.ToLower(gfile.Ext(file.Filename))
	if _, ok := allowedExts[ext]; !ok {
		return nil, gerror.New("仅支持 jpg / jpeg / png / gif / webp / bmp 格式的图片")
	}

	if file.Size <= 0 {
		return nil, gerror.New("上传的文件为空")
	}

	if file.Size > MaxFileSize {
		return nil, gerror.New("单张图片不能超过 10MB")
	}

	rootDir := StorageDir(ctx)
	subDir := gtime.Now().Format("Ymd")
	fullDir := gfile.Join(rootDir, subDir)

	// 随机重命名：既避免同名覆盖，也规避原始文件名带来的路径穿越风险
	fileName, err := file.Save(fullDir, true)
	if err != nil {
		return nil, gerror.New("保存图片失败，请稍后重试")
	}

	relativePath := path.Join("/", URLPathPrefix(ctx), subDir, fileName)

	return &SavedFile{
		Path: relativePath,
		Url:  AbsoluteURL(ctx, relativePath),
		Name: gfile.Basename(file.Filename),
		Size: file.Size,
	}, nil
}

func SaveContractFile(ctx context.Context, file *ghttp.UploadFile) (*SavedFile, error) {
	if file == nil {
		return nil, gerror.New("请选择要上传的合同附件")
	}
	ext := strings.ToLower(gfile.Ext(file.Filename))
	if _, ok := allowedContractExts[ext]; !ok {
		return nil, gerror.New("合同附件仅支持 PDF、Word、Excel、JPG、PNG 格式")
	}
	if file.Size <= 0 {
		return nil, gerror.New("上传的文件为空")
	}
	if file.Size > MaxContractFileSize {
		return nil, gerror.New("单个合同附件不能超过 20MB")
	}

	rootDir := StorageDir(ctx)
	subDir := gtime.Now().Format("Ymd")
	fullDir := gfile.Join(rootDir, "contracts", subDir)
	fileName, err := file.Save(fullDir, true)
	if err != nil {
		return nil, gerror.New("保存合同附件失败，请稍后重试")
	}
	relativePath := path.Join("/", URLPathPrefix(ctx), "contracts", subDir, fileName)
	return &SavedFile{Path: relativePath, Url: AbsoluteURL(ctx, relativePath), Name: gfile.Basename(file.Filename), Size: file.Size}, nil
}

func ContractUploadLimitText() string {
	return fmt.Sprintf("每个附件最大 %dMB", MaxContractFileSize/(1024*1024))
}

// AbsoluteURL 相对路径 → 可访问的绝对地址
//
// 优先取配置 app.uploadUrlPrefix（切 OSS / 配 CDN 时填这里），
// 未配置则按当前请求的 scheme://host 推导，本地开发无需任何配置。
func AbsoluteURL(ctx context.Context, relativePath string) string {
	relativePath = strings.TrimSpace(relativePath)
	if relativePath == "" {
		return ""
	}

	// 已经是完整地址（例如以后存的 OSS 链接）直接返回
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}

	prefix := g.Cfg().MustGet(ctx, "app.uploadUrlPrefix").String()
	if prefix == "" {
		prefix = requestOrigin(ctx)
	}

	return strings.TrimSuffix(prefix, "/") + relativePath
}

// AbsoluteURLs 批量把相对路径转成绝对地址，空值自动跳过
func AbsoluteURLs(ctx context.Context, relativePaths []string) []string {
	result := make([]string, 0, len(relativePaths))

	for _, item := range relativePaths {
		if url := AbsoluteURL(ctx, item); url != "" {
			result = append(result, url)
		}
	}

	return result
}

// requestOrigin 从当前请求推导 http(s)://host
func requestOrigin(ctx context.Context) string {
	request := g.RequestFromCtx(ctx)
	if request == nil {
		return ""
	}

	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}

	// 兼容反向代理场景
	if proto := request.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	return scheme + "://" + request.Host
}
