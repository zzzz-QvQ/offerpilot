# OfferPilot 项目规划

## 1. 项目结构规划

第一阶段目标是完成前端工程初始化与页面壳子，因此目录设计以“清晰分层、方便扩展、便于后续接入真实后端”为核心。

```text
OfferPilot/
├─ frontend/
│  ├─ public/
│  ├─ src/
│  │  ├─ assets/                 # 静态资源：图片、图标、全局样式变量
│  │  ├─ components/             # 通用组件
│  │  │  ├─ business/            # 业务通用组件，如题卡、面试记录卡片
│  │  │  ├─ charts/              # ECharts 图表组件封装
│  │  │  ├─ common/              # Button、Empty、Tag 等二次封装
│  │  │  └─ layout/              # 顶栏、侧边栏、面包屑、页容器
│  │  ├─ composables/            # 组合式逻辑封装
│  │  ├─ config/                 # 环境配置、导航配置、页面常量
│  │  ├─ constants/              # 常量枚举、字典映射
│  │  ├─ layouts/                # 页面布局
│  │  │  ├─ MainLayout.vue
│  │  │  ├─ AuthLayout.vue
│  │  │  └─ EmptyLayout.vue
│  │  ├─ router/                 # 路由入口与模块路由
│  │  │  ├─ index.ts
│  │  │  └─ modules/
│  │  ├─ services/               # 网络请求与 API 模块
│  │  │  ├─ http/
│  │  │  └─ modules/
│  │  ├─ stores/                 # Pinia store
│  │  ├─ styles/                 # reset、主题变量、全局样式
│  │  ├─ types/                  # 全局 TypeScript 类型
│  │  ├─ utils/                  # 工具函数
│  │  ├─ views/                  # 页面级视图
│  │  │  ├─ dashboard/
│  │  │  ├─ question-bank/
│  │  │  ├─ mock-interview/
│  │  │  ├─ review/
│  │  │  ├─ project-packaging/
│  │  │  ├─ profile/
│  │  │  └─ system/
│  │  ├─ App.vue
│  │  └─ main.ts
│  ├─ .env.development
│  ├─ .env.production
│  ├─ index.html
│  ├─ package.json
│  ├─ tsconfig.json
│  └─ vite.config.ts
├─ backend/                      # 后端目录先预留，不在第一阶段展开
├─ deploy/                       # Docker Compose、部署脚本预留
├─ docs/                         # 需求文档、接口协议、原型说明
└─ PROJECT_PLAN.md
```

### 结构原则

- `views/` 只放页面级容器，不堆叠复杂逻辑。
- `components/business/` 放跨页面复用的业务组件，避免直接散落在各页面目录。
- `services/` 统一管理请求层，后续接入 Axios、SSE、Mock 时不会污染页面。
- `stores/` 只存跨组件共享状态，不把局部 UI 状态全部塞进 store。
- `composables/` 承担页面逻辑复用和状态编排，避免页面组件过胖。
- `types/` 优先定义领域模型类型，减少接口接入后的重构成本。

## 2. 页面与模块拆分建议

第一阶段只做页面壳子，建议围绕核心使用路径做一级导航。

### 一级页面建议

1. `Dashboard` 首页 / 控制台
   - 展示训练入口、学习进度、最近模拟、能力概览
   - 第一阶段可先做卡片式概览壳子

2. `QuestionBank` 题库学习
   - 题目分类导航
   - 题目列表区
   - 题目详情区
   - 学习记录 / 收藏 / 掌握状态

3. `MockInterview` 模拟面试
   - 模拟面试首页
   - 面试会话页
   - 多轮追问记录区
   - 结果摘要区

4. `Review` 复盘分析
   - 历史记录列表
   - 单次复盘详情
   - 能力雷达图 / 趋势图
   - 问题归因与改进建议区

5. `ProjectPackaging` 项目包装
   - 项目经历列表
   - 项目优化建议
   - 表达话术整理
   - 高频追问点整理

6. `Profile` 个人中心
   - 用户基本信息
   - 求职方向偏好
   - 训练目标设定

### 页面拆分原则

- 一级页面负责路由和布局，不直接承载细碎内容。
- 复杂页面采用“页面容器 + 区块组件”模式。
- 图表、筛选器、卡片列表等优先拆成独立组件，后续便于复用。

### 建议的视图拆分示例

```text
views/
├─ dashboard/
│  └─ DashboardView.vue
├─ question-bank/
│  ├─ QuestionBankView.vue
│  ├─ components/
│  │  ├─ QuestionCategoryTree.vue
│  │  ├─ QuestionListPanel.vue
│  │  └─ QuestionDetailPanel.vue
├─ mock-interview/
│  ├─ MockInterviewHomeView.vue
│  ├─ MockInterviewSessionView.vue
│  └─ components/
├─ review/
│  ├─ ReviewListView.vue
│  ├─ ReviewDetailView.vue
│  └─ components/
└─ project-packaging/
   ├─ ProjectPackagingView.vue
   └─ components/
```

## 3. Store 设计建议

Pinia 按“领域状态”拆分，不按页面拆分。

### 推荐 store 划分

