import { computed, ref } from 'vue';
import { fetchGetFinFundAccounts } from '@/service/api/fin-basic';

/**
 * 资金账户下拉选项 Hook。
 *
 * 收款 / 付款都要回答「钱进了 / 出了哪个账户」，系统据此自动生成
 * 「借：银行存款 / 贷：主营业务收入」这类借贷分录。
 *
 * 选项一次拉取后在财务模块内共享，避免各页面 + 抽屉重复请求。
 */
const fundAccounts = ref<Api.FinBasic.FundAccount[]>([]);
const loading = ref(false);

export function useFundAccounts() {
  /** 启用中的资金账户（停用的账户不能再用于新业务） */
  const enabledFundAccounts = computed(() => fundAccounts.value.filter(item => item.enabled));

  /** 下拉选项：现金 / 银行 + 开户行 + 尾号 */
  const fundAccountOptions = computed(() =>
    enabledFundAccounts.value.map(item => ({
      label: formatFundAccountLabel(item),
      value: item.id
    }))
  );

  /** 是否已配置可用的资金账户（未配置时说明财务还没初始化） */
  const hasFundAccounts = computed(() => enabledFundAccounts.value.length > 0);

  /** 只有一个账户时自动选中，减少录入成本 */
  const defaultFundAccountId = computed(() => (enabledFundAccounts.value.length === 1 ? enabledFundAccounts.value[0].id : null));

  /**
   * 加载资金账户下拉选项
   *
   * @param force 强制重新拉取（默认已有数据时直接复用）
   */
  async function loadFundAccounts(force = false) {
    if (!force && fundAccounts.value.length) {
      return fundAccounts.value;
    }

    loading.value = true;

    try {
      const { data, error } = await fetchGetFinFundAccounts();

      if (!error && data) {
        fundAccounts.value = data.records ?? [];
      }
    } finally {
      loading.value = false;
    }

    return fundAccounts.value;
  }

  /** 根据账户ID取展示名 */
  function getFundAccountName(id?: number | null) {
    if (!id) return '';

    const account = fundAccounts.value.find(item => item.id === id);

    return account ? formatFundAccountLabel(account) : '';
  }

  return {
    fundAccounts: enabledFundAccounts,
    fundAccountOptions,
    fundAccountLoading: loading,
    hasFundAccounts,
    defaultFundAccountId,
    loadFundAccounts,
    getFundAccountName
  };
}

/** 账户展示名：中国银行 · 尾号 1234 */
export function formatFundAccountLabel(account: Api.FinBasic.FundAccount) {
  const suffix = account.bankName ? ` · ${account.bankName}` : '';
  const last4 = account.accountLast4 ? ` 尾号 ${account.accountLast4}` : '';

  return `${account.name}${suffix}${last4}`;
}
