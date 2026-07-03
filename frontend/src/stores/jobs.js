import { defineStore } from 'pinia'
import { ref } from 'vue'
import { jobsApi } from '../api/jobs.js'

const PAGE_SIZE = 5 // Infinite scroll loads 5 records at a time

export const useJobStore = defineStore('jobs', () => {
  // State
  const jobs = ref([])
  const loading = ref(false)
  const loadingMore = ref(false)
  const error = ref(null)
  const searchQuery = ref('')
  const selectedQualification = ref('All')
  const selectedJob = ref(null)
  const detailModalOpen = ref(false)
  const qualifications = ['All', '10th Pass', '12th Pass', 'Graduate', 'Post Graduate']

  // Pagination state
  const currentPage = ref(1)
  const totalJobs = ref(0)
  const hasNextPage = ref(false)

  // Actions
  async function fetchJobs() {
    loading.value = true
    error.value = null
    currentPage.value = 1
    jobs.value = []
    try {
      const resp = await jobsApi.fetchPublicJobs({
        qualification: selectedQualification.value,
        search: searchQuery.value,
        page: 1,
        limit: PAGE_SIZE
      })
      jobs.value = resp.data || []
      totalJobs.value = resp.meta?.total ?? 0
      hasNextPage.value = resp.meta?.hasNext ?? false
    } catch (err) {
      error.value = 'Could not fetch government jobs.'
      jobs.value = []
    } finally {
      loading.value = false
    }
  }

  async function loadMoreJobs() {
    if (loadingMore.value || !hasNextPage.value) return
    loadingMore.value = true
    try {
      const nextPage = currentPage.value + 1
      const resp = await jobsApi.fetchPublicJobs({
        qualification: selectedQualification.value,
        search: searchQuery.value,
        page: nextPage,
        limit: PAGE_SIZE
      })
      jobs.value = [...jobs.value, ...(resp.data || [])]
      currentPage.value = nextPage
      hasNextPage.value = resp.meta?.hasNext ?? false
      totalJobs.value = resp.meta?.total ?? totalJobs.value
    } catch (err) {
      // fail silently on load-more
    } finally {
      loadingMore.value = false
    }
  }

  async function openDetails(job) {
    loading.value = true
    try {
      const data = await jobsApi.fetchJobDetails(job.id)
      selectedJob.value = data
      detailModalOpen.value = true
    } catch (err) {
      selectedJob.value = job
      detailModalOpen.value = true
    } finally {
      loading.value = false
    }
  }

  function closeDetails() {
    detailModalOpen.value = false
    selectedJob.value = null
  }

  /**
   * Track job application internally + redirect to official portal.
   * Returns the official apply_link from the backend.
   */
  async function applyToJob(jobId) {
    try {
      const resp = await jobsApi.applyToJob(jobId)
      return resp.apply_link
    } catch (err) {
      throw err
    }
  }

  return {
    jobs,
    loading,
    loadingMore,
    error,
    searchQuery,
    selectedQualification,
    selectedJob,
    detailModalOpen,
    qualifications,
    currentPage,
    totalJobs,
    hasNextPage,
    fetchJobs,
    loadMoreJobs,
    openDetails,
    closeDetails,
    applyToJob
  }
})
