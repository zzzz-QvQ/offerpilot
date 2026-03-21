import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { interviewApi } from '@/services/modules/interview';
import type {
  InterviewMessageItem,
  InterviewQuestionItem,
  InterviewScoreItem,
  InterviewStatusItem,
  KnowledgeHitItem,
} from '@/types/interview';

export const useInterviewStore = defineStore('interview', () => {
  const sessionId = ref('');
  const loading = ref(false);
  const submitting = ref(false);
  const currentQuestionId = ref('');
  const questions = ref<InterviewQuestionItem[]>([]);
  const messages = ref<InterviewMessageItem[]>([]);
  const statusItems = ref<InterviewStatusItem[]>([]);
  const scoreItems = ref<InterviewScoreItem[]>([]);
  const knowledgeHits = ref<KnowledgeHitItem[]>([]);

  const currentQuestion = computed(() => {
    return questions.value.find((item) => item.id === currentQuestionId.value) ?? null;
  });

  const fetchSession = async () => {
    loading.value = true;

    try {
      const data = await interviewApi.fetchInterviewSession();
      sessionId.value = data.sessionId;
      currentQuestionId.value = data.currentQuestionId;
      questions.value = data.questions;
      messages.value = data.messages;
      statusItems.value = data.statusItems;
      scoreItems.value = data.scoreItems;
      knowledgeHits.value = data.knowledgeHits;
    } finally {
      loading.value = false;
    }
  };

  const switchQuestion = (questionId: string) => {
    currentQuestionId.value = questionId;
    questions.value = questions.value.map((item) => ({
      ...item,
      status: item.id === questionId ? '进行中' : item.status === '进行中' ? '待开始' : item.status,
    }));
  };

  const appendCandidateMessage = (content: string) => {
    messages.value.push({
      id: `candidate-${Date.now()}`,
      role: 'candidate',
      content,
      time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
    });
  };

  const appendInterviewerMessage = (content: string) => {
    messages.value.push({
      id: `interviewer-${Date.now()}`,
      role: 'interviewer',
      content,
      time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
    });
  };

  const submitAnswer = async (content: string) => {
    if (!content.trim()) {
      return;
    }

    appendCandidateMessage(content);
    submitting.value = true;

    try {
      const result = await interviewApi.submitInterviewAnswer(content);
      appendInterviewerMessage(result.interviewerReply);
    } finally {
      submitting.value = false;
    }
  };

  return {
    loading,
    submitting,
    sessionId,
    currentQuestionId,
    questions,
    messages,
    statusItems,
    scoreItems,
    knowledgeHits,
    currentQuestion,
    fetchSession,
    switchQuestion,
    submitAnswer,
  };
});