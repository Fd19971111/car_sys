import { ref } from 'vue';
import {
  fetchBatchDeleteCustomer,
  fetchCreateCustomer,
  fetchDeleteCustomer,
  fetchGetCustomerDetail,
  fetchGetCustomerList,
  fetchGetCustomerStats,
  fetchUpdateCustomer
} from '@/service/api/customer';

/** 客户模块数据 Hook：全部数据走后端 /customer 接口（GoFrame + MySQL） */
export function useCustomerData() {
  /** 列表加载状态 */
  const loading = ref(false);
  /** 当前页客户数据 */
  const records = ref<Api.Customer.Item[]>([]);
  /** 客户总数（过滤后） */
  const total = ref(0);
  /** 概览统计 */
  const stats = ref<Api.Customer.Stats>({ total: 0, following: 0, deal: 0, lost: 0, dealAmount: 0 });

  /** 获取客户分页列表 */
  async function loadList(params: Api.Customer.SearchParams) {
    loading.value = true;

    try {
      const { data, error } = await fetchGetCustomerList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        status: params.status ?? null,
        source: params.source || null
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
    const { data, error } = await fetchGetCustomerStats();

    if (!error && data) {
      stats.value = data;
    }
  }

  /** 获取客户详情（含关联项目 / 负责人 / 最近进度），失败返回 null */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetCustomerDetail(id);

    return error || !data ? null : data;
  }

  /** 新增客户，返回是否成功 */
  async function createCustomer(form: Api.Customer.FormData) {
    const { error } = await fetchCreateCustomer(form);

    return !error;
  }

  /** 更新客户，返回是否成功 */
  async function updateCustomer(id: number, form: Api.Customer.FormData) {
    const { error } = await fetchUpdateCustomer({ ...form, id });

    return !error;
  }

  /** 删除客户（支持批量），返回是否成功 */
  async function deleteCustomers(ids: number[]) {
    if (!ids.length) {
      return false;
    }

    const { error } = ids.length > 1 ? await fetchBatchDeleteCustomer(ids) : await fetchDeleteCustomer(ids[0]);

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
    createCustomer,
    updateCustomer,
    deleteCustomers
  };
}
