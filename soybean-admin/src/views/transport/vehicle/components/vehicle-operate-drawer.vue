<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules, SelectOption } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { getDriverStatusMeta } from '../../driver/components/driver-meta';
import { VEHICLE_CATEGORY_OPTIONS, VEHICLE_STATUS_OPTIONS } from './vehicle-meta';
import { useVehicleData } from './use-vehicle-data';

defineOptions({ name: 'VehicleOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Transport.Vehicle | null;
  driverOptions: Api.Transport.DriverOption[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 500));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createVehicle, updateVehicle } = useVehicleData();
const title = computed(() => (props.operateType === 'add' ? '新增车辆' : '编辑车辆'));
const submitting = ref(false);

interface FormModel {
  plateNumber: string;
  brand: string;
  category: Api.Transport.VehicleCategory;
  seatCount: number | null;
  status: Api.Transport.VehicleStatus;
  driverId: number | null;
  remark: string;
}

function createDefaultModel(): FormModel {
  return {
    plateNumber: '',
    brand: '',
    category: 'comfort',
    seatCount: 5,
    status: 'active',
    driverId: null,
    remark: ''
  };
}

const model = ref<FormModel>(createDefaultModel());

const driverSelectOptions = computed<SelectOption[]>(() =>
  props.driverOptions.map(item => ({
    label: `${item.name} · ${item.phone}（${getDriverStatusMeta(item.serviceStatus).label}）`,
    value: item.id
  }))
);

const rules: FormRules = {
  plateNumber: [{ required: true, message: '请输入车牌号', trigger: ['input', 'blur'] }],
  category: [{ required: true, message: '请选择服务车型', trigger: ['change', 'blur'] }],
  seatCount: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 1 && Number(value ?? 0) <= 8,
      message: '可乘人数必须在 1-8 之间',
      trigger: ['change', 'blur']
    }
  ],
  status: [{ required: true, message: '请选择车辆状态', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const { plateNumber, brand, category, seatCount, status, driverId, remark } = props.rowData;
    model.value = {
      plateNumber,
      brand,
      category,
      seatCount,
      status,
      driverId: driverId || null,
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

  const payload: Api.Transport.VehicleFormData = {
    plateNumber: model.value.plateNumber.trim(),
    brand: model.value.brand.trim(),
    category: model.value.category,
    seatCount: Number(model.value.seatCount ?? 0),
    status: model.value.status,
    driverId: model.value.driverId ?? 0,
    remark: model.value.remark.trim()
  };

  submitting.value = true;
  try {
    const success =
      props.operateType === 'add' ? await createVehicle(payload) : await updateVehicle(props.rowData?.id ?? 0, payload);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增车辆成功' : '车辆信息已更新');
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
        <NFormItem label="车牌号" path="plateNumber">
          <NInput v-model:value="model.plateNumber" maxlength="20" placeholder="如：贵A·12345" clearable />
        </NFormItem>
        <NFormItem label="品牌车型" path="brand">
          <NInput v-model:value="model.brand" maxlength="100" placeholder="如：比亚迪 汉 EV" clearable />
        </NFormItem>
        <NFormItem label="服务车型" path="category">
          <NSelect v-model:value="model.category" :options="VEHICLE_CATEGORY_OPTIONS" placeholder="请选择服务车型" />
        </NFormItem>
        <NFormItem label="可乘人数" path="seatCount">
          <NInputNumber v-model:value="model.seatCount" class="w-full" :min="1" :max="8" :precision="0" />
        </NFormItem>
        <NFormItem label="车辆状态" path="status">
          <NSelect v-model:value="model.status" :options="VEHICLE_STATUS_OPTIONS" placeholder="请选择车辆状态" />
        </NFormItem>
        <NFormItem label="当前司机" path="driverId">
          <NSelect
            v-model:value="model.driverId"
            :options="driverSelectOptions"
            placeholder="可不绑定，稍后在车辆管理中指派"
            filterable
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
