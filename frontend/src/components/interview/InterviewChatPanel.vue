<template>
  <el-card class="interview-panel interview-chat-panel" shadow="never">
    <template #header>
      <div class="interview-panel__header interview-chat-panel__header">
        <div>
          <div class="interview-chat-panel__title">Conversation</div>
          <div class="interview-chat-panel__subtitle">{{ currentQuestion?.title ?? 'Select a question to inspect the current round.' }}</div>
        </div>
        <el-tag v-if="streaming" size="small" type="primary">Streaming</el-tag>
      </div>
    </template>

    <el-skeleton v-if="loading" :rows="8" animated />
    <el-empty v-else-if="!items.length" description="No interview messages yet." />

    <div v-else class="interview-chat-panel__messages">
      <div
        v-for="item in items"
        :key="item.id"
        class="interview-chat-panel__message"
        :class="`is-${item.role}`"
      >
        <div class="interview-chat-panel__role">{{ item.role === 'interviewer' ? 'Interviewer' : 'You' }}</div>
        <div class="interview-chat-panel__bubble">
          {{ item.content || (item.isStreaming ? 'Generating response...' : '') }}
          <span v-if="item.isStreaming" class="interview-chat-panel__cursor" />
        </div>
        <div class="interview-chat-panel__time">{{ item.time }}</div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { InterviewMessageItem, InterviewQuestionItem } from '@/types/interview';

defineProps<{
  items: InterviewMessageItem[];
  currentQuestion: InterviewQuestionItem | null;
  loading: boolean;
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

.interview-chat-panel {
  display: flex;
  flex-direction: column;
}

.interview-chat-panel :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.interview-chat-panel__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.interview-chat-panel__title {
  font-size: 16px;
  font-weight: 600;
}

.interview-chat-panel__subtitle {
  margin-top: 8px;
  color: var(--color-text-secondary);
  line-height: 1.6;
}

.interview-chat-panel__messages {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 420px;
}

.interview-chat-panel__message {
  max-width: 88%;
}

.interview-chat-panel__message.is-candidate {
  align-self: flex-end;
}

.interview-chat-panel__role {
  margin-bottom: 6px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.interview-chat-panel__bubble {
  padding: 14px 16px;
  border-radius: 14px;
  line-height: 1.8;
  background: #f3f6fb;
  color: var(--color-text-primary);
}

.interview-chat-panel__message.is-candidate .interview-chat-panel__bubble {
  background: #2563eb;
  color: #ffffff;
}

.interview-chat-panel__time {
  margin-top: 6px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.interview-chat-panel__cursor {
  display: inline-block;
  width: 8px;
  height: 1em;
  margin-left: 4px;
  vertical-align: middle;
  background: currentColor;
  animation: blink 1s steps(1) infinite;
}

@keyframes blink {
  50% {
    opacity: 0;
  }
}
</style>