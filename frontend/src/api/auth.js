import { client } from './client.js'

export const authApi = {
  login: (credentials) => client.post('/api/auth/login', credentials),
  register: (payload) => client.post('/api/auth/register', payload),
  fetchProfile: () => client.get('/api/user/profile'),
  updateProfile: (profileData) => client.put('/api/user/profile', profileData)
}
