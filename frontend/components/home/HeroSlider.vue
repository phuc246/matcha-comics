<template>
  <section class="hero-waka" aria-label="Truyện nổi bật">
    <!-- Dynamic Blurred Background -->
    <div class="hero-backdrop" :style="{ backgroundImage: `url(${current.coverUrl})` }"></div>
    <div class="hero-backdrop-overlay"></div>

    <div class="hero-container">
      
      <!-- Left: Information -->
      <div class="hero-info" :key="`info-${activeIdx}`">
        <div class="hero-badge-wrap animate-fade-up">
          <span class="badge badge-hot" v-if="current.isHot">ĐỀ XUẤT</span>
          <span class="badge badge-type">{{ current.type === 'novel' ? 'Tiểu thuyết' : 'Truyện tranh' }}</span>
        </div>
        
        <h1 class="hero-title animate-fade-up stagger-1">{{ current.title }}</h1>
        
        <p class="hero-desc animate-fade-up stagger-2" v-html="current.description"></p>
        
        <div class="hero-actions animate-fade-up stagger-3">
          <NuxtLink :to="getReadLink(current)" class="btn-read">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/>
              <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>
            </svg>
            Đọc truyện
          </NuxtLink>
          <NuxtLink :to="getDetailLink(current)" class="btn-detail" aria-label="Xem chi tiết">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/>
            </svg>
          </NuxtLink>
        </div>
      </div>

      <!-- Right: Cover Carousel -->
      <div class="hero-carousel">
        <div class="carousel-track">
          <div 
            v-for="(slide, i) in slides" 
            :key="slide.id"
            class="carousel-item"
            :class="getItemClass(i)"
            @click="goTo(i)"
          >
            <div class="carousel-card">
              <img :src="slide.coverUrl" :alt="slide.title" class="carousel-img" loading="lazy" />
              <div class="carousel-overlay" v-if="i !== activeIdx" />
            </div>
          </div>
        </div>

        <!-- Dots indicator -->
        <div class="carousel-dots">
          <button
            v-for="(_, i) in slides" :key="i"
            class="dot"
            :class="{ active: i === activeIdx }"
            @click="goTo(i)"
            :aria-label="`Slide ${i + 1}`"
          />
        </div>
      </div>

    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  slides: Array<{
    id: number | string
    title: string
    slug: string
    coverUrl: string
    type: 'comic' | 'novel'
    description?: string
    isHot?: boolean
  }>
}>()

const DURATION = 6000
const activeIdx = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const current = computed(() => props.slides[activeIdx.value])

const getReadLink = (item: any) =>
  item.type === 'novel' ? `/truyen-chu/${item.slug}/chapter-1` : `/truyen-tranh/${item.slug}/chapter-1`

const getDetailLink = (item: any) =>
  item.type === 'novel' ? `/truyen-chu/${item.slug}` : `/truyen-tranh/${item.slug}`

const goTo = (i: number) => {
  activeIdx.value = i
  resetTimer()
}

const next = () => goTo((activeIdx.value + 1) % props.slides.length)
const prev = () => goTo((activeIdx.value - 1 + props.slides.length) % props.slides.length)

// Determine the position of the item relative to activeIdx
const getItemClass = (i: number) => {
  const diff = i - activeIdx.value
  const len = props.slides.length
  
  if (diff === 0) return 'active'
  if (diff === 1 || diff === -len + 1) return 'next'
  if (diff === -1 || diff === len - 1) return 'prev'
  
  // For items further away, just hide them or push them further back
  if (diff === 2 || diff === -len + 2) return 'next-2'
  if (diff === -2 || diff === len - 2) return 'prev-2'
  
  return 'hidden'
}

const resetTimer = () => {
  if (timer) clearInterval(timer)
  timer = setInterval(next, DURATION)
}

onMounted(() => resetTimer())
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<style scoped>
.hero-waka {
  position: relative;
  background: var(--bg-primary);
  padding: 120px 0 80px;
  min-height: 70vh; /* Full screen feel */
  display: flex;
  align-items: center;
  overflow: hidden;
}

.hero-backdrop {
  position: absolute;
  inset: -10%;
  background-size: cover;
  background-position: center;
  filter: blur(40px);
  opacity: 0.3;
  transition: background-image 0.5s ease;
  z-index: 0;
}

