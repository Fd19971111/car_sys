package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"engineServ/internal/consts"
	defaultDao "engineServ/model/default/dao"
	defaultDo "engineServ/model/default/model/do"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// 客户状态：与 customer.status 注释保持一致（customer 表结构固定，不新增/修改列）
const (
	statusFollowing = 1 // 跟进中
	statusDeal      = 2 // 已成交
	statusLost      = 3 // 已流失
)

// Customer 客户服务。
type Customer struct{}

func New() *Customer {
	return &Customer{}
}

// List 客户分页列表：关键词（姓名/手机号/房屋地址）+ 状态 + 来源过滤，按创建时间倒序。
// 同时返回每个客户名下未删除项目的数量映射。
func (c *Customer) List(
	ctx context.Context,
	keyword string,
	source string,
	status *int,
	current int,
	size int,
) (list []*defaultEntity.Customer, total int, err error) {
	cols := defaultDao.Customer.Columns()
	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		// 用括号包裹 OR 条件，避免与后面的 AND 过滤出现优先级问题
		kwCondition := fmt.Sprintf("(%s LIKE ? OR %s LIKE ? OR %s LIKE ?)", cols.Name, cols.Phone, cols.Address)
		model = model.Where(kwCondition, like, like, like)
	}

	// 注意：前端未选择状态时会带 `status=` 空串请求，GoFrame 会把它绑定成「指向 0 的非空指针」。
	// 若只判断 status != nil，SQL 会变成 WHERE status = 0，而客户状态只有 1/2/3，结果必然为空。
	// 因此这里统一用 > 0 作为「有筛选」的判据。
	if status != nil && *status > 0 {
		model = model.Where(cols.Status, *status)
	}

	source = strings.TrimSpace(source)
	if source != "" {
		model = model.Where(cols.Source, source)
	}

	err = model.
		Page(current, size).
		OrderDesc(cols.CreatedAt).
		OrderDesc(cols.Id).
		ScanAndCount(&list, &total, false)
	if err != nil {
		glog.Errorf(ctx, "查询客户分页列表失败: %v", err)
		return nil, 0, gerror.New("查询客户列表失败，请稍后重试")
	}

	return list, total, nil
}

// Stats 客户概览统计：各状态客户数量 + 已成交客户名下项目合同总金额。
func (c *Customer) Stats(ctx context.Context) (total, following, deal, lost int, dealAmount float64, err error) {
	cols := defaultDao.Customer.Columns()

	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx))
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	total, err = model.Count()
	if err != nil {
		glog.Errorf(ctx, "统计客户总数失败: %v", err)
		return 0, 0, 0, 0, 0, gerror.New("获取客户统计失败，请稍后重试")
	}

	countByStatus := func(status int) int {
		statusModel := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(cols.Status, status)
		if !isSuperUser(ctx) {
			statusModel = statusModel.Where("owner_user_id", currentUserId(ctx))
		}
		count, countErr := statusModel.Count()
		if countErr != nil {
			glog.Errorf(ctx, "统计客户状态(%d)失败: %v", status, countErr)
			err = gerror.New("获取客户统计失败，请稍后重试")
		}

		return count
	}

	following = countByStatus(statusFollowing)
	if err != nil {
		return
	}

	deal = countByStatus(statusDeal)
	if err != nil {
		return
	}

	lost = countByStatus(statusLost)
	if err != nil {
		return
	}

	// 已成交客户（customer.status=2）名下已完成订单的最终费用合计
	dealValue, err := g.DB().GetValue(ctx, `SELECT COALESCE(SUM(o.final_fare),0)
		FROM ride_order o
		INNER JOIN customer cu ON cu.id=o.customer_id
		WHERE o.company_id=? AND o.deleted_at IS NULL AND o.status='completed'
		  AND cu.company_id=? AND cu.deleted_at IS NULL AND cu.status=?`,
		currentCompanyId(ctx), currentCompanyId(ctx), statusDeal)
	if err != nil {
		glog.Errorf(ctx, "统计成交订单金额失败: %v", err)
		return 0, 0, 0, 0, 0, gerror.New("获取客户统计失败，请稍后重试")
	}
	dealAmount = dealValue.Float64()

	return total, following, deal, lost, dealAmount, nil
}

