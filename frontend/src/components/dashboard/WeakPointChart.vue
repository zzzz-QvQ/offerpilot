<template>
  <el-card class="dashboard-panel" shadow="never">
    <template #header>
      <div class="dashboard-panel__header">
        <span>薄弱知识点图表</span>
      </div>
    </template>

    <div ref="chartRef" class="weak-point-chart"></div>
  </el-card>
</template>

<script setup lang="ts">
import { watch } from 'vue';

import { useEChart } from '@/composables/useEChart';
import type { WeakPointItem } from '@/types/dashboard';

const props = defineProps<{
  items: WeakPointItem[];
}>();

const { chartRef, setChartOption } = useEChart();

const renderChart = () => {
  void setChartOption({
    tooltip: {
      trigger: 'axis',
    },
    grid: {
      top: 16,
      left: 12,
      right: 12,
      bottom: 8,
      containLabel: true,
    },
    xAxis: {
      type: 'value',
      max: 100,
      splitLine: {
        lineStyle: {
          color: '#e5e7eb',
        },
      },
    },
    yAxis: {
      type: 'category',
      data: props.items.map((item) => item.name),
      axisTick: {
        show: false,
      },
    },
    series: [
      {
        type: 'bar',
        data: props.items.map((item) => item.value),
        barWidth: 18,
        itemStyle: {
          color: '#2563eb',
          borderRadius: [0, 8, 8, 0],
        },
        label: {
          show: true,
          position: 'right',
          color: '#111827',
        },
      },
    ],
  });
};

watch(() => props.items, renderChart, { deep: true, immediate: true });
</script>

<style scoped lang="scss">
.dashboard-panel {
  height: 100%;
}

.dashboard-panel__header {
  font-weight: 600;
}

.weak-point-chart {
  width: 100%;
  height: 320px;
}
</style>