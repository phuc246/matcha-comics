<template>
  <div class="novel-detail-page">
    <div v-if="loading" class="page-loader">
      <div class="loader-spinner"></div>
    </div>
    
    <template v-else-if="novel">
      <!-- Backdrop Banner -->
      <div class="detail-backdrop">
        <img :src="novel.coverUrl" :alt="novel.title" class="backdrop-img" />
        <div class="backdrop-overlay"></div>
      </div>

      <div class="container detail-container">
        <!-- Breadcrumb -->
        <div class="breadcrumb">
          <NuxtLink to="/">Trang chủ</NuxtLink>
          <span class="separator">/</span>
          <NuxtLink to="/truyen-chu">Truyện chữ</NuxtLink>
          <span class="separator">/</span>
          <span class="current">{{ novel.title }}</span>
        </div>

        <!-- Header Info -->
        <div class="detail-header">
          <div class="detail-cover-wrap">
            <div class="novel-cover-3d">
              <img :src="novel.coverUrl" :alt="novel.title" class="detail-cover" />
              <div class="rank-badge" v-if="novel.isHot">🔥 HOT</div>
            </div>
          </div>
          
          <div class="detail-info animate-slide-left">
            <h1 class="novel-title">{{ novel.title }}</h1>
            
            <div class="novel-meta">
              <div class="meta-item">
                <span class="meta-label">Tác giả:</span>
                <span class="meta-value">Đang cập nhật</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">Trạng thái:</span>
                <span class="meta-value status-badge" :class="novel.status">{{ novel.status === 'completed' ? 'Hoàn thành' : 'Đang ra' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">Cập nhật:</span>
                <span class="meta-value">{{ novel.updatedAt }}</span>
              </div>
            </div>

            <div class="novel-stats">
              <div class="stat"><span class="icon">👁</span> {{ formatViews(novel.views) }} Lượt xem</div>
              <div class="stat"><span class="icon">📄</span> {{ novel.chapterCount }} Chương</div>
              <div class="stat"><span class="icon">⭐</span> {{ novel.rating }}/10 Đánh giá</div>
            </div>

            <div class="novel-genres">
              <NuxtLink v-for="g in novel.genres" :key="g" :to="`/the-loai/${g.toLowerCase()}`" class="genre-tag">
                {{ g }}
              </NuxtLink>
            </div>

            <div class="novel-desc-box">
              <p class="novel-desc" :class="{ 'expanded': descExpanded }">{{ novel.description }}</p>
              <button class="btn-expand" @click="descExpanded = !descExpanded">
                {{ descExpanded ? 'Thu gọn' : 'Xem thêm' }}
              </button>
            </div>

            <div class="detail-actions">
              <NuxtLink :to="`/truyen-chu/${novel.slug}/chapter-1`" class="btn-read-first">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M8 5v14l11-7z"/></svg>
                Đọc Từ Đầu
              </NuxtLink>
              <NuxtLink :to="`/truyen-chu/${novel.slug}/chapter-${novel.latestChapter}`" class="btn-read-latest">
                Chương Mới Nhất
              </NuxtLink>
              <button class="btn-bookmark" aria-label="Lưu truyện">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/></svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Main Content area (Tabs & Chapters) -->
        <div class="detail-main">
          <div class="tabs">
            <button class="tab-btn active">Danh sách chương</button>
            <button class="tab-btn">Bình luận (124)</button>
          </div>

          <div class="tab-content chapters-section">
            <div class="chapters-header">
              <h3>Chương mới nhất</h3>
            </div>
            <div class="chapters-list">
              <NuxtLink
                v-for="i in Math.min(20, novel.latestChapter)"
                :key="i"
                :to="`/truyen-chu/${novel.slug}/chapter-${novel.latestChapter - i + 1}`"
                class="chapter-item"
              >
                <span class="chapter-name">Chương {{ novel.latestChapter - i + 1 }}</span>
                <span class="chapter-time">Mới đây</span>
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </template>
    
    <div v-else class="container empty-state">
      <h2>Truyện không tồn tại hoặc đã bị xóa</h2>
      <NuxtLink to="/" class="btn-home">Về trang chủ</NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { mockNovels } from '~/data/mock'

const route = useRoute()
const loading = ref(true)
const novel = ref<any>(null)
const descExpanded = ref(false)

const formatViews = (v?: number) => {
  if (!v) return '0'
  if (v >= 1_000_000) return `${(v / 1_000_000).toFixed(1)}M`
  if (v >= 1_000) return `${(v / 1_000).toFixed(0)}K`
  return v.toString()
}

onMounted(async () => {
  // Simulate API fetch
  await new Promise(r => setTimeout(r, 400))
  const slug = route.params.slug as string
  novel.value = mockNovels.find(c => c.slug === slug)
  loading.value = false
})

useHead(() => ({
  title: novel.value ? `${novel.value.title} - MatchaComic` : 'Đang tải...',
}))
</script>

<style scoped>
.novel-detail-page {
  position: relative;
  min-height: 100vh;
  padding-bottom: 60px;
}

/* Loader */
.page-loader { height: 60vh; display: flex; align-items: center; justify-content: center; }
.loader-spinner { width: 40px; height: 40px; border: 3px solid var(--border-card); border-top-color: var(--accent-primary); border-radius: 50%; animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Backdrop */
.detail-backdrop { position: absolute; top: 0; left: 0; right: 0; height: 500px; z-index: 0; overflow: hidden; }
.backdrop-img { width: 100%; height: 100%; object-fit: cover; object-position: center 20%; filter: blur(25px) saturate(1.5) brightness(0.6); transform: scale(1.1); }
.backdrop-overlay { position: absolute; inset: 0; background: linear-gradient(to bottom, rgba(13,15,20,0.3) 0%, rgba(13,15,20,0.8) 50%, var(--bg-primary) 100%); }

/* Container */
.detail-container { position: relative; z-index: 10; padding-top: 30px; }

/* Breadcrumb */
.breadcrumb { display: flex; gap: 8px; font-size: 0.85rem; color: var(--text-muted); margin-bottom: 30px; }
.breadcrumb a { transition: var(--transition-fast); }
.breadcrumb a:hover { color: var(--accent-primary); }
.breadcrumb .current { color: var(--text-primary); font-weight: 500; }

/* Header Info */
.detail-header { display: flex; gap: 40px; margin-bottom: 50px; background: rgba(20,23,32,0.45); backdrop-filter: blur(16px); padding: 30px; border-radius: var(--radius-lg); border: 1px solid rgba(255,255,255,0.05); box-shadow: 0 20px 40px rgba(0,0,0,0.5); }
@media (max-width: 768px) { .detail-header { flex-direction: column; padding: 20px; gap: 24px; align-items: center; text-align: center; } }

/* Cover */
.detail-cover-wrap { flex-shrink: 0; }
.novel-cover-3d { position: relative; width: 220px; aspect-ratio: 2/3; border-radius: var(--radius-md); box-shadow: 0 15px 35px rgba(0,0,0,0.6), var(--shadow-glow); border: 2px solid var(--border-active); }
@media (max-width: 768px) { .novel-cover-3d { width: 180px; } }
.detail-cover { width: 100%; height: 100%; object-fit: cover; border-radius: calc(var(--radius-md) - 2px); }
.rank-badge { position: absolute; top: -10px; right: -10px; background: linear-gradient(135deg, #FF4B2B, #FF416C); color: #fff; font-size: 0.75rem; font-weight: 800; padding: 4px 10px; border-radius: 20px; box-shadow: 0 4px 10px rgba(255,75,43,0.4); }

/* Info */
.detail-info { flex: 1; display: flex; flex-direction: column; gap: 16px; }
.novel-title { font-family: 'Montserrat', sans-serif; font-size: clamp(1.8rem, 4vw, 2.8rem); font-weight: 900; color: var(--text-primary); line-height: 1.1; text-shadow: 0 2px 10px rgba(0,0,0,0.5); }

.novel-meta { display: flex; flex-wrap: wrap; gap: 20px; font-size: 0.9rem; }
@media (max-width: 768px) { .novel-meta { justify-content: center; } }
.meta-item { display: flex; gap: 6px; }
.meta-label { color: var(--text-muted); }
.meta-value { color: var(--text-primary); font-weight: 600; }
.status-badge.completed { color: #4ade80; }
.status-badge.ongoing { color: var(--accent-primary); }

.novel-stats { display: flex; flex-wrap: wrap; gap: 20px; padding: 12px 20px; background: rgba(0,0,0,0.2); border-radius: var(--radius-sm); border: 1px solid var(--border-card); }
@media (max-width: 768px) { .novel-stats { justify-content: center; } }
.stat { display: flex; align-items: center; gap: 6px; font-size: 0.9rem; color: var(--text-secondary); font-weight: 500; }
.stat .icon { font-size: 1.1rem; color: var(--accent-gold); }

.novel-genres { display: flex; flex-wrap: wrap; gap: 8px; }
@media (max-width: 768px) { .novel-genres { justify-content: center; } }
.genre-tag { padding: 6px 14px; background: var(--bg-tertiary); border: 1px solid var(--border-subtle); border-radius: 100px; font-size: 0.8rem; font-weight: 600; color: var(--text-secondary); transition: var(--transition-fast); }
.genre-tag:hover { background: rgba(156,167,100,0.15); color: var(--accent-primary); border-color: var(--accent-primary); transform: translateY(-2px); }

.novel-desc-box { margin-top: 4px; }
.novel-desc { font-size: 0.95rem; line-height: 1.7; color: var(--text-secondary); display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
.novel-desc.expanded { -webkit-line-clamp: unset; }
.btn-expand { margin-top: 8px; background: none; border: none; color: var(--accent-primary); font-size: 0.85rem; font-weight: 600; cursor: pointer; padding: 0; transition: color 0.2s; }
.btn-expand:hover { color: var(--accent-hover); text-decoration: underline; }

.detail-actions { display: flex; flex-wrap: wrap; gap: 12px; margin-top: auto; padding-top: 10px; }
@media (max-width: 768px) { .detail-actions { justify-content: center; } }
.btn-read-first { display: flex; align-items: center; gap: 8px; padding: 12px 28px; background: linear-gradient(135deg, var(--accent-primary), var(--accent-hover)); color: var(--bg-primary); font-weight: 800; font-size: 1rem; border-radius: var(--radius-md); transition: var(--transition-base); box-shadow: 0 4px 15px rgba(156,167,100,0.3); }
.btn-read-first:hover { transform: translateY(-3px); box-shadow: 0 8px 25px rgba(156,167,100,0.5); }
.btn-read-latest { padding: 12px 28px; background: var(--bg-tertiary); border: 1px solid var(--border-active); color: var(--text-primary); font-weight: 700; font-size: 1rem; border-radius: var(--radius-md); transition: var(--transition-base); }
.btn-read-latest:hover { background: rgba(156,167,100,0.1); color: var(--accent-primary); transform: translateY(-3px); }
.btn-bookmark { width: 48px; height: 48px; display: flex; align-items: center; justify-content: center; background: var(--bg-tertiary); border: 1px solid var(--border-card); border-radius: var(--radius-md); color: var(--text-secondary); cursor: pointer; transition: var(--transition-base); }
.btn-bookmark:hover { color: #ef4444; border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.05); }

/* Main Section (Chapters) */
.detail-main { background: var(--bg-secondary); border-radius: var(--radius-lg); border: 1px solid var(--border-card); overflow: hidden; }
.tabs { display: flex; border-bottom: 1px solid var(--border-card); background: rgba(0,0,0,0.2); }
.tab-btn { flex: 1; padding: 18px; background: none; border: none; font-size: 1rem; font-weight: 700; color: var(--text-muted); cursor: pointer; transition: var(--transition-base); position: relative; }
.tab-btn:hover { color: var(--text-primary); background: rgba(255,255,255,0.02); }
.tab-btn.active { color: var(--accent-gold); }
.tab-btn.active::after { content: ''; position: absolute; bottom: -1px; left: 0; right: 0; height: 2px; background: var(--accent-gold); box-shadow: 0 0 10px var(--accent-primary); }

.chapters-section { padding: 30px; }
.chapters-header { margin-bottom: 20px; display: flex; justify-content: space-between; align-items: center; }
.chapters-header h3 { font-size: 1.2rem; font-weight: 700; color: var(--text-primary); }

.chapters-list { display: flex; flex-direction: column; gap: 8px; max-height: 400px; overflow-y: auto; padding-right: 8px; }
.chapters-list::-webkit-scrollbar { width: 6px; }
.chapters-list::-webkit-scrollbar-track { background: rgba(0, 0, 0, 0.2); border-radius: 4px; }
.chapters-list::-webkit-scrollbar-thumb { background: var(--border-active); border-radius: 4px; }
.chapters-list::-webkit-scrollbar-thumb:hover { background: var(--accent-primary); }
.chapter-item { display: flex; justify-content: space-between; align-items: center; padding: 14px 16px; background: var(--bg-tertiary); border: 1px solid var(--border-subtle); border-radius: var(--radius-sm); transition: var(--transition-fast); }
.chapter-item:hover { background: rgba(156,167,100,0.08); border-color: var(--accent-primary); transform: translateX(4px); }
.chapter-name { font-size: 0.9rem; font-weight: 600; color: var(--text-secondary); }
.chapter-item:hover .chapter-name { color: var(--accent-primary); }
.chapter-time { font-size: 0.75rem; color: var(--text-muted); }

/* Empty state */
.empty-state { text-align: center; padding: 100px 20px; }
.empty-state h2 { font-size: 1.5rem; color: var(--text-secondary); margin-bottom: 20px; }
.btn-home { display: inline-block; padding: 10px 24px; background: var(--bg-tertiary); border: 1px solid var(--border-card); border-radius: 100px; color: var(--text-primary); transition: var(--transition-base); }
.btn-home:hover { background: var(--accent-primary); color: #000; }
</style>
