import { request } from '../request';

/**
 * 获取会计科目列表与统计
 *
 * @param params 查询参数（关键词 / 类别 / 启用状态）
 */
export function fetchGetFinAccounts(params?: Api.FinBasic.AccountSearchParams) {
  return request<Api.FinBasic.AccountListResult>({
    url: '/fin-basic/accounts',
    method: 'get',
    params
  });
}

/** 新增或修改会计科目 */
export function fetchSaveFinAccount(data: Api.FinBasic.AccountFormData) {
  return request<{ id: number }>({
    url: '/fin-basic/account/save',
    method: 'post',
    data
  });
}

/** 删除会计科目（系统预置或已被引用的科目不可删除） */
export function fetchDeleteFinAccount(id: number) {
  return request({
    url: '/fin-basic/account/delete',
    method: 'post',
    data: { id }
  });
}

/** 启用或停用会计科目 */
export function fetchToggleFinAccount(id: number, enabled: boolean) {
  return request({
    url: '/fin-basic/account/toggle',
    method: 'post',
    data: { id, enabled }
  });
}

/**
 * 获取资金账户列表
 *
 * @param params 查询参数（关键词 / 账户类型）
 */
export function fetchGetFinFundAccounts(params?: Api.FinBasic.FundAccountSearchParams) {
  return request<Api.FinBasic.FundAccountListResult>({
    url: '/fin-basic/fund-accounts',
    method: 'get',
    params
  });
}

/** 新增或修改资金账户 */
export function fetchSaveFinFundAccount(data: Api.FinBasic.FundAccountFormData) {
  return request<{ id: number }>({
    url: '/fin-basic/fund-account/save',
    method: 'post',
    data
  });
}

/** 删除资金账户 */
export function fetchDeleteFinFundAccount(id: number) {
  return request({
    url: '/fin-basic/fund-account/delete',
    method: 'post',
    data: { id }
  });
}

/** 启用或停用资金账户 */
export function fetchToggleFinFundAccount(id: number, enabled: boolean) {
  return request({
    url: '/fin-basic/fund-account/toggle',
    method: 'post',
    data: { id, enabled }
  });
}
