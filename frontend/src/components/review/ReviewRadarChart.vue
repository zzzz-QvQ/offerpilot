<template>
  <el-card class="review-chart-card" shadow="never">
    <template #header>
      <div class="review-chart-card__header">多维评分雷达图</div>
    </template>
    <el-empty v-if="!items.length" description="暂无评分数据" />
    <div v-else ref="chartRef" class="review-chart"></div>
  </el-card>
</template>

<script setup lang="ts">
import { watch } from 'vue';

import { useEChart } from '@/composables/useEChart';
import type { ReviewRadarItem } from '@/types/review';

const props = defineProps<{
  items: ReviewRadarItem[];
}>();

const { chartRef, setChartOption } = useEChart();

const renderChart = () => {
  if (!props.items.length) {
    return;
  }

  void setChartOption({
    tooltip: {},
    radar: {
      radius: '62%',
      indicator: props.items.map((item) => ({ name: item.name, max: item.fullMark })),
      splitArea: { areaStyle: { color: ['#f8fbff', '#f2f7ff'] } },
    },
    series: [{
      type: 'radar',
      data: [{
        value: props.items.map((item) => item.score),
        areaStyle: { color: 'rgba(37, 99, 235, 0.18)' },
        lineStyle: { color: '#2563eb' },
        itemStyle: { color: '#2563eb' },
      }],
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