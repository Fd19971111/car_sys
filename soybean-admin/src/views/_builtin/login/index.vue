<script setup lang="ts">
import { computed } from 'vue';
import type { Component } from 'vue';
import { mixColor } from '@sa/color';
import { loginModuleRecord } from '@/constants/app';
import { useAppStore } from '@/store/modules/app';
import { useThemeStore } from '@/store/modules/theme';
import { $t } from '@/locales';
import loginBg from '@/assets/imgs/ride-hailing-bg.jpg';
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

/** Keep the login panel subtly tied to the transport-blue brand. */
const panelBg = computed(() =>
  themeStore.darkMode ? 'rgb(10, 25, 43)' : mixColor('#ffffff', themeStore.themeColor, 0.08)
);

/** 左侧品牌区的能力标签 */
const highlights = ['实时订单', '司机与车辆', '路线调度', '营收对账'];
</script>

<template>
  <div class="login-page">
    <!-- 左侧品牌区：实景装修图 + 品牌主张，小屏隐藏 -->
    <aside class="login-brand">
      <img :src="loginBg" class="login-brand-bg" alt="黔车会 QCH 运营车辆" />
      <div class="login-brand-overlay"></div>

      <div class="login-brand-content">
        <div class="flex-y-center gap-14px">
          <SystemLogo class="size-52px" />
          <div>
            <div class="text-24px text-white font-600 tracking-wide">黔车会</div>
            <div class="text-11px text-white/60 tracking-[3px]">QIAN CHE HUI</div>
          </div>
        </div>

        <div class="mt-auto">
          <h1 class="text-34px text-white font-600 leading-[1.4]">
            让每一程<br />
            都准时、安全、可追踪
          </h1>
          <p class="mt-16px text-14px text-white/75 leading-7">
            车辆、司机、订单与结算，一套系统掌握每一次出行服务。
          </p>
          <div class="mt-26px flex flex-wrap gap-8px">
            <span v-for="item in highlights" :key="item" class="login-brand-tag">{{ item }}</span>
          </div>
        </div>
      </div>

      <div class="login-brand-footer">黔车会 · 打车运营管理平台</div>
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
            <p class="mt-8px text-12px text-gray-400">打车运营管理平台</p>
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

        <footer class="pt-36px text-center text-11px text-gray-400">© {{ new Date().getFullYear() }} 黔车会 · 保留所有权利</footer>
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
  flex: 7 1 0;
  min-width: 0;
  overflow: hidden;
}

.login-brand-bg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: 68% center;
  filter: brightness(0.82) saturate(0.82) contrast(1.05);
  transform: scale(1.08);
}

.login-brand-overlay {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 88% 92%, rgb(5 24 46 / 62%) 0%, rgb(5 24 46 / 0%) 34%),
    linear-gradient(90deg, rgb(5 24 46 / 62%) 0%, rgb(5 24 46 / 8%) 58%, rgb(5 24 46 / 18%) 100%),
    linear-gradient(180deg, rgb(5 24 46 / 58%) 0%, rgb(5 24 46 / 10%) 40%, rgb(5 24 46 / 88%) 100%);
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
  flex: 0 0 clamp(320px, 26vw, 400px);
  min-width: 0;
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
  max-width: 340px;
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
