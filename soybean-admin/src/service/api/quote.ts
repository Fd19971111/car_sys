import { request } from '../request';

/**
 * 获取报价明细分页列表
 *
 * @param params 查询参数（分页 + 关键词 + 所属项目）
 */
export function fetchGetQuoteList(params?: Api.Quote.SearchParams) {
  return request<Api.Quote.List>({
    url: '/quote/list',
    method: 'get',
    params
  });
}

/** 获取报价概览统计（条目数 / 报价总额 / 涉及项目数 / 平均单条金额） */
export function fetchGetQuoteStats() {
  return request<Api.Quote.Stats>({
    url: '/quote/stats',
    method: 'get'
  });
}

/**
 * 获取报价明细详情（含所属项目合同金额与报价合计）
 *
 * @param id 报价明细ID
 */
export function fetchGetQuoteDetail(id: number) {
  return request<Api.Quote.Detail>({
    url: '/quote/detail',
    method: 'get',
    params: { id }
  });
}

/**
 * 新增报价明细
 *
 * @param data 报价明细信息
 */
export function fetchCreateQuote(data: Api.Quote.FormData) {
  return request<{ id: number; warning?: string }>({
    url: '/quote/create',
    method: 'post',
    data
  });
}

/**
 * 更新报价明细
 *
 * @param data 报价明细信息（含 id）
 */
export function fetchUpdateQuote(data: Api.Quote.FormData & { id: number }) {
  return request<{ id: number; warning?: string }>({
    url: '/quote/update',
    method: 'post',
    data
  });
}

/**
 * 删除报价明细（后端为软删除）
 *
 * @param id 报价明细ID
 */
export function fetchDeleteQuote(id: number) {
  return request({
    url: '/quote/delete',
    method: 'post',
    data: { id }
  });
}

/**
 * 批量删除报价明细
 *
 * @param ids 报价明细ID集合
 */
export function fetchBatchDeleteQuote(ids: number[]) {
  return request({
    url: '/quote/batchDelete',
    method: 'post',
    data: { ids }
  });
}
