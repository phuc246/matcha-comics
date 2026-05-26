<template>
  <div class="catalog-page">
    <!-- Page Hero -->
    <div class="catalog-hero">
      <div class="container">
        <div class="breadcrumb">
          <NuxtLink to="/">Trang chủ</NuxtLink>
          <span class="separator">/</span>
          <NuxtLink to="/the-loai">Thể loại</NuxtLink>
          <span class="separator">/</span>
          <span class="current">{{ genreName }}</span>
        </div>
        <h1 class="catalog-title">{{ genreIcon }} Thể loại: {{ genreName }}</h1>
        <p class="catalog-desc">Khám phá danh sách truyện thuộc thể loại {{ genreName }} cập nhật mới nhất, hấp dẫn nhất.</p>
      </div>
    </div>

    <div class="container catalog-body">
      <!-- Grid -->
      <div class="catalog-grid">
        <template v-if="loading">
          <ComicCardSkeleton v-for="i in 12" :key="i" />
        </template>
        <template v-else-if="comics.length > 0">
          <ComicCard3D
            v-for="(comic, i) in comics"
            :key="comic.id"
            :comic="comic"
            :show-info="true"
            class="animate-fade-up"
            :style="{ animationDelay: `${(i % 12) * 0.05}s` }"
          />
        </template>
        <div v-else class="empty-state">
          <p>Không tìm thấy truyện nào thuộc thể loại này.</p>
        </div>
      </div>

      <!-- Pagination -->
      <div class="pagination-wrap" v-if="comics.length > 0">
        <button class="page-btn" disabled>«</button>
        <button class="page-btn active">1</button>
        <button class="page-btn">2</button>
        <button class="page-btn">3</button>
        <span class="page-dots">...</span>
        <button class="page-btn">»</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { mockComics, mockNovels } from '~/data/mock'

const route = useRoute()
const loading = ref(true)
const comics = ref<any[]>([])
const genresList = ref<any[]>([])
const { get } = useApi()

const genreSlug = computed(() => route.params.slug as string)

const emojiMap: Record<string, string> = {
  'action': '⚔️',
  'hanh-dong': '⚔️',
  'romance': '💕',
  'lang-man': '💕',
  'ngon-tinh': '💞',
  'fantasy': '🧙',
  'ky-ao': '🧙',
  'horror': '👻',
  'kinh-di': '👻',
  'tien-hiep': '⛅',
  'kiem-hiep': '🗡️',
  'comedy': '😂',
  'hai-huoc': '😂',
  'drama': '🎭',
  'kich-tinh': '🎭',
  'adventure': '🧭',
  'phieu-luu': '🧭',
}

const getGenreIcon = (slug: string) => {
  return emojiMap[slug.toLowerCase()] || '🏷️'
}

const genreInfo = computed(() => genresList.value.find(g => g.slug.toLowerCase() === genreSlug.value.toLowerCase()))
const genreName = computed(() => genreInfo.value?.name || genreSlug.value)
const genreIcon = computed(() => getGenreIcon(genreSlug.value))

const fetchGenres = async () => {
  try {
    const data = await get<any[]>('/genres')
    if (data) genresList.value = data
  } catch (err) {
    console.error('Error fetching genres list in slug:', err)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    // 1. Fetch real stories from database
    const allStories = await get<any[]>('/stories')
    if (allStories && allStories.length > 0) {
      comics.value = allStories.filter(c => {
        return c.genres?.some((g: any) => g.slug.toLowerCase() === genreSlug.value.toLowerCase())
      })
    }
    
    // 2. If no real stories in DB or none match, fallback to mock data to keep the screen loaded with content
    if (comics.value.length === 0) {
      const allMock = [...mockComics, ...mockNovels]
      comics.value = allMock.filter(c => {
        return c.genres?.map((g: string) => g.toLowerCase().replace(/ /g, '-')).includes(genreSlug.value.toLowerCase().replace(/ /g, '-'))
      })
      if (comics.value.length === 0) {
        comics.value = allMock.slice(0, 8)
      }
    }
  } catch (error) {
    console.error('Error fetching stories for genre:', error)
    // Fallback on error
    const allMock = [...mockComics, ...mockNovels]
    comics.value = allMock.slice(0, 8)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await fetchGenres()
  await fetchData()
})

watch(() => route.params.slug, async () => {
  await fetchData()
})

useHead(() => ({
  title: `Thể loại ${genreName.value} - MatchaComic`,
}))
</script>

<style scoped>
.catalog-page { padding-bottom: 60px; }

.catalog-hero { background: var(--bg-secondary); padding: 40px 0 50px; border-bottom: 1px solid var(--border-card); position: relative; overflow: hidden; }
.catalog-hero::before { content: ''; position: absolute; inset: 0; background: radial-gradient(circle at 80% 0%, rgba(156,167,100,0.05) 0%, transparent 50%); pointer-events: none; }

.breadcrumb { display: flex; gap: 8px; font-size: 0.85rem; color: var(--text-muted); margin-bottom: 24px; }
.breadcrumb a { transition: var(--transition-fast); }
.breadcrumb a:hover { color: var(--accent-primary); }
.breadcrumb .current { color: var(--text-primary); font-weight: 500; }

.catalog-title { font-family: 'Montserrat', sans-serif; font-size: 2.2rem; font-weight: 800; color: var(--text-primary); margin-bottom: 12px; }
.catalog-desc { font-size: 0.95rem; color: var(--text-secondary); max-width: 600px; }

.catalog-body { margin-top: 40px; }

.catalog-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 24px; }
@media (min-width: 768px) { .catalog-grid { grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); } }
@media (min-width: 1024px) { .catalog-grid { grid-template-columns: repeat(6, 1fr); } }

.empty-state { grid-column: 1 / -1; text-align: center; padding: 60px 0; color: var(--text-muted); }

.pagination-wrap { display: flex; justify-content: center; gap: 8px; margin-top: 50px; }
.page-btn { min-width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; background: var(--bg-secondary); border: 1px solid var(--border-card); border-radius: var(--radius-sm); color: var(--text-secondary); font-weight: 600; cursor: pointer; transition: var(--transition-base); }
.page-btn:hover:not(:disabled) { background: var(--bg-tertiary); border-color: var(--border-active); color: var(--accent-primary); }
.page-btn.active { background: var(--accent-primary); color: #000; border-color: var(--accent-primary); }
.page-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.page-dots { display: flex; align-items: center; justify-content: center; color: var(--text-muted); padding: 0 5px; }
</style>
