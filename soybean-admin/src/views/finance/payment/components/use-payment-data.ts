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

/** 收款记录数据 hook。 */
export function usePaymentData() {
  const loading = ref(false);
  const records = ref<Api.Payment.Item[]>([]);
  const total = ref(0);
  const stats = ref<Api.Payment.Stats>({ totalCount: 0, totalAmount: 0, monthAmount: 0, monthCount: 0 });

  async function loadList(params: Api.Payment.SearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetPaymentList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        paymentType: params.paymentType?.trim() || null
      });
      if (!error && data) {
        records.value = data.records ?? [];
        total.value = data.total ?? 0;
      }
    } finally {
      loading.value = false;
    }
  }

  async function loadStats() {
    const { data, error } = await fetchGetPaymentStats();
    if (!error && data) stats.value = data;
  }

  async function getDetail(id: number) {
    const { data, error } = await fetchGetPaymentDetail(id);
    return error || !data ? null : data;
  }

  async function createPayment(form: Api.Payment.FormData) {
    const { data, error } = await fetchCreatePayment(form);
    return error || !data ? null : data;
  }

  async function updatePayment(id: number, form: Api.Payment.FormData) {
    const { data, error } = await fetchUpdatePayment({ ...form, id });
    return error || !data ? null : data;
  }

  async function deletePayments(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeletePayment(ids) : await fetchDeletePayment(ids[0]);
    return !error;
  }

  return { loading, records, total, stats, loadList, loadStats, getDetail, createPayment, updatePayment, deletePayments };
}
