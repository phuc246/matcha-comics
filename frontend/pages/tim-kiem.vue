<template>
  <div class="search-page">
    <div class="catalog-hero">
      <div class="container">
        <h1 class="catalog-title">🔍 Tìm Kiếm</h1>
        <p class="catalog-desc">Kết quả tìm kiếm cho: <strong>"{{ searchQuery }}"</strong></p>
      </div>
    </div>

    <div class="container catalog-body">
      <!-- Tabs -->
      <div class="type-tabs" v-if="searchQuery">
        <button class="type-tab" :class="{ active: currentType === 'comic' }" @click="currentType = 'comic'">
          Truyện Tranh ({{ comicResults.length }})
        </button>
        <button class="type-tab" :class="{ active: currentType === 'novel' }" @click="currentType = 'novel'">
          Truyện Chữ ({{ novelResults.length }})
        </button>
      </div>

      <!-- Grid -->
      <div class="comic-grid" v-if="!loading && currentResults.length > 0">
        <ComicCard3D
          v-for="(item, i) in currentResults"
          :key="item.id"
          :comic="item"
          class="animate-fade-up"
          :style="{ animationDelay: `${(i % 12) * 0.04}s` }"
        />
      </div>
      <div class="comic-grid" v-else-if="loading">
        <ComicCardSkeleton v-for="i in 6" :key="i" />
      </div>

      <!-- Empty state -->
      <div v-if="!loading && currentResults.length === 0" class="empty-state">
        <p class="empty-icon">📭</p>
        <p class="empty-text">Không tìm thấy kết quả nào phù hợp với "{{ searchQuery }}"</p>
        <NuxtLink to="/" class="btn-home mt-4">Về trang chủ</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { mockComics, mockNovels } from '~/data/mock'

const route = useRoute()
const loading = ref(true)
const currentType = ref('comic')

const searchQuery = computed(() => (route.query.q as string) || '')

useHead(() => ({ title: `Tìm kiếm: ${searchQuery.value} — Matcha Comic` }))

const comicResults = ref<any[]>([])
const novelResults = ref<any[]>([])

const performSearch = () => {
  const q = searchQuery.value.toLowerCase()
  if (!q) {
    comicResults.value = []
    novelResults.value = []
    return
  }
  
  comicResults.value = mockComics.filter(c => c.title.toLowerCase().includes(q))
  novelResults.value = mockNovels.filter(n => n.title.toLowerCase().includes(q))
  
  // Auto switch to novel tab if no comics found but novels found
  if (comicResults.value.length === 0 && novelResults.value.length > 0) {
    currentType.value = 'novel'
  } else {
    currentType.value = 'comic'
  }
}

const currentResults = computed(() => {
  return currentType.value === 'comic' ? comicResults.value : novelResults.value
})

watch(searchQuery, () => {
  loading.value = true
  setTimeout(() => {
    performSearch()
    loading.value = false
  }, 300)
})

onMounted(async () => {
  await new Promise(r => setTimeout(r, 400))
  performSearch()
  loading.value = false
})
</script>

<style scoped>
.catalog-hero { padding: 40px 0 28px; background: linear-gradient(180deg, var(--bg-secondary), var(--bg-primary)); border-bottom: 1px solid var(--border-subtle); }
.catalog-title { font-size: clamp(1.6rem, 4vw, 2.4rem); font-weight: 900; color: var(--text-primary); }
.catalog-desc { font-size: 0.95rem; color: var(--text-muted); margin-top: 8px; }
.catalog-body { padding-top: 28px; padding-bottom: 60px; min-height: 50vh; }

.type-tabs { display: flex; gap: 12px; margin-bottom: 30px; border-bottom: 1px solid var(--border-card); padding-bottom: 16px; }
.type-tab { padding: 10px 24px; border-radius: 100px; font-weight: 700; font-size: 0.95rem; background: var(--bg-secondary); border: 1px solid var(--border-card); color: var(--text-muted); cursor: pointer; transition: all 0.2s; }
.type-tab:hover { color: var(--text-primary); border-color: var(--border-subtle); }
.type-tab.active { background: var(--accent-primary); color: #000; border-color: transparent; }

.empty-state { text-align: center; padding: 80px 0; }
.empty-icon { font-size: 3rem; margin-bottom: 12px; }
.empty-text { color: var(--text-muted); font-size: 1rem; }
.mt-4 { margin-top: 16px; }
.btn-home { display: inline-block; padding: 10px 24px; background: var(--bg-tertiary); border: 1px solid var(--border-card); border-radius: 100px; color: var(--text-primary); transition: var(--transition-base); }
.btn-home:hover { background: var(--accent-primary); color: #000; }
</style>
