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

/** 付款记录模块数据 Hook：全部数据走后端 /payout 接口（GoFrame + MySQL） */
export function usePayoutData() {
  /** 列表加载状态 */
  const loading = ref(false);
  /** 当前页付款记录 */
  const records = ref<Api.Payout.Item[]>([]);
  /** 付款记录总数（过滤后） */
  const total = ref(0);
  /** 概览统计 */
  const stats = ref<Api.Payout.Stats>({ totalCount: 0, totalAmount: 0, monthAmount: 0, projectCount: 0 });

  /** 获取付款记录分页列表 */
  async function loadList(params: Api.Payout.SearchParams) {
    loading.value = true;

    try {
      const { data, error } = await fetchGetPayoutList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        projectId: params.projectId ?? null,
        payoutType: params.payoutType?.trim() || null,
        source: params.source ?? null
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
    const { data, error } = await fetchGetPayoutStats();

    if (!error && data) {
      stats.value = data;
    }
  }

  /** 获取付款记录详情，失败返回 null */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetPayoutDetail(id);

    return error || !data ? null : data;
  }

  /**
   * 新增付款记录
   *
   * @returns 成功返回记录ID与凭证号，失败返回 null
   */
  async function createPayout(form: Api.Payout.FormData) {
    const { data, error } = await fetchCreatePayout(form);

    return error || !data ? null : data;
  }

  /**
   * 更新付款记录
   *
   * @returns 成功返回记录ID与凭证号，失败返回 null
   */
  async function updatePayout(id: number, form: Api.Payout.FormData) {
    const { data, error } = await fetchUpdatePayout({ ...form, id });

    return error || !data ? null : data;
  }

  /** 删除付款记录（支持批量），返回是否成功 */
  async function deletePayouts(ids: number[]) {
    if (!ids.length) {
      return false;
    }

    const { error } = ids.length > 1 ? await fetchBatchDeletePayout(ids) : await fetchDeletePayout(ids[0]);

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
    createPayout,
    updatePayout,
    deletePayouts
  };
}
