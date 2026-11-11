<script setup lang="ts">
import { computed, h, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NTag } from 'naive-ui';
import { fetchGetProjectProfitDetail } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { getProjectStatusMeta } from '@/views/project/project-list/components/project-meta';

defineOptions({
  name: 'ProjectProfitDetailDrawer'
});

interface Props {
  projectId?: number | null;
}

const props = defineProps<Props>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 760));

const loading = ref(false);
const detail = ref<Api.FinReport.ProjectProfitDetail | null>(null);

const statusMeta = computed(() => getProjectStatusMeta(detail.value?.status ?? 1));

interface MetricItem {
  key: string;
  label: string;
  value: string;
  strong?: boolean;
  color?: string;
}

const metrics = computed<MetricItem[]>(() => {
  const item = detail.value;
  if (!item) return [];

  return [
    { key: 'contract', label: '合同金额', value: formatMoney(item.contractAmount) },
    { key: 'change', label: '变更净额', value: `${item.changeNet >= 0 ? '+' : ''}${formatMoney(item.changeNet)}` },
    { key: 'settlement', label: '结算金额', value: formatMoney(item.settlementAmount), strong: true },
    { key: 'received', label: '已收款', value: formatMoney(item.receivedAmount), color: 'success' },
    { key: 'unpaid', label: '待收款', value: formatMoney(item.unpaidAmount), color: 'warning' },
    { key: 'cost', label: '成本合计', value: formatMoney(item.costTotal), color: 'error' },
    {
      key: 'costRate',
      label: '成本占结算',
      value: `${item.settlementAmount > 0 ? Math.round((item.costTotal / item.settlementAmount) * 1000) / 10 : 0}%`
    },
    { key: 'materialUsed', label: '材料耗用（参考）', value: formatMoney(item.materialUsed) },
    {
      key: 'gross',
      label: '预计毛利',
      value: formatMoney(item.grossProfit),
      color: item.grossProfit >= 0 ? 'success' : 'error',
      strong: true
    },
    {
      key: 'grossRate',
      label: '毛利率',
      value: `${item.grossRate}%`,
      color: item.grossProfit >= 0 ? 'success' : 'error'
    },
    {
      key: 'cash',
      label: '现金结余',
      value: formatMoney(item.cashProfit),
      color: item.cashProfit >= 0 ? 'success' : 'error'
    }
  ];
});

const receiptColumns: DataTableColumns<Api.FinReport.ProjectProfitReceipt> = [
  { key: 'paymentDate', title: '日期', width: 100, render: row => h('span', { class: 'font-mono' }, row.paymentDate) },
  { key: 'paymentType', title: '款项类型', width: 90 },
  { key: 'amount', title: '金额', width: 110, align: 'right', render: row => formatMoney(row.amount) },
  { key: 'remark', title: '备注', ellipsis: { tooltip: true }, render: row => row.remark || '-' }
];

const payoutColumns: DataTableColumns<Api.FinReport.ProjectProfitPayout> = [
  { key: 'payoutDate', title: '日期', width: 100, render: row => h('span', { class: 'font-mono' }, row.payoutDate) },
  { key: 'payoutType', title: '类型', width: 84 },
  { key: 'vendorName', title: '供应商', width: 110, ellipsis: { tooltip: true }, render: row => row.vendorName || '-' },
  { key: 'amount', title: '金额', width: 110, align: 'right', render: row => formatMoney(row.amount) },
  { key: 'remark', title: '备注', ellipsis: { tooltip: true }, render: row => row.remark || '-' }
];

const materialColumns: DataTableColumns<Api.FinReport.ProjectProfitMaterial> = [
  { key: 'usedAt', title: '日期', width: 100, render: row => h('span', { class: 'font-mono' }, row.usedAt) },
  { key: 'materialName', title: '材料', minWidth: 110, ellipsis: { tooltip: true } },
  {
    key: 'quantity',
    title: '数量',
    width: 90,
    align: 'right',
    render: row => `${row.quantity} ${row.unit || ''}`
  },
  { key: 'unitPrice', title: '单价', width: 90, align: 'right', render: row => formatMoney(row.unitPrice) },
  { key: 'amount', title: '金额', width: 110, align: 'right', render: row => formatMoney(row.amount) },
  { key: 'receiver', title: '领用人', width: 90, render: row => row.receiver || '-' }
];

