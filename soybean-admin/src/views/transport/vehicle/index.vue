<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import type { DataTableRowKey } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import VehicleOperateDrawer from './components/vehicle-operate-drawer.vue';
import VehicleSearchBar from './components/vehicle-search-bar.vue';
import VehicleTable from './components/vehicle-table.vue';
import { useVehicleData } from './components/use-vehicle-data';

defineOptions({ name: 'TransportVehicle' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
  status: Api.Transport.VehicleStatus | null;
}

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

const { loading, records, total, driverOptions, loadList, deleteVehicles, loadDriverOptions } = useVehicleData();

const filters = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  status: null
});

const checkedRowKeys = ref<DataTableRowKey[]>([]);
const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Transport.Vehicle | null>(null);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 辆车`,
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

const tablePagination = computed(() => ({ ...pagination }));

async function getData() {
  filters.current = pagination.page;
  filters.size = pagination.pageSize;
  await loadList(filters);
  pagination.itemCount = total.value;
}

async function handleSearch() {
  pagination.page = 1;
  await getData();
}

async function handleReset() {
  filters.keyword = '';
  filters.status = null;
  pagination.page = 1;
  await getData();
}

async function handleRefresh() {
  await getData();
}

function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

function handleEdit(row: Api.Transport.Vehicle) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

async function handleDelete(row: Api.Transport.Vehicle) {
  const success = await deleteVehicles([row.id]);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  await afterDeleted(1);
}

async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(Number);
  if (!ids.length) return;
  const success = await deleteVehicles(ids);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  checkedRowKeys.value = [];
  await afterDeleted(ids.length);
}

async function afterDeleted(deletedCount: number) {
  const maxPage = Math.max(1, Math.ceil((total.value - deletedCount) / pagination.pageSize));
  if (pagination.page > maxPage) pagination.page = maxPage;
  await getData();
}

async function handleSubmitted() {
  if (operateType.value === 'add') pagination.page = 1;
  await getData();
}

onMounted(async () => {
  await Promise.all([getData(), loadDriverOptions()]);
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <NCard :bordered="false" class="card-wrapper" size="small">
        <div class="mb-12px flex flex-wrap items-end justify-between gap-8px">
          <div>
            <h1 class="text-16px font-semibold">车辆管理</h1>
            <p class="mt-4px text-13px text-gray-500">维护运营车辆、服务车型与司机绑定</p>
          </div>
          <span class="text-12px text-gray-400">共 {{ total }} 辆车</span>
        </div>

        <NSpace vertical :size="12">
          <VehicleSearchBar
            v-model:keyword="filters.keyword"
            v-model:status="filters.status"
            :loading="loading"
            :checked-count="checkedRowKeys.length"
            @search="handleSearch"
            @reset="handleReset"
            @refresh="handleRefresh"
            @add="handleAdd"
            @batch-delete="handleBatchDelete"
          />

          <VehicleTable
            v-model:checked-row-keys="checkedRowKeys"
            :data="records"
            :loading="loading"
            :pagination="tablePagination"
            @edit="handleEdit"
            @delete="handleDelete"
          />
        </NSpace>
      </NCard>
    </NSpace>

    <VehicleOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      :driver-options="driverOptions"
      @submitted="handleSubmitted"
    />
  </div>
</template>
