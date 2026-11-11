<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { REG_PHONE } from '@/constants/reg';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { GENDER_OPTIONS, POSITION_OPTIONS, STAFF_STATUS_OPTIONS } from './staff-meta';
import { useStaffData } from './use-staff-data';

defineOptions({ name: 'StaffOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Staff.Item | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 460));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createStaff, updateStaff, departmentOptions, loadDepartmentOptions } = useStaffData();

const title = computed(() => (props.operateType === 'add' ? '新增人员' : '编辑人员'));
const submitting = ref(false);

type StaffFormModel = Omit<Api.Staff.FormData, 'departmentId'> & { departmentId: number | null };

function createDefaultModel(): StaffFormModel {
  return {
    name: '',
    phone: '',
    gender: '男',
    position: '',
    departmentId: null,
    joinDate: null as any,
    status: 1,
    remark: ''
  };
}

const model = ref<StaffFormModel>(createDefaultModel());

const rules: FormRules = {
  name: [
    { required: true, message: '请输入姓名', trigger: ['input', 'blur'] },
    { min: 2, max: 20, message: '姓名为 2-20 个字符', trigger: ['input', 'blur'] }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: ['input', 'blur'] },
    { pattern: REG_PHONE, message: '手机号格式不正确', trigger: ['input', 'blur'] }
  ],
  departmentId: [{ required: true, type: 'number', message: '请选择所属部门', trigger: ['change', 'blur'], validator: (_rule, value: number | null) => typeof value === 'number' && value > 0 }],
  status: [{ required: true, type: 'number', message: '请选择状态', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const { name, phone, gender, position, departmentId, joinDate, status, remark } = props.rowData;
    model.value = { name, phone, gender, position, departmentId, joinDate: joinDate || null as any, status, remark };
  } else {
    model.value = createDefaultModel();
  }
  await loadDepartmentOptions();
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
    const formData: Api.Staff.FormData = {
      ...model.value,
      departmentId: model.value.departmentId as number
    };
    const success =
      props.operateType === 'add'
        ? await createStaff(formData)
        : await updateStaff(props.rowData?.id ?? 0, formData);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增人员成功' : '人员信息已更新');
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
        <NFormItem label="姓名" path="name">
          <NInput v-model:value="model.name" placeholder="请输入姓名" maxlength="20" clearable />
        </NFormItem>
        <NFormItem label="手机号" path="phone">
          <NInput v-model:value="model.phone" placeholder="请输入 11 位手机号" maxlength="11" clearable />
        </NFormItem>
        <NFormItem label="性别" path="gender">
          <NRadioGroup v-model:value="model.gender">
            <NRadioButton v-for="item in GENDER_OPTIONS" :key="item.value" :value="item.value">{{ item.label }}</NRadioButton>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="职位" path="position">
          <NSelect v-model:value="model.position" :options="POSITION_OPTIONS" placeholder="请选择职位" filterable tag clearable />
        </NFormItem>
        <NFormItem label="所属部门" path="departmentId">
          <NSelect v-model:value="model.departmentId" :options="departmentOptions" placeholder="请选择部门" filterable />
        </NFormItem>
        <NFormItem label="入职日期" path="joinDate">
          <NDatePicker v-model:formatted-value="model.joinDate" class="w-full" type="date" value-format="yyyy-MM-dd" placeholder="请选择入职日期" clearable />
        </NFormItem>
        <NFormItem label="状态" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadioButton v-for="item in STAFF_STATUS_OPTIONS" :key="item.value" :value="item.value">{{ item.label }}</NRadioButton>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="4" maxlength="200" show-count placeholder="补充说明" />
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