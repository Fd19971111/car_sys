/** 会计科目类别元信息：颜色用于列表标签与统计卡片 */
export const ACCOUNT_CATEGORY_META: Record<string, { label: string; color: string; bg: string }> = {
  ASSET: { label: '资产', color: '#1677ff', bg: 'rgba(22, 119, 255, 0.12)' },
  LIABILITY: { label: '负债', color: '#fa8c16', bg: 'rgba(250, 140, 22, 0.14)' },
  EQUITY: { label: '权益', color: '#722ed1', bg: 'rgba(114, 46, 209, 0.14)' },
  REVENUE: { label: '收入', color: '#18a058', bg: 'rgba(24, 160, 88, 0.14)' },
  EXPENSE: { label: '费用', color: '#d03050', bg: 'rgba(208, 48, 80, 0.14)' }
};

export const ACCOUNT_CATEGORY_OPTIONS = Object.entries(ACCOUNT_CATEGORY_META).map(([value, meta]) => ({
  label: meta.label,
  value
}));

/** 余额方向：借贷 */
export const ACCOUNT_DIRECTION_META: Record<string, { label: string; color: string; bg: string }> = {
  DEBIT: { label: '借', color: '#18a058', bg: 'rgba(24, 160, 88, 0.14)' },
  CREDIT: { label: '贷', color: '#fa8c16', bg: 'rgba(250, 140, 22, 0.14)' }
};

/** 资金账户类型 */
export const FUND_TYPE_META: Record<string, { label: string; icon: string; color: string; bg: string }> = {
  CASH: { label: '现金', icon: 'mdi:cash', color: '#18a058', bg: 'rgba(24, 160, 88, 0.14)' },
  BANK: { label: '银行', icon: 'mdi:bank-outline', color: '#1677ff', bg: 'rgba(22, 119, 255, 0.12)' }
};

export const FUND_TYPE_OPTIONS = Object.entries(FUND_TYPE_META).map(([value, meta]) => ({
  label: meta.label,
  value
}));

/** 根据科目类别推导余额方向（与后端 directionByCategory 保持一致） */
export function deriveDirection(categoryCode: string): 'DEBIT' | 'CREDIT' {
  return categoryCode === 'ASSET' || categoryCode === 'EXPENSE' ? 'DEBIT' : 'CREDIT';
}

export function getCategoryMeta(code: string) {
  return ACCOUNT_CATEGORY_META[code] ?? { label: code || '-', color: '#909399', bg: 'rgba(144, 147, 153, 0.12)' };
}

export function getDirectionMeta(code: string) {
  return ACCOUNT_DIRECTION_META[code] ?? { label: '-', color: '#909399', bg: 'rgba(144, 147, 153, 0.12)' };
}

export function getFundTypeMeta(type: string) {
  return FUND_TYPE_META[type] ?? { label: type || '-', icon: 'mdi:wallet-outline', color: '#909399', bg: 'rgba(144, 147, 153, 0.12)' };
}
