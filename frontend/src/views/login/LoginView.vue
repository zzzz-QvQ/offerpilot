<template>
  <div class="page-shell login-page">
    <el-card class="login-card" shadow="never">
      <div class="login-card__header">
        <h1 class="login-card__title">欢迎登录 OfferPilot</h1>
        <p class="login-card__desc">前端面试智能训练平台</p>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        size="large"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
            type="email"
            autocomplete="username"
          />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            placeholder="请输入密码"
            type="password"
            show-password
            autocomplete="current-password"
            @keyup.enter="handleSubmit"
          />
        </el-form-item>

        <el-form-item class="login-card__action">
          <el-button type="primary" :loading="submitting" class="login-card__button" @click="handleSubmit">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';
import { useRouter } from 'vue-router';

import { useUserStore } from '@/stores/modules/user';
import type { LoginParams } from '@/types/user';

const router = useRouter();
const userStore = useUserStore();
const formRef = ref<FormInstance>();
const submitting = ref(false);

const form = reactive<LoginParams>({
  email: '',
  password: '',
});

const rules: FormRules<LoginParams> = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入合法的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' },
  ],
};

const handleSubmit = async () => {
  if (!formRef.value || submitting.value) {
    return;
  }

  const valid = await formRef.value.validate().catch(() => false);

  if (!valid) {
    return;
  }

  submitting.value = true;

  try {
    await userStore.login(form);
    ElMessage.success('登录成功');
    await router.push('/dashboard');
  } catch (error) {
    const message = error instanceof Error ? error.message : '登录失败';
    ElMessage.error(message);
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped lang="scss">
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background:
    radial-gradient(circle at top left, rgba(37, 99, 235, 0.14), transparent 28%),
    linear-gradient(135deg, #eef4ff 0%, #f8fafc 55%, #eef2ff 100%);
}

.login-card {
  width: 100%;
  max-width: 440px;
  border: 1px solid rgba(37, 99, 235, 0.12);
}

.login-card__header {
  margin-bottom: 24px;
}

.login-card__title {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.login-card__desc {
  margin: 8px 0 0;
  color: var(--color-text-secondary);
}

.login-card__action {
  margin-bottom: 0;
}

.login-card__button {
  width: 100%;
}
</style>