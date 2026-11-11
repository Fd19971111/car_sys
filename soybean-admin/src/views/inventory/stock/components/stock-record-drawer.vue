<script setup lang="ts">
import { computed, h, reactive, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NTag } from 'naive-ui';
import { fetchGetStockRecordList } from '@/service/api';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { exportExcel, fetchAllPages, printTable, type PrintExportColumn } from '@/utils/print-export';

defineOptions({ name: 'StockRecordDrawer' });

interface Props {
  /** 默认筛选的项目（0 表示全部） */
  defaultProjectId?: number;
}

const props = defineProps<Props>();

const visible = defineModel<boolean>('visible', { default: false });

const loading = ref(false);
const records = ref<Api.Inventory.StockRecord[]>([]);
const recordType = ref<string | null>(null);

const pagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [20, 50],
  onUpdatePage(page: number) {
    pagination.page = page;
    void loadList();
  },
  onUpdatePageSize(size: number) {
    pagination.pageSize = size;
    pagination.page = 1;
    void loadList();
  }
});

const typeOptions: { label: string; value: string }[] = [
  { label: '全部', value: '' },
  { label: '入库', value: 'IN' },
  { label: '出库', value: 'OUT' }
];

const typeTagType: Record<string, 'success' | 'warning'> = {
  IN: 'success',
  OUT: 'warning'
};

const columns: DataTableColumns<Api.Inventory.StockRecord> = [
  { key: 'bizDate', title: '日期', width: 100, render: row => row.bizDate || '-' },
  {
    key: 'recordType',
    title: '类型',
    width: 90,
    render: row =>
      h(NTag, { size: 'small', bordered: false, type: typeTagType[row.recordType] ?? 'default' }, { default: () => (row.recordType === 'IN' ? '入库' : '出库') })
  },
  {
    key: 'materialName',
    title: '材料',
    minWidth: 160,
    ellipsis: { tooltip: true },
    render: row =>
      h('div', { class: 'flex-col gap-1px' }, [
        h('span', { class: 'font-medium' }, row.materialName || '-'),
        row.spec ? h('span', { class: 'text-11px text-gray-400' }, row.spec) : null
      ])
  },
  { key: 'projectName', title: '项目', width: 130, ellipsis: { tooltip: true }, render: row => row.projectName || '-' },
  { key: 'warehouseName', title: '仓库', width: 90, ellipsis: { tooltip: true }, render: row => row.warehouseName || '-' },
  { key: 'quantity', title: '数量', width: 80, align: 'right', render: row => `${row.quantity}${row.unit ? ` ${row.unit}` : ''}` },
  { key: 'unitPrice', title: '单价', width: 90, align: 'right', render: row => formatMoney(row.unitPrice) },
  { key: 'amount', title: '金额', width: 100, align: 'right', render: row => formatMoney(row.amount) },
  {
    key: 'finance',
    title: '财务单据',
    width: 150,
    render: row => {
      if (row.bizType === 'ADJUST') {
        return h(NTag, { size: 'small', bordered: false }, { default: () => '补录（无财务单据）' });
      }
      if (row.payoutId > 0) {
        return h(
          'span',
          { class: 'text-12px text-gray-600' },
          `付款单${row.payoutVoucherNo ? ` · 凭证 ${row.payoutVoucherNo}` : ''}`
        );
      }
      if (row.paymentId > 0) {
        return h(
          'span',
          { class: 'text-12px text-gray-600' },
          `收款单${row.paymentVoucherNo ? ` · 凭证 ${row.paymentVoucherNo}` : ''}`
        );
      }
      return h('span', { class: 'text-12px text-gray-400' }, '-');
    }
  },
  { key: 'remark', title: '备注', minWidth: 120, ellipsis: { tooltip: true }, render: row => row.remark || '-' },
  { key: 'vendorName', title: '供应商', width: 120, ellipsis: { tooltip: true }, render: row => row.vendorName || '-' }
];

const totalText = computed(() => `共 ${pagination.itemCount} 条流水`);

/** 打印 / 导出处理中 */
const printing = ref(false);
const exporting = ref(false);

