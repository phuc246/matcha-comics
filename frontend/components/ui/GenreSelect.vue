<template>
  <div class="genre-select-container" ref="selectContainer">
    <label class="select-label">
      {{ label || 'Thể loại' }}
      <span v-if="required" class="required">*</span>
    </label>
    
    <div class="custom-select" :class="{ 'is-open': isOpen }">
      <div class="select-header" @click="isOpen = !isOpen">
        <div class="selected-genres">
          <span v-if="modelValue.length === 0">Chọn thể loại...</span>
          <span v-else v-for="id in modelValue" :key="id" class="genre-tag">
            {{ availableGenres.find(g => g.id === id)?.name || id }}
            <button type="button" class="remove-genre" @click.stop="removeGenre(id)">×</button>
          </span>
        </div>
        <span class="chevron">▼</span>
      </div>
      
      <div class="select-dropdown" v-show="isOpen">
        <div class="genres-checkbox-list">
          <label v-for="genre in availableGenres" :key="genre.id" class="checkbox-item">
            <input 
              type="checkbox" 
              :value="genre.id" 
              :checked="modelValue.includes(genre.id)"
              @change="toggleGenre(genre.id)"
            />
            <span>{{ genre.name }}</span>
          </label>
        </div>
      </div>
    </div>
    <small class="help-text">Chọn tối đa {{ maxLimit }} thể loại</small>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: number[]
  label?: string
  required?: boolean
  maxLimit?: number
}>(), {
  required: false,
  maxLimit: 7,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number[]): void
}>()

const isOpen = ref(false)
const selectContainer = ref<HTMLElement | null>(null)
const availableGenres = ref<any[]>([])
const { get } = useApi()

const toggleGenre = (genreId: number) => {
  const selected = [...props.modelValue]
  const index = selected.indexOf(genreId)
  
  if (index > -1) {
    selected.splice(index, 1)
  } else {
    if (selected.length >= props.maxLimit) {
      alert(`Chỉ được chọn tối đa ${props.maxLimit} thể loại!`)
      return
    }
    selected.push(genreId)
  }
  emit('update:modelValue', selected)
}

const removeGenre = (genreId: number) => {
  const selected = props.modelValue.filter(id => id !== genreId)
  emit('update:modelValue', selected)
}

const handleClickOutside = (e: MouseEvent) => {
  if (selectContainer.value && !selectContainer.value.contains(e.target as Node)) {
    isOpen.value = false
  }
}

onMounted(async () => {
  document.addEventListener('click', handleClickOutside)
  try {
    const data = await get<any[]>('/genres')
    availableGenres.value = data || []
  } catch (err) {
    console.error('Error fetching genres in GenreSelect component:', err)
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.genre-select-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.select-label {
  font-size: 0.85rem;
  color: #A8A8B3;
  font-weight: 600;
  display: flex;
  justify-content: space-between;
}

.required {
  color: #ef4444;
}

/* Custom Select Dropdown Box */
.custom-select {
  position: relative;
  width: 100%;
}

.select-header {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 10px 12px;
  color: #fff;
  font-size: 0.9rem;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: border-color 0.2s;
  min-height: 46px;
}

.select-header:hover {
  border-color: #9CA764;
}

.custom-select.is-open .select-header {
  border-color: #9CA764;
}

.selected-genres {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

.genre-tag {
  background: rgba(156, 167, 100, 0.2);
  color: #9CA764;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
  border: 1px solid rgba(156, 167, 100, 0.3);
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.remove-genre {
  background: none;
  border: none;
  color: #9CA764;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 0;
  line-height: 1;
}

.remove-genre:hover {
  color: #ef4444;
}

.select-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  width: 100%;
  background: #1F2330;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 12px;
  z-index: 100;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
  max-height: 250px;
  overflow-y: auto;
}

.genres-checkbox-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 2px 0;
}

.checkbox-item:hover {
  color: #9CA764;
}

.checkbox-item input {
  accent-color: #9CA764;
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.chevron {
  font-size: 0.7rem;
  color: #A8A8B3;
  margin-left: 8px;
}

.help-text {
  font-size: 0.75rem;
  color: #9CA764;
  margin-top: 4px;
}
</style>
