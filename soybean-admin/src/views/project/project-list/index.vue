<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import ProjectDetailDrawer from './components/project-detail-drawer.vue';
import ProjectOperateDrawer from './components/project-operate-drawer.vue';
import {
  PROJECT_STATUS_OPTIONS,
  formatDate,
  formatMoney,
  formatMoneyInWan,
  getProjectProgress,
  getProjectStatusMeta,
  getProjectStatusStyle
} from './components/project-meta';
import { useProjectData } from './components/use-project-data';

defineOptions({
  name: 'ProjectProjectList'
});

interface SearchParams {
  current: number;
  size: number;
  /** 关键词：项目名称 / 施工地址 / 客户姓名 */
  keyword: string;
  /** 项目状态 */
  status: Api.Project.Status | null;
  /** 所属客户 */
  customerId: number | null;
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
  /** 点击卡片后应用的筛选状态；为空表示「全部」 */
  status: Api.Project.Status | null;
}

const appStore = useAppStore();

const isMobile = computed(() => appStore.isMobile);

const { loading, records, total, stats, options, loadList, loadStats, loadOptions, deleteProjects } = useProjectData();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  status: null,
  customerId: null
});

/** 表格勾选的项目 */
const checkedRowKeys = ref<DataTableRowKey[]>([]);

const detailVisible = ref(false);
const detailId = ref<number | null>(null);

const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Project.Item | null>(null);

/** 概览卡片数据 */
const statCards = computed<StatCard[]>(() => [
  {
    key: 'total',
    label: '项目总数',
    value: String(stats.value.total ?? 0),
    icon: 'mdi:clipboard-text-outline',
    color: 'primary',
    status: null
  },
  {
    key: 'pending',
    label: '待开工',
    value: String(stats.value.pending ?? 0),
    icon: 'mdi:calendar-clock',
    color: 'info',
    status: 1
  },
  {
    key: 'building',
    label: '施工中',
    value: String(stats.value.building ?? 0),
    icon: 'mdi:hammer-wrench',
    color: 'warning',
    status: 2
  },
  {
    key: 'finished',
    label: '已完工',
    value: String(stats.value.finished ?? 0),
    icon: 'mdi:check-circle-outline',
    color: 'success',
    status: 3
  },
  {
    key: 'settled',
    label: '已结算',
    value: String(stats.value.settled ?? 0),
    icon: 'mdi:cash-check',
    color: 'primary',
    status: 4
  },
  {
    key: 'amount',
    label: '合同总金额',
    value: formatMoneyInWan(stats.value.totalAmount),
    icon: 'mdi:currency-cny',
    color: 'primary',
    status: null
  }
]);

/** 客户筛选下拉选项 */
const customerFilterOptions = computed(() => options.value.customers);

