<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import type { UploadCustomRequestOptions, UploadFileInfo } from 'naive-ui';
import { fetchGetPrintSetting, fetchSavePrintSetting } from '@/service/api/print';
import { fetchUploadAdminImage } from '@/service/api/upload';
import { invalidatePrintSetting, renderPrintHtml, saveHtmlAsPdf, type PrintExportColumn } from '@/utils/print-export';

defineOptions({ name: 'SettingsPrint' });

const loading = ref(false);
const saving = ref(false);

const model = reactive<Api.Print.PrintSettingForm>({
  title: '',
  subTitle: '',
  footer: '',
  showDate: true,
  showSeal: false,
  sealImage: ''
});

/** 印章图片预览地址 */
const sealPreviewUrl = ref('');

const sealFileList = ref<UploadFileInfo[]>([]);

/** A4 预览的示例数据 */
const previewColumns = ['材料名称', '规格型号', '单位', '数量', '单价', '金额'];
const previewRows = [
  ['多层实木板', '18mm 环保E1', '张', '20', '128.00', '2,560.00'],
  ['不锈钢合页', '4寸 缓冲', '个', '24', '8.50', '204.00'],
  ['水性木器漆', '哑光白色 5L', '桶', '2', '368.00', '736.00']
];

// ===================== 打印预览 =====================

const previewVisible = ref(false);
const previewLoading = ref(false);
const previewHtml = ref('');
const previewIframeRef = ref<HTMLIFrameElement | null>(null);
const savingPdf = ref(false);

/** 预览用的示例单据列定义（与打印输出同构） */
const samplePrintColumns: PrintExportColumn<string[]>[] = [
  { title: '材料名称', value: cells => cells[0], width: 16 },
  { title: '规格型号', value: cells => cells[1], width: 14 },
  { title: '单位', value: cells => cells[2], width: 6, align: 'right' },
  { title: '数量', value: cells => cells[3], width: 8, align: 'right' },
  { title: '单价', value: cells => cells[4], width: 10, align: 'right' },
  { title: '金额', value: cells => cells[5], width: 12, align: 'right' }
];

/** 生成示例单据的打印 HTML 并在弹窗中预览 */
async function handlePrintPreview() {
  previewLoading.value = true;

  try {
    previewHtml.value = await renderPrintHtml({
      title: '材料清单',
      subtitle: '预览示例数据',
      columns: samplePrintColumns,
      rows: previewRows,
      summary: [{ label: '合计金额', value: '¥ 3,500.00' }],
      note: '以上为示例数据，实际打印内容以各单据页面的数据为准。'
    });
    previewVisible.value = true;
  } finally {
    previewLoading.value = false;
  }
}

/** 预览弹窗内直接调起打印 */
function handlePreviewPrint() {
  previewIframeRef.value?.contentWindow?.focus();
  previewIframeRef.value?.contentWindow?.print();
}

/** 预览文档另存为 PDF 文件 */
async function handlePreviewSavePdf() {
  if (!previewHtml.value) {
    return;
  }

  savingPdf.value = true;

  try {
    await saveHtmlAsPdf(previewHtml.value, model.title || '材料清单');
    window.$message?.success('PDF 已保存');
  } catch {
    window.$message?.error('PDF 生成失败，请稍后重试');
  } finally {
    savingPdf.value = false;
  }
}

const headTitle = computed(() => model.title || '材料清单（示例单据）');
const subLine = computed(() => [model.subTitle, '筛选：全部项目 · 全部仓库'].filter(Boolean).join('　|　'));
const dateText = computed(() => (model.showDate ? '打印日期：2026年11月05日' : ''));

async function loadSetting() {
  loading.value = true;

  try {
    const { data, error } = await fetchGetPrintSetting();
    if (!error && data?.setting) {
      model.title = data.setting.title;
      model.subTitle = data.setting.subTitle;
      model.footer = data.setting.footer;
      model.showDate = data.setting.showDate;
      model.showSeal = data.setting.showSeal;
      model.sealImage = data.setting.sealImage;
      sealPreviewUrl.value = data.setting.sealImageUrl ?? '';
    }
  } finally {
    loading.value = false;
  }
}

