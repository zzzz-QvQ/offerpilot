import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { interviewApi } from '@/api/interview';
import type {
  InterviewMessageItem,
  InterviewQuestionItem,
  InterviewScoreItem,
  InterviewSessionDetail,
  InterviewStatusItem,
  KnowledgeHitItem,
} from '@/types/interview';

export const useInterviewStore = defineStore('interview', () => {
  const sessionId = ref('');
  const loading = ref(false);
  const submitting = ref(false);
  const finishing = ref(false);
  const error = ref('');
  const currentQuestionId = ref('');
  const questions = ref<InterviewQuestionItem[]>([]);
  const messages = ref<InterviewMessageItem[]>([]);
  const statusItems = ref<InterviewStatusItem[]>([]);
  const scoreItems = ref<InterviewScoreItem[]>([]);
  const knowledgeHits = ref<KnowledgeHitItem[]>([]);

  const currentQuestion = computed(() => {
    return questions.value.find((item) => item.id === currentQuestionId.value) ?? null;
  });

  const applySessionDetail = (detail: InterviewSessionDetail) => {
    sessionId.value = detail.sessionId;
    currentQuestionId.value = detail.currentQuestionId;
    questions.value = detail.questions;
    messages.value = detail.messages;
    statusItems.value = detail.statusItems;
    scoreItems.value = detail.scoreItems;
    knowledgeHits.value = detail.knowledgeHits;
  };

  const createSession = async () => {
    const response = await interviewApi.createSession();
    sessionId.value = response.data.sessionId;
    return response.data.sessionId;
  };

  const fetchSession = async () => {
    loading.value = true;
    error.value = '';

    try {
      if (!sessionId.value) {
        await createSession();
      }

      const response = await interviewApi.getSession(sessionId.value);
      applySessionDetail(response.data);
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载面试会话失败';
      questions.value = [];
      messages.value = [];
      statusItems.value = [];
      scoreItems.value = [];
      knowledgeHits.value = [];
    } finally {
      loading.value = false;
    }
  };

  const switchQuestion = (questionId: string) => {
    currentQuestionId.value = questionId;
  };

  const submitAnswer = async (content: string) => {
    if (!content.trim() || !sessionId.value) {
      return;
    }

    submitting.value = true;
    error.value = '';

    try {
      const response = await interviewApi.submitMessage(sessionId.value, { content });
      const result = response.data;
      currentQuestionId.value = result.currentQuestionId;
      questions.value = result.questions;
      messages.value = [...messages.value, result.message];
      statusItems.value = result.statusItems;
      scoreItems.value = result.scoreItems;
      knowledgeHits.value = result.knowledgeHits;
    } catch (err) {
      error.value = err instanceof Error ? err.message : '提交回答失败';
      throw err;
    } finally {
      submitting.value = false;
    }
  };

  const finishSession = async () => {
    if (!sessionId.value) {
      return;
    }

    finishing.value = true;

    try {
      await interviewApi.finishSession(sessionId.value);
    } finally {
      finishing.value = false;
    }
  };

  return {
    loading,
    submitting,
    finishing,
    error,
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
    finishSession,
  };
});