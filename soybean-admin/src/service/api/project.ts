import { request } from '../request';

/**
 * 获取项目分页列表
 *
 * @param params 查询参数（分页 + 关键词 + 状态 + 所属客户）
 */
export function fetchGetProjectList(params?: Api.Project.SearchParams) {
  return request<Api.Project.List>({
    url: '/project/list',
    method: 'get',
    params
  });
}

/** 获取项目概览统计（总数 / 待开工 / 施工中 / 已完工 / 已结算 / 合同总额） */
export function fetchGetProjectStats() {
  return request<Api.Project.Stats>({
    url: '/project/stats',
    method: 'get'
  });
}

/** 获取项目表单下拉选项（客户 / 负责人） */
export function fetchGetProjectOptions() {
  return request<Api.Project.Options>({
    url: '/project/options',
    method: 'get'
  });
}

/** 获取项目下拉选项（供报价 / 收款等模块选择项目，label 为「项目名称 · 客户姓名」） */
export function fetchGetProjectSelectOptions() {
  return request<{ projects: Api.Project.SelectOption[] }>({
    url: '/project/selectOptions',
    method: 'get'
  });
}

/**
 * 获取项目详情（含客户信息与施工进度时间轴）
 *
 * @param id 项目ID
 */
export function fetchGetProjectDetail(id: number) {
  return request<Api.Project.Detail>({
    url: '/project/detail',
    method: 'get',
    params: { id }
  });
}

/**
 * 新增项目
 *
 * @param data 项目信息
 */
export function fetchCreateProject(data: Api.Project.FormData) {
  return request<Api.Project.Item>({
    url: '/project/create',
    method: 'post',
    data
  });
}

/**
 * 更新项目
 *
 * @param data 项目信息（含 id）
 */
export function fetchUpdateProject(data: Api.Project.FormData & { id: number }) {
  return request<Api.Project.Item>({
    url: '/project/update',
    method: 'post',
    data
  });
}

/**
 * 删除项目（后端为软删除）
 *
 * @param id 项目ID
 */
export function fetchDeleteProject(id: number) {
  return request({
    url: '/project/delete',
    method: 'post',
    data: { id }
  });
}

/**
 * 批量删除项目
 *
 * @param ids 项目ID集合
 */
export function fetchBatchDeleteProject(ids: number[]) {
  return request({
    url: '/project/batchDelete',
    method: 'post',
    data: { ids }
  });
}
