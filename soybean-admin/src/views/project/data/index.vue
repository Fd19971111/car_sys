<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NImage, NPopconfirm, NTag } from 'naive-ui';
import { fetchGetProjectSelectOptions } from '@/service/api';
import { formatDate, formatMoney } from '../project-list/components/project-meta';
import ChangeOperateDrawer from './components/change-operate-drawer.vue';
import ContractOperateDrawer from './components/contract-operate-drawer.vue';
import SiteLogOperateDrawer from './components/site-log-operate-drawer.vue';
import { useProjectData } from './components/use-project-data';

defineOptions({
  name: 'ProjectData'
});

const {
  logLoading,
  logs,
  logTotal,
  loadLogs,
  deleteLog,
  changeLoading,
  changes,
  changeTotal,
  addTotal,
  reduceTotal,
  loadChanges,
  deleteChange,
  contractLoading,
  contracts,
  contractTotal,
  loadContracts,
  deleteContract
} = useProjectData();

/** 项目下拉 */
const projectOptions = ref<Api.Project.SelectOption[]>([]);
const projectLoading = ref(false);
const projectId = ref<number | null>(null);
const selectedProject = computed(() => projectOptions.value.find(item => item.value === projectId.value));

// ===================== 施工日志 =====================
const logPage = ref(1);
const logPageSize = ref(10);
const logPagination = computed(() => ({
  page: logPage.value,
  pageSize: logPageSize.value,
  itemCount: logTotal.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: (info: { itemCount?: number }) => `共 ${info.itemCount ?? 0} 条日志`,
  onUpdatePage: (page: number) => {
    logPage.value = page;
    void loadLogData();
  },
  onUpdatePageSize: (size: number) => {
    logPageSize.value = size;
    logPage.value = 1;
    void loadLogData();
  }
}));

const logDrawerVisible = ref(false);
const logOperateType = ref<NaiveUI.TableOperateType>('add');
const editingLog = ref<Api.ProjectData.SiteLog | null>(null);

// ===================== 变更记录 =====================
const changeType = ref<Api.ProjectData.ChangeType | null>(null);
const changeKeyword = ref('');
const changePagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: (info: { itemCount?: number }) => `共 ${info.itemCount ?? 0} 条变更`,
  onUpdatePage(page: number) {
    changePagination.page = page;
    void loadChangeData();
  },
  onUpdatePageSize(pageSize: number) {
    changePagination.pageSize = pageSize;
    changePagination.page = 1;
    void loadChangeData();
  }
});

const changeDrawerVisible = ref(false);
const changeOperateType = ref<NaiveUI.TableOperateType>('add');
const editingChange = ref<Api.ProjectData.Change | null>(null);

const contractPagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: (info: { itemCount?: number }) => `共 ${info.itemCount ?? 0} 份合同`,
  onUpdatePage(page: number) {
    contractPagination.page = page;
    void loadContractData();
  },
  onUpdatePageSize(pageSize: number) {
    contractPagination.pageSize = pageSize;
    contractPagination.page = 1;
    void loadContractData();
  }
});

const contractDrawerVisible = ref(false);
const contractOperateType = ref<NaiveUI.TableOperateType>('add');
const editingContract = ref<Api.ProjectData.Contract | null>(null);

const contractColumns: DataTableColumns<Api.ProjectData.Contract> = [
  { key: 'name', title: '合同名称', minWidth: 180, ellipsis: { tooltip: true } },
  { key: 'contractNo', title: '合同编号', width: 150, render: row => row.contractNo || '—' },
  {
    key: 'amount',
    title: '合同金额',
    width: 140,
    align: 'right',
    render: row => h('span', { class: 'font-mono font-medium' }, formatMoney(row.amount))
  },
  { key: 'signedAt', title: '签订日期', width: 120 },
  {
    key: 'attachments',
    title: '合同附件',
    minWidth: 200,
    render: row =>
      row.attachments?.length
        ? h(
            'div',
            { class: 'flex flex-col items-start gap-4px' },
            row.attachments.map((file, index) =>
              h('a', {
                key: `${file.path}-${index}`,
                href: file.url || file.path,
                target: '_blank',
                rel: 'noopener noreferrer',
                class: 'text-primary truncate max-w-220px',
                title: file.name,
                innerText: file.name
              })
            )
          )
        : h('span', { class: 'text-gray-400' }, '—')
  },
  { key: 'remark', title: '备注', minWidth: 160, ellipsis: { tooltip: true }, render: row => row.remark || '—' },
  {
    key: 'operate',
    title: '操作',
    width: 110,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-2px' }, [
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => openContractEdit(row) }, { default: () => '编辑' }),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleContractDelete(row) },
          {
            trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除这份合同吗？'
          }
        )
      ])
  }
];

