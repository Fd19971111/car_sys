<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import dayjs from 'dayjs';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '../components/finance-meta';
import VoucherDetailDrawer from './components/voucher-detail-drawer.vue';
import {
  exportVoucherExcel,
  getPeriodStatusMeta,
  getSourceTypeLabel,
  getVoucherStatusMeta,
  getVoucherTypeMeta
} from './components/voucher-meta';
import { useVoucherData } from './components/use-voucher-data';

defineOptions({
  name: 'FinanceVoucher'
});

interface SearchParams {
  current: number;
  size: number;
  /** 会计期间筛选，null 为全部 */
  periodCode: string | null;
  /** 凭证状态筛选 */
  statusCode: Api.FinVoucher.VoucherStatus | null;
  /** 凭证号 / 摘要关键字 */
  keyword: string;
}

/** 表格分页 state */
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

const {
  loading,
  records,
  total,
  periods,
  periodLoading,
  initialized,
  currentPeriod,
  startPeriod,
  loadList,
  loadPeriods,
  closePeriod,
  reopenPeriod,
  exportVouchers
} = useVoucherData();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  periodCode: null,
  statusCode: null,
  keyword: ''
});

const detailVisible = ref(false);
const detailId = ref<number | null>(null);

const closing = ref(false);
const exporting = ref(false);

/** 当前选中的会计期间 */
const selectedPeriod = computed(() =>
  periods.value.find(item => item.periodCode === searchParams.periodCode)
);

const selectedStatusMeta = computed(() => getPeriodStatusMeta(selectedPeriod.value?.statusCode));

/** 期间下拉选项 */
const periodOptions = computed(() =>
  periods.value.map(item => ({
    label: `${item.periodCode}（${getPeriodStatusMeta(item.statusCode).label}）`,
    value: item.periodCode
  }))
);

const statusOptions = [
  { label: '已记账', value: 'POSTED' },
  { label: '已作废', value: 'VOID' }
];

