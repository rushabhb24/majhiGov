import { client } from './client.js'

/**
 * Builds a query string from params, appending page/limit for pagination.
 * All API responses for list endpoints now return { data: [...], meta: { page, limit, total, hasNext } }
 */
function buildSchemesQuery(params = {}) {
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
  if (params.page) {
    queryParts.push(`page=${params.page}`)
  }
  if (params.limit) {
    queryParts.push(`limit=${params.limit}`)
  }
  const qs = queryParts.length > 0 ? '?' + queryParts.join('&') : ''
  return `/api/schemes${qs}`
}

export const schemesApi = {
  /**
   * Returns paginated response: { data: Scheme[], meta: { page, limit, total, hasNext } }
   * Pass params.page and params.limit (default 5) for infinite scroll pagination.
   */
  fetchPublicSchemes: (params = {}) => client.get(buildSchemesQuery(params)),
  fetchSchemeDetails: (id) => client.get(`/api/schemes/${id}`),
  checkEligibility: (payload) => client.post('/api/eligibility-check', payload),
  translateText: (text, target) => client.get(`/api/translate?q=${encodeURIComponent(text)}&target=${target}`)
}
