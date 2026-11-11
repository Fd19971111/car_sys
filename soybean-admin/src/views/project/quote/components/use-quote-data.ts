import { ref } from 'vue';
import {
  fetchBatchDeleteQuote,
  fetchCreateQuote,
  fetchDeleteQuote,
  fetchGetQuoteDetail,
  fetchGetQuoteList,
  fetchGetQuoteStats,
  fetchUpdateQuote
} from '@/service/api/quote';

/** 报价明细模块数据 Hook：全部数据走后端 /quote 接口（GoFrame + MySQL） */
export function useQuoteData() {
  /** 列表加载状态 */
  const loading = ref(false);
  /** 当前页报价明细 */
  const records = ref<Api.Quote.Item[]>([]);
  /** 报价明细总数（过滤后） */
  const total = ref(0);
  /** 概览统计 */
  const stats = ref<Api.Quote.Stats>({ totalCount: 0, totalAmount: 0, projectCount: 0, avgAmount: 0 });

  /** 获取报价明细分页列表 */
  async function loadList(params: Api.Quote.SearchParams) {
    loading.value = true;

    try {
      const { data, error } = await fetchGetQuoteList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        projectId: params.projectId ?? null
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
    const { data, error } = await fetchGetQuoteStats();

    if (!error && data) {
      stats.value = data;
    }
  }

  /** 获取报价明细详情，失败返回 null */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetQuoteDetail(id);

    return error || !data ? null : data;
  }

  /** 新增报价明细，返回是否成功及超额提示 */
  async function createQuote(form: Api.Quote.FormData) {
    const { data, error } = await fetchCreateQuote(form);

    if (error || !data) {
      return { success: false, warning: '' };
    }

    return { success: true, warning: data.warning ?? '' };
  }

  /** 更新报价明细，返回是否成功及超额提示 */
  async function updateQuote(id: number, form: Api.Quote.FormData) {
    const { data, error } = await fetchUpdateQuote({ ...form, id });

    if (error || !data) {
      return { success: false, warning: '' };
    }

    return { success: true, warning: data.warning ?? '' };
  }

  /** 删除报价明细（支持批量），返回是否成功 */
  async function deleteQuotes(ids: number[]) {
    if (!ids.length) {
      return false;
    }

    const { error } = ids.length > 1 ? await fetchBatchDeleteQuote(ids) : await fetchDeleteQuote(ids[0]);

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
    createQuote,
    updateQuote,
    deleteQuotes
  };
}