/** 分页配置 */
const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 张凭证`,
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

const columns: DataTableColumns<Api.FinVoucher.Item> = [
  {
    key: 'voucherNo',
    title: '凭证号',
    width: 150,
    fixed: 'left',
    render: row => h('span', { class: 'font-mono' }, row.voucherNo)
  },
  {
    key: 'voucherDate',
    title: '凭证日期',
    width: 110,
    render: row => h('span', { class: 'font-mono' }, row.voucherDate)
  },
  {
    key: 'voucherTypeCode',
    title: '类型',
    width: 108,
    render: row => {
      const meta = getVoucherTypeMeta(row.voucherTypeCode);

      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'summary',
    title: '摘要',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render: row => row.summary || '-'
  },
  {
    key: 'debitTotal',
    title: '借方合计',
    width: 130,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, formatMoney(row.debitTotal))
  },
  {
    key: 'creditTotal',
    title: '贷方合计',
    width: 130,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, formatMoney(row.creditTotal))
  },
  {
    key: 'statusCode',
    title: '状态',
    width: 92,
    render: row => {
      const meta = getVoucherStatusMeta(row.statusCode);

      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'sourceType',
    title: '来源',
    width: 108,
    render: row => getSourceTypeLabel(row.sourceType)
  },
  {
    key: 'operate',
    title: '操作',
    width: 80,
    align: 'center',
    fixed: 'right',
    render: row =>
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
      )
  }
];

function rowProps(row: Api.FinVoucher.Item) {
  return {
    style: 'cursor: pointer;',
    onClick: () => handleDetail(row)
  };
}

/** 加载凭证列表 */
async function getData() {
  searchParams.current = pagination.page;
  searchParams.size = pagination.pageSize;

  await loadList({
    current: searchParams.current,
    size: searchParams.size,
    periodCode: searchParams.periodCode,
    statusCode: searchParams.statusCode,
    keyword: searchParams.keyword
  });

  pagination.itemCount = total.value;
}

/** 刷新期间 + 列表 */
async function refreshAll() {
  await loadPeriods();
  await getData();
}

/** 查询 */
async function handleSearch() {
  pagination.page = 1;

  await getData();
}

/** 重置查询条件 */
async function handleReset() {
  searchParams.periodCode = null;
  searchParams.statusCode = null;
  searchParams.keyword = '';
  pagination.page = 1;

  await getData();
}

/** 点击期间快捷筛选 */
async function handleSelectPeriod(periodCode: string | null) {
  searchParams.periodCode = periodCode;
  pagination.page = 1;

  await getData();
}

/** 查看凭证详情 */
function handleDetail(row: Api.FinVoucher.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

/** 月结 */
async function handleClose() {
  const period = selectedPeriod.value;
  if (!period) {
    return;
  }

  closing.value = true;

  try {
    const success = await closePeriod(period.periodCode);

    if (!success) {
      window.$message?.error('结账失败，请稍后重试');

      return;
    }

    window.$message?.success(`${period.periodCode} 已结账，该期间凭证已锁定`);
    await refreshAll();
  } finally {
    closing.value = false;
  }
}

/** 反结账 */
async function handleReopen() {
  const period = selectedPeriod.value;
  if (!period) {
    return;
  }

  closing.value = true;

  try {
    const success = await reopenPeriod(period.periodCode);

    if (!success) {
      window.$message?.error('反结账失败，请稍后重试');

      return;
    }

    window.$message?.success(`${period.periodCode} 已反结账，可以继续修改该期间凭证`);
    await refreshAll();
  } finally {
    closing.value = false;
  }
}

/** 导出 Excel */
async function handleExport() {
  exporting.value = true;

  try {
    const rows = await exportVouchers({
      periodCode: searchParams.periodCode,
      statusCode: searchParams.statusCode,
      keyword: searchParams.keyword
    });

    if (!rows) {
      window.$message?.error('导出失败，请稍后重试');

      return;
    }
    if (!rows.length) {
      window.$message?.warning('当前筛选条件下没有可导出的凭证');

      return;
    }

    const periodLabel = searchParams.periodCode || '全部期间';
    exportVoucherExcel(rows, `记账凭证-${periodLabel}-${dayjs().format('YYYYMMDDHHmm')}`);
    window.$message?.success(`已导出 ${rows.length} 条分录`);
  } finally {
    exporting.value = false;
  }
}

onMounted(async () => {
  await loadPeriods();

  // 默认定位到当前会计期间，没有凭证时退回到最近一个期间
  const current = periods.value.find(item => item.periodCode === currentPeriod.value);
  const withVoucher = [...periods.value].reverse().find(item => item.voucherCount > 0);

  searchParams.periodCode = current?.periodCode ?? withVoucher?.periodCode ?? null;
  pagination.page = 1;

  await getData();
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <NAlert v-if="!initialized" type="warning" :bordered="false">
        财务还没有完成初始化，暂时无法查看凭证。请先到「财务基础资料」完成初始化配置。
      </NAlert>

      <!-- 会计期间 / 月结 -->
      <NCard :bordered="false" class="card-wrapper" size="small" :loading="periodLoading">
        <template #header>
          <div class="flex-y-center gap-8px">
            <span>会计期间</span>
            <NTag v-if="selectedPeriod" :type="selectedStatusMeta.tagType" size="small" :bordered="false">
              {{ selectedStatusMeta.label }}
            </NTag>
          </div>
        </template>
        <template #header-extra>
          <NSpace :size="8">
            <NPopconfirm v-if="selectedPeriod?.canClose" @positive-click="handleClose">
              <template #trigger>
                <NButton type="primary" ghost size="small" :loading="closing">
                  <template #icon>
                    <SvgIcon icon="mdi:lock-check-outline" />
                  </template>
                  月结
                </NButton>
              </template>
              确认对 {{ selectedPeriod?.periodCode }} 月结吗？结账后该期间不能再新增或修改凭证。
            </NPopconfirm>
            <NPopconfirm v-if="selectedPeriod?.canReopen" @positive-click="handleReopen">
              <template #trigger>
                <NButton size="small" :loading="closing">
                  <template #icon>
                    <SvgIcon icon="mdi:lock-open-variant-outline" />
                  </template>
                  反结账
                </NButton>
              </template>
              确认反结账 {{ selectedPeriod?.periodCode }} 吗？反结账后可以继续修改该期间凭证。
            </NPopconfirm>
            <NButton size="small" :loading="periodLoading" @click="refreshAll">
              <template #icon>
                <SvgIcon icon="mdi:refresh" />
              </template>
              刷新
            </NButton>
          </NSpace>
        </template>

        <div
          v-if="selectedPeriod"
          class="mt-10px flex flex-wrap items-center gap-16px text-12px text-gray-500"
        >
          <span>已记账 {{ selectedPeriod.voucherCount }} 张</span>
          <span>借方 {{ formatMoney(selectedPeriod.debitTotal) }}</span>
          <span>贷方 {{ formatMoney(selectedPeriod.creditTotal) }}</span>
          <span v-if="selectedPeriod.beforeStart" class="text-warning">早于财务启用期间，仅作查看</span>
        </div>

        <p class="mt-8px text-12px text-gray-400">
          月结会锁定该期间的全部凭证；空期间会在月结时自动关闭。反结账只支持最近一个有凭证的已结账期间。
        </p>
      </NCard>

      <!-- 凭证列表 -->
      <NCard :bordered="false" class="card-wrapper" size="small">
        <NSpace vertical :size="12">
          <div class="flex flex-wrap items-center gap-8px">
            <NSelect
              v-model:value="searchParams.periodCode"
              class="w-190px flex-shrink-0"
              :options="periodOptions"
              placeholder="会计期间"
              clearable
            />
            <NSelect
              v-model:value="searchParams.statusCode"
              class="w-120px flex-shrink-0"
              :options="statusOptions"
              placeholder="凭证状态"
              clearable
            />
            <NInput
              v-model:value="searchParams.keyword"
              class="w-220px min-w-140px flex-shrink"
              placeholder="搜索凭证号 / 摘要"
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
              <NButton size="small" type="primary" ghost :loading="exporting" @click="handleExport">
                <template #icon>
                  <SvgIcon icon="mdi:file-delimited-outline" />
                </template>
                导出 Excel
              </NButton>
              <NButton size="small" :loading="loading" @click="getData">
                <template #icon>
                  <SvgIcon icon="mdi:refresh" />
                </template>
                刷新
              </NButton>
            </div>
          </div>

          <NDataTable
            :columns="columns"
            :data="records"
            :loading="loading"
            :pagination="tablePagination"
            :row-key="row => row.id"
            :row-props="rowProps"
            :scroll-x="1160"
            remote
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </NSpace>

    <VoucherDetailDrawer v-model:visible="detailVisible" :voucher-id="detailId" />
  </div>
</template>

<style scoped>
.period-strip {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.period-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 11px;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.65);
  cursor: pointer;
  background-color: rgb(var(--base-text-color) / 0.04);
  border: 1px solid transparent;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.period-chip:hover {
  background-color: rgb(var(--primary-color) / 0.08);
}

.period-chip--active {
  font-weight: 600;
  color: rgb(var(--primary-color));
  background-color: rgb(var(--primary-color) / 0.12);
  border-color: rgb(var(--primary-color) / 0.4);
}

.period-chip__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.period-chip__dot--open {
  background-color: rgb(var(--success-color));
}

.period-chip__dot--closed {
  background-color: rgb(var(--base-text-color) / 0.3);
}

.period-chip__dot--future {
  background-color: rgb(var(--info-color));
}

.period-chip__code {
  font-variant-numeric: tabular-nums;
}

.period-chip__count {
  padding: 0 5px;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.45);
  background-color: rgb(var(--base-text-color) / 0.06);
  border-radius: 999px;
}
</style>
