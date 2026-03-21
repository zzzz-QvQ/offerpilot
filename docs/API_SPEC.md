# OfferPilot API 契约文档

本文档基于当前前端页面、Pinia store、`src/api/*` 和 mock 数据结构整理，目标是为后端第一版接口提供稳定契约。

## 1. 通用约定

### 1.1 Base URL

- 开发环境：`/api`
- 生产环境：由网关或反向代理统一转发

### 1.2 认证方式

- 登录成功后返回 `token`
- 前端通过 `Authorization: Bearer <token>` 自动附带到后续请求

### 1.3 通用响应结构

```ts
interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
  traceId?: string;
}
```

### 1.4 分页结构

```ts
interface PageParams {
  page: number;
  pageSize: number;
}

interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
}
```

## 2. 接口列表

### 2.1 登录认证

#### `POST /auth/login`
说明：邮箱密码登录。

请求参数：

```ts
interface LoginParams {
  email: string;
  password: string;
}
```

响应参数：

```ts
interface UserProfile {
  id: string;
  email: string;
  nickname: string;
}

interface LoginResponse {
  token: string;
  user: UserProfile;
}
```

响应示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "jwt-token",
    "user": {
      "id": "u_1001",
      "email": "demo@example.com",
      "nickname": "OfferPilot User"
    }
  },
  "traceId": "trace-xxx"
}
```

#### `POST /auth/logout`
说明：退出登录。

请求参数：无。

响应参数：

```ts
type LogoutResponse = void;
```

#### `GET /auth/profile`
说明：获取当前登录用户信息。

响应参数：

```ts
type ProfileResponse = UserProfile;
```

---

### 2.2 Dashboard 数据

#### `GET /dashboard/overview`
说明：获取工作台首页聚合数据。

响应参数：

```ts
interface DashboardStatItem {
  label: string;
  value: string;
  hint: string;
}

interface RecentSessionItem {
  id: string;
  title: string;
  score: number;
  durationText: string;
  finishedAt: string;
  tag: string;
}

interface WeakPointItem {
  name: string;
  value: number;
}

interface RecommendTaskItem {
  id: string;
  title: string;
  description: string;
  level: '高优先级' | '中优先级' | '巩固建议';
}

interface DashboardOverview {
  statCards: DashboardStatItem[];
  recentSessions: RecentSessionItem[];
  weakPoints: WeakPointItem[];
  recommendTasks: RecommendTaskItem[];
}
```

建议说明：
- `statCards` 可以由后端直接返回，减少前端映射成本
- 如果后续需要国际化，可改成“结构化字段 + 前端文案映射”

---

### 2.3 题库查询

#### `GET /questions`
说明：按关键词、分类、频率分页查询题库。

请求参数：

```ts
type QuestionFrequency = '高频' | '中频';

interface QuestionListParams {
  keyword?: string;
  category?: string;
  frequency?: QuestionFrequency;
  page?: number;
  pageSize?: number;
}
```

响应参数：

```ts
interface QuestionItem {
  id: string;
  title: string;
  category: string;
  frequency: QuestionFrequency;
  isFavorite: boolean;
  summary: string;
  difficulty: '基础' | '进阶';
  answerPoints: string[];
}

type QuestionListResponse = PageResult<QuestionItem>;
```

#### `GET /questions/:questionId`
说明：获取单题详情。

响应参数：

```ts
type QuestionDetailResponse = QuestionItem;
```

#### `POST /questions/:questionId/favorite`
说明：切换收藏状态。

请求参数建议：

```ts
interface ToggleFavoriteParams {
  isFavorite?: boolean;
}
```

响应参数建议：

```ts
interface ToggleFavoriteResponse {
  questionId: string;
  isFavorite: boolean;
}
```

---

### 2.4 面试会话

#### `POST /interviews`
说明：创建新的模拟面试会话。

请求参数建议：

```ts
interface CreateInterviewSessionParams {
  mode?: '综合模拟' | '专项训练' | '项目拷打';
  category?: string;
  questionCount?: number;
}
```

响应参数：

```ts
interface CreateInterviewSessionResponse {
  sessionId: string;
}
```

#### `GET /interviews/:sessionId`
说明：获取会话详情与当前快照。

响应参数：

```ts
interface InterviewQuestionItem {
  id: string;
  title: string;
  status: '待开始' | '进行中' | '已完成';
}

