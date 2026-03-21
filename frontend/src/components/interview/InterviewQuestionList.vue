<template>
  <el-card class="interview-panel" shadow="never">
    <template #header>
      <div class="interview-panel__header">
        <span>问题列表</span>
      </div>
    </template>

    <el-skeleton v-if="loading" :rows="6" animated />
    <el-empty v-else-if="!items.length" description="暂无问题数据" />

    <div v-else class="interview-question-list">
      <div
        v-for="item in items"
        :key="item.id"
        class="interview-question-list__item"
        :class="{ 'is-active': item.id === activeId }"
        @click="emit('select', item.id)"
      >
        <div class="interview-question-list__title">{{ item.title }}</div>
        <el-tag size="small" effect="plain">{{ item.status }}</el-tag>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { InterviewQuestionItem } from '@/types/interview';

defineProps<{
  items: InterviewQuestionItem[];
  activeId: string;
  loading: boolean;
}>();

const emit = defineEmits<{
  select: [questionId: string];
}>();
</script>

<style scoped lang="scss">
.interview-panel {
  height: 100%;
}

.interview-panel__header {
  font-weight: 600;
}

.interview-question-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interview-question-list__item {
  padding: 14px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  cursor: pointer;
  background: #fcfdff;
  transition: all 0.2s ease;
}

.interview-question-list__item:hover,
.interview-question-list__item.is-active {
  border-color: rgba(37, 99, 235, 0.28);
  background: #eef4ff;
}

.interview-question-list__title {
  margin-bottom: 10px;
  font-size: 14px;
  line-height: 1.7;
  color: var(--color-text-primary);
}
</style>