const changeColumns: DataTableColumns<Api.FinReport.ProjectProfitChange> = [
  { key: 'changeDate', title: '日期', width: 100, render: row => h('span', { class: 'font-mono' }, row.changeDate) },
  {
    key: 'changeType',
    title: '类型',
    width: 74,
    render: row =>
      h(
        NTag,
        { size: 'small', bordered: false, type: row.changeType === 'ADD' ? 'success' : 'warning' },
        { default: () => (row.changeType === 'ADD' ? '增项' : '减项') }
      )
  },
  { key: 'title', title: '变更事项', minWidth: 120, ellipsis: { tooltip: true } },
  {
    key: 'amount',
    title: '金额',
    width: 110,
    align: 'right',
    render: row => `${row.changeType === 'ADD' ? '+' : '-'}${formatMoney(row.amount)}`
  },
  { key: 'reason', title: '原因', ellipsis: { tooltip: true }, render: row => row.reason || '-' }
];

watch(
  [visible, () => props.projectId],
  async ([isVisible]) => {
    if (!isVisible) {
      return;
    }

    if (!props.projectId) {
      detail.value = null;

      return;
    }

    loading.value = true;

    try {
      detail.value = await fetchGetProjectProfitDetail(props.projectId).then(res => (res.error ? null : res.data));
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
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth">
    <NDrawerContent title="项目成本与毛利" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="项目不存在或已被删除" />

        <template v-else-if="detail">
          <div class="mb-14px flex-y-center justify-between gap-8px">
            <div class="min-w-0 flex-col gap-2px">
              <div class="flex-y-center gap-8px">
                <span class="truncate text-16px font-medium">{{ detail.projectName }}</span>
                <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
              </div>
              <span class="text-12px text-gray-400">
                {{ detail.customerName || '未关联客户' }} · {{ detail.startDate || '-' }} ~ {{ detail.endDate || '-' }}
              </span>
            </div>
            <span class="text-12px text-gray-400">回款率 {{ detail.receivedRate }}%</span>
          </div>

          <div class="metric-grid">
            <div v-for="item in metrics" :key="item.key" class="metric">
              <span class="metric__label">{{ item.label }}</span>
              <span
                class="metric__value"
                :class="[item.strong ? 'metric__value--strong' : '', item.color ? `metric__value--${item.color}` : '']"
              >
                {{ item.value }}
              </span>
            </div>
          </div>

          <NCard :bordered="false" size="small" class="mt-14px">
            <template #header>
              <span class="text-13px font-medium">成本构成（按付款归集）</span>
            </template>
            <NEmpty v-if="!detail.costItems.length" size="small" description="还没有归集到成本，登记项目付款后会自动统计" />
            <div v-else class="cost-list">
              <div v-for="item in detail.costItems" :key="item.type" class="cost-row">
                <span class="cost-row__type">{{ item.type }}</span>
                <div class="cost-row__bar">
                  <div class="cost-row__bar-inner" :style="{ width: `${item.percent}%` }" />
                </div>
                <span class="cost-row__amount">{{ formatMoney(item.amount) }}</span>
                <span class="cost-row__percent">{{ item.percent }}%</span>
              </div>
              <div class="cost-total">
                <span>成本合计</span>
                <strong>{{ formatMoney(detail.costTotal) }}</strong>
              </div>
              <div v-if="detail.materialUsed > 0" class="material-compare">
                <span>材料款付款 {{ formatMoney(detail.materialPaid) }}</span>
                <span>材料耗用 {{ formatMoney(detail.materialUsed) }}（数量 × 单价）</span>
                <NTag v-if="detail.materialUsed > detail.materialPaid + 0.01" type="warning" size="small" :bordered="false">
                  耗用高于付款，请核对
                </NTag>
              </div>
            </div>
          </NCard>

          <NCard :bordered="false" size="small" class="mt-14px">
            <template #header>
              <span class="text-13px font-medium">明细</span>
            </template>
            <NTabs type="line" animated>
              <NTabPane :tab="`收款（${detail.receipts.length}）`" name="receipts">
                <NDataTable
                  :columns="receiptColumns"
                  :data="detail.receipts"
                  :row-key="row => row.id"
                  :pagination="false"
                  size="small"
                  :max-height="260"
                />
              </NTabPane>
              <NTabPane :tab="`付款成本（${detail.payouts.length}）`" name="payouts">
                <NDataTable
                  :columns="payoutColumns"
                  :data="detail.payouts"
                  :row-key="row => row.id"
                  :pagination="false"
                  size="small"
                  :max-height="260"
                />
              </NTabPane>
              <NTabPane :tab="`材料领用（${detail.materials.length}）`" name="materials">
                <NDataTable
                  :columns="materialColumns"
                  :data="detail.materials"
                  :row-key="row => row.id"
                  :pagination="false"
                  size="small"
                  :max-height="260"
                />
              </NTabPane>
              <NTabPane :tab="`变更（${detail.changes.length}）`" name="changes">
                <NDataTable
                  :columns="changeColumns"
                  :data="detail.changes"
                  :row-key="row => row.id"
                  :pagination="false"
                  size="small"
                  :max-height="260"
                />
              </NTabPane>
            </NTabs>
          </NCard>

          <p class="mt-10px text-12px text-gray-400">
            成本按项目付款归集（材料款/人工/施工/设计/管理/销售/其他），公司公共支出不摊入项目。
            材料耗用金额 = 领用数量 × 材料单价，仅用于核对材料款付款是否合理，不重复计入成本。
            在建项目的毛利为预计值，成本尚未完全发生时仅供参考。
          </p>
        </template>
      </NSpin>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="handleClose">关闭</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
.metric-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(130px, 1fr));
  gap: 8px;
}

.metric {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 9px 11px;
  background-color: rgb(var(--base-text-color) / 0.03);
  border-radius: 8px;
}

.metric__label {
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.45);
}

