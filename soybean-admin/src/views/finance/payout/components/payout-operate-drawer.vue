<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '../../components/finance-meta';
import { useFundAccounts } from '../../components/use-fund-accounts';
import { useProjectOptions } from '../../components/use-project-options';
import { fetchGetSupplierList } from '@/service/api/supplier';
import { fetchGetPendingPurchaseOrders } from '@/service/api';
import { PAYOUT_TYPE_ACCOUNT_HINT, PAYOUT_TYPE_OPTIONS } from './payout-meta';
import { usePayoutData } from './use-payout-data';

defineOptions({
  name: 'PayoutOperateDrawer'
});

interface Props {
  /** 操作类型：add=新增 edit=编辑 */
  operateType: NaiveUI.TableOperateType;
  /** 编辑时的付款记录数据 */
  rowData?: Api.Payout.Item | null;
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

const { fundAccountOptions, hasFundAccounts, defaultFundAccountId, loadFundAccounts, getFundAccountName } =
  useFundAccounts();

const { createPayout, updatePayout } = usePayoutData();

const supplierOptions = ref<{ label: string; value: number }[]>([]);
const purchaseOrderOptions = ref<{ label: string; value: number }[]>([]);
const purchaseOrders = ref<Api.ProjectData.PurchaseOrder[]>([]);

async function loadSupplierOptions() {
  const { data, error } = await fetchGetSupplierList({ current: 1, size: 100 });
  if (!error && data) {
    supplierOptions.value = (data.records ?? []).map(item => ({ label: item.name, value: item.id }));
  }
}

async function loadPurchaseOrders() {
  const { data, error } = await fetchGetPendingPurchaseOrders();
  if (!error && data) {
    purchaseOrders.value = data.records ?? [];
    purchaseOrderOptions.value = purchaseOrders.value.map(item => ({
      label: `${item.orderNo} · ${item.projectName} · ${item.materialName} · ${formatMoney(item.amount)}`,
      value: item.id
    }));
  }
}

const title = computed(() => (props.operateType === 'add' ? '新增付款记录' : '编辑付款记录'));

const submitting = ref(false);
const payoutSource = ref<'purchase' | 'manual'>('manual');

function handleSourceChange(source: 'purchase' | 'manual') {
  payoutSource.value = source;
  if (source === 'manual') {
    model.value.purchaseOrderId = null;
  }
}

function createDefaultModel(): Api.Payout.FormData {
  return {
    projectId: props.defaultProjectId ?? null,
    amount: null,
    payoutType: '材料款',
    payoutDate: dayjs().format('YYYY-MM-DD'),
    vendorId: null,
    purchaseOrderId: null,
    fundAccountId: defaultFundAccountId.value,
    remark: ''
  };
}

const model = ref<Api.Payout.FormData>(createDefaultModel());

function handlePurchaseOrderChange(id: number | null) {
  model.value.purchaseOrderId = id;
  if (!id) return;
  const order = purchaseOrders.value.find(item => item.id === id);
  if (!order) return;
  model.value.projectId = order.projectId;
  model.value.amount = order.amount;
  model.value.payoutType = '材料款';
  model.value.remark = `采购单 ${order.orderNo}：${order.materialName}${order.spec ? `（${order.spec}）` : ''}${order.remark ? `；${order.remark}` : ''}`.slice(0, 255);
  model.value.vendorId = supplierOptions.value.find(item => item.label === order.supplier)?.value ?? null;
}

/** 选中项目的合同金额（用于对照付款） */
const selectedProjectAmount = computed(
  () => projectOptions.value.find(item => item.value === model.value.projectId)?.amount ?? 0
);

/** 借方科目提示：让不懂财务的人也知道这笔钱记到哪里 */
const debitAccountHint = computed(() => PAYOUT_TYPE_ACCOUNT_HINT[model.value.payoutType] ?? '成本 / 费用科目');

/** 编辑历史记录时，如果原本就没有凭证，保存后仍然不会记账 */
const isLegacyRecord = computed(() => props.operateType === 'edit' && !props.rowData?.voucherNo);

/** 付款账户是自动记账的必要条件，只有配置了资金账户时才必填 */
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
  vendorId: [
    {
      required: true,
      type: 'number',
      message: '请选择供应商',
      trigger: ['change', 'blur']
    }
  ],
  fundAccountId: [
    {
      required: hasFundAccounts.value,
      type: 'number',
      message: '请选择付款账户',
      trigger: ['change', 'blur']
    }
  ]
}));

