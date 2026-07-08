import { client } from './client.js'
import { API_BASE_URL } from '../config.js'

export const aiApi = {
  analyzeResumeFile: (file) => {
    const formData = new FormData()
    formData.append('resume', file)
    return fetch(`${API_BASE_URL}/api/ai/resume-analyze`, {
      method: 'POST',
      credentials: 'include',
      body: formData
    }).then(async (r) => {
      if (!r.ok) {
        const text = await r.text()
        throw new Error(text || 'Resume upload failed')
      }
      return r.json()
    })
  },
  analyzeResumeText: (text) => client.post('/api/ai/resume-analyze', { resume_text: text }),
  getCareerAdvice: () => client.post('/api/ai/career-advisor', {}),
  analyzeSkillGap: (payload) => client.post('/api/ai/skill-gap', payload),
  getInterviewQuestions: (payload) => client.post('/api/ai/interview-prep', payload),
  generateCoverLetter: (payload) => client.post('/api/ai/cover-letter', payload),
  smartSearch: (query) => client.post('/api/ai/smart-search', { query })
}
