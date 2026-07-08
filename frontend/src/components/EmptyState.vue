<script setup>
defineProps({
  title:       { type: String, default: 'Nothing here yet' },
  description: { type: String, default: '' },
  actionLabel: { type: String, default: '' },
  icon:        { type: String, default: 'search' } // search | bookmark | file | jobs
})
defineEmits(['action'])
</script>

<template>
  <div class="empty-state">
    <div class="empty-illustration">
      <!-- Search icon -->
      <svg v-if="icon === 'search'" xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"/>
        <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        <line x1="8" y1="11" x2="14" y2="11" stroke-width="1.5"/>
        <line x1="11" y1="8" x2="11" y2="14" stroke-width="1.5"/>
      </svg>

      <!-- Bookmark icon -->
      <svg v-else-if="icon === 'bookmark'" xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
      </svg>

      <!-- File icon -->
      <svg v-else-if="icon === 'file'" xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
        <polyline points="14 2 14 8 20 8"/>
        <line x1="16" y1="13" x2="8" y2="13"/>
        <line x1="16" y1="17" x2="8" y2="17"/>
        <polyline points="10 9 9 9 8 9"/>
      </svg>

      <!-- Jobs icon -->
      <svg v-else xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <rect x="2" y="7" width="20" height="14" rx="2" ry="2"/>
        <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
      </svg>
    </div>

    <h3 class="empty-title">{{ title }}</h3>
    <p v-if="description" class="empty-desc">{{ description }}</p>

    <button v-if="actionLabel" class="empty-action-btn" @click="$emit('action')">
      {{ actionLabel }}
    </button>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 32px;
  text-align: center;
  animation: fadeUp 0.4s ease both;
}

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(16px); }
  to   { opacity: 1; transform: translateY(0); }
}

.empty-illustration {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background: var(--clr-primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  color: var(--clr-primary);
}

.empty-title {
  font-family: var(--font-heading);
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--clr-text-main);
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 0.9rem;
  color: var(--clr-text-muted);
  max-width: 360px;
  line-height: 1.6;
  margin-bottom: 24px;
}

.empty-action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 10px 24px;
  background: var(--clr-primary);
  color: white;
  border: none;
  border-radius: 9999px;
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: 0 4px 14px var(--clr-primary-glow);
}

.empty-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 22px var(--clr-primary-glow);
}
</style>
