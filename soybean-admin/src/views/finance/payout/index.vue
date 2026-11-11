<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import { fetchGetPayoutList } from '@/service/api/payout';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';
import { formatDate, formatMoney, formatMoneyInWan } from '../components/finance-meta';
import { useProjectOptions } from '../components/use-project-options';
import PayoutDetailDrawer from './components/payout-detail-drawer.vue';
import PayoutOperateDrawer from './components/payout-operate-drawer.vue';
import { PAYOUT_TYPE_OPTIONS, getPayoutTypeLabel, getPayoutTypeMeta } from './components/payout-meta';
import { usePayoutData } from './components/use-payout-data';

defineOptions({
  name: 'FinancePayout'
});

interface SearchParams {
  current: number;
  size: number;
  /** 关键词：项目名称 / 客户姓名 / 收款方 / 费用类型 / 备注 */
  keyword: string;
  /** 所属项目 */
  projectId: number | null;
  /** 费用类型 */
  payoutType: string | null;
  /** 付款来源 */
  source: 'purchase' | 'manual' | null;
}

/** 表格分页（只保留本页用到的字段，避免与 PaginationProps 的可选类型纠缠） */
interface PaginationState {
  page: number;
  pageSize: number;
  itemCount: number;
  showSizePicker: boolean;
  pageSizes: number[];
  prefix: (info: { itemCount?: number }) => string;
  onUpdatePage: (page: number) => void;
  onUpdatePageSize: (pageSize: number) => void;
}

/** 顶部概览卡片 */
interface StatCard {
  key: string;
  label: string;
  value: string;
  icon: string;
  color: string;
}

const appStore = useAppStore();

const isMobile = computed(() => appStore.isMobile);

const { loading, records, total, stats, loadList, loadStats, deletePayouts } = usePayoutData();

const { projectOptions, loadProjectOptions } = useProjectOptions();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  projectId: null,
  payoutType: null,
  source: null
});

/** 表格勾选的付款记录 */
const checkedRowKeys = ref<DataTableRowKey[]>([]);

const detailVisible = ref(false);
const detailId = ref<number | null>(null);

const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Payout.Item | null>(null);

/** 概览卡片数据 */
const statCards = computed<StatCard[]>(() => [
  {
    key: 'count',
    label: '付款笔数',
    value: String(stats.value.totalCount ?? 0),
    icon: 'mdi:receipt-text-outline',
    color: 'primary'
  },
  {
    key: 'amount',
    label: '付款总额',
    value: formatMoneyInWan(stats.value.totalAmount),
    icon: 'mdi:cash-minus',
    color: 'error'
  },
  {
    key: 'month',
    label: '本月付款',
    value: formatMoneyInWan(stats.value.monthAmount),
    icon: 'mdi:calendar-month-outline',
    color: 'info'
  },
  {
    key: 'project',
    label: '涉及项目',
    value: String(stats.value.projectCount ?? 0),
    icon: 'mdi:home-city-outline',
    color: 'warning'
  }
]);

