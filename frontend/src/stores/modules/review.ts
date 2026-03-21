import { defineStore } from 'pinia';
import { ref } from 'vue';

import { reviewApi } from '@/api/review';
import type {
  ReviewHistoryItem,
  ReviewRadarItem,
  ReviewReportResponse,
  ReviewSuggestionItem,
  ReviewSummaryItem,
  ReviewTrendItem,
  ReviewWeaknessItem,
} from '@/types/review';

const applyReport = (
  report: ReviewReportResponse,
  target: {
    overallScore: typeof overallScore;
    summaryItems: typeof summaryItems;
    radarItems: typeof radarItems;
    weaknessItems: typeof weaknessItems;
    trendItems: typeof trendItems;
    suggestionItems: typeof suggestionItems;
  },
) => {
  target.overallScore.value = report.overallScore;
  target.summaryItems.value = report.summaryItems;
  target.radarItems.value = report.radarItems;
  target.weaknessItems.value = report.weaknessItems;
  target.trendItems.value = report.trendItems;
  target.suggestionItems.value = report.suggestionItems;
};

const overallScore = ref(0);
const summaryItems = ref<ReviewSummaryItem[]>([]);
const radarItems = ref<ReviewRadarItem[]>([]);
const weaknessItems = ref<ReviewWeaknessItem[]>([]);
const trendItems = ref<ReviewTrendItem[]>([]);
const suggestionItems = ref<ReviewSuggestionItem[]>([]);

export const useReviewStore = defineStore('review', () => {
  const loading = ref(false);
  const error = ref('');
  const currentSessionId = ref('');
  const historyItems = ref<ReviewHistoryItem[]>([]);

  const clearReport = () => {
    overallScore.value = 0;
    summaryItems.value = [];
    radarItems.value = [];
    weaknessItems.value = [];
    trendItems.value = [];
    suggestionItems.value = [];
  };

  const fetchReport = async (sessionId: string) => {
    const response = await reviewApi.getReport(sessionId);
    currentSessionId.value = sessionId;
    applyReport(response.data, {
      overallScore,
      summaryItems,
      radarItems,
      weaknessItems,
      trendItems,
      suggestionItems,
    });
  };

  const fetchReviewData = async () => {
    loading.value = true;
    error.value = '';

    try {
      const historyResponse = await reviewApi.getHistory();
      historyItems.value = historyResponse.data;

      if (!historyItems.value.length) {
        clearReport();
        currentSessionId.value = '';
        return;
      }

      await fetchReport(historyItems.value[0].sessionId);
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载复盘数据失败';
      historyItems.value = [];
      clearReport();
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    error,
    currentSessionId,
    historyItems,
    overallScore,
    summaryItems,
    radarItems,
    weaknessItems,
    trendItems,
    suggestionItems,
    fetchReviewData,
    fetchReport,
  };
});