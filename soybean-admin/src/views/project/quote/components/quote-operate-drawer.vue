<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { useProjectOptions } from '@/views/finance/components/use-project-options';
import { QUOTE_ITEM_PRESET_OPTIONS, QUOTE_UNIT_OPTIONS } from './quote-meta';
import { useQuoteData } from './use-quote-data';

defineOptions({
  name: 'QuoteOperateDrawer'
});

interface Props {
  /** 操作类型：add=新增 edit=编辑 */
  operateType: NaiveUI.TableOperateType;
  /** 编辑时的报价明细数据 */
  rowData?: Api.Quote.Item | null;
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

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 560));

const { formRef, validate, restoreValidation } = useNaiveForm();

const { projectOptions, loadProjectOptions } = useProjectOptions();

const { createQuote, updateQuote } = useQuoteData();

const title = computed(() => (props.operateType === 'add' ? '新增报价明细' : '编辑报价明细'));

const submitting = ref(false);

function createDefaultModel(): Api.Quote.FormData {
  return {
    projectId: props.defaultProjectId ?? null,
    itemName: '',
    unit: '平米',
    price: null,
    quantity: 1,
    remark: ''
  };
}

const model = ref<Api.Quote.FormData>(createDefaultModel());

/** 小计预览：单价 × 数量 */
const previewTotal = computed(() => {
  const price = Number(model.value.price ?? 0);
  const quantity = Number(model.value.quantity ?? 0);

  return Math.round(price * quantity * 100) / 100;
});

/** 选中项目的合同金额（用于对照报价） */
const selectedProjectAmount = computed(
  () => projectOptions.value.find(item => item.value === model.value.projectId)?.amount ?? 0
);

const rules: FormRules = {
  projectId: [{ required: true, type: 'number', message: '请选择所属项目', trigger: ['change', 'blur'] }],
  itemName: [
    { required: true, message: '请输入报价项名称', trigger: ['input', 'blur'] },
    { max: 100, message: '报价项名称最多 100 个字符', trigger: ['input', 'blur'] }
  ],
  price: [{ required: true, type: 'number', message: '请输入单价', trigger: ['change', 'blur'] }],
  quantity: [{ required: true, type: 'number', message: '请输入数量', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { projectId, itemName, unit, price, quantity, remark } = props.rowData;

    model.value = {
      projectId: projectId || null,
      itemName,
      unit: unit || '项',
      price,
      quantity,
      remark
    };
  } else {
    model.value = createDefaultModel();
  }

  await nextTick();

  restoreValidation();

  // 项目下拉可能尚未加载（例如直接进入本页后立刻点新增），强制刷新以包含刚新建的项目
  await loadProjectOptions(true);
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
        ? await createQuote(model.value)
        : await updateQuote(props.rowData?.id ?? 0, model.value);

    if (!result.success) {
      return;
    }

    window.$message?.success(props.operateType === 'add' ? '新增报价明细成功' : '报价明细已更新');

    if (result.warning) {
      window.$message?.warning(result.warning);
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
        <NFormItem label="报价项名称" path="itemName">
          <NSelect
            v-model:value="model.itemName"
            :options="QUOTE_ITEM_PRESET_OPTIONS"
            placeholder="请选择或输入报价项名称，如：水电改造"
            filterable
            tag
            clearable
          />
        </NFormItem>
        <NFormItem label="单位" path="unit">
          <NSelect
            v-model:value="model.unit"
            :options="QUOTE_UNIT_OPTIONS"
            placeholder="请选择或输入单位"
            filterable
            tag
            clearable
          />
        </NFormItem>
        <NFormItem label="单价" path="price">
          <NInputNumber
            v-model:value="model.price"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            placeholder="请输入单价"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="数量" path="quantity">
          <NInputNumber
            v-model:value="model.quantity"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            placeholder="请输入数量"
          />
        </NFormItem>
        <NFormItem label="小计">
          <span class="text-16px font-medium text-primary">{{ formatMoney(previewTotal) }}</span>
          <span class="ml-8px text-12px text-gray-400">（后端按「单价 × 数量」自动计算）</span>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput
            v-model:value="model.remark"
            type="textarea"
            :rows="3"
            maxlength="255"
            show-count
            placeholder="如：不含瓷砖、按实际面积结算"
          />
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
