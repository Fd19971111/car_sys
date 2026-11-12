<script setup lang="ts">
import { h } from 'vue';
import type { DataTableColumns, DataTableRowKey, PaginationProps } from 'naive-ui';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { getVehicleCategoryMeta, getVehicleStatusMeta } from './vehicle-meta';

defineOptions({ name: 'VehicleTable' });

interface Props {
  data: Api.Transport.Vehicle[];
  loading: boolean;
  checkedRowKeys: DataTableRowKey[];
  pagination: PaginationProps;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'update:checkedRowKeys', value: DataTableRowKey[]): void;
  (e: 'edit', row: Api.Transport.Vehicle): void;
  (e: 'delete', row: Api.Transport.Vehicle): void;
}>();

const columns: DataTableColumns<Api.Transport.Vehicle> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'plateNumber',
    title: '车辆',
    width: 190,
    fixed: 'left',
    render: row =>
      h('div', { class: 'flex-col' }, [
        h('span', { class: 'font-medium font-mono' }, row.plateNumber),
        h('span', { class: 'text-12px text-gray-400' }, row.brand || `ID ${row.id}`)
      ])
  },
  {
    key: 'category',
    title: '服务车型',
    width: 120,
    render: row => {
      const meta = getVehicleCategoryMeta(row.category);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'seatCount',
    title: '可乘人数',
    width: 100,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, `${row.seatCount ?? 0} 人`)
  },
  {
    key: 'status',
    title: '车辆状态',
    width: 110,
    render: row => {
      const meta = getVehicleStatusMeta(row.status);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'driverName',
    title: '当前司机',
    width: 140,
    render: row => (row.driverName ? h('span', {}, row.driverName) : h('span', { class: 'text-gray-400' }, '未绑定'))
  },
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
          default: () => '确认删除该车辆吗？'
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
    :row-key="(row: Api.Transport.Vehicle) => row.id"
    :scroll-x="1000"
    remote
    size="small"
    striped
    @update:checked-row-keys="handleCheckedRowKeys"
  />
</template>
