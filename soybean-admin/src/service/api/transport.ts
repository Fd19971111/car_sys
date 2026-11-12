import { request } from '../request';

// ===================== 司机管理 =====================

/** 获取司机分页列表 */
export function fetchGetDriverList(params: Api.Transport.DriverSearchParams) {
  return request<Api.Transport.DriverList>({ url: '/transport/driver/list', method: 'get', params });
}

/** 获取司机详情 */
export function fetchGetDriverDetail(id: number) {
  return request<Api.Transport.Driver>({ url: '/transport/driver/detail', method: 'get', params: { id } });
}

/** 新增司机 */
export function fetchCreateDriver(data: Api.Transport.DriverFormData) {
  return request<{ id: number }>({ url: '/transport/driver/create', method: 'post', data });
}

/** 更新司机 */
export function fetchUpdateDriver(data: Api.Transport.DriverFormData & { id: number }) {
  return request({ url: '/transport/driver/update', method: 'post', data });
}

/** 删除司机 */
export function fetchDeleteDriver(id: number) {
  return request({ url: '/transport/driver/delete', method: 'post', data: { id } });
}

/** 批量删除司机 */
export function fetchBatchDeleteDriver(ids: number[]) {
  return request({ url: '/transport/driver/batchDelete', method: 'post', data: { ids } });
}

/** 获取司机下拉选项 */
export function fetchGetDriverOptions() {
  return request<{ items: Api.Transport.DriverOption[] }>({ url: '/transport/driver/options', method: 'get' });
}

// ===================== 车辆管理 =====================

/** 获取车辆分页列表 */
export function fetchGetVehicleList(params: Api.Transport.VehicleSearchParams) {
  return request<Api.Transport.VehicleList>({ url: '/transport/vehicle/list', method: 'get', params });
}

/** 获取车辆详情 */
export function fetchGetVehicleDetail(id: number) {
  return request<Api.Transport.Vehicle>({ url: '/transport/vehicle/detail', method: 'get', params: { id } });
}

/** 新增车辆 */
export function fetchCreateVehicle(data: Api.Transport.VehicleFormData) {
  return request<{ id: number }>({ url: '/transport/vehicle/create', method: 'post', data });
}

/** 更新车辆 */
export function fetchUpdateVehicle(data: Api.Transport.VehicleFormData & { id: number }) {
  return request({ url: '/transport/vehicle/update', method: 'post', data });
}

/** 删除车辆 */
export function fetchDeleteVehicle(id: number) {
  return request({ url: '/transport/vehicle/delete', method: 'post', data: { id } });
}

/** 批量删除车辆 */
export function fetchBatchDeleteVehicle(ids: number[]) {
  return request({ url: '/transport/vehicle/batchDelete', method: 'post', data: { ids } });
}
