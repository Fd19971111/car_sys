<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import {
  formatDate,
  formatDateTime,
  formatMoney,
  getCustomerSourceLabel,
  getCustomerStatusMeta,
  getProjectStatusMeta
} from './customer-meta';
import { useCustomerData } from './use-customer-data';

defineOptions({
  name: 'CustomerDetailDrawer'
});

interface Props {
  /** 当前查看的客户ID */
  customerId?: number | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'edit', row: Api.Customer.Item): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 560));

const { getDetail } = useCustomerData();

const loading = ref(false);
const detail = ref<Api.Customer.Item | null>(null);

const statusMeta = computed(() => getCustomerStatusMeta(detail.value?.status ?? 1));

const projects = computed(() => detail.value?.projects ?? []);

/** 该客户名下的合同总金额 */
const projectAmount = computed(() => projects.value.reduce((sum, project) => sum + (project.amount ?? 0), 0));

const avatarStyle = computed(() => ({
  backgroundColor: `rgb(var(--${statusMeta.value.color}-color) / 0.12)`,
  color: `rgb(var(--${statusMeta.value.color}-color))`
}));

watch(
  [visible, () => props.customerId],
  async ([isVisible]) => {
    if (!isVisible) {
      return;
    }

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
  if (detail.value) {
    emit('edit', detail.value);
  }
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth">
    <NDrawerContent title="客户详情" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="客户信息不存在或已被删除" />

        <template v-else-if="detail">
          <!-- 主要信息一览 -->
          <div class="flex-y-center gap-12px pb-16px">
            <div class="h-48px w-48px flex-center rounded-full text-18px font-medium" :style="avatarStyle">
              {{ detail.name.slice(0, 1) }}
            </div>
            <div class="flex-col gap-4px">
              <div class="flex-y-center gap-8px">
                <span class="text-16px font-medium">{{ detail.name }}</span>
                <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
              </div>
              <span class="text-12px text-gray-400">
                {{ getCustomerSourceLabel(detail.source) }} · 创建于 {{ formatDateTime(detail.createdAt) }}
              </span>
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
            <NDescriptionsItem v-if="projectAmount" label="合同金额">
              {{ formatMoney(projectAmount) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="创建时间">{{ formatDateTime(detail.createdAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="更新时间">{{ formatDateTime(detail.updatedAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="备注">{{ detail.remark || '-' }}</NDescriptionsItem>
          </NDescriptions>

          <!-- 关联装修项目 -->
          <div class="mt-20px">
            <div class="mb-12px flex-y-center justify-between">
              <span class="text-14px font-medium">关联装修项目</span>
              <span class="text-12px text-gray-400">共 {{ projects.length }} 个</span>
            </div>
            <NEmpty v-if="!projects.length" size="small" description="暂无关联项目（客户成交后会显示在这里）" />
            <div v-else class="flex-col gap-12px">
              <NCard v-for="project in projects" :key="project.id" size="small" :bordered="false" class="card-wrapper">
                <div class="flex-y-center justify-between gap-8px">
                  <span class="font-medium">{{ project.name }}</span>
                  <NTag :type="getProjectStatusMeta(project.status).tagType" size="small" :bordered="false">
                    {{ getProjectStatusMeta(project.status).label }}
                  </NTag>
                </div>
                <div class="mt-8px flex flex-wrap items-center gap-x-16px gap-y-4px text-12px text-gray-500">
                  <span>合同金额：{{ formatMoney(project.amount) }}</span>
                  <span>负责人：{{ project.managerName || '-' }}</span>
                  <span>工期：{{ formatDate(project.startDate) }} ~ {{ formatDate(project.endDate) }}</span>
                </div>
                <div v-if="project.latestProgress" class="mt-8px text-12px text-gray-500">
                  最近进度：{{ project.latestProgress }}
                  <span v-if="project.latestProgressAt">（{{ formatDate(project.latestProgressAt) }}）</span>
                </div>
              </NCard>
            </div>
          </div>
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

<style scoped></style>
