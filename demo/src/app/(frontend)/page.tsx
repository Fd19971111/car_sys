import config from '@payload-config'
import Link from 'next/link'
import { connection } from 'next/server'
import { getPayload } from 'payload'

type IconName =
  | 'grid'
  | 'order'
  | 'route'
  | 'customer'
  | 'driver'
  | 'car'
  | 'wallet'
  | 'chart'
  | 'settings'
  | 'bell'
  | 'search'
  | 'plus'
  | 'chevron'
  | 'arrow'
  | 'pin'
  | 'clock'
  | 'users'
  | 'more'

function Icon({ name, className = 'size-4' }: { name: IconName; className?: string }) {
  const paths: Record<IconName, React.ReactNode> = {
    grid: (
      <>
        <rect x="3" y="3" width="7" height="7" rx="2" />
        <rect x="14" y="3" width="7" height="7" rx="2" />
        <rect x="3" y="14" width="7" height="7" rx="2" />
        <rect x="14" y="14" width="7" height="7" rx="2" />
      </>
    ),
    order: (
      <>
        <path d="M6 3h12v18l-3-2-3 2-3-2-3 2V3Z" />
        <path d="M9 8h6M9 12h6" />
      </>
    ),
    route: (
      <>
        <circle cx="6" cy="18" r="3" />
        <circle cx="18" cy="6" r="3" />
        <path d="M8.2 15.8 15.8 8.2" />
      </>
    ),
    customer: (
      <>
        <circle cx="12" cy="8" r="4" />
        <path d="M4 21c.8-4 3.5-6 8-6s7.2 2 8 6" />
      </>
    ),
    driver: (
      <>
        <circle cx="12" cy="7" r="3" />
        <path d="M5 21v-3a7 7 0 0 1 14 0v3M8 14l4 4 4-4" />
      </>
    ),
    car: (
      <>
        <path d="m5 16-2-1v-3l2-1 2-5h10l2 5 2 1v3l-2 1H5Z" />
        <path d="M7 16v2M17 16v2M7 11h10" />
      </>
    ),
    wallet: (
      <>
        <path d="M4 6h15v13H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h13" />
        <path d="M14 10h7v5h-7a2.5 2.5 0 0 1 0-5Z" />
      </>
    ),
    chart: (
      <>
        <path d="M4 20V10M10 20V4M16 20v-7M22 20H2" />
      </>
    ),
    settings: (
      <>
        <circle cx="12" cy="12" r="3" />
        <path d="M19.4 15a1.7 1.7 0 0 0 .3 1.9l.1.1-2.8 2.8-.1-.1a1.7 1.7 0 0 0-1.9-.3 1.7 1.7 0 0 0-1 1.6v.2h-4V21a1.7 1.7 0 0 0-1-1.6 1.7 1.7 0 0 0-1.9.3l-.1.1L4.2 17l.1-.1a1.7 1.7 0 0 0 .3-1.9A1.7 1.7 0 0 0 3 14H2.8v-4H3a1.7 1.7 0 0 0 1.6-1 1.7 1.7 0 0 0-.3-1.9L4.2 7 7 4.2l.1.1a1.7 1.7 0 0 0 1.9.3A1.7 1.7 0 0 0 10 3v-.2h4V3a1.7 1.7 0 0 0 1 1.6 1.7 1.7 0 0 0 1.9-.3l.1-.1L19.8 7l-.1.1a1.7 1.7 0 0 0-.3 1.9 1.7 1.7 0 0 0 1.6 1h.2v4H21a1.7 1.7 0 0 0-1.6 1Z" />
      </>
    ),
    bell: (
      <>
        <path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9M10 21h4" />
      </>
    ),
    search: (
      <>
        <circle cx="11" cy="11" r="7" />
        <path d="m20 20-4-4" />
      </>
    ),
    plus: <path d="M12 5v14M5 12h14" />,
    chevron: <path d="m9 18 6-6-6-6" />,
    arrow: <path d="M5 12h14m-5-5 5 5-5 5" />,
    pin: (
      <>
        <path d="M20 10c0 5-8 11-8 11S4 15 4 10a8 8 0 1 1 16 0Z" />
        <circle cx="12" cy="10" r="2.5" />
      </>
    ),
    clock: (
      <>
        <circle cx="12" cy="12" r="9" />
        <path d="M12 7v5l3 2" />
      </>
    ),
    users: (
      <>
        <circle cx="9" cy="8" r="3" />
        <path d="M3 20c.5-4 2.5-6 6-6s5.5 2 6 6M16 5a3 3 0 0 1 0 6M17 14c2.5.3 3.8 2.3 4 5" />
      </>
    ),
    more: (
      <>
        <circle cx="5" cy="12" r="1" fill="currentColor" />
        <circle cx="12" cy="12" r="1" fill="currentColor" />
        <circle cx="19" cy="12" r="1" fill="currentColor" />
      </>
    ),
  }
  return (
    <svg
      className={className}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden
    >
      {paths[name]}
    </svg>
  )
}

