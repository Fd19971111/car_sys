<script setup lang="ts">
import { VEHICLE_STATUS_OPTIONS } from './vehicle-meta';

defineOptions({ name: 'VehicleSearchBar' });

interface Props {
  keyword: string;
  status: Api.Transport.VehicleStatus | null;
  loading: boolean;
  checkedCount: number;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'update:keyword', value: string): void;
  (e: 'update:status', value: Api.Transport.VehicleStatus | null): void;
  (e: 'search'): void;
  (e: 'reset'): void;
  (e: 'refresh'): void;
  (e: 'add'): void;
  (e: 'batch-delete'): void;
}>();
</script>

<template>
  <div class="flex flex-wrap items-center gap-8px">
    <NInput
      :value="keyword"
      class="w-240px min-w-160px flex-shrink"
      placeholder="搜索车牌号 / 品牌 / 司机姓名"
      clearable
      @update:value="value => emit('update:keyword', value)"
      @keyup.enter="emit('search')"
    >
      <template #prefix>
        <SvgIcon icon="mdi:magnify" class="text-gray-400" />
      </template>
    </NInput>

    <NSelect
      :value="status"
      class="w-130px flex-shrink-0"
      :options="VEHICLE_STATUS_OPTIONS"
      placeholder="车辆状态"
      clearable
      @update:value="value => emit('update:status', value)"
    />

    <NButton class="flex-shrink-0" type="primary" ghost size="small" @click="emit('search')">
      <template #icon>
        <SvgIcon icon="mdi:magnify" />
      </template>
      查询
    </NButton>
    <NButton class="flex-shrink-0" size="small" @click="emit('reset')">
      <template #icon>
        <SvgIcon icon="mdi:backspace-outline" />
      </template>
      重置
    </NButton>

    <div class="ml-auto flex flex-shrink-0 items-center gap-8px">
      <NButton size="small" type="primary" ghost @click="emit('add')">
        <template #icon>
          <SvgIcon icon="mdi:car-plus" />
        </template>
        新增车辆
      </NButton>
      <NPopconfirm @positive-click="emit('batch-delete')">
        <template #trigger>
          <NButton size="small" ghost type="error" :disabled="!checkedCount">
            <template #icon>
              <SvgIcon icon="mdi:delete-outline" />
            </template>
            批量删除
          </NButton>
        </template>
        确认删除选中的 {{ checkedCount }} 辆车吗？
      </NPopconfirm>
      <NButton size="small" @click="emit('refresh')">
        <template #icon>
          <SvgIcon icon="mdi:refresh" :class="{ 'animate-spin': loading }" />
        </template>
        刷新
      </NButton>
    </div>
  </div>
</template>
