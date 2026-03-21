<template>
  <el-container class="main-layout">
    <el-aside class="main-layout__aside" width="240px">
      <div class="main-layout__brand">
        <div class="main-layout__brand-title">OfferPilot</div>
        <div class="main-layout__brand-subtitle">Frontend Interview Trainer</div>
      </div>

      <el-menu
        :default-active="activePath"
        class="main-layout__menu"
        router
      >
        <el-menu-item
          v-for="item in menuItems"
          :key="item.path"
          :index="item.path"
        >
          <span>{{ item.label }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-layout__container">
      <el-header class="main-layout__header">
        <div class="main-layout__header-left">
          <h1 class="main-layout__title">{{ pageTitle }}</h1>
        </div>

        <div class="main-layout__header-right">
          <span class="main-layout__nickname">{{ displayName }}</span>
          <el-button link type="primary" @click="handleLogout">退出</el-button>
        </div>
      </el-header>

      <el-main class="main-layout__main">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { useUserStore } from '@/stores/modules/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const menuItems = [
  { path: '/dashboard', label: '工作台首页' },
  { path: '/question-bank', label: '题库学习' },
  { path: '/interview', label: '模拟面试' },
  { path: '/review', label: '复盘中心' },
  { path: '/project-polish', label: '项目包装' },
  { path: '/settings', label: '设置' },
];

const activePath = computed(() => route.path);
const pageTitle = computed(() => String(route.meta.title ?? 'OfferPilot'));
const displayName = computed(() => userStore.userInfo?.nickname || '未登录用户');

const handleLogout = async () => {
  userStore.logout();
  await router.push('/login');
};
</script>

<style scoped lang="scss">
.main-layout {
  min-height: 100vh;
  background: var(--color-page-bg);
}

.main-layout__aside {
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--color-border);
  background: linear-gradient(180deg, #ffffff 0%, #f8fbff 100%);
}

.main-layout__brand {
  padding: 24px 20px 20px;
  border-bottom: 1px solid rgba(229, 231, 235, 0.8);
}

.main-layout__brand-title {
  font-size: 22px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--color-text-primary);
}

.main-layout__brand-subtitle {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-secondary);
  letter-spacing: 0.04em;
}

.main-layout__menu {
  flex: 1;
  padding: 12px;
  border-right: none;
  background: transparent;
}

.main-layout__menu :deep(.el-menu-item) {
  height: 44px;
  margin-bottom: 6px;
  border-radius: 10px;
  color: var(--color-text-secondary);
}

.main-layout__menu :deep(.el-menu-item:hover) {
  background: #eef4ff;
  color: var(--color-brand);
}

.main-layout__menu :deep(.el-menu-item.is-active) {
  background: #e8f0ff;
  color: var(--color-brand);
  font-weight: 600;
}

.main-layout__container {
  min-width: 0;
}

.main-layout__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 72px;
  padding: 0 24px;
  border-bottom: 1px solid var(--color-border);
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
}

.main-layout__header-left,
.main-layout__header-right {
  display: flex;
  align-items: center;
}

.main-layout__title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.main-layout__header-right {
  gap: 12px;
}

.main-layout__nickname {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.main-layout__main {
  padding: 24px;
}
</style>