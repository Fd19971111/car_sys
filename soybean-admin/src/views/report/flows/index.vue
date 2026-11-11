<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NTag } from 'naive-ui';
import dayjs from 'dayjs';
import { fetchGetFinFundAccounts, fetchGetFinReportFlows } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';
import ReportToolbar from '../components/report-toolbar.vue';
import { FLOW_DIRECTION_META } from '../components/report-meta';
import { useReportData } from '../components/use-report-data';
import { useReportYear } from '../components/use-report-year';

defineOptions({
  name: 'ReportFlows'
});

interface SearchParams {
  /** 月份，0 表示全年 */
  month: number;
  /** 收支方向，空字符串表示全部 */
  directionCode: '' | Api.FinReport.FlowDirection;
  /** 资金账户ID，0 表示全部 */
  fundAccountId: number;
  /** 备注 / 科目 / 账户关键字 */
  keyword: string;
}

const appStore = useAppStore();

const isMobile = computed(() => appStore.isMobile);

const { year } = useReportYear();

const { flowsLoading, flows, flowTotal, flowSummary, loadFlows } = useReportData();

const searchParams = reactive<SearchParams>({
  month: 0,
  directionCode: '',
  fundAccountId: 0,
  keyword: ''
});

const fundAccountOptions = ref<CommonType.Option<number, string>[]>([]);

const monthOptions = [
  { label: '全年', value: 0 },
  ...Array.from({ length: 12 }, (_, index) => ({ label: `${index + 1} 月`, value: index + 1 }))
];

const directionOptions = [
  { label: '全部', value: '' },
  { label: '只看收入', value: 'INCOME' },
  { label: '只看支出', value: 'EXPENSE' }
];

/** 分页配置 */
const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: (info: { itemCount?: number }) => `共 ${info.itemCount ?? 0} 条流水`,
  onUpdatePage(page: number) {
    pagination.page = page;

    void load();
  },
  onUpdatePageSize(pageSize: number) {
    pagination.pageSize = pageSize;
    pagination.page = 1;

    void load();
  }
});

const tablePagination = computed(() => ({
  ...pagination,
  pageSlot: isMobile.value ? 3 : 9,
  prefix: isMobile.value ? undefined : pagination.prefix
}));

const summaryCards = computed(() => [
  {
    key: 'income',
    label: '筛选范围收入',
    value: formatMoney(flowSummary.value.income),
    color: 'success',
    icon: 'mdi:arrow-down-left'
  },
  {
    key: 'expense',
    label: '筛选范围支出',
    value: formatMoney(flowSummary.value.expense),
    color: 'error',
    icon: 'mdi:arrow-up-right'
  },
  {
    key: 'net',
    label: '净额',
    value: formatMoney(flowSummary.value.net),
    color: 'primary',
    icon: 'mdi:scale-balance'
  }
]);

