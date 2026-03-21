<template>
  <el-card class="review-chart-card" shadow="never">
    <template #header>
      <div class="review-chart-card__header">知识模块得分柱状图</div>
    </template>
    <el-empty v-if="!items.length" description="暂无薄弱点数据" />
    <div v-else ref="chartRef" class="review-chart"></div>
  </el-card>
</template>

<script setup lang="ts">
import { watch } from 'vue';

import { useEChart } from '@/composables/useEChart';
import type { ReviewWeaknessItem } from '@/types/review';

const props = defineProps<{
  items: ReviewWeaknessItem[];
}>();

const { chartRef, setChartOption } = useEChart();

const renderChart = () => {
  if (!props.items.length) {
    return;
  }

  void setChartOption({
    tooltip: { trigger: 'axis' },
    grid: { top: 20, left: 16, right: 16, bottom: 12, containLabel: true },
    xAxis: {
      type: 'category',
      data: props.items.map((item) => item.name),
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      max: 100,
      splitLine: { lineStyle: { color: '#e5e7eb' } },
    },
    series: [{
      type: 'bar',
      barWidth: 28,
      data: props.items.map((item) => item.score),
      itemStyle: { color: '#60a5fa', borderRadius: [8, 8, 0, 0] },
    }],
  });
};

watch(() => props.items, renderChart, { deep: true, immediate: true });
</script>

<style scoped lang="scss">
.review-chart-card {
  height: 100%;
}

.review-chart-card__header {
  font-weight: 600;
}

.review-chart {
  width: 100%;
  height: 320px;
}
</style>