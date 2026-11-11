<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import type { FormRules } from 'naive-ui';
import dayjs from 'dayjs';
import { useNaiveForm } from '@/hooks/common/form';
import { useAppStore } from '@/store/modules/app';
import { PROJECT_STATUS_OPTIONS } from './project-meta';
import { useProjectData } from './use-project-data';

defineOptions({
  name: 'ProjectOperateDrawer'
});

interface Props {
  /** 操作类型：add=新增 edit=编辑 */
  operateType: NaiveUI.TableOperateType;
  /** 编辑时的项目数据 */
  rowData?: Api.Project.Item | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 520));

const { formRef, validate, restoreValidation } = useNaiveForm();

const { options, loadOptions, createProject, updateProject } = useProjectData();

const title = computed(() => (props.operateType === 'add' ? '新增项目' : '编辑项目'));

const submitting = ref(false);

/** 客户下拉选项 */
const customerOptions = computed(() => options.value.customers);

/** 负责人下拉选项 */
const managerOptions = computed(() => options.value.managers);

function createDefaultModel(): Api.Project.FormData {
  return {
    customerId: null,
    name: '',
    address: '',
    managerId: null,
    status: 1,
    amount: null,
    startDate: null,
    endDate: null,
    remark: ''
  };
}

const model = ref<Api.Project.FormData>(createDefaultModel());

const rules: FormRules = {
  customerId: [{ required: true, type: 'number', message: '请选择所属客户', trigger: ['change', 'blur'] }],
  name: [
    { required: true, message: '请输入项目名称', trigger: ['input', 'blur'] },
    { min: 2, max: 50, message: '项目名称长度为 2-50 个字符', trigger: ['input', 'blur'] }
  ],
  status: [{ required: true, type: 'number', message: '请选择项目状态', trigger: ['change', 'blur'] }],
  endDate: [
    {
      validator: (_rule, value: string | null) => {
        if (!value || !model.value.startDate) return true;

        return !dayjs(value).isBefore(dayjs(model.value.startDate), 'day');
      },
      message: '预计完工日期不能早于开工日期',
      trigger: ['change', 'blur']
    }
  ]
};

/** 预计完工日期可选范围：不早于开工日期 */
function disabledEndDate(timestamp: number) {
  if (!model.value.startDate) return false;

  return timestamp < dayjs(model.value.startDate).startOf('day').valueOf();
}

watch(visible, async value => {
  if (!value) {
    return;
  }

  if (props.operateType === 'edit' && props.rowData) {
    const { customerId, name, address, managerId, status, amount, startDate, endDate, remark } = props.rowData;

    model.value = {
      customerId: customerId || null,
      name,
      address,
      managerId: managerId || null,
      status,
      amount: amount ?? null,
      startDate: startDate || null,
      endDate: endDate || null,
      remark
    };
  } else {
    model.value = createDefaultModel();
  }

  await nextTick();

  restoreValidation();

  // 下拉选项可能尚未加载（例如直接进入本页后立刻点新增），这里兜底拉一次
  if (!customerOptions.value.length && !managerOptions.value.length) {
    await loadOptions();
  }
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
        ? await createProject(model.value)
        : await updateProject(props.rowData?.id ?? 0, model.value);

    if (!success) {
      window.$message?.error('保存失败，请稍后重试');

      return;
    }

    window.$message?.success(props.operateType === 'add' ? '新增项目成功' : '项目信息已更新');

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
        <NFormItem label="所属客户" path="customerId">
          <NSelect
            v-model:value="model.customerId"
            :options="customerOptions"
            placeholder="请选择所属客户（成交客户）"
            filterable
            clearable
          />
        </NFormItem>
        <NFormItem label="项目名称" path="name">
          <NInput v-model:value="model.name" placeholder="如：阳光美地 8-901 全屋整装" maxlength="50" clearable />
        </NFormItem>
        <NFormItem label="施工地址" path="address">
          <NInput v-model:value="model.address" placeholder="如：阳光美地 8栋2单元901" maxlength="60" clearable />
        </NFormItem>
        <NFormItem label="项目负责人" path="managerId">
          <NSelect
            v-model:value="model.managerId"
            :options="managerOptions"
            placeholder="请选择项目负责人"
            filterable
            clearable
          />
        </NFormItem>
        <NFormItem label="项目状态" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadioButton v-for="item in PROJECT_STATUS_OPTIONS" :key="item.value" :value="item.value">
              {{ item.label }}
            </NRadioButton>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="合同金额" path="amount">
          <NInputNumber
            v-model:value="model.amount"
            class="w-full"
            :min="0"
            :precision="2"
            :show-button="false"
            placeholder="请输入合同金额（元）"
          >
            <template #prefix>¥</template>
          </NInputNumber>
        </NFormItem>
        <NFormItem label="开工日期" path="startDate">
          <NDatePicker
            v-model:formatted-value="model.startDate"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="请选择开工日期"
            clearable
          />
        </NFormItem>
        <NFormItem label="预计完工" path="endDate">
          <NDatePicker
            v-model:formatted-value="model.endDate"
            class="w-full"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="请选择预计完工日期"
            :disabled-date="disabledEndDate"
            clearable
          />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput
            v-model:value="model.remark"
            type="textarea"
            :rows="4"
            maxlength="200"
            show-count
            placeholder="记录施工范围、特殊要求、付款节点等"
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
