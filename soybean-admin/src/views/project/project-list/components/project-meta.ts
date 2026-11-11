import dayjs from 'dayjs';

/** 状态等展示元信息 */
interface ProjectStatusMeta {
  /** 展示文案 */
  label: string;
  /** n-tag 的类型 */
  tagType: NaiveUI.ThemeColor;
  /** 主题色名，配合 `rgb(var(--{color}-color))` 使用 */
  color: string;
  /** iconify 图标名 */
  icon: string;
  /** 状态说明（详情页展示） */
  desc: string;
}

/**
 * 项目状态元信息：1=待开工 2=施工中 3=已完工 4=已结算
 *
 * - 与 customer 模块的 `PROJECT_STATUS_META`（详情抽屉里展示客户关联项目）保持同一套配色
 */
export const PROJECT_STATUS_META: Record<Api.Project.Status, ProjectStatusMeta> = {
  1: { label: '待开工', tagType: 'default', color: 'info', icon: 'mdi:calendar-clock', desc: '已签约，等待进场施工' },
  2: { label: '施工中', tagType: 'info', color: 'info', icon: 'mdi:hammer-wrench', desc: '正在现场施工，可持续记录进度' },
  3: { label: '已完工', tagType: 'success', color: 'success', icon: 'mdi:check-circle-outline', desc: '工程已完工，等待结算' },
  4: { label: '已结算', tagType: 'primary', color: 'primary', icon: 'mdi:cash-check', desc: '款项已结清，项目归档' }
};

/** 项目状态下拉选项 */
export const PROJECT_STATUS_OPTIONS: CommonType.Option<Api.Project.Status, string>[] = [
  { value: 1, label: '待开工' },
  { value: 2, label: '施工中' },
  { value: 3, label: '已完工' },
  { value: 4, label: '已结算' }
];

/** 获取项目状态元信息（容错兜底） */
export function getProjectStatusMeta(status: number): ProjectStatusMeta {
  return (
    PROJECT_STATUS_META[status as Api.Project.Status] ?? {
      label: '未知状态',
      tagType: 'default',
      color: 'info',
      icon: 'mdi:help-circle-outline',
      desc: '状态数据异常'
    }
  );
}

/** 项目状态标签的圆形头像内联样式（列表首列用） */
export function getProjectStatusStyle(status: number) {
  const meta = getProjectStatusMeta(status);

  return {
    backgroundColor: `rgb(var(--${meta.color}-color) / 0.12)`,
    color: `rgb(var(--${meta.color}-color))`
  };
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

  return `¥${Number(value).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

/** 格式化金额（万元）：28.60 万，用于统计卡片 */
export function formatMoneyInWan(value?: number | null) {
  if (!value) return '0.00 万';

  return `${(value / 10000).toFixed(2)} 万`;
}

/**
 * 工期文案：`2024-03-01 ~ 2024-09-30（共 214 天）`
 *
 * - 缺少日期时退化为 `-`
 */
export function formatDuration(startDate?: string | null, endDate?: string | null) {
  if (!startDate && !endDate) {
    return '-';
  }

  const start = startDate ? dayjs(startDate) : null;
  const end = endDate ? dayjs(endDate) : null;

  if (start?.isValid() && end?.isValid()) {
    return `${start.format('YYYY-MM-DD')} ~ ${end.format('YYYY-MM-DD')}（共 ${end.diff(start, 'day') + 1} 天）`;
  }

  return `${formatDate(startDate)} ~ ${formatDate(endDate)}`;
}

/**
 * 施工进度百分比
 *
 * - 待开工 0%，施工中按工期推进比例（10% ~ 90%），已完工 / 已结算 100%
 */
export function getProjectProgress(status: number, startDate?: string | null, endDate?: string | null) {
  if (status === 1) return 0;
  if (status >= 3) return 100;

  const start = startDate ? dayjs(startDate) : null;
  const end = endDate ? dayjs(endDate) : null;

  if (!start?.isValid() || !end?.isValid() || end.isBefore(start)) {
    return 50;
  }

  const total = end.diff(start, 'day') + 1;
  const passed = dayjs().diff(start, 'day') + 1;
  const percent = Math.round((passed / total) * 100);

  return Math.min(90, Math.max(10, percent));
}
