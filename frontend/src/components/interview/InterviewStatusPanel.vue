<template>
  <el-card class="interview-panel" shadow="never">
    <template #header>
      <div class="interview-panel__header interview-status-panel__header">
        <span>状态与评分</span>
        <el-tag v-if="streaming" size="small" type="primary">实时中</el-tag>
      </div>
    </template>

    <div class="interview-status-panel">
      <div class="interview-status-panel__section">
        <div class="interview-status-panel__section-title">Agent 状态</div>
        <div class="interview-status-panel__timeline">
          <div
            v-for="stage in stageDefinitions"
            :key="stage.key"
            class="interview-status-panel__timeline-item"
            :class="timelineItemClass(stage.key)"
          >
            <div class="interview-status-panel__timeline-dot" />
            <div class="interview-status-panel__timeline-content">
              <div class="interview-status-panel__timeline-title">{{ stage.label }}</div>
              <div class="interview-status-panel__timeline-desc">{{ stage.description }}</div>
            </div>
          </div>
        </div>
      </div>

      <div class="interview-status-panel__section">
        <div class="interview-status-panel__section-title">会话状态</div>
        <el-empty v-if="!statusItems.length" description="当前还没有状态数据" />
        <div v-else class="interview-status-panel__status-list">
          <div v-for="item in statusItems" :key="item.label" class="interview-status-panel__status-item">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
          </div>
        </div>
      </div>

      <div class="interview-status-panel__section">
        <div class="interview-status-panel__section-title">实时评分</div>
        <el-empty v-if="!scoreItems.length" description="当前还没有评分结果" />
        <div v-else class="interview-status-panel__score-list">
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
import { computed } from 'vue';

import type { InterviewAgentStage, InterviewScoreItem, InterviewStatusItem } from '@/types/interview';

const props = defineProps<{
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
  streaming: boolean;
  agentStage: InterviewAgentStage;
}>();

const stageDefinitions: Array<{
  key: InterviewAgentStage;
  label: string;
  description: string;
}> = [
  { key: 'session_created', label: '会话已创建', description: '当前面试会话已经准备就绪。' },
  { key: 'generating_question', label: '生成问题', description: '正在组织当前轮次的提问内容。' },
  { key: 'retrieving_knowledge', label: '检索知识', description: '正在召回相关题库知识与参考信息。' },
  { key: 'scoring', label: '评分分析', description: '正在评估上一条回答的表现。' },
  { key: 'generating_followup', label: '生成追问', description: '正在生成反馈与下一条追问。' },
  { key: 'completed', label: '当前轮结束', description: '当前这一轮的输出已经完成。' },
];

const currentStageIndex = computed(() => {
  return stageDefinitions.findIndex((item) => item.key === props.agentStage);
});

const timelineItemClass = (stage: InterviewAgentStage) => {
  const stageIndex = stageDefinitions.findIndex((item) => item.key === stage);

  if (stageIndex < currentStageIndex.value) {
    return 'is-completed';
  }

  if (stageIndex === currentStageIndex.value) {
    return 'is-active';
  }

  return 'is-pending';
};
</script>

<style scoped lang="scss">
.interview-panel {
  height: 100%;
}

.interview-panel__header {
  font-weight: 600;
}

.interview-status-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
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

.interview-status-panel__timeline,
.interview-status-panel__status-list,
.interview-status-panel__score-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interview-status-panel__timeline-item,
.interview-status-panel__status-item,
.interview-status-panel__score-item {
  padding: 12px 14px;
  border-radius: 10px;
  background: #f8fafc;
}

.interview-status-panel__timeline-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  border: 1px solid transparent;
  transition: border-color 0.2s ease, background-color 0.2s ease, opacity 0.2s ease;
}

.interview-status-panel__timeline-item.is-pending {
  opacity: 0.52;
}

.interview-status-panel__timeline-item.is-active {
  border-color: rgba(37, 99, 235, 0.25);
  background: #eef4ff;
  box-shadow: 0 0 0 1px rgba(37, 99, 235, 0.06) inset;
}

.interview-status-panel__timeline-item.is-completed {
  background: #f0fdf4;
}

.interview-status-panel__timeline-dot {
  flex: 0 0 10px;
  width: 10px;
  height: 10px;
  margin-top: 6px;
  border-radius: 999px;
  background: #cbd5e1;
}

.interview-status-panel__timeline-item.is-active .interview-status-panel__timeline-dot {
  background: #2563eb;
  box-shadow: 0 0 0 6px rgba(37, 99, 235, 0.12);
}

.interview-status-panel__timeline-item.is-completed .interview-status-panel__timeline-dot {
  background: #16a34a;
}

.interview-status-panel__timeline-content {
  min-width: 0;
}

.interview-status-panel__timeline-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.interview-status-panel__timeline-desc {
  margin-top: 6px;
  line-height: 1.6;
  color: var(--color-text-secondary);
}

.interview-status-panel__status-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: var(--color-text-secondary);
}

.interview-status-panel__status-item strong,
.interview-status-panel__score-row strong {
  color: var(--color-text-primary);
}

.interview-status-panel__score-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  color: var(--color-text-secondary);
}
</style>
