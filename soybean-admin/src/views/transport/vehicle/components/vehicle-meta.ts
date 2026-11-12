export type VehicleTagType = 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error';

interface VehicleMeta {
  label: string;
  tagType: VehicleTagType;
}

const VEHICLE_CATEGORY_META: Record<string, VehicleMeta> = {
  economy: { label: '经济型', tagType: 'default' },
  comfort: { label: '舒适型', tagType: 'info' },
  business: { label: '商务型', tagType: 'warning' },
  'six-seat': { label: '六座商务', tagType: 'success' }
};

const VEHICLE_STATUS_META: Record<string, VehicleMeta> = {
  active: { label: '运营中', tagType: 'success' },
  maintenance: { label: '保养中', tagType: 'warning' },
  inactive: { label: '已停用', tagType: 'default' }
};

export const VEHICLE_CATEGORY_OPTIONS = [
  { label: '经济型', value: 'economy' },
  { label: '舒适型', value: 'comfort' },
  { label: '商务型', value: 'business' },
  { label: '六座商务', value: 'six-seat' }
];

export const VEHICLE_STATUS_OPTIONS = [
  { label: '运营中', value: 'active' },
  { label: '保养中', value: 'maintenance' },
  { label: '已停用', value: 'inactive' }
];

export function getVehicleCategoryMeta(category: string): VehicleMeta {
  return VEHICLE_CATEGORY_META[category] ?? { label: category || '-', tagType: 'default' };
}

export function getVehicleStatusMeta(status: string): VehicleMeta {
  return VEHICLE_STATUS_META[status] ?? { label: status || '-', tagType: 'default' };
}
