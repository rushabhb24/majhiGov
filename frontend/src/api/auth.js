import { client } from './client.js'

export const authApi = {
  login:                 (credentials) => client.post('/api/auth/login', credentials),
  register:              (payload) => client.post('/api/auth/register', payload),
  logout:                () => client.post('/api/auth/logout', {}),
  fetchProfile:          () => client.get('/api/user/profile'),
  updateProfile:         (profileData) => client.put('/api/user/profile', profileData),
  fetchNotifications:    () => client.get('/api/user/notifications'),
  markNotificationsRead: () => client.put('/api/user/notifications/read', {}),
  getRecommendations:    () => client.get('/api/user/recommendations')
}
