import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';

import { useInterviewStore } from '@/stores/modules/interview';

export const useInterviewSession = () => {
  const store = useInterviewStore();
  const {
    loading,
    submitting,
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

  onMounted(() => {
    if (!store.sessionId) {
      void store.fetchSession();
    }
  });

  return {
    loading,
    submitting,
    questions,
    messages,
    statusItems,
    scoreItems,
    knowledgeHits,
    currentQuestion,
    sendMessage,
    selectQuestion,
  };
};