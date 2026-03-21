<template>
  <el-card shadow="never" class="project-input-form">
    <template #header>
      <div class="project-input-form__header">项目原始信息</div>
    </template>

    <el-form ref="formRef" :model="localForm" :rules="rules" label-position="top">
      <el-form-item label="项目名称" prop="projectName">
        <el-input v-model="localForm.projectName" placeholder="请输入项目名称" />
      </el-form-item>

      <el-form-item label="角色定位" prop="role">
        <el-input v-model="localForm.role" placeholder="请输入你在项目中的角色" />
      </el-form-item>

      <el-form-item label="技术栈" prop="techStack">
        <el-input v-model="localForm.techStack" placeholder="请输入主要技术栈" />
      </el-form-item>

      <el-form-item label="项目背景" prop="projectBackground">
        <el-input v-model="localForm.projectBackground" type="textarea" :rows="4" resize="none" />
      </el-form-item>

      <el-form-item label="职责与贡献" prop="responsibility">
        <el-input v-model="localForm.responsibility" type="textarea" :rows="5" resize="none" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="loading" @click="handleSubmit">生成项目包装</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import { ref } from 'vue';

import type { ProjectPolishFormData } from '@/types/project';

const props = defineProps<{
  modelValue: ProjectPolishFormData;
  loading: boolean;
}>();

const emit = defineEmits<{
  submit: [payload: ProjectPolishFormData];
}>();

const formRef = ref<FormInstance>();
const localForm = reactive<ProjectPolishFormData>({ ...props.modelValue });

watch(
  () => props.modelValue,
  (value) => {
    Object.assign(localForm, value);
  },
  { deep: true },
);

const rules: FormRules<ProjectPolishFormData> = {
  projectName: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  role: [{ required: true, message: '请输入角色定位', trigger: 'blur' }],
  techStack: [{ required: true, message: '请输入技术栈', trigger: 'blur' }],
  projectBackground: [{ required: true, message: '请输入项目背景', trigger: 'blur' }],
  responsibility: [{ required: true, message: '请输入职责与贡献', trigger: 'blur' }],
};

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false);

  if (!valid) {
    return;
  }

  emit('submit', { ...localForm });
};
</script>

<style scoped lang="scss">
.project-input-form {
  height: 100%;
}

.project-input-form__header {
  font-weight: 600;
}
</style>