# 亿民家 · 工程项目管理系统（ERP）

面向中小型工程施工 / 装饰 / 安装类企业的项目型管理系统。以「项目」为主线，把报价、合同、材料、采购、出入库、收付款、会计凭证和经营报表串成一条链：**业务单据一录入，财务凭证自动生成，报表实时可查**。

前后端分离，源码完整交付，可私有化部署到自己的服务器。

---

## 功能清单

| 模块 | 主要功能 |
| --- | --- |
| 项目管理 | 项目列表 / 详情、项目资料（合同附件、变更、工地日志、材料清单）、报价管理 |
| 库存管理 | 材料清单、仓库管理、库存台账、采购单、材料领用 |
| 财务管理 | 收款、付款、资金账户、会计科目、会计期间结账 / 反结账、记账凭证、凭证明细导出 |
| 基础资料 | 客户、供应商、财务初始化配置、数据字典 |
| 报表 | 经营总览、资金流水、科目余额、项目利润（毛利） |
| 组织与权限 | 部门、员工、账号、菜单管理、部门级菜单权限、动态路由 |
| 首页 | 业务概览、待办、最新进度、统计卡片 |
| 业主端接口 | 小程序侧登录、项目概览、施工进度、图片上传（`/apiv1/mp`，独立 `customer_token` 鉴权） |

### 核心设计：业财一体

- `internal/logic/fin/voucher` 提供统一记账入口。收付款等业务通过 `CreateWithTx` 与自身业务数据写在**同一个数据库事务**里，保证业务与账务同时成功或同时回滚。
- 自动完成借贷平衡校验、会计期间状态校验、金额格式校验、凭证号生成，并同步科目期间余额与资金流水。
- 因此收款、付款、采购、成本录入后，科目余额表、利润表、资金流水无需人工二次录入。
- 凭证金额以字符串入参，后端校验后计算，避免浮点误差。
- 采购支持无票 / 含税价 / 未税价，税额与未税额单独保存；当前付款凭证仍按应付总额计入成本，未拆分可抵扣进项税额。

---

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.25 + GoFrame v2.10 + MySQL 5.7 / 8.0 |
| 前端 | Vue 3.5 + TypeScript + Vite 8 + Naive UI + UnoCSS + Pinia + ECharts（基于 soybean-admin 2.2） |
| 接口 | REST，统一前缀 `/apiv1`，统一响应包装，GoFrame 自动生成 OpenAPI / Swagger |
| 鉴权 | 后台 `site_user_token`、业主端 `customer_token` 双链路独立中间件 |

**环境要求**：Go >= 1.25、MySQL >= 5.7、Node >= 20.19、pnpm >= 10.5。

---

## 目录结构

```
yiminjia/
├── backend/                        # Go 后端（module: engineServ）
│   ├── api/<module>/v1/            # 请求/响应结构 + 路由元信息
│   ├── internal/
│   │   ├── cmd/cmd.go              # 路由分组与三条鉴权链路
│   │   ├── controller/             # 控制器
│   │   ├── logic/api/              # 业务逻辑
│   │   ├── logic/fin/voucher/      # 复式记账引擎
│   │   └── logic/middleware/       # 鉴权、日志
│   ├── model/default/              # gf gen dao 生成的 DAO / 实体
│   ├── manifest/config/config.yaml # 运行配置
│   ├── manifest/sql/yiminjia.sql   # 建表 + 初始化数据 + 演示数据
│   └── storage/file/               # 本地上传文件（图片 / 合同附件）
└── soybean-admin/                  # Vue3 前端（pnpm workspace）
    ├── src/views/                  # 业务页面
    ├── src/service/api/            # 接口封装
    └── .env / .env.test / .env.prod
```

---

## 本地启动

### 1. 导入数据库

