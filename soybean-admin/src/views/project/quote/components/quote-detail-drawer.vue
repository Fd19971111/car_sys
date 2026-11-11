<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { formatDate, formatDateTime, formatMoney, formatNumber } from '@/views/finance/components/finance-meta';
import { getQuotePercent } from './quote-meta';
import { useQuoteData } from './use-quote-data';

defineOptions({
  name: 'QuoteDetailDrawer'
});

interface Props {
  /** 当前查看的报价明细ID */
  quoteId?: number | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'edit', row: Api.Quote.Detail): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 620));

const { getDetail } = useQuoteData();

const loading = ref(false);
const detail = ref<Api.Quote.Detail | null>(null);

/** 报价合计占合同金额的比例 */
const quotePercent = computed(() => getQuotePercent(detail.value?.quoteAmount, detail.value?.projectAmount));

watch(
  [visible, () => props.quoteId],
  async ([isVisible]) => {
    if (!isVisible) {
      return;
    }

    if (!props.quoteId) {
      detail.value = null;

      return;
    }

    loading.value = true;

    try {
      detail.value = await getDetail(props.quoteId);
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
    <NDrawerContent title="报价明细详情" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="报价明细不存在或已被删除" />

        <template v-else-if="detail">
          <!-- 主要信息一览 -->
          <div class="flex-y-center gap-12px pb-16px">
            <div class="h-48px w-48px flex-center rounded-full bg-primary/12 text-18px font-medium text-primary">
              {{ detail.itemName.slice(0, 1) }}
            </div>
            <div class="min-w-0 flex-col gap-4px">
              <span class="truncate text-16px font-medium">{{ detail.itemName }}</span>
              <span class="text-12px text-gray-400">
                {{ detail.projectName || '未关联项目' }} · {{ formatDateTime(detail.createdAt) }}
              </span>
            </div>
          </div>

          <!-- 小计大字展示 -->
          <div class="mb-16px rounded-8px bg-primary/8 px-16px py-12px">
            <div class="text-12px text-gray-500">小计（单价 × 数量）</div>
            <div class="mt-4px text-24px font-medium text-primary">{{ formatMoney(detail.total) }}</div>
            <div class="mt-4px text-12px text-gray-400">
              {{ formatMoney(detail.price) }} × {{ formatNumber(detail.quantity) }}{{ detail.unit || '' }}
            </div>
          </div>

          <!-- 报价占合同金额比例 -->
          <div v-if="detail.projectAmount > 0" class="mb-16px flex-col gap-6px">
            <div class="flex-y-center justify-between text-12px">
              <span class="text-gray-500">该项目报价合计 / 合同金额</span>
              <span class="font-medium">
                {{ formatMoney(detail.quoteAmount) }} / {{ formatMoney(detail.projectAmount) }}（{{ quotePercent }}%）
              </span>
            </div>
            <div class="h-8px w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
              <div class="h-full rounded-full bg-primary transition-all" :style="{ width: `${quotePercent}%` }" />
            </div>
          </div>

          <NDescriptions :column="1" label-placement="left" :label-width="116" bordered size="small">
            <NDescriptionsItem label="报价项名称">{{ detail.itemName }}</NDescriptionsItem>
            <NDescriptionsItem label="所属项目">{{ detail.projectName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="客户">{{ detail.customerName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="单位">{{ detail.unit || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="单价">{{ formatMoney(detail.price) }}</NDescriptionsItem>
            <NDescriptionsItem label="数量">{{ formatNumber(detail.quantity) }}</NDescriptionsItem>
            <NDescriptionsItem label="小计">{{ formatMoney(detail.total) }}</NDescriptionsItem>
            <NDescriptionsItem v-if="detail.projectAmount > 0" label="合同金额">
              {{ formatMoney(detail.projectAmount) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="创建时间">{{ formatDateTime(detail.createdAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="更新时间">{{ formatDate(detail.updatedAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="备注">{{ detail.remark || '-' }}</NDescriptionsItem>
          </NDescriptions>
        </template>
      </NSpin>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="handleClose">关闭</NButton>
          <NButton type="primary" :disabled="!detail" @click="handleEdit">编辑报价</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
