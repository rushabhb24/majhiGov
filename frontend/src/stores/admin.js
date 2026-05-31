import { defineStore } from 'pinia'
import { ref } from 'vue'
import { adminApi } from '../api/admin.js'

export const useAdminStore = defineStore('admin', () => {
  // State
  const analytics = ref(null)
  const schemes = ref([])
  const categories = ref([])
  const users = ref([])
  const notifications = ref([])
  const applications = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Actions
  async function fetchAnalytics() {
    loading.value = true
    error.value = null
    try {
      analytics.value = await adminApi.fetchAnalytics()
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
      schemes.value = await adminApi.fetchAllSchemes()
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
      const result = await adminApi.createScheme(payload)
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
      const result = await adminApi.updateScheme(schemeId, payload)
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
      const result = await adminApi.toggleSchemeStatus(schemeId)
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
      categories.value = await adminApi.fetchAllCategories()
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
      const result = await adminApi.createCategory(payload)
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
      const result = await adminApi.deleteCategory(categoryId)
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
      users.value = await adminApi.fetchAllUsers()
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
      const result = await adminApi.toggleUserStatus(userId)
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
      const result = await adminApi.createAdmin(payload)
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
      const result = await adminApi.sendNotification(payload)
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
      notifications.value = await adminApi.fetchNotifications()
    } catch (err) {
      console.error(err)
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  async function fetchAllApplications() {
    loading.value = true
    error.value = null
    try {
      applications.value = await adminApi.fetchAllApplications()
    } catch (err) {
      console.error(err)
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  async function updateApplicationStatus(applicationId, status, notes = '') {
    loading.value = true
    error.value = null
    try {
      const result = await adminApi.updateApplicationStatus({
        application_id: Number(applicationId),
        status,
        notes
      })
      await fetchAllApplications()
      return result
    } catch (err) {
      console.error(err)
      error.value = err.message
      throw err
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
    applications,
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
    fetchNotifications,
    applications,
    fetchAllApplications,
    updateApplicationStatus
  }
})
