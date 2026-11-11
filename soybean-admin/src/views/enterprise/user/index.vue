<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NTag, NSwitch, NPopconfirm, NInput, NSelect, NModal, NForm, NFormItem } from 'naive-ui';
import {
  fetchCreateUserAccount,
  fetchChangeUserPwd,
  fetchGetAvailableStaff,
  fetchGetUserList,
  fetchResetUserPwd,
  fetchToggleUser,
  fetchUpdateUserAccount
} from '@/service/api/user';

defineOptions({ name: 'EnterpriseUser' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
  status: number | null;
  role: number | null;
}

type CreateAccountFormModel = Omit<Api.User.CreateAccountForm, 'staffId'> & { staffId: number | null };

const loading = ref(false);
const records = ref<Api.User.Item[]>([]);
const total = ref(0);
const modalVisible = ref(false);
const form = reactive<CreateAccountFormModel>({ staffId: null, username: '', password: '' });
const formSubmitting = ref(false);
const staffOptions = ref<Api.User.StaffOption[]>([]);
const editModalVisible = ref(false);
const editForm = reactive<Api.User.UpdateAccountForm>({ id: 0, username: '', role: 2 });
const editSubmitting = ref(false);
const changePwdModalVisible = ref(false);
const changePwdSubmitting = ref(false);
const changePwdUserId = ref(0);
const changePwdUsername = ref('');
const newPassword = ref('');
const confirmPassword = ref('');

const searchParams = reactive<SearchParams>({ current: 1, size: 10, keyword: '', status: null, role: null });

const pagination = reactive({
  page: 1, pageSize: 10, itemCount: 0, showSizePicker: true, pageSizes: [10, 20, 50],
  prefix: (info: { itemCount?: number }) => `共 ${info.itemCount ?? 0} 个账号`,
  onUpdatePage(page: number) { pagination.page = page; void getData(); },
  onUpdatePageSize(size: number) { pagination.pageSize = size; pagination.page = 1; void getData(); }
});

const columns: DataTableColumns<Api.User.Item> = [
  {
    title: '用户', key: 'nickname', minWidth: 160,
    render: row => h('div', { class: 'flex-y-center gap-8px' }, [
      h('div', { class: 'h-28px w-28px flex-center rounded-full text-12px font-medium', style: { backgroundColor: 'rgb(var(--primary-color)/0.12)', color: 'rgb(var(--primary-color))' } }, (row.nickname || row.username).slice(0, 1)),
      h('div', { class: 'flex-col' }, [h('span', { class: 'font-medium' }, row.nickname || row.username), h('span', { class: 'text-12px text-gray-400' }, row.username)])
    ])
  },
  { title: '手机号', key: 'phone', width: 130 },
  { title: '所属部门', key: 'departmentName', minWidth: 130, render: row => row.departmentName || '-' },
  {
    title: '角色', key: 'role', width: 90,
    render: row => h(NTag, { size: 'small', bordered: false, type: row.role === 1 ? 'primary' : 'default' }, { default: () => row.roleName })
  },
  {
    title: '状态', key: 'status', width: 80,
    render: row => h(NSwitch, { value: row.status === 1, size: 'small', onUpdateValue: (v: boolean) => void toggleStatus(row, v) })
  },
  { title: '最后登录', key: 'lastLoginAt', width: 140, render: row => row.lastLoginAt || '-' },
  {
    title: '操作', key: 'actions', width: 280, align: 'center', fixed: 'right',
    render: row => h('div', { class: 'flex-center gap-2' }, [
      h(NButton, { size: 'tiny', quaternary: true, onClick: () => openChangePwdModal(row) }, { default: () => '修改密码' }),
      h(NButton, { size: 'tiny', quaternary: true, onClick: () => void resetPwd(row) }, { default: () => '重置密码' }),
      h(NButton, { size: 'tiny', quaternary: true, onClick: () => openEditModal(row) }, { default: () => '编辑' }),
      h(NPopconfirm, { onPositiveClick: () => void toggleStatus(row, row.status !== 1) }, {
        trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: row.status === 1 ? 'error' : 'success' }, { default: () => row.status === 1 ? '停用' : '启用' }),
        default: () => `确认${row.status === 1 ? '停用' : '启用'}该账号？`
      })
    ])
  }
];

async function getData() {
  loading.value = true;
  try {
    const { data } = await fetchGetUserList({ current: pagination.page, size: pagination.pageSize, keyword: searchParams.keyword || null, status: searchParams.status, role: searchParams.role });
    if (data) { records.value = data.records; total.value = data.total; pagination.itemCount = data.total; }
  } finally { loading.value = false; }
}

async function resetPwd(row: Api.User.Item) {
  const { data } = await fetchResetUserPwd(row.id);
  if (data) window.$message?.success(`密码已重置为: ${data.password}`);
}

async function toggleStatus(row: Api.User.Item, toEnable: boolean) {
  const { error } = await fetchToggleUser(row.id, toEnable ? 1 : 0);
  if (!error) {
    row.status = toEnable ? 1 : 0;
    window.$message?.success(toEnable ? '账号已启用' : '账号已停用');
  }
}

async function openCreateModal() {
  const { data } = await fetchGetAvailableStaff();
  staffOptions.value = data?.items ?? [];
  form.staffId = null;
  form.username = '';
  form.password = '';
  modalVisible.value = true;
}

function openEditModal(row: Api.User.Item) {
  editForm.id = row.id;
  editForm.username = row.username;
  editForm.role = row.role;
  editModalVisible.value = true;
}

function openChangePwdModal(row: Api.User.Item) {
  changePwdUserId.value = row.id;
  changePwdUsername.value = row.nickname || row.username;
  newPassword.value = '';
  confirmPassword.value = '';
  changePwdModalVisible.value = true;
}

