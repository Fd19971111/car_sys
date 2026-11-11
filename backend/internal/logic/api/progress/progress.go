// Package progress 施工进度（project_log）写入：后台与小程序端共用同一套落库逻辑。
package progress

import (
	"context"
	"strings"

	defaultDao "engineServ/model/default/dao"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	// MaxImages 单条进度最多关联的图片数。
	//
	// project_log.images 是 VARCHAR(1000)，按相对路径（约 32 字符/张）估算，
	// 9 张约 300 字符，留足余量；表结构固定，因此在这里做上限约束。
	MaxImages = 9

	// MaxContentLength 进度描述长度上限，与 project_log.content VARCHAR(500) 对齐
	MaxContentLength = 500

	// maxImagesFieldLength images 字段的实际长度上限
	maxImagesFieldLength = 1000
)

// Form 新增施工进度入参
type Form struct {
	ProjectId uint
	UserId    uint
	Content   string
	Images    []string
	// LogDate 进度业务日期（YYYY-MM-DD，进度实际发生日期）。
	// 为空时取当天；补录历史进度时传入真实日期。
	LogDate string
}

type Progress struct{}

func New() *Progress {
	return &Progress{}
}

// Create 新增一条施工进度记录，返回记录ID。
func (p *Progress) Create(ctx context.Context, form *Form) (int64, error) {
	if form.ProjectId == 0 {
		return 0, gerror.New("请选择项目")
	}

	content := strings.TrimSpace(form.Content)
	if content == "" {
		return 0, gerror.New("请填写进度描述")
	}

	if len([]rune(content)) > MaxContentLength {
		return 0, gerror.Newf("进度描述最多 %d 个字符", MaxContentLength)
	}

	images := normalizeImages(form.Images)
	if len(images) > MaxImages {
		return 0, gerror.Newf("单条进度最多上传 %d 张图片", MaxImages)
	}

	imagesStr := strings.Join(images, ",")
	if len(imagesStr) > maxImagesFieldLength {
		return 0, gerror.New("图片地址过长，请减少图片数量后重试")
	}

	// 业务日期：为空取当天，格式必须为 YYYY-MM-DD
	logDate := strings.TrimSpace(form.LogDate)
	if logDate == "" {
		logDate = gtime.Now().Format("Y-m-d")
	} else if gtime.New(logDate).Format("Y-m-d") != logDate {
		return 0, gerror.New("进度日期格式不正确，应为 YYYY-MM-DD")
	}

	// 校验项目存在，避免挂到已删除 / 不存在的项目上
	projectCols := defaultDao.Project.Columns()

	count, err := defaultDao.Project.Ctx(ctx).Where(projectCols.Id, form.ProjectId).Count()
	if err != nil {
		glog.Errorf(ctx, "校验进度所属项目失败, projectId=%d: %v", form.ProjectId, err)
		return 0, gerror.New("保存施工进度失败，请稍后重试")
	}

	if count == 0 {
		return 0, gerror.New("所选项目不存在或已被删除")
	}

	cols := defaultDao.ProjectLog.Columns()

	result, err := defaultDao.ProjectLog.Ctx(ctx).Data(g.Map{
		cols.ProjectId: form.ProjectId,
		cols.UserId:    form.UserId,
		cols.Content:   content,
		cols.LogDate:   logDate,
		cols.Images:    imagesStr,
	}).Insert()
	if err != nil {
		glog.Errorf(ctx, "新增施工进度失败, projectId=%d: %v", form.ProjectId, err)
		return 0, gerror.New("保存施工进度失败，请稍后重试")
	}

	id, err := result.LastInsertId()
	if err != nil {
		glog.Errorf(ctx, "读取施工进度ID失败: %v", err)
		return 0, gerror.New("保存施工进度失败，请稍后重试")
	}

	return id, nil
}

// normalizeImages 去空、去重并保持原有顺序
func normalizeImages(images []string) []string {
	result := make([]string, 0, len(images))
	seen := make(map[string]struct{}, len(images))

	for _, item := range images {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}

		if _, ok := seen[trimmed]; ok {
			continue
		}

		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}
