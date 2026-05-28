<script setup>
import { computed } from 'vue'

const props = defineProps({
  schemes: {
    type: Array,
    required: true
  },
  categories: {
    type: Array,
    required: true
  },
  filterCategory: {
    type: String,
    required: true
  },
  filterType: {
    type: String,
    required: true
  },
  filterStatus: {
    type: String,
    required: true
  }
})

const emit = defineEmits([
  'update:filterCategory',
  'update:filterType',
  'update:filterStatus',
  'view-scheme',
  'edit-scheme',
  'delete-scheme'
])

const categoryEmojiMap = computed(() => {
  const map = {}
  props.categories.forEach(c => {
    map[c.name] = c.icon
  })
  // Fallbacks
  map['Farmers'] = '🌾'
  map['Students'] = '🎓'
  map['Women'] = '👩'
  map['Business Owners'] = '💼'
  map['Business'] = '💼'
  map['Senior Citizens'] = '👴'
  return map
})

const getCategoryEmoji = (catName) => {
  return categoryEmojiMap.value[catName] || '🌾'
}
</script>

<template>
  <div class="schemes-tab">
    
    <!-- Filter Row (above table) -->
    <div class="filter-row">
      <div class="filter-group">
        <select 
          :value="filterCategory" 
          @change="emit('update:filterCategory', $event.target.value)"
        >
          <option value="All">All Categories</option>
          <option value="Farmers">Farmers</option>
          <option value="Students">Students</option>
          <option value="Women">Women</option>
          <option value="Business Owners">Business Owners</option>
          <option value="Senior Citizens">Senior Citizens</option>
        </select>

        <select 
          :value="filterType" 
          @change="emit('update:filterType', $event.target.value)"
        >
          <option value="All">All Types</option>
          <option value="central">Central</option>
          <option value="state">State</option>
        </select>

        <select 
          :value="filterStatus" 
          @change="emit('update:filterStatus', $event.target.value)"
        >
          <option value="All">All Statuses</option>
          <option value="Active">Active</option>
          <option value="Inactive">Inactive</option>
          <option value="Expiring">Expiring</option>
        </select>
      </div>
    </div>

    <!-- Schemes Table (full width card) -->
    <div class="card mt-3">
      <div class="card-body p-0">
        <table class="data-table">
          <thead>
            <tr>
              <th>Scheme Name</th>
              <th>Category</th>
              <th>Type</th>
              <th>Deadline</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in schemes" :key="s.id">
              <td>
                <div class="scheme-title-text">{{ s.title }}</div>
                <div class="scheme-subtitle-text">{{ s.benefits || s.description }}</div>
              </td>
              <td>
                <span class="cat-pill">
                  {{ getCategoryEmoji(s.category_name) }} {{ s.category_name }}
                </span>
              </td>
              <td>
                <span :class="['badge', s.government_level === 'central' ? 'central' : 'state']">
                  {{ s.government_level === 'central' ? 'Central' : 'State' + (s.state ? ' - ' + s.state : '') }}
                </span>
              </td>
              <td>{{ s.application_end_date }}</td>
              <td>
                <span :class="['badge', s.is_active ? 'active' : 'inactive']">
                  {{ s.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td>
                <div class="table-actions">
                  <button class="action-btn" title="View" @click="emit('view-scheme', s)">
                    <i class="ti ti-eye"></i>
                  </button>
                  <button class="action-btn" title="Edit" @click="emit('edit-scheme', s)">
                    <i class="ti ti-edit"></i>
                  </button>
                  <button class="action-btn danger-hover" title="Delete" @click="emit('delete-scheme', s)">
                    <i class="ti ti-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="schemes.length === 0">
              <td colspan="6" style="text-align: center; padding: 30px; color: #64748b;">
                No schemes match the selected filters.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>
</template>

<style scoped>
.schemes-tab {
  width: 100%;
}

.filter-row {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}

.filter-group {
  display: flex;
  gap: 10px;
}

.filter-group select {
  padding: 7px 10px;
  border: 0.5px solid rgba(0, 0, 0, 0.08); /* var(--border) */
  border-radius: 6px; /* var(--radius) */
  font-size: 13px;
  background-color: #ffffff;
  color: #0f172a;
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
}

.filter-group select:focus {
  border-color: #1a3a6b;
}

.mt-3 { margin-top: 12px; }

.card {
  background-color: #ffffff;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  overflow: hidden;
  width: 100%;
}

.p-0 { padding: 0 !important; }

/* Table Styling */
.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  background-color: #f8fafc; /* var(--bg2) */
  text-align: left;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #64748b;
  padding: 8px 12px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
}

.data-table td {
  padding: 10px 12px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
  font-size: 13px;
  vertical-align: middle;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table tr:hover td {
  background-color: #f8fafc;
}

.scheme-title-text {
  font-weight: 500;
  color: #0f172a;
}

.scheme-subtitle-text {
  font-size: 11px;
  color: #64748b;
  margin-top: 2px;
}

.cat-pill {
  font-size: 13px;
  color: #0f172a;
}

/* Badges */
.badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 8px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
  line-height: 1;
}

.badge.active { background-color: #f0fdf4; color: #16a34a; }
.badge.inactive { background-color: #f8fafc; color: #64748b; }
.badge.central { background-color: #e8eef8; color: #1a3a6b; }
.badge.state { background-color: #fff4ed; color: #f97316; }

/* Table Actions */
.table-actions {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: 0.5px solid rgba(0, 0, 0, 0.08); /* var(--border) */
  background-color: #ffffff;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 15px;
  padding: 0;
  box-sizing: border-box;
}

.action-btn i {
  font-size: 15px !important;
}

.action-btn:hover {
  background-color: #f8fafc;
  color: #0f172a;
}

.action-btn.danger-hover:hover {
  background-color: #fef2f2;
  color: #dc2626;
  border-color: #dc2626;
}
</style>
