# PostMe - REST 请求工具技术设计文档

## 1. 项目概述

PostMe 是一款基于 Wails 构建的轻量级 REST 请求发送工具，类似于 Postman/Insomnia，但更简洁。

## 2. 技术栈

| 层级 | 技术选型 | 版本 |
|------|----------|------|
| 框架 | Wails | v2 |
| 后端 | Go | 1.25.6 |
| 前端 | Vue 3 + TypeScript + Composition API | Vue 3.4+ |
| UI 样式 | Tailwind CSS | 3.x |
| UI 组件 | Headless UI | 1.x |
| 图标 | Heroicons | 2.x |
| 状态管理 | Pinia | 2.x |
| 代码编辑器 | CodeMirror | 6.x |
| 数据库 | SQLite (modernc.org/sqlite) | - |
| SQL 增强 | sqlx | 1.x |

## 3. 项目目录结构

```
postme/
├── main.go                     # 应用入口
├── app.go                      # 主应用结构体
├── wails.json                  # Wails 配置
├── go.mod
├── go.sum
│
├── data/                       # 数据目录（与应用同级，便于迁移）
│   └── postme.db               # SQLite 数据库
│
├── internal/
│   ├── models/                 # 数据模型
│   │   ├── request.go          # 请求模型
│   │   ├── response.go         # 响应模型
│   │   ├── collection.go       # 集合模型
│   │   ├── folder.go           # 文件夹模型
│   │   ├── history.go          # 历史记录模型
│   │   ├── environment.go      # 环境变量模型
│   │   └── app_state.go        # 应用状态模型
│   │
│   ├── database/               # 数据库层
│   │   ├── db.go               # 数据库连接管理
│   │   ├── migrations.go       # 数据库迁移
│   │   └── repository/         # 数据访问层
│   │       ├── request_repo.go
│   │       ├── collection_repo.go
│   │       ├── folder_repo.go
│   │       ├── history_repo.go
│   │       ├── environment_repo.go
│   │       └── app_state_repo.go
│   │
│   ├── services/               # 业务逻辑层
│   │   ├── http_client.go      # HTTP 请求执行
│   │   ├── request_service.go  # 请求管理服务
│   │   ├── collection_service.go
│   │   ├── environment_service.go
│   │   └── history_service.go
│   │
│   └── handlers/               # Wails 绑定处理器（暴露给前端）
│       ├── request_handler.go
│       ├── collection_handler.go
│       ├── environment_handler.go
│       ├── history_handler.go
│       └── app_state_handler.go
│
├── frontend/                   # 前端代码
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   │
│   │   ├── components/         # UI 组件
│   │   │   ├── request/        # 请求相关
│   │   │   │   ├── RequestPanel.vue
│   │   │   │   ├── MethodSelect.vue
│   │   │   │   ├── UrlInput.vue
│   │   │   │   ├── HeadersEditor.vue
│   │   │   │   ├── ParamsEditor.vue
│   │   │   │   └── BodyEditor.vue
│   │   │   │
│   │   │   ├── response/       # 响应相关
│   │   │   │   ├── ResponsePanel.vue
│   │   │   │   ├── ResponseHeaders.vue
│   │   │   │   └── ResponseBody.vue
│   │   │   │
│   │   │   ├── sidebar/        # 侧边栏
│   │   │   │   ├── Sidebar.vue
│   │   │   │   ├── CollectionTree.vue
│   │   │   │   └── HistoryList.vue
│   │   │   │
│   │   │   ├── tabs/           # 标签页
│   │   │   │   ├── TabBar.vue
│   │   │   │   └── TabItem.vue
│   │   │   │
│   │   │   ├── modals/         # 弹窗
│   │   │   │   ├── ConfirmModal.vue
│   │   │   │   ├── SaveRequestModal.vue
│   │   │   │   ├── EnvironmentModal.vue
│   │   │   │   └── InputModal.vue
│   │   │   │
│   │   │   └── common/         # 通用组件
│   │   │       ├── KeyValueEditor.vue
│   │   │       ├── TabGroup.vue
│   │   │       ├── IconButton.vue
│   │   │       └── Toast.vue
│   │   │
│   │   ├── stores/             # Pinia 状态管理
│   │   │   ├── tabs.ts         # Tab 状态
│   │   │   ├── request.ts      # 当前请求状态
│   │   │   ├── response.ts     # 响应缓存（仅内存）
│   │   │   ├── collection.ts   # 集合数据
│   │   │   ├── environment.ts  # 环境变量
│   │   │   ├── history.ts      # 历史记录
│   │   │   └── appState.ts     # 应用状态
│   │   │
│   │   ├── composables/        # 组合式函数
│   │   │   ├── useRequest.ts
│   │   │   ├── useKeyValue.ts
│   │   │   └── useModal.ts
│   │   │
│   │   ├── types/              # TypeScript 类型
│   │   │   └── index.ts
│   │   │
│   │   └── assets/
│   │       └── styles/
│   │           └── main.css
│   │
│   ├── wailsjs/                # Wails 自动生成
│   │   └── go/
│   │       ├── handlers/
│   │       └── models.ts
│   │
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── postcss.config.js
│
└── build/                      # 构建输出
```

## 4. 数据模型设计

### 4.1 数据库表结构

