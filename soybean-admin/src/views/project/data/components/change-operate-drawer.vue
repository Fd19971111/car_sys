<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import ProjectImageUpload from './image-upload.vue';
import { useProjectData } from './use-project-data';

defineOptions({
  name: 'ChangeOperateDrawer'
});

interface Props {
  projectId: number | null;
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.ProjectData.Change | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 560));

const { formRef, validate, restoreValidation } = useNaiveForm();

const { createChange, updateChange } = useProjectData();

const title = computed(() => (props.operateType === 'add' ? '新增变更记录' : '编辑变更记录'));

const submitting = ref(false);

function createDefaultModel(): Api.ProjectData.ChangeForm {
  return {
    projectId: props.projectId,
    changeType: 'ADD',
    title: '',
    amount: null,
    changeDate: dayjs().format('YYYY-MM-DD'),
    reason: '',
    images: [],
    remark: ''
  };
}

const model = ref<Api.ProjectData.ChangeForm>(createDefaultModel());

const imageUrls = ref<string[]>([]);

const rules: FormRules = {
  changeType: [{ required: true, message: '请选择变更类型', trigger: ['change', 'blur'] }],
  title: [{ required: true, message: '请填写变更事项', trigger: ['input', 'blur'] }],
  amount: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 0,
      message: '变更金额不能为负数',
      trigger: ['change', 'blur']
    }
  ],
  changeDate: [{ required: true, message: '请选择变更日期', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { projectId, changeType, title: changeTitle, amount, changeDate, reason, images, remark } = props.rowData;

    model.value = {
      projectId,
      changeType,
      title: changeTitle,
      amount,
      changeDate,
      reason,
      images: [...(images ?? [])],
      remark
    };
    imageUrls.value = props.rowData.imageUrls ?? [];
  } else {
    model.value = createDefaultModel();
    imageUrls.value = [];
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

  if (!model.value.projectId) {
    window.$message?.warning('请先选择项目');

    return;
  }

  submitting.value = true;

  try {
    const success =
      props.operateType === 'add'
        ? await createChange(model.value)
        : await updateChange(props.rowData?.id ?? 0, model.value);

    if (!success) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    window.$message?.success(props.operateType === 'add' ? '变更记录已保存' : '变更记录已更新');
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
        label-width="92"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="变更类型" path="changeType">
          <NRadioGroup v-model:value="model.changeType">
            <NRadioButton value="ADD">增项（加钱）</NRadioButton>
            <NRadioButton value="REDUCE">减项（减钱）</NRadioButton>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="变更事项" path="title">
          <NInput v-model:value="model.title" maxlength="100" placeholder="如：增加阳台柜、减少吊顶面积" clearable />
        </NFormItem>
        <NFormItem label="变更金额" path="amount">
          <NInputNumber
            v-model:value="model.amount"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            placeholder="请输入变更金额"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="变更日期" path="changeDate">
          <NDatePicker
            v-model:formatted-value="model.changeDate"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
            clearable
          />
        </NFormItem>
        <NFormItem label="变更原因" path="reason">
          <NInput
            v-model:value="model.reason"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-count
            placeholder="如：业主现场要求增加、原方案尺寸调整"
          />
        </NFormItem>
        <NFormItem label="附件照片">
          <div class="w-full flex-col gap-6px">
            <ProjectImageUpload v-model="model.images" :image-urls="imageUrls" :max="9" />
            <span class="text-12px text-gray-400">可上传变更前后的现场照片、签字单等，最多 9 张。</span>
          </div>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" maxlength="255" placeholder="选填" clearable />
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

<style scoped></style>
