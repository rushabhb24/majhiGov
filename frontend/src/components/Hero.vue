<script setup>
defineProps({
  t: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['startCheck']);
</script>

<template>
  <section class="hero animate-fade-in">
    <h1>{{ t.heroTitle }}</h1>
    <p>{{ t.heroSubtitle }}</p>
    <div class="hero-actions">
      <button class="btn btn-accent btn-hero-cta" @click="emit('startCheck')">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 11 12 14 22 4"></polyline><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path></svg>
        {{ t.btnCheckElig }}
      </button>
    </div>
  </section>
</template>

<style scoped>
.animate-fade-in {
  animation: slideIn 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}
.btn-hero-cta {
  padding: 14px 32px;
  font-size: 1.05rem;
}
</style>
