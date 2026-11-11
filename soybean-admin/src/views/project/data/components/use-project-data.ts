import { ref } from 'vue';
import {
  fetchCreateChange,
  fetchCreateContract,
  fetchCreateMaterial,
  fetchCreateMaterialUsage,
  fetchCreateSiteLog,
  fetchDeleteChange,
  fetchDeleteContract,
  fetchDeleteMaterial,
  fetchDeleteMaterialUsage,
  fetchDeleteSiteLog,
  fetchGetChangeList,
  fetchGetContractList,
  fetchGetMaterialList,
  fetchGetMaterialUsageList,
  fetchGetSiteLogList,
  fetchUpdateChange,
  fetchUpdateContract,
  fetchUpdateMaterial,
  fetchUpdateSiteLog
} from '@/service/api';

/** 项目资料（材料 / 领用 / 施工日志 / 变更）数据 Hook */
export function useProjectData() {
  const contractLoading = ref(false);
  const contracts = ref<Api.ProjectData.Contract[]>([]);
  const contractTotal = ref(0);

  async function loadContracts(params: Api.ProjectData.ContractListParams) {
    contractLoading.value = true;
    try {
      const { data, error } = await fetchGetContractList(params);
      if (!error && data) {
        contracts.value = data.records ?? [];
        contractTotal.value = data.total ?? 0;
      }
    } finally {
      contractLoading.value = false;
    }
  }

  async function createContract(form: Api.ProjectData.ContractForm) {
    const { error } = await fetchCreateContract(form);
    return !error;
  }

  async function updateContract(id: number, form: Api.ProjectData.ContractForm) {
    const { error } = await fetchUpdateContract({ ...form, id });
    return !error;
  }

  async function deleteContract(id: number) {
    const { error } = await fetchDeleteContract(id);
    return !error;
  }

  // ===================== 材料清单 =====================
  const materialLoading = ref(false);
  const materials = ref<Api.ProjectData.Material[]>([]);
  const materialTotal = ref(0);

  async function loadMaterials(params: Api.ProjectData.MaterialSearchParams) {
    materialLoading.value = true;

    try {
      const { data, error } = await fetchGetMaterialList({
        projectId: params.projectId,
        keyword: params.keyword?.trim() || null,
        current: params.current || 1,
        size: params.size || 10
      });

      if (!error && data) {
        materials.value = data.records ?? [];
        materialTotal.value = data.total ?? 0;
      }
    } finally {
      materialLoading.value = false;
    }
  }

  async function createMaterial(form: Api.ProjectData.MaterialForm) {
    const { error } = await fetchCreateMaterial(form);

    return !error;
  }

  async function updateMaterial(id: number, form: Api.ProjectData.MaterialForm) {
    const { error } = await fetchUpdateMaterial({ ...form, id });

    return !error;
  }

  async function deleteMaterial(id: number) {
    const { error } = await fetchDeleteMaterial(id);

    return !error;
  }

  // ===================== 材料领用 =====================
  const usageLoading = ref(false);
  const usages = ref<Api.ProjectData.Usage[]>([]);
  const usageTotal = ref(0);

  async function loadUsages(params: Api.ProjectData.UsageSearchParams) {
    usageLoading.value = true;

    try {
      const { data, error } = await fetchGetMaterialUsageList({
        projectId: params.projectId,
        materialId: params.materialId ?? null,
        current: params.current || 1,
        size: params.size || 10
      });

      if (!error && data) {
        usages.value = data.records ?? [];
        usageTotal.value = data.total ?? 0;
      }
    } finally {
      usageLoading.value = false;
    }
  }

  async function createUsage(form: Api.ProjectData.UsageForm) {
    const { error } = await fetchCreateMaterialUsage(form);

    return !error;
  }

  async function deleteUsage(id: number) {
    const { error } = await fetchDeleteMaterialUsage(id);

    return !error;
  }

  // ===================== 施工日志 =====================
  const logLoading = ref(false);
  const logs = ref<Api.ProjectData.SiteLog[]>([]);
  const logTotal = ref(0);

  async function loadLogs(params: Api.ProjectData.SiteLogSearchParams) {
    logLoading.value = true;

    try {
      const { data, error } = await fetchGetSiteLogList({
        projectId: params.projectId,
        current: params.current || 1,
        size: params.size || 10
      });

      if (!error && data) {
        logs.value = data.records ?? [];
        logTotal.value = data.total ?? 0;
      }
    } finally {
      logLoading.value = false;
    }
  }

  async function createLog(form: Api.ProjectData.SiteLogForm) {
    const { error } = await fetchCreateSiteLog(form);

    return !error;
  }

  async function updateLog(id: number, form: Api.ProjectData.SiteLogForm) {
    const { error } = await fetchUpdateSiteLog({ ...form, id });

    return !error;
  }

  async function deleteLog(id: number) {
    const { error } = await fetchDeleteSiteLog(id);

    return !error;
  }

  // ===================== 变更记录 =====================
  const changeLoading = ref(false);
  const changes = ref<Api.ProjectData.Change[]>([]);
  const changeTotal = ref(0);
  const addTotal = ref(0);
  const reduceTotal = ref(0);

  async function loadChanges(params: Api.ProjectData.ChangeSearchParams) {
    changeLoading.value = true;

    try {
      const { data, error } = await fetchGetChangeList({
        projectId: params.projectId,
        changeType: params.changeType ?? null,
        keyword: params.keyword?.trim() || null,
        current: params.current || 1,
        size: params.size || 10
      });

      if (!error && data) {
        changes.value = data.records ?? [];
        changeTotal.value = data.total ?? 0;
        addTotal.value = data.addTotal ?? 0;
        reduceTotal.value = data.reduceTotal ?? 0;
      }
    } finally {
      changeLoading.value = false;
    }
  }

  async function createChange(form: Api.ProjectData.ChangeForm) {
    const { error } = await fetchCreateChange(form);

    return !error;
  }

  async function updateChange(id: number, form: Api.ProjectData.ChangeForm) {
    const { error } = await fetchUpdateChange({ ...form, id });

    return !error;
  }

  async function deleteChange(id: number) {
    const { error } = await fetchDeleteChange(id);

    return !error;
  }

  return {
    contractLoading,
    contracts,
    contractTotal,
    loadContracts,
    createContract,
    updateContract,
    deleteContract,
    materialLoading,
    materials,
    materialTotal,
    loadMaterials,
    createMaterial,
    updateMaterial,
    deleteMaterial,

    usageLoading,
    usages,
    usageTotal,
    loadUsages,
    createUsage,
    deleteUsage,

    logLoading,
    logs,
    logTotal,
    loadLogs,
    createLog,
    updateLog,
    deleteLog,

    changeLoading,
    changes,
    changeTotal,
    addTotal,
    reduceTotal,
    loadChanges,
    createChange,
    updateChange,
    deleteChange
  };
}
