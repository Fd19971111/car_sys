import { request } from '../request';

export function fetchMenuList() {
  return request<Api.Menu.List>({ url: '/menu/list' });
}

export function fetchDepartments() {
  return request<Api.Menu.DepartmentList>({ url: '/menu/departments' });
}

/** 快捷新增部门（菜单模块，仅需名称与编码；完整 CRUD 走 department 模块的 fetchCreateDepartment） */
export function fetchQuickCreateDepartment(name: string, code: string) {
  return request<{ id: number }>({ url: '/menu/createDepartment', method: 'post', data: { name, code } });
}

export function fetchDepartmentUsers() {
  return request<Api.Menu.UserList>({ url: '/menu/users' });
}

export function fetchSetUserDepartment(userId: number, departmentId: number) {
  return request({ url: '/menu/setUserDepartment', method: 'post', data: { userId, departmentId } });
}

export function fetchDepartmentMenus(departmentId: number) {
  return request<Api.Menu.DepartmentMenuList>({ url: '/menu/departmentMenus', params: { departmentId } });
}

export function fetchSaveMenu(data: Api.Menu.Form) {
  return request<{ id: number }>({ url: '/menu/save', method: 'post', data });
}

export function fetchDeleteMenu(id: number) {
  return request({ url: '/menu/delete', method: 'post', data: { id } });
}

export function fetchSetDepartmentMenus(departmentId: number, menuIds: number[]) {
  return request({ url: '/menu/setDepartmentMenus', method: 'post', data: { departmentId, menuIds } });
}