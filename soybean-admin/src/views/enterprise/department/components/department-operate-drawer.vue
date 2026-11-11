<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { DEPT_STATUS_OPTIONS } from './department-meta';
import { useDepartmentData } from './use-department-data';

defineOptions({ name: 'DepartmentOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Department.Item | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 420));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createDepartment, updateDepartment } = useDepartmentData();

const title = computed(() => (props.operateType === 'add' ? '新增部门' : '编辑部门'));
const submitting = ref(false);

function createDefaultModel(): Api.Department.FormData {
  return {
    name: '',
    code: '',
    leaderId: null,
    status: 1,
    sortOrder: 1,
    remark: ''
  };
}

const model = ref<Api.Department.FormData>(createDefaultModel());

const rules: FormRules = {
  name: [
    { required: true, message: '请输入部门名称', trigger: ['input', 'blur'] },
    { min: 2, max: 30, message: '部门名称为 2-30 个字符', trigger: ['input', 'blur'] }
  ],
  code: [
    { required: true, message: '请输入部门编码', trigger: ['input', 'blur'] },
    { min: 2, max: 20, message: '编码为 2-20 个字符', trigger: ['input', 'blur'] }
  ]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const { name, code, status, sortOrder } = props.rowData;
    model.value = { name, code, leaderId: null, status, sortOrder, remark: '' };
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
    const success =
      props.operateType === 'add'
        ? await createDepartment(model.value)
        : await updateDepartment(props.rowData?.id ?? 0, model.value);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增部门成功' : '部门信息已更新');
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
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="100" require-mark-placement="right-hanging">
        <NFormItem label="部门名称" path="name">
          <NInput v-model:value="model.name" placeholder="如：设计部" maxlength="30" clearable />
        </NFormItem>
        <NFormItem label="部门编码" path="code">
          <NInput v-model:value="model.code" placeholder="如：design" maxlength="20" clearable />
        </NFormItem>
        <NFormItem label="状态" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadioButton v-for="item in DEPT_STATUS_OPTIONS" :key="item.value" :value="item.value">{{ item.label }}</NRadioButton>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="排序" path="sortOrder">
          <NInputNumber v-model:value="model.sortOrder" class="w-full" :min="0" :max="9999" />
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