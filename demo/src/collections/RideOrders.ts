import type { Access, CollectionConfig } from 'payload'

const adminOnly: Access = ({ req }) => req.user?.collection === 'admins'

export const RideOrders: CollectionConfig = {
  slug: 'ride-orders',
  labels: {
    singular: '行程订单',
    plural: '订单管理',
  },
  access: {
    create: adminOnly,
    delete: adminOnly,
    read: adminOnly,
    update: adminOnly,
  },
  admin: {
    group: '业务管理',
    useAsTitle: 'orderNumber',
    defaultColumns: ['orderNumber', 'customer', 'status', 'quotedFare', 'scheduledAt'],
  },
  fields: [
    {
      name: 'orderNumber',
      label: '订单编号',
      type: 'text',
      required: true,
      unique: true,
      index: true,
    },
    {
      name: 'customer',
      label: '客户',
      type: 'relationship',
      relationTo: 'customers',
      required: true,
    },
    {
      name: 'status',
      label: '订单状态',
      type: 'select',
      required: true,
      defaultValue: 'quoting',
      options: [
        { label: '待报价', value: 'quoting' },
        { label: '待客户确认', value: 'awaiting-confirmation' },
        { label: '待派单', value: 'dispatching' },
        { label: '司机前往中', value: 'driver-en-route' },
        { label: '行程中', value: 'in-progress' },
        { label: '已完成', value: 'completed' },
        { label: '已取消', value: 'cancelled' },
      ],
    },
    {
      type: 'row',
      fields: [
        {
          name: 'passengerCount',
          label: '乘车人数',
          type: 'number',
          required: true,
          min: 1,
          max: 8,
        },
        {
          name: 'vehicleCategory',
          label: '车型',
          type: 'select',
          required: true,
          options: [
            { label: '经济型', value: 'economy' },
            { label: '舒适型', value: 'comfort' },
            { label: '商务型', value: 'business' },
            { label: '六座商务', value: 'six-seat' },
          ],
        },
      ],
    },
    { name: 'pickupAddress', label: '上车地点', type: 'text', required: true },
    { name: 'destinationAddress', label: '目的地点', type: 'text', required: true },
    {
      type: 'row',
      fields: [
        {
          name: 'estimatedDistanceKm',
          label: '预计里程（公里）',
          type: 'number',
          required: true,
          min: 0,
        },
        {
          name: 'estimatedDurationMin',
          label: '预计时长（分钟）',
          type: 'number',
          required: true,
          min: 0,
        },
      ],
    },
    {
      type: 'row',
      fields: [
        { name: 'quotedFare', label: '客户报价（元）', type: 'number', required: true, min: 0 },
        { name: 'finalFare', label: '最终费用（元）', type: 'number', min: 0 },
      ],
    },
    {
      name: 'customerResponse',
      label: '客户反馈',
      type: 'select',
      defaultValue: 'pending',
      options: [
        { label: '等待回复', value: 'pending' },
        { label: '接受报价', value: 'accepted' },
        { label: '拒绝报价', value: 'rejected' },
      ],
    },
    { name: 'scheduledAt', label: '预约用车时间', type: 'date', required: true },
    { name: 'driver', label: '指派司机', type: 'relationship', relationTo: 'drivers' },
    { name: 'vehicle', label: '指派车辆', type: 'relationship', relationTo: 'vehicles' },
    { name: 'completedAt', label: '行程结束时间', type: 'date' },
    { name: 'internalNotes', label: '内部备注', type: 'textarea' },
  ],
}
