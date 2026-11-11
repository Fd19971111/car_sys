// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProjectLog is the golang structure of table project_log for DAO operations like Where/Data.
type ProjectLog struct {
	g.Meta    `orm:"table:project_log, do:true"`
	Id        any         //
	ProjectId any         // 项目ID
	UserId    any         // 记录人ID（site_user.id）
	Content   any         // 进度内容
	LogDate   any         // 进度业务日期（YYYY-MM-DD，进度实际发生日期）
	Images    any         // 图片URL，多个用逗号分隔
	CreatedAt *gtime.Time // 创建时间
}
