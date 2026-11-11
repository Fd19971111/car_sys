<script setup lang="ts">
import { computed, h, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NTag } from 'naive-ui';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '../../components/finance-meta';
import { getSourceTypeLabel, getVoucherStatusMeta, getVoucherTypeMeta } from './voucher-meta';
import { useVoucherData } from './use-voucher-data';

defineOptions({
  name: 'VoucherDetailDrawer'
});

interface Props {
  /** 当前查看的凭证ID */
  voucherId?: number | null;
}

const props = defineProps<Props>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 780));

const { getDetail } = useVoucherData();

const loading = ref(false);
const detail = ref<Api.FinVoucher.Detail | null>(null);

const statusMeta = computed(() => getVoucherStatusMeta(detail.value?.statusCode));
const typeMeta = computed(() => getVoucherTypeMeta(detail.value?.voucherTypeCode));

const columns: DataTableColumns<Api.FinVoucher.Entry> = [
  { key: 'lineNo', title: '行号', width: 56, align: 'center' },
  {
    key: 'directionCode',
    title: '方向',
    width: 64,
    align: 'center',
    render: row =>
      h(
        NTag,
        { size: 'small', bordered: false, type: row.directionCode === 'DEBIT' ? 'primary' : 'warning' },
        { default: () => (row.directionCode === 'DEBIT' ? '借' : '贷') }
      )
  },
  {
    key: 'accountName',
    title: '会计科目',
    minWidth: 150,
    ellipsis: { tooltip: true },
    render: row => `${row.accountCode} ${row.accountName}`
  },
  {
    key: 'summary',
    title: '分录摘要',
    minWidth: 150,
    ellipsis: { tooltip: true },
    render: row => row.summary || '-'
  },
  {
    key: 'amount',
    title: '借方金额',
    width: 120,
    align: 'right',
    render: row =>
      row.directionCode === 'DEBIT' ? h('span', { class: 'font-mono' }, formatMoney(row.amount)) : h('span', { class: 'text-gray-300' }, '—')
  },
  {
    key: 'creditAmount',
    title: '贷方金额',
    width: 120,
    align: 'right',
    render: row =>
      row.directionCode === 'CREDIT' ? h('span', { class: 'font-mono' }, formatMoney(row.amount)) : h('span', { class: 'text-gray-300' }, '—')
  }
];

watch(
  [visible, () => props.voucherId],
  async ([isVisible]) => {
    if (!isVisible) {
      return;
    }

    if (!props.voucherId) {
      detail.value = null;

      return;
    }

    loading.value = true;

    try {
      detail.value = await getDetail(props.voucherId);
    } finally {
      loading.value = false;
    }
  },
  { immediate: true }
);

function handleClose() {
  visible.value = false;
  detail.value = null;
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth">
    <NDrawerContent title="记账凭证" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="凭证不存在" />

        <template v-else-if="detail">
          <div class="mb-14px flex-y-center gap-8px">
            <span class="text-16px font-medium">{{ detail.voucherNo }}</span>
            <NTag :type="typeMeta.tagType" size="small" :bordered="false">{{ typeMeta.label }}</NTag>
            <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
          </div>

          <NDescriptions :column="1" label-placement="left" :label-width="96" bordered size="small">
            <NDescriptionsItem label="凭证日期">{{ detail.voucherDate }}</NDescriptionsItem>
            <NDescriptionsItem label="会计期间">{{ detail.periodCode }}</NDescriptionsItem>
            <NDescriptionsItem label="凭证摘要">{{ detail.summary || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="来源">
              {{ getSourceTypeLabel(detail.sourceType) }}
              <span v-if="detail.sourceId" class="text-gray-400">#{{ detail.sourceId }}</span>
            </NDescriptionsItem>
            <NDescriptionsItem label="制单时间">{{ detail.createdAt || '-' }}</NDescriptionsItem>
          </NDescriptions>

          <div class="mt-16px overflow-hidden rounded-8px border border-gray-200 dark:border-gray-700">
            <NDataTable
              :columns="columns"
              :data="detail.entries"
              :row-key="row => row.lineNo"
              :pagination="false"
              :scroll-x="720"
              size="small"
              :bordered="false"
            />
            <div class="flex-y-center justify-end gap-24px bg-gray-50 px-12px py-8px text-12px dark:bg-gray-800">
              <span>借方合计 <strong class="font-mono">{{ formatMoney(detail.debitTotal) }}</strong></span>
              <span>贷方合计 <strong class="font-mono">{{ formatMoney(detail.creditTotal) }}</strong></span>
            </div>
          </div>

          <p class="mt-10px text-12px text-gray-400">
            凭证由业务单据自动生成：收款 / 付款登记或修改时同步更新，作废后不再参与报表统计。
          </p>
        </template>
      </NSpin>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="handleClose">关闭</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