interface InterviewMessageItem {
  id: string;
  role: 'interviewer' | 'candidate';
  content: string;
  time: string;
}

interface InterviewStatusItem {
  label: string;
  value: string;
}

interface InterviewScoreItem {
  label: string;
  score: number;
}

interface KnowledgeHitItem {
  name: string;
  level: '高' | '中' | '低';
  summary: string;
}

interface InterviewSessionDetail {
  sessionId: string;
  currentQuestionId: string;
  questions: InterviewQuestionItem[];
  messages: InterviewMessageItem[];
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
  knowledgeHits: KnowledgeHitItem[];
}
```

#### `POST /interviews/answer`
说明：提交当前题目的回答。

请求参数：

```ts
interface SubmitAnswerParams {
  sessionId: string;
  questionId: string;
  content: string;
}
```

响应参数：

```ts
interface SubmitAnswerResponse {
  message: InterviewMessageItem;
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
}
```

#### 预留：`GET /interviews/:sessionId/stream`
说明：后续流式问答可改为 SSE。

建议返回事件：

```ts
interface InterviewStreamEvent {
  event: 'start' | 'question' | 'follow_up' | 'score' | 'summary' | 'done' | 'error';
  data: Record<string, unknown>;
}
```

---

### 2.5 复盘报告

#### `GET /reviews/overview`
说明：获取复盘中心聚合数据。

响应参数：

```ts
interface ReviewSummaryItem {
  label: string;
  value: string;
  hint: string;
}

interface ReviewRadarItem {
  name: string;
  score: number;
  fullMark: number;
}

interface ReviewWeaknessItem {
  name: string;
  score: number;
}

interface ReviewTrendItem {
  date: string;
  score: number;
}

interface ReviewSuggestionItem {
  title: string;
  description: string;
  type: '建议' | '薄弱点';
}

interface ReviewOverview {
  overallScore: number;
  summaryItems: ReviewSummaryItem[];
  radarItems: ReviewRadarItem[];
  weaknessItems: ReviewWeaknessItem[];
  trendItems: ReviewTrendItem[];
  suggestionItems: ReviewSuggestionItem[];
}
```

建议说明：
- `trendItems` 当前按日期聚合，后续可以扩展时间粒度，如 `day/week/month`
- `overallScore` 建议与 `summaryItems` 中的展示文案分离，避免前端做字符串解析

---

### 2.6 项目包装

#### `POST /projects/polish`
说明：输入项目原始信息，生成包装结果。

请求参数：

```ts
interface ProjectPolishFormData {
  projectName: string;
  role: string;
  techStack: string;
  projectBackground: string;
  responsibility: string;
}
```

响应参数：

```ts
interface FollowupQuestionItem {
  question: string;
  answer: string;
}

interface ProjectPolishOutput {
  resumeDescription: string;
  interviewDescription: string;
  highlights: string[];
  difficulties: string[];
  followups: FollowupQuestionItem[];
}
```

#### `GET /projects/history`
说明：获取历史生成记录。

响应参数：

```ts
interface ProjectHistoryItem {
  id: string;
  projectName: string;
  createdAt: string;
  summary: string;
}

type ProjectHistoryResponse = ProjectHistoryItem[];
```

建议扩展：
- 后续可补 `GET /projects/history/:id` 查看历史详情
- 如果历史记录量变大，建议改分页

---

### 2.7 设置

#### `GET /settings`
说明：获取当前用户设置。

#### `PUT /settings`
说明：更新当前用户设置。

请求/响应参数：

```ts
interface SettingsState {
  defaultInterviewMode: '综合模拟' | '专项训练' | '项目拷打';
  followUpIntensity: '低' | '中' | '高';
  scoringStrictness: '宽松' | '标准' | '严格';
  modelName: string;
  streamOutput: boolean;
}
```

更新接口请求建议：

```ts
type UpdateSettingsParams = Partial<SettingsState>;
```

---

## 3. 错误码建议

建议使用统一业务码，`HTTP Status` 保留标准语义，业务失败通过 `code` 表达。

```ts
enum ApiErrorCode {
  SUCCESS = 0,
  BAD_REQUEST = 40000,
  UNAUTHORIZED = 40100,
  FORBIDDEN = 40300,
  NOT_FOUND = 40400,
  RATE_LIMITED = 42900,
  INTERNAL_ERROR = 50000,

