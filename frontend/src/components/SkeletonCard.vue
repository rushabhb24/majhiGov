<template>
  <div class="skeleton-card">
    <div class="skeleton-top-row">
      <div class="skeleton-line w-24 h-6"></div>
      <div class="skeleton-circle"></div>
    </div>
    <div class="skeleton-title"></div>
    <div class="skeleton-line"></div>
    <div class="skeleton-line w-75"></div>
    <div class="skeleton-benefit"></div>
    <div class="skeleton-footer">
      <div class="skeleton-btn"></div>
      <div class="skeleton-btn"></div>
    </div>
  </div>
</template>

<style scoped>
@keyframes shimmer {
  0%   { background-position: -200% 0; }
  100% { background-position:  200% 0; }
}

.skeleton-card {
  background: var(--clr-surface);
  border-radius: 16px;
  padding: 24px;
  border: 1px solid var(--clr-border);
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-height: 260px;
}

/* Shimmer mixin for all skeleton elements */
[class*="skeleton-line"],
[class*="skeleton-title"],
[class*="skeleton-circle"],
[class*="skeleton-benefit"],
[class*="skeleton-btn"],
.skeleton-top-row [class*="skeleton-"] {
  background: linear-gradient(
    90deg,
    var(--clr-surface-alt) 25%,
    rgba(255, 255, 255, 0.08) 50%,
    var(--clr-surface-alt) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s ease-in-out infinite;
  border-radius: 6px;
}

.skeleton-top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.skeleton-line {
  height: 14px;
  width: 100%;
  border-radius: 6px;
}

.w-24 { width: 96px; }
.w-75 { width: 75%; }
.h-6  { height: 24px; }

.skeleton-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  flex-shrink: 0;
}

.skeleton-title {
  height: 22px;
  width: 85%;
  border-radius: 6px;
}

.skeleton-benefit {
  height: 72px;
  border-radius: 10px;
}

.skeleton-footer {
  display: flex;
  gap: 12px;
  margin-top: auto;
  padding-top: 8px;
}

.skeleton-btn {
  height: 40px;
  border-radius: 10px;
  flex: 1;
}
</style>
