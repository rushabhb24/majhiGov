<script setup>
defineProps({
  categories: {
    type: Array,
    required: true
  },
  newCategory: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['add-category', 'delete-category'])
</script>

<template>
  <div class="categories-tab">
    
    <!-- Two equal columns side by side -->
    <div class="two-col-grid">
      
      <!-- LEFT CARD — "Add New Category" -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">Add New Category</div>
        </div>
        <div class="card-body">
          <form @submit.prevent="emit('add-category')">
            <div class="form-group">
              <label class="form-label">Category Name (English)</label>
              <input 
                type="text" 
                class="form-input" 
                v-model="newCategory.name" 
                placeholder="e.g. Farmers" 
                required 
              />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Hindi Name</label>
                <input 
                  type="text" 
                  class="form-input" 
                  v-model="newCategory.name_hi" 
                  placeholder="किसान" 
                />
              </div>
              <div class="form-group">
                <label class="form-label">Marathi Name</label>
                <input 
                  type="text" 
                  class="form-input" 
                  v-model="newCategory.name_mr" 
                  placeholder="शेतकरी" 
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Icon (Emoji)</label>
              <input 
                type="text" 
                class="form-input" 
                v-model="newCategory.icon" 
                placeholder="🌾" 
                required 
              />
            </div>

            <div class="form-group">
              <label class="form-label">Description</label>
              <input 
                type="text" 
                class="form-input" 
                v-model="newCategory.description" 
                placeholder="Schemes for farmers and agriculture" 
                required 
              />
            </div>

            <button type="submit" class="submit-btn">
              <i class="ti ti-plus"></i>
              <span>Add Category</span>
            </button>
          </form>
        </div>
      </div>

      <!-- RIGHT CARD — "Existing Categories" -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">Existing Categories</div>
        </div>
        <div class="card-body">
          <div class="categories-list">
            <div class="category-item" v-for="c in categories" :key="c.id">
              <div class="cat-left">
                <span class="cat-emoji">{{ c.icon }}</span>
                <div class="cat-texts">
                  <div class="cat-name">{{ c.name }}</div>
                  <div class="cat-subtitle">{{ c.name_hi }} · {{ c.name_mr }}</div>
                </div>
              </div>
              <div class="cat-actions">
                <button class="action-btn" title="Edit" @click="$emit('edit-category', c)">
                  <i class="ti ti-edit"></i>
                </button>
                <button class="action-btn danger-hover" title="Delete" @click="emit('delete-category', c.id)">
                  <i class="ti ti-trash"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

  </div>
</template>

<style scoped>
.categories-tab {
  width: 100%;
}

.two-col-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  align-items: start;
}

.card {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  padding: 14px 16px;
  border-bottom: 0.5px solid var(--border);
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
}

.card-body {
  padding: 16px;
}

/* Forms */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.form-label {
  font-size: 13px;
  color: var(--text);
}

.form-input {
  padding: 8px 10px;
  border: 0.5px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  background-color: var(--bg);
  color: var(--text);
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
  width: 100%;
}

.form-input:focus {
  border-color: var(--primary);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.submit-btn {
  background-color: var(--primary);
  color: var(--clr-text-light);
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
}

/* Categories List */
.categories-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-item {
  background-color: var(--bg2);
  border-radius: 6px; /* var(--radius) */
  padding: 8px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
}

.cat-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cat-emoji {
  font-size: 20px;
  line-height: 1;
}

.cat-texts {
  display: flex;
  flex-direction: column;
}

.cat-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
  line-height: 1.3;
}

.cat-subtitle {
  font-size: 11px;
  color: var(--text2);
  margin-top: 1px;
  line-height: 1.2;
}

.cat-actions {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: 0.5px solid var(--border);
  background-color: var(--bg);
  color: var(--text2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 14px;
  padding: 0;
  box-sizing: border-box;
}

.action-btn i {
  font-size: 14px !important;
}

.action-btn:hover {
  background-color: var(--bg2);
  color: var(--text);
}

.action-btn.danger-hover:hover {
  background-color: var(--danger-bg);
  color: var(--danger);
  border-color: var(--danger);
}
</style>
