<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import {
  fetchDeleteStock,
  fetchGetProjectSelectOptions,
  fetchGetStockList,
  fetchGetWarehouseList
} from '@/service/api';
import { formatMoney } from '@/views/project/project-list/components/project-meta';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';
import StockOperateDrawer from './components/stock-operate-drawer.vue';
import StockMovementDrawer from './components/stock-movement-drawer.vue';
import StockRecordDrawer from './components/stock-record-drawer.vue';

defineOptions({ name: 'InventoryStock' });

type MovementMode = 'in' | 'out' | 'adjust';

const loading = ref(false);
const records = ref<Api.Inventory.Stock[]>([]);
const summary = ref<Api.Inventory.StockSummary>({ kinds: 0, amount: 0, shortage: 0 });

/** 项目筛选：null 表示全部项目（后端 projectId=0） */
const projectId = ref<number | null>(null);
const projectOptions = ref<Api.Project.SelectOption[]>([]);
/** 仓库筛选：null 表示全部仓库（后端 warehouseId=0） */
const warehouseId = ref<number | null>(null);
const warehouseOptions = ref<{ label: string; value: number }[]>([]);
const keyword = ref('');

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onUpdatePage(page: number) {
    pagination.page = page;
    void loadList();
  },
  onUpdatePageSize(size: number) {
    pagination.pageSize = size;
    pagination.page = 1;
    void loadList();
  }
});

const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingRow = ref<Api.Inventory.Stock | null>(null);

/** 出入库 / 补录抽屉 */
const movementVisible = ref(false);
const movementMode = ref<MovementMode>('in');
const movementMaterial = ref<Api.Inventory.Stock | null>(null);

/** 出入库流水抽屉 */
const recordVisible = ref(false);

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.Inventory.Stock>[] = [
  { title: '所属项目', value: row => row.projectName || '-', width: 18 },
  { title: '仓库', value: row => row.warehouseName || '-', width: 10 },
  { title: '材料名称', value: row => row.name, width: 16 },
  { title: '规格型号', value: row => row.spec || '-', width: 12 },
  { title: '单位', value: row => row.unit || '-', width: 6 },
  { title: '累计入库', value: row => String(row.quantity), align: 'right' },
  { title: '已领用', value: row => String(row.usedQuantity), align: 'right' },
  { title: '当前结存', value: row => String(row.remainQuantity), align: 'right' },
  { title: '单价', value: row => formatMoney(row.unitPrice), align: 'right' },
  { title: '金额', value: row => formatMoney(row.amount), align: 'right' },
  { title: '供应商', value: row => row.supplier || '-', width: 12 },
  { title: '最后变动', value: row => row.updatedAt || '-', width: 16 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const project = projectOptions.value.find(item => item.value === projectId.value)?.label ?? '全部项目';
  const warehouse = warehouseOptions.value.find(item => item.value === warehouseId.value)?.label ?? '全部仓库';
  const parts = [`项目：${project}`, `仓库：${warehouse}`];

  if (keyword.value.trim()) {
    parts.push(`关键词：${keyword.value.trim()}`);
  }

  return parts.join('　|　');
}

/** 拉取当前筛选下的全部库存（分页聚合） */
async function collectStockRows() {
  return fetchAllPages<Api.Inventory.Stock>(async current => {
    const { data, error } = await fetchGetStockList({
      projectId: projectId.value ?? 0,
      warehouseId: warehouseId.value ?? 0,
      keyword: keyword.value,
      current,
      size: 100
    });

    if (error || !data) {
      throw new Error('数据加载失败');
    }

    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

function stockSummary(rows: Api.Inventory.Stock[]) {
  return [
    { label: '库存种类', value: `${rows.length} 种` },
    { label: '库存总金额', value: formatMoney(rows.reduce((sum, row) => sum + row.amount, 0)) }
  ];
}

async function handlePrint() {
  printing.value = true;

  try {
    const rows = await collectStockRows();
    await printTable({
      title: '库存清单',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: stockSummary(rows)
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
    const rows = await collectStockRows();
    exportExcel({
      filename: '库存清单',
      title: '库存清单',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: stockSummary(rows)
    });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

const columns: DataTableColumns<Api.Inventory.Stock> = [
  {
    key: 'projectName',
    title: '所属项目',
    width: 150,
    ellipsis: { tooltip: true },
    render: row => row.projectName || '-'
  },
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
  {
    key: 'name',
    title: '材料',
    minWidth: 180,
    render: row =>
      h('div', { class: 'flex-col gap-1px' }, [
        h('span', { class: 'font-medium' }, row.name),
        row.spec ? h('span', { class: 'text-11px text-gray-400' }, row.spec) : null
      ])
  },
  { key: 'unit', title: '单位', width: 70, render: row => row.unit || '-' },
  {
    key: 'quantity',
    title: () => h('span', null, [h('span', null, '累计入库 '), h('span', { class: 'text-11px text-gray-400' }, '(入库)')]),
    width: 110,
    align: 'right',
    render: row => String(row.quantity)
  },
  {
    key: 'usedQuantity',
    title: () => h('span', null, [h('span', null, '已领用 '), h('span', { class: 'text-11px text-gray-400' }, '(出库)')]),
    width: 110,
    align: 'right',
    render: row => String(row.usedQuantity)
  },
  {
    key: 'remainQuantity',
    title: '当前结存',
    width: 100,
    align: 'right',
    render: row =>
      h(
        'span',
        { class: `font-mono font-medium ${row.remainQuantity <= 0 ? 'text-error' : 'text-success'}` },
        String(row.remainQuantity)
      )
  },
  { key: 'unitPrice', title: '单价', width: 110, align: 'right', render: row => formatMoney(row.unitPrice) },
  { key: 'amount', title: '金额', width: 120, align: 'right', render: row => formatMoney(row.amount) },
  { key: 'supplier', title: '供应商', width: 130, ellipsis: { tooltip: true }, render: row => row.supplier || '-' },
  {
    key: 'updatedAt',
    title: '最后变动',
    width: 150,
    render: row => h('span', { class: 'text-12px text-gray-500' }, row.updatedAt || '-')
  },
  {
    key: 'operate',
    title: '操作',
    width: 200,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-2px' }, [
        h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => openMovement('in', row) }, { default: () => '入库' }),
        h(
          NButton,
          { size: 'tiny', quaternary: true, type: 'warning', disabled: row.remainQuantity <= 0, onClick: () => openMovement('out', row) },
          { default: () => '出库' }
        ),
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => openEdit(row) }, { default: () => '编辑' }),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(row) },
          {
            trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除该库存吗？已有领用记录会保留。'
          }
        )
      ])
  }
];

