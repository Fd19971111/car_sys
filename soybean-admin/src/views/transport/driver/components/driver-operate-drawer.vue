<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import { REG_PHONE } from '@/constants/reg';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { DRIVER_LICENSE_TYPE_OPTIONS, DRIVER_STATUS_OPTIONS } from './driver-meta';
import { useDriverData } from './use-driver-data';

defineOptions({ name: 'DriverOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Transport.Driver | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 500));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createDriver, updateDriver } = useDriverData();
const title = computed(() => (props.operateType === 'add' ? '新增司机' : '编辑司机'));
const submitting = ref(false);

interface FormModel {
  name: string;
  phone: string;
  licenseNumber: string;
  licenseType: string | null;
  serviceStatus: Api.Transport.DriverServiceStatus;
  rating: number | null;
  completedTrips: number | null;
  joinedAt: string | null;
  remark: string;
}

function createDefaultModel(): FormModel {
  return {
    name: '',
    phone: '',
    licenseNumber: '',
    licenseType: 'C1',
    serviceStatus: 'available',
    rating: 5,
    completedTrips: 0,
    joinedAt: null,
    remark: ''
  };
}

const model = ref<FormModel>(createDefaultModel());

const rules: FormRules = {
  name: [{ required: true, message: '请输入司机姓名', trigger: ['input', 'blur'] }],
  phone: [
    { required: true, message: '请输入手机号', trigger: ['input', 'blur'] },
    { pattern: REG_PHONE, message: '手机号格式不正确', trigger: ['input', 'blur'] }
  ],
  serviceStatus: [{ required: true, message: '请选择接单状态', trigger: ['change', 'blur'] }],
  rating: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 0 && Number(value ?? 0) <= 5,
      message: '服务评分必须在 0-5 之间',
      trigger: ['change', 'blur']
    }
  ]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const { name, phone, licenseNumber, licenseType, serviceStatus, rating, completedTrips, joinedAt, remark } = props.rowData;
    model.value = {
      name,
      phone,
      licenseNumber,
      licenseType: licenseType || null,
      serviceStatus,
      rating,
      completedTrips,
      joinedAt: joinedAt || null,
      remark
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
    window.$message?.warning('请检查表单填写是否正确');
    return;
  }

  const payload: Api.Transport.DriverFormData = {
    name: model.value.name.trim(),
    phone: model.value.phone.trim(),
    licenseNumber: model.value.licenseNumber.trim(),
    licenseType: model.value.licenseType ?? '',
    serviceStatus: model.value.serviceStatus,
    rating: Number(model.value.rating ?? 0),
    completedTrips: Number(model.value.completedTrips ?? 0),
    joinedAt: model.value.joinedAt || null,
    remark: model.value.remark.trim()
  };

  submitting.value = true;
  try {
    const success =
      props.operateType === 'add' ? await createDriver(payload) : await updateDriver(props.rowData?.id ?? 0, payload);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增司机成功' : '司机信息已更新');
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
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="96" require-mark-placement="right-hanging">
        <NFormItem label="司机姓名" path="name">
          <NInput v-model:value="model.name" maxlength="50" placeholder="请输入司机姓名" clearable />
        </NFormItem>
        <NFormItem label="手机号" path="phone">
          <NInput v-model:value="model.phone" maxlength="11" placeholder="请输入 11 位手机号" clearable />
        </NFormItem>
        <NFormItem label="驾驶证号" path="licenseNumber">
          <NInput v-model:value="model.licenseNumber" maxlength="50" placeholder="请输入驾驶证号" clearable />
        </NFormItem>
        <NFormItem label="准驾车型" path="licenseType">
          <NSelect v-model:value="model.licenseType" :options="DRIVER_LICENSE_TYPE_OPTIONS" placeholder="请选择准驾车型" clearable />
        </NFormItem>
        <NFormItem label="接单状态" path="serviceStatus">
          <NSelect v-model:value="model.serviceStatus" :options="DRIVER_STATUS_OPTIONS" placeholder="请选择接单状态" />
        </NFormItem>
        <NFormItem label="服务评分" path="rating">
          <NInputNumber v-model:value="model.rating" class="w-full" :min="0" :max="5" :step="0.1" :precision="1" />
        </NFormItem>
        <NFormItem label="完成订单" path="completedTrips">
          <NInputNumber v-model:value="model.completedTrips" class="w-full" :min="0" :precision="0" />
        </NFormItem>
        <NFormItem label="加入日期" path="joinedAt">
          <NDatePicker
            v-model:formatted-value="model.joinedAt"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="请选择加入日期"
            clearable
          />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput v-model:value="model.remark" type="textarea" :rows="3" maxlength="255" show-count placeholder="补充说明" />
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
