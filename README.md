# 工作安排管理系统

一个基于 **Go + React** 的工作安排管理 Web 应用。命令行启动，浏览器访问，**单文件部署**，无需任何运行时依赖。

## 目录

- [技术路线](#技术路线)
- [快速开始](#快速开始)
- [使用方法](#使用方法)
  - [构建与运行](#构建与运行)
  - [开发模式](#开发模式)
  - [交叉编译](#交叉编译)
  - [后台运行与开机自启](#后台运行与开机自启)
- [功能说明](#功能说明)
  - [1. 工作安排管理](#1-工作安排管理)
  - [2. 新增记录](#2-新增记录)
  - [3. 编辑记录](#3-编辑记录)
  - [4. 筛选查询](#4-筛选查询)
  - [5. 导出数据](#5-导出数据)
  - [6. 导入数据](#6-导入数据)
  - [7. 复制功能](#7-复制功能)
  - [8. 删除记录](#8-删除记录)
  - [9. 开机自启动](#9-开机自启动)
- [系统架构](#系统架构)
- [API 文档](#api-文档)
  - [通用说明](#通用说明)
  - [工作安排 CRUD](#工作安排-crud)
  - [导出接口](#导出接口)
  - [导入接口](#导入接口)
  - [参考数据接口](#参考数据接口)
  - [设置接口](#设置接口)
- [数据库设计](#数据库设计)
- [项目结构](#项目结构)
- [部署说明](#部署说明)

---

## 技术路线

| 层级 | 技术选型 | 版本 |
|------|---------|------|
| 后端语言 | Go | 1.25+ |
| HTTP 路由 | `net/http`（标准库，零外部依赖） | — |
| 数据库 | SQLite（`modernc.org/sqlite`，纯 Go 实现，无 CGo） | v1.51 |
| Excel 处理 | `excelize/v2` | v2.10 |
| 前端框架 | React + TypeScript | 18.2 |
| 构建工具 | Vite | 3.x |
| UI 组件库 | Ant Design + @ant-design/icons | 6.4 |
| 日期处理 | dayjs | 1.11 |
| 进程管理 | 平台原生 API（macOS LaunchAgent / Windows Registry） | — |

### 设计原则

- **单文件部署**：前端构建产物通过 Go 1.16+ `//go:embed` 嵌入二进制，编译后仅一个可执行文件
- **零运行时依赖**：SQLite 使用纯 Go 实现，无需 C 编译器或系统库
- **跨平台**：Go 原生交叉编译 + 平台构建标签（`//go:build darwin` / `//go:build windows`）
- **标准库优先**：HTTP 路由使用 `net/http`，不使用第三方 Web 框架

---

## 快速开始

```bash
# 一键构建并启动
make start

# 浏览器访问
open http://localhost:8080
```

**前置条件**：Go 1.25+ 和 Node.js（`make start` 会自动执行 `npm install && npm run build`）。

---

## 使用方法

### 构建与运行

```bash
# ========== Makefile（推荐）==========

make build          # 构建前端 + 编译 Go → 产物 ./work-manager
make run            # 直接运行已编译的二进制
make start          # build + run 一键执行
make frontend       # 仅构建前端
make clean          # 清理二进制、数据库文件、前端构建产物

# ========== 手动逐步执行 ==========

# 1. 构建前端
cd frontend && npm install && npm run build && cd ..

# 2. 编译 Go（Go 编译器自动嵌入 frontend/dist）
go build -o work-manager .

# 3. 运行
./work-manager
```

**自定义端口**：

```bash
PORT=9090 ./work-manager    # 默认 8080
```

### 开发模式

开发时前后端分离运行以支持热重载（HMR）：

```bash
# 终端 1：启动 Go 后端（监听 :8080）
go run .

# 终端 2：启动 Vite 前端开发服务器（监听 :5173）
cd frontend && npm run dev

# 浏览器访问 http://localhost:5173
# Vite 自动将 /api/* 请求代理到 :8080
```

运行 `make dev` 可查看上述指令。

### 交叉编译

在任意平台构建目标平台二进制：

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o work-manager .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o work-manager .

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o work-manager.exe .

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o work-manager .
```

产物为单文件，复制到目标机器直接运行即可。

### 后台运行与开机自启

**手动后台运行**：

```bash
# Linux / macOS
nohup ./work-manager > /dev/null 2>&1 &

# Windows（CMD）
start /B work-manager.exe
```

**开机自启动**：启动应用后，在 Web 界面左侧菜单进入「设置」→ 开启「开机自启动」开关。底层实现：

| 平台 | 实现方式 |
|------|---------|
| macOS | 写入 `~/Library/LaunchAgents/com.chaitin.workmanager.plist` 并通过 `launchctl bootstrap` 注册 |
| Windows | 写入注册表 `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` |

---

## 功能说明

### 1. 工作安排管理

主界面以表格形式展示所有工作安排记录，按**日期降序、更新时间降序**排列。每条记录包含以下字段：

| 字段 | JSON Key | 类型 | 说明 | 可选值 |
|------|----------|------|------|--------|
| 内部 ID | `id` | int64 | 数据库自增主键 | 自动生成，不可编辑 |
| 项目 ID | `project_id` | int64 | 用户自行填写的项目/工单 ID | 手动输入 |
| 日期 | `date` | string | 工作日期 | YYYY-MM-DD |
| 客户名称 | `customer` | string | 客户公司/联系人 | 自由文本 |
| 项目名称 | `project` | string | 所属项目 | 自由文本 |
| 工作类型 | `work_type` | string | 工作分类 | `测试` `交付` `售后` |
| 工作地点 | `location` | string | 工作方式 | `远程` `现场` |
| 伙伴 | `partner` | string | 是否有生态伙伴参与 | `是` `否` |
| 工作内容 | `content` | string | 具体工作描述 | 自由文本 |
| 工作耗时 | `duration` | float64 | 小时数 | 数字，支持 0.5 |
| 工作进度 | `progress` | string | 当前状态 | `未开始` `进行中` `已完成` `已暂停` `已取消` |
| 备注 | `notes` | string | 额外说明 | 自由文本 |
| 创建时间 | `created_at` | string | 自动生成 | 数据库本地时间 |
| 更新时间 | `updated_at` | string | 自动维护 | 每次更新自动刷新 |

表格中关键字段使用 Ant Design `Tag` 进行颜色编码：

- **工作类型**：测试=蓝色、交付=绿色、售后=橙色
- **工作地点**：远程=蓝色、现场=绿色
- **伙伴**：是=绿色、否=默认
- **工作进度**：未开始=默认、进行中=蓝色、已完成=绿色、已暂停=橙色、已取消=红色

支持分页浏览，每页可选 20 / 50 / 100 / 200 条。

### 2. 新增记录

点击「新增」按钮，弹出 Modal 表单：

- **日期**自动填充当天
- **工作类型**默认 `测试`、**工作地点**默认 `远程`、**伙伴**默认 `否`、**工作进度**默认 `未开始`
- 表单采用双列网格布局（ID / 日期 / 耗时 / 客户 / 项目 / 类型 / 地点 / 伙伴 / 进度），内容区域为多行文本域（3 行），备注为 2 行文本域
- 前端和后端均进行枚举值校验，非法数据会被拒绝
- 提交后记录出现在列表顶部

### 3. 编辑记录

点击行末「编辑」按钮，弹出与新增相同的 Modal 表单，预填当前数据。提交后列表即时刷新，`updated_at` 自动更新为当前本地时间。

### 4. 筛选查询

页面顶部提供筛选栏，支持以下条件组合筛选：

| 筛选项 | 匹配方式 | 说明 |
|--------|---------|------|
| 日期范围 | `>=` / `<=` | 选择起止日期（Ant Design RangePicker） |
| 客户名称 | `LIKE %xxx%` | 模糊匹配 |
| 项目名称 | `LIKE %xxx%` | 模糊匹配 |
| 工作类型 | `=` | 精确匹配（下拉选择） |
| 工作进度 | `LIKE %xxx%` | 模糊匹配（下拉选择） |

点击「查询」应用筛选，点击「重置」清除所有条件并加载全部数据。

**实现细节**：后端使用动态 SQL 构建，`WHERE 1=1` 基础上逐个追加 `AND` 条件，参数化查询防注入。

### 5. 导出数据

点击「导出数据」下拉菜单：

- **导出 Excel (.xlsx)**：生成带样式的 Excel 文件
  - 表头：蓝色背景（`#4472C4`）+ 白色加粗字体 + 居中 + 边框
  - 数据行：11px 字体 + 边框 + 垂直居中
  - 自适应列宽（12~40 字符范围）
  - 导出列：项目 ID、日期、客户名称、项目名称、工作类型、工作地点、伙伴、工作内容、工作耗时(h)、工作进度、备注、创建时间、更新时间
- **导出 JSON (.json)**：生成格式化 JSON 数组（2 空格缩进）

导出会**应用当前筛选条件**——筛选后再导出，只导出筛选结果。未设置任何筛选条件时导出全部数据。

### 6. 导入数据

点击「导入数据」按钮，进入四步导入流程：

1. **选择文件**：拖拽上传区域，接受 `.xlsx` 或 `.json` 文件（最大 10MB）
2. **系统解析**：后端解析文件内容
   - **Excel 解析**：读取第一个工作表，按中文表头名称自动映射列（支持 `ID`→`project_id`、`工作耗时(h)`→`工作耗时` 兼容处理），必需列：日期、客户名称、工作类型、工作地点
   - **JSON 解析**：反序列化对象数组
   - 空字段自动填充默认值（伙伴→`否`、进度→`未开始`）
   - 逐行校验（日期非空、客户非空、枚举值合法）
3. **预览结果**：展示解析出的记录表格 + 有效/无效统计
4. **确认导入**：点击确认后批量写入，使用**数据库事务**确保原子性

**导入失败处理**：逐行校验，合法行写入、非法行跳过并记录错误原因（行号 + 具体错误），不影响已成功的行。

### 7. 复制功能

点击行末「复制」按钮，自动将记录格式化为规范文本并写入系统剪贴板（`navigator.clipboard.writeText`）：

| 条件 | 格式 |
|------|------|
| 伙伴 = 是 | `【工作地点】【工作类型】【生态】项目名称-工作内容` |
| 伙伴 = 否 | `【工作地点】【工作类型】项目名称-工作内容` |

示例：
```
【远程】【测试】【生态】XX安全评估项目-漏洞扫描与渗透测试
【现场】【售后】YY防火墙项目-设备上架与策略调试
```

操作成功后显示 `message.success` 提示。

### 8. 删除记录

点击行末「删除」按钮，弹出 `Popconfirm` 确认对话框（"确定要删除这条记录吗？"），确认后执行删除。

**注意：删除操作不可恢复**，后端执行物理删除（`DELETE FROM work_arrangements WHERE id = ?`）。

### 9. 开机自启动

左侧菜单点击「设置」进入设置页面，提供开机自启动开关：

- 调用后端 API 进行系统级注册
- macOS：通过 LaunchAgent plist（`~/Library/LaunchAgents/com.chaitin.workmanager.plist`）
- Windows：通过注册表 Run 键（`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`）
- 开关加载时自动查询当前系统状态

---

## 系统架构

```
┌──────────────────────────────────────────────────┐
│                    浏览器                          │
│  React 18 SPA (Ant Design 6)                      │
│  Vite 构建 → 嵌入 Go 二进制                       │
└────────────────────┬─────────────────────────────┘
                     │ HTTP (localhost:8080)
┌────────────────────▼─────────────────────────────┐
│                  Go 后端                           │
│                                                    │
│  main.go          → 入口、embed 静态资源            │
│  server.go        → 路由注册、CORS、工具函数        │
│                                                    │
│  handlers/        → HTTP 层（请求解析/响应）        │
│    ├── work_arrangement.go  CRUD + 参考数据        │
│    ├── export.go            导出接口               │
│    ├── import_handler.go    导入接口               │
│    └── settings.go          设置接口               │
│                                                    │
│  services/        → 业务逻辑层                      │
│    ├── work_arrangement_svc.go  验证、CRUD、筛选   │
│    ├── export_svc.go            Excel/JSON 生成    │
│    └── import_svc.go            文件解析、校验      │
│                                                    │
│  models/          → 数据模型、枚举常量              │
│    └── work_arrangement.go                        │
│                                                    │
│  db/              → 数据库层                        │
│    ├── db.go          连接管理、WAL 模式           │
│    └── migrations.go  建表、索引、向后兼容迁移     │
│                                                    │
│  platform/        → 平台抽象（build tags）          │
│    ├── autostart.go         接口定义               │
│    ├── autostart_darwin.go  macOS LaunchAgent      │
│    └── autostart_windows.go Windows Registry       │
└────────────────────┬─────────────────────────────┘
                     │
┌────────────────────▼─────────────────────────────┐
│     SQLite 数据库                                  │
│     ~/.work_manager/work_manager.db               │
│     WAL 模式 + 外键约束                            │
│     单连接（SetMaxOpenConns=1）                    │
└──────────────────────────────────────────────────┘
```

**分层职责**：

| 层 | 职责 | 关键约束 |
|----|------|---------|
| `main.go` | 进程入口：初始化数据库、加载前端资源、启动 HTTP 服务、优雅关闭 | — |
| `server.go` | 路由注册：URL 路由、HTTP 方法分发、CORS 头注入 | 标准库 `net/http`，无框架 |
| `handlers/` | HTTP 适配：解析请求参数（路径/Query/JSON Body/Multipart）、调用 Service、写响应 | 不包含业务逻辑 |
| `services/` | 业务逻辑：字段校验、动态 SQL 构建、事务管理、Excel/JSON 读写 | 无 HTTP 依赖，可单独测试 |
| `models/` | 数据结构：结构体定义、枚举常量、校验函数、格式化方法 | 纯数据，无外部依赖 |
| `db/` | 数据持久化：连接管理、Schema 迁移、索引管理 | WAL 模式，单写连接 |
| `platform/` | 系统集成：开机自启的平台原生实现 | `//go:build` 条件编译 |

---

## API 文档

### 通用说明

- **Base URL**：`http://localhost:8080/api`
- **Content-Type**：`application/json`（除文件上传使用 `multipart/form-data`）
- **CORS**：所有 `/api/*` 路由返回 `Access-Control-Allow-Origin: *`
- **错误响应格式**：`{ "error": "错误描述" }`
- **成功响应**：数据直接作为 JSON 返回（无包装）

### 工作安排 CRUD

#### 获取全部记录（支持筛选）

```
GET /api/work-arrangements
```

**Query 参数**（均为可选）：

| 参数 | 类型 | 匹配方式 | 示例 |
|------|------|---------|------|
| `date_from` | string | `>=` | `2026-06-01` |
| `date_to` | string | `<=` | `2026-06-30` |
| `customer` | string | `LIKE %xxx%` | `中国移动` |
| `project` | string | `LIKE %xxx%` | `安全评估` |
| `work_type` | string | `=` | `测试` |
| `progress` | string | `LIKE %xxx%` | `完成` |

**响应** `200 OK`：

```json
[
  {
    "id": 1,
    "project_id": 20240001,
    "date": "2026-06-06",
    "customer": "中国移动",
    "project": "网络安全评估",
    "work_type": "测试",
    "location": "远程",
    "partner": "是",
    "content": "漏洞扫描与渗透测试",
    "duration": 8.5,
    "progress": "已完成",
    "notes": "已提交报告",
    "created_at": "2026-06-06 10:30:00",
    "updated_at": "2026-06-06 18:00:00"
  }
]
```

#### 获取单条记录

```
GET /api/work-arrangements/{id}
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 数据库内部主键 |

**响应** `200 OK`：单条 `WorkArrangement` 对象（格式同上）

**错误**：`400 Bad Request`（ID 非法）、`404 Not Found`（记录不存在）

#### 新增记录

```
POST /api/work-arrangements
```

**请求体**（JSON）：

```json
{
  "project_id": 20240001,
  "date": "2026-06-06",
  "customer": "中国移动",
  "project": "网络安全评估",
  "work_type": "测试",
  "location": "远程",
  "partner": "否",
  "content": "漏洞扫描",
  "duration": 4.0,
  "progress": "进行中",
  "notes": ""
}
```

必填字段：`date`、`customer`、`work_type`、`location`、`partner`、`progress`（后端会校验枚举值）。

**响应** `201 Created`：新创建的完整 `WorkArrangement` 对象（含自动生成的 `id`、`created_at`、`updated_at`）

**错误**：`400 Bad Request`（JSON 解析失败或字段校验不通过）

#### 更新记录

```
PUT /api/work-arrangements/{id}
```

**路径参数**：`id` — 数据库内部主键

**请求体**（JSON）：同新增（服务端会从路径覆盖 `id` 字段）

**响应** `200 OK`：更新后的完整 `WorkArrangement` 对象（`updated_at` 已刷新）

**错误**：`400 Bad Request`、`404 Not Found`

#### 删除记录

```
DELETE /api/work-arrangements/{id}
```

**响应** `204 No Content`

**错误**：`400 Bad Request`、`404 Not Found`

#### 批量新增

```
POST /api/work-arrangements/bulk
```

**请求体**（JSON 数组）：

```json
[
  {
    "project_id": 20240001,
    "date": "2026-06-06",
    "customer": "客户A",
    "project": "项目A",
    "work_type": "测试",
    "location": "远程",
    "partner": "否",
    "content": "工作内容",
    "duration": 4.0,
    "progress": "进行中",
    "notes": ""
  }
]
```

**响应** `200 OK`：

```json
{
  "created": 5,
  "skipped": 2,
  "errors": [
    "行 3: 无效类型 '未知'",
    "行 7: 客户为空"
  ]
}
```

所有操作在一个数据库事务中完成，成功行写入、失败行跳过并记录原因。

### 导出接口

#### 导出 Excel

```
GET /api/export/excel
```

Query 参数同 `GET /api/work-arrangements` 的筛选参数。若无筛选参数则导出全部。

**响应** `200 OK`：
- `Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- `Content-Disposition: attachment; filename=工作安排.xlsx`

#### 导出 JSON

```
GET /api/export/json
```

**响应** `200 OK`：
- `Content-Type: application/json`
- `Content-Disposition: attachment; filename=工作安排.json`
- Body：格式化 JSON 数组（2 空格缩进）

### 导入接口

#### 解析上传文件

```
POST /api/import/parse
```

**请求**：`multipart/form-data`，字段名 `file`，最大 10MB。

支持文件格式：
- `.xlsx` — Excel（读取第一个工作表，按中文表头映射列）
- `.json` — JSON 数组

**响应** `200 OK`：

```json
{
  "Records": [
    {
      "project_id": 20240001,
      "date": "2026-06-06",
      "customer": "客户A",
      "project": "项目A",
      "work_type": "测试",
      "location": "远程",
      "partner": "否",
      "content": "工作内容",
      "duration": 4.0,
      "progress": "进行中",
      "notes": ""
    }
  ],
  "Created": 3,
  "Skipped": 1,
  "Errors": ["行 5: 日期不能为空"]
}
```

注意：此接口仅**解析预览**，不写入数据库。`Created`/`Skipped` 指解析层面的通过/跳过计数。

#### 确认导入

```
POST /api/import/confirm
```

**请求体**（JSON）：`WorkArrangement` 对象数组

**响应** `200 OK`：

```json
{
  "created": 3,
  "skipped": 1,
  "errors": ["行 5: 日期不能为空"]
}
```

此接口执行**实际写入**数据库，带有事务保护。与 `bulk` 接口共享同一 `BulkCreate` 服务方法。

### 参考数据接口

#### 获取客户列表

```
GET /api/reference/customers
```

**响应** `200 OK`：`["客户A", "客户B", "中国移动"]`（去重、按字母排序、排除空字符串）

#### 获取项目列表

```
GET /api/reference/projects
```

**响应** `200 OK`：`["项目A", "项目B", "网络安全评估"]`（同上）

### 设置接口

#### 查询开机自启状态

```
GET /api/settings/autostart
```

**响应** `200 OK`：`{ "enabled": true }`

#### 设置开机自启

```
PUT /api/settings/autostart
```

**请求体**（JSON）：
```json
{ "enabled": true }
```

**响应** `200 OK`：`{ "enabled": true }`

**错误**：`400 Bad Request`（JSON 解析失败）、`500 Internal Server Error`（系统操作失败，如权限不足）

---

## 数据库设计

### 部署位置

```
~/.work_manager/work_manager.db
```

数据库文件位于用户主目录下的 `.work_manager` 隐藏目录，**跨构建/重启持久化**。

### 连接配置

- **驱动**：`modernc.org/sqlite`（纯 Go 实现，无需 CGo）
- **Journal Mode**：WAL（Write-Ahead Logging，并发读性能更好）
- **外键约束**：开启（`_foreign_keys=on`）
- **连接池**：`MaxOpenConns=1`（SQLite 单写模式最佳实践）

### 表结构

**`work_arrangements`**：

```sql
CREATE TABLE IF NOT EXISTS work_arrangements (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id  INTEGER NOT NULL DEFAULT 0,
    date        TEXT NOT NULL,
    customer    TEXT NOT NULL DEFAULT '',
    project     TEXT NOT NULL DEFAULT '',
    work_type   TEXT NOT NULL CHECK(work_type IN ('测试','交付','售后')),
    location    TEXT NOT NULL CHECK(location IN ('远程','现场')),
    partner     TEXT NOT NULL CHECK(partner IN ('是','否')),
    content     TEXT NOT NULL DEFAULT '',
    duration    REAL NOT NULL DEFAULT 0,
    progress    TEXT NOT NULL CHECK(progress IN ('未开始','进行中','已完成','已暂停','已取消')),
    notes       TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);
```

**索引**：

```sql
CREATE INDEX IF NOT EXISTS idx_date     ON work_arrangements(date);
CREATE INDEX IF NOT EXISTS idx_customer ON work_arrangements(customer);
CREATE INDEX IF NOT EXISTS idx_project  ON work_arrangements(project);
CREATE INDEX IF NOT EXISTS idx_work_type ON work_arrangements(work_type);
CREATE INDEX IF NOT EXISTS idx_progress ON work_arrangements(progress);
```

### 向后兼容迁移

系统启动时自动检测并执行 Schema 升级：

1. 通过 `PRAGMA table_info(work_arrangements)` 检查 `project_id` 列是否存在
2. 若不存在（旧版本数据库），执行 `ALTER TABLE ADD COLUMN project_id INTEGER NOT NULL DEFAULT 0`
3. 新安装直接执行 `CREATE TABLE IF NOT EXISTS` + 索引创建

---

## 项目结构

```
chaitin_job/
├── main.go                           # 进程入口：DB 初始化、前端资源嵌入、HTTP 启动、信号处理
├── server.go                         # 路由注册、CORS 中间件、JSON 工具函数、服务初始化
├── go.mod                            # Go 模块定义 (chaitin-job/work-manager)
├── go.sum                            # 依赖锁定
├── Makefile                          # 构建自动化 (build/run/start/dev/clean)
├── README.md                         # 本文件
├── requirements.md                   # 原始需求文档
│
├── handlers/                         # HTTP 处理层
│   ├── work_arrangement.go           # CRUD 接口 + 参考数据接口
│   ├── export.go                     # 导出接口（Excel / JSON + 筛选参数解析）
│   ├── import_handler.go             # 导入接口（文件上传解析 + 确认写入）
│   └── settings.go                   # 设置接口（开机自启查询/设置）
│
├── services/                         # 业务逻辑层
│   ├── work_arrangement_svc.go       # 核心 CRUD、动态筛选、批量创建（事务）、去重查询
│   ├── export_svc.go                 # Excel 生成（含样式）、JSON 序列化、Writer 接口
│   └── import_svc.go                 # 文件解析（Excel/JSON）、列映射、逐行校验
│
├── models/                           # 数据模型层
│   └── work_arrangement.go           # 结构体、枚举常量、校验函数、格式化复制文本
│
├── db/                               # 数据库层
│   ├── db.go                         # 连接管理（~/.work_manager/work_manager.db）
│   └── migrations.go                 # Schema 迁移、索引创建、向后兼容升级
│
├── platform/                         # 平台抽象层（build tags）
│   ├── autostart.go                  # 接口定义 + 常量
│   ├── autostart_darwin.go           # macOS LaunchAgent 实现
│   └── autostart_windows.go          # Windows Registry 实现
│
└── frontend/                         # React 前端
    ├── index.html                    # Vite 入口 HTML
    ├── package.json                  # 依赖定义
    ├── tsconfig.json                 # TypeScript 配置
    ├── vite.config.ts                # Vite 配置（含 /api 代理）
    └── src/
        ├── main.tsx                  # React 18 入口
        ├── App.tsx                   # 根组件：路由、全局状态、主题配置
        ├── style.css                 # 全局样式
        ├── vite-env.d.ts             # Vite 类型声明
        ├── api/
        │   ├── client.ts             # HTTP 客户端封装（apiRequest / downloadFile）
        │   └── workArrangements.ts    # API 函数（CRUD / 导入导出 / 设置）
        ├── types/
        │   └── workArrangement.ts     # TypeScript 类型定义 + 枚举常量
        ├── hooks/
        │   └── useWorkArrangements.ts # 数据管理 Hook（状态 + 操作）
        └── components/
            ├── Layout/
            │   └── AppLayout.tsx      # 整体布局（侧边栏 + 顶栏 + 内容区）
            └── WorkArrangement/
                ├── FilterBar.tsx      # 筛选条件栏
                ├── WorkTable.tsx      # 数据表格（含颜色标签、分页、操作按钮）
                ├── WorkForm.tsx       # 新增/编辑 Modal 表单
                ├── EditableCell.tsx   # 可编辑单元格组件
                └── ImportModal.tsx    # 导入数据 Modal（上传→解析→预览→确认）
```

---

## 部署说明

### 单机部署

1. 在开发机上执行交叉编译（目标平台）
2. 将产物 `work-manager`（或 `work-manager.exe`）复制到目标机器
3. 直接运行，首次启动自动在 `~/.work_manager/` 下创建数据库
4. 访问 `http://localhost:8080`

### 后台运行

```bash
# Linux / macOS
nohup ./work-manager > /tmp/work-manager.log 2>&1 &

# 或配合 screen / tmux
screen -dmS work-manager ./work-manager

# Windows（CMD）
start /B work-manager.exe
```

### 配置项

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `PORT` | `8080` | HTTP 监听端口 |

### 日志

应用日志输出到 stdout/stderr，无文件日志。如使用 `nohup` 启动，建议重定向至文件。

### 数据备份

直接复制 `~/.work_manager/work_manager.db` 文件即可。SQLite 使用 WAL 模式时，备份前建议先关闭应用或执行 checkpoint。

---

> 更多信息参考 `requirements.md`（原始需求文档）