/** 财务单据列的文本形式 */
function financeText(row: Api.Inventory.StockRecord): string {
  if (row.bizType === 'ADJUST') return '补录（无财务单据）';
  if (row.payoutId > 0) return `付款单${row.payoutVoucherNo ? ` · 凭证 ${row.payoutVoucherNo}` : ''}`;
  if (row.paymentId > 0) return `收款单${row.paymentVoucherNo ? ` · 凭证 ${row.paymentVoucherNo}` : ''}`;

  return '-';
}

/** 打印 / 导出列定义 */
const printExportColumns: PrintExportColumn<Api.Inventory.StockRecord>[] = [
  { title: '日期', value: row => row.bizDate || '-', width: 12 },
  { title: '类型', value: row => (row.recordType === 'IN' ? '入库' : '出库'), width: 6, align: 'center' },
  { title: '材料名称', value: row => row.materialName || '-', width: 16 },
  { title: '规格型号', value: row => row.spec || '-', width: 12 },
  { title: '项目', value: row => row.projectName || '-', width: 14 },
  { title: '仓库', value: row => row.warehouseName || '-', width: 9 },
  { title: '数量', value: row => `${row.quantity}${row.unit ? ` ${row.unit}` : ''}`, align: 'right' },
  { title: '单价', value: row => formatMoney(row.unitPrice), align: 'right' },
  { title: '金额', value: row => formatMoney(row.amount), align: 'right' },
  { title: '财务单据', value: financeText, width: 18 },
  { title: '备注', value: row => row.remark || '-', width: 16 }
];

/** 拉取全部流水（分页聚合） */
async function collectRecordRows() {
  return fetchAllPages<Api.Inventory.StockRecord>(async current => {
    const { data, error } = await fetchGetStockRecordList({
      projectId: props.defaultProjectId ?? 0,
      recordType: recordType.value || null,
      current,
      size: 100
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
    const rows = await collectRecordRows();
    await printTable({ title: '出入库流水', columns: printExportColumns, rows });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    printing.value = false;
  }
}

async function handleExport() {
  exporting.value = true;

  try {
    const rows = await collectRecordRows();
    exportExcel({ filename: '出入库流水', title: '出入库流水', columns: printExportColumns, rows });
  } catch {
    window.$message?.error('数据加载失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

async function loadList() {
  loading.value = true;

  try {
    const { data, error } = await fetchGetStockRecordList({
      projectId: props.defaultProjectId ?? 0,
      recordType: recordType.value || null,
      current: pagination.page,
      size: pagination.pageSize
    });

    if (!error && data) {
      records.value = data.records ?? [];
      pagination.itemCount = data.total ?? 0;
    }
  } finally {
    loading.value = false;
  }
}

function handleTypeChange() {
  pagination.page = 1;
  void loadList();
}

watch(visible, value => {
  if (value) {
    pagination.page = 1;
    void loadList();
  }
});
</script>

<template>
  <NDrawer v-model:show="visible" :width="1000" :mask-closable="true">
    <NDrawerContent title="出入库流水" closable :native-scrollbar="false" @close="visible = false">
      <div class="mb-12px flex items-center justify-between">
        <NRadioGroup v-model:value="recordType" size="small" @update:value="handleTypeChange">
          <NRadioButton v-for="item in typeOptions" :key="String(item.value)" :value="item.value" :label="item.label" />
        </NRadioGroup>

        <NSpace :size="8">
          <span class="text-12px text-gray-400">{{ totalText }}</span>
          <NButton size="small" :loading="printing" @click="handlePrint">
            <template #icon><SvgIcon icon="mdi:printer-outline" /></template>
            打印
          </NButton>
          <NButton size="small" :loading="exporting" @click="handleExport">
            <template #icon><SvgIcon icon="mdi:microsoft-excel" /></template>
            导出
          </NButton>
        </NSpace>
      </div>

      <NDataTable
        :columns="columns"
        :data="records"
        :loading="loading"
        :pagination="pagination"
        :scroll-x="1200"
        size="small"
        remote
        :row-key="(row: Api.Inventory.StockRecord) => row.id"
      />

      <template #footer>
        <NButton @click="visible = false">关闭</NButton>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
