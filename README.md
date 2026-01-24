# PostMe

<div align="center">

一款基于 Wails 构建的轻量级 REST API 客户端，类似 Postman/Insomnia，但更简洁高效。

</div>

## ✨ 特性

- **多标签页** - 同时打开多个请求，支持拖放排序
- **集合管理** - 使用集合和文件夹组织请求
- **环境变量** - 支持 `{{变量名}}` 语法，轻松切换开发/生产环境
- **历史记录** - 自动保存最近 100 条请求记录
- **代码编辑器** - 基于 CodeMirror 6，支持语法高亮、代码折叠、查找替换
- **多种请求体** - JSON、XML、Form Data、x-www-form-urlencoded、Binary
- **系统代理** - 自动读取 Windows 系统代理设置
- **状态持久化** - 窗口状态、标签页、未保存修改全部自动恢复
- **亮暗主题** - 支持亮色、暗色、跟随系统三种模式

## 📦 技术栈

| 层级 | 技术 |
|------|------|
| 框架 | [Wails](https://wails.io/) v2 |
| 后端 | Go 1.25+ |
| 前端 | Vue 3 + TypeScript + Pinia |
| UI | Tailwind CSS + Headless UI |
| 编辑器 | CodeMirror 6 |
| 数据库 | SQLite (modernc.org/sqlite) |

## 🚀 快速开始

### 环境要求

- **Go**: 1.21 或更高版本
- **Node.js**: 18 或更高版本
- **Wails CLI**: v2

### 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

验证安装：

```bash
wails doctor
```

确保所有依赖项都显示 ✓。

### 克隆项目

```bash
git clone https://github.com/SoulTraitor/postme.git
cd postme
```

### 开发模式

```bash
wails dev
```

这将启动热重载开发服务器，前后端修改都会自动刷新。

### 生产构建

```bash
wails build
```

构建产物位于 `build/bin/postme.exe`。

## 🔧 可选配置

### UPX 压缩

[UPX](https://upx.github.io/) 是一个可执行文件压缩工具，可以显著减小最终产物体积（通常压缩 50-70%）。

#### 安装 UPX

**Windows (Scoop):**
```bash
scoop install upx
```

**Windows (Chocolatey):**
```bash
choco install upx
```

**Windows (手动安装):**
1. 从 [UPX Releases](https://github.com/upx/upx/releases) 下载最新版本
2. 解压并将 `upx.exe` 所在目录添加到系统 PATH

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install upx

# Arch
sudo pacman -S upx
```

**macOS:**
```bash
brew install upx
```

#### 使用 UPX 构建

```bash
wails build -upx
```

或指定压缩级别（1-9，默认 9，越高压缩率越大但更慢）：

```bash
wails build -upx -upxflags="-9"
```

### NSIS 安装包

Wails 支持使用 NSIS 创建 Windows 安装程序：

1. 安装 [NSIS](https://nsis.sourceforge.io/Download)
2. 构建安装包：

```bash
wails build -nsis
```

## 📁 项目结构

```
postme/
├── main.go                 # 应用入口
├── wails.json              # Wails 配置
├── internal/
│   ├── models/             # 数据模型
│   ├── database/           # SQLite 数据库层
│   ├── services/           # 业务逻辑
│   └── handlers/           # Wails 绑定（暴露给前端）
├── frontend/
│   ├── src/
│   │   ├── components/     # Vue 组件
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── composables/    # 组合式函数
│   │   └── types/          # TypeScript 类型
│   └── wailsjs/            # Wails 自动生成的绑定
└── build/                  # 构建输出
```

## ⌨️ 快捷键

| 快捷键 | 功能 |
|--------|------|
| `Ctrl+Enter` | 发送请求 |
| `Ctrl+S` | 保存请求 |
| `Ctrl+T` | 新建标签页 |
| `Ctrl+W` | 关闭当前标签页 |
| `Ctrl+B` | 切换侧边栏 |
| `Ctrl+\` | 切换分栏方向 |
| `Ctrl+F` | 查找（编辑器内） |
| `Ctrl+H` | 替换（编辑器内） |
| `Ctrl+Shift+F` | 格式化 JSON |

## 📄 License

MIT
