<script setup lang="ts">
import { computed, ref } from 'vue';
import dayjs from 'dayjs';
import { useRouterPush } from '@/hooks/common/router';
import { formatMoney } from '@/views/finance/components/finance-meta';
import ReportCompositionChart from './composition-chart.vue';
import ReportTrendChart from './trend-chart.vue';
import { collapseItems, formatMoneyCompact, formatPercent } from './report-meta';

defineOptions({
  name: 'ReportOverviewPanel'
});

interface Props {
  /** 报表总览数据 */
  overview: Api.FinReport.Overview | null;
  loading: boolean;
  /** 当前统计年份 */
  year: number;
}

const props = defineProps<Props>();

const { routerPushByKey } = useRouterPush();

const compositionType = ref<'income' | 'expense'>('income');

const summary = computed(() => props.overview?.summary);

const monthly = computed(() => props.overview?.monthly ?? []);

/** 是否有任何可展示的财务数据 */
const hasAnyData = computed(() => {
  const s = summary.value;

  return Boolean(
    s &&
      (s.income > 0 ||
        s.expense > 0 ||
        s.totalBalance !== 0 ||
        s.contractAmount > 0)
  );
});

const hasTrendData = computed(() => monthly.value.some(item => item.income > 0 || item.expense > 0));

/** 有收款 / 合同数据，但还没有任何记账凭证 */
const pendingVoucher = computed(() => {
  const s = summary.value;

  return Boolean(s && s.contractAmount > 0 && s.income === 0 && s.expense === 0 && s.totalBalance === 0);
});

const compositionItems = computed(() => {
  if (!props.overview) return [];

  return compositionType.value === 'income' ? props.overview.incomeItems : props.overview.expenseItems;
});

const compositionTotal = computed(() => compositionItems.value.reduce((total, item) => total + item.amount, 0));

const compositionList = computed(() => collapseItems(compositionItems.value));

const fundAccounts = computed(() => props.overview?.fundAccounts ?? []);

const yearLabel = computed(() => (props.year === dayjs().year() ? '本年' : `${props.year} 年`));

const yearReceivedLabel = computed(() => (props.year === dayjs().year() ? '本年回款' : `${props.year} 年回款`));

/** 进入收款管理，方便业务人员把线下收款补录进来 */
function goPayment() {
  void routerPushByKey('finance_payment');
}

/** 进入付款管理，登记供应商 / 人工等支出 */
function goPayout() {
  void routerPushByKey('finance_payout');
}

/** 期间状态的中文说明 */
const periodText = computed(() => {
  const status = props.overview?.periodStatus;

  if (status === 'OPEN') return '记账中';
  if (status === 'CLOSED') return '已结账';
  if (status === 'FUTURE') return '未启用';

  return '未启用';
});

const quarterTip = computed(() => {
  const s = summary.value;

  if (!s || s.income <= 0) {
    return '收入大于 0 时，这里会显示结余率';
  }

  return `每 100 元收入，结余 ${s.netRate.toFixed(1)} 元`;
});

const statCards = computed(() => {
  const s = summary.value;

  return [
    {
      key: 'income',
      label: `${yearLabel.value}收入`,
      value: formatMoneyCompact(s?.income),
      hint: s && s.income > 0 ? `月均 ${formatMoneyCompact(s.income / 12)}` : '还没有收入记录',
      icon: 'mdi:trending-up',
      color: '#18a058',
      bg: 'rgba(24, 160, 88, 0.12)'
    },
    {
      key: 'expense',
      label: `${yearLabel.value}支出`,
      value: formatMoneyCompact(s?.expense),
      hint: s && s.expense > 0 ? `月均 ${formatMoneyCompact(s.expense / 12)}` : '还没有支出记录',
      icon: 'mdi:trending-down',
      color: '#e34d59',
      bg: 'rgba(227, 77, 89, 0.12)'
    },
    {
      key: 'net',
      label: `${yearLabel.value}结余`,
      value: formatMoneyCompact(s?.net),
      hint: quarterTip.value,
      icon: 'mdi:scale-balance',
      color: '#1677ff',
      bg: 'rgba(22, 119, 255, 0.12)'
    },
    {
      key: 'balance',
      label: '当前资金余额',
      value: formatMoneyCompact(s?.totalBalance),
      hint: `现金 ${formatMoneyCompact(s?.cashBalance)} · 银行 ${formatMoneyCompact(s?.bankBalance)}`,
      icon: 'mdi:wallet-outline',
      color: '#722ed1',
      bg: 'rgba(114, 46, 209, 0.12)'
    }
  ];
});

