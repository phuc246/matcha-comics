<template>
  <footer class="app-footer">
    <div class="glow-line" />
    <div class="container footer-inner">
      <div class="footer-brand">
        <NuxtLink to="/" class="footer-logo">
          <span class="logo-icon">🍵</span>
          <span>Matcha<span class="logo-accent">Comic</span></span>
        </NuxtLink>
        <p class="footer-desc">Đọc truyện tranh &amp; tiểu thuyết online miễn phí. Cập nhật nhanh nhất, chất lượng cao nhất.</p>
        <div class="footer-stats">
          <div class="stat-item"><span class="stat-num">5,000+</span><span class="stat-label">Bộ truyện</span></div>
          <div class="stat-divider" />
          <div class="stat-item"><span class="stat-num">50K+</span><span class="stat-label">Chương</span></div>
          <div class="stat-divider" />
          <div class="stat-item"><span class="stat-num">Free</span><span class="stat-label">Hoàn toàn</span></div>
        </div>
      </div>

      <div class="footer-links-group" v-for="group in linkGroups" :key="group.title">
        <h4 class="footer-group-title">{{ group.title }}</h4>
        <div class="footer-links">
          <NuxtLink v-for="link in group.links" :key="link.path" :to="link.path" class="footer-link">
            {{ link.label }}
          </NuxtLink>
        </div>
      </div>
    </div>

    <div class="footer-bottom">
      <div class="container footer-bottom-inner">
        <p class="footer-copy">© 2025 MatchaComic. Tất cả quyền được bảo lưu.</p>
        <div class="footer-tags">
          <span class="footer-tag">Truyện tranh</span>
          <span class="footer-tag">Tiểu thuyết</span>
          <span class="footer-tag">Manga</span>
          <span class="footer-tag">Manhwa</span>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const linkGroups = [
  {
    title: 'Truyện Tranh',
    links: [
      { label: 'Tất cả truyện', path: '/truyen-tranh' },
      { label: 'Mới cập nhật', path: '/truyen-tranh?sort=latest' },
      { label: 'Top xem nhiều', path: '/truyen-tranh?sort=views' },
      { label: 'Hoàn thành', path: '/truyen-tranh?status=completed' },
    ],
  },
  {
    title: 'Truyện Chữ',
    links: [
      { label: 'Tất cả tiểu thuyết', path: '/truyen-chu' },
      { label: 'Mới cập nhật', path: '/truyen-chu?sort=latest' },
      { label: 'Top xem nhiều', path: '/truyen-chu?sort=views' },
      { label: 'Hoàn thành', path: '/truyen-chu?status=completed' },
    ],
  },
  {
    title: 'Thể Loại',
    links: [
      { label: 'Action', path: '/the-loai/action' },
      { label: 'Romance', path: '/the-loai/romance' },
      { label: 'Fantasy', path: '/the-loai/fantasy' },
      { label: 'Tiên hiệp', path: '/the-loai/tien-hiep' },
      { label: 'Ngôn tình', path: '/the-loai/ngon-tinh' },
    ],
  },
]
</script>

<style scoped>
.app-footer { background: var(--bg-secondary); margin-top: 80px; }
.footer-inner { display: grid; grid-template-columns: 1.5fr 1fr 1fr 1fr; gap: 48px; padding: 56px 0 40px; }
@media (max-width: 1024px) { .footer-inner { grid-template-columns: 1fr 1fr; } }
@media (max-width: 640px) { .footer-inner { grid-template-columns: 1fr; gap: 32px; } }
.footer-brand { display: flex; flex-direction: column; gap: 16px; }
.footer-logo { display: flex; align-items: center; gap: 8px; font-family: 'Montserrat', sans-serif; font-weight: 900; font-size: 1.3rem; color: var(--text-primary); }
.logo-icon { font-size: 1.4rem; filter: drop-shadow(0 0 6px rgba(156,167,100,0.7)); }
.logo-accent { color: var(--accent-primary); }
.footer-desc { font-size: 0.875rem; color: var(--text-muted); line-height: 1.6; max-width: 280px; }
.footer-stats { display: flex; align-items: center; gap: 16px; }
.stat-item { display: flex; flex-direction: column; }
.stat-num { font-family: 'Montserrat', sans-serif; font-weight: 800; font-size: 1.1rem; color: var(--accent-primary); }
.stat-label { font-size: 0.72rem; color: var(--text-muted); }
.stat-divider { width: 1px; height: 32px; background: var(--border-subtle); }
.footer-group-title { font-family: 'Montserrat', sans-serif; font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.08em; color: var(--accent-primary); margin-bottom: 16px; }
.footer-links { display: flex; flex-direction: column; gap: 8px; }
.footer-link { font-size: 0.875rem; color: var(--text-muted); transition: var(--transition-fast); }
.footer-link:hover { color: var(--text-primary); transform: translateX(4px); display: block; }
.footer-bottom { border-top: 1px solid var(--border-subtle); padding: 20px 0; }
.footer-bottom-inner { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; }
.footer-copy { font-size: 0.82rem; color: var(--text-muted); }
.footer-tags { display: flex; gap: 8px; flex-wrap: wrap; }
.footer-tag { font-size: 0.72rem; padding: 3px 10px; border-radius: 20px; background: var(--bg-tertiary); border: 1px solid var(--border-card); color: var(--text-muted); }
</style>
