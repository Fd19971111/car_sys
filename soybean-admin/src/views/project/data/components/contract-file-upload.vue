<script setup lang="ts">
import { ref, watch } from 'vue';
import type { UploadCustomRequestOptions, UploadFileInfo } from 'naive-ui';
import { fetchUploadAdminContractFile } from '@/service/api';

defineOptions({ name: 'ContractFileUpload' });

interface ContractFileInfo extends UploadFileInfo {
  path?: string;
}

const props = withDefaults(defineProps<{ modelValue?: Api.ProjectData.ContractAttachment[] }>(), {
  modelValue: () => []
});

const emit = defineEmits<{
  (event: 'update:modelValue', value: Api.ProjectData.ContractAttachment[]): void;
}>();

const fileList = ref<ContractFileInfo[]>([]);
const lastEmitted = ref('');

function syncFromModel() {
  fileList.value = (props.modelValue ?? []).map((file, index) => ({
    id: `contract-${index}-${file.path}`,
    name: file.name,
    status: 'finished',
    url: file.url,
    path: file.path
  }));
}

watch(
  () => props.modelValue,
  value => {
    if (JSON.stringify(value ?? []) !== lastEmitted.value) syncFromModel();
  },
  { immediate: true, deep: true }
);

function emitModel() {
  const attachments = fileList.value
    .filter(file => file.path)
    .map(file => ({ path: file.path!, name: file.name, url: file.url || undefined }));
  lastEmitted.value = JSON.stringify(attachments);
  emit('update:modelValue', attachments);
}

function handleUpdate(list: UploadFileInfo[]) {
  fileList.value = list as ContractFileInfo[];
  emitModel();
}

async function customRequest({ file, onFinish, onError }: UploadCustomRequestOptions) {
  if (!file.file || file.file.size > 20 * 1024 * 1024) {
    window.$message?.error('请选择 20MB 以内的文件');
    onError();
    return;
  }

  const { data, error } = await fetchUploadAdminContractFile(file.file);
  if (error || !data?.files?.length) {
    onError();
    return;
  }

  const uploaded = data.files[0];
  const target = fileList.value.find(item => item.id === file.id);
  if (target) {
    target.path = uploaded.path;
    target.name = uploaded.name;
    target.url = uploaded.url;
    target.status = 'finished';
  }

  onFinish();
  emitModel();
}
</script>

<template>
  <div class="w-full">
    <NUpload
      :file-list="fileList"
      accept=".pdf,.doc,.docx,.xls,.xlsx,.jpg,.jpeg,.png"
      :max="9"
      :custom-request="customRequest"
      @update:file-list="handleUpdate"
    >
      <NButton size="small" secondary>
        <template #icon><SvgIcon icon="mdi:paperclip" /></template>
        添加附件
      </NButton>
    </NUpload>
    <span class="mt-6px block text-12px text-gray-400">PDF、Word、Excel 或 JPG/PNG，最多 9 个，每个不超过 20MB</span>
  </div>
</template>
