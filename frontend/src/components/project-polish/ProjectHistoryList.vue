<template>
  <el-card shadow="never" class="project-history-list">
    <template #header>
      <div class="project-history-list__header">历史记录</div>
    </template>

    <el-skeleton v-if="loading" :rows="4" animated />
    <el-empty v-else-if="!items.length" description="暂无历史记录" />

    <div v-else class="project-history-list__list">
      <div
        v-for="item in items"
        :key="item.id"
        class="project-history-list__item"
        :class="{ 'is-active': item.id === activeId }"
        @click="emit('select', item.id)"
      >
        <h3>{{ item.projectName }}</h3>
        <div class="project-history-list__time">{{ item.createdAt }}</div>
        <p>{{ item.summary }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { ProjectHistoryItem } from '@/types/project';

defineProps<{
  items: ProjectHistoryItem[];
  loading?: boolean;
  activeId?: string;
}>();

const emit = defineEmits<{
  select: [id: string];
}>();
</script>

<style scoped lang="scss">
.project-history-list {
  height: 100%;
}

.project-history-list__header {
  font-weight: 600;
}

.project-history-list__list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.project-history-list__item {
  padding: 14px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fcfdff;
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease;
}

.project-history-list__item:hover,
.project-history-list__item.is-active {
  border-color: rgba(37, 99, 235, 0.28);
  background: #eef4ff;
}

.project-history-list__item h3 {
  margin: 0;
  font-size: 15px;
}

.project-history-list__time {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.project-history-list__item p {
  margin: 10px 0 0;
  line-height: 1.7;
  color: var(--color-text-secondary);
}
</style>