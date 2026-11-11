// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectLog is the golang structure for table project_log.
type ProjectLog struct {
	Id        uint        `json:"id"        orm:"id"         description:""`                    //
	ProjectId uint        `json:"projectId" orm:"project_id" description:"项目ID"`                // 项目ID
	UserId    uint        `json:"userId"    orm:"user_id"    description:"记录人ID（site_user.id）"` // 记录人ID（site_user.id）
	Content   string      `json:"content"   orm:"content"    description:"进度内容"`                // 进度内容
	LogDate   *gtime.Time `json:"logDate"   orm:"log_date"    description:"进度业务日期（YYYY-MM-DD）"`   // 进度业务日期（YYYY-MM-DD）
	Images    string      `json:"images"    orm:"images"     description:"图片URL，多个用逗号分隔"`       // 图片URL，多个用逗号分隔
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                // 创建时间
}
