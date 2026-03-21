import type { InterviewSessionSnapshot } from '@/types/interview';

const mockSnapshot: InterviewSessionSnapshot = {
  sessionId: 'mock-session-1',
  currentQuestionId: 'iq-1',
  questions: [
    { id: 'iq-1', title: '请介绍一下 Vue 3 响应式系统的核心机制。', status: '进行中' },
    { id: 'iq-2', title: '说说浏览器事件循环与微任务、宏任务。', status: '待开始' },
    { id: 'iq-3', title: '前端性能优化你通常怎么系统性回答？', status: '待开始' },
    { id: 'iq-4', title: 'TypeScript 泛型在项目中的典型应用场景有哪些？', status: '待开始' },
  ],
  messages: [
    {
      id: 'msg-1',
      role: 'interviewer',
      content: '先从响应式原理开始，请你用面试回答的方式概述 Vue 3 相比 Vue 2 的变化。',
      time: '10:00',
    },
    {
      id: 'msg-2',
      role: 'candidate',
      content: 'Vue 3 主要使用 Proxy 重写了响应式系统，相比 Vue 2 基于 Object.defineProperty 的方案，对新增属性、数组和深层对象场景支持更自然。',
      time: '10:01',
    },
    {
      id: 'msg-3',
      role: 'interviewer',
      content: '那你继续展开说一下 ref、reactive 以及依赖收集触发更新的过程。',
      time: '10:01',
    },
  ],
  statusItems: [
    { label: '会话状态', value: '模拟中' },
    { label: '当前题号', value: '第 1 / 4 题' },
    { label: '已回答轮次', value: '2 轮' },
    { label: '预计剩余时间', value: '18 分钟' },
  ],
  scoreItems: [
    { label: '表达完整度', score: 84 },
    { label: '知识准确性', score: 79 },
    { label: '结构清晰度', score: 81 },
  ],
  knowledgeHits: [
    { name: 'Proxy 与 Reflect', level: '高', summary: '已覆盖响应式核心实现差异。' },
    { name: 'ref / reactive', level: '高', summary: '已主动提到 API 使用边界。' },
    { name: 'effect / track / trigger', level: '中', summary: '提到依赖收集，但还可以更体系化。' },
    { name: 'computed / watch', level: '低', summary: '尚未展开和响应式系统的关联。' },
  ],
};

const fetchInterviewSession = async (): Promise<InterviewSessionSnapshot> => {
  await new Promise((resolve) => {
    window.setTimeout(resolve, 300);
  });

  return mockSnapshot;
};

const submitInterviewAnswer = async (content: string) => {
  await new Promise((resolve) => {
    window.setTimeout(resolve, 400);
  });

  return {
    interviewerReply: `收到你的回答：${content}。接下来请进一步说明响应式依赖是如何建立和触发更新的。`,
  };
};

export const interviewApi = {
  fetchInterviewSession,
  submitInterviewAnswer,
};