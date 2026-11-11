<script setup lang="ts">
import { computed } from 'vue';
import { formatMoneyInWan } from '@/views/finance/components/finance-meta';
import { toWan } from './home-meta';

defineOptions({ name: 'HomeStatOverview' });

interface Props {
  /** 客户概览统计 */
  customer: Api.Customer.Stats;
  /** 项目概览统计 */
  project: Api.Project.Stats;
  /** 收款概览统计 */
  payment: Api.Payment.Stats;
}

/** 首页核心指标卡 */
interface StatCard {
  key: string;
  label: string;
  /** 主数值（参与数字滚动） */
  value: number;
  /** 主数值小数位 */
  decimals: number;
  /** 主数值单位，单独展示、不参与滚动 */
  unit: string;
  icon: string;
  /** 主题色名，配合 rgb(var(--{color}-color)) 使用 */
  color: string;
  /** 底部两段补充信息 */
  meta: { label: string; value: string }[];
}

const props = defineProps<Props>();

const statCards = computed<StatCard[]>(() => [
  {
    key: 'customer',
    label: '客户总数',
    value: props.customer.total ?? 0,
    decimals: 0,
    unit: '位',
    icon: 'mdi:account-multiple-outline',
    color: 'primary',
    meta: [
      { label: '跟进中', value: String(props.customer.following ?? 0) },
      { label: '已成交', value: String(props.customer.deal ?? 0) }
    ]
  },
  {
    key: 'project',
    label: '项目总数',
    value: props.project.total ?? 0,
    decimals: 0,
    unit: '个',
    icon: 'mdi:home-city-outline',
    color: 'primary',
    meta: [
      { label: '施工中', value: String(props.project.building ?? 0) },
      { label: '待开工', value: String(props.project.pending ?? 0) }
    ]
  },
  {
    key: 'contract',
    label: '合同总额',
    value: toWan(props.project.totalAmount),
    decimals: 2,
    unit: '万',
    icon: 'mdi:file-sign',
    color: 'warning',
    meta: [
      { label: '已完工', value: String(props.project.finished ?? 0) },
      { label: '已结算', value: String(props.project.settled ?? 0) }
    ]
  },
  {
    key: 'paid',
    label: '已回款',
    value: toWan(props.payment.totalAmount),
    decimals: 2,
    unit: '万',
    icon: 'mdi:cash-check',
    color: 'success',
    meta: [
      { label: '本月', value: formatMoneyInWan(props.payment.monthAmount) },
      { label: '笔数', value: String(props.payment.totalCount ?? 0) }
    ]
  }
]);
</script>

<template>
  <div class="kpi-grid">
    <article
      v-for="card in statCards"
      :key="card.key"
      class="kpi"
      :style="{
        '--accent-rgb': `var(--${card.color}-color)`,
        '--accent': `rgb(var(--${card.color}-color))`,
        '--accent-soft': `rgb(var(--${card.color}-color) / 0.12)`
      }"
    >
      <header class="kpi__head">
        <span class="kpi__icon">
          <SvgIcon :icon="card.icon" />
        </span>
        <span class="kpi__label">{{ card.label }}</span>
      </header>

      <p class="kpi__value">
        <CountTo
          :end-value="card.value"
          :decimals="card.decimals"
          separator=","
          transition="easeOutCubic"
          :duration="900"
        />
        <small class="kpi__unit">{{ card.unit }}</small>
      </p>

      <footer class="kpi__meta">
        <span v-for="item in card.meta" :key="item.label" class="kpi__meta-item">
          <em>{{ item.label }}</em>
          <b>{{ item.value }}</b>
        </span>
      </footer>
    </article>
  </div>
</template>

<style scoped>
.kpi-grid {
  display: grid;
  gap: 16px;
  /* 按可用宽度自动换行，避免窄屏四列挤压导致数值被裁切 */
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

@media (max-width: 520px) {
  .kpi-grid {
    gap: 12px;
    grid-template-columns: minmax(0, 1fr);
  }
}

.kpi {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 16px 10px;
  overflow: hidden;
  background:
    radial-gradient(110px 76px at 100% 0%, var(--accent-soft), transparent 72%),
    rgb(var(--container-bg-color));
  border: 1px solid rgb(var(--base-text-color) / 0.06);
  border-radius: 10px;
  box-shadow: 0 1px 2px rgb(0 21 41 / 0.04);
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease,
    border-color 0.25s ease;
}

/* 顶部主色高光，四张卡统一节奏 */
.kpi::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, var(--accent), transparent 62%);
  opacity: 0.55;
}

.kpi:hover {
  transform: translateY(-2px);
  border-color: rgb(var(--accent-rgb) / 0.28);
  box-shadow: 0 10px 24px rgb(var(--accent-rgb) / 0.14);
}

html.dark .kpi:hover {
  box-shadow: 0 10px 24px rgb(0 0 0 / 0.4);
}

.kpi__head {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.kpi__label {
  overflow: hidden;
  font-size: 12px;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.6);
}

.kpi__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 8px;
  font-size: 15px;
  color: var(--accent);
  background-color: var(--accent-soft);
}

.kpi__value {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  line-height: 1.1;
  letter-spacing: -0.02em;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color));
}

.kpi__unit {
  font-size: 11px;
  font-weight: 500;
  color: rgb(var(--base-text-color) / 0.45);
}

.kpi__meta {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: baseline;
  gap: 14px;
  margin-top: auto;
  padding-top: 8px;
  border-top: 1px solid rgb(var(--base-text-color) / 0.06);
}

.kpi__meta-item {
  display: inline-flex;
  align-items: baseline;
  gap: 5px;
  min-width: 0;
}

.kpi__meta-item + .kpi__meta-item {
  padding-left: 14px;
  border-left: 1px solid rgb(var(--base-text-color) / 0.06);
}

.kpi__meta-item em {
  overflow: hidden;
  font-size: 11px;
  font-style: normal;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.42);
}

.kpi__meta-item b {
  overflow: hidden;
  font-size: 12px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.82);
}
</style>
