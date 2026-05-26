<template>
  <div class="novel-reader" :class="`theme-${readerSettings.theme}`">
    <ContentProtection />
    <!-- Floating Header -->
    <header v-if="novel" class="reader-floating-header" :class="{ 'header-hidden': !showHeader }">
      <div class="header-left">
        <button class="back-btn" @click="goBack" aria-label="Quay lại">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m15 18-6-6 6-6"/>
          </svg>
        </button>
        <div class="reader-title-info">
          <h1 class="story-title">{{ novel.title }}</h1>
          
          <!-- Chapter Selector Dropdown -->
          <div class="dropdown-wrap">
            <button class="chapter-selector" @click.stop="showChapters = !showChapters">
              Chương {{ currentChapterNum }}
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 9l6 6 6-6"/></svg>
            </button>
            <div v-if="showChapters" class="chapter-dropdown custom-scrollbar">
              <NuxtLink
                v-for="c in sortedChapters"
                :key="c.id"
                :to="`/truyen-chu/${novel.slug}/chapter-${c.number}`"
                class="dropdown-item"
                :class="{ active: c.number === currentChapterNum }"
                @click="showChapters = false"
              >
                Chương {{ c.number }} {{ c.title ? ': ' + c.title : '' }}
              </NuxtLink>
            </div>
          </div>
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

    <div v-if="loading" class="reader-loader">
      <div class="loader-spinner"></div>
    </div>
    
    <div v-else-if="novel" class="reader-container" :style="{ fontSize: `${readerSettings.fontSize}px`, fontFamily: readerSettings.fontFamily }">
      <div class="reader-header-inline">
        <h1>{{ novel.title }}</h1>
        <h2>Chương {{ currentChapterNum }}{{ chapterDetail?.title ? `: ${chapterDetail.title}` : '' }}</h2>
      </div>

      <!-- Settings toolbar -->
      <div class="reader-settings-bar">
        <div class="setting-group">
          <button @click="changeFontSize(-2)" class="setting-btn">A-</button>
          <span class="setting-value">{{ readerSettings.fontSize }}px</span>
          <button @click="changeFontSize(2)" class="setting-btn">A+</button>
        </div>
        <div class="setting-group theme-toggles">
          <button @click="setTheme('dark')" class="theme-btn theme-dark" :class="{ active: readerSettings.theme === 'dark' }" title="Tối"></button>
          <button @click="setTheme('light')" class="theme-btn theme-light" :class="{ active: readerSettings.theme === 'light' }" title="Sáng"></button>
          <button @click="setTheme('sepia')" class="theme-btn theme-sepia" :class="{ active: readerSettings.theme === 'sepia' }" title="Vàng nhạt"></button>
        </div>
      </div>

      <!-- Content -->
      <div class="chapter-content" v-html="chapterContent"></div>

      <!-- Chapter Navigation -->
      <div class="chapter-navigation">
        <NuxtLink
          v-if="prevChapter"
          :to="`/truyen-chu/${novel.slug}/chapter-${prevChapter}`"
          class="nav-btn prev-btn"
        >
          Chương trước
        </NuxtLink>
        <button v-else class="nav-btn disabled">Chương trước</button>

        <NuxtLink
          v-if="nextChapter"
          :to="`/truyen-chu/${novel.slug}/chapter-${nextChapter}`"
          class="nav-btn next-btn primary"
        >
          Chương sau
        </NuxtLink>
        <button v-else class="nav-btn disabled">Chương sau</button>
      </div>
    </div>
    
    <div v-else class="error-state">
      <p>Không thể tải chương này.</p>
      <NuxtLink to="/truyen-chu" class="back-link">Về danh sách truyện</NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

definePageMeta({
  layout: false
})

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const novel = ref<any>(null)
const currentChapterNum = ref(1)

const chapterDetail = ref<any>(null)
const chapterContent = computed(() => {
  if (!chapterDetail.value || !chapterDetail.value.servers) return ''
  // For novel, content is stored in the first server's content directly as HTML string
  return chapterDetail.value.servers[0]?.content || '<p>Chương này chưa có nội dung.</p>'
})

const readerSettings = reactive({
  fontSize: 20,
  fontFamily: "'Inter', sans-serif",
  theme: 'dark' // dark | light | sepia
})

const changeFontSize = (delta: number) => {
  const newVal = readerSettings.fontSize + delta
  if (newVal >= 14 && newVal <= 36) {
    readerSettings.fontSize = newVal
  }
}

const setTheme = (theme: string) => {
  readerSettings.theme = theme
}

const sortedChapters = computed(() => {
  if (!novel.value?.chapters) return []
  return [...novel.value.chapters].sort((a, b) => a.number - b.number)
})

const prevChapter = computed(() => {
  if (!novel.value?.chapters) return null
  const idx = sortedChapters.value.findIndex(c => c.number === currentChapterNum.value)
  return idx > 0 ? sortedChapters.value[idx - 1].number : null
})

