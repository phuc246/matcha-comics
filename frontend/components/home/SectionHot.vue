<template>
  <section class="page-section section-hot">
    <div class="container">
      <div class="section-header">
        <h2 class="section-title">
          <span class="fire-icon animate-fire">🔥</span> Đang Hot
        </h2>
        <NuxtLink to="/truyen-tranh?sort=hot" class="section-view-all">Xem tất cả →</NuxtLink>
      </div>
      <div class="hot-grid">
        <template v-if="loading">
          <ComicCardSkeleton v-for="i in 6" :key="i" />
        </template>
        <template v-else>
          <ComicCard3D
            v-for="(comic, i) in comics"
            :key="comic.id"
            :comic="{ ...comic, isHot: true }"
            rarity="hot"
            class="animate-scale-in"
            :style="{ animationDelay: `${i * 0.07}s` }"
          />
        </template>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{ comics: any[]; loading?: boolean }>()
</script>

<style scoped>
.section-hot { position: relative; }
.section-hot::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(ellipse at 50% 0%, rgba(255,77,77,0.04) 0%, transparent 70%);
  pointer-events: none;
}
.fire-icon { display: inline-block; }
.hot-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 20px;
}
@media (min-width: 640px) { .hot-grid { grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); } }
@media (min-width: 1024px) { .hot-grid { grid-template-columns: repeat(6, 1fr); } }
</style>
