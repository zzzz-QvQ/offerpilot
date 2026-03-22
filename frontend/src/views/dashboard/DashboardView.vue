<template>
  <div class="dashboard-view">
    <el-alert
      v-if="dashboardStore.error"
      :title="dashboardStore.error"
      type="error"
      show-icon
      :closable="false"
    />

    <el-skeleton v-if="dashboardStore.loading" :rows="10" animated />

    <template v-else>
      <el-empty
        v-if="!dashboardStore.statCards.length && !dashboardStore.recentSessions.length && !dashboardStore.weakPoints.length"
        description="暂无 Dashboard 数据"
      />

      <template v-else>
        <el-row :gutter="16" class="dashboard-view__stats">
          <el-col v-for="item in dashboardStore.statCards" :key="item.label" :xs="24" :sm="12" :lg="6">
            <DashboardStatCard :label="item.label" :value="item.value" :hint="item.hint" />
          </el-col>
        </el-row>

        <el-row :gutter="16" class="dashboard-view__content">
          <el-col :xs="24" :xl="14">
            <RecentSessionList :items="dashboardStore.recentSessions" :loading="dashboardStore.loading" />
          </el-col>
          <el-col :xs="24" :xl="10">
            <RecommendTaskPanel :items="dashboardStore.recommendTasks" :loading="dashboardStore.loading" />
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="24">
            <WeakPointChart :items="dashboardStore.weakPoints" :loading="dashboardStore.loading" />
          </el-col>
        </el-row>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, nextTick } from 'vue';

import DashboardStatCard from '@/components/dashboard/DashboardStatCard.vue';
import RecentSessionList from '@/components/dashboard/RecentSessionList.vue';
import RecommendTaskPanel from '@/components/dashboard/RecommendTaskPanel.vue';
import WeakPointChart from '@/components/dashboard/WeakPointChart.vue';
import { useDashboardStore } from '@/stores/modules/dashboard';

const dashboardStore = useDashboardStore();

onMounted(async () => {
  await nextTick();

  if (!dashboardStore.loading && !dashboardStore.statCards.length && !dashboardStore.error) {
    void dashboardStore.fetchDashboardData();
  }
});
</script>

<style scoped lang="scss">
.dashboard-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dashboard-view__stats,
.dashboard-view__content {
  margin-bottom: 0;
}
</style>
