import { defineStore } from 'pinia'
import { ref } from 'vue'
import { aiApi } from '../api/ai.js'

export const useAiStore = defineStore('ai', () => {
  const resumeResult = ref(null)
  const careerResult = ref(null)
  const skillGapResult = ref(null)
  const interviewResult = ref(null)
  const coverLetterResult = ref(null)

  const loading = ref(false)
  const error = ref(null)
  const activeTab = ref('resume')

  async function analyzeResume(fileOrText) {
    loading.value = true
    error.value = null
    resumeResult.value = null
    try {
      if (typeof fileOrText === 'string') {
        resumeResult.value = await aiApi.analyzeResumeText(fileOrText)
      } else {
        resumeResult.value = await aiApi.analyzeResumeFile(fileOrText)
      }
    } catch (err) {
      error.value = err.message || 'Failed to analyze resume'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getCareerAdvice() {
    loading.value = true
    error.value = null
    careerResult.value = null
    try {
      careerResult.value = await aiApi.getCareerAdvice()
    } catch (err) {
      error.value = err.message || 'Failed to get career roadmap advice'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function analyzeSkillGap(payload) {
    loading.value = true
    error.value = null
    skillGapResult.value = null
    try {
      skillGapResult.value = await aiApi.analyzeSkillGap(payload)
    } catch (err) {
      error.value = err.message || 'Failed to analyze skill gap'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getInterviewQuestions(payload) {
    loading.value = true
    error.value = null
    interviewResult.value = null
    try {
      interviewResult.value = await aiApi.getInterviewQuestions(payload)
    } catch (err) {
      error.value = err.message || 'Failed to generate interview questions'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function generateCoverLetter(payload) {
    loading.value = true
    error.value = null
    coverLetterResult.value = null
    try {
      coverLetterResult.value = await aiApi.generateCoverLetter(payload)
    } catch (err) {
      error.value = err.message || 'Failed to generate cover letter'
      throw err
    } finally {
      loading.value = false
    }
  }

  function resetAll() {
    resumeResult.value = null
    careerResult.value = null
    skillGapResult.value = null
    interviewResult.value = null
    coverLetterResult.value = null
    error.value = null
  }

  return {
    resumeResult,
    careerResult,
    skillGapResult,
    interviewResult,
    coverLetterResult,
    loading,
    error,
    activeTab,
    analyzeResume,
    getCareerAdvice,
    analyzeSkillGap,
    getInterviewQuestions,
    generateCoverLetter,
    resetAll
  }
})