type MenuItem = {
  label: string
  icon: IconName
  href: string
  active?: boolean
  badge?: string
}

const menuGroups: { label: string; items: MenuItem[] }[] = [
  {
    label: '工作台',
    items: [{ label: '运营总览', icon: 'grid' as const, active: true, href: '#overview' }],
  },
  {
    label: '订单中心',
    items: [
      { label: '全部订单', icon: 'order' as const, href: '#orders', badge: '6' },
      { label: '行程监控', icon: 'route' as const, href: '#monitor' },
    ],
  },
  {
    label: '业务管理',
    items: [
      { label: '客户管理', icon: 'customer' as const, href: '/admin/collections/customers' },
      { label: '司机管理', icon: 'driver' as const, href: '/admin/collections/drivers' },
      { label: '车辆管理', icon: 'car' as const, href: '/admin/collections/vehicles' },
    ],
  },
  {
    label: '财务与数据',
    items: [
      { label: '财务对账', icon: 'wallet' as const, href: '#finance' },
      { label: '数据分析', icon: 'chart' as const, href: '#analytics' },
    ],
  },
  {
    label: '系统',
    items: [
      { label: '系统配置', icon: 'settings' as const, href: '/admin/collections/pricing-rules' },
    ],
  },
]

const statusMeta: Record<string, { label: string; className: string }> = {
  quoting: { label: '待报价', className: 'bg-amber-50 text-amber-700 ring-amber-600/10' },
  'awaiting-confirmation': {
    label: '待确认',
    className: 'bg-violet-50 text-violet-700 ring-violet-600/10',
  },
  dispatching: { label: '待派单', className: 'bg-blue-50 text-blue-700 ring-blue-600/10' },
  'driver-en-route': { label: '前往接驾', className: 'bg-cyan-50 text-cyan-700 ring-cyan-600/10' },
  'in-progress': {
    label: '行程中',
    className: 'bg-emerald-50 text-emerald-700 ring-emerald-600/10',
  },
  completed: { label: '已完成', className: 'bg-slate-100 text-slate-600 ring-slate-500/10' },
  cancelled: { label: '已取消', className: 'bg-rose-50 text-rose-700 ring-rose-600/10' },
}

function money(value: number | null | undefined) {
  return `¥${(value ?? 0).toLocaleString('zh-CN', { minimumFractionDigits: 2 })}`
}

function time(value: string | null | undefined) {
  if (!value) return '待安排'
  return new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }).format(new Date(value))
}

