# OfferPilot

OfferPilot 是一个面向前端求职者的前端面试智能训练平台，目标是帮助用户围绕“刷题、模拟面试、复盘分析、项目包装”形成一条完整训练闭环。

项目当前已具备前后端分层架构、真实 API 接入、Interview 模块 SSE 流式输出，以及基于 LLM / RAG 的核心能力基础，适合继续演进为可演示、可扩展、可联调的正式项目。

## 项目简介

OfferPilot 聚焦前端求职场景，围绕候选人在真实面试准备过程中最常见的几个环节提供支持：

- 题库学习：按分类、频率、关键词检索前端题目，查看题目详情与收藏状态
- 模拟面试：创建面试会话，提交回答，实时接收反馈、追问、评分状态和知识命中结果
- 复盘分析：基于真实会话记录生成复盘报告、维度评分、薄弱点与学习建议
- 项目包装：将原始项目材料整理为更适合简历表达、面试表达和追问准备的内容
- 设置管理：配置模型名称、面试模式、追问强度、评分严格度等偏好项

## 核心功能模块

### 1. 登录认证

- 邮箱 + 密码登录
- 获取当前用户信息
- 本地登录态恢复
- 退出登录

### 2. Dashboard 工作台

- 训练数据概览
- 最近训练记录
- 薄弱知识点图表
- 推荐训练任务面板

### 3. 题库学习

- 关键词搜索
- 分类筛选
- 高频 / 中频筛选
- 题目详情查看
- 收藏与取消收藏

### 4. 模拟面试

- 创建会话
- 获取会话详情
- 提交回答
- SSE 实时流式反馈
- Agent 状态展示
- RAG reference 知识命中展示
- 会话结束

### 5. 复盘分析

- 总评分展示
- 多维评分雷达图
- 薄弱知识点柱状图
- 历史趋势折线图
- 学习建议与薄弱点总结

### 6. 项目包装

- 输入项目原始材料
- 生成简历版项目描述
- 生成面试版项目描述
- 输出技术亮点、难点与解决方案、追问建议
- 查看历史记录与历史详情

### 7. 设置页

- 默认面试模式
- 追问强度
- 评分严格程度
- 模型名称
- 流式输出开关

## 技术栈说明

### 前端

- Vue 3
- TypeScript
- Vite
- Pinia
- Vue Router
- Element Plus
- Axios
- ECharts
- dayjs
- @vueuse/core

### 后端

- Go
- Gin
- GORM
- MySQL
- JWT
- SSE
- Redis（技术栈规划中，当前可继续接入缓存 / 会话 / 黑名单能力）
- Milvus
- Docker Compose（本地依赖启动已补充）

### AI / 检索能力

- OpenAI-compatible LLM 接口
- OpenAI-compatible Embedding 接口
- Milvus 向量检索
- Interview / Review / Project 共享 LLM service
- Question 知识向量化与 RAG 检索基础能力

## 项目目录结构

```text
OfferPilot/
├─ frontend/                # Vue 3 前端工程
│  ├─ src/
│  │  ├─ api/               # API 接口定义
│  │  ├─ components/        # 页面组件
│  │  ├─ composables/       # 组合式逻辑
│  │  ├─ layouts/           # 主布局
│  │  ├─ router/            # 路由配置
│  │  ├─ stores/            # Pinia 状态管理
│  │  ├─ styles/            # 全局样式
│  │  ├─ types/             # TypeScript 类型定义
│  │  ├─ utils/             # 请求封装与工具方法
│  │  └─ views/             # 页面入口
│  ├─ .env.example
│  └─ package.json
├─ backend/                 # Go 后端工程
│  ├─ cmd/server/           # 服务启动入口
│  ├─ internal/
│  │  ├─ config/            # 配置读取
│  │  ├─ database/          # 数据库初始化
│  │  ├─ handler/           # HTTP Handler
│  │  ├─ middleware/        # 中间件
│  │  ├─ model/             # 数据模型
│  │  ├─ pkg/               # 通用基础设施
│  │  ├─ repository/        # 数据访问层
│  │  ├─ router/            # 路由注册
│  │  └─ service/           # 业务逻辑 / LLM / RAG
│  ├─ .env.example
│  └─ go.mod
├─ docs/                    # 设计说明、接口文档、重构与收尾记录
├─ .env.example             # Docker Compose 示例变量
├─ docker-compose.yml       # 本地依赖服务启动配置
└─ README.md
```

## 前端启动方式

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 配置环境变量

```bash
cp .env.example .env.local
```

