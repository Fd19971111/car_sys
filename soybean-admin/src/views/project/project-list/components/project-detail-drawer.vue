<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { formatDate, formatDuration, formatMoney, getProjectProgress, getProjectStatusMeta } from './project-meta';
import { useProjectData } from './use-project-data';

defineOptions({
  name: 'ProjectDetailDrawer'
});

interface Props {
  /** 当前查看的项目ID */
  projectId?: number | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'edit', row: Api.Project.Detail): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { default: false });

const appStore = useAppStore();

const drawerWidth = computed(() => (appStore.isMobile ? '100%' : 620));

const { getDetail } = useProjectData();

const loading = ref(false);
const detail = ref<Api.Project.Detail | null>(null);

const statusMeta = computed(() => getProjectStatusMeta(detail.value?.status ?? 1));

const progress = computed(() => getProjectProgress(detail.value?.status ?? 1, detail.value?.startDate, detail.value?.endDate));

const logs = computed(() => detail.value?.logs ?? []);

const avatarStyle = computed(() => ({
  backgroundColor: `rgb(var(--${statusMeta.value.color}-color) / 0.12)`,
  color: `rgb(var(--${statusMeta.value.color}-color))`
}));

/** NTimeline 不支持 primary，这里做一次映射 */
const timelineType = computed<'default' | 'info' | 'success' | 'warning' | 'error'>(() =>
  statusMeta.value.tagType === 'primary' ? 'info' : statusMeta.value.tagType
);

watch(
  [visible, () => props.projectId],
  async ([isVisible]) => {
    if (!isVisible) {
      return;
    }

    if (!props.projectId) {
      detail.value = null;

      return;
    }

    loading.value = true;

    try {
      detail.value = await getDetail(props.projectId);
    } finally {
      loading.value = false;
    }
  },
  { immediate: true }
);

function handleClose() {
  visible.value = false;
  detail.value = null;
}

function handleEdit() {
  if (detail.value) {
    emit('edit', detail.value);
  }
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="drawerWidth">
    <NDrawerContent title="项目详情" closable :native-scrollbar="false" @close="handleClose">
      <NSpin :show="loading">
        <NEmpty v-if="!loading && !detail" size="small" class="py-40px" description="项目不存在或已被删除" />

        <template v-else-if="detail">
          <!-- 主要信息一览 -->
          <div class="flex-y-center gap-12px pb-16px">
            <div class="h-48px w-48px flex-center rounded-full text-18px font-medium" :style="avatarStyle">
              {{ detail.name.slice(0, 1) }}
            </div>
            <div class="min-w-0 flex-col gap-4px">
              <div class="flex-y-center gap-8px">
                <span class="truncate text-16px font-medium">{{ detail.name }}</span>
                <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
              </div>
              <span class="text-12px text-gray-400">
                {{ detail.customerName || '未关联客户' }} · 创建于 {{ formatDate(detail.createdAt) }}
              </span>
            </div>
          </div>

          <NDivider class="mt-0" />

          <!-- 施工进度条 -->
          <div class="mb-16px flex-col gap-6px">
            <div class="flex-y-center justify-between text-12px">
              <span class="text-gray-500">施工进度</span>
              <span class="font-medium">{{ progress }}%</span>
            </div>
            <div class="h-8px w-full overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700">
              <div
                class="h-full rounded-full transition-all"
                :style="{ width: `${progress}%`, backgroundColor: `rgb(var(--${statusMeta.color}-color))` }"
              />
            </div>
            <span class="text-12px text-gray-400">{{ statusMeta.desc }}</span>
          </div>

          <NDescriptions :column="1" label-placement="left" :label-width="116" bordered size="small">
            <NDescriptionsItem label="项目名称">{{ detail.name }}</NDescriptionsItem>
            <NDescriptionsItem label="所属客户">
              <span>{{ detail.customerName || '-' }}</span>
              <a v-if="detail.customerPhone" class="ml-8px text-primary" :href="`tel:${detail.customerPhone}`">
                {{ detail.customerPhone }}
              </a>
            </NDescriptionsItem>
            <NDescriptionsItem label="施工地址">{{ detail.address || detail.customerAddress || '-' }}</NDescriptionsItem>
            <NDescriptionsItem label="项目负责人">{{ detail.managerName || '未指派' }}</NDescriptionsItem>
            <NDescriptionsItem label="项目状态">
              <NTag :type="statusMeta.tagType" size="small" :bordered="false">{{ statusMeta.label }}</NTag>
            </NDescriptionsItem>
            <NDescriptionsItem label="合同金额">{{ formatMoney(detail.amount) }}</NDescriptionsItem>
            <NDescriptionsItem label="开工日期">{{ formatDate(detail.startDate) }}</NDescriptionsItem>
            <NDescriptionsItem label="预计完工">{{ formatDate(detail.endDate) }}</NDescriptionsItem>
            <NDescriptionsItem label="工期">{{ formatDuration(detail.startDate, detail.endDate) }}</NDescriptionsItem>
            <NDescriptionsItem label="创建时间">{{ formatDate(detail.createdAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="更新时间">{{ formatDate(detail.updatedAt) }}</NDescriptionsItem>
            <NDescriptionsItem label="备注">{{ detail.remark || '-' }}</NDescriptionsItem>
          </NDescriptions>

          <!-- 施工进度记录 -->
          <div class="mt-20px">
            <div class="mb-12px flex-y-center justify-between">
              <span class="text-14px font-medium">施工进度记录</span>
              <span class="text-12px text-gray-400">共 {{ logs.length }} 条</span>
            </div>
            <NEmpty v-if="!logs.length" size="small" description="暂无进度记录" />
            <NTimeline v-else>
              <NTimelineItem v-for="log in logs" :key="log.id" :type="timelineType" :time="log.logDate || log.createdAt">
                <div class="flex-col gap-4px">
                  <span class="text-14px">{{ log.content }}</span>
                  <span v-if="log.userName" class="text-12px text-gray-400">记录人：{{ log.userName }}</span>
                  <div v-if="log.images?.length" class="mt-4px flex flex-wrap gap-8px">
                    <NImage
                      v-for="(image, index) in log.images"
                      :key="index"
                      :src="image"
                      :width="64"
                      :height="64"
                      object-fit="cover"
                      class="rounded-6px"
                    />
                  </div>
                </div>
              </NTimelineItem>
            </NTimeline>
          </div>
        </template>
      </NSpin>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="handleClose">关闭</NButton>
          <NButton type="primary" :disabled="!detail" @click="handleEdit">编辑项目</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
