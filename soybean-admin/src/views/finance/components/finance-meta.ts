import dayjs from 'dayjs';

/**
 * 财务模块（报价 / 收款）共用的展示格式化工具。
 *
 * 报价与收款两个页面都涉及金额与日期展示，这里统一实现，避免各写一份。
 */

/** 格式化金额：¥286,000.00 */
export function formatMoney(value?: number | null) {
  if (value === null || value === undefined) return '-';

  return `¥${Number(value).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

/** 格式化金额（万元）：28.60 万，用于统计卡片 */
export function formatMoneyInWan(value?: number | null) {
  if (!value) return '0.00 万';

  return `${(value / 10000).toFixed(2)} 万`;
}

/** 格式化日期：YYYY-MM-DD */
export function formatDate(value?: string | null) {
  if (!value) return '-';

  const date = dayjs(value);

  return date.isValid() ? date.format('YYYY-MM-DD') : value;
}

/** 格式化日期时间：YYYY-MM-DD HH:mm */
export function formatDateTime(value?: string | null) {
  if (!value) return '-';

  const date = dayjs(value);

  return date.isValid() ? date.format('YYYY-MM-DD HH:mm') : value;
}

/** 数字展示：去掉多余的 0，如 80.00 → 80、12.50 → 12.5 */
export function formatNumber(value?: number | null) {
  if (value === null || value === undefined) return '-';

  return String(Number(value));
}

/** 计算收款进度百分比，合同金额为 0 时返回 0 */
export function getPaidPercent(paidAmount?: number | null, projectAmount?: number | null) {
  if (!projectAmount || projectAmount <= 0) return 0;

  const percent = Math.round(((paidAmount ?? 0) / projectAmount) * 100);

  return Math.min(100, Math.max(0, percent));
}