const summaryCards = [
  { key: 'kinds', label: '库存种类', icon: 'mdi:package-variant-closed', color: '#2080f0' },
  { key: 'amount', label: '库存总金额', icon: 'mdi:cash-multiple', color: '#18a058' },
  { key: 'shortage', label: '结存预警', icon: 'mdi:alert-outline', color: '#d03050' }
] as const;

function summaryValue(key: (typeof summaryCards)[number]['key']) {
  if (key === 'amount') return formatMoney(summary.value.amount);
  return String(summary.value[key]);
}

async function loadProjectOptions() {
  const { data, error } = await fetchGetProjectSelectOptions();
  if (!error && data) projectOptions.value = data.projects ?? [];
}

async function loadWarehouseOptions() {
  const { data, error } = await fetchGetWarehouseList({ current: 1, size: 100 });
  if (!error && data) {
    warehouseOptions.value = (data.records ?? []).map(item => ({ label: item.name, value: item.id }));
  }
}

async function loadList() {
  loading.value = true;
  try {
    const { data, error } = await fetchGetStockList({
      projectId: projectId.value ?? 0,
      warehouseId: warehouseId.value ?? 0,
      keyword: keyword.value,
      current: pagination.page,
      size: pagination.pageSize
    });
    if (!error && data) {
      records.value = data.records ?? [];
      summary.value = data.summary ?? { kinds: 0, amount: 0, shortage: 0 };
      pagination.itemCount = data.total ?? 0;
    }
  } finally {
    loading.value = false;
  }
}

function openEdit(row: Api.Inventory.Stock) {
  operateType.value = 'edit';
  editingRow.value = row;
  drawerVisible.value = true;
}

/** 打开出库 / 入库 / 补录抽屉；row 为空表示从顶级按钮进入 */
function openMovement(mode: MovementMode, row?: Api.Inventory.Stock | null) {
  movementMode.value = mode;
  movementMaterial.value = row ?? null;
  movementVisible.value = true;
}

async function handleDelete(row: Api.Inventory.Stock) {
  const { error } = await fetchDeleteStock(row.id);
  if (error) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  await loadList();
}

