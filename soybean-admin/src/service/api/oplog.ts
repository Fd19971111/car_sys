import { request } from '../request';

/**
 * 获取操作日志列表（收付款等关键业务动作的审计记录）
 *
 * @param params 查询参数（模块 / 动作 / 关键词 / 日期区间 / 分页）
 */
export function fetchGetOperationLogList(params?: Api.OpLog.SearchParams) {
  return request<Api.OpLog.List>({
    url: '/operation-log/list',
    method: 'get',
    params
  });
}
