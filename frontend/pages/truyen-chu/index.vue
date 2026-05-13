<template>
  <div class="catalog-page">
    <!-- Page Hero -->
    <div class="catalog-hero">
      <div class="container">
        <h1 class="catalog-title">📖 Matcha Comics Phòng Triển Lãm (Đang xây dựng)</h1>
        <p class="catalog-desc">Tiểu thuyết, Tiên hiệp, Ngôn tình — Nơi lưu giữ những kiệt tác nghệ thuật</p>
        <div class="hero-spacer"></div>
      </div>
    </div>

    <div class="container catalog-body">
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
        <div class="filter-right">
          <select class="filter-select" v-model="currentStatus">
            <option value="">Tất cả trạng thái</option>
            <option value="ongoing">Đang tiến hành</option>
            <option value="completed">Đã hoàn thành</option>
            <option value="hiatus">Tạm dừng</option>
          </select>
          <select class="filter-select" v-model="currentGenre">
            <option value="">Tất cả thể loại</option>
            <option v-for="g in genres" :key="g.slug" :value="g.slug">{{ g.name }}</option>
          </select>
        </div>
      </div>

      <!-- Results count -->
      <div class="results-meta">
        <span class="results-count">{{ filteredNovels.length }} bộ truyện</span>
      </div>

      <!-- Grid -->
      <div class="comic-grid" v-if="!loading">
        <NuxtLink
          v-for="(novel, i) in paginatedNovels"
          :key="novel.id"
          :to="`/truyen-chu/${novel.slug}`"
          class="novel-card-link animate-fade-up"
          :style="{ animationDelay: `${(i % 12) * 0.04}s` }"
        >
          <Book3D
            :title="novel.title"
            :cover-url="novel.coverUrl"
            size="md"
          />
          <div class="novel-info">
            <h3 class="novel-title">{{ novel.title }}</h3>
            <p class="novel-meta">Ch.{{ novel.latestChapter }}</p>
          </div>
        </NuxtLink>
      </div>
      <div class="comic-grid" v-else>
        <ComicCardSkeleton v-for="i in 18" :key="i" />
      </div>

      <!-- Empty state -->
      <div v-if="!loading && filteredNovels.length === 0" class="empty-state">
        <p class="empty-icon">📭</p>
        <p class="empty-text">Không tìm thấy truyện nào</p>
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
import { useRoute, useRouter } from 'vue-router'
import { mockNovels, mockGenres } from '~/data/mock'

useHead({ title: 'Truyện Chữ — Matcha Comic' })

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const novels = ref<any[]>([])
const genres = ref(mockGenres)
const currentSort = ref((route.query.sort as string) || 'latest')
const currentStatus = ref((route.query.status as string) || '')
const currentGenre = ref((route.query.genre as string) || '')
const currentPage = ref(1)
const perPage = 18

const sortOptions = [
  { label: 'Mới nhất', value: 'latest', icon: '🆕' },
  { label: 'Xem nhiều', value: 'views', icon: '👁' },
  { label: 'Đánh giá', value: 'rating', icon: '⭐' },
  { label: 'Mới cập nhật', value: 'updated', icon: '🔄' },
]

const setSort = (v: string) => { currentSort.value = v; currentPage.value = 1 }

const filteredNovels = computed(() => {
  let list = novels.value.filter(c => c.type === 'novel')
  if (currentStatus.value) list = list.filter(c => c.status === currentStatus.value)
  if (currentGenre.value) list = list.filter(c => c.genres?.some((g: string) => g.toLowerCase().replace(/\s/g,'-') === currentGenre.value))
  if (currentSort.value === 'views') list.sort((a, b) => b.views - a.views)
  else if (currentSort.value === 'rating') list.sort((a, b) => b.rating - a.rating)
  return list
})

const totalPages = computed(() => Math.ceil(filteredNovels.value.length / perPage))
const paginatedNovels = computed(() => filteredNovels.value.slice((currentPage.value - 1) * perPage, currentPage.value * perPage))
const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

watch([currentSort, currentStatus, currentGenre], () => { currentPage.value = 1 })

onMounted(async () => {
  await new Promise(r => setTimeout(r, 500))
  novels.value = mockNovels
  loading.value = false
})
</script>

<style scoped>
.catalog-hero { padding: 40px 0 28px; background: linear-gradient(180deg, var(--bg-secondary), var(--bg-primary)); border-bottom: 1px solid var(--border-subtle); }
.catalog-title { font-size: clamp(1.6rem, 4vw, 2.4rem); font-weight: 900; color: var(--text-primary); }
.catalog-desc { font-size: 0.95rem; color: var(--text-muted); margin-top: 8px; }
.hero-spacer { height: 120px; }
.catalog-body { padding-top: 28px; padding-bottom: 60px; }
.filter-bar { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 20px; flex-wrap: wrap; }
.filter-left { display: flex; gap: 8px; flex-wrap: wrap; }
.filter-btn { padding: 7px 14px; border-radius: var(--radius-sm); border: 1px solid var(--border-card); background: var(--bg-secondary); color: var(--text-secondary); font-size: 0.84rem; cursor: pointer; transition: var(--transition-base); }
.filter-btn:hover { border-color: var(--border-subtle); color: var(--text-primary); }
.filter-btn.active { background: linear-gradient(135deg, var(--accent-dark), var(--accent-primary)); color: var(--bg-primary); border-color: transparent; font-weight: 600; }
.filter-right { display: flex; gap: 8px; flex-wrap: wrap; }
.filter-select { padding: 7px 12px; border-radius: var(--radius-sm); border: 1px solid var(--border-card); background: var(--bg-secondary); color: var(--text-secondary); font-size: 0.84rem; cursor: pointer; outline: none; transition: var(--transition-base); }
.filter-select:focus { border-color: var(--border-active); color: var(--text-primary); }
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

.comic-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 32px 16px; justify-items: center; }
.novel-card-link { display: flex; flex-direction: column; align-items: center; gap: 12px; text-decoration: none; width: 100%; }
.novel-info { text-align: center; padding: 0 4px; }
.novel-title { font-size: 0.85rem; font-weight: 700; color: var(--text-primary); line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; margin-bottom: 4px; }
.novel-meta { font-size: 0.75rem; color: var(--accent-primary); font-weight: 600; }
</style>
