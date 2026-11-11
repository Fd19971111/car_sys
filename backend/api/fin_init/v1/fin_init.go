package v1

import "github.com/gogf/gf/v2/frame/g"

type Account struct {
	Code         string `json:"code" v:"required|max-length:32" dc:"科目编码"`
	Name         string `json:"name" v:"required|max-length:100" dc:"科目名称"`
	CategoryCode string `json:"categoryCode" v:"required|in:ASSET,LIABILITY,EQUITY,REVENUE,EXPENSE" dc:"科目类别代码"`
	Enabled      bool   `json:"enabled" dc:"是否启用"`
}

type FundAccount struct {
	Type              string `json:"type" v:"required|in:CASH,BANK" dc:"资金账户类型：CASH/BANK"`
	Name              string `json:"name" v:"required|max-length:100" dc:"账户显示名称"`
	BankName          string `json:"bankName" v:"max-length:100" dc:"开户银行名称"`
	AccountLast4      string `json:"accountLast4" v:"max-length:4" dc:"银行卡末四位"`
	LedgerAccountCode string `json:"ledgerAccountCode" v:"required|max-length:32" dc:"关联资产科目编码"`
}

type OpeningBalance struct {
	AccountCode string `json:"accountCode" v:"required|max-length:32" dc:"科目编码"`
	Amount      string `json:"amount" dc:"科目期初余额，精确到分"`
}

type SaveReq struct {
	g.Meta          `path:"/fin-init/save" method:"post" sm:"保存财务初始化" tags:"Finance"`
	Year            int              `json:"year" v:"required|between:2000,2100" dc:"会计年度"`
	StartMonth      int              `json:"startMonth" v:"required|between:1,12" dc:"启用月份"`
	Accounts        []Account        `json:"accounts" v:"required" dc:"启用科目清单"`
	FundAccounts    []FundAccount    `json:"fundAccounts" dc:"现金及银行资金账户"`
	OpeningBalances []OpeningBalance `json:"openingBalances" dc:"启用期间各科目期初余额"`
	Reconfigure     bool             `json:"reconfigure" dc:"是否由管理员重新配置"`
}

type SaveRes struct {
	Initialized bool `json:"initialized" dc:"是否完成初始化"`
}
