<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { fetchSaveFinFundAccount } from '@/service/api';
import { FUND_TYPE_OPTIONS } from './basic-meta';

defineOptions({
  name: 'FinanceBasicFundDrawer'
});

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.FinBasic.FundAccount | null;
  /** 可关联的资产类科目（启用中） */
  ledgerOptions: { label: string; value: string }[];
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 460));

const { formRef, validate, restoreValidation } = useNaiveForm();

const title = computed(() => (props.operateType === 'add' ? '新增资金账户' : '编辑资金账户'));

const submitting = ref(false);

function createDefaultModel(): Api.FinBasic.FundAccountFormData {
  return {
    type: 'BANK',
    name: '',
    bankName: '',
    accountLast4: '',
    ledgerAccountCode: '',
    enabled: true
  };
}

const model = ref<Api.FinBasic.FundAccountFormData>(createDefaultModel());

const isBank = computed(() => model.value.type === 'BANK');

/** 编辑时若原关联科目已停用，仍展示出来以便用户知晓 */
const mergedLedgerOptions = computed(() => {
  const options = [...props.ledgerOptions];

  if (props.rowData && !options.some(item => item.value === props.rowData!.ledgerAccountCode)) {
    options.unshift({
      label: `${props.rowData.ledgerAccountCode} ${props.rowData.ledgerAccountName}（已停用）`,
      value: props.rowData.ledgerAccountCode
    });
  }

  return options;
});

const ledgerStale = computed(
  () => Boolean(props.rowData) && !props.ledgerOptions.some(item => item.value === props.rowData?.ledgerAccountCode)
);

const rules: FormRules = {
  type: [{ required: true, message: '请选择账户类型', trigger: ['change', 'blur'] }],
  name: [
    { required: true, message: '请输入账户名称', trigger: ['input', 'blur'] },
    { max: 100, message: '账户名称不能超过100个字符', trigger: ['input', 'blur'] }
  ],
  bankName: [
    {
      validator: (_rule, value: string) => {
        if (!isBank.value) {
          return true;
        }

        return Boolean(value && value.trim()) || new Error('银行账户必须填写开户银行');
      },
      trigger: ['input', 'blur']
    }
  ],
  accountLast4: [
    {
      validator: (_rule, value: string) => {
        if (!value) {
          return true;
        }

        return /^\d{4}$/.test(value) || new Error('卡号末四位必须是4位数字');
      },
      trigger: ['input', 'blur']
    }
  ],
  ledgerAccountCode: [{ required: true, message: '请选择关联科目', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { id, type, name, bankName, accountLast4, ledgerAccountCode, enabled } = props.rowData;

    model.value = { id, type, name, bankName, accountLast4, ledgerAccountCode, enabled };
  } else {
    model.value = createDefaultModel();
  }

  await nextTick();

  restoreValidation();
});

watch(
  () => model.value.type,
  type => {
    if (type === 'CASH') {
      model.value.bankName = '';
      model.value.accountLast4 = '';
    }
  }
);

async function handleSubmit() {
  try {
    await validate();
  } catch {
    window.$message?.warning('请检查表单填写是否正确');

    return;
  }

  submitting.value = true;

  try {
    const { error } = await fetchSaveFinFundAccount({ ...model.value });

    if (error) {
      return;
    }

    window.$message?.success(props.operateType === 'add' ? '资金账户已新增' : '资金账户已更新');

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
      <NAlert v-if="ledgerStale" class="mb-16px" type="warning" :bordered="false" title="关联科目已停用">
        原关联科目当前不可用，请重新选择一个启用中的资产类科目。
      </NAlert>

      <NForm ref="formRef" :model="model" :rules="rules" label-placement="top" require-mark-placement="right-hanging">
        <NFormItem label="账户类型" path="type">
          <NRadioGroup v-model:value="model.type">
            <NRadioButton v-for="item in FUND_TYPE_OPTIONS" :key="item.value" :value="item.value">
              {{ item.label }}账户
            </NRadioButton>
          </NRadioGroup>
        </NFormItem>

        <NFormItem label="账户名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" show-count placeholder="如 工商银行基本户、备用金" />
        </NFormItem>

        <template v-if="isBank">
          <NFormItem label="开户银行" path="bankName">
            <NInput v-model:value="model.bankName" maxlength="100" placeholder="请输入开户银行名称" />
          </NFormItem>
          <NFormItem label="卡号末四位" path="accountLast4">
            <NInput v-model:value="model.accountLast4" maxlength="4" placeholder="仅保存末四位，如 1234" />
          </NFormItem>
        </template>

        <NFormItem label="关联科目" path="ledgerAccountCode">
          <NSelect
            v-model:value="model.ledgerAccountCode"
            :options="mergedLedgerOptions"
            filterable
            placeholder="请选择资产类科目，如 1002 银行存款"
          />
        </NFormItem>
        <div class="form-tip">资金账户的收支会通过该科目参与借贷记账与资金流水统计</div>

        <NFormItem label="启用状态">
          <div class="flex-y-center gap-12px">
            <NSwitch v-model:value="model.enabled" />
            <span class="text-13px text-gray-400">{{ model.enabled ? '启用后可用于记账' : '停用后不再产生新流水' }}</span>
          </div>
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">
            {{ props.operateType === 'add' ? '新增账户' : '保存修改' }}
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
.form-tip {
  padding: 10px 14px;
  margin: -6px 0 18px;
  font-size: 12px;
  line-height: 1.6;
  color: rgb(var(--base-text-color) / 0.5);
  background-color: rgb(var(--primary-color) / 0.05);
  border-radius: 8px;
}
</style>
