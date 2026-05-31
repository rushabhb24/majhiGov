import { client } from './client.js'

export const schemesApi = {
  fetchPublicSchemes: (params = {}) => {
    let endpoint = '/api/schemes'
    const queryParts = []
    
    if (params.category && params.category !== 'All') {
      queryParts.push(`category=${encodeURIComponent(params.category)}`)
    }
    if (params.search) {
      queryParts.push(`search=${encodeURIComponent(params.search)}`)
    }
    if (params.sort_by) {
      queryParts.push(`sort_by=${encodeURIComponent(params.sort_by)}`)
    }
    
    if (queryParts.length > 0) {
      endpoint += '?' + queryParts.join('&')
    }
    
    return client.get(endpoint)
  },
  fetchSchemeDetails: (id) => client.get(`/api/schemes/${id}`),
  checkEligibility: (payload) => client.post('/api/eligibility-check', payload),
  translateText: (text, target) => client.get(`/api/translate?q=${encodeURIComponent(text)}&target=${target}`)
}