```sql
-- 集合（顶层容器）
CREATE TABLE collections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 文件夹（只能在集合下，不能嵌套）
CREATE TABLE folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
);

-- 请求
CREATE TABLE requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_id INTEGER NOT NULL,
    folder_id INTEGER,
    name TEXT NOT NULL,
    method TEXT NOT NULL DEFAULT 'GET',
    url TEXT NOT NULL DEFAULT '',
    headers TEXT DEFAULT '[]',
    params TEXT DEFAULT '[]',
    body TEXT DEFAULT '',
    body_type TEXT DEFAULT 'none',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
    FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
);

-- 历史记录（最多保留 100 条，超出自动删除最早记录）
CREATE TABLE history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id INTEGER,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    request_headers TEXT,
    request_body TEXT,
    status_code INTEGER,
    response_headers TEXT,
    response_body TEXT,
    duration_ms INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE SET NULL
);

-- 环境
CREATE TABLE environments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    variables TEXT DEFAULT '[]',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 全局变量
CREATE TABLE global_variables (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    variables TEXT DEFAULT '[]',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用状态（单行）
CREATE TABLE app_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    window_width INTEGER DEFAULT 1200,
    window_height INTEGER DEFAULT 800,
    window_x INTEGER,
    window_y INTEGER,
    window_maximized INTEGER DEFAULT 0,
    sidebar_open INTEGER DEFAULT 1,
    sidebar_width INTEGER DEFAULT 260,
    layout_direction TEXT DEFAULT 'horizontal',
    split_ratio INTEGER DEFAULT 50,
    theme TEXT DEFAULT 'system',
    active_env_id INTEGER,
    request_timeout REAL DEFAULT 30,
    auto_locate_sidebar INTEGER DEFAULT 1,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 侧边栏展开状态
CREATE TABLE sidebar_state (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_type TEXT NOT NULL,
    item_id INTEGER NOT NULL,
    expanded INTEGER DEFAULT 0,
    UNIQUE(item_type, item_id)
);

-- Tab 会话
CREATE TABLE tab_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tab_id TEXT UNIQUE NOT NULL,
    request_id INTEGER,
    title TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    is_active INTEGER DEFAULT 0,
    is_dirty INTEGER DEFAULT 0,
    method TEXT DEFAULT 'GET',
    url TEXT DEFAULT '',
    headers TEXT DEFAULT '[]',
    params TEXT DEFAULT '[]',
    body TEXT DEFAULT '',
    body_type TEXT DEFAULT 'none',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE SET NULL
);
```

### 4.2 Go 数据模型

```go
// 请求
type Request struct {
    ID           int64      `json:"id" db:"id"`
    CollectionID int64      `json:"collectionId" db:"collection_id"`
    FolderID     *int64     `json:"folderId" db:"folder_id"`
    Name         string     `json:"name" db:"name"`
    Method       string     `json:"method" db:"method"`
    URL          string     `json:"url" db:"url"`
    Headers      []KeyValue `json:"headers"`
    Params       []KeyValue `json:"params"`
    Body         string     `json:"body" db:"body"`
    BodyType     string     `json:"bodyType" db:"body_type"`
    SortOrder    int        `json:"sortOrder" db:"sort_order"`
    CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
    UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

type KeyValue struct {
    Key     string `json:"key"`
    Value   string `json:"value"`
    Enabled bool   `json:"enabled"`
}

// 响应
type Response struct {
    StatusCode int               `json:"statusCode"`
    Status     string            `json:"status"`
    Headers    map[string]string `json:"headers"`
    Body       string            `json:"body"`
    Size       int64             `json:"size"`
    Duration   int64             `json:"duration"`
}

// 集合
type Collection struct {
    ID          int64     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    SortOrder   int       `json:"sortOrder" db:"sort_order"`
    CreatedAt   time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// 文件夹
type Folder struct {
    ID           int64     `json:"id" db:"id"`
    CollectionID int64     `json:"collectionId" db:"collection_id"`
    Name         string    `json:"name" db:"name"`
    SortOrder    int       `json:"sortOrder" db:"sort_order"`
    CreatedAt    time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// 环境
type Environment struct {
    ID        int64      `json:"id" db:"id"`
    Name      string     `json:"name" db:"name"`
    Variables []Variable `json:"variables"`
    CreatedAt time.Time  `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

type Variable struct {
    Key    string `json:"key"`
    Value  string `json:"value"`
    Secret bool   `json:"secret"`
}
```

### 4.3 层级结构规则

- 集合（Collection）：顶层容器
- 文件夹（Folder）：只能在集合下，不能嵌套
- 请求（Request）：可以直接在集合下，也可以在文件夹下

最大层级：集合 → 文件夹 → 请求（3层）

## 5. UI 布局设计

### 5.0 窗口框架

应用使用无边框窗口（Frameless）配合自定义标题栏：

```go
// main.go - Wails 窗口配置
Frameless: true,
Windows: &windows.Options{
    WebviewIsTransparent: false,
    WindowIsTranslucent:  false,
},
```

自定义标题栏组件 `TitleBar.vue`：
- 可拖动区域用于移动窗口（通过 `--wails-draggable: drag`）
- 窗口控制按钮（最小化、最大化/还原、关闭）
- 集成菜单按钮、环境选择器、主题切换、设置按钮

### 5.0.1 窗口控制按钮

自定义标题栏包含三个窗口控制按钮：

| 按钮 | 图标 | 功能 |
|------|------|------|
| 最小化 | `—` | 最小化窗口 |
| 最大化/还原 | `□` / `⧉` | 最大化时显示叠加方块图标，恢复后显示单方块 |
| 关闭 | `✕` | 关闭应用 |

最大化按钮通过监听 Wails 窗口事件动态切换图标：
- `wails:window-maximised` → 显示还原图标
- `wails:window-restored` / `wails:window-unmaximised` → 显示最大化图标

### 5.1 整体布局

采用多标签页 + 可切换分栏布局：

```
┌──────────────────────────────────────────────────────────────────────┐
│  ☰ PostMe                    [ env ▼ ]  [☀/🌙] [⚙]      ─   □   ✕  │
├──────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐ ┌──────────────┐ ┌───┐                             │
│  │● GET /users ✕│ │ POST /login ✕│ │ + │                             │
│  └──────────────┘ └──────────────┘ └───┘                             │
├────────────────────────────────┬─────────────────────────────────────┤
│        请求编辑区              │            响应展示区               │
│                                │                                     │
│  默认左右分栏                  │  可切换为上下分栏                   │
└────────────────────────────────┴─────────────────────────────────────┘

