<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NPopconfirm } from 'naive-ui';
import { fetchDeleteMaterialCatalog, fetchGetMaterialCatalogList } from '@/service/api';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';
import MaterialCatalogOperateDrawer from './components/material-catalog-operate-drawer.vue';

defineOptions({ name: 'InventoryMaterials' });

const loading = ref(false);
const records = ref<Api.Inventory.MaterialCatalog[]>([]);
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
const editingRow = ref<Api.Inventory.MaterialCatalog | null>(null);
const printing = ref(false);
const exporting = ref(false);

const columns: DataTableColumns<Api.Inventory.MaterialCatalog> = [
  {
    key: 'name',
    title: '材料',
    minWidth: 200,
    render: row =>
      h('div', { class: 'flex-col gap-1px' }, [
        h('span', { class: 'font-medium' }, row.name),
        row.spec ? h('span', { class: 'text-11px text-gray-400' }, row.spec) : null
      ])
  },
  { key: 'unit', title: '单位', width: 90, render: row => row.unit || '-' },
  { key: 'unitPrice', title: '默认单价', width: 130, align: 'right', render: row => formatMoney(row.unitPrice) },
  { key: 'remark', title: '备注', minWidth: 160, ellipsis: { tooltip: true }, render: row => row.remark || '-' },
  {
    key: 'operate',
    title: '操作',
    width: 110,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-2px' }, [
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => openEdit(row) }, { default: () => '编辑' }),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(row) },
          {
            trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除该材料档案吗？已用于项目材料的档案不能删除。'
          }
        )
      ])
  }
];

const printExportColumns: PrintExportColumn<Api.Inventory.MaterialCatalog>[] = [
  { title: '材料名称', value: row => row.name, width: 18 },
  { title: '规格型号', value: row => row.spec || '-', width: 16 },
  { title: '单位', value: row => row.unit || '-', width: 8 },
  { title: '默认单价', value: row => formatMoney(row.unitPrice), align: 'right' },
  { title: '备注', value: row => row.remark || '-', width: 20 }
];

async function loadList() {
  loading.value = true;
  try {
    const { data, error } = await fetchGetMaterialCatalogList({
      keyword: keyword.value.trim() || null,
      current: pagination.page,
      size: pagination.pageSize
    });
    if (!error && data) {
      records.value = data.records ?? [];
      pagination.itemCount = data.total ?? 0;
    }
  } finally {
    loading.value = false;
  }
}

async function collectRows() {
  return fetchAllPages<Api.Inventory.MaterialCatalog>(async current => {
    const { data, error } = await fetchGetMaterialCatalogList({
      keyword: keyword.value.trim() || null,
      current,
      size: 100
    });
    if (error || !data) throw new Error('数据加载失败');
    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

function openAdd() {
  operateType.value = 'add';
  editingRow.value = null;
  drawerVisible.value = true;
}

function openEdit(row: Api.Inventory.MaterialCatalog) {
  operateType.value = 'edit';
  editingRow.value = row;
  drawerVisible.value = true;
}

async function handleDelete(row: Api.Inventory.MaterialCatalog) {
  const { error } = await fetchDeleteMaterialCatalog(row.id);
  if (error) return;
  window.$message?.success('材料档案已删除');
  await loadList();
}

async function search() {
  pagination.page = 1;
  await loadList();
}

async function handlePrint() {
  printing.value = true;
  try {
    const rows = await collectRows();
    await printTable({ title: '材料档案', columns: printExportColumns, rows });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    printing.value = false;
  }
}

async function handleExport() {
  exporting.value = true;
  try {
    const rows = await collectRows();
    exportExcel({ filename: '材料档案', title: '材料档案', columns: printExportColumns, rows });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

onMounted(loadList);
</script>

<template>
  <div class="p-16px">
    <NCard :bordered="false" size="small" class="card-wrapper">
      <div class="mb-12px flex flex-wrap items-center gap-8px">
        <NInput
          v-model:value="keyword"
          class="w-320px max-w-full"
          placeholder="搜索材料名称 / 规格"
          clearable
          @keyup.enter="search"
        />
        <NButton size="small" type="primary" ghost @click="search">查询</NButton>
        <div class="ml-auto flex items-center gap-8px">
          <NButton size="small" type="primary" ghost @click="openAdd">
            <template #icon><SvgIcon icon="mdi:plus" /></template>
            新增材料
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
      <NDataTable
        :columns="columns"
        :data="records"
        :loading="loading"
        :pagination="pagination"
        :row-key="row => row.id"
        :scroll-x="920"
        remote
        size="small"
        striped
      />
    </NCard>

    <MaterialCatalogOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingRow"
      @submitted="loadList"
    />
  </div>
</template>

<style scoped>
.card-wrapper {
  border-radius: 8px;
}
</style>
