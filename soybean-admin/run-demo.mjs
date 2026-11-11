/**
 * 亿民家 · 自动演示视频脚本 (Playwright)
 * ========================================
 * 自动走完 登录 → 项目列表 → 收款 → 凭证 → 利润报表 全流程，输出一个 MP4 视频。
 *
 * 用法（在 soybean-admin 目录下执行）：
 *   npx playwright install chromium   # 首次运行，装浏览器内核
 *   node ../.screenshots/run-demo.mjs
 *
 * 依赖：Node >= 18、npx、本机已有 Playwright（你这台已经有 1.65）。
 */
import { chromium } from 'playwright';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const FRONTEND_URL = 'http://127.0.0.1:9528';
const ADMIN_USER = 'admin';
const ADMIN_PASS = '123456';

/** 每步之间留的"展示秒数"——给录屏/后期留出旁白空间 */
const HOLD = {
  intro: 3000,
  step: 2500,
  small: 1200,
  titleEnd: 1500
};

/** 小工具：等待 + 可选的高亮描边（方便后期剪视频时知道录到哪一步） */
async function highlight(page, selector, holdMs = HOLD.small) {
  try {
    await page.evaluate(sel => {
      const el = document.querySelector(sel);
      if (el) {
        el.style.transition = 'outline 0.2s';
        el.style.outline = '3px solid #C05621';
        el.style.outlineOffset = '2px';
      }
    }, selector);
    await page.waitForTimeout(holdMs);
    await page.evaluate(sel => {
      const el = document.querySelector(sel);
      if (el) el.style.outline = '';
    }, selector);
  } catch {
    await page.waitForTimeout(holdMs);
  }
}

async function safeClick(page, selector, holdMs = HOLD.small, desc = '') {
  console.log(`  ▶ click   ${desc || selector}`);
  await page.locator(selector).first.waitFor({ state: 'visible', timeout: 15000 });
  await page.locator(selector).first.click();
  await page.waitForTimeout(holdMs);
}

async function waitLoading(page, holdMs = 800) {
  try {
    await page.locator('.n-spin, .n-spin-container').last.waitFor({ state: 'hidden', timeout: 8000 });
  } catch {
    /* no loader, ignore */
  }
  await page.waitForTimeout(holdMs);
}

