import { defineStore } from 'pinia'
import { ref } from 'vue'
import { API_BASE_URL } from '../config.js'
import { useAuthStore } from './auth'

export const useAdminStore = defineStore('admin', () => {
  const authStore = useAuthStore()

  // State
  const analytics = ref(null)
  const schemes = ref([])
  const categories = ref([])
  const users = ref([])
  const notifications = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Request helper
  async function adminFetch(endpoint, options = {}) {
    if (!authStore.token) {
      throw new Error('No authentication session found')
    }

    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authStore.token}`,
      ...(options.headers || {})
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers
    })

    if (!response.ok) {
      const errText = await response.text()
      let errMsg = 'API request failed'
      try {
        const parsed = JSON.parse(errText)
        errMsg = parsed.error || errMsg
      } catch (e) {
        errMsg = errText || errMsg
      }
      throw new Error(errMsg)
    }

    if (response.status === 204) return null
    return await response.json()
  }

  // Actions
  async function fetchAnalytics() {
    loading.value = true
    error.value = null
    try {
      analytics.value = await adminFetch('/api/admin/analytics')
    } catch (err) {
      console.error(err)
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  async function fetchAllSchemes() {
    loading.value = true
    error.value = null
    try {
      schemes.value = await adminFetch('/api/admin/schemes')
    } catch (err) {
      console.error(err)
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  async function createScheme(payload) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch('/api/admin/schemes', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      await fetchAllSchemes()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateScheme(schemeId, payload) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch(`/api/admin/schemes/${schemeId}`, {
        method: 'PUT',
        body: JSON.stringify(payload)
      })
      await fetchAllSchemes()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function toggleSchemeStatus(schemeId) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch(`/api/admin/schemes/${schemeId}`, {
        method: 'DELETE'
      })
      await fetchAllSchemes()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchAllCategories() {
    loading.value = true
    error.value = null
    try {
      categories.value = await adminFetch('/api/admin/categories')
    } catch (err) {
      console.error(err)
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  async function createCategory(payload) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch('/api/admin/categories', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      await fetchAllCategories()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteCategory(categoryId) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch(`/api/admin/categories/${categoryId}`, {
        method: 'DELETE'
      })
      await fetchAllCategories()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchAllUsers() {
    loading.value = true
    error.value = null
    try {
      users.value = await adminFetch('/api/admin/users')
    } catch (err) {
      console.error(err)
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  async function toggleUserStatus(userId) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch('/api/admin/users/toggle-active', {
        method: 'POST',
        body: JSON.stringify({ user_id: userId })
      })
      await fetchAllUsers()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createAdmin(payload) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch('/api/admin/users/admin', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      await fetchAllUsers()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function sendNotification(payload) {
    loading.value = true
    error.value = null
    try {
      const result = await adminFetch('/api/admin/notifications', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      await fetchNotifications()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchNotifications() {
    loading.value = true
    error.value = null
    try {
      notifications.value = await adminFetch('/api/admin/notifications')
    } catch (err) {
      console.error(err)
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  return {
    analytics,
    schemes,
    categories,
    users,
    notifications,
    loading,
    error,
    fetchAnalytics,
    fetchAllSchemes,
    createScheme,
    updateScheme,
    toggleSchemeStatus,
    fetchAllCategories,
    createCategory,
    deleteCategory,
    fetchAllUsers,
    toggleUserStatus,
    createAdmin,
    sendNotification,
    fetchNotifications
  }
})
