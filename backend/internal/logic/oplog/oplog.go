// Package oplog 操作日志：把关键业务动作落库，供「数据报表 → 操作日志」追溯。
//
// 使用方式统一为一行调用，业务方法成功后执行：
//
//	oplog.New().Record(ctx, &oplog.Entry{
//		Module:  oplog.ModuleReceipt,
//		Action:  oplog.ActionCreate,
//		Title:   "新增收款：张伟明 · 锦绣路全屋翻新（定金）",
//		Amount:  form.Amount,
//		Detail:  map[string]interface{}{"voucherNo": voucherNo},
//	})
//
// 日志失败只写错误日志，绝不回滚或阻断业务。
package oplog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"engineServ/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// 业务模块
const (
	ModuleReceipt  = "收款"
	ModulePayout   = "付款"
	ModuleProject  = "项目"
	ModuleStock    = "库存"
	ModuleCustomer = "客户"
)

// 操作动作
const (
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

// ActionLabels 动作中文名，列表与详情共用。
var ActionLabels = map[string]string{
	ActionCreate: "新增",
	ActionUpdate: "修改",
	ActionDelete: "删除",
}

// Entry 一条待写入的操作日志。
type Entry struct {
	// Module 业务模块，取 Module* 常量，如「收款」
	Module string
	// Action 动作，取 Action* 常量
	Action string
	// TargetType 对象类型，如 payment / payout，用于按单据反查操作历史
	TargetType string
	// TargetId 对象ID
	TargetId int64
	// Title 操作摘要，展示在列表主列
	Title string
	// Amount 涉及金额，无金额填 0
	Amount float64
	// Detail 变更详情，任意结构，落库为 JSON
	Detail interface{}
}

// Item 操作日志列表行。
type Item struct {
	Id         int64   `json:"id"         dc:"日志ID"`
	CreatedAt  string  `json:"createdAt"  dc:"操作时间 YYYY-MM-DD HH:mm:ss"`
	UserId     int64   `json:"userId"     dc:"操作人ID"`
	UserName   string  `json:"userName"   dc:"操作人"`
	Module     string  `json:"module"     dc:"业务模块"`
	Action     string  `json:"action"     dc:"动作 CREATE/UPDATE/DELETE"`
	ActionName string  `json:"actionName" dc:"动作中文名"`
	TargetType string  `json:"targetType" dc:"对象类型"`
	TargetId   int64   `json:"targetId"   dc:"对象ID"`
	Title      string  `json:"title"      dc:"操作摘要"`
	Amount     float64 `json:"amount"     dc:"涉及金额"`
	Detail     string  `json:"detail"     dc:"变更详情 JSON"`
	Ip         string  `json:"ip"         dc:"操作来源IP"`
}

// Filter 操作日志查询条件。
type Filter struct {
	Module    string
	Action    string
	Keyword   string
	StartDate string
	EndDate   string
	Current   int
	Size      int
}

// Service 操作日志服务。
type Service struct{}

// New 创建操作日志服务。
func New() *Service { return &Service{} }

// Record 写入一条操作日志。
//
// 约定：业务事务提交成功后调用；失败只记录错误日志、不向上抛错——
// 审计日志不应该成为业务成功与否的变量。
func (s *Service) Record(ctx context.Context, entry *Entry) {
	if entry == nil || strings.TrimSpace(entry.Title) == "" {
		return
	}

	userId := gconv.Uint(ctx.Value(consts.CtxUserIdKey))

	_, err := g.DB().Model("sys_operation_log").Ctx(ctx).Data(g.Map{
		"company_id":  currentCompanyId(ctx),
		"user_id":     userId,
		"user_name":   s.operatorName(ctx, userId),
		"module":      strings.TrimSpace(entry.Module),
		"action":      strings.ToUpper(strings.TrimSpace(entry.Action)),
		"target_type": strings.TrimSpace(entry.TargetType),
		"target_id":   entry.TargetId,
		"title":       limitRunes(entry.Title, 255),
		"amount":      entry.Amount,
		"detail":      encodeDetail(entry.Detail),
		"ip":          clientIp(ctx),
	}).Insert()
	if err != nil {
		glog.Errorf(ctx, "写入操作日志失败, module=%s, action=%s, title=%s: %v",
			entry.Module, entry.Action, entry.Title, err)
	}
}

// List 分页查询操作日志。
func (s *Service) List(ctx context.Context, filter *Filter) ([]*Item, int, []string, error) {
	where := []string{"company_id=?"}
	args := []interface{}{currentCompanyId(ctx)}

	if module := strings.TrimSpace(filter.Module); module != "" {
		where = append(where, "module=?")
		args = append(args, module)
	}
	if action := strings.TrimSpace(filter.Action); action != "" {
		where = append(where, "action=?")
		args = append(args, strings.ToUpper(action))
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		where = append(where, "(title LIKE ? OR user_name LIKE ? OR detail LIKE ?)")
		like := "%" + keyword + "%"
		args = append(args, like, like, like)
	}
	if start := strings.TrimSpace(filter.StartDate); start != "" {
		where = append(where, "created_at>=?")
		args = append(args, start+" 00:00:00")
	}
	if end := strings.TrimSpace(filter.EndDate); end != "" {
		where = append(where, "created_at<=?")
		args = append(args, end+" 23:59:59")
	}

	whereSQL := strings.Join(where, " AND ")

	total, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM sys_operation_log WHERE "+whereSQL, args...)
	if err != nil {
		glog.Errorf(ctx, "统计操作日志失败: %v", err)
		return nil, 0, nil, gerror.New("查询操作日志失败，请稍后重试")
	}

	modules, err := s.modules(ctx)
	if err != nil {
		// 模块下拉仅用于筛选，失败不影响列表
		glog.Errorf(ctx, "查询操作日志模块失败: %v", err)
		modules = []string{}
	}

	if total.Int() == 0 {
		return []*Item{}, 0, modules, nil
	}

	current, size := normalizePage(filter.Current, filter.Size)
	queryArgs := append(append([]interface{}{}, args...), (current-1)*size, size)

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT id,user_id,user_name,module,action,target_type,target_id,title,amount,ip,
			COALESCE(detail,'') AS detail,
			DATE_FORMAT(created_at,'%%Y-%%m-%%d %%H:%%i:%%s') AS created_at
		FROM sys_operation_log
		WHERE %s
		ORDER BY id DESC
		LIMIT ?,?`, whereSQL), queryArgs...)
	if err != nil {
		glog.Errorf(ctx, "查询操作日志失败: %v", err)
		return nil, 0, nil, gerror.New("查询操作日志失败，请稍后重试")
	}

	items := make([]*Item, 0, len(records))
	for _, row := range records {
		action := row["action"].String()
		items = append(items, &Item{
			Id:         row["id"].Int64(),
			CreatedAt:  row["created_at"].String(),
			UserId:     row["user_id"].Int64(),
			UserName:   row["user_name"].String(),
			Module:     row["module"].String(),
			Action:     action,
			ActionName: actionLabel(action),
			TargetType: row["target_type"].String(),
			TargetId:   row["target_id"].Int64(),
			Title:      row["title"].String(),
			Amount:     round2(row["amount"].Float64()),
			Detail:     row["detail"].String(),
			Ip:         row["ip"].String(),
		})
	}

	return items, total.Int(), modules, nil
}

// modules 已出现过的业务模块，供前端筛选下拉。
func (s *Service) modules(ctx context.Context) ([]string, error) {
	values, err := g.DB().GetAll(ctx, `SELECT DISTINCT module FROM sys_operation_log
		WHERE company_id=? AND module<>'' ORDER BY module`, currentCompanyId(ctx))
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(values))
	for _, row := range values {
		result = append(result, row["module"].String())
	}

	return result, nil
}

// operatorName 操作人姓名：优先 site_user.nickname，取不到时退回登录名。
func (s *Service) operatorName(ctx context.Context, userId uint) string {
	if userId == 0 {
		return "系统"
	}

	value, err := g.DB().Model("site_user").Ctx(ctx).Where("id", userId).Value("nickname")
	if err == nil && !value.IsEmpty() {
		return value.String()
	}

	if name := gconv.String(ctx.Value(consts.CtxUsernameKey)); name != "" {
		return name
	}

	return "系统"
}

// ===================== 内部辅助 =====================

func actionLabel(action string) string {
	if label, ok := ActionLabels[action]; ok {
		return label
	}

	return action
}

func encodeDetail(detail interface{}) string {
	if detail == nil {
		return ""
	}

	if text, ok := detail.(string); ok {
		return text
	}

	buf, err := json.Marshal(detail)
	if err != nil {
		return ""
	}

	return string(buf)
}

func clientIp(ctx context.Context) string {
	if request := ghttp.RequestFromCtx(ctx); request != nil {
		return request.GetClientIp()
	}

	return ""
}

// limitRunes 按字符截断，避免超长摘要写库报错。
func limitRunes(text string, max int) string {
	runes := []rune(strings.TrimSpace(text))
	if len(runes) <= max {
		return string(runes)
	}

	return string(runes[:max])
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func normalizePage(current, size int) (int, int) {
	if current <= 0 {
		current = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 200 {
		size = 200
	}

	return current, size
}

func round2(value float64) float64 {
	return float64(int64(value*100+0.5)) / 100
}
