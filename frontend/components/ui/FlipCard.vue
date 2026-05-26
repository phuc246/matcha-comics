<template>
  <div class="flip-card-container">
    <div class="flip-card-inner">
      <!-- Front Side -->
      <div class="flip-card-front">
        <img :src="comic.coverUrl || '/images/no-cover.jpg'" :alt="comic.title" class="cover-image" loading="lazy" />
        <div class="front-overlay">
          <div class="badges">
            <span v-if="comic.isHot" class="badge badge-hot">🔥</span>
            <span v-if="comic.status === 'completed'" class="badge badge-completed">FULL</span>
          </div>
          <div class="front-info">
            <h3 class="comic-title">{{ comic.title }}</h3>
            <p class="comic-meta">Ch.{{ comic.latestChapter }}</p>
          </div>
        </div>
      </div>
      
      <!-- Back Side -->
      <div class="flip-card-back">
        <div class="back-content">
          <h3 class="back-title">{{ comic.title }}</h3>
          
          <div class="stats">
            <span>👁 {{ formatViews(comic.views) }}</span>
            <span v-if="comic.rating">⭐ {{ comic.rating }}</span>
          </div>
          
          <div class="genres">
            <span v-for="g in comic.genres?.slice(0, 3)" :key="g" class="genre-tag">{{ g }}</span>
          </div>
          
          <p class="description">{{ comic.description || 'Chưa có mô tả cho truyện này...' }}</p>
          
          <NuxtLink :to="`/truyen-tranh/${comic.slug}`" class="btn-read">
            Đọc Ngay
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Comic {
  id: number | string
  title: string
  slug: string
  coverUrl?: string
  latestChapter?: number
  views?: number
  rating?: number
  status?: string
  isHot?: boolean
  description?: string
  genres?: string[]
}

defineProps<{ comic: Comic }>()

const formatViews = (v?: number) => {
  if (!v) return '0'
  if (v >= 1_000_000) return `${(v/1_000_000).toFixed(1)}M`
  if (v >= 1_000) return `${(v/1_000).toFixed(0)}K`
  return v.toString()
}
</script>

<style scoped>
.flip-card-container {
  background-color: transparent;
  width: 100%;
  aspect-ratio: 2/3;
  perspective: 1000px; /* Remove this if you don't want the 3D effect */
  cursor: pointer;
}

.flip-card-inner {
  position: relative;
  width: 100%;
  height: 100%;
  text-align: center;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  transform-style: preserve-3d;
  border-radius: var(--radius-md);
  box-shadow: 0 4px 8px rgba(0,0,0,0.5);
}

.flip-card-container:hover .flip-card-inner {
  transform: rotateY(180deg);
  box-shadow: 0 16px 32px rgba(0,0,0,0.8), 0 0 20px rgba(156,167,100,0.3);
}

.flip-card-front, .flip-card-back {
  position: absolute;
  width: 100%;
  height: 100%;
  -webkit-backface-visibility: hidden; /* Safari */
  backface-visibility: hidden;
  border-radius: var(--radius-md);
  overflow: hidden;
}

/* Front Side */
.flip-card-front {
  background-color: var(--bg-card);
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.front-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(13,15,20,0.95) 0%, rgba(13,15,20,0.4) 40%, transparent 100%);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 12px;
}

.badges { display: flex; gap: 4px; flex-wrap: wrap; align-items: center; }
.badge { font-size: 0.65rem; font-weight: 800; padding: 3px 6px; border-radius: 4px; color: #fff; }
.badge-hot {
  background: linear-gradient(135deg, #ff4d4d, #cc0000);
  width: 22px;
  height: 22px;
  padding: 0;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  box-shadow: 0 0 10px rgba(255, 77, 77, 0.4);
}
.badge-completed {
  background: rgba(0, 0, 0, 0.45);
  color: rgba(255, 255, 255, 0.55);
  border: 1px solid rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(4px);
}

.front-info { text-align: left; }
.comic-title { font-size: 0.9rem; font-weight: 700; color: #fff; line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; text-shadow: 0 2px 4px rgba(0,0,0,0.8); }
.comic-meta { font-size: 0.75rem; color: var(--accent-primary); font-weight: 600; margin-top: 4px; }

/* Back Side */
.flip-card-back {
  background-color: var(--bg-secondary);
  color: white;
  transform: rotateY(180deg);
  border: 1px solid var(--border-card);
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  height: 100%;
}

.back-title {
  font-family: 'Montserrat', sans-serif;
  font-size: 1rem;
  font-weight: 800;
  color: var(--accent-primary);
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}

.stats {
  display: flex;
  justify-content: center;
  gap: 12px;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.genres {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 4px;
}
.genre-tag { font-size: 0.65rem; background: var(--bg-tertiary); padding: 2px 8px; border-radius: 12px; border: 1px solid var(--border-subtle); color: var(--text-secondary); }

.description {
  font-size: 0.8rem;
  color: var(--text-secondary);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
  text-align: left;
}

.btn-read {
  display: block;
  width: 100%;
  padding: 8px 0;
  background: var(--accent-primary);
  color: var(--bg-primary);
  text-align: center;
  border-radius: var(--radius-sm);
  font-weight: 700;
  font-size: 0.85rem;
  transition: var(--transition-base);
  margin-top: auto;
}
.btn-read:hover { background: var(--accent-hover); transform: translateY(-2px); }
</style>
