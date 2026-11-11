<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';
import type { RouteKey } from '@elegant-router/types';
import { useRouterPush } from '@/hooks/common/router';
import { formatDate, formatMoney } from '@/views/finance/components/finance-meta';
import type { OverdueProject } from './use-home-data';

defineOptions({ name: 'HomeTodoPanel' });

interface Props {
  /** 待开工项目明细 */
  pendingProjects: Api.Project.Item[];
  /** 待开工总数（来自项目统计，比明细更准确） */
  pendingCount: number;
  /** 逾期未完工项目（已按逾期天数倒序） */
  overdueProjects: OverdueProject[];
  /** 跟进中的客户明细 */
  followingCustomers: Api.Customer.Item[];
  /** 跟进中客户总数（来自客户统计） */
  followingCount: number;
}

/** 每组待办最多展示的条数 */
const DISPLAY_LIMIT = 3;
/** 待跟进客户只展示最早的四个 */
const FOLLOWING_LIMIT = 4;

/** 单条待办 */
interface TodoItem {
  id: number;
  title: string;
  /** 次要信息：客户 / 日期 / 电话等 */
  sub: string;
  /** 右侧补充信息 */
  extra?: string;
  /** 右侧强调标记（逾期天数等） */
  badge?: string;
}

/** 一组待办 */
interface TodoGroup {
  key: string;
  title: string;
  icon: string;
  /** 主题色名，配合 rgb(var(--{color}-color)) 使用 */
  color: string;
  /** 该类待办总数（取自统计接口） */
  count: number;
  items: TodoItem[];
  /** 去处理该组待办的页面 */
  route: RouteKey;
  emptyText: string;
  /** 超过展示上限时的补充提示 */
  moreText?: string;
}

const props = defineProps<Props>();

const { routerPushByKey } = useRouterPush();

/** 按日期升序取最早的，日期为空的排最后 */
function sortByEarliest<T>(list: T[], getDate: (item: T) => string | null | undefined) {
  function toTime(item: T) {
    const value = getDate(item);

    return value ? dayjs(value).valueOf() : Number.POSITIVE_INFINITY;
  }

  return [...list].sort((prev, next) => toTime(prev) - toTime(next));
}

const groups = computed<TodoGroup[]>(() => {
  const pendingItems = sortByEarliest(props.pendingProjects, item => item.startDate)
    .slice(0, DISPLAY_LIMIT)
    .map<TodoItem>(item => ({
      id: item.id,
      title: item.name,
      sub: `${item.customerName || '未指定客户'} · 计划开工 ${formatDate(item.startDate)}`,
      extra: formatMoney(item.amount)
    }));

  /** 逾期列表已按逾期天数倒序（等价于预计完工日最早在前），这里直接截取 */
  const overdueItems = props.overdueProjects.slice(0, DISPLAY_LIMIT).map<TodoItem>(item => ({
    id: item.id,
    title: item.name,
    sub: `${item.customerName || '未指定客户'} · 预计完工 ${formatDate(item.endDate)}`,
    badge: `逾期 ${item.overdueDays} 天`
  }));

  const followingItems = sortByEarliest(props.followingCustomers, item => item.createdAt)
    .slice(0, FOLLOWING_LIMIT)
    .map<TodoItem>(item => ({
      id: item.id,
      title: item.name,
      sub: `${item.phone || '未留电话'} · ${item.source || '来源未填'}`,
      extra: `${item.projectCount ?? 0} 个项目`
    }));

  return [
    {
      key: 'pending',
      title: '待开工项目',
      icon: 'mdi:calendar-clock',
      color: 'info',
      count: props.pendingCount,
      items: pendingItems,
      route: 'project_project-list',
      emptyText: '暂无待开工项目',
      moreText: props.pendingCount > pendingItems.length ? `共 ${props.pendingCount} 个，仅列出最早 ${pendingItems.length} 个` : undefined
    },
    {
      key: 'overdue',
      title: '逾期未完工',
      icon: 'mdi:alert-circle-outline',
      color: 'error',
      count: props.overdueProjects.length,
      items: overdueItems,
      route: 'project_project-list',
      emptyText: '没有逾期项目',
      moreText:
        props.overdueProjects.length > overdueItems.length
          ? `共 ${props.overdueProjects.length} 个，仅列出逾期最久 ${overdueItems.length} 个`
          : undefined
    },
    {
      key: 'following',
      title: '待跟进客户',
      icon: 'mdi:account-clock-outline',
      color: 'warning',
      count: props.followingCount,
      items: followingItems,
      route: 'settings_customer',
      emptyText: '暂无待跟进客户',
      moreText: props.followingCount > followingItems.length ? `共 ${props.followingCount} 位，仅列出最早 ${followingItems.length} 位` : undefined
    }
  ];
});

/** 跳转到处理该组待办的页面 */
function go(route: RouteKey) {
  void routerPushByKey(route);
}

/** 三组待办合计 */
const todoTotal = computed(() => groups.value.reduce((sum, group) => sum + group.count, 0));
</script>

