<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'MonitorMap' });

interface Props {
  record: Api.Order.Item | null;
  activeCount: number;
  updatedAt: string;
}

const props = defineProps<Props>();

const driverInitial = computed(() => props.record?.driverName?.slice(0, 1) || '—');
</script>

<template>
  <div class="monitor-map">
    <div class="monitor-map__head">
      <div>
        <h2 class="text-14px font-semibold">行程监控</h2>
        <p class="mt-3px text-11px text-gray-400">当前 {{ activeCount }} 个进行中订单</p>
      </div>
      <span class="flex items-center gap-5px text-11px font-medium text-success">
        <span class="h-6px w-6px animate-pulse rounded-full bg-success" />
        实时
      </span>
    </div>

    <div class="monitor-map__canvas">
      <div class="map-grid absolute inset-0" />
      <div class="absolute left-[9%] top-[26%] h-6px w-[86%] -rotate-6 rounded-full bg-white/90 shadow-sm" />
      <div class="absolute -left-[8%] top-[62%] h-6px w-[92%] rotate-12 rounded-full bg-white/90 shadow-sm" />
      <div class="absolute left-[48%] top-[-10%] h-[120%] w-6px rotate-6 rounded-full bg-white/90 shadow-sm" />
      <svg class="absolute inset-0 h-full w-full" viewBox="0 0 360 220" preserveAspectRatio="none">
        <path
          d="M35 172 C78 138, 96 158, 132 120 S212 92, 240 112 S294 88, 328 52"
          fill="none"
          stroke="rgb(var(--primary-color))"
          stroke-width="4"
          stroke-linecap="round"
          stroke-dasharray="2 9"
        />
      </svg>

      <div class="absolute left-[6%] top-[72%] flex h-28px w-28px items-center justify-center rounded-full bg-white text-primary shadow-md ring-4 ring-primary/10">
        <SvgIcon icon="mdi:map-marker" class="text-16px" />
      </div>
      <div class="absolute right-[5%] top-[13%] flex h-30px w-30px items-center justify-center rounded-9px bg-primary text-white shadow-lg">
        <SvgIcon icon="mdi:car" class="text-17px" />
      </div>
      <div class="absolute left-[45%] top-[38%] rounded-8px bg-white px-10px py-8px shadow-lg">
        <p class="text-9px text-gray-400">预计到达</p>
        <p class="text-12px font-semibold">{{ record ? `${record.estimatedDurationMin} 分钟` : '暂无行程' }}</p>
      </div>
    </div>

    <div class="mt-12px flex items-center gap-10px">
      <div class="flex h-36px w-36px items-center justify-center rounded-full bg-#172033 text-12px font-semibold text-white">
        {{ driverInitial }}
      </div>
      <div class="min-w-0 flex-1">
        <p class="truncate text-12px font-semibold">
          {{ record?.driverName || '待指派司机' }} · {{ record?.vehiclePlate || '待指派车辆' }}
        </p>
        <p class="mt-3px truncate text-10px text-gray-400">
          {{ record ? `${record.pickupAddress} → ${record.destinationAddress}` : '当前没有进行中的行程' }}
        </p>
      </div>
      <span class="rounded-6px bg-primary/10 px-7px py-4px text-10px font-medium text-primary">
        {{ record?.orderNumber || '空闲' }}
      </span>
    </div>

    <p class="mt-8px text-right text-10px text-gray-300">刷新于 {{ updatedAt || '-' }}</p>
  </div>
</template>

<style scoped>
.monitor-map {
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 12px;
  background: rgb(var(--container-bg-color));
  padding: 14px;
}

.monitor-map__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.monitor-map__canvas {
  position: relative;
  height: 260px;
  margin-top: 12px;
  overflow: hidden;
  border-radius: 10px;
  background: #eef3f7;
}

.map-grid {
  background-image:
    linear-gradient(rgb(255 255 255 / 0.55) 1px, transparent 1px),
    linear-gradient(90deg, rgb(255 255 255 / 0.55) 1px, transparent 1px);
  background-size: 32px 32px;
}
</style>