标题栏按钮说明：
- ☰ : 展开/收起侧边栏
- [ env ▼ ] : 环境切换下拉
- [☀] 或 [🌙] : 主题切换（仅在亮暗之间切换，跟随系统只能在设置中选择）
- [⚙] : 设置

### 5.2.1 主题切换按钮行为

按钮图标反映当前实际显示的主题：
- 亮色主题下显示 ☀
- 暗色主题下显示 🌙
- 跟随系统时，图标随系统主题自动变化

点击在亮暗之间切换：亮色 ↔ 暗色

**注意**：跟随系统选项只能在设置弹窗中选择，标题栏按钮不会切换到跟随系统模式。

Tooltip 提示：
- 亮色时："切换到暗色模式"
- 暗色时："切换到亮色模式"
```

### 5.2 布局控制

| 操作 | 说明 |
|------|------|
| 左上角 `☰` | 展开/收起左侧集合面板 |
| 响应区 `⇅` 按钮 | 切换左右/上下分栏 |
| 拖拽分隔条 | 调整分栏比例 |
| `Ctrl+B` | 切换侧边栏 |
| `Ctrl+\` | 切换分栏方向 |

### 5.3 侧边栏（从左侧滑出）

采用 Tab 切换模式，默认显示 Collections：

```
┌─────────────────────────────────────┐
│ ┌────────────┬────────────┐      ‹  │
│ │Collections │  History   │         │
│ └────────────┴────────────┘         │
├─────────────────────────────────────┤
│  🔍 Search...                + 📁   │
├─────────────────────────────────────┤
│                                     │
│  ▼ 📁 Users API                     │
│      ▼ 📂 Authentication            │
│          POST /login                │
│          POST /logout               │
│      ▶ 📂 Management                │
│  ▶ 📁 Orders API                    │
│  ▶ 📁 Uncategorized                 │
│                                     │
└─────────────────────────────────────┘

切换到 History Tab：

┌─────────────────────────────────────┐
│ ┌────────────┬────────────┐      ‹  │
│ │ Collections│  History   │         │
│ └────────────┴────────────┘         │
├─────────────────────────────────────┤
│  🔍 Search...                   🗑  │
├─────────────────────────────────────┤
│                                     │
│  Today                              │
│    10:30  GET   /users         200  │
│    10:28  POST  /login         200  │
│    10:25  GET   /orders        404  │
│                                     │
│  Yesterday                          │
│    18:20  PUT   /users/123     200  │
│                                     │
└─────────────────────────────────────┘
```

### 5.4 History 记录限制

| 规则 | 值 |
|------|-----|
| 最大记录数 | 100 条 |
| 超出处理 | 自动删除最早的记录 |
| 清空按钮 | `🗑` 一键清空（需二次确认） |

### 5.4.1 History 记录时机

每次成功执行请求后自动保存历史记录，记录内容包括：

| 字段 | 说明 |
|------|------|
| requestId | 关联的请求ID（可选，临时请求无此字段） |
| method | HTTP 方法 |
| url | 完整请求 URL（包含替换后的变量） |
| requestHeaders | 请求头 JSON |
| requestBody | 请求体内容 |
| statusCode | 响应状态码 |
| responseHeaders | 响应头 JSON |
| responseBody | 响应体内容 |
| durationMs | 请求耗时（毫秒） |
| createdAt | 记录时间 |

### 5.5 Tab 与侧边栏联动

- 切换 Tab 时，侧边栏自动展开对应文件夹并高亮请求
- 双击侧边栏请求，打开/切换到对应 Tab
- 单击侧边栏请求，预览模式（斜体 Tab，切换时不保留）

#### 5.5.1 Tab-Sidebar 同步实现

当 `activeTab` 变化时，自动定位并高亮侧边栏中对应的请求：

```typescript
// App.vue - 监听 activeTab 变化
watch(() => tabsStore.activeTab, async (tab) => {
  if (!tab?.requestId || !appState.autoLocateSidebar) return
  
  // 1. 查找请求所在路径
  const path = collectionStore.findRequestPath(tab.requestId)
  if (!path) return
  
  // 2. 展开父级集合和文件夹
  appState.setSidebarItemExpanded('collection', path.collectionId, true)
  if (path.folderId) {
    appState.setSidebarItemExpanded('folder', path.folderId, true)
  }
  
  // 3. 高亮请求项
  appState.highlightedRequestId = path.requestId
})
```

RequestItem.vue 根据 `highlightedRequestId` 显示高亮样式。

### 5.6 Tab 状态标识

| 状态 | 显示样式 |
|------|----------|
| 新建未保存 | 斜体 `Untitled` |
| 已保存无修改 | 正常 `GET /users` |
| 已保存有修改 | ● 圆点 + 正常 `● GET /users` |

#### 5.6.1 Dirty 状态检测

Tab 的脏状态通过比较当前内容与 `originalState` 计算：

```typescript
// tabs.ts - 计算 dirty 状态
function computeDirty(tab: Tab): boolean {
  if (!tab.originalState) return true
  const orig = tab.originalState
  return (
    tab.method !== orig.method ||
    tab.url !== orig.url ||
    tab.body !== orig.body ||
    tab.bodyType !== orig.bodyType ||
    JSON.stringify(tab.headers) !== JSON.stringify(orig.headers) ||
    JSON.stringify(tab.params) !== JSON.stringify(orig.params)
  )
}
```

- `originalState` 在打开请求或保存后更新
- 切换 HTTP 方法（如 POST → GET）正确触发脏状态

### 5.7 Tab 拖放重排序

支持通过拖放重新排列标签页顺序：

```typescript
// TabBar.vue - 拖放实现
function onDragStart(e: DragEvent, index: number) {
  dragIndex.value = index
  e.dataTransfer!.effectAllowed = 'move'
}

