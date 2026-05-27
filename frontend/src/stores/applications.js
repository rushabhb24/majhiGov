import { defineStore } from 'pinia'
import { ref } from 'vue'
import { API_BASE_URL } from '../config.js'

export const useApplicationStore = defineStore('applications', () => {
  // State
  const applications = ref([])
  const applySubmitting = ref(false)
  const refreshing = ref(false)

  // Actions
  async function fetchApplications() {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    if (!authStore.token) return

    try {
      const response = await fetch(`${API_BASE_URL}/api/user/applications`, {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      })
      if (response.ok) {
        const data = await response.json()
        applications.value = data || []
      }
    } catch (err) {
      console.error('Failed to fetch user applications:', err)
    }
  }

  async function refreshApplications() {
    refreshing.value = true
    await fetchApplications()
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()
    uiStore.showToast('Application statuses refreshed!', 'success')
    refreshing.value = false
  }

  /**
   * Apply via official link:
   * 1. Open the official apply link in a new browser tab
   * 2. Record the application in the backend database
   */
  async function applyViaOfficialLink(scheme) {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    if (!authStore.token) {
      authStore.openAuthModal('login')
      uiStore.showToast('Please login to apply for schemes!', 'info')
      return
    }

    applySubmitting.value = true

    try {
      // 1. Record the application in our database
      const payload = {
        scheme_id: Number(scheme.id),
        notes: 'Applied via official portal link'
      }

      const response = await fetch(`${API_BASE_URL}/api/user/apply`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${authStore.token}`
        },
        body: JSON.stringify(payload)
      })

      if (!response.ok) {
        const errText = await response.text()
        // If already applied, still open the link
        if (errText.includes('already have a pending')) {
          uiStore.showToast('You already applied for this scheme. Opening official portal...', 'info')
        } else {
          throw new Error(errText || 'Failed to record application')
        }
      } else {
        uiStore.showToast('Application recorded! Redirecting to official portal...', 'success')
        fetchApplications()
      }

      // 2. Open the official apply link in a new tab
      const applyUrl = scheme.apply_link || scheme.official_website
      if (applyUrl) {
        window.open(applyUrl, '_blank', 'noopener,noreferrer')
      } else {
        uiStore.showToast('Official apply link not available for this scheme.', 'warning')
      }
    } catch (err) {
      console.error(err)
      uiStore.showToast(err.message || 'Error recording application.', 'danger')
    } finally {
      applySubmitting.value = false
    }
  }

  function clearApplications() {
    applications.value = []
  }

  /**
   * Get the status step number for progress visualization
   * pending = 1, under_review = 2, approved/rejected = 3
   */
  function getStatusStep(status) {
    if (status === 'pending') return 1
    if (status === 'under_review') return 2
    if (status === 'approved' || status === 'rejected') return 3
    return 1
  }

  return {
    applications,
    applySubmitting,
    refreshing,
    fetchApplications,
    refreshApplications,
    applyViaOfficialLink,
    clearApplications,
    getStatusStep
  }
})
