<template>
  <el-card class="interview-panel" shadow="never">
    <template #header>
      <div class="interview-panel__header">
        <span>知识命中区</span>
      </div>
    </template>

    <div class="knowledge-reference-card">
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
  '高': 'success',
  '中': 'warning',
  '低': 'danger',
};

defineProps<{
  items: KnowledgeHitItem[];
}>();
</script>

<style scoped lang="scss">
.interview-panel {
  height: 100%;
}

.interview-panel__header {
  font-weight: 600;
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