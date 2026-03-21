import { computed, onMounted } from 'vue';
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
      store.setStreamError('A live stream is already running for the current session.');
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
    if (!store.sessionId && !store.loading) {
      void store.fetchSession();
    }
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