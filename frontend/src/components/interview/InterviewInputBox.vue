<template>
  <el-card class="interview-input-box" shadow="never">
    <el-input
      v-model="inputValue"
      type="textarea"
      :rows="4"
      resize="none"
      placeholder="请输入你的回答，当前为本地 mock 会话，后续可接入流式面试服务。"
    />

    <div class="interview-input-box__actions">
      <el-button type="primary" :loading="submitting" @click="handleSubmit">发送回答</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage } from 'element-plus';

const emit = defineEmits<{
  submit: [content: string];
}>();

defineProps<{
  submitting: boolean;
}>();

const inputValue = ref('');

const handleSubmit = () => {
  const value = inputValue.value.trim();

  if (!value) {
    ElMessage.warning('请输入回答内容');
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
  margin-top: 14px;
}
</style>