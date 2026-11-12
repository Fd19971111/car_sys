<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { formatMoney } from '../components/finance-meta';
import DriverBalanceTable from './components/driver-balance-table.vue';
import SettlementDrawer from './components/settlement-drawer.vue';
import SettlementTable from './components/settlement-table.vue';
import { useDriverFinanceData } from './components/use-driver-finance-data';

defineOptions({ name: 'FinanceDriverSettlement' });
const { loading, drivers, settlements, load, settle } = useDriverFinanceData();
const selectedDriver = ref<Api.DriverFinance.DriverBalance | null>(null);
const drawerVisible = ref(false);
const submitting = ref(false);
const totalPayable = computed(() => drivers.value.reduce((sum, driver) => sum + driver.payable, 0));

function openSettlement(driver: Api.DriverFinance.DriverBalance) {
  selectedDriver.value = driver;
  drawerVisible.value = true;
}

async function submitSettlement(form: Api.DriverFinance.SettlementForm) {
  submitting.value = true;
  try {
    const result = await settle(form);
    if (!result) {
      window.$message?.error('结算失败，请检查待结余额和财务账户');
      return;
    }
    window.$message?.success(`结算成功 ${result.settlementNo}，凭证 ${result.voucherNo}`);
    drawerVisible.value = false;
  } finally {
    submitting.value = false;
  }
}

onMounted(load);
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <div class="flex flex-wrap items-end justify-between gap-12px">
        <div><h2 class="text-18px font-semibold">司机结算</h2><p class="mt-4px text-13px text-gray-500">按已完单确认的司机应得金额结算，支持分次付款。</p></div>
        <div class="text-right"><span class="text-12px text-gray-500">司机待结合计</span><strong class="ml-8px text-20px text-error">{{ formatMoney(totalPayable) }}</strong></div>
      </div>
      <NCard title="司机待结余额" :bordered="false" size="small" class="card-wrapper"><DriverBalanceTable :records="drivers" :loading="loading" @settle="openSettlement" /></NCard>
      <NCard title="最近结算记录" :bordered="false" size="small" class="card-wrapper"><SettlementTable :records="settlements" :loading="loading" /></NCard>
    </NSpace>
    <SettlementDrawer v-model:visible="drawerVisible" :driver="selectedDriver" :submitting="submitting" @submit="submitSettlement" />
  </div>
</template>