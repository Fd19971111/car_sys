<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import SupplierOperateDrawer from './components/supplier-operate-drawer.vue';
import { useSupplierData } from './components/use-supplier-data';

defineOptions({ name: 'SettingsSupplier' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
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

const appStore = useAppStore();
const isMobile = computed(() => appStore.isMobile);
const { loading, records, total, loadList, deleteSuppliers } = useSupplierData();
const searchParams = reactive<SearchParams>({ current: 1, size: 10, keyword: '' });
const checkedRowKeys = ref<DataTableRowKey[]>([]);
const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Supplier.Item | null>(null);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 家供应商`,
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

const columns: DataTableColumns<Api.Supplier.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'name',
    title: '供应商名称',
    minWidth: 180,
    fixed: 'left',
    render: row => h('div', { class: 'flex-col gap-1px' }, [
      h('span', { class: 'font-medium' }, row.name),
      h('span', { class: 'text-11px text-gray-400' }, `ID ${row.id}`)
    ])
  },
  { key: 'contactName', title: '联系人', width: 120, render: row => row.contactName || '—' },
  { key: 'phone', title: '联系电话', width: 150, render: row => row.phone || '—' },
  { key: 'address', title: '地址', minWidth: 200, ellipsis: { tooltip: true }, render: row => row.address || '—' },
  { key: 'remark', title: '备注', minWidth: 180, ellipsis: { tooltip: true }, render: row => row.remark || '—' },
  {
    key: 'operate',
    title: '操作',
    width: 130,
    align: 'center',
    fixed: 'right',
    render: row => h('div', { class: 'flex-center gap-4px' }, [
      h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
      h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
        trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
        default: () => '确认删除该供应商吗？'
      })
    ])
  }
];

async function getData() {
  searchParams.current = pagination.page;
  searchParams.size = pagination.pageSize;
  await loadList(searchParams);
  pagination.itemCount = total.value;
}

async function handleSearch() {
  pagination.page = 1;
  await getData();
}

async function handleReset() {
  searchParams.keyword = '';
  pagination.page = 1;
  await getData();
}

function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

function handleEdit(row: Api.Supplier.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

async function handleDelete(row: Api.Supplier.Item) {
  const success = await deleteSuppliers([row.id]);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  await afterDeleted(1);
}

async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(Number);
  const success = await deleteSuppliers(ids);
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
    <NCard :bordered="false" class="card-wrapper" size="small">
      <NSpace vertical :size="12">
        <div class="flex flex-wrap items-center gap-8px">
          <NInput
            v-model:value="searchParams.keyword"
            class="w-260px min-w-160px"
            placeholder="搜索供应商 / 联系人 / 电话 / 地址"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix><SvgIcon icon="mdi:magnify" class="text-gray-400" /></template>
          </NInput>
          <NButton type="primary" ghost size="small" @click="handleSearch">
            <template #icon><SvgIcon icon="mdi:magnify" /></template>查询
          </NButton>
          <NButton size="small" @click="handleReset">
            <template #icon><SvgIcon icon="mdi:backspace-outline" /></template>重置
          </NButton>
          <div class="ml-auto flex items-center gap-8px">
            <NButton size="small" type="primary" ghost @click="handleAdd">
              <template #icon><SvgIcon icon="mdi:domain-plus" /></template>新增供应商
            </NButton>
            <NPopconfirm @positive-click="handleBatchDelete">
              <template #trigger>
                <NButton size="small" ghost type="error" :disabled="!checkedRowKeys.length">
                  <template #icon><SvgIcon icon="mdi:delete-outline" /></template>批量删除
                </NButton>
              </template>
              确认删除选中的 {{ checkedRowKeys.length }} 家供应商吗？
            </NPopconfirm>
            <NButton size="small" @click="getData">
              <template #icon><SvgIcon icon="mdi:refresh" :class="{ 'animate-spin': loading }" /></template>刷新
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
          :scroll-x="1050"
          remote
          size="small"
          striped
        />
      </NSpace>
    </NCard>
    <SupplierOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      @submitted="handleSubmitted"
    />
  </div>
</template>