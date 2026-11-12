import { onMounted, onUnmounted, ref } from 'vue';
import { fetchGetOrderMonitor } from '@/service/api/order';

/** 行程监控数据 hook：15 秒轮询进行中订单。 */
export function useOrderMonitor() {
  const loading = ref(false);
  const autoRefresh = ref(true);
  const data = ref<Api.Order.Monitor>({
    activeCount: 0,
    dispatching: 0,
    driverEnRoute: 0,
    inProgress: 0,
    records: [],
    updatedAt: ''
  });
  let timer: ReturnType<typeof setInterval> | null = null;

  async function load() {
    loading.value = true;
    try {
      const { data: result, error } = await fetchGetOrderMonitor();
      if (!error && result) data.value = result;
    } finally {
      loading.value = false;
    }
  }

  function startTimer() {
    stopTimer();
    if (!autoRefresh.value) return;
    timer = setInterval(() => void load(), 15000);
  }

  function stopTimer() {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  }

  function toggleAutoRefresh(value: boolean) {
    autoRefresh.value = value;
    if (value) startTimer();
    else stopTimer();
  }

  onMounted(async () => {
    await load();
    startTimer();
  });

  onUnmounted(stopTimer);

  return { loading, autoRefresh, data, load, toggleAutoRefresh };
}