<template>
  <NCard :bordered="false" class="card-wrapper home-card" size="small">
    <template #header>
      <div class="card-head">
        <span class="card-head__icon card-head__icon--warning">
          <SvgIcon icon="mdi:clipboard-text-clock-outline" />
        </span>
        <span class="card-head__title">待办中心</span>
      </div>
    </template>
    <template #header-extra>
      <span v-if="todoTotal" class="todo-total">
        <i class="todo-total__dot" aria-hidden="true" />
        共 {{ todoTotal }} 项待处理
      </span>
    </template>

    <div class="todo">
      <section
        v-for="group in groups"
        :key="group.key"
        class="group"
        :style="{
          '--accent-rgb': `var(--${group.color}-color)`,
          '--accent': `rgb(var(--${group.color}-color))`,
          '--accent-soft': `rgb(var(--${group.color}-color) / 0.12)`
        }"
      >
        <header class="group__head">
          <span class="group__icon">
            <SvgIcon :icon="group.icon" />
          </span>
          <span class="group__title">{{ group.title }}</span>
          <span class="group__count">{{ group.count }}</span>
          <button class="group__more" type="button" @click="go(group.route)">
            查看
            <SvgIcon icon="mdi:chevron-right" />
          </button>
        </header>

        <ul v-if="group.items.length" class="group__list">
          <li v-for="item in group.items" :key="item.id" class="item" @click="go(group.route)">
            <span class="item__main">
              <span class="item__title">{{ item.title }}</span>
              <span class="item__sub">{{ item.sub }}</span>
            </span>
            <span v-if="item.badge" class="item__badge">{{ item.badge }}</span>
            <span v-else-if="item.extra" class="item__extra">{{ item.extra }}</span>
          </li>
        </ul>

        <p v-else class="group__empty">{{ group.emptyText }}</p>

        <p v-if="group.moreText" class="group__foot">{{ group.moreText }}</p>
      </section>
    </div>
  </NCard>
</template>

<style scoped>
.home-card {
  display: flex;
  min-height: 0;
  flex-direction: column;
}

.home-card :deep(.n-card__content) {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  flex-direction: column;
}

/* 与其他首页卡片保持一致的标题样式 */
.card-head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-head__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 8px;
  font-size: 16px;
  color: rgb(var(--primary-color));
  background-color: rgb(var(--primary-color) / 0.1);
}

.card-head__icon--warning {
  color: rgb(var(--warning-color));
  background-color: rgb(var(--warning-color) / 0.12);
}

.card-head__title {
  font-size: 15px;
  font-weight: 500;
}

/* 卡片标题右侧的待办合计 */
.todo-total {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 2px 9px;
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
  color: rgb(var(--warning-color));
  background-color: rgb(var(--warning-color) / 0.1);
}

.todo-total__dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background-color: rgb(var(--warning-color));
}

.todo {
  display: flex;
  flex: 1;
  min-height: 0;
  flex-direction: column;
  gap: 10px;
  padding-right: 4px;
  overflow-y: auto;
}

.group {
  position: relative;
  padding: 10px 12px 10px 14px;
  border: 1px solid rgb(var(--base-text-color) / 0.06);
  border-radius: 10px;
  background-color: rgb(var(--base-text-color) / 0.02);
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease;
}

/* 左侧主色竖条，三组待办一眼区分 */
.group::before {
  content: '';
  position: absolute;
  top: 10px;
  bottom: 10px;
  left: 0;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: linear-gradient(180deg, var(--accent), rgb(var(--accent-rgb) / 0.18));
}

.group:hover {
  border-color: rgb(var(--accent-rgb) / 0.25);
  background-color: rgb(var(--accent-rgb) / 0.035);
}

.group__head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.group__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 7px;
  font-size: 15px;
  color: var(--accent);
  background-color: var(--accent-soft);
}

.group__title {
  font-size: 14px;
  font-weight: 500;
  color: rgb(var(--base-text-color) / 0.85);
}

.group__count {
  min-width: 22px;
  padding: 1px 7px;
  border-radius: 100px;
  font-size: 12px;
  font-weight: 600;
  text-align: center;
  font-variant-numeric: tabular-nums;
  color: var(--accent);
  background-color: var(--accent-soft);
}

.group__more {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  margin-left: auto;
  padding: 0;
  font-size: 12px;
  color: rgb(var(--accent-rgb));
  background: none;
  border: none;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.group__more:hover {
  opacity: 0.7;
}

.group__list {
  display: flex;
  flex-direction: column;
  margin: 8px 0 0;
  padding: 0;
  list-style: none;
}

.item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  border-radius: 7px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.item:hover {
  background-color: rgb(var(--accent-rgb) / 0.07);
}

.item__main {
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}

.item__title {
  overflow: hidden;
  font-size: 13px;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.85);
}

.item__sub {
  overflow: hidden;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.45);
}

.item__extra {
  flex-shrink: 0;
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.55);
}

.item__badge {
  flex-shrink: 0;
  padding: 2px 8px;
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
  color: var(--accent);
  background-color: var(--accent-soft);
}

.group__empty,
.group__foot {
  margin: 10px 0 0;
  padding-left: 10px;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.35);
}

.group__foot {
  margin-top: 6px;
}
</style>
