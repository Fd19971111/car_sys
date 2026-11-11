/** 报价单位预设选项（表单支持直接输入自定义单位） */
export const QUOTE_UNIT_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '平米', label: '平米' },
  { value: '米', label: '米' },
  { value: '项', label: '项' },
  { value: '套', label: '套' },
  { value: '个', label: '个' },
  { value: '樘', label: '樘' },
  { value: '组', label: '组' },
  { value: '吨', label: '吨' },
  { value: 'kg', label: 'kg' },
  { value: '工日', label: '工日' }
];

/** 报价项名称常用预设（表单支持直接输入自定义名称） */
export const QUOTE_ITEM_PRESET_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '拆墙', label: '拆墙' },
  { value: '砌墙', label: '砌墙' },
  { value: '水电改造', label: '水电改造' },
  { value: '防水处理', label: '防水处理' },
  { value: '贴砖', label: '贴砖' },
  { value: '木工吊顶', label: '木工吊顶' },
  { value: '墙面乳胶漆', label: '墙面乳胶漆' },
  { value: '主材代购', label: '主材代购' },
  { value: '开关插座安装', label: '开关插座安装' },
  { value: '保洁清运', label: '保洁清运' }
];

/** 报价合计占合同金额的比例（0-100），合同金额为 0 时返回 0 */
export function getQuotePercent(quoteAmount?: number | null, projectAmount?: number | null) {
  if (!projectAmount || projectAmount <= 0) return 0;

  const percent = Math.round(((quoteAmount ?? 0) / projectAmount) * 100);

  return Math.min(100, Math.max(0, percent));
}