```sql
CREATE DATABASE IF NOT EXISTS yiminjia DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

将 `backend/manifest/sql/yiminjia.sql` 导入该库（命令行或 Navicat 均可）。文件包含 30 张业务表、财务科目 / 字典初始化数据，以及一批演示业务数据。

已有数据库升级税务字段时，单独执行 `manifest/sql/upgrade_purchase_tax.sql` 一次；新库不需要重复执行该升级脚本。

### 2. 修改配置

编辑 `backend/manifest/config/config.yaml`：

- `database.default.link`（以及 `database.chat.link`、`gfcli.gen.dao[0].link`）改为自己的 MySQL 账号、地址与库名。
- `logger.path` 改为本机可写目录，原值是 Linux 路径，Windows 下需调整。
- `server.address` 默认为 `:8007`，按需修改并同步前端 `.env*` 中的地址。
- **`app.authkey`、`app.qwen.apiKey` 为示例值**：对外分发或部署前务必替换，且不要把含真实密钥的配置文件提交到公开仓库。
- `server.sessionPath` 指向的 `storage/sessions` 目录需有写权限。

### 3. 启动后端

```bash
cd backend
go mod download
go run .
```

- 服务地址：`http://127.0.0.1:8007`
- 接口文档：`http://127.0.0.1:8007/swagger`
- OpenAPI：`http://127.0.0.1:8007/api.json`
- 也可编译为单文件：`go build -o engineServ .`（Windows 下加 `.exe`）

### 4. 启动前端

```bash
cd soybean-admin
pnpm install
pnpm dev
```

打开 `http://localhost:9528`。`.env.test` 默认指向 `http://127.0.0.1:8007/apiv1`；路由为动态模式（`VITE_AUTH_ROUTE_MODE=dynamic`），登录后由后端 `/route/getUserRoutes` 下发菜单。

### 5. 登录

导入原始 `yiminjia.sql` 后，登录页预填的演示账号为 `admin` / `123456`（如改过初始化数据，以 `site_user` 表中的实际值为准）。

> **正式使用前请立即修改密码并清空演示数据**：SQL 内含演示公司、客户、项目、凭证等业务数据，不要直接用于生产。

---

## 生产部署

1. 前端 `pnpm build`，产物在 `dist/`，由 Nginx 托管。
2. `.env.prod` 中 `VITE_SERVICE_BASE_URL=/apiv1`，需在 Nginx 把 `/apiv1` 反向代理到后端 `127.0.0.1:8007`。
3. 后端以 `go build` 产物常驻运行；附件默认存 `storage/file`，通过 `app.uploadDir` 改物理目录、`app.uploadUrlPath` 改访问前缀、`app.uploadUrlPrefix` 配域名 / CDN（见 `internal/cmd/cmd.go` 注释）。
4. 开启 MySQL 定时备份；上线前在「财务基础资料」核对会计科目、会计期间与资金账户初始余额。

---

## 二次开发指引

- **新增接口**：在 `api/<module>/v1/` 定义 Req / Res 与 `g.Meta`，在 `internal/controller/<module>/` 绑定，在 `internal/logic/api/<module>/` 写业务，最后到 `internal/cmd/cmd.go` 注册。
- **改表结构**：改 MySQL 后用 GoFrame CLI 的 `gf gen dao` 重新生成 `model/default` 下的 DAO 与实体。
- **接入自动记账**：参考 `internal/logic/fin/voucher/voucher.go` 文件头的完整示例，在自己的事务里调用 `CreateWithTx(ctx, tx, &voucher.CreateInput{...})`，无需自己处理凭证号、期间、余额。
- **新增页面**：在 `soybean-admin/src/views/` 下建页面与 `src/service/api/` 接口封装，再到数据库 `sys_menu` 配置菜单与部门权限（前端动态路由会自动加载）。
- **鉴权链路**：`/apiv1` 下分三条——登录 / 刷新令牌免鉴权、业主端 `CustomerAuth`、后台 `Auth`；新增接口务必挂到正确的分组。

---

## 交付边界

- 本仓库为**源码交付**：默认不含部署实施、二次开发、数据迁移、培训与售后，如需请另行约定。
- 授权用于内部业务系统或二次开发；转售、贴牌、去除版权声明需另行获得授权。
- 第三方组件遵循各自开源协议：前端 soybean-admin（MIT）、GoFrame、Naive UI、ECharts 等；二次分发请保留相应版权声明与 LICENSE。
- 当前版本**不包含**：发票 / 税务申报、银行流水对账、固定资产、工资社保、多租户隔离审计、小程序前端源码（后端业主端接口已就绪）。
