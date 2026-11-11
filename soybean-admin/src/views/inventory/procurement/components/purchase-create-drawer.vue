<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { fetchCreatePurchaseOrder, fetchGetMaterialList, fetchGetSupplierList } from '@/service/api';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { calculatePurchaseTax } from '@/utils/purchase-tax';

defineOptions({
  name: 'PurchaseCreateDrawer'
});

interface Props {
  /** 页面当前筛选的项目，仅作为新增时的默认项目 */
  defaultProjectId: number | null;
  /** 项目下拉选项 */
  projectOptions: Api.Project.SelectOption[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted', projectId: number): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));

const { formRef, validate, restoreValidation } = useNaiveForm();

const submitting = ref(false);
const materialLoading = ref(false);
const materials = ref<Api.ProjectData.Material[]>([]);
const vendorOptions = ref<CommonType.Option<number, string>[]>([]);

const projectSelectOptions = computed<CommonType.Option<number, string>[]>(() =>
  props.projectOptions.map(item => ({ label: item.label, value: item.value }))
);

function createDefaultModel(): Api.ProjectData.PurchaseOrderForm {
  return {
    projectId: props.defaultProjectId ?? 0,
    materialId: 0,
    vendorId: null,
    purchaseDate: dayjs().format('YYYY-MM-DD'),
    quantity: 1,
    unitPrice: 0,
    taxMode: 'NO_INVOICE',
    taxRate: 0,
    remark: ''
  };
}

const model = ref<Api.ProjectData.PurchaseOrderForm>(createDefaultModel());

watch(
  () => model.value.taxMode,
  mode => {
    if (mode === 'NO_INVOICE') model.value.taxRate = 0;
  }
);

const materialOptions = computed(() =>
  materials.value.map(item => ({
    label: `${item.name}${item.spec ? ` · ${item.spec}` : ''}`,
    value: item.id
  }))
);

const selectedMaterial = computed(() => materials.value.find(item => item.id === model.value.materialId) ?? null);

const taxAmounts = computed(() =>
  calculatePurchaseTax(
    Number(model.value.quantity ?? 0),
    Number(model.value.unitPrice ?? 0),
    model.value.taxMode,
    Number(model.value.taxRate ?? 0)
  )
);

const unitPriceLabel = computed(() =>
  model.value.taxMode === 'TAX_INCLUDED'
    ? '采购单价（含税）'
    : model.value.taxMode === 'TAX_EXCLUDED'
      ? '采购单价（未税）'
      : '采购单价'
);

/** 已选项目但该项目没有材料清单 */
const noMaterialHint = computed(
  () => !!model.value.projectId && !materialLoading.value && materials.value.length === 0
);

/** 校验器：数字必须大于 0（0 不是有效选择值） */
function positiveNumber(message: string) {
  return (_rule: unknown, value: number | null) => (Number(value ?? 0) > 0 ? true : new Error(message));
}

const rules: FormRules = {
  projectId: [{ validator: positiveNumber('请选择项目'), trigger: ['change', 'blur'] }],
  materialId: [{ validator: positiveNumber('请选择材料'), trigger: ['change', 'blur'] }],
  vendorId: [{ validator: positiveNumber('请选择供应商'), trigger: ['change', 'blur'] }],
  purchaseDate: [{ required: true, message: '请选择采购日期', trigger: ['change', 'blur'] }],
  quantity: [{ validator: positiveNumber('采购数量必须大于 0'), trigger: ['change', 'blur'] }],
  unitPrice: [
    {
      validator: (_rule: unknown, value: number | null) =>
        Number(value ?? 0) >= 0 ? true : new Error('采购单价不能为负数'),
      trigger: ['change', 'blur']
    }
  ],
  taxRate: [
    {
      validator: (_rule: unknown, value: number | null) => Number(value ?? 0) >= 0 && Number(value ?? 0) <= 100,
      message: '税率必须在 0% 到 100% 之间',
      trigger: ['change', 'blur']
    }
  ]
};

/** 加载指定项目的材料清单（切换项目时重新拉取） */
async function loadMaterialOptions(projectId: number) {
  if (!projectId) {
    materials.value = [];

    return;
  }

  materialLoading.value = true;

  try {
    const { data, error } = await fetchGetMaterialList({ projectId, current: 1, size: 200 });

    materials.value = !error && data ? (data.records ?? []) : [];
  } finally {
    materialLoading.value = false;
  }
}

/** 加载供应商选项 */
async function loadVendorOptions() {
  const { data, error } = await fetchGetSupplierList({ current: 1, size: 200 });

  if (!error && data) {
    vendorOptions.value = (data.records ?? []).map(item => ({ label: item.name, value: item.id }));
  }
}

/** 切换项目：清空已选材料与单价，并重新加载项目材料 */
async function handleProjectChange(value: number | null) {
  model.value.projectId = value ?? 0;
  model.value.materialId = 0;
  model.value.unitPrice = 0;
  model.value.taxMode = 'NO_INVOICE';
  model.value.taxRate = 0;

  await loadMaterialOptions(model.value.projectId);
}

/** 选择材料：带出清单单价 */
function handleMaterialChange(value: number | null) {
  const material = materials.value.find(item => item.id === value);
  model.value.materialId = value ?? 0;
  model.value.unitPrice = material?.unitPrice ?? 0;
  model.value.taxMode = material?.taxMode || 'NO_INVOICE';
  model.value.taxRate = material?.taxRate ?? 0;
}

watch(visible, async value => {
  if (!value) {
    return;
  }

  model.value = createDefaultModel();

  await nextTick();

  restoreValidation();

  await Promise.all([loadVendorOptions(), loadMaterialOptions(model.value.projectId)]);
});

async function handleSubmit() {
  try {
    await validate();
  } catch {
    window.$message?.warning('请完善项目、材料、供应商、采购日期、数量和单价');

    return;
  }

  submitting.value = true;

  try {
    const { error } = await fetchCreatePurchaseOrder({
      ...model.value,
      projectId: Number(model.value.projectId),
      materialId: Number(model.value.materialId)
    });

    if (error) {
      return;
    }

    window.$message?.success('采购单已创建，请到付款管理处理付款');

    emit('submitted', Number(model.value.projectId));

    visible.value = false;
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth" :mask-closable="false">
    <NDrawerContent title="新建采购单" closable :native-scrollbar="false" @close="visible = false">
      <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        label-width="86"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="项目" path="projectId">
          <NSelect
            :value="model.projectId || null"
            :options="projectSelectOptions"
            filterable
            placeholder="请选择项目"
            @update:value="handleProjectChange"
          />
        </NFormItem>
        <NFormItem label="材料" path="materialId">
          <NSelect
            v-model:value="model.materialId"
            :options="materialOptions"
            :loading="materialLoading"
            :disabled="!model.projectId || noMaterialHint"
            filterable
            placeholder="请选择项目材料"
            @update:value="handleMaterialChange"
          />
        </NFormItem>
        <NFormItem v-if="noMaterialHint" label=" ">
          <span class="text-12px text-warning">该项目暂无材料清单，请先到「材料清单」中新增材料</span>
        </NFormItem>
        <NFormItem label="供应商" path="vendorId">
          <NSelect v-model:value="model.vendorId" :options="vendorOptions" filterable placeholder="请选择供应商" />
        </NFormItem>
        <NFormItem label="采购日期" path="purchaseDate">
          <NDatePicker
            v-model:formatted-value="model.purchaseDate"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
          />
        </NFormItem>
        <div v-if="selectedMaterial" class="mb-16px rounded-6px bg-gray-100 p-12px text-13px dark:bg-gray-800">
          <div>规格：{{ selectedMaterial.spec || '-' }}　单位：{{ selectedMaterial.unit || '-' }}</div>
          <div class="mt-4px text-gray-500">材料清单单价：{{ formatMoney(selectedMaterial.unitPrice) }}</div>
        </div>
        <NFormItem label="采购数量" path="quantity">
          <NInputNumber
            v-model:value="model.quantity"
            class="w-full"
            :min="0.01"
            :precision="2"
            :show-button="false"
            placeholder="请输入采购数量"
          />
        </NFormItem>
        <NFormItem :label="unitPriceLabel" path="unitPrice">
          <NInputNumber v-model:value="model.unitPrice" class="w-full" :min="0" :precision="2" :show-button="false">
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
        <NFormItem label="税率" path="taxRate">
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
        <NFormItem label="未税金额">
          <NInput :value="formatMoney(taxAmounts.netAmount)" disabled />
        </NFormItem>
        <NFormItem label="税额">
          <NInput :value="formatMoney(taxAmounts.taxAmount)" disabled />
        </NFormItem>
        <NFormItem label="应付总额">
          <NInput :value="formatMoney(taxAmounts.totalAmount)" disabled />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" :disabled="noMaterialHint" @click="handleSubmit">
            创建待付款单
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
