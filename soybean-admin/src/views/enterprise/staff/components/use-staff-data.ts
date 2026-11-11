import { ref } from 'vue';
import {
  fetchBatchDeleteStaff,
  fetchCreateStaff,
  fetchDeleteStaff,
  fetchGetDepartmentOptions,
  fetchGetStaffDetail,
  fetchGetStaffList,
  fetchUpdateStaff
} from '@/service/api/staff';

/** 人员模块数据 Hook */
export function useStaffData() {
  const loading = ref(false);
  const records = ref<Api.Staff.Item[]>([]);
  const total = ref(0);
  const departmentOptions = ref<Api.Staff.DepartmentOption[]>([]);

  /** 获取人员分页列表 */
  async function loadList(params: Api.Staff.SearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetStaffList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null,
        departmentId: params.departmentId ?? null,
        status: params.status ?? null
      });
      if (error || !data) return;
      records.value = data.records ?? [];
      total.value = data.total ?? 0;
    } finally {
      loading.value = false;
    }
  }

  /** 获取人员详情 */
  async function getDetail(id: number) {
    const { data, error } = await fetchGetStaffDetail(id);
    return error || !data ? null : data;
  }

  /** 新增人员 */
  async function createStaff(form: Api.Staff.FormData) {
    const { error } = await fetchCreateStaff(form);
    return !error;
  }

  /** 更新人员 */
  async function updateStaff(id: number, form: Api.Staff.FormData) {
    const { error } = await fetchUpdateStaff({ ...form, id });
    return !error;
  }

  /** 删除人员（支持批量） */
  async function deleteStaff(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeleteStaff(ids) : await fetchDeleteStaff(ids[0]);
    return !error;
  }

  /** 加载部门下拉选项 */
  async function loadDepartmentOptions() {
    const { data, error } = await fetchGetDepartmentOptions();
    if (!error && data) departmentOptions.value = data.items ?? [];
  }

  return {
    loading,
    records,
    total,
    departmentOptions,
    loadList,
    getDetail,
    createStaff,
    updateStaff,
    deleteStaff,
    loadDepartmentOptions
  };
}