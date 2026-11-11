<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { NAlert, NButton, NCheckbox, NInput, NInputNumber, NSelect, NSteps, NStep, NTag } from 'naive-ui';
import { fetchSaveFinanceInit } from '@/service/api';
import { useAuthStore } from '@/store/modules/auth';

defineOptions({ name: 'FinanceSetup' });

type Category = 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'EXPENSE';
interface Account {
  code: string;
  name: string;
  categoryCode: Category;
  enabled: boolean;
  system: boolean;
}
interface FundAccount {
  type: 'CASH' | 'BANK';
  name: string;
  bankName: string;
  accountLast4: string;
  ledgerAccountCode: string;
}
interface Draft {
  version: number;
  step: number;
  year: number;
  startMonth: number;
  accounts: Account[];
  fundAccounts: FundAccount[];
  openingAmounts: Record<string, number | null>;
}

const authStore = useAuthStore();
const route = useRoute();
const router = useRouter();
const currentYear = new Date().getFullYear();
const categoryOptions = [
  { label: '资产', value: 'ASSET' },
  { label: '负债', value: 'LIABILITY' },
  { label: '权益', value: 'EQUITY' },
  { label: '收入', value: 'REVENUE' },
  { label: '费用', value: 'EXPENSE' }
];
const defaultAccounts: Account[] = [
  { code: '1001', name: '库存现金', categoryCode: 'ASSET', enabled: true, system: true },
  { code: '1002', name: '银行存款', categoryCode: 'ASSET', enabled: true, system: true },
  { code: '1122', name: '应收账款', categoryCode: 'ASSET', enabled: true, system: true },
  { code: '1221', name: '其他应收款', categoryCode: 'ASSET', enabled: false, system: true },
  { code: '1405', name: '库存商品', categoryCode: 'ASSET', enabled: false, system: true },
  { code: '1601', name: '固定资产', categoryCode: 'ASSET', enabled: false, system: true },
  { code: '2202', name: '应付账款', categoryCode: 'LIABILITY', enabled: true, system: true },
  { code: '2211', name: '应付职工薪酬', categoryCode: 'LIABILITY', enabled: false, system: true },
  { code: '2221', name: '应交税费', categoryCode: 'LIABILITY', enabled: false, system: true },
  { code: '2241', name: '其他应付款', categoryCode: 'LIABILITY', enabled: false, system: true },
  { code: '4001', name: '实收资本', categoryCode: 'EQUITY', enabled: true, system: true },
  { code: '4103', name: '本年利润', categoryCode: 'EQUITY', enabled: true, system: true },
  { code: '5001', name: '主营业务收入', categoryCode: 'REVENUE', enabled: true, system: true },
  { code: '5051', name: '其他业务收入', categoryCode: 'REVENUE', enabled: false, system: true },
  { code: '5401', name: '主营业务成本', categoryCode: 'EXPENSE', enabled: true, system: true },
  { code: '5601', name: '销售费用', categoryCode: 'EXPENSE', enabled: true, system: true },
  { code: '5602', name: '管理费用', categoryCode: 'EXPENSE', enabled: true, system: true },
  { code: '5603', name: '财务费用', categoryCode: 'EXPENSE', enabled: false, system: true }
];
const draftKey = `finance-init-draft-${authStore.userInfo.userId}`;
const step = ref(1);
const saving = ref(false);
const accounts = ref<Account[]>(structuredClone(defaultAccounts));
const fundAccounts = ref<FundAccount[]>([]);
const openingAmounts = ref<Record<string, number | null>>({});
const draft = reactive<Draft>({ version: 4, step: 1, year: currentYear, startMonth: 1, accounts: [], fundAccounts: [], openingAmounts: {} });
const isReconfigure = computed(() => route.query.reconfigure === '1');
const enabledAccounts = computed(() => accounts.value.filter(account => account.enabled));
const enabledAssetAccounts = computed(() => enabledAccounts.value.filter(account => account.categoryCode === 'ASSET'));
const openingTotals = computed(() => {
  return enabledAccounts.value.reduce(
    (totals, account) => {
      const amountCents = Math.round((openingAmounts.value[account.code] ?? 0) * 100);
      if (directionFor(account.categoryCode) === 'DEBIT') totals.debit += amountCents;
      else totals.credit += amountCents;
      return totals;
    },
    { debit: 0, credit: 0 }
  );
});
const openingBalanced = computed(() => openingTotals.value.debit === openingTotals.value.credit);
const openingAmountText = (value: number) => (value / 100).toFixed(2);
const directionFor = (category: Category) =>
  category === 'ASSET' || category === 'EXPENSE' ? 'DEBIT' : 'CREDIT';
