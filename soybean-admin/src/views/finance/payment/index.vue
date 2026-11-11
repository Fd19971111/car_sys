<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import { fetchGetPaymentList } from '@/service/api/payment';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';
import { formatDate, formatMoney, formatMoneyInWan } from '../components/finance-meta';
import { useProjectOptions } from '../components/use-project-options';
import PaymentDetailDrawer from './components/payment-detail-drawer.vue';
import PaymentOperateDrawer from './components/payment-operate-drawer.vue';
import { PAYMENT_TYPE_OPTIONS, getPaymentTypeMeta, getPaymentTypeLabel } from './components/payment-meta';
import { usePaymentData } from './components/use-payment-data';

defineOptions({
  name: 'FinancePayment'
});

interface SearchParams {
  current: number;
  size: number;
  /** 关键词：项目名称 / 客户姓名 / 款项类型 / 备注 */
  keyword: string;
  /** 所属项目 */
  projectId: number | null;
  /** 款项类型 */
  paymentType: string | null;
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
  /** 点击卡片后应用的款项类型筛选；为空表示「全部」 */
  paymentType: string | null;
}

const appStore = useAppStore();

const isMobile = computed(() => appStore.isMobile);

const { loading, records, total, stats, loadList, loadStats, deletePayments } = usePaymentData();

const { projectOptions, loadProjectOptions } = useProjectOptions();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  projectId: null,
  paymentType: null
});

/** 表格勾选的收款记录 */
const checkedRowKeys = ref<DataTableRowKey[]>([]);

const detailVisible = ref(false);
const detailId = ref<number | null>(null);

const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Payment.Item | null>(null);

/** 概览卡片数据 */
const statCards = computed<StatCard[]>(() => [
  {
    key: 'count',
    label: '收款笔数',
    value: String(stats.value.totalCount ?? 0),
    icon: 'mdi:receipt-text-outline',
    color: 'primary',
    paymentType: null
  },
  {
    key: 'amount',
    label: '收款总额',
    value: formatMoneyInWan(stats.value.totalAmount),
    icon: 'mdi:cash-multiple',
    color: 'success',
    paymentType: null
  },
  {
    key: 'month',
    label: '本月收款',
    value: formatMoneyInWan(stats.value.monthAmount),
    icon: 'mdi:calendar-month-outline',
    color: 'info',
    paymentType: null
  },
  {
    key: 'project',
    label: '涉及项目',
    value: String(stats.value.projectCount ?? 0),
    icon: 'mdi:home-city-outline',
    color: 'warning',
    paymentType: null
  }
]);