1. `appStore`
   - 全局 UI 状态
   - 侧边栏折叠状态
   - 当前主题
   - 全局 loading 状态

2. `userStore`
   - 用户信息
   - 登录态占位
   - 用户偏好配置
   - 求职方向、目标岗位等

3. `questionBankStore`
   - 当前分类
   - 筛选条件
   - 当前题目
   - 收藏 / 掌握状态的本地占位数据

4. `mockInterviewStore`
   - 当前面试配置
   - 面试会话状态
   - 消息流列表
   - 当前追问轮次
   - 结果摘要占位

5. `reviewStore`
   - 历史复盘列表
   - 当前复盘详情
   - 能力图表数据

6. `projectPackagingStore`
   - 项目列表
   - 当前项目详情
   - 项目亮点与风险点草稿

### Store 使用原则

- store 存共享状态和领域状态，不存纯展示型临时状态。
- 异步请求调用放在 action 中，但具体请求函数由 `services/` 提供。
- 页面内一次性状态，如弹窗显隐、tab 切换，优先放组件内或 composable。
- 每个 store 预留 `reset` 方法，方便后续退出登录或切换会话时清理状态。

## 4. Composables 设计建议

组合式函数用于承接“页面逻辑复用”和“UI 交互编排”。

### 推荐 composables 划分

- `usePagination`
  - 统一分页状态与翻页逻辑

- `useTableSelection`
  - 表格选中态处理

- `useQuestionFilters`
  - 题库筛选条件管理

- `useMockInterviewSession`
  - 面试消息列表、输入态、滚动定位、状态切换

- `useSseInterviewStream`
  - 后续对接 SSE 的流式消息消费封装
  - 第一阶段先只保留接口和 mock 适配位

- `useChartOptions`
  - 复盘模块的图表 option 生成器

- `useAsyncState`
  - 通用异步状态封装：`loading / error / run`

- `usePageTitle`
  - 页面标题与面包屑联动

### Composables 原则

- 一个 composable 只负责一类能力，不做“大而全”的页面总管。
- 依赖路由、store、API 时保持显式注入或清晰引用，避免隐式耦合。
- 对外暴露状态时优先返回只读或清晰命名的方法，减少误用。

## 5. API 模块预留方案

第一阶段不接真实后端，但请求层必须提前设计好，确保后续切换成本低。

### 目录建议

```text
services/
├─ http/
│  ├─ client.ts          # Axios 实例
│  ├─ interceptors.ts    # 请求/响应拦截器
│  ├─ request.ts         # 通用 request 封装
│  └─ types.ts           # 通用响应结构定义
├─ modules/
│  ├─ auth.ts
│  ├─ user.ts
│  ├─ question-bank.ts
│  ├─ mock-interview.ts
│  ├─ review.ts
│  └─ project-packaging.ts
└─ sse/
   └─ interview-stream.ts
```

### 设计建议

- `http/client.ts`
  - 创建 Axios 实例
  - 管理 baseURL、超时、通用 header

- `http/interceptors.ts`
  - 请求头注入 token
  - 统一错误处理
  - 统一业务码处理

- `http/request.ts`
  - 二次封装 `get/post/put/delete`
  - 可扩展重试、取消请求、泛型推导

- `modules/*.ts`
  - 每个领域一个 API 文件
  - 页面只调用模块函数，不直接写请求地址

- `sse/interview-stream.ts`
  - 专门处理模拟面试流式响应
  - 后续可以统一封装 EventSource 或 fetch stream 兼容逻辑

### Mock 预留建议

- 第一阶段可先在 `services/modules/` 返回本地 Promise 假数据。
- 若后续页面联调变多，再引入 `mock/` 目录统一维护假数据。
- 页面层不要感知数据来源是真接口还是 mock。

### 建议的接口返回结构

```ts
interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  traceId?: string;
}
```

流式 SSE 可独立设计：

```ts
interface InterviewStreamEvent {
  event: 'start' | 'question' | 'follow_up' | 'summary' | 'done' | 'error';
  data: Record<string, unknown>;
}
```

## 6. 第一阶段落地建议

当前阶段建议只做下面这些内容：

1. 初始化前端工程
   - Vue 3 + TypeScript + Vite
   - Pinia
   - Vue Router
   - Element Plus
   - Axios
   - ECharts

2. 完成基础工程骨架
   - 目录创建
   - 路由骨架
   - Layout 骨架
   - 全局样式与主题变量

3. 完成核心页面壳子
   - Dashboard
   - 题库学习
   - 模拟面试
   - 复盘分析
   - 项目包装
   - 个人中心

4. 预留数据层能力
   - store 空实现
   - API 模块占位
   - composables 占位

## 7. 下一步实施顺序建议

建议按以下顺序继续推进：

1. 初始化 `frontend/` 工程
2. 搭建 `layouts + router + views` 页面壳子
3. 建立 `stores + services + composables` 基础空文件
4. 补齐导航、面包屑、基础页容器组件
5. 再进入第二轮细化各页面 UI

---

这份规划刻意控制在“可执行骨架”层面，避免在第一阶段过早写入业务细节，便于后续逐步迭代。
