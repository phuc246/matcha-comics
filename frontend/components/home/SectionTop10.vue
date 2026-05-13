<template>
  <section class="page-section">
    <div class="container">
      <div class="section-header">
        <h2 class="section-title">🏆 Top 10 Tuần</h2>
        <NuxtLink to="/truyen-tranh?sort=views" class="section-view-all">Xem tất cả →</NuxtLink>
      </div>
      <div class="top10-layout">
        <!-- Left: Top 3 large -->
        <div class="top3-col">
          <template v-if="loading">
            <ComicCardSkeleton v-for="i in 3" :key="i" class="top3-skeleton" />
          </template>
          <template v-else>
            <div v-for="(comic, i) in comics.slice(0,3)" :key="comic.id" class="top3-card-wrap">
              <ComicCard3D :comic="comic" :rank="i + 1" rarity="gold" class="animate-scale-in" :style="{ animationDelay: `${i*0.08}s` }" />
            </div>
          </template>
        </div>

        <!-- Right: Rank 4-10 list -->
        <div class="rank-list">
          <template v-if="loading">
            <div v-for="i in 7" :key="i" class="rank-row-skeleton">
              <div class="skeleton rank-thumb-sk" />
              <div class="rank-info-sk">
                <div class="skeleton rank-title-sk" />
                <div class="skeleton rank-sub-sk" />
              </div>
            </div>
          </template>
          <template v-else>
            <NuxtLink
              v-for="(comic, i) in comics.slice(3, 10)"
              :key="comic.id"
              :to="`/truyen-tranh/${comic.slug}`"
              class="rank-row animate-slide-left"
              :style="{ animationDelay: `${(i+3)*0.06}s` }"
            >
              <span class="rank-num rank-other">{{ i + 4 }}</span>
              <img :src="comic.coverUrl || '/images/no-cover.jpg'" :alt="comic.title" class="rank-thumb" loading="lazy" />
              <div class="rank-info">
                <p class="rank-title">{{ comic.title }}</p>
                <p class="rank-sub">
                  <span class="rank-views">👁 {{ formatViews(comic.views) }}</span>
                  <span class="rank-chapter">Ch.{{ comic.latestChapter }}</span>
                </p>
              </div>
              <div class="rank-trend">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="var(--accent-primary)" stroke-width="2.5"><path d="m18 15-6-6-6 6"/></svg>
              </div>
            </NuxtLink>
          </template>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{ comics: any[]; loading?: boolean }>()

const formatViews = (v?: number) => {
  if (!v) return '0'
  if (v >= 1_000_000) return `${(v/1_000_000).toFixed(1)}M`
  if (v >= 1_000) return `${(v/1_000).toFixed(0)}K`
  return v.toString()
}
</script>

<style scoped>
.top10-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 32px; }
@media (max-width: 768px) { .top10-layout { grid-template-columns: 1fr; } }

.top3-col { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; align-content: start; }
.top3-card-wrap { position: relative; }

/* Rank list */
.rank-list { display: flex; flex-direction: column; gap: 4px; justify-content: center; }
.rank-row { display: flex; align-items: center; gap: 12px; padding: 10px 12px; border-radius: var(--radius-sm); background: var(--bg-secondary); border: 1px solid var(--border-card); transition: var(--transition-base); }
.rank-row:hover { background: var(--bg-tertiary); border-color: var(--border-subtle); transform: translateX(4px); }
.rank-num { width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-family: 'Montserrat', sans-serif; font-weight: 900; font-size: 0.78rem; flex-shrink: 0; background: var(--bg-tertiary); border: 1px solid var(--border-card); color: var(--text-muted); }
.rank-thumb { width: 38px; height: 52px; object-fit: cover; border-radius: var(--radius-xs); flex-shrink: 0; }
.rank-info { flex: 1; min-width: 0; }
.rank-title { font-size: 0.85rem; font-weight: 600; color: var(--text-primary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rank-sub { display: flex; gap: 10px; margin-top: 3px; }
.rank-views,.rank-chapter { font-size: 0.72rem; color: var(--text-muted); }
.rank-trend { flex-shrink: 0; opacity: 0.7; }

/* Skeletons */
.rank-row-skeleton { display: flex; align-items: center; gap: 12px; padding: 10px 12px; border-radius: var(--radius-sm); }
.rank-thumb-sk { width: 38px; height: 52px; border-radius: var(--radius-xs); flex-shrink: 0; }
.rank-info-sk { flex: 1; display: flex; flex-direction: column; gap: 6px; }
.rank-title-sk { height: 13px; width: 80%; border-radius: 4px; }
.rank-sub-sk { height: 10px; width: 50%; border-radius: 4px; }
.top3-skeleton { aspect-ratio: 2/3; border-radius: var(--radius-md); }
</style>
