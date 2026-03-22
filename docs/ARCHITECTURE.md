# OfferPilot 架构说明

本文档用于说明 OfferPilot 的整体技术架构与核心链路，适合在项目答辩、技术面试或团队内部讲解时使用。

文档目标不是做学术化的系统设计论文，而是用清晰、可复述的方式讲明这个项目是如何分层、如何通信、如何接入 LLM / RAG，以及当前边界在哪里。

## 一、整体模块划分

OfferPilot 可以分成 4 层来理解：

```text
┌──────────────────────────────────────────────┐
│                Frontend Web App             │
│ 登录 / Dashboard / 题库 / Interview / Review / Project / Settings │
└──────────────────────────────────────────────┘
                      │
                      │ REST / SSE
                      ▼
┌──────────────────────────────────────────────┐
│                Backend API Layer            │
│ Auth / Dashboard / Question / Interview / Review / Project        │
└──────────────────────────────────────────────┘
                      │
          ┌───────────┴───────────┐
          ▼                       ▼
┌───────────────────┐   ┌────────────────────────┐
│   Business Logic   │   │  AI / Retrieval Layer  │
│ Service / Repo     │   │ LLM / Embedding / RAG  │
└───────────────────┘   └────────────────────────┘
          │                       │
          ▼                       ▼
┌───────────────────┐   ┌────────────────────────┐
│ MySQL / Redis     │   │ Milvus / OpenAI-compatible APIs │
└───────────────────┘   └────────────────────────┘
```

从产品视角看，系统包含 7 个主要业务模块：

- 登录认证
- Dashboard 工作台
- 题库学习
- 模拟面试
- 复盘分析
- 项目包装
- 设置页

从技术视角看，系统由以下几个子系统组成：

- 前端展示层
- 后端 REST / SSE API 层
- 业务服务层
- 数据存储层
- LLM / Embedding / 向量检索层

## 二、前端架构说明

前端基于 Vue 3 + TypeScript + Vite，整体采用“页面层 + 组件层 + 状态层 + 组合式逻辑层 + API 层”的结构。

### 1. 前端分层结构

```text
views            页面入口
components       可复用 UI 组件
stores           Pinia 状态管理
composables      页面逻辑编排与交互收口
api              真实接口定义
utils/request    axios 请求封装
types            前端类型定义
router           路由系统
layouts          主布局
```

### 2. 分层职责

#### `views/`

页面只负责：

- 页面级布局组合
- 接入 store / composable
- 串联组件
- 承载页面级 loading / empty / error 视图

原则上不把复杂业务逻辑直接堆进单个页面文件。

#### `components/`

组件负责单一展示职责，例如：

- Dashboard 的统计卡片、列表、图表
- Interview 的问答面板、输入框、状态区、知识卡片
- Review 的雷达图、柱状图、趋势图、建议面板
- Project 的输入表单、结果区、历史记录区

这种做法保证页面不会膨胀成超大文件，也便于后续单独重构某个区块。

#### `stores/`

Pinia store 负责：

- 页面核心状态
- 请求后的数据缓存
- loading / error / empty 所需状态
- 与后端接口返回数据对齐后的标准状态结构

当前每个主要模块都有自己的 store，例如：

- `user`
- `dashboard`
- `question-bank`
- `interview`
- `review`
- `project`
- `settings`

#### `composables/`

composable 负责收口“页面交互逻辑”，比如：

- 初始化页面数据
- 提交动作
- 页面卸载时清理状态
- SSE 启停与事件分发
- 复制、切换、筛选等交互逻辑

这层存在的价值，是避免：

- 页面直接写大量逻辑
- 组件之间彼此耦合
- store 同时承担太多 UI 交互职责

#### `api/` 与 `utils/request.ts`

API 层负责：

- 定义真实接口模块
- 对齐后端契约
- 与 `axios` 请求实例配合

`utils/request.ts` 负责：

- baseURL 注入
- token 自动附带
- 请求 / 响应拦截器
- 统一错误提示

### 3. 前端架构关键词

如果需要用答辩式表达，可以这样总结：

- 组件负责展示
- composable 负责交互编排
- store 负责状态收口
- api 负责接口边界
- 页面只负责组织结构

## 三、后端架构说明

后端基于 Go + Gin + GORM，整体采用典型的分层设计：

```text
handler      HTTP 请求入口
service      业务逻辑层
repository   数据访问层
model        数据模型层
middleware   中间件
router       路由注册
config       配置加载
pkg          通用基础设施
```

### 1. `handler`

负责：

- 参数解析
- 调用 service
- 统一响应返回
- 状态码控制

原则上不在 handler 中堆复杂 prompt 或业务流程。

### 2. `service`

负责：

