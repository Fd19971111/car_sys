<script setup lang="ts">
import { computed } from 'vue';
import { useRouterPush } from '@/hooks/common/router';
import { formatMoneyInWan } from '@/views/finance/components/finance-meta';
import { HOME_PROJECT_STATUS, toPercent } from './home-meta';

defineOptions({ name: 'HomeBusinessOverview' });

interface Props {
  /** 项目概览统计 */
  project: Api.Project.Stats;
  /** 收款概览统计 */
  payment: Api.Payment.Stats;
  /** 报价概览统计 */
  quote: Api.Quote.Stats;
}

const props = defineProps<Props>();

const { routerPushByKey } = useRouterPush();

/** 合同总额 - 已收款，为负时按 0 展示 */
const unpaidAmount = computed(() => Math.max(0, (props.project.totalAmount ?? 0) - (props.payment.totalAmount ?? 0)));

/** 回款率 */
const paidRate = computed(() => toPercent(props.payment.totalAmount, props.project.totalAmount));

/** 合同金额是否有效：为 0 时不展示回款占比条 */
const hasAmount = computed(() => (props.project.totalAmount ?? 0) > 0);

/** 各状态数量 + 占比 */
const statusItems = computed(() =>
  HOME_PROJECT_STATUS.map(item => {
    const count = props.project[item.key] ?? 0;

    return { ...item, count, percent: toPercent(count, props.project.total) };
  })
);

/** 分布条上只画有数据的状态，避免出现一排细小色块 */
const statusSegments = computed(() => statusItems.value.filter(item => item.count > 0));

/** 报价概览四项 */
const quoteItems = computed(() => [
  { key: 'amount', label: '报价总额', value: formatMoneyInWan(props.quote.totalAmount) },
  { key: 'count', label: '明细条数', value: `${props.quote.totalCount ?? 0} 条` },
  { key: 'avg', label: '单条均价', value: formatMoneyInWan(props.quote.avgAmount) },
  { key: 'project', label: '覆盖项目', value: `${props.quote.projectCount ?? 0} 个` }
]);

/** 跳转报价管理 */
function toQuote() {
  void routerPushByKey('project_quote');
}
</script>

<template>
  <NCard :bordered="false" class="card-wrapper home-card" size="small">
    <template #header>
      <div class="card-head">
        <span class="card-head__icon">
          <SvgIcon icon="mdi:chart-donut-variant" />
        </span>
        <span class="card-head__title">经营概览</span>
      </div>
    </template>
    <template #header-extra>
      <NButton text size="tiny" type="primary" @click="toQuote">报价明细</NButton>
    </template>

    <div class="biz">
      <!-- 回款进度 -->
      <section class="biz__col">
        <h4 class="biz__title">回款进度</h4>

        <div class="biz__ring-row">
          <div class="ring" :style="{ '--ring-percent': paidRate }">
            <div class="ring__inner">
              <span class="ring__value">
                {{ paidRate }}
                <small>%</small>
              </span>
              <span class="ring__label">回款率</span>
            </div>
          </div>

          <dl class="amounts">
            <div class="amount">
              <dt>合同总额</dt>
              <dd>{{ formatMoneyInWan(project.totalAmount) }}</dd>
            </div>
            <div class="amount">
              <dt><i class="amount__dot amount__dot--paid" />已回款</dt>
              <dd class="amount--paid">{{ formatMoneyInWan(payment.totalAmount) }}</dd>
            </div>
            <div class="amount">
              <dt><i class="amount__dot amount__dot--unpaid" />待回款</dt>
              <dd class="amount--unpaid">{{ formatMoneyInWan(unpaidAmount) }}</dd>
            </div>
          </dl>
        </div>

        <div
          v-if="hasAmount"
          class="paid-bar"
          :title="`已回款 ${formatMoneyInWan(payment.totalAmount)} / 合同总额 ${formatMoneyInWan(project.totalAmount)}`"
        >
          <span class="paid-bar__paid" :style="{ flexGrow: payment.totalAmount ?? 0 }" />
          <span class="paid-bar__unpaid" :style="{ flexGrow: unpaidAmount }" />
        </div>
      </section>

      <!-- 项目状态分布 -->
      <section class="biz__col">
        <h4 class="biz__title">项目状态分布</h4>

        <div class="dist">
          <span
            v-for="item in statusSegments"
            :key="item.key"
            class="dist__seg"
            :style="{ flexGrow: item.count, backgroundColor: item.color }"
            :title="`${item.label} ${item.count} 个`"
          />
        </div>

        <ul class="legend">
          <li v-for="item in statusItems" :key="item.key" class="legend__item">
            <span class="legend__dot" :style="{ backgroundColor: item.color }" />
            <span class="legend__label">{{ item.label }}</span>
            <span class="legend__value">{{ item.count }}</span>
            <span class="legend__percent">{{ item.percent }}%</span>
          </li>
        </ul>
      </section>

      <!-- 报价概览 -->
      <section class="biz__col">
        <h4 class="biz__title">报价概览</h4>

        <ul class="quote">
          <li v-for="item in quoteItems" :key="item.key" class="quote__item">
            <span class="quote__label">{{ item.label }}</span>
            <b class="quote__value">{{ item.value }}</b>
          </li>
        </ul>
      </section>
    </div>
  </NCard>
