<script setup>
import { ref, computed } from 'vue'
import { useAiStore } from '../stores/ai'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'
import { aiApi } from '../api/ai.js'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppLabel from '../components/ui/AppLabel.vue'
import AppTextarea from '../components/ui/AppTextarea.vue'
import AppSelect from '../components/ui/AppSelect.vue'

const aiStore = useAiStore()
const authStore = useAuthStore()
const uiStore = useUiStore()

const resumeText = ref('')
const selectedFile = ref(null)
const dragActive = ref(false)

// Tab state mapping
const tabs = [
  { id: 'resume', name: 'Resume Analyzer', icon: '📄' },
  { id: 'roadmap', name: 'AI Career Advisor', icon: 'Advisor' },
  { id: 'skillgap', name: 'Skill Gap Analysis', icon: '📊' },
  { id: 'interview', name: 'Interview Prep', icon: '🎤' },
  { id: 'coverletter', name: 'Cover Letter Maker', icon: '✍️' }
]

// Tab 3 form
const skillGapForm = ref({
  job_title: '',
  job_description: '',
  required_skills: ''
})

// Tab 4 form
const interviewForm = ref({
  job_title: '',
  job_description: '',
  experience_years: 0,
  interview_type: 'mixed'
})
const visibleAnswers = ref({})

// Tab 5 form
const coverLetterForm = ref({
  company_name: '',
  job_title: '',
  job_description: '',
  tone: 'professional'
})

// Conversational AI Chat states
const chatMessages = ref([
  { sender: 'assistant', text: 'Hello! I am your AI Career Assistant. Ask me anything about schemes, eligibility, or career roadmaps.' }
])
const chatInput = ref('')
const chatLoading = ref(false)
const chatSuggestions = ['Schemes for Students', 'Farmers subsidies', 'Resume tips', 'OBC Scholarship']

// Drag and drop handlers
function handleDragOver(e) {
  e.preventDefault()
  dragActive.value = true
}

function handleDragLeave() {
  dragActive.value = false
}

function handleDrop(e) {
  e.preventDefault()
  dragActive.value = false
  if (e.dataTransfer.files && e.dataTransfer.files[0]) {
    const file = e.dataTransfer.files[0]
    validateAndSetFile(file)
  }
}

function handleFileSelect(e) {
  if (e.target.files && e.target.files[0]) {
    validateAndSetFile(e.target.files[0])
  }
}

function validateAndSetFile(file) {
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  if (ext !== '.pdf' && ext !== '.txt') {
    uiStore.showToast('Please upload a PDF or TXT file only', 'warning')
    return
  }
  selectedFile.value = file
  uiStore.showToast(`Selected file: ${file.name}`, 'info')
}

function clearSelectedFile() {
  selectedFile.value = null
}

async function triggerResumeAnalysis() {
  if (!selectedFile.value && !resumeText.value) {
    uiStore.showToast('Please upload a file or paste resume text', 'warning')
    return
  }
  try {
    if (selectedFile.value) {
      await aiStore.analyzeResume(selectedFile.value)
    } else {
      await aiStore.analyzeResume(resumeText.value)
    }
    uiStore.showToast('Resume analyzed successfully!', 'success')
  } catch (err) {
    uiStore.showToast(err.message || 'Analysis failed', 'danger')
  }
}

async function triggerCareerAdvisor() {
  try {
    await aiStore.getCareerAdvice()
    uiStore.showToast('Roadmap generated!', 'success')
  } catch (err) {
    uiStore.showToast(err.message || 'Generation failed', 'danger')
  }
}

async function triggerSkillGap() {
  if (!skillGapForm.value.job_title || !skillGapForm.value.job_description) {
    uiStore.showToast('Job Title and Description are required', 'warning')
    return
  }
  try {
    const skillsList = skillGapForm.value.required_skills
      ? skillGapForm.value.required_skills.split(',').map(s => s.trim()).filter(s => s !== '')
      : []
    await aiStore.analyzeSkillGap({
      job_title: skillGapForm.value.job_title,
      job_description: skillGapForm.value.job_description,
      required_skills: skillsList
    })
    uiStore.showToast('Skill gap analysis complete!', 'success')
  } catch (err) {
    uiStore.showToast(err.message || 'Analysis failed', 'danger')
  }
}