const netChange = computed(() => addTotal.value - reduceTotal.value);

const changeColumns: DataTableColumns<Api.ProjectData.Change> = [
  {
    key: 'changeDate',
    title: '变更日期',
    width: 110,
    render: row => h('span', { class: 'font-mono' }, row.changeDate)
  },
  {
    key: 'changeType',
    title: '类型',
    width: 84,
    render: row =>
      h(
        NTag,
        { size: 'small', bordered: false, type: row.changeType === 'ADD' ? 'success' : 'warning' },
        { default: () => (row.changeType === 'ADD' ? '增项' : '减项') }
      )
  },
  {
    key: 'title',
    title: '变更事项',
    minWidth: 180,
    ellipsis: { tooltip: true }
  },
  {
    key: 'amount',
    title: '金额',
    width: 130,
    align: 'right',
    render: row =>
      h(
        'span',
        { class: `font-mono font-medium ${row.changeType === 'ADD' ? 'text-success' : 'text-warning'}` },
        `${row.changeType === 'ADD' ? '+' : '-'}${formatMoney(row.amount)}`
      )
  },
  {
    key: 'reason',
    title: '变更原因',
    minWidth: 160,
    ellipsis: { tooltip: true },
    render: row => row.reason || '-'
  },
  {
    key: 'images',
    title: '附件',
    width: 120,
    render: row =>
      row.imageUrls?.length
        ? h(
            'div',
            { class: 'flex gap-4px' },
            row.imageUrls
              .slice(0, 3)
              .map((url, index) =>
                h(NImage, { key: index, src: url, width: 32, height: 32, objectFit: 'cover', class: 'rounded-4px' })
              )
          )
        : h('span', { class: 'text-gray-400' }, '—')
  },
  {
    key: 'operate',
    title: '操作',
    width: 110,
    align: 'center',
    fixed: 'right',
    render: row =>
      h('div', { class: 'flex-center gap-2px' }, [
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => openChangeEdit(row) }, { default: () => '编辑' }),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleChangeDelete(row) },
          {
            trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除该变更记录吗？'
          }
        )
      ])
  }
];

// ===================== 数据加载 =====================
async function loadProjectOptions() {
  projectLoading.value = true;

  try {
    const { data, error } = await fetchGetProjectSelectOptions();

    if (!error && data) {
      projectOptions.value = data.projects ?? [];

      if (!projectId.value && projectOptions.value.length) {
        projectId.value = projectOptions.value[0].value;
      }
    }
  } finally {
    projectLoading.value = false;
  }
}

async function loadLogData() {
  if (!projectId.value) return;

  await loadLogs({
    projectId: projectId.value,
    current: logPage.value,
    size: logPageSize.value
  });
}

async function loadChangeData() {
  if (!projectId.value) return;

  await loadChanges({
    projectId: projectId.value,
    changeType: changeType.value,
    keyword: changeKeyword.value,
    current: changePagination.page,
    size: changePagination.pageSize
  });
  changePagination.itemCount = changeTotal.value;
}

async function loadContractData() {
  if (!projectId.value) return;
  await loadContracts({
    projectId: projectId.value,
    current: contractPagination.page,
    size: contractPagination.pageSize
  });
  contractPagination.itemCount = contractTotal.value;
}

async function loadAll() {
  if (!projectId.value) return;

  await Promise.all([loadLogData(), loadChangeData(), loadContractData()]);
}

function resetPaging() {
  logPage.value = 1;
  changePagination.page = 1;
  contractPagination.page = 1;
}

watch(projectId, async () => {
  changeType.value = null;
  changeKeyword.value = '';
  resetPaging();
  await loadAll();
});

// ===================== 日志操作 =====================
function openLogAdd() {
  logOperateType.value = 'add';
  editingLog.value = null;
  logDrawerVisible.value = true;
}

function openLogEdit(row: Api.ProjectData.SiteLog) {
  logOperateType.value = 'edit';
  editingLog.value = row;
  logDrawerVisible.value = true;
}

async function handleLogDelete(row: Api.ProjectData.SiteLog) {
  const success = await deleteLog(row.id);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');
  await loadLogData();
}

async function handleLogSubmitted() {
  await loadLogData();
}

// ===================== 变更操作 =====================
function openChangeAdd() {
  changeOperateType.value = 'add';
  editingChange.value = null;
  changeDrawerVisible.value = true;
}

function openChangeEdit(row: Api.ProjectData.Change) {
  changeOperateType.value = 'edit';
  editingChange.value = row;
  changeDrawerVisible.value = true;
}

