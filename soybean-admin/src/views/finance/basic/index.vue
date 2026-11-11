<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NPopconfirm, NSwitch } from 'naive-ui';
import { useAuthStore } from '@/store/modules/auth';
import {
  fetchDeleteFinAccount,
  fetchDeleteFinFundAccount,
  fetchGetFinAccounts,
  fetchGetFinFundAccounts,
  fetchToggleFinAccount,
  fetchToggleFinFundAccount
} from '@/service/api';
import AccountOperateDrawer from './components/account-operate-drawer.vue';
import FundAccountOperateDrawer from './components/fund-account-operate-drawer.vue';
import {
  ACCOUNT_CATEGORY_OPTIONS,
  FUND_TYPE_OPTIONS,
  getCategoryMeta,
  getDirectionMeta,
  getFundTypeMeta
} from './components/basic-meta';

defineOptions({
  name: 'FinanceBasic'
});

type TabKey = 'account' | 'fund';

const authStore = useAuthStore();

const canManage = computed(() => authStore.userInfo.role === 1);

const activeTab = ref<TabKey>('account');

const accountsLoading = ref(false);
const fundsLoading = ref(false);

const accounts = ref<Api.FinBasic.Account[]>([]);
const fundAccounts = ref<Api.FinBasic.FundAccount[]>([]);

const summary = ref<Api.FinBasic.AccountSummary>({
  total: 0,
  enabled: 0,
  disabled: 0,
  system: 0,
  fundAccountCount: 0
});

interface AccountFilters {
  keyword: string;
  categoryCode: Api.FinBasic.CategoryCode | null;
  enabled: 'enabled' | 'disabled' | null;
}

interface FundFilters {
  keyword: string;
  type: Api.FinBasic.FundAccountType | null;
}

const accountFilters = reactive<AccountFilters>({
  keyword: '',
  categoryCode: null,
  enabled: null
});

const fundFilters = reactive<FundFilters>({
  keyword: '',
  type: null
});

const accountStatusOptions = [
  { label: '启用中', value: 'enabled' },
  { label: '已停用', value: 'disabled' }
];

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: (info: { itemCount?: number }) => `共 ${info.itemCount ?? 0} 条`
});

const togglingAccountId = ref<number | null>(null);
const togglingFundId = ref<number | null>(null);

const accountDrawerVisible = ref(false);
const fundDrawerVisible = ref(false);
const operateType = ref<NaiveUI.TableOperateType>('add');
const editingAccount = ref<Api.FinBasic.Account | null>(null);
const editingFund = ref<Api.FinBasic.FundAccount | null>(null);

const statCards = computed(() => {
  const { total, enabled, disabled, system, fundAccountCount } = summary.value;
  const enabledPercent = total ? Math.round((enabled / total) * 100) : 0;

  return [
    {
      key: 'total',
      label: '会计科目',
      value: String(total),
      hint: `系统预置 ${system} 个`,
      icon: 'mdi:book-open-variant',
      color: '#1677ff',
      bg: 'rgba(22, 119, 255, 0.12)'
    },
    {
      key: 'enabled',
      label: '启用科目',
      value: String(enabled),
      hint: `占全部科目 ${enabledPercent}%`,
      icon: 'mdi:check-decagram-outline',
      color: '#18a058',
      bg: 'rgba(24, 160, 88, 0.14)'
    },
    {
      key: 'disabled',
      label: '停用科目',
      value: String(disabled),
      hint: '停用后不参与新记账',
      icon: 'mdi:archive-outline',
      color: '#fa8c16',
      bg: 'rgba(250, 140, 22, 0.14)'
    },
    {
      key: 'fund',
      label: '资金账户',
      value: String(fundAccountCount),
      hint: '现金与银行账户',
      icon: 'mdi:wallet-outline',
      color: '#722ed1',
      bg: 'rgba(114, 46, 209, 0.14)'
    }
  ];
});

/** 资金账户可关联的资产类科目 */
const ledgerOptions = computed(() =>
  accounts.value
    .filter(item => item.enabled && item.categoryCode === 'ASSET')
    .map(item => ({ label: `${item.code} ${item.name}`, value: item.code }))
);