function onDragOver(e: DragEvent, index: number) {
  e.preventDefault()
  dropIndex.value = index
}

function onDrop(e: DragEvent, index: number) {
  e.preventDefault()
  if (dragIndex.value !== null && dragIndex.value !== index) {
    tabsStore.reorderTabs(dragIndex.value, index)
  }
  resetDragState()
}
```

- 拖动时显示插入指示线
- 放置后更新 `sortOrder` 并持久化到数据库

## 6. 配色方案

### 6.1 暗色主题 (Dark Mode)

```
背景层级:
  bg-base        #1a1a1a   最底层背景
  bg-surface     #262626   卡片/面板背景
  bg-elevated    #333333   悬浮/弹窗背景
  bg-hover       #3d3d3d   悬停状态

文字颜色:
  text-primary   #f5f5f5   主要文字
  text-secondary #a3a3a3   次要文字
  text-muted     #737373   占位符/禁用

强调色:
  accent         #d97706   主强调色（按钮、链接）
  accent-hover   #b45309   强调色悬停
  accent-subtle  #78350f   强调色背景（淡）

边框:
  border         #404040   默认边框
  border-focus   #d97706   聚焦边框
```

### 6.2 亮色主题 (Light Mode)

```
背景层级:
  bg-base        #ffffff   最底层背景
  bg-surface     #fafafa   卡片/面板背景
  bg-elevated    #f5f5f5   悬浮/弹窗背景
  bg-hover       #e5e5e5   悬停状态

文字颜色:
  text-primary   #171717   主要文字
  text-secondary #525252   次要文字
  text-muted     #a3a3a3   占位符/禁用

强调色:
  accent         #d97706   主强调色
  accent-hover   #b45309   强调色悬停
  accent-subtle  #fef3c7   强调色背景（淡）

边框:
  border         #e5e5e5   默认边框
  border-focus   #d97706   聚焦边框
```

### 6.3 HTTP 方法颜色

```
GET      #22c55e   绿色
POST     #3b82f6   蓝色
PUT      #f59e0b   橙黄色
PATCH    #8b5cf6   紫色
DELETE   #ef4444   红色
OPTIONS  #6b7280   灰色
HEAD     #6b7280   灰色
```

### 6.4 状态码颜色

```
2xx 成功       #22c55e   绿色
3xx 重定向     #3b82f6   蓝色
4xx 客户端错误 #f59e0b   橙色
5xx 服务端错误 #ef4444   红色
```

### 6.5 圆角规范

```
radius-sm    4px    小按钮、标签
radius-md    8px    按钮、输入框、卡片
radius-lg    12px   面板、弹窗、Tab
```

## 7. 弹窗设计

### 7.1 规则

- 所有弹窗使用自定义组件，禁止使用系统原生弹窗（alert/confirm/prompt）
- 所有删除操作必须二次确认
- 弹窗跟随主题切换颜色

### 7.2 弹窗类型

| 类型 | 用途 |
|------|------|
| confirm | 普通确认操作 |
| danger | 危险操作（删除） |
| input | 输入内容（重命名） |
| select | 选择项目（移动到） |
| info | 信息提示 |

### 7.3 按钮颜色规则

| 场景 | 确认按钮 | 取消按钮 |
|------|----------|----------|
| 普通确认 | `#d97706` 强调色 | `#404040` 中性 |
| 危险操作 | `#ef4444` 红色 | `#404040` 中性 |

### 7.4 交互规范

| 交互 | 行为 |
|------|------|
| 点击遮罩 | 关闭弹窗（危险操作除外） |
| ESC 键 | 关闭弹窗 |
| Enter 键 | 触发确认（输入框除外） |
| 打开时 | 焦点移到取消按钮（防止误操作） |

### 7.4.1 环境管理弹窗

通过点击环境选择器右侧的齿轮按钮打开，支持完整的环境 CRUD 操作：

```
┌─────────────────────────────────────────────────────────────────┐
│  管理环境                                                    ✕   │
├─────────────────────────────────────────────────────────────────┤
│  ┌───────────────────┐  ┌─────────────────────────────────────┐ │
│  │ 环境列表          │  │ 变量编辑                            │ │
│  │ ─────────────────│  │ ───────────────────────────────────│ │
│  │ ▸ Development    │  │  Key            Value         [🗑] │ │
│  │   Production     │  │  ┌──────────┐  ┌────────────┐       │ │
│  │   Staging        │  │  │ BASE_URL │  │ http://... │       │ │
│  │                  │  │  └──────────┘  └────────────┘       │ │
│  │                  │  │  ┌──────────┐  ┌────────────┐       │ │
│  │  [+ 新建环境]     │  │  │ API_KEY  │  │ ********   │  [👁] │ │
│  └───────────────────┘  │  └──────────┘  └────────────┘       │ │
│                         │                                     │ │
│                         │  [+ 添加变量]                        │ │
│                         └─────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                                              [ 取消 ]  [ 保存 ] │
└─────────────────────────────────────────────────────────────────┘
```

功能：
- 左侧环境列表支持选择、重命名、删除
- 右侧编辑选中环境的变量
- 变量支持 Secret 模式（密码遮挡，点击眼睛图标显示）
- 新建环境自动选中并进入编辑
- 重复变量名检测（红色边框高亮，保存时 Toast 提示）
- 支持全局变量（所有环境共享）

### 7.4.2 删除确认

所有删除操作都需要二次确认：

| 对象 | 确认消息 |
|------|----------|
| 集合 | 删除集合将同时删除所有文件夹和请求 |
| 文件夹 | 删除文件夹将同时删除其中的所有请求 |
| 请求 | 确认删除请求 |
| 环境 | 删除环境将丢失所有变量 |
| 环境变量 | 仅当变量名非空时确认删除 |

