import { ref } from 'vue';
import {
  fetchBatchDeleteDepartment,
  fetchCreateDepartment,
  fetchDeleteDepartment,
  fetchGetDepartmentDetail,
  fetchGetDepartmentList,
  fetchUpdateDepartment
} from '@/service/api/department';

/** 部门模块数据 Hook */
export function useDepartmentData() {
  const loading = ref(false);
  const records = ref<Api.Department.Item[]>([]);
  const total = ref(0);

  /** 获取部门分页列表 */
  async function loadList(params: Api.Department.SearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetDepartmentList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        status: params.status ?? null
      });
      if (error || !data) return;
      records.value = data.records ?? [];
      total.value = data.total ?? 0;
    } finally {
      loading.value = false;
    }
  }

  /** 获取部门详情 */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetDepartmentDetail(id);
    return error || !data ? null : data;
  }

  /** 新增部门 */
  async function createDepartment(form: Api.Department.FormData) {
    const { error } = await fetchCreateDepartment(form);
    return !error;
  }

  /** 更新部门 */
  async function updateDepartment(id: number, form: Api.Department.FormData) {
    const { error } = await fetchUpdateDepartment({ ...form, id });
    return !error;
  }

  /** 删除部门（支持批量） */
  async function deleteDepartment(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeleteDepartment(ids) : await fetchDeleteDepartment(ids[0]);
    return !error;
  }

  return {
    loading,
    records,
    total,
    loadList,
    getDetail,
    createDepartment,
    updateDepartment,
    deleteDepartment
  };
}