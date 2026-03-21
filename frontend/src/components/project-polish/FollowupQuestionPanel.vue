<template>
  <el-card shadow="never" class="followup-question-panel">
    <template #header>
      <div class="followup-question-panel__header">可能追问与建议回答</div>
    </template>

    <el-skeleton v-if="loading" :rows="4" animated />
    <el-empty v-else-if="!items.length" description="暂无追问内容" />

    <div v-else class="followup-question-panel__list">
      <div v-for="item in items" :key="item.question" class="followup-question-panel__item">
        <h3>{{ item.question }}</h3>
        <p>{{ item.answer }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { FollowupQuestionItem } from '@/types/project';

defineProps<{
  items: FollowupQuestionItem[];
  loading?: boolean;
}>();
</script>

<style scoped lang="scss">
.followup-question-panel {
  height: 100%;
}

.followup-question-panel__header {
  font-weight: 600;
}

.followup-question-panel__list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.followup-question-panel__item {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fcfdff;
}

.followup-question-panel__item h3 {
  margin: 0;
  font-size: 15px;
}

.followup-question-panel__item p {
  margin: 10px 0 0;
  line-height: 1.8;
  color: var(--color-text-secondary);
}
</style>