/** 加载会计科目 */
async function loadAccounts() {
  accountsLoading.value = true;

  try {
    const { data, error } = await fetchGetFinAccounts({
      keyword: accountFilters.keyword.trim() || null,
      categoryCode: accountFilters.categoryCode,
      enabled: accountFilters.enabled === null ? null : accountFilters.enabled === 'enabled'
    });

    if (!error && data) {
      accounts.value = data.records ?? [];
      summary.value = data.summary;
    }
  } finally {
    accountsLoading.value = false;
  }
}

/** 加载资金账户 */
async function loadFunds() {
  fundsLoading.value = true;

  try {
    const { data, error } = await fetchGetFinFundAccounts({
      keyword: fundFilters.keyword.trim() || null,
      type: fundFilters.type
    });

    if (!error && data) {
      fundAccounts.value = data.records ?? [];
    }
  } finally {
    fundsLoading.value = false;
  }
}

async function refreshAll() {
  await Promise.all([loadAccounts(), loadFunds()]);
}

/** 会计科目查询 / 重置 */
async function handleAccountSearch() {
  await loadAccounts();
}

async function handleAccountReset() {
  accountFilters.keyword = '';
  accountFilters.categoryCode = null;
  accountFilters.enabled = null;

  await loadAccounts();
}

/** 资金账户查询 / 重置 */
async function handleFundSearch() {
  await loadFunds();
}

async function handleFundReset() {
  fundFilters.keyword = '';
  fundFilters.type = null;

  await loadFunds();
}

/** 切换科目启用状态 */
async function handleToggleAccount(row: Api.FinBasic.Account, value: boolean) {
  togglingAccountId.value = row.id;

  const { error } = await fetchToggleFinAccount(row.id, value);

  togglingAccountId.value = null;

  if (error) {
    return;
  }

  window.$message?.success(value ? '科目已启用' : '科目已停用');

  await Promise.all([loadAccounts(), loadFunds()]);
}

/** 切换资金账户启用状态 */
async function handleToggleFund(row: Api.FinBasic.FundAccount, value: boolean) {
  togglingFundId.value = row.id;

  const { error } = await fetchToggleFinFundAccount(row.id, value);

  togglingFundId.value = null;

  if (error) {
    return;
  }

  window.$message?.success(value ? '资金账户已启用' : '资金账户已停用');

  await loadFunds();
}

/** 删除会计科目 */
async function handleDeleteAccount(row: Api.FinBasic.Account) {
  const { error } = await fetchDeleteFinAccount(row.id);

  if (error) {
    return;
  }

  window.$message?.success('会计科目已删除');

  await refreshAll();
}

/** 删除资金账户 */
async function handleDeleteFund(row: Api.FinBasic.FundAccount) {
  const { error } = await fetchDeleteFinFundAccount(row.id);

  if (error) {
    return;
  }

  window.$message?.success('资金账户已删除');

  await refreshAll();
}

function handleCreate() {
  operateType.value = 'add';

  if (activeTab.value === 'account') {
    editingAccount.value = null;
    accountDrawerVisible.value = true;

    return;
  }

  editingFund.value = null;
  fundDrawerVisible.value = true;
}

function switchTab(tab: TabKey) {
  activeTab.value = tab;
  pagination.page = 1;
}

function handleEditAccount(row: Api.FinBasic.Account) {
  operateType.value = 'edit';
  editingAccount.value = row;
  accountDrawerVisible.value = true;
}

function handleEditFund(row: Api.FinBasic.FundAccount) {
  operateType.value = 'edit';
  editingFund.value = row;
  fundDrawerVisible.value = true;
}

function handleSubmitted() {
  void refreshAll();
}

