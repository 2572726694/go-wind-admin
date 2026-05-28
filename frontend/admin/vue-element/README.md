<div align="center">
  <img alt="go-wind-admin element-admin" width="80" height="80" src="./src/assets/images/logo.png">
  <h1>GoWind Admin - Element Plus</h1>

  <img src="https://img.shields.io/badge/Vue-3.5.30-brightgreen.svg"/>
  <img src="https://img.shields.io/badge/Vite-8.0.0-green.svg"/>
  <img src="https://img.shields.io/badge/Element Plus-2.13.5-blue.svg"/>
  <img src="https://img.shields.io/badge/TypeScript-5.9.3-blue.svg"/>
  <img src="https://img.shields.io/badge/license-MIT-green.svg"/>
</div>

## 项目简介

[GoWind Admin - Element Plus](https://github.com/tx7do/go-wind-admin/tree/main/frontend/admin/vue-element) 是一个基于 Vue3、Vite 8、TypeScript 和 Element-Plus 构建的企业级后台管理模板。项目在 vue3-element-admin 基础上进行定制开发，适配 Go-Wind 后端服务。

**核心技术栈：**

- 🎯 Vue 3.5.30+ - Composition API
- ⚡ Vite 8.0.0+ - 下一代前端构建工具
- 💎 TypeScript 5.9.3+ - 类型安全
- 🎨 Element Plus 2.13.5+ - UI 组件库
- 📊 VXE Table 4.6.25+ - 高级表格组件
- 📝 TipTap 3.20.0+ - 富文本编辑器
- 🎨 UnoCSS - 原子化 CSS 引擎

**主要特性：**

- ✅ **系统管理**：用户、角色、菜单、部门、字典、系统配置、通知公告
- ✅ **权限管理**：动态路由、按钮权限、角色权限、数据权限
- ✅ **多租户**：支持多租户模式和租户隔离（可配置开关）
- ✅ **应用模块**：操作日志、内部消息、租户管理、权限管理等
- ✅ **基础设施**：国际化、多布局、暗黑模式、全屏、水印、接口文档
- ✅ **开发者工具**：代码生成器、Mock 服务器、ESLint + Prettier + Stylelint

## 功能模块

### 🔐 核心功能

- **登录认证**：账号密码登录、Token 管理、多端认证
- **个人中心**：基本信息设置、密码修改、安全设置、邮箱绑定、消息通知
- **工作台**：数据分析、统计图表、业务概览

### 👥 组织与人员管理 (OPM)

- **组织架构**：部门管理、树形结构展示、层级维护
- **岗位管理**：岗位配置、岗位分配、职责管理
- **用户管理**：用户列表、新增编辑、角色分配、状态控制、详情查看

### 🛡️ 权限管理

- **权限点管理**：按钮权限、接口权限、权限码配置
- **角色管理**：角色列表、权限配置、数据权限、角色分配

### ⚙️ 系统配置

- **菜单管理**：菜单配置、路由映射、图标选择、权限标识
- **API 管理**：接口定义、路由配置、权限绑定
- **字典管理**：数据字典、字典项管理、字典分类
- **文件管理**：文件上传、存储管理、文件预览
- **任务管理**：定时任务、任务调度、执行日志
- **登录策略**：安全策略、登录限制、密码策略
- **语言包管理**：国际化配置、多语言支持、语言包管理

### 📋 审计日志

- **登录日志**：登录记录、登录时间、登录 IP、登录状态
- **API 日志**：接口调用、请求参数、响应时间
- **操作日志**：操作记录、变更内容、操作人追踪
- **数据访问日志**：数据查询、数据修改、数据导出
- **权限审计日志**：权限变更、权限使用、权限异常

### 🏭 租户管理

- **租户管理**：租户配置、租户状态、有效期管理、资源配额

### 💬 内部消息

- **消息分类**：消息分类管理、消息类型配置
- **消息管理**：消息发送、消息列表、消息已读标记

## 环境要求

| 环境类型        | 版本要求                                                         | 备注                       |
|-------------|--------------------------------------------------------------|--------------------------|
| **Node.js** | `^20.19.0` 或 `>=22.12.0`                                     | 推荐使用 LTS 版本（主版本为偶数）      |
| **包管理器**    | `pnpm >= 8.0.0`                                              | 项目使用 pnpm 作为包管理器         |
| **开发工具**    | [Visual Studio Code](https://code.visualstudio.com/Download) | 推荐安装 Vue、TypeScript 相关插件 |

## 快速开始

```bash
# 安装 pnpm（如未安装）
npm install pnpm -g

# 设置镜像源（可选，国内用户推荐）
pnpm config set registry https://registry.npmmirror.com

# 安装依赖
pnpm install

# 启动开发服务器
pnpm run dev

# 访问地址：http://localhost:3000
```

## 构建部署

```bash
# 生产环境构建
pnpm run build

# 构建完成后生成 dist 目录
# 将 dist 目录部署到 Nginx 或其他静态服务器
```

### Nginx 配置示例

```nginx
server {
    listen      80;
    server_name localhost;

    location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
    }

    # 反向代理配置
    location /admin/v1/ {
        # 替换为您的后端 API 地址
        proxy_pass http://localhost:7788/;
    }
}
```

## 配置说明

### 环境变量

项目提供了两个环境配置文件：

**开发环境** (`.env.development`)：

```bash
# 应用端口
VITE_APP_PORT=3000

# 代理前缀
VITE_APP_BASE_API=/admin/v1

# 接口地址（默认指向本地后端）
VITE_APP_API_URL=http://localhost:7788

# SSE 实时推送地址
VITE_APP_SSE_URL=http://localhost:7789/events

# 启用 Mock 服务（true:开启 false:关闭）
VITE_MOCK_DEV_SERVER=false

# 多租户开关（true:开启 false:关闭）
VITE_APP_TENANT_ENABLED=false
```

**生产环境** (`.env.production`)：

```bash
# 接口地址
VITE_APP_API_URL=https://api.demo.admin.gowind.cloud

# SSE 地址
VITE_APP_SSE_URL=https://sse.demo.admin.gowind.cloud/events

# 多租户开关
VITE_APP_TENANT_ENABLED=true
```

## 项目结构

```
├── src/
│   ├── api/              # API 接口定义
│   ├── assets/           # 静态资源（图标、图片等）
│   ├── components/       # 公共组件
│   │   ├── DictSelect/  # 字典选择器
│   │   ├── DictTag/     # 字典标签
│   │   ├── Editor/      # 富文本编辑器
│   │   ├── Pagination/  # 分页组件
│   │   └── ...          # 其他组件
│   ├── composables/     # 组合式函数
│   ├── constants/       # 常量配置
│   ├── core/            # 核心模块（权限、配置、存储）
│   ├── directives/      # 自定义指令
│   ├── i18n/            # 国际化配置
│   ├── layouts/         # 布局组件
│   ├── plugins/         # 插件配置（VXE Table、ECharts）
│   ├── router/          # 路由配置
│   ├── stores/          # Pinia 状态管理
│   ├── styles/          # 全局样式
│   ├── types/           # TypeScript 类型定义
│   ├── utils/           # 工具函数
│   └── pages/           # 页面视图
│       ├── app/         # 应用模块（日志、消息、权限等）
│       └── core/        # 核心功能（登录、个人中心等）
├── mock/                # Mock 数据
├── public/              # 公共资源
└── vite.config.ts       # Vite 配置文件
```

## 技术亮点

### 1. 自动化导入

使用 `unplugin-auto-import` 和 `unplugin-vue-components` 实现：

- Vue Composition API 自动导入
- Element Plus 组件和函数自动导入
- @vueuse/core 函数自动导入
- 无需手动 import，提升开发效率

### 2. 现代化 CSS 方案

- **UnoCSS**：原子化 CSS 引擎，支持即时按需生成样式
- **SCSS**：预处理器，支持变量和混入
- **CSS Modules**：组件级样式隔离

### 3. 类型安全

- **TypeScript 严格模式**：全面的类型定义
- **自动类型生成**：API 和组件类型自动生成
- **Zod 验证**：运行时数据验证

### 4. 高性能优化

- **依赖预加载**：预加载常用组件，提升首次渲染速度
- **代码分割**：路由级别的代码分割
- **懒加载**：图片和组件懒加载
- **打包优化**：Terser 压缩、移除 console 和 debugger

### 5. 开发者体验

- **ESLint + Prettier + Stylelint**：全方位的代码质量检查
- **Husky + Commitlint**：Git 提交规范强制
- **Mock Server**：本地 Mock 开发，不依赖后端
- **热更新**：Vite HMR，秒级刷新

## 脚本命令

| 命令                          | 说明              |
|-----------------------------|-----------------|
| `pnpm run dev`              | 启动开发服务器         |
| `pnpm run build`            | 类型检查并构建生产版本     |
| `pnpm run build-only`       | 仅构建，不进行类型检查     |
| `pnpm run preview`          | 预览生产构建          |
| `pnpm run type-check`       | TypeScript 类型检查 |
| `pnpm run lint`             | 执行所有代码规范检查      |
| `pnpm run lint:eslint`      | ESLint 代码检查     |
| `pnpm run lint:prettier`    | Prettier 代码格式化  |
| `pnpm run lint:stylelint`   | Stylelint 样式检查  |
| `pnpm run lint:lint-staged` | 提交前代码检查         |
| `pnpm run commit`           | 使用 cz-git 交互式提交 |

## 开发规范

### 代码规范

项目使用严格的代码规范，配置文件包括：

- `.eslintrc` / `eslint.config.ts`：JavaScript/TypeScript 代码规范
- `.prettierrc.yaml`：代码格式化配置
- `.stylelintrc.cjs`：CSS/SCSS 样式规范
- `.editorconfig`：编辑器基础配置

### 提交规范

项目采用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 使用交互式提交
git commit
# 或使用命令
pnpm run commit
```

**提交格式：** `<type>(<scope>): <subject>`

**Type 类型：**

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档变更
- `style`: 代码格式调整
- `refactor`: 重构代码
- `test`: 测试相关
- `chore`: 构建/工具链相关

## 常见问题

### 1. 启动报错或依赖问题

```bash
# 清除缓存并重新安装
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

### 2. 浏览器访问空白

升级浏览器至最新版本，低版本浏览器内核可能不支持可选链操作符 `?.` 等新语法。

### 3. IDE 爆红但能正常运行

重启 VSCode，或等待 TypeScript 服务初始化完成。

### 4. 同步更新后启动失败

```bash
# 清除缓存并重新安装
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

### 5. Mock 数据不生效

检查 `.env.development` 中 `VITE_MOCK_DEV_SERVER` 是否为 `true`，重启开发服务器。

### 6. 接口请求失败

检查 `.env.development` 中的 `VITE_APP_API_URL` 配置是否正确，确保后端服务已启动。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 特别感谢

- [Vue.js](https://vuejs.org/) - 优雅的渐进式 JavaScript 框架
- [Vite](https://vitejs.dev/) - 下一代前端开发工具
- [Element Plus](https://element-plus.org/) - 基于 Vue 3 的组件库
- [Vue3-Element-Admin](https://gitee.com/panjiachen/vue-element-admin) - 灵感来源

---

**如果这个项目对你有帮助，请给一个 ⭐️ Star 支持！**
