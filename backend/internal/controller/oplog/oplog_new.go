package oplog

import (
	oplogApi "engineServ/api/oplog"
	oplogLogic "engineServ/internal/logic/oplog"
)

type ControllerV1 struct {
	oplog *oplogLogic.Service
}

func NewV1() oplogApi.IOpLogV1 {
	return &ControllerV1{oplog: oplogLogic.New()}
}
