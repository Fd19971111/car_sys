<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { useSupplierData } from './use-supplier-data';

defineOptions({ name: 'SupplierOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Supplier.Item | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{ (e: 'submitted'): void }>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 460));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createSupplier, updateSupplier } = useSupplierData();
const title = computed(() => (props.operateType === 'add' ? '新增供应商' : '编辑供应商'));
const submitting = ref(false);

function createDefaultModel(): Api.Supplier.FormData {
  return { name: '', contactName: '', phone: '', address: '', remark: '' };
}

const model = ref<Api.Supplier.FormData>(createDefaultModel());
const rules: FormRules = {
  name: [{ required: true, message: '请输入供应商名称', trigger: ['input', 'blur'] }]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const { name, contactName, phone, address, remark } = props.rowData;
    model.value = { name, contactName, phone, address, remark };
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
    const success = props.operateType === 'add'
      ? await createSupplier(model.value)
      : await updateSupplier(props.rowData?.id ?? 0, model.value);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增供应商成功' : '供应商信息已更新');
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
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="104">
        <NFormItem label="供应商名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" placeholder="请输入公司或供应商名称" clearable />
        </NFormItem>
        <NFormItem label="联系人" path="contactName">
          <NInput v-model:value="model.contactName" maxlength="50" placeholder="选填" clearable />
        </NFormItem>
        <NFormItem label="联系电话" path="phone">
          <NInput v-model:value="model.phone" maxlength="30" placeholder="选填" clearable />
        </NFormItem>
        <NFormItem label="地址" path="address">
          <NInput v-model:value="model.address" maxlength="255" placeholder="选填" clearable />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="4" maxlength="255" show-count />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="visible = false">取消</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSubmit">保存</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>