const accountColumns: DataTableColumns<Api.FinBasic.Account> = [
  {
    key: 'code',
    title: '科目编码',
    width: 130,
    fixed: 'left',
    render: row => h('span', { class: 'font-mono text-13px font-medium' }, row.code)
  },
  {
    key: 'name',
    title: '科目名称',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render: row =>
      h('div', { class: 'flex-y-center gap-8px' }, [
        h('span', { class: 'truncate' }, row.name),
        row.isSystem
          ? h('span', { class: 'mini-badge' }, '系统')
          : h('span', { class: 'mini-badge mini-badge--custom' }, '自定义')
      ])
  },
  {
    key: 'categoryCode',
    title: '科目类别',
    width: 110,
    render: row => {
      const meta = getCategoryMeta(row.categoryCode);

      return h('span', { class: 'soft-tag', style: { color: meta.color, backgroundColor: meta.bg } }, meta.label);
    }
  },
  {
    key: 'directionCode',
    title: '余额方向',
    width: 100,
    render: row => {
      const meta = getDirectionMeta(row.directionCode);

      return h('span', { class: 'soft-tag', style: { color: meta.color, backgroundColor: meta.bg } }, meta.label);
    }
  },
  {
    key: 'enabled',
    title: '状态',
    width: 110,
    render: row =>
      h(NSwitch, {
        size: 'small',
        value: row.enabled,
        loading: togglingAccountId.value === row.id,
        disabled: !canManage.value,
        onUpdateValue: (value: boolean) => handleToggleAccount(row, value)
      })
  },
  {
    key: 'used',
    title: '引用情况',
    width: 100,
    render: row =>
      h(
        'span',
        {
          class: row.used ? 'soft-tag soft-tag--warning' : 'soft-tag soft-tag--muted'
        },
        row.used ? '已引用' : '未引用'
      )
  },
  {
    key: 'operate',
    title: '操作',
    width: 150,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-4px' }, [
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'primary',
            disabled: !canManage.value,
            onClick: () => handleEditAccount(row)
          },
          { default: () => '编辑' }
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDeleteAccount(row) },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'tiny',
                  quaternary: true,
                  type: 'error',
                  disabled: !canManage.value || row.isSystem || row.used,
                  title: row.isSystem ? '系统预置科目不可删除' : row.used ? '该科目已被引用，不可删除' : ''
                },
                { default: () => '删除' }
              ),
            default: () => '确认删除该会计科目吗？删除后不可恢复。'
          }
        )
      ])
  }
];

const fundColumns: DataTableColumns<Api.FinBasic.FundAccount> = [
  {
    key: 'type',
    title: '类型',
    width: 100,
    render: row => {
      const meta = getFundTypeMeta(row.type);

      return h('span', { class: 'soft-tag', style: { color: meta.color, backgroundColor: meta.bg } }, [
        h('span', { class: 'mr-6px inline-block h-6px w-6px rounded-full align-middle', style: { backgroundColor: meta.color } }),
        meta.label
      ]);
    }
  },
  {
    key: 'name',
    title: '账户名称',
    minWidth: 200,
    ellipsis: { tooltip: true },
    render: row => h('span', { class: 'font-medium' }, row.name)
  },
  {
    key: 'bankName',
    title: '开户银行',
    minWidth: 180,
    ellipsis: { tooltip: true },
    render: row => row.bankName || h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'accountLast4',
    title: '卡号',
    width: 120,
    render: row =>
      row.accountLast4
        ? h('span', { class: 'font-mono text-13px' }, `•••• ${row.accountLast4}`)
        : h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'ledgerAccount',
    title: '关联科目',
    minWidth: 190,
    render: row =>
      h('div', { class: 'flex-col' }, [
        h('span', { class: 'font-mono text-12px text-gray-400' }, row.ledgerAccountCode),
        h('span', { class: 'truncate text-13px' }, row.ledgerAccountName)
      ])
  },
  {
    key: 'enabled',
    title: '状态',
    width: 110,
    render: row =>
      h(NSwitch, {
        size: 'small',
        value: row.enabled,
        loading: togglingFundId.value === row.id,
        disabled: !canManage.value,
        onUpdateValue: (value: boolean) => handleToggleFund(row, value)
      })
  },
  {
    key: 'operate',
    title: '操作',
    width: 150,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-4px' }, [
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'primary',
            disabled: !canManage.value,
            onClick: () => handleEditFund(row)
          },
          { default: () => '编辑' }
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDeleteFund(row) },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'tiny',
                  quaternary: true,
                  type: 'error',
                  disabled: !canManage.value
                },
                { default: () => '删除' }
              ),
            default: () => '确认删除该资金账户吗？历史资金流水不会被删除。'
          }
        )
      ])
  }
];

