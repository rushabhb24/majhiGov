import { createI18n } from 'vue-i18n'
import en from '../locales/en.js'
import hi from '../locales/hi.js'
import mr from '../locales/mr.js'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    hi,
    mr
  }
})

export default i18n