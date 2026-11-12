<script setup lang="ts">
import { computed, onMounted } from 'vue';
import dayjs from 'dayjs';
import type { RouteKey } from '@elegant-router/types';
import { useRouterPush } from '@/hooks/common/router';
import { useAuthStore } from '@/store/modules/auth';
import { formatOrderMoney, getOrderStatusMeta } from '@/views/order/components/order-meta';
import { useHomeData } from './modules/use-home-data';
import HomeOrderOverview from './modules/order-overview.vue';

const authStore = useAuthStore();
const { routerPushByKey } = useRouterPush();

const {
  loading,
  loaded,
  customerStats,
  orderStats,
  paymentStats,
  recentOrders,
  isEmpty,
  canViewCustomer,
  canViewOrder,
  canViewPayment,
  canViewPayout,
  loadData
} = useHomeData();

const WEEKDAYS = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];
const today = dayjs();
const todayText = `${today.format('YYYY年M月D日')} ${WEEKDAYS[today.day()]}`;
const displayName = computed(() => authStore.userInfo.nickname || authStore.userInfo.userName || '管理员');
const avatarText = computed(() => displayName.value.slice(0, 1));

const greeting = computed(() => {
  const hour = dayjs().hour();
  if (hour < 6) return '夜深了';
  if (hour < 9) return '早上好';
  if (hour < 12) return '上午好';
  if (hour < 14) return '中午好';
  if (hour < 18) return '下午好';
  return '晚上好';
});

const quickActions = computed(() =>
  [
    {
      key: 'customer',
      label: '客户管理',
      icon: 'mdi:account-multiple-outline',
      color: 'primary',
      route: 'settings_customer' as RouteKey,
      allowed: canViewCustomer.value
    },
    {
      key: 'order',
      label: '全部订单',
      icon: 'mdi:clipboard-text-clock-outline',
      color: 'info',
      route: 'order_all' as RouteKey,
      allowed: canViewOrder.value
    },
    {
      key: 'monitor',
      label: '行程监控',
      icon: 'mdi:map-marker-path',
      color: 'warning',
      route: 'order_monitor' as RouteKey,
      allowed: canViewOrder.value
    },
    {
      key: 'payment',
      label: '收款管理',
      icon: 'mdi:cash-multiple',
      color: 'success',
      route: 'finance_payment' as RouteKey,
      allowed: canViewPayment.value
    },
    {
      key: 'payout',
      label: '付款管理',
      icon: 'mdi:cash-minus',
      color: 'error',
      route: 'finance_payout' as RouteKey,
      allowed: canViewPayout.value
    }
  ].filter(action => action.allowed)
);

function go(route: RouteKey) {
  void routerPushByKey(route);
}

onMounted(() => {
  void loadData();
});
</script>