async function handleChangeDelete(row: Api.ProjectData.Change) {
  const success = await deleteChange(row.id);

  if (!success) {
    window.$message?.error('删除失败，请稍后重试');

    return;
  }

  window.$message?.success('删除成功');
  await loadChangeData();
}

async function handleChangeSubmitted() {
  await loadChangeData();
}

async function handleChangeSearch() {
  changePagination.page = 1;
  await loadChangeData();
}

function openContractAdd() {
  contractOperateType.value = 'add';
  editingContract.value = null;
  contractDrawerVisible.value = true;
}

function openContractEdit(row: Api.ProjectData.Contract) {
  contractOperateType.value = 'edit';
  editingContract.value = row;
  contractDrawerVisible.value = true;
}

async function handleContractDelete(row: Api.ProjectData.Contract) {
  const success = await deleteContract(row.id);
  if (!success) {
    window.$message?.error('删除失败，请稍后重试');
    return;
  }
  window.$message?.success('合同已删除');
  await loadContractData();
}

async function handleContractSubmitted() {
  await loadContractData();
}

onMounted(async () => {
  await loadProjectOptions();
  await loadAll();
});
</script>

<template>
  <div class="p-16px">
    <NSpace vertical :size="16">
      <!-- 项目选择 -->
      <NCard :bordered="false" class="card-wrapper" size="small">
        <div class="flex flex-wrap items-center gap-12px">
          <NSelect
            v-model:value="projectId"
            class="w-360px max-w-full"
            :options="projectOptions"
            :loading="projectLoading"
            filterable
            placeholder="请选择项目"
          />
          <template v-if="selectedProject">
            <NTag size="small" :bordered="false">合同 {{ formatMoney(selectedProject.amount) }}</NTag>
            <span class="text-12px text-gray-400">
              施工日志照片、增项/减项都按项目归档；施工日志照片会同步到业主端小程序
            </span>
          </template>
          <span v-else class="text-12px text-gray-400">请选择项目后登记资料</span>
        </div>
      </NCard>

      <NEmpty v-if="!projectId" description="请先选择项目" class="py-60px" />

      <NCard v-else :bordered="false" class="card-wrapper" size="small">
        <NTabs type="line" animated>
          <!-- 施工日志 -->
          <NTabPane name="log" :tab="`施工日志（${logTotal}）`">
            <div class="mb-12px flex flex-wrap items-center gap-8px">
              <span class="text-12px text-gray-400">现场照片会随日志一起展示给业主</span>
              <div class="ml-auto flex items-center gap-8px">
                <NButton size="small" type="primary" ghost @click="openLogAdd">
                  <template #icon>
                    <SvgIcon icon="mdi:camera-plus-outline" />
                  </template>
                  新增日志
                </NButton>
                <NButton size="small" :loading="logLoading" @click="loadLogData">
                  <template #icon>
                    <SvgIcon icon="mdi:refresh" />
                  </template>
                  刷新
                </NButton>
              </div>
            </div>

            <NSpin :show="logLoading">
              <NEmpty
                v-if="!logLoading && !logs.length"
                description="还没有施工日志，点「新增日志」记录现场照片"
                class="py-40px"
              />

              <div v-else class="log-list">
                <article v-for="log in logs" :key="log.id" class="log-item">
                  <div class="log-item__head">
                    <div class="flex-col gap-2px">
                      <span class="text-13px font-medium">{{ log.logDate || formatDate(log.createdAt) }}</span>
                      <span class="text-11px text-gray-400">记录人：{{ log.userName || '-' }} · 记录于 {{ formatDate(log.createdAt) }}</span>
                    </div>
                    <div class="flex items-center gap-4px">
                      <NButton size="tiny" quaternary @click="openLogEdit(log)">编辑</NButton>
                      <NPopconfirm @positive-click="handleLogDelete(log)">
                        <template #trigger>
                          <NButton size="tiny" quaternary type="error">删除</NButton>
                        </template>
                        确认删除这条施工日志吗？业主端将同步不可见。
                      </NPopconfirm>
                    </div>
                  </div>
                  <p class="log-item__content">{{ log.content }}</p>
                  <div v-if="log.imageUrls?.length" class="mt-8px flex flex-wrap gap-8px">
                    <NImage
                      v-for="(url, index) in log.imageUrls"
                      :key="index"
                      :src="url"
                      :width="96"
                      :height="96"
                      object-fit="cover"
                      class="rounded-6px"
                    />
                  </div>
                </article>
              </div>
            </NSpin>

            <div v-if="logTotal > 0" class="mt-12px flex justify-end">
              <NPagination v-bind="logPagination" />
            </div>
          </NTabPane>

          <!-- 变更记录 -->
          <NTabPane name="change" :tab="`变更记录（${changeTotal}）`">
            <div class="mb-12px flex flex-wrap gap-12px">
              <NCard :bordered="false" size="small" class="flex-1">
                <span class="text-12px text-gray-400">增项合计</span>
                <div class="text-16px font-medium text-success">+{{ formatMoney(addTotal) }}</div>
              </NCard>
              <NCard :bordered="false" size="small" class="flex-1">
                <span class="text-12px text-gray-400">减项合计</span>
                <div class="text-16px font-medium text-warning">-{{ formatMoney(reduceTotal) }}</div>
              </NCard>
              <NCard :bordered="false" size="small" class="flex-1">
                <span class="text-12px text-gray-400">净变更</span>
                <div class="text-16px font-medium">{{ netChange >= 0 ? '+' : '' }}{{ formatMoney(netChange) }}</div>
              </NCard>
            </div>

            <div class="mb-12px flex flex-wrap items-center gap-8px">
              <NSelect
                v-model:value="changeType"
                class="w-120px flex-shrink-0"
                :options="[
                  { label: '增项', value: 'ADD' },
                  { label: '减项', value: 'REDUCE' }
                ]"
                placeholder="变更类型"
                clearable
              />
              <NInput
                v-model:value="changeKeyword"
                class="w-220px min-w-140px flex-shrink"
                placeholder="搜索变更事项 / 原因"
                clearable
                @keyup.enter="handleChangeSearch"
              >
                <template #prefix>
                  <SvgIcon icon="mdi:magnify" class="text-gray-400" />
                </template>
              </NInput>
              <NButton size="small" type="primary" ghost @click="handleChangeSearch">
                <template #icon>
                  <SvgIcon icon="mdi:magnify" />
                </template>
                查询
              </NButton>
              <div class="ml-auto flex items-center gap-8px">
                <NButton size="small" type="primary" ghost @click="openChangeAdd">
                  <template #icon>
                    <SvgIcon icon="mdi:plus" />
                  </template>
                  新增变更
                </NButton>
                <NButton size="small" :loading="changeLoading" @click="loadChangeData">
                  <template #icon>
                    <SvgIcon icon="mdi:refresh" />
                  </template>
                  刷新
                </NButton>
              </div>
            </div>

            <NDataTable
              :columns="changeColumns"
              :data="changes"
              :loading="changeLoading"
              :pagination="changePagination"
              :row-key="row => row.id"
              :scroll-x="940"
              remote
              size="small"
              striped
            />
          </NTabPane>

          <NTabPane name="contract" :tab="`合同（${contractTotal}）`">
            <div class="mb-12px flex items-center gap-8px">
              <span class="text-12px text-gray-400">合同文件按项目归档</span>
              <div class="ml-auto flex items-center gap-8px">
                <NButton size="small" type="primary" ghost @click="openContractAdd">
                  <template #icon><SvgIcon icon="mdi:plus" /></template>
                  新增合同
                </NButton>
                <NButton size="small" :loading="contractLoading" @click="loadContractData">
                  <template #icon><SvgIcon icon="mdi:refresh" /></template>
                  刷新
                </NButton>
              </div>
            </div>
            <NDataTable
              :columns="contractColumns"
              :data="contracts"
              :loading="contractLoading"
              :pagination="contractPagination"
              :row-key="row => row.id"
              :scroll-x="980"
              remote
              size="small"
              striped
            />
          </NTabPane>
        </NTabs>
      </NCard>
    </NSpace>

    <!-- 抽屉 -->
    <SiteLogOperateDrawer
      v-model:visible="logDrawerVisible"
      :project-id="projectId"
      :operate-type="logOperateType"
      :row-data="editingLog"
      @submitted="handleLogSubmitted"
    />

    <ChangeOperateDrawer
      v-model:visible="changeDrawerVisible"
      :project-id="projectId"
      :operate-type="changeOperateType"
      :row-data="editingChange"
      @submitted="handleChangeSubmitted"
    />

    <ContractOperateDrawer
      v-model:visible="contractDrawerVisible"
      :project-id="projectId"
      :operate-type="contractOperateType"
      :row-data="editingContract"
      @submitted="handleContractSubmitted"
    />
  </div>
</template>

<style scoped>
.card-wrapper {
  border-radius: 12px;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-item {
  padding: 12px 14px;
  background-color: rgb(var(--base-text-color) / 0.02);
  border: 1px solid rgb(var(--base-text-color) / 0.06);
  border-radius: 10px;
}

.log-item__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.log-item__content {
  margin: 8px 0 0;
  font-size: 13px;
  line-height: 1.7;
  color: rgb(var(--base-text-color) / 0.8);
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
