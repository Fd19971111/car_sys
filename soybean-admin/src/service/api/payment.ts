import { request } from '../request';

/**
 * 获取收款记录分页列表
 *
 * @param params 查询参数（分页 + 关键词 + 所属项目 + 款项类型）
 */
export function fetchGetPaymentList(params?: Api.Payment.SearchParams) {
  return request<Api.Payment.List>({
    url: '/payment/list',
    method: 'get',
    params
  });
}

/** 获取收款概览统计（笔数 / 收款总额 / 本月收款 / 涉及项目数） */
export function fetchGetPaymentStats() {
  return request<Api.Payment.Stats>({
    url: '/payment/stats',
    method: 'get'
  });
}

/**
 * 获取收款记录详情（含所属项目合同金额与累计已收款）
 *
 * @param id 收款记录ID
 */
export function fetchGetPaymentDetail(id: number) {
  return request<Api.Payment.Detail>({
    url: '/payment/detail',
    method: 'get',
    params: { id }
  });
}

/**
 * 新增收款记录
 *
 * @param data 收款记录信息
 */
export function fetchCreatePayment(data: Api.Payment.FormData) {
  return request<{ id: number; voucherNo: string }>({
    url: '/payment/create',
    method: 'post',
    data
  });
}

/**
 * 更新收款记录
 *
 * @param data 收款记录信息（含 id）
 */
export function fetchUpdatePayment(data: Api.Payment.FormData & { id: number }) {
  return request<{ id: number; voucherNo: string }>({
    url: '/payment/update',
    method: 'post',
    data
  });
}

/**
 * 删除收款记录（后端为软删除）
 *
 * @param id 收款记录ID
 */
export function fetchDeletePayment(id: number) {
  return request({
    url: '/payment/delete',
    method: 'post',
    data: { id }
  });
}

/**
 * 批量删除收款记录
 *
 * @param ids 收款记录ID集合
 */
export function fetchBatchDeletePayment(ids: number[]) {
  return request({
    url: '/payment/batchDelete',
    method: 'post',
    data: { ids }
  });
}
