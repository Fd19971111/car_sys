<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, DataTableRowKey } from 'naive-ui';
import { NButton, NPopconfirm, NSwitch, NTag } from 'naive-ui';
import { jsonClone } from '@sa/utils';
import { useAppStore } from '@/store/modules/app';
import PricingRuleOperateDrawer from './components/pricing-rule-operate-drawer.vue';
import { usePricingRuleData } from './components/use-pricing-rule-data';

defineOptions({ name: 'SettingsPricing' });

interface SearchParams {
  current: number;
  size: number;
  keyword: string;
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

type TagType = 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error';

/** 车型色板：让不同车型一眼可区分 */
const CAR_TYPE_TAG_TYPES: Record<string, TagType> = {
  经济型: 'default',
  舒适型: 'info',
  商务型: 'warning',
  六座商务: 'success',
  七座商务: 'success',
  豪华型: 'error'
};

const appStore = useAppStore();
const isMobile = computed(() => appStore.isMobile);
const { loading, records, total, loadList, updateRule, deleteRules } = usePricingRuleData();
const searchParams = reactive<SearchParams>({ current: 1, size: 10, keyword: '' });
const checkedRowKeys = ref<DataTableRowKey[]>([]);
const drawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData = ref<Api.PricingRule.Item | null>(null);
const togglingIds = ref<number[]>([]);

const pagination = reactive<PaginationState>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: info => `共 ${info.itemCount ?? 0} 条计价规则`,
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

const tablePagination = computed(() => ({
  ...pagination,
  pageSlot: isMobile.value ? 3 : 9,
  prefix: isMobile.value ? undefined : pagination.prefix
}));

function formatMoney(value: number) {
  return `¥${Number(value ?? 0).toFixed(2)}`;
}

const columns: DataTableColumns<Api.PricingRule.Item> = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    key: 'name',
    title: '规则名称',
    minWidth: 200,
    fixed: 'left',
    render: row =>
      h('div', { class: 'flex-col gap-1px' }, [
        h('span', { class: 'font-medium' }, row.name),
        h('span', { class: 'text-11px text-gray-400' }, `ID ${row.id}`)
      ])
  },
  {
    key: 'carType',
    title: '适用车型',
    width: 130,
    render: row =>
      h(NTag, { size: 'small', bordered: false, type: CAR_TYPE_TAG_TYPES[row.carType] ?? 'default' }, {
        default: () => row.carType || '—'
      })
  },
  {
    key: 'basePrice',
    title: '起步价（元）',
    width: 130,
    align: 'right',
    render: row => h('span', { class: 'font-mono font-medium' }, formatMoney(row.basePrice))
  },
  {
    key: 'perKmPrice',
    title: '超里程单价（元/公里）',
    width: 180,
    align: 'right',
    render: row => h('span', { class: 'font-mono' }, formatMoney(row.perKmPrice))
  },
  {
    key: 'enabled',
    title: '启用规则',
    width: 110,
    align: 'center',
    render: row =>
      h(NSwitch, {
        value: row.enabled,
        size: 'small',
        loading: togglingIds.value.includes(row.id),
        onUpdateValue: (value: boolean) => handleToggleEnabled(row, value)
      })
  },
  {
    key: 'operate',
    title: '操作',
    width: 130,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-4px' }, [
        h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => handleEdit(row) }, {
          default: () => '编辑'
        }),
        h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
          trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
          default: () => '确认删除该计价规则吗？'
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
  pagination.page = 1;
  await getData();
}

function handleAdd() {
  operateType.value = 'add';
  editingData.value = null;
  drawerVisible.value = true;
}

function handleEdit(row: Api.PricingRule.Item) {
  operateType.value = 'edit';
  editingData.value = jsonClone(row);
  drawerVisible.value = true;
}

async function handleToggleEnabled(row: Api.PricingRule.Item, value: boolean) {
  if (togglingIds.value.includes(row.id)) return;
  togglingIds.value = [...togglingIds.value, row.id];
  try {
    const success = await updateRule(row.id, {
      name: row.name,
      carType: row.carType,
      basePrice: row.basePrice,
      perKmPrice: row.perKmPrice,
      enabled: value
    });
    if (!success) {
      window.$message?.error('状态更新失败，请稍后重试');
      return;
    }
    row.enabled = value;
    window.$message?.success(value ? '规则已启用' : '规则已停用');
  } finally {
    togglingIds.value = togglingIds.value.filter(id => id !== row.id);
  }
}

async function handleDelete(row: Api.PricingRule.Item) {
  const success = await deleteRules([row.id]);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('删除成功');
  await afterDeleted(1);
}

async function handleBatchDelete() {
  const ids = checkedRowKeys.value.map(Number);
  const success = await deleteRules(ids);
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

onMounted(() => void getData());
</script>

<template>
  <div class="p-16px">
    <NCard :bordered="false" class="card-wrapper" size="small">
      <NSpace vertical :size="12">
        <div class="flex flex-wrap items-center gap-8px">
          <NInput
            v-model:value="searchParams.keyword"
            class="w-260px min-w-160px"
            placeholder="搜索规则名称 / 适用车型"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix><SvgIcon icon="mdi:magnify" class="text-gray-400" /></template>
          </NInput>
          <NButton type="primary" ghost size="small" @click="handleSearch">
            <template #icon><SvgIcon icon="mdi:magnify" /></template>查询
          </NButton>
          <NButton size="small" @click="handleReset">
            <template #icon><SvgIcon icon="mdi:backspace-outline" /></template>重置
          </NButton>
          <div class="ml-auto flex items-center gap-8px">
            <NButton size="small" type="primary" ghost @click="handleAdd">
              <template #icon><SvgIcon icon="mdi:plus-circle-outline" /></template>新增计价规则
            </NButton>
            <NPopconfirm @positive-click="handleBatchDelete">
              <template #trigger>
                <NButton size="small" ghost type="error" :disabled="!checkedRowKeys.length">
                  <template #icon><SvgIcon icon="mdi:delete-outline" /></template>批量删除
                </NButton>
              </template>
              确认删除选中的 {{ checkedRowKeys.length }} 条计价规则吗？
            </NPopconfirm>
            <NButton size="small" @click="getData">
              <template #icon><SvgIcon icon="mdi:refresh" :class="{ 'animate-spin': loading }" /></template>刷新
            </NButton>
          </div>
        </div>
        <NDataTable
          v-model:checked-row-keys="checkedRowKeys"
          :columns="columns"
          :data="records"
          :loading="loading"
          :pagination="tablePagination"
          :row-key="row => row.id"
          :scroll-x="1000"
          remote
          size="small"
          striped
        />
      </NSpace>
    </NCard>
    <PricingRuleOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :row-data="editingData"
      @submitted="handleSubmitted"
    />
  </div>
</template>
