<script setup lang="ts">
import { computed, h, onMounted, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NTag } from 'naive-ui';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { exportExcel, printTable, type PrintExportColumn } from '@/utils/print-export';
import ReportToolbar from '../components/report-toolbar.vue';
import { CATEGORY_META, CATEGORY_ORDER, getCategoryMeta } from '../components/report-meta';
import { useReportData } from '../components/use-report-data';
import { useReportYear } from '../components/use-report-year';

defineOptions({
  name: 'ReportBalance'
});

const { year } = useReportYear();
const { overviewLoading, overview, loadOverview } = useReportData();

/** 类别筛选，null 表示全部 */
const categoryFilter = ref<Api.FinReport.CategoryCode | null>(null);

const categoryOptions = CATEGORY_ORDER.map(code => ({ label: CATEGORY_META[code].label, value: code }));

const balances = computed(() => overview.value?.accountBalances ?? []);

/** 各科目类别的余额合计 */
const categoryTotals = computed(() => {
  const totals: Record<Api.FinReport.CategoryCode, number> = {
    ASSET: 0,
    LIABILITY: 0,
    EQUITY: 0,
    REVENUE: 0,
    EXPENSE: 0
  };

  balances.value.forEach(item => {
    totals[item.categoryCode] = (totals[item.categoryCode] ?? 0) + item.balance;
  });

  return totals;
});

/** 按类别筛选后的科目 */
const filteredBalances = computed(() => {
  if (!categoryFilter.value) return balances.value;

  return balances.value.filter(item => item.categoryCode === categoryFilter.value);
});

/** 试算平衡：资产 + 支出 = 负债 + 权益 + 收入 */
const debitSide = computed(() => categoryTotals.value.ASSET + categoryTotals.value.EXPENSE);
const creditSide = computed(() => categoryTotals.value.LIABILITY + categoryTotals.value.EQUITY + categoryTotals.value.REVENUE);
const balanced = computed(() => Math.abs(debitSide.value - creditSide.value) < 0.011);

/** 类别卡片 */
const categoryCards = computed(() =>
  CATEGORY_ORDER.map(code => ({
    key: code,
    label: CATEGORY_META[code].label,
    plain: CATEGORY_META[code].plain,
    color: CATEGORY_META[code].color,
    bg: CATEGORY_META[code].bg,
    icon: CATEGORY_META[code].icon,
    value: formatMoney(categoryTotals.value[code])
  }))
);

