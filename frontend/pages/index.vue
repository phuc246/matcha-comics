<template>
  <div>
    <!-- Hero Banner -->
    <HeroSlider :slides="slides" />

    <!-- Sections -->
    <div class="divider" style="margin: 0;" />

    <SectionTop10 :comics="topComics" :loading="loadingTop" />

    <div class="glow-line" />

    <SectionHot :comics="hotComics" :loading="loadingHot" />

    <div class="glow-line" />

    <SectionLatest
      :comics="latestComics"
      :loading="loadingLatest"
      title="Truyện Tranh Mới Nhất"
      icon="🎨"
      view-all-link="/truyen-tranh?sort=latest"
    />

    <div class="glow-line" />

    <SectionLatest
      :comics="latestNovels"
      :loading="loadingNovels"
      title="Tiểu Thuyết Mới Nhất"
      icon="📖"
      view-all-link="/truyen-chu?sort=latest"
    />

    <div class="glow-line" />

    <SectionCompleted :comics="completedComics" :loading="false" />


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { mockComics, mockNovels, mockSlides, mockGenres } from '~/data/mock'

useHead({
  title: 'Matcha Comic - Đọc Truyện Tranh & Tiểu Thuyết Online Miễn Phí',
  meta: [
    { name: 'description', content: 'Đọc truyện tranh manga, manhwa và tiểu thuyết online miễn phí. Cập nhật nhanh nhất với hàng nghìn bộ truyện chất lượng cao.' },
  ],
})

const slides = ref(mockSlides)
const genres = ref(mockGenres)

// Simulated loading states
const loadingTop = ref(true)
const loadingHot = ref(true)
const loadingLatest = ref(true)
const loadingNovels = ref(true)

const topComics = ref<any[]>([])
const hotComics = ref<any[]>([])
const latestComics = ref<any[]>([])
const latestNovels = ref<any[]>([])
const completedComics = ref<any[]>([])
const infiniteComics = ref<any[]>([])
const loadingInfinite = ref(false)
const allLoaded = ref(false)
const scrollTrigger = ref<HTMLElement | null>(null)

const loadMore = async () => {
  if (loadingInfinite.value || allLoaded.value) return
  loadingInfinite.value = true
  
  // Simulate API delay
  await new Promise(r => setTimeout(r, 800))
  
  const nextItems = mockComics.slice(0, 6) // Mock loading more
  infiniteComics.value.push(...nextItems)
  
  if (infiniteComics.value.length > 30) {
    allLoaded.value = true
  }
  loadingInfinite.value = false
}

const { get } = useApi()

