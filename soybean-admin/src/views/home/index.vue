<script setup lang="ts">
import { computed, onMounted } from 'vue';
import dayjs from 'dayjs';
import type { RouteKey } from '@elegant-router/types';
import { useRouterPush } from '@/hooks/common/router';
import { useAuthStore } from '@/store/modules/auth';
import HomeBusinessOverview from './modules/business-overview.vue';
import HomeLatestProgress from './modules/latest-progress.vue';
import HomeTodoPanel from './modules/todo-panel.vue';
import { toWan } from './modules/home-meta';
import { useHomeData } from './modules/use-home-data';

const authStore = useAuthStore();

const { routerPushByKey } = useRouterPush();

const {
  loading,
  loaded,
  customerStats,
  projectStats,
  quoteStats,
  paymentStats,
  pendingProjects,
  overdueProjects,
  followingList,
  latestProgresses,
  isEmpty,
  canViewCustomer,
  canViewProject,
  canViewQuote,
  canViewPayment,
  canViewPayout,
  loadData
} = useHomeData();

const WEEKDAYS = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];

const today = dayjs();

/** 今天的日期 + 星期，例如：2026年3月8日 周日 */
const todayText = `${today.format('YYYY年M月D日')} ${WEEKDAYS[today.day()]}`;

/** 当前登录用户展示名：优先昵称 */
const displayName = computed(() => authStore.userInfo.nickname || authStore.userInfo.userName || '管理员');

/** 头像文字 */
const avatarText = computed(() => displayName.value.slice(0, 1));

/** 按当前时间问候 */
const greeting = computed(() => {
  const hour = dayjs().hour();

  if (hour < 6) return '夜深了';
  if (hour < 9) return '早上好';
  if (hour < 12) return '上午好';
  if (hour < 14) return '中午好';
  if (hour < 18) return '下午好';

  return '晚上好';
});

/** 本月回款（万元），横幅里高亮展示 */
const monthAmountInWan = computed(() => toWan(paymentStats.value.monthAmount));

/** 待办合计：待开工 + 逾期 + 待跟进 */
const todoTotal = computed(
  () => (projectStats.value.pending ?? 0) + overdueProjects.value.length + (customerStats.value.following ?? 0)
);

/** 常用入口 */
const quickActions = computed(() =>
  [
    {
      key: 'customer',
      label: '客户管理',
      icon: 'mdi:account-multiple-outline',
      color: 'primary',
      route: 'settings_customer' as RouteKey,
      allowed: canViewCustomer.value
    },
    {
      key: 'project',
      label: '项目列表',
      icon: 'mdi:home-city-outline',
      color: 'primary',
      route: 'project_project-list' as RouteKey,
      allowed: canViewProject.value
    },
    {
      key: 'quote',
      label: '报价管理',
      icon: 'mdi:file-document-outline',
      color: 'warning',
      route: 'project_quote' as RouteKey,
      allowed: canViewQuote.value
    },
    {
      key: 'payment',
      label: '收款管理',
      icon: 'mdi:cash-multiple',
      color: 'success',
      route: 'finance_payment' as RouteKey,
      allowed: canViewPayment.value
    },
    {
      key: 'payout',
      label: '付款管理',
      icon: 'mdi:cash-minus',
      color: 'error',
      route: 'finance_payout' as RouteKey,
      allowed: canViewPayout.value
    }
  ].filter(action => action.allowed)
);

/** 手动刷新首页数据 */
function refresh() {
  void loadData();
}

function go(route: RouteKey) {
  void routerPushByKey(route);
}

onMounted(() => {
  void loadData();
});
</script>