const directionLabel = (category: Category) => (directionFor(category) === 'DEBIT' ? '借方' : '贷方');
const categoryLabel = (category: Category) => categoryOptions.find(item => item.value === category)?.label ?? category;
const monthOptions = Array.from({ length: 12 }, (_, index) => ({ label: `${index + 1} 月`, value: index + 1 }));

function saveDraft() {
  draft.version = 4;
  draft.step = step.value;
  draft.accounts = accounts.value.map(account => ({ ...account }));
  draft.fundAccounts = fundAccounts.value.map(account => ({ ...account }));
  draft.openingAmounts = { ...openingAmounts.value };
  sessionStorage.setItem(draftKey, JSON.stringify(draft));
}

watch([step, accounts, fundAccounts, openingAmounts, () => draft.year, () => draft.startMonth], saveDraft, { deep: true });

onMounted(() => {
  const saved = sessionStorage.getItem(draftKey);
  if (!saved) return;
  try {
    const parsed = JSON.parse(saved) as Draft;
    if (parsed.year && parsed.startMonth && Array.isArray(parsed.accounts)) {
      const restoredStep = parsed.version === 4
        ? parsed.step
        : parsed.version === 3 || parsed.version === 2
          ? parsed.step >= 4 ? parsed.step + 1 : parsed.step
          : parsed.step >= 3 ? parsed.step + 2 : parsed.step;
      step.value = Math.min(6, Math.max(1, restoredStep || 1));
      draft.version = 4;
      draft.year = parsed.year;
      draft.startMonth = parsed.startMonth;
      accounts.value = parsed.accounts;
      fundAccounts.value = parsed.fundAccounts ?? [];
      openingAmounts.value = parsed.openingAmounts ?? {};
    }
  } catch {
    sessionStorage.removeItem(draftKey);
  }
});

function addAccount() {
  accounts.value.push({ code: '', name: '', categoryCode: 'ASSET', enabled: true, system: false });
}

function addFundAccount(type: FundAccount['type']) {
  const defaultCode = type === 'CASH' ? '1001' : '1002';
  const ledgerAccountCode = enabledAssetAccounts.value.some(account => account.code === defaultCode)
    ? defaultCode
    : enabledAssetAccounts.value[0]?.code ?? '';
  fundAccounts.value.push({
    type,
    name: type === 'CASH' ? '库存现金' : '',
    bankName: '',
    accountLast4: '',
    ledgerAccountCode
  });
}

function removeAccount(account: Account) {
  accounts.value = accounts.value.filter(item => item !== account);
}

function nextStep() {
  if (step.value === 1 && (!draft.year || draft.startMonth < 1 || draft.startMonth > 12)) {
    window.$message?.error('请检查会计年度和启用月份');
    return;
  }
  if (step.value === 2 && !enabledAccounts.value.length) {
    window.$message?.error('请至少启用一个科目');
    return;
  }
  if (
    step.value === 3 &&
    fundAccounts.value.some(
      account =>
        !account.name.trim() ||
        !enabledAssetAccounts.value.some(asset => asset.code === account.ledgerAccountCode) ||
        (account.type === 'BANK' && !account.bankName.trim()) ||
        (account.accountLast4 && !/^\d{4}$/.test(account.accountLast4))
    )
  ) {
    window.$message?.error('请补全资金账户名称、开户银行和资产科目；银行卡末四位需为4位数字');
    return;
  }
  if (step.value === 4 && !openingBalanced.value) {
    window.$message?.error(`期初借贷不平衡，借方比贷方${openingTotals.value.debit > openingTotals.value.credit ? '多' : '少'} ${openingAmountText(Math.abs(openingTotals.value.debit - openingTotals.value.credit))} 元`);
    return;
  }
  step.value = Math.min(6, step.value + 1);
}

