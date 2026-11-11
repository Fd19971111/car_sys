import dayjs from 'dayjs';

/** 报表页统一使用的展示工具与文案 */

/** 图表配色（与系统主题色保持一致的低饱和企业色） */
export const CHART_COLORS = ['#1677ff', '#18a058', '#f0a020', '#e34d59', '#722ed1', '#13c2c2', '#eb2f96', '#8c8c8c'];

/** 科目类别：代码 + 财务术语 + 大白话解释 */
export interface CategoryMeta {
  code: Api.FinReport.CategoryCode;
  /** 财务术语 */
  label: string;
  /** 给非财务人员看的一句话解释 */
  plain: string;
  color: string;
  bg: string;
  icon: string;
}

export const CATEGORY_META: Record<Api.FinReport.CategoryCode, CategoryMeta> = {
  ASSET: {
    code: 'ASSET',
    label: '资产',
    plain: '公司手里拥有的钱和物',
    color: '#1677ff',
    bg: 'rgba(22, 119, 255, 0.12)',
    icon: 'mdi:bank-outline'
  },
  LIABILITY: {
    code: 'LIABILITY',
    label: '负债',
    plain: '欠供应商、员工或银行的钱',
    color: '#f0a020',
    bg: 'rgba(240, 160, 32, 0.14)',
    icon: 'mdi:hand-coin-outline'
  },
  EQUITY: {
    code: 'EQUITY',
    label: '所有者权益',
    plain: '真正属于老板自己的部分',
    color: '#722ed1',
    bg: 'rgba(114, 46, 209, 0.12)',
    icon: 'mdi:account-key-outline'
  },
  REVENUE: {
    code: 'REVENUE',
    label: '收入',
    plain: '做业务赚进来的钱',
    color: '#18a058',
    bg: 'rgba(24, 160, 88, 0.12)',
    icon: 'mdi:trending-up'
  },
  EXPENSE: {
    code: 'EXPENSE',
    label: '支出',
    plain: '材料、人工、房租等花出去的钱',
    color: '#e34d59',
    bg: 'rgba(227, 77, 89, 0.12)',
    icon: 'mdi:trending-down'
  }
};

/** 科目余额页签中类别的展示顺序 */
export const CATEGORY_ORDER: Api.FinReport.CategoryCode[] = ['ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE'];

export function getCategoryMeta(code: string): CategoryMeta {
  return CATEGORY_META[code as Api.FinReport.CategoryCode] ?? CATEGORY_META.ASSET;
}

/** 资金方向：收入 / 支出 */
export const FLOW_DIRECTION_META: Record<Api.FinReport.FlowDirection, { label: string; color: string; bg: string }> = {
  INCOME: { label: '收入', color: '#18a058', bg: 'rgba(24, 160, 88, 0.12)' },
  EXPENSE: { label: '支出', color: '#e34d59', bg: 'rgba(227, 77, 89, 0.12)' }
};

/**
 * 金额紧凑展示：超过 1 万用「万」，方便大额一眼看懂
 *
 * 例如：¥1,234.00 / ¥12.34 万
 */
export function formatMoneyCompact(value?: number | null) {
  const amount = Number(value ?? 0);

  if (Math.abs(amount) >= 10000) {
    return `¥${(amount / 10000).toFixed(2)} 万`;
  }

  return `¥${amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

/** 坐标轴金额：超过 1 万显示为「万」 */
export function formatAxisMoney(value: number) {
  if (Math.abs(value) >= 10000) {
    return `${(value / 10000).toFixed(1)}万`;
  }

  return String(value);
}

/** 百分比展示，最多一位小数 */
export function formatPercent(value?: number | null) {
  return `${Number(value ?? 0).toFixed(1).replace(/\.0$/, '')}%`;
}

/** 流水时间：MM-DD HH:mm */
export function formatFlowTime(value?: string | null) {
  if (!value) return '-';

  const date = dayjs(value);

  return date.isValid() ? date.format('MM-DD HH:mm') : value;
}

/** 流水日期：YYYY-MM-DD */
export function formatFlowDate(value?: string | null) {
  if (!value) return '-';

  const date = dayjs(value);

  return date.isValid() ? date.format('YYYY-MM-DD') : value;
}

/** 在「收支构成」中按占比排序后取前 N 项，其余合并为「其他」 */
export function collapseItems(items: Api.FinReport.AccountAmount[], top = 6) {
  if (items.length <= top) {
    return items;
  }

  const sorted = [...items].sort((a, b) => b.amount - a.amount);
  const head = sorted.slice(0, top);
  const rest = sorted.slice(top);
  const restAmount = rest.reduce((sum, item) => sum + item.amount, 0);
  const restPercent = rest.reduce((sum, item) => sum + item.percent, 0);

  if (restAmount <= 0) {
    return head;
  }

  return [...head, { code: 'OTHER', name: '其他', amount: restAmount, percent: Number(restPercent.toFixed(1)) }];
}
