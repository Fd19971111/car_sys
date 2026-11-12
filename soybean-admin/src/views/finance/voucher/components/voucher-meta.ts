import * as XLSX from 'xlsx';

/** 记账凭证的展示文案与导出工具 */

interface TagMeta {
  label: string;
  tagType: NaiveUI.ThemeColor;
}

/** 凭证类型 */
export const VOUCHER_TYPE_META: Record<string, TagMeta> = {
  GENERAL: { label: '通用记账凭证', tagType: 'default' },
  RECEIPT: { label: '收款凭证', tagType: 'success' },
  PAYMENT: { label: '付款凭证', tagType: 'warning' },
  TRANSFER: { label: '转账凭证', tagType: 'info' }
};

export function getVoucherTypeMeta(code?: string | null): TagMeta {
  return VOUCHER_TYPE_META[code ?? ''] ?? { label: code || '记账凭证', tagType: 'default' };
}

/** 凭证状态 */
export const VOUCHER_STATUS_META: Record<string, TagMeta> = {
  DRAFT: { label: '草稿', tagType: 'info' },
  POSTED: { label: '已记账', tagType: 'success' },
  VOID: { label: '已作废', tagType: 'default' }
};

export function getVoucherStatusMeta(code?: string | null): TagMeta {
  return VOUCHER_STATUS_META[code ?? ''] ?? { label: code || '-', tagType: 'default' };
}

/** 会计期间状态 */
export const PERIOD_STATUS_META: Record<string, TagMeta> = {
  OPEN: { label: '记账中', tagType: 'success' },
  CLOSED: { label: '已结账', tagType: 'default' },
  FUTURE: { label: '未启用', tagType: 'info' }
};

export function getPeriodStatusMeta(code?: string | null): TagMeta {
  return PERIOD_STATUS_META[code ?? ''] ?? { label: code || '-', tagType: 'default' };
}

/** 凭证来源 */
export const SOURCE_TYPE_LABEL: Record<string, string> = {
  MANUAL: '手工录入',
  PROJECT_PAYMENT: '收款',
  PROJECT_PAYOUT: '付款'
};

export function getSourceTypeLabel(type?: string | null) {
  return SOURCE_TYPE_LABEL[type ?? ''] ?? type ?? '-';
}

/**
 * 导出凭证分录明细为 Excel 文件。
 *
 * @param rows 导出行
 * @param filename 文件名（不含扩展名）
 */
export function exportVoucherExcel(rows: Api.FinVoucher.ExportRow[], filename: string) {
  const headers = [
    '凭证号',
    '凭证日期',
    '会计期间',
    '凭证类型',
    '凭证摘要',
    '状态',
    '来源',
    '来源ID',
    '行号',
    '科目编码',
    '科目名称',
    '分录摘要',
    '借贷方向',
    '金额'
  ];

  const lines = rows.map(row => [
    row.voucherNo,
    row.voucherDate,
    row.periodCode,
    getVoucherTypeMeta(row.voucherTypeCode).label,
    row.voucherSummary,
    getVoucherStatusMeta(row.statusCode).label,
    getSourceTypeLabel(row.sourceType),
    row.sourceId || '',
    row.lineNo,
    row.accountCode,
    row.accountName,
    row.entrySummary,
    row.directionCode === 'DEBIT' ? '借' : '贷',
    row.amount
  ]);

  const sheet = XLSX.utils.aoa_to_sheet([headers, ...lines]);

  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(workbook, sheet, '凭证分录');

  XLSX.writeFile(workbook, `${filename}.xlsx`);
}
