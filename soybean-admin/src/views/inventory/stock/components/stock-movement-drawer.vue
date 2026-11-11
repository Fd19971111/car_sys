<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormInst, FormItemRule, FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { fetchGetMaterialList, fetchGetSupplierList, fetchStockAdjust, fetchStockIn, fetchStockOut } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { useFundAccounts } from '@/views/finance/components/use-fund-accounts';

defineOptions({ name: 'StockMovementDrawer' });

type MovementMode = 'in' | 'out' | 'adjust';

interface Props {
  mode: MovementMode;
  /** 行级快捷入口选中的库存行（材料），顶级按钮进入时为 null */
  defaultMaterial?: Api.Inventory.Stock | null;
  /** 项目下拉选项 */
  projectOptions: Api.Project.SelectOption[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));

const { fundAccountOptions, hasFundAccounts, defaultFundAccountId, loadFundAccounts } = useFundAccounts();

const formRef = ref<FormInst | null>(null);
const submitting = ref(false);

/** 项目材料清单（下拉选择，材料是带单价 / 单位属性的 item） */
const materials = ref<Api.ProjectData.Material[]>([]);
const materialLoading = ref(false);
const vendorOptions = ref<CommonType.Option<number, string>[]>([]);

const isAdjust = computed(() => props.mode === 'adjust');
const adjustDirection = ref<'IN' | 'OUT'>('IN');

const title = computed(() => {
  if (props.mode === 'in') return '材料入库';
  if (props.mode === 'out') return '材料出库';
  return '库存补录';
});

const submitText = computed(() => {
  if (props.mode === 'in') return '确认入库';
  if (props.mode === 'out') return '确认出库';
  return '确认补录';
});

/** 是否展示单价输入（入库 / 出库 / 补录入库），0 表示沿用材料清单单价 */
const showUnitPrice = computed(() => isAdjust.value ? adjustDirection.value === 'IN' : true);
/** 单价字段文案 */
const unitPriceLabel = computed(() =>
  props.mode === 'out' || (isAdjust.value && adjustDirection.value === 'OUT') ? '出库单价' : '入库单价'
);
/** 是否展示领用人（出库 / 补录出库） */
const showReceiver = computed(() => props.mode === 'out' || (isAdjust.value && adjustDirection.value === 'OUT'));
/**
 * 是否展示财务联动
 * 出库（项目领用）只扣库存、计入项目材料成本，不生成任何财务单据，因此不展示财务联动
 */
const showFinance = computed(() => props.mode === 'in');

interface MovementModel {
  projectId: number | null;
  materialId: number | null;
  quantity: number | null;
  unitPrice: number | null;
  generateFinance: boolean;
  vendorId: number | null;
  fundAccountId: number | null;
  receiver: string;
  bizDate: string;
  remark: string;
}

function createDefaultModel(): MovementModel {
  return {
    projectId: props.defaultMaterial?.projectId ?? null,
    materialId: props.defaultMaterial?.id ?? null,
    quantity: null,
    unitPrice: props.defaultMaterial?.unitPrice ?? null,
    generateFinance: true,
    vendorId: null,
    fundAccountId: null,
    receiver: '',
    bizDate: dayjs().format('YYYY-MM-DD'),
    remark: ''
  };
}

const model = ref<MovementModel>(createDefaultModel());

const projectSelectOptions = computed(() =>
  props.projectOptions.map(item => ({ label: item.label, value: item.value }))
);

const materialOptions = computed(() =>
  materials.value.map(item => ({
    label: `${item.name}${item.spec ? ` · ${item.spec}` : ''}`,
    value: item.id
  }))
);

const selectedMaterial = computed(() => materials.value.find(item => item.id === model.value.materialId) ?? null);

const orderAmount = computed(() =>
  Number((Number(model.value.quantity ?? 0) * Number(model.value.unitPrice ?? selectedMaterial.value?.unitPrice ?? 0)).toFixed(2))
);

const financeHint = '入库后将自动生成一张材料款付款单并记账';

const noMaterialHint = computed(
  () => !!model.value.projectId && !materialLoading.value && materials.value.length === 0
);

/** 校验器：数字必须大于 0 */
function positiveNumber(message: string): FormItemRule['validator'] {
  return (_rule, value) => (Number(value ?? 0) > 0 ? true : new Error(message));
}

