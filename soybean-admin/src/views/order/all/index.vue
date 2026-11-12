<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import type { DataTableRowKey } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { formatOrderMoney } from '../components/order-meta';
import { useOrderData } from '../components/use-order-data';
import OrderDetailDrawer from './components/order-detail-drawer.vue';
import OrderOperateDrawer from './components/order-operate-drawer.vue';
import OrderSearchBar from './components/order-search-bar.vue';
import OrderTable from './components/order-table.vue';

defineOptions({ name: 'OrderAll' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
  status: Api.Order.Status | null;
  customerResponse: Api.Order.CustomerResponse | null;
  dateFrom: string | null;
  dateTo: string | null;
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

const {
  loading,
  records,
  total,
  stats,
  customerOptions,
  driverOptions,
  vehicleOptions,
  loadList,
  loadStats,
  loadFormOptions,
  deleteOrders
} = useOrderData();

const filters = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  status: null,
  customerResponse: null,
  dateFrom: null,
  dateTo: null
});

const checkedRowKeys = ref<DataTableRowKey[]>([]);
const drawerVisible = ref(false);
const detailVisible = ref(false);
const detailId = ref<number | null>(null);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Order.Item | null>(null);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 个订单`,
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

const statCards = computed(() => [
  { label: '订单总数', value: stats.value.total, suffix: '单', icon: 'mdi:clipboard-text-outline', color: 'primary' },
  { label: '今日预约', value: stats.value.today, suffix: '单', icon: 'mdi:calendar-today-outline', color: 'info' },
  { label: '进行中', value: stats.value.active, suffix: '单', icon: 'mdi:car-connected', color: 'warning' },
  { label: '已完成收入', value: formatOrderMoney(stats.value.revenue), suffix: '', icon: 'mdi:cash-check', color: 'success' }
]);

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
  filters.customerResponse = null;
  filters.dateFrom = null;
  filters.dateTo = null;
  pagination.page = 1;
  await getData();
}

async function handleRefresh() {
  await Promise.all([getData(), loadStats()]);
}

function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

function handleEdit(row: Api.Order.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
  detailVisible.value = false;
}

function handleDetail(row: Api.Order.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

async function handleDelete(row: Api.Order.Item) {
  const success = await deleteOrders([row.id]);
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
  const success = await deleteOrders(ids);
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
  await Promise.all([getData(), loadStats()]);
}

async function handleSubmitted() {
  if (operateType.value === 'add') pagination.page = 1;
  await Promise.all([getData(), loadStats()]);
}

onMounted(async () => {
  await Promise.all([getData(), loadStats(), loadFormOptions()]);
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <div class="grid grid-cols-2 gap-12px lg:grid-cols-4">
        <div v-for="card in statCards" :key="card.label" class="card-wrapper flex items-center gap-12px p-16px">
          <div
            class="flex h-42px w-42px flex-shrink-0 items-center justify-center rounded-10px"
            :style="{ backgroundColor: `rgb(var(--${card.color}-color) / 0.12)`, color: `rgb(var(--${card.color}-color))` }"
          >
            <SvgIcon :icon="card.icon" class="text-22px" />
          </div>
          <div class="min-w-0">
            <p class="text-12px text-gray-400">{{ card.label }}</p>
            <p class="mt-2px truncate text-20px font-semibold">{{ card.value }}<small class="ml-2px text-12px font-normal text-gray-400">{{ card.suffix }}</small></p>
          </div>
        </div>
      </div>

      <NCard :bordered="false" class="card-wrapper" size="small">
        <div class="mb-12px flex flex-wrap items-end justify-between gap-8px">
          <div>
            <h1 class="text-16px font-semibold">全部订单</h1>
            <p class="mt-4px text-13px text-gray-500">管理报价、派单、行程与客户反馈</p>
          </div>
          <span class="text-12px text-gray-400">共 {{ total }} 个订单</span>
        </div>

        <NSpace vertical :size="12">
          <OrderSearchBar
            v-model:keyword="filters.keyword"
            v-model:status="filters.status"
            v-model:customer-response="filters.customerResponse"
            v-model:date-from="filters.dateFrom"
            v-model:date-to="filters.dateTo"
            :loading="loading"
            :checked-count="checkedRowKeys.length"
            @search="handleSearch"
            @reset="handleReset"
            @refresh="handleRefresh"
            @add="handleAdd"
            @batch-delete="handleBatchDelete"
          />

          <OrderTable
            v-model:checked-row-keys="checkedRowKeys"
            :data="records"
            :loading="loading"
            :pagination="tablePagination"
            @detail="handleDetail"
            @edit="handleEdit"
            @delete="handleDelete"
          />
        </NSpace>
      </NCard>
    </NSpace>

    <OrderOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      :customer-options="customerOptions"
      :driver-options="driverOptions"
      :vehicle-options="vehicleOptions"
      @submitted="handleSubmitted"
    />

    <OrderDetailDrawer v-model:visible="detailVisible" :order-id="detailId" @edit="handleEdit" />
  </div>
</template>
