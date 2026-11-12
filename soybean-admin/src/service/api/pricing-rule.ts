import { request } from '../request';

export function fetchGetPricingRuleList(params?: Api.PricingRule.SearchParams) {
  return request<Api.PricingRule.List>({ url: '/pricing-rule/list', method: 'get', params });
}

export function fetchCreatePricingRule(data: Api.PricingRule.FormData) {
  return request<{ id: number }>({ url: '/pricing-rule/create', method: 'post', data });
}

export function fetchUpdatePricingRule(data: Api.PricingRule.FormData & { id: number }) {
  return request<{ id: number }>({ url: '/pricing-rule/update', method: 'post', data });
}

export function fetchDeletePricingRule(id: number) {
  return request({ url: '/pricing-rule/delete', method: 'post', data: { id } });
}

export function fetchBatchDeletePricingRule(ids: number[]) {
  return request({ url: '/pricing-rule/batchDelete', method: 'post', data: { ids } });
}