const rules = computed<FormRules>(() => {
  const quantityRules: FormItemRule[] = [
    {
      validator: positiveNumber(props.mode === 'in' ? '入库数量必须大于 0' : '数量必须大于 0'),
      trigger: ['change', 'blur']
    }
  ];

  if (props.mode === 'out' || (isAdjust.value && adjustDirection.value === 'OUT')) {
    quantityRules.push({
      validator: (_rule, value) => {
        const remain = selectedMaterial.value?.remainQuantity;

        if (remain !== undefined && remain !== null && Number(value) > remain) {
          return new Error(`数量不能超过当前结存 ${remain}`);
        }

        return true;
      },
      trigger: ['change', 'blur']
    });
  }

  const base: FormRules = {
    projectId: [{ validator: positiveNumber('请选择项目'), trigger: ['change', 'blur'] }],
    materialId: [{ validator: positiveNumber('请选择材料'), trigger: ['change', 'blur'] }],
    quantity: quantityRules,
    bizDate: [{ required: true, message: '请选择日期', trigger: ['change', 'blur'] }]
  };

  if (showFinance.value && model.value.generateFinance) {
    base.fundAccountId = [{ validator: positiveNumber('请选择资金账户'), trigger: ['change', 'blur'] }];
    base.vendorId = [{ validator: positiveNumber('请选择供应商'), trigger: ['change', 'blur'] }];
  }

  return base;
});