const nextChapter = computed(() => {
  if (!novel.value?.chapters) return null
  const idx = sortedChapters.value.findIndex(c => c.number === currentChapterNum.value)
  return (idx >= 0 && idx < sortedChapters.value.length - 1) ? sortedChapters.value[idx + 1].number : null
})

const { get } = useApi()

const loadData = async () => {
  const slug = route.params.slug as string
  const chapterParam = route.params.chapter as string
  
  const match = chapterParam.match(/chapter-(\d+(\.\d+)?)/)
  if (match) {
    currentChapterNum.value = parseFloat(match[1])
  }

  loading.value = true
  try {
    const [storyRes, chapterRes] = await Promise.all([
      get<any>(`/stories/${slug}`),
      get<any>(`/stories/${slug}/chapters/${currentChapterNum.value}`)
    ])
    
    if (storyRes) {
      novel.value = storyRes
    }
    
    if (chapterRes) {
      chapterDetail.value = chapterRes
    }
  } catch (err) {
    console.error('Error loading novel chapter:', err)
  } finally {
    loading.value = false
  }
}

// Scroll and dropdown state
const showHeader = ref(true)
const showChapters = ref(false)
let lastScroll = 0

const handleScroll = () => {
  const currentScroll = window.scrollY
  if (currentScroll <= 0) {
    showHeader.value = true
    return
  }
  if (currentScroll > lastScroll && currentScroll > 50) {
    showHeader.value = false // scrolling down
    showChapters.value = false
  } else if (currentScroll < lastScroll) {
    showHeader.value = true // scrolling up
  }
  lastScroll = currentScroll
}

const goBack = () => {
  if (window.history.length > 2) {
    router.back()
  } else {
    router.push(`/truyen-chu/${novel.value?.slug || ''}`)
  }
}

const closeDropdown = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.dropdown-wrap')) {
    showChapters.value = false
  }
}

onMounted(() => {
  loadData()
  window.addEventListener('scroll', handleScroll, { passive: true })
  window.addEventListener('click', closeDropdown)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  window.removeEventListener('click', closeDropdown)
})

watch(() => route.params.chapter, () => {
  loadData()
  window.scrollTo(0, 0)
})

useHead(() => ({
  title: novel.value ? `${novel.value.title} - Chương ${currentChapterNum.value}` : 'Đang tải...',
}))
</script>

<style scoped>
.novel-reader {
  width: 100%;
  min-height: 100vh;
  transition: background-color 0.3s, color 0.3s;
}

/* Floating Header */
.reader-floating-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 100;
  transition: transform 0.3s ease, background-color 0.3s, border-color 0.3s, color 0.3s;
}

.header-hidden {
  transform: translateY(-100%);
}

/* Match Header styling with active theme */
.theme-dark .reader-floating-header {
  background: rgba(13, 15, 20, 0.95);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  color: #a8a8b3;
}

.theme-light .reader-floating-header {
  background: rgba(248, 249, 250, 0.95);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  color: #333333;
}

.theme-sepia .reader-floating-header {
  background: rgba(241, 231, 208, 0.95);
  border-bottom: 1px solid rgba(92, 75, 55, 0.15);
  color: #5c4b37;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  min-width: 0;
}

.back-btn {
  background: none;
  border: none;
  color: inherit;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 6px;
  border-radius: 50%;
  transition: background 0.2s;
}

.theme-dark .back-btn:hover { background: rgba(255,255,255,0.08); }
.theme-light .back-btn:hover { background: rgba(0,0,0,0.05); }
.theme-sepia .back-btn:hover { background: rgba(0,0,0,0.05); }

.reader-title-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.story-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: inherit;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin: 0;
}

/* Chapter Selector Dropdown styling */
.dropdown-wrap {
  position: relative;
}

.chapter-selector {
  background: none;
  border: none;
  color: #9CA764;
  font-size: 0.8rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 0;
  cursor: pointer;
}

.chapter-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  border-radius: 8px;
  width: 220px;
  max-height: 250px;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgba(0,0,0,0.3);
  z-index: 101;
  margin-top: 8px;
  border: 1px solid;
  transition: background-color 0.3s, color 0.3s, border-color 0.3s;
}

.theme-dark .chapter-dropdown {
  background: #14171f;
  border-color: rgba(255, 255, 255, 0.1);
}
.theme-light .chapter-dropdown {
  background: #ffffff;
  border-color: rgba(0, 0, 0, 0.1);
}
.theme-sepia .chapter-dropdown {
  background: #ebdcb9;
  border-color: rgba(92, 75, 55, 0.15);
}

.dropdown-item {
  display: block;
  padding: 10px 14px;
  text-decoration: none;
  font-size: 0.85rem;
  border-bottom: 1px solid;
  transition: all 0.2s;
  color: inherit;
}