  AUTH_INVALID_CREDENTIALS = 40101,
  AUTH_TOKEN_EXPIRED = 40102,

  QUESTION_NOT_FOUND = 40401,
  INTERVIEW_SESSION_NOT_FOUND = 40411,
  INTERVIEW_INVALID_STATE = 40011,
  REVIEW_REPORT_NOT_FOUND = 40421,
  PROJECT_POLISH_FAILED = 50031,
  SETTINGS_INVALID_VALUE = 40041,
}
```

错误响应示例：

```json
{
  "code": 40101,
  "message": "邮箱或密码错误",
  "data": null,
  "traceId": "trace-xxx"
}
```

建议：
- 登录失效统一使用 `40102`
- 参数校验失败统一使用 `40000`
- 资源不存在按领域细分，方便前端埋点和排查

## 4. TypeScript 类型定义建议

建议后端契约优先与前端现有类型对齐，并补一层“接口 DTO 类型”。

### 4.1 通用 DTO

```ts
interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
  traceId?: string;
}

interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
}
```

### 4.2 推荐目录

```text
frontend/src/types/
├─ api.ts
├─ user.ts
├─ dashboard.ts
├─ question-bank.ts
├─ interview.ts
├─ review.ts
├─ project.ts
└─ settings.ts
```

### 4.3 建议补充的后端返回字段

当前前端能跑通，但后续建议后端补这些结构化字段：

- `DashboardStatItem.key`
  便于前端识别统计项，而不是只依赖文案
- `RecentSessionItem.finishedAtTs`
  便于前端自己格式化时间
- `QuestionItem.updatedAt`
  便于后续排序和缓存
- `InterviewMessageItem.createdAt`
  当前只有 `time` 文本，真实接口建议返回标准时间戳
- `ProjectHistoryItem.createdAtTs`
  当前只有展示文案，不利于真实排序

## 5. mock 替换为真实接口的迁移说明

当前前端同时存在两层：

- `src/services/modules/*`
  当前承载 mock 逻辑
- `src/api/*`
  当前承载真实接口签名骨架

### 5.1 推荐迁移步骤

1. 保留 `src/api/*` 作为真实接口入口，不改文件职责。
2. 将 `src/services/modules/*` 逐步改成“数据源适配层”。
   - 开发期调用 mock
   - 联调期切换调用 `src/api/*`
3. store 和 composable 不直接依赖 axios，只依赖 `services/modules/*`。
4. 联调完成后，如果 mock 不再需要：
   - 可以删除 `services/modules/*` 中的假实现
   - 或保留为开发开关，如 `USE_MOCK=true`

### 5.2 推荐改造方式

以登录为例：

```ts
// services/modules/auth.ts
import { authApi } from '@/api/auth';

export const authService = {
  login(params: LoginParams) {
    if (useMock) {
      return mockLogin(params);
    }

    return authApi.login(params).then((res) => res.data);
  },
};
```

说明：
- 页面、store 不需要感知 mock 还是 real
- 切换真实接口时，只改 service 层
- 这也是当前前端结构最平滑的演进方式

### 5.3 联调时需要优先替换的模块顺序

建议顺序：

1. 登录认证
2. 设置
3. Dashboard 聚合接口
4. 题库查询
5. 面试会话
6. 复盘报告
7. 项目包装

原因：
- 登录和设置依赖最轻，适合先验证请求封装和 token 注入
- Dashboard/题库是典型查询接口，适合验证聚合结构与分页
- 面试会话和项目包装后续可能涉及流式、异步生成，应该后置

## 6. 实施建议

- 后端第一版优先保证字段名稳定，不要过早做复杂嵌套
- 时间字段建议统一补标准时间戳，前端自行格式化
- 枚举值建议固定，不要混用展示文案和内部编码
- 对聚合型页面，优先提供页面级聚合接口，避免前端一次加载拆成多次请求