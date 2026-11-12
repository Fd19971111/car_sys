<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { useAppStore } from '@/store/modules/app';
import { useFundAccounts } from '../../components/use-fund-accounts';
import { formatMoney } from '../../components/finance-meta';

const props = defineProps<{ driver: Api.DriverFinance.DriverBalance | null; submitting: boolean }>();
const emit = defineEmits<{ submit: [form: Api.DriverFinance.SettlementForm] }>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 480));
const { fundAccountOptions, defaultFundAccountId, loadFundAccounts } = useFundAccounts();
const model = ref<Api.DriverFinance.SettlementForm>({ driverId: 0, amount: 0, settlementDate: dayjs().format('YYYY-MM-DD'), fundAccountId: 0, remark: '' });

watch(visible, async value => {
  if (!value || !props.driver) return;
  await loadFundAccounts();
  model.value = { driverId: props.driver.driverId, amount: props.driver.payable, settlementDate: dayjs().format('YYYY-MM-DD'), fundAccountId: defaultFundAccountId.value ?? 0, remark: '' };
  await nextTick();
});

async function submit() {
  if (!props.driver || model.value.amount <= 0 || !model.value.fundAccountId) return;
  emit('submit', { ...model.value });
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth" :mask-closable="false">
    <NDrawerContent title="结算司机车费" closable @close="visible = false">
      <NForm v-if="driver" label-placement="left" label-width="100">
        <NFormItem label="司机"><span>{{ driver.driverName }} · 待结 {{ formatMoney(driver.payable) }}</span></NFormItem>
        <NFormItem label="结算金额"><NInputNumber v-model:value="model.amount" class="w-full" :min="0.01" :max="driver.payable" :precision="2"><template #prefix>¥</template></NInputNumber></NFormItem>
        <NFormItem label="付款日期"><NDatePicker v-model:formatted-value="model.settlementDate" class="w-full" type="date" value-format="yyyy-MM-dd" /></NFormItem>
        <NFormItem label="付款账户"><NSelect v-model:value="model.fundAccountId" :options="fundAccountOptions" filterable /></NFormItem>
        <NFormItem label="备注"><NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count /></NFormItem>
      </NForm>
      <template #footer><NSpace justify="end"><NButton @click="visible = false">取消</NButton><NButton type="primary" :loading="submitting" @click="submit">确认结算</NButton></NSpace></template>
    </NDrawerContent>
  </NDrawer>
</template>

