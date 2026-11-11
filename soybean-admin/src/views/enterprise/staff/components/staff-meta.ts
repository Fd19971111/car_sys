import dayjs from 'dayjs';

/** 人员状态元信息 */
export const STAFF_STATUS_META: Record<number, { label: string; tagType: NaiveUI.ThemeColor }> = {
  1: { label: '在职', tagType: 'success' },
  2: { label: '离职', tagType: 'error' }
};

/** 人员状态下拉选项 */
export const STAFF_STATUS_OPTIONS: CommonType.Option<number, string>[] = [
  { value: 1, label: '在职' },
  { value: 2, label: '离职' }
];

/** 性别选项 */
export const GENDER_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '男', label: '男' },
  { value: '女', label: '女' }
];

/** 职位预设选项 */
export const POSITION_OPTIONS: CommonType.Option<string, string>[] = [
  { value: '设计师', label: '设计师' },
  { value: '项目监理', label: '项目监理' },
  { value: '施工队长', label: '施工队长' },
  { value: '工人', label: '工人' },
  { value: '财务', label: '财务' },
  { value: '行政', label: '行政' },
  { value: '其他', label: '其他' }
];

/** 获取人员状态元信息 */
export function getStaffStatusMeta(status: number) {
  return (
    STAFF_STATUS_META[status] ?? {
      label: '未知',
      tagType: 'default' as NaiveUI.ThemeColor
    }
  );
}

/** 格式化日期：YYYY-MM-DD */
export function formatDate(value?: string | null) {
  if (!value) return '-';
  const date = dayjs(value);
  return date.isValid() ? date.format('YYYY-MM-DD') : '-';
}

/** 格式化日期时间：YYYY-MM-DD HH:mm */
export function formatDateTime(value?: string | null) {
  if (!value) return '-';
  const date = dayjs(value);
  return date.isValid() ? date.format('YYYY-MM-DD HH:mm') : value;
}