# GO后端

本项目是基于 Go + `gow` 脚手架 构建的**微服务后端架构**，支持：

- 多协议：RESTful API + gRPC
- 多存储：MySQL/PostgreSQL、Redis、MinIO
- 服务治理：Etcd 注册发现、Jaeger 链路追踪
- 代码生成：Ent ORM、Wire 依赖注入、Protobuf 自动生成、TypeScript 客户端
- 容器化：Docker + Docker Compose 一键部署
- 文档自动生成：Swagger UI、OpenAPI v3
- CRUD自动生成：基于 SQL的DDL语句 一键生成 CRUD 接口和实现

## 前置环境要求

### 基础环境

| 工具             | 版本要求  | 说明          |
|----------------|-------|-------------|
| Go             | 1.21+ | 核心开发语言      |
| Docker         | 20.0+ | 容器化部署       |
| Docker Compose | v2+   | 服务编排        |
| Make           | 3.8+  | 构建工具        |
| PM2            | 可选    | 进程管理（物理机部署） |

### 中间件（自动部署）

| 中间件        | 版本要求          | 说明                                 |
|------------|---------------|------------------------------------|
| Redis      | 6.0+          | 缓存 / 会话存储 / 任务队列（asynq）           |
| Postgresql | 14+           | 关系型数据库，支持 JSON / 事务                |
| MinIO      | RELEASE.2024+ | 兼容 S3 协议的对象存储服务                    |
| Jaeger     | 1.40+         | 分布式链路追踪，可视化请求耗时                    |

## 内置API文档

### Swagger UI