Windows 下可手动复制 `frontend/.env.example` 为 `.env.local`。

### 3. 启动开发环境

```bash
npm run dev
```

### 4. 构建生产包

```bash
npm run build
```

### 5. 本地预览构建结果

```bash
npm run preview
```

默认情况下，前端会通过 `src/utils/request.ts` 使用 `VITE_API_BASE_URL` 调用后端接口。

## 后端启动方式

### 1. 准备环境变量

```bash
cd backend
cp .env.example .env
```

Windows 下可手动复制或直接基于 `.env.example` 新建 `.env`。

### 2. 准备 MySQL

如果你已经使用下方 Docker Compose 启动依赖服务，会自动得到一个本地 MySQL 实例。

如果你使用自建 MySQL，请先创建数据库，例如：

```sql
CREATE DATABASE offerpilot CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 安装依赖并启动

```bash
go mod tidy
go run ./cmd/server
```

如果使用的是当前仓库结构，也可直接执行：

```bash
go run ./cmd/server/main.go
```

### 4. 编译检查

```bash
go build ./...
```

## Docker / Docker Compose 启动方式

当前仓库已经补充根目录 [docker-compose.yml](docker-compose.yml)，用于本地快速启动核心依赖服务：

- MySQL
- Redis
- Milvus standalone
- Milvus 依赖组件：etcd、MinIO

### 1. 准备 Compose 环境变量

```bash
cp .env.example .env
```

Windows 下可手动复制根目录 `.env.example` 为 `.env`。

### 2. 启动依赖服务

```bash
docker compose up -d
```

### 3. 查看服务状态

```bash
docker compose ps
```

### 4. 停止并移除容器

```bash
docker compose down
```

### 5. 如需同时清理数据卷

```bash
docker compose down -v
```

说明：

- 当前 compose 主要负责本地依赖服务，不包含 `frontend` 和 `backend` 容器。
- 原因是仓库当前尚未提交正式可复用的 `frontend/Dockerfile` 与 `backend/Dockerfile`。
- 推荐工作流是：先用 compose 启动 MySQL / Redis / Milvus，再本地运行前后端源码。

## 环境变量说明

前后端与 compose 的示例环境变量文件如下：

- [frontend/.env.example](frontend/.env.example)
- [backend/.env.example](backend/.env.example)
- [.env.example](.env.example)

### 前端环境变量

| 变量名 | 示例值 | 说明 |
| --- | --- | --- |
| `VITE_PORT` | `5173` | Vite 本地开发端口 |
| `VITE_APP_TITLE` | `OfferPilot` | 前端应用展示名称 |
| `VITE_API_BASE_URL` | `http://127.0.0.1:8080` | 后端 API 基础地址 |
| `VITE_ENABLE_DEVTOOLS` | `false` | 预留的开发开关示例值 |

### 后端环境变量

| 变量名 | 示例值 | 说明 |
| --- | --- | --- |
| `APP_PORT` | `8080` | 后端服务端口 |
| `APP_JWT_SECRET` | `replace_with_a_long_random_secret` | JWT 签名密钥 |
| `APP_JWT_EXPIRE_HOURS` | `72` | JWT 过期时间（小时） |
| `MYSQL_DSN` | `root:password@tcp(127.0.0.1:3306)/offerpilot?...` | MySQL 连接串 |
| `LLM_BASE_URL` | `https://api.openai.com/v1` | OpenAI-compatible LLM 接口地址 |
| `LLM_API_KEY` | 空 | LLM 调用密钥 |
| `LLM_MODEL` | `gpt-4o-mini` | 默认 LLM 模型名 |
| `EMBEDDING_BASE_URL` | `https://api.openai.com/v1` | Embedding 接口地址 |
| `EMBEDDING_API_KEY` | 空 | Embedding 调用密钥 |
| `EMBEDDING_MODEL` | `text-embedding-3-small` | 默认 Embedding 模型名 |
| `MILVUS_BASE_URL` | `http://127.0.0.1:19530` | Milvus 服务地址 |
| `MILVUS_TOKEN` | 空 | Milvus 访问令牌 |
| `MILVUS_DATABASE` | `default` | Milvus 数据库名 |
| `MILVUS_COLLECTION` | `question_knowledge` | 向量集合名 |
| `MILVUS_VECTOR_DIM` | `1536` | 向量维度 |
| `ENABLE_QUESTION_INDEXING` | `false` | 是否在启动时执行题库向量化入库 |

### Compose 环境变量

