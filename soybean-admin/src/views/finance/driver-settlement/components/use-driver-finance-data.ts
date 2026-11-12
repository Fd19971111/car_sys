import { ref } from 'vue';
import { fetchCreateDriverSettlement, fetchGetDriverBalances, fetchGetDriverSettlements } from '@/service/api';

export function useDriverFinanceData() {
  const loading = ref(false);
  const drivers = ref<Api.DriverFinance.DriverBalance[]>([]);
  const settlements = ref<Api.DriverFinance.Settlement[]>([]);
  const total = ref(0);

  async function load() {
    loading.value = true;
    try {
      const [driverResult, settlementResult] = await Promise.all([
        fetchGetDriverBalances(),
        fetchGetDriverSettlements({ current: 1, size: 50 })
      ]);
      if (!driverResult.error && driverResult.data) drivers.value = driverResult.data.records ?? [];
      if (!settlementResult.error && settlementResult.data) {
        settlements.value = settlementResult.data.records ?? [];
        total.value = settlementResult.data.total ?? 0;
      }
    } finally {
      loading.value = false;
    }
  }

  async function settle(data: Api.DriverFinance.SettlementForm) {
    const { data: result, error } = await fetchCreateDriverSettlement(data);
    if (error || !result) return null;
    await load();
    return result;
  }

  return { loading, drivers, settlements, total, load, settle };
}