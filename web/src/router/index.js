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
        name: 'UserDashboard',
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
      }
    ]
  },
  {
    path: '/favorites',
    name: 'Favorites',
    component: () => import('@/layouts/UserCenterLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        component: () => import('@/views/Favorites.vue')
      }
    ]
  },
  {
    path: '/notifications',
    name: 'Notifications',
    component: () => import('@/layouts/UserCenterLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        component: () => import('@/views/Notifications.vue')
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
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/admin/Login.vue')
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAdminAuth: true },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue')
      },
      {
        path: 'articles',
        name: 'AdminArticles',
        component: () => import('@/views/admin/Articles.vue')
      },
      {
        path: 'articles/create',
        name: 'AdminArticleCreate',
        component: () => import('@/views/admin/ArticleEdit.vue')
      },
      {
        path: 'articles/:id/edit',
        name: 'AdminArticleEdit',
        component: () => import('@/views/admin/ArticleEdit.vue')
      },
      {
        path: 'works',
        name: 'AdminWorks',
        component: () => import('@/views/admin/Works.vue')
      },
      {
        path: 'categories',
        name: 'AdminCategories',
        component: () => import('@/views/admin/Categories.vue')
      },
      {
        path: 'tags',
        name: 'AdminTags',
        component: () => import('@/views/admin/Tags.vue')
      },
      {
        path: 'comments',
        name: 'AdminComments',
        component: () => import('@/views/admin/Comments.vue')
      },
      {
        path: 'links',
        name: 'AdminLinks',
        component: () => import('@/views/admin/Links.vue')
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: () => import('@/views/admin/Settings.vue')
      },
      {
        path: 'users',
        name: 'AdminUsers',
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
  const userStore = useUserStore()
  
  // 管理后台路由守卫
  if (to.meta.requiresAdminAuth) {
    const adminToken = localStorage.getItem('admin_token')
    if (!adminToken) {
      next('/admin/login')
      return
    }
  }
  
  // 用户路由守卫
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router