| 变量名 | 示例值 | 说明 |
| --- | --- | --- |
| `MYSQL_PORT` | `3306` | 暴露给宿主机的 MySQL 端口 |
| `MYSQL_ROOT_PASSWORD` | `password` | MySQL root 密码示例值 |
| `MYSQL_DATABASE` | `offerpilot` | MySQL 默认数据库名 |
| `MYSQL_USER` | `offerpilot` | MySQL 业务用户名 |
| `MYSQL_PASSWORD` | `offerpilot123` | MySQL 业务用户密码示例值 |
| `REDIS_PORT` | `6379` | Redis 端口 |
| `MINIO_API_PORT` | `9000` | MinIO API 端口 |
| `MINIO_CONSOLE_PORT` | `9001` | MinIO 控制台端口 |
| `MINIO_ROOT_USER` | `minioadmin` | MinIO 用户名示例值 |
| `MINIO_ROOT_PASSWORD` | `minioadmin` | MinIO 密码示例值 |
| `MILVUS_PORT` | `19530` | Milvus gRPC / SDK 端口 |
| `MILVUS_HTTP_PORT` | `9091` | Milvus HTTP 健康检查端口 |

提示：

- 示例文件中不包含任何真实密钥。
- 请将示例文件复制为本地 `.env` 或 `.env.local` 后再填写真实值。
- 生产环境建议通过部署平台的密钥管理能力注入配置，而不是直接提交到仓库。

## 演示流程建议

建议按下面顺序演示项目能力：

1. 使用 Docker Compose 启动 MySQL / Redis / Milvus
2. 启动后端服务
3. 启动前端服务
4. 登录系统
   - 使用测试账号登录，验证认证接口与登录态恢复
5. 查看 Dashboard
   - 展示训练概览、最近训练记录、薄弱项图表与推荐任务
6. 进入题库学习
   - 搜索题目、切换分类与频率、查看题目详情、操作收藏
7. 进入模拟面试
   - 创建会话
   - 提交回答
   - 演示 SSE 实时流式输出
   - 展示 Agent 状态与知识命中 reference
8. 查看复盘分析
   - 展示总体评分、多维评分、趋势图和学习建议
9. 进入项目包装
   - 输入项目原始信息
   - 生成简历版 / 面试版描述和追问建议
10. 打开设置页
   - 演示模型与面试偏好配置能力

## 当前开发进度与后续规划

### 当前进度

- 已完成前端工程初始化与页面壳子搭建
- 已完成登录、Dashboard、题库、模拟面试、复盘、项目包装、设置页的前端实现
- 已完成认证、Dashboard、题库、Interview、Review、Project 的后端 API 实现
- 已完成 Interview SSE 流式接口接入
- 已完成统一 LLM service、Embedding service、Milvus 检索基础能力
- 已完成 Question 知识向量化入库能力
- 已完成 Interview 模块的真实 LLM + RAG 闭环
- 已完成 Project Polish 的真实 LLM 生成
- 已完成 Review 基于真实会话数据的分析逻辑升级
- 已补充本地依赖服务的 Docker Compose 启动能力

### 后续规划

- 补齐 frontend / backend Dockerfile，实现整套容器化启动
- 引入 Redis 用于缓存、会话状态、SSE 辅助状态或黑名单能力
- 补充更稳定的结构化 LLM 输出校验与异常恢复
- 补充后台管理能力，例如题库管理、知识入库管理、模型配置管理
- 优化前端打包体积与按需加载策略
- 增强复盘报告持久化与多轮会话聚合分析能力
- 完善测试、监控、日志与部署流程

## 截图占位说明

当前先预留截图位置，后续建议补充以下页面截图：

```text
/docs/screenshots/
├─ login.png
├─ dashboard.png
├─ question-bank.png
├─ interview.png
├─ review.png
├─ project-polish.png
└─ settings.png
```

GitHub README 中可使用如下方式展示：

```md
## Screenshots

### Login
![Login](docs/screenshots/login.png)

### Dashboard
![Dashboard](docs/screenshots/dashboard.png)
```

当前仓库尚未提交正式截图资源，README 先保留说明位。

## License

当前仓库暂未指定正式开源协议。

如果后续计划公开发布，建议根据项目用途选择合适协议，例如：

- MIT
- Apache-2.0
- GPL-3.0

在正式确定 License 之前，建议默认按“保留所有权利，仅供学习与演示使用”处理。

---

如果你准备继续推进部署或开源发布，下一步建议优先补齐以下内容：

1. frontend / backend Dockerfile
2. 初始化 SQL / seed 数据
3. 截图资源
4. LICENSE 文件
5. 面向外部读者的快速体验账号与演示说明
