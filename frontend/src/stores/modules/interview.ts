import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { interviewApi } from '@/api/interview';
import type {
  DoneEventPayload,
  InterviewAgentStage,
  InterviewMessageItem,
  InterviewQuestionItem,
  InterviewScoreItem,
  InterviewSessionDetail,
  InterviewStatusItem,
  KnowledgeHitItem,
  ReferenceEventPayload,
  StateEventPayload,
} from '@/types/interview';

const STREAMING_MESSAGE_ID = '__streaming_interviewer_message__';
const CANDIDATE_TEMP_PREFIX = 'candidate-';

export const useInterviewStore = defineStore('interview', () => {
  const sessionId = ref('');
  const loading = ref(false);
  const submitting = ref(false);
  const finishing = ref(false);
  const isStreaming = ref(false);
  const error = ref('');
  const streamError = ref('');
  const streamStoppedByUser = ref(false);
  const currentQuestionId = ref('');
  const agentStage = ref<InterviewAgentStage>('session_created');
  const questions = ref<InterviewQuestionItem[]>([]);
  const messages = ref<InterviewMessageItem[]>([]);
  const statusItems = ref<InterviewStatusItem[]>([]);
  const scoreItems = ref<InterviewScoreItem[]>([]);
  const knowledgeHits = ref<KnowledgeHitItem[]>([]);

  const currentQuestion = computed(() => {
    return questions.value.find((item) => item.id === currentQuestionId.value) ?? null;
  });

  const clearStreamFeedback = () => {
    streamError.value = '';
    streamStoppedByUser.value = false;
  };

  const applySessionDetail = (detail: InterviewSessionDetail) => {
    sessionId.value = detail.sessionId;
    currentQuestionId.value = detail.currentQuestionId;
    questions.value = detail.questions;
    messages.value = detail.messages;
    statusItems.value = detail.statusItems;
    scoreItems.value = detail.scoreItems;
    knowledgeHits.value = detail.knowledgeHits;
    agentStage.value = detail.agentStage ?? 'session_created';
  };

  const createSession = async () => {
    const response = await interviewApi.createSession();
    sessionId.value = response.data.sessionId;
    agentStage.value = 'session_created';
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
      error.value = err instanceof Error ? err.message : 'Failed to load interview session.';
      questions.value = [];
      messages.value = [];
      statusItems.value = [];
      scoreItems.value = [];
      knowledgeHits.value = [];
      agentStage.value = 'session_created';
    } finally {
      loading.value = false;
    }
  };

  const switchQuestion = (questionId: string) => {
    currentQuestionId.value = questionId;
  };

  const appendCandidateMessage = (content: string) => {
    messages.value = [
      ...messages.value,
      {
        id: `candidate-${Date.now()}`,
        role: 'candidate',
        content,
        time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      },
    ];
  };

  const rollbackPendingRound = () => {
    messages.value = messages.value.filter((item, index, list) => {
      if (item.id === STREAMING_MESSAGE_ID) {
        return false;
      }

      if (
        item.id.startsWith(CANDIDATE_TEMP_PREFIX)
        && item.role === 'candidate'
        && index === list.length - 1
      ) {
        return false;
      }

      return true;
    });
  };

  const startStreamingReply = () => {
    clearStreamFeedback();
    messages.value = messages.value.filter((item) => item.id !== STREAMING_MESSAGE_ID);
    messages.value = [
      ...messages.value,
      {
        id: STREAMING_MESSAGE_ID,
        role: 'interviewer',
        content: '',
        time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        isStreaming: true,
      },
    ];
    isStreaming.value = true;
    agentStage.value = 'generating_question';
  };

  const appendStreamDelta = (content: string) => {
    const target = messages.value.find((item) => item.id === STREAMING_MESSAGE_ID);
    if (!target) {
      startStreamingReply();
      return appendStreamDelta(content);
    }

    target.content += content;
  };

  const applyStreamState = (payload: StateEventPayload) => {
    if (payload.currentQuestionId) {
      currentQuestionId.value = payload.currentQuestionId;
    }
    if (payload.questions) {
      questions.value = payload.questions;
    }
    if (payload.statusItems) {
      statusItems.value = payload.statusItems;
    }
    if (payload.scoreItems) {
      scoreItems.value = payload.scoreItems;
    }
    if (payload.agentStage) {
      agentStage.value = payload.agentStage;
    }
  };

  const applyStreamReference = (payload: ReferenceEventPayload) => {
    if (payload.knowledgeHits) {
      knowledgeHits.value = payload.knowledgeHits;
    }
  };

  const finishStreamingReply = (payload?: DoneEventPayload) => {
    const target = messages.value.find((item) => item.id === STREAMING_MESSAGE_ID);

    if (payload?.message) {
      messages.value = messages.value.filter((item) => item.id !== STREAMING_MESSAGE_ID);
      messages.value = [...messages.value, { ...payload.message, isStreaming: false }];
    } else if (target) {
      target.isStreaming = false;
    }

    if (payload) {
      applyStreamState(payload);
      applyStreamReference(payload);
    }

    clearStreamFeedback();
    agentStage.value = payload?.agentStage ?? 'completed';
    isStreaming.value = false;
    submitting.value = false;
  };

  const interruptStreaming = (reason: 'manual' | 'error' = 'manual') => {
    const target = messages.value.find((item) => item.id === STREAMING_MESSAGE_ID);
    if (target && !target.content.trim()) {
      messages.value = messages.value.filter((item) => item.id !== STREAMING_MESSAGE_ID);
    } else if (target) {
      target.isStreaming = false;
    }

    if (reason === 'manual') {
      streamStoppedByUser.value = true;
      streamError.value = '';
    }

    agentStage.value = 'completed';
    isStreaming.value = false;
    submitting.value = false;
  };

  const submitAnswer = async (content: string) => {
    if (!content.trim() || !sessionId.value || isStreaming.value) {
      return false;
    }

    submitting.value = true;
    error.value = '';
    clearStreamFeedback();
    agentStage.value = 'generating_question';
    appendCandidateMessage(content);
    startStreamingReply();

    try {
      await interviewApi.submitMessage(sessionId.value, { content });
      return true;
    } catch (err) {
      rollbackPendingRound();
      error.value = err instanceof Error ? err.message : 'Failed to submit answer.';
      submitting.value = false;
      agentStage.value = 'completed';
      throw err;
    }
  };

  const finishSession = async () => {
    if (!sessionId.value) {
      return;
    }

    finishing.value = true;

    try {
      await interviewApi.finishSession(sessionId.value);
      agentStage.value = 'completed';
    } finally {
      finishing.value = false;
    }
  };

  const setError = (message: string) => {
    error.value = message;
  };

  const setStreamError = (message: string) => {
    streamError.value = message;
    streamStoppedByUser.value = false;
  };

  return {
    loading,
    submitting,
    finishing,
    isStreaming,
    error,
    streamError,
    streamStoppedByUser,
    sessionId,
    currentQuestionId,
    agentStage,
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
    applySessionDetail,
    appendStreamDelta,
    applyStreamState,
    applyStreamReference,
    finishStreamingReply,
    interruptStreaming,
    rollbackPendingRound,
    clearStreamFeedback,
    setError,
    setStreamError,
  };
});