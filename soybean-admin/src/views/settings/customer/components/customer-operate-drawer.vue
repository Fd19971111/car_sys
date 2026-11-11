<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { REG_PHONE } from '@/constants/reg';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { CUSTOMER_SOURCE_OPTIONS, CUSTOMER_STATUS_OPTIONS } from './customer-meta';
import { useCustomerData } from './use-customer-data';

defineOptions({
  name: 'CustomerOperateDrawer'
});

interface Props {
  /** 操作类型：add=新增 edit=编辑 */
  operateType: NaiveUI.TableOperateType;
  /** 编辑时的客户数据 */
  rowData?: Api.Customer.Item | null;
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

const { createCustomer, updateCustomer } = useCustomerData();

const title = computed(() => (props.operateType === 'add' ? '新增客户' : '编辑客户'));

const submitting = ref(false);

function createDefaultModel(): Api.Customer.FormData {
  return {
    name: '',
    phone: '',
    address: '',
    source: '转介绍',
    status: 1,
    remark: ''
  };
}

const model = ref<Api.Customer.FormData>(createDefaultModel());

const rules: FormRules = {
  name: [
    { required: true, message: '请输入客户姓名', trigger: ['input', 'blur'] },
    { min: 2, max: 20, message: '客户姓名长度为 2-20 个字符', trigger: ['input', 'blur'] }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: ['input', 'blur'] },
    { pattern: REG_PHONE, message: '手机号格式不正确', trigger: ['input', 'blur'] }
  ],
  status: [{ required: true, type: 'number', message: '请选择客户状态', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { name, phone, address, source, status, remark } = props.rowData;

    model.value = { name, phone, address, source, status, remark };
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
        ? await createCustomer(model.value)
        : await updateCustomer(props.rowData?.id ?? 0, model.value);

    if (!success) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    window.$message?.success(props.operateType === 'add' ? '新增客户成功' : '客户信息已更新');

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
        label-width="100"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="客户姓名" path="name">
          <NInput v-model:value="model.name" placeholder="请输入客户姓名" maxlength="20" clearable />
        </NFormItem>
        <NFormItem label="手机号" path="phone">
          <NInput v-model:value="model.phone" placeholder="请输入 11 位手机号" maxlength="11" clearable />
        </NFormItem>
        <NFormItem label="房屋地址" path="address">
          <NInput v-model:value="model.address" placeholder="如：阳光美地 8栋2单元901" maxlength="60" clearable />
        </NFormItem>
        <NFormItem label="客户来源" path="source">
          <NSelect
            v-model:value="model.source"
            :options="CUSTOMER_SOURCE_OPTIONS"
            placeholder="请选择或输入客户来源"
            filterable
            tag
            clearable
          />
        </NFormItem>
        <NFormItem label="客户状态" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadioButton v-for="item in CUSTOMER_STATUS_OPTIONS" :key="item.value" :value="item.value">
              {{ item.label }}
            </NRadioButton>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput
            v-model:value="model.remark"
            type="textarea"
            :rows="4"
            maxlength="200"
            show-count
            placeholder="记录客户需求、预算、跟进要点"
          />
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
