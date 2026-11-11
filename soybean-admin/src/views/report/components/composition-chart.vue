<script setup lang="ts">
import { watch } from 'vue';
import type { ECOption } from '@/hooks/common/echarts';
import { useEcharts } from '@/hooks/common/echarts';
import { CHART_COLORS, formatMoneyCompact, formatPercent } from './report-meta';

defineOptions({
  name: 'ReportCompositionChart'
});

interface Props {
  /** 收入或支出的科目构成 */
  items: Api.FinReport.AccountAmount[];
  /** 同类合计金额 */
  total: number;
}

const props = defineProps<Props>();

/** 构建环形图：只保留金额大于 0 的科目 */
function buildOption(): ECOption {
  const data = props.items
    .filter(item => item.amount > 0)
    .map((item, index) => ({
      name: item.name,
      value: item.amount,
      itemStyle: { color: CHART_COLORS[index % CHART_COLORS.length] }
    }));

  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(0, 0, 0, 0.76)',
      borderWidth: 0,
      textStyle: { color: '#fff', fontSize: 12 },
      valueFormatter: (value: unknown) => formatMoneyCompact(Number(value ?? 0))
    },
    series: [
      {
        type: 'pie',
        radius: ['62%', '82%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: true,
        padAngle: 2,
        itemStyle: { borderRadius: 6, borderColor: 'transparent', borderWidth: 1 },
        label: { show: false },
        labelLine: { show: false },
        emphasis: { scale: true, scaleSize: 4 },
        data
      }
    ]
  };
}

const { domRef, updateOptions } = useEcharts(buildOption, {
  onRender: instance => instance.hideLoading(),
  onUpdated: instance => instance.hideLoading()
});

watch(
  () => props.items,
  () => {
    void updateOptions(() => buildOption());
  },
  { deep: true }
);
</script>

<template>
  <div class="composition">
    <div class="composition__chart">
      <div ref="domRef" class="composition__canvas"></div>
      <div class="composition__center">
        <span class="composition__center-label">合计</span>
        <strong class="composition__center-value">{{ formatMoneyCompact(total) }}</strong>
      </div>
    </div>

    <ul class="composition__list">
      <li v-for="(item, index) in items" :key="item.code" class="composition__item">
        <span class="composition__dot" :style="{ backgroundColor: CHART_COLORS[index % CHART_COLORS.length] }"></span>
        <span class="composition__name" :title="`${item.code} ${item.name}`">{{ item.name }}</span>
        <span class="composition__percent">{{ formatPercent(item.percent) }}</span>
        <span class="composition__amount">{{ formatMoneyCompact(item.amount) }}</span>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.composition {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.composition__chart {
  position: relative;
  width: 100%;
  height: 190px;
}

.composition__canvas {
  width: 100%;
  height: 100%;
}

.composition__center {
  position: absolute;
  top: 50%;
  left: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  pointer-events: none;
  transform: translate(-50%, -50%);
}

.composition__center-label {
  font-size: 11px;
  color: rgb(var(--base-text-color) / 0.45);
}

.composition__center-value {
  font-size: 16px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.85);
}

.composition__list {
  display: flex;
  flex-direction: column;
  gap: 2px;
  max-height: 148px;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  list-style: none;
}

.composition__item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 6px;
  font-size: 12px;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.composition__item:hover {
  background-color: rgb(var(--base-text-color) / 0.04);
}

.composition__dot {
  flex-shrink: 0;
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

.composition__name {
  overflow: hidden;
  color: rgb(var(--base-text-color) / 0.75);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.composition__percent {
  flex-shrink: 0;
  margin-left: auto;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--base-text-color) / 0.4);
}

.composition__amount {
  flex-shrink: 0;
  min-width: 72px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  text-align: right;
  color: rgb(var(--base-text-color) / 0.85);
}
</style>
