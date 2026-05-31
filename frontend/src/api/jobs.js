import { client } from './client.js'

export const jobsApi = {
  fetchPublicJobs: (params = {}) => {
    let endpoint = '/api/jobs'
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
    
    if (queryParts.length > 0) {
      endpoint += '?' + queryParts.join('&')
    }
    
    return client.get(endpoint)
  },
  fetchJobDetails: (id) => client.get(`/api/jobs/${id}`),
  
  // Admin endpoints
  adminFetchAllJobs: () => client.get('/api/admin/jobs'),
  adminCreateJob: (payload) => client.post('/api/admin/jobs', payload),
  adminUpdateJob: (id, payload) => client.put(`/api/admin/jobs/${id}`, payload),
  adminDeleteJob: (id) => client.delete(`/api/admin/jobs/${id}`)
}