/** 分页配置 */
const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 笔收款`,
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

const columns: DataTableColumns<Api.Payment.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'paymentDate',
    title: '收款日期',
    width: 130,
    fixed: 'left',
    sorter: (a, b) => (a.paymentDate > b.paymentDate ? 1 : -1),
    render: row => {
      const meta = getPaymentTypeMeta(row.paymentType);

      return h('div', { class: 'flex-y-center gap-8px' }, [
        h('div', {
          class: 'h-8px w-8px flex-shrink-0 rounded-full',
          style: { backgroundColor: `rgb(var(--${meta.color}-color))` }
        }),
        h('span', { class: 'font-mono' }, formatDate(row.paymentDate))
      ]);
    }
  },
  {
    key: 'projectName',
    title: '所属项目',
    minWidth: 170,
    ellipsis: { tooltip: true },
    render: row => row.projectName || '-'
  },
  {
    key: 'customerName',
    title: '客户',
    width: 110,
    ellipsis: { tooltip: true },
    render: row => row.customerName || '-'
  },
  {
    key: 'paymentType',
    title: '款项类型',
    width: 110,
    render: row => {
      const meta = getPaymentTypeMeta(row.paymentType);

      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => getPaymentTypeLabel(row.paymentType) });
    }
  },
  {
    key: 'amount',
    title: '收款金额',
    width: 140,
    align: 'right',
    sorter: (a, b) => a.amount - b.amount,
    render: row => h('span', { class: 'font-mono font-medium text-success' }, formatMoney(row.amount))
  },
  {
    key: 'fundAccountName',
    title: '收款账户',
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
            default: () => '确认删除该条收款记录吗？'
          }
        )
      ])
  }
];

function rowProps(row: Api.Payment.Item) {
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

/** 加载收款列表 */
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
  searchParams.paymentType = null;
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
const printExportColumns: PrintExportColumn<Api.Payment.Item>[] = [
  { title: '收款日期', value: row => formatDate(row.paymentDate), width: 12 },
  { title: '所属项目', value: row => row.projectName || '-', width: 18 },
  { title: '客户', value: row => row.customerName || '-', width: 10 },
  { title: '款项类型', value: row => getPaymentTypeLabel(row.paymentType), width: 10 },
  { title: '收款金额', value: row => formatMoney(row.amount), align: 'right' },
  { title: '收款账户', value: row => row.fundAccountName || '—', width: 12 },
  { title: '记账凭证', value: row => row.voucherNo || '未记账', width: 14 },
  { title: '备注', value: row => row.remark || '—', width: 16 },
  { title: '登记时间', value: row => formatDate(row.createdAt), width: 16 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const project = projectOptions.value.find(item => item.value === searchParams.projectId)?.label;
  const parts = [`项目：${project ?? '全部项目'}`, `款项类型：${searchParams.paymentType || '全部'}`];

  if (searchParams.keyword.trim()) {
    parts.push(`关键词：${searchParams.keyword.trim()}`);
  }

  return parts.join('　|　');
}

/** 拉取当前筛选下的全部收款记录（分页聚合） */
async function collectPaymentRows() {
  return fetchAllPages<Api.Payment.Item>(async current => {
    const { data, error } = await fetchGetPaymentList({
      keyword: searchParams.keyword.trim() || null,
      projectId: searchParams.projectId,
      paymentType: searchParams.paymentType,
      current,
      size: 100
    });

    if (error || !data) {
      throw new Error('数据加载失败');
    }

    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

function paymentSummary(rows: Api.Payment.Item[]) {
  return [
    { label: '收款笔数', value: `${rows.length} 笔` },
    { label: '合计金额', value: formatMoney(rows.reduce((sum, row) => sum + row.amount, 0)) }
  ];
}

async function handlePrint() {
  printing.value = true;

  try {
    const rows = await collectPaymentRows();
    await printTable({
      title: '收款单列表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: paymentSummary(rows)
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
    const rows = await collectPaymentRows();
    exportExcel({
      filename: '收款单列表',
      title: '收款单列表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: paymentSummary(rows)
    });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

/** 点击概览卡片：按款项类型快速筛选（本页卡片暂不带筛选条件） */
async function handleStatCardClick(card: StatCard) {
  searchParams.paymentType = card.paymentType;
  pagination.page = 1;

  await getData();
}

/** 概览卡片是否处于选中态 */
function isStatCardActive(card: StatCard) {
  return card.paymentType !== null && searchParams.paymentType === card.paymentType;
}

/** 查看收款详情 */
function handleDetail(row: Api.Payment.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

/** 新增收款 */
function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

/** 编辑收款 */
function handleEdit(row: Api.Payment.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

/** 从详情抽屉进入编辑 */
function handleEditFromDetail(row: Api.Payment.Item) {
  detailVisible.value = false;

  handleEdit(row);
}

/** 删除收款 */
async function handleDelete(row: Api.Payment.Item) {
  const success = await deletePayments([row.id]);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');

  await afterDeleted(1);
}

/** 批量删除收款 */
async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(key => Number(key));
  const success = await deletePayments(ids);

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
          class="card-wrapper min-w-160px flex-1 border-2 border-transparent transition-all hover:shadow-md"
          :style="isStatCardActive(card) ? { borderColor: 'rgb(var(--primary-color))' } : undefined"
          size="small"
          @click="handleStatCardClick(card)"
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
              placeholder="搜索项目名称 / 客户姓名 / 款项类型 / 备注"
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
              v-model:value="searchParams.paymentType"
              class="w-120px flex-shrink-0"
              :options="PAYMENT_TYPE_OPTIONS"
              placeholder="款项类型"
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
                  <SvgIcon icon="mdi:cash-plus" />
                </template>
                新增收款
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
                确认删除选中的 {{ checkedRowKeys.length }} 笔收款吗？
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
            :scroll-x="1470"
            remote
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </NSpace>

    <!-- 新增 / 编辑 -->
    <PaymentOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      :default-project-id="searchParams.projectId"
      @submitted="handleSubmitted"
    />

    <!-- 收款详情 -->
    <PaymentDetailDrawer v-model:visible="detailVisible" :payment-id="detailId" @edit="handleEditFromDetail" />
  </div>
</template>

<style scoped></style>
