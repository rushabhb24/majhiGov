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
   * Records the application in the backend database.
   * Note: The caller (App.vue handleApplyAction) opens the official URL synchronously
   * before calling this function to avoid popup blockers.
   */
  async function applyViaOfficialLink(scheme) {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    if (!authStore.token) return

    // 2. Record the application in the background (non-blocking)
    applySubmitting.value = true
    try {
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

      if (response.ok) {
        uiStore.showToast('Application recorded! Track your status in My Applications.', 'success')
        fetchApplications()
      } else {
        const errData = await response.text()
        if (response.status === 409 || errData.includes('already')) {
          uiStore.showToast('You already applied for this scheme. Opening official portal...', 'info')
        } else {
          console.error('Failed to record application:', errData)
        }
      }
    } catch (err) {
      console.error('Error recording application:', err)
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
