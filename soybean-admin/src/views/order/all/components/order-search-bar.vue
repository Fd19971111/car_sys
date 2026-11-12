<script setup lang="ts">
import { CUSTOMER_RESPONSE_OPTIONS, ORDER_STATUS_OPTIONS } from '../../components/order-meta';

defineOptions({ name: 'OrderSearchBar' });

interface Props {
  keyword: string;
  status: Api.Order.Status | null;
  customerResponse: Api.Order.CustomerResponse | null;
  dateFrom: string | null;
  dateTo: string | null;
  loading: boolean;
  checkedCount: number;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'update:keyword', value: string): void;
  (e: 'update:status', value: Api.Order.Status | null): void;
  (e: 'update:customerResponse', value: Api.Order.CustomerResponse | null): void;
  (e: 'update:dateFrom', value: string | null): void;
  (e: 'update:dateTo', value: string | null): void;
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
      class="w-260px min-w-170px flex-shrink"
      placeholder="搜索订单号 / 客户 / 手机号 / 上下车地点 / 车牌"
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
      class="w-140px flex-shrink-0"
      :options="ORDER_STATUS_OPTIONS"
      placeholder="订单状态"
      clearable
      @update:value="value => emit('update:status', value)"
    />
    <NSelect
      :value="customerResponse"
      class="w-130px flex-shrink-0"
      :options="CUSTOMER_RESPONSE_OPTIONS"
      placeholder="客户反馈"
      clearable
      @update:value="value => emit('update:customerResponse', value)"
    />
    <NDatePicker
      :formatted-value="dateFrom"
      class="w-140px flex-shrink-0"
      type="date"
      value-format="yyyy-MM-dd"
      placeholder="预约开始"
      clearable
      @update:formatted-value="value => emit('update:dateFrom', value)"
    />
    <NDatePicker
      :formatted-value="dateTo"
      class="w-140px flex-shrink-0"
      type="date"
      value-format="yyyy-MM-dd"
      placeholder="预约结束"
      clearable
      @update:formatted-value="value => emit('update:dateTo', value)"
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
          <SvgIcon icon="mdi:plus-circle-outline" />
        </template>
        新增订单
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
        确认删除选中的 {{ checkedCount }} 个订单吗？
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
