/** 部门状态元信息 */
export const DEPT_STATUS_META: Record<number, { label: string; tagType: NaiveUI.ThemeColor }> = {
  1: { label: '启用', tagType: 'success' },
  2: { label: '停用', tagType: 'error' }
};

/** 部门状态下拉选项 */
export const DEPT_STATUS_OPTIONS: CommonType.Option<number, string>[] = [
  { value: 1, label: '启用' },
  { value: 2, label: '停用' }
];

/** 获取部门状态元信息 */
export function getDeptStatusMeta(status: number) {
  return (
    DEPT_STATUS_META[status] ?? {
      label: '未知',
      tagType: 'default' as NaiveUI.ThemeColor
    }
  );
}