/** 印章上传：单张图片，上传成功后记录相对路径并预览 */
async function handleSealUpload({ file, onFinish, onError }: UploadCustomRequestOptions) {
  if (!file.file) {
    onError();

    return;
  }

  const { data, error } = await fetchUploadAdminImage(file.file);

  if (error || !data?.files?.length) {
    window.$message?.error('印章图片上传失败，请重试');
    onError();

    return;
  }

  const uploaded = data.files[0];
  model.sealImage = uploaded.path;
  sealPreviewUrl.value = uploaded.url;

  sealFileList.value = [
    { id: 'seal', name: uploaded.name, status: 'finished', url: uploaded.url } as UploadFileInfo
  ];
  onFinish();
}

function handleSealRemove() {
  model.sealImage = '';
  sealPreviewUrl.value = '';
  sealFileList.value = [];
}

async function handleSave() {
  saving.value = true;

  try {
    const { error } = await fetchSavePrintSetting({ ...model });

    if (error) {
      window.$message?.error('保存失败，请稍后重试');
      return;
    }

    invalidatePrintSetting();
    window.$message?.success('保存成功，打印单据时将使用此设置');
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  void loadSetting();
});
</script>

<template>
  <div class="min-h-500px">
    <NGrid :cols="2" :x-gap="16" responsive="screen" item-responsive>
      <!-- 左侧：设置表单 -->
      <NGi span="2 m:1">
        <NCard title="打印设置" :bordered="false" class="card-wrapper">
          <template #header-extra>
            <span class="text-12px text-gray-400">打印 PDF 单据时使用的公共样式</span>
          </template>

          <NSpin :show="loading">
            <NForm label-placement="left" label-width="96" require-mark-placement="right-hanging">
              <NFormItem label="抬头标题">
                <NInput v-model:value="model.title" maxlength="120" placeholder="如：XX 装饰工程有限公司，留空则只显示单据名" />
              </NFormItem>

              <NFormItem label="副标题">
                <NInput v-model:value="model.subTitle" maxlength="120" placeholder="如：YIMINJIA DESIGN · 设计工程部（选填）" />
              </NFormItem>

              <NFormItem label="页脚文字">
                <NInput v-model:value="model.footer" maxlength="255" placeholder="如：本单据由项目管理系统生成" />
              </NFormItem>

              <NFormItem label="打印日期">
                <div class="w-full">
                  <NSwitch v-model:value="model.showDate" />
                  <div class="mt-4px text-12px text-gray-400">开启后单据右上角显示当日日期</div>
                </div>
              </NFormItem>

              <NFormItem label="电子章">
                <div class="w-full">
                  <NSwitch v-model:value="model.showSeal" />
                  <div class="mt-4px text-12px text-gray-400">开启后单据右下角加盖电子章</div>
                </div>
              </NFormItem>

              <NFormItem v-if="model.showSeal" label="印章图片">
                <div class="w-full">
                  <NUpload
                    v-model:file-list="sealFileList"
                    :max="1"
                    accept="image/png,image/jpeg"
                    list-type="image-card"
                    :default-upload="true"
                    :custom-request="handleSealUpload"
                    @remove="handleSealRemove"
                  >
                    <div v-if="!sealFileList.length" class="flex-col items-center justify-center p-8px text-12px text-gray-400">
                      <SvgIcon icon="mdi:stamp" class="text-20px" />
                      <span>上传印章</span>
                    </div>
                  </NUpload>
                  <div class="text-12px text-gray-400">建议使用透明底 PNG 图片，打印效果更自然</div>
                </div>
              </NFormItem>
            </NForm>

            <div class="mt-8px">
              <NSpace :size="12">
                <NButton type="primary" :loading="saving" @click="handleSave">保存设置</NButton>
                <NButton :loading="previewLoading" @click="handlePrintPreview">
                  <template #icon><SvgIcon icon="mdi:file-eye-outline" /></template>
                  打印预览
                </NButton>
              </NSpace>
            </div>
          </NSpin>
        </NCard>
      </NGi>

      <!-- 右侧：A4 效果预览 -->
      <NGi span="2 m:1">
        <NCard title="打印效果预览" :bordered="false" class="card-wrapper">
          <div class="flex justify-center">
            <div class="a4-preview">
              <div class="print-header">
                <div class="title">{{ headTitle }}</div>
                <div v-if="subLine" class="sub">{{ subLine }}</div>
                <div v-if="dateText" class="date">{{ dateText }}</div>
              </div>

              <table>
                <thead>
                  <tr>
                    <th v-for="col in previewColumns" :key="col">{{ col }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(row, index) in previewRows" :key="index">
                    <td v-for="(cell, cellIndex) in row" :key="cellIndex" :class="cellIndex >= 2 ? 'col-right' : ''">{{ cell }}</td>
                  </tr>
                </tbody>
              </table>

              <div class="summary">
                <div class="summary-row"><span class="summary-label">合计金额</span><span>¥ 3,500.00</span></div>
              </div>

              <div v-if="model.footer" class="print-footer">{{ model.footer }}</div>

              <img v-if="model.showSeal && sealPreviewUrl" class="print-seal" :src="sealPreviewUrl" alt="seal" />
            </div>
          </div>
        </NCard>
      </NGi>
    </NGrid>

    <!-- 打印预览弹窗：展示与实际打印完全一致的 A4 文档 -->
    <NModal v-model:show="previewVisible" preset="card" title="打印预览（示例单据）" class="w-880px max-w-92vw">
      <template #header-extra>
        <span class="text-12px text-gray-400">确认无误后可直接打印，打印时选择「另存为 PDF」即可导出 PDF 文件</span>
      </template>

      <NSpin :show="previewLoading">
        <div class="flex justify-center">
          <iframe
            ref="previewIframeRef"
            class="h-68vh w-full border-0 bg-white"
            :srcdoc="previewHtml"
            title="打印预览"
          />
        </div>
      </NSpin>

      <template #footer>
        <NSpace justify="end" :size="12">
          <NButton @click="previewVisible = false">关闭</NButton>
          <NButton :loading="savingPdf" @click="handlePreviewSavePdf">
            <template #icon><SvgIcon icon="mdi:file-pdf-box" /></template>
            另存为 PDF
          </NButton>
          <NButton type="primary" @click="handlePreviewPrint">
            <template #icon><SvgIcon icon="mdi:printer-outline" /></template>
            直接打印
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.a4-preview {
  width: 100%;
  max-width: 480px;
  padding: 28px 24px 40px;
  background: #fff;
  box-shadow: 0 2px 12px rgb(0 0 0 / 8%);
  border-radius: 4px;
  font-size: 11px;
  color: #1f2937;
  position: relative;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
}

.print-header {
  text-align: center;
  margin-bottom: 12px;
}

.print-header .title {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 2px;
}

.print-header .sub {
  margin-top: 4px;
  font-size: 10px;
  color: #6b7280;
}

.print-header .date {
  text-align: right;
  font-size: 10px;
  color: #6b7280;
  margin-top: -14px;
}

.a4-preview table {
  width: 100%;
  border-collapse: collapse;
}

.a4-preview th,
.a4-preview td {
  border: 1px solid #94a3b8;
  padding: 4px 6px;
}

.a4-preview th {
  background: #f1f5f9;
  font-weight: 600;
}

.a4-preview td.col-right {
  text-align: right;
}

.summary {
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 3px;
  align-items: flex-end;
}

.summary-label {
  color: #6b7280;
  margin-right: 10px;
}

.print-footer {
  margin-top: 18px;
  text-align: center;
  font-size: 10px;
  color: #9ca3af;
}

.print-seal {
  position: absolute;
  right: 40px;
  bottom: 36px;
  width: 84px;
  height: 84px;
  object-fit: contain;
  opacity: 0.75;
  transform: rotate(-12deg);
  pointer-events: none;
}
</style>
