import type { Access, CollectionConfig } from 'payload'

const adminOnly: Access = ({ req }) => req.user?.collection === 'admins'

export const PricingRules: CollectionConfig = {
  slug: 'pricing-rules',
  labels: {
    singular: '计价规则',
    plural: '计价配置',
  },
  access: {
    create: adminOnly,
    delete: adminOnly,
    read: adminOnly,
    update: adminOnly,
  },
  admin: {
    group: '系统配置',
    useAsTitle: 'name',
    defaultColumns: ['name', 'category', 'baseFare', 'pricePerKm', 'enabled'],
  },
  fields: [
    { name: 'name', label: '规则名称', type: 'text', required: true },
    {
      name: 'category',
      label: '适用车型',
      type: 'select',
      required: true,
      options: [
        { label: '经济型', value: 'economy' },
        { label: '舒适型', value: 'comfort' },
        { label: '商务型', value: 'business' },
        { label: '六座商务', value: 'six-seat' },
      ],
    },
    { name: 'baseFare', label: '起步价（元）', type: 'number', required: true, min: 0 },
    { name: 'includedKm', label: '起步里程（公里）', type: 'number', required: true, min: 0 },
    { name: 'pricePerKm', label: '超里程单价（元/公里）', type: 'number', required: true, min: 0 },
    { name: 'pricePerMinute', label: '时长费（元/分钟）', type: 'number', required: true, min: 0 },
    { name: 'minimumFare', label: '最低消费（元）', type: 'number', required: true, min: 0 },
    { name: 'enabled', label: '启用规则', type: 'checkbox', defaultValue: true },
  ],
}
