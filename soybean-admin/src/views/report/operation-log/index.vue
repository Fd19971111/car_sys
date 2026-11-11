<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NTag } from 'naive-ui';
import { fetchGetOperationLogList } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';

defineOptions({
  name: 'ReportOperationLog'
});

interface SearchParams {
  /** 业务模块，null 表示全部 */
  module: string | null;
  /** 动作，'' 表示全部 */
  action: '' | Api.OpLog.Action;
  keyword: string;
  /** 开始日期 */
  startDate: string | null;
  /** 结束日期 */
  endDate: string | null;
}

const appStore = useAppStore();

const isMobile = computed(() => appStore.isMobile);

const loading = ref(false);
const records = ref<Api.OpLog.Item[]>([]);
const total = ref(0);
const modules = ref<string[]>([]);

const searchParams = reactive<SearchParams>({
  module: null,
  action: '',
  keyword: '',
  startDate: null,
  endDate: null
});

/** 动作下拉 */
const actionOptions = [
  { label: '全部动作', value: '' },
  { label: '新增', value: 'CREATE' },
  { label: '修改', value: 'UPDATE' },
  { label: '删除', value: 'DELETE' }
];

/** 模块下拉：接口返回的已用模块 + 常用模块兜底 */
const moduleOptions = computed(() => {
  const preset = ['收款', '付款', '项目', '库存', '客户', '供应商'];
  const merged = Array.from(new Set([...modules.value, ...preset]));

  return [{ label: '全部模块', value: '' }, ...merged.map(item => ({ label: item, value: item }))];
});

/** 动作标签配色 */
const ACTION_TYPE: Record<Api.OpLog.Action, NaiveUI.ThemeColor> = {
  CREATE: 'success',
  UPDATE: 'info',
  DELETE: 'error'
};

/** 分页配置 */
const pagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [20, 50, 100],
  prefix: (info: { itemCount?: number }) => `共 ${info.itemCount ?? 0} 条操作`,
  onUpdatePage(page: number) {
    pagination.page = page;

    void load();
  },
  onUpdatePageSize(pageSize: number) {
    pagination.pageSize = pageSize;
    pagination.page = 1;

    void load();
  }
});

const tablePagination = computed(() => ({
  ...pagination,
  pageSlot: isMobile.value ? 3 : 9,
  prefix: isMobile.value ? undefined : pagination.prefix
}));

const columns: DataTableColumns<Api.OpLog.Item> = [
  {
    key: 'createdAt',
    title: '操作时间',
    width: 155,
    render: row => h('span', { class: 'font-mono' }, row.createdAt)
  },
  {
    key: 'userName',
    title: '操作人',
    width: 100,
    ellipsis: { tooltip: true },
    render: row => row.userName || '系统'
  },
  {
    key: 'module',
    title: '模块',
    width: 90,
    render: row => h(NTag, { size: 'small', bordered: false }, { default: () => row.module || '-' })
  },
  {
    key: 'action',
    title: '动作',
    width: 84,
    render: row =>
      h(
        NTag,
        { size: 'small', bordered: false, type: ACTION_TYPE[row.action] ?? 'default' },
        { default: () => row.actionName || row.action }
      )
  },
  {
    key: 'title',
    title: '操作内容',
    minWidth: 320,
    ellipsis: { tooltip: true },
    render: row => row.title || '-'
  },
  {
    key: 'amount',
    title: '涉及金额',
    width: 130,
    align: 'right',
    render: row =>
      row.amount
        ? h('span', { class: 'font-mono font-medium' }, formatMoney(row.amount))
        : h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'ip',
    title: '来源IP',
    width: 130,
    render: row => h('span', { class: 'font-mono text-12px text-gray-500' }, row.ip || '—')
  }
];

/** 详情抽屉 */
const detailVisible = ref(false);
const detailRow = ref<Api.OpLog.Item | null>(null);

/** 详情里格式化后的 JSON */
const detailText = computed(() => {
  const raw = detailRow.value?.detail;

  if (!raw) {
    return '';
  }

  try {
    return JSON.stringify(JSON.parse(raw), null, 2);
  } catch {
    return raw;
  }
});

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

