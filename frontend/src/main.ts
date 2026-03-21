import { createApp } from 'vue';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';

import App from './App.vue';
import { router } from './router';
import { pinia } from './stores';
import { useUserStore } from './stores/modules/user';
import './styles/index.scss';

const bootstrap = async () => {
  const app = createApp(App);

  app.use(pinia);
  app.use(router);
  app.use(ElementPlus);

  const userStore = useUserStore();
  await userStore.restoreSession();

  app.mount('#app');
};

void bootstrap();