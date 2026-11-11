package oplog

import (
	"context"

	v1 "engineServ/api/oplog/v1"
)

// IOpLogV1 操作日志接口（登录后可查看）。
type IOpLogV1 interface {
	List(ctx context.Context, req *v1.OpLogListReq) (res *v1.OpLogListRes, err error)
}