async function submit() {
  if (saving.value) return;
  saving.value = true;
  const { error } = await fetchSaveFinanceInit({
    year: draft.year,
    startMonth: draft.startMonth,
    accounts: enabledAccounts.value.map(({ code, name, categoryCode }) => ({ code, name, categoryCode, enabled: true })),
    fundAccounts: fundAccounts.value.map(account => ({ ...account })),
    openingBalances: enabledAccounts.value.map(account => ({
      accountCode: account.code,
      amount: (openingAmounts.value[account.code] ?? 0).toFixed(2)
    })),
    reconfigure: isReconfigure.value
  });
  saving.value = false;
  if (error) return;
  sessionStorage.removeItem(draftKey);
  authStore.userInfo.finInitialized = true;
  window.$message?.success('财务基础配置已保存');
  await router.push({ name: 'home' });
}
</script>

<template>
  <main class="setup-page">
    <header class="setup-header">
      <div>
        <p class="eyebrow">FINANCE FOUNDATION</p>
        <h1>{{ isReconfigure ? '财务基础配置' : '开始使用财务管理' }}</h1>
        <p class="intro">先确定记账期间和常用科目，系统会自动安排借贷方向。现在只保存基础设置，不会生成资金流水或凭证。</p>
      </div>
      <div class="header-mark"><SvgIcon icon="mdi:scale-balance" /></div>
    </header>

    <div class="stepper-wrap">
      <NSteps :current="step" :status="'process'" size="small" class="stepper">
        <NStep title="会计期间" description="确定从何时开始记录" />
        <NStep title="科目清单" description="选择需要使用的科目" />
        <NStep title="资金账户" description="登记现金和银行卡" />
        <NStep title="期初余额" description="录入启用日账面余额" />
        <NStep title="借贷方向" description="了解系统的自动规则" />
        <NStep title="核对提交" description="确认后一次保存" />
      </NSteps>
    </div>

    <section class="step-content">
      <div v-if="step === 1" class="step-panel">
        <div class="section-heading">
          <span class="step-index">01</span>
          <div><h2>选择会计期间</h2><p>会计期间是财务记录归属的月份。默认从今年开始，历史月份不会自动补录。</p></div>
        </div>
        <div class="period-fields">
          <label><span>会计年度</span><NInputNumber v-model:value="draft.year" :min="2000" :max="2100" :precision="0" /></label>
          <label><span>启用月份</span><NSelect v-model:value="draft.startMonth" :options="monthOptions" /></label>
        </div>
        <NAlert type="info" :bordered="false">从 {{ draft.year }} 年 {{ draft.startMonth }} 月开始启用；当年后续月份会预先创建，之后可按月结账。</NAlert>
      </div>

      <div v-else-if="step === 2" class="step-panel">
        <div class="section-heading">
          <span class="step-index">02</span>
          <div><h2>选择常用科目</h2><p>科目就像分类标签，用来说明每笔账属于现金、应收款、收入或费用。勾选本公司实际需要的项目，也可调整编码或新增。</p></div>
        </div>
        <div class="account-toolbar"><span>已启用 {{ enabledAccounts.length }} 项</span><NButton size="small" secondary @click="addAccount">新增科目</NButton></div>
        <div class="account-list">
          <div v-for="(account, index) in accounts" :key="`${account.code}-${index}`" class="account-row">
            <NCheckbox v-model:checked="account.enabled" />
            <NInput v-model:value="account.code" :disabled="!account.enabled" placeholder="编码" class="code-input" />
            <NInput v-model:value="account.name" :disabled="!account.enabled" placeholder="科目名称" />
            <NSelect v-model:value="account.categoryCode" :options="categoryOptions" :disabled="!account.enabled" class="category-select" />
            <NButton v-if="!account.system" quaternary type="error" @click="removeAccount(account)">移除</NButton>
          </div>
        </div>
      </div>

      <div v-else-if="step === 3" class="step-panel">
        <div class="section-heading">
          <span class="step-index">03</span>
          <div><h2>配置现金与银行账户</h2><p>在这里登记日常收付款使用的具体账户，并关联对应的资产科目。银行卡仅记录开户行和可选的末四位，不需要输入完整卡号或密码。</p></div>
        </div>
        <div class="fund-accounts">
          <div class="fund-heading">
            <div><h3>资金账户清单</h3><p>可添加多个现金账户或银行卡，也可以暂时跳过，之后再补充。</p></div>
            <div class="fund-actions"><NButton size="small" secondary @click="addFundAccount('CASH')">添加现金账户</NButton><NButton size="small" secondary @click="addFundAccount('BANK')">添加银行卡</NButton></div>
          </div>
          <div v-if="!fundAccounts.length" class="fund-empty">尚未添加资金账户。</div>
          <div class="fund-rows">
            <div v-for="(account, index) in fundAccounts" :key="index" class="fund-row" :class="{ 'fund-row--cash': account.type === 'CASH' }">
              <NTag :bordered="false" :type="account.type === 'BANK' ? 'info' : 'success'">{{ account.type === 'BANK' ? '银行卡' : '现金' }}</NTag>
              <NInput v-model:value="account.name" :placeholder="account.type === 'BANK' ? '账户名称，如：公司基本户' : '账户名称，如：前台备用金'" />
              <NInput v-if="account.type === 'BANK'" v-model:value="account.bankName" placeholder="开户银行" />
              <NInput v-if="account.type === 'BANK'" v-model:value="account.accountLast4" maxlength="4" placeholder="卡号末四位（可选）" />
              <NSelect v-model:value="account.ledgerAccountCode" :options="enabledAssetAccounts.map(item => ({ label: `${item.code} ${item.name}`, value: item.code }))" placeholder="关联资产科目" />
              <NButton quaternary type="error" @click="fundAccounts.splice(index, 1)">移除</NButton>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="step === 4" class="step-panel">
        <div class="section-heading">
          <span class="step-index">04</span>
          <div><h2>录入启用日的期初余额</h2><p>填写正式启用记账前，各科目账面上已有的余额。没有余额的科目填 0；现金、银行、应收应付、固定资产等请按同一日期的账面数据录入。</p></div>
        </div>
        <NAlert type="info" :bordered="false" class="opening-help">金额自动按科目类别计入正常方向：资产/费用计借方，负债/权益/收入计贷方。请核对借贷合计相等；不平衡时不能继续。</NAlert>
        <div class="opening-table">
          <div class="opening-table__head"><span>科目</span><span>类别</span><span>余额方向</span><span>期初余额（元）</span></div>
          <div v-for="account in enabledAccounts" :key="account.code" class="opening-row">
            <div><strong>{{ account.name }}</strong><small>{{ account.code }}</small></div>
            <span>{{ categoryLabel(account.categoryCode) }}</span>
            <span>{{ directionLabel(account.categoryCode) }}</span>
            <NInputNumber v-model:value="openingAmounts[account.code]" :min="0" :precision="2" :max="9999999999999999" placeholder="0.00" />
          </div>
        </div>
        <div class="opening-totals" :class="{ 'opening-totals--unbalanced': !openingBalanced }">
          <div><span>借方合计</span><strong>¥ {{ openingAmountText(openingTotals.debit) }}</strong></div>
          <span class="balance-equality">{{ openingBalanced ? '=' : '≠' }}</span>
          <div><span>贷方合计</span><strong>¥ {{ openingAmountText(openingTotals.credit) }}</strong></div>
          <NTag :bordered="false" :type="openingBalanced ? 'success' : 'error'">{{ openingBalanced ? '平衡' : `差额 ¥ ${openingAmountText(Math.abs(openingTotals.debit - openingTotals.credit))}` }}</NTag>
        </div>
      </div>

      <div v-else-if="step === 5" class="step-panel">
        <div class="section-heading">
          <span class="step-index">05</span>
          <div><h2>确认借贷方向规则</h2><p>借和贷是账目左右两边的记账标记，不代表收入或支出。方向由科目类别自动决定，您无需手动设置。</p></div>
        </div>
        <div class="direction-grid">
          <div class="direction-card debit"><span class="direction-symbol">借</span><div><h3>借方</h3><p>资产增加，费用发生</p></div><NTag :bordered="false" type="info">自动</NTag></div>
          <div class="direction-card credit"><span class="direction-symbol">贷</span><div><h3>贷方</h3><p>负债、权益或收入增加</p></div><NTag :bordered="false" type="success">自动</NTag></div>
        </div>
        <div class="direction-table">
          <div class="direction-table__head"><span>科目</span><span>类别</span><span>记账方向</span></div>
          <div v-for="account in enabledAccounts" :key="account.code" class="direction-table__row">
            <strong>{{ account.name }}</strong><span>{{ categoryLabel(account.categoryCode) }}</span><span>{{ directionLabel(account.categoryCode) }}</span>
          </div>
        </div>
      </div>

      <div v-else class="step-panel">
        <div class="section-heading">
          <span class="step-index">06</span>
          <div><h2>核对后完成配置</h2><p>提交后配置会保存到系统，之后登录将直接进入首页。此操作不会创建任何凭证、分录或资金流水。</p></div>
        </div>
        <div class="review-summary">
          <div class="review-item"><span>会计年度</span><strong>{{ draft.year }} 年</strong></div>
          <div class="review-item"><span>启用月份</span><strong>{{ draft.startMonth }} 月</strong></div>
          <div class="review-item"><span>启用科目</span><strong>{{ enabledAccounts.length }} 项</strong></div>
        </div>
        <div class="review-accounts">
          <div v-for="account in enabledAccounts" :key="account.code" class="review-account">
            <span class="review-account__code">{{ account.code }}</span><strong>{{ account.name }}</strong><span>{{ categoryLabel(account.categoryCode) }} · {{ directionLabel(account.categoryCode) }} ¥{{ openingAmountText(Math.round((openingAmounts[account.code] ?? 0) * 100)) }}</span><NTag size="small" :bordered="false" :type="directionFor(account.categoryCode) === 'DEBIT' ? 'info' : 'success'">{{ directionLabel(account.categoryCode) }}</NTag>
          </div>
        </div>
        <div class="review-balance-total"><span>期初余额借贷校验</span><strong>借 ¥{{ openingAmountText(openingTotals.debit) }} = 贷 ¥{{ openingAmountText(openingTotals.credit) }}</strong></div>
        <div v-if="fundAccounts.length" class="review-funds">
          <h3>现金与银行账户</h3>
          <div v-for="(account, index) in fundAccounts" :key="index" class="review-fund">
            <NTag size="small" :bordered="false" :type="account.type === 'BANK' ? 'info' : 'success'">{{ account.type === 'BANK' ? '银行卡' : '现金' }}</NTag>
            <strong>{{ account.name }}</strong>
            <span>{{ account.type === 'BANK' ? `${account.bankName}${account.accountLast4 ? ` · 末四位 ${account.accountLast4}` : ''}` : '现金账户' }}</span>
            <span>{{ account.ledgerAccountCode }}</span>
          </div>
        </div>
        <NAlert type="warning" :bordered="false">本次只建立会计期间和科目设置，不会生成任何凭证。配置完成后，登记收款 / 付款时会按借贷规则自动生成凭证，并同步资金流水与科目余额。</NAlert>
      </div>
    </section>

    <footer class="wizard-footer">
      <NButton quaternary @click="step > 1 ? step-- : router.push({ name: 'home' })">{{ step > 1 ? '上一步' : '暂不配置' }}</NButton>
      <div class="footer-right"><span>第 {{ step }} 步，共 6 步</span><NButton v-if="step < 6" type="primary" @click="nextStep">继续</NButton><NButton v-else type="primary" :loading="saving" @click="submit">确认并保存</NButton></div>
    </footer>
  </main>
