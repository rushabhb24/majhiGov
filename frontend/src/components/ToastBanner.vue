<script setup>
defineProps({
  show: {
    type: Boolean,
    required: true
  },
  message: {
    type: String,
    required: true
  },
  type: {
    type: String,
    default: 'success' // success, info, danger
  }
});
</script>

<template>
  <Transition name="toast-fade">
    <div v-if="show" :class="['toast-banner', type]">
      <div class="toast-content">
        <span class="toast-icon">
          <span v-if="type === 'success'">✔️</span>
          <span v-else-if="type === 'danger'">⚠️</span>
          <span v-else>💡</span>
        </span>
        <span>{{ message }}</span>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.toast-fade-enter-active {
  animation: slideIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.toast-fade-leave-active {
  animation: slideIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) reverse;
}
</style>
