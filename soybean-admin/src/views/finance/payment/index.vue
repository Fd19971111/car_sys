<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import { formatDate, formatMoney, formatMoneyInWan } from '../components/finance-meta';
import PaymentDetailDrawer from './components/payment-detail-drawer.vue';
import PaymentOperateDrawer from './components/payment-operate-drawer.vue';
import { PAYMENT_TYPE_OPTIONS, getPaymentTypeLabel, getPaymentTypeMeta } from './components/payment-meta';
import { usePaymentData } from './components/use-payment-data';

defineOptions({ name: 'FinancePayment' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
  paymentType: string | null;
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
const { loading, records, total, stats, loadList, loadStats, deletePayments } = usePaymentData();

const searchParams = reactive<SearchParams>({ current: 1, size: 10, keyword: '', paymentType: null });
const checkedRowKeys = ref<DataTableRowKey[]>([]);
const detailVisible = ref(false);
const detailId = ref<number | null>(null);
const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Payment.Item | null>(null);

const statCards = computed(() => [
  { label: '收款笔数', value: String(stats.value.totalCount ?? 0), icon: 'mdi:receipt-text-outline', color: 'primary' },
  { label: '收款总额', value: formatMoneyInWan(stats.value.totalAmount), icon: 'mdi:cash-multiple', color: 'success' },
  { label: '本月收款', value: formatMoneyInWan(stats.value.monthAmount), icon: 'mdi:calendar-month-outline', color: 'info' },
  { label: '本月笔数', value: String(stats.value.monthCount ?? 0), icon: 'mdi:counter', color: 'warning' }
]);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 笔收款`,
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

const columns: DataTableColumns<Api.Payment.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'paymentDate',
    title: '收款日期',
    width: 130,
    fixed: 'left',
    render: row => {
      const meta = getPaymentTypeMeta(row.paymentType);
      return h('div', { class: 'flex-y-center gap-8px' }, [
        h('div', { class: 'h-8px w-8px flex-shrink-0 rounded-full', style: { backgroundColor: `rgb(var(--${meta.color}-color))` } }),
        h('span', { class: 'font-mono' }, formatDate(row.paymentDate))
      ]);
    }
  },
  {
    key: 'paymentType',
    title: '款项类型',
    width: 120,
    render: row => {
      const meta = getPaymentTypeMeta(row.paymentType);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => getPaymentTypeLabel(row.paymentType) });
    }
  },
  {
    key: 'amount',
    title: '收款金额',
    width: 150,
    align: 'right',
    render: row => h('span', { class: 'font-mono font-medium text-success' }, formatMoney(row.amount))
  },
  { key: 'fundAccountName', title: '收款账户', width: 160, ellipsis: { tooltip: true }, render: row => row.fundAccountName || '-' },
  {
    key: 'voucherNo',
    title: '记账凭证',
    width: 160,
    render: row =>
      row.voucherNo
        ? h(NTag, { size: 'small', bordered: false, type: 'success' }, { default: () => row.voucherNo })
        : h(NTag, { size: 'small', bordered: false }, { default: () => '未记账' })
  },
  { key: 'remark', title: '备注', minWidth: 160, ellipsis: { tooltip: true }, render: row => row.remark || '-' },
  { key: 'createdAt', title: '登记时间', width: 140, render: row => formatDate(row.createdAt) },
  {
    key: 'operate',
    title: '操作',
    width: 150,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-4px' }, [
        h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => handleDetail(row) }, { default: () => '查看' }),
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
          trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
          default: () => '确认删除该收款记录吗？'
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
  await Promise.all([getData(), loadStats()]);
}

async function handleReset() {
  searchParams.keyword = '';
  searchParams.paymentType = null;
  pagination.page = 1;
  await getData();
}

function handleDetail(row: Api.Payment.Item) {
  detailId.value = row.id;
  detailVisible.value = true;
}

function handleEdit(row: Api.Payment.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  detailVisible.value = false;
  drawerVisible.value = true;
}

function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

async function handleDelete(row: Api.Payment.Item) {
  const success = await deletePayments([row.id]);
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
  const success = await deletePayments(ids);
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
  await Promise.all([getData(), loadStats()]);
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
          <div>
            <p class="text-12px text-gray-400">{{ card.label }}</p>
            <p class="mt-2px text-20px font-semibold">{{ card.value }}</p>
          </div>
        </div>
      </div>

      <NCard :bordered="false" class="card-wrapper" size="small">
        <NSpace vertical :size="12">
          <div class="flex flex-wrap items-center gap-8px">
            <NInput v-model:value="searchParams.keyword" class="w-220px min-w-150px" placeholder="搜索款项类型 / 备注" clearable @keyup.enter="handleSearch">
              <template #prefix><SvgIcon icon="mdi:magnify" class="text-gray-400" /></template>
            </NInput>
            <NSelect v-model:value="searchParams.paymentType" class="w-130px" :options="PAYMENT_TYPE_OPTIONS" placeholder="款项类型" clearable />
            <NButton type="primary" ghost size="small" @click="handleSearch">
              <template #icon><SvgIcon icon="mdi:magnify" /></template>查询
            </NButton>
            <NButton size="small" @click="handleReset">
              <template #icon><SvgIcon icon="mdi:backspace-outline" /></template>重置
            </NButton>
            <div class="ml-auto flex items-center gap-8px">
              <NButton size="small" type="primary" ghost @click="handleAdd">
                <template #icon><SvgIcon icon="mdi:plus-circle-outline" /></template>新增收款
              </NButton>
              <NPopconfirm @positive-click="handleBatchDelete">
                <template #trigger>
                  <NButton size="small" ghost type="error" :disabled="!checkedRowKeys.length">
                    <template #icon><SvgIcon icon="mdi:delete-outline" /></template>批量删除
                  </NButton>
                </template>
                确认删除选中的 {{ checkedRowKeys.length }} 笔收款吗？
              </NPopconfirm>
              <NButton size="small" @click="handleSearch">
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
            :scroll-x="1200"
            remote
            size="small"
            striped
          />
        </NSpace>
      </NCard>
    </NSpace>

    <PaymentOperateDrawer v-model:visible="drawerVisible" :operate-type="operateType" :row-data="editingData" @submitted="handleSubmitted" />
    <PaymentDetailDrawer v-model:visible="detailVisible" :payment-id="detailId" @edit="handleEdit" />
  </div>
</template>
