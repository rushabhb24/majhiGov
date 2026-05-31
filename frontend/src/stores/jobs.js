import { defineStore } from 'pinia'
import { ref } from 'vue'
import { jobsApi } from '../api/jobs.js'

export const useJobStore = defineStore('jobs', () => {
  // State
  const jobs = ref([])
  const loading = ref(false)
  const error = ref(null)
  const searchQuery = ref('')
  const selectedQualification = ref('All')
  const selectedJob = ref(null)
  const detailModalOpen = ref(false)
  const qualifications = ['All', '10th Pass', '12th Pass', 'Graduate']

  // Actions
  async function fetchJobs() {
    loading.value = true
    error.value = null
    try {
      const data = await jobsApi.fetchPublicJobs({
        qualification: selectedQualification.value,
        search: searchQuery.value
      })
      jobs.value = data || []
    } catch (err) {
      console.error(err)
      error.value = 'Could not fetch government jobs.'
      jobs.value = []
    } finally {
      loading.value = false
    }
  }

  async function openDetails(job) {
    loading.value = true
    try {
      const data = await jobsApi.fetchJobDetails(job.id)
      selectedJob.value = data
      detailModalOpen.value = true
    } catch (err) {
      console.error(err)
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

  return {
    jobs,
    loading,
    error,
    searchQuery,
    selectedQualification,
    selectedJob,
    detailModalOpen,
    qualifications,
    fetchJobs,
    openDetails,
    closeDetails
  }
})