- 模块业务编排
- LLM / RAG 链路组织
- SSE 事件输出流程控制
- 输入输出结构整理
- 与 repository / 外部服务协作

当前系统中最关键的业务复杂度主要都收口在 service 层。

### 3. `repository`

负责：

- 封装数据库读写
- 保持与 GORM 的交互边界清晰
- 让 service 不直接拼接复杂 SQL

### 4. `model`

定义核心数据模型，例如：

- `users`
- `training_sessions`
- `questions`
- `question_favorites`
- `interview_sessions`
- `interview_messages`
- `review_reports`
- `project_polish_records`

### 5. `middleware`

当前主要承担：

- JWT 鉴权
- 认证上下文注入

未来可以继续扩展：

- 日志
- traceId
- 限流
- 权限

### 6. 后端架构关键词

可以概括为：

- handler 只接请求
- service 编排业务
- repository 管数据访问
- middleware 处理横切逻辑
- LLM / 检索能力也通过 service 抽象统一收口

## 四、REST + SSE 双通道通信说明

OfferPilot 的通信模式不是单一 REST，而是“REST + SSE”组合。

### 1. REST 通道的职责

REST 主要用于：

- 登录认证
- 数据获取
- 列表查询
- 创建资源
- 提交动作
- 获取历史记录

例如：

- `POST /api/auth/login`
- `GET /api/dashboard/summary`
- `GET /api/questions`
- `POST /api/interview/session`
- `POST /api/interview/session/:id/message`
- `GET /api/review/:sessionId`
- `POST /api/project/polish`

REST 的特点是：

- 一次请求，一次返回
- 更适合结构化数据获取
- 更适合页面初始化与普通 CRUD

### 2. SSE 通道的职责

SSE 当前主要用于 Interview 模块：

- `GET /api/interview/session/:id/stream`

它负责把模型生成过程拆成连续事件返回给前端，而不是等完整文本全部生成后一次性返回。

当前事件协议包括：

- `state`
- `reference`
- `delta`
- `done`

### 3. 为什么要 REST + SSE 双通道

因为 Interview 的真实交互有两种不同性质的任务：

- 创建会话、提交回答、结束会话：适合 REST
- 模型生成追问、评分状态变化、知识命中流式输出：适合 SSE

所以当前设计的边界是：

- REST 负责“触发动作”和“获取快照”
- SSE 负责“持续输出过程”

这也是一个比较适合讲给面试官的设计点，因为它说明系统不是简单聊天接口，而是有明确通信边界设计。

## 五、LLM / Embedding / Milvus / RAG 的链路说明

当前项目已经具备一套可复用的 AI 能力层，核心组件包括：

- `LLMService`
- `EmbeddingService`
- `QuestionRetrievalService`
- Milvus vector store 封装

### 1. LLM service

统一负责调用 OpenAI-compatible 模型接口，提供两类能力：

- 普通文本生成
- 流式文本生成

这样 Interview、Review、Project 三个模块不需要各自重复实现模型调用细节。

### 2. Embedding service

负责把文本转换成向量，用于知识入库和相似度检索。

当前主要服务于题库 questions 的知识入库能力。

### 3. Milvus 封装

Milvus 主要提供三类基础能力：

- collection 初始化
- 向量写入
- 相似度检索

当前 question 知识文档会被拼成统一文本后写入向量库。

### 4. RAG 链路说明

当前 RAG 的典型路径如下：

```text
用户问题 / 用户回答 / 当前题目上下文
        ↓
拼接检索 query
        ↓
Embedding 向量化
        ↓
Milvus 相似度检索
        ↓
返回题库知识片段 reference
        ↓
将 reference 拼接到 prompt 中
        ↓
调用 LLM 生成反馈 / 追问 / 分析结果
```

### 5. 为什么这样设计

这样做的好处是：

- LLM 不只依赖通用世界知识
- 输出能结合项目私有题库内容
- 能让 Interview 与 Review 的结果更贴近“前端面试训练”场景
- 后续也方便将 Project、Review、题库管理进一步接入知识系统

## 六、Interview 模块 Agent 闭环说明

Interview 是当前项目最核心、最能体现系统复杂度的模块。

它不是简单的“提问 + 回答”，而是一条带 Agent 状态、RAG 检索、评分分析、SSE 输出的闭环链路。

### 1. 闭环路径

```text
创建会话
  ↓
LLM 生成首题
  ↓
用户提交回答
  ↓
检索相关知识 reference
  ↓
结合 reference + 当前回答做评分 / 分析
  ↓
LLM 生成反馈与追问
  ↓
SSE 将状态、reference、文本增量返回前端
  ↓
前端完成本轮展示
```

### 2. Interview 事件流

