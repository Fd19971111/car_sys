/**
 * 通用打印（PDF）与 Excel 导出工具。
 *
 * - 打印：读取「基础设置-打印设置」的抬头 / 页脚 / 日期 / 电子章，
 *   生成 A4 打印文档后调用浏览器打印（另存为 PDF 即可）。
 * - 导出：基于 xlsx 生成 .xlsx 文件。
 */
import dayjs from 'dayjs';
import html2canvas from 'html2canvas';
import { jsPDF } from 'jspdf';
import * as XLSX from 'xlsx';
import { fetchGetPrintSetting } from '@/service/api/print';

/** 打印 / 导出通用列定义（value 输出纯文本） */
export interface PrintExportColumn<T> {
  title: string;
  value: (row: T) => string;
  /** 对齐方式，默认左对齐，金额类用 right */
  align?: 'left' | 'right' | 'center';
  /** Excel 列宽（字符数），默认 12 */
  width?: number;
}

/** 打印 / 导出的合计行 */
export interface PrintExportSummary {
  label: string;
  value: string;
}

interface TableOptions<T> {
  /** 单据名，如「库存清单」 */
  title: string;
  /** 副标题（表格上方说明，如筛选条件） */
  subtitle?: string;
  columns: PrintExportColumn<T>[];
  rows: T[];
  /** 表格下方的合计行 */
  summary?: PrintExportSummary[];
  /** 附加说明（打印时显示在合计下方） */
  note?: string;
}

// ===================== 打印设置缓存 =====================

let settingCache: Api.Print.PrintSetting | null = null;

/** 获取打印设置（带缓存，保存设置后调用 invalidatePrintSetting 清缓存） */
export async function getPrintSetting(force = false): Promise<Api.Print.PrintSetting> {
  if (settingCache && !force) {
    return settingCache;
  }

  const { data, error } = await fetchGetPrintSetting();
  settingCache = !error && data?.setting
    ? data.setting
    : { title: '', subTitle: '', footer: '', showDate: true, showSeal: false, sealImage: '' };

  return settingCache;
}

/** 保存打印设置后清除缓存 */
export function invalidatePrintSetting() {
  settingCache = null;
}

// ===================== 打印（PDF） =====================

