<script setup lang="ts">
import { h } from 'vue';
import type { DataTableColumns, DataTableRowKey, PaginationProps } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import {
  formatOrderDateTime,
  formatOrderMoney,
  getCustomerResponseMeta,
  getOrderStatusMeta,
  getVehicleCategoryMeta
} from '../../components/order-meta';

defineOptions({ name: 'OrderTable' });

interface Props {
  data: Api.Order.Item[];
  loading: boolean;
  checkedRowKeys: DataTableRowKey[];
  pagination: PaginationProps;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'update:checkedRowKeys', value: DataTableRowKey[]): void;
  (e: 'detail', row: Api.Order.Item): void;
  (e: 'edit', row: Api.Order.Item): void;
  (e: 'delete', row: Api.Order.Item): void;
}>();

const columns: DataTableColumns<Api.Order.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'orderNumber',
    title: '订单',
    width: 200,
    fixed: 'left',
    render: row =>
      h('div', { class: 'flex-col gap-1px' }, [
        h('span', { class: 'font-medium font-mono' }, row.orderNumber),
        h('span', { class: 'text-12px text-gray-400' }, `${row.customerName || '未关联客户'} · ${row.customerPhone || '-'}`)
      ])
  },
  {
    key: 'route',
    title: '行程',
    minWidth: 280,
    render: row =>
      h('div', { class: 'flex-col gap-2px' }, [
        h('div', { class: 'flex items-start gap-6px' }, [
          h('span', { class: 'mt-5px h-6px w-6px flex-shrink-0 rounded-full bg-success' }),
          h('span', { class: 'text-13px' }, row.pickupAddress || '-')
        ]),
        h('div', { class: 'flex items-start gap-6px' }, [
          h('span', { class: 'mt-5px h-6px w-6px flex-shrink-0 rounded-full bg-error' }),
          h('span', { class: 'text-13px text-gray-500' }, row.destinationAddress || '-')
        ])
      ])
  },
  {
    key: 'scheduledAt',
    title: '预约时间',
    width: 160,
    render: row => formatOrderDateTime(row.scheduledAt)
  },
  {
    key: 'status',
    title: '订单状态',
    width: 120,
    render: row => {
      const meta = getOrderStatusMeta(row.status);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'customerResponse',
    title: '客户反馈',
    width: 110,
    render: row => {
      const meta = getCustomerResponseMeta(row.customerResponse);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'vehicleCategory',
    title: '车型',
    width: 100,
    render: row => {
      const meta = getVehicleCategoryMeta(row.vehicleCategory);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'fare',
    title: '费用',
    width: 120,
    align: 'right',
    render: row =>
      h('div', { class: 'flex-col' }, [
        h('span', { class: 'font-mono font-medium' }, formatOrderMoney(row.finalFare || row.quotedFare)),
        h('span', { class: 'text-11px text-gray-400' }, `报价 ${formatOrderMoney(row.quotedFare)}`)
      ])
  },
  {
    key: 'driver',
    title: '司机 / 车辆',
    width: 150,
    render: row =>
      h('div', { class: 'flex-col' }, [
        h('span', {}, row.driverName || '未指派司机'),
        h('span', { class: 'text-12px text-gray-400 font-mono' }, row.vehiclePlate || '未指派车辆')
      ])
  },
  {
    key: 'operate',
    title: '操作',
    width: 170,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-2px' }, [
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => emit('detail', row) }, { default: () => '详情' }),
        h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => emit('edit', row) }, { default: () => '编辑' }),
        h(NPopconfirm, { onPositiveClick: () => emit('delete', row) }, {
          trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
          default: () => '确认删除该订单吗？'
        })
      ])
  }
];

function handleCheckedRowKeys(value: DataTableRowKey[]) {
  emit('update:checkedRowKeys', value);
}
</script>

<template>
  <NDataTable
    :checked-row-keys="checkedRowKeys"
    :columns="columns"
    :data="data"
    :loading="loading"
    :pagination="pagination"
    :row-key="(row: Api.Order.Item) => row.id"
    :scroll-x="1450"
    remote
    size="small"
    striped
    @update:checked-row-keys="handleCheckedRowKeys"
  />
</template>
