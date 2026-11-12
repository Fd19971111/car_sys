import { ref } from 'vue';
import {
  fetchBatchDeletePricingRule,
  fetchCreatePricingRule,
  fetchDeletePricingRule,
  fetchGetPricingRuleList,
  fetchUpdatePricingRule
} from '@/service/api/pricing-rule';

export function usePricingRuleData() {
  const loading = ref(false);
  const records = ref<Api.PricingRule.Item[]>([]);
  const total = ref(0);

  async function loadList(params: Api.PricingRule.SearchParams) {
    loading.value = true;
    try {
      const { data, error } = await fetchGetPricingRuleList({
        current: params.current || 1,
        size: params.size || 10,
        keyword: params.keyword?.trim() || null
      });
      if (!error && data) {
        records.value = data.records ?? [];
        total.value = data.total ?? 0;
      }
    } finally {
      loading.value = false;
    }
  }

  async function createRule(form: Api.PricingRule.FormData) {
    const { error } = await fetchCreatePricingRule(form);
    return !error;
  }

  async function updateRule(id: number, form: Api.PricingRule.FormData) {
    const { error } = await fetchUpdatePricingRule({ ...form, id });
    return !error;
  }

  async function deleteRules(ids: number[]) {
    if (!ids.length) return false;
    const { error } = ids.length > 1 ? await fetchBatchDeletePricingRule(ids) : await fetchDeletePricingRule(ids[0]);
    return !error;
  }

  return { loading, records, total, loadList, createRule, updateRule, deleteRules };
}
