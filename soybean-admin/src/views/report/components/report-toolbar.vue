<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';

defineOptions({
  name: 'ReportToolbar'
});

interface Props {
  /** 当前报表的数据口径说明 */
  hint: string;
  /** 当前统计年份 */
  year: number;
  /** 数据加载中 */
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
});

const emit = defineEmits<{
  (e: 'update:year', value: number): void;
  (e: 'refresh'): void;
}>();

const yearOptions = computed(() => {
  const current = dayjs().year();

  return [current, current - 1, current - 2].map(item => ({ label: `${item} 年`, value: item }));
});
</script>

<template>
  <div class="report-toolbar">
    <p class="report-toolbar__hint">
      <SvgIcon icon="mdi:information-outline" />
      <span>{{ props.hint }}</span>
    </p>

    <div class="report-toolbar__actions">
      <NSelect
        :value="props.year"
        :options="yearOptions"
        size="small"
        class="w-110px"
        @update:value="emit('update:year', $event)"
      />
      <NButton size="small" :loading="props.loading" @click="emit('refresh')">
        <template #icon>
          <SvgIcon icon="mdi:refresh" />
        </template>
        刷新
      </NButton>
    </div>
  </div>
</template>

<style scoped>
.report-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 10px 16px;
  padding: 0 2px;
}

.report-toolbar__hint {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
  font-size: 12px;
  color: rgb(var(--base-text-color) / 0.45);
}

.report-toolbar__actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 8px;
}

@media (max-width: 640px) {
  .report-toolbar {
    align-items: flex-start;
  }

  .report-toolbar__actions {
    width: 100%;
  }

  .report-toolbar__actions :deep(.n-select) {
    flex: 1;
  }
}
</style>
