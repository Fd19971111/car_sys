package auth

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
	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

func New() *Auth {
	return &Auth{}
}

// Login 站点用户登录：校验 site_user 账号密码，签发 site_user_token。
// 返回 userId 与 token（refreshToken 与 token 同值，刷新时在原记录上轮换）。
func (a *Auth) Login(ctx context.Context, userName, password, client string) (userId uint, token string, err error) {
	client = normalizeClient(client)

	// 1. 查询用户（deleted_at 由 ORM 软删除特性自动过滤）
	var user defaultEntity.SiteUser
	err = defaultDao.SiteUser.Ctx(ctx).
		Where(defaultDao.SiteUser.Columns().Username, userName).
		Scan(&user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		glog.Errorf(ctx, "查询site_user失败: %v", err)
		return 0, "", gerror.New("登录失败，请稍后重试")
	}
	if user.Id == 0 {
		return 0, "", gerror.New("用户名或密码错误")
	}

	// 2. 校验状态与密码
	if user.Status != 1 {
		return 0, "", gerror.New("账号已被禁用，请联系管理员")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return 0, "", gerror.New("用户名或密码错误")
	}

	// 3. 签发令牌：同一用户同一客户端只保留最新一条（新登录踢掉旧会话）
	token = randomToken()
	expireAt := gtime.NewFromTime(time.Now().Add(a.tokenExpire(ctx)))
	ip := g.RequestFromCtx(ctx).GetClientIp()

	_, err = defaultDao.SiteUserToken.Ctx(ctx).
		Where(defaultDao.SiteUserToken.Columns().UserId, user.Id).
		Where(defaultDao.SiteUserToken.Columns().Client, client).
		Delete()
	if err != nil {
		glog.Errorf(ctx, "清理旧token失败: %v", err)
		return 0, "", gerror.New("登录失败，请稍后重试")
	}

	_, err = defaultDao.SiteUserToken.Ctx(ctx).Data(defaultDo.SiteUserToken{
		UserId:   user.Id,
		Token:    token,
		Client:   client,
		Ip:       ip,
		ExpireAt: expireAt,
	}).Insert()
	if err != nil {
		glog.Errorf(ctx, "写入site_user_token失败: %v", err)
		return 0, "", gerror.New("登录失败，请稍后重试")
	}

	// 4. 回写最后登录时间/IP
	_, _ = defaultDao.SiteUser.Ctx(ctx).
		Where(defaultDao.SiteUser.Columns().Id, user.Id).
		Data(defaultDo.SiteUser{
			LastLoginAt: gtime.Now(),
			LastLoginIp: ip,
		}).Update()

	return user.Id, token, nil
}

// RefreshToken 刷新令牌：校验 refreshToken 对应的 token 记录仍有效，然后轮换令牌并滑动续期。
func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (token string, err error) {
	_, tokenRow, err := a.validateToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	newToken := randomToken()
	expireAt := gtime.NewFromTime(time.Now().Add(a.tokenExpire(ctx)))

	_, err = defaultDao.SiteUserToken.Ctx(ctx).
		Where(defaultDao.SiteUserToken.Columns().Id, tokenRow.Id).
		Data(defaultDo.SiteUserToken{
			Token:    newToken,
			ExpireAt: expireAt,
		}).Update()
	if err != nil {
		glog.Errorf(ctx, "轮换token失败: %v", err)
		return "", authFailedError("登录已失效，请重新登录")
	}

	return newToken, nil
}

// Authenticate 供中间件调用：校验请求头中的 Bearer 令牌，返回当前登录用户。
// 令牌缺失/无效/用户被禁用 -> 8888（前端登出）；令牌过期 -> 9999（前端静默刷新）。
func (a *Auth) Authenticate(ctx context.Context) (*defaultEntity.SiteUser, error) {
	tokenString := strings.TrimPrefix(g.RequestFromCtx(ctx).Header.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		return nil, authFailedError("请先登录")
	}

	user, _, err := a.validateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// validateToken 校验令牌：查 site_user_token 记录 + 过期时间 + site_user 用户状态。
func (a *Auth) validateToken(ctx context.Context, tokenString string) (
	user *defaultEntity.SiteUser,
	tokenRow *defaultEntity.SiteUserToken,
	err error,
) {
	var row defaultEntity.SiteUserToken
	_ = defaultDao.SiteUserToken.Ctx(ctx).
		Where(defaultDao.SiteUserToken.Columns().Token, tokenString).
		Scan(&row)
	if row.Id == 0 {
		return nil, nil, authFailedError("登录已失效，请重新登录")
	}

	// 令牌已过期：提示前端走 refreshToken 续期流程
	if row.ExpireAt == nil || row.ExpireAt.Before(gtime.Now()) {
		return nil, nil, tokenExpiredError("登录已过期，请重新登录")
	}

	var siteUser defaultEntity.SiteUser
	_ = defaultDao.SiteUser.Ctx(ctx).
		Where(defaultDao.SiteUser.Columns().Id, row.UserId).
		Where(defaultDao.SiteUser.Columns().Status, 1).
		Scan(&siteUser)
	if siteUser.Id == 0 {
		return nil, nil, authFailedError("用户不存在或已被禁用")
	}

	return &siteUser, &row, nil
}

// GetUserInfo 获取登录用户信息及前端角色编码。
func (a *Auth) GetUserInfo(ctx context.Context, userId uint) (
	user *defaultEntity.SiteUser,
	roles []string,
	buttons []string,
	err error,
) {
	var siteUser defaultEntity.SiteUser
	_ = defaultDao.SiteUser.Ctx(ctx).
		Where(defaultDao.SiteUser.Columns().Id, userId).
		Where(defaultDao.SiteUser.Columns().Status, 1).
		Scan(&siteUser)
	if siteUser.Id == 0 {
		return nil, nil, nil, gerror.New("用户不存在或已被禁用")
	}

	// 角色映射：管理员 R_SUPER（静态路由超级角色），其余 R_USER
	if siteUser.Role == 1 {
		roles = []string{"R_SUPER"}
	} else {
		roles = []string{"R_USER"}
	}
	buttons = []string{}

	return &siteUser, roles, buttons, nil
}

// tokenExpire 读取 app.tokenExpire 配置，缺省 7 天。
func (a *Auth) tokenExpire(ctx context.Context) time.Duration {
	d := g.Cfg().MustGet(ctx, "app.tokenExpire").Duration()
	if d <= 0 {
		return consts.DefaultTokenExpire
	}
	return d
}

// normalizeClient 规范化客户端类型，非 app 一律按 web 处理。
func normalizeClient(client string) string {
	client = strings.ToLower(strings.TrimSpace(client))
	if client == "app" {
		return "app"
	}
	return "web"
}

// randomToken 生成 64 位十六进制随机令牌（32 字节随机数）。
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
