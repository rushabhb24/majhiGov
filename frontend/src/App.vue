<script setup>
import { ref, onMounted, computed, watch } from 'vue';

// Import Reusable Modular Components
import Header from './components/Header.vue';
import Hero from './components/Hero.vue';
import SchemeExplorer from './components/SchemeExplorer.vue';
import SchemeCard from './components/SchemeCard.vue';
import EligibilityWizard from './components/EligibilityWizard.vue';
import EligibilityResults from './components/EligibilityResults.vue';
import SchemeDetailsModal from './components/SchemeDetailsModal.vue';
import ToastBanner from './components/ToastBanner.vue';

// Configuration
const API_BASE_URL = 'http://localhost:8080';

// Global Reactive States
const currentLanguage = ref('en'); // 'en', 'hi', 'mr'
const theme = ref('dark'); // 'dark', 'light'
const activeTab = ref('explorer'); // 'explorer', 'eligibility', 'saved'
const schemes = ref([]);
const loading = ref(false);
const error = ref(null);
const searchQu = ref('');
const selectedCategory = ref('All');
const sortBy = ref('date_desc');
const ruralMode = ref(false);

// Saved Bookmarks state
const savedSchemeIds = ref([]);

// Selected Relational Scheme state for modal
const selectedScheme = ref(null);
const detailModalOpen = ref(false);

// Categories definitions
const categories = ['All', 'Farmers', 'Students', 'Women', 'Senior Citizens', 'Business Owners'];

