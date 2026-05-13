<template>
  <section class="page-section section-latest">
    <div class="container">
      <div class="section-header">
        <h2 class="section-title">{{ icon }} {{ title }}</h2>
        <NuxtLink :to="viewAllLink" class="section-view-all">Xem tất cả →</NuxtLink>
      </div>
      <div class="comic-grid">
        <template v-if="loading">
          <ComicCardSkeleton v-for="i in 12" :key="i" />
        </template>
        <template v-else>
          <ComicCard3D
            v-for="(comic, i) in comics"
            :key="comic.id"
            :comic="comic"
            :show-info="false"
            class="animate-fade-up"
            :style="{ animationDelay: `${i * 0.05}s` }"
          />
        </template>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  comics: any[]
  loading?: boolean
  title?: string
  icon?: string
  viewAllLink?: string
}>(), {
  loading: false,
  title: 'Mới Nhất',
  icon: '🆕',
  viewAllLink: '/truyen-tranh?sort=latest',
})
</script>
