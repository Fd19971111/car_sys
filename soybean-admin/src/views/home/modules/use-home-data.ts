import { computed, ref } from 'vue';
import dayjs from 'dayjs';
import type { RouteKey } from '@elegant-router/types';
import { fetchGetCustomerList, fetchGetCustomerStats } from '@/service/api/customer';
import { fetchGetPaymentStats } from '@/service/api/payment';
import { fetchGetProjectList, fetchGetProjectStats } from '@/service/api/project';
import { fetchGetQuoteStats } from '@/service/api/quote';
import { useRouteStore } from '@/store/modules/route';

/** 逾期未完工项目：在项目信息上补充逾期天数 */
export interface OverdueProject extends Api.Project.Item {
  /** 逾期天数（> 0） */
  overdueDays: number;
}

/** 首页「最近施工进度」条目 */
export interface ProgressItem {
  /** 项目ID */
  projectId: number;
  /** 项目名称 */
  projectName: string;
  /** 项目状态：1=待开工 2=施工中 3=已完工 4=已结算 */
  status: number;
  /** 所属客户姓名 */
  customerName: string;
  /** 进度内容 */
  content: string;
  /** 记录时间（YYYY-MM-DD HH:mm:ss） */
  time: string;
  /** 合同金额（元） */
  amount: number;
  /** 项目负责人 */
  managerName: string;
}

/**
 * 首页数据 Hook：全部取自易民家业务接口，不再有演示数据。
 *
 * - 四个统计块来自 `/customer/stats`、`/project/stats`、`/quote/stats`、`/payment/stats`
 * - 待办与最近进度派生自项目 / 客户列表（后端暂无「全局最近进度」「逾期项目」专用接口）
 */
export function useHomeData() {
  /** 一次拉取的条数：用于派生待办与最近进度，避免首页发多次分页请求 */
  const OVERVIEW_SIZE = 100;
  const routeStore = useRouteStore();

  function hasRouteAccess(routeKeys: RouteKey[]) {
    function includesRoute(menus: App.Global.Menu[]): boolean {
      return menus.some(menu => routeKeys.includes(menu.routeKey) || (menu.children ? includesRoute(menu.children) : false));
    }

    return includesRoute(routeStore.menus);
  }

  const canViewCustomer = computed(() => hasRouteAccess(['settings_customer']));
  const canViewProject = computed(() => hasRouteAccess(['project_project-list', 'project_data']));
  const canViewQuote = computed(() => hasRouteAccess(['project_quote']));
  const canViewPayment = computed(() => hasRouteAccess(['finance_payment']));
  const canViewPayout = computed(() => hasRouteAccess(['finance_payout']));

  const loading = ref(false);
  /** 是否已成功加载过一次（用于区分「加载中」与「确实没有数据」） */
  const loaded = ref(false);

  const customerStats = ref<Api.Customer.Stats>({ total: 0, following: 0, deal: 0, lost: 0, dealAmount: 0 });

  const projectStats = ref<Api.Project.Stats>({
    total: 0,
    pending: 0,
    building: 0,
    finished: 0,
    settled: 0,
    totalAmount: 0
  });

  const quoteStats = ref<Api.Quote.Stats>({ totalCount: 0, totalAmount: 0, projectCount: 0, avgAmount: 0 });

  const paymentStats = ref<Api.Payment.Stats>({ totalCount: 0, totalAmount: 0, monthAmount: 0, projectCount: 0 });

  /** 最近创建的项目（列表按创建时间倒序） */
  const projects = ref<Api.Project.Item[]>([]);
  /** 跟进中的客户（按创建时间倒序） */
  const followingCustomers = ref<Api.Customer.Item[]>([]);

  /** 逾期未完工项目：施工中，且预计完工日已过 */
  const overdueProjects = computed<OverdueProject[]>(() => {
    const today = dayjs().startOf('day');

    return projects.value
      .filter(item => item.status === 2 && Boolean(item.endDate))
      .map(item => ({
        ...item,
        overdueDays: today.diff(dayjs(item.endDate as string).startOf('day'), 'day')
      }))
      .filter(item => item.overdueDays > 0)
      .sort((prev, next) => next.overdueDays - prev.overdueDays);
  });

  /** 待开工项目（统计总数以 `/project/stats` 的 pending 为准，这里给出可展示的明细） */
  const pendingProjects = computed(() => projects.value.filter(item => item.status === 1));

  /** 待跟进客户 */
  const followingList = computed(() => followingCustomers.value);

  /** 最近施工进度：按进度记录时间倒序 */
  const latestProgresses = computed<ProgressItem[]>(() =>
    projects.value
      .filter(item => Boolean(item.latestProgress) && Boolean(item.latestProgressAt))
      .map(item => ({
        projectId: item.id,
        projectName: item.name,
        status: item.status,
        customerName: item.customerName,
        content: item.latestProgress as string,
        time: item.latestProgressAt as string,
        amount: item.amount ?? 0,
        managerName: item.managerName ?? ''
      }))
      .sort((prev, next) => dayjs(next.time).valueOf() - dayjs(prev.time).valueOf())
  );

  /** 首页还没有任何业务数据（客户 / 项目都为空） */
  const isEmpty = computed(
    () =>
      (canViewCustomer.value || canViewProject.value) &&
      (!canViewCustomer.value || customerStats.value.total === 0) &&
      (!canViewProject.value || projectStats.value.total === 0)
  );

  /** 加载首页全部数据：任一接口失败都不影响其它区块展示 */
  async function loadData() {
    loading.value = true;

    try {
      const [customerStatsRes, projectStatsRes, quoteStatsRes, paymentStatsRes, projectListRes, customerListRes] =
        await Promise.all([
          canViewCustomer.value ? fetchGetCustomerStats() : Promise.resolve(null),
          canViewProject.value ? fetchGetProjectStats() : Promise.resolve(null),
          canViewQuote.value ? fetchGetQuoteStats() : Promise.resolve(null),
          canViewPayment.value ? fetchGetPaymentStats() : Promise.resolve(null),
          canViewProject.value ? fetchGetProjectList({ current: 1, size: OVERVIEW_SIZE }) : Promise.resolve(null),
          canViewCustomer.value
            ? fetchGetCustomerList({ current: 1, size: OVERVIEW_SIZE, status: 1 })
            : Promise.resolve(null)
        ]);

      if (customerStatsRes && !customerStatsRes.error && customerStatsRes.data) {
        customerStats.value = customerStatsRes.data;
      }

      if (projectStatsRes && !projectStatsRes.error && projectStatsRes.data) {
        projectStats.value = projectStatsRes.data;
      }

      if (quoteStatsRes && !quoteStatsRes.error && quoteStatsRes.data) {
        quoteStats.value = quoteStatsRes.data;
      }

      if (paymentStatsRes && !paymentStatsRes.error && paymentStatsRes.data) {
        paymentStats.value = paymentStatsRes.data;
      }

      if (projectListRes && !projectListRes.error && projectListRes.data) {
        projects.value = projectListRes.data.records ?? [];
      }

      if (customerListRes && !customerListRes.error && customerListRes.data) {
        followingCustomers.value = customerListRes.data.records ?? [];
      }

      loaded.value = true;
    } finally {
      loading.value = false;
    }
  }

  return {
    loading,
    loaded,
    canViewCustomer,
    canViewProject,
    canViewQuote,
    canViewPayment,
    canViewPayout,
    customerStats,
    projectStats,
    quoteStats,
    paymentStats,
    projects,
    pendingProjects,
    overdueProjects,
    followingList,
    latestProgresses,
    isEmpty,
    loadData
  };
}
