<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { NDescriptions, NDescriptionsItem, NTag } from 'naive-ui';
import { fetchGetOrderDetail } from '@/service/api/order';
import { useAppStore } from '@/store/modules/app';
import {
  formatOrderDateTime,
  formatOrderMoney,
  getCustomerResponseMeta,
  getOrderStatusMeta,
  getVehicleCategoryMeta
} from '../../components/order-meta';

defineOptions({ name: 'OrderDetailDrawer' });

interface Props {
  orderId: number | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'edit', row: Api.Order.Item): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 560));
const loading = ref(false);
const order = ref<Api.Order.Item | null>(null);

watch(
  () => [visible.value, props.orderId] as const,
  async ([show, id]) => {
    if (!show || !id) return;
    loading.value = true;
    try {
      const { data, error } = await fetchGetOrderDetail(id);
      if (!error && data) order.value = data;
    } finally {
      loading.value = false;
    }
  }
);

const statusMeta = computed(() => getOrderStatusMeta(order.value?.status ?? ''));
const responseMeta = computed(() => getCustomerResponseMeta(order.value?.customerResponse ?? ''));
const categoryMeta = computed(() => getVehicleCategoryMeta(order.value?.vehicleCategory ?? ''));
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth" placement="right">
    <NDrawerContent title="订单详情" closable :native-scrollbar="false" @close="visible = false">
      <NSpin :show="loading">
        <div v-if="order" class="flex-col gap-16px">
          <div class="rounded-10px bg-#f5f7fa p-14px">
            <div class="flex items-center justify-between gap-8px">
              <div>
                <p class="font-mono text-16px font-semibold">{{ order.orderNumber }}</p>
                <p class="mt-4px text-12px text-gray-400">创建于 {{ formatOrderDateTime(order.createdAt) }}</p>
              </div>
              <div class="flex items-center gap-6px">
                <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
                <NTag :type="responseMeta.tagType" size="small" :bordered="false">{{ responseMeta.label }}</NTag>
              </div>
            </div>
          </div>

          <div class="rounded-10px border border-#eef0f3 p-14px">
            <p class="mb-10px text-13px font-semibold">行程信息</p>
            <div class="flex-col gap-8px">
              <div class="flex items-start gap-8px">
                <span class="mt-6px h-7px w-7px flex-shrink-0 rounded-full bg-success" />
                <span class="text-13px">{{ order.pickupAddress }}</span>
              </div>
              <div class="ml-3px h-14px w-1px bg-#e5e7eb" />
              <div class="flex items-start gap-8px">
                <span class="mt-6px h-7px w-7px flex-shrink-0 rounded-full bg-error" />
                <span class="text-13px">{{ order.destinationAddress }}</span>
              </div>
            </div>
            <div class="mt-12px flex flex-wrap gap-16px text-12px text-gray-500">
              <span>预约：{{ formatOrderDateTime(order.scheduledAt) }}</span>
              <span>预计 {{ order.estimatedDistanceKm }} 公里 / {{ order.estimatedDurationMin }} 分钟</span>
              <span>乘车 {{ order.passengerCount }} 人</span>
              <span>车型：{{ categoryMeta.label }}</span>
            </div>
          </div>

          <NDescriptions :column="2" bordered size="small" label-placement="left">
            <NDescriptionsItem label="客户">{{ order.customerName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="手机号">{{ order.customerPhone || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="司机">{{ order.driverName || '未指派' }}</NDescriptionsItem>
            <NDescriptionsItem label="司机电话">{{ order.driverPhone || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="车辆">{{ order.vehiclePlate || '未指派' }}</NDescriptionsItem>
            <NDescriptionsItem label="客户报价">{{ formatOrderMoney(order.quotedFare) }}</NDescriptionsItem>
            <NDescriptionsItem label="最终费用">{{ formatOrderMoney(order.finalFare) }}</NDescriptionsItem>
            <NDescriptionsItem label="司机应得">{{ formatOrderMoney(order.driverFare) }}</NDescriptionsItem>
            <NDescriptionsItem label="结束时间">{{ formatOrderDateTime(order.completedAt) }}</NDescriptionsItem>
          </NDescriptions>

          <div class="rounded-10px border border-#eef0f3 p-14px">
            <p class="mb-6px text-13px font-semibold">内部备注</p>
            <p class="text-13px leading-22px text-gray-500">{{ order.internalNotes || '暂无备注' }}</p>
          </div>
        </div>
      </NSpin>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="visible = false">关闭</NButton>
          <NButton v-if="order" type="primary" @click="emit('edit', order)">编辑订单</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>
