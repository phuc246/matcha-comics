<template>
  <div class="comic-reader" :class="{ 'menu-hidden': !showMenu }">
    <ContentProtection />
    <!-- Loader -->
    <div v-if="loading" class="reader-loader">
      <div class="spinner"></div>
      <p>Đang tải chương truyện...</p>
    </div>
    
    <template v-else-if="comic">
      <!-- Top Bar (Floating) -->
      <header class="reader-header" :class="{ 'is-hidden': !showMenu }">
        <div class="header-left">
          <button class="btn-icon" @click="router.back()">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
          </button>
          <div class="reader-title">
            <h1>{{ comic.title }}</h1>
            
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
                  :to="`/truyen-tranh/${comic.slug}/chapter-${c.number}`"
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
          <NuxtLink to="/" class="btn-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
          </NuxtLink>
        </div>
      </header>

      <!-- Main Content -->
      <main class="reader-main" @click="handleScreenClick">
        <div class="pages-list">
          <!-- Spacer top -->
          <div class="reader-spacer"></div>
          
          <div v-for="(img, index) in chapterImages" :key="index" class="page-item page-item-protected">
            <img 
              :src="img" 
              :alt="`Trang ${index + 1}`" 
              class="page-img"
              loading="lazy"
            />
            <!-- Protection Shield -->
            <div class="protection-shield"></div>
          </div>
          
          <!-- Spacer bottom -->
          <div class="reader-spacer"></div>
        </div>

        <!-- End of Chapter Navigation -->
        <div class="reader-footer-nav container">
          <div class="nav-controls">
            <NuxtLink
              v-if="prevChapter"
              :to="`/truyen-tranh/${comic.slug}/chapter-${prevChapter}`"
              class="nav-btn prev"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 18l-6-6 6-6"/></svg>
              <span>Chương Trước</span>
            </NuxtLink>
            <button v-else class="nav-btn disabled">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 18l-6-6 6-6"/></svg>
              <span>Chương Trước</span>
            </button>

            <div class="chapter-picker">
              {{ currentChapterNum }} / {{ comic.chapters?.length || '?' }}
            </div>

            <NuxtLink
              v-if="nextChapter"
              :to="`/truyen-tranh/${comic.slug}/chapter-${nextChapter}`"
              class="nav-btn next active"
            >
              <span>Chương Sau</span>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 18l6-6-6-6"/></svg>
            </NuxtLink>
            <NuxtLink v-else :to="`/truyen-tranh/${comic.slug}`" class="nav-btn next active">
              <span>Hết Truyện</span>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
            </NuxtLink>
          </div>
        </div>
      </main>

      <!-- Bottom Floating Menu -->
      <nav class="reader-bottom-bar" :class="{ 'is-hidden': !showMenu }">
        <div class="bottom-nav-inner">
           <!-- Server Switcher -->
           <div class="server-switcher" v-if="availableServers.length > 1">
             <span class="server-label">Server:</span>
             <div class="server-btns">
               <button 
                v-for="(srv, idx) in availableServers" 
                :key="idx"
                class="srv-btn"
                :class="{ active: currentServerIndex === idx }"
                @click.stop="switchServer(idx)"
               >
                 {{ srv.name.replace('Server ', 'S') }}
               </button>
             </div>
           </div>

           <div class="bottom-controls">
             <NuxtLink v-if="prevChapter" :to="`/truyen-tranh/${comic.slug}/chapter-${prevChapter}`" class="bottom-item">
               <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 18l-6-6 6-6"/></svg>
             </NuxtLink>
             <div class="bottom-item current" @click.stop="showChapters = !showChapters">
               <span>Ch.{{ currentChapterNum }}</span>
             </div>
             <NuxtLink v-if="nextChapter" :to="`/truyen-tranh/${comic.slug}/chapter-${nextChapter}`" class="bottom-item">
               <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 18l6-6-6-6"/></svg>
             </NuxtLink>
           </div>
        </div>
      </nav>
    </template>
    
    <div v-else class="reader-error">
      <div class="error-box">
        <span class="error-icon">📭</span>
        <h2>Không tìm thấy chương</h2>
        <NuxtLink to="/" class="btn-back">Quay lại</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

