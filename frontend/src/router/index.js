import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'explorer',
    component: () => import('../views/ExplorerView.vue')
  },
  {
    path: '/eligibility',
    name: 'eligibility',
    component: () => import('../views/EligibilityView.vue')
  },
  {
    path: '/saved',
    name: 'saved',
    component: () => import('../views/SavedView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/jobs',
    name: 'jobs',
    component: () => import('../views/JobsView.vue')
  },
  {
    path: '/ai-assistant',
    name: 'ai-assistant',
    component: () => import('../views/AIAssistantView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/career-resources',
    name: 'career-resources',
    component: () => import('../views/CareerResourcesView.vue')
  },
  {
    path: '/companies',
    name: 'companies',
    component: () => import('../views/CompaniesView.vue')
  },
  {
    path: '/internships',
    name: 'internships',
    component: () => import('../views/JobsView.vue')
  },
  {
    path: '/applications',
    name: 'applications',
    component: () => import('../views/ApplicationsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('../views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin-dashboard',
    name: 'admin-login',
    component: () => import('../views/AdminLoginView.vue')
  },
  {
    path: '/admin/dashboard',
    name: 'admin',
    component: () => import('../views/AdminView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// Navigation guard for auth-protected routes
router.beforeEach((to) => {
  const authStore = useAuthStore()

  // Enforce that logged-in administrators are restricted to the Admin Dashboard
  if (authStore.isAdmin) {
    if (to.path !== '/admin/dashboard') {
      return '/admin/dashboard'
    }
  }

  // Prevent non-admin logged-in users from accessing the admin login page
  if (authStore.isLoggedIn && !authStore.isAdmin && to.path === '/admin-dashboard') {
    return '/'
  }

  // Strict prefix guard for admin routes
  if (to.path.startsWith('/admin/')) {
    if (!authStore.isAdmin) {
      return '/admin-dashboard'
    }
  }

  if (to.meta.requiresAuth) {
    if (!authStore.isLoggedIn) {
      authStore.openAuthModal('login')
      return false
    }
  }

  if (to.meta.requiresAdmin) {
    if (!authStore.isAdmin) {
      return '/admin-dashboard'
    }
  }
})

export default router
