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
      <div v-for="item in items" :key="item.name" class="knowledge-reference-card__item">
        <div class="knowledge-reference-card__top">
          <strong>{{ item.name }}</strong>
          <el-tag size="small" :type="tagTypeMap[item.level]">{{ item.level }}</el-tag>
        </div>
        <p>{{ item.summary }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { KnowledgeHitItem } from '@/types/interview';

const tagTypeMap: Record<KnowledgeHitItem['level'], 'danger' | 'warning' | 'success'> = {
  high: 'success',
  medium: 'warning',
  low: 'danger',
};

defineProps<{
  items: KnowledgeHitItem[];
  streaming: boolean;
}>();
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
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.knowledge-reference-card__item p {
  margin: 10px 0 0;
  line-height: 1.7;
  color: var(--color-text-secondary);
}
</style>