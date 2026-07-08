import { defineStore } from 'pinia'
import { ref } from 'vue'
import { API_BASE_URL } from '../config.js'

export const usePrivateJobStore = defineStore('privateJobs', () => {
  const jobs = ref([])
  const loading = ref(false)
  const loadingMore = ref(false)
  const error = ref(null)
  const activeTab = ref('all') // 'all', 'govt', 'private', 'internship', 'walkin', 'hackathon'
  
  const page = ref(1)
  const hasNext = ref(false)
  const total = ref(0)
  
  async function fetchJobs(search = '', location = '', reset = true) {
    if (reset) {
      page.value = 1
      jobs.value = []
    }
    
    // 'govt' is managed by the main jobs store, so we skip it here
    if (activeTab.value === 'govt') return
    
    loading.value = reset
    loadingMore.value = !reset
    error.value = null
    
    try {
      let url = `${API_BASE_URL}/api/private-jobs?page=${page.value}&limit=6`
      if (activeTab.value !== 'all') {
        url += `&job_type=${activeTab.value}`
      }
      if (search) {
        url += `&search=${encodeURIComponent(search)}`
      }
      if (location) {
        url += `&location=${encodeURIComponent(location)}`
      }
      
      const resp = await fetch(url, { credentials: 'include' })
      if (!resp.ok) throw new Error('Failed to fetch private jobs')
      
      const res = await resp.json()
      if (reset) {
        jobs.value = res.data || []
      } else {
        jobs.value = [...jobs.value, ...(res.data || [])]
      }
      
      total.value = res.meta?.total || 0
      hasNext.value = res.meta?.hasNext || false
    } catch (err) {
      error.value = err.message || 'Error loading private jobs'
    } finally {
      loading.value = false
      loadingMore.value = false
    }
  }

  async function loadMore(search = '', location = '') {
    if (loadingMore.value || !hasNext.value) return
    page.value++
    await fetchJobs(search, location, false)
  }

  async function applyToPrivateJob(jobId) {
    try {
      const resp = await fetch(`${API_BASE_URL}/api/user/apply-private-job`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ private_job_id: jobId })
      })
      if (!resp.ok) {
        const text = await resp.text()
        throw new Error(text || 'Failed to apply')
      }
      return await resp.json() // returns {success, apply_link, message}
    } catch (err) {
      throw err
    }
  }

  return {
    jobs,
    loading,
    loadingMore,
    error,
    activeTab,
    hasNext,
    total,
    fetchJobs,
    loadMore,
    applyToPrivateJob
  }
})
