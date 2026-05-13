<template>
  <NuxtLink :to="linkTo" class="comic-card-wrapper">
    <div class="hover-3d comic-card" :class="`rarity-${rarity}`">
      <!-- Cover -->
      <div class="comic-card__cover">
        <img
          :src="comic.coverUrl || '/images/no-cover.jpg'"
          :alt="comic.title"
          class="cover-img"
          loading="lazy"
          @error="onImgError"
        />
        <div class="comic-card__overlay">
          <div class="comic-card__badges">
            <span v-if="comic.isHot" class="badge badge-hot">🔥 HOT</span>
            <span v-if="comic.isNew" class="badge badge-new">MỚI</span>
            <span v-if="comic.status === 'completed'" class="badge badge-completed">FULL</span>
          </div>
          <div class="comic-card__info">
            <h3 class="comic-card__title">{{ comic.title }}</h3>
            <p class="comic-card__chapter">
              <span class="chapter-dot" />
              Ch.{{ comic.latestChapter || 0 }} · {{ timeAgo(comic.updatedAt) }}
            </p>
            <div class="comic-card__stats">
              <span>👁 {{ formatViews(comic.views) }}</span>
              <span v-if="comic.rating">⭐ {{ comic.rating }}</span>
            </div>
          </div>
        </div>
        <!-- Glow shimmer on hover -->
        <div class="comic-card__glow" />
      </div>
      <!-- 9-zone hover grid for CSS :has() 3D tilt -->
      <div class="hover-zones" aria-hidden="true">
        <span v-for="i in 9" :key="i" :class="`zone-${i}`" />
      </div>
    </div>

    <!-- Rank badge (for Top10) -->
    <div v-if="rank" class="rank-badge" :class="rankClass">{{ rank }}</div>

    <!-- Below card info (optional) -->
    <div v-if="showInfo" class="card-below-info">
      <p class="below-title">{{ comic.title }}</p>
      <p class="below-chapter">Ch.{{ comic.latestChapter }}</p>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Comic {
  id: number | string
  title: string
  slug: string
  coverUrl?: string
  type: 'comic' | 'novel'
  latestChapter?: number
  updatedAt?: string
  views?: number
  rating?: number
  status?: 'ongoing' | 'completed' | 'hiatus'
  isHot?: boolean
  isNew?: boolean
}

const props = withDefaults(defineProps<{
  comic: Comic
  rank?: number
  showInfo?: boolean
  rarity?: 'normal' | 'hot' | 'gold'
}>(), {
  showInfo: false,
  rarity: 'normal',
})

const linkTo = computed(() =>
  props.comic.type === 'novel'
    ? `/truyen-chu/${props.comic.slug}`
    : `/truyen-tranh/${props.comic.slug}`
)

const rankClass = computed(() => {
  if (!props.rank) return ''
  if (props.rank === 1) return 'rank-1'
  if (props.rank === 2) return 'rank-2'
  if (props.rank === 3) return 'rank-3'
  return 'rank-other'
})

const formatViews = (v?: number) => {
  if (!v) return '0'
  if (v >= 1_000_000) return `${(v / 1_000_000).toFixed(1)}M`
  if (v >= 1_000) return `${(v / 1_000).toFixed(0)}K`
  return v.toString()
}

const timeAgo = (dateStr?: string) => {
  if (!dateStr) return 'Vừa cập nhật'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 7) return date.toLocaleDateString('vi-VN')
  if (days > 0) return `${days} ngày trước`
  if (hours > 0) return `${hours} giờ trước`
  if (minutes > 0) return `${minutes} phút trước`
  return 'Vừa xong'
}

const onImgError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.src = '/images/no-cover.jpg'
}
</script>

<style scoped>
.comic-card-wrapper {
  display: block;
  position: relative;
  user-select: none;
}

/* ---- 3D tilt via CSS :has() + hover zones ---- */
.hover-3d {
  --rx: 0deg;
  --ry: 0deg;
  transform-style: preserve-3d;
  transform: perspective(700px) rotateX(var(--rx)) rotateY(var(--ry));
  transition: transform 0.12s ease, box-shadow 0.25s ease;
  border-radius: var(--radius-md);
  position: relative;
  aspect-ratio: 2/3;
  overflow: hidden;
  background: var(--bg-card);
  border: 1px solid var(--border-card);
}

