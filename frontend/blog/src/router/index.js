import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue')
      },
      {
        path: 'blog',
        name: 'Blog',
        component: () => import('@/views/Blog.vue')
      },
      {
        path: 'blog/:id',
        name: 'BlogDetail',
        component: () => import('@/views/BlogDetail.vue')
      },
      {
        path: 'works',
        name: 'Works',
        component: () => import('@/views/Works.vue')
      },
      {
        path: 'works/:id',
        name: 'WorkDetail',
        component: () => import('@/views/WorkDetail.vue')
      },
      {
        path: 'about',
        name: 'About',
        component: () => import('@/views/About.vue')
      },
      {
        path: 'links',
        name: 'Links',
        component: () => import('@/views/Links.vue')
      },
      {
        path: 'users/:id',
        name: 'UserProfile',
        component: () => import('@/views/UserProfile.vue')
      },
      {
        path: 'users/:id/follows',
        name: 'FollowList',
        component: () => import('@/views/FollowList.vue')
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
    {
      path: '/dashboard',
      component: () => import('@/layouts/UserCenterLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('@/views/user/Dashboard.vue')
        },
        {
          path: 'articles',
          name: 'MyArticles',
          component: () => import('@/views/user/MyArticles.vue')
        },
        {
          path: 'articles/create',
          name: 'CreateArticle',
          component: () => import('@/views/user/ArticleEdit.vue')
        },
        {
          path: 'articles/:id/edit',
          name: 'EditArticle',
          component: () => import('@/views/user/ArticleEdit.vue')
        },
      {
        path: 'works',
        name: 'MyWorks',
        component: () => import('@/views/user/MyWorks.vue')
      },
      {
        path: 'works/create',
        name: 'CreateWork',
        component: () => import('@/views/user/WorkEdit.vue')
      },
      {
        path: 'works/:id/edit',
        name: 'EditWork',
        component: () => import('@/views/user/WorkEdit.vue')
      },
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/views/user/Notifications.vue')
      },
      {
        path: 'comments',
        name: 'MyComments',
        component: () => import('@/views/user/MyComments.vue')
      }
    ]
  },
  {
    path: '/favorites',
    component: () => import('@/layouts/UserCenterLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Favorites',
        component: () => import('@/views/Favorites.vue')
      }
    ]
  },
  {
    path: '/profile',
    component: () => import('@/layouts/UserCenterLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: 'edit',
        name: 'ProfileEdit',
        component: () => import('@/views/ProfileEdit.vue')
      }
    ]
  },
  // 管理后台路由已移至独立项目 frontend/admin/
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 用户路由守卫
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router

