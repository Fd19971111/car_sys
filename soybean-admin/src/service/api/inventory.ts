import { request } from '../request';

// ===================== 库存管理 =====================

/** 获取公司级材料档案分页列表 */
export function fetchGetMaterialCatalogList(params: Api.Inventory.MaterialCatalogSearchParams) {
  return request<Api.Inventory.MaterialCatalogList>({
    url: '/inventory/materialCatalog/list',
    method: 'get',
    params
  });
}

/** 新增公司级材料档案 */
export function fetchCreateMaterialCatalog(data: Api.Inventory.MaterialCatalogForm) {
  return request<{ id: number }>({
    url: '/inventory/materialCatalog/create',
    method: 'post',
    data
  });
}

/** 更新公司级材料档案 */
export function fetchUpdateMaterialCatalog(data: Api.Inventory.MaterialCatalogForm & { id: number }) {
  return request<{ id: number }>({
    url: '/inventory/materialCatalog/update',
    method: 'post',
    data
  });
}

/** 删除公司级材料档案 */
export function fetchDeleteMaterialCatalog(id: number) {
  return request({
    url: '/inventory/materialCatalog/delete',
    method: 'post',
    data: { id }
  });
}

/** 获取库存分页列表（跨项目，projectId 为 0 时表示全部项目） */
export function fetchGetStockList(params: Api.Inventory.StockSearchParams) {
  return request<Api.Inventory.StockList>({
    url: '/inventory/stock/list',
    method: 'get',
    params
  });
}

/** 新增库存 */
export function fetchCreateStock(data: Api.Inventory.StockForm) {
  return request<{ id: number }>({
    url: '/inventory/stock/create',
    method: 'post',
    data
  });
}

/** 更新库存 */
export function fetchUpdateStock(data: Api.Inventory.StockForm & { id: number }) {
  return request<{ id: number }>({
    url: '/inventory/stock/update',
    method: 'post',
    data
  });
}

/** 删除库存（软删除，领用记录保留） */
export function fetchDeleteStock(id: number) {
  return request({
    url: '/inventory/stock/delete',
    method: 'post',
    data: { id }
  });
}

/** 材料入库：库存增加，可选同时生成一张付款单 */
export function fetchStockIn(data: Api.Inventory.StockInForm) {
  return request<{ id: number; payoutId: number; voucherNo: string }>({
    url: '/inventory/stock/in',
    method: 'post',
    data
  });
}

/** 材料出库（项目领用）：写入领用记录并扣减结存，不生成任何财务单据 */
export function fetchStockOut(data: Api.Inventory.StockOutForm) {
  return request<{ id: number }>({
    url: '/inventory/stock/out',
    method: 'post',
    data
  });
}

/** 库存补录：直接调整库存，不生成财务单据 */
export function fetchStockAdjust(data: Api.Inventory.StockAdjustForm) {
  return request<{ id: number }>({
    url: '/inventory/stock/adjust',
    method: 'post',
    data
  });
}

/** 出入库流水分页列表 */
export function fetchGetStockRecordList(params: Api.Inventory.StockRecordSearchParams) {
  return request<Api.Inventory.StockRecordList>({
    url: '/inventory/stock/records',
    method: 'get',
    params
  });
}

// ===================== 仓库管理 =====================

/** 获取仓库分页列表 */
export function fetchGetWarehouseList(params?: Api.Inventory.WarehouseSearchParams) {
  return request<Api.Inventory.WarehouseList>({
    url: '/inventory/warehouse/list',
    method: 'get',
    params
  });
}

/** 新增仓库 */
export function fetchCreateWarehouse(data: Api.Inventory.WarehouseForm) {
  return request<{ id: number }>({
    url: '/inventory/warehouse/create',
    method: 'post',
    data
  });
}

/** 更新仓库 */
export function fetchUpdateWarehouse(data: Api.Inventory.WarehouseForm & { id: number }) {
  return request<{ id: number }>({
    url: '/inventory/warehouse/update',
    method: 'post',
    data
  });
}

/** 删除仓库（软删除，默认仓库不可删除） */
export function fetchDeleteWarehouse(id: number) {
  return request({
    url: '/inventory/warehouse/delete',
    method: 'post',
    data: { id }
  });
}

/** 批量删除仓库（软删除，包含默认仓库时整体拒绝） */
export function fetchBatchDeleteWarehouse(ids: number[]) {
  return request({
    url: '/inventory/warehouse/batchDelete',
    method: 'post',
    data: { ids }
  });
}
