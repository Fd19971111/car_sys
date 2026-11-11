import { ref } from 'vue';
import { fetchGetFinReportFlows, fetchGetFinReportOverview } from '@/service/api';

/** 报表页数据 Hook：总览一次性拉取，资金流水按筛选条件分页拉取 */
export function useReportData() {
  const overviewLoading = ref(false);
  const overview = ref<Api.FinReport.Overview | null>(null);

  const flowsLoading = ref(false);
  const flows = ref<Api.FinReport.FlowItem[]>([]);
  const flowTotal = ref(0);
  const flowSummary = ref<Api.FinReport.FlowListResult['summary']>({ income: 0, expense: 0, net: 0 });

  /** 加载报表总览 */
  async function loadOverview(year: number) {
    overviewLoading.value = true;

    try {
      const { data, error } = await fetchGetFinReportOverview({ year });

      if (!error && data) {
        overview.value = data;
      }
    } finally {
      overviewLoading.value = false;
    }
  }

  /** 加载资金流水明细 */
  async function loadFlows(params: Api.FinReport.FlowSearchParams) {
    flowsLoading.value = true;

    try {
      const { data, error } = await fetchGetFinReportFlows(params);

      if (!error && data) {
        flows.value = data.records ?? [];
        flowTotal.value = data.total ?? 0;
        flowSummary.value = data.summary ?? { income: 0, expense: 0, net: 0 };
      }
    } finally {
      flowsLoading.value = false;
    }
  }

  return {
    overviewLoading,
    overview,
    flowsLoading,
    flows,
    flowTotal,
    flowSummary,
    loadOverview,
    loadFlows
  };
}
