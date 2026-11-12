export type DriverStatusTagType = 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error';

interface DriverStatusMeta {
  label: string;
  tagType: DriverStatusTagType;
}

const DRIVER_STATUS_META: Record<string, DriverStatusMeta> = {
  available: { label: '可接单', tagType: 'success' },
  busy: { label: '服务中', tagType: 'info' },
  offline: { label: '已离线', tagType: 'default' },
  suspended: { label: '暂停服务', tagType: 'warning' }
};

export const DRIVER_STATUS_OPTIONS = [
  { label: '可接单', value: 'available' },
  { label: '服务中', value: 'busy' },
  { label: '已离线', value: 'offline' },
  { label: '暂停服务', value: 'suspended' }
];

export const DRIVER_LICENSE_TYPE_OPTIONS = ['A1', 'A2', 'A3', 'B1', 'B2', 'C1', 'C2'].map(item => ({
  label: item,
  value: item
}));

export function getDriverStatusMeta(status: string): DriverStatusMeta {
  return DRIVER_STATUS_META[status] ?? { label: status || '-', tagType: 'default' };
}

export function formatDriverDate(value: string): string {
  return value || '-';
}
