<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { usePricingRuleData } from './use-pricing-rule-data';

defineOptions({ name: 'PricingRuleOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.PricingRule.Item | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{ (e: 'submitted'): void }>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 480));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createRule, updateRule } = usePricingRuleData();
const title = computed(() => (props.operateType === 'add' ? '新增计价规则' : '编辑计价规则'));
const submitting = ref(false);

const carTypeOptions = ['经济型', '舒适型', '商务型', '豪华型', '六座商务', '七座商务'].map(item => ({
  label: item,
  value: item
}));

interface FormModel {
  name: string;
  carType: string | null;
  basePrice: number | null;
  perKmPrice: number | null;
  enabled: boolean;
}

function createDefaultModel(): FormModel {
  return { name: '', carType: null, basePrice: null, perKmPrice: null, enabled: true };
}

const model = ref<FormModel>(createDefaultModel());

const rules: FormRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: ['input', 'blur'] }],
  carType: [{ required: true, message: '请选择或输入适用车型', trigger: ['change', 'blur'] }],
  basePrice: [{ required: true, type: 'number', message: '请输入起步价', trigger: ['input', 'blur'] }],
  perKmPrice: [{ required: true, type: 'number', message: '请输入超里程单价', trigger: ['input', 'blur'] }]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const { name, carType, basePrice, perKmPrice, enabled } = props.rowData;
    model.value = { name, carType, basePrice, perKmPrice, enabled };
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
  const payload: Api.PricingRule.FormData = {
    name: model.value.name.trim(),
    carType: model.value.carType ?? '',
    basePrice: Number(model.value.basePrice ?? 0),
    perKmPrice: Number(model.value.perKmPrice ?? 0),
    enabled: model.value.enabled
  };
  submitting.value = true;
  try {
    const success =
      props.operateType === 'add' ? await createRule(payload) : await updateRule(props.rowData?.id ?? 0, payload);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增计价规则成功' : '计价规则已更新');
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
      <div class="mb-16px rounded-8px bg-#f7f8fa px-12px py-10px text-12px leading-18px text-gray-500">
        起步价为乘客上车后的基础费用；超里程单价为超出起步里程后每公里加收的费用。启用后该规则用于订单计费。
      </div>
      <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        label-align="right"
        label-width="auto"
      >
        <NFormItem label="规则名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" placeholder="如：经济型标准计价" clearable />
        </NFormItem>
        <NFormItem label="适用车型" path="carType">
          <NSelect
            v-model:value="model.carType"
            :options="carTypeOptions"
            placeholder="请选择车型，也可输入自定义车型"
            filterable
            tag
            clearable
          />
        </NFormItem>
        <NFormItem label="起步价（元）" path="basePrice">
          <NInputNumber
            v-model:value="model.basePrice"
            class="w-full"
            :min="0"
            :precision="2"
            :step="1"
            placeholder="请输入起步价"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="超里程单价（元/公里）" path="perKmPrice">
          <NInputNumber
            v-model:value="model.perKmPrice"
            class="w-full"
            :min="0"
            :precision="2"
            :step="0.1"
            placeholder="请输入超里程单价"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="启用规则" path="enabled">
          <NSwitch v-model:value="model.enabled" />
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

<style scoped>
/* 表单标签统一单行展示，避免「超里程单价（元/公里）」这类长标签换行 */
:deep(.n-form-item-label) {
  white-space: nowrap;
}
</style>
