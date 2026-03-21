import type { QuestionItem } from '@/types/question-bank';

const mockQuestions: QuestionItem[] = [
  {
    id: 'q-1',
    title: 'Vue 3 中 ref 和 reactive 的区别是什么？',
    category: 'Vue',
    frequency: '高频',
    isFavorite: true,
    summary: '考察响应式系统理解与使用边界。',
    difficulty: '基础',
    answerPoints: ['响应式包装方式不同', '模板与 script 中使用方式不同', '复杂对象场景下的选择'],
  },
  {
    id: 'q-2',
    title: '浏览器从输入 URL 到页面渲染发生了什么？',
    category: '浏览器',
    frequency: '高频',
    isFavorite: false,
    summary: '覆盖网络请求、解析渲染与 JS 执行链路。',
    difficulty: '进阶',
    answerPoints: ['DNS 与 TCP/TLS', 'HTTP 请求与响应', 'HTML 解析、CSSOM、Render Tree'],
  },
  {
    id: 'q-3',
    title: 'TypeScript 中 any、unknown、never 有什么区别？',
    category: 'TypeScript',
    frequency: '高频',
    isFavorite: true,
    summary: '考察类型系统边界与安全性。',
    difficulty: '基础',
    answerPoints: ['类型安全差异', '断言与收窄', 'never 的穷尽检查场景'],
  },
  {
    id: 'q-4',
    title: '前端性能优化通常从哪些维度切入？',
    category: '性能优化',
    frequency: '中频',
    isFavorite: false,
    summary: '围绕资源、渲染、缓存和监控。',
    difficulty: '进阶',
    answerPoints: ['资源体积与请求数量', '渲染阻塞与长任务', '缓存与性能监控'],
  },
  {
    id: 'q-5',
    title: 'HTTP 缓存强缓存和协商缓存的区别是什么？',
    category: '网络',
    frequency: '高频',
    isFavorite: false,
    summary: '考察缓存命中策略与响应头理解。',
    difficulty: '基础',
    answerPoints: ['Expires 与 Cache-Control', 'ETag 与 Last-Modified', '状态码 304'],
  },
  {
    id: 'q-6',
    title: 'Webpack 和 Vite 的核心差异是什么？',
    category: '工程化',
    frequency: '中频',
    isFavorite: true,
    summary: '聚焦开发阶段与构建阶段机制差异。',
    difficulty: '进阶',
    answerPoints: ['Dev Server 启动方式', 'ESM 与预构建', '生产构建差异'],
  },
];

const fetchQuestionList = async (): Promise<QuestionItem[]> => {
  await new Promise((resolve) => {
    window.setTimeout(resolve, 300);
  });

  return mockQuestions;
};

export const questionBankApi = {
  fetchQuestionList,
};