import { ref } from 'vue';
import {
  fetchBatchDeleteVehicle,
  fetchCreateVehicle,
  fetchDeleteVehicle,
  fetchGetDriverOptions,
  fetchGetVehicleList,
  fetchUpdateVehicle
} from '@/service/api/transport';

/** 车辆模块数据 hook：列表、增删改、司机下拉统一走 /transport 接口。 */
export function useVehicleData() {
  const loading = ref(false);
  const records = ref<Api.Transport.Vehicle[]>([]);
  const total = ref(0);
  const driverOptions = ref<Api.Transport.DriverOption[]>([]);

  async function loadList(params: Api.Transport.VehicleSearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetVehicleList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        status: params.status ?? null
      });
      if (!error && data) {
        records.value = data.records ?? [];
        total.value = data.total ?? 0;
      }
    } finally {
      loading.value = false;
    }
  }

  async function createVehicle(form: Api.Transport.VehicleFormData) {
    const { error } = await fetchCreateVehicle(form);
    return !error;
  }

  async function updateVehicle(id: number, form: Api.Transport.VehicleFormData) {
    const { error } = await fetchUpdateVehicle({ ...form, id });
    return !error;
  }

  async function deleteVehicles(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeleteVehicle(ids) : await fetchDeleteVehicle(ids[0]);
    return !error;
  }

  async function loadDriverOptions() {
    const { data, error } = await fetchGetDriverOptions();
    if (!error && data) driverOptions.value = data.items ?? [];
    return driverOptions.value;
  }

  return { loading, records, total, driverOptions, loadList, createVehicle, updateVehicle, deleteVehicles, loadDriverOptions };
}