(async () => {
  const outDir = path.join(__dirname, 'demo-artifact');

  console.log('🎬 启动浏览器（1920×1080，带视频录制）…');
  const browser = await chromium.launch({
    headless: true,
    args: ['--start-maximized', '--disable-dev-shm-usage']
  });
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    deviceScaleFactor: 1,
    recordVideo: { dir: outDir, size: { width: 1920, height: 1080 } }
  });
  const page = await context.newPage();

  // 吞掉无用的错误日志，只打印关键信息
  page.on('pageerror', e => console.log('  ⚠ pageerror:', e.message.slice(0, 80)));

  // ==================== 00:00 – 登录页 ====================
  console.log('\n[1/6] 登录页');
  await page.goto(`${FRONTEND_URL}/login`, { waitUntil: 'networkidle', timeout: 30000 });
  await page.waitForTimeout(HOLD.intro);

  // 填账号密码
  const usernameInput = page.locator('input[placeholder*="用户名"], input[type="text"]').first;
  const passwordInput = page.locator('input[type="password"]').first;
  await usernameInput.fill(ADMIN_USER);
  await page.waitForTimeout(300);
  await passwordInput.fill(ADMIN_PASS);
  await page.waitForTimeout(500);

  // 点登录按钮
  await safeClick(page, 'button:has-text("登录")', 3000, '登录按钮');
  await page.waitForURL(`${FRONTEND_URL}/home`, { timeout: 15000 }).catch(() => null);
  await waitLoading(page, 1500);

  // ==================== 00:12 – 首页 ====================
  console.log('\n[2/6] 首页总览');
  // 把侧边栏展开一下（有些布局默认收起）
  try {
    await safeClick(page, 'button:has([data-icon="mdi:menu"])', 800, '展开侧边栏');
  } catch {
    /* already open */
  }
  await highlight(page, '.n-card', 2500);

  // ==================== 00:22 – 项目列表 ====================
  console.log('\n[3/6] 项目列表 → 锦绣路');
  // 点左侧菜单 项目管理
  await safeClick(page, 'text=项目管理', HOLD.small, '项目管理菜单');
  await safeClick(page, 'text=项目列表', HOLD.step, '项目列表子菜单');
  await waitLoading(page, 1200);

  // 在搜索框输入 "锦绣"
  const kwInput = page.locator('input[placeholder*="搜索"], input[placeholder*="项目"]').first;
  await kwInput.fill('锦绣');
  await kwInput.press('Enter');
  await page.waitForTimeout(HOLD.small);
  await waitLoading(page, 1000);
  await highlight(page, 'td:has-text("锦绣路")', 2200);

  // 点该行的"查看"按钮，打开详情抽屉
  const row = page.locator('tr:has-text("锦绣路")').first;
  await row.locator('button:has-text("查看")').click();
  await page.waitForTimeout(HOLD.step);

  // 抽屉里滚到合同区，给个镜头
  await page.mouse.wheel(0, 600);
  await page.waitForTimeout(HOLD.step);

  // 关掉抽屉
  await safeClick(page, '.n-drawer button.n-base-close, .n-drawer .n-drawer-content header button', HOLD.small, '关闭详情');

  // ==================== 01:20 – 收款列表 ====================
  console.log('\n[4/6] 收款管理');
  await safeClick(page, 'text=财务管理', HOLD.small, '财务管理菜单');
  await safeClick(page, 'text=收款管理', HOLD.step, '收款管理子菜单');
  await waitLoading(page, 1200);

  // 扫一眼已有的收款（定金、中期款、尾款这些标签）
  await highlight(page, '.n-tag', 2500);
  await page.waitForTimeout(HOLD.small);

  // 选一笔已有的收款，点"查看"看自动生成的凭证号关联
  const paymentRow = page.locator('tr:has-text("锦绣")').first;
  try {
    await paymentRow.locator('button:has-text("查看")').click();
    await page.waitForTimeout(HOLD.step);
    // 抽屉里找到"凭证"相关的字段高亮
    await highlight(page, '.n-form-item:has-text("凭证"), .n-form-item:has-text("关联凭证")', HOLD.step);
    await safeClick(page, '.n-drawer button.n-base-close', HOLD.small, '关闭详情');
  } catch {
    /* detail drawer variant not found, skip */
  }

  // ==================== 01:48 – 凭证列表 ====================
  console.log('\n[5/6] 记账凭证');
  await safeClick(page, 'text=记账凭证', HOLD.step, '记账凭证菜单');
  await waitLoading(page, 1500);

  // 点最新那条的查看
  try {
    await page.locator('button:has-text("查看")').first.click();
    await page.waitForTimeout(HOLD.step);
    // 高亮借贷两行
    await highlight(page, 'td:has-text("借")', HOLD.small);
    await highlight(page, 'td:has-text("贷")', HOLD.small);
    await safeClick(page, '.n-drawer button.n-base-close', HOLD.small, '关闭凭证详情');
  } catch {
    /* noop */
  }

  // ==================== 02:18 – 项目利润报表 ====================
  console.log('\n[6/6] 项目利润报表');
  await safeClick(page, 'text=报表中心', HOLD.small, '报表中心菜单');
  await safeClick(page, 'text=项目利润', HOLD.step, '项目利润子菜单');
  await waitLoading(page, 1500);

  // 先扫顶栏五张卡片
  for (let i = 0; i < 5; i++) {
    const card = page.locator('.n-card').nth(i);
    await highlight(page, card, 1200);
  }

  // 搜索 "锦绣"
  const profitKw = page.locator('input[placeholder*="项目"], input[placeholder*="关键字"]').first;
  await profitKw.fill('锦绣');
  await profitKw.press('Enter');
  await page.waitForTimeout(HOLD.small);
  await waitLoading(page, 800);

  // 高亮 毛利率那列
  await highlight(page, 'td:has-text("%")', 2500);
  await page.waitForTimeout(HOLD.step);

  // 点查看利润明细
  const profitRow = page.locator('tr:has-text("锦绣路")').first;
  try {
    await profitRow.locator('button:has-text("查看")').click();
    await page.waitForTimeout(HOLD.step * 2);
    await page.mouse.wheel(0, 500);
    await page.waitForTimeout(HOLD.step);
    await safeClick(page, '.n-drawer button.n-base-close', HOLD.small, '关闭利润详情');
  } catch {
    /* noop */
  }

  // 收尾静止 3 秒给落版
  console.log('\n✅ 演示流程走完，收尾 3 秒静止…');
  await page.waitForTimeout(3000);

  // 收尾
  const video = await page.video();
  const videoPath = await video.path();
  await video.saveAs(path.join(outDir, 'yiminjia-demo.webm'));
  await context.close();
  await browser.close();

  console.log('\n🎉 视频已生成：');
  console.log(`   ${path.join(outDir, 'yiminjia-demo.webm')}`);
  console.log('\n💡 想转成 mp4 给客户看？装个 ffmpeg 后跑：');
  console.log(`   ffmpeg -i "${path.join(outDir, 'yiminjia-demo.webm')}" -c:v libx264 -pix_fmt yuv420p -r 30 "${path.join(outDir, 'yiminjia-demo.mp4')}"`);
})().catch(err => {
  console.error('❌ 演示脚本挂了：', err.message);
  process.exit(1);
});
