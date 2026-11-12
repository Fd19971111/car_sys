/** 金额（元）→ 万元数值（保留两位小数），用于数字滚动展示。 */
export function toWan(value?: number | null) {
  return Number(((value ?? 0) / 10000).toFixed(2));
}
