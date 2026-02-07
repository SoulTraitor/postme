# PostMe - REST 请求工具技术设计文档

## 1. 项目概述

PostMe 是一款基于 Wails 构建的轻量级 REST 请求发送工具，类似于 Postman/Insomnia，但更简洁。

## 2. 技术栈

| 层级 | 技术选型 | 版本 |
|------|----------|------|
| 框架 | Wails | v2 |
| 后端 | Go | 1.21+ |
| 前端 | Vue 3 + TypeScript + Composition API | Vue 3.4+ |
| UI 样式 | Tailwind CSS | 3.x |
| UI 组件 | Headless UI | 1.x |
| 图标 | Heroicons | 2.x |
| 状态管理 | Pinia | 2.x |
| 代码编辑器 | CodeMirror | 6.x |
| 数据库 | SQLite (modernc.org/sqlite) | - |
| 字体 | JetBrains Mono | - |

## 3. 项目目录结构

```
postme/
├── main.go                     # 应用入口
├── app.go                      # 主应用结构体
├── wails.json                  # Wails 配置
├── data/postme.db              # SQLite 数据库
├── internal/
│   ├── models/                 # 数据模型
│   ├── database/               # 数据库层
│   │   └── repository/         # 数据访问层
│   ├── services/               # 业务逻辑层
│   └── handlers/               # Wails 绑定处理器
├── frontend/
│   ├── src/
│   │   ├── components/         # UI 组件
│   │   ├── stores/             # Pinia 状态管理
│   │   ├── composables/        # 组合式函数
│   │   └── types/              # TypeScript 类型
│   └── public/favicon.ico      # 应用图标
└── build/                      # 构建输出
```

## 4. 数据模型

### 4.1 核心表结构

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
    collection_id INTEGER NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 请求