/** 生成 A4 打印 HTML */
function buildPrintHtml<T>(opts: TableOptions<T>, setting: Api.Print.PrintSetting): string {
  const headTitle = setting.title || opts.title;
  const subLine = [setting.subTitle, opts.subtitle].filter(Boolean).join('　|　');
  const dateText = setting.showDate ? dayjs().format('YYYY年MM月DD日') : '';

  const ths = opts.columns.map(col => `<th>${escapeHtml(col.title)}</th>`).join('');
  const trs = opts.rows
    .map(row => {
      const tds = opts.columns
        .map(col => `<td class="col-${col.align || 'left'}">${escapeHtml(col.value(row))}</td>`)
        .join('');
      return `<tr>${tds}</tr>`;
    })
    .join('');

  const summaryHtml = (opts.summary ?? [])
    .map(
      item =>
        `<div class="summary-row"><span class="summary-label">${escapeHtml(item.label)}</span><span>${escapeHtml(item.value)}</span></div>`
    )
    .join('');

  const noteHtml = opts.note ? `<div class="note">${escapeHtml(opts.note)}</div>` : '';
  const footerHtml = setting.footer ? `<div class="print-footer">${escapeHtml(setting.footer)}</div>` : '';
  const sealHtml =
    setting.showSeal && setting.sealImageUrl
      ? `<img class="print-seal" src="${escapeHtml(setting.sealImageUrl)}" alt="seal" />`
      : '';

  return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8" />
<title>${escapeHtml(headTitle)}</title>
<style>
  @page { size: A4 portrait; margin: 14mm 12mm; }
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
    color: #1f2937;
    font-size: 12px;
    position: relative;
  }
  .print-header { text-align: center; margin-bottom: 14px; }
  .print-header .title { font-size: 20px; font-weight: 700; letter-spacing: 2px; }
  .print-header .sub { margin-top: 4px; font-size: 11px; color: #6b7280; }
  .print-header .date { text-align: right; font-size: 11px; color: #6b7280; margin-top: -16px; }
  table { width: 100%; border-collapse: collapse; }
  th, td { border: 1px solid #94a3b8; padding: 5px 8px; word-break: break-all; }
  th { background: #f1f5f9; font-weight: 600; text-align: center; }
  td.col-right { text-align: right; }
  td.col-center { text-align: center; }
  tr:nth-child(even) td { background: #fafafa; }
  .summary { margin-top: 10px; display: flex; flex-direction: column; gap: 4px; align-items: flex-end; }
  .summary-row { font-size: 12px; }
  .summary-label { color: #6b7280; margin-right: 12px; }
  .note { margin-top: 8px; font-size: 11px; color: #6b7280; }
  .print-footer { margin-top: 22px; text-align: center; font-size: 11px; color: #9ca3af; }
  .print-seal {
    position: absolute;
    right: 60px;
    bottom: 60px;
    width: 110px;
    height: 110px;
    object-fit: contain;
    opacity: 0.75;
    transform: rotate(-12deg);
    pointer-events: none;
  }
</style>
</head>
<body>
  <div class="print-header">
    <div class="title">${escapeHtml(headTitle)}</div>
    ${subLine ? `<div class="sub">${escapeHtml(subLine)}</div>` : ''}
    ${dateText ? `<div class="date">打印日期：${escapeHtml(dateText)}</div>` : ''}
  </div>
  <table>
    <thead><tr>${ths}</tr></thead>
    <tbody>${trs || '<tr><td colspan="' + opts.columns.length + '" class="col-center">暂无数据</td></tr>'}</tbody>
  </table>
  ${summaryHtml ? `<div class="summary">${summaryHtml}</div>` : ''}
  ${noteHtml}
  ${footerHtml}
  ${sealHtml}
</body>
</html>`;
}

/** 通过隐藏 iframe 调起浏览器打印（打印预览中另存为 PDF） */
function printHtml(html: string) {
  const iframe = document.createElement('iframe');
  iframe.style.cssText = 'position:fixed;right:0;bottom:0;width:0;height:0;border:0;visibility:hidden;';
  document.body.appendChild(iframe);

  const doc = iframe.contentWindow?.document;
  if (!doc) {
    document.body.removeChild(iframe);

    return;
  }

  doc.open();
  doc.write(html);
  doc.close();

  const win = iframe.contentWindow!;
  const cleanup = () => {
    window.setTimeout(() => iframe.parentNode?.removeChild(iframe), 500);
  };

  if (win.document.readyState === 'complete') {
    win.focus();
    win.print();
    win.onafterprint = cleanup;
  } else {
    iframe.onload = () => {
      win.focus();
      win.print();
      win.onafterprint = cleanup;
    };
  }
}

/** 生成打印文档 HTML（读取打印设置），供打印调起与打印预览共用 */
export async function renderPrintHtml<T>(opts: TableOptions<T>): Promise<string> {
  const setting = await getPrintSetting();

  return buildPrintHtml(opts, setting);
}

/** 打印表格：读取打印设置并调起浏览器打印 */
export async function printTable<T>(opts: TableOptions<T>) {
  const html = await renderPrintHtml(opts);
  printHtml(html);
}

// ===================== PDF 导出 =====================

/**
 * 将打印文档 HTML 直接保存为 PDF 文件（A4，多页自动分页）。
 *
 * 原理：在离屏 iframe 中渲染文档，截图后按 A4 比例切片写入 jsPDF。
 */
export async function saveHtmlAsPdf(html: string, filename: string) {
  // 离屏 iframe 渲染文档（794×1123 ≈ A4 @96dpi）
  const iframe = document.createElement('iframe');
  iframe.style.cssText = 'position:fixed;left:-10000px;top:0;width:794px;height:1123px;border:0;visibility:hidden;';
  document.body.appendChild(iframe);

  const doc = iframe.contentWindow?.document;
  if (!doc) {
    document.body.removeChild(iframe);
    throw new Error('当前浏览器不支持导出 PDF');
  }

  doc.open();
  doc.write(html);
  doc.close();

  // 屏幕渲染时 @page 边距不生效，手动补齐与打印一致的页边距
  const style = doc.createElement('style');
  style.textContent = 'body { padding: 14mm 12mm; }';
  doc.head.appendChild(style);

  // 等待文档与图片渲染完成
  await new Promise<void>(resolve => {
    const done = () => resolve();
    iframe.onload = done;
    window.setTimeout(done, 500);
  });

  try {
    const canvas = await html2canvas(doc.body, {
      scale: 2,
      useCORS: true,
      backgroundColor: '#ffffff',
      windowWidth: 794
    });

    const pdf = new jsPDF('p', 'mm', 'a4');
    const pageWidth = pdf.internal.pageSize.getWidth();
    const pageHeight = pdf.internal.pageSize.getHeight();
    const pxPerMm = canvas.width / pageWidth;
    const pageHeightPx = pageHeight * pxPerMm;

    let rendered = 0;
    let pageIndex = 0;

    while (rendered < canvas.height - 1) {
      const sliceHeight = Math.min(pageHeightPx, canvas.height - rendered);
      const slice = document.createElement('canvas');
      slice.width = canvas.width;
      slice.height = Math.ceil(sliceHeight);

      const ctx = slice.getContext('2d');
      if (!ctx) {
        throw new Error('PDF 生成失败');
      }

      ctx.fillStyle = '#ffffff';
      ctx.fillRect(0, 0, slice.width, slice.height);
      ctx.drawImage(canvas, 0, rendered, canvas.width, sliceHeight, 0, 0, canvas.width, sliceHeight);

      if (pageIndex > 0) {
        pdf.addPage();
      }
      pdf.addImage(slice.toDataURL('image/jpeg', 0.92), 'JPEG', 0, 0, pageWidth, sliceHeight / pxPerMm);

      rendered += sliceHeight;
      pageIndex += 1;
    }

    pdf.save(`${filename}-${dayjs().format('YYYYMMDD')}.pdf`);
  } finally {
    iframe.parentNode?.removeChild(iframe);
  }
}

// ===================== Excel 导出 =====================

/** 导出 Excel（.xlsx） */
export function exportExcel<T>(opts: TableOptions<T> & { filename: string }) {
  const header = opts.columns.map(col => col.title);
  const body = opts.rows.map(row => opts.columns.map(col => col.value(row)));
  const summary = (opts.summary ?? []).map(item => [item.label, item.value]);

  const aoa: (string | number)[][] = [header, ...body];
  if (summary.length > 0) {
    aoa.push([], ...summary);
  }

  const sheet = XLSX.utils.aoa_to_sheet(aoa);
  sheet['!cols'] = opts.columns.map(col => ({ wch: col.width ?? 14 }));

  const book = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(book, sheet, opts.title.slice(0, 30) || 'Sheet1');
  XLSX.writeFile(book, `${opts.filename}-${dayjs().format('YYYYMMDD')}.xlsx`);
}

// ===================== 工具 =====================

/** 分页拉取全部数据（size=100，最多 50 页），用于打印 / 导出当前筛选的完整数据 */
export async function fetchAllPages<T>(
  fetcher: (current: number) => Promise<{ records: T[]; total: number }>
): Promise<T[]> {
  const first = await fetcher(1);
  const all: T[] = [...first.records];
  const totalPages = Math.min(Math.ceil(first.total / 100), 50);

  for (let current = 2; current <= totalPages; current += 1) {
    const page = await fetcher(current); // eslint-disable-line no-await-in-loop
    all.push(...page.records);
  }

  return all;
}

function escapeHtml(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
