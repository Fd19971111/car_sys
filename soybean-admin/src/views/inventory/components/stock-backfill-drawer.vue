<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { fetchGetWarehouseList, fetchStockInPurchaseOrder, fetchStockOutMaterialUsage } from '@/service/api';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';

defineOptions({ name: 'StockBackfillDrawer' });

interface Props {
  /** in=采购单入库 out=领用单出库 */
  mode: 'in' | 'out';
  /** 单据ID（采购单ID / 领用记录ID） */
  recordId: number;
  materialName: string;
  quantity: number;
  unit?: string;
  /** 默认仓库（一般取材料所属仓库） */
  defaultWarehouseId?: number;
}

const props = withDefaults(defineProps<Props>(), { unit: '', defaultWarehouseId: 0 });

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 460));

const { formRef, validate, restoreValidation } = useNaiveForm();

const submitting = ref(false);
const warehouseLoading = ref(false);
const warehouseOptions = ref<CommonType.Option<number, string>[]>([]);

const model = ref<{ warehouseId: number | null; remark: string }>({ warehouseId: null, remark: '' });

const title = computed(() => (props.mode === 'in' ? '采购单入库' : '领用单出库'));

const quantityText = computed(() => `${props.quantity}${props.unit ? ` ${props.unit}` : ''}`);

const rules: FormRules = {
  warehouseId: [{ required: true, type: 'number', message: '请选择仓库', trigger: ['change', 'blur'] }]
};

async function loadWarehouseOptions() {
  warehouseLoading.value = true;

  try {
    const { data, error } = await fetchGetWarehouseList({ current: 1, size: 100 });

    if (!error && data) {
      warehouseOptions.value = (data.records ?? []).map(item => ({ label: item.name, value: item.id }));
    }
  } finally {
    warehouseLoading.value = false;
  }
}

watch(visible, async value => {
  if (!value) {
    return;
  }

  model.value = {
    warehouseId: props.defaultWarehouseId > 0 ? props.defaultWarehouseId : null,
    remark: ''
  };

  await nextTick();

  restoreValidation();

  await loadWarehouseOptions();

  if (!model.value.warehouseId && warehouseOptions.value.length === 1) {
    model.value.warehouseId = warehouseOptions.value[0].value;
  }
});

async function handleSubmit() {
  try {
    await validate();
  } catch {
    window.$message?.warning('请选择仓库');

    return;
  }

  if (!props.recordId || !model.value.warehouseId) {
    return;
  }

  submitting.value = true;

  try {
    const payload = { id: props.recordId, warehouseId: model.value.warehouseId, remark: model.value.remark };
    const { error } =
      props.mode === 'in' ? await fetchStockInPurchaseOrder(payload) : await fetchStockOutMaterialUsage(payload);

    if (error) {
      return;
    }

    window.$message?.success(props.mode === 'in' ? '入库完成' : '出库完成');

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
        label-width="86"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="材料">
          <span class="text-13px">{{ materialName || '-' }}</span>
        </NFormItem>
        <NFormItem label="数量">
          <span class="text-13px font-mono">{{ quantityText }}</span>
        </NFormItem>
        <NFormItem label="仓库" path="warehouseId">
          <NSelect
            v-model:value="model.warehouseId"
            :options="warehouseOptions"
            :loading="warehouseLoading"
            filterable
            placeholder="请选择仓库"
          />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count />
        </NFormItem>
        <NAlert type="info" :bordered="false" :show-icon="true">
          完成后会生成一条{{ mode === 'in' ? '入库' : '出库' }}流水，该单据不再允许重复出入库。
        </NAlert>
      </NForm>
      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">
            确认{{ mode === 'in' ? '入库' : '出库' }}
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