CREATE TABLE requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_id INTEGER NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    folder_id INTEGER REFERENCES folders(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    method TEXT NOT NULL DEFAULT 'GET',
    url TEXT NOT NULL DEFAULT '',
    headers TEXT DEFAULT '[]',
    params TEXT DEFAULT '[]',
    body TEXT DEFAULT '',
    body_type TEXT DEFAULT 'none',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 历史记录（最多 100 条）
CREATE TABLE history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id INTEGER REFERENCES requests(id) ON DELETE SET NULL,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    request_headers TEXT,
    request_body TEXT,
    status_code INTEGER,
    response_headers TEXT,
    response_body TEXT,
    duration_ms INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 环境变量
CREATE TABLE environments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    variables TEXT DEFAULT '[]',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
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
    use_system_proxy INTEGER DEFAULT 1,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tab 会话持久化
CREATE TABLE tab_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tab_id TEXT UNIQUE NOT NULL,
    request_id INTEGER REFERENCES requests(id) ON DELETE SET NULL,
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
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 4.2 层级结构

- **集合（Collection）**：顶层容器
- **文件夹（Folder）**：只能在集合下，不能嵌套
- **请求（Request）**：可以直接在集合下，也可以在文件夹下

最大层级：集合 → 文件夹 → 请求（3层）

## 5. UI 布局

### 5.1 窗口结构

无边框窗口 + 自定义标题栏：

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
│  默认左右分栏，可切换上下分栏  │  可拖动分隔线调整比例               │
└────────────────────────────────┴─────────────────────────────────────┘
```

### 5.2 标题栏

| 元素 | 功能 |
|------|------|
| ☰ | 展开/收起侧边栏 |
| [ env ▼ ] | 环境切换下拉 |
| ☀/🌙 | 主题切换（亮色 ↔ 暗色） |
| ⚙ | 设置 |
| □/⧉ | 最大化/还原（图标随状态切换） |

### 5.3 侧边栏

Tab 切换模式：Collections / History

- **Collections**：树形结构展示集合、文件夹、请求
- **History**：按时间分组显示历史记录（最多 100 条）

### 5.4 Tab 状态

| 状态 | 显示样式 |
|------|----------|
| 新建未保存 | 斜体 `Untitled` |
| 已保存无修改 | 正常显示 |
| 已保存有修改 | ● 圆点（脉冲动画） |
| 预览模式 | 斜体 + 半透明 |
| 活动 Tab | 底部橙色指示条 + 阴影提升 |

### 5.5 分隔线交互

- 拖动调整分栏比例（20%-80%）
- 拖动时显示橙色发光效果和比例提示
- 双击重置为 50%

## 6. 视觉设计

### 6.1 配色方案

**暗色主题**：
```
背景: #1a1a1a → #262626 → #333333
文字: #f5f5f5 / #a3a3a3 / #737373
强调: #d97706 (橙色)
边框: #404040
```

**亮色主题**：
```
背景: #ffffff → #fafafa → #f5f5f5
文字: #171717 / #525252 / #a3a3a3
强调: #d97706 (橙色)
边框: #e5e5e5
```

### 6.2 HTTP 方法颜色

| 方法 | 颜色 |
|------|------|
| GET | #22c55e 绿色 |
| POST | #3b82f6 蓝色 |
| PUT | #f59e0b 橙黄色 |
| PATCH | #8b5cf6 紫色 |
| DELETE | #ef4444 红色 |
| OPTIONS/HEAD | #6b7280 灰色 |

### 6.3 状态码徽章

圆角徽章样式，带图标：

| 范围 | 颜色 | 图标 |
|------|------|------|
| 2xx | 绿色 | CheckCircleIcon |
| 3xx | 蓝色 | - |
| 4xx | 黄色 | XCircleIcon |
| 5xx | 红色 | XCircleIcon |

### 6.4 字体

**JetBrains Mono**（等宽字体）应用于：
- CodeMirror 编辑器（13px，禁用连字）
- URL 输入框
- 键值对编辑器
- 响应头值
- 状态栏数值

### 6.5 动效规范

| 元素 | 动效 |
|------|------|
| Send 按钮 | 悬停上移 2px + 阴影增强 |
| Toast 通知 | 右侧滑入 + 缩放（300ms） |
| 模态框 | 缩放 90%→100% + 背景模糊（300ms） |
| Tab 切换 | 平滑过渡（200ms） |
| 成功图标 | 单次弹跳动画 |

## 7. 组件设计

### 7.1 弹窗类型

| 类型 | 确认按钮颜色 | 用途 |
|------|-------------|------|
| confirm | 橙色 #d97706 | 普通确认 |
| danger | 红色 #ef4444 | 删除操作 |
| input | 橙色 | 输入内容（重命名） |
| select | 橙色 | 选择项目（移动到） |

**交互规范**：
- ESC 键关闭弹窗
- 点击遮罩关闭（危险操作除外）
- 嵌套弹窗时，ESC 只关闭最上层

### 7.2 Toast 通知

右上角显示，3 秒后自动消失：

| 类型 | 颜色 | 触发场景 |
|------|------|---------|
| success | 绿色 | 保存成功、删除成功、复制成功 |
| error | 红色 | 请求失败、保存失败 |
| warning | 橙色 | 变量名重复 |
| info | 蓝色 | 信息提示 |

### 7.3 右键菜单

**集合/文件夹**：新建请求、新建文件夹、重命名、复制、移动到、删除

**请求**：在新标签页打开、复制、重命名、删除

### 7.4 键值对编辑器

- 行悬停高亮
- 禁用项半透明
- 删除按钮点击缩放反馈
- 新增行淡入动画
- 文本溢出时悬停显示完整内容浮层

### 7.5 空状态设计

响应面板空状态：
- 渐变光晕背景
- 大图标（16x16）
- 标题 + 描述 + 快捷键提示

## 8. 编辑器功能

### 8.1 CodeMirror 配置

- 语法高亮（JSON、XML、HTML）
- JSON 格式化
- 撤销/重做
- 查找/替换
- 自动补全（括号、引号配对）
- 代码折叠
- 行号显示
- 自动折行（避免水平滚动）

### 8.2 环境变量

语法：`{{变量名}}`

| 状态 | 显示 |
|------|------|
| 变量存在 | 蓝色高亮 |
| 变量未定义 | 红色 + 波浪下划线 |

## 9. 快捷键

| 快捷键 | 功能 |
|--------|------|
| `Ctrl+Enter` | 发送请求 |
| `Ctrl+S` | 保存请求 |
| `Ctrl+T` | 新建标签页 |
| `Ctrl+W` | 关闭标签页 |
| `Ctrl+Tab` | 切换到下一个标签页 |
| `Ctrl+Shift+Tab` | 切换到上一个标签页 |
| `Ctrl+B` | 切换侧边栏 |
| `Ctrl+\` | 切换分栏方向 |
| `Ctrl+F` | 查找 |
| `Ctrl+H` | 替换 |
| `Ctrl+Shift+F` | 格式化 JSON |

## 10. 请求体类型

| 类型 | Content-Type |
|------|--------------|
| none | - |
| json | application/json |
| xml | application/xml |
| text | text/plain |
| form-data | multipart/form-data |
| x-www-form-urlencoded | application/x-www-form-urlencoded |
| binary | application/octet-stream |

## 11. 网络配置

### 11.1 系统代理

自动读取 Windows 系统代理设置，可在设置中开关。

### 11.2 TLS 指纹

使用 uTLS 模拟 Chrome 120 指纹，避免被 Cloudflare 等 WAF 拦截。

### 11.3 默认请求头

| 请求头 | 默认值 |
|--------|--------|
| User-Agent | Chrome 120 |
| Accept | `*/*` |
| Accept-Language | `en-US,en;q=0.9` |
| Accept-Encoding | `gzip, deflate, br` |

## 12. 状态持久化

### 12.1 存储策略

| 状态 | 存储位置 | 生命周期 |
|------|----------|----------|
| 窗口位置/大小 | SQLite | 永久 |
| 布局设置 | SQLite | 永久 |
| Tab 列表和内容 | SQLite | 永久 |
| 侧边栏展开状态 | SQLite | 永久 |
| 响应数据 | 内存 | 应用运行期间 |

### 12.2 保存时机

- 窗口大小改变（防抖 500ms）
- Tab 内容编辑（防抖 500ms）
- 展开/收起集合
- 切换/关闭 Tab

### 12.3 启动优化

采用优先级加载策略：
1. **Priority 1**（~50ms）：appState + tabs → 立即显示 UI
2. **Priority 2**（~500ms）：collections → 侧边栏显示
3. **Priority 3**（后台）：environments、history

## 13. 配置常量

```go
const (
    MaxHistoryRecords     = 100   // 历史记录最大数
    WindowResizeDebounce  = 500   // ms
    ContentEditDebounce   = 500   // ms
    DefaultWindowWidth    = 1200
    DefaultWindowHeight   = 800
    DefaultSidebarWidth   = 260
    DefaultSplitRatio     = 50    // %
    DefaultRequestTimeout = 30.0  // 秒
)
```

## 14. 设置项

| 分类 | 设置项 | 默认值 |
|------|--------|--------|
| 请求 | 超时时间 | 30 秒 |
| 界面 | 主题 | 跟随系统 |
| 界面 | 切换Tab自动定位 | 开启 |
| 网络 | 使用系统代理 | 开启 |

---

*文档最后更新：2026-02-02*