/* Zone-based tilt using :has() */
.hover-3d:has(.zone-1:hover) { --rx: 12deg;  --ry: -12deg; }
.hover-3d:has(.zone-2:hover) { --rx: 12deg;  --ry: 0deg;   }
.hover-3d:has(.zone-3:hover) { --rx: 12deg;  --ry: 12deg;  }
.hover-3d:has(.zone-4:hover) { --rx: 0deg;   --ry: -12deg; }
.hover-3d:has(.zone-5:hover) { --rx: 0deg;   --ry: 0deg;   }
.hover-3d:has(.zone-6:hover) { --rx: 0deg;   --ry: 12deg;  }
.hover-3d:has(.zone-7:hover) { --rx: -12deg; --ry: -12deg; }
.hover-3d:has(.zone-8:hover) { --rx: -12deg; --ry: 0deg;   }
.hover-3d:has(.zone-9:hover) { --rx: -12deg; --ry: 12deg;  }

.hover-3d:hover {
  box-shadow: 0 24px 64px rgba(0,0,0,0.75), 0 0 24px rgba(156,167,100,0.25);
}

/* Rarity borders */
.rarity-hot  { border-color: rgba(255,77,77,0.4); }
.rarity-gold { border-color: rgba(255,215,0,0.35); }

/* Cover layer */
.comic-card__cover {
  position: absolute;
  inset: 0;
  z-index: 1;
}

.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s cubic-bezier(0.4,0,0.2,1);
}

.hover-3d:hover .cover-img { transform: scale(1.07); }

/* Gradient overlay */
.comic-card__overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to top,
    rgba(8,10,16,0.97) 0%,
    rgba(8,10,16,0.65) 42%,
    rgba(8,10,16,0.08) 68%,
    transparent 100%
  );
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 10px;
  z-index: 2;
}

.comic-card__badges { display: flex; gap: 4px; flex-wrap: wrap; }

.comic-card__info { display: flex; flex-direction: column; gap: 5px; }

.comic-card__title {
  font-family: 'Montserrat', sans-serif;
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-shadow: 0 1px 6px rgba(0,0,0,0.9);
}

.comic-card__chapter {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.7rem;
  color: var(--accent-primary);
  font-weight: 600;
}

.chapter-dot {
  width: 5px;
  height: 5px;
  background: var(--accent-primary);
  border-radius: 50%;
  animation: glowPulse 2s ease-in-out infinite;
}

.comic-card__stats {
  display: flex;
  gap: 8px;
  font-size: 0.67rem;
  color: var(--text-secondary);
}

/* Glow shimmer */
.comic-card__glow {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(156,167,100,0.18), transparent 60%);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 3;
  pointer-events: none;
}
.hover-3d:hover .comic-card__glow { opacity: 1; }

/* 9-zone grid */
.hover-zones {
  position: absolute;
  inset: 0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  z-index: 10;
}
.hover-zones span { display: block; cursor: pointer; }

/* Below card info */
.card-below-info { padding: 6px 2px 0; }
.below-title { font-size: 0.8rem; font-weight: 600; color: var(--text-primary); line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; }
.below-chapter { font-size: 0.7rem; color: var(--accent-primary); margin-top: 2px; }

/* Rank badge */
.rank-badge {
  position: absolute;
  top: -6px;
  left: -6px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Montserrat', sans-serif;
  font-weight: 900;
  font-size: 0.78rem;
  z-index: 20;
  border: 2px solid var(--bg-primary);
}
.rank-1 { background: linear-gradient(135deg, #FFD700, #FFA500); color: #000; }
.rank-2 { background: linear-gradient(135deg, #C0C0C0, #909090); color: #000; }
.rank-3 { background: linear-gradient(135deg, #CD7F32, #A0522D); color: #fff; }
.rank-other { background: var(--bg-tertiary); color: var(--text-secondary); }
</style>