/** 回款进度百分比，限制在 0-100 */
const receivedRate = computed(() => Math.min(100, Math.max(0, summary.value?.receivedRate ?? 0)));
</script>

<template>
  <NSpin :show="loading" description="正在汇总财务数据...">
    <div class="overview">
      <NAlert v-if="overview && !overview.initialized" type="warning" :bordered="false" class="tip">
        财务还没有完成初始化，报表暂时无法统计。请先在「财务初始化」中完成基础配置。
      </NAlert>

      <NAlert v-else-if="overview && !hasAnyData" type="info" :bordered="false" class="tip">
        还没有财务数据。完成记账凭证后，这里会自动汇总收支、结余与资金余额。
      </NAlert>

      <NAlert v-else-if="overview && pendingVoucher" type="info" :bordered="false" class="tip">
        <div class="tip__row">
          <span>
            报表按记账凭证统计。历史收款记录还没有生成凭证，所以收入、支出和资金余额暂显示为 0；从现在起，新登记的收款与付款会按借贷规则自动记账，并在这里自动汇总。
          </span>
          <NButton text type="primary" size="tiny" @click="goPayment">
            去登记收款
            <template #icon>
              <SvgIcon icon="mdi:arrow-right" />
            </template>
          </NButton>
          <NButton text type="primary" size="tiny" @click="goPayout">
            去登记付款
            <template #icon>
              <SvgIcon icon="mdi:arrow-right" />
            </template>
          </NButton>
        </div>
      </NAlert>

      <section class="stat-grid">
        <article
          v-for="card in statCards"
          :key="card.key"
          class="stat-card"
          :style="{ '--accent': card.color, '--accent-soft': card.bg }"
        >
          <div class="stat-card__icon">
            <SvgIcon :icon="card.icon" />
          </div>
          <div class="stat-card__text">
            <span class="stat-card__label">{{ card.label }}</span>
            <strong class="stat-card__value">{{ card.value }}</strong>
            <span class="stat-card__hint">{{ card.hint }}</span>
          </div>
        </article>
      </section>

      <section class="chart-grid">
        <NCard :bordered="false" class="panel" size="small">
          <div class="panel__head">
            <div>
              <h3 class="panel__title">收支趋势</h3>
              <p class="panel__subtitle">每个月收入多少、花掉多少，一眼看清</p>
            </div>
            <span class="period-chip">
              <span class="period-chip__dot" :class="{ 'period-chip__dot--closed': overview?.periodStatus !== 'OPEN' }"></span>
              {{ overview?.currentPeriod || '—' }} · {{ periodText }}
            </span>
          </div>

          <ReportTrendChart v-if="hasTrendData" :data="monthly" />
          <div v-else class="empty-state">
            <SvgIcon icon="mdi:chart-bar-stacked" />
            <p>当前年份还没有收支数据</p>
          </div>
        </NCard>

        <NCard :bordered="false" class="panel" size="small">
          <div class="panel__head">
            <div>
              <h3 class="panel__title">收支构成</h3>
              <p class="panel__subtitle">钱主要从哪来、又花到哪里去了</p>
            </div>
            <div class="mini-switch">
              <button
                type="button"
                class="mini-switch__btn"
                :class="{ 'mini-switch__btn--active': compositionType === 'income' }"
                @click="compositionType = 'income'"
              >
                收入
              </button>
              <button
                type="button"
                class="mini-switch__btn"
                :class="{ 'mini-switch__btn--active': compositionType === 'expense' }"
                @click="compositionType = 'expense'"
              >
                支出
              </button>
            </div>
          </div>

          <ReportCompositionChart v-if="compositionList.length" :items="compositionList" :total="compositionTotal" />
          <div v-else class="empty-state empty-state--compact">
            <SvgIcon :icon="compositionType === 'income' ? 'mdi:cash-plus' : 'mdi:cash-minus'" />
            <p>{{ compositionType === 'income' ? '当前年份还没有收入' : '当前年份还没有支出' }}</p>
          </div>
        </NCard>
      </section>

      <section class="detail-grid">
        <NCard :bordered="false" class="panel" size="small">
          <div class="panel__head">
            <div>
              <h3 class="panel__title">项目回款</h3>
              <p class="panel__subtitle">合同签了多少、客户已经付了多少</p>
            </div>
            <span class="panel__badge">{{ formatPercent(receivedRate) }}</span>
          </div>

          <div class="progress">
            <div class="progress__bar" :style="{ width: `${receivedRate}%` }"></div>
          </div>

          <div class="collection">
            <div class="collection__item">
              <span>合同总额</span>
              <strong>{{ formatMoney(summary?.contractAmount) }}</strong>
            </div>
            <div class="collection__item">
              <span>累计回款</span>
              <strong class="amount--success">{{ formatMoney(summary?.receivedAmount) }}</strong>
            </div>
            <div class="collection__item">
              <span>{{ yearReceivedLabel }}</span>
              <strong class="amount--info">{{ formatMoney(summary?.yearReceived) }}</strong>
            </div>
            <div class="collection__item">
              <span>待回款</span>
              <strong class="amount--error">{{ formatMoney(summary?.unpaidAmount) }}</strong>
            </div>
          </div>

          <p class="panel__note">回款率 = 累计回款 ÷ 合同总额；回款按实际收款日期统计。</p>
        </NCard>

        <NCard :bordered="false" class="panel" size="small">
          <div class="panel__head">
            <div>
              <h3 class="panel__title">钱在哪里</h3>
              <p class="panel__subtitle">现金与银行账户的当前余额</p>
            </div>
            <span class="panel__badge panel__badge--muted">{{ fundAccounts.length }} 个账户</span>
          </div>

          <ul v-if="fundAccounts.length" class="fund-list">
            <li v-for="item in fundAccounts" :key="item.id" class="fund-item">
              <span class="fund-item__icon" :class="item.type === 'CASH' ? 'fund-item__icon--cash' : 'fund-item__icon--bank'">
                <SvgIcon :icon="item.type === 'CASH' ? 'mdi:cash' : 'mdi:bank-outline'" />
              </span>
              <div class="fund-item__text">
                <strong>{{ item.name }}</strong>
                <span>{{ item.bankName ? `${item.bankName} · 尾号 ${item.accountLast4 || '—'}` : item.accountName }}</span>
              </div>
              <strong class="fund-item__balance">{{ formatMoney(item.balance) }}</strong>
            </li>
          </ul>

          <div v-else class="empty-state empty-state--compact">
            <SvgIcon icon="mdi:wallet-outline" />
            <p>还没有启用的资金账户</p>
          </div>

          <p v-if="fundAccounts.length && summary && summary.totalBalance === 0" class="panel__note">
            还没有记账凭证，账户余额会在记账后更新；当前显示 0 不代表账户没有钱。
          </p>
        </NCard>
      </section>
    </div>
  </NSpin>
