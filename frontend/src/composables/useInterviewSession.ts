import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';

import { useInterviewStore } from '@/stores/modules/interview';

export const useInterviewSession = () => {
  const store = useInterviewStore();
  const {
    loading,
    submitting,
    finishing,
    error,
    questions,
    messages,
    statusItems,
    scoreItems,
    knowledgeHits,
    currentQuestion,
  } = storeToRefs(store);

  const sendMessage = async (content: string) => {
    await store.submitAnswer(content);
  };

  const selectQuestion = (questionId: string) => {
    store.switchQuestion(questionId);
  };

  const finishInterview = async () => {
    await store.finishSession();
  };

  onMounted(() => {
    if (!store.sessionId && !store.loading) {
      void store.fetchSession();
    }
  });

  return {
    loading,
    submitting,
    finishing,
    error,
    questions,
    messages,
    statusItems,
    scoreItems,
    knowledgeHits,
    currentQuestion,
    sendMessage,
    selectQuestion,
    finishInterview,
  };
};