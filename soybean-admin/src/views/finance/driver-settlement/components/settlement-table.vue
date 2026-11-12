<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui';
import { formatDate, formatMoney } from '../../components/finance-meta';

defineProps<{ records: Api.DriverFinance.Settlement[]; loading: boolean }>();

const columns: DataTableColumns<Api.DriverFinance.Settlement> = [
  { key: 'settlementDate', title: '结算日期', width: 130, render: row => formatDate(row.settlementDate) },
  { key: 'settlementNo', title: '结算单号', width: 180 },
  { key: 'driverName', title: '司机', width: 120 },
  { key: 'amount', title: '结算金额', width: 140, align: 'right', render: row => formatMoney(row.amount) },
  { key: 'fundAccountName', title: '付款账户', width: 150 },
  { key: 'voucherNo', title: '记账凭证', width: 160, render: row => row.voucherNo || '未记账' },
  { key: 'remark', title: '备注', minWidth: 160, ellipsis: { tooltip: true } }
];
</script>

<template>
  <NDataTable :columns="columns" :data="records" :loading="loading" :row-key="row => row.id" :scroll-x="900" :pagination="{ pageSize: 10 }" size="small" striped />
</template>