// Static translation dictionary for 3 languages
const t = computed(() => {
  const dictionary = {
    en: {
      explorer: 'Scheme Explorer',
      eligibility: 'Smart Eligibility',
      saved: 'Saved Schemes',
      ruralMode: 'Rural (Accessible) Mode',
      normalMode: 'Normal Mode',
      heroTitle: 'Sarkari Yojana Portal 🇮🇳',
      heroSubtitle: 'Discover central and state government schemes matching your profile and check eligibility in seconds.',
      btnCheckElig: 'Check Eligibility Now',
      searchLabel: 'Search Government Schemes / Scheme Dhundhein',
      searchPlaceholder: 'Jaise: Kisan scholarship loan...',
      sortByLabel: 'Sort By / Kram',
      recentFirst: 'Sabse Nayi (Recent First)',
      titleAlphabetical: 'Naam Ke Anusar (A-Z)',
      chooseCategory: 'Choose Category:',
      loading: 'Loading schemes from database...',
      retry: 'Retry Connection',
      benefitsLabel: 'Benefits',
      lastDateLabel: 'Last Application Date',
      viewDetails: 'Benefits and Documents',
      applyLink: 'Official Apply Link',
      personalProfile: '1. Personal Profile / Vyakti Gata Jankari',
      incomeOccupation: '2. Income & Occupation / Kamai aur Kaam',
      educationSpecial: '3. Education & Special Status / Shiksha aur Divyangta',
      ageLabel: 'Age (Umar)',
      genderLabel: 'Gender (Ling)',
      maleOpt: 'Male (Purush)',
      femaleOpt: 'Female (Mahila)',
      otherOpt: 'Other',
      stateLabel: 'State (Rajya)',
      casteLabel: 'Caste Category (Varg)',
      incomeLabel: 'Annual Family Income (Saalana Kamai)',
      occupationLabel: 'Occupation (Aap kya karte hain)',
      employeeTypeLabel: 'Employment Type (Job Ka Prakaar)',
      educationLabel: 'Highest Education Level (Shiksha)',
      disabilityLabel: 'Are you differently-abled? (Kya aap Divyang hain?)',
      next: 'Aage Badhein',
      back: 'Piche Jayein',
      calculate: 'Scheme Check Karein!',
      resultsTitle: 'Results / Aapki Yogyata',
      eligibleTitle: 'Eligible Schemes (Isme Apply Karein)',
      notEligibleTitle: 'Not Eligible Schemes',
      noEligible: 'No eligible schemes found for your current profile.',
      reasonsLabel: 'Why you qualify:',
      notEligReasonsLabel: 'Disqualifying reasons:',
      noDetailsChecked: 'Left panel me diye gaye form ko bharein aur check karein ki kaunsi scheme aapke liye bani hai.',
      notCheckedIntro: 'Abhi tak check nahi kiya gaya',
      savedTitle: 'Aapki Saved Schemes',
      savedSubtitle: 'Your bookmarked government schemes.',
      noSaved: 'No saved schemes yet. Go to Explorer and save some.',
      mandatoryBadge: 'Mandatory',
      optionalBadge: 'Optional',
      modalOverview: 'Scheme Overview',
      modalBenefits: 'Benefits',
      modalDocs: 'Required Documents',
      modalSteps: 'Step-by-Step Application Process',
      modalFaqs: 'Frequently Asked Questions',
      saveSchemeBtn: 'Save Scheme',
      removeSavedBtn: 'Remove Saved',
      directApplyBtn: 'Direct Official Apply Link ↗',
      allCategory: 'All',
      farmerCategory: 'Farmers',
      studentCategory: 'Students',
      womenCategory: 'Women',
      seniorCategory: 'Senior Citizens',
      businessCategory: 'Business Owners',
      toastSaved: 'Scheme saved to your profile!',
      toastRemoved: 'Scheme removed from your profile.',
      toastSuccess: 'Eligibility calculated successfully!',
      toastFail: 'Could not connect to Go backend.'
    },
    hi: {
      explorer: 'योजना खोजें',
      eligibility: 'योग्यता जांचें',
      saved: 'सुरक्षित योजनाएं',
      ruralMode: 'ग्रामीण (सुलभ) मोड',
      normalMode: 'सामान्य मोड',
      heroTitle: 'सरकारी योजना पोर्टल 🇮🇳',
      heroSubtitle: 'अपने प्रोफाइल के अनुसार सही केंद्रीय एवं राज्य सरकारी योजनाएं खोजें और पात्रता की जांच करें।',
      btnCheckElig: 'पात्रता की जांच करें',
      searchLabel: 'सरकारी योजनाएं खोजें',
      searchPlaceholder: 'योजना का नाम, कीवर्ड दर्ज करें...',
      sortByLabel: 'क्रमबद्ध करें',
      recentFirst: 'नवीनतम पहले',
      titleAlphabetical: 'नाम के अनुसार (A-Z)',
      chooseCategory: 'श्रेणी चुनें:',
      loading: 'डेटाबेस से योजनाएं लोड हो रही हैं...',
      retry: 'पुनः प्रयास करें',
      benefitsLabel: 'लाभ',
      lastDateLabel: 'अंतिम तिथि',
      viewDetails: 'फायदे और दस्तावेज देखें',
      applyLink: 'आधिकारिक आवेदन लिंक',
      personalProfile: '1. व्यक्तिगत जानकारी / Personal Profile',
      incomeOccupation: '2. आय और व्यवसाय / Income & Occupation',
      educationSpecial: '3. शिक्षा और विशेष स्थिति / Education & Disability',
      ageLabel: 'आयु (उम्र)',
      genderLabel: 'लिंग',
      maleOpt: 'पुरुष',
      femaleOpt: 'महिला',
      otherOpt: 'अन्य',
      stateLabel: 'राज्य',
      casteLabel: 'जाति श्रेणी (वर्ग)',
      incomeLabel: 'वार्षिक पारिवारिक आय',
      occupationLabel: 'व्यवसाय (आप क्या करते हैं)',
      employeeTypeLabel: 'रोजगार का प्रकार (नौकरी)',
      educationLabel: 'उच्चतम शैक्षणिक स्तर',
      disabilityLabel: 'क्या आप दिव्यांग हैं?',
      next: 'आगे बढ़ें',
      back: 'पीछे जाएं',
      calculate: 'पात्रता की जांच करें!',
      resultsTitle: 'परिणाम / आपकी योग्यता',
      eligibleTitle: 'योग्य योजनाएं (इसमें आवेदन करें)',
      notEligibleTitle: 'अयोग्य योजनाएं',
      noEligible: 'आपकी वर्तमान जानकारी के अनुसार कोई योग्य योजना नहीं मिली।',
      reasonsLabel: 'आप क्यों पात्र हैं:',
      notEligReasonsLabel: 'पात्र न होने के कारण:',
      noDetailsChecked: 'बाएं पैनल में दिए गए फॉर्म को भरें और जांचें कि कौन सी योजनाएं आपके लिए उपयुक्त हैं।',
      notCheckedIntro: 'अभी तक जांच नहीं की गई',
      savedTitle: 'आपकी सुरक्षित योजनाएं',
      savedSubtitle: 'आपके द्वारा सुरक्षित की गई सरकारी योजनाएं।',
      noSaved: 'अभी तक कोई योजना सुरक्षित नहीं की गई है। योजना खोजें और सुरक्षित करें।',
      mandatoryBadge: 'अनिवार्य',
      optionalBadge: 'वैकल्पिक',
      modalOverview: 'योजना का विवरण',
      modalBenefits: 'मिलने वाले लाभ',
      modalDocs: 'आवश्यक दस्तावेज',
      modalSteps: 'आवेदन करने की चरण-दर-चरण प्रक्रिया',
      modalFaqs: 'अक्सर पूछे जाने वाले प्रश्न',
      saveSchemeBtn: 'योजना सुरक्षित करें',
      removeSavedBtn: 'सुरक्षित सूची से हटाएं',
      directApplyBtn: 'सीधा आधिकारिक आवेदन लिंक ↗',
      allCategory: 'सभी',
      farmerCategory: 'किसान',
      studentCategory: 'विद्यार्थी',
      womenCategory: 'महिलाएं',
      seniorCategory: 'वरिष्ठ नागरिक',
      businessCategory: 'व्यवसायी',
      toastSaved: 'योजना सुरक्षित कर ली गई है!',
      toastRemoved: 'योजना सुरक्षित सूची से हटा दी गई है।',
      toastSuccess: 'पात्रता की गणना सफलतापूर्वक की गई!',
      toastFail: 'गो (Go) बैकएंड से कनेक्ट नहीं हो सका।'
    },
    mr: {
      explorer: 'योजना शोधा',
      eligibility: 'पात्रता तपासा',
      saved: 'जतन केलेल्या योजना',
      ruralMode: 'ग्रामीण (सुलभ) मोड',
      normalMode: 'सामान्य मोड',
      heroTitle: 'शासकीय योजनांचे पोर्टल 🇮🇳',
      heroSubtitle: 'तुमच्या प्रोफाइलनुसार योग्य केंद्रीय व राज्य शासकीय योजना शोधा आणि पात्रता तपासा.',
      btnCheckElig: 'पात्रता तपासा',
      searchLabel: 'शासकीय योजना शोधा',
      searchPlaceholder: 'योजनेचे नाव, मुख्य शब्द प्रविष्ट करा...',
      sortByLabel: 'क्रमवारी लावा',
      recentFirst: 'नवीनतम आधी',
      titleAlphabetical: 'नावानुसार (A-Z)',
      chooseCategory: 'श्रेणी निवडा:',
      loading: 'डेटाबेस मधून योजना लोड होत आहेत...',
      retry: 'पुन्हा प्रयत्न करा',
      benefitsLabel: 'मिळणारे फायदे',
      lastDateLabel: 'अंतिम तारीख',
      viewDetails: 'फायदे आणि कागदपत्रे पहा',
      applyLink: 'अधिकृत अर्ज लिंक',
      personalProfile: '1. वैयक्तिक माहिती / Personal Profile',
      incomeOccupation: '2. उत्पन्न आणि व्यवसाय / Income & Occupation',
      educationSpecial: '3. शिक्षण आणि विशेष स्थिती / Education & Disability',
      ageLabel: 'वय',
      genderLabel: 'लिंग',
      maleOpt: 'पुरुष',
      femaleOpt: 'महिला',
      otherOpt: 'इतर',
      stateLabel: 'राज्य',
      casteLabel: 'जातीचा प्रवर्ग (वर्ग)',
      incomeLabel: 'वार्षिक कौटुंबिक उत्पन्न',
      occupationLabel: 'व्यवसाय (तुम्ही काय करता)',
      employeeTypeLabel: 'रोजगाराचा प्रकार (नोकरी)',
      educationLabel: 'उच्चतम शैक्षणिक पातळी',
      disabilityLabel: 'तुम्ही दिव्यांग आहात का?',
      next: 'पुढे जा',
      back: 'मागे जा',
      calculate: 'पात्रता तपासा!',
      resultsTitle: 'निकाल / तुमची पात्रता',
      eligibleTitle: 'पात्र योजना (यात अर्ज करा)',
      notEligibleTitle: 'अपात्र योजना',
      noEligible: 'तुमच्या सध्याच्या प्रोफाइलनुसार एकही पात्र योजना आढळली नाही.',
      reasonsLabel: 'तुम्ही का पात्र आहात:',
      notEligReasonsLabel: 'अपात्रतेची कारणे:',
      noDetailsChecked: 'डाव्या बाजूला दिलेल्या फॉर्ममध्ये माहिती भरा आणि तुम्ही कोणत्या योजनांसाठी पात्र आहात ते तपासा.',
      notCheckedIntro: 'अद्याप तपासणी केली नाही',
      savedTitle: 'तुमच्या जतन केलेल्या योजना',
      savedSubtitle: 'तुम्ही जतन करून ठेवलेल्या शासकीय योजना.',
      noSaved: 'अद्याप कोणतीही योजना जतन केलेली नाही. योजना शोधून जतन करा.',
      mandatoryBadge: 'अनिवार्य',
      optionalBadge: 'पर्यायी',
      modalOverview: 'योजनेची माहिती',
      modalBenefits: 'मिळणारे फायदे',
      modalDocs: 'आवश्यक कागदपत्रे',
      modalSteps: 'अर्ज करण्याची टप्प्याटप्प्याने प्रक्रिया',
      modalFaqs: 'नेहमी विचारले जाणारे प्रश्न',
      saveSchemeBtn: 'योजना जतन करा',
      removeSavedBtn: 'जतन सूचीतून काढा',
      directApplyBtn: 'थेट अधिकृत अर्ज लिंक ↗',
      allCategory: 'सर्व',
      farmerCategory: 'शेतकरी',
      studentCategory: 'विद्यार्थी',
      womenCategory: 'महिला',
      seniorCategory: 'ज्येष्ठ नागरिक',
      businessCategory: 'व्यावसायिक',
      toastSaved: 'योजना यशस्वीरित्या जतन केली!',
      toastRemoved: 'योजना जतन सूचीतून काढली.',
      toastSuccess: 'पात्रतेची गणना यशस्वी झाली!',
      toastFail: 'गो (Go) बॅकएंडशी संपर्क होऊ शकला नाही.'
    }
  };
  return dictionary[currentLanguage.value];
});