const printExportColumns: PrintExportColumn<Api.OpLog.Item>[] = [
  { title: '操作时间', value: row => row.createdAt, width: 20 },
  { title: '操作人', value: row => row.userName || '系统', width: 12 },
  { title: '模块', value: row => row.module || '-', align: 'center', width: 10 },
  { title: '动作', value: row => row.actionName || row.action, align: 'center', width: 8 },
  { title: '操作内容', value: row => row.title || '-', width: 46 },
  { title: '涉及金额', value: row => (row.amount ? formatMoney(row.amount) : '—'), align: 'right', width: 14 }
];

/** 打印 / 导出用的筛选说明 */
function filterSubtitle() {
  const parts = [
    `模块：${searchParams.module || '全部'}`,
    `动作：${actionOptions.find(item => item.value === searchParams.action)?.label ?? '全部'}`
  ];

  if (searchParams.startDate || searchParams.endDate) {
    parts.push(`日期：${searchParams.startDate || '不限'} ~ ${searchParams.endDate || '不限'}`);
  }

  if (searchParams.keyword.trim()) {
    parts.push(`关键词：${searchParams.keyword.trim()}`);
  }

  return parts.join('　|　');
}

/** 拉取当前筛选下的全部日志（分页聚合，打印 / 导出用） */
async function collectRows() {
  return fetchAllPages<Api.OpLog.Item>(async current => {
    const { data, error } = await fetchGetOperationLogList({
      current,
      size: 100,
      module: searchParams.module,
      action: searchParams.action || null,
      keyword: searchParams.keyword,
      startDate: searchParams.startDate,
      endDate: searchParams.endDate
    });

    if (error || !data) {
      throw new Error('数据加载失败');
    }

    return { records: data.records ?? [], total: data.total ?? 0 };
  });
}

