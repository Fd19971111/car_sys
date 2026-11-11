<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { useProjectData } from './use-project-data';
import ContractFileUpload from './contract-file-upload.vue';

defineOptions({ name: 'ContractOperateDrawer' });

const props = defineProps<{
  projectId: number | null;
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.ProjectData.Contract | null;
}>();

const emit = defineEmits<{ (event: 'submitted'): void }>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 560));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createContract, updateContract } = useProjectData();
const title = computed(() => (props.operateType === 'add' ? '新增合同' : '编辑合同'));
const submitting = ref(false);

function createDefaultModel(): Api.ProjectData.ContractForm {
  return {
    projectId: props.projectId,
    name: '',
    contractNo: '',
    amount: null,
    signedAt: dayjs().format('YYYY-MM-DD'),
    remark: '',
    attachments: []
  };
}

const model = ref<Api.ProjectData.ContractForm>(createDefaultModel());
const rules: FormRules = {
  name: [{ required: true, message: '请填写合同名称', trigger: ['input', 'blur'] }],
  amount: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 0,
      message: '合同金额不能为负数',
      trigger: ['change', 'blur']
    }
  ],
  signedAt: [{ required: true, message: '请选择签订日期', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) return;

  if (props.operateType === 'edit' && props.rowData) {
    const { projectId, name, contractNo, amount, signedAt, remark, attachments } = props.rowData;
    model.value = {
      projectId,
      name,
      contractNo,
      amount,
      signedAt,
      remark,
      attachments: [...(attachments ?? [])]
    };
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
    window.$message?.warning('请检查合同信息是否填写完整');
    return;
  }
  if (!model.value.projectId) {
    window.$message?.warning('请先选择项目');
    return;
  }

  submitting.value = true;
  try {
    const success =
      props.operateType === 'add'
        ? await createContract(model.value)
        : await updateContract(props.rowData?.id ?? 0, model.value);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }

    window.$message?.success(props.operateType === 'add' ? '合同已保存' : '合同已更新');
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
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="top" require-mark-placement="right-hanging">
        <NFormItem label="合同名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" placeholder="请输入合同名称" />
        </NFormItem>
        <NFormItem label="合同编号" path="contractNo">
          <NInput v-model:value="model.contractNo" maxlength="100" placeholder="选填" />
        </NFormItem>
        <NFormItem label="合同金额" path="amount">
          <NInputNumber v-model:value="model.amount" class="w-full" :min="0" :precision="2" :show-button="false">
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="签订日期" path="signedAt">
          <NDatePicker
            v-model:formatted-value="model.signedAt"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
            clearable
          />
        </NFormItem>
        <NFormItem label="合同附件">
          <ContractFileUpload v-model="model.attachments" />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="500" show-count />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton :disabled="submitting" @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">保存</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>