export default async function HomePage() {
  await connection()
  const payload = await getPayload({ config })
  const [orders, customers, drivers, vehicles] = await Promise.all([
    payload.find({
      collection: 'ride-orders',
      limit: 6,
      sort: '-createdAt',
      depth: 1,
      overrideAccess: true,
    }),
    payload.count({ collection: 'customers', overrideAccess: true }),
    payload.find({ collection: 'drivers', limit: 50, overrideAccess: true }),
    payload.find({ collection: 'vehicles', limit: 50, overrideAccess: true }),
  ])

  const activeOrders = orders.docs.filter((item) =>
    ['dispatching', 'driver-en-route', 'in-progress'].includes(item.status),
  )
  const monitorOrder = activeOrders[0]
  const monitorDriver =
    monitorOrder && typeof monitorOrder.driver === 'object' ? monitorOrder.driver : null
  const monitorVehicle =
    monitorOrder && typeof monitorOrder.vehicle === 'object' ? monitorOrder.vehicle : null
  const availableDrivers = drivers.docs.filter((item) => item.serviceStatus === 'available').length
  const activeVehicles = vehicles.docs.filter((item) => item.status === 'active').length
  const todayRevenue = orders.docs
    .filter((item) => item.status === 'completed')
    .reduce((total, item) => total + (item.finalFare ?? item.quotedFare ?? 0), 0)

  return (
    <main className="min-h-screen bg-[#f5f7fa] text-[#172033]">
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-[236px] flex-col border-r border-slate-200/80 bg-white lg:flex">
        <div className="flex h-[72px] items-center gap-3 border-b border-slate-100 px-6">
          <div className="grid size-9 place-items-center rounded-xl bg-[#216bfe] text-white shadow-[0_7px_20px_rgba(33,107,254,.24)]">
            <Icon name="route" className="size-5" />
          </div>
          <div>
            <p className="text-[17px] font-bold tracking-tight">轻行智控</p>
            <p className="text-[10px] font-medium tracking-[.16em] text-slate-400">RIDE OPS</p>
          </div>
        </div>
        <nav className="scrollbar-none flex-1 overflow-y-auto px-3 py-5">
          {menuGroups.map((group) => (
            <details
              className="sidebar-group mb-2"
              key={group.label}
              open={['工作台', '订单中心', '业务管理'].includes(group.label)}
            >
              <summary className="flex h-9 cursor-pointer list-none items-center justify-between rounded-lg px-3 text-[11px] font-semibold tracking-[.08em] text-slate-400 transition hover:bg-slate-50 hover:text-slate-700">
                <span>{group.label}</span>
                <Icon
                  name="chevron"
                  className="sidebar-group-chevron size-3.5 transition-transform duration-200"
                />
              </summary>
              <div className="mt-1 space-y-1 border-l border-slate-100 pl-2">
                {group.items.map((item) => (
                  <Link
                    className={`group flex h-10 items-center gap-3 rounded-lg px-3 text-[13px] font-medium transition ${item.active ? 'bg-[#edf4ff] text-[#216bfe]' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-950'}`}
                    href={item.href}
                    key={item.label}
                  >
                    <Icon name={item.icon} className="size-[18px]" />
                    <span className="flex-1">{item.label}</span>
                    {'badge' in item && item.badge && (
                      <span className="rounded-full bg-blue-50 px-2 py-0.5 text-[10px] text-blue-600">
                        {item.badge}
                      </span>
                    )}
                    {!item.active && (
                      <Icon
                        name="chevron"
                        className="size-3.5 text-slate-300 opacity-0 transition group-hover:opacity-100"
                      />
                    )}
                  </Link>
                ))}
              </div>
            </details>
          ))}
        </nav>
        <div className="m-3 rounded-xl border border-slate-100 bg-slate-50/80 p-3.5">
          <div className="flex items-center gap-3">
            <div className="grid size-9 place-items-center rounded-full bg-[#172033] text-xs font-semibold text-white">
              林
            </div>
            <div className="min-w-0 flex-1">
              <p className="truncate text-xs font-semibold">林小行</p>
              <p className="text-[10px] text-slate-400">超级管理员</p>
            </div>
            <Icon name="more" className="size-4 text-slate-400" />
          </div>
        </div>
      </aside>

      <div className="lg:pl-[236px]">
        <header className="sticky top-0 z-20 flex h-[72px] items-center justify-between border-b border-slate-200/80 bg-white/90 px-5 backdrop-blur-xl sm:px-8">
          <div>
            <h1 className="text-[17px] font-bold tracking-tight">运营总览</h1>
            <p className="hidden text-[11px] text-slate-400 sm:block">实时掌握今日业务与运力状态</p>
          </div>
          <div className="flex items-center gap-2.5">
            <Link
              className="hidden h-9 w-56 items-center gap-2 rounded-lg border border-slate-200 bg-slate-50 px-3 text-slate-400 transition hover:border-blue-200 hover:text-[#216bfe] md:flex"
              href="/admin/collections/ride-orders"
            >
              <Icon name="search" className="size-4" />
              <span className="flex-1 text-xs">搜索订单 / 客户</span>
              <span className="rounded border border-slate-200 bg-white px-1.5 py-0.5 text-[9px]">
                ⌘K
              </span>
            </Link>
            <Link
              className="relative grid size-9 place-items-center rounded-lg border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-50"
              aria-label="通知"
              href="/admin/collections/ride-orders?where[status][equals]=awaiting-confirmation"
            >
              <Icon name="bell" />
              <span className="absolute right-2 top-2 size-1.5 rounded-full bg-rose-500 ring-2 ring-white" />
            </Link>
            <Link
              className="flex h-9 items-center gap-2 rounded-lg bg-[#216bfe] px-3.5 text-xs font-semibold text-white shadow-[0_5px_16px_rgba(33,107,254,.22)] transition hover:bg-[#165ce8]"
              href="/admin/collections/ride-orders/create"
            >
              <Icon name="plus" className="size-4" />
              新建订单
            </Link>
          </div>
        </header>

        <div className="mx-auto max-w-[1540px] space-y-6 p-5 sm:p-8" id="overview">
          <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            {[
              {
                label: '今日订单',
                value: orders.totalDocs,
                suffix: '笔',
                note: '较昨日 +12.5%',
                tone: 'blue',
                icon: 'order' as const,
                bars: [14, 20, 17, 27, 25, 34, 40],
              },
              {
                label: '进行中',
                value: activeOrders.length,
                suffix: '单',
                note: '运行状态良好',
                tone: 'green',
                icon: 'route' as const,
                bars: [17, 14, 20, 24, 18, 28, 31],
              },
              {
                label: '今日流水',
                value: money(todayRevenue),
                suffix: '',
                note: '已完成订单收入',
                tone: 'purple',
                icon: 'wallet' as const,
                bars: [12, 18, 24, 22, 31, 28, 38],
              },
              {
                label: '在线运力',
                value: availableDrivers,
                suffix: ` / ${activeVehicles} 辆`,
                note: '司机 / 运营车辆',
                tone: 'orange',
                icon: 'car' as const,
                bars: [28, 25, 32, 24, 35, 37, 40],
              },
            ].map((card) => (
              <article
                className="rounded-2xl border border-slate-200/75 bg-white p-5 shadow-[0_2px_12px_rgba(15,23,42,.035)]"
                key={card.label}
              >
                <div className="flex items-start justify-between">
                  <div className={`metric-icon ${card.tone}`}>
                    <Icon name={card.icon} className="size-[19px]" />
                  </div>
                  <div className="flex h-9 items-end gap-1">
                    {card.bars.map((height, index) => (
                      <span
                        className={`w-1.5 rounded-full metric-bar ${card.tone}`}
                        style={{ height }}
                        key={index}
                      />
                    ))}
                  </div>
                </div>
                <p className="mt-5 text-xs font-medium text-slate-500">{card.label}</p>
                <div className="mt-1 flex items-baseline gap-1.5">
                  <strong className="text-[28px] font-bold tracking-[-.04em] text-slate-900">
                    {card.value}
                  </strong>
                  <span className="text-xs text-slate-400">{card.suffix}</span>
                </div>
                <p
                  className={`mt-2 text-[11px] ${card.tone === 'blue' ? 'text-emerald-600' : 'text-slate-400'}`}
                >
                  {card.note}
                </p>
              </article>
            ))}
          </section>

          <section className="grid gap-5 xl:grid-cols-[minmax(0,1.65fr)_minmax(320px,.75fr)]">
            <article
              className="overflow-hidden rounded-2xl border border-slate-200/75 bg-white shadow-[0_2px_12px_rgba(15,23,42,.035)]"
              id="orders"
            >
              <div className="flex items-center justify-between border-b border-slate-100 px-5 py-4">
                <div>
                  <h2 className="text-sm font-bold">实时订单</h2>
                  <p className="mt-1 text-[11px] text-slate-400">最近更新的行程订单</p>
                </div>
                <Link
                  className="flex items-center gap-1 text-xs font-semibold text-[#216bfe]"
                  href="/admin/collections/ride-orders"
                >
                  查看全部 <Icon name="arrow" className="size-3.5" />
                </Link>
              </div>
              {orders.docs.length ? (
                <div className="overflow-x-auto">
                  <table className="w-full min-w-[760px] text-left">
                    <thead>
                      <tr className="bg-slate-50/70 text-[10px] font-semibold tracking-wide text-slate-400">
                        <th className="px-5 py-3">订单 / 客户</th>
                        <th className="px-4 py-3">行程路线</th>
                        <th className="px-4 py-3">预约时间</th>
                        <th className="px-4 py-3">金额</th>
                        <th className="px-4 py-3">状态</th>
                        <th className="px-5 py-3 text-right">操作</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {orders.docs.map((order) => {
                        const customer = typeof order.customer === 'object' ? order.customer : null
                        const meta = statusMeta[order.status] ?? statusMeta.quoting
                        return (
                          <tr className="group hover:bg-slate-50/60" key={order.id}>
                            <td className="px-5 py-4">
                              <p className="text-xs font-semibold text-slate-800">
                                {order.orderNumber}
                              </p>
                              <p className="mt-1 text-[11px] text-slate-400">
                                {customer?.name ?? '客户'} · {order.passengerCount}人
                              </p>
                            </td>
                            <td className="max-w-[270px] px-4 py-4">
                              <div className="flex items-center gap-2 text-[11px]">
                                <span className="size-1.5 shrink-0 rounded-full bg-[#216bfe]" />
                                <span className="truncate text-slate-600">
                                  {order.pickupAddress}
                                </span>
                              </div>
                              <div className="ml-[2.5px] h-2 border-l border-dashed border-slate-300" />
                              <div className="flex items-center gap-2 text-[11px]">
                                <span className="size-1.5 shrink-0 rounded-full bg-emerald-500" />
                                <span className="truncate text-slate-600">
                                  {order.destinationAddress}
                                </span>
                              </div>
                            </td>
                            <td className="px-4 py-4 text-[11px] text-slate-500">
                              {time(order.scheduledAt)}
                            </td>
                            <td className="px-4 py-4 text-xs font-semibold text-slate-800">
                              {money(order.finalFare ?? order.quotedFare)}
                            </td>
                            <td className="px-4 py-4">
                              <span
                                className={`inline-flex rounded-full px-2.5 py-1 text-[10px] font-semibold ring-1 ring-inset ${meta.className}`}
                              >
                                {meta.label}
                              </span>
                            </td>
                            <td className="px-5 py-4 text-right">
                              <Link
                                className="inline-grid size-7 place-items-center rounded-md text-slate-400 hover:bg-slate-100 hover:text-slate-700"
                                href={`/admin/collections/ride-orders/${order.id}`}
                              >
                                <Icon name="more" className="size-4" />
                              </Link>
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="grid min-h-72 place-items-center p-8 text-center">
                  <div>
                    <div className="mx-auto grid size-12 place-items-center rounded-2xl bg-blue-50 text-[#216bfe]">
                      <Icon name="order" className="size-6" />
                    </div>
                    <p className="mt-4 text-sm font-semibold">还没有行程订单</p>
                    <p className="mt-1 text-xs text-slate-400">
                      创建首个订单后，实时状态会显示在这里
                    </p>
                    <Link
                      className="mt-4 inline-flex items-center gap-1 text-xs font-semibold text-[#216bfe]"
                      href="/admin/collections/ride-orders/create"
                    >
                      创建订单 <Icon name="arrow" className="size-3.5" />
                    </Link>
                  </div>
                </div>
              )}
            </article>

            <article
              className="rounded-2xl border border-slate-200/75 bg-white p-5 shadow-[0_2px_12px_rgba(15,23,42,.035)]"
              id="monitor"
            >
              <div className="flex items-start justify-between">
                <div>
                  <h2 className="text-sm font-bold">行程监控</h2>
                  <p className="mt-1 text-[11px] text-slate-400">
                    当前 {activeOrders.length} 个进行中订单
                  </p>
                </div>
                <span className="flex items-center gap-1.5 text-[10px] font-medium text-emerald-600">
                  <span className="size-1.5 animate-pulse rounded-full bg-emerald-500" />
                  实时
                </span>
              </div>
              <div className="map-grid relative mt-4 h-[222px] overflow-hidden rounded-xl bg-[#eef3f7]">
                <div className="absolute left-[12%] top-[28%] h-2 w-[85%] -rotate-12 rounded-full bg-white shadow-sm" />
                <div className="absolute -left-[8%] top-[58%] h-2 w-[95%] rotate-[18deg] rounded-full bg-white shadow-sm" />
                <div className="absolute left-[48%] top-[-10%] h-[125%] w-2 rotate-[12deg] rounded-full bg-white shadow-sm" />
                <svg className="absolute inset-0 size-full" viewBox="0 0 360 220">
                  <path
                    d="M35 168 C75 135, 95 154, 130 118 S210 94, 238 112 S292 90, 326 53"
                    fill="none"
                    stroke="#216bfe"
                    strokeWidth="4"
                    strokeLinecap="round"
                    strokeDasharray="2 9"
                  />
                </svg>
                <div className="absolute left-[7%] top-[70%] grid size-7 place-items-center rounded-full bg-white text-[#216bfe] shadow-md ring-4 ring-blue-500/10">
                  <Icon name="pin" className="size-4" />
                </div>
                <div className="absolute right-[6%] top-[14%] grid size-8 place-items-center rounded-lg bg-[#216bfe] text-white shadow-lg shadow-blue-500/25">
                  <Icon name="car" className="size-4" />
                </div>
                <div className="absolute left-[46%] top-[39%] rounded-lg bg-white px-2.5 py-2 shadow-lg">
                  <p className="text-[9px] text-slate-400">预计到达</p>
                  <p className="text-xs font-bold">
                    {monitorOrder ? `${monitorOrder.estimatedDurationMin} 分钟` : '暂无行程'}
                  </p>
                </div>
              </div>
              <div className="mt-4 flex items-center gap-3">
                <div className="grid size-9 place-items-center rounded-full bg-[#172033] text-[10px] font-bold text-white">
                  {monitorDriver?.name?.slice(0, 1) ?? '—'}
                </div>
                <div className="min-w-0 flex-1">
                  <p className="text-xs font-semibold">
                    {monitorDriver?.name ?? '待指派司机'} ·{' '}
                    {monitorVehicle?.plateNumber ?? '待指派车辆'}
                  </p>
                  <p className="mt-1 truncate text-[10px] text-slate-400">
                    {monitorOrder
                      ? `${monitorOrder.pickupAddress} → ${monitorOrder.destinationAddress}`
                      : '当前没有进行中的订单'}
                  </p>
                </div>
                <span className="rounded-md bg-blue-50 px-2 py-1 text-[10px] font-semibold text-blue-600">
                  {monitorOrder ? (statusMeta[monitorOrder.status]?.label ?? '处理中') : '空闲'}
                </span>
              </div>
            </article>
          </section>

          <section className="grid gap-5 lg:grid-cols-3">
            <article
              className="rounded-2xl border border-slate-200/75 bg-white p-5 shadow-[0_2px_12px_rgba(15,23,42,.035)]"
              id="analytics"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-sm font-bold">订单趋势</h2>
                  <p className="mt-1 text-[11px] text-slate-400">近 7 日订单量</p>
                </div>
                <span className="rounded-lg border border-slate-200 px-2.5 py-1.5 text-[10px] text-slate-500">
                  本周
                </span>
              </div>
              <div className="relative mt-6 flex h-32 items-end justify-between gap-3 border-b border-slate-100 px-1">
                {[38, 57, 44, 72, 60, 88, 76].map((height, index) => (
                  <div className="group flex h-full flex-1 items-end" key={index}>
                    <div
                      className={`w-full rounded-t-md ${index === 5 ? 'bg-[#216bfe]' : 'bg-blue-100'} transition group-hover:bg-blue-300`}
                      style={{ height: `${height}%` }}
                    />
                  </div>
                ))}
              </div>
              <div className="mt-2 flex justify-between px-1 text-[9px] text-slate-400">
                {['一', '二', '三', '四', '五', '六', '日'].map((item) => (
                  <span key={item}>周{item}</span>
                ))}
              </div>
            </article>

            <article
              className="rounded-2xl border border-slate-200/75 bg-white p-5 shadow-[0_2px_12px_rgba(15,23,42,.035)]"
              id="finance"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-sm font-bold">业务概况</h2>
                  <p className="mt-1 text-[11px] text-slate-400">客户与运力资产</p>
                </div>
                <Link className="text-[10px] font-semibold text-[#216bfe]" href="/admin">
                  管理数据
                </Link>
              </div>
              <div className="mt-5 grid grid-cols-3 gap-3">
                {[
                  {
                    label: '客户',
                    value: customers.totalDocs,
                    icon: 'users' as const,
                    tone: 'bg-blue-50 text-blue-600',
                  },
                  {
                    label: '司机',
                    value: drivers.totalDocs,
                    icon: 'driver' as const,
                    tone: 'bg-emerald-50 text-emerald-600',
                  },
                  {
                    label: '车辆',
                    value: vehicles.totalDocs,
                    icon: 'car' as const,
                    tone: 'bg-violet-50 text-violet-600',
                  },
                ].map((item) => (
                  <div
                    className="rounded-xl border border-slate-100 bg-slate-50/60 p-3"
                    key={item.label}
                  >
                    <div className={`grid size-8 place-items-center rounded-lg ${item.tone}`}>
                      <Icon name={item.icon} className="size-4" />
                    </div>
                    <p className="mt-4 text-xl font-bold">{item.value}</p>
                    <p className="mt-1 text-[10px] text-slate-400">{item.label}总数</p>
                  </div>
                ))}
              </div>
              <div className="mt-4 flex items-center justify-between rounded-xl bg-[#172033] px-4 py-3 text-white">
                <div>
                  <p className="text-[10px] text-slate-400">当前可用司机</p>
                  <p className="mt-1 text-lg font-bold">
                    {availableDrivers}
                    <span className="ml-1 text-[10px] font-normal text-slate-400">人在线</span>
                  </p>
                </div>
                <div className="flex -space-x-2">
                  {['林', '王', '陈'].map((item) => (
                    <span
                      className="grid size-8 place-items-center rounded-full border-2 border-[#172033] bg-slate-600 text-[9px]"
                      key={item}
                    >
                      {item}
                    </span>
                  ))}
                </div>
              </div>
            </article>

            <article className="rounded-2xl border border-slate-200/75 bg-white p-5 shadow-[0_2px_12px_rgba(15,23,42,.035)]">
              <div>
                <h2 className="text-sm font-bold">标准服务流程</h2>
                <p className="mt-1 text-[11px] text-slate-400">从需求到订单闭环</p>
              </div>
              <div className="mt-5 space-y-0">
                {[
                  ['01', '记录用车需求', '客户、人数、起终点'],
                  ['02', '自动计算报价', '里程、时长、车型规则'],
                  ['03', '客户确认报价', '记录接受或拒绝反馈'],
                  ['04', '派单并跟踪行程', '司机接驾至行程完成'],
                  ['05', '订单归档', '费用、客户记录写入数据库'],
                ].map((step, index) => (
                  <div className="flex gap-3" key={step[0]}>
                    <div className="flex flex-col items-center">
                      <span
                        className={`grid size-7 place-items-center rounded-full text-[9px] font-bold ${index === 0 ? 'bg-[#216bfe] text-white' : 'bg-slate-100 text-slate-500'}`}
                      >
                        {step[0]}
                      </span>
                      {index < 4 && <span className="h-7 w-px bg-slate-200" />}
                    </div>
                    <div className="-mt-0.5">
                      <p className="text-[11px] font-semibold text-slate-700">{step[1]}</p>
                      <p className="mt-0.5 text-[9px] text-slate-400">{step[2]}</p>
                    </div>
                  </div>
                ))}
              </div>
            </article>
          </section>
        </div>
      </div>
    </main>
  )
}