async function triggerInterviewPrep() {
  if (!interviewForm.value.job_title) {
    uiStore.showToast('Job Title is required', 'warning')
    return
  }
  try {
    visibleAnswers.value = {}
    await aiStore.getInterviewQuestions({
      job_title: interviewForm.value.job_title,
      job_description: interviewForm.value.job_description,
      experience_years: Number(interviewForm.value.experience_years),
      interview_type: interviewForm.value.interview_type
    })
    uiStore.showToast('Interview prep pack generated!', 'success')
  } catch (err) {
    uiStore.showToast(err.message || 'Generation failed', 'danger')
  }
}

async function triggerCoverLetter() {
  if (!coverLetterForm.value.company_name || !coverLetterForm.value.job_title) {
    uiStore.showToast('Company Name and Job Title are required', 'warning')
    return
  }
  try {
    await aiStore.generateCoverLetter(coverLetterForm.value)
    uiStore.showToast('Cover letter created!', 'success')
  } catch (err) {
    uiStore.showToast(err.message || 'Generation failed', 'danger')
  }
}

function copyToClipboard(text) {
  navigator.clipboard.writeText(text)
  uiStore.showToast('Copied to clipboard!', 'success')
}

function downloadTxt(filename, text) {
  const blob = new Blob([text], { type: 'text/plain;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', filename)
  link.click()
  URL.revokeObjectURL(url)
}

function toggleAnswer(key) {
  visibleAnswers.value[key] = !visibleAnswers.value[key]
}

// AI Assistant Chat submit logic
async function sendChatMessage(text) {
  const msg = text || chatInput.value
  if (!msg.trim()) return
  
  chatMessages.value.push({ sender: 'user', text: msg })
  chatInput.value = ''
  chatLoading.value = true
  
  try {
    const res = await aiApi.smartSearch(msg)
    let reply = `Based on your request, I found some relevant schemes and jobs:\n\n`
    if (res && res.schemes && res.schemes.length > 0) {
      reply += `**Schemes:**\n`
      res.schemes.slice(0, 3).forEach(s => {
        reply += `- **${s.title}**: ${s.description.substring(0, 80)}... [Learn More](${s.apply_link || s.official_website})\n`
      })
    }
    if (res && res.jobs && res.jobs.length > 0) {
      reply += `\n**Jobs:**\n`
      res.jobs.slice(0, 3).forEach(j => {
        reply += `- **${j.title}** at ${j.organization} [Apply on Official Website](${j.official_website})\n`
      })
    }
    if ((!res.schemes || res.schemes.length === 0) && (!res.jobs || res.jobs.length === 0)) {
      reply = `I searched active schemes and jobs for "${msg}" but couldn't find a direct match. Try asking about general educational grants, farmers subsidies, or specific tech roles!`
    }
    chatMessages.value.push({ sender: 'assistant', text: reply })
  } catch (err) {
    chatMessages.value.push({ sender: 'assistant', text: 'Sorry, I encountered an error searching schemes. Please try again.' })
  } finally {
    chatLoading.value = false
  }
}

function clearChatHistory() {
  chatMessages.value = [
    { sender: 'assistant', text: 'Chat history cleared. How can I help you today?' }
  ]
}
</script>

<template>
  <div class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- Hero Banner -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-8">
      <h1 class="tw-font-heading tw-font-bold tw-text-2xl tw-text-foreground tw-m-0">
        AI Career Co-Pilot &amp; Advisor
      </h1>
      <p class="tw-text-xs tw-text-muted-foreground tw-mt-2 tw-m-0 tw-line-height-[1.5]">
        Optimize resumes, align skill profiles, draft cover letters, and prepare for interviews using our deep Gemini-powered career modules.
      </p>
    </div>

    <!-- main split layout -->
    <div class="tw-grid tw-grid-cols-1 lg:tw-grid-cols-12 tw-gap-8">
      
      <!-- Left side: Career tools -->
      <div class="lg:tw-col-span-7 tw-flex tw-flex-col tw-gap-6">
        
        <!-- Tab selector -->
        <div class="tw-flex tw-gap-1.5 tw-bg-muted/60 tw-p-1 tw-rounded-xl tw-overflow-x-auto tw-scrollbar-none">
          <button
            v-for="t in tabs"
            :key="t.id"
            class="tw-border-none tw-outline-none tw-flex-1 tw-flex tw-items-center tw-justify-center tw-gap-1.5 tw-py-2 tw-px-3 tw-text-xs tw-font-bold tw-font-heading tw-rounded-lg tw-transition-all tw-cursor-pointer"
            :class="aiStore.activeTab === t.id ? 'tw-bg-card tw-text-foreground tw-shadow-sm' : 'tw-text-muted-foreground hover:tw-text-foreground tw-bg-transparent'"
            @click="aiStore.activeTab = t.id"
          >
            <span>{{ t.name }}</span>
          </button>
        </div>

        <!-- Tool panel card -->
        <div class="glass tw-p-6 tw-rounded-2xl">
          
          <!-- Global Loader -->
          <div v-if="aiStore.loading" class="tw-text-center tw-py-12 tw-flex tw-flex-col tw-items-center tw-gap-3">
            <div class="tw-w-8 tw-h-8 tw-border-3 tw-border-muted tw-border-t-primary tw-rounded-full tw-animate-spin"></div>
            <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">AI is analyzing...</h3>
            <p class="tw-text-xs tw-text-muted-foreground tw-m-0">Please wait while we run our Generative Intelligence models</p>
          </div>

          <!-- Resume Analyzer -->
          <div v-else-if="aiStore.activeTab === 'resume'" class="tw-flex tw-flex-col tw-gap-6">
            <div>
              <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">Upload Resume</h3>
              <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">Upload your resume in PDF or TXT format to calculate ATS score, keyword matches, and improvement tips.</p>
            </div>

            <div 
              class="tw-border-2 tw-border-dashed tw-border-border tw-rounded-2xl tw-p-8 tw-text-center tw-cursor-pointer hover:tw-border-primary/50 tw-transition-colors"
              :class="{ 'tw-border-primary/50': dragActive }"
              @dragover="handleDragOver"
              @dragleave="handleDragLeave"
              @drop="handleDrop"
              @click="triggerResumeUpload"
            >
              <input type="file" ref="resumeInput" class="tw-hidden" accept=".pdf,.txt" @change="handleFileSelect" />
              <div v-if="!selectedFile" class="tw-flex tw-flex-col tw-items-center tw-gap-1.5">
                <span class="tw-text-2xl">☁️</span>
                <span class="tw-text-xs tw-font-bold tw-text-foreground">Drag &amp; drop file here or Browse File</span>
              </div>
              <div v-else class="tw-flex tw-items-center tw-justify-between tw-bg-muted/50 tw-p-2 tw-rounded-xl tw-text-xs">
                <div class="tw-flex tw-items-center tw-gap-2">
                  <span>📄</span>
                  <span class="tw-font-bold tw-text-foreground">{{ selectedFile.name }}</span>
                </div>
                <button class="tw-border-none tw-bg-transparent tw-text-destructive tw-cursor-pointer" @click.stop="clearSelectedFile">Remove</button>
              </div>
            </div>

            <div class="tw-text-center tw-text-[10px] tw-font-bold tw-text-muted-foreground tw-tracking-widest tw-uppercase tw-my-1">— OR PASTE TEXT —</div>

            <div>
              <AppTextarea 
                v-model="resumeText" 
                placeholder="Paste your plain resume text here..." 
                rows="6"
                :disabled="selectedFile !== null"
              />
            </div>

            <AppButton variant="primary" size="sm" @click="triggerResumeAnalysis">
              Analyze Resume
            </AppButton>

            <!-- Results -->
            <div v-if="aiStore.resumeResult" class="tw-border-t tw-border-border/50 tw-pt-6 tw-flex tw-flex-col tw-gap-4">
              <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">Analysis Report</h3>
              
              <div class="tw-flex tw-items-center tw-gap-4 tw-bg-muted/40 tw-p-4 tw-rounded-xl">
                <div class="tw-w-16 tw-h-16 tw-rounded-full tw-bg-primary tw-text-white tw-font-display tw-font-black tw-text-xl tw-flex tw-items-center tw-justify-center">
                  {{ aiStore.resumeResult.ats_score }}%
                </div>
                <div>
                  <span class="tw-text-xs tw-text-muted-foreground tw-block">ATS Match Score</span>
                  <AppBadge :tone="aiStore.resumeResult.strength === 'strong' ? 'success' : 'info'">
                    {{ aiStore.resumeResult.strength }} Resume
                  </AppBadge>
                </div>
              </div>

              <div>
                <AppLabel>Summary</AppLabel>
                <div class="tw-text-xs tw-text-muted-foreground tw-line-height-[1.5]">{{ aiStore.resumeResult.summary }}</div>
              </div>

              <div>
                <AppLabel>Skills Identified</AppLabel>
                <div class="tw-flex tw-flex-wrap tw-gap-1.5 tw-mt-1">
                  <AppBadge v-for="skill in aiStore.resumeResult.skills_found" :key="skill" tone="neutral">
                    {{ skill }}
                  </AppBadge>
                </div>
              </div>

              <div>
                <AppLabel>Missing Keywords</AppLabel>
                <div class="tw-flex tw-flex-wrap tw-gap-1.5 tw-mt-1">
                  <AppBadge v-for="kw in aiStore.resumeResult.missing_keywords" :key="kw" tone="danger">
                    {{ kw }}
                  </AppBadge>
                </div>
              </div>

              <div>
                <AppLabel>Suggestions for Improvement</AppLabel>
                <ul class="tw-list-disc tw-pl-4 tw-text-xs tw-text-muted-foreground tw-flex tw-flex-col tw-gap-1 tw-mt-1">
                  <li v-for="sug in aiStore.resumeResult.suggestions" :key="sug">{{ sug }}</li>
                </ul>
              </div>
            </div>
          </div>

          <!-- AI Career Advisor -->
          <div v-else-if="aiStore.activeTab === 'roadmap'" class="tw-flex tw-flex-col tw-gap-6">
            <div>
              <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">AI Career Advisor</h3>
              <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">Generate a custom timeline, roadmap, and training certificates based on your user profile statistics.</p>
            </div>
            
            <AppButton variant="primary" size="sm" @click="triggerCareerAdvisor">
              Generate My Career Roadmap
            </AppButton>

            <div v-if="aiStore.careerResult" class="tw-border-t tw-border-border/50 tw-pt-6 tw-flex tw-flex-col tw-gap-5">
              
              <!-- Roles -->
              <div>
                <AppLabel>Recommended Roles</AppLabel>
                <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-3 tw-mt-2">
                  <div v-for="role in aiStore.careerResult.suitable_roles" :key="role.title" class="tw-p-3 tw-bg-muted/40 tw-rounded-xl tw-border tw-border-border">
                    <div class="tw-flex tw-justify-between tw-items-center tw-mb-1.5">
                      <span class="tw-font-heading tw-font-bold tw-text-xs tw-text-foreground">{{ role.title }}</span>
                      <AppBadge tone="success">{{ role.match_score }}% Match</AppBadge>
                    </div>
                    <p class="tw-text-[11px] tw-text-muted-foreground tw-m-0">{{ role.description }}</p>
                  </div>
                </div>
              </div>

              <!-- Learning roadmap -->
              <div>
                <AppLabel>Chronological Learning Roadmap</AppLabel>
                <div class="tw-flex tw-flex-col tw-gap-4 tw-mt-2">
                  <div v-for="(phase, idx) in aiStore.careerResult.roadmap" :key="idx" class="tw-flex tw-gap-3">
                    <div class="tw-w-6 tw-h-6 tw-rounded-full tw-bg-primary/10 tw-text-primary tw-font-bold tw-text-xs tw-flex tw-items-center tw-justify-center tw-flex-shrink-0">
                      {{ idx + 1 }}
                    </div>
                    <div>
                      <strong class="tw-text-xs tw-text-foreground">{{ phase.phase }} ({{ phase.duration }})</strong>
                      <ul class="tw-list-disc tw-pl-4 tw-text-[11px] tw-text-muted-foreground tw-mt-1">
                        <li v-for="milestone in phase.milestones" :key="milestone">{{ milestone }}</li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Accreditations & Resources -->
              <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
                <div class="tw-p-4 tw-rounded-xl tw-bg-muted/40 tw-border tw-border-border">
                  <h5 class="tw-font-heading tw-font-bold tw-text-xs tw-text-foreground tw-mt-0 tw-mb-2">Accredited Certifications</h5>
                  <ul class="tw-list-none tw-p-0 tw-m-0 tw-flex tw-flex-col tw-gap-2 tw-text-[11px] tw-text-muted-foreground">
                    <li v-for="cert in aiStore.careerResult.certifications" :key="cert.name">
                      <strong>{{ cert.name }}</strong> ({{ cert.provider }})
                    </li>
                  </ul>
                </div>
                <div class="tw-p-4 tw-rounded-xl tw-bg-muted/40 tw-border tw-border-border">
                  <h5 class="tw-font-heading tw-font-bold tw-text-xs tw-text-foreground tw-mt-0 tw-mb-2">Recommended Resources</h5>
                  <ul class="tw-list-none tw-p-0 tw-m-0 tw-flex tw-flex-col tw-gap-2 tw-text-[11px] tw-text-muted-foreground">
                    <li v-for="res in aiStore.careerResult.learning_resources" :key="res.name">
                      <a :href="res.url" target="_blank" class="tw-text-primary tw-font-bold tw-no-underline hover:tw-underline">{{ res.name }}</a>
                      <span class="tw-bg-muted tw-px-1.5 tw-py-0.5 tw-rounded tw-ml-2">{{ res.type }}</span>
                    </li>
                  </ul>
                </div>
              </div>

            </div>
          </div>

          <!-- Skill Gap -->
          <div v-else-if="aiStore.activeTab === 'skillgap'" class="tw-flex tw-flex-col tw-gap-4">
            <div>
              <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">Skill Gap Analysis</h3>
              <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">Insert a job post title and description to compare your qualifications against the listing requirements.</p>
            </div>
            
            <div>
              <AppLabel for="gap-title">Job Title</AppLabel>
              <AppInput id="gap-title" v-model="skillGapForm.job_title" placeholder="e.g. Frontend Developer" />
            </div>

            <div>
              <AppLabel for="gap-skills">Required Skills (Comma separated)</AppLabel>
              <AppInput id="gap-skills" v-model="skillGapForm.required_skills" placeholder="React, Pinia, CSS, Node" />
            </div>

            <div>
              <AppLabel for="gap-desc">Job Description</AppLabel>
              <AppTextarea id="gap-desc" v-model="skillGapForm.job_description" placeholder="Paste job description..." rows="5" />
            </div>

            <AppButton variant="primary" size="sm" @click="triggerSkillGap">
              Analyze Skill Gap
            </AppButton>

            <div v-if="aiStore.skillGapResult" class="tw-border-t tw-border-border/50 tw-pt-6 tw-flex tw-flex-col tw-gap-4">
              <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">Comparison Report</h3>

              <div class="tw-flex tw-items-center tw-gap-4 tw-bg-muted/40 tw-p-4 tw-rounded-xl">
                <div class="tw-w-16 tw-h-16 tw-rounded-full tw-bg-primary tw-text-white tw-font-display tw-font-black tw-text-xl tw-flex tw-items-center tw-justify-center">
                  {{ aiStore.skillGapResult.match_score }}%
                </div>
                <div>
                  <span class="tw-text-xs tw-text-muted-foreground tw-block">Required Skills Match</span>
                  <span class="tw-text-xs tw-text-primary tw-font-bold">Skill Gap Map ready</span>
                </div>
              </div>

              <div>
                <AppLabel>Matching Skills</AppLabel>
                <div class="tw-flex tw-flex-wrap tw-gap-1.5 tw-mt-1">
                  <AppBadge v-for="skill in aiStore.skillGapResult.matching_skills" :key="skill" tone="success">
                    {{ skill }}
                  </AppBadge>
                </div>
              </div>

              <div>
                <AppLabel>Missing Skills to Learn</AppLabel>
                <div class="tw-flex tw-flex-col tw-gap-2 tw-mt-1">
                  <div v-for="skill in aiStore.skillGapResult.missing_skills" :key="skill.skill" class="tw-flex tw-justify-between tw-items-center tw-p-2.5 tw-bg-muted/40 tw-rounded-xl tw-text-xs">
                    <div>
                      <strong class="tw-text-foreground">{{ skill.skill }}</strong>
                      <span class="tw-text-[10px] tw-text-muted-foreground tw-ml-2">Est: {{ skill.estimated_learning_time }}</span>
                    </div>
                    <AppBadge :tone="skill.priority === 'high' ? 'danger' : 'info'">
                      {{ skill.priority }} Priority
                    </AppBadge>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Interview Prep -->
          <div v-else-if="aiStore.activeTab === 'interview'" class="tw-flex tw-flex-col tw-gap-4">
            <div>
              <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">Interview Practice Generator</h3>
              <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">Generate custom technical, HR, or coding questions tailored to the position details.</p>
            </div>

            <div>
              <AppLabel for="int-title">Job Title</AppLabel>
              <AppInput id="int-title" v-model="interviewForm.job_title" placeholder="e.g. Frontend Engineer" />
            </div>

            <div>
              <AppLabel for="int-exp">Years of Experience</AppLabel>
              <AppInput id="int-exp" type="number" v-model="interviewForm.experience_years" min="0" max="30" />
            </div>

            <div>
              <AppLabel for="int-type">Interview Type</AppLabel>
              <AppSelect id="int-type" v-model="interviewForm.interview_type">
                <option value="hr">HR / Fit Round</option>
                <option value="technical">Technical</option>
                <option value="behavioral">Behavioral (STAR)</option>
                <option value="mixed">Mixed Selection Pack</option>
              </AppSelect>
            </div>

            <AppButton variant="primary" size="sm" @click="triggerInterviewPrep">
              Generate Questions
            </AppButton>

            <div v-if="aiStore.interviewResult" class="tw-border-t tw-border-border/50 tw-pt-6 tw-flex tw-flex-col tw-gap-4">
              <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">Practice Question Set</h3>
              
              <div v-if="aiStore.interviewResult.technical_questions && aiStore.interviewResult.technical_questions.length > 0" class="tw-flex tw-flex-col tw-gap-3">
                <div v-for="(q, idx) in aiStore.interviewResult.technical_questions" :key="'tech-'+idx" class="tw-p-4 tw-rounded-xl tw-bg-muted/40 tw-border tw-border-border tw-text-xs">
                  <div class="tw-font-semibold tw-text-foreground">Q{{ idx + 1 }}: {{ q.question }} ({{ q.difficulty }})</div>
                  <button class="tw-bg-transparent tw-border-none tw-text-primary tw-font-bold tw-cursor-pointer tw-mt-2 hover:tw-underline" @click="toggleAnswer('tech-'+idx)">
                    {{ visibleAnswers['tech-'+idx] ? 'Hide Answer' : 'Show Answer' }}
                  </button>
                  <p v-if="visibleAnswers['tech-'+idx]" class="tw-mt-2 tw-text-muted-foreground tw-m-0">{{ q.expected_answer }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Cover Letter Maker -->
          <div v-else-if="aiStore.activeTab === 'coverletter'" class="tw-flex tw-flex-col tw-gap-4">
            <div>
              <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">Cover Letter Draft Creator</h3>
              <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">Draft a fully tailored cover letter matching your profile and company details.</p>
            </div>

            <div>
              <AppLabel for="cl-comp">Company Name</AppLabel>
              <AppInput id="cl-comp" v-model="coverLetterForm.company_name" placeholder="e.g. Infosys Ltd" />
            </div>

            <div>
              <AppLabel for="cl-role">Job Title</AppLabel>
              <AppInput id="cl-role" v-model="coverLetterForm.job_title" placeholder="e.g. Software Associate" />
            </div>

            <div>
              <AppLabel for="cl-tone">Writing Tone</AppLabel>
              <AppSelect id="cl-tone" v-model="coverLetterForm.tone">
                <option value="professional">Professional & Formal</option>
                <option value="enthusiastic">Enthusiastic & Driven</option>
                <option value="creative">Creative & Out-of-box</option>
              </AppSelect>
            </div>

            <AppButton variant="primary" size="sm" @click="triggerCoverLetter">
              Generate Cover Letter
            </AppButton>

            <div v-if="aiStore.coverLetterResult" class="tw-border-t tw-border-border/50 tw-pt-6 tw-flex tw-flex-col tw-gap-3">
              <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">Draft Cover Letter</h3>
              <div class="tw-p-4 tw-rounded-xl tw-bg-muted/40 tw-border tw-border-border tw-text-xs">
                <div><strong>Subject:</strong> {{ aiStore.coverLetterResult.subject_line }}</div>
                <pre class="tw-whitespace-pre-wrap tw-font-body tw-text-muted-foreground tw-mt-3 tw-m-0">{{ aiStore.coverLetterResult.cover_letter_text }}</pre>
              </div>
            </div>
          </div>

        </div>

      </div>

      <!-- Right side: Chat interface -->
      <div class="lg:tw-col-span-5 tw-flex tw-flex-col tw-h-[580px] glass tw-rounded-2xl tw-p-4">
        
        <!-- Header -->
        <div class="tw-flex tw-justify-between tw-items-center tw-pb-3 tw-border-b tw-border-border/50">
          <div class="tw-flex tw-items-center tw-gap-2">
            <div class="tw-w-2 tw-h-2 tw-rounded-full tw-bg-emerald-400 tw-animate-pulse"></div>
            <span class="tw-font-heading tw-font-bold tw-text-xs tw-text-foreground">AI Career Advisor Chat</span>
          </div>
          <button 
            class="tw-bg-transparent tw-border-none tw-text-muted-foreground hover:tw-text-destructive tw-cursor-pointer tw-text-xs tw-font-bold"
            @click="clearChatHistory"
          >
            Clear History
          </button>
        </div>

        <!-- Chat messages scroller -->
        <div class="tw-flex-1 tw-overflow-y-auto tw-my-3 tw-flex tw-flex-col tw-gap-3 tw-pr-1">
          <div 
            v-for="(msg, idx) in chatMessages" 
            :key="idx" 
            class="tw-p-3 tw-rounded-2xl tw-max-w-[85%] tw-text-xs tw-line-height-[1.5]"
            :class="msg.sender === 'user' ? 'tw-bg-primary tw-text-white tw-self-end tw-rounded-tr-none' : 'tw-bg-muted/60 tw-text-foreground tw-self-start tw-rounded-tl-none'"
          >
            <!-- Render basic markdown-like list formatting -->
            <div class="tw-whitespace-pre-line">{{ msg.text }}</div>
          </div>

          <!-- Streaming loader -->
          <div v-if="chatLoading" class="tw-self-start tw-bg-muted/60 tw-p-3 tw-rounded-2xl tw-rounded-tl-none tw-flex tw-items-center tw-gap-1.5 tw-text-xs">
            <span class="tw-w-1.5 tw-h-1.5 tw-bg-muted-foreground tw-rounded-full tw-animate-bounce"></span>
            <span class="tw-w-1.5 tw-h-1.5 tw-bg-muted-foreground tw-rounded-full tw-animate-bounce" style="animation-delay: 0.2s"></span>
            <span class="tw-w-1.5 tw-h-1.5 tw-bg-muted-foreground tw-rounded-full tw-animate-bounce" style="animation-delay: 0.4s"></span>
          </div>
        </div>

        <!-- Suggestions Chips -->
        <div class="tw-flex tw-flex-wrap tw-gap-1.5 tw-mb-3">
          <button 
            v-for="chip in chatSuggestions" 
            :key="chip"
            class="tw-border tw-border-border tw-bg-card tw-text-muted-foreground hover:tw-text-foreground hover:tw-border-primary/30 tw-text-[10px] tw-font-semibold tw-px-2 tw-py-1 tw-rounded-full tw-cursor-pointer tw-transition-colors"
            @click="sendChatMessage(chip)"
          >
            {{ chip }}
          </button>
        </div>

        <!-- Message input form -->
        <form @submit.prevent="sendChatMessage()" class="tw-flex tw-gap-2 tw-border-t tw-border-border/50 tw-pt-3">
          <input 
            v-model="chatInput" 
            type="text" 
            placeholder="Ask about schemes, eligibility, jobs..." 
            class="tw-flex-1 tw-px-3 tw-py-2 tw-text-xs tw-bg-muted/50 tw-text-foreground tw-border tw-border-border tw-rounded-lg tw-outline-none focus:tw-border-primary tw-transition-colors"
            :disabled="chatLoading"
          />
          <AppButton 
            type="submit" 
            variant="primary" 
            size="sm"
            :disabled="chatLoading"
          >
            Send
          </AppButton>
        </form>

      </div>

    </div>

  </div>
</template>
