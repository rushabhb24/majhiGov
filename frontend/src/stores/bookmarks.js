import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useSchemeStore } from './schemes.js'
import { useJobStore } from './jobs.js'
import { bookmarksApi } from '../api/bookmarks.js'

export const useBookmarkStore = defineStore('bookmarks', () => {
  // State
  const savedSchemeIds = ref([])
  const savedJobIds = ref([])
  const savedSchemesData = ref([])
  const savedJobsData = ref([])

  // Getters
  const bookmarkedSchemes = computed(() => {
    if (savedSchemesData.value.length > 0) {
      return savedSchemesData.value
    }
    const schemeStore = useSchemeStore()
    return schemeStore.schemes.filter(s => savedSchemeIds.value.includes(s.id))
  })

  const bookmarkedJobs = computed(() => {
    if (savedJobsData.value.length > 0) {
      return savedJobsData.value
    }
    const jobStore = useJobStore()
    return jobStore.jobs.filter(j => savedJobIds.value.includes(j.id))
  })

  // Actions
  async function fetchSavedSchemes() {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    if (!authStore.userProfile) return

    try {
      const data = await bookmarksApi.fetchSavedSchemes()
      if (Array.isArray(data)) {
        savedSchemeIds.value = data
        localStorage.setItem('yojana_saved_ids', JSON.stringify(data))
        
        // Clean up cached schemes data to match IDs
        savedSchemesData.value = savedSchemesData.value.filter(s => data.includes(s.id))
        localStorage.setItem('yojana_saved_schemes_data', JSON.stringify(savedSchemesData.value))
      }
    } catch (err) {
      console.error('Failed to fetch bookmarks from server:', err)
    }
  }

  async function toggleBookmark(schemeId) {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    if (!authStore.userProfile) {
      authStore.openAuthModal('login')
      uiStore.showToast('Please login to save schemes!', 'info')
      return
    }

    try {
      await bookmarksApi.toggleSavedScheme(schemeId)

      const index = savedSchemeIds.value.indexOf(schemeId)
      if (index === -1) {
        savedSchemeIds.value.push(schemeId)
        
        // Save full scheme object to localStorage cache
        const schemeStore = useSchemeStore()
        const schemeObj = schemeStore.schemes.find(s => s.id === schemeId)
        if (schemeObj && !savedSchemesData.value.some(s => s.id === schemeId)) {
          savedSchemesData.value.push(schemeObj)
        }
        
        uiStore.showToast('Scheme saved to your profile!', 'success')
      } else {
        savedSchemeIds.value.splice(index, 1)
        const dIndex = savedSchemesData.value.findIndex(s => s.id === schemeId)
        if (dIndex !== -1) {
          savedSchemesData.value.splice(dIndex, 1)
        }
        uiStore.showToast('Scheme removed from your profile.', 'info')
      }
      localStorage.setItem('yojana_saved_ids', JSON.stringify(savedSchemeIds.value))
      localStorage.setItem('yojana_saved_schemes_data', JSON.stringify(savedSchemesData.value))
    } catch (err) {
      console.error('Error toggling bookmark:', err)
      uiStore.showToast('Failed to sync bookmark with account.', 'danger')
    }
  }

  async function toggleJobBookmark(jobId) {
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    const index = savedJobIds.value.indexOf(jobId)
    if (index === -1) {
      savedJobIds.value.push(jobId)
      
      // Save full job object to localStorage cache
      const jobStore = useJobStore()
      let jobObj = jobStore.jobs.find(j => j.id === jobId)
      if (!jobObj) {
        const { usePrivateJobStore } = await import('./privateJobs.js')
        const privateJobStore = usePrivateJobStore()
        jobObj = privateJobStore.jobs?.find(j => j.id === jobId)
      }
      if (jobObj && !savedJobsData.value.some(j => j.id === jobId)) {
        savedJobsData.value.push(jobObj)
      }
      
      uiStore.showToast('Job saved to your profile!', 'success')
    } else {
      savedJobIds.value.splice(index, 1)
      const dIndex = savedJobsData.value.findIndex(j => j.id === jobId)
      if (dIndex !== -1) {
        savedJobsData.value.splice(dIndex, 1)
      }
      uiStore.showToast('Job removed from your profile.', 'info')
    }
    localStorage.setItem('yojana_saved_job_ids', JSON.stringify(savedJobIds.value))
    localStorage.setItem('yojana_saved_jobs_data', JSON.stringify(savedJobsData.value))
  }

  function loadBookmarks() {
    try {
      const data = localStorage.getItem('yojana_saved_ids')
      if (data) {
        savedSchemeIds.value = JSON.parse(data)
      }
      const schemesData = localStorage.getItem('yojana_saved_schemes_data')
      if (schemesData) {
        savedSchemesData.value = JSON.parse(schemesData)
      }
      const jobData = localStorage.getItem('yojana_saved_job_ids')
      if (jobData) {
        savedJobIds.value = JSON.parse(jobData)
      }
      const jobsData = localStorage.getItem('yojana_saved_jobs_data')
      if (jobsData) {
        savedJobsData.value = JSON.parse(jobsData)
      }
    } catch (e) {
      console.error('Failed to load bookmarks:', e)
    }
  }

  function clearBookmarks() {
    savedSchemeIds.value = []
    savedJobIds.value = []
    savedSchemesData.value = []
    savedJobsData.value = []
  }

  return {
    savedSchemeIds,
    savedJobIds,
    savedSchemesData,
    savedJobsData,
    bookmarkedSchemes,
    bookmarkedJobs,
    fetchSavedSchemes,
    toggleBookmark,
    toggleJobBookmark,
    loadBookmarks,
    clearBookmarks
  }
})
