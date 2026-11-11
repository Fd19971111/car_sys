import { request } from '../request';

/**
 * 获取部门分页列表
 *
 * @param params 查询参数（分页 + 关键词 + 状态）
 */
export function fetchGetDepartmentList(params?: Api.Department.SearchParams) {
  return request<Api.Department.List>({
    url: '/department/list',
    method: 'get',
    params
  });
}

/**
 * 获取部门详情
 *
 * @param id 部门ID
 */
export function fetchGetDepartmentDetail(id: number) {
  return request<Api.Department.Item>({
    url: '/department/detail',
    method: 'get',
    params: { id }
  });
}

/**
 * 新增部门
 *
 * @param data 部门信息
 */
export function fetchCreateDepartment(data: Api.Department.FormData) {
  return request<Api.Department.Item>({
    url: '/department/create',
    method: 'post',
    data
  });
}

/**
 * 更新部门
 *
 * @param data 部门信息（含 id）
 */
export function fetchUpdateDepartment(data: Api.Department.FormData & { id: number }) {
  return request<Api.Department.Item>({
    url: '/department/update',
    method: 'post',
    data
  });
}

/**
 * 删除部门
 *
 * @param id 部门ID
 */
export function fetchDeleteDepartment(id: number) {
  return request({
    url: '/department/delete',
    method: 'post',
    data: { id }
  });
}

/**
 * 批量删除部门
 *
 * @param ids 部门ID集合
 */
export function fetchBatchDeleteDepartment(ids: number[]) {
  return request({
    url: '/department/batchDelete',
    method: 'post',
    data: { ids }
  });
}