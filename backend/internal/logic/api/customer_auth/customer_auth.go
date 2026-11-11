// Package customer_auth 业主端（小程序）登录与令牌校验。
//
// 与后台用户体系的区别：
//   - 业主不是 site_user，账号主体是 customer 表的一行，以「手机号」作为登录凭据
//   - 令牌单独存 customer_token，与 site_user_token 完全隔离，互不踢线
package customer_auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"engineServ/internal/consts"
	defaultDao "engineServ/model/default/dao"
	defaultDo "engineServ/model/default/model/do"
	defaultEntity "engineServ/model/default/model/entity"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

// 客户端类型：小程序端固定 mp，用来和其它端区分会话
const clientMiniProgram = "mp"

type CustomerAuth struct{}

func New() *CustomerAuth {
	return &CustomerAuth{}
}

// Login 业主登录：校验手机号 + 验证码，签发 customer_token。
//
// 验证码当前由配置 app.mpLoginCode 提供（开发阶段固定值），
// 接入短信服务后，把 verifyCode 换成「从短信平台校验」即可，其余流程不变。
func (a *CustomerAuth) Login(
	ctx context.Context,
	phone, code string,
) (token string, expireAt *gtime.Time, customer *defaultEntity.Customer, err error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "", nil, nil, gerror.New("请输入手机号")
	}

	if err := a.verifyCode(ctx, code); err != nil {
		return "", nil, nil, err
	}

	customer, err = a.getCustomerByPhone(ctx, phone)
	if err != nil {
		return "", nil, nil, err
	}

	// 同名手机号只保留一条会话：换设备登录会把旧令牌顶掉
	_, err = defaultDao.CustomerToken.Ctx(ctx).
		Where(defaultDao.CustomerToken.Columns().CustomerId, customer.Id).
		Where(defaultDao.CustomerToken.Columns().Client, clientMiniProgram).
		Delete()
	if err != nil {
		glog.Errorf(ctx, "清理旧业主令牌失败: %v", err)
		return "", nil, nil, gerror.New("登录失败，请稍后重试")
	}

	token = randomToken()
	expireAt = gtime.NewFromTime(time.Now().Add(a.tokenExpire(ctx)))
	ip := g.RequestFromCtx(ctx).GetClientIp()

	_, err = defaultDao.CustomerToken.Ctx(ctx).Data(defaultDo.CustomerToken{
		CustomerId: customer.Id,
		Token:      token,
		Client:     clientMiniProgram,
		Ip:         ip,
		ExpireAt:   expireAt,
	}).Insert()
	if err != nil {
		glog.Errorf(ctx, "写入customer_token失败: %v", err)
		return "", nil, nil, gerror.New("登录失败，请稍后重试")
	}

	return token, expireAt, customer, nil
}

// Logout 退出登录：删除当前令牌。
func (a *CustomerAuth) Logout(ctx context.Context, token string) error {
	_, err := defaultDao.CustomerToken.Ctx(ctx).
		Where(defaultDao.CustomerToken.Columns().Token, token).
		Delete()
	if err != nil {
		glog.Errorf(ctx, "退出登录失败: %v", err)
		return gerror.New("退出登录失败，请稍后重试")
	}

	return nil
}

