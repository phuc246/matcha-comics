<template>
  <div class="genre-page">
    <div class="container page-header">
      <h1 class="page-title">📋 Tất cả thể loại</h1>
      <p class="page-desc">Khám phá hàng ngàn bộ truyện theo các thể loại khác nhau</p>
    </div>

    <div class="container">
      <div class="genre-grid">
        <NuxtLink
          v-for="(genre, i) in genres"
          :key="genre.slug"
          :to="`/the-loai/${genre.slug}`"
          class="genre-card animate-fade-up"
          :style="{ animationDelay: `${i * 0.04}s` }"
        >
          <span class="genre-icon">{{ genre.icon }}</span>
          <div class="genre-info">
            <span class="genre-name">{{ genre.name }}</span>
            <span class="genre-count">{{ genre.count }} truyện</span>
          </div>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockGenres } from '~/data/mock'

useHead({ title: 'Thể Loại — Matcha Comic' })

const genres = ref(mockGenres)
</script>

<style scoped>
.genre-page { padding: 40px 0 80px; }
.page-header { text-align: center; margin-bottom: 40px; }
.page-title { font-size: clamp(1.8rem, 4vw, 2.4rem); font-weight: 900; color: var(--text-primary); margin-bottom: 8px; }
.page-desc { color: var(--text-muted); font-size: 1rem; }

.genre-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.genre-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-card);
  border-radius: var(--radius-lg);
  transition: var(--transition-base);
}

.genre-card:hover {
  background: var(--bg-tertiary);
  border-color: var(--accent-primary);
  transform: translateY(-4px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.5);
}

.genre-icon {
  font-size: 2.5rem;
  transition: transform 0.3s;
}

.genre-card:hover .genre-icon {
  transform: scale(1.1) rotate(5deg);
}

.genre-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.genre-name {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-primary);
}

.genre-count {
  font-size: 0.8rem;
  color: var(--text-muted);
}
</style>