<template>
  <div class="flex flex-col gap-16px p-16px">
    <section class="hero-card">
      <div class="flex flex-wrap items-center justify-between gap-16px">
        <div class="flex items-center gap-14px">
          <span class="flex h-48px w-48px items-center justify-center rounded-14px bg-primary text-19px font-semibold text-white shadow-lg">
            {{ avatarText }}
          </span>
          <div>
            <h2 class="m-0 text-20px font-semibold">{{ greeting }}，{{ displayName }}</h2>
            <p class="mt-4px flex items-center gap-6px text-12px text-gray-400">
              <SvgIcon icon="mdi:calendar-blank-outline" class="text-14px" />
              <span>{{ todayText }}</span>
              <i class="h-3px w-3px rounded-full bg-gray-300" />
              <span>订单运营概览</span>
            </p>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-18px">
          <div>
            <p class="text-11px text-gray-400">本月回款</p>
            <p class="mt-2px text-18px font-semibold text-success">¥{{ (paymentStats.monthAmount ?? 0).toFixed(2) }}</p>
          </div>
          <div>
            <p class="text-11px text-gray-400">今日订单</p>
            <p class="mt-2px text-18px font-semibold">{{ orderStats.today }}</p>
          </div>
          <div>
            <p class="text-11px text-gray-400">进行中</p>
            <p class="mt-2px text-18px font-semibold text-warning">{{ orderStats.active }}</p>
          </div>
          <div>
            <p class="text-11px text-gray-400">待跟进客户</p>
            <p class="mt-2px text-18px font-semibold">{{ customerStats.following }}</p>
          </div>
          <NButton size="small" :loading="loading" @click="loadData">
            <template #icon><SvgIcon icon="mdi:refresh" /></template>
            刷新
          </NButton>
        </div>
      </div>

      <div class="mt-16px flex flex-wrap gap-8px border-t border-#eef0f3 pt-14px">
        <button
          v-for="action in quickActions"
          :key="action.key"
          class="quick-entry"
          type="button"
          :style="{ '--entry-color': `rgb(var(--${action.color}-color))`, '--entry-soft': `rgb(var(--${action.color}-color) / 0.1)` }"
          @click="go(action.route)"
        >
          <SvgIcon :icon="action.icon" class="text-18px" />
          <span>{{ action.label }}</span>
          <SvgIcon icon="mdi:arrow-right" class="ml-auto text-14px opacity-50" />
        </button>
      </div>
    </section>

    <NAlert v-if="loaded && isEmpty" :bordered="false" type="info">
      还没有订单和客户数据：先到「客户管理」录入客户，再到「全部订单」创建第一笔行程订单。
    </NAlert>

    <NSpin :show="loading && !loaded">
      <div class="grid grid-cols-1 gap-16px xl:grid-cols-3">
        <div class="xl:col-span-2">
          <HomeOrderOverview :stats="orderStats" :orders="recentOrders" />
        </div>

        <div class="flex-col gap-16px">
          <div class="rounded-12px border border-#eef0f3 bg-white p-16px">
            <p class="text-14px font-semibold">业务概览</p>
            <p class="mt-2px text-11px text-gray-400">客户与财务关键指标</p>
            <div class="mt-14px grid grid-cols-2 gap-12px">
              <div class="rounded-10px bg-#f7f8fa p-12px">
                <p class="text-11px text-gray-400">客户总数</p>
                <p class="mt-4px text-20px font-semibold">{{ customerStats.total }}</p>
              </div>
              <div class="rounded-10px bg-#f7f8fa p-12px">
                <p class="text-11px text-gray-400">成交客户</p>
                <p class="mt-4px text-20px font-semibold">{{ customerStats.deal }}</p>
              </div>
              <div class="rounded-10px bg-#f7f8fa p-12px">
                <p class="text-11px text-gray-400">订单收入</p>
                <p class="mt-4px text-16px font-semibold text-success">{{ formatOrderMoney(orderStats.revenue) }}</p>
              </div>
              <div class="rounded-10px bg-#f7f8fa p-12px">
                <p class="text-11px text-gray-400">本月收款</p>
                <p class="mt-4px text-16px font-semibold text-primary">¥{{ (paymentStats.monthAmount ?? 0).toFixed(2) }}</p>
              </div>
            </div>
          </div>

          <div class="rounded-12px border border-#eef0f3 bg-white p-16px">
            <p class="text-14px font-semibold">订单状态</p>
            <div class="mt-12px flex-col gap-10px">
              <div v-for="item in orderStats.statusCounts" :key="item.status" class="flex items-center justify-between text-12px">
                <span class="text-gray-500">{{ getOrderStatusMeta(item.status).label }}</span>
                <span class="font-medium">{{ item.count }}</span>
              </div>
              <p v-if="!orderStats.statusCounts.length" class="text-12px text-gray-400">暂无订单数据</p>
            </div>
          </div>
        </div>
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
.hero-card {
  position: relative;
  overflow: hidden;
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 14px;
  padding: 18px 20px 14px;
  background:
    radial-gradient(900px 300px at 92% -40%, rgb(var(--primary-color) / 0.14), transparent 62%),
    rgb(var(--container-bg-color));
  box-shadow: 0 10px 30px rgb(0 21 41 / 0.035);
}

.quick-entry {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 150px;
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 9px;
  padding: 9px 12px;
  font-size: 12px;
  font-weight: 500;
  color: var(--entry-color);
  background: var(--entry-soft);
  cursor: pointer;
  transition:
    transform 0.15s,
    box-shadow 0.15s;
}

.quick-entry:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgb(0 21 41 / 0.08);
}
</style>
