<script setup lang="ts">
import { NTag } from 'naive-ui';
import { formatOrderDateTime, getOrderStatusMeta } from '../../components/order-meta';

defineOptions({ name: 'MonitorList' });

interface Props {
  records: Api.Order.Item[];
  selectedId: number | null;
  loading: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'select', row: Api.Order.Item): void;
}>();
</script>

<template>
  <div class="monitor-list">
    <div class="mb-10px flex items-center justify-between">
      <p class="text-13px font-semibold">进行中订单</p>
      <span class="text-11px text-gray-400">{{ records.length }} 个</span>
    </div>

    <NSpin :show="loading">
      <div v-if="records.length" class="flex-col gap-8px">
        <button
          v-for="row in records"
          :key="row.id"
          class="monitor-item"
          :class="{ 'monitor-item--active': selectedId === row.id }"
          type="button"
          @click="emit('select', row)"
        >
          <div class="flex items-center justify-between gap-8px">
            <span class="font-mono text-12px font-semibold">{{ row.orderNumber }}</span>
            <NTag :type="getOrderStatusMeta(row.status).tagType" size="small" :bordered="false">
              {{ getOrderStatusMeta(row.status).label }}
            </NTag>
          </div>
          <p class="mt-6px truncate text-12px">{{ row.pickupAddress }} → {{ row.destinationAddress }}</p>
          <div class="mt-6px flex items-center justify-between gap-8px text-11px text-gray-400">
            <span class="truncate">{{ row.driverName || '待指派司机' }} · {{ row.vehiclePlate || '待派车' }}</span>
            <span class="flex-shrink-0">{{ formatOrderDateTime(row.scheduledAt) }}</span>
          </div>
        </button>
      </div>
      <div v-else class="flex-col items-center gap-6px py-40px text-center">
        <SvgIcon icon="mdi:map-marker-off-outline" class="text-32px text-gray-300" />
        <p class="text-13px text-gray-400">暂无进行中的行程</p>
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
.monitor-list {
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 12px;
  background: rgb(var(--container-bg-color));
  padding: 14px;
}

.monitor-item {
  width: 100%;
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 10px;
  padding: 10px 12px;
  text-align: left;
  background: transparent;
  cursor: pointer;
  transition:
    border-color 0.15s,
    background 0.15s;
}

.monitor-item:hover {
  border-color: rgb(var(--primary-color) / 0.35);
  background: rgb(var(--primary-color) / 0.03);
}

.monitor-item--active {
  border-color: rgb(var(--primary-color) / 0.55);
  background: rgb(var(--primary-color) / 0.06);
}
</style>
