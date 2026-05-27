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
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// Navigation guard for auth-protected routes
router.beforeEach((to) => {
  if (to.meta.requiresAuth) {
    const authStore = useAuthStore()
    if (!authStore.isLoggedIn) {
      authStore.openAuthModal('login')
      return false
    }
  }
})

export default router
