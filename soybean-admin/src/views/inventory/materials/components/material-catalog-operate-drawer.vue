<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { fetchCreateMaterialCatalog, fetchUpdateMaterialCatalog } from '@/service/api';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Inventory.MaterialCatalog | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{ (e: 'submitted'): void }>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));
const { formRef, validate, restoreValidation } = useNaiveForm();
const submitting = ref(false);
const title = computed(() => (props.operateType === 'add' ? '新增材料档案' : '编辑材料档案'));

function createDefaultModel(): Api.Inventory.MaterialCatalogForm {
  return { name: '', spec: '', unit: '', unitPrice: null, taxMode: 'NO_INVOICE', taxRate: 0, remark: '' };
}

const model = ref<Api.Inventory.MaterialCatalogForm>(createDefaultModel());

watch(
  () => model.value.taxMode,
  mode => {
    if (mode === 'NO_INVOICE') model.value.taxRate = 0;
  }
);
const rules: FormRules = {
  name: [{ required: true, message: '请填写材料名称', trigger: ['input', 'blur'] }],
  unitPrice: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 0,
      message: '默认单价不能为负数',
      trigger: ['change', 'blur']
    }
  ],
  taxRate: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 0 && Number(value ?? 0) <= 100,
      message: '税率必须在 0% 到 100% 之间',
      trigger: ['change', 'blur']
    }
  ]
};

watch(visible, async value => {
  if (!value) return;

  if (props.operateType === 'edit' && props.rowData) {
    const { name, spec, unit, unitPrice, taxMode, taxRate, remark } = props.rowData;
    model.value = { name, spec, unit, unitPrice, taxMode: taxMode || 'NO_INVOICE', taxRate: taxRate ?? 0, remark };
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
    const { error } =
      props.operateType === 'add'
        ? await fetchCreateMaterialCatalog(model.value)
        : await fetchUpdateMaterialCatalog({ ...model.value, id: props.rowData?.id ?? 0 });

    if (error) return;

    window.$message?.success(props.operateType === 'add' ? '材料档案已新增' : '材料档案已更新');
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
        label-width="96"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="材料名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" placeholder="如：瓷砖、乳胶漆" clearable />
        </NFormItem>
        <NFormItem label="规格/型号" path="spec">
          <NInput v-model:value="model.spec" maxlength="100" placeholder="如：800×800、18L" clearable />
        </NFormItem>
        <NFormItem label="单位" path="unit">
          <NInput v-model:value="model.unit" maxlength="20" placeholder="如：片、桶、平米" clearable />
        </NFormItem>
        <NFormItem label="默认单价" path="unitPrice">
          <NInputNumber
            v-model:value="model.unitPrice"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            placeholder="请输入默认单价"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="税务方式" path="taxMode">
          <NSelect
            v-model:value="model.taxMode"
            :options="[
              { label: '无票/不计税', value: 'NO_INVOICE' },
              { label: '含税价', value: 'TAX_INCLUDED' },
              { label: '未税价', value: 'TAX_EXCLUDED' }
            ]"
          />
        </NFormItem>
        <NFormItem label="默认税率" path="taxRate">
          <NInputNumber
            v-model:value="model.taxRate"
            class="w-full"
            :min="0"
            :max="100"
            :precision="2"
            :disabled="model.taxMode === 'NO_INVOICE'"
            :show-button="false"
          >
            <template #suffix>%</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">提交</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>
