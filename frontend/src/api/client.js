import { API_BASE_URL } from '../config.js'
import { useAuthStore } from '../stores/auth.js'

async function request(endpoint, options = {}) {
  const authStore = useAuthStore()
  
  const headers = {
    'Content-Type': 'application/json',
    ...(authStore.token ? { 'Authorization': `Bearer ${authStore.token}` } : {}),
    ...(options.headers || {})
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers
  })

  if (!response.ok) {
    const errText = await response.text()
    let errMsg = 'API Request Failed'
    try {
      const parsed = JSON.parse(errText)
      errMsg = parsed.error || errMsg
    } catch (e) {
      errMsg = errText || errMsg
    }
    throw new Error(errMsg)
  }

  if (response.status === 204) return null
  return await response.json()
}

export const client = {
  get: (endpoint, options) => request(endpoint, { ...options, method: 'GET' }),
  post: (endpoint, body, options) => request(endpoint, { ...options, method: 'POST', body: JSON.stringify(body) }),
  put: (endpoint, body, options) => request(endpoint, { ...options, method: 'PUT', body: JSON.stringify(body) }),
  delete: (endpoint, options) => request(endpoint, { ...options, method: 'DELETE' })
}
