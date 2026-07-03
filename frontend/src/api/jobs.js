import { client } from './client.js'

function buildJobsQuery(params = {}) {
  const queryParts = []
  if (params.qualification && params.qualification !== 'All') {
    queryParts.push(`qualification=${encodeURIComponent(params.qualification)}`)
  }
  if (params.search) {
    queryParts.push(`search=${encodeURIComponent(params.search)}`)
  }
  if (params.organization) {
    queryParts.push(`organization=${encodeURIComponent(params.organization)}`)
  }
  if (params.page) {
    queryParts.push(`page=${params.page}`)
  }
  if (params.limit) {
    queryParts.push(`limit=${params.limit}`)
  }
  const qs = queryParts.length > 0 ? '?' + queryParts.join('&') : ''
  return `/api/jobs${qs}`
}

export const jobsApi = {
  /**
   * Returns paginated response: { data: Job[], meta: { page, limit, total, hasNext } }
   */
  fetchPublicJobs: (params = {}) => client.get(buildJobsQuery(params)),
  fetchJobDetails: (id) => client.get(`/api/jobs/${id}`),

  /**
   * Log internal apply tracking + receive official apply_link for external redirect.
   * Dual action: track in DB AND redirect user to govt portal.
   */
  applyToJob: (jobId) => client.post('/api/user/apply-job', { job_id: jobId }),
  fetchJobApplications: () => client.get('/api/user/job-applications'),

  // Admin endpoints
  adminFetchAllJobs: () => client.get('/api/admin/jobs'),
  adminCreateJob: (payload) => client.post('/api/admin/jobs', payload),
  adminUpdateJob: (id, payload) => client.put(`/api/admin/jobs/${id}`, payload),
  adminDeleteJob: (id) => client.delete(`/api/admin/jobs/${id}`)
}