async function updateAccount() {
  const username = editForm.username.trim();
  if (username.length < 4 || username.length > 20 || ![1, 2].includes(editForm.role)) {
    window.$message?.warning('请检查用户名和角色');
    return;
  }
  editSubmitting.value = true;
  try {
    const { error } = await fetchUpdateUserAccount({ ...editForm, username });
    if (!error) {
      window.$message?.success('账号信息已更新');
      editModalVisible.value = false;
      await getData();
    }
  } finally {
    editSubmitting.value = false;
  }
}

async function createAccount() {
  const staffId = form.staffId;
  if (typeof staffId !== 'number' || staffId <= 0 || !form.username || !form.password) {
    window.$message?.warning('请完善所有字段');
    return;
  }
  formSubmitting.value = true;
  try {
    const { error } = await fetchCreateUserAccount({ ...form, staffId });
    if (!error) {
      window.$message?.success('账号创建成功');
      modalVisible.value = false;
      await getData();
    }
  } finally { formSubmitting.value = false; }
}

async function changePwd() {
  if (newPassword.value.length < 6 || newPassword.value.length > 20) {
    window.$message?.warning('密码长度需为 6-20 个字符');
    return;
  }
  if (newPassword.value !== confirmPassword.value) {
    window.$message?.warning('两次输入的密码不一致');
    return;
  }
  changePwdSubmitting.value = true;
  try {
    const { error } = await fetchChangeUserPwd(changePwdUserId.value, newPassword.value);
    if (!error) {
      window.$message?.success('密码修改成功');
      changePwdModalVisible.value = false;
    }
  } finally {
    changePwdSubmitting.value = false;
  }
}

onMounted(() => void getData());
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <NCard :bordered="false" class="card-wrapper" size="small">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-16px font-semibold">用户管理</h1>
            <p class="mt-4px text-13px text-gray-500">登录账号与权限管理</p>
          </div>
          <NButton type="primary" size="small" @click="openCreateModal">
            <template #icon><SvgIcon icon="mdi:account-key-outline" /></template>
            创建账号
          </NButton>
        </div>
      </NCard>

      <NCard :bordered="false" class="card-wrapper" size="small">
        <NSpace vertical :size="12">
          <div class="flex flex-wrap items-center gap-8px">
            <NInput v-model:value="searchParams.keyword" class="w-200px" placeholder="搜索姓名/用户名/手机号" clearable @keyup.enter="getData()">
              <template #prefix><SvgIcon icon="mdi:magnify" class="text-gray-400" /></template>
            </NInput>
            <NSelect v-model:value="searchParams.role" class="w-120px" :options="[{ value: 1, label: '管理员' }, { value: 2, label: '员工' }]" placeholder="角色" clearable />
            <NSelect v-model:value="searchParams.status" class="w-100px" :options="[{ value: 1, label: '正常' }, { value: 0, label: '禁用' }]" placeholder="状态" clearable />
            <NButton type="primary" ghost size="small" @click="pagination.page = 1; getData()">查询</NButton>
            <NButton size="small" @click="searchParams.keyword = ''; searchParams.status = null; searchParams.role = null; pagination.page = 1; getData()">重置</NButton>
          </div>
          <NDataTable :columns="columns" :data="records" :loading="loading" :pagination="pagination" :row-key="row => row.id" :scroll-x="900" remote size="small" striped />
        </NSpace>
      </NCard>
    </NSpace>

    <NModal v-model:show="modalVisible" preset="card" title="创建登录账号" class="w-[min(440px,calc(100vw-32px))]">
      <NForm :model="form" label-placement="left" label-width="96">
        <NFormItem label="选择员工" required>
          <NSelect v-model:value="form.staffId" :options="staffOptions.map(item => ({ label: `${item.name}（${item.phone}）`, value: item.id }))" placeholder="选择待开通员工" filterable />
        </NFormItem>
        <NFormItem label="登录用户名" required>
          <NInput v-model:value="form.username" placeholder="4-20 个字符" maxlength="20" />
        </NFormItem>
        <NFormItem label="登录密码" required>
          <NInput v-model:value="form.password" type="password" placeholder="6-20 个字符" maxlength="20" show-password-on="click" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="modalVisible = false">取消</NButton>
          <NButton type="primary" :loading="formSubmitting" @click="createAccount">创建</NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal v-model:show="changePwdModalVisible" preset="card" :title="`修改密码 - ${changePwdUsername}`" class="w-[min(440px,calc(100vw-32px))]">
      <NForm label-placement="left" label-width="96">
        <NFormItem label="新密码" required>
          <NInput v-model:value="newPassword" type="password" placeholder="6-20 个字符" maxlength="20" show-password-on="click" />
        </NFormItem>
        <NFormItem label="确认密码" required>
          <NInput v-model:value="confirmPassword" type="password" placeholder="再次输入新密码" maxlength="20" show-password-on="click" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="changePwdModalVisible = false">取消</NButton>
          <NButton type="primary" :loading="changePwdSubmitting" @click="changePwd">确认修改</NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal v-model:show="editModalVisible" preset="card" title="编辑账号" class="w-[min(440px,calc(100vw-32px))]">
      <NForm :model="editForm" label-placement="left" label-width="96">
        <NFormItem label="登录用户名" required>
          <NInput v-model:value="editForm.username" placeholder="4-20 个字符" maxlength="20" />
        </NFormItem>
        <NFormItem label="角色" required>
          <NSelect v-model:value="editForm.role" :options="[{ value: 1, label: '管理员' }, { value: 2, label: '员工' }]" placeholder="请选择角色" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="editModalVisible = false">取消</NButton>
          <NButton type="primary" :loading="editSubmitting" @click="updateAccount">保存</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped></style>