onMounted(() => {
  void refreshAll();
});
</script>

<template>
  <div class="finance-basic">
    <section class="stat-grid">
      <NCard v-for="card in statCards" :key="card.key" :bordered="false" class="card-wrapper stat-card" size="small">
        <div class="stat-card__body">
          <div class="stat-card__icon" :style="{ color: card.color, backgroundColor: card.bg }">
            <SvgIcon :icon="card.icon" />
          </div>
          <div class="stat-card__text">
            <span class="stat-card__label">{{ card.label }}</span>
            <strong class="stat-card__value">{{ card.value }}</strong>
            <span class="stat-card__hint">{{ card.hint }}</span>
          </div>
        </div>
      </NCard>
    </section>

    <NAlert v-if="!canManage" class="readonly-tip" type="info" :bordered="false">
      当前账号仅可查看财务基础资料，新增、修改与停用请联系公司超级管理员。
    </NAlert>

    <NCard :bordered="false" class="card-wrapper main-card" size="small">
      <div class="toolbar">
        <div class="tab-group">
          <button
            type="button"
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'account' }"
            @click="switchTab('account')"
          >
            <SvgIcon icon="mdi:book-open-variant" />
            会计科目
            <span class="tab-count">{{ summary.total }}</span>
          </button>
          <button
            type="button"
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'fund' }"
            @click="switchTab('fund')"
          >
            <SvgIcon icon="mdi:wallet-outline" />
            资金账户
            <span class="tab-count">{{ summary.fundAccountCount }}</span>
          </button>
        </div>

        <div v-if="activeTab === 'account'" class="filter-group">
          <NInput
            v-model:value="accountFilters.keyword"
            class="w-230px min-w-150px"
            placeholder="搜索科目编码 / 名称"
            clearable
            @keyup.enter="handleAccountSearch"
          >
            <template #prefix>
              <SvgIcon icon="mdi:magnify" class="text-gray-400" />
            </template>
          </NInput>
          <NSelect
            v-model:value="accountFilters.categoryCode"
            class="w-120px"
            :options="ACCOUNT_CATEGORY_OPTIONS"
            placeholder="科目类别"
            clearable
          />
          <NSelect
            v-model:value="accountFilters.enabled"
            class="w-110px"
            :options="accountStatusOptions"
            placeholder="状态"
            clearable
          />
          <NButton type="primary" ghost size="small" @click="handleAccountSearch">
            <template #icon>
              <SvgIcon icon="mdi:magnify" />
            </template>
            查询
          </NButton>
          <NButton size="small" @click="handleAccountReset">
            <template #icon>
              <SvgIcon icon="mdi:backspace-outline" />
            </template>
            重置
          </NButton>
          <NButton v-if="canManage" type="primary" size="small" @click="handleCreate">
            <template #icon>
              <SvgIcon icon="mdi:book-plus-outline" />
            </template>
            新增科目
          </NButton>
        </div>

        <div v-else class="filter-group">
          <NInput
            v-model:value="fundFilters.keyword"
            class="w-230px min-w-150px"
            placeholder="搜索账户 / 银行 / 科目"
            clearable
            @keyup.enter="handleFundSearch"
          >
            <template #prefix>
              <SvgIcon icon="mdi:magnify" class="text-gray-400" />
            </template>
          </NInput>
          <NSelect
            v-model:value="fundFilters.type"
            class="w-110px"
            :options="FUND_TYPE_OPTIONS"
            placeholder="账户类型"
            clearable
          />
          <NButton type="primary" ghost size="small" @click="handleFundSearch">
            <template #icon>
              <SvgIcon icon="mdi:magnify" />
            </template>
            查询
          </NButton>
          <NButton size="small" @click="handleFundReset">
            <template #icon>
              <SvgIcon icon="mdi:backspace-outline" />
            </template>
            重置
          </NButton>
          <NButton v-if="canManage" type="primary" size="small" @click="handleCreate">
            <template #icon>
              <SvgIcon icon="mdi:wallet-plus-outline" />
            </template>
            新增资金账户
          </NButton>
        </div>
      </div>

      <NDataTable
        v-if="activeTab === 'account'"
        :columns="accountColumns"
        :data="accounts"
        :loading="accountsLoading"
        :pagination="pagination"
        :row-key="row => row.id"
        :scroll-x="1080"
        size="small"
        striped
        class="data-table"
      >
        <template #empty>
          <div class="empty-state">
            <SvgIcon icon="mdi:database-off-outline" />
            <p>暂无符合条件的会计科目</p>
          </div>
        </template>
      </NDataTable>

      <NDataTable
        v-else
        :columns="fundColumns"
        :data="fundAccounts"
        :loading="fundsLoading"
        :pagination="pagination"
        :row-key="row => row.id"
        :scroll-x="1000"
        size="small"
        striped
        class="data-table"
      >
        <template #empty>
          <div class="empty-state">
            <SvgIcon icon="mdi:wallet-outline" />
            <p>暂无符合条件资金账户</p>
          </div>
        </template>
      </NDataTable>
    </NCard>

    <AccountOperateDrawer
      v-model:visible="accountDrawerVisible"
      :operate-type="operateType"
      :row-data="editingAccount"
      @submitted="handleSubmitted"
    />

    <FundAccountOperateDrawer
      v-model:visible="fundDrawerVisible"
      :operate-type="operateType"
      :row-data="editingFund"
      :ledger-options="ledgerOptions"
      @submitted="handleSubmitted"
    />
  </div>
