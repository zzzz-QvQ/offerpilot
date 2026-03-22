# OfferPilot 发布前检查清单

本文档用于在本地演示、项目答辩或阶段性交付前做最后一轮收口检查。

## 一、问题清单（按严重程度排序）

### 1. 当前最主要的非阻塞问题：本地完整演示依赖真实环境准备

项目虽然已经具备完整前后端代码和本地依赖启动能力，但以下条件如果没有提前准备好，会直接影响完整演示效果：

- MySQL 中需要有测试用户和基础业务数据
- LLM / Embedding 相关环境变量需要填写真实可用值
- Milvus 需要正常启动，且题库知识最好已完成向量化入库

结论：这不是代码阻塞，但属于演示前必须确认的前置条件。

### 2. 前端生产构建存在较大 chunk 告警

前端已经可以正常构建，但 `ECharts` 相关产物较大，Vite 构建时仍有 chunk size warning。

结论：不阻塞演示，不影响当前交付，但如果后续做线上部署，建议继续优化分包策略。

### 3. 当前 Docker Compose 只覆盖依赖服务，不包含前后端容器

仓库已经提供 `docker-compose.yml`，但目前只负责本地启动：

- MySQL
- Redis
- Milvus 及其依赖

`frontend` 和 `backend` 仍需要本地源码启动。

结论：不阻塞本地演示，但如果目标是“一键容器化启动全项目”，还需要补 `Dockerfile`。

## 二、检查结果总览

| 检查项 | 结果 | 说明 |
| --- | --- | --- |
| 前端是否能启动 | 基本通过 | 已通过 `npm run build`，本地开发启动依赖前端环境变量配置 |
| 后端是否能启动 | 基本通过 | 已通过 `go build ./...`，实际运行依赖数据库 / 模型 / 向量库配置 |
| 关键页面是否可访问 | 通过 | 登录、Dashboard、题库、Interview、Review、Project、Settings 路由都已存在 |
| Interview 是否能完整跑通 | 有条件通过 | 代码链路完整，实际演示依赖认证、LLM、Embedding、Milvus 和题库数据 |
| Review 是否能展示结果 | 有条件通过 | 页面和接口链路已接通，依赖 review / interview 相关真实数据 |
| Project Polish 是否能生成结果 | 有条件通过 | 接口链路已具备，依赖可用的 LLM 配置 |
| 核心文档是否齐全 | 通过 | README / DEMO / ARCHITECTURE / INTERVIEW_PACKAGE 已齐全 |
| 是否仍有明显阻碍演示的问题 | 有 3 项前置依赖风险 | 主要是数据、模型配置和向量库准备，不是结构性代码问题 |

## 三、已执行的检查与小修复

本轮实际完成了以下检查：

- 前端构建检查：`npm.cmd run build`
- 后端编译检查：`go build ./...`
- 路由检查：确认关键页面路由已注册
- 文档检查：确认下列文档已存在
  - `README.md`
  - `docs/DEMO_FLOW.md`
  - `docs/ARCHITECTURE.md`
  - `docs/INTERVIEW_PACKAGE.md`

本轮顺手修复了 3 个影响演示体验的小问题：

- 修复 `Review` 页面空状态文案乱码
- 修复统一请求层的错误提示乱码
- 修复 `Interview` 页面里残留的英文按钮 / 提示文案

## 四、详细发布前检查清单

### 1. 前端检查

- [x] `frontend/.env.example` 已提供
- [x] 前端请求层已配置 `VITE_API_BASE_URL`
- [x] `npm run build` 可通过
- [x] 登录页可访问
- [x] Dashboard 可访问
- [x] 题库页可访问
- [x] Interview 页可访问
- [x] Review 页可访问
- [x] Project Polish 页可访问
- [x] Settings 页可访问
- [x] 页面基础 loading / error / empty 状态已补齐
- [ ] 如需现场演示，请提前确认 `.env.local` 指向正确后端地址

### 2. 后端检查

- [x] `backend/.env.example` 已提供
- [x] 后端编译可通过 `go build ./...`
- [x] 认证模块接口已存在
- [x] Dashboard 接口已存在
- [x] 题库接口已存在
- [x] Interview REST + SSE 接口已存在
- [x] Review 接口已存在
- [x] Project Polish 接口已存在
- [ ] 如需本地启动，请提前准备 `.env` 中的 MySQL / LLM / Embedding / Milvus 配置

### 3. 本地依赖检查

- [x] 已提供根目录 `docker-compose.yml`
- [x] 已提供根目录 `.env.example`
- [x] Compose 已覆盖 MySQL / Redis / Milvus
- [ ] 演示前确认 `docker compose up -d` 后各服务健康
- [ ] 演示前确认 MySQL 中存在测试用户
- [ ] 演示前确认题库数据已导入
- [ ] 演示前确认题库知识已完成向量化入库（如需展示 RAG reference）

### 4. Interview 演示专项检查

- [x] 创建会话接口已存在
- [x] 获取会话详情接口已存在
- [x] 提交回答接口已存在
- [x] 结束会话接口已存在
- [x] SSE 流式接口已存在
- [x] 前端已支持 `state / reference / delta / done`
- [x] 前端已支持停止流式输出
- [x] 前端已支持避免重复开启多个 EventSource
- [ ] 演示前确认 LLM 与 Embedding 配置真实可用
- [ ] 演示前确认 Milvus 检索链路可用

### 5. Review 演示专项检查

- [x] Review 页面图表与建议区已具备
- [x] Review 已接真实分析路径
- [ ] 演示前确认至少存在一条可复盘的 interview session
- [ ] 演示前确认 Review 接口能返回真实结构化结果

### 6. Project Polish 演示专项检查

- [x] Project Polish 页面已完成
- [x] 表单提交、结果展示、历史查看链路已具备
- [x] 已接真实 LLM 生成逻辑
- [ ] 演示前确认 LLM 配置可用

### 7. 文档检查

- [x] `README.md`
- [x] `docs/DEMO_FLOW.md`
- [x] `docs/ARCHITECTURE.md`
- [x] `docs/INTERVIEW_PACKAGE.md`
- [x] 环境变量示例文件已补齐
- [x] Compose 启动说明已补齐

## 五、当前是否具备演示条件

### 结论

当前项目已经具备“代码层面可交付、文档层面可讲解、演示路径层面可执行”的状态。

如果以下 4 项前置条件准备完成，就可以进入正式演示：

1. 本地依赖服务已启动
2. 后端环境变量已配置完成
3. 至少有测试用户和基础题库数据
4. LLM / Embedding / Milvus 链路可正常工作

## 六、正式演示前最后建议

建议在正式演示前再做一次最小流程自测：

1. 登录
2. 进入 Dashboard
3. 打开题库，查看一题详情
4. 发起一次 Interview，提交一条回答，观察 SSE 输出
5. 打开 Review，确认能显示结果
6. 打开 Project Polish，确认能生成结果

如果这 6 步都能走通，就说明当前项目已经达到较好的答辩 / 面试演示状态。
