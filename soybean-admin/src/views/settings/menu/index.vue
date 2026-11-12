<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { NButton, NPopconfirm, NSelect, NSwitch, NInputNumber, NCheckbox, NModal, NForm, NFormItem, NInput, useMessage } from 'naive-ui';
import { fetchDeleteMenu, fetchMenuList, fetchSaveMenu } from '@/service/api/menu';

defineOptions({ name: 'SettingsMenu' });

const message = useMessage();
const loading = ref(false);
const menus = ref<Api.Menu.Item[]>([]);
const editorVisible = ref(false);
const editingId = ref<number | undefined>();
const fixedParentTitle = ref('');
const parentLocked = ref(false);
const expandedDirs = ref<Set<number>>(new Set());
const form = reactive<Api.Menu.Form>({
  parentId: 0, title: '', routeName: '', routePath: '', component: '', icon: '', sortOrder: 1, enabled: true, hideInMenu: false
});

const parentOptions = computed(() => [
  { label: '顶级菜单', value: 0 },
  ...menus.value.filter(item => item.component === 'layout.base').map(item => ({ label: item.title, value: item.id }))
]);

const routeOptions = [
  { label: '首页', value: 'home', path: '/home', component: 'layout.base$view.home' },
  { label: '订单中心', value: 'order', path: '/order', component: 'layout.base' },
  { label: '全部订单', value: 'order_all', path: '/order/all', component: 'view.order_all' },
  { label: '行程监控', value: 'order_monitor', path: '/order/monitor', component: 'view.order_monitor' },
  { label: '运力管理', value: 'transport', path: '/transport', component: 'layout.base' },
  { label: '司机管理', value: 'transport_driver', path: '/transport/driver', component: 'view.transport_driver' },
  { label: '车辆管理', value: 'transport_vehicle', path: '/transport/vehicle', component: 'view.transport_vehicle' },
  { label: '财务管理', value: 'finance', path: '/finance', component: 'layout.base' },
  { label: '收款管理', value: 'finance_payment', path: '/finance/payment', component: 'view.finance_payment' },
  { label: '付款管理', value: 'finance_payout', path: '/finance/payout', component: 'view.finance_payout' },
  { label: '记账凭证', value: 'finance_voucher', path: '/finance/voucher', component: 'view.finance_voucher' },
  { label: '财务基础资料', value: 'finance_basic', path: '/finance/basic', component: 'view.finance_basic' },
  { label: '报表', value: 'report', path: '/report', component: 'layout.base' },
  { label: '经营总览', value: 'report_overview', path: '/report/overview', component: 'view.report_overview' },
  { label: '资金流水', value: 'report_flows', path: '/report/flows', component: 'view.report_flows' },
  { label: '科目余额', value: 'report_balance', path: '/report/balance', component: 'view.report_balance' },
  { label: '系统设置', value: 'settings', path: '/settings', component: 'layout.base' },
  { label: '客户管理', value: 'settings_customer', path: '/settings/customer', component: 'view.settings_customer' },
  { label: '计价配置', value: 'settings_pricing', path: '/settings/pricing', component: 'view.settings_pricing' },
  { label: '菜单管理', value: 'settings_menu', path: '/settings/menu', component: 'view.settings_menu' },
  { label: '企业管理', value: 'enterprise', path: '/enterprise', component: 'layout.base' },
  { label: '人员管理', value: 'enterprise_staff', path: '/enterprise/staff', component: 'view.enterprise_staff' },
  { label: '部门管理', value: 'enterprise_department', path: '/enterprise/department', component: 'view.enterprise_department' },
  { label: '用户管理', value: 'enterprise_user', path: '/enterprise/user', component: 'view.enterprise_user' }
];

interface MenuTreeItem extends Api.Menu.Item {
  children?: MenuTreeItem[];
}

