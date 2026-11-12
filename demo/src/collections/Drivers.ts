import type { Access, CollectionConfig } from 'payload'

const adminOnly: Access = ({ req }) => req.user?.collection === 'admins'

export const Drivers: CollectionConfig = {
  slug: 'drivers',
  labels: {
    singular: '司机',
    plural: '司机管理',
  },
  access: {
    create: adminOnly,
    delete: adminOnly,
    read: adminOnly,
    update: adminOnly,
  },
  admin: {
    group: '运力管理',
    useAsTitle: 'name',
    defaultColumns: ['name', 'phone', 'serviceStatus', 'rating'],
  },
  fields: [
    { name: 'name', label: '司机姓名', type: 'text', required: true },
    { name: 'phone', label: '手机号', type: 'text', required: true, unique: true, index: true },
    {
      name: 'serviceStatus',
      label: '接单状态',
      type: 'select',
      required: true,
      defaultValue: 'available',
      options: [
        { label: '可接单', value: 'available' },
        { label: '服务中', value: 'busy' },
        { label: '已离线', value: 'offline' },
        { label: '暂停服务', value: 'suspended' },
      ],
    },
    { name: 'rating', label: '服务评分', type: 'number', min: 0, max: 5, defaultValue: 5 },
    { name: 'completedTrips', label: '完成订单', type: 'number', min: 0, defaultValue: 0 },
    { name: 'licenseNumber', label: '驾驶证号', type: 'text' },
    { name: 'joinedAt', label: '加入日期', type: 'date' },
  ],
}