onMounted(async () => {
  try {
    // 1. Fetch Featured (Hero Slider) - Trending in the last 24h
    const featuredData = await get<any[]>('/stories?limit=5&sort=trending')
    if (featuredData && featuredData.length > 0) {
      // Merge with mock if not enough
      const realIds = new Set(featuredData.map(s => s.slug))
      const combinedSlides = [
        ...featuredData.map(s => ({ ...s, isHot: true })),
        ...mockSlides.filter(s => !realIds.has(s.slug))
      ]
      slides.value = combinedSlides.slice(0, 5)
    } else {
      slides.value = mockSlides
    }

    // 2. Fetch Top Comics - Ranking
    const topData = await get<any[]>('/stories?type=comic&limit=10&sort=views')
    if (topData && topData.length > 0) {
      const realSlugs = new Set(topData.map(s => s.slug))
      topComics.value = [
        ...topData,
        ...mockComics.filter(c => !realSlugs.has(c.slug))
      ].slice(0, 10)
    } else {
      topComics.value = [...mockComics].sort((a, b) => b.views - a.views).slice(0, 10)
    }
    loadingTop.value = false

    // 3. Fetch Hot Comics
    const hotData = await get<any[]>('/stories?limit=6&sort=views')
    if (hotData && hotData.length > 0) {
      const realSlugs = new Set(hotData.map(s => s.slug))
      hotComics.value = [
        ...hotData,
        ...mockComics.filter(c => !realSlugs.has(c.slug))
      ].slice(0, 6)
    } else {
      hotComics.value = mockComics.filter(c => c.isHot).slice(0, 6)
    }
    loadingHot.value = false

    // 4. Fetch Latest Comics
    const latestC = await get<any[]>('/stories?type=comic&limit=12')
    if (latestC && latestC.length > 0) {
      const realSlugs = new Set(latestC.map(s => s.slug))
      latestComics.value = [
        ...latestC,
        ...mockComics.filter(c => !realSlugs.has(c.slug))
      ].slice(0, 12)
    } else {
      latestComics.value = mockComics.slice(0, 12)
    }
    loadingLatest.value = false

    // 5. Fetch Latest Novels
    const latestN = await get<any[]>('/stories?type=novel&limit=12')
    if (latestN && latestN.length > 0) {
      const realSlugs = new Set(latestN.map(s => s.slug))
      latestNovels.value = [
        ...latestN,
        ...mockNovels.filter(n => !realSlugs.has(n.slug))
      ].slice(0, 12)
    } else {
      latestNovels.value = mockNovels.slice(0, 12)
    }
    loadingNovels.value = false

    // 6. Fetch Completed
    const completed = await get<any[]>('/stories?status=completed&limit=10')
    if (completed && completed.length > 0) {
      completedComics.value = completed
    } else {
      completedComics.value = mockComics.filter(c => c.status === 'completed').slice(0, 10)
    }
    
  } catch (err) {
    console.error('Error loading homepage data:', err)
    // Fallback to mock on error
    slides.value = mockSlides
    topComics.value = mockComics.slice(0, 10)
    hotComics.value = mockComics.slice(0, 6)
    latestComics.value = mockComics.slice(0, 12)
    latestNovels.value = mockNovels.slice(0, 12)
    completedComics.value = mockComics.slice(0, 10)
  }

  // Initial load for infinite
  loadMore()

  // Intersection Observer for Infinite Scroll
  const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) {
      loadMore()
    }
  }, { threshold: 0.1 })

  if (scrollTrigger.value) {
    observer.observe(scrollTrigger.value)
  }
})
</script>

<style scoped>
.genre-section { position: relative; }
.genre-section::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(ellipse at 50% 100%, rgba(156,167,100,0.03) 0%, transparent 60%);
  pointer-events: none;
}
.genre-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
  gap: 16px;
}
@media (min-width: 768px) { .genre-grid { grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); } }

.genre-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 24px 12px;
  background: rgba(20, 23, 32, 0.4);
  backdrop-filter: blur(10px);
  border: 1px solid var(--border-card);
  border-radius: var(--radius-lg);
  text-align: center;
  transition: var(--transition-base);
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.genre-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--accent-primary), transparent);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.genre-card:hover {
  background: rgba(28, 32, 48, 0.6);
  border-color: var(--border-active);
  transform: translateY(-4px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.5), 0 0 20px rgba(156,167,100,0.15);
}

.genre-card:hover::before {
  opacity: 1;
}

.genre-icon { font-size: 2.2rem; filter: drop-shadow(0 4px 6px rgba(0,0,0,0.4)); transition: transform 0.3s ease; }
.genre-card:hover .genre-icon { transform: scale(1.1) rotate(5deg); }
.genre-name { font-size: 0.85rem; font-weight: 700; color: var(--text-primary); letter-spacing: 0.02em; }
.genre-count { font-size: 0.72rem; color: var(--text-muted); background: rgba(255,255,255,0.05); padding: 2px 8px; border-radius: 12px; }

.scroll-trigger {
  padding: 40px 0;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100px;
}

.loader-spinner-small {
  width: 24px;
  height: 24px;
  border: 2px solid rgba(156,167,100,0.2);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.end-text {
  color: var(--text-muted);
  font-size: 0.9rem;
  font-style: italic;
}
</style>
