<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import { formatDate, formatMoney, formatMoneyInWan, formatNumber } from '@/views/finance/components/finance-meta';
import { useProjectOptions } from '@/views/finance/components/use-project-options';
import QuoteDetailDrawer from './components/quote-detail-drawer.vue';
import QuoteOperateDrawer from './components/quote-operate-drawer.vue';
import { useQuoteData } from './components/use-quote-data';

defineOptions({
  name: 'ProjectQuote'
});

interface SearchParams {
  current: number;
  size: number;
  /** 关键词：报价项名称 / 项目名称 / 客户姓名 */
  keyword: string;
  /** 所属项目 */
  projectId: number | null;
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
  /** 点击卡片后应用的筛选项目；为空表示「全部」 */
  projectId: number | null;
}

const appStore = useAppStore();

const isMobile = computed(() => appStore.isMobile);

const { loading, records, total, stats, loadList, loadStats, deleteQuotes } = useQuoteData();

const { projectOptions, loadProjectOptions } = useProjectOptions();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  projectId: null
});

/** 表格勾选的报价明细 */
const checkedRowKeys = ref<DataTableRowKey[]>([]);

const detailVisible = ref(false);
const detailId = ref<number | null>(null);

const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Quote.Item | null>(null);

/** 概览卡片数据 */
const statCards = computed<StatCard[]>(() => [
  {
    key: 'count',
    label: '报价条目',
    value: String(stats.value.totalCount ?? 0),
    icon: 'mdi:format-list-bulleted',
    color: 'primary',
    projectId: null
  },
  {
    key: 'amount',
    label: '报价总额',
    value: formatMoneyInWan(stats.value.totalAmount),
    icon: 'mdi:file-document-edit-outline',
    color: 'info',
    projectId: null
  },
  {
    key: 'project',
    label: '涉及项目',
    value: String(stats.value.projectCount ?? 0),
    icon: 'mdi:home-city-outline',
    color: 'success',
    projectId: null
  },
  {
    key: 'avg',
    label: '平均单条',
    value: formatMoney(stats.value.avgAmount),
    icon: 'mdi:scale-balance',
    color: 'warning',
    projectId: null
  }
]);

/** 分页配置 */
const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 条报价`,
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

const columns: DataTableColumns<Api.Quote.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'itemName',
    title: '报价项名称',
    width: 210,
    fixed: 'left',
    render: row =>
      h('div', { class: 'flex-y-center gap-8px' }, [
        h(
          'div',
          {
            class: 'h-28px w-28px flex-center flex-shrink-0 rounded-full text-12px font-medium',
            style: {
              backgroundColor: 'rgb(var(--primary-color) / 0.12)',
              color: 'rgb(var(--primary-color))'
            }
          },
          row.itemName.slice(0, 1)
        ),
        h('div', { class: 'min-w-0 flex-col' }, [
          h('span', { class: 'truncate font-medium' }, row.itemName),
          h('span', { class: 'text-12px text-gray-400' }, `ID ${row.id}`)
        ])
      ])
  },
  {
    key: 'projectName',
    title: '所属项目',
    minWidth: 160,
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
    key: 'unit',
    title: '单位',
    width: 80,
    render: row => row.unit || '-'
  },
  {
    key: 'price',
    title: '单价',
    width: 110,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, formatMoney(row.price))
  },
  {
    key: 'quantity',
    title: '数量',
    width: 90,
    align: 'right',
    render: row => formatNumber(row.quantity)
  },
  {
    key: 'total',
    title: '小计',
    width: 130,
    align: 'right',
    render: row => h('span', { class: 'font-mono font-medium text-primary' }, formatMoney(row.total))
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
    title: '创建时间',
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
            default: () => '确认删除该条报价吗？'
          }
        )
      ])
  }
];

function rowProps(row: Api.Quote.Item) {
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

/** 加载报价列表 */
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
  pagination.page = 1;

  await getData();
}

/** 刷新列表（同时刷新项目下拉，避免新增项目后选项过期） */
async function handleRefresh() {
  await Promise.all([refreshAll(), loadProjectOptions(true)]);
}

/** 点击概览卡片：按项目快速筛选 */
async function handleStatCardClick(card: StatCard) {
  searchParams.projectId = card.projectId;
  pagination.page = 1;

  await getData();
}

/** 概览卡片是否处于选中态（本页卡片不参与项目筛选，恒为未选中） */
function isStatCardActive(card: StatCard) {
  return card.projectId !== null && searchParams.projectId === card.projectId;
}

/** 查看报价详情 */
function handleDetail(row: Api.Quote.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

/** 新增报价 */
function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

/** 编辑报价 */
function handleEdit(row: Api.Quote.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

/** 从详情抽屉进入编辑 */
function handleEditFromDetail(row: Api.Quote.Item) {
  detailVisible.value = false;

  handleEdit(row);
}

/** 删除报价 */
async function handleDelete(row: Api.Quote.Item) {
  const success = await deleteQuotes([row.id]);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');

  await afterDeleted(1);
}

/** 批量删除报价 */
async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(key => Number(key));
  const success = await deleteQuotes(ids);

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
              placeholder="搜索报价项名称 / 项目名称 / 客户姓名"
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
                  <SvgIcon icon="mdi:file-document-plus-outline" />
                </template>
                新增报价
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
                确认删除选中的 {{ checkedRowKeys.length }} 条报价吗？
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
            :scroll-x="1420"
            remote
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </NSpace>

    <!-- 新增 / 编辑 -->
    <QuoteOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      :default-project-id="searchParams.projectId"
      @submitted="handleSubmitted"
    />

    <!-- 报价详情 -->
    <QuoteDetailDrawer v-model:visible="detailVisible" :quote-id="detailId" @edit="handleEditFromDetail" />
  </div>
</template>

<style scoped></style>
