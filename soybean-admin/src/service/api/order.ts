import { request } from '../request';

/** 获取订单分页列表 */
export function fetchGetOrderList(params: Api.Order.SearchParams) {
  return request<Api.Order.List>({ url: '/order/list', method: 'get', params });
}

/** 获取订单详情 */
export function fetchGetOrderDetail(id: number) {
  return request<Api.Order.Item>({ url: '/order/detail', method: 'get', params: { id } });
}

/** 新增订单 */
export function fetchCreateOrder(data: Api.Order.FormData) {
  return request<{ id: number; orderNumber: string }>({ url: '/order/create', method: 'post', data });
}

/** 更新订单 */
export function fetchUpdateOrder(data: Api.Order.FormData & { id: number }) {
  return request({ url: '/order/update', method: 'post', data });
}

/** 删除订单 */
export function fetchDeleteOrder(id: number) {
  return request({ url: '/order/delete', method: 'post', data: { id } });
}

/** 批量删除订单 */
export function fetchBatchDeleteOrder(ids: number[]) {
  return request({ url: '/order/batchDelete', method: 'post', data: { ids } });
}

/** 订单统计 */
export function fetchGetOrderStats() {
  return request<Api.Order.Stats>({ url: '/order/stats', method: 'get' });
}

/** 行程监控 */
export function fetchGetOrderMonitor() {
  return request<Api.Order.Monitor>({ url: '/order/monitor', method: 'get' });
}
