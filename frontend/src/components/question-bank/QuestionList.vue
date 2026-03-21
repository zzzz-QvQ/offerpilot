<template>
  <el-card shadow="never" class="question-list-card">
    <template #header>
      <div class="question-list-card__header">
        <span>题目列表</span>
        <span class="question-list-card__count">共 {{ items.length }} 题</span>
      </div>
    </template>

    <el-empty v-if="!items.length && !loading" description="暂无匹配题目" />
    <el-skeleton v-else-if="loading" :rows="6" animated />

    <div v-else class="question-list">
      <div v-for="item in items" :key="item.id" class="question-list__item" @click="emit('view-detail', item.id)">
        <div class="question-list__top">
          <h3 class="question-list__title">{{ item.title }}</h3>
          <el-button
            text
            class="question-list__favorite-button"
            :loading="favoriteLoadingId === item.id"
            @click.stop="emit('toggle-favorite', item)"
          >
            <el-icon class="question-list__favorite" :class="{ 'is-active': item.isFavorite }">
              <StarFilled />
            </el-icon>
          </el-button>
        </div>

        <div class="question-list__meta">
          <el-tag size="small" effect="plain">{{ item.category }}</el-tag>
          <el-tag size="small" :type="item.frequency === '高频' ? 'danger' : 'warning'">{{ item.frequency }}</el-tag>
          <el-tag size="small" type="info">{{ item.difficulty }}</el-tag>
        </div>

        <p class="question-list__summary">{{ item.summary }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { StarFilled } from '@element-plus/icons-vue';

import type { QuestionItem } from '@/types/question-bank';

defineProps<{
  items: QuestionItem[];
  loading: boolean;
  favoriteLoadingId?: string;
}>();

const emit = defineEmits<{
  'view-detail': [questionId: string];
  'toggle-favorite': [question: QuestionItem];
}>();
</script>

<style scoped lang="scss">
.question-list-card {
  height: 100%;
}

.question-list-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.question-list-card__count {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.question-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.question-list__item {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  cursor: pointer;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.question-list__item:hover {
  border-color: rgba(37, 99, 235, 0.28);
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.06);
  transform: translateY(-1px);
}

.question-list__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.question-list__title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.question-list__favorite-button {
  padding: 0;
}

.question-list__favorite {
  margin-top: 2px;
  color: #cbd5e1;
}

.question-list__favorite.is-active {
  color: #f59e0b;
}

.question-list__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.question-list__summary {
  margin: 12px 0 0;
  color: var(--color-text-secondary);
  line-height: 1.7;
}
</style>