import { ref } from 'vue';
import {
  fetchBatchDeleteSupplier,
  fetchCreateSupplier,
  fetchDeleteSupplier,
  fetchGetSupplierList,
  fetchUpdateSupplier
} from '@/service/api/supplier';

export function useSupplierData() {
  const loading = ref(false);
  const records = ref<Api.Supplier.Item[]>([]);
  const total = ref(0);

  async function loadList(params: Api.Supplier.SearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetSupplierList({
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

  async function createSupplier(form: Api.Supplier.FormData) {
    const { error } = await fetchCreateSupplier(form);
    return !error;
  }

  async function updateSupplier(id: number, form: Api.Supplier.FormData) {
    const { error } = await fetchUpdateSupplier({ ...form, id });
    return !error;
  }

  async function deleteSuppliers(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeleteSupplier(ids) : await fetchDeleteSupplier(ids[0]);
    return !error;
  }

  return { loading, records, total, loadList, createSupplier, updateSupplier, deleteSuppliers };
}