### 7.5 Toast 通知

非阻断式提示，右上角显示，3秒后自动消失：

```
✓ 操作成功   #22c55e
✕ 操作失败   #ef4444
⚠ 警告提示   #f59e0b
ℹ 信息提示   #3b82f6
```

## 8. 编辑器功能

### 8.1 使用 CodeMirror 6

支持功能：
- 语法高亮（JSON、XML、HTML）
- JSON 格式化（美化/压缩）
- 撤销/重做（Ctrl+Z / Ctrl+Shift+Z）
- 查找/替换（Ctrl+F / Ctrl+H）
- 自动补全（括号、引号配对）
- 代码折叠
- 行号显示

### 8.1.1 CodeMirror 快捷键处理

CodeMirror 编辑器配置了自定义快捷键映射，避免与全局快捷键冲突：

```typescript
keymap.of([
  {
    key: 'Ctrl-Enter',  // 发送请求，阻止换行
    run: () => { emitKeyboardAction('send'); return true }
  },
  {
    key: 'Ctrl-h',      // 打开搜索替换面板
    run: (view) => { openSearchPanel(view); return true }
  },
  {
    key: 'Ctrl-Shift-f', // 打开搜索面板
    run: (view) => { openSearchPanel(view); return true }
  },
])
```

搜索扩展通过 `search({ top: true })` 配置在顶部显示搜索面板。

### 8.1.2 亮色模式语法高亮

亮色模式下 CodeMirror 使用高饱和度颜色，以适配 WebView2 渲染特性：

| 元素 | 颜色 | Hex |
|------|------|-----|
| 属性名 | 绿色 | #22c55e |
| 字符串 | 橙色 | #f97316 |
| 数字 | 蓝色 | #0d6efd |
| 布尔/null | 黄色 | #eab308 |

**注意**：WebView2 子像素渲染在模态框遮罩下会对彩色文本产生轻微光晕，绿色最不明显，橙色/黄色略有蓝色光晕但可接受。

### 8.2 环境变量支持

- 使用 `{{变量名}}` 语法
- 编辑器中变量高亮显示
- 变量存在且有值 → 蓝色
- 变量未定义 → 红色 + 波浪下划线警告
- 鼠标悬停显示当前值

### 8.3 变量类型

| 类型 | 作用域 | 说明 |
|------|--------|------|
| 环境变量 | 当前环境 | 切换环境时值改变 |
| 全局变量 | 所有环境 | 跨环境共享 |
| 临时变量 | 当前会话 | 从响应中提取，关闭后丢失 |

## 9. 快捷键

| 快捷键 | 作用域 | 功能 |
|--------|--------|------|
| `Ctrl+S` | 全局 | 保存当前请求 |
| `Ctrl+Enter` | 全局 | 发送请求 |
| `Ctrl+Z` | 编辑器 | 撤销 |
| `Ctrl+Shift+Z` | 编辑器 | 重做 |
| `Ctrl+F` | 编辑器 | 查找 |
| `Ctrl+H` | 编辑器 | 替换（仅请求体） |
| `Ctrl+Shift+F` | 编辑器 | 格式化 JSON |
| `Ctrl+B` | 全局 | 切换侧边栏 |
| `Ctrl+\` | 全局 | 切换分栏方向 |
| `Ctrl+T` | 全局 | 新建标签页 |
| `Ctrl+W` | 全局 | 关闭当前标签页 |

### 9.1 快捷键实现架构

全局快捷键通过事件总线模式实现，避免组件间的紧耦合：

```typescript
// useKeyboardActions.ts - 事件发射器
const listeners = new Map<string, Function[]>()

export function onKeyboardAction(action: string, callback: Function) {
  if (!listeners.has(action)) listeners.set(action, [])
  listeners.get(action)!.push(callback)
  return () => { /* cleanup */ }
}

