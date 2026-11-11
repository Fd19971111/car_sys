<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { useWarehouseData } from './use-warehouse-data';

defineOptions({ name: 'WarehouseOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Inventory.Warehouse | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{ (e: 'submitted'): void }>();
const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 460));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createWarehouse, updateWarehouse } = useWarehouseData();
const title = computed(() => (props.operateType === 'add' ? '新增仓库' : '编辑仓库'));
const submitting = ref(false);

function createDefaultModel(): Api.Inventory.WarehouseForm {
  return { name: '', code: '', address: '', remark: '', isDefault: false, sort: 0 };
}

const model = ref<Api.Inventory.WarehouseForm>(createDefaultModel());
const rules: FormRules = {
  name: [{ required: true, message: '请输入仓库名称', trigger: ['input', 'blur'] }]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const { name, code, address, remark, isDefault, sort } = props.rowData;
    model.value = { name, code, address, remark, isDefault, sort };
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
    const form = { ...model.value, sort: Number(model.value.sort ?? 0) };
    const success = props.operateType === 'add'
      ? await createWarehouse(form)
      : await updateWarehouse(props.rowData?.id ?? 0, form);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增仓库成功' : '仓库信息已更新');
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
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="92">
        <NFormItem label="仓库名称" path="name">
          <NInput v-model:value="model.name" maxlength="100" placeholder="如：默认仓库 / 一号仓库" clearable />
        </NFormItem>
        <NFormItem label="仓库编号" path="code">
          <NInput v-model:value="model.code" maxlength="50" placeholder="选填，如：WH-002" clearable />
        </NFormItem>
        <NFormItem label="仓库地址" path="address">
          <NInput v-model:value="model.address" maxlength="255" placeholder="选填" clearable />
        </NFormItem>
        <NFormItem label="排序" path="sort">
          <NInputNumber v-model:value="model.sort" class="w-full" :min="0" :max="9999" :precision="0" placeholder="数字越小越靠前" />
        </NFormItem>
        <NFormItem label="设为默认" path="isDefault">
          <div class="flex-col">
            <NSwitch v-model:value="model.isDefault" />
            <span class="mt-4px text-12px text-gray-400">新入库材料的归属仓库；同一时间只能有一个默认仓库</span>
          </div>
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count />
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
