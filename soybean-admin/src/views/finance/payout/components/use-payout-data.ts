import { ref } from 'vue';
import {
  fetchBatchDeletePayout,
  fetchCreatePayout,
  fetchDeletePayout,
  fetchGetPayoutDetail,
  fetchGetPayoutList,
  fetchGetPayoutStats,
  fetchUpdatePayout
} from '@/service/api/payout';

/** 付款记录数据 hook。 */
export function usePayoutData() {
  const loading = ref(false);
  const records = ref<Api.Payout.Item[]>([]);
  const total = ref(0);
  const stats = ref<Api.Payout.Stats>({ totalCount: 0, totalAmount: 0, monthAmount: 0, monthCount: 0 });

  async function loadList(params: Api.Payout.SearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetPayoutList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        payoutType: params.payoutType?.trim() || null
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
    const { data, error } = await fetchGetPayoutStats();
    if (!error && data) stats.value = data;
  }

  async function getDetail(id: number) {
    const { data, error } = await fetchGetPayoutDetail(id);
    return error || !data ? null : data;
  }

  async function createPayout(form: Api.Payout.FormData) {
    const { data, error } = await fetchCreatePayout(form);
    return error || !data ? null : data;
  }

  async function updatePayout(id: number, form: Api.Payout.FormData) {
    const { data, error } = await fetchUpdatePayout({ ...form, id });
    return error || !data ? null : data;
  }

  async function deletePayouts(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeletePayout(ids) : await fetchDeletePayout(ids[0]);
    return !error;
  }

  return { loading, records, total, stats, loadList, loadStats, getDetail, createPayout, updatePayout, deletePayouts };
}
