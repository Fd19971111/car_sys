/**
 * 首页展示用的元信息与纯函数工具。
 *
 * 首页只做展示，不改动任何业务数据，因此这里只放「金额换算 / 取比例 / 配色」这类无副作用的东西。
 */

import dayjs from 'dayjs';

/**
 * 施工进度的业务日期标签。
 *
 * 进度支持补录历史（业务日期 logDate），不能用「相对时间」，
 * 统一展示为：今天 / 昨天 / MM-DD（跨年时 YYYY-MM-DD）。
 */
export function progressTimeLabel(value?: string | null): string {
  if (!value) return '-';

  const date = dayjs(value);

  if (!date.isValid()) return value;

  const today = dayjs();

  if (date.isSame(today, 'day')) return '今天';
  if (date.isSame(today.subtract(1, 'day'), 'day')) return '昨天';
  if (date.isSame(today, 'year')) return date.format('MM-DD');

  return date.format('YYYY-MM-DD');
}

/** 金额（元）→ 万元数值（保留两位小数），用于数字滚动与占比计算 */
export function toWan(value?: number | null) {
  return Number(((value ?? 0) / 10000).toFixed(2));
}

/** 占比：分母为 0 时返回 0，结果取 0 ~ 100 的整数 */
export function toPercent(part?: number | null, total?: number | null) {
  if (!total || total <= 0) return 0;

  return Math.min(100, Math.max(0, Math.round(((part ?? 0) / total) * 100)));
}

/** 首页「项目状态分布」条目 */
export interface HomeStatusItem {
  /** 取值来自 Api.Project.Stats 的同名字段 */
  key: 'pending' | 'building' | 'finished' | 'settled';
  label: string;
  /** 完整的 CSS 颜色表达式（主题色都以 "r g b" 形式暴露，需自行包 rgb()） */
  color: string;
  icon: string;
}

/**
 * 首页「项目状态分布」的配色
 *
 * 两个约束：
 * 1. 不能复用 `project-meta.ts` 的 `PROJECT_STATUS_META`——那份配色里「待开工 / 施工中」同为 info 色，放在同一根分布条上分不清；
 * 2. 本主题开启了 `isInfoFollowPrimary`，`--info-color` 与 `--primary-color` 取值完全相同，
 *    因此四个状态改用 warning / primary / success 三个色相，加上一个自适应深浅的第四色。
 *
 * 「已结算」用的是 `--home-settled-color`：它在 `business-overview.vue` 里定义，
 * 浅色主题取 primary-800（深靛蓝），深色主题换成 primary-300，避免深色底上糊成一团。
 */
export const HOME_PROJECT_STATUS: HomeStatusItem[] = [
  { key: 'pending', label: '待开工', color: 'rgb(var(--warning-color))', icon: 'mdi:calendar-clock' },
  { key: 'building', label: '施工中', color: 'rgb(var(--primary-color))', icon: 'mdi:hammer-wrench' },
  { key: 'finished', label: '已完工', color: 'rgb(var(--success-color))', icon: 'mdi:check-circle-outline' },
  { key: 'settled', label: '已结算', color: 'rgb(var(--home-settled-color))', icon: 'mdi:cash-check' }
];
