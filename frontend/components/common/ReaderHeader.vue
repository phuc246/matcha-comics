<template>
  <header class="reader-header" :class="{ 'header-hidden': !visible }">
    <div class="header-left">
      <button class="back-btn" @click="goBack" aria-label="Quay lại">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m15 18-6-6 6-6"/>
        </svg>
      </button>
      <div class="reader-title-info">
        <h1 class="story-title">{{ title || 'Đang tải...' }}</h1>
        <p class="chapter-title" v-if="chapter">{{ chapter }}</p>
      </div>
    </div>
    
    <div class="header-right">
      <NuxtLink to="/" class="home-btn" aria-label="Trang chủ">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
        </svg>
      </NuxtLink>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

defineProps<{
  title?: string
  chapter?: string
}>()

const router = useRouter()
const visible = ref(true)
let lastScroll = 0

const handleScroll = () => {
  const currentScroll = window.scrollY
  if (currentScroll <= 0) {
    visible.value = true
    return
  }
  if (currentScroll > lastScroll && currentScroll > 50) {
    visible.value = false // scrolling down
  } else if (currentScroll < lastScroll) {
    visible.value = true // scrolling up
  }
  lastScroll = currentScroll
}

const goBack = () => {
  if (window.history.length > 2) {
    router.back()
  } else {
    router.push('/')
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.reader-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: rgba(13, 15, 20, 0.95);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 100;
  transition: transform 0.3s ease;
}

.header-hidden {
  transform: translateY(-100%);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  min-width: 0; /* for text truncation */
}

.back-btn {
  background: none;
  border: none;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 4px;
  border-radius: 50%;
  transition: background 0.2s;
}

.back-btn:hover {
  background: rgba(255,255,255,0.1);
}

.reader-title-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.story-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin: 0;
}

.chapter-title {
  font-size: 0.75rem;
  color: var(--accent-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.home-btn {
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  border-radius: 50%;
  transition: all 0.2s;
}

.home-btn:hover {
  color: var(--accent-primary);
  background: rgba(156,167,100,0.1);
}
</style>
