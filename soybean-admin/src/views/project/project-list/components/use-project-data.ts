import { ref } from 'vue';
import {
  fetchBatchDeleteProject,
  fetchCreateProject,
  fetchDeleteProject,
  fetchGetProjectDetail,
  fetchGetProjectList,
  fetchGetProjectOptions,
  fetchGetProjectStats,
  fetchUpdateProject
} from '@/service/api/project';

/** 项目模块数据 Hook：全部数据走后端 /project 接口（GoFrame + MySQL） */
export function useProjectData() {
  /** 列表加载状态 */
  const loading = ref(false);
  /** 当前页项目数据 */
  const records = ref<Api.Project.Item[]>([]);
  /** 项目总数（过滤后） */
  const total = ref(0);
  /** 概览统计 */
  const stats = ref<Api.Project.Stats>({
    total: 0,
    pending: 0,
    building: 0,
    finished: 0,
    settled: 0,
    totalAmount: 0
  });
  /** 表单下拉选项（客户 / 负责人） */
  const options = ref<Api.Project.Options>({ customers: [], managers: [] });

  /** 获取项目分页列表 */
  async function loadList(params: Api.Project.SearchParams) {
    loading.value = true;

    try {
      const { data, error } = await fetchGetProjectList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        status: params.status ?? null,
        customerId: params.customerId ?? null
      });

      if (error || !data) {
        return;
      }

      records.value = data.records ?? [];
      total.value = data.total ?? 0;
    } finally {
      loading.value = false;
    }
  }

  /** 获取概览统计 */
  async function loadStats() {
    const { data, error } = await fetchGetProjectStats();

    if (!error && data) {
      stats.value = data;
    }
  }

  /** 获取表单下拉选项（客户 / 负责人） */
  async function loadOptions() {
    const { data, error } = await fetchGetProjectOptions();

    if (!error && data) {
      options.value = {
        customers: data.customers ?? [],
        managers: data.managers ?? []
      };
    }
  }

  /** 获取项目详情（含客户信息与进度时间轴），失败返回 null */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetProjectDetail(id);

    return error || !data ? null : data;
  }

  /** 新增项目，返回是否成功 */
  async function createProject(form: Api.Project.FormData) {
    const { error } = await fetchCreateProject(form);

    return !error;
  }

  /** 更新项目，返回是否成功 */
  async function updateProject(id: number, form: Api.Project.FormData) {
    const { error } = await fetchUpdateProject({ ...form, id });

    return !error;
  }

  /** 删除项目（支持批量），返回是否成功 */
  async function deleteProjects(ids: number[]) {
    if (!ids.length) {
      return false;
    }

    const { error } = ids.length > 1 ? await fetchBatchDeleteProject(ids) : await fetchDeleteProject(ids[0]);

    return !error;
  }

  return {
    loading,
    records,
    total,
    stats,
    options,
    loadList,
    loadStats,
    loadOptions,
    getDetail,
    createProject,
    updateProject,
    deleteProjects
  };
}
