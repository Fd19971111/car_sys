<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { formatDate, formatDateTime, formatMoney } from '../../components/finance-meta';
import { getPayoutTypeLabel, getPayoutTypeMeta } from './payout-meta';
import { usePayoutData } from './use-payout-data';

defineOptions({ name: 'PayoutDetailDrawer' });

interface Props {
  payoutId?: number | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'edit', row: Api.Payout.Detail): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 560));
const { getDetail } = usePayoutData();
const loading = ref(false);
const detail = ref<Api.Payout.Detail | null>(null);
const typeMeta = computed(() => getPayoutTypeMeta(detail.value?.payoutType));

watch(
  [visible, () => props.payoutId],
  async ([isVisible]) => {
    if (!isVisible) return;
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
  if (detail.value) emit('edit', detail.value);
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth">
    <NDrawerContent title="付款记录详情" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="付款记录不存在或已被删除" />
        <template v-else-if="detail">
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
                <NTag :type="typeMeta.tagType" size="small" :bordered="false">{{ getPayoutTypeLabel(detail.payoutType) }}</NTag>
              </div>
              <span class="text-12px text-gray-400">付款日期 {{ formatDate(detail.payoutDate) }}</span>
            </div>
          </div>

          <NDescriptions :column="1" label-placement="left" :label-width="116" bordered size="small">
            <NDescriptionsItem label="费用类型">{{ getPayoutTypeLabel(detail.payoutType) }}</NDescriptionsItem>
            <NDescriptionsItem label="本次付款">{{ formatMoney(detail.amount) }}</NDescriptionsItem>
            <NDescriptionsItem label="收款方">{{ detail.vendorName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="付款日期">{{ formatDate(detail.payoutDate) }}</NDescriptionsItem>
            <NDescriptionsItem label="付款账户">{{ detail.fundAccountName || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="登记时间">{{ formatDateTime(detail.createdAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="备注">{{ detail.remark || '-' }}</NDescriptionsItem>
          </NDescriptions>

          <section class="ledger">
            <div class="ledger__head">
              <span class="ledger__title">自动记账</span>
              <NTag v-if="detail.voucherNo" type="success" size="small" :bordered="false">已生成凭证 {{ detail.voucherNo }}</NTag>
              <NTag v-else type="default" size="small" :bordered="false">未记账</NTag>
            </div>
            <div v-if="detail.ledgerEntries.length" class="ledger__lines">
              <div v-for="(entry, index) in detail.ledgerEntries" :key="index" class="ledger__line">
                <span class="ledger__direction" :class="entry.direction === 'DEBIT' ? 'ledger__direction--debit' : 'ledger__direction--credit'">
                  {{ entry.direction === 'DEBIT' ? '借' : '贷' }}
                </span>
                <span class="ledger__account">{{ entry.accountName }}</span>
                <span class="ledger__amount">{{ formatMoney(entry.amount) }}</span>
              </div>
            </div>
            <p v-else class="ledger__hint">该记录暂未生成记账凭证（历史数据或财务尚未初始化）。</p>
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
  justify-content: space-between;
  gap: 8px;
  padding: 10px 12px;
  background: rgb(var(--base-text-color) / 0.03);
}
.ledger__title {
  font-size: 13px;
  font-weight: 600;
}
.ledger__lines {
  padding: 4px 12px;
}
.ledger__line {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
  border-bottom: 1px dashed rgb(var(--base-text-color) / 0.08);
  font-size: 13px;
}
.ledger__line:last-child {
  border-bottom: none;
}
.ledger__direction {
  flex-shrink: 0;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
}
.ledger__direction--debit {
  background: rgb(var(--primary-color) / 0.1);
  color: rgb(var(--primary-color));
}
.ledger__direction--credit {
  background: rgb(var(--success-color) / 0.1);
  color: rgb(var(--success-color));
}
.ledger__account {
  flex: 1;
}
.ledger__amount {
  font-family: 'SF Mono', 'Fira Code', monospace;
}
.ledger__hint {
  padding: 12px;
  font-size: 12px;
  line-height: 20px;
  color: rgb(var(--base-text-color) / 0.45);
}
</style>
