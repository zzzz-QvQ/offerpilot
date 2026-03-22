<template>
  <el-card class="interview-input-box" shadow="never">
    <el-input
      v-model="inputValue"
      type="textarea"
      :rows="4"
      resize="none"
      placeholder="输入你对当前题目的回答"
      :disabled="submitting || streaming"
    />

    <div class="interview-input-box__actions">
      <el-button v-if="streaming" plain @click="emit('stop')">停止输出</el-button>
      <el-button type="primary" :loading="submitting" :disabled="streaming || !inputValue.trim()" @click="handleSubmit">
        提交回答
      </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage } from 'element-plus';

const emit = defineEmits<{
  submit: [content: string];
  stop: [];
}>();

defineProps<{
  submitting: boolean;
  streaming: boolean;
}>();

const inputValue = ref('');

const handleSubmit = () => {
  const value = inputValue.value.trim();

  if (!value) {
    ElMessage.warning('请先输入本轮回答内容');
    return;
  }

  emit('submit', value);
  inputValue.value = '';
};
</script>

<style scoped lang="scss">
.interview-input-box__actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 14px;
}
</style>
