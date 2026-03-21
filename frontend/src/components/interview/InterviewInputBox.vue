<template>
  <el-card class="interview-input-box" shadow="never">
    <el-input
      v-model="inputValue"
      type="textarea"
      :rows="4"
      resize="none"
      placeholder="Enter your answer for the current round."
      :disabled="submitting || streaming"
    />

    <div class="interview-input-box__actions">
      <el-button v-if="streaming" plain @click="emit('stop')">Stop</el-button>
      <el-button type="primary" :loading="submitting" :disabled="streaming" @click="handleSubmit">Send Answer</el-button>
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
    ElMessage.warning('Please enter your answer first.');
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