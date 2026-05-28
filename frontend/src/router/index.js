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
  console.log("Router: Navigating to:", to.path, "requiresAuth:", !!to.meta.requiresAuth, "requiresAdmin:", !!to.meta.requiresAdmin, "isLoggedIn:", authStore.isLoggedIn, "isAdmin:", authStore.isAdmin)

  // Strict prefix guard for admin routes
  if (to.path.startsWith('/admin/')) {
    if (!authStore.isAdmin) {
      console.log("Router: Access Denied to admin prefix. Redirecting to Admin Login Page.")
      return '/admin-dashboard'
    }
  }

  if (to.meta.requiresAuth) {
    if (!authStore.isLoggedIn) {
      console.log("Router: Access Denied. User not logged in. Opening login modal.")
      authStore.openAuthModal('login')
      return false
    }
  }

  if (to.meta.requiresAdmin) {
    if (!authStore.isAdmin) {
      console.log("Router: Access Denied. User is not an admin. Redirecting to Admin Login Page.")
      return '/admin-dashboard'
    }
  }

  console.log("Router: Navigation allowed to:", to.path)
})

export default router
