<script setup lang="ts">
import { h } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton } from 'naive-ui';
import { formatMoney } from '../../components/finance-meta';

defineProps<{ records: Api.DriverFinance.DriverBalance[]; loading: boolean }>();
const emit = defineEmits<{ settle: [driver: Api.DriverFinance.DriverBalance] }>();

const columns: DataTableColumns<Api.DriverFinance.DriverBalance> = [
  { key: 'driverName', title: '司机', minWidth: 140 },
  { key: 'phone', title: '手机号', width: 140 },
  { key: 'tripCount', title: '完单数', width: 90, align: 'right' },
  { key: 'earnedTotal', title: '累计应得', width: 140, align: 'right', render: row => formatMoney(row.earnedTotal) },
  { key: 'settledTotal', title: '累计已结', width: 140, align: 'right', render: row => formatMoney(row.settledTotal) },
  { key: 'payable', title: '待结余额', width: 140, align: 'right', render: row => h('strong', { class: row.payable > 0 ? 'text-error' : '' }, formatMoney(row.payable)) },
  { key: 'action', title: '操作', width: 100, align: 'center', render: row => h(NButton, { size: 'small', type: 'primary', secondary: true, disabled: row.payable <= 0, onClick: () => emit('settle', row) }, { default: () => '结算' }) }
];
</script>

<template>
  <NDataTable :columns="columns" :data="records" :loading="loading" :row-key="row => row.driverId" :scroll-x="900" size="small" striped />
</template>