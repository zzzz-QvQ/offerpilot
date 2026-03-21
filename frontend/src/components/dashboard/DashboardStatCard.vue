<template>
  <el-card class="dashboard-stat-card" shadow="never">
    <div class="dashboard-stat-card__label">{{ label }}</div>
    <div class="dashboard-stat-card__value">{{ value }}</div>
    <div class="dashboard-stat-card__hint">{{ hint }}</div>
  </el-card>
</template>

<script setup lang="ts">
defineProps<{
  label: string;
  value: string;
  hint: string;
}>();
</script>

<style scoped lang="scss">
.dashboard-stat-card {
  border: 1px solid rgba(37, 99, 235, 0.08);
}

.dashboard-stat-card__label {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.dashboard-stat-card__value {
  margin-top: 12px;
  font-size: 30px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--color-text-primary);
}

.dashboard-stat-card__hint {
  margin-top: 10px;
  font-size: 13px;
  color: var(--color-text-secondary);
}
</style>