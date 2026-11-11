<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { fetchSaveFinAccount } from '@/service/api';
import { ACCOUNT_CATEGORY_OPTIONS, deriveDirection, getDirectionMeta } from './basic-meta';

defineOptions({
  name: 'FinanceBasicAccountDrawer'
});

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.FinBasic.Account | null;
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

const title = computed(() => (props.operateType === 'add' ? '新增会计科目' : '编辑会计科目'));

const submitting = ref(false);

function createDefaultModel(): Api.FinBasic.AccountFormData {
  return {
    code: '',
    name: '',
    categoryCode: 'ASSET',
    enabled: true
  };
}

const model = ref<Api.FinBasic.AccountFormData>(createDefaultModel());

const isSystem = computed(() => props.operateType === 'edit' && Boolean(props.rowData?.isSystem));

const direction = computed(() => getDirectionMeta(deriveDirection(model.value.categoryCode)));

const rules: FormRules = {
  code: [
    { required: true, message: '请输入科目编码', trigger: ['input', 'blur'] },
    {
      pattern: /^[0-9A-Za-z][0-9A-Za-z._-]{0,31}$/,
      message: '编码只能由字母、数字、点、中划线或下划线组成，最长32位',
      trigger: ['input', 'blur']
    }
  ],
  name: [
    { required: true, message: '请输入科目名称', trigger: ['input', 'blur'] },
    { max: 100, message: '科目名称不能超过100个字符', trigger: ['input', 'blur'] }
  ],
  categoryCode: [{ required: true, message: '请选择科目类别', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { id, code, name, categoryCode, enabled } = props.rowData;

    model.value = { id, code, name, categoryCode, enabled };
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
    const { error } = await fetchSaveFinAccount({ ...model.value });

    if (error) {
      return;
    }

    window.$message?.success(props.operateType === 'add' ? '会计科目已新增' : '会计科目已更新');

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
      <NAlert v-if="isSystem" class="mb-16px" type="info" :bordered="false" title="系统预置科目">
        编码与类别已锁定，仅支持修改名称和启用状态。
      </NAlert>

      <NForm ref="formRef" :model="model" :rules="rules" label-placement="top" require-mark-placement="right-hanging">
        <div class="grid grid-cols-1 gap-x-16px sm:grid-cols-2">
          <NFormItem label="科目编码" path="code">
            <NInput
              v-model:value="model.code"
              :disabled="isSystem"
              maxlength="32"
              placeholder="如 1002、6602.01"
            />
          </NFormItem>
          <NFormItem label="科目类别" path="categoryCode">
            <NSelect
              v-model:value="model.categoryCode"
              :options="ACCOUNT_CATEGORY_OPTIONS"
              :disabled="isSystem"
              placeholder="请选择科目类别"
            />
          </NFormItem>
        </div>

        <NFormItem label="科目名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" show-count placeholder="请输入科目名称" />
        </NFormItem>

        <div class="direction-tip">
          <div class="direction-tip__label">余额方向</div>
          <div class="direction-tip__value" :style="{ color: direction.color, backgroundColor: direction.bg }">
            {{ direction.label }}
          </div>
          <span class="direction-tip__desc">由科目类别自动确定，无需手动选择</span>
        </div>

        <NFormItem label="启用状态">
          <div class="flex-y-center gap-12px">
            <NSwitch v-model:value="model.enabled" />
            <span class="text-13px text-gray-400">{{ model.enabled ? '启用后可用于新凭证' : '停用后不再参与新记账' }}</span>
          </div>
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">
            {{ props.operateType === 'add' ? '新增科目' : '保存修改' }}
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
.direction-tip {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  margin-bottom: 18px;
  border-radius: 10px;
  background-color: rgb(var(--primary-color) / 0.05);
}

.direction-tip__label {
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.55);
}

.direction-tip__value {
  min-width: 32px;
  padding: 2px 10px;
  font-size: 13px;
  font-weight: 600;
  text-align: center;
  border-radius: 999px;
}

.direction-tip__desc {
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.4);
}
</style>
