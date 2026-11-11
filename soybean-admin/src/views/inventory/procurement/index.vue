<script setup lang="ts">
import { computed, h, onMounted, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NTag } from 'naive-ui';
import { fetchGetProjectSelectOptions, fetchGetPurchaseOrderList } from '@/service/api';
import { useAuthStore } from '@/store/modules/auth';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';
import StockBackfillDrawer from '../components/stock-backfill-drawer.vue';
import PurchaseCreateDrawer from './components/purchase-create-drawer.vue';

defineOptions({ name: 'InventoryProcurement' });

const authStore = useAuthStore();
const currentUserId = computed(() => authStore.userInfo.userId);

const projectOptions = ref<Api.Project.SelectOption[]>([]);
const projectId = ref<number | null>(null);
const projectLoading = ref(false);
const orders = ref<Api.ProjectData.PurchaseOrder[]>([]);
const loading = ref(false);
const drawerVisible = ref(false);

/** 采购单入库抽屉 */
const stockVisible = ref(false);
const stockRow = ref<Api.ProjectData.PurchaseOrder | null>(null);

/** 只有本人创建（历史数据不限）且未入库的采购单才能补入库 */
function canStock(row: Api.ProjectData.PurchaseOrder) {
  return !row.stockedAt && (!row.createdBy || row.createdBy === currentUserId.value);
}

function openStock(row: Api.ProjectData.PurchaseOrder) {
  stockRow.value = row;
  stockVisible.value = true;
}

const columns: DataTableColumns<Api.ProjectData.PurchaseOrder> = [
  { key: 'orderNo', title: '采购单号', width: 190 },
  {
    key: 'materialName',
    title: '材料',
    minWidth: 180,
    render: row =>
      h('div', { class: 'flex-col gap-2px' }, [
        h('span', { class: 'font-medium' }, row.materialName),
        h('span', { class: 'text-12px text-gray-400' }, row.spec || '无规格')
      ])
  },
  { key: 'quantity', title: '采购数量', width: 100, align: 'right', render: row => `${row.quantity} ${row.unit}` },
  { key: 'unitPrice', title: '采购单价', width: 120, align: 'right', render: row => formatMoney(row.unitPrice) },
  {
    key: 'taxMode',
    title: '税务方式',
    width: 120,
    render: row => ({ NO_INVOICE: '无票/不计税', TAX_INCLUDED: `含税 ${row.taxRate}%`, TAX_EXCLUDED: `未税 ${row.taxRate}%` })[row.taxMode]
  },
  { key: 'taxAmount', title: '税额', width: 110, align: 'right', render: row => formatMoney(row.taxAmount) },
  {
    key: 'amount',
    title: '采购金额',
    width: 130,
    align: 'right',
    render: row => h('span', { class: 'font-mono font-medium' }, formatMoney(row.amount))
  },
  { key: 'supplier', title: '供应商', minWidth: 140, render: row => row.supplier || '-' },
  {
    key: 'warehouseName',
    title: '仓库',
    width: 110,
    ellipsis: { tooltip: true },
    render: row =>
      row.warehouseName
        ? h(NTag, { size: 'small', bordered: false, round: true }, { default: () => row.warehouseName })
        : '-'
  },
  { key: 'createdAt', title: '创建时间', width: 170 },
  {
    key: 'status',
    title: '状态',
    width: 100,
    render: row =>
      h(
        NTag,
        {
          size: 'small',
          bordered: false,
          type: row.status === 'CONFIRMED' ? 'warning' : row.status === 'PAID' ? 'success' : 'default'
        },
        { default: () => ({ DRAFT: '草稿', CONFIRMED: '待付款', PAID: '已付款' })[row.status] ?? row.status }
      )
  },
  {
    key: 'stockStatus',
    title: '出入库',
    width: 100,
    render: row =>
      h(
        NTag,
        { size: 'small', bordered: false, type: row.stockedAt ? 'success' : 'warning' },
        { default: () => (row.stockedAt ? '已入库' : '未入库') }
      )
  },
  {
    key: 'operate',
    title: '操作',
    width: 110,
    align: 'center',
    fixed: 'right',
    render: row => {
      if (row.stockedAt) {
        return h('span', { class: 'text-12px text-gray-400' }, '已完成');
      }

      if (!canStock(row)) {
        return h('span', { class: 'text-12px text-gray-400' }, '非本人创建');
      }

      return h(
        NButton,
        { size: 'tiny', type: 'primary', quaternary: true, onClick: () => openStock(row) },
        { default: () => '入库' }
      );
    }
  }
];

async function loadProjectOptions() {
  projectLoading.value = true;

  try {
    const { data, error } = await fetchGetProjectSelectOptions();

    if (!error && data) projectOptions.value = data.projects ?? [];
  } finally {
    projectLoading.value = false;
  }
}

