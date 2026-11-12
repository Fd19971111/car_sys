<script setup lang="ts">
import { computed, ref } from 'vue';
import { useOrderMonitor } from './components/use-order-monitor';
import MonitorList from './components/monitor-list.vue';
import MonitorMap from './components/monitor-map.vue';
import OrderDetailDrawer from '../all/components/order-detail-drawer.vue';

defineOptions({ name: 'OrderMonitor' });

const { loading, autoRefresh, data, load, toggleAutoRefresh } = useOrderMonitor();

const selectedId = ref<number | null>(null);
const detailVisible = ref(false);

const selectedOrder = computed(() => {
  if (!data.value.records.length) return null;
  return data.value.records.find(item => item.id === selectedId.value) ?? data.value.records[0];
});

const summaryCards = computed(() => [
  { label: '进行中', value: data.value.activeCount, color: 'primary', icon: 'mdi:car-connected' },
  { label: '待派单', value: data.value.dispatching, color: 'warning', icon: 'mdi:clipboard-clock-outline' },
  { label: '司机前往中', value: data.value.driverEnRoute, color: 'info', icon: 'mdi:car-arrow-right' },
  { label: '行程中', value: data.value.inProgress, color: 'success', icon: 'mdi:map-marker-path' }
]);

function handleSelect(row: Api.Order.Item) {
  selectedId.value = row.id;
}

function openDetail() {
  detailVisible.value = true;
}
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <div class="grid grid-cols-2 gap-12px lg:grid-cols-4">
        <div v-for="card in summaryCards" :key="card.label" class="card-wrapper flex items-center gap-12px p-16px">
          <div
            class="flex h-42px w-42px flex-shrink-0 items-center justify-center rounded-10px"
            :style="{ backgroundColor: `rgb(var(--${card.color}-color) / 0.12)`, color: `rgb(var(--${card.color}-color))` }"
          >
            <SvgIcon :icon="card.icon" class="text-22px" />
          </div>
          <div>
            <p class="text-12px text-gray-400">{{ card.label }}</p>
            <p class="mt-2px text-20px font-semibold">{{ card.value }}<small class="ml-2px text-12px font-normal text-gray-400">单</small></p>
          </div>
        </div>
      </div>

      <div class="flex flex-wrap items-center justify-between gap-10px">
        <div>
          <h1 class="text-16px font-semibold">行程监控</h1>
          <p class="mt-4px text-13px text-gray-500">实时查看待派单、司机前往中与行程中的订单</p>
        </div>
        <div class="flex items-center gap-10px">
          <span class="text-12px text-gray-400">自动刷新</span>
          <NSwitch :value="autoRefresh" size="small" @update:value="toggleAutoRefresh" />
          <NButton size="small" :loading="loading" @click="load">
            <template #icon>
              <SvgIcon icon="mdi:refresh" />
            </template>
            刷新
          </NButton>
          <NButton size="small" type="primary" ghost :disabled="!selectedOrder" @click="openDetail">
            订单详情
          </NButton>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-16px xl:grid-cols-3">
        <div class="xl:col-span-2">
          <MonitorMap :record="selectedOrder" :active-count="data.activeCount" :updated-at="data.updatedAt" />
        </div>
        <MonitorList :records="data.records" :selected-id="selectedOrder?.id ?? null" :loading="loading" @select="handleSelect" />
      </div>
    </NSpace>

    <OrderDetailDrawer v-model:visible="detailVisible" :order-id="selectedOrder?.id ?? null" />
  </div>
</template>
