import { createRouter, createWebHistory } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const routes = [
  {
    path: '/login',
    name: 'AdminLogin',
    component: () => import('@/views/admin/Login.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/admin/Dashboard.vue')
      },
      {
        path: 'articles',
        name: 'Articles',
        component: () => import('@/views/admin/Articles.vue')
      },
      {
        path: 'articles/create',
        name: 'ArticleCreate',
        component: () => import('@/views/admin/ArticleEdit.vue')
      },
      {
        path: 'articles/:id/edit',
        name: 'ArticleEdit',
        component: () => import('@/views/admin/ArticleEdit.vue')
      },
      {
        path: 'works',
        name: 'Works',
        component: () => import('@/views/admin/Works.vue')
      },
      {
        path: 'categories',
        name: 'Categories',
        component: () => import('@/views/admin/Categories.vue')
      },
      {
        path: 'tags',
        name: 'Tags',
        component: () => import('@/views/admin/Tags.vue')
      },
      {
        path: 'comments',
        name: 'Comments',
        component: () => import('@/views/admin/Comments.vue')
      },
      {
        path: 'links',
        name: 'Links',
        component: () => import('@/views/admin/Links.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/admin/Settings.vue')
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/admin/Users.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const adminStore = useAdminStore()
  
  if (to.meta.requiresAuth && !adminStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router
