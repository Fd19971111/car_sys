import type { RouteMeta } from 'vue-router';
import ElegantVueRouter from '@elegant-router/vue/vite';
import type { RouteKey } from '@elegant-router/types';

export function setupElegantRouter() {
  return ElegantVueRouter({
    layouts: {
      base: 'src/layouts/base-layout/index.vue',
      blank: 'src/layouts/blank-layout/index.vue'
    },
    routePathTransformer(routeName, routePath) {
      const key = routeName as RouteKey;

      if (key === 'login') {
        const modules: UnionKey.LoginModule[] = ['pwd-login', 'code-login', 'register', 'reset-pwd', 'bind-wechat'];

        const moduleReg = modules.join('|');

        return `/login/:module(${moduleReg})?`;
      }

      return routePath;
    },
    onRouteMetaGen(routeName) {
      const key = routeName as RouteKey;

      const constantRoutes: RouteKey[] = ["login", "403", "404", "500"];

      const hideInMenuRoutes: RouteKey[] = ["finance_setup"];
      const iconMap: Partial<Record<RouteKey, string>> = {
        settings: "mdi:cog",
        settings_customer: "mdi:account-tie",
        settings_menu: "mdi:menu",
        settings_pricing: "mdi:taxi",
        enterprise: "mdi:domain",
        enterprise_department: "mdi:office-building-outline",
        enterprise_staff: "mdi:account-group-outline",
        enterprise_user: "mdi:shield-account-outline",
        project: "mdi:home-city-outline",
        "project_project-list":"mdi:clipboard-text-outline",
        project_data: "mdi:folder-open-outline",
        project_quote: "mdi:file-document-edit-outline",
        inventory: "mdi:package-variant-closed",
        inventory_materials: "mdi:format-list-bulleted",
        inventory_procurement: "mdi:cart-outline",
        inventory_usages: "mdi:clipboard-list-outline",
        finance: "mdi:wallet-outline",
        finance_payment: "mdi:cash-multiple",
        finance_payout: "mdi:cash-minus",
        finance_voucher: "mdi:book-check-outline",
        finance_basic: "mdi:book-open-variant",
        report: "mdi:chart-box-outline",
        report_overview: "mdi:view-dashboard-outline",
        report_flows: "mdi:swap-horizontal",
        report_balance: "mdi:book-open-page-variant-outline",
        report_profit: "mdi:chart-line",
        home: "mdi:monitor-dashboard",
      };
      const orderMap: Partial<Record<RouteKey, number>> = {
        home: 1,
        project: 2,
        "project_project-list": 1,
        project_data: 2,
        project_quote: 3,
        finance: 3,
        finance_basic: 1,
        finance_payment: 2,
        finance_payout: 3,
        finance_voucher: 4,
        settings: 4,
        settings_customer: 1,
        settings_menu: 2,
        settings_pricing: 2,
        enterprise: 6,
        enterprise_staff: 1,
        enterprise_department: 2,
        enterprise_user: 3,
        inventory: 5,
        inventory_materials: 1,
        inventory_usages: 2,
        inventory_procurement: 3,
        report: 4,
        report_overview: 1,
        report_flows: 2,
        report_balance: 3,
        report_profit: 4,
      };

      // keep the shape of enterprise children consistent so TypeScript can infer precise literal types
      const isAdminMap: Partial<Record<RouteKey, boolean>> = {

      };

      const meta: Partial<RouteMeta> = {
        title: key,
        i18nKey: `route.${key}` as App.I18n.I18nKey,
        icon: iconMap[key] ?? "",
        order: orderMap[key] ?? null,
        isAdmin: isAdminMap[key],
      };

      if (constantRoutes.includes(key)) {
        meta.constant = true;
      }

      if (hideInMenuRoutes.includes(key)) {
        meta.hideInMenu = true;
      }

      return meta;
    }
  });
}
