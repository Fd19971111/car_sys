# 项目规则

## 前端页面组件化（强制）

新增或改造前端页面时，**禁止把所有功能堆在 `index.vue` 里**。

- `index.vue` 只负责页面编排、路由参数、跨组件状态和接口调用顺序。
- 表格列定义、搜索栏、操作抽屉、表单、数据请求 hook、状态元数据必须拆到同目录的 `components/` 下。
- 数据请求与增删改逻辑放进 `components/use-xxx-data.ts`。
- 列表列定义放进 `components/xxx-table.vue`，不要在 `index.vue` 里写大段 `h()` 渲染。
- 搜索/筛选/操作按钮放进 `components/xxx-search-bar.vue`，通过 props / emits 与页面通信。
- 新增/编辑表单放进 `components/xxx-operate-drawer.vue`。
- 状态文案、颜色、选项放进 `components/xxx-meta.ts`。
- 单个 `.vue` 文件超过约 300 行时必须继续拆分。

参考实现：`soybean-admin/src/views/transport/driver/`、`soybean-admin/src/views/transport/vehicle/`。
