import { computed, ref } from 'vue';
import type { RouteKey } from '@elegant-router/types';
import { fetchGetCustomerStats } from '@/service/api/customer';
import { fetchGetOrderList, fetchGetOrderStats } from '@/service/api/order';
import { fetchGetPaymentStats } from '@/service/api/payment';
import { useRouteStore } from '@/store/modules/route';

/** 首页数据 hook：订单中心 + 客户 + 财务概览。 */
export function useHomeData() {
  const routeStore = useRouteStore();

  function hasRouteAccess(routeKeys: RouteKey[]) {
    function includesRoute(menus: App.Global.Menu[]): boolean {
      return menus.some(menu => routeKeys.includes(menu.routeKey) || (menu.children ? includesRoute(menu.children) : false));
    }
    return includesRoute(routeStore.menus);
  }

  const canViewCustomer = computed(() => hasRouteAccess(['settings_customer']));
  const canViewOrder = computed(() => hasRouteAccess(['order_all', 'order_monitor']));
  const canViewPayment = computed(() => hasRouteAccess(['finance_payment']));
  const canViewPayout = computed(() => hasRouteAccess(['finance_payout']));

  const loading = ref(false);
  const loaded = ref(false);

  const customerStats = ref<Api.Customer.Stats>({ total: 0, following: 0, deal: 0, lost: 0, dealAmount: 0 });
  const orderStats = ref<Api.Order.Stats>({
    total: 0,
    today: 0,
    active: 0,
    completed: 0,
    cancelled: 0,
    revenue: 0,
    statusCounts: []
  });
  const paymentStats = ref<Api.Payment.Stats>({ totalCount: 0, totalAmount: 0, monthAmount: 0, monthCount: 0 });
  const recentOrders = ref<Api.Order.Item[]>([]);

  const isEmpty = computed(() => orderStats.value.total === 0 && customerStats.value.total === 0);

  async function loadData() {
    loading.value = true;
    try {
      const [customerRes, orderStatsRes, paymentRes, orderListRes] = await Promise.all([
        canViewCustomer.value ? fetchGetCustomerStats() : Promise.resolve(null),
        canViewOrder.value ? fetchGetOrderStats() : Promise.resolve(null),
        canViewPayment.value ? fetchGetPaymentStats() : Promise.resolve(null),
        canViewOrder.value ? fetchGetOrderList({ current: 1, size: 6 }) : Promise.resolve(null)
      ]);

      if (customerRes && !customerRes.error && customerRes.data) customerStats.value = customerRes.data;
      if (orderStatsRes && !orderStatsRes.error && orderStatsRes.data) orderStats.value = orderStatsRes.data;
      if (paymentRes && !paymentRes.error && paymentRes.data) paymentStats.value = paymentRes.data;
      if (orderListRes && !orderListRes.error && orderListRes.data) recentOrders.value = orderListRes.data.records ?? [];

      loaded.value = true;
    } finally {
      loading.value = false;
    }
  }

  return {
    loading,
    loaded,
    canViewCustomer,
    canViewOrder,
    canViewPayment,
    canViewPayout,
    customerStats,
    orderStats,
    paymentStats,
    recentOrders,
    isEmpty,
    loadData
  };
}
