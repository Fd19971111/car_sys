import { ref } from 'vue';
import {
  fetchBatchDeleteDriver,
  fetchCreateDriver,
  fetchDeleteDriver,
  fetchGetDriverList,
  fetchGetDriverOptions,
  fetchUpdateDriver
} from '@/service/api/transport';

/** 司机模块数据 hook：列表、增删改、下拉选项统一走 /transport/driver 接口。 */
export function useDriverData() {
  const loading = ref(false);
  const records = ref<Api.Transport.Driver[]>([]);
  const total = ref(0);
  const driverOptions = ref<Api.Transport.DriverOption[]>([]);

  async function loadList(params: Api.Transport.DriverSearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetDriverList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        serviceStatus: params.serviceStatus ?? null
      });
      if (!error && data) {
        records.value = data.records ?? [];
        total.value = data.total ?? 0;
      }
    } finally {
      loading.value = false;
    }
  }

  async function createDriver(form: Api.Transport.DriverFormData) {
    const { error } = await fetchCreateDriver(form);
    return !error;
  }

  async function updateDriver(id: number, form: Api.Transport.DriverFormData) {
    const { error } = await fetchUpdateDriver({ ...form, id });
    return !error;
  }

  async function deleteDrivers(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeleteDriver(ids) : await fetchDeleteDriver(ids[0]);
    return !error;
  }

  async function loadDriverOptions() {
    const { data, error } = await fetchGetDriverOptions();
    if (!error && data) driverOptions.value = data.items ?? [];
    return driverOptions.value;
  }

  return {
    loading,
    records,
    total,
    driverOptions,
    loadList,
    createDriver,
    updateDriver,
    deleteDrivers,
    loadDriverOptions
  };
}