/** 分页配置 */
const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 个项目`,
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

const columns: DataTableColumns<Api.Project.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'name',
    title: '项目名称',
    width: 210,
    fixed: 'left',
    render: row => {
      const meta = getProjectStatusMeta(row.status);

      return h('div', { class: 'flex-y-center gap-8px' }, [
        h(
          'div',
          {
            class: 'h-28px w-28px flex-center flex-shrink-0 rounded-full text-12px font-medium',
            style: getProjectStatusStyle(row.status),
            title: meta.label
          },
          row.name.slice(0, 1)
        ),
        h('div', { class: 'min-w-0 flex-col' }, [
          h('span', { class: 'truncate font-medium' }, row.name),
          h('span', { class: 'text-12px text-gray-400' }, `ID ${row.id}`)
        ])
      ]);
    }
  },
  {
    key: 'customerName',
    title: '所属客户',
    width: 120,
    ellipsis: { tooltip: true },
    render: row => row.customerName || '-'
  },
  {
    key: 'address',
    title: '施工地址',
    minWidth: 180,
    ellipsis: { tooltip: true },
    render: row => row.address || '-'
  },
  {
    key: 'managerName',
    title: '负责人',
    width: 100,
    render: row => row.managerName || '-'
  },
  {
    key: 'status',
    title: '项目状态',
    width: 96,
    render: row => {
      const meta = getProjectStatusMeta(row.status);

      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'progress',
    title: '施工进度',
    width: 150,
    render: row => {
      const percent = getProjectProgress(row.status, row.startDate, row.endDate);
      const meta = getProjectStatusMeta(row.status);

      return h('div', { class: 'flex-col gap-4px' }, [
        h('div', { class: 'h-6px w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700' }, [
          h('div', {
            class: 'h-full rounded-full transition-all',
            style: {
              width: `${percent}%`,
              backgroundColor: `rgb(var(--${meta.color}-color))`
            }
          })
        ]),
        h('span', { class: 'text-12px text-gray-400' }, `${percent}%`)
      ]);
    }
  },
  {
    key: 'amount',
    title: '合同金额',
    width: 130,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, formatMoney(row.amount))
  },
  {
    key: 'period',
    title: '工期',
    width: 200,
    render: row => h('span', { class: 'text-13px text-gray-500' }, `${formatDate(row.startDate)} ~ ${formatDate(row.endDate)}`)
  },
  {
    key: 'latestProgress',
    title: '最近进度',
    minWidth: 180,
    ellipsis: { tooltip: true },
    render: row => row.latestProgress || h('span', { class: 'text-gray-400' }, '暂无进度记录')
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
            default: () => '确认删除该项目吗？'
          }
        )
      ])
  }
];

function rowProps(row: Api.Project.Item) {
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

/** 加载项目列表 */
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
  searchParams.status = null;
  searchParams.customerId = null;
  pagination.page = 1;

  await getData();
}

/** 刷新列表 */
async function handleRefresh() {
  await refreshAll();
}

/** 点击概览卡片：按状态快速筛选 */
async function handleStatCardClick(card: StatCard) {
  searchParams.status = card.status;
  pagination.page = 1;

  await getData();
}

/** 概览卡片是否处于选中态（「项目总数 / 合同总金额」卡片不参与状态筛选） */
function isStatCardActive(card: StatCard) {
  return card.status !== null && searchParams.status === card.status;
}

/** 查看项目详情（详情抽屉按 id 从后端拉取完整信息） */
function handleDetail(row: Api.Project.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

/** 新增项目 */
function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

/** 编辑项目 */
function handleEdit(row: Api.Project.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

/** 从详情抽屉进入编辑（payload 为后端返回的最新详情） */
function handleEditFromDetail(row: Api.Project.Item) {
  detailVisible.value = false;

  handleEdit(row);
}

/** 删除项目 */
async function handleDelete(row: Api.Project.Item) {
  const success = await deleteProjects([row.id]);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');

  await afterDeleted(1);
}

/** 批量删除项目 */
async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(key => Number(key));
  const success = await deleteProjects(ids);

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
  await Promise.all([refreshAll(), loadOptions()]);
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <!-- 概览统计：点击卡片按状态快速筛选 -->
      <div class="flex flex-wrap gap-12px">
        <NCard
          v-for="card in statCards"
          :key="card.key"
          :bordered="false"
          class="card-wrapper min-w-160px flex-1 cursor-pointer border-2 border-transparent transition-all hover:shadow-md"
          :style="isStatCardActive(card) ? { borderColor: 'rgb(var(--primary-color))' } : undefined"
          size="small"
          @click="handleStatCardClick(card)"
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

      <!-- 搜索 / 筛选 / 操作按钮：全部排在同一行 -->
      <NCard :bordered="false" class="card-wrapper" size="small">
        <NSpace vertical :size="12">
          <div class="flex flex-wrap items-center gap-8px">
            <NInput
              v-model:value="searchParams.keyword"
              class="w-220px min-w-140px flex-shrink"
              placeholder="搜索项目名称 / 施工地址 / 客户姓名"
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
              :options="PROJECT_STATUS_OPTIONS"
              placeholder="项目状态"
              clearable
            />
            <NSelect
              v-model:value="searchParams.customerId"
              class="w-140px flex-shrink-0"
              :options="customerFilterOptions"
              placeholder="所属客户"
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
                  <SvgIcon icon="mdi:folder-plus-outline" />
                </template>
                新增项目
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
                确认删除选中的 {{ checkedRowKeys.length }} 个项目吗？
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
            :scroll-x="1720"
            remote
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </NSpace>

    <!-- 新增 / 编辑 -->
    <ProjectOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      @submitted="handleSubmitted"
    />

    <!-- 项目详情 -->
    <ProjectDetailDrawer v-model:visible="detailVisible" :project-id="detailId" @edit="handleEditFromDetail" />
  </div>
</template>

<style scoped></style>
