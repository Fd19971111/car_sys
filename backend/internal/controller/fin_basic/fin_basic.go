package fin_basic

import (
	"context"

	finBasicApi "engineServ/api/fin_basic"
	v1 "engineServ/api/fin_basic/v1"
	"engineServ/internal/consts"
	finBasicLogic "engineServ/internal/logic/api/fin_basic"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// ControllerV1 财务基础资料控制器。
type ControllerV1 struct{}

// NewV1 创建财务基础资料控制器。
func NewV1() finBasicApi.IFinBasicV1 { return &ControllerV1{} }

// requireAdmin 写操作仅限公司超级管理员。
func requireAdmin(ctx context.Context) error {
	if gconv.Uint(ctx.Value(consts.CtxRoleKey)) != 1 {
		return gerror.New("仅公司超级管理员可维护财务基础资料")
	}
	return nil
}

// Accounts 会计科目列表。
func (c *ControllerV1) Accounts(ctx context.Context, req *v1.AccountListReq) (*v1.AccountListRes, error) {
	items, summary, err := finBasicLogic.New().Accounts(ctx, &finBasicLogic.AccountFilter{
		Keyword:      req.Keyword,
		CategoryCode: req.CategoryCode,
		Enabled:      req.Enabled,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*v1.AccountItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.AccountItem{
			Id:            item.Id,
			Code:          item.Code,
			Name:          item.Name,
			CategoryCode:  item.CategoryCode,
			DirectionCode: item.DirectionCode,
			Enabled:       item.Enabled,
			IsSystem:      item.IsSystem,
			SortOrder:     item.SortOrder,
			Used:          item.Used,
		})
	}

	return &v1.AccountListRes{
		Records: records,
		Summary: v1.AccountSummary{
			Total:            summary.Total,
			Enabled:          summary.Enabled,
			Disabled:         summary.Disabled,
			System:           summary.System,
			FundAccountCount: summary.FundAccountCount,
		},
	}, nil
}

// SaveAccount 新增或修改会计科目。
func (c *ControllerV1) SaveAccount(ctx context.Context, req *v1.AccountSaveReq) (*v1.AccountSaveRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	id, err := finBasicLogic.New().SaveAccount(ctx, &finBasicLogic.AccountForm{
		Id:           req.Id,
		Code:         req.Code,
		Name:         req.Name,
		CategoryCode: req.CategoryCode,
		Enabled:      req.Enabled,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AccountSaveRes{Id: id}, nil
}

// DeleteAccount 删除会计科目。
func (c *ControllerV1) DeleteAccount(ctx context.Context, req *v1.AccountDeleteReq) (*v1.AccountDeleteRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if err := finBasicLogic.New().DeleteAccount(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.AccountDeleteRes{}, nil
}

// ToggleAccount 启用或停用会计科目。
func (c *ControllerV1) ToggleAccount(ctx context.Context, req *v1.AccountToggleReq) (*v1.AccountToggleRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if err := finBasicLogic.New().ToggleAccount(ctx, req.Id, req.Enabled); err != nil {
		return nil, err
	}
	return &v1.AccountToggleRes{}, nil
}

// FundAccounts 资金账户列表。
func (c *ControllerV1) FundAccounts(ctx context.Context, req *v1.FundAccountListReq) (*v1.FundAccountListRes, error) {
	items, err := finBasicLogic.New().FundAccounts(ctx, &finBasicLogic.FundAccountFilter{
		Keyword: req.Keyword,
		Type:    req.Type,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*v1.FundAccountItem, 0, len(items))
	for _, item := range items {
		records = append(records, &v1.FundAccountItem{
			Id:                item.Id,
			Type:              item.Type,
			Name:              item.Name,
			BankName:          item.BankName,
			AccountLast4:      item.AccountLast4,
			LedgerAccountId:   item.LedgerAccountId,
			LedgerAccountCode: item.LedgerAccountCode,
			LedgerAccountName: item.LedgerAccountName,
			Enabled:           item.Enabled,
		})
	}
	return &v1.FundAccountListRes{Records: records}, nil
}

// SaveFundAccount 新增或修改资金账户。
func (c *ControllerV1) SaveFundAccount(ctx context.Context, req *v1.FundAccountSaveReq) (*v1.FundAccountSaveRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	id, err := finBasicLogic.New().SaveFundAccount(ctx, &finBasicLogic.FundAccountForm{
		Id:                req.Id,
		Type:              req.Type,
		Name:              req.Name,
		BankName:          req.BankName,
		AccountLast4:      req.AccountLast4,
		LedgerAccountCode: req.LedgerAccountCode,
		Enabled:           req.Enabled,
	})
	if err != nil {
		return nil, err
	}
	return &v1.FundAccountSaveRes{Id: id}, nil
}

// DeleteFundAccount 删除资金账户。
func (c *ControllerV1) DeleteFundAccount(ctx context.Context, req *v1.FundAccountDeleteReq) (*v1.FundAccountDeleteRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if err := finBasicLogic.New().DeleteFundAccount(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.FundAccountDeleteRes{}, nil
}

// ToggleFundAccount 启用或停用资金账户。
func (c *ControllerV1) ToggleFundAccount(ctx context.Context, req *v1.FundAccountToggleReq) (*v1.FundAccountToggleRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if err := finBasicLogic.New().ToggleFundAccount(ctx, req.Id, req.Enabled); err != nil {
		return nil, err
	}
	return &v1.FundAccountToggleRes{}, nil
}
