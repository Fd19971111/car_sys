<script setup lang="ts">
import { DRIVER_STATUS_OPTIONS } from './driver-meta';

defineOptions({ name: 'DriverSearchBar' });

interface Props {
  keyword: string;
  serviceStatus: Api.Transport.DriverServiceStatus | null;
  loading: boolean;
  checkedCount: number;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'update:keyword', value: string): void;
  (e: 'update:serviceStatus', value: Api.Transport.DriverServiceStatus | null): void;
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
      class="w-220px min-w-150px flex-shrink"
      placeholder="搜索姓名 / 手机号 / 驾驶证号"
      clearable
      @update:value="value => emit('update:keyword', value)"
      @keyup.enter="emit('search')"
    >
      <template #prefix>
        <SvgIcon icon="mdi:magnify" class="text-gray-400" />
      </template>
    </NInput>

    <NSelect
      :value="serviceStatus"
      class="w-120px flex-shrink-0"
      :options="DRIVER_STATUS_OPTIONS"
      placeholder="接单状态"
      clearable
      @update:value="value => emit('update:serviceStatus', value)"
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
          <SvgIcon icon="mdi:account-plus-outline" />
        </template>
        新增司机
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
        确认删除选中的 {{ checkedCount }} 位司机吗？
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