当前 SSE 事件协议为：

- `state`
  - 表示当前 Agent 所处阶段
- `reference`
  - 表示本轮检索命中的知识片段
- `delta`
  - 表示流式文本增量
- `done`
  - 表示本轮结束并返回最终快照

### 3. Agent 阶段说明

当前前后端已经围绕以下阶段对齐：

- `session_created`
- `generating_question`
- `retrieving_knowledge`
- `scoring`
- `generating_followup`
- `completed`

这让 Interview 页面右侧状态面板可以实时反映系统在做什么。

### 4. 前端如何消费这个闭环

前端 Interview 模块的职责分配是：

- 页面负责三栏布局组织
- `useInterviewSession` 负责发起 REST 与 SSE
- `useSSEStream` 负责 EventSource 封装
- `interview store` 负责统一收口消息、状态、评分、reference
- 组件只负责展示对应区块

这也是一个很适合答辩时强调的点：

- 状态变化是事件驱动的
- 页面不是直接操作 EventSource
- RAG 检索结果和文本流是分通道进入 store 的

## 七、Review 与 Project Polish 的生成链路说明

虽然 Interview 是最复杂的链路，但 Review 和 Project 也已经接入真实生成逻辑。

## 1. Review 生成链路

Review 的目标不是简单展示统计数据，而是基于真实会话内容生成可读分析。

### 输入来源

- `interview_messages`
- 会话基本信息
- 可用时的 reference 检索结果

### 输出内容

- overall summary
- overall score
- dimension scores
- weak points
- suggestions
- 趋势展示所需结构

### 生成流程

```text
读取指定会话
   ↓
加载 interview messages
   ↓
整理会话上下文
   ↓
可选补充 reference 检索结果
   ↓
调用 LLM 生成结构化分析
   ↓
转换为前端图表和面板可直接消费的结构
```

### 设计价值

- 不是前端本地拼文案
- 不是单纯规则评分
- 结果与真实面试内容强相关

## 2. Project Polish 生成链路

Project Polish 主要服务于“把原始项目材料转成更适合求职表达的内容”。

### 输入来源

用户输入的项目原始材料，例如：

- projectName
- background
- techStack
- responsibilities
- difficulties
- achievements

### 输出内容

- `resumeVersion`
- `interviewVersion`
- `highlights`
- `followups`

### 生成流程

```text
接收项目原始输入
   ↓
整理 prompt 上下文
   ↓
调用统一 LLM service
   ↓
解析结构化输出
   ↓
写入历史记录
   ↓
返回给前端结果页
```

### 设计价值

这个模块的亮点在于：

- 不是纯模板拼接
- 结果可以直接面向简历和面试准备场景使用
- 历史记录机制方便用户反复打磨项目材料

## 八、当前项目的边界与后续可扩展方向

当前项目已经具备比较完整的闭环，但边界也很明确。

## 1. 当前边界

### 已完成的核心部分

- 前后端主业务页面和 API 基本贯通
- Interview 已形成真实 LLM + RAG + SSE 闭环
- Review 已接入真实会话分析路径
- Project Polish 已接入真实 LLM 生成
- Question 已具备向量化与检索基础能力

### 当前仍然偏“工程原型 / 可演示版本”的部分

- Docker 化还只覆盖依赖服务，未覆盖前后端容器
- Redis 目前更多是技术栈预留，还可继续接入真实缓存 / 会话能力
- 管理后台尚未补齐
- 更复杂的权限体系、审计日志、监控与测试还未完整建设
- LLM 输出校验和重试策略仍可增强

## 2. 可扩展方向

### 工程方向

- 增加 `frontend` / `backend` Dockerfile
- 增加 CI/CD、测试、日志、监控
- 增加 traceId、限流、告警

### 产品方向

- 多种面试模式
- 自定义岗位画像
- 更多题库来源
- 训练计划与成长轨迹
- 项目包装结果版本对比

### AI 方向

- 更强的结构化输出约束
- 更完整的 Agent 状态编排
- 多轮 reference 聚合
- Review 报告持久化增强
- 更细粒度评分模型

## 九、适合答辩时的一段总结

如果需要在答辩或面试中用一段话概括项目架构，可以这样表述：

> OfferPilot 采用前后端分层架构，前端通过页面、组件、store 和 composable 进行职责拆分，后端通过 handler、service、repository、model 分层组织业务逻辑。在通信层面，普通数据交互使用 REST，模拟面试的生成过程使用 SSE 进行流式输出。在 AI 能力层面，系统通过统一 LLM service、Embedding service 与 Milvus 检索能力构建了 RAG 基础设施，并将其用于 Interview、Review、Project 等模块，形成了从学习、面试、复盘到项目包装的完整训练闭环。
