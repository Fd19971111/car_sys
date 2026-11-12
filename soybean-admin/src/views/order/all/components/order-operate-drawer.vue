<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules, SelectOption } from 'naive-ui';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import {
  CUSTOMER_RESPONSE_OPTIONS,
  ORDER_STATUS_OPTIONS,
  VEHICLE_CATEGORY_OPTIONS,
  formatOrderDateTime
} from '../../components/order-meta';
import { useOrderData } from '../../components/use-order-data';

defineOptions({ name: 'OrderOperateDrawer' });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.Order.Item | null;
  customerOptions: SelectOption[];
  driverOptions: SelectOption[];
  vehicleOptions: SelectOption[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'submitted'): void;
}>();

const visible = defineModel<boolean>('visible', { default: false });
const appStore = useAppStore();
const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 620));
const { formRef, validate, restoreValidation } = useNaiveForm();
const { createOrder, updateOrder } = useOrderData();
const title = computed(() => (props.operateType === 'add' ? '新增订单' : '编辑订单'));
const submitting = ref(false);

interface FormModel {
  customerId: number | null;
  status: Api.Order.Status;
  passengerCount: number | null;
  vehicleCategory: Api.Order.VehicleCategory;
  pickupAddress: string;
  destinationAddress: string;
  estimatedDistanceKm: number | null;
  estimatedDurationMin: number | null;
  quotedFare: number | null;
  finalFare: number | null;
  driverFare: number | null;
  customerResponse: Api.Order.CustomerResponse;
  scheduledAt: string | null;
  driverId: number | null;
  vehicleId: number | null;
  completedAt: string | null;
  internalNotes: string;
}

function createDefaultModel(): FormModel {
  return {
    customerId: null,
    status: 'quoting',
    passengerCount: 1,
    vehicleCategory: 'comfort',
    pickupAddress: '',
    destinationAddress: '',
    estimatedDistanceKm: 0,
    estimatedDurationMin: 0,
    quotedFare: 0,
    finalFare: 0,
    driverFare: 0,
    customerResponse: 'pending',
    scheduledAt: null,
    driverId: null,
    vehicleId: null,
    completedAt: null,
    internalNotes: ''
  };
}

const model = ref<FormModel>(createDefaultModel());

const rules: FormRules = {
  customerId: [{ required: true, type: 'number', message: '请选择客户', trigger: ['change', 'blur'] }],
  status: [{ required: true, message: '请选择订单状态', trigger: ['change', 'blur'] }],
  vehicleCategory: [{ required: true, message: '请选择车型', trigger: ['change', 'blur'] }],
  pickupAddress: [{ required: true, message: '请填写上车地点', trigger: ['input', 'blur'] }],
  destinationAddress: [{ required: true, message: '请填写目的地点', trigger: ['input', 'blur'] }],
  passengerCount: [
    {
      validator: (_rule, value: number | null) => Number(value ?? 0) >= 1 && Number(value ?? 0) <= 8,
      message: '乘车人数必须在 1-8 之间',
      trigger: ['change', 'blur']
    }
  ],
  scheduledAt: [{ required: true, message: '请选择预约用车时间', trigger: ['change', 'blur'] }],
  customerResponse: [{ required: true, message: '请选择客户反馈', trigger: ['change', 'blur'] }]
};

