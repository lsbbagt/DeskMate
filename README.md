# DeskMate - 桌面伙伴

<div align="center">

![DeskMate Logo](build/appicon.png)

**简洁高效的桌面效率工具**

[![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Wails Version](https://img.shields.io/badge/Wails-2.11.0-00ADD8?style=flat)](https://wails.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

## 📖 项目简介

DeskMate 是一款为**学生、开发者和知识工作者**设计的桌面效率工具，提供待办事项管理、代码模板收藏、Markdown编辑等功能，帮助你更好地组织和管理工作流程。

### ✨ 核心特性

- 📝 **待办事项** - 管理任务，设置截止日期，保持高效
- 💻 **代码模板** - 收藏和管理常用代码片段
- 📄 **Markdown编辑器** - 支持数学公式、代码高亮的强大编辑器
- ⭐ **收藏夹** - 快速访问常用网页和资源

### 🎯 设计理念

- **简洁** - 精简界面，专注于核心功能
- **高效** - 快速访问，流畅操作
- **稳定** - 可靠的数据存储，不丢失任何内容
- **跨平台** - 支持 Windows、macOS、Linux

## 🚀 快速开始

### 前置要求

- Go 1.25.5 或更高版本
- Node.js 16.0 或更高版本
- Wails CLI v2.11.0

### 安装

```bash
# 克隆仓库
git clone https://github.com/lsbbagt/DeskMate.git
cd DeskMate

# 安装 Wails CLI（如果尚未安装）
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装前端依赖
cd frontend
npm install
cd ..

# 运行开发模式
wails dev

# 构建生产版本
wails build
```

### 运行

```bash
# 开发模式（支持热重载）
wails dev

# 生产构建
wails build
./build/bin/DeskMate.exe
```

## 📦 项目结构

```
DeskMate/
├── build/              # 构建配置和资源
│   ├── appicon.ico     # 应用图标
│   └── bin/            # 编译输出
├── frontend/           # 前端代码
│   ├── src/
│   │   ├── components/ # Vue 组件
│   │   ├── views/      # 页面视图
│   │   ├── stores/     # Pinia 状态管理
│   │   └── utils/      # 工具函数
│   └── package.json
├── app.go              # 后端核心逻辑
├── main.go             # 应用入口
├── wails.json          # Wails 配置
└── _CONTEXT.md         # 项目上下文文档
```

## 🛠️ 技术栈

### 后端
- **Go 1.25.5** - 核心逻辑
- **Wails v2.11.0** - 桌面应用框架

### 前端
- **Vue 3** - 渐进式 JavaScript 框架
- **Vuetify 3** - Material Design 组件库
- **TypeScript** - 类型安全
- **Pinia** - 状态管理
- **Markdown-it** - Markdown 渲染
- **KaTeX** - 数学公式渲染
- **Highlight.js** - 代码高亮

## 📝 功能模块

### 1. 待办事项 (Todo)
- ✅ 创建、编辑、删除任务
- ✅ 设置截止日期和描述
- ✅ 任务完成标记
- ✅ 本地数据持久化

### 2. 代码模板 (Code Templates)
- ✅ 按文件夹组织代码片段
- ✅ 支持多种编程语言
- ✅ 语法高亮显示
- ✅ 快速复制和打开

### 3. Markdown 编辑器
- ✅ 实时预览
- ✅ 数学公式支持 (KaTeX)
- ✅ 代码高亮 (Highlight.js)
- ✅ 文件夹管理
- ✅ 导入/导出功能

### 4. 收藏夹 (Bookmarks)
- ✅ 网页链接收藏
- ✅ 分类管理
- ✅ 快速访问

## 🔧 配置

应用数据存储在用户主目录：
- Windows: `C:\Users\<用户名>\.deskmate\`
- macOS: `~/.deskmate/`
- Linux: `~/.deskmate/`

## 📊 更新日志

### v0.7 (2026-03-29)
- ✅ 待办事项编辑功能
- ✅ 应用重命名为 DeskMate
- ✅ 优化编译配置
- ✅ 清理项目文件

### v0.6 (2026-03-29)
- ✅ 文件系统同步功能
- ✅ Markdown 编辑器优化
- ✅ 预览区折叠功能

### v0.5 (2026-03-29)
- ✅ Markdown 编辑器模块
- ✅ 数学公式渲染
- ✅ 代码高亮

<details>
<summary>查看更多版本历史</summary>

### v0.4 (2026-03-29)
- ✅ 代码模板模块
- ✅ 文件夹管理

### v0.3 (2026-03-29)
- ✅ 待办事项模块
- ✅ 本地数据持久化

### v0.2 (2026-03-29)
- ✅ 收藏夹模块
- ✅ 基础架构搭建

### v0.1 (2026-03-29)
- ✅ 项目初始化
- ✅ 基础UI框架
</details>

## 🤝 贡献指南

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 👨‍💻 作者

**lsbbagt**
- GitHub: [@lsbbagt](https://github.com/lsbbagt)
- Email: 2771474297@qq.com

## 🙏 致谢

- [Wails](https://wails.io/) - 优秀的 Go 桌面应用框架
- [Vuetify](https://vuetifyjs.com/) - 强大的 Material Design 组件库
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架

---

<div align="center">
如果这个项目对你有帮助，请给一个 ⭐️ 支持一下！
</div>