definePageMeta({ layout: false })

const route = useRoute()
const router = useRouter()
const { get } = useApi()

const loading = ref(true)
const comic = ref<any>(null)
const currentChapterNum = ref(1)
const chapterImages = ref<string[]>([])
const showMenu = ref(true)
const showChapters = ref(false)

const currentServerIndex = ref(0)
const availableServers = ref<any[]>([])

const sortedChapters = computed(() => {
  if (!comic.value?.chapters) return []
  return [...comic.value.chapters].sort((a, b) => a.number - b.number)
})

const prevChapter = computed(() => {
  if (!comic.value?.chapters) return null
  const idx = sortedChapters.value.findIndex(c => c.number === currentChapterNum.value)
  return idx > 0 ? sortedChapters.value[idx - 1].number : null
})

const nextChapter = computed(() => {
  if (!comic.value?.chapters) return null
  const idx = sortedChapters.value.findIndex(c => c.number === currentChapterNum.value)
  return (idx >= 0 && idx < sortedChapters.value.length - 1) ? sortedChapters.value[idx + 1].number : null
})

const handleScreenClick = (e: MouseEvent) => {
  if (showChapters.value) {
    showChapters.value = false
    return
  }
  const y = e.clientY
  const height = window.innerHeight
  if (y > height * 0.2 && y < height * 0.8) {
    showMenu.value = !showMenu.value
  }
}

