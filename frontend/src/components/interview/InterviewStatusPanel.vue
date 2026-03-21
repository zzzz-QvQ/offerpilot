<template>
  <el-card class="interview-panel" shadow="never">
    <template #header>
      <div class="interview-panel__header">
        <span>状态区与评分区</span>
      </div>
    </template>

    <div class="interview-status-panel">
      <div class="interview-status-panel__section">
        <div class="interview-status-panel__section-title">会话状态</div>
        <div class="interview-status-panel__status-list">
          <div v-for="item in statusItems" :key="item.label" class="interview-status-panel__status-item">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
          </div>
        </div>
      </div>

      <div class="interview-status-panel__section">
        <div class="interview-status-panel__section-title">即时评分</div>
        <div class="interview-status-panel__score-list">
          <div v-for="item in scoreItems" :key="item.label" class="interview-status-panel__score-item">
            <div class="interview-status-panel__score-row">
              <span>{{ item.label }}</span>
              <strong>{{ item.score }}</strong>
            </div>
            <el-progress :percentage="item.score" :show-text="false" :stroke-width="8" />
          </div>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { InterviewScoreItem, InterviewStatusItem } from '@/types/interview';

defineProps<{
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
}>();
</script>

<style scoped lang="scss">
.interview-panel {
  height: 100%;
}

.interview-panel__header {
  font-weight: 600;
}

.interview-status-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.interview-status-panel__section-title {
  margin-bottom: 12px;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.interview-status-panel__status-list,
.interview-status-panel__score-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interview-status-panel__status-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-radius: 10px;
  background: #f8fafc;
  color: var(--color-text-secondary);
}

.interview-status-panel__status-item strong,
.interview-status-panel__score-row strong {
  color: var(--color-text-primary);
}

.interview-status-panel__score-item {
  padding: 12px 14px;
  border-radius: 10px;
  background: #f8fafc;
}

.interview-status-panel__score-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  color: var(--color-text-secondary);
}
</style>