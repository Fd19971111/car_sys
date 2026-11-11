<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NTag } from 'naive-ui';
import { fetchGetProjectProfitList } from '@/service/api';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { exportExcel, printTable, type PrintExportColumn } from '@/utils/print-export';
import {
  PROJECT_STATUS_OPTIONS,
  getProjectStatusMeta
} from '@/views/project/project-list/components/project-meta';
import ProjectProfitDetailDrawer from './components/profit-detail-drawer.vue';

defineOptions({
  name: 'ReportProfit'
});

interface SearchParams {
  keyword: string;
  status: Api.Project.Status | null;
  onlyLoss: boolean;
}

const searchParams = reactive<SearchParams>({
  keyword: '',
  status: null,
  onlyLoss: false
});

const loading = ref(false);
const records = ref<Api.FinReport.ProjectProfitItem[]>([]);
const summary = ref<Api.FinReport.ProjectProfitSummary | null>(null);

const detailVisible = ref(false);
const detailProjectId = ref<number | null>(null);

const summaryCards = computed(() => {  const s = summary.value;
  if (!s) return [];

  return [
    {
      key: 'settlement',
      label: '结算总额',
      value: formatMoney(s.settlementTotal),
      hint: '合同 + 变更净额',
      color: 'primary',
      icon: 'mdi:file-document-outline'
    },
    {
      key: 'received',
      label: '已收款',
      value: formatMoney(s.receivedTotal),
      hint: `待收 ${formatMoney(s.unpaidTotal)}`,
      color: 'success',
      icon: 'mdi:cash-check'
    },
    {
      key: 'cost',
      label: '成本合计',
      value: formatMoney(s.costTotal),
      hint: '按项目付款归集',
      color: 'error',
      icon: 'mdi:cash-minus'
    },
    {
      key: 'profit',
      label: '预计毛利',
      value: formatMoney(s.grossProfit),
      hint: `整体毛利率 ${s.grossRate}%`,
      color: s.grossProfit >= 0 ? 'success' : 'error',
      icon: 'mdi:chart-line'
    },
    {
      key: 'loss',
      label: '亏损项目',
      value: String(s.lossCount),
      hint: s.lossCount > 0 ? '需要重点复核' : '暂无亏损项目',
      color: s.lossCount > 0 ? 'warning' : 'success',
      icon: 'mdi:alert-circle-outline'
    }
  ];
});

/** 毛利率标签颜色 */
function grossTagType(rate: number) {
  if (rate < 0) return 'error';
  if (rate < 10) return 'warning';

  return 'success';
}

/** 成本占结算比：太低说明成本还没录完，毛利会偏高 */
const costRate = computed(() => {
  const s = summary.value;
  if (!s || s.settlementTotal <= 0) return 0;

  return Math.round((s.costTotal / s.settlementTotal) * 1000) / 10;
});

/** 现金结余合计 = 已收款 - 成本 */
const cashProfitTotal = computed(() => {
  const s = summary.value;
  if (!s) return 0;

  return Math.round((s.receivedTotal - s.costTotal) * 100) / 100;
});