const menuTree = computed<MenuTreeItem[]>(() => {
  const nodes = new Map<number, MenuTreeItem>();
  const roots: MenuTreeItem[] = [];
  menus.value.forEach(item => nodes.set(item.id, { ...item, children: [] }));
  nodes.forEach(node => {
    const parent = nodes.get(node.parentId);
    if (parent) parent.children!.push(node);
    else roots.push(node);
  });
  const sortBranch = (items: MenuTreeItem[]) => {
    items.sort((left, right) => left.sortOrder - right.sortOrder || left.id - right.id);
    items.forEach(item => sortBranch(item.children ?? []));
  };
  sortBranch(roots);
  return roots;
});

const isDir = (item: MenuTreeItem) => item.component === 'layout.base';
const isDirShowing = (item: MenuTreeItem) => expandedDirs.value.has(item.id);

function toggleDir(id: number) {
  if (expandedDirs.value.has(id)) expandedDirs.value.delete(id);
  else expandedDirs.value.add(id);
  expandedDirs.value = new Set(expandedDirs.value);
}

async function load() {
  loading.value = true;
  const { data, error } = await fetchMenuList();
  if (!error && data) {
    menus.value = data.items;
  }
  loading.value = false;
}

const selectedRoute = ref<string | null>(null);

function openEditor(item?: Api.Menu.Item, parent?: Api.Menu.Item) {
  editingId.value = item?.id;
  parentLocked.value = Boolean(parent);
  fixedParentTitle.value = parent?.title ?? '';
  Object.assign(form, item ?? {
    parentId: parent?.id ?? 0, title: '', routeName: '', routePath: '', component: '', icon: '', sortOrder: 1, enabled: true, hideInMenu: false
  });
  selectedRoute.value = item?.routeName ?? null;
  editorVisible.value = true;
}

function applyRoute(routeName: string | null) {
  selectedRoute.value = routeName;
  const route = routeOptions.find(item => item.value === routeName);
  if (!route) return;
  form.routeName = route.value;
  form.routePath = route.path;
  form.component = route.component;
  if (!form.title) form.title = route.label;
}

async function save(values: Api.Menu.Form = { ...form, id: editingId.value }) {
  if (!values.title.trim() || !values.routeName) {
    message.warning('请填写菜单名称并选择页面');
    return;
  }
  const { error } = await fetchSaveMenu(values);
  if (!error) {
    message.success('菜单已保存');
    editorVisible.value = false;
    await load();
  }
}

function toggleEnabled(item: MenuTreeItem) {
  item.enabled = !item.enabled;
  save({ ...item });
}

async function remove(id: number) {
  const { error } = await fetchDeleteMenu(id);
  if (!error) {
    message.success('菜单已删除');
    await load();
  }
}

function getIcon(item: MenuTreeItem): string {
  return item.icon || (isDir(item) ? 'mdi:folder-outline' : 'mdi:circle-small');
}

onMounted(() => void load());
</script>

