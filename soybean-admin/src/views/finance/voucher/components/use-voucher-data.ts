import { ref } from 'vue';
import {
  fetchCloseVoucherPeriod,
  fetchExportVouchers,
  fetchGetVoucherDetail,
  fetchGetVoucherList,
  fetchGetVoucherPeriods,
  fetchReopenVoucherPeriod
} from '@/service/api';

/** 记账凭证模块数据 Hook */
export function useVoucherData() {
  /** 列表加载状态 */
  const loading = ref(false);
  /** 当前页凭证 */
  const records = ref<Api.FinVoucher.Item[]>([]);
  /** 凭证总数 */
  const total = ref(0);

  /** 会计期间 */
  const periods = ref<Api.FinVoucher.Period[]>([]);
  const periodLoading = ref(false);
  const initialized = ref(false);
  const currentPeriod = ref('');
  const startPeriod = ref('');

  /** 获取凭证分页列表 */
  async function loadList(params: Api.FinVoucher.SearchParams) {
    loading.value = true;

    try {
      const { data, error } = await fetchGetVoucherList({
        current: params.current || 1,
        size: params.size || 10,
        periodCode: params.periodCode?.trim() || null,
        statusCode: params.statusCode ?? null,
        keyword: params.keyword?.trim() || null
      });

      if (error || !data) {
        return;
      }

      records.value = data.records ?? [];
      total.value = data.total ?? 0;
    } finally {
      loading.value = false;
    }
  }

  /** 加载会计期间 */
  async function loadPeriods() {
    periodLoading.value = true;

    try {
      const { data, error } = await fetchGetVoucherPeriods();

      if (!error && data) {
        periods.value = data.records ?? [];
        initialized.value = data.initialized;
        currentPeriod.value = data.currentPeriod;
        startPeriod.value = data.startPeriod;
      }
    } finally {
      periodLoading.value = false;
    }
  }

  /** 获取凭证详情，失败返回 null */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetVoucherDetail(id);

    return error || !data ? null : data;
  }

  /** 会计期间结账 */
  async function closePeriod(periodCode: string) {
    const { error } = await fetchCloseVoucherPeriod(periodCode);

    return !error;
  }

  /** 会计期间反结账 */
  async function reopenPeriod(periodCode: string) {
    const { error } = await fetchReopenVoucherPeriod(periodCode);

    return !error;
  }

  /** 导出凭证分录明细，失败返回 null */
  async function exportVouchers(params: Api.FinVoucher.ExportParams) {
    const { data, error } = await fetchExportVouchers(params);

    return error || !data ? null : data.records;
  }

  return {
    loading,
    records,
    total,
    periods,
    periodLoading,
    initialized,
    currentPeriod,
    startPeriod,
    loadList,
    loadPeriods,
    getDetail,
    closePeriod,
    reopenPeriod,
    exportVouchers
  };
}