const columns: DataTableColumns<Api.FinReport.ProjectProfitItem> = [
  {
    key: 'projectName',
    title: '项目',
    minWidth: 190,
    fixed: 'left',
    ellipsis: { tooltip: true },
    render: row =>
      h('div', { class: 'flex-col gap-1px' }, [
        h('span', { class: 'font-medium' }, row.projectName),
        h('span', { class: 'text-11px text-gray-400' }, row.customerName || '未关联客户')
      ])
  },
  {
    key: 'status',
    title: '状态',
    width: 88,
    render: row => {
      const meta = getProjectStatusMeta(row.status);

      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'contractAmount',
    title: '合同金额',
    width: 120,
    align: 'right',
    sorter: (a, b) => a.contractAmount - b.contractAmount,
    render: row => h('span', { class: 'font-mono' }, formatMoney(row.contractAmount))
  },
  {
    key: 'changeNet',
    title: '变更净额',
    width: 110,
    align: 'right',
    sorter: (a, b) => a.changeNet - b.changeNet,
    render: row =>
      h(
        'span',
        { class: `font-mono ${row.changeNet < 0 ? 'text-warning' : row.changeNet > 0 ? 'text-success' : 'text-gray-400'}` },
        `${row.changeNet >= 0 ? '+' : ''}${formatMoney(row.changeNet)}`
      )
  },
  {
    key: 'settlementAmount',
    title: '结算金额',
    width: 125,
    align: 'right',
    sorter: (a, b) => a.settlementAmount - b.settlementAmount,
    render: row => h('span', { class: 'font-mono font-medium' }, formatMoney(row.settlementAmount))
  },
  {
    key: 'receivedAmount',
    title: '已收款',
    width: 120,
    align: 'right',
    sorter: (a, b) => a.receivedAmount - b.receivedAmount,
    render: row => h('span', { class: 'font-mono text-success' }, formatMoney(row.receivedAmount))
  },
  {
    key: 'receivedRate',
    title: '回款率',
    width: 84,
    align: 'right',
    sorter: (a, b) => a.receivedRate - b.receivedRate,
    render: row => h('span', { class: 'font-mono text-12px text-gray-500' }, `${row.receivedRate}%`)
  },
  {
    key: 'costTotal',
    title: '成本合计',
    width: 120,
    align: 'right',
    sorter: (a, b) => a.costTotal - b.costTotal,
    render: row => h('span', { class: 'font-mono text-error' }, formatMoney(row.costTotal))
  },
  {
    key: 'materialUsed',
    title: '材料耗用（参考）',
    width: 130,
    align: 'right',
    sorter: (a, b) => a.materialUsed - b.materialUsed,
    render: row => h('span', { class: 'font-mono text-12px text-gray-400' }, formatMoney(row.materialUsed))
  },
  {
    key: 'grossProfit',
    title: '预计毛利',
    width: 130,
    align: 'right',
    sorter: (a, b) => a.grossProfit - b.grossProfit,
    render: row =>
      h(
        'span',
        { class: `font-mono font-medium ${row.grossProfit >= 0 ? 'text-success' : 'text-error'}` },
        formatMoney(row.grossProfit)
      )
  },
  {
    key: 'grossRate',
    title: '毛利率',
    width: 90,
    align: 'right',
    sorter: (a, b) => a.grossRate - b.grossRate,
    render: row =>
      h(
        NTag,
        { size: 'small', bordered: false, type: grossTagType(row.grossRate) },
        { default: () => `${row.grossRate}%` }
      )
  },
  {
    key: 'operate',
    title: '操作',
    width: 90,
    align: 'center',
    fixed: 'right',
    render: row =>
      h(
        NButton,
        {
          size: 'tiny',
          quaternary: true,
          type: 'primary',
          onClick: (event: MouseEvent) => {
            event.stopPropagation();
            openDetail(row);
          }
        },
        { default: () => '成本明细' }
      )
  }
];

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.FinReport.ProjectProfitItem>[] = [
  { title: '项目名称', value: row => row.projectName, width: 24 },
  { title: '客户', value: row => row.customerName || '未关联客户', width: 12 },
  { title: '状态', value: row => getProjectStatusMeta(row.status).label, align: 'center', width: 10 },
  { title: '合同金额', value: row => formatMoney(row.contractAmount), align: 'right', width: 16 },
  {
    title: '变更净额',
    value: row => `${row.changeNet >= 0 ? '+' : ''}${formatMoney(row.changeNet)}`,
    align: 'right',
    width: 16
  },
  { title: '结算金额', value: row => formatMoney(row.settlementAmount), align: 'right', width: 16 },
  { title: '已收款', value: row => formatMoney(row.receivedAmount), align: 'right', width: 16 },
  { title: '回款率', value: row => `${row.receivedRate}%`, align: 'right', width: 9 },
  { title: '成本合计', value: row => formatMoney(row.costTotal), align: 'right', width: 16 },
  { title: '预计毛利', value: row => formatMoney(row.grossProfit), align: 'right', width: 16 },
  { title: '毛利率', value: row => `${row.grossRate}%`, align: 'right', width: 9 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const status = PROJECT_STATUS_OPTIONS.find(item => item.value === searchParams.status)?.label;
  const parts = ['项目成本与毛利', `项目状态：${status ?? '全部'}`];

  if (searchParams.onlyLoss) {
    parts.push('只看亏损');
  }

  if (searchParams.keyword.trim()) {
    parts.push(`关键词：${searchParams.keyword.trim()}`);
  }

  return parts.join('　|　');
}

/** 合计行：优先用接口汇总，没有汇总时按当前行现算 */
function profitPrintSummary(rows: Api.FinReport.ProjectProfitItem[]) {
  const s = summary.value;

  return [
    { label: '项目数', value: `${s?.projectCount ?? rows.length} 个` },
    { label: '结算总额', value: formatMoney(s?.settlementTotal ?? rows.reduce((sum, row) => sum + row.settlementAmount, 0)) },
    { label: '已收款', value: formatMoney(s?.receivedTotal ?? rows.reduce((sum, row) => sum + row.receivedAmount, 0)) },
    { label: '待收合计', value: formatMoney(s?.unpaidTotal ?? rows.reduce((sum, row) => sum + row.unpaidAmount, 0)) },
    { label: '成本合计', value: formatMoney(s?.costTotal ?? rows.reduce((sum, row) => sum + row.costTotal, 0)) },
    { label: '预计毛利', value: formatMoney(s?.grossProfit ?? rows.reduce((sum, row) => sum + row.grossProfit, 0)) },
    { label: '整体毛利率', value: `${s?.grossRate ?? 0}%` },
    { label: '亏损项目', value: `${s?.lossCount ?? rows.filter(row => row.grossProfit < 0).length} 个` }
  ];
}

async function handlePrint() {
  printing.value = true;

  try {
    await printTable({
      title: '项目成本与毛利表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows: records.value,
      summary: profitPrintSummary(records.value),
      note: '在建项目成本尚未完全发生，预计毛利偏高属正常现象，项目完工后更准确。'
    });
  } catch {
    window.$message?.error('打印失败，请稍后重试');
  } finally {
    printing.value = false;
  }
}

function handleExport() {
  exporting.value = true;

  try {
    exportExcel({
      filename: '项目成本与毛利表',
      title: '项目成本与毛利表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows: records.value,
      summary: profitPrintSummary(records.value),
      note: '在建项目成本尚未完全发生，预计毛利偏高属正常现象，项目完工后更准确。'
    });
    window.$message?.success(`已导出 ${records.value.length} 个项目`);
  } catch {
    window.$message?.error('导出失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

function rowProps(row: Api.FinReport.ProjectProfitItem) {
  return {
    style: 'cursor: pointer;',
    onClick: () => openDetail(row)
  };
}

async function load() {
  loading.value = true;

  try {
    const { data, error } = await fetchGetProjectProfitList({
      keyword: searchParams.keyword?.trim() || null,
      status: searchParams.status,
      onlyLoss: searchParams.onlyLoss
    });

    if (!error && data) {
      records.value = data.records ?? [];
      summary.value = data.summary ?? null;
    }
  } finally {
    loading.value = false;
  }
}

async function handleSearch() {
  await load();
}

async function handleReset() {
  searchParams.keyword = '';
  searchParams.status = null;
  searchParams.onlyLoss = false;

  await load();
}

function openDetail(row: Api.FinReport.ProjectProfitItem) {
  detailProjectId.value = row.projectId;
  detailVisible.value = true;
}

onMounted(load);
</script>

<template>
  <div class="report-page">
    <NCard :bordered="false" class="panel" size="small">
      <div class="flex flex-wrap items-start justify-between gap-12px">
        <div class="min-w-0 flex-col gap-4px">
          <span class="text-16px font-medium">项目成本与毛利</span>
          <span class="text-12px text-gray-400">
            按项目归集收款、付款成本、材料耗用与变更，算出每个项目的预计毛利和毛利率。在建项目成本尚未完全发生，
            预计毛利偏高属于正常现象，项目完工后更准确；点任意一行看成本明细。
          </span>
        </div>
        <div class="flex flex-shrink-0 items-center gap-8px">
          <NButton size="small" :loading="printing" @click="handlePrint">
            <template #icon>
              <SvgIcon icon="mdi:printer-outline" />
            </template>
            打印
          </NButton>
          <NButton size="small" :loading="exporting" @click="handleExport">
            <template #icon>
              <SvgIcon icon="mdi:microsoft-excel" />
            </template>
            导出
          </NButton>
          <NButton size="small" :loading="loading" @click="load">
            <template #icon>
              <SvgIcon icon="mdi:refresh" />
            </template>
            刷新
          </NButton>
        </div>
      </div>
    </NCard>

    <div class="summary-grid">
      <NCard v-for="card in summaryCards" :key="card.key" :bordered="false" size="small" class="summary-card">
        <div class="flex-y-center gap-12px">
          <div
            class="h-38px w-38px flex-center flex-shrink-0 rounded-10px text-18px"
            :style="{
              backgroundColor: `rgb(var(--${card.color}-color) / 0.12)`,
              color: `rgb(var(--${card.color}-color))`
            }"
          >
            <SvgIcon :icon="card.icon" />
          </div>
          <div class="min-w-0 flex-col">
            <span class="truncate text-12px text-gray-400">{{ card.label }}</span>
            <span class="truncate text-17px font-medium">{{ card.value }}</span>
            <span class="truncate text-11px text-gray-400">{{ card.hint }}</span>
          </div>
        </div>
      </NCard>
    </div>

    <div v-if="summary" class="flex flex-wrap items-center gap-16px px-2px text-12px text-gray-500">
      <span>项目数 {{ summary.projectCount }}</span>
      <span>待收合计 {{ formatMoney(summary.unpaidTotal) }}</span>
      <span>现金结余 {{ formatMoney(cashProfitTotal) }}</span>
      <span>成本占结算 {{ costRate }}%</span>
      <span>材料耗用合计 {{ formatMoney(summary.materialUsedTotal) }}（参考）</span>
      <span>公司公共支出 {{ formatMoney(summary.companyCost) }}（未摊入项目成本）</span>
    </div>

    <NCard :bordered="false" class="panel" size="small">
      <NSpace vertical :size="12">
        <div class="flex flex-wrap items-center gap-8px">
          <NInput
            v-model:value="searchParams.keyword"
            class="w-230px min-w-140px flex-shrink"
            placeholder="搜索项目名称 / 客户姓名"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <SvgIcon icon="mdi:magnify" class="text-gray-400" />
            </template>
          </NInput>
          <NSelect
            v-model:value="searchParams.status"
            class="w-130px flex-shrink-0"
            :options="PROJECT_STATUS_OPTIONS"
            placeholder="项目状态"
            clearable
          />
          <div class="flex flex-shrink-0 items-center gap-6px">
            <NSwitch v-model:value="searchParams.onlyLoss" size="small" />
            <span class="text-12px text-gray-500">只看亏损</span>
          </div>
          <NButton class="flex-shrink-0" type="primary" ghost size="small" @click="handleSearch">
            <template #icon>
              <SvgIcon icon="mdi:magnify" />
            </template>
            查询
          </NButton>
          <NButton class="flex-shrink-0" size="small" @click="handleReset">
            <template #icon>
              <SvgIcon icon="mdi:backspace-outline" />
            </template>
            重置
          </NButton>
        </div>

        <NDataTable
          :columns="columns"
          :data="records"
          :loading="loading"
          :row-key="row => row.projectId"
          :row-props="rowProps"
          :pagination="{ pageSize: 20, showSizePicker: true, pageSizes: [10, 20, 50] }"
          :scroll-x="1500"
          size="small"
          striped
        />
      </NSpace>
    </NCard>

    <ProjectProfitDetailDrawer v-model:visible="detailVisible" :project-id="detailProjectId" />
  </div>
</template>

<style scoped>
.report-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(190px, 1fr));
  gap: 12px;
}

.panel,
.summary-card {
  border-radius: 12px;
}

@media (max-width: 768px) {
  .report-page {
    padding: 12px;
  }
}
</style>
