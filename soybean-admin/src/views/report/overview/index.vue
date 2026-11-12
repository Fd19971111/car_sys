<script setup lang="ts">
import { onMounted, watch } from 'vue';
import ReportToolbar from '../components/report-toolbar.vue';
import ReportOverviewPanel from '../components/overview-panel.vue';
import { useReportData } from '../components/use-report-data';
import { useReportYear } from '../components/use-report-year';

defineOptions({
  name: 'ReportOverview'
});

const { year } = useReportYear();
const { overviewLoading, overview, loadOverview } = useReportData();

async function load() {
  await loadOverview(year.value);
}

watch(year, load);

onMounted(load);
</script>

<template>
  <div class="report-page">
    <ReportToolbar
      v-model:year="year"
      hint="总账收支按已记账凭证统计；行程和司机结算按业务记录统计"
      :loading="overviewLoading"
      @refresh="load"
    />

    <ReportOverviewPanel :overview="overview" :loading="overviewLoading" :year="year" />
  </div>
</template>

<style scoped>
.report-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
}

@media (max-width: 768px) {
  .report-page {
    padding: 12px;
  }
}
</style>
