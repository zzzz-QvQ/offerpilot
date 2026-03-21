import { defineStore } from 'pinia';
import { ref } from 'vue';

import type {
  ReviewRadarItem,
  ReviewSuggestionItem,
  ReviewSummaryItem,
  ReviewTrendItem,
  ReviewWeaknessItem,
} from '@/types/review';

export const useReviewStore = defineStore('review', () => {
  const overallScore = ref(83);

  const summaryItems = ref<ReviewSummaryItem[]>([
    { label: '总体评分', value: '83 分', hint: '最近一轮模拟面试综合表现' },
    { label: '最佳模块', value: 'Vue / 组件设计', hint: '表达清晰，案例支撑较完整' },
    { label: '待加强模块', value: '性能优化', hint: '回答结构完整度偏弱' },
  ]);

  const radarItems = ref<ReviewRadarItem[]>([
    { name: '知识准确性', score: 84, fullMark: 100 },
    { name: '表达完整度', score: 79, fullMark: 100 },
    { name: '结构清晰度', score: 81, fullMark: 100 },
    { name: '追问应对', score: 76, fullMark: 100 },
    { name: '项目结合度', score: 88, fullMark: 100 },
  ]);

  const weaknessItems = ref<ReviewWeaknessItem[]>([
    { name: '性能优化', score: 62 },
    { name: '浏览器原理', score: 69 },
    { name: '工程化', score: 74 },
    { name: '网络基础', score: 78 },
    { name: 'TypeScript', score: 85 },
  ]);

  const trendItems = ref<ReviewTrendItem[]>([
    { date: '03-15', score: 72 },
    { date: '03-16', score: 75 },
    { date: '03-17', score: 79 },
    { date: '03-18', score: 77 },
    { date: '03-19', score: 81 },
    { date: '03-20', score: 83 },
  ]);

  const suggestionItems = ref<ReviewSuggestionItem[]>([
    {
      title: '补强性能优化回答模板',
      description: '建议按“发现问题、定位瓶颈、优化手段、效果验证”四段式组织答案。',
      type: '建议',
    },
    {
      title: '浏览器渲染链路不够完整',
      description: '在回答时常忽略 CSSOM、布局与合成线程之间的衔接。',
      type: '薄弱点',
    },
    {
      title: '加强追问场景演练',
      description: '建议增加多轮追问训练，重点提升临场补充和延展能力。',
      type: '建议',
    },
    {
      title: '工程化案例表达偏抽象',
      description: '需要补充真实项目中的构建优化、分包与 CI/CD 实践案例。',
      type: '薄弱点',
    },
  ]);

  return {
    overallScore,
    summaryItems,
    radarItems,
    weaknessItems,
    trendItems,
    suggestionItems,
  };
});