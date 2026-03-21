<template>
  <el-drawer :model-value="visible" title="题目详情" size="560px" @close="emit('close')">
    <template v-if="question">
      <div class="question-detail-drawer__header">
        <h2 class="question-detail-drawer__title">{{ question.title }}</h2>
        <div class="question-detail-drawer__meta">
          <el-tag effect="plain">{{ question.category }}</el-tag>
          <el-tag :type="question.frequency === '高频' ? 'danger' : 'warning'">{{ question.frequency }}</el-tag>
          <el-tag type="info">{{ question.difficulty }}</el-tag>
        </div>
      </div>

      <section class="question-detail-drawer__section">
        <h3>题目简介</h3>
        <p>{{ question.summary }}</p>
      </section>

      <section class="question-detail-drawer__section">
        <h3>答题要点</h3>
        <ul class="question-detail-drawer__points">
          <li v-for="point in question.answerPoints" :key="point">{{ point }}</li>
        </ul>
      </section>

      <section class="question-detail-drawer__section">
        <h3>收藏状态</h3>
        <p>{{ question.isFavorite ? '已收藏' : '未收藏' }}</p>
      </section>
    </template>

    <el-empty v-else description="请选择题目查看详情" />
  </el-drawer>
</template>

<script setup lang="ts">
import type { QuestionItem } from '@/types/question-bank';

defineProps<{
  visible: boolean;
  question: QuestionItem | null;
}>();

const emit = defineEmits<{
  close: [];
}>();
</script>

<style scoped lang="scss">
.question-detail-drawer__header {
  margin-bottom: 20px;
}

.question-detail-drawer__title {
  margin: 0;
  font-size: 22px;
  line-height: 1.5;
  color: var(--color-text-primary);
}

.question-detail-drawer__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.question-detail-drawer__section + .question-detail-drawer__section {
  margin-top: 24px;
}

.question-detail-drawer__section h3 {
  margin: 0 0 10px;
  font-size: 16px;
}

.question-detail-drawer__section p {
  margin: 0;
  line-height: 1.8;
  color: var(--color-text-secondary);
}

.question-detail-drawer__points {
  margin: 0;
  padding-left: 18px;
  color: var(--color-text-secondary);
}

.question-detail-drawer__points li + li {
  margin-top: 8px;
}
</style>