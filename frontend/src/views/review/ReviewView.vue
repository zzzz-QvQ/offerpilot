<template>
  <div class="review-view">
    <el-alert
      v-if="reviewStore.error"
      :title="reviewStore.error"
      type="error"
      show-icon
      closable="false"
    />

    <el-skeleton v-if="reviewStore.loading" :rows="10" animated />

    <el-empty
      v-else-if="!reviewStore.summaryItems.length && !reviewStore.historyItems.length"
      description="暂无复盘数据"
    />

    <template v-else>
      <el-row :gutter="16">
        <el-col :xs="24" :lg="8">
          <ReviewSummaryCard :score="reviewStore.overallScore" :items="reviewStore.summaryItems" />
        </el-col>
        <el-col :xs="24" :lg="16">
          <ReviewRadarChart :items="reviewStore.radarItems" />
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :xs="24" :xl="14">
          <ReviewWeaknessBarChart :items="reviewStore.weaknessItems" />
        </el-col>
        <el-col :xs="24" :xl="10">
          <ReviewSuggestionPanel :items="reviewStore.suggestionItems" />
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :span="24">
          <ReviewTrendLineChart :items="reviewStore.trendItems" />
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';

import ReviewRadarChart from '@/components/review/ReviewRadarChart.vue';
import ReviewSuggestionPanel from '@/components/review/ReviewSuggestionPanel.vue';
import ReviewSummaryCard from '@/components/review/ReviewSummaryCard.vue';
import ReviewTrendLineChart from '@/components/review/ReviewTrendLineChart.vue';
import ReviewWeaknessBarChart from '@/components/review/ReviewWeaknessBarChart.vue';
import { useReviewStore } from '@/stores/modules/review';

const reviewStore = useReviewStore();

onMounted(() => {
  if (!reviewStore.loading && !reviewStore.historyItems.length && !reviewStore.error) {
    void reviewStore.fetchReviewData();
  }
});
</script>

<style scoped lang="scss">
.review-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
