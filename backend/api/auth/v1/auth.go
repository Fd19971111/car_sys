package v1

import "github.com/gogf/gf/v2/frame/g"

// LoginReq 登录请求
type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" sm:"站点用户登录" tags:"Auth"`
	UserName string `json:"userName" v:"required" dc:"登录账号"`
	Password string `json:"password" v:"required" dc:"登录密码"`
	Client   string `json:"client" dc:"客户端类型：web/app，默认web"`
}

// LoginRes 登录响应
type LoginRes struct {
	UserId       uint   `json:"userId" dc:"用户ID"`
	Token        string `json:"token" dc:"访问令牌"`
	RefreshToken string `json:"refreshToken" dc:"刷新令牌（当前与token相同，过期后调用刷新接口轮换）"`
}

// RefreshTokenReq 刷新令牌请求
type RefreshTokenReq struct {
	g.Meta       `path:"/auth/refreshToken" method:"post" sm:"刷新令牌" tags:"Auth"`
	RefreshToken string `json:"refreshToken" v:"required" dc:"登录时返回的refreshToken"`
}

// RefreshTokenRes 刷新令牌响应
type RefreshTokenRes struct {
	Token        string `json:"token" dc:"新的访问令牌"`
	RefreshToken string `json:"refreshToken" dc:"新的刷新令牌"`
}

// GetUserInfoReq 获取用户信息请求
type GetUserInfoReq struct {
	g.Meta `path:"/auth/getUserInfo" method:"get" sm:"获取用户信息" tags:"Auth"`
}

// GetUserInfoRes 获取用户信息响应
type GetUserInfoRes struct {
	UserId         uint     `json:"userId" dc:"用户ID"`
	UserName       string   `json:"userName" dc:"登录账号"`
	Nickname       string   `json:"nickname" dc:"昵称/姓名"`
	Avatar         string   `json:"avatar" dc:"头像URL"`
	Phone          string   `json:"phone" dc:"手机号"`
	Role           uint     `json:"role" dc:"角色：1=管理员 2=员工 3=商家"`
	MerchantId     uint     `json:"merchantId" dc:"所属商家ID"`
	CompanyId      uint     `json:"companyId" dc:"所属公司ID"`
	FinInitialized bool     `json:"finInitialized" dc:"财务初始化是否完成"`
	Roles          []string `json:"roles" dc:"前端角色编码列表"`
	Buttons        []string `json:"buttons" dc:"按钮权限列表"`
}
