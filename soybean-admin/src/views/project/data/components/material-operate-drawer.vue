<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { fetchGetMaterialCatalogList, fetchGetWarehouseList } from '@/service/api';
import { formatMoney } from '../../project-list/components/project-meta';
import { useProjectData } from './use-project-data';

defineOptions({
  name: 'MaterialOperateDrawer'
});

interface Props {
  projectId: number | null;
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.ProjectData.Material | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));

const { formRef, validate, restoreValidation } = useNaiveForm();

const { createMaterial, updateMaterial } = useProjectData();

const title = computed(() => (props.operateType === 'add' ? '新增材料' : '编辑材料'));

const submitting = ref(false);

function createDefaultModel(): Api.ProjectData.MaterialForm {
  return {
    projectId: props.projectId,
    catalogId: null,
    warehouseId: null,
    name: '',
    spec: '',
    unit: '',
    quantity: null,
    unitPrice: null,
    taxMode: 'NO_INVOICE',
    taxRate: 0,
    supplier: '',
    remark: ''
  };
}

const model = ref<Api.ProjectData.MaterialForm>(createDefaultModel());

watch(
  () => model.value.taxMode,
  mode => {
    if (mode === 'NO_INVOICE') model.value.taxRate = 0;
  }
);

/** 金额预览 = 数量 × 单价 */
const amountPreview = computed(() => Number(model.value.quantity ?? 0) * Number(model.value.unitPrice ?? 0));

/** 仓库下拉选项（来源：库存查询-仓库管理） */
const warehouseOptions = ref<{ label: string; value: number }[]>([]);
const materialCatalogOptions = ref<Api.Inventory.MaterialCatalog[]>([]);

const materialCatalogSelectOptions = computed(() =>
  materialCatalogOptions.value.map(item => ({
    label: `${item.name}${item.spec ? ` · ${item.spec}` : ''}`,
    value: item.id
  }))
);

async function loadMaterialCatalogOptions() {
  const { data, error } = await fetchGetMaterialCatalogList({ current: 1, size: 200 });
  if (!error && data) materialCatalogOptions.value = data.records ?? [];
}

function handleMaterialCatalogChange(catalogId: number | null) {
  model.value.catalogId = catalogId;
  const catalog = materialCatalogOptions.value.find(item => item.id === catalogId);
  if (!catalog) return;

  model.value.name = catalog.name;
  model.value.spec = catalog.spec;
  model.value.unit = catalog.unit;
  model.value.unitPrice = catalog.unitPrice;
  model.value.taxMode = catalog.taxMode || 'NO_INVOICE';
  model.value.taxRate = catalog.taxRate ?? 0;
}

/** 加载仓库选项（来源：库存查询-仓库管理） */
async function loadWarehouseOptions() {
  const { data, error } = await fetchGetWarehouseList({ current: 1, size: 100 });

  if (!error && data) {
    warehouseOptions.value = (data.records ?? []).map(item => ({ label: item.name, value: item.id }));
  }
}

const unitOptions = [
  { label: '项', value: '项' },
  { label: '个', value: '个' },
  { label: '套', value: '套' },
  { label: '米', value: '米' },
  { label: '平米', value: '平米' },
  { label: '桶', value: '桶' },
  { label: '张', value: '张' },
  { label: '袋', value: '袋' }
];

const rules: FormRules = {
  catalogId: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) > 0,
      message: '请选择全局材料档案',
      trigger: ['change', 'blur']
    }
  ],
  name: [{ required: true, message: '请填写材料名称', trigger: ['input', 'blur'] }],
  warehouseId: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) > 0,
      message: '请选择仓库',
      trigger: ['change', 'blur']
    }
  ],
  quantity: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 0,
      message: '数量不能为负数',
      trigger: ['change', 'blur']
    }
  ],
  unitPrice: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 0,
      message: '单价不能为负数',
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
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { catalogId, warehouseId, name, spec, unit, quantity, unitPrice, taxMode, taxRate, remark } = props.rowData;

    model.value = {
      projectId: props.rowData.projectId,
      catalogId: catalogId || null,
      warehouseId,
      name,
      spec,
      unit,
      quantity,
      unitPrice,
      taxMode: taxMode || 'NO_INVOICE',
      taxRate: taxRate ?? 0,
      supplier: '',
      remark
    };
  } else {
    model.value = createDefaultModel();
  }

  await nextTick();
  restoreValidation();

  await Promise.all([loadWarehouseOptions(), loadMaterialCatalogOptions()]);
});

async function handleSubmit() {
  try {
    await validate();
  } catch {
    window.$message?.warning('请检查表单填写是否正确');

    return;
  }

  if (!model.value.projectId) {
    window.$message?.warning('请先选择项目');

    return;
  }

  if (!model.value.warehouseId) {
    window.$message?.warning('请选择所属仓库');

    return;
  }

  submitting.value = true;

  try {
    const success =
      props.operateType === 'add'
        ? await createMaterial(model.value)
        : await updateMaterial(props.rowData?.id ?? 0, model.value);

    if (!success) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    window.$message?.success(props.operateType === 'add' ? '新增材料成功' : '材料已更新');
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
        <NFormItem label="所属仓库" path="warehouseId">
          <NSelect
            v-model:value="model.warehouseId"
            :options="warehouseOptions"
            placeholder="请选择所属仓库"
            filterable
          />
        </NFormItem>
        <NFormItem label="材料档案" path="catalogId">
          <NSelect
            :value="model.catalogId"
            :options="materialCatalogSelectOptions"
            placeholder="请选择公司级材料档案"
            filterable
            clearable
            @update:value="handleMaterialCatalogChange"
          />
        </NFormItem>
        <NFormItem label="材料名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" placeholder="如：瓷砖、乳胶漆" :disabled="!!model.catalogId" clearable />
        </NFormItem>
        <NFormItem label="规格/型号" path="spec">
          <NInput v-model:value="model.spec" maxlength="100" placeholder="如：800×800、18L" :disabled="!!model.catalogId" clearable />
        </NFormItem>
        <NFormItem label="单位" path="unit">
          <NSelect v-model:value="model.unit" :options="unitOptions" placeholder="请选择或输入单位" :disabled="!!model.catalogId" filterable tag clearable />
        </NFormItem>
        <NFormItem label="清单数量" path="quantity">
          <NInputNumber
            v-model:value="model.quantity"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            :disabled="!!model.catalogId"
            placeholder="计划使用的数量"
          />
        </NFormItem>
        <NFormItem label="单价" path="unitPrice">
          <NInputNumber
            v-model:value="model.unitPrice"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            placeholder="请输入单价"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="金额">
          <span class="text-14px font-medium">{{ formatMoney(amountPreview) }}</span>
          <span class="ml-8px text-12px text-gray-400">数量 × 单价，自动计算</span>
        </NFormItem>
        <NFormItem label="税务方式" path="taxMode">
          <NSelect
            v-model:value="model.taxMode"
            :disabled="!!model.catalogId"
            :options="[
              { label: '无票/不计税', value: 'NO_INVOICE' },
              { label: '含税价', value: 'TAX_INCLUDED' },
              { label: '未税价', value: 'TAX_EXCLUDED' }
            ]"
          />
        </NFormItem>
        <NFormItem label="税率" path="taxRate">
          <NInputNumber
            v-model:value="model.taxRate"
            class="w-full"
            :min="0"
            :max="100"
            :precision="2"
            :disabled="!!model.catalogId || model.taxMode === 'NO_INVOICE'"
            :show-button="false"
          >
            <template #suffix>%</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="2" maxlength="255" show-count />
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

<style scoped></style>