export function emitKeyboardAction(action: string) {
  listeners.get(action)?.forEach(cb => cb())
}
```

- App.vue 监听全局键盘事件并调用 `emitKeyboardAction`
- RequestPanel.vue 通过 `onKeyboardAction` 订阅并执行保存/发送

## 10. 状态保存

### 10.1 存储位置

| 状态 | 存储位置 | 生命周期 |
|------|----------|----------|
| 窗口大小/位置 | SQLite | 永久 |
| 布局设置 | SQLite | 永久 |
| 侧边栏展开状态 | SQLite | 永久 |
| Tab 列表和请求内容 | SQLite | 永久 |
| 响应缓存 | 内存 (Pinia) | 应用运行期间 |

### 10.1.1 窗口位置验证

恢复窗口位置时进行有效性检查，避免窗口显示在不可见区域：

| 验证项 | 有效范围 | 无效时处理 |
|--------|----------|------------|
| 坐标 X/Y | -100 ~ 10000 | 使用系统默认位置 |
| 窗口尺寸 | ≥ 800×600 | 不保存（保留上次有效值） |
| 最小化状态 | - | 不保存位置和尺寸 |

### 10.2 响应缓存规则

- 仅在应用运行期间保留
- 关闭 Tab 时清除对应响应
- 关闭应用时清空所有响应
- 重新打开应用后响应区显示空状态

### 10.3 状态保存时机

| 事件 | 保存内容 |
|------|----------|
| 窗口大小改变 | 窗口尺寸（防抖 500ms） |
| 展开/收起集合 | 侧边栏状态 |
| 切换 Tab | 当前激活 Tab |
| 编辑请求 | Tab 内容（防抖 1s） |
| 新建/关闭 Tab | Tab 列表 |
| 关闭软件 | 全量保存 |

### 10.4 启动恢复流程

1. 读取 app_state，恢复窗口大小/位置
2. 读取 sidebar_state，恢复展开/收起状态
3. 读取 tab_sessions，恢复所有 Tab（包括未保存的修改）
4. 激活上次的 Tab
5. 响应区显示空状态（需重新发送请求）

## 11. 新建请求保存交互

### 11.1 状态流转

1. 点击 `+` 新建空白请求
2. Tab 显示斜体 `Untitled`，● 表示未保存
3. 用户编辑请求内容
4. `Ctrl+S` 或点击保存图标
5. 弹出保存弹窗，选择名称和位置

### 11.2 保存弹窗

- 请求名称输入框（默认根据 URL 生成，如 `GET /users`）
- 保存位置选择（集合/文件夹树形选择）
- 支持新建集合

### 11.3 关闭未保存 Tab

弹窗询问：是否保存对 "xxx" 的更改？
- 不保存：直接关闭
- 取消：取消关闭操作
- 保存：保存后关闭

## 12. 右键菜单

### 12.1 集合/文件夹右键菜单

- 新建请求
- 新建文件夹（仅集合）
- 重命名
- 复制
- 移动到...
- 导出
- 删除

### 12.2 请求右键菜单

- 打开
- 在新标签页打开
- 重命名
- 复制
- 移动到...
- 复制为 cURL
- 删除

## 12.3 拖放功能

支持通过拖放重新组织集合、文件夹和请求：

### 12.3.1 可拖放的元素

| 元素类型 | 可拖放 | 可作为放置目标 |
|----------|--------|----------------|
| 集合 | ✓ | ✓（重排序） |
| 文件夹 | ✓ | ✓（重排序、接收请求） |
| 请求 | ✓ | ✓（重排序） |

### 12.3.2 拖放规则

| 操作 | 行为 |
|------|------|
| 拖动集合到集合 | 重新排序集合顺序 |
| 拖动文件夹到文件夹 | 在同一集合内重新排序文件夹 |
| 拖动请求到集合 | 移动请求到集合根目录 |
| 拖动请求到文件夹 | 移动请求到指定文件夹 |
| 拖动请求到请求 | 在同一容器内重新排序请求 |

### 12.3.3 视觉反馈

- 拖动时元素半透明
- 放置目标高亮显示（环形边框）
- 插入位置显示指示线（上/下边框）

#### 拖放位置计算

请求拖放时根据光标位置动态计算插入位置（上方/下方）：

```typescript
function onDragOver(e: DragEvent, targetType: string, targetId: number) {
  e.preventDefault()
  
  // 计算光标在目标元素中的相对位置
  if (targetType === 'request') {
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
    const y = e.clientY - rect.top
    const position = y < rect.height / 2 ? 'before' : 'after'
    dropTarget.value = { type: targetType, id: targetId, position }
  }
}
```

这确保向下移动请求时能正确插入到目标下方。

### 12.3.4 跨容器移动

请求可以跨集合/文件夹移动：
1. 拖动请求到目标位置
2. 自动更新请求的 `collectionId` 和 `folderId`
3. 自动更新排序顺序

文件夹可以跨集合移动：
1. 拖动文件夹到目标集合
2. 自动更新文件夹的 `collectionId`
3. 文件夹内的请求一同移动

## 13. 请求体类型支持

### 13.1 支持的 Body 类型

| 类型 | Content-Type | 说明 |
|------|--------------|------|
| none | - | 无请求体 |
| json | application/json | JSON 格式 |
| xml | application/xml | XML 格式 |
| text | text/plain | 纯文本 |
| form-data | multipart/form-data | 表单数据，支持键值对编辑 |

### 13.2 自动请求头

切换 Body 类型时自动添加/更新 Content-Type 请求头：
- 如果已存在 Content-Type 头，更新其值
- 如果不存在，添加新的 Content-Type 头
- 选择 none 时移除 Content-Type 头

Content-Length 由 Go http.Client 自动处理，无需手动设置。

### 13.3 Form-Data 编辑

Form-Data 使用键值对编辑器，数据以 JSON 数组格式存储在 body 字段中：

```json
[
  {"key": "username", "value": "admin", "enabled": true},
  {"key": "password", "value": "123456", "enabled": true}
]
```

发送请求时自动转换为 multipart/form-data 格式。

## 14. 系统代理支持

应用自动支持系统代理设置：

| 环境变量 | 用途 |
|----------|------|
| HTTP_PROXY | HTTP 请求代理 |
| HTTPS_PROXY | HTTPS 请求代理 |
| NO_PROXY | 不使用代理的主机列表 |

Go 的 `http.DefaultTransport` 会自动读取这些环境变量并应用代理设置。

## 15. 数据存储位置

### 15.1 SQLite 数据库

数据库文件存放在应用同级目录的 data 文件夹：

```
<exe_dir>/data/postme.db
```

### 15.2 WebView2 用户数据

WebView2 的缓存、Cookie 等数据存放在用户配置目录：

```
Windows: %APPDATA%/postme/
```

通过 `WebviewUserDataPath` 配置，避免默认使用 `postme.exe` 作为文件夹名。

## 16. 请求控制

### 16.1 请求取消

- Send 按钮在请求进行中变为 Cancel 按钮（红色）
- 点击 Cancel 或按 ESC 取消当前请求
- 取消后响应区显示 "请求已取消"
- Go 后端使用 context.WithCancel 实现

### 16.2 请求超时

- 默认超时时间 30 秒
- 直接输入数值，单位秒，允许小数（如 1.5）
- 不能输入负数，最小值为 0
- 0 表示不限制超时
- 超时后响应区显示 "请求超时 (Xs)"
- Go 后端使用 context.WithTimeout 实现

### 16.3 响应区状态

| 状态 | 图标 | 文字 |
|------|------|------|
| 空闲 | 📤 | 点击 Send 发送请求 |
| 请求中 | ● (转圈) | 正在发送请求... |
| 已取消 | ⊘ | 请求已取消 |
| 超时 | ⏱ | 请求超时 (30s) |
| 错误 | ✕ | 连接失败 / 错误信息 |

## 17. 设置功能

### 17.1 设置入口

标题栏右侧添加设置按钮 `⚙`

### 17.2 设置项

| 分类 | 设置项 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| 请求 | 请求超时时间 | 数字输入框 | 30 | 单位秒，允许小数，0表示不限制 |
| 界面 | 主题 | 下拉 | 跟随系统 | 亮色/暗色/跟随系统 |
| 界面 | 切换Tab自动定位 | 开关 | 开启 | - |

### 17.3 设置弹窗

```
┌─────────────────────────────────────────────────────────────────┐
│  设置                                                       ✕   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  请求                                                           │
│  ───────────────────────────────────────────────────────────    │
│                                                                 │
│  请求超时时间                                                   │
│  ┌───────────────────────────────┐                              │
│  │  30                           │  秒                          │
│  └───────────────────────────────┘                              │
│  输入 0 表示不限制超时                                          │
│                                                                 │
│  界面                                                           │
│  ───────────────────────────────────────────────────────────    │
│                                                                 │
│  主题                                                           │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  跟随系统                                            ▼  │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  切换 Tab 时自动定位侧边栏                                      │
│  ┌────┐                                                         │
│  │ ✓  │  开启                                                   │
│  └────┘                                                         │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                                              [ 取消 ]  [ 保存 ] │
└─────────────────────────────────────────────────────────────────┘
```

### 17.4 数据模型

```sql
-- app_state 表添加字段
ALTER TABLE app_state ADD COLUMN request_timeout REAL DEFAULT 30;
ALTER TABLE app_state ADD COLUMN auto_locate_sidebar INTEGER DEFAULT 1;
```

## 18. 配置常量

```go
const (
    // History 最大记录数
    MaxHistoryRecords = 100

    // 状态保存防抖时间
    WindowResizeDebounce = 500  // ms
    ContentEditDebounce  = 1000 // ms

    // 默认窗口尺寸
    DefaultWindowWidth  = 1200
    DefaultWindowHeight = 800

    // 默认侧边栏宽度
    DefaultSidebarWidth = 260

    // 默认分栏比例
    DefaultSplitRatio = 50 // 百分比

    // 请求超时
    DefaultRequestTimeout = 30.0 // 秒，0 表示不限制，支持小数
)
```

## 19. Tab 会话持久化

应用重启时自动恢复所有打开的标签页（不包括响应数据）：

### 19.1 保存时机

- Tab 内容变更时自动保存（防抖 500ms）
- Tab 列表变更时自动保存
- 切换活动 Tab 时保存激活状态

### 19.2 恢复内容

| 内容 | 保存 | 说明 |
|------|------|------|
| Tab ID | ✓ | 唯一标识 |
| 请求 ID | ✓ | 关联的保存请求 |
| 标题 | ✓ | Tab 显示名称 |
| 排序顺序 | ✓ | Tab 位置 |
| 激活状态 | ✓ | 最后激活的 Tab |
| 脏状态 | ✓ | 是否有未保存修改 |
| 请求内容 | ✓ | method, url, headers, params, body, bodyType |
| 响应数据 | ✗ | 仅内存保留 |

## 20. Tab 交互增强

### 20.1 双击固定 Tab

- 单击侧栏请求：预览模式打开（斜体）
- 双击侧栏请求：固定模式打开（常规样式）
- 对于已打开的 Tab，双击切换并固定

### 20.2 关闭未保存确认

关闭有未保存更改的 Tab 时显示确认对话框：
- 通过关闭按钮关闭
- 通过 Ctrl+W 快捷键关闭

### 20.3 Tab 同步

- 侧栏重命名请求：Tab 标题实时更新
- 侧栏删除请求：自动关闭对应 Tab
- 切换 Tab：侧栏自动展开并高亮（可在设置中关闭）

### 20.4 活动 Tab 高亮

活动 Tab 底部显示强调色指示条（2px），在亮暗主题下都有明显的视觉区分。

## 21. 系统代理支持（增强）

### 21.1 Windows 系统代理读取

应用可以读取 Windows 系统代理设置（Internet 选项中配置的代理）：

```go
// 从 Windows 注册表读取代理设置
key, _ := registry.OpenKey(registry.CURRENT_USER,
    `Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
    registry.QUERY_VALUE)

proxyEnable, _, _ := key.GetIntegerValue("ProxyEnable")
proxyServer, _, _ := key.GetStringValue("ProxyServer")
```

### 21.2 代理设置项

在设置弹窗中添加代理开关：

| 设置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| 使用系统代理 | 开关 | 开启 | 是否使用 Windows 系统代理 |

### 21.3 代理优先级

1. Windows 系统代理（Internet 选项）
2. 环境变量（HTTP_PROXY, HTTPS_PROXY）
3. 无代理

## 22. Cloudflare 兼容性

为避免请求被 Cloudflare 等 WAF 拦截，HTTP 客户端进行了以下优化：

### 22.1 TLS 配置

- 使用现代 TLS 版本（1.2-1.3）
- 配置浏览器常用密码套件顺序
- 设置 EC 曲线优先级

### 22.2 默认请求头

未设置时自动添加以下请求头：

| 请求头 | 默认值 |
|--------|--------|
| User-Agent | Chrome 120 用户代理字符串 |
| Accept | `*/*` |
| Accept-Language | `en-US,en;q=0.9` |

### 22.3 HTTP/2 支持

启用 HTTP/2 协议支持，提高与现代服务器的兼容性。

## 23. 右键菜单增强

### 23.1 菜单排他性

同时只能显示一个右键菜单，打开新菜单时自动关闭其他菜单。

### 23.2 请求菜单项

| 菜单项 | 功能 |
|--------|------|
| 在新标签页打开 | 打开请求到新 Tab |
| 复制 | 复制请求（名称加 "(copy)" 后缀） |
| 重命名 | 修改请求名称 |
| 删除 | 删除请求（需确认） |

### 23.3 复制请求功能

复制后的请求：
- 名称后添加 "(copy)" 后缀
- 保留所有请求设置（method, url, headers, params, body）
- 自动打开新 Tab 编辑

## 24. 表单自动提交

KeyValueEditor 组件（用于请求头、请求参数、表单数据）支持自动添加行：

- 在新行输入内容后，焦点离开时自动添加
- Tab 在同一行的字段间切换不会触发添加
- 无需手动点击 "+" 按钮

## 25. 数据模型更新

### 25.1 app_state 表新增字段

```sql
ALTER TABLE app_state ADD COLUMN use_system_proxy INTEGER DEFAULT 1;
```

### 25.2 AppState 模型更新

```go
type AppState struct {
    // ... 其他字段
    UseSystemProxy bool `json:"useSystemProxy" db:"use_system_proxy"`
}
```

## 26. 双击交互增强

### 26.1 侧边栏请求双击

双击侧边栏中的请求项将固定该 Tab（从预览模式转为永久模式）。

实现机制：
- 延迟单击处理（200ms）以区分单击和双击
- 双击时取消单击定时器
- 双击时将预览 Tab 转为永久 Tab

### 26.2 标题栏双击

双击标题栏可最大化/还原窗口，与 Windows 原生行为一致。

- 双击拖拽区域切换最大化状态
- 按钮区域（`wails-no-drag`）不触发

### 26.3 最大化按钮图标

窗口最大化时，最大化按钮显示还原图标（两个重叠方框）：

| 状态 | 图标 |
|------|------|
| 正常窗口 | `StopIcon`（单个方框） |
| 最大化 | `Square2StackIcon`（两个重叠方框） |

通过 Wails 窗口事件监听状态变化：
- `wails:window-maximised`
- `wails:window-restored`
- `wails:window-unmaximised`

## 27. 系统代理完整实现

### 27.1 初始化时应用

应用启动时根据设置自动应用代理配置：

```typescript
// App.vue loadData()
const state = await api.getAppState()
await api.setUseSystemProxy(state.useSystemProxy)
```

### 27.2 设置更改时应用

设置保存时同步更新 HTTP 客户端：

```typescript
// SettingsModal.vue save()
await api.setUseSystemProxy(localSettings.useSystemProxy)
```

### 27.3 后端 API

```go
// RequestHandler
func (h *RequestHandler) SetUseSystemProxy(useProxy bool) {
    if h.httpClient != nil {
        h.httpClient.SetUseSystemProxy(useProxy)
    }
}
```

## 28. uTLS 突破 Cloudflare 检测

### 28.1 TLS 指纹伪装

使用 `github.com/refraction-networking/utls` 库模拟 Chrome 120 的 TLS 指纹：

```go
import utls "github.com/refraction-networking/utls"

tlsConn := utls.UClient(conn, &utls.Config{
    ServerName: hostname,
}, utls.HelloChrome_120)
```

### 28.2 HTTP/2 支持

根据 ALPN 协商结果选择协议：

```go
alpn := tlsConn.ConnectionState().NegotiatedProtocol
if alpn == "h2" {
    // 使用 HTTP/2
    h2Transport := &http2.Transport{}
    h2Conn, _ := h2Transport.NewClientConn(tlsConn)
    return h2Conn.RoundTrip(req)
}
// 否则使用 HTTP/1.1
```

### 28.3 增强的请求头

添加更多浏览器特征请求头：

| 请求头 | 值 |
|--------|-----|
| Accept-Encoding | `gzip, deflate, br` |
| Sec-Fetch-Dest | `empty` |
| Sec-Fetch-Mode | `cors` |
| Sec-Fetch-Site | `cross-site` |

## 29. 请求体类型扩展

### 29.1 支持的请求体类型

| 类型 | 标签 | Content-Type | 用途 |
|------|------|--------------|------|
| none | None | - | 无请求体 |
| form-data | Form Data | multipart/form-data | 文件上传、二进制数据 |
| x-www-form-urlencoded | URL Encoded | application/x-www-form-urlencoded | 简单文本表单 |
| binary | Binary | application/octet-stream | 二进制文件上传 |
| json | JSON | application/json | JSON 数据 |
| xml | XML | application/xml | XML 数据 |
| text | Text | text/plain | 纯文本 |

### 29.2 x-www-form-urlencoded 编码

```go
if req.BodyType == "x-www-form-urlencoded" && req.Body != "" {
    var formItems []models.KeyValue
    json.Unmarshal([]byte(req.Body), &formItems)
    formData := url.Values{}
    for _, item := range formItems {
        if item.Enabled && item.Key != "" {
            formData.Add(item.Key, item.Value)
        }
    }
    bodyReader = strings.NewReader(formData.Encode())
    contentType = "application/x-www-form-urlencoded"
}
```

### 29.3 二进制文件上传

前端使用 Wails 文件对话框选择文件：

```typescript
async function selectBinaryFile() {
    const filePath = await window.runtime.OpenFileDialog({
        Title: 'Select File',
        Filters: [{ DisplayName: 'All Files (*.*)', Pattern: '*.*' }],
    })
    if (filePath) {
        emit('update:body', filePath)  // body 存储文件路径
    }
}
```

后端读取并发送文件：

```go
if req.BodyType == "binary" && req.Body != "" {
    file, err := os.Open(req.Body)  // body 是文件路径
    if err != nil {
        return nil, err
    }
    defer file.Close()
    fileContent, _ := io.ReadAll(file)
    bodyReader = bytes.NewReader(fileContent)
    contentType = "application/octet-stream"
}
```