const switchServer = (idx: number) => {
  currentServerIndex.value = idx
  const server = availableServers.value[idx]
  if (server && server.content) {
    try {
      chapterImages.value = JSON.parse(server.content)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch (e) { console.error(e) }
  }
}

const loadData = async () => {
  const slug = route.params.slug as string
  const chapterParam = route.params.chapter as string
  const match = chapterParam.match(/chapter-(\d+)/)
  if (match) currentChapterNum.value = parseInt(match[1])

  loading.value = true
  showChapters.value = false
  try {
    const story = await get<any>(`/stories/${slug}`)
    if (story) comic.value = story

    const chapterData = await get<any>(`/stories/${slug}/chapters/${currentChapterNum.value}`)
    if (chapterData && chapterData.servers) {
      availableServers.value = chapterData.servers.filter((s: any) => s.content && s.content !== '[]')
      if (availableServers.value.length > 0) {
        switchServer(0)
      }
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
    setTimeout(() => { showMenu.value = false }, 3000)
  }
}

onMounted(() => {
  loadData()
  window.addEventListener('click', () => { showChapters.value = false })
})

watch(() => route.params.chapter, () => {
  loadData()
  window.scrollTo(0, 0)
})
</script>

<style scoped>
.comic-reader {
  min-height: 100vh;
  background: #08090b;
  color: #fff;
  font-family: 'Inter', sans-serif;
}

.reader-loader { height: 100vh; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 20px; }
.spinner { width: 50px; height: 50px; border: 4px solid rgba(156,167,100,0.1); border-top-color: #9CA764; border-radius: 50%; animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Header */
.reader-header { position: fixed; top: 0; left: 0; right: 0; height: 64px; background: rgba(13,15,20,0.95); backdrop-filter: blur(20px); display: flex; align-items: center; justify-content: space-between; padding: 0 16px; z-index: 1000; border-bottom: 1px solid rgba(255,255,255,0.08); transition: transform 0.3s ease; }
.reader-header.is-hidden { transform: translateY(-100%); }

.header-left { display: flex; align-items: center; gap: 16px; flex: 1; }
.reader-title { display: flex; flex-direction: column; }
.reader-title h1 { font-size: 0.9rem; font-weight: 700; margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 150px; color: #fff; }

.dropdown-wrap { position: relative; }
.chapter-selector { background: none; border: none; color: #9CA764; font-size: 0.8rem; font-weight: 700; display: flex; align-items: center; gap: 4px; padding: 4px 0; cursor: pointer; }
.chapter-dropdown { position: absolute; top: 100%; left: 0; background: #14171f; border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; width: 220px; max-height: 300px; overflow-y: auto; box-shadow: 0 10px 30px rgba(0,0,0,0.5); z-index: 1001; margin-top: 8px; }
.dropdown-item { display: block; padding: 12px 16px; color: #A8A8B3; text-decoration: none; font-size: 0.85rem; border-bottom: 1px solid rgba(255,255,255,0.05); }
.dropdown-item:hover { background: rgba(255,255,255,0.05); color: #fff; }
.dropdown-item.active { color: #9CA764; background: rgba(156,167,100,0.1); }

/* Main */
.reader-main { display: flex; flex-direction: column; align-items: center; }
.pages-list { width: 100%; max-width: 900px; display: flex; flex-direction: column; }
.reader-spacer { height: 80px; width: 100%; }
.page-item { width: 100%; line-height: 0; background: #000; }
.page-img { width: 100%; height: auto; display: block; }

/* Navigation */
.reader-footer-nav { width: 100%; max-width: 600px; padding: 40px 20px 120px; }
.nav-controls { display: flex; align-items: center; justify-content: space-between; gap: 20px; }
.nav-btn { flex: 1; display: flex; align-items: center; justify-content: center; gap: 8px; padding: 16px; background: #1a1d26; border: 1px solid rgba(255,255,255,0.1); border-radius: 14px; color: #fff; font-weight: 700; text-decoration: none; }
.nav-btn.active { background: #9CA764; color: #000; border: none; }
.nav-btn.disabled { opacity: 0.3; cursor: not-allowed; }

/* Bottom Bar */
.reader-bottom-bar { position: fixed; bottom: 0; left: 0; right: 0; background: rgba(13,15,20,0.95); backdrop-filter: blur(20px); z-index: 1000; border-top: 1px solid rgba(255,255,255,0.08); padding: 8px 16px 20px; transition: transform 0.3s ease; }
.reader-bottom-bar.is-hidden { transform: translateY(100%); }

.bottom-nav-inner { max-width: 600px; margin: 0 auto; display: flex; flex-direction: column; gap: 12px; }

.server-switcher { display: flex; align-items: center; justify-content: center; gap: 12px; padding-bottom: 8px; border-bottom: 1px solid rgba(255,255,255,0.05); }
.server-label { font-size: 0.75rem; color: #5C5C6B; font-weight: 600; }
.server-btns { display: flex; gap: 8px; }
.srv-btn { background: #1a1d26; border: 1px solid rgba(255,255,255,0.1); color: #A8A8B3; padding: 6px 12px; border-radius: 6px; font-size: 0.75rem; font-weight: 700; cursor: pointer; }
.srv-btn.active { background: #9CA764; color: #000; border: none; }

.bottom-controls { display: flex; align-items: center; justify-content: space-around; }
.bottom-item { color: #A8A8B3; cursor: pointer; display: flex; align-items: center; justify-content: center; width: 44px; height: 44px; }
.bottom-item.current { color: #9CA764; font-weight: 800; font-size: 0.9rem; width: auto; padding: 0 20px; }

@media (max-width: 600px) {
  .nav-btn span { display: none; }
  .reader-spacer { height: 64px; }
}

/* Custom scrollbar for dropdown */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 10px; }

/* Image Protection Overlay */
.page-item-protected {
  position: relative;
  width: 100%;
}
.protection-shield {
  position: absolute;
  inset: 0;
  background: transparent;
  z-index: 10;
  user-select: none;
  -webkit-touch-callout: none;
}
</style>