watch(visible, async value => {
  if (!value) {
    return;
  }

  // 资金账户 / 项目下拉可能尚未加载（例如直接进入本页后立刻点新增）
  await Promise.all([loadProjectOptions(true), loadFundAccounts(), loadSupplierOptions(), loadPurchaseOrders()]);

  if (props.operateType === 'edit' && props.rowData) {
    const { projectId, amount, payoutType, payoutDate, vendorId, purchaseOrderId, fundAccountId, remark } = props.rowData;

    model.value = {
      projectId: projectId || null,
      amount,
      payoutType: payoutType || '材料款',
      payoutDate: payoutDate || dayjs().format('YYYY-MM-DD'),
      vendorId: vendorId || null,
      purchaseOrderId: purchaseOrderId || null,
      fundAccountId: fundAccountId || defaultFundAccountId.value,
      remark
    };
    payoutSource.value = purchaseOrderId ? 'purchase' : 'manual';
  } else {
    model.value = createDefaultModel();
    payoutSource.value = 'manual';
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
        ? await createPayout(model.value)
        : await updatePayout(props.rowData?.id ?? 0, model.value);

    if (!result) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    if (result.voucherNo) {
      window.$message?.success(
        `${props.operateType === 'add' ? '新增' : '更新'}付款记录成功，已自动生成凭证 ${result.voucherNo}`
      );
    } else {
      window.$message?.success(
        `${props.operateType === 'add' ? '新增' : '更新'}付款记录成功（财务未初始化，暂未生成凭证）`
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
        <NFormItem v-if="operateType === 'add'" label="付款来源">
          <NRadioGroup :value="payoutSource" @update:value="handleSourceChange">
            <NRadioButton value="purchase">采购付款</NRadioButton>
            <NRadioButton value="manual">其他付款</NRadioButton>
          </NRadioGroup>
        </NFormItem>
        <NFormItem v-if="operateType === 'add' && payoutSource === 'purchase'" label="采购单">
          <NSelect
            :value="model.purchaseOrderId ?? null"
            :options="purchaseOrderOptions"
            filterable
            clearable
            placeholder="选择待付款采购单"
            @update:value="handlePurchaseOrderChange"
          />
        </NFormItem>
        <NFormItem label="所属项目" path="projectId">
          <NSelect
            v-model:value="model.projectId"
            :options="projectOptions"
                        :disabled="Boolean(model.purchaseOrderId)"
            placeholder="不选表示公司公共支出（如房租、办公费）"
            filterable
            clearable
          />
        </NFormItem>
        <NFormItem v-if="selectedProjectAmount > 0" label=" ">
          <span class="text-12px text-gray-400">该项目合同金额：{{ formatMoney(selectedProjectAmount) }}</span>
        </NFormItem>
        <NFormItem label="付款金额" path="amount">
          <NInputNumber
            v-model:value="model.amount"
            class="w-full"
            :min="0.01"
            :precision="2"
            :show-button="false"
                        :disabled="Boolean(model.purchaseOrderId)"
            placeholder="请输入付款金额"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="费用类型" path="payoutType">
          <NSelect
            v-model:value="model.payoutType"
            :options="PAYOUT_TYPE_OPTIONS"
            placeholder="请选择或输入费用类型"
            filterable
            tag
            clearable
          />
        </NFormItem>
        <NFormItem label="供应商" path="vendorId">
          <NSelect
            v-model:value="model.vendorId"
            :options="supplierOptions"
            placeholder="请选择供应商"
            filterable
            clearable
          />
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
          <NSelect
            v-model:value="model.fundAccountId"
            :options="fundAccountOptions"
            :disabled="!hasFundAccounts"
            placeholder="请选择钱从哪个账户付出"
            clearable
          />
        </NFormItem>
        <NFormItem v-if="hasFundAccounts && isLegacyRecord" label=" ">
          <span class="text-12px text-gray-400">
            该记录是历史数据，保存后仍保持「未记账」；新登记的付款才会自动生成凭证。
          </span>
        </NFormItem>
        <NFormItem v-else-if="hasFundAccounts" label=" ">
          <span class="text-12px text-gray-400">
            保存后按「借：{{ debitAccountHint }} / 贷：{{ getFundAccountName(model.fundAccountId) || '付款账户' }}」自动生成付款凭证
          </span>
        </NFormItem>
        <NFormItem v-else label=" ">
          <span class="text-12px text-warning">
            还没有可用的资金账户，本次付款不会生成记账凭证。请先到「财务基础资料」配置现金 / 银行账户。
          </span>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput
            v-model:value="model.remark"
            type="textarea"
            :rows="3"
            maxlength="255"
            show-count
            placeholder="如：月结转账、结算人工费尾款"
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">
            {{ operateType === 'add' ? '保存付款' : '保存修改' }}
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