<template>
  <div class="menu-page">
    <!-- 头部 -->
    <div class="menu-header">
      <div class="menu-header__info">
        <h1 class="menu-header__title">菜单管理</h1>
        <p class="menu-header__desc">配置系统导航菜单与路由绑定</p>
      </div>
      <NButton type="primary" @click="openEditor()">
        <template #icon><SvgIcon icon="mdi:plus" /></template>
        新增目录
      </NButton>
    </div>

    <!-- 树列表 -->
    <div class="menu-tree-card">
      <div v-if="!menuTree.length && !loading" class="menu-empty">
        <SvgIcon icon="mdi:menu-open" class="text-40px text-gray-300 mb-12px" />
        <span class="text-14px text-gray-400">暂无菜单，点击「新增目录」开始配置</span>
      </div>
      <div v-else class="menu-tree" :class="{ 'menu-tree--loading': loading }">
        <template v-for="dir in menuTree" :key="dir.id">
          <!-- 目录行 -->
          <div class="dir-row" :class="{ 'dir-row--disabled': !dir.enabled }">
            <button class="dir-toggle" @click="toggleDir(dir.id)">
              <SvgIcon :icon="isDirShowing(dir) ? 'mdi:chevron-down' : 'mdi:chevron-right'" class="text-18px" />
            </button>
            <SvgIcon :icon="getIcon(dir)" class="dir-icon" />
            <span class="dir-title">{{ dir.title }}</span>
            <span class="dir-path">{{ dir.routePath }}</span>
            <NPopconfirm @positive-click="() => toggleEnabled(dir)">
              <template #trigger>
                <NSwitch :value="dir.enabled" size="small" @click.stop />
              </template>
              确认{{ dir.enabled ? '禁用' : '启用' }}此目录？
            </NPopconfirm>
            <span class="dir-children-hint">{{ dir.children?.length || 0 }} 个子菜单</span>
            <div class="dir-actions">
              <NButton size="tiny" quaternary type="primary" @click="openEditor(undefined, dir)">+ 子菜单</NButton>
              <NButton size="tiny" quaternary @click="openEditor(dir)">编辑</NButton>
              <NPopconfirm @positive-click="() => remove(dir.id)">
                <template #trigger><NButton size="tiny" quaternary type="error">删除</NButton></template>
                确认删除此目录及其子菜单？
              </NPopconfirm>
            </div>
          </div>

          <!-- 子菜单列表 -->
          <div v-show="isDirShowing(dir)" class="children-list">
            <div
              v-for="child in dir.children ?? []"
              :key="child.id"
              class="child-row"
              :class="{ 'child-row--disabled': !child.enabled, 'child-row--hidden': child.hideInMenu }"
            >
              <span class="child-indent" />
              <SvgIcon :icon="getIcon(child)" class="child-icon" />
              <span class="child-title">{{ child.title }}</span>
              <span class="child-path">{{ child.routePath }}</span>
              <span v-if="child.hideInMenu" class="child-badge">隐藏</span>
              <div class="child-switches">
                <NPopconfirm @positive-click="() => toggleEnabled(child)">
                  <template #trigger><NSwitch :value="child.enabled" size="small" /></template>
                  确认{{ child.enabled ? '禁用' : '启用' }}？
                </NPopconfirm>
              </div>
              <div class="child-actions">
                <NButton size="tiny" quaternary @click="openEditor(child)">编辑</NButton>
                <NPopconfirm @positive-click="() => remove(child.id)">
                  <template #trigger><NButton size="tiny" quaternary type="error">删除</NButton></template>
                  确认删除此菜单？
                </NPopconfirm>
              </div>
            </div>
            <div v-if="!(dir.children?.length)" class="children-empty">
              <span class="text-13px text-gray-400">暂无子菜单</span>
              <NButton size="tiny" quaternary type="primary" @click="openEditor(undefined, dir)">添加一个</NButton>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <NModal v-model:show="editorVisible" preset="card" :title="editingId ? '编辑菜单' : '新增菜单'" class="w-[min(540px,calc(100vw-32px))]">
      <NForm :model="form" label-placement="top">
        <div class="grid grid-cols-1 gap-x-4 sm:grid-cols-2">
          <NFormItem label="菜单名称" required><NInput v-model:value="form.title" placeholder="如：客户管理" /></NFormItem>
          <NFormItem v-if="parentLocked" label="所属目录"><NInput :value="fixedParentTitle" disabled /></NFormItem>
          <NFormItem v-else-if="editingId" label="父级目录"><NSelect v-model:value="form.parentId" :options="parentOptions" /></NFormItem>
          <NFormItem label="页面路由" required>
            <NSelect
              v-model:value="selectedRoute"
              :options="routeOptions.filter(item => form.parentId === 0 ? item.component.startsWith('layout.base') : item.component.startsWith('view.')).map(item => ({ label: item.label, value: item.value }))"
              placeholder="选择已存在的页面"
              @update:value="applyRoute"
            />
          </NFormItem>
          <NFormItem label="图标"><NInput v-model:value="form.icon" placeholder="mdi:account-outline" /></NFormItem>
          <NFormItem label="排序"><NInputNumber v-model:value="form.sortOrder" class="w-full" :min="0" /></NFormItem>
          <div class="flex items-center gap-6 pb-6">
            <NCheckbox v-model:checked="form.enabled">启用</NCheckbox>
            <NCheckbox v-model:checked="form.hideInMenu">隐藏导航</NCheckbox>
          </div>
        </div>
        <div class="flex justify-end gap-2">
          <NButton @click="editorVisible = false">取消</NButton>
          <NButton type="primary" @click="save()">保存</NButton>
        </div>
      </NForm>
    </NModal>
  </div>
