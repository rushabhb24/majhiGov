import { client } from './client.js'

export const bookmarksApi = {
  fetchSavedSchemes: () => client.get('/api/user/saved'),
  toggleSavedScheme: (schemeId) => client.post('/api/user/saved', { scheme_id: Number(schemeId) })
}
