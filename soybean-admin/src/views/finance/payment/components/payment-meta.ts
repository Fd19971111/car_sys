/** 款项类型元信息 */
interface PaymentTypeMeta {
  /** n-tag 的类型 */
  tagType: NaiveUI.ThemeColor;
  /** 主题色名，配合 `rgb(var(--{color}-color))` 使用 */
  color: string;
}

/** 款项类型预设（表单支持直接输入自定义类型） */
export const PAYMENT_TYPE_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '行程车费', label: '行程车费' },
  { value: '线下补收', label: '线下补收' },
  { value: '其他收款', label: '其他收款' }
];

/** 款项类型配色：定金(橙) → 首期款(蓝) → 中期款(紫) → 尾款(绿) → 全款(绿) → 材料款(蓝绿) */
export const PAYMENT_TYPE_META: Record<string, PaymentTypeMeta> = {
  行程车费: { tagType: 'success', color: 'success' },
  线下补收: { tagType: 'info', color: 'info' },
  其他收款: { tagType: 'default', color: 'info' }
};

/** 获取款项类型元信息（容错兜底，支持自定义类型） */
export function getPaymentTypeMeta(type?: string | null): PaymentTypeMeta {
  if (!type) {
    return { tagType: 'default', color: 'info' };
  }

  return PAYMENT_TYPE_META[type] ?? { tagType: 'default', color: 'info' };
}

/** 款项类型展示文案 */
export function getPaymentTypeLabel(type?: string | null) {
  return type || '未填写';
}