// Eligibility Wizard Demographic Profile state
const eligibilityProfile = ref({
  age: 25,
  gender: 'Male',
  state: 'Maharashtra',
  caste: 'General',
  annual_income: 180000,
  occupation: 'Unemployed',
  employee_type: 'Unemployed',
  education_level: '12th Pass',
  is_disabled: false
});

const checkingEligibility = ref(false);
const eligibilityResults = ref(null);
const eligibilityChecked = ref(false);
const eligibilityStep = ref(1);

// Toast Notification Banners state
const toast = ref({
  show: false,
  message: '',
  type: 'success'
});

function showToast(msg, type = 'success') {
  toast.value.message = msg;
  toast.value.type = type;
  toast.value.show = true;
  setTimeout(() => {
    toast.value.show = false;
  }, 3000);
}

// Fetch active schemes from Go relational backend
async function fetchSchemes() {
  loading.value = true;
  error.value = null;
  try {
    const url = new URL(`${API_BASE_URL}/api/schemes`);
    if (selectedCategory.value !== 'All') {
      url.searchParams.append('category', selectedCategory.value);
    }
    if (searchQu.value) {
      url.searchParams.append('search', searchQu.value);
    }
    url.searchParams.append('sort_by', sortBy.value);

    const response = await fetch(url.toString());
    if (!response.ok) throw new Error('Failed to load schemes from server.');
    const data = await response.json();
    schemes.value = data || [];
  } catch (err) {
    console.error(err);
    error.value = t.value.toastFail;
    schemes.value = [];
  } finally {
    loading.value = false;
  }
}