.theme-dark .dropdown-item { border-bottom-color: rgba(255,255,255,0.05); }
.theme-light .dropdown-item { border-bottom-color: rgba(0,0,0,0.05); }
.theme-sepia .dropdown-item { border-bottom-color: rgba(92, 75, 55, 0.08); }

.dropdown-item:hover {
  background: rgba(156, 167, 100, 0.15);
}

.dropdown-item.active {
  color: #9CA764 !important;
  font-weight: 700;
  background: rgba(156, 167, 100, 0.1);
}

.header-right {
  display: flex;
  align-items: center;
}

.home-btn {
  color: inherit;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  border-radius: 50%;
  transition: all 0.2s;
}

.theme-dark .home-btn:hover { background: rgba(255,255,255,0.08); }
.theme-light .home-btn:hover { background: rgba(0,0,0,0.05); }
.theme-sepia .home-btn:hover { background: rgba(0,0,0,0.05); }

/* Themes */
.theme-dark {
  background-color: #0d0f14;
  color: #a8a8b3;
}
.theme-light {
  background-color: #f8f9fa;
  color: #333333;
}
.theme-sepia {
  background-color: #f1e7d0;
  color: #5c4b37;
}

.reader-loader { height: 80vh; display: flex; align-items: center; justify-content: center; }
.loader-spinner { width: 40px; height: 40px; border: 3px solid rgba(156,167,100,0.3); border-top-color: #9CA764; border-radius: 50%; animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.reader-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 80px 20px 40px !important;
}

.reader-header-inline {
  text-align: center;
  margin-bottom: 40px;
}
.reader-header-inline h1 {
  font-size: 1.5em;
  font-weight: 800;
  margin-bottom: 8px;
  font-family: 'Montserrat', sans-serif;
}

.reader-header-inline h1,
.reader-header-inline h2 {
  color: inherit !important;
}

.reader-header-inline h2 {
  font-size: 1.1em;
  opacity: 0.8;
  font-weight: 500;
}

.reader-settings-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-radius: 12px;
  margin-bottom: 40px;
  font-size: 1rem;
}

.theme-dark .reader-settings-bar { background: rgba(255,255,255,0.05); }
.theme-light .reader-settings-bar { background: rgba(0,0,0,0.05); }
.theme-sepia .reader-settings-bar { background: rgba(0,0,0,0.05); }

.setting-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.setting-btn {
  background: none;
  border: 1px solid currentColor;
  opacity: 0.5;
  border-radius: 4px;
  padding: 4px 12px;
  cursor: pointer;
  color: inherit;
  transition: all 0.2s;
}
.setting-btn:hover { opacity: 1; }

.theme-toggles { gap: 8px; }
.theme-btn {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: transform 0.2s;
}
.theme-btn.active { transform: scale(1.2); border-color: #9CA764; }
.theme-dark.theme-btn { background: #0d0f14; border: 1px solid #333; }
.theme-light.theme-btn { background: #f8f9fa; border: 1px solid #ccc; }
.theme-sepia.theme-btn { background: #f1e7d0; border: 1px solid #d4c5b0; }

.chapter-content {
  line-height: 1.8;
  margin-bottom: 60px;
}

.chapter-content :deep(*) {
  color: inherit !important;
  background-color: transparent !important;
}

.chapter-content :deep(p) {
  margin-bottom: 1.5em;
}

.chapter-navigation {
  display: flex;
  gap: 16px;
  padding: 40px 0 80px;
  justify-content: center;
}

.nav-btn {
  padding: 14px 28px;
  border-radius: 8px;
  font-weight: 600;
  text-align: center;
  transition: all 0.2s;
  flex: 1;
  max-width: 200px;
  font-size: 1rem;
  border: 1px solid currentColor;
  background: transparent;
  color: inherit;
  opacity: 0.8;
}

.nav-btn:hover:not(.disabled) {
  opacity: 1;
}

.nav-btn.primary {
  background: #9CA764;
  color: #000;
  border: none;
  opacity: 1;
}

.nav-btn.primary:hover:not(.disabled) {
  background: #B8C878;
}

.nav-btn.disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

/* Custom scrollbar for dropdown */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 10px; }

@media (max-width: 600px) {
  .reader-container {
    padding: 68px 12px 24px !important;
  }
  .reader-header-inline {
    margin-bottom: 24px;
  }
  .reader-header-inline h1 {
    font-size: 1.25em;
  }
  .reader-header-inline h2 {
    font-size: 0.95em;
  }
  .reader-settings-bar {
    padding: 8px 12px;
    margin-bottom: 24px;
    font-size: 0.85rem;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: center;
  }
  .setting-group {
    gap: 8px;
  }
  .setting-btn {
    padding: 2px 8px;
  }
  .chapter-content {
    line-height: 1.6;
    margin-bottom: 40px;
  }
  .chapter-navigation {
    padding: 20px 0 60px;
    gap: 10px;
  }
  .nav-btn {
    padding: 10px 16px;
    font-size: 0.85rem;
  }
}
</style>