</template>

<style scoped>
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

.home-card {
  display: flex;
  min-height: 0;
  flex-direction: column;
}

.home-card :deep(.n-card__content) {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  flex-direction: column;
}

.biz {
  display: grid;
  gap: 20px;
  grid-template-columns: minmax(0, 1.15fr) minmax(0, 1fr) minmax(0, 1fr);
  /* 「已结算」的分布色：深色主题下换成浅色，否则深靛蓝压在深灰底上几乎看不出 */
  --home-settled-color: var(--primary-800-color);
}

html.dark .biz {
  --home-settled-color: var(--primary-300-color);
}

/* 三栏之间加细分隔线 */
.biz__col + .biz__col {
  padding-left: 20px;
  border-left: 1px solid rgb(var(--base-text-color) / 0.06);
}

.biz__title {
  margin: 0 0 8px;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.02em;
  color: rgb(var(--base-text-color) / 0.45);
}

.biz__ring-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
}

/* ============ 回款率圆环 ============ */
.ring {
  position: relative;
  flex-shrink: 0;
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: conic-gradient(
    from -90deg,
    rgb(var(--success-color)) calc(var(--ring-percent) * 1%),
    rgb(var(--base-text-color) / 0.07) 0
  );
  box-shadow: 0 0 0 1px rgb(var(--success-color) / 0.08);
}

.ring__inner {
  position: absolute;
  inset: 7px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: rgb(var(--container-bg-color));
  box-shadow: inset 0 0 0 1px rgb(var(--base-text-color) / 0.04);
}

.ring__value {
  font-size: 17px;
  font-weight: 600;
  line-height: 1.1;
  letter-spacing: -0.02em;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color));
}

.ring__value small {
  font-size: 11px;
  font-weight: 500;
  color: rgb(var(--base-text-color) / 0.45);
}

.ring__label {
  margin-top: 1px;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.45);
}

/* ============ 金额明细 ============ */
.amounts {
  display: flex;
  flex: 1;
  min-width: 126px;
  flex-direction: column;
  gap: 6px;
  margin: 0;
}

.amount {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 10px;
}

.amount dt {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.45);
}

.amount__dot {
  width: 6px;
  height: 6px;
  flex-shrink: 0;
  border-radius: 50%;
}

.amount__dot--paid {
  background-color: rgb(var(--success-color));
}

.amount__dot--unpaid {
  background-color: rgb(var(--error-color) / 0.7);
}

.amount dd {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.85);
}

.amount--paid {
  color: rgb(var(--success-color));
}

.amount--unpaid {
  color: rgb(var(--error-color));
}

/* ============ 回款占比条 ============ */
.paid-bar {
  display: flex;
  height: 5px;
  margin-top: 10px;
  overflow: hidden;
  border-radius: 100px;
  background-color: rgb(var(--base-text-color) / 0.07);
}

.paid-bar__paid {
  flex-basis: 0;
  background: linear-gradient(90deg, rgb(var(--success-color)), rgb(var(--success-color) / 0.68));
}

.paid-bar__unpaid {
  flex-basis: 0;
  background-color: rgb(var(--error-color) / 0.32);
}

/* ============ 状态分布条 ============ */
.dist {
  display: flex;
  height: 8px;
  overflow: hidden;
  border-radius: 100px;
  background-color: rgb(var(--base-text-color) / 0.07);
}

.dist__seg {
  min-width: 4px;
  flex-basis: 0;
  transition: flex-grow 0.3s ease;
}

.legend {
  display: grid;
  gap: 6px 14px;
  margin: 10px 0 0;
  padding: 0;
  list-style: none;
  grid-template-columns: repeat(auto-fit, minmax(104px, 1fr));
}

.legend__item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.legend__dot {
  width: 8px;
  height: 8px;
  flex-shrink: 0;
  border-radius: 3px;
}

.legend__label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.6);
}

.legend__value {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.85);
}

.legend__percent {
  width: 38px;
  text-align: right;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.4);
}

/* ============ 报价概览 ============ */
.quote {
  display: grid;
  gap: 6px;
  margin: 0;
  padding: 0;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  list-style: none;
}

.quote__item {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
  padding: 6px 8px;
  border: 1px solid rgb(var(--base-text-color) / 0.06);
  border-radius: 8px;
  background-color: rgb(var(--base-text-color) / 0.02);
}

.quote__label {
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.42);
}

.quote__value {
  overflow: hidden;
  font-size: 14px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.86);
}

/* 卡片宽度不足时三栏改为纵向排列，避免文字与圆环挤压 */
@container home-main (max-width: 760px) {
  .biz {
    grid-template-columns: minmax(0, 1fr);
    gap: 14px;
  }

  .biz__col + .biz__col {
    padding-top: 14px;
    padding-left: 0;
    border-top: 1px solid rgb(var(--base-text-color) / 0.06);
    border-left: none;
  }
}
</style>
