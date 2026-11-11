import { ref } from 'vue';
import {
  fetchBatchDeleteWarehouse,
  fetchCreateWarehouse,
  fetchDeleteWarehouse,
  fetchGetWarehouseList,
  fetchUpdateWarehouse
} from '@/service/api';

export function useWarehouseData() {
  const loading = ref(false);
  const records = ref<Api.Inventory.Warehouse[]>([]);
  const total = ref(0);

  async function loadList(params: Api.Inventory.WarehouseSearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetWarehouseList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null
      });
      if (!error && data) {
        records.value = data.records ?? [];
        total.value = data.total ?? 0;
      }
    } finally {
      loading.value = false;
    }
  }

  async function createWarehouse(form: Api.Inventory.WarehouseForm) {
    const { error } = await fetchCreateWarehouse(form);
    return !error;
  }

  async function updateWarehouse(id: number, form: Api.Inventory.WarehouseForm) {
    const { error } = await fetchUpdateWarehouse({ ...form, id });
    return !error;
  }

  async function deleteWarehouses(ids: number[]) {
    if (!ids.length) return false;
    const { error, data } = ids.length > 1
      ? await fetchBatchDeleteWarehouse(ids)
      : await fetchDeleteWarehouse(ids[0]);
    if (error) {
      // 默认仓库等业务拒绝信息由请求拦截器统一弹出，这里只返回失败
      return false;
    }
    void data;
    return true;
  }

  return { loading, records, total, loadList, createWarehouse, updateWarehouse, deleteWarehouses };
}
