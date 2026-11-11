<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { fetchGetMaterialUsageList, fetchGetProjectSelectOptions } from '@/service/api';
import { useAuthStore } from '@/store/modules/auth';
import { useProjectData } from '@/views/project/data/components/use-project-data';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';
import StockBackfillDrawer from '../components/stock-backfill-drawer.vue';
import UsageCreateDrawer from './components/usage-create-drawer.vue';

defineOptions({ name: 'InventoryUsages' });

const authStore = useAuthStore();
const currentUserId = computed(() => authStore.userInfo.userId);

const { materials, loadMaterials, usages, usageLoading, usageTotal, loadUsages, deleteUsage } = useProjectData();
const projectOptions = ref<Api.Project.SelectOption[]>([]);
const projectId = ref<number | null>(null);
const projectLoading = ref(false);
const materialId = ref<number | null>(null);
const createVisible = ref(false);

/** 领用出库抽屉 */
const stockVisible = ref(false);
const stockRow = ref<Api.ProjectData.Usage | null>(null);

/** 只有本人创建（历史数据不限）且未出库的领用记录才能补出库 */
function canStock(row: Api.ProjectData.Usage) {
  return !row.stockedAt && (!row.createdBy || row.createdBy === currentUserId.value);
}

function openStock(row: Api.ProjectData.Usage) {
  stockRow.value = row;
  stockVisible.value = true;
}
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

const materialOptions = ref<{ label: string; value: number }[]>([]);