watch(visible, async value => {
  if (!value) return;
  if (props.operateType === 'edit' && props.rowData) {
    const {
      customerId,
      status,
      passengerCount,
      vehicleCategory,
      pickupAddress,
      destinationAddress,
      estimatedDistanceKm,
      estimatedDurationMin,
      quotedFare,
      finalFare,
      driverFare,
      customerResponse,
      scheduledAt,
      driverId,
      vehicleId,
      completedAt,
      internalNotes
    } = props.rowData;
    model.value = {
      customerId: customerId || null,
      status,
      passengerCount,
      vehicleCategory,
      pickupAddress,
      destinationAddress,
      estimatedDistanceKm,
      estimatedDurationMin,
      quotedFare,
      finalFare,
      driverFare,
      customerResponse,
      scheduledAt: scheduledAt || null,
      driverId: driverId || null,
      vehicleId: vehicleId || null,
      completedAt: completedAt || null,
      internalNotes
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

  const payload: Api.Order.FormData = {
    customerId: Number(model.value.customerId ?? 0),
    status: model.value.status,
    passengerCount: Number(model.value.passengerCount ?? 1),
    vehicleCategory: model.value.vehicleCategory,
    pickupAddress: model.value.pickupAddress.trim(),
    destinationAddress: model.value.destinationAddress.trim(),
    estimatedDistanceKm: Number(model.value.estimatedDistanceKm ?? 0),
    estimatedDurationMin: Number(model.value.estimatedDurationMin ?? 0),
    quotedFare: Number(model.value.quotedFare ?? 0),
    finalFare: Number(model.value.finalFare ?? 0),
    driverFare: Number(model.value.driverFare ?? 0),
    customerResponse: model.value.customerResponse,
    scheduledAt: model.value.scheduledAt ?? '',
    driverId: Number(model.value.driverId ?? 0),
    vehicleId: Number(model.value.vehicleId ?? 0),
    completedAt: model.value.completedAt || null,
    internalNotes: model.value.internalNotes.trim()
  };

  submitting.value = true;
  try {
    const success =
      props.operateType === 'add' ? await createOrder(payload) : await updateOrder(props.rowData?.id ?? 0, payload);
    if (!success) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }
    window.$message?.success(props.operateType === 'add' ? '新增订单成功' : '订单已更新');
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
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="104" require-mark-placement="right-hanging">
        <NFormItem label="客户" path="customerId">
          <NSelect v-model:value="model.customerId" :options="customerOptions" placeholder="请选择客户" filterable />
        </NFormItem>
        <div class="grid grid-cols-1 gap-x-4 sm:grid-cols-2">
          <NFormItem label="订单状态" path="status">
            <NSelect v-model:value="model.status" :options="ORDER_STATUS_OPTIONS" />
          </NFormItem>
          <NFormItem label="客户反馈" path="customerResponse">
            <NSelect v-model:value="model.customerResponse" :options="CUSTOMER_RESPONSE_OPTIONS" />
          </NFormItem>
          <NFormItem label="乘车人数" path="passengerCount">
            <NInputNumber v-model:value="model.passengerCount" class="w-full" :min="1" :max="8" :precision="0" />
          </NFormItem>
          <NFormItem label="车型" path="vehicleCategory">
            <NSelect v-model:value="model.vehicleCategory" :options="VEHICLE_CATEGORY_OPTIONS" />
          </NFormItem>
        </div>

        <NFormItem label="上车地点" path="pickupAddress">
          <NInput v-model:value="model.pickupAddress" maxlength="255" placeholder="请输入上车地点" clearable />
        </NFormItem>
        <NFormItem label="目的地点" path="destinationAddress">
          <NInput v-model:value="model.destinationAddress" maxlength="255" placeholder="请输入目的地点" clearable />
        </NFormItem>

        <div class="grid grid-cols-1 gap-x-4 sm:grid-cols-2">
          <NFormItem label="预计里程" path="estimatedDistanceKm">
            <NInputNumber v-model:value="model.estimatedDistanceKm" class="w-full" :min="0" :precision="2">
              <template #suffix>公里</template>
            </NInputNumber>
          </NFormItem>
          <NFormItem label="预计时长" path="estimatedDurationMin">
            <NInputNumber v-model:value="model.estimatedDurationMin" class="w-full" :min="0" :precision="0">
              <template #suffix>分钟</template>
            </NInputNumber>
          </NFormItem>
          <NFormItem label="客户报价" path="quotedFare">
            <NInputNumber v-model:value="model.quotedFare" class="w-full" :min="0" :precision="2">
              <template #prefix>¥</template>
            </NInputNumber>
          </NFormItem>
          <NFormItem label="最终费用" path="finalFare">
            <NInputNumber v-model:value="model.finalFare" class="w-full" :min="0" :precision="2">
              <template #prefix>¥</template>
            </NInputNumber>
          </NFormItem>
          <NFormItem label="司机应得" path="driverFare">
            <NInputNumber v-model:value="model.driverFare" class="w-full" :min="0" :max="model.finalFare ?? undefined" :precision="2">
              <template #prefix>¥</template>
            </NInputNumber>
          </NFormItem>
        </div>

        <NFormItem label="预约用车" path="scheduledAt">
          <NDatePicker
            v-model:formatted-value="model.scheduledAt"
            class="w-full"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            placeholder="请选择预约用车时间"
            clearable
          />
        </NFormItem>

        <div class="grid grid-cols-1 gap-x-4 sm:grid-cols-2">
          <NFormItem label="指派司机" path="driverId">
            <NSelect v-model:value="model.driverId" :options="driverOptions" placeholder="可不指派" filterable clearable />
          </NFormItem>
          <NFormItem label="指派车辆" path="vehicleId">
            <NSelect v-model:value="model.vehicleId" :options="vehicleOptions" placeholder="可不指派" filterable clearable />
          </NFormItem>
        </div>

        <NFormItem v-if="model.status === 'completed'" label="结束时间" path="completedAt">
          <NDatePicker
            v-model:formatted-value="model.completedAt"
            class="w-full"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            placeholder="留空则保存时自动记录"
            clearable
          />
        </NFormItem>

        <NFormItem label="内部备注" path="internalNotes">
          <NInput v-model:value="model.internalNotes" type="textarea" :rows="3" maxlength="500" show-count placeholder="仅内部可见" />
        </NFormItem>

        <div v-if="operateType === 'edit' && rowData" class="rounded-8px bg-#f7f8fa px-12px py-10px text-12px text-gray-500">
          创建时间：{{ formatOrderDateTime(rowData.createdAt) }} · 更新时间：{{ formatOrderDateTime(rowData.updatedAt) }}
        </div>
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
