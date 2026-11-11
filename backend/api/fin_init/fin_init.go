package fin_init

import (
	"context"

	v1 "engineServ/api/fin_init/v1"
)

type IFinInitV1 interface {
	Save(ctx context.Context, req *v1.SaveReq) (res *v1.SaveRes, err error)
}
