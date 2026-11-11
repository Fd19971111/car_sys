import { request } from '../request';

/**
 * 获取付款记录分页列表
 *
 * @param params 查询参数（分页 + 关键词 + 所属项目 + 费用类型）
 */
export function fetchGetPayoutList(params?: Api.Payout.SearchParams) {
  return request<Api.Payout.List>({
    url: '/payout/list',
    method: 'get',
    params
  });
}

/** 获取付款概览统计（笔数 / 付款总额 / 本月付款 / 涉及项目数） */
export function fetchGetPayoutStats() {
  return request<Api.Payout.Stats>({
    url: '/payout/stats',
    method: 'get'
  });
}

/**
 * 获取付款记录详情（含所属项目合同金额、累计已付款与记账分录）
 *
 * @param id 付款记录ID
 */
export function fetchGetPayoutDetail(id: number) {
  return request<Api.Payout.Detail>({
    url: '/payout/detail',
    method: 'get',
    params: { id }
  });
}

/**
 * 新增付款记录（后端会自动生成付款凭证）
 *
 * @param data 付款记录信息
 */
export function fetchCreatePayout(data: Api.Payout.FormData) {
  return request<{ id: number; voucherNo: string }>({
    url: '/payout/create',
    method: 'post',
    data
  });
}

/**
 * 更新付款记录（后端会作废旧凭证并重新记账）
 *
 * @param data 付款记录信息（含 id）
 */
export function fetchUpdatePayout(data: Api.Payout.FormData & { id: number }) {
  return request<{ id: number; voucherNo: string }>({
    url: '/payout/update',
    method: 'post',
    data
  });
}

/**
 * 删除付款记录（后端为软删除，同时作废凭证）
 *
 * @param id 付款记录ID
 */
export function fetchDeletePayout(id: number) {
  return request({
    url: '/payout/delete',
    method: 'post',
    data: { id }
  });
}

/**
 * 批量删除付款记录
 *
 * @param ids 付款记录ID集合
 */
export function fetchBatchDeletePayout(ids: number[]) {
  return request({
    url: '/payout/batchDelete',
    method: 'post',
    data: { ids }
  });
}
