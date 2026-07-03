import { API_BASE_URL } from '../config.js'

/**
 * Centralized fetch wrapper.
 * - Always sends `credentials: 'include'` so httpOnly cookies are forwarded automatically.
 * - No Authorization header needed — JWT lives in a server-managed httpOnly cookie.
 * - Handles JSON parsing, non-2xx errors, and 204 No Content.
 */
async function request(endpoint, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {})
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include' // Required: sends httpOnly yojana_auth cookie on every request
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
  get:    (endpoint, options) => request(endpoint, { ...options, method: 'GET' }),
  post:   (endpoint, body, options) => request(endpoint, { ...options, method: 'POST', body: JSON.stringify(body) }),
  put:    (endpoint, body, options) => request(endpoint, { ...options, method: 'PUT', body: JSON.stringify(body) }),
  delete: (endpoint, options) => request(endpoint, { ...options, method: 'DELETE' })
}
