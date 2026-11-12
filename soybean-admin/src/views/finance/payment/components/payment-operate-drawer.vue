<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { useFundAccounts } from '../../components/use-fund-accounts';
import { PAYMENT_TYPE_OPTIONS } from './payment-meta';
import { usePaymentData } from './use-payment-data';

defineOptions({ name: 'PaymentOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Payment.Item | null;
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
const { createPayment, updatePayment } = usePaymentData();
const title = computed(() => (props.operateType === 'add' ? '新增收款记录' : '编辑收款记录'));
const submitting = ref(false);

function createDefaultModel(): Api.Payment.FormData {
  return {
    rideOrderId: null,
    amount: null,
    paymentType: '行程车费',
    paymentDate: dayjs().format('YYYY-MM-DD'),
    fundAccountId: defaultFundAccountId.value,
    remark: ''
  };
}

const model = ref<Api.Payment.FormData>(createDefaultModel());

const rules = computed<FormRules>(() => ({
  amount: [
    { required: true, type: 'number', message: '请输入收款金额', trigger: ['change', 'blur'] },
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) > 0,
      message: '收款金额必须大于 0',
      trigger: ['change', 'blur']
    }
  ],
  paymentType: [{ required: true, message: '请选择款项类型', trigger: ['change', 'blur'] }],
  paymentDate: [{ required: true, message: '请选择收款日期', trigger: ['change', 'blur'] }],
  fundAccountId: [{ required: hasFundAccounts.value, type: 'number', message: '请选择收款账户', trigger: ['change', 'blur'] }]
}));

watch(visible, async value => {
  if (!value) return;
  await loadFundAccounts();
  if (props.operateType === 'edit' && props.rowData) {
    const { rideOrderId, amount, paymentType, paymentDate, fundAccountId, remark } = props.rowData;
    model.value = {
      rideOrderId: rideOrderId || null,
      amount,
      paymentType: paymentType || '行程车费',
      paymentDate: paymentDate || dayjs().format('YYYY-MM-DD'),
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
      props.operateType === 'add' ? await createPayment(model.value) : await updatePayment(props.rowData?.id ?? 0, model.value);
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
        <NFormItem label="行程订单ID" path="rideOrderId">
          <NInputNumber v-model:value="model.rideOrderId" class="w-full" :min="1" :precision="0" clearable placeholder="可留空；填写订单详情中的ID以关联收款" />
        </NFormItem>
        <NFormItem label="收款金额" path="amount">
          <NInputNumber v-model:value="model.amount" class="w-full" :min="0.01" :precision="2" :show-button="false" placeholder="请输入收款金额">
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="款项类型" path="paymentType">
          <NSelect v-model:value="model.paymentType" :options="PAYMENT_TYPE_OPTIONS" placeholder="请选择或输入款项类型" filterable tag clearable />
        </NFormItem>
        <NFormItem label="收款日期" path="paymentDate">
          <NDatePicker
            v-model:formatted-value="model.paymentDate"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="请选择收款日期"
            clearable
          />
        </NFormItem>
        <NFormItem label="收款账户" path="fundAccountId">
          <NSelect v-model:value="model.fundAccountId" :options="fundAccountOptions" placeholder="请选择收款账户" filterable clearable />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count placeholder="如：微信转账" />
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
