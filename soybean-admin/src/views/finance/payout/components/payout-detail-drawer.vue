<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { formatDate, formatDateTime, formatMoney } from '../../components/finance-meta';
import { getPayoutTypeMeta, getPayoutTypeLabel } from './payout-meta';
import { usePayoutData } from './use-payout-data';

defineOptions({
  name: 'PayoutDetailDrawer'
});

interface Props {
  /** 当前查看的付款记录ID */
  payoutId?: number | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'edit', row: Api.Payout.Detail): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 620));

const { getDetail } = usePayoutData();

const loading = ref(false);
const detail = ref<Api.Payout.Detail | null>(null);

const typeMeta = computed(() => getPayoutTypeMeta(detail.value?.payoutType));

/** 项目累计付款占合同金额的比例（用于成本控制参考） */
const paidPercent = computed(() => {
  const projectAmount = detail.value?.projectAmount ?? 0;
  const paid = detail.value?.projectPaid ?? 0;

  if (projectAmount <= 0) return 0;

  return Math.min(100, Math.max(0, Math.round((paid / projectAmount) * 100)));
});

watch(
  [visible, () => props.payoutId],
  async ([isVisible]) => {
    if (!isVisible) {
      return;
    }

    if (!props.payoutId) {
      detail.value = null;

      return;
    }

    loading.value = true;

    try {
      detail.value = await getDetail(props.payoutId);
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
    <NDrawerContent title="付款记录详情" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="付款记录不存在或已被删除" />

        <template v-else-if="detail">
          <!-- 主要信息一览 -->
          <div class="flex-y-center gap-12px pb-16px">
            <div
              class="h-48px w-48px flex-center rounded-full text-18px font-medium"
              :style="{
                backgroundColor: `rgb(var(--${typeMeta.color}-color) / 0.12)`,
                color: `rgb(var(--${typeMeta.color}-color))`
              }"
            >
              {{ getPayoutTypeLabel(detail.payoutType).slice(0, 1) }}
            </div>
            <div class="min-w-0 flex-col gap-4px">
              <div class="flex-y-center gap-8px">
                <span class="text-16px font-medium">{{ formatMoney(detail.amount) }}</span>
                <NTag :type="typeMeta.tagType" size="small" :bordered="false">
                  {{ getPayoutTypeLabel(detail.payoutType) }}
                </NTag>
              </div>
              <span class="text-12px text-gray-400">
                供应商：{{ detail.vendorName || '未填写' }} · {{ formatDate(detail.payoutDate) }}
              </span>
            </div>
          </div>

          <!-- 该项目累计付款情况 -->
          <div v-if="detail.projectId > 0 && detail.projectAmount > 0" class="mb-16px flex-col gap-6px">
            <div class="flex-y-center justify-between text-12px">
              <span class="text-gray-500">该项目累计付款（占合同）</span>
              <span class="font-medium">
                {{ formatMoney(detail.projectPaid) }} / {{ formatMoney(detail.projectAmount) }}（{{ paidPercent }}%）
              </span>
            </div>
            <div class="h-8px w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
              <div
                class="h-full rounded-full transition-all"
                :style="{ width: `${paidPercent}%`, backgroundColor: 'rgb(var(--info-color))' }"
              />
            </div>
          </div>

          <NDescriptions :column="1" label-placement="left" :label-width="116" bordered size="small">
            <NDescriptionsItem label="所属项目">
              {{ detail.projectName || (detail.projectId > 0 ? '-' : '公司公共支出') }}
            </NDescriptionsItem>
            <NDescriptionsItem label="客户">{{ detail.customerName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="费用类型">
              <NTag :type="typeMeta.tagType" size="small" :bordered="false">
                {{ getPayoutTypeLabel(detail.payoutType) }}
              </NTag>
            </NDescriptionsItem>
            <NDescriptionsItem label="供应商">{{ detail.vendorName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="本次付款">{{ formatMoney(detail.amount) }}</NDescriptionsItem>
            <NDescriptionsItem label="付款日期">{{ formatDate(detail.payoutDate) }}</NDescriptionsItem>
            <NDescriptionsItem label="付款账户">{{ detail.fundAccountName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem v-if="detail.projectAmount > 0" label="合同金额">
              {{ formatMoney(detail.projectAmount) }}
            </NDescriptionsItem>
            <NDescriptionsItem v-if="detail.projectId > 0" label="累计已付">
              {{ formatMoney(detail.projectPaid) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="登记时间">{{ formatDateTime(detail.createdAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="备注">{{ detail.remark || '-' }}</NDescriptionsItem>
          </NDescriptions>

          <!-- 自动记账结果 -->
          <section class="ledger">
            <div class="ledger__head">
              <span class="ledger__title">自动记账</span>
              <NTag v-if="detail.voucherNo" type="success" size="small" :bordered="false">
                已生成凭证 {{ detail.voucherNo }}
              </NTag>
              <NTag v-else type="default" size="small" :bordered="false">未记账</NTag>
            </div>

            <div v-if="detail.ledgerEntries.length" class="ledger__lines">
              <div v-for="(entry, index) in detail.ledgerEntries" :key="index" class="ledger__line">
                <span
                  class="ledger__direction"
                  :class="entry.direction === 'DEBIT' ? 'ledger__direction--debit' : 'ledger__direction--credit'"
                >
                  {{ entry.direction === 'DEBIT' ? '借' : '贷' }}
                </span>
                <span class="ledger__account">{{ entry.accountName }}</span>
                <span class="ledger__amount">{{ formatMoney(entry.amount) }}</span>
              </div>
            </div>

            <p v-else class="ledger__hint">
              该记录暂未生成记账凭证（财务尚未初始化）。新登记的付款会自动按「借：成本 / 费用科目 / 贷：付款账户」记账。
            </p>
          </section>
        </template>
      </NSpin>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="handleClose">关闭</NButton>
          <NButton type="primary" :disabled="!detail" @click="handleEdit">编辑记录</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
.ledger {
  margin-top: 16px;
  overflow: hidden;
  border: 1px solid rgb(var(--base-text-color) / 0.08);
  border-radius: 10px;
}

.ledger__head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background-color: rgb(var(--base-text-color) / 0.03);
}

.ledger__title {
  font-size: 13px;
  font-weight: 600;
  color: rgb(var(--base-text-color) / 0.8);
}

.ledger__lines {
  display: flex;
  flex-direction: column;
}

.ledger__line {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  font-size: 13px;
  border-top: 1px solid rgb(var(--base-text-color) / 0.06);
}

.ledger__direction {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  font-size: 12px;
  font-weight: 600;
  border-radius: 6px;
}

.ledger__direction--debit {
  color: rgb(var(--primary-color));
  background-color: rgb(var(--primary-color) / 0.12);
}

.ledger__direction--credit {
  color: #d97706;
  background-color: rgb(217 119 6 / 0.12);
}

.ledger__account {
  color: rgb(var(--base-text-color) / 0.75);
}

.ledger__amount {
  margin-left: auto;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.85);
}

.ledger__hint {
  padding: 10px 12px;
  margin: 0;
  font-size: 12px;
  line-height: 1.7;
  color: rgb(var(--base-text-color) / 0.45);
}
</style>
