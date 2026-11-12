/** 费用类型元信息 */
interface PayoutTypeMeta {
  /** n-tag 的类型 */
  tagType: NaiveUI.ThemeColor;
  /** 主题色名，配合 `rgb(var(--{color}-color))` 使用 */
  color: string;
}

/** 费用类型预设（表单支持直接输入自定义类型） */
export const PAYOUT_TYPE_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '燃油费', label: '燃油费' },
  { value: '维修保养', label: '维修保养' },
  { value: '保险费', label: '保险费' },
  { value: '路桥停车', label: '路桥停车' },
  { value: '车辆租赁', label: '车辆租赁' },
  { value: '办公管理费', label: '办公管理费' },
  { value: '其他运营支出', label: '其他运营支出' }
];

/** 费用类型配色：材料(蓝) → 人工(橙) → 施工(青) → 设计(绿) → 管理(蓝) → 销售(橙) */
export const PAYOUT_TYPE_META: Record<string, PayoutTypeMeta> = {
  燃油费: { tagType: 'warning', color: 'warning' },
  维修保养: { tagType: 'info', color: 'info' },
  保险费: { tagType: 'primary', color: 'primary' },
  路桥停车: { tagType: 'success', color: 'success' },
  车辆租赁: { tagType: 'primary', color: 'primary' },
  办公管理费: { tagType: 'info', color: 'info' },
  其他运营支出: { tagType: 'default', color: 'info' }
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
 * 运营类费用默认计入主营业务成本，办公管理费计入管理费用。
 */
export const PAYOUT_TYPE_ACCOUNT_HINT: Record<string, string> = {
  燃油费: '主营业务成本',
  维修保养: '主营业务成本',
  保险费: '主营业务成本',
  路桥停车: '主营业务成本',
  车辆租赁: '主营业务成本',
  办公管理费: '管理费用',
  其他运营支出: '主营业务成本'
};
