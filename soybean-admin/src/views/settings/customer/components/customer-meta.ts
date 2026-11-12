import dayjs from 'dayjs';

/** 状态 / 来源等展示元信息 */
interface CustomerStatusMeta {
  /** 展示文案 */
  label: string;
  /** n-tag 的类型 */
  tagType: NaiveUI.ThemeColor;
  /** 主题色名，配合 `rgb(var(--{color}-color))` 使用 */
  color: string;
  /** iconify 图标名 */
  icon: string;
}

/** 客户状态元信息：1=跟进中 2=已成交 3=已流失 */
export const CUSTOMER_STATUS_META: Record<Api.Customer.Status, CustomerStatusMeta> = {
  1: { label: '跟进中', tagType: 'warning', color: 'warning', icon: 'mdi:progress-clock' },
  2: { label: '已成交', tagType: 'success', color: 'success', icon: 'mdi:handshake-outline' },
  3: { label: '已流失', tagType: 'error', color: 'error', icon: 'mdi:account-off-outline' }
};

/** 客户状态下拉选项 */
export const CUSTOMER_STATUS_OPTIONS: CommonType.Option<Api.Customer.Status, string>[] = [
  { value: 1, label: '跟进中' },
  { value: 2, label: '已成交' },
  { value: 3, label: '已流失' }
];

/** 客户来源预设选项（表单里支持直接输入自定义来源） */
export const CUSTOMER_SOURCE_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '转介绍', label: '转介绍' },
  { value: '老客户推荐', label: '老客户推荐' },
  { value: '线上推广', label: '线上推广' },
  { value: '门店到访', label: '门店到访' },
  { value: '电话咨询', label: '电话咨询' },
  { value: '其他', label: '其他' }
];

/** 获取客户状态元信息（容错兜底） */
export function getCustomerStatusMeta(status: number): CustomerStatusMeta {
  return (
    CUSTOMER_STATUS_META[status as Api.Customer.Status] ?? {
      label: '未知状态',
      tagType: 'default',
      color: 'primary',
      icon: 'mdi:help-circle-outline'
    }
  );
}

/** 客户来源展示文案 */
export function getCustomerSourceLabel(source?: string | null) {
  return source || '未填写';
}

/** 格式化日期时间：YYYY-MM-DD HH:mm */
export function formatDateTime(value?: string | null) {
  if (!value) return '-';

  const date = dayjs(value);

  return date.isValid() ? date.format('YYYY-MM-DD HH:mm') : value;
}

/** 格式化日期：YYYY-MM-DD */
export function formatDate(value?: string | null) {
  if (!value) return '-';

  const date = dayjs(value);

  return date.isValid() ? date.format('YYYY-MM-DD') : value;
}

/** 格式化金额：¥286,000.00 */
export function formatMoney(value?: number | null) {
  if (value === null || value === undefined) return '-';

  return `¥${value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

/** 格式化金额（万元）：28.60 万 */
export function formatMoneyInWan(value?: number | null) {
  if (!value) return '0 万';

  return `${(value / 10000).toFixed(2)} 万`;
}
