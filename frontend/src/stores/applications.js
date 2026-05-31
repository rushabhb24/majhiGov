import { defineStore } from 'pinia'
import { ref } from 'vue'
import { applicationsApi } from '../api/applications.js'

export const useApplicationStore = defineStore('applications', () => {
  // State
  const applications = ref([])
  const applySubmitting = ref(false)
  const refreshing = ref(false)
  const applyModalOpen = ref(false)
  const applyNotes = ref('')
  const applyingScheme = ref(null)

  // Actions
  async function fetchApplications() {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    if (!authStore.token) return

    try {
      const data = await applicationsApi.fetchUserApplications()
      applications.value = data || []
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
   */
  async function applyViaOfficialLink(scheme) {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    if (!authStore.token) return

    applySubmitting.value = true
    try {
      const payload = {
        scheme_id: Number(scheme.id),
        notes: 'Applied via official portal link'
      }

      await applicationsApi.applyForScheme(payload)
      uiStore.showToast('Application recorded! Track your status in My Applications.', 'success')
      await fetchApplications()
    } catch (err) {
      if (err.message.includes('already')) {
        uiStore.showToast('You already applied for this scheme. Opening official portal...', 'info')
      } else {
        console.error('Error recording application:', err)
      }
    } finally {
      applySubmitting.value = false
    }
  }

  function clearApplications() {
    applications.value = []
  }

  function openApplyModal(scheme) {
    applyingScheme.value = scheme
    applyNotes.value = ''
    applyModalOpen.value = true
  }

  function closeApplyModal() {
    applyModalOpen.value = false
    applyingScheme.value = null
    applyNotes.value = ''
  }

  async function submitApplyForm() {
    if (!applyingScheme.value) return

    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    if (!authStore.token) return

    applySubmitting.value = true
    try {
      const payload = {
        scheme_id: Number(applyingScheme.value.id),
        notes: applyNotes.value
      }

      await applicationsApi.applyForScheme(payload)
      uiStore.showToast('Application submitted successfully!', 'success')
      closeApplyModal()
      
      // Redirect or track status
      await fetchApplications()
      
      // Delay open redirect
      if (applyingScheme.value.apply_link) {
        setTimeout(() => {
          window.open(applyingScheme.value.apply_link, '_blank')
        }, 800)
      }
    } catch (err) {
      uiStore.showToast(err.message || 'Failed to submit application.', 'danger')
    } finally {
      applySubmitting.value = false
    }
  }

  function getStatusStep(status) {
    if (status === 'approved' || status === 'rejected') return 3
    if (status === 'pending') return 2
    return 1
  }

  return {
    applications,
    applySubmitting,
    refreshing,
    applyModalOpen,
    applyNotes,
    applyingScheme,
    fetchApplications,
    refreshApplications,
    applyViaOfficialLink,
    clearApplications,
    openApplyModal,
    closeApplyModal,
    submitApplyForm,
    getStatusStep
  }
})

