<template>
  <div class="settings-view">
    <el-row :gutter="16">
      <el-col :xs="24" :xl="12">
        <el-card shadow="never" class="settings-card">
          <template #header>
            <div class="settings-card__header">模型配置表单</div>
          </template>

          <el-form label-position="top">
            <el-form-item label="模型名称">
              <el-input
                :model-value="settingsStore.settings.modelName"
                placeholder="请输入模型名称"
                @update:model-value="handleChange('modelName', $event)"
              />
            </el-form-item>

            <el-form-item label="流式输出开关">
              <el-switch
                :model-value="settingsStore.settings.streamOutput"
                @update:model-value="handleChange('streamOutput', $event)"
              />
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :xl="12">
        <el-card shadow="never" class="settings-card">
          <template #header>
            <div class="settings-card__header">面试偏好配置表单</div>
          </template>

          <el-form label-position="top">
            <el-form-item label="默认面试模式">
              <el-select
                :model-value="settingsStore.settings.defaultInterviewMode"
                class="settings-card__control"
                @update:model-value="handleChange('defaultInterviewMode', $event)"
              >
                <el-option label="综合模拟" value="综合模拟" />
                <el-option label="专项训练" value="专项训练" />
                <el-option label="项目拷打" value="项目拷打" />
              </el-select>
            </el-form-item>

            <el-form-item label="追问强度">
              <el-segmented
                :model-value="settingsStore.settings.followUpIntensity"
                :options="['低', '中', '高']"
                @update:model-value="handleChange('followUpIntensity', $event)"
              />
            </el-form-item>

            <el-form-item label="评分严格程度">
              <el-radio-group
                :model-value="settingsStore.settings.scoringStrictness"
                @update:model-value="handleChange('scoringStrictness', $event)"
              >
                <el-radio-button label="宽松" value="宽松" />
                <el-radio-button label="标准" value="标准" />
                <el-radio-button label="严格" value="严格" />
              </el-radio-group>
            </el-form-item>

            <el-form-item>
              <el-button @click="settingsStore.resetSettings()">重置为默认配置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from '@/stores/modules/settings';
import type { SettingsState } from '@/types/settings';

const settingsStore = useSettingsStore();

const handleChange = <K extends keyof SettingsState>(key: K, value: SettingsState[K]) => {
  settingsStore.updateSettings({ [key]: value } as Partial<SettingsState>);
};
</script>

<style scoped lang="scss">
.settings-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.settings-card {
  height: 100%;
}

.settings-card__header {
  font-weight: 600;
}

.settings-card__control {
  width: 100%;
}
</style>
