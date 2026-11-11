package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// AccountItem 会计科目列表项。
type AccountItem struct {
	Id            int64  `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	CategoryCode  string `json:"categoryCode"`
	DirectionCode string `json:"directionCode"`
	Enabled       bool   `json:"enabled"`
	IsSystem      bool   `json:"isSystem"`
	SortOrder     int    `json:"sortOrder"`
	Used          bool   `json:"used"` // 是否已被凭证、余额或资金账户引用
}

// AccountSummary 会计科目统计。
type AccountSummary struct {
	Total            int `json:"total"`
	Enabled          int `json:"enabled"`
	Disabled         int `json:"disabled"`
	System           int `json:"system"`
	FundAccountCount int `json:"fundAccountCount"`
}

// AccountListReq 查询会计科目。
type AccountListReq struct {
	g.Meta       `path:"/fin-basic/accounts" method:"get" tags:"财务基础资料" summary:"会计科目列表"`
	Keyword      string `json:"keyword" dc:"编码或名称关键字"`
	CategoryCode string `json:"categoryCode" dc:"科目类别"`
	Enabled      *bool  `json:"enabled" dc:"启用状态，不传为全部"`
}

// AccountListRes 会计科目列表结果。
type AccountListRes struct {
	Records []*AccountItem `json:"records"`
	Summary AccountSummary `json:"summary"`
}

// AccountSaveReq 新增或修改会计科目。
type AccountSaveReq struct {
	g.Meta       `path:"/fin-basic/account/save" method:"post" tags:"财务基础资料" summary:"保存会计科目"`
	Id           int64  `json:"id" dc:"科目ID，为0表示新增"`
	Code         string `json:"code" v:"required|max-length:32" dc:"科目编码"`
	Name         string `json:"name" v:"required|max-length:100" dc:"科目名称"`
	CategoryCode string `json:"categoryCode" v:"required|in:ASSET,LIABILITY,EQUITY,REVENUE,EXPENSE" dc:"科目类别"`
	Enabled      bool   `json:"enabled" dc:"是否启用"`
}

// AccountSaveRes 保存结果。
type AccountSaveRes struct {
	Id int64 `json:"id"`
}

// AccountDeleteReq 删除会计科目。
type AccountDeleteReq struct {
	g.Meta `path:"/fin-basic/account/delete" method:"post" tags:"财务基础资料" summary:"删除会计科目"`
	Id     int64 `json:"id" v:"required" dc:"科目ID"`
}

// AccountDeleteRes 删除结果。
type AccountDeleteRes struct{}

// AccountToggleReq 启用或停用会计科目。
type AccountToggleReq struct {
	g.Meta  `path:"/fin-basic/account/toggle" method:"post" tags:"财务基础资料" summary:"启用停用会计科目"`
	Id      int64 `json:"id" v:"required" dc:"科目ID"`
	Enabled bool  `json:"enabled" dc:"目标状态"`
}

// AccountToggleRes 操作结果。
type AccountToggleRes struct{}

// FundAccountItem 资金账户列表项。
type FundAccountItem struct {
	Id                int64  `json:"id"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	BankName          string `json:"bankName"`
	AccountLast4      string `json:"accountLast4"`
	LedgerAccountId   int64  `json:"ledgerAccountId"`
	LedgerAccountCode string `json:"ledgerAccountCode"`
	LedgerAccountName string `json:"ledgerAccountName"`
	Enabled           bool   `json:"enabled"`
}

// FundAccountListReq 查询资金账户。
type FundAccountListReq struct {
	g.Meta  `path:"/fin-basic/fund-accounts" method:"get" tags:"财务基础资料" summary:"资金账户列表"`
	Keyword string `json:"keyword" dc:"账户名称、银行或科目关键字"`
	Type    string `json:"type" dc:"账户类型 CASH/BANK"`
}

// FundAccountListRes 资金账户列表结果。
type FundAccountListRes struct {
	Records []*FundAccountItem `json:"records"`
}

// FundAccountSaveReq 新增或修改资金账户。
type FundAccountSaveReq struct {
	g.Meta            `path:"/fin-basic/fund-account/save" method:"post" tags:"财务基础资料" summary:"保存资金账户"`
	Id                int64  `json:"id" dc:"资金账户ID，为0表示新增"`
	Type              string `json:"type" v:"required|in:CASH,BANK" dc:"账户类型"`
	Name              string `json:"name" v:"required|max-length:100" dc:"账户名称"`
	BankName          string `json:"bankName" v:"max-length:100" dc:"开户银行"`
	AccountLast4      string `json:"accountLast4" v:"max-length:4" dc:"卡号末四位"`
	LedgerAccountCode string `json:"ledgerAccountCode" v:"required" dc:"关联资产科目编码"`
	Enabled           bool   `json:"enabled" dc:"是否启用"`
}

// FundAccountSaveRes 保存结果。
type FundAccountSaveRes struct {
	Id int64 `json:"id"`
}

// FundAccountDeleteReq 删除资金账户。
type FundAccountDeleteReq struct {
	g.Meta `path:"/fin-basic/fund-account/delete" method:"post" tags:"财务基础资料" summary:"删除资金账户"`
	Id     int64 `json:"id" v:"required" dc:"资金账户ID"`
}

// FundAccountDeleteRes 删除结果。
type FundAccountDeleteRes struct{}

// FundAccountToggleReq 启用或停用资金账户。
type FundAccountToggleReq struct {
	g.Meta  `path:"/fin-basic/fund-account/toggle" method:"post" tags:"财务基础资料" summary:"启用停用资金账户"`
	Id      int64 `json:"id" v:"required" dc:"资金账户ID"`
	Enabled bool  `json:"enabled" dc:"目标状态"`
}

// FundAccountToggleRes 操作结果。
type FundAccountToggleRes struct{}
