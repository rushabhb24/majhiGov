import { createApp } from 'vue'
import { createPinia } from 'pinia'
import i18n from './i18n'
import router from './router'
import App from './App.vue'
import './tailwind.css'
import './style.css'

const app = createApp(App)
app.use(createPinia())
app.use(i18n)
app.use(router)
app.mount('#app')

// Unregister active Service Workers in dev to prevent caching conflicts & extension script interception bugs
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.getRegistrations().then((registrations) => {
    for (const registration of registrations) {
      registration.unregister().then(() => {
        // Clear caches to remove any stale/cached index.html or API responses
        caches.keys().then((keys) => {
          return Promise.all(keys.map(key => caches.delete(key)))
        })
      })
    }
  })
}
