<template>
  <div class="project-polish-view">
    <div class="project-polish-view__grid">
      <section class="project-polish-view__left">
        <ProjectInputForm :model-value="form" :loading="loading" @submit="submitForm" />
        <ProjectHistoryList :items="historyList" />
      </section>

      <section class="project-polish-view__right">
        <ProjectPolishResult :result="result" @copy="copyText" />
        <FollowupQuestionPanel :items="result?.followups ?? []" />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import FollowupQuestionPanel from '@/components/project-polish/FollowupQuestionPanel.vue';
import ProjectHistoryList from '@/components/project-polish/ProjectHistoryList.vue';
import ProjectInputForm from '@/components/project-polish/ProjectInputForm.vue';
import ProjectPolishResult from '@/components/project-polish/ProjectPolishResult.vue';
import { useProjectPolish } from '@/composables/useProjectPolish';

const { loading, form, result, historyList, submitForm, copyText } = useProjectPolish();
</script>

<style scoped lang="scss">
.project-polish-view__grid {
  display: grid;
  grid-template-columns: 420px minmax(0, 1fr);
  gap: 16px;
  align-items: start;
}

.project-polish-view__left,
.project-polish-view__right {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}
</style>