- [Admin Swagger UI](http://localhost:7788/docs/)

### openapi.yaml

- [Admin openapi.yaml](http://localhost:7788/docs/openapi.yaml)

## 项目目录结构

```text
myproject/
├── api/                    # API定义/Protobuf
├── app/                    # 微服务主目录
│   └── admin/              # 单个服务
│       └── service/        # 服务启动目录
│              ├── cmd/         # 启动入口
│              ├── configs/     # 配置文件
│              └── internal/    # 业务核心代码
│                     ├── data/     # 数据访问层（Ent ORM）
│                     ├── server/   # gRPC/HTTP服务层
│                     └── service/  # 业务逻辑层
├── scripts/             # 部署脚本
├── go.mod               # Go模块依赖
├── go.sum               # 依赖校验
└── Makefile             # 构建命令
```

## 工具体系：gow 与 make

项目有两套命令行工具，定位不同，按需选择：

| 工具 | 定位 | 使用场景 |
|------|------|----------|
| `gow` | 脚手架快捷工具 | 新建项目、快速运行单个服务、代码生成速记 |
| `make` | 标准构建系统 | 安装依赖工具、批量代码生成、编译打包、Docker 部署 |

> **日常开发建议**：运行服务用 `air`（见下文热更新章节），代码生成和构建用 `make`，脚手架操作才用 `gow`。

---

### gow 脚手架工具

`gow` 是项目框架自带的 CLI，适合**项目搭建阶段**和**快速调试单个服务**。

#### 安装

```bash
go install github.com/tx7do/go-wind-toolkit/gowind/cmd/gow@latest

# 验证
gow version
```

> Windows 用户需确保 `$(go env GOPATH)\bin` 在 PATH 中：
> ```powershell
> $env:PATH += ";$(go env GOPATH)\bin"
> ```

#### 常用命令

```bash
# 项目搭建（首次创建时使用）
gow new myproject -m github.com/yourname/myproject
gow add service admin          # 新增微服务模块

# 运行服务（等价于 go run，适合临时调试）
cd app/admin/service && gow run
gow run admin                  # 根目录下指定服务名

# 代码生成（快捷方式，底层调用 ent/wire/buf）
gow ent admin                  # 生成 Ent ORM 代码
gow wire admin                 # 生成 Wire 依赖注入代码
gow api                        # 生成 Protobuf Go 代码
```

---

### make 构建系统

`make` 是项目的**标准构建入口**，覆盖从环境初始化到部署的完整流程。

项目有两层 Makefile，作用范围不同：

| 层级 | 执行目录 | 作用范围 |
|------|----------|----------|
| 项目级 `Makefile` | `backend/` | 全局：所有服务、工具安装、Docker |
| 服务级 `app.mk` | `backend/app/admin/service/` | 单个服务：运行、编译、代码生成 |

#### 首次搭建：安装开发工具链

```bash
cd backend
make init        # 一次性安装所有 CLI 工具和 protoc 插件
```

这会安装 `buf`、`ent`、`wire`、`protoc-gen-go`、`golangci-lint` 等工具，后续代码生成依赖它们。

#### 日常开发：代码生成与运行

在**服务目录**下执行（以 admin 为例）：

```bash
cd backend/app/admin/service

make run         # 运行服务（自动生成 API 代码后启动）
make gen         # 一键生成全部代码：ent + wire + api + openapi
make ent         # 仅生成 Ent ORM 代码
make wire        # 仅生成 Wire 依赖注入代码
make api         # 仅生成 Protobuf Go 代码
make openapi     # 仅生成 OpenAPI v3 文档
make build       # 编译二进制到 ./bin/
```

在**项目根目录**下执行（批量操作所有服务）：

```bash
cd backend

make api         # 生成所有服务的 API 代码
make gen         # 生成所有服务的全部代码
make build       # 编译所有服务
make ts          # 生成前端 TypeScript 客户端代码
```

#### Docker 与部署

```bash
cd backend

# 容器编排
make compose-up              # 启动所有服务 + 中间件
make compose-down            # 停止所有服务
make compose-up-without-service  # 仅启动中间件（开发推荐）

# 脚本部署（跨平台）
make docker-libs             # 仅启动中间件容器
make docker-up               # 启动全部（应用 + 中间件）
make pm2-deploy              # 用 PM2 部署微服务（物理机）

# 镜像构建
make docker                  # 构建所有服务的 Docker 镜像
```

#### 测试与质量

```bash
cd backend

make test        # 运行所有测试
make cover       # 运行测试并生成覆盖率报告
make lint        # 运行 golangci-lint 静态检查
make vet         # 运行 go vet 分析
```

## 开发模式：Air 热更新

开发时修改 `.go` / `.yaml` 文件后自动编译重启，无需手动停止服务。

### 安装

```bash
go install github.com/air-verse/air@latest
```

### 使用

```bash
cd app/admin/service
air
```

启动后 air 会监听 `cmd/`、`internal/` 目录，文件保存后自动触发重新编译和重启。

### 配置说明

配置文件位于 `app/admin/service/.air.toml`，主要配置项：

| 配置项 | 说明 |
|------|------|
| `include_dir` | 监听目录：`cmd`、`internal` |
| `include_ext` | 监听文件类型：`.go`、`.yaml` |
| `exclude_regex` | 排除文件：`_test.go`、`wire_gen.go` |
| `stop_on_error` | 编译失败时停止旧进程，修复后保存即恢复 |
| `delay` | 文件变动后延迟 1000ms 触发，避免频繁重启 |

> **注意**：`backend/pkg/` 和 `backend/api/gen/go/` 不在监听范围内，修改这两个目录后需手动重启 air。

## 一键部署整个项目

部署项目有两种方法：

1. 三方中间件和微服务都运行在Docker之下；
2. 三方中间件运行在Docker下，微服务通过PM2管理运行在OS下。

#### 1. 三方中间件和微服务都运行在Docker之下

```bash
cd ./scripts/docker
chmod +x *.sh
./full_deploy.sh
```

#### 2. 三方中间件运行在Docker下，微服务运行在OS下

先部署三方中间件：

```bash
cd ./scripts/docker
chmod +x *.sh
./libs_only.sh
```

接着部署PM2管理的微服务：

```bash
cd ./scripts/deploy
chmod +x *.sh
./pm2_service.sh
```

## 常见问题 FAQ

### go.sum 依赖验证失败

有时候，有的依赖包在`go mod tidy`之后会出现`go.sum`验证不通过的问题，解决方法：

```bash
go clean -modcache
go mod tidy
```

如果还不行的话，可以删除`go.sum`文件，然后重新执行`go mod tidy`。

### 中间件连接失败（postgres/redis/minio）

**原因**：容器内域名解析需要配置 hosts

我们需要修改`hosts`文件，修改需要管理员权限，其配置文件路径在：

- Linux：`/etc/hosts`
- MacOS: `/private/etc/hosts`
- Windows: `C:\Windows\System32\drivers\etc\hosts`

增加以下内容：

```ini
127.0.0.1 postgres
127.0.0.1 redis
127.0.0.1 minio
127.0.0.1 consul
127.0.0.1 jaeger
```

> 注意：如果注册中心使用Consul，consul的地址填写为`consul`会返回`502`，使用`localhost`或者`127.0.0.1`都可以。
> ```yaml
> registry:
> type: "consul"
>
> consul:
> address: "localhost:8500"
> ```

# Redis 认证失败 / 命令报错

- **原因**：Redis 版本过低或密码配置不正确
- **解决**：确认 Redis ≥ 6.0，并检查 `configs/data.yaml` 中 `redis.password` 配置是否与本地 Redis 一致