</template>

<style scoped>
.setup-page { --ink: #20302e; --muted: #64736f; --line: #dce5e1; --teal: #157b70; --wash: #f1f7f4; max-width: 1080px; margin: 0 auto; padding: 34px clamp(16px, 4vw, 48px) 28px; color: var(--ink); }
.setup-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 24px; padding: 12px 0 28px; border-bottom: 1px solid var(--line); }
.eyebrow { margin: 0 0 10px; color: var(--teal); font-size: 11px; font-weight: 700; letter-spacing: 1.4px; }
h1 { margin: 0; font-size: 28px; font-weight: 650; }
.intro { max-width: 700px; margin: 10px 0 0; color: var(--muted); line-height: 1.7; }
.header-mark { display: grid; width: 52px; height: 52px; flex: 0 0 auto; place-items: center; border: 1px solid #c9ded8; border-radius: 8px; background: var(--wash); color: var(--teal); font-size: 25px; }
.stepper-wrap { padding: 26px 0 24px; border-bottom: 1px solid var(--line); }
.stepper { max-width: 850px; margin: auto; }
.step-content { min-height: 400px; padding: 28px 0; }
.step-panel { max-width: 850px; margin: auto; animation: enter 260ms ease both; }
.section-heading { display: flex; gap: 16px; margin-bottom: 28px; }
.step-index { color: #a7b9b2; font-size: 14px; font-weight: 700; font-variant-numeric: tabular-nums; }
h2 { margin: 0; font-size: 21px; font-weight: 620; }
.section-heading p { max-width: 700px; margin: 8px 0 0; color: var(--muted); line-height: 1.65; }
.period-fields { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 20px; max-width: 600px; margin: 32px 0 24px 32px; }
.period-fields label { display: grid; gap: 9px; color: var(--muted); font-size: 13px; }
.account-toolbar { display: flex; align-items: center; justify-content: space-between; margin: 0 0 10px 32px; color: var(--muted); font-size: 13px; }
.account-list { overflow: hidden; border: 1px solid var(--line); border-radius: 6px; }
.account-row { display: grid; grid-template-columns: 24px 105px minmax(150px, 1fr) 145px auto; align-items: center; gap: 10px; min-height: 54px; padding: 7px 12px; border-bottom: 1px solid #edf1ef; }
.account-row:last-child { border-bottom: 0; }
.fund-accounts { margin-top: 26px; border-top: 1px solid var(--line); padding-top: 20px; }
.fund-heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; }
.fund-heading h3 { margin: 0; font-size: 15px; font-weight: 620; }
.fund-heading p { max-width: 650px; margin: 6px 0 0; color: var(--muted); font-size: 12px; line-height: 1.6; }
.fund-actions { display: flex; flex: 0 0 auto; gap: 8px; }
.fund-empty { padding: 18px 12px; color: var(--muted); font-size: 13px; }
.fund-rows { min-width: 0; }
.fund-row { display: grid; grid-template-columns: 72px minmax(120px, 1.1fr) minmax(110px, 1fr) minmax(130px, 1.1fr) minmax(135px, 1.1fr) max-content; align-items: center; gap: 9px; padding: 10px 0; border-bottom: 1px solid #edf1ef; }
.fund-row--cash { grid-template-columns: 72px minmax(150px, 1fr) minmax(150px, 1fr) max-content; }
.fund-row :deep(.n-button) { white-space: nowrap; }
.opening-help { margin-bottom: 16px; }
.opening-table { overflow: hidden; border: 1px solid var(--line); border-radius: 6px; }
.opening-table__head,.opening-row { display: grid; grid-template-columns: minmax(170px, 1.4fr) minmax(90px, .7fr) minmax(90px, .7fr) minmax(160px, 1fr); align-items: center; gap: 14px; padding: 11px 14px; border-bottom: 1px solid #edf1ef; }
.opening-table__head { color: var(--muted); font-size: 12px; background: var(--wash); }
.opening-row { min-height: 60px; font-size: 13px; }.opening-row>div:first-child { display: grid; gap: 4px; }.opening-row small { color: var(--muted); font-size: 11px; }.opening-row :deep(.n-input-number) { width: 100%; }
.opening-totals { display: flex; align-items: center; justify-content: flex-end; gap: 22px; margin-top: 16px; padding: 15px 18px; border: 1px solid #b9d8c8; border-radius: 6px; background: #f0f8f3; color: #286744; }
.opening-totals--unbalanced { border-color: #e5b9b5; background: #fff5f4; color: #a3342b; }
.opening-totals>div { display: grid; gap: 4px; }.opening-totals>div span { font-size: 12px; }.opening-totals>div strong { font-size: 16px; font-variant-numeric: tabular-nums; }.balance-equality { font-size: 20px; font-weight: 700; }
.review-balance-total { display: flex; justify-content: space-between; margin: -8px 0 20px 32px; padding: 12px; border: 1px solid #b9d8c8; border-radius: 6px; background: #f0f8f3; color: #286744; font-size: 13px; }
.code-input { font-variant-numeric: tabular-nums; }
.direction-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 14px; margin-bottom: 25px; }
.direction-card { display: flex; align-items: center; gap: 14px; min-height: 94px; padding: 16px; border: 1px solid var(--line); border-radius: 6px; background: #fff; }
.direction-card h3 { margin: 0; font-size: 15px; }.direction-card p { margin: 5px 0 0; color: var(--muted); font-size: 13px; }.direction-card :deep(.n-tag) { margin-left: auto; }
.direction-symbol { display: grid; width: 42px; height: 42px; place-items: center; border-radius: 6px; background: #e9f2f7; color: #32657c; font-size: 20px; font-weight: 700; }.credit .direction-symbol { background: #e7f4ed; color: #317454; }
.direction-table { border-top: 1px solid var(--line); }.direction-table__head,.direction-table__row { display: grid; grid-template-columns: 1fr 1fr 120px; gap: 12px; padding: 12px 14px; border-bottom: 1px solid #edf1ef; }.direction-table__head { color: var(--muted); font-size: 12px; }.direction-table__row { font-size: 13px; }.direction-table__row span:last-child { font-weight: 600; }
.review-summary { display: grid; grid-template-columns: repeat(3, 1fr); margin: 0 0 20px 32px; border: 1px solid var(--line); border-radius: 6px; }.review-item { display: grid; gap: 7px; padding: 15px 18px; border-right: 1px solid var(--line); }.review-item:last-child { border: 0; }.review-item span,.review-account__code,.review-account>span:nth-last-child(2) { color: var(--muted); font-size: 12px; }.review-item strong { font-size: 16px; }
.review-accounts { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); margin: 0 0 22px 32px; border-top: 1px solid var(--line); }.review-account { display: grid; grid-template-columns: 52px minmax(0, 1fr) 45px auto; align-items: center; gap: 8px; min-height: 44px; padding: 6px 10px; border-bottom: 1px solid #edf1ef; font-size: 13px; }.review-account__code { font-variant-numeric: tabular-nums; }
.review-funds { margin: 0 0 22px 32px; border-top: 1px solid var(--line); }.review-funds h3 { margin: 14px 0 8px; font-size: 14px; }.review-fund { display: grid; grid-template-columns: 76px minmax(130px, 1fr) minmax(150px, 1fr) 80px; align-items: center; gap: 10px; min-height: 42px; border-bottom: 1px solid #edf1ef; font-size: 13px; }.review-fund span:not(.n-tag) { color: var(--muted); font-size: 12px; }
.wizard-footer { display: flex; align-items: center; justify-content: space-between; padding-top: 16px; border-top: 1px solid var(--line); }.footer-right { display: flex; align-items: center; gap: 18px; }.footer-right>span { color: var(--muted); font-size: 12px; }
@keyframes enter { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
@media (max-width: 700px) { .setup-page { padding-top: 20px; }.setup-header { padding-bottom: 20px; }h1 { font-size: 23px; }.header-mark { width: 42px; height: 42px; }.stepper-wrap { overflow-x: auto; }.stepper { min-width: 760px; }.period-fields { margin-left: 0; }.account-toolbar { margin-left: 0; }.account-row { grid-template-columns: 24px 75px minmax(110px, 1fr) 95px; gap: 6px; padding: 7px; }.account-row :deep(.n-button) { grid-column: 2 / -1; justify-self: end; }.fund-heading { flex-direction: column; }.fund-actions { flex-wrap: wrap; }.fund-row { grid-template-columns: 72px repeat(2, minmax(0, 1fr)); gap: 8px; }.fund-row--cash { grid-template-columns: 72px minmax(0, 1fr) max-content; }.fund-row--cash :deep(.n-button) { grid-column: 3; justify-self: end; }.opening-table__head { display: none; }.opening-row { grid-template-columns: minmax(0, 1fr) minmax(110px, .8fr); gap: 8px; padding: 10px; }.opening-row>span { font-size: 12px; }.opening-totals { justify-content: space-between; gap: 8px; padding: 12px; }.opening-totals :deep(.n-tag) { display: none; }.review-summary { margin-left: 0; }.review-accounts,.review-funds { grid-template-columns: 1fr; margin-left: 0; }.review-account { grid-template-columns: 48px minmax(0, 1fr) 38px auto; }.review-balance-total { gap: 10px; flex-direction: column; margin-left: 0; }.review-fund { grid-template-columns: 76px minmax(0, 1fr); padding: 8px 0; }.wizard-footer { gap: 10px; }.footer-right { gap: 9px; } }
</style>