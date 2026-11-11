import { request } from '../request';

/**
 * 获取记账凭证分页列表
 *
 * @param params 查询参数（期间 / 状态 / 关键字 / 分页）
 */
export function fetchGetVoucherList(params?: Api.FinVoucher.SearchParams) {
  return request<Api.FinVoucher.List>({
    url: '/fin-voucher/list',
    method: 'get',
    params
  });
}

/**
 * 获取凭证详情（含借贷分录）
 *
 * @param id 凭证ID
 */
export function fetchGetVoucherDetail(id: number) {
  return request<Api.FinVoucher.Detail>({
    url: '/fin-voucher/detail',
    method: 'get',
    params: { id }
  });
}

/** 获取会计期间列表（含是否可结账 / 反结账） */
export function fetchGetVoucherPeriods() {
  return request<Api.FinVoucher.PeriodList>({
    url: '/fin-voucher/periods',
    method: 'get'
  });
}

/**
 * 会计期间结账
 *
 * @param periodCode 会计期间 YYYY-MM
 */
export function fetchCloseVoucherPeriod(periodCode: string) {
  return request<{ periodCode: string; statusCode: string }>({
    url: '/fin-voucher/close',
    method: 'post',
    data: { periodCode }
  });
}

/**
 * 会计期间反结账
 *
 * @param periodCode 会计期间 YYYY-MM
 */
export function fetchReopenVoucherPeriod(periodCode: string) {
  return request<{ periodCode: string; statusCode: string }>({
    url: '/fin-voucher/reopen',
    method: 'post',
    data: { periodCode }
  });
}

/**
 * 导出凭证分录明细（用于生成 CSV）
 *
 * @param params 导出条件（期间 / 状态 / 关键字）
 */
export function fetchExportVouchers(params?: Api.FinVoucher.ExportParams) {
  return request<Api.FinVoucher.ExportResult>({
    url: '/fin-voucher/export',
    method: 'get',
    params
  });
}
