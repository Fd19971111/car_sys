/** 费用类型元信息 */
interface PayoutTypeMeta {
  /** n-tag 的类型 */
  tagType: NaiveUI.ThemeColor;
  /** 主题色名，配合 `rgb(var(--{color}-color))` 使用 */
  color: string;
}

/** 费用类型预设（表单支持直接输入自定义类型） */
export const PAYOUT_TYPE_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '材料款', label: '材料款' },
  { value: '人工费', label: '人工费' },
  { value: '施工费', label: '施工费' },
  { value: '设计费', label: '设计费' },
  { value: '管理费', label: '管理费' },
  { value: '销售费', label: '销售费' },
  { value: '其他', label: '其他' }
];

/** 费用类型配色：材料(蓝) → 人工(橙) → 施工(青) → 设计(绿) → 管理(蓝) → 销售(橙) */
export const PAYOUT_TYPE_META: Record<string, PayoutTypeMeta> = {
  材料款: { tagType: 'primary', color: 'primary' },
  人工费: { tagType: 'warning', color: 'warning' },
  施工费: { tagType: 'info', color: 'info' },
  设计费: { tagType: 'success', color: 'success' },
  管理费: { tagType: 'primary', color: 'primary' },
  销售费: { tagType: 'warning', color: 'warning' },
  其他: { tagType: 'default', color: 'info' }
};

/** 获取费用类型元信息（容错兜底，支持自定义类型） */
export function getPayoutTypeMeta(type?: string | null): PayoutTypeMeta {
  if (!type) {
    return { tagType: 'default', color: 'info' };
  }

  return PAYOUT_TYPE_META[type] ?? { tagType: 'default', color: 'info' };
}

/** 费用类型展示文案 */
export function getPayoutTypeLabel(type?: string | null) {
  return type || '未填写';
}

/**
 * 费用类型对应的记账借方科目说明，用于表单里向用户解释「这笔钱记到哪里」。
 *
 * 与后端 booking 包的映射保持一致：材料/人工/施工/设计计入主营业务成本，
 * 管理费计入管理费用，销售费计入销售费用。
 */
export const PAYOUT_TYPE_ACCOUNT_HINT: Record<string, string> = {
  材料款: '主营业务成本',
  人工费: '主营业务成本',
  施工费: '主营业务成本',
  设计费: '主营业务成本',
  管理费: '管理费用',
  销售费: '销售费用',
  其他: '主营业务成本'
};
