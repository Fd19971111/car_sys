<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '../../components/finance-meta';
import { useFundAccounts } from '../../components/use-fund-accounts';
import { useProjectOptions } from '../../components/use-project-options';
import { PAYMENT_TYPE_OPTIONS } from './payment-meta';
import { usePaymentData } from './use-payment-data';

defineOptions({
  name: 'PaymentOperateDrawer'
});

interface Props {
  /** 操作类型：add=新增 edit=编辑 */
  operateType: NaiveUI.TableOperateType;
  /** 编辑时的收款记录数据 */
  rowData?: Api.Payment.Item | null;
  /** 新增时默认选中的项目 */
  defaultProjectId?: number | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 540));

const { formRef, validate, restoreValidation } = useNaiveForm();

const { projectOptions, loadProjectOptions } = useProjectOptions();

const { fundAccountOptions, hasFundAccounts, defaultFundAccountId, loadFundAccounts } = useFundAccounts();

const { createPayment, updatePayment } = usePaymentData();

const title = computed(() => (props.operateType === 'add' ? '新增收款记录' : '编辑收款记录'));

const submitting = ref(false);

function createDefaultModel(): Api.Payment.FormData {
  return {
    projectId: props.defaultProjectId ?? null,
    amount: null,
    paymentType: '定金',
    paymentDate: dayjs().format('YYYY-MM-DD'),
    fundAccountId: defaultFundAccountId.value,
    remark: ''
  };
}

const model = ref<Api.Payment.FormData>(createDefaultModel());

/** 选中项目的合同金额（用于对照收款） */
const selectedProjectAmount = computed(
  () => projectOptions.value.find(item => item.value === model.value.projectId)?.amount ?? 0
);

/** 编辑历史记录时，如果原本就没有凭证，保存后仍然不会记账 */
const isLegacyRecord = computed(() => props.operateType === 'edit' && !props.rowData?.voucherNo);

/** 收款账户是自动记账的必要条件，只有配置了资金账户时才必填 */
const rules = computed<FormRules>(() => ({
  projectId: [{ required: true, type: 'number', message: '请选择所属项目', trigger: ['change', 'blur'] }],
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
  fundAccountId: [
    {
      required: hasFundAccounts.value,
      type: 'number',
      message: '请选择收款账户',
      trigger: ['change', 'blur']
    }
  ]
}));

watch(visible, async value => {
  if (!value) {
    return;
  }

  // 资金账户 / 项目下拉可能尚未加载（例如直接进入本页后立刻点新增），项目强制刷新以包含刚新建的项目
  await Promise.all([loadProjectOptions(true), loadFundAccounts()]);

  if (props.operateType === 'edit' && props.rowData) {
    const { projectId, amount, paymentType, paymentDate, fundAccountId, remark } = props.rowData;

    // 历史记录可能没有资金账户，此时回退到默认账户，交给用户确认
    model.value = {
      projectId: projectId || null,
      amount,
      paymentType: paymentType || '定金',
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
      props.operateType === 'add'
        ? await createPayment(model.value)
        : await updatePayment(props.rowData?.id ?? 0, model.value);

    if (!result) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    if (result.voucherNo) {
      window.$message?.success(
        `${props.operateType === 'add' ? '新增' : '更新'}收款记录成功，已自动生成凭证 ${result.voucherNo}`
      );
    } else {
      window.$message?.success(
        `${props.operateType === 'add' ? '新增' : '更新'}收款记录成功（财务未初始化，暂未生成凭证）`
      );
    }

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
      <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        label-width="100"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="所属项目" path="projectId">
          <NSelect
            v-model:value="model.projectId"
            :options="projectOptions"
            placeholder="请选择所属项目"
            filterable
            clearable
          />
        </NFormItem>
        <NFormItem v-if="selectedProjectAmount > 0" label=" ">
          <span class="text-12px text-gray-400">该项目合同金额：{{ formatMoney(selectedProjectAmount) }}</span>
        </NFormItem>
        <NFormItem label="收款金额" path="amount">
          <NInputNumber
            v-model:value="model.amount"
            class="w-full"
            :min="0.01"
            :precision="2"
            :show-button="false"
            placeholder="请输入收款金额"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="款项类型" path="paymentType">
          <NSelect
            v-model:value="model.paymentType"
            :options="PAYMENT_TYPE_OPTIONS"
            placeholder="请选择或输入款项类型"
            filterable
            tag
            clearable
          />
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
          <NSelect
            v-model:value="model.fundAccountId"
            :options="fundAccountOptions"
            :disabled="!hasFundAccounts"
            placeholder="请选择钱进入了哪个账户"
            clearable
          />
        </NFormItem>
        <NFormItem v-if="hasFundAccounts && isLegacyRecord" label=" ">
          <span class="text-12px text-gray-400">
            该记录是历史数据，保存后仍保持「未记账」；新登记的收款才会自动生成凭证。
          </span>
        </NFormItem>
        <NFormItem v-else-if="hasFundAccounts" label=" ">
          <span class="text-12px text-gray-400">
            保存后按「借：收款账户 / 贷：主营业务收入」自动生成收款凭证
          </span>
        </NFormItem>
        <NFormItem v-else label=" ">
          <span class="text-12px text-warning">
            还没有可用的资金账户，本次收款不会生成记账凭证。请先到「财务基础资料」配置现金 / 银行账户。
          </span>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput
            v-model:value="model.remark"
            type="textarea"
            :rows="3"
            maxlength="255"
            show-count
            placeholder="如：微信转账、银行转账尾号 1234"
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">
            {{ operateType === 'add' ? '保存收款' : '保存修改' }}
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
