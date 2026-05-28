<script setup>
import { computed } from 'vue'

const props = defineProps({
  analytics: {
    type: Object,
    default: () => ({})
  }
})

const stats = computed(() => {
  const sa = props.analytics || {}
  const totalApp = sa.total_applications || 24
  const pendingApp = sa.pending_applications || 7
  const approvedApp = sa.approved_applications || 12
  const rejectedApp = sa.rejected_applications || 5

  return [
    { num: totalApp, label: 'Total Applications', desc: 'Total database count', icon: 'ti-checklist', class: 'blue' },
    { num: pendingApp, label: 'Pending Approvals', desc: 'Awaiting review', icon: 'ti-clock', class: 'orange' },
    { num: approvedApp, label: 'Approved Applications', desc: 'Successfully verified', icon: 'ti-check', class: 'green' },
    { num: rejectedApp, label: 'Rejected Applications', desc: 'Ineligible submissions', icon: 'ti-close', class: 'red' }
  ]
})

const rankings = computed(() => {
  return [
    { rank: 1, name: 'PM Kisan Samman Nidhi', count: '14,200', fill: 100 },
    { rank: 2, name: 'NSP Post-Matric Scholarship', count: '11,080', fill: 78 },
    { rank: 3, name: 'PM Mudra Yojana', count: '8,520', fill: 60 },
    { rank: 4, name: 'Ladli Behna Yojana', count: '6,840', fill: 48 },
    { rank: 5, name: 'Ayushman Bharat Yojana', count: '4,970', fill: 35 }
  ]
})
</script>

<template>
  <div class="analytics-tab">
    
    <!-- Stats Row — 4 cards -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.label">
        <div :class="['stat-icon-box', stat.class]">
          <i :class="['ti', stat.icon]"></i>
        </div>
        <div class="stat-number">{{ stat.num }}</div>
        <div class="stat-label">{{ stat.label }}</div>
        <div class="stat-desc text-success">
          {{ stat.desc }}
        </div>
      </div>
    </div>

    <!-- Top 5 Schemes Card (full width) -->
    <div class="card mt-4">
      <div class="card-header">
        <div class="card-title">Top 5 Most Applied Schemes</div>
      </div>
      <div class="card-body">
        <div class="ranking-list">
          <div class="ranking-row" v-for="r in rankings" :key="r.rank">
            <span class="rank-number">{{ r.rank }}</span>
            <span class="rank-name">{{ r.name }}</span>
            <div class="bar-container">
              <!-- Rank 1 is orange (accent), 2-5 are navy (primary) -->
              <div 
                :class="['bar-fill', r.rank === 1 ? 'accent' : 'primary']" 
                :style="{ width: r.fill + '%' }"
              ></div>
            </div>
            <span class="rank-count">{{ r.count }}</span>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.analytics-tab {
  width: 100%;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  width: 100%;
}

.stat-card {
  background-color: #ffffff; /* var(--bg) */
  border: 0.5px solid rgba(0, 0, 0, 0.08); /* var(--border) */
  border-radius: 8px; /* var(--radius) */
  padding: 14px 16px;
  box-sizing: border-box;
}

.stat-icon-box {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.stat-icon-box.blue { background-color: #e8eef8; color: #1a3a6b; }
.stat-icon-box.green { background-color: #f0fdf4; color: #16a34a; }
.stat-icon-box.orange { background-color: #fff4ed; color: #f97316; }
.stat-icon-box.red { background-color: #fef2f2; color: #dc2626; }

.stat-icon-box i {
  font-size: 18px !important;
}

.stat-number {
  font-size: 22px;
  font-weight: 500;
  color: #0f172a;
  margin-top: 10px;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.stat-desc {
  font-size: 11px;
  font-weight: 500;
  margin-top: 8px;
}

.text-success { color: #16a34a; }

.mt-4 { margin-top: 16px; }

.card {
  background-color: #ffffff;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  overflow: hidden;
  width: 100%;
}

.card-header {
  padding: 14px 16px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: #0f172a;
}

.card-body {
  padding: 16px;
}

/* Rankings */
.ranking-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.ranking-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.rank-number {
  font-size: 13px;
  font-weight: 500;
  width: 20px;
  color: #0f172a;
  flex-shrink: 0;
}

.rank-name {
  font-size: 13px;
  font-weight: 500;
  color: #0f172a;
  width: 200px;
  flex-shrink: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bar-container {
  flex-grow: 1;
  height: 4px;
  background-color: #f8fafc; /* var(--bg2) */
  border-radius: 2px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 2px;
}

.bar-fill.accent {
  background-color: #f97316; /* var(--accent) */
}

.bar-fill.primary {
  background-color: #1a3a6b; /* var(--primary) */
}

.rank-count {
  width: 60px;
  text-align: right;
  font-weight: 500;
  font-size: 13px;
  color: #64748b;
  flex-shrink: 0;
}
</style>
