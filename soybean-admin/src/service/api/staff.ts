import { request } from '../request';

/**
 * 获取人员分页列表
 *
 * @param params 查询参数（分页 + 关键词 + 部门 + 状态）
 */
export function fetchGetStaffList(params?: Api.Staff.SearchParams) {
  return request<Api.Staff.List>({
    url: '/staff/list',
    method: 'get',
    params
  });
}

/**
 * 获取人员详情
 *
 * @param id 人员ID
 */
export function fetchGetStaffDetail(id: number) {
  return request<Api.Staff.Item>({
    url: '/staff/detail',
    method: 'get',
    params: { id }
  });
}

/**
 * 新增人员
 *
 * @param data 人员信息
 */
export function fetchCreateStaff(data: Api.Staff.FormData) {
  return request<Api.Staff.Item>({
    url: '/staff/create',
    method: 'post',
    data
  });
}

/**
 * 更新人员
 *
 * @param data 人员信息（含 id）
 */
export function fetchUpdateStaff(data: Api.Staff.FormData & { id: number }) {
  return request<Api.Staff.Item>({
    url: '/staff/update',
    method: 'post',
    data
  });
}

/**
 * 删除人员
 *
 * @param id 人员ID
 */
export function fetchDeleteStaff(id: number) {
  return request({
    url: '/staff/delete',
    method: 'post',
    data: { id }
  });
}

/**
 * 批量删除人员
 *
 * @param ids 人员ID集合
 */
export function fetchBatchDeleteStaff(ids: number[]) {
  return request({
    url: '/staff/batchDelete',
    method: 'post',
    data: { ids }
  });
}

/**
 * 获取部门下拉选项（新增/编辑人员时选择部门用）
 */
export function fetchGetDepartmentOptions() {
  return request<{ items: Api.Staff.DepartmentOption[] }>({
    url: '/department/options',
    method: 'get'
  });
}