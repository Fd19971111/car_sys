package quote

import (
	"context"

	v1 "engineServ/api/quote/v1"
)

func (c *ControllerV1) BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (res *v1.BatchDeleteRes, err error) {
	if err = c.quote.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}

	return &v1.BatchDeleteRes{}, nil
}
