package inventory

import (
	inventoryApi "engineServ/api/inventory"
	inventoryLogic "engineServ/internal/logic/api/inventory"
)

// ControllerV1 库存管理控制器。
type ControllerV1 struct {
	stock *inventoryLogic.Service
}

// NewV1 创建库存管理控制器。
func NewV1() inventoryApi.IInventoryV1 {
	return &ControllerV1{stock: inventoryLogic.New()}
}
