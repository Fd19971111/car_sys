<script setup lang="ts">
import { h } from 'vue';
import type { DataTableColumns, DataTableRowKey, PaginationProps } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { formatDriverDate, getDriverStatusMeta } from './driver-meta';

defineOptions({ name: 'DriverTable' });

interface Props {
  data: Api.Transport.Driver[];
  loading: boolean;
  checkedRowKeys: DataTableRowKey[];
  pagination: PaginationProps;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'update:checkedRowKeys', value: DataTableRowKey[]): void;
  (e: 'edit', row: Api.Transport.Driver): void;
  (e: 'delete', row: Api.Transport.Driver): void;
}>();

const columns: DataTableColumns<Api.Transport.Driver> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'name',
    title: '司机',
    width: 190,
    fixed: 'left',
    render: row =>
      h('div', { class: 'flex-y-center gap-8px' }, [
        h(
          'div',
          {
            class: 'h-28px w-28px flex-center flex-shrink-0 rounded-full text-12px font-medium',
            style: { backgroundColor: 'rgb(var(--primary-color) / 0.12)', color: 'rgb(var(--primary-color))' }
          },
          row.name.slice(0, 1)
        ),
        h('div', { class: 'flex-col' }, [
          h('span', { class: 'font-medium' }, row.name),
          h('span', { class: 'text-12px text-gray-400' }, row.licenseNumber ? `驾驶证 ${row.licenseNumber}` : `ID ${row.id}`)
        ])
      ])
  },
  { key: 'phone', title: '手机号', width: 130, render: row => h('span', { class: 'font-mono' }, row.phone) },
  {
    key: 'licenseType',
    title: '准驾车型',
    width: 100,
    render: row => row.licenseType || '-'
  },
  {
    key: 'serviceStatus',
    title: '接单状态',
    width: 110,
    render: row => {
      const meta = getDriverStatusMeta(row.serviceStatus);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'rating',
    title: '服务评分',
    width: 100,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, Number(row.rating ?? 0).toFixed(1))
  },
  {
    key: 'completedTrips',
    title: '完成订单',
    width: 100,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, row.completedTrips ?? 0)
  },
  { key: 'joinedAt', title: '加入日期', width: 120, render: row => formatDriverDate(row.joinedAt) },
  {
    key: 'operate',
    title: '操作',
    width: 130,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-4px' }, [
        h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => emit('edit', row) }, { default: () => '编辑' }),
        h(NPopconfirm, { onPositiveClick: () => emit('delete', row) }, {
          trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
          default: () => '确认删除该司机吗？车辆绑定会自动解除。'
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
    :row-key="(row: Api.Transport.Driver) => row.id"
    :scroll-x="1120"
    remote
    size="small"
    striped
    @update:checked-row-keys="handleCheckedRowKeys"
  />
</template>
