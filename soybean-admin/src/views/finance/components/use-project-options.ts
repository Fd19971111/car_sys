import { computed, ref } from 'vue';
import { fetchGetProjectSelectOptions } from '@/service/api/project';

/**
 * 项目下拉选项 Hook。
 *
 * 报价 / 收款两个模块都要选择「所属项目」，选项一次拉取后在模块内共享，
 * 避免两个页面 + 两个抽屉各自重复请求。
 */
const projectOptions = ref<Api.Project.SelectOption[]>([]);
const loading = ref(false);

/** 项目ID → 选项 的映射，便于反查合同金额等信息 */
export const projectOptionMap = computed(() => {
  const map = new Map<number, Api.Project.SelectOption>();

  projectOptions.value.forEach(item => map.set(item.value, item));

  return map;
});

export function useProjectOptions() {
  /**
   * 加载项目下拉选项
   *
   * @param force 强制重新拉取（默认已有数据时直接复用）
   */
  async function loadProjectOptions(force = false) {
    if (!force && projectOptions.value.length) {
      return projectOptions.value;
    }

    loading.value = true;

    try {
      const { data, error } = await fetchGetProjectSelectOptions();

      if (!error && data) {
        projectOptions.value = data.projects ?? [];
      }
    } finally {
      loading.value = false;
    }

    return projectOptions.value;
  }

  /** 根据项目ID取项目名称 */
  function getProjectName(projectId?: number | null) {
    if (!projectId) return '';

    return projectOptionMap.value.get(projectId)?.label ?? '';
  }

  return {
    projectOptions,
    projectOptionMap,
    projectLoading: loading,
    loadProjectOptions,
    getProjectName
  };
}
