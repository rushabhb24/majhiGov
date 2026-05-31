import { client } from './client.js'

export const adminApi = {
  fetchAnalytics: () => client.get('/api/admin/analytics'),
  
  // Schemes Admin
  fetchAllSchemes: () => client.get('/api/admin/schemes'),
  createScheme: (payload) => client.post('/api/admin/schemes', payload),
  updateScheme: (schemeId, payload) => client.put(`/api/admin/schemes/${schemeId}`, payload),
  toggleSchemeStatus: (schemeId) => client.delete(`/api/admin/schemes/${schemeId}`),
  
  // Categories Admin
  fetchAllCategories: () => client.get('/api/admin/categories'),
  createCategory: (payload) => client.post('/api/admin/categories', payload),
  deleteCategory: (categoryId) => client.delete(`/api/admin/categories/${categoryId}`),
  
  // Users Admin
  fetchAllUsers: () => client.get('/api/admin/users'),
  toggleUserStatus: (userId) => client.post('/api/admin/users/toggle-active', { user_id: Number(userId) }),
  createAdmin: (payload) => client.post('/api/admin/users/admin', payload),
  
  // Notifications Admin
  fetchNotifications: () => client.get('/api/admin/notifications'),
  sendNotification: (payload) => client.post('/api/admin/notifications', payload),
  
  // Applications Admin [NEW]
  fetchAllApplications: () => client.get('/api/admin/applications'),
  updateApplicationStatus: (payload) => client.post('/api/admin/applications/status', payload)
}
