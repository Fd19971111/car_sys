import { request } from '../request';

export function fetchGetSupplierList(params?: Api.Supplier.SearchParams) {
  return request<Api.Supplier.List>({ url: '/vendor/list', method: 'get', params });
}

export function fetchCreateSupplier(data: Api.Supplier.FormData) {
  return request<{ id: number }>({ url: '/vendor/create', method: 'post', data });
}

export function fetchUpdateSupplier(data: Api.Supplier.FormData & { id: number }) {
  return request<{ id: number }>({ url: '/vendor/update', method: 'post', data });
}

export function fetchDeleteSupplier(id: number) {
  return request({ url: '/vendor/delete', method: 'post', data: { id } });
}

export function fetchBatchDeleteSupplier(ids: number[]) {
  return request({ url: '/vendor/batchDelete', method: 'post', data: { ids } });
}