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
        path: 'articles/:id',
        name: 'ArticleView',
        component: () => import('@/views/admin/ArticleView.vue')
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
      },
      {
        path: 'ads',
        name: 'Ads',
        component: () => import('@/views/admin/Ads.vue')
      },
      {
        path: 'knowledge',
        name: 'KnowledgeOverview',
        component: () => import('@/views/admin/KnowledgeOverview.vue')
      },
      { path: 'knowledge/workspaces', name: 'KnowledgeWorkspaces', component: () => import('@/views/admin/KnowledgeWorkspaces.vue') },
      { path: 'knowledge/docs', name: 'KnowledgeDocs', component: () => import('@/views/admin/KnowledgeDocs.vue') },
      { path: 'knowledge/quota', name: 'KnowledgeQuota', component: () => import('@/views/admin/KnowledgeQuota.vue') },
      { path: 'knowledge/shares', name: 'KnowledgeShares', component: () => import('@/views/admin/KnowledgeShares.vue') },
      { path: 'knowledge/usage', name: 'KnowledgeUsage', component: () => import('@/views/admin/KnowledgeUsage.vue') },
      { path: 'knowledge/audit-logs', name: 'KnowledgeAuditLogs', component: () => import('@/views/admin/KnowledgeAuditLogs.vue') }
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
