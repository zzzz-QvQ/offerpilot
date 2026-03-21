# Frontend Refactor Notes

## 审查结论

### 发现的问题

1. 请求层与用户状态层重复维护存储 key。
   - `frontend/src/utils/request.ts`
   - `frontend/src/stores/modules/user.ts`
   这会让 token 存储规则在后续接真实 API、切换持久化策略时产生双点维护成本。

2. 图表组件存在明显重复逻辑。
   - `frontend/src/components/dashboard/WeakPointChart.vue`
   - `frontend/src/components/review/ReviewRadarChart.vue`
   - `frontend/src/components/review/ReviewWeaknessBarChart.vue`
   - `frontend/src/components/review/ReviewTrendLineChart.vue`
   这几处都重复了 ECharts 初始化、resize 监听、销毁清理逻辑，后续一旦要加主题切换、容器刷新或懒加载，维护成本会放大。

3. 当前 mock 服务层和未来真实 API 层并存，但边界还没有完全统一。
   - `frontend/src/services/modules/*`
   - `frontend/src/api/*`
   当前结构可用，但后续建议让 `services` 变成“场景适配层”，只负责切换 mock / real api，而不是继续增长为第二套接口定义。

4. `SettingsView` 仍承担了完整表单展示与状态映射。
   - `frontend/src/views/settings/SettingsView.vue`
   这个页面目前规模尚可，但如果后续继续增加模型参数、偏好项和实验配置，建议拆成 `ModelSettingsForm` 与 `InterviewPreferenceForm` 两个组件。

5. composables 覆盖面还不均衡。
   - 已有 `useQuestionBank.ts`、`useInterviewSession.ts`、`useProjectPolish.ts`
   - `Dashboard`、`Review`、`Settings` 仍以“页面直接读 store”为主
   当前不算错误，但如果这些页面再增加异步状态、筛选或局部交互，建议及时抽出 composable，避免页面逐步变重。

## 本次增量优化

1. 抽离了统一存储 key 常量。
   - 新增 `frontend/src/constants/storage.ts`
   - `request.ts` 和 `user store` 改为复用同一来源

2. 抽离了通用图表 composable。
   - 新增 `frontend/src/composables/useEChart.ts`
   - 复用到 Dashboard 与 Review 的 4 个图表组件中

3. 保持现有页面、store 和 mock 数据结构不变。
   - 这次只消除了重复逻辑，没有改动业务行为和页面结构

## 后续建议

1. 为 `settings` 页面补一个轻量 composable 或拆分表单组件。
2. 让 `services/modules/*` 逐步转成 mock adapter，统一调用 `src/api/*` 的接口签名。
3. 将常用领域常量继续从 store 中抽离，例如面试模式、评分等级、问题频率枚举。