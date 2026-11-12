import { request } from '../request';

/** 获取财务报表总览（收支、结余、资金余额、构成、回款、科目余额） */
export function fetchGetFinReportOverview(params?: Api.FinReport.OverviewParams) {
  return request<Api.FinReport.Overview>({ url: '/fin-report/overview', method: 'get', params });
}

/** 获取资金流水明细 */
export function fetchGetFinReportFlows(params?: Api.FinReport.FlowSearchParams) {
  return request<Api.FinReport.FlowListResult>({ url: '/fin-report/flows', method: 'get', params });
}
