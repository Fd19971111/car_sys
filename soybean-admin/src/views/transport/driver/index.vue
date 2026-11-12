<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import type { DataTableRowKey } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import DriverOperateDrawer from './components/driver-operate-drawer.vue';
import DriverSearchBar from './components/driver-search-bar.vue';
import DriverTable from './components/driver-table.vue';
import { useDriverData } from './components/use-driver-data';

defineOptions({ name: 'TransportDriver' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
  serviceStatus: Api.Transport.DriverServiceStatus | null;
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

const { loading, records, total, loadList, deleteDrivers } = useDriverData();

const filters = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  serviceStatus: null
});

const checkedRowKeys = ref<DataTableRowKey[]>([]);
const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Transport.Driver | null>(null);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 位司机`,
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
  filters.serviceStatus = null;
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

function handleEdit(row: Api.Transport.Driver) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

async function handleDelete(row: Api.Transport.Driver) {
  const success = await deleteDrivers([row.id]);
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
  const success = await deleteDrivers(ids);
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

onMounted(() => void getData());
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <NCard :bordered="false" class="card-wrapper" size="small">
        <div class="mb-12px flex flex-wrap items-end justify-between gap-8px">
          <div>
            <h1 class="text-16px font-semibold">司机管理</h1>
            <p class="mt-4px text-13px text-gray-500">维护司机档案、接单状态与服务评分</p>
          </div>
          <span class="text-12px text-gray-400">共 {{ total }} 位司机</span>
        </div>

        <NSpace vertical :size="12">
          <DriverSearchBar
            v-model:keyword="filters.keyword"
            v-model:service-status="filters.serviceStatus"
            :loading="loading"
            :checked-count="checkedRowKeys.length"
            @search="handleSearch"
            @reset="handleReset"
            @refresh="handleRefresh"
            @add="handleAdd"
            @batch-delete="handleBatchDelete"
          />

          <DriverTable
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

    <DriverOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      @submitted="handleSubmitted"
    />
  </div>
</template>