/** 分页配置 */
const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 笔付款`,
  onUpdatePage(page) {
    pagination.page = page;

    void getData();
  },
  onUpdatePageSize(pageSize) {
    pagination.pageSize = pageSize;
    pagination.page = 1;

    void getData();
  }
});

const tablePagination = computed(() => ({
  ...pagination,
  pageSlot: isMobile.value ? 3 : 9,
  prefix: isMobile.value ? undefined : pagination.prefix
}));

const columns: DataTableColumns<Api.Payout.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'payoutDate',
    title: '付款日期',
    width: 130,
    fixed: 'left',
    sorter: (a, b) => (a.payoutDate > b.payoutDate ? 1 : -1),
    render: row => {
      const meta = getPayoutTypeMeta(row.payoutType);

      return h('div', { class: 'flex-y-center gap-8px' }, [
        h('div', {
          class: 'h-8px w-8px flex-shrink-0 rounded-full',
          style: { backgroundColor: `rgb(var(--${meta.color}-color))` }
        }),
        h('span', { class: 'font-mono' }, formatDate(row.payoutDate))
      ]);
    }
  },
  {
    key: 'projectName',
    title: '所属项目',
    minWidth: 170,
    ellipsis: { tooltip: true },
    render: row => row.projectName || h('span', { class: 'text-gray-400' }, '公司支出')
  },
  {
    key: 'vendorName',
    title: '供应商',
    width: 130,
    ellipsis: { tooltip: true },
    render: row => row.vendorName || h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'payoutType',
    title: '费用类型',
    width: 110,
    render: row => {
      const meta = getPayoutTypeMeta(row.payoutType);

      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => getPayoutTypeLabel(row.payoutType) });
    }
  },
  {
    key: 'amount',
    title: '付款金额',
    width: 140,
    align: 'right',
    sorter: (a, b) => a.amount - b.amount,
    render: row => h('span', { class: 'font-mono font-medium text-error' }, formatMoney(row.amount))
  },
  {
    key: 'fundAccountName',
    title: '付款账户',
    width: 150,
    ellipsis: { tooltip: true },
    render: row => row.fundAccountName || h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'voucherNo',
    title: '记账凭证',
    width: 150,
    render: row =>
      row.voucherNo
        ? h(NTag, { size: 'small', bordered: false, type: 'success' }, { default: () => row.voucherNo })
        : h(NTag, { size: 'small', bordered: false }, { default: () => '未记账' })
  },
  {
    key: 'remark',
    title: '备注',
    minWidth: 160,
    ellipsis: { tooltip: true },
    render: row => row.remark || h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'createdAt',
    title: '登记时间',
    width: 120,
    render: row => formatDate(row.createdAt)
  },
  {
    key: 'operate',
    title: '操作',
    width: 160,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-4px' }, [
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'primary',
            onClick: (event: MouseEvent) => {
              event.stopPropagation();
              handleDetail(row);
            }
          },
          { default: () => '查看' }
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            onClick: (event: MouseEvent) => {
              event.stopPropagation();
              handleEdit(row);
            }
          },
          { default: () => '编辑' }
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(row) },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'tiny',
                  quaternary: true,
                  type: 'error',
                  onClick: (event: MouseEvent) => event.stopPropagation()
                },
                { default: () => '删除' }
              ),
            default: () => '确认删除该条付款记录吗？对应的记账凭证会一起作废。'
          }
        )
      ])
  }
];

function rowProps(row: Api.Payout.Item) {
  return {
    style: 'cursor: pointer;',
    onClick: (event: MouseEvent) => {
      const target = event.target as HTMLElement | null;

      // 勾选列 / 展开列上的点击不打开详情
      if (target?.closest('.n-data-table-td--selection, .n-data-table-td--expand')) {
        return;
      }

      handleDetail(row);
    }
  };
}

/** 加载付款列表 */
async function getData() {
  searchParams.current = pagination.page;
  searchParams.size = pagination.pageSize;

  await loadList(searchParams);

  pagination.itemCount = total.value;
}

/** 刷新列表 + 概览统计 */
async function refreshAll() {
  await Promise.all([getData(), loadStats()]);
}

/** 查询 */
async function handleSearch() {
  pagination.page = 1;

  await getData();
}

/** 重置查询条件 */
async function handleReset() {
  searchParams.keyword = '';
  searchParams.projectId = null;
  searchParams.payoutType = null;
  searchParams.source = null;
  pagination.page = 1;

  await getData();
}

/** 刷新列表（同时刷新项目下拉，避免新增项目后选项过期） */
async function handleRefresh() {
  await Promise.all([refreshAll(), loadProjectOptions(true)]);
}

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.Payout.Item>[] = [
  { title: '付款日期', value: row => formatDate(row.payoutDate), width: 12 },
  { title: '所属项目', value: row => row.projectName || '公司支出', width: 18 },
  { title: '供应商', value: row => row.vendorName || '—', width: 12 },
  { title: '费用类型', value: row => getPayoutTypeLabel(row.payoutType), width: 10 },
  { title: '付款金额', value: row => formatMoney(row.amount), align: 'right' },
  { title: '付款账户', value: row => row.fundAccountName || '—', width: 12 },
  { title: '记账凭证', value: row => row.voucherNo || '未记账', width: 14 },
  { title: '备注', value: row => row.remark || '—', width: 16 },
  { title: '登记时间', value: row => formatDate(row.createdAt), width: 16 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const project = projectOptions.value.find(item => item.value === searchParams.projectId)?.label;
  const parts = [`项目：${project ?? '全部（含公司支出）'}`, `费用类型：${searchParams.payoutType || '全部'}`];

  if (searchParams.keyword.trim()) {
    parts.push(`关键词：${searchParams.keyword.trim()}`);
  }

  return parts.join('　|　');
}

/** 拉取当前筛选下的全部付款记录（分页聚合） */
async function collectPayoutRows() {
  return fetchAllPages<Api.Payout.Item>(async current => {
    const { data, error } = await fetchGetPayoutList({
      keyword: searchParams.keyword.trim() || null,
      projectId: searchParams.projectId,
      payoutType: searchParams.payoutType,
      source: searchParams.source,
      current,
      size: 100
    });

    if (error || !data) {
      throw new Error('数据加载失败');
    }

    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

function payoutSummary(rows: Api.Payout.Item[]) {
  return [
    { label: '付款笔数', value: `${rows.length} 笔` },
    { label: '合计金额', value: formatMoney(rows.reduce((sum, row) => sum + row.amount, 0)) }
  ];
}

async function handlePrint() {
  printing.value = true;

  try {
    const rows = await collectPayoutRows();
    await printTable({
      title: '付款单列表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: payoutSummary(rows)
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
    const rows = await collectPayoutRows();
    exportExcel({
      filename: '付款单列表',
      title: '付款单列表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: payoutSummary(rows)
    });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

/** 查看付款详情 */
function handleDetail(row: Api.Payout.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

/** 新增付款 */
function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

/** 编辑付款 */
function handleEdit(row: Api.Payout.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

/** 从详情抽屉进入编辑 */
function handleEditFromDetail(row: Api.Payout.Item) {
  detailVisible.value = false;

  handleEdit(row);
}

/** 删除付款 */
async function handleDelete(row: Api.Payout.Item) {
  const success = await deletePayouts([row.id]);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');

  await afterDeleted(1);
}

/** 批量删除付款 */
async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(key => Number(key));
  const success = await deletePayouts(ids);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');

  checkedRowKeys.value = [];

  await afterDeleted(ids.length);
}

/** 删除后刷新：当前页被删空时自动回到上一页 */
async function afterDeleted(deletedCount: number) {
  const maxPage = Math.max(1, Math.ceil((total.value - deletedCount) / pagination.pageSize));

  if (pagination.page > maxPage) {
    pagination.page = maxPage;
  }

  await refreshAll();
}

/** 新增 / 编辑提交成功后刷新 */
async function handleSubmitted() {
  if (operateType.value === 'add') {
    pagination.page = 1;
  }

  await refreshAll();
}

onMounted(async () => {
  await Promise.all([refreshAll(), loadProjectOptions(true)]);
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <!-- 概览统计 -->
      <div class="flex flex-wrap gap-12px">
        <NCard
          v-for="card in statCards"
          :key="card.key"
          :bordered="false"
          class="card-wrapper min-w-160px flex-1"
          size="small"
        >
          <div class="flex-y-center gap-12px">
            <div
              class="h-40px w-40px flex-center flex-shrink-0 rounded-10px text-20px"
              :style="{
                backgroundColor: `rgb(var(--${card.color}-color) / 0.12)`,
                color: `rgb(var(--${card.color}-color))`
              }"
            >
              <SvgIcon :icon="card.icon" />
            </div>
            <div class="min-w-0 flex-col">
              <span class="truncate text-12px text-gray-400">{{ card.label }}</span>
              <span class="truncate text-20px font-medium">{{ card.value }}</span>
            </div>
          </div>
        </NCard>
      </div>

      <!-- 搜索 / 筛选 / 操作按钮 -->
      <NCard :bordered="false" class="card-wrapper" size="small">
        <NSpace vertical :size="12">
          <div class="flex flex-wrap items-center gap-8px">
            <NInput
              v-model:value="searchParams.keyword"
              class="w-240px min-w-140px flex-shrink"
              placeholder="搜索项目 / 客户 / 收款方 / 费用类型 / 备注"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <SvgIcon icon="mdi:magnify" class="text-gray-400" />
              </template>
            </NInput>
            <NSelect
              v-model:value="searchParams.projectId"
              class="w-200px flex-shrink-0"
              :options="projectOptions"
              placeholder="所属项目"
              filterable
              clearable
            />
            <NSelect
              v-model:value="searchParams.payoutType"
              class="w-120px flex-shrink-0"
              :options="PAYOUT_TYPE_OPTIONS"
              placeholder="费用类型"
              clearable
            />
            <NSelect
              v-model:value="searchParams.source"
              class="w-140px flex-shrink-0"
              :options="[{ value: 'purchase', label: '采购付款' }, { value: 'manual', label: '其他付款' }]"
              placeholder="付款来源"
              clearable
            />
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
              <NButton size="small" type="primary" ghost @click="handleAdd">
                <template #icon>
                  <SvgIcon icon="mdi:cash-minus" />
                </template>
                新增付款
              </NButton>
              <NPopconfirm @positive-click="handleBatchDelete">
                <template #trigger>
                  <NButton size="small" ghost type="error" :disabled="!checkedRowKeys.length">
                    <template #icon>
                      <SvgIcon icon="mdi:delete-outline" />
                    </template>
                    批量删除
                  </NButton>
                </template>
                确认删除选中的 {{ checkedRowKeys.length }} 笔付款吗？对应凭证会一起作废。
              </NPopconfirm>
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
              <NButton size="small" @click="handleRefresh">
                <template #icon>
                  <SvgIcon icon="mdi:refresh" :class="{ 'animate-spin': loading }" />
                </template>
                刷新
              </NButton>
            </div>
          </div>

          <NDataTable
            v-model:checked-row-keys="checkedRowKeys"
            :columns="columns"
            :data="records"
            :loading="loading"
            :pagination="tablePagination"
            :row-key="row => row.id"
            :row-props="rowProps"
            :scroll-x="1580"
            remote
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </NSpace>

    <!-- 新增 / 编辑 -->
    <PayoutOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      :default-project-id="searchParams.projectId"
      @submitted="handleSubmitted"
    />

    <!-- 付款详情 -->
    <PayoutDetailDrawer v-model:visible="detailVisible" :payout-id="detailId" @edit="handleEditFromDetail" />
  </div>
</template>

<style scoped></style>
