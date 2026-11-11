<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import type { UploadCustomRequestOptions, UploadFileInfo } from 'naive-ui';
import { fetchUploadAdminImage } from '@/service/api';

defineOptions({
  name: 'ProjectImageUpload'
});

interface Props {
  /** 图片相对路径（写库用） */
  modelValue?: string[];
  /** 图片绝对地址（编辑回显用） */
  imageUrls?: string[];
  /** 最多上传张数 */
  max?: number;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => [],
  imageUrls: () => [],
  max: 9,
  disabled: false
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: string[]): void;
}>();

interface PathFileInfo extends UploadFileInfo {
  path?: string;
}

const fileList = ref<PathFileInfo[]>([]);
const lastEmitted = ref('');

function fileName(path: string, index: number) {
  return path.split('/').pop() || `图片${index + 1}`;
}

function initFileList() {
  const paths = props.modelValue ?? [];
  const urls = props.imageUrls ?? [];

  fileList.value = paths.map((path, index) => ({
    id: `existing-${index}-${path}`,
    name: fileName(path, index),
    status: 'finished',
    url: urls[index] || path,
    path
  }));
}

watch(
  () => [props.modelValue, props.imageUrls],
  () => {
    if (JSON.stringify(props.modelValue ?? []) === lastEmitted.value) {
      return;
    }

    initFileList();
  },
  { immediate: true, deep: true }
);

function collectPaths(list: PathFileInfo[]) {
  return list.map(item => item.path).filter((path): path is string => Boolean(path));
}

function emitModel() {
  const paths = collectPaths(fileList.value);
  lastEmitted.value = JSON.stringify(paths);
  emit('update:modelValue', paths);
}

function handleUpdate(list: UploadFileInfo[]) {
  fileList.value = list as PathFileInfo[];
  emitModel();
}

async function customRequest({ file, onFinish, onError }: UploadCustomRequestOptions) {
  if (!file.file) {
    onError();

    return;
  }

  const { data, error } = await fetchUploadAdminImage(file.file);

  if (error || !data?.files?.length) {
    onError();

    return;
  }

  const uploaded = data.files[0];
  const target = fileList.value.find(item => item.id === file.id);

  if (target) {
    target.path = uploaded.path;
    target.url = uploaded.url;
    target.status = 'finished';
  }

  onFinish();
  emitModel();
}
</script>

<template>
  <NUpload
    :file-list="fileList"
    list-type="image-card"
    :max="max"
    accept="image/*"
    :disabled="disabled"
    :custom-request="customRequest"
    @update:file-list="handleUpdate"
  >
    <div class="upload-tip">
      <SvgIcon icon="mdi:plus" />
      <span>上传</span>
    </div>
  </NUpload>
</template>

<style scoped>
.upload-tip {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  width: 100%;
  height: 100%;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.45);
}
</style>