.metric__value {
  overflow: hidden;
  font-size: 14px;
  font-variant-numeric: tabular-nums;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.8);
}

.metric__value--strong {
  font-weight: 600;
  color: rgb(var(--base-text-color));
}

.metric__value--success {
  color: rgb(var(--success-color));
}

.metric__value--warning {
  color: rgb(var(--warning-color));
}

.metric__value--error {
  color: rgb(var(--error-color));
}

.cost-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cost-row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
}

.cost-row__type {
  width: 58px;
  flex-shrink: 0;
  color: rgb(var(--base-text-color) / 0.65);
}

.cost-row__bar {
  flex: 1;
  height: 8px;
  overflow: hidden;
  background-color: rgb(var(--base-text-color) / 0.06);
  border-radius: 999px;
}

.cost-row__bar-inner {
  height: 100%;
  background: linear-gradient(90deg, rgb(var(--error-color)), rgb(var(--warning-color)));
  border-radius: 999px;
  transition: width 0.4s ease;
}

.cost-row__amount {
  width: 100px;
  flex-shrink: 0;
  font-variant-numeric: tabular-nums;
  text-align: right;
}

.cost-row__percent {
  width: 46px;
  flex-shrink: 0;
  color: rgb(var(--base-text-color) / 0.4);
  text-align: right;
}

.cost-total {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 8px;
  margin-top: 2px;
  font-size: 13px;
  border-top: 1px solid rgb(var(--base-text-color) / 0.08);
}

.cost-total strong {
  font-variant-numeric: tabular-nums;
  color: rgb(var(--error-color));
}

.material-compare {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  padding-top: 8px;
  margin-top: 2px;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.45);
  border-top: 1px dashed rgb(var(--base-text-color) / 0.08);
}
</style>
