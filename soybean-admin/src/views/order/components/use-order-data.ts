import { ref } from 'vue';
import { fetchGetCustomerList } from '@/service/api/customer';
import {
  fetchBatchDeleteOrder,
  fetchCreateOrder,
  fetchDeleteOrder,
  fetchGetOrderList,
  fetchGetOrderStats,
  fetchUpdateOrder
} from '@/service/api/order';
import { fetchGetDriverOptions, fetchGetVehicleList } from '@/service/api/transport';

/** 订单中心数据 hook：列表、统计、表单下拉与增删改。 */
export function useOrderData() {
  const loading = ref(false);
  const records = ref<Api.Order.Item[]>([]);
  const total = ref(0);
  const stats = ref<Api.Order.Stats>({
    total: 0,
    today: 0,
    active: 0,
    completed: 0,
    cancelled: 0,
    revenue: 0,
    statusCounts: []
  });
  const customerOptions = ref<{ label: string; value: number }[]>([]);
  const driverOptions = ref<{ label: string; value: number }[]>([]);
  const vehicleOptions = ref<{ label: string; value: number }[]>([]);

  async function loadList(params: Api.Order.SearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetOrderList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        status: params.status ?? null,
        customerResponse: params.customerResponse ?? null,
        dateFrom: params.dateFrom || null,
        dateTo: params.dateTo || null
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
    const { data, error } = await fetchGetOrderStats();
    if (!error && data) stats.value = data;
  }

  async function loadFormOptions() {
    const [customers, drivers, vehicles] = await Promise.all([
      fetchGetCustomerList({ current: 1, size: 200 }),
      fetchGetDriverOptions(),
      fetchGetVehicleList({ current: 1, size: 200 })
    ]);

    if (!customers.error && customers.data) {
      customerOptions.value = (customers.data.records ?? []).map(item => ({
        label: `${item.name} · ${item.phone}`,
        value: item.id
      }));
    }
    if (!drivers.error && drivers.data) {
      driverOptions.value = (drivers.data.items ?? []).map(item => ({
        label: `${item.name} · ${item.phone}`,
        value: item.id
      }));
    }
    if (!vehicles.error && vehicles.data) {
      vehicleOptions.value = (vehicles.data.records ?? []).map(item => ({
        label: `${item.plateNumber}${item.brand ? ` · ${item.brand}` : ''}`,
        value: item.id
      }));
    }
  }

  async function createOrder(form: Api.Order.FormData) {
    const { error } = await fetchCreateOrder(form);
    return !error;
  }

  async function updateOrder(id: number, form: Api.Order.FormData) {
    const { error } = await fetchUpdateOrder({ ...form, id });
    return !error;
  }

  async function deleteOrders(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeleteOrder(ids) : await fetchDeleteOrder(ids[0]);
    return !error;
  }

  return {
    loading,
    records,
    total,
    stats,
    customerOptions,
    driverOptions,
    vehicleOptions,
    loadList,
    loadStats,
    loadFormOptions,
    createOrder,
    updateOrder,
    deleteOrders
  };
}
