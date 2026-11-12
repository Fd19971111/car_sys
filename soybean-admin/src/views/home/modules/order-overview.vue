<script setup lang="ts">
import { useRouterPush } from '@/hooks/common/router';
import { formatOrderDateTime, formatOrderMoney, getOrderStatusMeta } from '@/views/order/components/order-meta';

defineOptions({ name: 'HomeOrderOverview' });

interface Props {
  stats: Api.Order.Stats;
  orders: Api.Order.Item[];
}

defineProps<Props>();

const { routerPushByKey } = useRouterPush();

const statCards = [
  { key: 'total', label: '订单总数', icon: 'mdi:clipboard-text-outline', color: 'primary' },
  { key: 'today', label: '今日预约', icon: 'mdi:calendar-today-outline', color: 'info' },
  { key: 'active', label: '进行中', icon: 'mdi:car-connected', color: 'warning' },
  { key: 'completed', label: '已完成', icon: 'mdi:check-circle-outline', color: 'success' }
] as const;
</script>

<template>
  <div class="flex-col gap-16px">
    <div class="grid grid-cols-2 gap-12px xl:grid-cols-4">
      <div v-for="card in statCards" :key="card.key" class="rounded-12px border border-#eef0f3 bg-white p-14px">
        <div class="flex items-center gap-8px">
          <SvgIcon :icon="card.icon" class="text-18px" :style="{ color: `rgb(var(--${card.color}-color))` }" />
          <span class="text-12px text-gray-400">{{ card.label }}</span>
        </div>
        <p class="mt-8px text-22px font-semibold">
          {{ card.key === 'total' ? stats.total : card.key === 'today' ? stats.today : card.key === 'active' ? stats.active : stats.completed }}
        </p>
      </div>
    </div>

    <div class="rounded-12px border border-#eef0f3 bg-white">
      <div class="flex items-center justify-between border-b border-#f0f2f5 px-16px py-12px">
        <div>
          <p class="text-14px font-semibold">最近订单</p>
          <p class="mt-2px text-11px text-gray-400">按预约时间倒序</p>
        </div>
        <NButton text type="primary" size="small" @click="routerPushByKey('order_all')">查看全部</NButton>
      </div>

      <div v-if="orders.length" class="divide-y divide-#f3f4f6">
        <div v-for="row in orders" :key="row.id" class="flex flex-wrap items-center gap-10px px-16px py-11px">
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-8px">
              <span class="font-mono text-12px font-semibold">{{ row.orderNumber }}</span>
              <NTag :type="getOrderStatusMeta(row.status).tagType" size="small" :bordered="false">
                {{ getOrderStatusMeta(row.status).label }}
              </NTag>
            </div>
            <p class="mt-4px truncate text-12px text-gray-500">{{ row.customerName }} · {{ row.pickupAddress }} → {{ row.destinationAddress }}</p>
          </div>
          <div class="text-right">
            <p class="text-12px font-medium">{{ formatOrderMoney(row.finalFare || row.quotedFare) }}</p>
            <p class="mt-2px text-10px text-gray-400">{{ formatOrderDateTime(row.scheduledAt) }}</p>
          </div>
        </div>
      </div>
      <div v-else class="flex-col items-center gap-6px py-40px text-center">
        <SvgIcon icon="mdi:clipboard-text-off-outline" class="text-32px text-gray-300" />
        <p class="text-13px text-gray-400">还没有订单，先去「全部订单」新增第一单</p>
        <NButton type="primary" ghost size="small" @click="routerPushByKey('order_all')">新增订单</NButton>
      </div>
    </div>
  </div>
</template>
