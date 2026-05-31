import { client } from './client.js'

export const eligibilityApi = {
  checkEligibility: (profileData) => client.post('/api/eligibility-check', profileData)
}
