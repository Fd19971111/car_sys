<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { useProjectData } from '@/views/project/data/components/use-project-data';

defineOptions({
  name: 'UsageCreateDrawer'
});

interface Props {
  /** 当前选中的项目 */
  projectId: number | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));

const { formRef, validate, restoreValidation } = useNaiveForm();

/** 材料列表用于下拉选择（每次打开抽屉重新拉取，保证新增材料能立即选到） */
const { materials, materialLoading, loadMaterials, createUsage } = useProjectData();

const submitting = ref(false);

function createDefaultModel(): Api.ProjectData.UsageForm {
  return {
    projectId: props.projectId,
    materialId: null,
    quantity: null,
    usedAt: dayjs().format('YYYY-MM-DD'),
    receiver: '',
    remark: ''
  };
}

const model = ref<Api.ProjectData.UsageForm>(createDefaultModel());

const materialOptions = computed(() =>
  materials.value.map(item => ({
    label: `${item.name}${item.spec ? ` · ${item.spec}` : ''}`,
    value: item.id
  }))
);

/** 选中材料：展示剩余数量、单位与单价，避免领超 */
const selectedMaterial = computed(() => materials.value.find(item => item.id === model.value.materialId) ?? null);

const rules: FormRules = {
  materialId: [{ required: true, type: 'number', message: '请选择领取的材料', trigger: ['change', 'blur'] }],
  quantity: [
    { required: true, type: 'number', message: '请填写领取数量', trigger: ['change', 'blur'] },
    {
      validator: (_rule, value: number | null) => {
        if (Number(value ?? 0) <= 0) {
          return new Error('领取数量必须大于 0');
        }

        const remain = selectedMaterial.value?.remainQuantity;

        if (remain !== undefined && remain !== null && Number(value) > remain) {
          return new Error(`领取数量不能超过剩余数量 ${remain}`);
        }

        return true;
      },
      trigger: ['change', 'blur']
    }
  ],
  usedAt: [{ required: true, message: '请选择领取日期', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  model.value = createDefaultModel();

  await nextTick();

  restoreValidation();

  if (props.projectId) {
    await loadMaterials({ projectId: props.projectId, current: 1, size: 200 });
  }
});

async function handleSubmit() {
  try {
    await validate();
  } catch {
    window.$message?.warning('请检查表单填写是否正确');

    return;
  }

  if (!props.projectId) {
    window.$message?.warning('请先选择项目');

    return;
  }

  submitting.value = true;

  try {
    const success = await createUsage({ ...model.value, projectId: props.projectId });

    if (!success) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    window.$message?.success('领取记录已保存');

    emit('submitted');

    visible.value = false;
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth" :mask-closable="false">
    <NDrawerContent title="新增领取" closable :native-scrollbar="false" @close="visible = false">
      <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        label-width="96"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="领取材料" path="materialId">
          <NSelect
            v-model:value="model.materialId"
            :options="materialOptions"
            :loading="materialLoading"
            placeholder="请选择本项目材料清单中的材料"
            filterable
            clearable
          />
        </NFormItem>
        <NFormItem v-if="selectedMaterial" label=" ">
          <span class="text-12px text-gray-400">
            剩余 {{ selectedMaterial.remainQuantity }} {{ selectedMaterial.unit }} · 单价
            {{ formatMoney(selectedMaterial.unitPrice) }}
          </span>
        </NFormItem>
        <NFormItem label="领取数量" path="quantity">
          <NInputNumber
            v-model:value="model.quantity"
            class="w-full"
            :min="0.01"
            :precision="2"
            :show-button="false"
            :placeholder="selectedMaterial ? `最多可领 ${selectedMaterial.remainQuantity} ${selectedMaterial.unit}` : '请输入领取数量'"
          >
            <template v-if="selectedMaterial?.unit" #suffix>{{ selectedMaterial.unit }}</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="领取日期" path="usedAt">
          <NDatePicker v-model:formatted-value="model.usedAt" class="w-full" type="date" value-format="yyyy-MM-dd" />
        </NFormItem>
        <NFormItem label="领取人" path="receiver">
          <NInput v-model:value="model.receiver" maxlength="50" placeholder="如：张师傅" clearable />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count placeholder="如：3 号楼客厅用" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">保存领取</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
