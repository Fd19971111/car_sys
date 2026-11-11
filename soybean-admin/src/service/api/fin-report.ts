import { request } from '../request';

/**
 * 获取财务报表总览（收支、结余、资金余额、构成、回款、科目余额）
 *
 * @param params 查询参数（统计年份）
 */
export function fetchGetFinReportOverview(params?: Api.FinReport.OverviewParams) {
  return request<Api.FinReport.Overview>({
    url: '/fin-report/overview',
    method: 'get',
    params
  });
}

/**
 * 获取资金流水明细
 *
 * @param params 查询参数（年份 / 月份 / 方向 / 账户 / 关键词 / 分页）
 */
export function fetchGetFinReportFlows(params?: Api.FinReport.FlowSearchParams) {
  return request<Api.FinReport.FlowListResult>({
    url: '/fin-report/flows',
    method: 'get',
    params
  });
}

/**
 * 获取项目成本 / 毛利列表
 *
 * @param params 查询参数（关键字 / 项目状态 / 只看亏损）
 */
export function fetchGetProjectProfitList(params?: Api.FinReport.ProjectProfitSearchParams) {
  return request<Api.FinReport.ProjectProfitList>({
    url: '/fin-report/project-profit/list',
    method: 'get',
    params
  });
}

/**
 * 获取单个项目的成本归集与毛利明细
 *
 * @param projectId 项目ID
 */
export function fetchGetProjectProfitDetail(projectId: number) {
  return request<Api.FinReport.ProjectProfitDetail>({
    url: '/fin-report/project-profit/detail',
    method: 'get',
    params: { projectId }
  });
}