</template>

<style scoped>
.overview {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tip {
  border-radius: 10px;
}

.tip__row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 12px;
}

.stat-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 18px;
  overflow: hidden;
  background-color: rgb(var(--container-bg-color));
  border: 1px solid rgb(var(--base-text-color) / 0.06);
  border-radius: 12px;
  box-shadow: 0 1px 2px rgb(0 21 41 / 0.04);
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease,
    border-color 0.25s ease;
}

.stat-card::after {
  content: '';
  position: absolute;
  top: -40px;
  right: -40px;
  width: 108px;
  height: 108px;
  background: var(--accent-soft);
  border-radius: 50%;
  pointer-events: none;
}

.stat-card:hover {
  transform: translateY(-2px);
  border-color: rgb(var(--base-text-color) / 0.1);
  box-shadow: 0 8px 24px rgb(0 21 41 / 0.1);
}

html.dark .stat-card:hover {
  box-shadow: 0 8px 24px rgb(0 0 0 / 0.45);
}

.stat-card__icon {
  position: relative;
  z-index: 1;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  font-size: 22px;
  color: var(--accent);
  background-color: var(--accent-soft);
  border-radius: 12px;
}

.stat-card__text {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.stat-card__label {
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.5);
}

.stat-card__value {
  font-size: 22px;
  font-weight: 600;
  line-height: 1.3;
  letter-spacing: -0.01em;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color));
}

.stat-card__hint {
  overflow: hidden;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.4);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chart-grid {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
  gap: 16px;
}

.detail-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 16px;
}

.panel {
  border-radius: 12px;
}

.panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.panel__title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: rgb(var(--base-text-color) / 0.9);
}

.panel__subtitle {
  margin: 3px 0 0;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.45);
}

.panel__badge {
  flex-shrink: 0;
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--primary-color));
  background-color: rgb(var(--primary-color) / 0.1);
  border-radius: 999px;
}

.panel__badge--muted {
  font-weight: 500;
  color: rgb(var(--base-text-color) / 0.55);
  background-color: rgb(var(--base-text-color) / 0.06);
}

.panel__note {
  margin: 12px 0 0;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.38);
}

.period-chip {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  gap: 6px;
  padding: 3px 10px;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.55);
  background-color: rgb(var(--base-text-color) / 0.05);
  border-radius: 999px;
}

.period-chip__dot {
  width: 6px;
  height: 6px;
  background-color: rgb(var(--success-color));
  border-radius: 50%;
}

.period-chip__dot--closed {
  background-color: rgb(var(--base-text-color) / 0.35);
}

.mini-switch {
  display: flex;
  flex-shrink: 0;
  gap: 2px;
  padding: 2px;
  background-color: rgb(var(--base-text-color) / 0.05);
  border-radius: 8px;
}

.mini-switch__btn {
  padding: 4px 12px;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.55);
  cursor: pointer;
  background: transparent;
  border: none;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.mini-switch__btn--active {
  font-weight: 600;
  color: rgb(var(--primary-color));
  background-color: rgb(var(--primary-color) / 0.12);
}

.progress {
  height: 8px;
  overflow: hidden;
  background-color: rgb(var(--base-text-color) / 0.06);
  border-radius: 999px;
}

.progress__bar {
  height: 100%;
  background: linear-gradient(90deg, #18a058, #36cfc9);
  border-radius: 999px;
  transition: width 0.6s ease;
}

.collection {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin-top: 14px;
}

.collection__item {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
  padding: 8px 10px;
  background-color: rgb(var(--base-text-color) / 0.03);
  border-radius: 8px;
}

.collection__item span {
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.45);
}

.collection__item strong {
  overflow: hidden;
  font-size: 13px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.85);
}

.amount--success {
  color: #18a058;
}

.amount--info {
  color: rgb(var(--primary-color));
}

.amount--error {
  color: #e34d59;
}

.latest-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.latest-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 8px;
  transition: background-color 0.2s ease;
}

.latest-item:hover {
  background-color: rgb(var(--base-text-color) / 0.03);
}

.latest-item__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  font-size: 16px;
  border-radius: 8px;
}

.latest-item__text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.latest-item__text strong {
  overflow: hidden;
  font-size: 13px;
  font-weight: 600;
  color: rgb(var(--base-text-color) / 0.82);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.latest-item__text span {
  overflow: hidden;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.42);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.latest-item__amount {
  flex-shrink: 0;
  margin-left: auto;
  font-size: 14px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.latest-item__amount--income {
  color: #18a058;
}

.latest-item__amount--expense {
  color: #e34d59;
}

.fund-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.fund-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  background-color: rgb(var(--base-text-color) / 0.03);
  border-radius: 8px;
}

.fund-item__icon {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  font-size: 16px;
  border-radius: 8px;
}

.fund-item__icon--cash {
  color: #18a058;
  background-color: rgb(24 160 88 / 0.12);
}

.fund-item__icon--bank {
  color: #1677ff;
  background-color: rgb(22 119 255 / 0.12);
}

.fund-item__text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.fund-item__text strong {
  overflow: hidden;
  font-size: 13px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--base-text-color) / 0.85);
}

.fund-item__text span {
  overflow: hidden;
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.42);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fund-item__balance {
  flex-shrink: 0;
  margin-left: auto;
  font-size: 14px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color));
}

.balance-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.balance-chip {
  display: inline-flex;
  align-items: baseline;
  gap: 6px;
  padding: 5px 12px;
  font-size: 12px;
  border-radius: 8px;
}

.balance-chip strong {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 72px 0;
  font-size: 34px;
  color: rgb(var(--base-text-color) / 0.22);
}

.empty-state--compact {
  padding: 34px 0;
  font-size: 30px;
}

.empty-state p {
  margin: 0;
  font-size: 13px;
  color: rgb(var(--base-text-color) / 0.45);
}

@media (max-width: 1100px) {
  .chart-grid,
  .detail-grid {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 640px) {
  .collection {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
