import { createRouter, createWebHistory } from 'vue-router';

import { staticRoutes } from './modules/static';

export const router = createRouter({
  history: createWebHistory(),
  routes: staticRoutes,
  scrollBehavior() {
    return { top: 0 };
  },
});
