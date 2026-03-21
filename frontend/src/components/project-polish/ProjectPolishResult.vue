<template>
  <el-card shadow="never" class="project-polish-result">
    <template #header>
      <div class="project-polish-result__header">
        <span>生成结果展示区</span>
      </div>
    </template>

    <el-empty v-if="!result" description="请先填写项目信息并生成内容" />

    <div v-else class="project-polish-result__content">
      <section class="project-polish-result__section">
        <div class="project-polish-result__section-top">
          <h3>简历版项目描述</h3>
          <el-button text type="primary" @click="emit('copy', result.resumeDescription)">复制</el-button>
        </div>
        <p>{{ result.resumeDescription }}</p>
      </section>

      <section class="project-polish-result__section">
        <div class="project-polish-result__section-top">
          <h3>面试版项目描述</h3>
          <el-button text type="primary" @click="emit('copy', result.interviewDescription)">复制</el-button>
        </div>
        <p>{{ result.interviewDescription }}</p>
      </section>

      <section class="project-polish-result__section">
        <div class="project-polish-result__section-top">
          <h3>技术亮点</h3>
          <el-button text type="primary" @click="emit('copy', result.highlights.join('\n'))">复制</el-button>
        </div>
        <ul>
          <li v-for="item in result.highlights" :key="item">{{ item }}</li>
        </ul>
      </section>

      <section class="project-polish-result__section">
        <div class="project-polish-result__section-top">
          <h3>难点与解决方案</h3>
          <el-button text type="primary" @click="emit('copy', result.difficulties.join('\n'))">复制</el-button>
        </div>
        <ul>
          <li v-for="item in result.difficulties" :key="item">{{ item }}</li>
        </ul>
      </section>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { ProjectPolishOutput } from '@/types/project';

defineProps<{
  result: ProjectPolishOutput | null;
}>();

const emit = defineEmits<{
  copy: [text: string];
}>();
</script>

<style scoped lang="scss">
.project-polish-result {
  height: 100%;
}

.project-polish-result__header {
  font-weight: 600;
}

.project-polish-result__content {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.project-polish-result__section {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fcfdff;
}

.project-polish-result__section-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.project-polish-result__section h3 {
  margin: 0;
  font-size: 16px;
}

.project-polish-result__section p,
.project-polish-result__section ul {
  margin: 12px 0 0;
  line-height: 1.8;
  color: var(--color-text-secondary);
}

.project-polish-result__section ul {
  padding-left: 18px;
}
</style>