</template>

<style scoped>
.finance-basic {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.stat-card__body {
  display: flex;
  align-items: center;
  gap: 14px;
}

.stat-card__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  font-size: 21px;
  border-radius: 12px;
}

.stat-card__text {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.stat-card__label {
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.5);
}

.stat-card__value {
  font-size: 22px;
  font-weight: 600;
  line-height: 1.3;
}

.stat-card__hint {
  overflow: hidden;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.38);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.readonly-tip {
  border-radius: 10px;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-bottom: 14px;
}

.tab-group {
  display: flex;
  gap: 4px;
  padding: 3px;
  background-color: rgb(var(--base-text-color) / 0.05);
  border-radius: 10px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.55);
  cursor: pointer;
  background: transparent;
  border: none;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.tab-btn:hover {
  color: rgb(var(--primary-color));
}

.tab-btn--active {
  font-weight: 600;
  color: rgb(var(--primary-color));
  background-color: rgb(var(--primary-color) / 0.12);
}

.tab-count {
  min-width: 20px;
  padding: 0 6px;
  font-size: 11px;
  font-weight: 500;
  line-height: 18px;
  text-align: center;
  background-color: rgb(var(--base-text-color) / 0.08);
  border-radius: 999px;
}

.filter-group {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.data-table {
  margin-top: 2px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 40px 0;
  font-size: 32px;
  color: rgb(var(--base-text-color) / 0.25);
}

.empty-state p {
  margin: 0;
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.45);
}

:deep(.mini-badge) {
  flex-shrink: 0;
  padding: 0 6px;
  font-size: 11px;
  line-height: 18px;
  color: rgb(var(--primary-color));
  background-color: rgb(var(--primary-color) / 0.1);
  border-radius: 4px;
}

:deep(.mini-badge--custom) {
  color: rgb(var(--base-text-color) / 0.5);
  background-color: rgb(var(--base-text-color) / 0.06);
}

:deep(.soft-tag) {
  display: inline-flex;
  align-items: center;
  padding: 1px 8px;
  font-size: 12px;
  line-height: 20px;
  border-radius: 6px;
}

:deep(.soft-tag--warning) {
  color: #fa8c16;
  background-color: rgb(250 140 22 / 14%);
}

:deep(.soft-tag--muted) {
  color: rgb(var(--base-text-color) / 0.45);
  background-color: rgb(var(--base-text-color) / 0.06);
}

@media (max-width: 640px) {
  .toolbar {
    align-items: stretch;
  }

  .tab-group {
    width: 100%;
  }

  .tab-btn {
    flex: 1;
    justify-content: center;
  }

  .filter-group {
    width: 100%;
  }

  .filter-group :deep(.n-input),
  .filter-group :deep(.n-select) {
    flex: 1;
    width: auto;
    min-width: 110px;
  }
}
</style>