<template>
  <div class="home">
    <!-- 顶部横幅：登录人 + 今日 + 关键指标 + 常用入口 -->
    <section class="hero">
      <div class="hero__main">
        <div class="hero__profile">
          <span class="hero__avatar">{{ avatarText }}</span>
          <div class="hero__text">
            <h2 class="hero__title">{{ greeting }}，{{ displayName }}</h2>
            <p class="hero__sub">
              <SvgIcon icon="mdi:calendar-blank-outline" class="hero__sub-icon" />
              <span>{{ todayText }}</span>
              <i class="hero__sub-sep" aria-hidden="true" />
              <span>易民家装修业务概览</span>
            </p>
          </div>
        </div>

        <div class="hero__metrics">
          <div class="metric">
            <span class="metric__icon metric__icon--success">
              <SvgIcon icon="mdi:cash-check" />
            </span>
            <div class="metric__body">
              <span class="metric__label">本月回款</span>
              <span class="metric__value">
                <CountTo
                  :end-value="monthAmountInWan"
                  :decimals="2"
                  separator=","
                  transition="easeOutCubic"
                  :duration="900"
                />
                <small>万</small>
              </span>
            </div>
          </div>

          <div class="metric">
            <span class="metric__icon metric__icon--warning">
              <SvgIcon icon="mdi:clipboard-text-clock-outline" />
            </span>
            <div class="metric__body">
              <span class="metric__label">待办事项</span>
              <span class="metric__value">
                <CountTo :end-value="todoTotal" :duration="900" transition="easeOutCubic" />
                <small>项</small>
              </span>
            </div>
          </div>

          <button class="hero__refresh" type="button" :disabled="loading" title="刷新数据" @click="refresh">
            <span class="hero__refresh-icon" :class="{ 'hero__refresh-icon--spinning': loading }">
              <SvgIcon icon="mdi:refresh" />
            </span>
          </button>
        </div>
      </div>

      <nav class="hero__nav">
        <button
          v-for="action in quickActions"
          :key="action.key"
          class="entry"
          type="button"
          :style="{
            '--entry-rgb': `var(--${action.color}-color)`,
            '--entry-soft': `rgb(var(--${action.color}-color) / 0.12)`
          }"
          @click="go(action.route)"
        >
          <span class="entry__icon">
            <SvgIcon :icon="action.icon" />
          </span>
          <span class="entry__label">{{ action.label }}</span>
          <span class="entry__arrow"><SvgIcon icon="mdi:arrow-right" /></span>
        </button>

        <button
          v-if="authStore.userInfo.role === 1 && authStore.userInfo.finInitialized"
          class="entry"
          type="button"
          :style="{
            '--entry-rgb': 'var(--primary-color)',
            '--entry-soft': 'rgb(var(--primary-color) / 0.12)'
          }"
          @click="routerPushByKey('finance_setup', { query: { reconfigure: '1' } })"
        >
          <span class="entry__icon"><SvgIcon icon="mdi:cog-sync-outline" /></span>
          <span class="entry__label">财务基础配置</span>
          <span class="entry__arrow"><SvgIcon icon="mdi:arrow-right" /></span>
        </button>
      </nav>
    </section>

    <!-- 没有任何客户 / 项目时给出引导，避免首页一片 0 看不出问题 -->
    <NAlert v-if="loaded && isEmpty" :bordered="false" type="info">
      还没有客户和项目数据：先去「客户管理」录入客户，再到「项目列表」建项目、记进度，首页会自动汇总。
    </NAlert>

    <NSpin class="home__spin" :show="loading && !loaded">
      <div class="home__stack">
        <div class="home__dashboard">
          <div class="home__main-col">
            <HomeBusinessOverview :project="projectStats" :payment="paymentStats" :quote="quoteStats" />
            <HomeLatestProgress :items="latestProgresses" />
          </div>

          <div class="home__aside-col">
            <HomeTodoPanel
              :pending-projects="pendingProjects"
              :pending-count="projectStats.pending"
              :overdue-projects="overdueProjects"
              :following-customers="followingList"
              :following-count="customerStats.following"
            />
          </div>
        </div>
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  width: 100%;
  min-width: 0;
  /* 首页铺满内容区，不产生外层滚动条，滚动交给内部列表 */
  height: 100%;
  min-height: 0;
  overflow: hidden;
  flex-direction: column;
  gap: 12px;
}