async function loadOrders() {
  if (!projectId.value) {
    orders.value = [];

    return;
  }

  loading.value = true;

  try {
    const { data, error } = await fetchGetPurchaseOrderList({ projectId: projectId.value, current: 1, size: 100 });

    if (!error && data) orders.value = data.records ?? [];
  } finally {
    loading.value = false;
  }
}

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.ProjectData.PurchaseOrder>[] = [
  { title: '采购单号', value: row => row.orderNo || '-', width: 18 },
  { title: '材料', value: row => row.materialName || '-', width: 16 },
  { title: '规格型号', value: row => row.spec || '无规格', width: 12 },
  { title: '采购数量', value: row => `${row.quantity} ${row.unit}`, align: 'right' },
  { title: '采购单价', value: row => formatMoney(row.unitPrice), align: 'right' },
  { title: '税务方式', value: row => ({ NO_INVOICE: '无票/不计税', TAX_INCLUDED: `含税 ${row.taxRate}%`, TAX_EXCLUDED: `未税 ${row.taxRate}%` })[row.taxMode] ?? '-' },
  { title: '税额', value: row => formatMoney(row.taxAmount), align: 'right' },
  { title: '采购金额', value: row => formatMoney(row.amount), align: 'right' },
  { title: '供应商', value: row => row.supplier || '-', width: 12 },
  { title: '仓库', value: row => row.warehouseName || '-', width: 10 },
  { title: '状态', value: row => ({ DRAFT: '草稿', CONFIRMED: '待付款', PAID: '已付款' })[row.status] ?? row.status },
  { title: '出入库', value: row => (row.stockedAt ? '已入库' : '未入库'), width: 8 },
  { title: '创建时间', value: row => row.createdAt || '-', width: 16 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const project = projectOptions.value.find(item => item.value === projectId.value)?.label ?? '-';

  return `项目：${project}`;
}

/** 拉取当前筛选下的全部采购单（分页聚合） */
async function collectOrderRows() {
  return fetchAllPages<Api.ProjectData.PurchaseOrder>(async current => {
    const { data, error } = await fetchGetPurchaseOrderList({ projectId: projectId.value ?? 0, current, size: 100 });

    if (error || !data) {
      throw new Error('数据加载失败');
    }

    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

function orderSummary(rows: Api.ProjectData.PurchaseOrder[]) {
  return [
    { label: '采购单数', value: `${rows.length} 单` },
    { label: '合计金额', value: formatMoney(rows.reduce((sum, row) => sum + row.amount, 0)) }
  ];
}

async function handlePrint() {
  if (!projectId.value) {
    window.$message?.warning('请先选择项目');

    return;
  }

  printing.value = true;

  try {
    const rows = await collectOrderRows();
    await printTable({
      title: '采购单列表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: orderSummary(rows)
    });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    printing.value = false;
  }
}

async function handleExport() {
  if (!projectId.value) {
    window.$message?.warning('请先选择项目');

    return;
  }

  exporting.value = true;

  try {
    const rows = await collectOrderRows();
    exportExcel({
      filename: '采购单列表',
      title: '采购单列表',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: orderSummary(rows)
    });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

function openCreate() {
  drawerVisible.value = true;
}

/** 新增完成后：切到对应项目查看新建的采购单 */
async function handleSubmitted(submittedProjectId: number) {
  if (submittedProjectId && submittedProjectId !== projectId.value) {
    projectId.value = submittedProjectId;

    return;
  }

  await loadOrders();
}

watch(projectId, () => void loadOrders());

onMounted(async () => {
  await loadProjectOptions();
  projectId.value = projectOptions.value[0]?.value ?? null;
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <NCard :bordered="false" size="small" class="card-wrapper">
        <div class="flex flex-wrap items-center gap-12px">
          <NSelect
            v-model:value="projectId"
            class="w-360px max-w-full"
            :options="projectOptions"
            :loading="projectLoading"
            filterable
            placeholder="请选择项目"
          />
          <span class="text-12px text-gray-400">
            采购单可在此选择仓库直接入库；付款及记账在付款管理处理
          </span>
        </div>
      </NCard>

      <NCard v-if="projectId" :bordered="false" size="small" class="card-wrapper">
        <div class="mb-12px flex items-center">
          <span class="text-13px text-gray-500">采购单 {{ orders.length }} 条</span>
          <div class="ml-auto flex gap-8px">
            <NButton size="small" type="primary" ghost @click="openCreate">
              <template #icon><SvgIcon icon="mdi:plus" /></template>
              新建采购单
            </NButton>
            <NButton size="small" :loading="printing" @click="handlePrint">
              <template #icon><SvgIcon icon="mdi:printer-outline" /></template>
              打印
            </NButton>
            <NButton size="small" :loading="exporting" @click="handleExport">
              <template #icon><SvgIcon icon="mdi:microsoft-excel" /></template>
              导出
            </NButton>
            <NButton size="small" :loading="loading" @click="loadOrders">
              <template #icon><SvgIcon icon="mdi:refresh" /></template>
              刷新
            </NButton>
          </div>
        </div>
        <NDataTable
          :columns="columns"
          :data="orders"
          :loading="loading"
          :row-key="row => row.id"
          :scroll-x="1400"
          size="small"
          striped
        />
      </NCard>
      <NEmpty v-else-if="!projectLoading" description="暂无项目" class="py-60px" />
    </NSpace>

    <PurchaseCreateDrawer
      v-model:visible="drawerVisible"
      :default-project-id="projectId"
      :project-options="projectOptions"
      @submitted="handleSubmitted"
    />

    <StockBackfillDrawer
      v-model:visible="stockVisible"
      mode="in"
      :record-id="stockRow?.id ?? 0"
      :material-name="stockRow?.materialName ?? ''"
      :quantity="stockRow?.quantity ?? 0"
      :unit="stockRow?.unit"
      :default-warehouse-id="stockRow?.warehouseId ?? 0"
      @submitted="loadOrders"
    />
  </div>
</template>

<style scoped>
.card-wrapper {
  border-radius: 8px;
}
</style>
