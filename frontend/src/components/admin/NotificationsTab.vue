<script setup>
import { computed } from 'vue'

const props = defineProps({
  broadcast: {
    type: Object,
    required: true
  },
  notifications: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['send-broadcast'])

const recentList = computed(() => {
  if (props.notifications.length > 0) {
    return props.notifications
  }

  // Exact mock list from md file
  return [
    { type: 'Deadline Reminder', title: 'Deadline Reminder — Ladli Behna Yojana', target: 'Sent to 4,230 women users', time_ago: '2 hours ago', dot: 'orange' },
    { type: 'New Scheme Alert', title: 'New Scheme Alert — PM Vishwakarma Yojana', target: 'Sent to 18,430 users', time_ago: '1 day ago', dot: 'green' },
    { type: 'Deadline Reminder', title: 'Deadline Reminder — NSP Scholarship 2024', target: 'Sent to 2,100 student users', time_ago: '3 days ago', dot: 'blue' }
  ]
})

function getDotClass(type) {
  if (type === 'Deadline Reminder') return 'orange'
  if (type === 'New Scheme Alert') return 'green'
  return 'blue'
}
</script>

<template>
  <div class="notifications-tab">
    
    <div class="two-col-grid">
      
      <!-- LEFT CARD — "Send Notification" -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">Send Notification</div>
        </div>
        <div class="card-body">
          <form @submit.prevent="emit('send-broadcast')">
            <div class="form-group">
              <label class="form-label">Send To</label>
              <select class="form-input" v-model="broadcast.send_to" required>
                <option value="All Users">All Users</option>
                <option value="All Farmers">All Farmers</option>
                <option value="All Students">All Students</option>
                <option value="Specific State">Specific State</option>
              </select>
            </div>

            <div class="form-group" v-if="broadcast.send_to === 'Specific State'">
              <label class="form-label">Target State</label>
              <input 
                type="text" 
                class="form-input" 
                v-model="broadcast.state" 
                placeholder="e.g. Maharashtra" 
                required 
              />
            </div>

            <div class="form-group">
              <label class="form-label">Title</label>
              <input 
                type="text" 
                class="form-input" 
                v-model="broadcast.title" 
                placeholder="Scheme deadline reminder" 
                required 
              />
            </div>

            <div class="form-group">
              <label class="form-label">Message</label>
              <textarea 
                class="form-input" 
                v-model="broadcast.message" 
                rows="4" 
                placeholder="Dear citizen, Ladli Behna Yojana..." 
                required
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Type</label>
              <select class="form-input" v-model="broadcast.type" required>
                <option value="Deadline Reminder">Deadline Reminder</option>
                <option value="New Scheme Alert">New Scheme Alert</option>
                <option value="System Update">System Update</option>
              </select>
            </div>

            <button type="submit" class="submit-btn">
              <i class="ti ti-send"></i>
              <span>Send Notification</span>
            </button>
          </form>
        </div>
      </div>

      <!-- RIGHT CARD — "Recent Notifications Sent" -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">Recent Notifications Sent</div>
        </div>
        <div class="card-body p-0">
          <div class="activity-feed">
            <div class="activity-item" v-for="(n, idx) in recentList" :key="idx">
              <span :class="['dot-indicator', getDotClass(n.type)]"></span>
              <div class="activity-content">
                <div class="activity-title">{{ n.title }}</div>
                <div class="activity-subtext">{{ n.target || ('Sent to users · ' + n.time_ago) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

  </div>
</template>

<style scoped>
.notifications-tab {
  width: 100%;
}

.two-col-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  align-items: start;
}

.card {
  background-color: #ffffff;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  overflow: hidden;
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

.p-0 { padding: 0 !important; }

/* Forms */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.form-label {
  font-size: 13px;
  color: #0f172a;
}

.form-input {
  padding: 8px 10px;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  font-size: 13px;
  background-color: #ffffff;
  color: #0f172a;
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
  width: 100%;
}

.form-input:focus {
  border-color: #1a3a6b;
}

textarea.form-input {
  resize: vertical;
}

.submit-btn {
  background-color: #1a3a6b; /* var(--primary) */
  color: #ffffff;
  border: none;
  border-radius: 6px;
  padding: 9px 18px;
  font-size: 13px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  width: 100%;
  font-family: inherit;
  box-sizing: border-box;
  margin-top: 8px;
}

.submit-btn:hover {
  opacity: 0.9;
}

.submit-btn i {
  font-size: 16px !important;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* Activity Feed */
.activity-feed {
  display: flex;
  flex-direction: column;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
  align-items: flex-start;
  box-sizing: border-box;
}

.activity-item:last-child {
  border-bottom: none;
}

.dot-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 5px;
  flex-shrink: 0;
}

.dot-indicator.orange { background-color: #f97316; } /* var(--accent) */
.dot-indicator.green { background-color: #16a34a; } /* var(--success) */
.dot-indicator.blue { background-color: #1a3a6b; } /* var(--primary) */

.activity-content {
  flex-grow: 1;
}

.activity-title {
  font-size: 13px;
  font-weight: 500;
  color: #0f172a;
}

.activity-subtext {
  font-size: 11px;
  color: #64748b;
  margin-top: 2px;
}
</style>
