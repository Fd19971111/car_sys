package project_data

import (
	projectDataApi "engineServ/api/project_data"
	projectDataLogic "engineServ/internal/logic/api/project_data"
)

// ControllerV1 项目资料控制器。
type ControllerV1 struct {
	data *projectDataLogic.Service
}

// NewV1 创建项目资料控制器。
func NewV1() projectDataApi.IProjectDataV1 {
	return &ControllerV1{data: projectDataLogic.New()}
}