function search() {
  pagination.page = 1;
  void loadList();
}

function reset() {
  projectId.value = null;
  warehouseId.value = null;
  keyword.value = '';
  pagination.page = 1;
  void loadList();
}

onMounted(async () => {
  await Promise.all([loadProjectOptions(), loadWarehouseOptions(), loadList()]);
});
</script>

<template>
  <div class="min-h-full bg-layout p-16px">
    <NSpace vertical :size="16">
      <!-- 库存汇总 -->
      <div class="grid grid-cols-1 gap-12px sm:grid-cols-3 lg:grid-cols-3">
        <NCard v-for="card in summaryCards" :key="card.key" :bordered="false" size="small" class="card-wrapper">
          <div class="flex items-center gap-12px">
            <div class="flex-center h-40px w-40px rounded-8px" :style="{ backgroundColor: `${card.color}1a` }">
              <SvgIcon :icon="card.icon" class="text-22px" :style="{ color: card.color }" />
            </div>
            <div class="flex-col">
              <span class="text-12px text-gray-500">{{ card.label }}</span>
              <span class="text-20px font-semibold" :style="{ color: card.color }">{{ summaryValue(card.key) }}</span>
            </div>
          </div>
        </NCard>
      </div>

      <!-- 筛选 + 列表 -->
      <NCard :bordered="false" size="small" class="card-wrapper">
        <div class="mb-12px flex flex-wrap items-center gap-8px">
          <NSelect
            v-model:value="projectId"
            class="w-260px max-w-full"
            :options="projectOptions"
            clearable
            filterable
            placeholder="全部项目"
          />
          <NSelect
            v-model:value="warehouseId"
            class="w-180px max-w-full"
            :options="warehouseOptions"
            clearable
            placeholder="全部仓库"
          />
          <NInput
            v-model:value="keyword"
            class="w-220px max-w-full"
            placeholder="搜索材料名称 / 规格 / 供应商"
            clearable
            @keyup.enter="search"
          />
          <NButton size="small" type="primary" ghost @click="search">查询</NButton>
          <NButton size="small" quaternary @click="reset">重置</NButton>
          <div class="ml-auto flex items-center gap-8px">
            <NButton size="small" type="primary" @click="openMovement('in')">
              <template #icon><SvgIcon icon="mdi:arrow-down-bold-circle-outline" /></template>
              入库
            </NButton>
            <NButton size="small" type="warning" @click="openMovement('out')">
              <template #icon><SvgIcon icon="mdi:arrow-up-bold-circle-outline" /></template>
              出库
            </NButton>
            <NButton size="small" @click="openMovement('adjust')">
              <template #icon><SvgIcon icon="mdi:history" /></template>
              补录
            </NButton>
            <NButton size="small" quaternary @click="recordVisible = true">
              <template #icon><SvgIcon icon="mdi:format-list-bulleted" /></template>
              流水
            </NButton>
            <NButton size="small" :loading="printing" @click="handlePrint">
              <template #icon><SvgIcon icon="mdi:printer-outline" /></template>
              打印
            </NButton>
            <NButton size="small" :loading="exporting" @click="handleExport">
              <template #icon><SvgIcon icon="mdi:microsoft-excel" /></template>
              导出
            </NButton>
            <NButton size="small" :loading="loading" @click="loadList">
              <template #icon><SvgIcon icon="mdi:refresh" /></template>
              刷新
            </NButton>
          </div>
        </div>

        <NAlert v-if="summary.shortage > 0" type="warning" :show-icon="true" class="mb-12px" :bordered="false">
          有 {{ summary.shortage }} 种材料当前结存为 0，请及时采购入库。
        </NAlert>

        <NDataTable
          :columns="columns"
          :data="records"
          :loading="loading"
          :pagination="pagination"
          :row-key="row => row.id"
          :scroll-x="1280"
          remote
          size="small"
          striped
        />
      </NCard>
    </NSpace>

    <StockOperateDrawer v-model:visible="drawerVisible" :operate-type="operateType" :row-data="editingRow" @submitted="loadList" />

    <StockMovementDrawer
      v-model:visible="movementVisible"
      :mode="movementMode"
      :default-material="movementMaterial"
      :project-options="projectOptions"
      @submitted="loadList"
    />

    <StockRecordDrawer v-model:visible="recordVisible" :default-project-id="projectId ?? 0" />
  </div>
</template>

<style scoped>
.card-wrapper {
  border-radius: 8px;
}
</style>
