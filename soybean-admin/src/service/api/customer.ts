import { request } from '../request';

/**
 * 获取客户分页列表
 *
 * @param params 查询参数（分页 + 关键词 + 状态 + 来源）
 */
export function fetchGetCustomerList(params?: Api.Customer.SearchParams) {
  return request<Api.Customer.List>({
    url: '/customer/list',
    method: 'get',
    params
  });
}

/** 获取客户概览统计（总数 / 跟进中 / 已成交 / 已流失 / 成交金额） */
export function fetchGetCustomerStats() {
  return request<Api.Customer.Stats>({
    url: '/customer/stats',
    method: 'get'
  });
}

/**
 * 获取客户详情（含关联装修项目、负责人、最近一条施工进度）
 *
 * @param id 客户ID
 */
export function fetchGetCustomerDetail(id: number) {
  return request<Api.Customer.Item>({
    url: '/customer/detail',
    method: 'get',
    params: { id }
  });
}

/**
 * 新增客户
 *
 * @param data 客户信息
 */
export function fetchCreateCustomer(data: Api.Customer.FormData) {
  return request<Api.Customer.Item>({
    url: '/customer/create',
    method: 'post',
    data
  });
}

/**
 * 更新客户
 *
 * @param data 客户信息（含 id）
 */
export function fetchUpdateCustomer(data: Api.Customer.FormData & { id: number }) {
  return request<Api.Customer.Item>({
    url: '/customer/update',
    method: 'post',
    data
  });
}

/**
 * 删除客户（后端为软删除）
 *
 * @param id 客户ID
 */
export function fetchDeleteCustomer(id: number) {
  return request({
    url: '/customer/delete',
    method: 'post',
    data: { id }
  });
}

/**
 * 批量删除客户
 *
 * @param ids 客户ID集合
 */
export function fetchBatchDeleteCustomer(ids: number[]) {
  return request({
    url: '/customer/batchDelete',
    method: 'post',
    data: { ids }
  });
}
