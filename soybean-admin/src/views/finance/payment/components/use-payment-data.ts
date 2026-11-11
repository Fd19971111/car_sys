import { ref } from 'vue';
import {
  fetchBatchDeletePayment,
  fetchCreatePayment,
  fetchDeletePayment,
  fetchGetPaymentDetail,
  fetchGetPaymentList,
  fetchGetPaymentStats,
  fetchUpdatePayment
} from '@/service/api/payment';

/** 收款记录模块数据 Hook：全部数据走后端 /payment 接口（GoFrame + MySQL） */
export function usePaymentData() {
  /** 列表加载状态 */
  const loading = ref(false);
  /** 当前页收款记录 */
  const records = ref<Api.Payment.Item[]>([]);
  /** 收款记录总数（过滤后） */
  const total = ref(0);
  /** 概览统计 */
  const stats = ref<Api.Payment.Stats>({ totalCount: 0, totalAmount: 0, monthAmount: 0, projectCount: 0 });

  /** 获取收款记录分页列表 */
  async function loadList(params: Api.Payment.SearchParams) {
    loading.value = true;

    try {
      const { data, error } = await fetchGetPaymentList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        projectId: params.projectId ?? null,
        paymentType: params.paymentType?.trim() || null
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

  /** 获取概览统计 */
  async function loadStats() {
    const { data, error } = await fetchGetPaymentStats();

    if (!error && data) {
      stats.value = data;
    }
  }

  /** 获取收款记录详情，失败返回 null */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetPaymentDetail(id);

    return error || !data ? null : data;
  }

  /**
   * 新增收款记录
   *
   * @returns 成功返回记录ID与凭证号，失败返回 null
   */
  async function createPayment(form: Api.Payment.FormData) {
    const { data, error } = await fetchCreatePayment(form);

    return error || !data ? null : data;
  }

  /**
   * 更新收款记录
   *
   * @returns 成功返回记录ID与凭证号，失败返回 null
   */
  async function updatePayment(id: number, form: Api.Payment.FormData) {
    const { data, error } = await fetchUpdatePayment({ ...form, id });

    return error || !data ? null : data;
  }

  /** 删除收款记录（支持批量），返回是否成功 */
  async function deletePayments(ids: number[]) {
    if (!ids.length) {
      return false;
    }

    const { error } = ids.length > 1 ? await fetchBatchDeletePayment(ids) : await fetchDeletePayment(ids[0]);

    return !error;
  }

  return {
    loading,
    records,
    total,
    stats,
    loadList,
    loadStats,
    getDetail,
    createPayment,
    updatePayment,
    deletePayments
  };
}
