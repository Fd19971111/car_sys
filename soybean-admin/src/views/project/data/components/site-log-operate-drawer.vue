<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import ProjectImageUpload from './image-upload.vue';
import { useProjectData } from './use-project-data';

defineOptions({
  name: 'SiteLogOperateDrawer'
});

interface Props {
  projectId: number | null;
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.ProjectData.SiteLog | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 560));

const { formRef, validate, restoreValidation } = useNaiveForm();

const { createLog, updateLog } = useProjectData();

const title = computed(() => (props.operateType === 'add' ? '新增施工日志' : '编辑施工日志'));

const submitting = ref(false);

function createDefaultModel(): Api.ProjectData.SiteLogForm {
  return {
    projectId: props.projectId,
    content: '',
    images: [],
    logDate: ''
  };
}

const model = ref<Api.ProjectData.SiteLogForm>(createDefaultModel());

const imageUrls = ref<string[]>([]);

const rules: FormRules = {
  content: [{ required: true, message: '请填写日志内容', trigger: ['input', 'blur'] }]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    model.value = {
      projectId: props.rowData.projectId,
      content: props.rowData.content,
      images: [...(props.rowData.images ?? [])],
      logDate: props.rowData.logDate || ''
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
    window.$message?.warning('请填写日志内容');

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
        ? await createLog(model.value)
        : await updateLog(props.rowData?.id ?? 0, model.value);

    if (!success) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    window.$message?.success(props.operateType === 'add' ? '施工日志已发布' : '施工日志已更新');
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
        label-placement="top"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="日志日期" path="logDate">
          <NDatePicker
            v-model:formatted-value="model.logDate"
            type="date"
            class="w-full"
            value-format="yyyy-MM-dd"
            placeholder="默认为今天，补录历史日志请选择实际发生日期"
            clearable
          />
        </NFormItem>
        <NFormItem label="日志内容" path="content">
          <NInput
            v-model:value="model.content"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-count
            placeholder="今天完成了什么、现场情况、验收结果等"
          />
        </NFormItem>
        <NFormItem label="现场照片">
          <div class="w-full flex-col gap-6px">
            <ProjectImageUpload v-model="model.images" :image-urls="imageUrls" :max="9" />
            <span class="text-12px text-gray-400">
              最多 9 张。照片会同步展示在业主端小程序的项目进度里，方便直接发给业主看。
            </span>
          </div>
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
