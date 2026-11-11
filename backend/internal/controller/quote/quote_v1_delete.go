package quote

import (
	"context"

	v1 "engineServ/api/quote/v1"
)

func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	if err = c.quote.Delete(ctx, req.Id); err != nil {
		return nil, err
	}

	return &v1.DeleteRes{}, nil
}
