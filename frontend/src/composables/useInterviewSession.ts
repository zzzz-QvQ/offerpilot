import { computed, onBeforeUnmount, onMounted } from 'vue';
import { storeToRefs } from 'pinia';

import { interviewApi } from '@/api/interview';
import { useSSEStream } from '@/composables/useSSEStream';
import { useInterviewStore } from '@/stores/modules/interview';
import type {
  DeltaEvent,
  DoneEvent,
  ReferenceEvent,
  StateEvent,
} from '@/types/interview';

export const useInterviewSession = () => {
  const store = useInterviewStore();
  const {
    loading,
    submitting,
    finishing,
    isStreaming,
    error,
    streamError,
    streamStoppedByUser,
    agentStage,
    questions,
    messages,
    statusItems,
    scoreItems,
    knowledgeHits,
    currentQuestion,
  } = storeToRefs(store);

  const stream = useSSEStream({
    onDelta: (event) => {
      const payload = (event as DeltaEvent).payload;
      if (typeof payload.content === 'string') {
        store.appendStreamDelta(payload.content);
      }
    },
    onState: (event) => {
      store.applyStreamState((event as StateEvent).payload);
    },
    onReference: (event) => {
      store.applyStreamReference((event as ReferenceEvent).payload);
    },
    onDone: (event) => {
      store.finishStreamingReply((event as DoneEvent).payload);
    },
    onError: (message) => {
      store.setStreamError(message);
      store.interruptStreaming('error');
    },
  });

  const displayError = computed(() => error.value || streamError.value);

  const sendMessage = async (content: string) => {
    if (store.isStreaming) {
      return;
    }

    const accepted = await store.submitAnswer(content);
    if (!accepted || !store.sessionId) {
      return;
    }

    const started = stream.start(interviewApi.getStreamUrl(store.sessionId));
    if (!started) {
      store.setStreamError('当前会话已经存在一个进行中的流式连接。');
      store.interruptStreaming('error');
    }
  };

  const stopStreaming = () => {
    stream.stop();
    store.interruptStreaming('manual');
  };

  const selectQuestion = (questionId: string) => {
    store.switchQuestion(questionId);
  };

  const finishInterview = async () => {
    if (store.isStreaming) {
      stopStreaming();
    }
    await store.finishSession();
  };

  onMounted(() => {
    store.clearStreamFeedback();
    if (!store.sessionId && !store.loading) {
      void store.fetchSession();
    }
  });

  onBeforeUnmount(() => {
    if (store.isStreaming) {
      stopStreaming();
    }
    store.clearStreamFeedback();
  });

  return {
    loading,
    submitting,
    finishing,
    isStreaming,
    displayError,
    streamStoppedByUser,
    agentStage,
    questions,
    messages,
    statusItems,
    scoreItems,
    knowledgeHits,
    currentQuestion,
    sendMessage,
    stopStreaming,
    selectQuestion,
    finishInterview,
  };
};
