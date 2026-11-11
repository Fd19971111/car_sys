import { request } from '../request';

export function fetchGetUserList(params?: Api.User.SearchParams) {
  return request<Api.User.List>({ url: '/user/list', method: 'get', params });
}

export function fetchResetUserPwd(id: number) {
  return request<{ password: string }>({ url: '/user/resetPwd', method: 'post', data: { id } });
}

export function fetchToggleUser(id: number, status: number) {
  return request({ url: '/user/toggle', method: 'post', data: { id, status } });
}

export function fetchCreateUserAccount(data: Api.User.CreateAccountForm) {
  return request<{ id: number }>({ url: '/user/create', method: 'post', data });
}

export function fetchGetAvailableStaff() {
  return request<{ items: Api.User.StaffOption[] }>({ url: '/user/availableStaff', method: 'get' });
}

export function fetchUpdateUserAccount(data: Api.User.UpdateAccountForm) {
  return request({ url: '/user/update', method: 'post', data });
}

export function fetchChangeUserPwd(id: number, password: string) {
  return request({ url: '/user/changePwd', method: 'post', data: { id, password } });
}