<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { fetchGetSupplierList } from '@/service/api/supplier';
import {
  fetchCreateStock,
  fetchGetProjectSelectOptions,
  fetchGetWarehouseList,
  fetchUpdateStock
} from '@/service/api';
import { formatMoney } from '@/views/project/project-list/components/project-meta';

defineOptions({
  name: 'StockOperateDrawer'
});

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Inventory.Stock | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));

const { formRef, validate, restoreValidation } = useNaiveForm();

const title = computed(() => (props.operateType === 'add' ? '新增库存' : '编辑库存'));

const submitting = ref(false);

/** 项目下拉选项（新增时必选） */
const projectOptions = ref<Api.Project.SelectOption[]>([]);

/** 仓库下拉选项（来源：仓库管理） */
const warehouseOptions = ref<{ label: string; value: number }[]>([]);

function createDefaultModel(): Api.Inventory.StockForm {
  return {
    projectId: null,
    warehouseId: null,
    name: '',
    spec: '',
    unit: '',
    quantity: null,
    unitPrice: null,
    supplier: '',
    remark: ''
  };
}

const model = ref<Api.Inventory.StockForm>(createDefaultModel());

/** 金额预览 = 数量 × 单价 */
const amountPreview = computed(() => Number(model.value.quantity ?? 0) * Number(model.value.unitPrice ?? 0));

/** 供应商下拉选项：库里存的是供应商名称，故 value 直接用名称 */
const supplierOptions = ref<CommonType.Option<string, string>[]>([]);

/** 供应商下拉的绑定值：空字符串统一转成 null，避免下拉框显示出空白项 */
const supplierValue = computed<string | null>({
  get: () => model.value.supplier || null,
  set: value => {
    model.value.supplier = value ?? '';
  }
});

/** 加载供应商选项（来源：供应商管理） */
async function loadSupplierOptions() {
  const { data, error } = await fetchGetSupplierList({ current: 1, size: 200 });

  if (!error && data) {
    supplierOptions.value = (data.records ?? []).map(item => ({ label: item.name, value: item.name }));
  }
}

/** 加载项目选项（来源：项目管理） */
async function loadProjectOptions() {
  const { data, error } = await fetchGetProjectSelectOptions();

  if (!error && data) {
    projectOptions.value = data.projects ?? [];
  }
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
  ]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { projectId, warehouseId, name, spec, unit, quantity, unitPrice, supplier, remark } = props.rowData;

    model.value = {
      projectId,
      warehouseId,
      name,
      spec,
      unit,
      quantity,
      unitPrice,
      supplier,
      remark
    };
  } else {
    model.value = createDefaultModel();
  }

  await nextTick();
  restoreValidation();

  await Promise.all([loadSupplierOptions(), loadProjectOptions(), loadWarehouseOptions()]);
});

async function handleSubmit() {
  try {
    await validate();
  } catch {
    window.$message?.warning('请检查表单填写是否正确');

    return;
  }

  if (!model.value.projectId) {
    window.$message?.warning('请先选择所属项目');

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
        ? await fetchCreateStock(model.value).then(res => !res.error)
        : await fetchUpdateStock({ ...model.value, id: props.rowData?.id ?? 0 }).then(res => !res.error);

    if (!success) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    window.$message?.success(props.operateType === 'add' ? '新增库存成功' : '库存已更新');
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
        <NFormItem label="所属项目" path="projectId">
          <NSelect
            v-model:value="model.projectId"
            :options="projectOptions"
            :disabled="operateType === 'edit'"
            placeholder="请选择所属项目"
            filterable
          />
        </NFormItem>
        <NFormItem label="所属仓库" path="warehouseId">
          <NSelect
            v-model:value="model.warehouseId"
            :options="warehouseOptions"
            placeholder="请选择所属仓库"
            filterable
          />
        </NFormItem>
        <NFormItem label="材料名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" placeholder="如：瓷砖、乳胶漆" clearable />
        </NFormItem>
        <NFormItem label="规格/型号" path="spec">
          <NInput v-model:value="model.spec" maxlength="100" placeholder="如：800×800、18L" clearable />
        </NFormItem>
        <NFormItem label="单位" path="unit">
          <NSelect v-model:value="model.unit" :options="unitOptions" placeholder="请选择或输入单位" filterable tag clearable />
        </NFormItem>
        <NFormItem label="库存数量" path="quantity">
          <NInputNumber
            v-model:value="model.quantity"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            placeholder="当前库存数量"
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
        <NFormItem label="供应商" path="supplier">
          <NSelect
            v-model:value="supplierValue"
            :options="supplierOptions"
            placeholder="请选择或输入供应商"
            filterable
            tag
            clearable
          />
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
