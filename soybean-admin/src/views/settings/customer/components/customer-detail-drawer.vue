<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { formatDateTime, getCustomerSourceLabel, getCustomerStatusMeta } from './customer-meta';
import { useCustomerData } from './use-customer-data';

defineOptions({ name: 'CustomerDetailDrawer' });

interface Props {
  customerId?: number | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'edit', row: Api.Customer.Item): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));
const { getDetail } = useCustomerData();
const loading = ref(false);
const detail = ref<Api.Customer.Item | null>(null);
const statusMeta = computed(() => getCustomerStatusMeta(detail.value?.status ?? 1));
const avatarStyle = computed(() => ({
  backgroundColor: `rgb(var(--${statusMeta.value.color}-color) / 0.12)`,
  color: `rgb(var(--${statusMeta.value.color}-color))`
}));

watch(
  [visible, () => props.customerId],
  async ([isVisible]) => {
    if (!isVisible) return;
    if (!props.customerId) {
      detail.value = null;
      return;
    }
    loading.value = true;
    try {
      detail.value = await getDetail(props.customerId);
    } finally {
      loading.value = false;
    }
  },
  { immediate: true }
);

function handleClose() {
  visible.value = false;
  detail.value = null;
}

function handleEdit() {
  if (detail.value) emit('edit', detail.value);
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth">
    <NDrawerContent title="客户详情" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="客户信息不存在或已被删除" />
        <template v-else-if="detail">
          <div class="flex-y-center gap-12px pb-16px">
            <div class="h-48px w-48px flex-center rounded-full text-18px font-medium" :style="avatarStyle">
              {{ detail.name.slice(0, 1) }}
            </div>
            <div class="flex-col gap-4px">
              <div class="flex-y-center gap-8px">
                <span class="text-16px font-medium">{{ detail.name }}</span>
                <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
              </div>
              <span class="text-12px text-gray-400">{{ getCustomerSourceLabel(detail.source) }} · 创建于 {{ formatDateTime(detail.createdAt) }}</span>
            </div>
          </div>

          <NDivider class="mt-0" />

          <NDescriptions :column="1" label-placement="left" :label-width="116" bordered size="small">
            <NDescriptionsItem label="客户姓名">{{ detail.name }}</NDescriptionsItem>
            <NDescriptionsItem label="手机号">
              <a class="text-primary" :href="`tel:${detail.phone}`">{{ detail.phone }}</a>
            </NDescriptionsItem>
            <NDescriptionsItem label="房屋地址">{{ detail.address || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="客户来源">{{ getCustomerSourceLabel(detail.source) }}</NDescriptionsItem>
            <NDescriptionsItem label="客户状态">
              <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
            </NDescriptionsItem>
            <NDescriptionsItem label="创建时间">{{ formatDateTime(detail.createdAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="更新时间">{{ formatDateTime(detail.updatedAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="备注">{{ detail.remark || '-' }}</NDescriptionsItem>
          </NDescriptions>
        </template>
      </NSpin>
      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="handleClose">关闭</NButton>
          <NButton type="primary" :disabled="!detail" @click="handleEdit">编辑客户</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>