// Bookmarks local storage integrations
function loadBookmarks() {
  const data = localStorage.getItem('yojana_saved_ids');
  if (data) {
    savedSchemeIds.value = JSON.parse(data);
  }
}

function toggleBookmark(schemeId) {
  const index = savedSchemeIds.value.indexOf(schemeId);
  if (index === -1) {
    savedSchemeIds.value.push(schemeId);
    showToast(t.value.toastSaved, 'success');
  } else {
    savedSchemeIds.value.splice(index, 1);
    showToast(t.value.toastRemoved, 'info');
  }
  localStorage.setItem('yojana_saved_ids', JSON.stringify(savedSchemeIds.value));
}

const bookmarkedSchemes = computed(() => {
  return schemes.value.filter(s => savedSchemeIds.value.includes(s.id));
});

// Dynamic Eligibility Computation
async function submitEligibility() {
  checkingEligibility.value = true;
  eligibilityResults.value = null;
  try {
    const response = await fetch(`${API_BASE_URL}/api/eligibility-check`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(eligibilityProfile.value)
    });
    if (!response.ok) throw new Error('Eligibility check failed.');
    const data = await response.json();
    eligibilityResults.value = data;
    eligibilityChecked.value = true;
    showToast(t.value.toastSuccess, 'success');
  } catch (err) {
    console.error(err);
    showToast(t.value.toastFail, 'danger');
  } finally {
    checkingEligibility.value = false;
  }
}

