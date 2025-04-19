import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      meta: {
        description: "首页",
      },
      component: () => import("../views/home.vue")
    },
    {
      path: '/list',
      name: 'list',
      meta: {
        description: "文件列表",
      },
      component: () => import("../views/list.vue")
    },
    {
      path: '/platform/settings',
      name: '',
      meta: {
        description: "配置",
      },
      component: () => import("../views/platform/configuration.vue")
    },
  ]
})

export default router
