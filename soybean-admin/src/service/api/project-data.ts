import { request } from '../request';

// ===================== 材料清单 =====================

/** 获取材料清单分页列表 */
export function fetchGetMaterialList(params: Api.ProjectData.MaterialSearchParams) {
  return request<Api.ProjectData.MaterialList>({
    url: '/project/material/list',
    method: 'get',
    params
  });
}

/** 新增材料 */
export function fetchCreateMaterial(data: Api.ProjectData.MaterialForm) {
  return request<{ id: number }>({
    url: '/project/material/create',
    method: 'post',
    data
  });
}

/** 更新材料 */
export function fetchUpdateMaterial(data: Api.ProjectData.MaterialForm & { id: number }) {
  return request<{ id: number }>({
    url: '/project/material/update',
    method: 'post',
    data
  });
}

/** 删除材料（软删除） */
export function fetchDeleteMaterial(id: number) {
  return request({
    url: '/project/material/delete',
    method: 'post',
    data: { id }
  });
}

// ===================== 材料领用 =====================

/** 获取材料领用记录 */
export function fetchGetMaterialUsageList(params: Api.ProjectData.UsageSearchParams) {
  return request<Api.ProjectData.UsageList>({
    url: '/project/material/usageList',
    method: 'get',
    params
  });
}

/** 新增材料领用记录 */
export function fetchCreateMaterialUsage(data: Api.ProjectData.UsageForm) {
  return request<{ id: number }>({
    url: '/project/material/usageCreate',
    method: 'post',
    data
  });
}

/** 删除材料领用记录 */
export function fetchDeleteMaterialUsage(id: number) {
  return request({
    url: '/project/material/usageDelete',
    method: 'post',
    data: { id }
  });
}

/** 领用单出库：选择仓库后直接完成出库 */
export function fetchStockOutMaterialUsage(data: { id: number; warehouseId: number; remark?: string }) {
  return request<{ id: number; recordId: number }>({
    url: '/project/material/usageStock',
    method: 'post',
    data
  });
}

/** 获取项目采购单 */
export function fetchGetPurchaseOrderList(params: { projectId: number; current?: number; size?: number }) {
  return request<Api.ProjectData.PurchaseOrderList>({
    url: '/project/purchaseOrder/list',
    method: 'get',
    params
  });
}

/** 获取付款管理可关联的待付款采购单 */
export function fetchGetPendingPurchaseOrders() {
  return request<{ records: Api.ProjectData.PurchaseOrder[] }>({
    url: '/project/purchaseOrder/pendingList',
    method: 'get'
  });
}

/** 创建材料采购单 */
export function fetchCreatePurchaseOrder(data: Api.ProjectData.PurchaseOrderForm) {
  return request<{ id: number }>({
    url: '/project/purchaseOrder/create',
    method: 'post',
    data
  });
}

/** 采购单入库：选择仓库后直接完成出入库 */
export function fetchStockInPurchaseOrder(data: { id: number; warehouseId: number; remark?: string }) {
  return request<{ id: number; recordId: number }>({
    url: '/project/purchaseOrder/stock',
    method: 'post',
    data
  });
}

// ===================== 施工日志 / 现场照片 =====================

/** 获取施工日志列表 */
export function fetchGetSiteLogList(params: Api.ProjectData.SiteLogSearchParams) {
  return request<Api.ProjectData.SiteLogList>({
    url: '/project/log/list',
    method: 'get',
    params
  });
}

/** 新增施工日志 */
export function fetchCreateSiteLog(data: Api.ProjectData.SiteLogForm) {
  return request<{ id: number }>({
    url: '/project/log/create',
    method: 'post',
    data
  });
}

/** 更新施工日志 */
export function fetchUpdateSiteLog(data: Api.ProjectData.SiteLogForm & { id: number }) {
  return request<{ id: number }>({
    url: '/project/log/update',
    method: 'post',
    data
  });
}

/** 删除施工日志 */
export function fetchDeleteSiteLog(id: number) {
  return request({
    url: '/project/log/delete',
    method: 'post',
    data: { id }
  });
}

// ===================== 变更记录 =====================

/** 获取变更记录列表（含增项 / 减项合计） */
export function fetchGetChangeList(params: Api.ProjectData.ChangeSearchParams) {
  return request<Api.ProjectData.ChangeList>({
    url: '/project/change/list',
    method: 'get',
    params
  });
}

/** 新增变更记录 */
export function fetchCreateChange(data: Api.ProjectData.ChangeForm) {
  return request<{ id: number }>({
    url: '/project/change/create',
    method: 'post',
    data
  });
}

/** 更新变更记录 */
export function fetchUpdateChange(data: Api.ProjectData.ChangeForm & { id: number }) {
  return request<{ id: number }>({
    url: '/project/change/update',
    method: 'post',
    data
  });
}

/** 删除变更记录 */
export function fetchDeleteChange(id: number) {
  return request({
    url: '/project/change/delete',
    method: 'post',
    data: { id }
  });
}

// ===================== 项目合同 =====================

export function fetchGetContractList(params: Api.ProjectData.ContractListParams) {
  return request<Api.ProjectData.ContractList>({ url: '/project/contract/list', method: 'get', params });
}

export function fetchCreateContract(data: Api.ProjectData.ContractForm) {
  return request<{ id: number }>({ url: '/project/contract/create', method: 'post', data });
}

export function fetchUpdateContract(data: Api.ProjectData.ContractForm & { id: number }) {
  return request<{ id: number }>({ url: '/project/contract/update', method: 'post', data });
}

export function fetchDeleteContract(id: number) {
  return request({ url: '/project/contract/delete', method: 'post', data: { id } });
}
