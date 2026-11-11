<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm, NSelect, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { fetchDepartments, fetchDepartmentUsers, fetchSetUserDepartment } from '@/service/api/menu';
import StaffOperateDrawer from './components/staff-operate-drawer.vue';
import { getStaffStatusMeta, formatDate } from './components/staff-meta';
import { useStaffData } from './components/use-staff-data';

defineOptions({ name: 'EnterpriseStaff' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
  departmentId: number | null;
  status: Api.Staff.Status | null;
}

interface PaginationState {
  page: number;
  pageSize: number;
  itemCount: number;
  showSizePicker: boolean;
  pageSizes: number[];
  prefix: (info: { itemCount?: number }) => string;
  onUpdatePage: (page: number) => void;
  onUpdatePageSize: (pageSize: number) => void;
}

const { loading, records, total, departmentOptions, loadList, deleteStaff, loadDepartmentOptions } = useStaffData();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  departmentId: null,
  status: null
});

const checkedRowKeys = ref<DataTableRowKey[]>([]);
const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Staff.Item | null>(null);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 人`,
  onUpdatePage(page) {
    pagination.page = page;
    void getData();
  },
  onUpdatePageSize(pageSize) {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    void getData();
  }
});

// ---- 人员归属相关 ----
const departments = ref<Api.Menu.Department[]>([]);
const users = ref<Api.Menu.DepartmentUser[]>([]);
const userLoading = ref(false);

const columns: DataTableColumns<Api.Staff.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'name',
    title: '人员',
    width: 150,
    fixed: 'left',
    render: row => h('div', { class: 'flex-y-center gap-8px' }, [
      h('div', {
        class: 'h-28px w-28px flex-center flex-shrink-0 rounded-full text-12px font-medium',
        style: { backgroundColor: 'rgb(var(--primary-color) / 0.12)', color: 'rgb(var(--primary-color))' }
      }, row.name.slice(0, 1)),
      h('div', { class: 'flex-col' }, [
        h('span', { class: 'font-medium' }, row.name),
        h('span', { class: 'text-12px text-gray-400' }, row.position || '未设置职位')
      ])
    ])
  },
  { key: 'phone', title: '手机号', width: 130, render: row => h('span', { class: 'font-mono' }, row.phone) },
  { key: 'gender', title: '性别', width: 70, render: row => row.gender || '-' },
  {
    key: 'departmentName',
    title: '所属部门',
    minWidth: 140,
    ellipsis: { tooltip: true },
    render: row => row.departmentName || '-'
  },
  { key: 'joinDate', title: '入职日期', width: 120, render: row => formatDate(row.joinDate) },
  {
    key: 'status',
    title: '状态',
    width: 90,
    render: row => {
      const meta = getStaffStatusMeta(row.status);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  {
    key: 'operate',
    title: '操作',
    width: 140,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-4px' }, [
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
          trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
          default: () => '确认删除该人员吗？'
        })
      ])
  }
];

const userColumns: DataTableColumns<Api.Menu.DepartmentUser> = [
  { title: '人员', key: 'name', minWidth: 150, render: row => row.nickname || row.username },
  { title: '账号', key: 'username', minWidth: 150 },
  {
    title: '角色',
    key: 'role',
    width: 160,
    render: row => (row.role === 1 ? '公司超级管理员' : '部门成员')
  },
  {
    title: '所属部门',
    key: 'departmentId',
    minWidth: 220,
    render: row =>
      h(NSelect, {
        value: row.departmentId,
        options: [{ label: '未分配', value: 0 }, ...departments.value.map(item => ({ label: item.name, value: item.id }))],
        disabled: row.role === 1,
        onUpdateValue: (value: number) => void updateUserDepartment(row, value)
      })
  }
];

async function getData() {
  searchParams.current = pagination.page;
  searchParams.size = pagination.pageSize;
  await loadList(searchParams);
  pagination.itemCount = total.value;
}

async function handleSearch() {
  pagination.page = 1;
  await getData();
}

async function handleReset() {
  searchParams.keyword = '';
  searchParams.departmentId = null;
  searchParams.status = null;
  pagination.page = 1;
  await getData();
}

async function handleRefresh() {
  await getData();
}

function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

function handleEdit(row: Api.Staff.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

async function handleDelete(row: Api.Staff.Item) {
  const success = await deleteStaff([row.id]);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  await afterDeleted(1);
}

async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(key => Number(key));
  const success = await deleteStaff(ids);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  checkedRowKeys.value = [];
  await afterDeleted(ids.length);
}

async function afterDeleted(deletedCount: number) {
  const maxPage = Math.max(1, Math.ceil((total.value - deletedCount) / pagination.pageSize));
  if (pagination.page > maxPage) pagination.page = maxPage;
  await getData();
}

async function handleSubmitted() {
  if (operateType.value === 'add') pagination.page = 1;
  await getData();
}

// ---- 人员归属方法 ----
async function loadUserData() {
  userLoading.value = true;
  try {
    const [deptResult, userResult] = await Promise.all([fetchDepartments(), fetchDepartmentUsers()]);
    if (!deptResult.error && deptResult.data) departments.value = deptResult.data.items;
    if (!userResult.error && userResult.data) users.value = userResult.data.items;
  } finally {
    userLoading.value = false;
  }
}

async function updateUserDepartment(user: Api.Menu.DepartmentUser, departmentId: number) {
  const { error } = await fetchSetUserDepartment(user.id, departmentId);
  if (!error) {
    user.departmentId = departmentId;
    window.$message?.success('人员所属部门已更新');
  }
}

onMounted(async () => {
  await Promise.all([getData(), loadDepartmentOptions()]);
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <NCard :bordered="false" class="card-wrapper" size="small">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-16px font-semibold">人员管理</h1>
            <p class="mt-4px text-13px text-gray-500">人员档案与部门归属</p>
          </div>
        </div>
      </NCard>

      <NCard :bordered="false" class="card-wrapper" size="small">
        <NTabs type="line" animated @update:value="(val: string) => { if (val === 'users') loadUserData(); }">
          <NTabPane name="staff" tab="人员档案">
            <NSpace vertical :size="12">
              <div class="flex flex-wrap items-center gap-8px">
                <NInput
                  v-model:value="searchParams.keyword"
                  class="w-200px min-w-140px flex-shrink"
                  placeholder="搜索姓名 / 手机号"
                  clearable
                  @keyup.enter="handleSearch"
                >
                  <template #prefix><SvgIcon icon="mdi:magnify" class="text-gray-400" /></template>
                </NInput>
                <NSelect
                  v-model:value="searchParams.departmentId"
                  class="w-160px flex-shrink-0"
                  :options="departmentOptions"
                  placeholder="所属部门"
                  clearable
                  filterable
                />
                <NSelect
                  v-model:value="searchParams.status"
                  class="w-100px flex-shrink-0"
                  :options="[{ value: 1, label: '在职' }, { value: 2, label: '离职' }]"
                  placeholder="状态"
                  clearable
                />
                <NButton class="flex-shrink-0" type="primary" ghost size="small" @click="handleSearch">
                  <template #icon><SvgIcon icon="mdi:magnify" /></template>
                  查询
                </NButton>
                <NButton class="flex-shrink-0" size="small" @click="handleReset">
                  <template #icon><SvgIcon icon="mdi:backspace-outline" /></template>
                  重置
                </NButton>
                <div class="ml-auto flex flex-shrink-0 items-center gap-8px">
                  <NButton size="small" type="primary" ghost @click="handleAdd">
                    <template #icon><SvgIcon icon="mdi:account-plus-outline" /></template>
                    新增人员
                  </NButton>
                  <NPopconfirm @positive-click="handleBatchDelete">
                    <template #trigger>
                      <NButton size="small" ghost type="error" :disabled="!checkedRowKeys.length">
                        <template #icon><SvgIcon icon="mdi:delete-outline" /></template>
                        批量删除
                      </NButton>
                    </template>
                    确认删除选中的 {{ checkedRowKeys.length }} 人吗？
                  </NPopconfirm>
                  <NButton size="small" @click="handleRefresh">
                    <template #icon><SvgIcon icon="mdi:refresh" :class="{ 'animate-spin': loading }" /></template>
                    刷新
                  </NButton>
                </div>
              </div>
              <NDataTable
                v-model:checked-row-keys="checkedRowKeys"
                :columns="columns"
                :data="records"
                :loading="loading"
                :pagination="pagination"
                :row-key="row => row.id"
                :scroll-x="1000"
                remote
                size="small"
                striped
              />
            </NSpace>
          </NTabPane>

          <NTabPane name="users" tab="人员归属">
            <NDataTable
              :columns="userColumns"
              :data="users"
              :loading="userLoading"
              :pagination="false"
              :row-key="row => row.id"
              :scroll-x="720"
            />
          </NTabPane>
        </NTabs>
      </NCard>
    </NSpace>

    <StaffOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      @submitted="handleSubmitted"
    />
  </div>
</template>

<style scoped></style>