// Fetch Relational details for modal
async function openDetails(scheme) {
  loading.value = true;
  try {
    const response = await fetch(`${API_BASE_URL}/api/schemes/${scheme.id}`);
    if (!response.ok) throw new Error('Could not fetch details.');
    const data = await response.json();
    selectedScheme.value = data;
    detailModalOpen.value = true;
  } catch (err) {
    console.error(err);
    // Fallback to loaded local list properties
    selectedScheme.value = scheme;
    detailModalOpen.value = true;
  } finally {
    loading.value = false;
  }
}

function closeDetails() {
  detailModalOpen.value = false;
  selectedScheme.value = null;
}

// Watchers
watch([selectedCategory, sortBy], () => {
  fetchSchemes();
});

// Search input debouncer
let searchTimeout;
watch(searchQu, () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    fetchSchemes();
  }, 400);
});

onMounted(() => {
  fetchSchemes();
  loadBookmarks();
});
</script>

<template>
  <div :class="['app-wrapper', { 'rural-mode': ruralMode }, theme]">
    <!-- Header component (Logo, selects, tabs, togglers) -->
    <Header 
      v-model:activeTab="activeTab"
      v-model:currentLanguage="currentLanguage"
      v-model:ruralMode="ruralMode"
      v-model:theme="theme"
      :saved-count="savedSchemeIds.length"
      :t="t"
    />

    <!-- Main Viewport Shell -->
    <main class="main-container">
      
      <!-- Premium Hero Headline banner -->
      <Hero :t="t" @start-check="activeTab = 'eligibility'" />

      <!-- TAB VIEW 1: SCHEME EXPLORER -->
      <SchemeExplorer 
        v-if="activeTab === 'explorer'"
        v-model:selectedCategory="selectedCategory"
        v-model:sortBy="sortBy"
        v-model:searchQu="searchQu"
        :schemes="schemes"
        :loading="loading"
        :error="error"
        :current-language="currentLanguage"
        :saved-scheme-ids="savedSchemeIds"
        :categories="categories"
        :t="t"
        @toggle-bookmark="toggleBookmark"
        @open-details="openDetails"
        @retry="fetchSchemes"
      />

      <!-- TAB VIEW 2: SMART ELIGIBILITY CHECKER -->
      <div v-else-if="activeTab === 'eligibility'" class="tab-content">
        <div class="grid-layout">
          <!-- step forms -->
          <EligibilityWizard 
            v-model:step="eligibilityStep"
            :profile="eligibilityProfile"
            :t="t"
            :checking="checkingEligibility"
            @submit="submitEligibility"
          />

          <!-- results lists -->
          <EligibilityResults 
            :checked="eligibilityChecked"
            :results="eligibilityResults"
            :current-language="currentLanguage"
            :t="t"
            @open-details="openDetails"
          />
        </div>
      </div>

      <!-- TAB VIEW 3: SAVED BOOKMARKS -->
      <div v-else-if="activeTab === 'saved'" class="tab-content animate-fade">
        <div class="card filter-panel">
          <h2 class="section-title">{{ t.savedTitle }}</h2>
          <p class="text-muted">{{ t.savedSubtitle }}</p>
        </div>

        <div v-if="bookmarkedSchemes.length === 0" class="empty-state text-center mt-4 card">
          <div class="empty-bookmarks-art">🔖</div>
          <h3>{{ t.noSaved }}</h3>
          <button class="btn btn-primary mt-4" @click="activeTab = 'explorer'">Explore Schemes</button>
        </div>

        <!-- Schemes Grid for Saved bookmarks using reusable SchemeCard -->
        <div v-else class="schemes-grid mt-4">
          <SchemeCard 
            v-for="scheme in bookmarkedSchemes" 
            :key="scheme.id"
            :scheme="scheme"
            :current-language="currentLanguage"
            :saved-scheme-ids="savedSchemeIds"
            :t="t"
            @toggle-bookmark="toggleBookmark"
            @open-details="openDetails"
          />
        </div>
      </div>
    </main>

    <!-- Details relational modal overlay (Acc FAQ + Docs lists) -->
    <SchemeDetailsModal 
      :scheme="selectedScheme"
      :current-language="currentLanguage"
      :saved-scheme-ids="savedSchemeIds"
      :open="detailModalOpen"
      :t="t"
      @close="closeDetails"
      @toggle-bookmark="toggleBookmark"
    />

    <!-- Frosted Notification banner alerts -->
    <ToastBanner 
      :show="toast.show"
      :message="toast.message"
      :type="toast.type"
    />
  </div>
</template>

<style scoped>
.animate-fade {
  animation: fadeIn 0.4s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}
.empty-bookmarks-art {
  font-size: 3.5rem;
  margin-bottom: 12px;
  filter: drop-shadow(0 6px 10px rgba(0,0,0,0.05));
}
</style>
