<template>
  <section class="page-section">
    <div class="container">
      <div class="section-header">
        <h2 class="section-title">✅ Đã Hoàn Thành</h2>
        <NuxtLink to="/truyen-tranh?status=completed" class="section-view-all">Xem tất cả →</NuxtLink>
      </div>
      <div class="completed-scroll">
        <div class="completed-track">
          <NuxtLink
            v-for="(comic, i) in comics"
            :key="comic.id"
            :to="`/truyen-tranh/${comic.slug}`"
            class="completed-card animate-fade-up"
            :style="{ animationDelay: `${i * 0.05}s` }"
          >
            <div class="completed-cover">
              <img :src="comic.coverUrl || '/images/no-cover.jpg'" :alt="comic.title" loading="lazy" />
              <span class="completed-badge">FULL</span>
            </div>
            <p class="completed-title">{{ comic.title }}</p>
            <p class="completed-meta">{{ comic.chapterCount }} chương</p>
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{ comics: any[]; loading?: boolean }>()
</script>

<style scoped>
.completed-scroll { overflow-x: auto; padding-bottom: 8px; }
.completed-scroll::-webkit-scrollbar { height: 4px; }
.completed-track { display: flex; gap: 16px; width: max-content; padding: 4px 0; }
.completed-card { display: flex; flex-direction: column; gap: 6px; width: 130px; flex-shrink: 0; }
.completed-cover { position: relative; aspect-ratio: 2/3; border-radius: var(--radius-sm); overflow: hidden; border: 1px solid var(--border-card); transition: var(--transition-base); }
.completed-cover:hover { transform: translateY(-4px); box-shadow: var(--shadow-glow); }
.completed-cover img { width: 100%; height: 100%; object-fit: cover; }
.completed-badge { position: absolute; bottom: 6px; right: 6px; background: rgba(91,141,239,0.9); color: #fff; font-size: 0.62rem; font-weight: 700; padding: 2px 7px; border-radius: 10px; }
.completed-title { font-size: 0.78rem; font-weight: 600; color: var(--text-primary); line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.completed-meta { font-size: 0.68rem; color: var(--text-muted); }
</style>
