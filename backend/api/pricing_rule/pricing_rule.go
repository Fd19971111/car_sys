package pricing_rule

import (
	"context"

	v1 "engineServ/api/pricing_rule/v1"
)

type IPricingRuleV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.SaveRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.SaveRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
	BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (res *v1.BatchDeleteRes, err error)
}
