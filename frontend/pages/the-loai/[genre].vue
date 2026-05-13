<template>
  <div class="genre-detail-page">
    <div class="catalog-hero">
      <div class="container">
        <h1 class="catalog-title">{{ genreName }}</h1>
        <p class="catalog-desc">Tuyển tập truyện tranh và tiểu thuyết thể loại {{ genreName.toLowerCase() }}</p>
      </div>
    </div>

    <div class="container catalog-body">
      <!-- Tabs -->
      <div class="type-tabs">
        <button class="type-tab" :class="{ active: currentType === 'comic' }" @click="currentType = 'comic'">Truyện Tranh</button>
        <button class="type-tab" :class="{ active: currentType === 'novel' }" @click="currentType = 'novel'">Truyện Chữ</button>
      </div>

      <!-- Filters -->
      <div class="filter-bar">
        <div class="filter-left">
          <button
            v-for="sort in sortOptions"
            :key="sort.value"
            class="filter-btn"
            :class="{ active: currentSort === sort.value }"
            @click="setSort(sort.value)"
          >
            {{ sort.icon }} {{ sort.label }}
          </button>
        </div>
      </div>

      <div class="results-meta">
        <span class="results-count">{{ filteredItems.length }} bộ truyện</span>
      </div>

      <!-- Grid -->
      <div class="comic-grid" v-if="!loading">
        <ComicCard3D
          v-for="(item, i) in paginatedItems"
          :key="item.id"
          :comic="item"
          class="animate-fade-up"
          :style="{ animationDelay: `${(i % 12) * 0.04}s` }"
        />
      </div>
      <div class="comic-grid" v-else>
        <ComicCardSkeleton v-for="i in 12" :key="i" />
      </div>

      <div v-if="!loading && filteredItems.length === 0" class="empty-state">
        <p class="empty-icon">📭</p>
        <p class="empty-text">Chưa có truyện nào thuộc thể loại này</p>
      </div>

      <!-- Pagination -->
      <div class="pagination" v-if="totalPages > 1">
        <button class="page-btn" :disabled="currentPage === 1" @click="currentPage--">←</button>
        <button
          v-for="p in visiblePages"
          :key="p"
          class="page-btn"
          :class="{ active: p === currentPage }"
          @click="currentPage = p"
        >{{ p }}</button>
        <button class="page-btn" :disabled="currentPage === totalPages" @click="currentPage++">→</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { mockComics, mockNovels, mockGenres } from '~/data/mock'

const route = useRoute()
const genreSlug = route.params.genre as string

const genreName = computed(() => {
  const g = mockGenres.find(x => x.slug === genreSlug)
  return g ? g.name : genreSlug.replace(/-/g, ' ')
})

useHead(() => ({ title: `${genreName.value} — Matcha Comic` }))

const loading = ref(true)
const currentType = ref('comic') // comic | novel
const currentSort = ref('latest')
const currentPage = ref(1)
const perPage = 18

const allComics = ref<any[]>([])
const allNovels = ref<any[]>([])

const sortOptions = [
  { label: 'Mới nhất', value: 'latest', icon: '🆕' },
  { label: 'Xem nhiều', value: 'views', icon: '👁' },
  { label: 'Đánh giá', value: 'rating', icon: '⭐' }
]

const setSort = (v: string) => { currentSort.value = v; currentPage.value = 1 }

const filteredItems = computed(() => {
  const sourceList = currentType.value === 'comic' ? allComics.value : allNovels.value
  let list = sourceList.filter(c => c.genres?.some((g: string) => g.toLowerCase().replace(/\s/g,'-') === genreSlug))
  
  if (currentSort.value === 'views') list.sort((a, b) => b.views - a.views)
  else if (currentSort.value === 'rating') list.sort((a, b) => b.rating - a.rating)
  
  return list
})

const totalPages = computed(() => Math.ceil(filteredItems.value.length / perPage))
const paginatedItems = computed(() => filteredItems.value.slice((currentPage.value - 1) * perPage, currentPage.value * perPage))
const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

watch([currentType, currentSort], () => { currentPage.value = 1 })

onMounted(async () => {
  await new Promise(r => setTimeout(r, 400))
  allComics.value = mockComics
  allNovels.value = mockNovels
  loading.value = false
})
</script>

<style scoped>
.catalog-hero { padding: 40px 0 28px; background: linear-gradient(180deg, var(--bg-secondary), var(--bg-primary)); border-bottom: 1px solid var(--border-subtle); text-transform: capitalize; }
.catalog-title { font-size: clamp(1.6rem, 4vw, 2.4rem); font-weight: 900; color: var(--text-primary); }
.catalog-desc { font-size: 0.95rem; color: var(--text-muted); margin-top: 8px; }
.catalog-body { padding-top: 28px; padding-bottom: 60px; }

.type-tabs { display: flex; gap: 12px; margin-bottom: 24px; border-bottom: 1px solid var(--border-card); padding-bottom: 16px; }
.type-tab { padding: 10px 24px; border-radius: 100px; font-weight: 700; font-size: 0.95rem; background: var(--bg-secondary); border: 1px solid var(--border-card); color: var(--text-muted); cursor: pointer; transition: all 0.2s; }
.type-tab:hover { color: var(--text-primary); border-color: var(--border-subtle); }
.type-tab.active { background: var(--accent-primary); color: #000; border-color: transparent; }

.filter-bar { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 20px; flex-wrap: wrap; }
.filter-left { display: flex; gap: 8px; flex-wrap: wrap; }
.filter-btn { padding: 7px 14px; border-radius: var(--radius-sm); border: 1px solid var(--border-card); background: var(--bg-secondary); color: var(--text-secondary); font-size: 0.84rem; cursor: pointer; transition: var(--transition-base); }
.filter-btn:hover { border-color: var(--border-subtle); color: var(--text-primary); }
.filter-btn.active { background: linear-gradient(135deg, var(--accent-dark), var(--accent-primary)); color: var(--bg-primary); border-color: transparent; font-weight: 600; }

.results-meta { margin-bottom: 16px; }
.results-count { font-size: 0.82rem; color: var(--text-muted); }
.empty-state { text-align: center; padding: 80px 0; }
.empty-icon { font-size: 3rem; margin-bottom: 12px; }
.empty-text { color: var(--text-muted); font-size: 1rem; }
.pagination { display: flex; justify-content: center; gap: 8px; margin-top: 48px; }
.page-btn { width: 40px; height: 40px; border-radius: var(--radius-sm); border: 1px solid var(--border-card); background: var(--bg-secondary); color: var(--text-secondary); cursor: pointer; font-size: 0.9rem; transition: var(--transition-base); }
.page-btn:hover:not(:disabled) { border-color: var(--border-active); color: var(--text-primary); }
.page-btn.active { background: var(--accent-primary); color: var(--bg-primary); border-color: transparent; font-weight: 700; }
.page-btn:disabled { opacity: 0.35; cursor: not-allowed; }
</style>
