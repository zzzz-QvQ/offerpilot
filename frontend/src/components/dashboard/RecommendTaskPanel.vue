<template>
  <el-card class="dashboard-panel" shadow="never">
    <template #header>
      <div class="dashboard-panel__header">
        <span>推荐训练任务</span>
      </div>
    </template>

    <div class="recommend-task-panel">
      <div v-for="task in items" :key="task.id" class="recommend-task-panel__item">
        <div class="recommend-task-panel__top">
          <h3 class="recommend-task-panel__title">{{ task.title }}</h3>
          <el-tag size="small" :type="tagTypeMap[task.level]">{{ task.level }}</el-tag>
        </div>
        <p class="recommend-task-panel__desc">{{ task.description }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { RecommendTaskItem } from '@/types/dashboard';

const tagTypeMap: Record<RecommendTaskItem['level'], 'danger' | 'warning' | 'success'> = {
  '高优先级': 'danger',
  '中优先级': 'warning',
  '巩固建议': 'success',
};

defineProps<{
  items: RecommendTaskItem[];
}>();
</script>

<style scoped lang="scss">
.dashboard-panel {
  height: 100%;
}

.dashboard-panel__header {
  font-weight: 600;
}

.recommend-task-panel {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.recommend-task-panel__item {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fcfdff;
}

.recommend-task-panel__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.recommend-task-panel__title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.recommend-task-panel__desc {
  margin: 10px 0 0;
  font-size: 13px;
  line-height: 1.7;
  color: var(--color-text-secondary);
}
</style>