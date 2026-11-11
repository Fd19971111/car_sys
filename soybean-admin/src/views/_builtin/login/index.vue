<script setup lang="ts">
import { computed } from 'vue';
import type { Component } from 'vue';
import { mixColor } from '@sa/color';
import { loginModuleRecord } from '@/constants/app';
import { useAppStore } from '@/store/modules/app';
import { useThemeStore } from '@/store/modules/theme';
import { $t } from '@/locales';
import loginBg from '@/assets/imgs/login-bg.jpg';
import PwdLogin from './modules/pwd-login.vue';
import CodeLogin from './modules/code-login.vue';
import Register from './modules/register.vue';
import ResetPwd from './modules/reset-pwd.vue';
import BindWechat from './modules/bind-wechat.vue';

interface Props {
  /** The login module */
  module?: UnionKey.LoginModule;
}

const props = defineProps<Props>();

const appStore = useAppStore();
const themeStore = useThemeStore();

interface LoginModule {
  label: App.I18n.I18nKey;
  component: Component;
}

const moduleMap: Record<UnionKey.LoginModule, LoginModule> = {
  'pwd-login': { label: loginModuleRecord['pwd-login'], component: PwdLogin },
  'code-login': { label: loginModuleRecord['code-login'], component: CodeLogin },
  register: { label: loginModuleRecord.register, component: Register },
  'reset-pwd': { label: loginModuleRecord['reset-pwd'], component: ResetPwd },
  'bind-wechat': { label: loginModuleRecord['bind-wechat'], component: BindWechat }
};

const activeModule = computed(() => moduleMap[props.module || 'pwd-login']);

/** 右侧登录区底色：浅色时用极淡的品牌色，深色时用深棕黑，避免照片干扰表单 */
const panelBg = computed(() =>
  themeStore.darkMode ? 'rgb(23, 18, 15)' : mixColor('#ffffff', themeStore.themeColor, 0.08)
);

/** 左侧品牌区的能力标签 */
const highlights = ['项目全流程', '材料与库存', '收付款对账', '自动记账凭证'];
</script>

<template>
  <div class="login-page">
    <!-- 左侧品牌区：实景装修图 + 品牌主张，小屏隐藏 -->
    <aside class="login-brand">
      <img :src="loginBg" class="login-brand-bg" alt="亿民家装修实景" />
      <div class="login-brand-overlay"></div>

      <div class="login-brand-content">
        <div class="flex-y-center gap-14px">
          <SystemLogo class="size-52px" />
          <div>
            <div class="text-24px text-white font-600 tracking-wide">亿民家</div>
            <div class="text-11px text-white/60 tracking-[3px]">YIMINJIA</div>
          </div>
        </div>

        <div class="mt-auto">
          <h1 class="text-34px text-white font-600 leading-[1.4]">
            把每一个工地<br />
            管得明白，算得清楚
          </h1>
          <p class="mt-16px text-14px text-white/75 leading-7">
            项目、材料、收付款、凭证、报表，一套系统串起装修工程全流程。
          </p>
          <div class="mt-26px flex flex-wrap gap-8px">
            <span v-for="item in highlights" :key="item" class="login-brand-tag">{{ item }}</span>
          </div>
        </div>
      </div>

      <div class="login-brand-footer">亿民家 · 装修工程数字化管理平台</div>
    </aside>

    <!-- 右侧登录区 -->
    <main class="login-panel" :style="{ backgroundColor: panelBg }">
      <div class="login-panel-tools">
        <ThemeSchemaSwitch
          :theme-schema="themeStore.themeScheme"
          :show-tooltip="false"
          class="text-18px"
          @switch="themeStore.toggleThemeScheme"
        />
        <LangSwitch
          v-if="themeStore.header.multilingual.visible"
          :lang="appStore.locale"
          :lang-options="appStore.localeOptions"
          :show-tooltip="false"
          @change-lang="appStore.changeLocale"
        />
      </div>

      <div class="login-panel-card">
        <header class="flex-y-center gap-12px">
          <SystemLogo class="size-44px" />
          <div>
            <h2 class="text-20px text-primary font-600 leading-none">{{ $t('system.title') }}</h2>
            <p class="mt-8px text-12px text-gray-400">装修工程管理平台</p>
          </div>
        </header>

        <main class="pt-30px">
          <h3 class="text-16px font-500">{{ $t(activeModule.label) }}</h3>
          <div class="pt-22px">
            <Transition :name="themeStore.page.animateMode" mode="out-in" appear>
              <component :is="activeModule.component" />
            </Transition>
          </div>
        </main>

        <footer class="pt-36px text-center text-11px text-gray-400">© {{ new Date().getFullYear() }} 亿民家 · 保留所有权利</footer>
      </div>
    </main>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

/* ---------- 左侧品牌区 ---------- */
.login-brand {
  position: relative;
  width: 52%;
  min-width: 520px;
  overflow: hidden;
}

.login-brand-bg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: 68% center;
  transform: scale(1.18);
}

.login-brand-overlay {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 88% 92%, rgb(24 14 8 / 62%) 0%, rgb(24 14 8 / 0%) 34%),
    linear-gradient(90deg, rgb(24 14 8 / 62%) 0%, rgb(24 14 8 / 8%) 58%, rgb(24 14 8 / 18%) 100%),
    linear-gradient(180deg, rgb(24 14 8 / 58%) 0%, rgb(24 14 8 / 10%) 40%, rgb(24 14 8 / 88%) 100%);
}

.login-brand-content {
  position: relative;
  z-index: 2;
  display: flex;
  height: 100%;
  flex-direction: column;
  padding: 40px 48px 68px;
}

.login-brand-tag {
  border: 1px solid rgb(255 255 255 / 18%);
  border-radius: 999px;
  padding: 6px 14px;
  font-size: 12px;
  color: rgb(255 255 255 / 88%);
  background: rgb(255 255 255 / 12%);
  backdrop-filter: blur(6px);
}

.login-brand-footer {
  position: absolute;
  z-index: 2;
  bottom: 24px;
  left: 48px;
  font-size: 12px;
  color: rgb(255 255 255 / 45%);
}

/* ---------- 右侧登录区 ---------- */
.login-panel {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
}

.login-panel-tools {
  position: absolute;
  top: 20px;
  right: 24px;
  display: flex;
  align-items: center;
  gap: 14px;
}

.login-panel-card {
  width: 100%;
  max-width: 380px;
}

/* ---------- 小屏适配 ---------- */
@media (max-width: 960px) {
  .login-brand {
    display: none;
  }

  .login-panel {
    padding: 40px 22px;
  }
}
</style>