async function handlePrint() {
  printing.value = true;

  try {
    const rows = await collectRows();
    await printTable({
      title: '操作日志',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows,
      summary: [
        { label: '操作笔数', value: `${rows.length} 条` },
        {
          label: '涉及金额',
          value: formatMoney(rows.reduce((sum, row) => sum + row.amount, 0))
        }
      ]
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
    const rows = await collectRows();
    exportExcel({
      filename: '操作日志',
      title: '操作日志',
      subtitle: filterSubtitle(),
      columns: printExportColumns,
      rows
    });
    window.$message?.success(`已导出 ${rows.length} 条操作`);
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

/** 加载操作日志 */
async function load() {
  loading.value = true;

  try {
    const { data, error } = await fetchGetOperationLogList({
      current: pagination.page,
      size: pagination.pageSize,
      module: searchParams.module,
      action: searchParams.action || null,
      keyword: searchParams.keyword,
      startDate: searchParams.startDate,
      endDate: searchParams.endDate
    });

    if (error || !data) {
      records.value = [];
      total.value = 0;

      return;
    }

    records.value = data.records ?? [];
    total.value = data.total ?? 0;
    pagination.itemCount = total.value;
    modules.value = data.modules ?? [];
  } finally {
    loading.value = false;
  }
}

/** 查询 */
async function handleSearch() {
  pagination.page = 1;

  await load();
}

/** 重置 */
async function handleReset() {
  searchParams.module = null;
  searchParams.action = '';
  searchParams.keyword = '';
  searchParams.startDate = null;
  searchParams.endDate = null;
  pagination.page = 1;

  await load();
}

/** 打开详情 */
function openDetail(row: Api.OpLog.Item) {
  detailRow.value = row;
  detailVisible.value = true;
}

/** 行事件：点击查看详情 */
function rowProps(row: Api.OpLog.Item) {
  return {
    style: 'cursor: pointer;',
    onClick: () => openDetail(row)
  };
}

onMounted(() => {
  void load();
});
</script>

<template>
  <div class="report-page">
    <NCard :bordered="false" class="filter-card" size="small">
      <NSpace vertical :size="12">
        <div class="flex flex-wrap items-center gap-8px">
          <NSelect
            v-model:value="searchParams.module"
            class="w-140px flex-shrink-0"
            :options="moduleOptions"
            placeholder="业务模块"
            clearable
          />
          <NSelect v-model:value="searchParams.action" class="w-130px flex-shrink-0" :options="actionOptions" />
          <NDatePicker
            v-model:formatted-value="searchParams.startDate"
            class="w-140px flex-shrink-0"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="开始日期"
            clearable
          />
          <NDatePicker
            v-model:formatted-value="searchParams.endDate"
            class="w-140px flex-shrink-0"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="结束日期"
            clearable
          />
          <NInput
            v-model:value="searchParams.keyword"
            class="w-220px min-w-140px flex-shrink"
            placeholder="搜索操作内容 / 操作人"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <SvgIcon icon="mdi:magnify" class="text-gray-400" />
            </template>
          </NInput>
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
            <NButton size="small" :loading="printing" @click="handlePrint">
              <template #icon>
                <SvgIcon icon="mdi:printer-outline" />
              </template>
              打印
            </NButton>
            <NButton size="small" :loading="exporting" @click="handleExport">
              <template #icon>
                <SvgIcon icon="mdi:microsoft-excel" />
              </template>
              导出
            </NButton>
            <NButton size="small" :loading="loading" @click="load">
              <template #icon>
                <SvgIcon icon="mdi:refresh" />
              </template>
              刷新
            </NButton>
          </div>
        </div>

        <NAlert type="default" :bordered="false" :show-icon="true">
          记录收付款等关键业务动作，可追溯「谁在什么时候改了什么」。当前已接入：收款、付款的新增 / 修改 / 删除。
        </NAlert>

        <NDataTable
          :columns="columns"
          :data="records"
          :loading="loading"
          :pagination="tablePagination"
          :row-key="row => row.id"
          :row-props="rowProps"
          :scroll-x="1020"
          remote
          size="small"
          striped
        />
      </NSpace>
    </NCard>

    <!-- 日志详情 -->
    <NDrawer v-model:show="detailVisible" :width="isMobile ? '100%' : 520">
      <NDrawerContent title="操作日志详情" closable :native-scrollbar="false" @close="detailVisible = false">
        <NDescriptions v-if="detailRow" :column="1" bordered size="small" label-placement="left">
          <NDescriptionsItem label="操作时间">{{ detailRow.createdAt }}</NDescriptionsItem>
          <NDescriptionsItem label="操作人">{{ detailRow.userName || '系统' }}</NDescriptionsItem>
          <NDescriptionsItem label="模块 / 动作">
            {{ detailRow.module || '-' }} / {{ detailRow.actionName || detailRow.action }}
          </NDescriptionsItem>
          <NDescriptionsItem label="操作内容">{{ detailRow.title }}</NDescriptionsItem>
          <NDescriptionsItem label="涉及金额">
            {{ detailRow.amount ? formatMoney(detailRow.amount) : '—' }}
          </NDescriptionsItem>
          <NDescriptionsItem label="关联单据">
            {{ detailRow.targetType ? `${detailRow.targetType} #${detailRow.targetId}` : '—' }}
          </NDescriptionsItem>
          <NDescriptionsItem label="来源IP">{{ detailRow.ip || '—' }}</NDescriptionsItem>
        </NDescriptions>

        <template v-if="detailText">
          <NDivider class="!my-12px">变更详情</NDivider>
          <pre class="detail-json">{{ detailText }}</pre>
        </template>
      </NDrawerContent>
    </NDrawer>
  </div>
</template>

<style scoped>
.report-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
}

.filter-card {
  border-radius: 12px;
}

.detail-json {
  margin: 0;
  padding: 12px;
  overflow: auto;
  font-family: theme('fontFamily.mono');
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  background-color: rgb(var(--base-color) / 0.06);
  border-radius: 8px;
}

@media (max-width: 768px) {
  .report-page {
    padding: 12px;
  }
}
</style>
