<template>
  <el-card shadow="never">
    <el-form label-position="top" class="question-filter-bar">
      <el-row :gutter="16">
        <el-col :xs="24" :md="10" :lg="12">
          <el-form-item label="搜索题目">
            <el-input
              :model-value="filters.keyword"
              placeholder="请输入题目关键词"
              clearable
              @update:model-value="handleKeywordChange"
            />
          </el-form-item>
        </el-col>

        <el-col :xs="24" :md="7" :lg="6">
          <el-form-item label="分类筛选">
            <el-select
              :model-value="filters.category"
              placeholder="全部分类"
              clearable
              class="question-filter-bar__control"
              @update:model-value="handleCategoryChange"
            >
              <el-option
                v-for="item in categories"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :xs="24" :md="7" :lg="6">
          <el-form-item label="频率筛选">
            <el-select
              :model-value="filters.frequency"
              placeholder="全部频率"
              clearable
              class="question-filter-bar__control"
              @update:model-value="handleFrequencyChange"
            >
              <el-option label="高频" value="高频" />
              <el-option label="中频" value="中频" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <div class="question-filter-bar__actions">
        <el-button @click="emit('reset')">重置筛选</el-button>
      </div>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import type { QuestionBankFilters, QuestionCategoryOption } from '@/types/question-bank';

const props = defineProps<{
  filters: QuestionBankFilters;
  categories: QuestionCategoryOption[];
}>();

const emit = defineEmits<{
  change: [payload: Partial<QuestionBankFilters>];
  reset: [];
}>();

const handleKeywordChange = (value: string | undefined) => {
  emit('change', { keyword: value ?? '' });
};

const handleCategoryChange = (value: string | undefined) => {
  emit('change', { category: value ?? '' });
};

const handleFrequencyChange = (value: QuestionBankFilters['frequency'] | undefined) => {
  emit('change', { frequency: value ?? '' });
};
</script>

<style scoped lang="scss">
.question-filter-bar__control {
  width: 100%;
}

.question-filter-bar__actions {
  display: flex;
  justify-content: flex-end;
}
</style>