const columns: DataTableColumns<Api.ProjectData.Usage> = [
  { key: 'usedAt', title: '领取日期', width: 120, render: row => h('span', { class: 'font-mono' }, row.usedAt) },
  { key: 'materialName', title: '材料', minWidth: 180, ellipsis: { tooltip: true } },
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
  { key: 'quantity', title: '领取数量', width: 110, align: 'right', render: row => h('span', { class: 'font-mono' }, String(row.quantity)) },
  { key: 'receiver', title: '领取人', width: 120, render: row => row.receiver || '-' },
  { key: 'remark', title: '备注', minWidth: 160, ellipsis: { tooltip: true }, render: row => row.remark || '-' },
  { key: 'createdAt', title: '登记时间', width: 170, render: row => row.createdAt || '-' },
  {
    key: 'stockStatus',
    title: '出入库',
    width: 100,
    render: row =>
      h(
        NTag,
        { size: 'small', bordered: false, type: row.stockedAt ? 'success' : 'warning' },
        { default: () => (row.stockedAt ? '已出库' : '未出库') }
      )
  },
  {
    key: 'operate',
    title: '操作',
    width: 140,
    align: 'center',
    render: row =>
      h('div', { class: 'flex-center gap-2px' }, [
        row.stockedAt || !canStock(row)
          ? null
          : h(
              NButton,
              { size: 'tiny', quaternary: true, type: 'warning', onClick: () => openStock(row) },
              { default: () => '出库' }
            ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(row) },
          {
            trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除这条领取记录吗？'
          }
        )
      ])
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

async function loadList() {
  if (!projectId.value) return;
  await loadUsages({
    projectId: projectId.value,
    materialId: materialId.value,
    current: pagination.page,
    size: pagination.pageSize
  });
  pagination.itemCount = usageTotal.value;
}

async function loadProjectData() {
  if (!projectId.value) return;
  pagination.page = 1;
  materialId.value = null;
  await loadMaterials({ projectId: projectId.value, current: 1, size: 100 });
  materialOptions.value = materials.value.map(item => ({
    label: item.spec ? `${item.name} · ${item.spec}` : item.name,
    value: item.id
  }));
  await loadList();
}

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.ProjectData.Usage>[] = [
  { title: '领取日期', value: row => row.usedAt || '-', width: 12 },
  { title: '材料', value: row => row.materialName || '-', width: 16 },
  { title: '仓库', value: row => row.warehouseName || '-', width: 10 },
  { title: '领取数量', value: row => String(row.quantity), align: 'right' },
  { title: '领取人', value: row => row.receiver || '-', width: 10 },
  { title: '备注', value: row => row.remark || '-', width: 16 },
  { title: '出入库', value: row => (row.stockedAt ? '已出库' : '未出库'), width: 8 },
  { title: '登记时间', value: row => row.createdAt || '-', width: 16 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const project = projectOptions.value.find(item => item.value === projectId.value)?.label ?? '-';
  const material = materialOptions.value.find(item => item.value === materialId.value)?.label;
  const parts = [`项目：${project}`, `材料：${material ?? '全部材料'}`];

  return parts.join('　|　');
}

/** 拉取当前筛选下的全部领取记录（分页聚合） */
async function collectUsageRows() {
  return fetchAllPages<Api.ProjectData.Usage>(async current => {
    const { data, error } = await fetchGetMaterialUsageList({
      projectId: projectId.value ?? 0,
      materialId: materialId.value ?? null,
      current,
      size: 100
    });

    if (error || !data) {
      throw new Error('数据加载失败');
    }

    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

function usageSummary(rows: Api.ProjectData.Usage[]) {
  return [
    { label: '领取笔数', value: `${rows.length} 笔` },
    { label: '领取数量合计', value: String(rows.reduce((sum, row) => sum + row.quantity, 0)) }
  ];
}

async function handlePrint() {
  if (!projectId.value) {
    window.$message?.warning('请先选择项目');

    return;
  }

  printing.value = true;

  try {
    const rows = await collectUsageRows();
    await printTable({
      title: '领取记录',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: usageSummary(rows)
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
    const rows = await collectUsageRows();
    exportExcel({
      filename: '领取记录',
      title: '领取记录',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: usageSummary(rows)
    });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

/** 打开「新增领取」抽屉 */
function openCreate() {
  if (!projectId.value) {
    window.$message?.warning('请先选择项目');

    return;
  }

  createVisible.value = true;
}

async function handleDelete(row: Api.ProjectData.Usage) {
  if (!(await deleteUsage(row.id))) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  await loadList();
}

watch(projectId, () => void loadProjectData());
watch(materialId, () => {
  pagination.page = 1;
  void loadList();
});

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
          <NSelect
            v-model:value="materialId"
            class="w-260px max-w-full"
            :options="materialOptions"
            clearable
            filterable
            placeholder="全部材料"
          />
          <NButton class="ml-auto" size="small" type="primary" ghost :disabled="!projectId" @click="openCreate">
            <template #icon><SvgIcon icon="mdi:plus" /></template>
            新增领取
          </NButton>
          <NButton size="small" :loading="printing" @click="handlePrint">
            <template #icon><SvgIcon icon="mdi:printer-outline" /></template>
            打印
          </NButton>
          <NButton size="small" :loading="exporting" @click="handleExport">
            <template #icon><SvgIcon icon="mdi:microsoft-excel" /></template>
            导出
          </NButton>
          <NButton size="small" :loading="usageLoading" @click="loadList">
            <template #icon><SvgIcon icon="mdi:refresh" /></template>
            刷新
          </NButton>
        </div>
      </NCard>

      <NCard v-if="projectId" :bordered="false" size="small" class="card-wrapper">
        <NDataTable
          :columns="columns"
          :data="usages"
          :loading="usageLoading"
          :pagination="pagination"
          :row-key="row => row.id"
          :scroll-x="1100"
          remote
          size="small"
          striped
        />
      </NCard>
      <NEmpty v-else-if="!projectLoading" description="暂无项目" class="py-60px" />
    </NSpace>

    <UsageCreateDrawer v-model:visible="createVisible" :project-id="projectId" @submitted="loadList" />

    <StockBackfillDrawer
      v-model:visible="stockVisible"
      mode="out"
      :record-id="stockRow?.id ?? 0"
      :material-name="stockRow?.materialName ?? ''"
      :quantity="stockRow?.quantity ?? 0"
      :unit="''"
      :default-warehouse-id="stockRow?.warehouseId ?? 0"
      @submitted="loadList"
    />
  </div>
</template>

<style scoped>
.card-wrapper {
  border-radius: 8px;
}
</style>