/* ============ 顶部横幅 ============ */
.hero {
  position: relative;
  flex-shrink: 0;
  padding: 18px 22px 16px;
  overflow: hidden;
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 14px;
  background:
    radial-gradient(900px 300px at 92% -40%, rgb(var(--primary-color) / 0.16), transparent 62%),
    radial-gradient(560px 240px at 6% 140%, rgb(var(--success-color) / 0.08), transparent 68%),
    rgb(var(--container-bg-color));
  box-shadow:
    0 1px 2px rgb(0 21 41 / 0.04),
    0 10px 30px rgb(0 21 41 / 0.035);
}

/* 顶部细高光条，呼应主色 */
.hero::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 46%;
  height: 2px;
  background: linear-gradient(90deg, rgb(var(--primary-color)) 0%, rgb(var(--primary-color) / 0.25) 72%, transparent 100%);
}

html.dark .hero {
  box-shadow:
    0 1px 2px rgb(0 0 0 / 0.3),
    0 10px 30px rgb(0 0 0 / 0.25);
}

.hero__main {
  position: relative;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
}

.hero__profile {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 14px;
}

.hero__avatar {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 14px;
  font-size: 19px;
  font-weight: 600;
  letter-spacing: 0.02em;
  color: #fff;
  background: linear-gradient(140deg, rgb(var(--primary-color)), rgb(var(--primary-color) / 0.68));
  box-shadow:
    0 6px 16px rgb(var(--primary-color) / 0.26),
    inset 0 1px 0 rgb(255 255 255 / 0.4);
}

.hero__text {
  min-width: 0;
}

.hero__title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  line-height: 1.35;
  letter-spacing: -0.01em;
  color: rgb(var(--base-text-color) / 0.92);
}

.hero__sub {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  margin: 4px 0 0;
  font-size: 12px;
  line-height: 1.5;
  color: rgb(var(--base-text-color) / 0.46);
}

.hero__sub-icon {
  font-size: 14px;
  opacity: 0.8;
}

.hero__sub-sep {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background-color: rgb(var(--base-text-color) / 0.25);
}

.hero__metrics {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 12px;
}

.metric {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px 8px 10px;
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 11px;
  background-color: rgb(var(--base-text-color) / 0.025);
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease;
}

.metric:hover {
  border-color: rgb(var(--base-text-color) / 0.12);
  background-color: rgb(var(--base-text-color) / 0.04);
}

.metric__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 9px;
  font-size: 18px;
}

.metric__icon--success {
  color: rgb(var(--success-color));
  background-color: rgb(var(--success-color) / 0.12);
}

.metric__icon--warning {
  color: rgb(var(--warning-color));
  background-color: rgb(var(--warning-color) / 0.12);
}

.metric__body {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.metric__label {
  font-size: 11px;
  line-height: 1.4;
  color: rgb(var(--base-text-color) / 0.45);
}

.metric__value {
  display: flex;
  align-items: baseline;
  gap: 3px;
  font-size: 20px;
  font-weight: 600;
  line-height: 1.2;
  letter-spacing: -0.02em;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.9);
}

.metric__value small {
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0;
  color: rgb(var(--base-text-color) / 0.45);
}

.hero__refresh {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  padding: 0;
  border: 1px solid rgb(var(--base-text-color) / 0.1);
  border-radius: 10px;
  color: rgb(var(--base-text-color) / 0.55);
  background-color: rgb(var(--base-text-color) / 0.03);
  cursor: pointer;
  transition:
    color 0.2s ease,
    background-color 0.2s ease,
    border-color 0.2s ease;
}

.hero__refresh:hover:not(:disabled) {
  color: rgb(var(--primary-color));
  border-color: rgb(var(--primary-color) / 0.35);
  background-color: rgb(var(--primary-color) / 0.06);
}

.hero__refresh:disabled {
  cursor: default;
  opacity: 0.7;
}

.hero__refresh-icon {
  display: flex;
  font-size: 18px;
}