// Authenticate 供中间件调用：校验 Bearer 令牌并返回业主ID。
//
// 令牌缺失/无效/业主被删除 -> 8888（前端清登录态跳登录页）
// 令牌已过期            -> 9999（前端同样跳登录页；业主端暂不做静默续期）
func (a *CustomerAuth) Authenticate(ctx context.Context) (uint, error) {
	tokenString := strings.TrimPrefix(g.RequestFromCtx(ctx).Header.Get("Authorization"), "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return 0, authFailedError("请先登录")
	}

	var row defaultEntity.CustomerToken
	_ = defaultDao.CustomerToken.Ctx(ctx).
		Where(defaultDao.CustomerToken.Columns().Token, tokenString).
		Scan(&row)

	if row.Id == 0 {
		return 0, authFailedError("登录已失效，请重新登录")
	}

	if row.ExpireAt == nil || row.ExpireAt.Before(gtime.Now()) {
		return 0, tokenExpiredError("登录已过期，请重新登录")
	}

	// 业主还在（未被删除）
	count, err := defaultDao.Customer.Ctx(ctx).
		Where(defaultDao.Customer.Columns().Id, row.CustomerId).
		Count()
	if err != nil {
		glog.Errorf(ctx, "校验业主失败, customerId=%d: %v", row.CustomerId, err)
		return 0, gerror.New("登录校验失败，请稍后重试")
	}

	if count == 0 {
		return 0, authFailedError("账号不存在，请联系您的装修顾问")
	}

	// 滑动续期：剩余有效期不足一半时自动延长，持续使用的业主不会被登出
	a.renewIfNeeded(ctx, &row)

	return row.CustomerId, nil
}

// CurrentCustomerId 从请求上下文取当前登录业主ID（由 CustomerAuth 中间件注入）。
func (a *CustomerAuth) CurrentCustomerId(ctx context.Context) uint {
	return g.RequestFromCtx(ctx).GetCtxVar(consts.CtxCustomerIdKey).Uint()
}

// verifyCode 校验验证码。
//
// 开发阶段从配置 app.mpLoginCode 读取固定值；接入短信后改成查短信平台 / 校验缓存中的动态码。
func (a *CustomerAuth) verifyCode(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return gerror.New("请输入验证码")
	}

	expected := g.Cfg().MustGet(ctx, "app.mpLoginCode", "123456").String()
	if code != expected {
		return gerror.New("验证码不正确")
	}

	return nil
}

// getCustomerByPhone 按手机号查业主，查不到给出可操作的提示。
func (a *CustomerAuth) getCustomerByPhone(ctx context.Context, phone string) (*defaultEntity.Customer, error) {
	var customer defaultEntity.Customer
	err := defaultDao.Customer.Ctx(ctx).
		Where(defaultDao.Customer.Columns().Phone, phone).
		Scan(&customer)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "按手机号查询业主失败: %v", err)
		return nil, gerror.New("登录失败，请稍后重试")
	}

	if customer.Id == 0 {
		return nil, gerror.New("该手机号暂未登记，请联系您的装修顾问")
	}

	return &customer, nil
}

// renewIfNeeded 剩余有效期不足一半时延长令牌
func (a *CustomerAuth) renewIfNeeded(ctx context.Context, row *defaultEntity.CustomerToken) {
	if row.ExpireAt == nil {
		return
	}

	ttl := a.tokenExpire(ctx)
	if time.Until(row.ExpireAt.Time) > ttl/2 {
		return
	}

	_, err := defaultDao.CustomerToken.Ctx(ctx).
		Where(defaultDao.CustomerToken.Columns().Id, row.Id).
		Data(defaultDo.CustomerToken{
			ExpireAt: gtime.NewFromTime(time.Now().Add(ttl)),
		}).Update()
	if err != nil {
		// 续期失败不影响本次请求
		glog.Errorf(ctx, "业主令牌续期失败, id=%d: %v", row.Id, err)
	}
}

// tokenExpire 读取 app.mpTokenExpire，缺省 30 天。
func (a *CustomerAuth) tokenExpire(ctx context.Context) time.Duration {
	d := g.Cfg().MustGet(ctx, "app.mpTokenExpire").Duration()
	if d <= 0 {
		return consts.DefaultMpTokenExpire
	}

	return d
}

func randomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}

func authFailedError(msg string) error {
	return gerror.NewCode(gcode.New(consts.CodeAuthFailed, "", nil), msg)
}

func tokenExpiredError(msg string) error {
	return gerror.NewCode(gcode.New(consts.CodeTokenExpired, "", nil), msg)
}
