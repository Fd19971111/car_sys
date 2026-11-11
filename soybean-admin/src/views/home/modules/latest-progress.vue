<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';
import { useRouterPush } from '@/hooks/common/router';
import { formatDateTime, formatMoney } from '@/views/finance/components/finance-meta';
import { getProjectStatusMeta } from '@/views/project/project-list/components/project-meta';
import { progressTimeLabel } from './home-meta';
import type { ProgressItem } from './use-home-data';

defineOptions({ name: 'HomeLatestProgress' });

interface Props {
  /** 全公司最近的施工进度（已按时间倒序） */
  items: ProgressItem[];
}

/** 最多展示的条数 */
const DISPLAY_LIMIT = 4;

const props = defineProps<Props>();

const { routerPushByKey } = useRouterPush();

/** 取开始时间最早的若干条：按进度时间升序，避免首页被进度刷屏 */
const preview = computed(() =>
  [...props.items]
    .sort((prev, next) => dayjs(prev.time).valueOf() - dayjs(next.time).valueOf())
    .slice(0, DISPLAY_LIMIT)
    .map(item => {
      const meta = getProjectStatusMeta(item.status);

      return {
        ...item,
        statusLabel: meta.label,
        statusIcon: meta.icon,
        statusRgb: `var(--${meta.color}-color)`
      };
    })
);

/** 跳转项目列表（进度都在项目详情里看） */
function toProjectList() {
  void routerPushByKey('project_project-list');
}
</script>

<template>
  <NCard :bordered="false" class="card-wrapper home-card" size="small">
    <template #header>
      <div class="card-head">
        <span class="card-head__icon">
          <SvgIcon icon="mdi:progress-wrench" />
        </span>
        <span class="card-head__title">最近施工进度</span>
      </div>
    </template>
    <template #header-extra>
      <NButton text size="tiny" type="primary" @click="toProjectList">全部项目</NButton>
    </template>

    <ol v-if="preview.length" class="feed">
      <li
        v-for="(item, index) in preview"
        :key="item.projectId"
        class="feed__item"
        :style="{ '--feed-rgb': item.statusRgb }"
        @click="toProjectList"
      >
        <span class="feed__rail" :class="{ 'feed__rail--head': index === 0 }" />

        <div class="feed__body">
          <div class="feed__head">
            <span class="feed__project">{{ item.projectName }}</span>
            <span v-if="item.customerName" class="feed__customer">{{ item.customerName }}</span>
            <span class="feed__status">
              <SvgIcon :icon="item.statusIcon" class="feed__status-icon" />
              {{ item.statusLabel }}
            </span>
          </div>
          <p class="feed__content">{{ item.content }}</p>
        </div>

        <div class="feed__side">
          <time class="feed__time" :title="formatDateTime(item.time)">{{ progressTimeLabel(item.time) }}</time>
          <span class="feed__amount">{{ formatMoney(item.amount) }}</span>
          <span class="feed__manager" :title="item.managerName || '未指定负责人'">
            {{ item.managerName || '未指定负责人' }}
          </span>
        </div>
      </li>
    </ol>

    <NEmpty v-else size="small" description="还没有施工进度记录" />
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

.card-head__title {
  font-size: 15px;
  font-weight: 500;
}

.feed {
  display: flex;
  flex: 1;
  min-height: 0;
  flex-direction: column;
  margin: 0;
  padding: 0 4px 0 0;
  overflow-y: auto;
  list-style: none;
}

.feed__item {
  position: relative;
  display: flex;
  gap: 14px;
  padding: 0 10px 12px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.feed__item:last-child {
  padding-bottom: 2px;
}

.feed__item:hover {
  background-color: rgb(var(--base-text-color) / 0.03);
}

/* 时间轴竖线 + 节点 */
.feed__rail {
  position: relative;
  flex-shrink: 0;
  width: 9px;
}

.feed__rail::before {
  content: '';
  position: absolute;
  top: 18px;
  bottom: -18px;
  left: 50%;
  width: 1px;
  background-color: rgb(var(--base-text-color) / 0.09);
  transform: translateX(-50%);
}

.feed__item:last-child .feed__rail::before {
  display: none;
}

.feed__rail::after {
  content: '';
  position: absolute;
  top: 6px;
  left: 50%;
  width: 9px;
  height: 9px;
  border: 2px solid rgb(var(--feed-rgb) / 0.45);
  border-radius: 50%;
  background-color: rgb(var(--container-bg-color));
  transform: translateX(-50%);
}

/* 首条（开始时间最早）用实心节点强调 */
.feed__rail--head::after {
  border-color: rgb(var(--feed-rgb));
  background-color: rgb(var(--feed-rgb));
  box-shadow: 0 0 0 3px rgb(var(--feed-rgb) / 0.14);
}

.feed__body {
  flex: 1;
  min-width: 0;
  padding-top: 1px;
}

.feed__head {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.feed__project {
  min-width: 0;
  overflow: hidden;
  font-size: 13px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.85);
}

.feed__customer {
  flex-shrink: 0;
  padding: 1px 7px;
  border-radius: 100px;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.55);
  background-color: rgb(var(--base-text-color) / 0.05);
}

.feed__status {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  gap: 3px;
  padding: 1px 7px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
  color: rgb(var(--feed-rgb));
  background-color: rgb(var(--feed-rgb) / 0.12);
}

.feed__status-icon {
  font-size: 12px;
}

/* 右侧信息列：时间 / 金额 / 负责人，填充宽卡片右侧空白 */
.feed__side {
  display: flex;
  flex-shrink: 0;
  flex-direction: column;
  align-items: flex-end;
  gap: 3px;
  padding-top: 1px;
}

.feed__time {
  flex-shrink: 0;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.4);
}

.feed__amount {
  font-size: 13px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.75);
}

.feed__manager {
  max-width: 140px;
  overflow: hidden;
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.4);
}

.feed__content {
  display: -webkit-box;
  overflow: hidden;
  margin: 4px 0 0;
  font-size: 13px;
  line-height: 1.6;
  color: rgb(var(--base-text-color) / 0.6);
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
</style>
