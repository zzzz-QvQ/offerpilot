<template>
  <el-card class="interview-panel" shadow="never">
    <template #header>
      <div class="interview-panel__header knowledge-reference-card__header">
        <span>Knowledge References</span>
        <el-tag v-if="streaming" size="small" type="info">Updating</el-tag>
      </div>
    </template>

    <el-empty v-if="!items.length" description="No knowledge references yet." />

    <div v-else class="knowledge-reference-card">
      <div v-for="item in items" :key="item.questionId" class="knowledge-reference-card__item">
        <div class="knowledge-reference-card__top">
          <div class="knowledge-reference-card__meta">
            <strong>{{ item.title }}</strong>
            <span class="knowledge-reference-card__category">{{ item.category }}</span>
          </div>
          <el-tag size="small" type="success">{{ formatScore(item.score) }}</el-tag>
        </div>

        <div class="knowledge-reference-card__id">Question ID: {{ item.questionId }}</div>
        <p>{{ item.snippet }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { KnowledgeReferenceItem } from '@/types/interview';

defineProps<{
  items: KnowledgeReferenceItem[];
  streaming: boolean;
}>();

const formatScore = (score: number) => {
  if (!Number.isFinite(score) || score <= 0) {
    return 'Match';
  }

  return score <= 1 ? score.toFixed(3) : Math.round(score).toString();
};
</script>

<style scoped lang="scss">
.interview-panel {
  height: 100%;
}

.interview-panel__header {
  font-weight: 600;
}

.knowledge-reference-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.knowledge-reference-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.knowledge-reference-card__item {
  padding: 14px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fcfdff;
}

.knowledge-reference-card__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.knowledge-reference-card__meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.knowledge-reference-card__meta strong {
  color: var(--color-text-primary);
}

.knowledge-reference-card__category,
.knowledge-reference-card__id,
.knowledge-reference-card__item p {
  color: var(--color-text-secondary);
}

.knowledge-reference-card__category,
.knowledge-reference-card__id {
  font-size: 12px;
}

.knowledge-reference-card__id {
  margin-top: 8px;
}

.knowledge-reference-card__item p {
  margin: 10px 0 0;
  line-height: 1.7;
}
</style>