// Detail 客户详情。
func (c *Customer) Detail(ctx context.Context, id uint) (customer *defaultEntity.Customer, err error) {
	return c.getCustomer(ctx, id)
}

// Create 新增客户，返回新客户ID。
func (c *Customer) Create(ctx context.Context, data *defaultDo.Customer) (id uint, err error) {
	insertData := g.Map{
		"company_id": currentCompanyId(ctx), "name": data.Name, "phone": data.Phone, "address": data.Address,
		"source": data.Source, "status": data.Status, "remark": data.Remark,
	}
	if !isSuperUser(ctx) {
		insertData["owner_user_id"] = currentUserId(ctx)
	}
	result, err := defaultDao.Customer.Ctx(ctx).Data(insertData).Insert()
	if err != nil {
		glog.Errorf(ctx, "新增客户失败: %v", err)
		return 0, gerror.New("新增客户失败，请稍后重试")
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		glog.Errorf(ctx, "读取新增客户ID失败: %v", err)
		return 0, gerror.New("新增客户失败，请稍后重试")
	}

	return uint(lastInsertId), nil
}

// Update 更新客户基础信息（仅允许修改 customer 表既有列）。
func (c *Customer) Update(ctx context.Context, id uint, data *defaultDo.Customer) error {
	if _, err := c.getCustomer(ctx, id); err != nil {
		return err
	}

	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Customer.Columns().Id, id)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	_, err := model.
		Data(data).
		Update()
	if err != nil {
		glog.Errorf(ctx, "更新客户失败, id=%v: %v", data.Id, err)
		return gerror.New("更新客户失败，请稍后重试")
	}

	return nil
}

// Delete 删除单个客户（ORM 软删除，写入 deleted_at）。
func (c *Customer) Delete(ctx context.Context, id uint) error {
	if _, err := c.getCustomer(ctx, id); err != nil {
		return err
	}

	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Customer.Columns().Id, id)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	_, err := model.Delete()
	if err != nil {
		glog.Errorf(ctx, "删除客户失败, id=%d: %v", id, err)
		return gerror.New("删除客户失败，请稍后重试")
	}

	return nil
}

// BatchDelete 批量删除客户（ORM 软删除）。
func (c *Customer) BatchDelete(ctx context.Context, ids []uint) error {
	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).WhereIn(defaultDao.Customer.Columns().Id, ids)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	_, err := model.Delete()
	if err != nil {
		glog.Errorf(ctx, "批量删除客户失败, ids=%v: %v", ids, err)
		return gerror.New("删除客户失败，请稍后重试")
	}

	return nil
}

// getCustomer 按主键查询未删除客户，不存在时返回业务错误。
func (c *Customer) getCustomer(ctx context.Context, id uint) (*defaultEntity.Customer, error) {
	var customer defaultEntity.Customer
	model := defaultDao.Customer.Ctx(ctx).Where("company_id", currentCompanyId(ctx)).Where(defaultDao.Customer.Columns().Id, id)
	if !isSuperUser(ctx) {
		model = model.Where("owner_user_id", currentUserId(ctx))
	}
	err := model.Scan(&customer)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询客户失败, id=%d: %v", id, err)
		return nil, gerror.New("获取客户信息失败，请稍后重试")
	}

	if customer.Id == 0 {
		return nil, gerror.New("客户不存在或已被删除")
	}

	return &customer, nil
}

func currentUserId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxUserIdKey))
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}

func isSuperUser(ctx context.Context) bool {
	role := gconv.Uint(ctx.Value(consts.CtxRoleKey))
	return role == 1 || role == 3
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}
