<script setup lang="ts">
import { watch } from 'vue';
import type { ECOption } from '@/hooks/common/echarts';
import { useEcharts } from '@/hooks/common/echarts';
import { formatMoney } from '@/views/finance/components/finance-meta';
import { formatAxisMoney } from './report-meta';

defineOptions({
  name: 'ReportTrendChart'
});

interface Props {
  /** 12 个月的收支数据 */
  data: Api.FinReport.MonthlyItem[];
}

const props = defineProps<Props>();

/** 构建收支趋势图：收入 / 支出柱状 + 结余折线 */
function buildOption(): ECOption {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      backgroundColor: 'rgba(0, 0, 0, 0.76)',
      borderWidth: 0,
      textStyle: { color: '#fff', fontSize: 12 },
      valueFormatter: (value: unknown) => formatMoney(Number(value ?? 0))
    },
    legend: {
      right: 0,
      top: 0,
      icon: 'roundRect',
      itemWidth: 10,
      itemHeight: 10,
      itemGap: 16,
      textStyle: { fontSize: 12, color: 'rgba(128, 128, 128, 0.95)' }
    },
    grid: { left: 4, right: 8, top: 42, bottom: 0, containLabel: true },
    xAxis: {
      type: 'category',
      data: props.data.map(item => item.label),
      axisTick: { show: false },
      axisLine: { lineStyle: { color: 'rgba(128, 128, 128, 0.2)' } },
      axisLabel: { fontSize: 11, color: 'rgba(128, 128, 128, 0.9)' }
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { type: 'dashed', color: 'rgba(128, 128, 128, 0.16)' } },
      axisLabel: {
        fontSize: 11,
        color: 'rgba(128, 128, 128, 0.8)',
        formatter: (value: unknown) => formatAxisMoney(Number(value))
      }
    },
    series: [
      {
        name: '收入',
        type: 'bar',
        barMaxWidth: 14,
        itemStyle: { color: '#18a058', borderRadius: [4, 4, 0, 0] },
        data: props.data.map(item => item.income)
      },
      {
        name: '支出',
        type: 'bar',
        barMaxWidth: 14,
        itemStyle: { color: '#e88080', borderRadius: [4, 4, 0, 0] },
        data: props.data.map(item => item.expense)
      },
      {
        name: '结余',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 5,
        lineStyle: { width: 2, color: '#1677ff' },
        itemStyle: { color: '#1677ff' },
        data: props.data.map(item => item.net)
      }
    ]
  };
}

const { domRef, updateOptions } = useEcharts(buildOption, {
  onRender: instance => instance.hideLoading(),
  onUpdated: instance => instance.hideLoading()
});

watch(
  () => props.data,
  () => {
    void updateOptions(() => buildOption());
  },
  { deep: true }
);
</script>

<template>
  <div ref="domRef" class="report-trend-chart"></div>
</template>

<style scoped>
.report-trend-chart {
  width: 100%;
  height: 300px;
}
</style>
