package driver_finance

import (
	"context"

	v1 "engineServ/api/driver_finance/v1"
)

type IDriverFinanceV1 interface {
	Drivers(ctx context.Context, req *v1.DriversReq) (res *v1.DriversRes, err error)
	Settlements(ctx context.Context, req *v1.SettlementsReq) (res *v1.SettlementsRes, err error)
	CreateSettlement(ctx context.Context, req *v1.CreateSettlementReq) (res *v1.CreateSettlementRes, err error)
}
