<template>
  <div class="interview-view">
    <div class="interview-view__grid">
      <section class="interview-view__left">
        <InterviewQuestionList
          :items="questions"
          :active-id="currentQuestion?.id ?? ''"
          :loading="loading"
          @select="selectQuestion"
        />
      </section>

      <section class="interview-view__center">
        <InterviewChatPanel
          :items="messages"
          :current-question="currentQuestion"
          :loading="loading"
        />
        <InterviewInputBox :submitting="submitting" @submit="sendMessage" />
      </section>

      <section class="interview-view__right">
        <InterviewStatusPanel :status-items="statusItems" :score-items="scoreItems" />
        <KnowledgeReferenceCard :items="knowledgeHits" />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import InterviewChatPanel from '@/components/interview/InterviewChatPanel.vue';
import InterviewInputBox from '@/components/interview/InterviewInputBox.vue';
import InterviewQuestionList from '@/components/interview/InterviewQuestionList.vue';
import InterviewStatusPanel from '@/components/interview/InterviewStatusPanel.vue';
import KnowledgeReferenceCard from '@/components/interview/KnowledgeReferenceCard.vue';
import { useInterviewSession } from '@/composables/useInterviewSession';

const {
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
} = useInterviewSession();
</script>

<style scoped lang="scss">
.interview-view {
  min-height: calc(100vh - 144px);
}

.interview-view__grid {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr) 320px;
  gap: 16px;
  align-items: start;
}

.interview-view__left,
.interview-view__right {
  min-width: 0;
}

.interview-view__center {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.interview-view__right {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>