/** 加载指定项目的材料清单 */
async function loadMaterialOptions(projectId: number | null) {
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

/** 加载供应商选项（付款单收款方） */
async function loadVendorOptions() {
  const { data, error } = await fetchGetSupplierList({ current: 1, size: 200 });

  if (!error && data) {
    vendorOptions.value = (data.records ?? []).map(item => ({ label: item.name, value: item.id }));
  }
}

function handleProjectChange(value: number | null) {
  model.value.projectId = value ?? null;
  model.value.materialId = null;
  model.value.unitPrice = null;

  void loadMaterialOptions(value);
}

function handleMaterialChange(value: number | null) {
  model.value.materialId = value ?? null;
  model.value.unitPrice = materials.value.find(item => item.id === value)?.unitPrice ?? null;
}

function handleAdjustDirectionChange(value: 'IN' | 'OUT') {
  adjustDirection.value = value;
  model.value.unitPrice = null;
}

watch(visible, async value => {
  if (!value) {
    return;
  }

  model.value = createDefaultModel();
  adjustDirection.value = 'IN';

  await nextTick();
  formRef.value?.restoreValidation();

  await Promise.all([loadFundAccounts(), loadVendorOptions()]);

  await loadMaterialOptions(model.value.projectId);

  // 只有一个资金账户时自动选中，减少录入成本
  model.value.fundAccountId = defaultFundAccountId.value;
});

async function handleSubmit() {
  try {
    await formRef.value?.validate();
  } catch {
    window.$message?.warning('请检查表单填写是否正确');
    return;
  }

  submitting.value = true;

  try {
    let success = false;

    if (props.mode === 'in') {
      const { error } = await fetchStockIn({
        projectId: model.value.projectId ?? 0,
        materialId: model.value.materialId ?? 0,
        quantity: Number(model.value.quantity ?? 0),
        unitPrice: Number(model.value.unitPrice ?? 0),
        generateFinance: model.value.generateFinance,
        vendorId: model.value.vendorId ?? 0,
        fundAccountId: model.value.fundAccountId ?? 0,
        bizDate: model.value.bizDate,
        remark: model.value.remark
      });

      success = !error;
    } else if (props.mode === 'out') {
      const { error } = await fetchStockOut({
        projectId: model.value.projectId ?? 0,
        materialId: model.value.materialId ?? 0,
        quantity: Number(model.value.quantity ?? 0),
        unitPrice: Number(model.value.unitPrice ?? 0),
        receiver: model.value.receiver,
        bizDate: model.value.bizDate,
        remark: model.value.remark
      });

      success = !error;
    } else {
      const { error } = await fetchStockAdjust({
        recordType: adjustDirection.value,
        projectId: model.value.projectId ?? 0,
        materialId: model.value.materialId ?? 0,
        quantity: Number(model.value.quantity ?? 0),
        unitPrice: Number(model.value.unitPrice ?? 0),
        receiver: model.value.receiver,
        bizDate: model.value.bizDate,
        remark: model.value.remark
      });

      success = !error;
    }

    if (!success) {
      window.$message?.error('操作失败，请稍后重试');
      return;
    }

    window.$message?.success(
      props.mode === 'in'
        ? `入库完成${model.value.generateFinance ? '，已生成付款单' : ''}`
        : props.mode === 'out'
          ? '出库完成，已计入项目材料成本'
          : '补录完成'
    );

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
      <NAlert v-if="isAdjust" type="warning" :show-icon="true" class="mb-16px">
        补录用于补记历史出入库数据：只调整库存数量，不会生成收付款单等财务数据。
      </NAlert>

      <NAlert v-if="props.mode === 'out'" type="info" :show-icon="true" class="mb-16px">
        项目领用只扣减库存并计入项目材料成本，不生成收款单 / 付款单等财务单据。
      </NAlert>

      <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        label-width="96"
        require-mark-placement="right-hanging"
      >
        <NFormItem v-if="isAdjust" label="补录方向">
          <NRadioGroup :value="adjustDirection" @update:value="handleAdjustDirectionChange">
            <NRadio value="IN">补录入库（数量增加）</NRadio>
            <NRadio value="OUT">补录出库（数量减少）</NRadio>
          </NRadioGroup>
        </NFormItem>

        <NFormItem label="所属项目" path="projectId">
          <NSelect
            v-model:value="model.projectId"
            :options="projectSelectOptions"
            :disabled="!!props.defaultMaterial"
            filterable
            placeholder="请选择项目"
            @update:value="handleProjectChange"
          />
        </NFormItem>

        <NFormItem :label="isAdjust ? '补录材料' : '选择材料'" path="materialId">
          <NSelect
            v-model:value="model.materialId"
            :options="materialOptions"
            :loading="materialLoading"
            :disabled="!!props.defaultMaterial"
            filterable
            placeholder="从材料清单中选择材料"
            @update:value="handleMaterialChange"
          />
        </NFormItem>

        <NFormItem v-if="noMaterialHint" label=" ">
          <span class="text-12px text-warning">该项目还没有材料，请先在「库存管理 - 材料清单」中添加</span>
        </NFormItem>

        <NFormItem v-if="selectedMaterial" label="材料信息">
          <span class="text-12px text-gray-500">
            当前结存 {{ selectedMaterial.remainQuantity }} {{ selectedMaterial.unit || '' }} · 清单单价
            {{ formatMoney(selectedMaterial.unitPrice) }} · 仓库
            {{ selectedMaterial.warehouseName || '-' }}
          </span>
        </NFormItem>

        <NFormItem :label="isAdjust ? (adjustDirection === 'IN' ? '入库数量' : '出库数量') : props.mode === 'in' ? '入库数量' : '出库数量'" path="quantity">
          <NInputNumber
            v-model:value="model.quantity"
            class="w-full"
            :min="0.01"
            :precision="2"
            :show-button="false"
            :placeholder="selectedMaterial ? `当前结存 ${selectedMaterial.remainQuantity} ${selectedMaterial.unit || ''}` : '请输入数量'"
          >
            <template v-if="selectedMaterial?.unit" #suffix>{{ selectedMaterial.unit }}</template>
          </NInputNumber>
        </NFormItem>

        <NFormItem v-if="showUnitPrice" :label="unitPriceLabel" path="unitPrice">
          <NInputNumber v-model:value="model.unitPrice" class="w-full" :min="0" :precision="2" :show-button="false" placeholder="默认使用材料清单单价，可手动修改">
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>

        <NFormItem v-if="showReceiver" label="领用人" path="receiver">
          <NInput v-model:value="model.receiver" maxlength="50" placeholder="如：张师傅（补录时可不填）" clearable />
        </NFormItem>

        <NFormItem :label="isAdjust ? '补录日期' : props.mode === 'in' ? '入库日期' : '出库日期'" path="bizDate">
          <NDatePicker v-model:formatted-value="model.bizDate" class="w-full" type="date" value-format="yyyy-MM-dd" />
        </NFormItem>

        <template v-if="showFinance">
          <NDivider class="!my-8px" />
          <NFormItem label="财务联动">
            <div class="w-full">
              <NSwitch v-model:value="model.generateFinance" />
              <div class="mt-4px text-12px text-gray-400">{{ financeHint }}</div>
            </div>
          </NFormItem>

          <template v-if="model.generateFinance">
            <NFormItem label="应付金额">
              <span class="font-medium">{{ formatMoney(orderAmount) }}</span>
            </NFormItem>

            <NFormItem label="供应商" path="vendorId">
              <NSelect v-model:value="model.vendorId" :options="vendorOptions" filterable placeholder="请选择供应商（付款单收款方）" />
            </NFormItem>

            <NFormItem label="资金账户" path="fundAccountId">
              <NSelect v-model:value="model.fundAccountId" :options="fundAccountOptions" filterable placeholder="请选择资金账户（用于自动记账）" />
            </NFormItem>

            <NFormItem v-if="!hasFundAccounts" label=" ">
              <span class="text-12px text-error">尚未配置可用资金账户，请先在「财务管理 - 财务初始化」中添加，或关闭财务联动</span>
            </NFormItem>
          </template>
        </template>

        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count placeholder="选填" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">{{ submitText }}</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
