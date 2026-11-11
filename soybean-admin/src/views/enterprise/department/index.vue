<script setup lang="ts">
import { h, onMounted, reactive, ref, computed } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NCheckbox, NPopconfirm, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { fetchDepartments, fetchDepartmentMenus, fetchMenuList, fetchQuickCreateDepartment, fetchSetDepartmentMenus } from '@/service/api/menu';
import DepartmentOperateDrawer from './components/department-operate-drawer.vue';
import { getDeptStatusMeta } from './components/department-meta';
import { useDepartmentData } from './components/use-department-data';

defineOptions({ name: 'EnterpriseDepartment' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
  status: Api.Department.Status | null;
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

const { loading, records, total, loadList, deleteDepartment } = useDepartmentData();

const searchParams = reactive<SearchParams>({
  current: 1,
  size: 10,
  keyword: '',
  status: null
});

const checkedRowKeys = ref<DataTableRowKey[]>([]);
const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.Department.Item | null>(null);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 个部门`,
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

// ---- 部门权限相关 ----
const departments = ref<Api.Menu.Department[]>([]);
const menus = ref<Api.Menu.Item[]>([]);
const selectedDepartment = ref<number | null>(null);
const grantedMenuIds = ref<number[]>([]);
const permissionLoading = ref(false);
const departmentName = ref('');
const departmentCode = ref('');

const columns: DataTableColumns<Api.Department.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'name',
    title: '部门名称',
    width: 180,
    fixed: 'left',
    render: row => h('div', { class: 'flex-col' }, [
      h('span', { class: 'font-medium' }, row.name),
      h('span', { class: 'text-12px text-gray-400' }, row.code)
    ])
  },
  {
    key: 'leaderName',
    title: '负责人',
    width: 120,
    render: row => row.leaderName || h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'memberCount',
    title: '成员数',
    width: 90,
    align: 'center',
    render: row => h('span', { class: 'font-mono text-primary' }, String(row.memberCount))
  },
  {
    key: 'status',
    title: '状态',
    width: 90,
    render: row => {
      const meta = getDeptStatusMeta(row.status);
      return h(NTag, { size: 'small', bordered: false, type: meta.tagType }, { default: () => meta.label });
    }
  },
  { key: 'sortOrder', title: '排序', width: 70, align: 'center', render: row => String(row.sortOrder) },
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
          default: () => '确认删除该部门吗？有关联人员时不可删除。'
        })
      ])
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

function handleEdit(row: Api.Department.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

async function handleDelete(row: Api.Department.Item) {
  const success = await deleteDepartment([row.id]);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  await afterDeleted(1);
}

async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(key => Number(key));
  const success = await deleteDepartment(ids);
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

// ---- 部门权限方法 ----
async function loadPermissionData() {
  const [deptResult, menuResult] = await Promise.all([fetchDepartments(), fetchMenuList()]);
  if (!deptResult.error && deptResult.data) {
    departments.value = deptResult.data.items;
    if (selectedDepartment.value == null && departments.value.length) selectedDepartment.value = departments.value[0].id;
  }
  if (!menuResult.error && menuResult.data) menus.value = menuResult.data.items;
  await loadDepartmentPermissions();
  collapseAll();
}

async function loadDepartmentPermissions() {
  if (!selectedDepartment.value) {
    grantedMenuIds.value = [];
    return;
  }
  permissionLoading.value = true;
  const { data, error } = await fetchDepartmentMenus(selectedDepartment.value);
  if (!error && data) grantedMenuIds.value = data.menuIds;
  permissionLoading.value = false;
}

async function savePermissions() {
  if (!selectedDepartment.value) return;
  const { error } = await fetchSetDepartmentMenus(selectedDepartment.value, grantedMenuIds.value);
  if (!error) window.$message?.success('部门权限已保存');
}

// ---- 菜单树 ----
interface MenuTreeNode {
  item: Api.Menu.Item;
  children: MenuTreeNode[];
}

const menuTree = computed<MenuTreeNode[]>(() => {
  if (!menus.value.length) return [];
  const map = new Map<number, MenuTreeNode>();
  const roots: MenuTreeNode[] = [];
  menus.value.forEach(item => map.set(item.id, { item, children: [] }));
  map.forEach(node => {
    const parent = map.get(node.item.parentId);
    if (parent) parent.children.push(node);
    else roots.push(node);
  });
  const sort = (nodes: MenuTreeNode[]) => {
    nodes.sort((a, b) => (a.item.sortOrder || 0) - (b.item.sortOrder || 0) || a.item.id - b.item.id);
    nodes.forEach(n => sort(n.children));
  };
  sort(roots);
  return roots;
});

function collectLeafIds(node: MenuTreeNode): number[] {
  if (!node.children.length) return [node.item.id];
  return node.children.flatMap(collectLeafIds);
}

function allLeafChecked(node: MenuTreeNode): boolean {
  const leaves = collectLeafIds(node);
  return leaves.every(id => grantedMenuIds.value.includes(id));
}

function someLeafChecked(node: MenuTreeNode): boolean {
  return collectLeafIds(node).some(id => grantedMenuIds.value.includes(id));
}

function toggleNode(node: MenuTreeNode) {
  const leafIds = collectLeafIds(node);
  const allIn = allLeafChecked(node);
  if (allIn) {
    grantedMenuIds.value = grantedMenuIds.value.filter(id => !leafIds.includes(id));
  } else {
    const add = leafIds.filter(id => !grantedMenuIds.value.includes(id));
    grantedMenuIds.value = [...grantedMenuIds.value, ...add];
  }
}

function toggleLeaf(id: number) {
  const idx = grantedMenuIds.value.indexOf(id);
  if (idx >= 0) {
    grantedMenuIds.value.splice(idx, 1);
  } else {
    grantedMenuIds.value.push(id);
  }
}

const allMenuIds = computed(() =>
  menus.value.map(m => m.id)
);

function selectAll() {
  grantedMenuIds.value = [...allMenuIds.value];
}

function clearAll() {
  grantedMenuIds.value = [];
}

const collapseExpanded = ref<string[]>([]);

function expandAll() {
  collapseExpanded.value = menuTree.value.map(n => String(n.item.id));
}

function collapseAll() {
  collapseExpanded.value = [];
}

function toggleCollapse(id: string) {
  const idx = collapseExpanded.value.indexOf(id);
  if (idx >= 0) {
    collapseExpanded.value.splice(idx, 1);
  } else {
    collapseExpanded.value.push(id);
  }
}

async function createDepartment() {
  if (!departmentName.value.trim() || !departmentCode.value.trim()) {
    window.$message?.warning('请填写部门名称和编码');
    return;
  }
  const { error } = await fetchQuickCreateDepartment(departmentName.value.trim(), departmentCode.value.trim());
  if (!error) {
    departmentName.value = '';
    departmentCode.value = '';
    window.$message?.success('部门已创建');
    await Promise.all([getData(), loadPermissionData()]);
  }
}

onMounted(async () => {
  await getData();
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <NCard :bordered="false" class="card-wrapper" size="small">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-16px font-semibold">部门管理</h1>
            <p class="mt-4px text-13px text-gray-500">部门列表与菜单访问权限</p>
          </div>
        </div>
      </NCard>

      <NCard :bordered="false" class="card-wrapper" size="small">
        <NTabs type="line" animated @update:value="(val: string) => { if (val === 'permissions') loadPermissionData(); }">
          <NTabPane name="list" tab="部门列表">
            <NSpace vertical :size="12">
              <div class="flex flex-wrap items-center gap-8px">
                <NInput
                  v-model:value="searchParams.keyword"
                  class="w-200px min-w-140px flex-shrink"
                  placeholder="搜索部门名称 / 编码"
                  clearable
                  @keyup.enter="handleSearch"
                >
                  <template #prefix><SvgIcon icon="mdi:magnify" class="text-gray-400" /></template>
                </NInput>
                <NSelect
                  v-model:value="searchParams.status"
                  class="w-100px flex-shrink-0"
                  :options="[{ value: 1, label: '启用' }, { value: 2, label: '停用' }]"
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
                    <template #icon><SvgIcon icon="mdi:office-building-plus-outline" /></template>
                    新增部门
                  </NButton>
                  <NPopconfirm @positive-click="handleBatchDelete">
                    <template #trigger>
                      <NButton size="small" ghost type="error" :disabled="!checkedRowKeys.length">
                        <template #icon><SvgIcon icon="mdi:delete-outline" /></template>
                        批量删除
                      </NButton>
                    </template>
                    确认删除选中的 {{ checkedRowKeys.length }} 个部门吗？
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
                :scroll-x="800"
                remote
                size="small"
                striped
              />
            </NSpace>
          </NTabPane>

          <NTabPane name="permissions" tab="部门权限">
            <div class="perm-layout">
              <!-- 顶部操作栏 -->
              <div class="perm-header">
                <div class="perm-header__left">
                  <NSelect
                    v-model:value="selectedDepartment"
                    class="perm-dept-select"
                    :options="departments.map(item => ({ label: item.name, value: item.id }))"
                    placeholder="选择部门"
                    @update:value="loadDepartmentPermissions"
                  />
                  <NButton type="primary" :disabled="!selectedDepartment" :loading="permissionLoading" @click="savePermissions">
                    保存授权
                  </NButton>
                </div>
                <div class="perm-header__right">
                  <span class="perm-counter">
                    <span class="perm-counter__num">{{ grantedMenuIds.length }}</span>
                    <span class="perm-counter__sep">/</span>
                    <span class="perm-counter__num">{{ allMenuIds.length }}</span>
                  </span>
                </div>
              </div>

              <!-- 快捷操作 -->
              <div class="perm-actions">
                <span class="perm-actions__label">快捷操作</span>
                <NButton size="tiny" quaternary @click="selectAll">
                  <template #icon><SvgIcon icon="mdi:check-all" /></template>
                  全选
                </NButton>
                <NButton size="tiny" quaternary @click="clearAll">
                  <template #icon><SvgIcon icon="mdi:close-circle-outline" /></template>
                  清空
                </NButton>
                <NButton size="tiny" quaternary @click="expandAll">
                  <template #icon><SvgIcon icon="mdi:unfold-more-horizontal" /></template>
                  展开
                </NButton>
                <NButton size="tiny" quaternary @click="collapseAll">
                  <template #icon><SvgIcon icon="mdi:unfold-less-horizontal" /></template>
                  收起
                </NButton>
              </div>

              <template v-if="!selectedDepartment">
                <NEmpty description="选择部门后开始配置菜单权限" class="py-60px" />
              </template>
              <template v-else-if="!menuTree.length">
                <NEmpty description="暂无菜单数据，请先在「菜单管理」中配置" class="py-60px" />
              </template>
              <template v-else>
                <!-- 菜单权限卡片组 -->
                <div class="perm-cards">
                  <div
                    v-for="dir in menuTree"
                    :key="dir.item.id"
                    class="perm-card"
                  >
                    <!-- 目录头 -->
                    <div class="perm-card__head">
                      <div class="perm-card__check" @click.stop>
                        <NCheckbox
                          :checked="allLeafChecked(dir)"
                          :indeterminate="!allLeafChecked(dir) && someLeafChecked(dir)"
                          size="small"
                          @update:checked="toggleNode(dir)"
                        />
                      </div>
                      <SvgIcon
                        :icon="dir.item.icon || 'mdi:folder-outline'"
                        class="perm-card__icon"
                      />
                      <span class="perm-card__title">{{ dir.item.title }}</span>
                      <button
                        class="perm-card__toggle"
                        @click="toggleCollapse(String(dir.item.id))"
                      >
                        <SvgIcon
                          :icon="collapseExpanded.includes(String(dir.item.id)) ? 'mdi:chevron-up' : 'mdi:chevron-down'"
                          class="text-16px"
                        />
                      </button>
                    </div>

                    <!-- 子菜单列表 -->
                    <div
                      v-show="collapseExpanded.includes(String(dir.item.id))"
                      class="perm-card__body"
                    >
                      <template v-if="!dir.children.length">
                        <span class="perm-card__empty">此目录下暂未配置子菜单</span>
                      </template>
                      <template v-else>
                        <div
                          v-for="child in dir.children"
                          :key="child.item.id"
                          class="perm-item"
                          :class="{ 'perm-item--active': grantedMenuIds.includes(child.item.id) }"
                          @click="toggleLeaf(child.item.id)"
                        >
                          <div class="perm-item__left" @click.stop>
                            <NCheckbox
                              :checked="grantedMenuIds.includes(child.item.id)"
                              size="small"
                              @update:checked="toggleLeaf(child.item.id)"
                            />
                          </div>
                          <SvgIcon
                            :icon="child.item.icon || 'mdi:circle-small'"
                            class="perm-item__icon"
                          />
                          <div class="perm-item__text">
                            <span class="perm-item__name">{{ child.item.title }}</span>
                            <span class="perm-item__path">{{ child.item.routePath }}</span>
                          </div>
                        </div>
                      </template>
                    </div>
                  </div>
                </div>
              </template>
            </div>
          </NTabPane>
        </NTabs>
      </NCard>
    </NSpace>

    <DepartmentOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      @submitted="handleSubmitted"
    />
  </div>
</template>

<style scoped>
/* ===== 部门权限 - 企业级精致布局 ===== */

.perm-layout {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ---- 顶部栏 ---- */
.perm-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.perm-header__left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.perm-dept-select {
  width: 220px;
}

.perm-header__right {
  display: flex;
  align-items: center;
}

.perm-counter {
  display: inline-flex;
  align-items: baseline;
  padding: 4px 14px;
  border-radius: 20px;
  background: rgb(var(--primary-color) / 0.06);
  font-variant-numeric: tabular-nums;
}

.perm-counter__num {
  font-size: 15px;
  font-weight: 600;
  color: rgb(var(--primary-color));
  line-height: 1;
}

.perm-counter__sep {
  margin: 0 3px;
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.3);
}

/* ---- 快捷操作 ---- */
.perm-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 6px 12px;
  background: rgb(var(--base-text-color) / 0.03);
  border-radius: 8px;
  border: 1px solid rgb(var(--base-text-color) / 0.06);
}

.perm-actions__label {
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.4);
  margin-right: 6px;
  letter-spacing: 0.3px;
}

/* ---- 菜单卡片组 ---- */
.perm-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ---- 单个菜单目录卡片 ---- */
.perm-card {
  border: 1px solid rgb(var(--base-text-color) / 0.08);
  border-radius: 10px;
  overflow: hidden;
  transition: border-color 0.2s, box-shadow 0.2s;
  background: rgb(var(--base-text-color) / 0.01);
}

.perm-card:hover {
  border-color: rgb(var(--primary-color) / 0.25);
  box-shadow: 0 1px 6px rgb(var(--primary-color) / 0.04);
}

/* ---- 目录头部 ---- */
.perm-card__head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  cursor: default;
  user-select: none;
}

.perm-card__check {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.perm-card__icon {
  font-size: 18px;
  color: rgb(var(--base-text-color) / 0.55);
  flex-shrink: 0;
}

.perm-card__title {
  font-size: 14px;
  font-weight: 600;
  color: rgb(var(--base-text-color));
  flex: 1;
  min-width: 0;
}

.perm-card__toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: rgb(var(--base-text-color) / 0.4);
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.15s, color 0.15s;
}

.perm-card__toggle:hover {
  background: rgb(var(--base-text-color) / 0.06);
  color: rgb(var(--base-text-color) / 0.7);
}

/* ---- 子菜单列表 ---- */
.perm-card__body {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px 16px 14px 42px;
  border-top: 1px solid rgb(var(--base-text-color) / 0.05);
}

.perm-card__empty {
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.35);
  padding: 4px 0;
}

/* ---- 单个菜单项 ---- */
.perm-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 12px 7px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s, box-shadow 0.15s;
  min-width: 0;
  border: 1px solid transparent;
}

.perm-item:hover {
  background: rgb(var(--base-text-color) / 0.04);
}

.perm-item--active {
  background: rgb(var(--primary-color) / 0.06);
  border-color: rgb(var(--primary-color) / 0.15);
  box-shadow: 0 0 0 1px rgb(var(--primary-color) / 0.06);
}

.perm-item--active:hover {
  background: rgb(var(--primary-color) / 0.09);
}

.perm-item__left {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.perm-item__icon {
  font-size: 15px;
  color: rgb(var(--base-text-color) / 0.4);
  flex-shrink: 0;
}

.perm-item--active .perm-item__icon {
  color: rgb(var(--primary-color) / 0.7);
}

.perm-item__text {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.perm-item__name {
  font-size: 13px;
  font-weight: 500;
  color: rgb(var(--base-text-color));
  line-height: 1.35;
}

.perm-item--active .perm-item__name {
  color: rgb(var(--primary-color));
}

.perm-item__path {
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.35);
  font-family: 'SF Mono', 'Fira Code', monospace;
  line-height: 1.3;
}

/* ---- 移动端 ---- */
@media (max-width: 640px) {
  .perm-dept-select {
    width: 160px;
  }

  .perm-card__body {
    padding-left: 28px;
  }

  .perm-item {
    padding: 6px 10px 6px 6px;
  }
}
</style>