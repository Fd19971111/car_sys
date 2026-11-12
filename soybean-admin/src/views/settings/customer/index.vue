<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import CustomerDetailDrawer from './components/customer-detail-drawer.vue';
import CustomerOperateDrawer from './components/customer-operate-drawer.vue';
import {
  CUSTOMER_SOURCE_OPTIONS,
  CUSTOMER_STATUS_OPTIONS,
  formatDateTime,
  getCustomerSourceLabel,
  getCustomerStatusMeta
} from './components/customer-meta';
import { useCustomerData } from './components/use-customer-data';

defineOptions({
  name: 'SettingsCustomer'
});

interface SearchParams {
  current: number;
  size: number;
  /** 关键词：姓名 / 手机号 / 房屋地址 */
  keyword: string;
  /** 客户状态 */
  status: Api.Customer.Status | null;
  /** 客户来源 */
  source: string | null;
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

const appStore = useAppStore();

const isMobile = computed(() => appStore.isMobile);

const { loading, records, total, loadList, deleteCustomers } = useCustomerData();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  status: null,
  source: null
});

/** 表格勾选的客户 */
const checkedRowKeys = ref<DataTableRowKey[]>([]);

const detailVisible = ref(false);
const detailId = ref<number | null>(null);

const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Customer.Item | null>(null);

/** 分页配置 */
const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 位客户`,
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

const columns: DataTableColumns<Api.Customer.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'name',
    title: '客户',
    width: 150,
    fixed: 'left',
    render: row => {
      const meta = getCustomerStatusMeta(row.status);

      return h('div', { class: 'flex-y-center gap-8px' }, [
        h(
          'div',
          {
            class: 'h-28px w-28px flex-center flex-shrink-0 rounded-full text-12px font-medium',
            style: {
              backgroundColor: `rgb(var(--${meta.color}-color) / 0.12)`,
              color: `rgb(var(--${meta.color}-color))`
            }
          },
          row.name.slice(0, 1)
        ),
        h('div', { class: 'flex-col' }, [
          h('span', { class: 'font-medium' }, row.name),
          h('span', { class: 'text-12px text-gray-400' }, `ID ${row.id}`)
        ])
      ]);
    }
  },
  {
    key: 'phone',
    title: '手机号',
    width: 130,
    render: row => h('span', { class: 'font-mono' }, row.phone)
  },
  {
    key: 'address',
    title: '房屋地址',
    minWidth: 200,
    ellipsis: { tooltip: true },
    render: row => row.address || '-'
  },
  {
    key: 'source',
    title: '客户来源',
    width: 110,
    render: row => h(NTag, { size: 'small', bordered: false }, { default: () => getCustomerSourceLabel(row.source) })
  },
  {
    key: 'status',
    title: '客户状态',
    width: 100,
    render: row => {
      const meta = getCustomerStatusMeta(row.status);

      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'createdAt',
    title: '创建时间',
    width: 150,
    render: row => formatDateTime(row.createdAt)
  },
  {
    key: 'operate',
    title: '操作',
    width: 170,
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
            default: () => '确认删除该客户吗？'
          }
        )
      ])
  }
];

function rowProps(row: Api.Customer.Item) {
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

/** 加载客户列表 */
async function getData() {
  searchParams.current = pagination.page;
  searchParams.size = pagination.pageSize;

  await loadList(searchParams);

  pagination.itemCount = total.value;
}

/** 查询 */
async function handleSearch() {
  pagination.page = 1;

  await getData();
}

/** 重置查询条件 */
async function handleReset() {
  searchParams.keyword = '';
  searchParams.status = null;
  searchParams.source = null;
  pagination.page = 1;

  await getData();
}

/** 刷新列表 */
async function handleRefresh() {
  await getData();
}

/** 查看客户详情（详情抽屉按 id 从后端拉取完整信息） */
function handleDetail(row: Api.Customer.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

/** 新增客户 */
function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

/** 编辑客户 */
function handleEdit(row: Api.Customer.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

/** 从详情抽屉进入编辑（payload 为后端返回的最新详情） */
function handleEditFromDetail(row: Api.Customer.Item) {
  detailVisible.value = false;

  handleEdit(row);
}

/** 删除客户 */
async function handleDelete(row: Api.Customer.Item) {
  const success = await deleteCustomers([row.id]);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');

  await afterDeleted(1);
}

/** 批量删除客户 */
async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(key => Number(key));
  const success = await deleteCustomers(ids);

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

  await getData();
}

/** 新增 / 编辑提交成功后刷新 */
async function handleSubmitted() {
  if (operateType.value === 'add') {
    pagination.page = 1;
  }

  await getData();
}

onMounted(async () => {
  await getData();
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <!-- 搜索 / 筛选 / 操作按钮：全部排在同一行 -->
      <NCard :bordered="false" class="card-wrapper" size="small">
        <NSpace vertical :size="12">
          <div class="flex flex-wrap items-center gap-8px">
            <NInput
              v-model:value="searchParams.keyword"
              class="w-200px min-w-140px flex-shrink"
              placeholder="搜索姓名 / 手机号 / 房屋地址"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <SvgIcon icon="mdi:magnify" class="text-gray-400" />
              </template>
            </NInput>
            <NSelect
              v-model:value="searchParams.status"
              class="w-110px flex-shrink-0"
              :options="CUSTOMER_STATUS_OPTIONS"
              placeholder="客户状态"
              clearable
            />
            <NSelect
              v-model:value="searchParams.source"
              class="w-120px flex-shrink-0"
              :options="CUSTOMER_SOURCE_OPTIONS"
              placeholder="客户来源"
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
                  <SvgIcon icon="mdi:account-plus-outline" />
                </template>
                新增客户
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
                确认删除选中的 {{ checkedRowKeys.length }} 位客户吗？
              </NPopconfirm>
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
            :scroll-x="1180"
            remote
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </NSpace>

    <!-- 新增 / 编辑 -->
    <CustomerOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      @submitted="handleSubmitted"
    />

    <!-- 客户详情 -->
    <CustomerDetailDrawer v-model:visible="detailVisible" :customer-id="detailId" @edit="handleEditFromDetail" />
  </div>
</template>

<style scoped></style>
