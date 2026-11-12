export type OrderTagType = 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error';

interface OrderMeta {
  label: string;
  tagType: OrderTagType;
  color: string;
}

const ORDER_STATUS_META: Record<string, OrderMeta> = {
  quoting: { label: '待报价', tagType: 'warning', color: 'warning' },
  'awaiting-confirmation': { label: '待客户确认', tagType: 'info', color: 'info' },
  dispatching: { label: '待派单', tagType: 'primary', color: 'primary' },
  'driver-en-route': { label: '司机前往中', tagType: 'info', color: 'info' },
  'in-progress': { label: '行程中', tagType: 'success', color: 'success' },
  completed: { label: '已完成', tagType: 'default', color: 'default' },
  cancelled: { label: '已取消', tagType: 'error', color: 'error' }
};

const CUSTOMER_RESPONSE_META: Record<string, OrderMeta> = {
  pending: { label: '等待回复', tagType: 'warning', color: 'warning' },
  accepted: { label: '接受报价', tagType: 'success', color: 'success' },
  rejected: { label: '拒绝报价', tagType: 'error', color: 'error' }
};

const VEHICLE_CATEGORY_META: Record<string, OrderMeta> = {
  economy: { label: '经济型', tagType: 'default', color: 'default' },
  comfort: { label: '舒适型', tagType: 'info', color: 'info' },
  business: { label: '商务型', tagType: 'warning', color: 'warning' },
  'six-seat': { label: '六座商务', tagType: 'success', color: 'success' }
};

export const ORDER_STATUS_OPTIONS = [
  { label: '待报价', value: 'quoting' },
  { label: '待客户确认', value: 'awaiting-confirmation' },
  { label: '待派单', value: 'dispatching' },
  { label: '司机前往中', value: 'driver-en-route' },
  { label: '行程中', value: 'in-progress' },
  { label: '已完成', value: 'completed' },
  { label: '已取消', value: 'cancelled' }
];

export const CUSTOMER_RESPONSE_OPTIONS = [
  { label: '等待回复', value: 'pending' },
  { label: '接受报价', value: 'accepted' },
  { label: '拒绝报价', value: 'rejected' }
];

export const VEHICLE_CATEGORY_OPTIONS = [
  { label: '经济型', value: 'economy' },
  { label: '舒适型', value: 'comfort' },
  { label: '商务型', value: 'business' },
  { label: '六座商务', value: 'six-seat' }
];

export const ACTIVE_ORDER_STATUSES: Api.Order.Status[] = ['dispatching', 'driver-en-route', 'in-progress'];

export function getOrderStatusMeta(status: string): OrderMeta {
  return ORDER_STATUS_META[status] ?? { label: status || '-', tagType: 'default', color: 'default' };
}

export function getCustomerResponseMeta(response: string): OrderMeta {
  return CUSTOMER_RESPONSE_META[response] ?? { label: response || '-', tagType: 'default', color: 'default' };
}

export function getVehicleCategoryMeta(category: string): OrderMeta {
  return VEHICLE_CATEGORY_META[category] ?? { label: category || '-', tagType: 'default', color: 'default' };
}

export function formatOrderMoney(value: number | null | undefined): string {
  return `¥${Number(value ?? 0).toFixed(2)}`;
}

export function formatOrderDateTime(value: string | null | undefined): string {
  return value || '-';
}
