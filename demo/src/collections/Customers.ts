import type { Access, CollectionConfig } from 'payload'

const adminOnly: Access = ({ req }) => req.user?.collection === 'admins'

export const Customers: CollectionConfig = {
  slug: 'customers',
  labels: {
    singular: '客户',
    plural: '客户管理',
  },
  access: {
    create: adminOnly,
    delete: adminOnly,
    read: adminOnly,
    update: adminOnly,
  },
  admin: {
    group: '业务管理',
    useAsTitle: 'name',
    defaultColumns: ['name', 'phone', 'totalTrips', 'lastTripAt'],
  },
  fields: [
    { name: 'name', label: '客户姓名', type: 'text', required: true },
    { name: 'phone', label: '手机号', type: 'text', required: true, unique: true, index: true },
    {
      name: 'source',
      label: '客户来源',
      type: 'select',
      defaultValue: 'mini-program',
      options: [
        { label: '微信小程序', value: 'mini-program' },
        { label: '客服代叫', value: 'operator' },
        { label: '企业客户', value: 'enterprise' },
      ],
    },
    { name: 'totalTrips', label: '累计行程', type: 'number', defaultValue: 0, min: 0 },
    { name: 'lastTripAt', label: '最近行程', type: 'date' },
    { name: 'notes', label: '客户备注', type: 'textarea' },
  ],
}