.hero__refresh-icon--spinning {
  animation: hero-spin 0.9s linear infinite;
}

@keyframes hero-spin {
  to {
    transform: rotate(360deg);
  }
}

/* ============ 常用入口 ============ */
.hero__nav {
  display: grid;
  gap: 10px;
  margin-top: 14px;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
}

.entry {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.74);
  background-color: rgb(var(--base-text-color) / 0.022);
  border: 1px solid rgb(var(--base-text-color) / 0.07);
  border-radius: 10px;
  cursor: pointer;
  transition:
    color 0.2s ease,
    background-color 0.2s ease,
    border-color 0.2s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.entry:hover {
  color: rgb(var(--base-text-color) / 0.92);
  background-color: var(--entry-soft);
  border-color: rgb(var(--entry-rgb) / 0.35);
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgb(var(--entry-rgb) / 0.12);
}

.entry__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 9px;
  font-size: 17px;
  color: rgb(var(--entry-rgb));
  background-color: var(--entry-soft);
}

.entry__label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.entry__arrow {
  display: flex;
  flex-shrink: 0;
  font-size: 14px;
  color: rgb(var(--entry-rgb));
  opacity: 0;
  transform: translateX(-4px);
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.entry:hover .entry__arrow {
  opacity: 0.75;
  transform: translateX(0);
}

/* ============ 主体 ============ */
.home__spin,
.home__spin :deep(.n-spin-container),
.home__spin :deep(.n-spin-content) {
  display: flex;
  width: 100%;
  min-width: 0;
  min-height: 0;
  flex: 1;
  flex-direction: column;
}

.home__stack {
  display: flex;
  width: 100%;
  min-width: 0;
  min-height: 0;
  height: 100%;
  flex-direction: column;
  gap: 12px;
}

.home__dashboard {
  display: grid;
  flex: 1;
  min-height: 0;
  align-items: stretch;
  gap: 12px;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
  grid-template-rows: minmax(0, 1fr);
}

.home__main-col {
  display: flex;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  gap: 12px;
  /* 作为容器，让「经营概览」按实际可用宽度切换栏数，而不是按窗口宽度 */
  container-type: inline-size;
  container-name: home-main;
}

.home__main-col > * {
  min-height: 0;
}

/* 让进度卡片吃掉剩余高度，并保证至少能完整展示两到三条进度 */
.home__main-col > *:last-child {
  flex: 1;
  min-height: 160px;
}

.home__aside-col {
  display: flex;
  min-width: 0;
  min-height: 0;
}

.home__aside-col > * {
  flex: 1;
  min-width: 0;
  min-height: 0;
}

/* 统一首页卡片外观：圆角、细边、轻投影 */
.home :deep(.home-card) {
  border: 1px solid rgb(var(--base-text-color) / 0.06);
  border-radius: 12px;
  box-shadow: 0 1px 2px rgb(0 21 41 / 0.04);
}

/* 压缩卡片头尾留白，把高度让给内容区 */
.home :deep(.home-card .n-card-header) {
  padding: 10px 16px 6px;
}

.home :deep(.home-card .n-card__content) {
  padding: 10px 16px 12px;
}

html.dark .home :deep(.home-card) {
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.3);
}

/* 窄屏：整体改为纵向排列，由主体区域内部滚动，卡片保持自然高度 */
@media (max-width: 1100px) {
  .home__dashboard {
    display: flex;
    flex-direction: column;
    overflow-y: auto;
  }

  .home__main-col,
  .home__aside-col {
    flex: 0 0 auto;
    gap: 12px;
  }

  .home__main-col > *:last-child,
  .home__aside-col > * {
    flex: 0 0 auto;
  }
}

@media (max-width: 720px) {
  .hero {
    padding: 14px 16px;
  }

  .hero__metrics {
    width: 100%;
    justify-content: space-between;
  }

  .home,
  .home__stack,
  .home__main-col,
  .home__dashboard {
    gap: 12px;
  }
}
</style>
