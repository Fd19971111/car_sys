<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { useFundAccounts } from '../../components/use-fund-accounts';
import { PAYOUT_TYPE_OPTIONS } from './payout-meta';
import { usePayoutData } from './use-payout-data';

defineOptions({ name: 'PayoutOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Payout.Item | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { fundAccountOptions, hasFundAccounts, defaultFundAccountId, loadFundAccounts } = useFundAccounts();
const { createPayout, updatePayout } = usePayoutData();
const title = computed(() => (props.operateType === 'add' ? '新增付款记录' : '编辑付款记录'));
const submitting = ref(false);

function createDefaultModel(): Api.Payout.FormData {
  return {
    amount: null,
    payoutType: '燃油费',
    payoutDate: dayjs().format('YYYY-MM-DD'),
    vendorName: '',
    fundAccountId: defaultFundAccountId.value,
    remark: ''
  };
}

const model = ref<Api.Payout.FormData>(createDefaultModel());

const rules = computed<FormRules>(() => ({
  amount: [
    { required: true, type: 'number', message: '请输入付款金额', trigger: ['change', 'blur'] },
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) > 0,
      message: '付款金额必须大于 0',
      trigger: ['change', 'blur']
    }
  ],
  payoutType: [{ required: true, message: '请选择费用类型', trigger: ['change', 'blur'] }],
  payoutDate: [{ required: true, message: '请选择付款日期', trigger: ['change', 'blur'] }],
  vendorName: [{ required: true, message: '请填写收款方', trigger: ['input', 'blur'] }],
  fundAccountId: [{ required: hasFundAccounts.value, type: 'number', message: '请选择付款账户', trigger: ['change', 'blur'] }]
}));

watch(visible, async value => {
  if (!value) return;
  await loadFundAccounts();
  if (props.operateType === 'edit' && props.rowData) {
    const { amount, payoutType, payoutDate, vendorName, fundAccountId, remark } = props.rowData;
    model.value = {
      amount,
      payoutType: payoutType || '燃油费',
      payoutDate: payoutDate || dayjs().format('YYYY-MM-DD'),
      vendorName: vendorName || '',
      fundAccountId: fundAccountId || defaultFundAccountId.value,
      remark
    };
  } else {
    model.value = createDefaultModel();
  }
  await nextTick();
  restoreValidation();
});

async function handleSubmit() {
  try {
    await validate();
  } catch {
    window.$message?.warning('请检查表单填写是否正确');
    return;
  }

  submitting.value = true;
  try {
    const result =
      props.operateType === 'add' ? await createPayout(model.value) : await updatePayout(props.rowData?.id ?? 0, model.value);
    if (!result) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(result.voucherNo ? `保存成功，已生成凭证 ${result.voucherNo}` : '保存成功（财务未初始化，暂未生成凭证）');
    emit('submitted');
    visible.value = false;
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth" :mask-closable="false">
    <NDrawerContent :title="title" closable :native-scrollbar="false" @close="visible = false">
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="100" require-mark-placement="right-hanging">
        <NFormItem label="付款金额" path="amount">
          <NInputNumber v-model:value="model.amount" class="w-full" :min="0.01" :precision="2" :show-button="false" placeholder="请输入付款金额">
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="费用类型" path="payoutType">
          <NSelect v-model:value="model.payoutType" :options="PAYOUT_TYPE_OPTIONS" placeholder="请选择或输入费用类型" filterable tag clearable />
        </NFormItem>
        <NFormItem label="收款方" path="vendorName">
          <NInput v-model:value="model.vendorName" maxlength="100" placeholder="加油站 / 维修厂 / 服务商" clearable />
        </NFormItem>
        <NFormItem label="付款日期" path="payoutDate">
          <NDatePicker
            v-model:formatted-value="model.payoutDate"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="请选择付款日期"
            clearable
          />
        </NFormItem>
        <NFormItem label="付款账户" path="fundAccountId">
          <NSelect v-model:value="model.fundAccountId" :options="fundAccountOptions" placeholder="请选择付款账户" filterable clearable />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count placeholder="如：车辆保养、停车费" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">保存</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>
