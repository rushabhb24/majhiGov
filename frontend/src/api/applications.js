import { client } from './client.js'

export const applicationsApi = {
  applyForScheme: (payload) => client.post('/api/user/apply', payload),
  fetchUserApplications: () => client.get('/api/user/applications')
}
