import { request } from '../request';

export function fetchGetDriverBalances() {
  return request<Api.DriverFinance.DriverList>({ url: '/driver-finance/drivers', method: 'get' });
}

export function fetchGetDriverSettlements(params: Api.DriverFinance.SettlementSearchParams) {
  return request<Api.DriverFinance.SettlementList>({ url: '/driver-finance/settlements', method: 'get', params });
}

export function fetchCreateDriverSettlement(data: Api.DriverFinance.SettlementForm) {
  return request<Api.DriverFinance.SettlementResult>({ url: '/driver-finance/settlement/create', method: 'post', data });
}