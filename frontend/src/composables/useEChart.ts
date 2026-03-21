import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue';
import * as echarts from 'echarts';
import type { EChartsOption } from 'echarts';

export const useEChart = () => {
  const chartRef = ref<HTMLDivElement>();
  let chart: echarts.ECharts | null = null;

  const ensureChart = async () => {
    await nextTick();

    if (!chartRef.value) {
      return null;
    }

    if (!chart) {
      chart = echarts.init(chartRef.value);
    }

    return chart;
  };

  const setChartOption = async (option: EChartsOption) => {
    const instance = await ensureChart();

    if (!instance) {
      return;
    }

    instance.setOption(option);
  };

  const resizeChart = () => {
    chart?.resize();
  };

  onMounted(() => {
    window.addEventListener('resize', resizeChart);
  });

  onBeforeUnmount(() => {
    window.removeEventListener('resize', resizeChart);
    chart?.dispose();
    chart = null;
  });

  return {
    chartRef,
    setChartOption,
    resizeChart,
  };
};