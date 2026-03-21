<template>
  <el-card class="dashboard-panel" shadow="never">
    <template #header>
      <div class="dashboard-panel__header">
        <span>最近训练记录</span>
      </div>
    </template>

    <el-skeleton v-if="loading" :rows="4" animated />
    <el-empty v-else-if="!items.length" description="暂无训练记录" />

    <div v-else class="recent-session-list">
      <div v-for="item in items" :key="item.id" class="recent-session-list__item">
        <div class="recent-session-list__main">
          <div class="recent-session-list__title-row">
            <h3 class="recent-session-list__title">{{ item.title }}</h3>
            <el-tag size="small" effect="plain">{{ item.tag }}</el-tag>
          </div>
          <div class="recent-session-list__meta">
            <span>{{ item.finishedAt }}</span>
            <span>{{ item.durationText }}</span>
          </div>
        </div>

        <div class="recent-session-list__score">{{ item.score }} 分</div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { RecentSessionItem } from '@/types/dashboard';

defineProps<{
  items: RecentSessionItem[];
  loading?: boolean;
}>();
</script>

<style scoped lang="scss">
.dashboard-panel {
  height: 100%;
}

.dashboard-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.recent-session-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.recent-session-list__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fbfdff;
}

.recent-session-list__main {
  min-width: 0;
  flex: 1;
}

.recent-session-list__title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.recent-session-list__title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.recent-session-list__meta {
  display: flex;
  gap: 16px;
  margin-top: 8px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.recent-session-list__score {
  flex-shrink: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--color-brand);
}
</style>