const columns: DataTableColumns<Api.FinReport.AccountBalance> = [
  {
    key: 'categoryCode',
    title: '类别',
    width: 110,
    render: row => {
      const meta = getCategoryMeta(row.categoryCode);

      return h('span', { class: 'category-tag', style: { color: meta.color, backgroundColor: meta.bg } }, meta.label);
    }
  },
  {
    key: 'code',
    title: '科目编码',
    width: 110,
    render: row => h('span', { class: 'font-mono' }, row.code)
  },
  {
    key: 'name',
    title: '科目名称',
    minWidth: 180,
    ellipsis: { tooltip: true }
  },
  {
    key: 'balance',
    title: '余额',
    width: 160,
    align: 'right',
    sorter: (a, b) => a.balance - b.balance,
    render: row => h('span', { class: 'font-mono font-medium' }, formatMoney(row.balance))
  }
];

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.FinReport.AccountBalance>[] = [
  { title: '类别', value: row => getCategoryMeta(row.categoryCode).label, align: 'center', width: 12 },
  { title: '科目编码', value: row => row.code, width: 12 },
  { title: '科目名称', value: row => row.name, width: 26 },
  { title: '余额', value: row => formatMoney(row.balance), align: 'right', width: 18 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const category = categoryFilter.value ? CATEGORY_META[categoryFilter.value].label : '全部类别';

  return [`年度：${year.value} 年`, `类别：${category}`, `科目数：${filteredBalances.value.length}`].join('　|　');
}

/** 合计行：各类别余额合计 */
function balanceSummary() {
  return CATEGORY_ORDER.map(code => ({
    label: `${CATEGORY_META[code].label}合计`,
    value: formatMoney(categoryTotals.value[code])
  }));
}

/** 打印底部的试算平衡说明 */
function balanceNote() {
  return `试算平衡：资产 + 支出 = ${formatMoney(debitSide.value)}；负债 + 权益 + 收入 = ${formatMoney(creditSide.value)}（${balanced.value ? '账面平衡' : '试算不平衡，请检查凭证'}）`;
}

async function handlePrint() {
  printing.value = true;

  try {
    await printTable({
      title: '科目余额表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows: filteredBalances.value,
      summary: balanceSummary(),
      note: balanceNote()
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
      filename: `科目余额表-${year.value}年`,
      title: '科目余额表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows: filteredBalances.value,
      summary: balanceSummary(),
      note: balanceNote()
    });
    window.$message?.success(`已导出 ${filteredBalances.value.length} 个科目`);
  } catch {
    window.$message?.error('导出失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

async function load() {
  await loadOverview(year.value);
}

watch(year, load);

onMounted(load);
</script>

<template>
  <div class="report-page">
    <ReportToolbar
      v-model:year="year"
      hint="科目余额 = 期初余额 + 已记账发生额，按借贷方向展示正常余额；只列有余额的科目"
      :loading="overviewLoading"
      @refresh="load"
    />

    <NAlert v-if="overview && !overview.initialized" type="warning" :bordered="false">
      财务还没有完成初始化，暂时无法统计科目余额。请先完成财务基础配置。
    </NAlert>

    <template v-else>
      <div class="category-grid">
        <NCard v-for="card in categoryCards" :key="card.key" :bordered="false" size="small" class="category-card">
          <div class="flex-y-center gap-12px">
            <div
              class="h-38px w-38px flex-center flex-shrink-0 rounded-10px text-18px"
              :style="{ color: card.color, backgroundColor: card.bg }"
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

      <NCard :bordered="false" class="panel" size="small">
        <div class="flex flex-wrap items-center justify-between gap-12px">
          <div class="flex-col gap-4px">
            <span class="text-14px font-medium">试算平衡</span>
            <span class="text-12px text-gray-400">
              资产 + 支出 = {{ formatMoney(debitSide) }}；负债 + 权益 + 收入 = {{ formatMoney(creditSide) }}
            </span>
          </div>
          <NTag :type="balanced ? 'success' : 'error'" :bordered="false" size="small">
            {{ balanced ? '账面平衡' : '试算不平衡，请检查凭证' }}
          </NTag>
        </div>
      </NCard>

      <NCard :bordered="false" class="panel" size="small">
        <NSpace vertical :size="12">
          <div class="flex flex-wrap items-center gap-8px">
            <NSelect
              v-model:value="categoryFilter"
              class="w-150px flex-shrink-0"
              :options="categoryOptions"
              placeholder="全部类别"
              clearable
            />
            <span class="text-12px text-gray-400">共 {{ filteredBalances.length }} 个科目有余额</span>
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
              <NButton size="small" :loading="overviewLoading" @click="load">
                <template #icon>
                  <SvgIcon icon="mdi:refresh" />
                </template>
                刷新
              </NButton>
            </div>
          </div>

          <NDataTable
            :columns="columns"
            :data="filteredBalances"
            :loading="overviewLoading"
            :row-key="row => row.code"
            :pagination="false"
            :scroll-x="660"
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </template>
  </div>
</template>

<style scoped>
.report-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}

.category-card,
.panel {
  border-radius: 12px;
}

.category-tag {
  display: inline-block;
  padding: 1px 8px;
  font-size: 12px;
  border-radius: 6px;
}

@media (max-width: 768px) {
  .report-page {
    padding: 12px;
  }
}
</style>