.hero-backdrop-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to right, rgba(10,12,16,1) 0%, rgba(10,12,16,0.85) 40%, rgba(10,12,16,0.4) 100%);
  z-index: 1;
}

.hero-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  max-width: 1800px;
  margin: 0 auto;
  padding: 0 4vw;
  z-index: 2;
}

/* Left: Info */
.hero-info {
  width: 45%;
  max-width: 650px;
  z-index: 10;
}

.hero-badge-wrap {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.badge-type {
  background: rgba(255,255,255,0.1);
  color: var(--text-secondary);
  border: 1px solid var(--border-subtle);
}

.hero-title {
  font-family: 'Montserrat', sans-serif;
  font-size: clamp(2.2rem, 4.5vw, 3.8rem);
  font-weight: 900;
  color: var(--text-primary);
  line-height: 1.15;
  margin-bottom: 20px;
  text-shadow: 0 4px 20px rgba(0,0,0,0.5);
}

.hero-desc {
  font-size: 1.05rem;
  color: rgba(255,255,255,0.85);
  line-height: 1.7;
  margin-bottom: 36px;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.hero-desc :deep(*) {
  color: inherit !important;
  background-color: transparent !important;
}

.hero-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}

.btn-read {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 28px;
  background: var(--accent-primary);
  color: var(--bg-primary);
  font-family: 'Montserrat', sans-serif;
  font-weight: 700;
  font-size: 0.95rem;
  border-radius: 100px;
  transition: var(--transition-base);
}

.btn-read:hover {
  background: var(--accent-hover);
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(156,167,100,0.3);
}

.btn-detail {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: 46px;
  border-radius: 50%;
  background: rgba(255,255,255,0.05);
  color: var(--text-primary);
  transition: var(--transition-base);
  border: 1px solid var(--border-card);
}

.btn-detail:hover {
  background: rgba(255,255,255,0.1);
  transform: translateY(-2px);
}

/* Right: Carousel */
.hero-carousel {
  width: 50%;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 480px;
}

.carousel-track {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  perspective: 1000px;
}

.carousel-item {
  position: absolute;
  width: 280px;
  aspect-ratio: 2/3;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  border-radius: 12px;
  z-index: 1;
}

/* 3D Positioning */
.carousel-item.active {
  transform: translateX(0) scale(1.15);
  z-index: 5;
  box-shadow: 0 20px 60px rgba(0,0,0,0.9), 0 0 20px rgba(255,255,255,0.05);
}

.carousel-item.prev {
  transform: translateX(-60%) scale(0.85);
  z-index: 4;
}

.carousel-item.next {
  transform: translateX(60%) scale(0.85);
  z-index: 4;
}

.carousel-item.prev-2 {
  transform: translateX(-110%) scale(0.65);
  z-index: 3;
}

.carousel-item.next-2 {
  transform: translateX(110%) scale(0.65);
  z-index: 3;
}

.carousel-item.hidden {
  opacity: 0;
  transform: translateX(0) scale(0.5);
  pointer-events: none;
}

/* Card content */
.carousel-card {
  width: 100%;
  height: 100%;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  background: var(--bg-card);
}

.carousel-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* Removed carousel-overlay logic to keep images sharp */
.carousel-overlay {
  display: none;
}

/* Removed nav buttons */

/* Dots */
.carousel-dots {
  position: absolute;
  bottom: -20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 6px;
  z-index: 10;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 3px;
  background: rgba(255,255,255,0.2);
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.dot.active {
  width: 20px;
  background: var(--text-primary);
}

/* Responsive */
@media (max-width: 1024px) {
  .hero-container {
    flex-direction: column;
    text-align: center;
    padding-top: 20px;
  }
  
  .hero-info {
    max-width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .hero-carousel {
    width: 100%;
    height: 350px;
  }
  
  .carousel-item {
    width: 200px;
  }
}

@media (max-width: 640px) {
  .carousel-item.prev, .carousel-item.next {
    opacity: 0.5;
  }
  .carousel-item.prev-2, .carousel-item.next-2 {
    opacity: 0;
  }
}
</style>
