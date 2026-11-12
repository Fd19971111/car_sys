// @unocss-include
import { getColorPalette, getRgb } from '@sa/color';
import { DARK_CLASS } from '@/constants/app';
import { localStg } from '@/utils/storage';
import { toggleHtmlClass } from '@/utils/common';
import { $t } from '@/locales';

export function setupLoading() {
  const themeColor = localStg.get('themeColor') || '#1677D2';
  const darkMode = localStg.get('darkMode') || false;
  const palette = getColorPalette(themeColor);

  const { r, g, b } = getRgb(themeColor);

  const primaryColor = `--primary-color: ${r} ${g} ${b}`;

  const svgCssVars = Array.from(palette.entries())
    .map(([key, value]) => `--logo-color-${key}: ${value}`)
    .join(';');

  const cssVars = `${primaryColor}; ${svgCssVars}`;

  if (darkMode) {
    toggleHtmlClass(DARK_CLASS).add();
  }

  const loadingClasses = [
    'left-0 top-0',
    'left-0 bottom-0 animate-delay-500',
    'right-0 top-0 animate-delay-1000',
    'right-0 bottom-0 animate-delay-1500'
  ];

  const dot = loadingClasses
    .map(item => {
      return `<div class="absolute w-16px h-16px bg-primary rounded-8px animate-pulse ${item}"></div>`;
    })
    .join('\n');

  const loading = `
<div class="fixed-center flex-col bg-layout" style="${cssVars}">
  <div class="w-128px h-128px">
    ${getLogoSvg()}
  </div>
  <div class="w-56px h-56px my-36px">
    <div class="relative h-full animate-spin">
      ${dot}
    </div>
  </div>
  <h2 class="text-28px font-500 text-primary">${$t('system.title')}</h2>
</div>`;

  const app = document.getElementById('app');

  if (app) {
    app.innerHTML = loading;
  }
}

function getLogoSvg() {
  return `<svg width="100%" height="100%" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="2" y="2" width="60" height="60" rx="16" fill="#1677D2" />
    <path d="m11 37 3.4-8.5a5 5 0 0 1 4.7-3.2h7.1l4.3-5.2a5 5 0 0 1 3.9-1.9h7.4a5 5 0 0 1 4.1 2.2l3.8 5.3 4 1.3a4.5 4.5 0 0 1 3.1 4.3v8.2a3.5 3.5 0 0 1-3.5 3.5h-1.1a5.8 5.8 0 0 1-11.3 0H27.1a5.8 5.8 0 0 1-11.3 0h-1.3A3.5 3.5 0 0 1 11 39.5V37Z" fill="#fff" />
    <path d="m28.7 25.3 3.6-4.4a2.8 2.8 0 0 1 2.2-1.1h7.3c.9 0 1.8.5 2.3 1.2l3 4.3H28.7Z" fill="#1677D2" />
    <circle cx="21.4" cy="42.5" r="3" fill="#1677D2" />
    <circle cx="45.1" cy="42.5" r="3" fill="#1677D2" />
  </svg>`;
}