const columns: DataTableColumns<Api.FinReport.FlowItem> = [
  {
    key: 'occurredAt',
    title: '发生时间',
    width: 150,
    render: row => h('span', { class: 'font-mono' }, dayjs(row.occurredAt).format('YYYY-MM-DD HH:mm'))
  },
  {
    key: 'directionCode',
    title: '收支',
    width: 88,
    render: row => {
      const meta = FLOW_DIRECTION_META[row.directionCode] ?? FLOW_DIRECTION_META.INCOME;

      return h(
        NTag,
        { size: 'small', bordered: false, type: row.directionCode === 'INCOME' ? 'success' : 'error' },
        { default: () => meta.label }
      );
    }
  },
  {
    key: 'amount',
    title: '金额',
    width: 140,
    align: 'right',
    render: row =>
      h(
        'span',
        { class: `font-mono font-medium ${row.directionCode === 'INCOME' ? 'text-success' : 'text-error'}` },
        `${row.directionCode === 'INCOME' ? '+' : '-'}${formatMoney(row.amount)}`
      )
  },
  {
    key: 'accountName',
    title: '会计科目',
    width: 140,
    ellipsis: { tooltip: true },
    render: row => row.accountName || '-'
  },
  {
    key: 'fundAccountName',
    title: '资金账户',
    width: 160,
    ellipsis: { tooltip: true },
    render: row => row.fundAccountName || row.accountName || '-'
  },
  {
    key: 'remark',
    title: '备注',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render: row => row.remark || h('span', { class: 'text-gray-400' }, '—')
  }
];

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.FinReport.FlowItem>[] = [
  { title: '发生时间', value: row => dayjs(row.occurredAt).format('YYYY-MM-DD HH:mm'), width: 18 },
  {
    title: '收支',
    value: row => (FLOW_DIRECTION_META[row.directionCode] ?? FLOW_DIRECTION_META.INCOME).label,
    align: 'center',
    width: 8
  },
  {
    title: '金额',
    value: row => `${row.directionCode === 'INCOME' ? '+' : '-'}${formatMoney(row.amount)}`,
    align: 'right',
    width: 16
  },
  { title: '会计科目', value: row => row.accountName || '-', width: 16 },
  { title: '资金账户', value: row => row.fundAccountName || row.accountName || '-', width: 18 },
  { title: '备注', value: row => row.remark || '—', width: 32 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const parts = [
    `年度：${year.value} 年`,
    `月份：${searchParams.month ? `${searchParams.month} 月` : '全年'}`,
    `收支：${searchParams.directionCode === 'INCOME' ? '只看收入' : searchParams.directionCode === 'EXPENSE' ? '只看支出' : '全部'}`
  ];

  if (searchParams.fundAccountId) {
    const account = fundAccountOptions.value.find(item => item.value === searchParams.fundAccountId)?.label;
    parts.push(`资金账户：${account ?? '已选账户'}`);
  }

  if (searchParams.keyword.trim()) {
    parts.push(`关键词：${searchParams.keyword.trim()}`);
  }

  return parts.join('　|　');
}

/** 拉取当前筛选下的全部流水（分页聚合，打印 / 导出用） */
async function collectFlowRows() {
  return fetchAllPages<Api.FinReport.FlowItem>(async current => {
    const { data, error } = await fetchGetFinReportFlows({
      year: year.value,
      month: searchParams.month || 0,
      directionCode: searchParams.directionCode || null,
      fundAccountId: searchParams.fundAccountId || 0,
      keyword: searchParams.keyword,
      current,
      size: 100
    });

    if (error || !data) {
      throw new Error('数据加载失败');
    }

    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

/** 合计行：收入 / 支出 / 净额 */
function flowPrintSummary(rows: Api.FinReport.FlowItem[]) {
  const income = rows.filter(row => row.directionCode === 'INCOME').reduce((sum, row) => sum + row.amount, 0);
  const expense = rows.filter(row => row.directionCode === 'EXPENSE').reduce((sum, row) => sum + row.amount, 0);

  return [
    { label: '流水笔数', value: `${rows.length} 笔` },
    { label: '收入合计', value: formatMoney(income) },
    { label: '支出合计', value: formatMoney(expense) },
    { label: '净额', value: formatMoney(income - expense) }
  ];
}

async function handlePrint() {
  printing.value = true;

  try {
    const rows = await collectFlowRows();
    await printTable({
      title: '资金流水表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: flowPrintSummary(rows)
    });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    printing.value = false;
  }
}

async function handleExport() {
  exporting.value = true;

  try {
    const rows = await collectFlowRows();
    exportExcel({
      filename: `资金流水表-${year.value}年`,
      title: '资金流水表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: flowPrintSummary(rows)
    });
    window.$message?.success(`已导出 ${rows.length} 条流水`);
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

/** 加载资金账户下拉 */
async function loadFundAccounts() {
  const { data, error } = await fetchGetFinFundAccounts();

  if (!error && data) {
    fundAccountOptions.value = (data.records ?? [])
      .filter(item => item.enabled)
      .map(item => ({ label: item.bankName ? `${item.name} · ${item.bankName}` : item.name, value: item.id }));
  }
}

/** 加载资金流水 */
async function load() {
  await loadFlows({
    year: year.value,
    month: searchParams.month || 0,
    directionCode: searchParams.directionCode || null,
    fundAccountId: searchParams.fundAccountId || 0,
    keyword: searchParams.keyword,
    current: pagination.page,
    size: pagination.pageSize
  });

  pagination.itemCount = flowTotal.value;
}

/** 查询 */
async function handleSearch() {
  pagination.page = 1;

  await load();
}

/** 重置 */
async function handleReset() {
  searchParams.month = 0;
  searchParams.directionCode = '';
  searchParams.fundAccountId = 0;
  searchParams.keyword = '';
  pagination.page = 1;

  await load();
}

watch(year, () => {
  pagination.page = 1;

  void load();
});

onMounted(async () => {
  await Promise.all([loadFundAccounts(), load()]);
});
</script>

<template>
  <div class="report-page">
    <ReportToolbar
      v-model:year="year"
      hint="资金流水由记账凭证自动生成，可按月份、收支类型、资金账户和关键字筛选"
      :loading="flowsLoading"
      @refresh="load"
    />

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
          </div>
        </div>
      </NCard>
    </div>

    <NCard :bordered="false" class="filter-card" size="small">
      <NSpace vertical :size="12">
        <div class="flex flex-wrap items-center gap-8px">
          <NSelect v-model:value="searchParams.month" class="w-110px flex-shrink-0" :options="monthOptions" />
          <NSelect v-model:value="searchParams.directionCode" class="w-130px flex-shrink-0" :options="directionOptions" />
          <NSelect
            v-model:value="searchParams.fundAccountId"
            class="w-190px flex-shrink-0"
            :options="fundAccountOptions"
            placeholder="资金账户"
            clearable
            filterable
          />
          <NInput
            v-model:value="searchParams.keyword"
            class="w-220px min-w-140px flex-shrink"
            placeholder="搜索备注 / 科目 / 账户"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <SvgIcon icon="mdi:magnify" class="text-gray-400" />
            </template>
          </NInput>
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
          <div class="ml-auto flex flex-shrink-0 items-center gap-8px">
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
            <NButton size="small" :loading="flowsLoading" @click="load">
              <template #icon>
                <SvgIcon icon="mdi:refresh" />
              </template>
              刷新
            </NButton>
          </div>
        </div>

        <NDataTable
          :columns="columns"
          :data="flows"
          :loading="flowsLoading"
          :pagination="tablePagination"
          :row-key="row => row.id"
          :scroll-x="960"
          remote
          size="small"
          striped
        />
      </NSpace>
    </NCard>
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

.summary-card,
.filter-card {
  border-radius: 12px;
}

@media (max-width: 768px) {
  .report-page {
    padding: 12px;
  }
}
</style>