</template>

<style scoped>
/* ===== 页面容器 ===== */
.menu-page {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ---- 头部 ---- */
.menu-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  padding: 16px 20px;
  border-radius: 10px;
  border: 1px solid rgb(var(--base-text-color) / 0.08);
  background: #fff;
}

.menu-header__title {
  font-size: 16px;
  font-weight: 600;
  line-height: 1.3;
}

.menu-header__desc {
  margin-top: 2px;
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.45);
}

/* ---- 树卡片 ---- */
.menu-tree-card {
  border-radius: 10px;
  border: 1px solid rgb(var(--base-text-color) / 0.08);
  background: #fff;
  overflow: hidden;
}

.menu-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 16px;
}

.menu-tree {
  opacity: 1;
  transition: opacity 0.2s;
}

.menu-tree--loading {
  opacity: 0.5;
  pointer-events: none;
}

/* ---- 目录行 ---- */
.dir-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid rgb(var(--base-text-color) / 0.05);
  transition: background 0.15s;
}

.dir-row:hover {
  background: rgb(var(--base-text-color) / 0.02);
}

.dir-row--disabled .dir-title,
.dir-row--disabled .dir-icon {
  opacity: 0.45;
}

.dir-toggle {
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

.dir-toggle:hover {
  background: rgb(var(--base-text-color) / 0.06);
  color: rgb(var(--base-text-color) / 0.7);
}

.dir-icon {
  font-size: 20px;
  color: rgb(var(--primary-color) / 0.65);
  flex-shrink: 0;
}

.dir-title {
  font-size: 14px;
  font-weight: 600;
  color: rgb(var(--base-text-color));
  flex-shrink: 0;
}

.dir-path {
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.35);
  font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', monospace;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dir-children-hint {
  margin-left: auto;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.35);
  white-space: nowrap;
  flex-shrink: 0;
}

.dir-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
}

/* ---- 子菜单列表 ---- */
.children-list {
  border-bottom: 1px solid rgb(var(--base-text-color) / 0.04);
}

.child-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px 10px 18px;
  transition: background 0.15s;
  border-top: 1px solid rgb(var(--base-text-color) / 0.03);
}

.child-row:hover {
  background: rgb(var(--base-text-color) / 0.02);
}

.child-row--disabled .child-title,
.child-row--disabled .child-icon {
  opacity: 0.4;
}

.child-indent {
  width: 20px;
  flex-shrink: 0;
}

.child-icon {
  font-size: 15px;
  color: rgb(var(--base-text-color) / 0.45);
  flex-shrink: 0;
}

.child-title {
  font-size: 13px;
  font-weight: 500;
  color: rgb(var(--base-text-color));
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.child-path {
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.3);
  font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', monospace;
  margin-left: auto;
  min-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex-shrink: 1;
}

.child-badge {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 4px;
  background: rgb(var(--base-text-color) / 0.06);
  color: rgb(var(--base-text-color) / 0.4);
  white-space: nowrap;
  flex-shrink: 0;
}

.child-switches {
  flex-shrink: 0;
}

.child-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
}

.children-empty {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px 12px 38px;
}
</style>