import type { Access, CollectionConfig } from 'payload'

const adminOnly: Access = ({ req }) => req.user?.collection === 'admins'

export const Vehicles: CollectionConfig = {
  slug: 'vehicles',
  labels: {
    singular: '车辆',
    plural: '车辆管理',
  },
  access: {
    create: adminOnly,
    delete: adminOnly,
    read: adminOnly,
    update: adminOnly,
  },
  admin: {
    group: '运力管理',
    useAsTitle: 'plateNumber',
    defaultColumns: ['plateNumber', 'brand', 'category', 'seatCount', 'status'],
  },
  fields: [
    {
      name: 'plateNumber',
      label: '车牌号',
      type: 'text',
      required: true,
      unique: true,
      index: true,
    },
    { name: 'brand', label: '品牌车型', type: 'text', required: true },
    {
      name: 'category',
      label: '服务车型',
      type: 'select',
      required: true,
      defaultValue: 'comfort',
      options: [
        { label: '经济型', value: 'economy' },
        { label: '舒适型', value: 'comfort' },
        { label: '商务型', value: 'business' },
        { label: '六座商务', value: 'six-seat' },
      ],
    },
    { name: 'seatCount', label: '可乘人数', type: 'number', required: true, min: 1, max: 8 },
    {
      name: 'status',
      label: '车辆状态',
      type: 'select',
      required: true,
      defaultValue: 'active',
      options: [
        { label: '运营中', value: 'active' },
        { label: '保养中', value: 'maintenance' },
        { label: '已停用', value: 'inactive' },
      ],
    },
    { name: 'driver', label: '当前司机', type: 'relationship', relationTo: 'drivers' },
  ],
}
