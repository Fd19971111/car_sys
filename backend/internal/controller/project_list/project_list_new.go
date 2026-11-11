package project_list

import (
	"strings"

	projectListApi "engineServ/api/project_list"
	projectListLogic "engineServ/internal/logic/api/project_list"

	"github.com/gogf/gf/v2/os/gtime"
)

type ControllerV1 struct {
	project *projectListLogic.Project
}

func NewV1() projectListApi.IProjectListV1 {
	return &ControllerV1{project: projectListLogic.New()}
}

// formatDateTime 格式化为 YYYY-MM-DD HH:mm:ss，空值返回空串。
func formatDateTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}

	return t.String()
}

// formatDate 格式化为 YYYY-MM-DD，空值返回空串。
func formatDate(t *gtime.Time) string {
	if t == nil {
		return ""
	}

	return t.Format("Y-m-d")
}

// parseImages 把 project_log.images（逗号分隔）拆成图片URL数组。
func parseImages(images string) []string {
	result := make([]string, 0)

	for _, item := range strings.Split(images, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
