import type { ProjectPolishFormData, ProjectPolishOutput } from '@/types/project';

const generateProjectPolish = async (payload: ProjectPolishFormData): Promise<ProjectPolishOutput> => {
  await new Promise((resolve) => {
    window.setTimeout(resolve, 400);
  });

  return {
    resumeDescription: `${payload.projectName} 是一个基于 ${payload.techStack} 构建的前端项目，我在项目中主要负责 ${payload.responsibility}，推动核心功能稳定交付并持续优化用户体验。`,
    interviewDescription: `在 ${payload.projectName} 中，我作为 ${payload.role}，围绕 ${payload.projectBackground} 展开建设，重点负责 ${payload.responsibility}。在推进过程中，我不仅完成了功能实现，还持续关注性能、可维护性和交付效率。`,
    highlights: [
      `围绕 ${payload.techStack} 完成核心模块拆分与组件沉淀，提升开发复用率。`,
      '通过明确状态流转与边界设计，降低复杂页面的联动成本。',
      '对关键交互链路进行优化，提升页面可用性与稳定性。',
    ],
    difficulties: [
      '复杂业务状态较多，容易造成组件耦合。通过状态分层和职责拆分降低维护成本。',
      '交付周期紧张。通过模块化开发与优先级梳理保证关键需求按时上线。',
      '项目表达容易停留在功能描述。通过补充指标、方案对比和落地结果增强说服力。',
    ],
    followups: [
      {
        question: '这个项目里你最有代表性的技术决策是什么？',
        answer: '我会重点讲状态分层和组件边界设计，因为它同时影响开发效率、可维护性和后续功能扩展。',
      },
      {
        question: '如果让你再做一次，你会优先重构哪个点？',
        answer: '我会优先补齐通用能力层，比如请求封装、领域状态和公共组件规范，这样后续扩展成本更低。',
      },
      {
        question: '你如何证明这个项目确实带来了结果？',
        answer: '我会补充上线效果、性能指标、交付效率提升或团队协作收益，而不是只讲做了什么。',
      },
    ],
  };
};

export const projectApi = {
  generateProjectPolish,
};