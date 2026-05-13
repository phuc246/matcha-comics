<template>
  <header class="app-header-wrapper" :class="{ 'header-scrolled': isScrolled }">
    <!-- Main Header -->
    <div class="app-header">
      <div class="container header-inner">
        <!-- Logo -->
        <NuxtLink to="/" class="header-logo">
          <div class="logo-icon-wrap">
            <span class="logo-icon">🍵</span>
          </div>
          <span class="logo-text">Matcha<span class="logo-accent">Comic</span></span>
        </NuxtLink>

        <!-- Desktop Nav -->
        <nav class="header-nav" v-show="!isMobile">
          <div class="nav-item" v-for="menu in menus" :key="menu.label">
            <NuxtLink :to="menu.path" class="nav-link" :class="{ active: isActive(menu.path) }">
              {{ menu.label }}
              <span v-if="menu.children" class="nav-arrow">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 9 6 6 6-6"/></svg>
              </span>
            </NuxtLink>
            <!-- Dropdown -->
            <div v-if="menu.children" class="nav-dropdown">
              <div class="dropdown-inner">
                <div class="dropdown-section" v-for="section in menu.children" :key="section.title">
                  <p class="dropdown-section-title">{{ section.title }}</p>
                  <div class="dropdown-links">
                    <NuxtLink
                      v-for="link in section.links"
                      :key="link.path"
                      :to="link.path"
                      class="dropdown-link"
                    >
                      <span class="dropdown-link-icon">{{ link.icon }}</span>
                      <span class="dropdown-link-text">{{ link.label }}</span>
                    </NuxtLink>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </nav>

        <!-- Right actions -->
        <div class="header-actions">
          <!-- Search -->
          <button class="action-btn search-btn" @click="toggleSearch" aria-label="Tìm kiếm">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
          </button>

          <!-- Admin link -->
          <NuxtLink to="/admin" class="admin-btn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2a3 3 0 0 0-3 3v1H6a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-3V5a3 3 0 0 0-3-3z"/>
            </svg>
            Admin
          </NuxtLink>

          <!-- Mobile menu toggle -->
          <button class="action-btn mobile-menu-btn" @click="toggleMobileMenu" v-if="isMobile" aria-label="Menu">
            <span class="hamburger" :class="{ open: mobileMenuOpen }">
              <span /><span /><span />
            </span>
          </button>
        </div>
      </div>

      <!-- Golden decorative line -->
      <div class="header-gold-line"></div>
    </div>

    <!-- Search overlay -->
    <Transition name="search-slide">
      <div v-if="searchOpen" class="search-overlay">
        <div class="container">
          <div class="search-wrapper">
            <svg class="search-icon-input" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
            </svg>
            <input
              ref="searchInput"
              v-model="searchQuery"
              type="text"
              placeholder="Nhập tên truyện, tác giả..."
              class="search-input"
              @keydown.enter="doSearch"
              @keydown.esc="toggleSearch"
            />
            <button class="search-close" @click="toggleSearch">✕</button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Mobile Drawer -->
    <Transition name="slide-down">
      <div v-if="mobileMenuOpen && isMobile" class="mobile-menu">
        <div v-for="menu in menus" :key="menu.label" class="mobile-menu-group">
          <p class="mobile-menu-title">{{ menu.label }}</p>
          <template v-if="menu.children">
            <template v-for="section in menu.children" :key="section.title">
              <NuxtLink
                v-for="link in section.links"
                :key="link.path"
                :to="link.path"
                class="mobile-menu-link"
                @click="mobileMenuOpen = false"
              >
                <span class="mobile-link-icon">{{ link.icon }}</span> {{ link.label }}
              </NuxtLink>
            </template>
          </template>
          <NuxtLink v-else :to="menu.path" class="mobile-menu-link" @click="mobileMenuOpen = false">
            {{ menu.label }}
          </NuxtLink>
        </div>
      </div>
    </Transition>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isScrolled = ref(false)
const searchOpen = ref(false)
const searchQuery = ref('')
const mobileMenuOpen = ref(false)
const searchInput = ref<HTMLInputElement>()
const isMobile = ref(false)

const menus = [
  {
    label: 'Truyện Tranh',
    path: '/truyen-tranh',
    children: [
      {
        title: 'Khám phá',
        links: [
          { label: 'Tất cả truyện', path: '/truyen-tranh', icon: '📚' },
          { label: 'Mới cập nhật', path: '/truyen-tranh?sort=latest', icon: '🆕' },
          { label: 'Top thịnh hành', path: '/truyen-tranh?sort=views', icon: '🔥' },
          { label: 'Đã hoàn thành', path: '/truyen-tranh?status=completed', icon: '✨' },
        ],
      },
      {
        title: 'Thể loại nổi bật',
        links: [
          { label: 'Hành động', path: '/the-loai/action', icon: '⚔️' },
          { label: 'Lãng mạn', path: '/the-loai/romance', icon: '💕' },
          { label: 'Kỳ ảo', path: '/the-loai/fantasy', icon: '🧙' },
          { label: 'Kinh dị', path: '/the-loai/horror', icon: '👻' },
          { label: 'Tất cả thể loại', path: '/the-loai', icon: '📋' },
        ],
      },
    ],
  },
  {
    label: 'Truyện Chữ',
    path: '/truyen-chu',
    children: [
      {
        title: 'Khám phá',
        links: [
          { label: 'Tất cả tiểu thuyết', path: '/truyen-chu', icon: '📖' },
          { label: 'Mới cập nhật', path: '/truyen-chu?sort=latest', icon: '🆕' },
          { label: 'Top thịnh hành', path: '/truyen-chu?sort=views', icon: '🔥' },
          { label: 'Đã hoàn thành', path: '/truyen-chu?status=completed', icon: '✨' },
        ],
      },
      {
        title: 'Thể loại nổi bật',
        links: [
          { label: 'Tiên hiệp', path: '/the-loai/tien-hiep', icon: '⛅' },
          { label: 'Kiếm hiệp', path: '/the-loai/kiem-hiep', icon: '🗡️' },
          { label: 'Ngôn tình', path: '/the-loai/ngon-tinh', icon: '💞' },
          { label: 'Tất cả thể loại', path: '/the-loai', icon: '📋' },
        ],
      },
    ],
  },
  { label: 'Thể Loại', path: '/the-loai' },
]

const isActive = (path: string) => route.path.startsWith(path)

const toggleSearch = async () => {
  searchOpen.value = !searchOpen.value
  if (searchOpen.value) {
    await nextTick()
    searchInput.value?.focus()
  } else {
    searchQuery.value = ''
  }
}

const doSearch = () => {
  if (searchQuery.value.trim()) {
    router.push(`/tim-kiem?q=${encodeURIComponent(searchQuery.value.trim())}`)
    toggleSearch()
  }
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const handleScroll = () => {
  isScrolled.value = window.scrollY > 10
}

const handleResize = () => {
  isMobile.value = window.innerWidth < 1024
  if (!isMobile.value) mobileMenuOpen.value = false
}

onMounted(() => {
  handleResize()
  window.addEventListener('scroll', handleScroll, { passive: true })
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.app-header-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: var(--z-sticky);
  transition: transform 0.3s ease;
}

/* Main Header */
.app-header {
  height: 72px;
  background: rgba(13, 15, 20, 0.75);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  position: relative;
  transition: var(--transition-base);
}

.header-scrolled .app-header {
  background: rgba(13, 15, 20, 0.95);
  box-shadow: 0 10px 40px rgba(0,0,0,0.6);
}

.header-inner {
  display: flex;
  align-items: center;
  height: 100%;
  gap: 40px;
}

.header-gold-line {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, var(--accent-primary) 50%, transparent 100%);
  opacity: 0.3;
}
.header-scrolled .header-gold-line {
  opacity: 0.8;
  height: 2px;
  background: linear-gradient(90deg, transparent 0%, var(--accent-gold) 50%, transparent 100%);
  box-shadow: 0 0 10px var(--accent-primary);
}

/* Logo */
.header-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.logo-icon-wrap {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--bg-tertiary), var(--bg-primary));
  border: 1px solid var(--border-active);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 15px rgba(156,167,100,0.2);
  transform: rotate(-5deg);
  transition: var(--transition-base);
}

.header-logo:hover .logo-icon-wrap {
  transform: rotate(0deg) scale(1.05);
  box-shadow: 0 0 25px rgba(156,167,100,0.4);
  border-color: var(--accent-gold);
}

.logo-icon {
  font-size: 1.2rem;
  line-height: 1;
}

.logo-text { 
  font-family: 'Montserrat', sans-serif;
  font-weight: 900;
  font-size: 1.4rem;
  letter-spacing: -0.02em;
  color: var(--text-primary);
}

.logo-accent { 
  color: transparent;
  background: linear-gradient(135deg, var(--accent-gold), var(--accent-primary));
  -webkit-background-clip: text;
  background-clip: text;
}

/* Nav */
.header-nav {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
}

.nav-item {
  position: relative;
  height: 72px; /* Full height to make hover easier */
  display: flex;
  align-items: center;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-secondary);
  border-radius: 100px;
  transition: var(--transition-base);
  white-space: nowrap;
}

.nav-link:hover,
.nav-link.active {
  color: var(--text-primary);
  background: rgba(156, 167, 100, 0.1);
}

.nav-link.active {
  color: var(--accent-gold);
}

.nav-arrow {
  display: flex;
  align-items: center;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--text-muted);
}

.nav-item:hover .nav-arrow {
  transform: rotate(180deg);
  color: var(--accent-primary);
}

/* Dropdown Mega Menu */
.nav-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 480px;
  background: rgba(20, 23, 32, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--border-card);
  border-top: 2px solid var(--accent-primary);
  border-radius: 0 0 var(--radius-lg) var(--radius-lg);
  box-shadow: 0 20px 40px rgba(0,0,0,0.8);
  opacity: 0;
  visibility: hidden;
  transform: translateY(10px);
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  pointer-events: none;
  z-index: 200;
}

.nav-item:hover .nav-dropdown {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
  pointer-events: all;
}

.dropdown-inner {
  display: flex;
  padding: 24px;
  gap: 32px;
}

.dropdown-section {
  flex: 1;
}

.dropdown-section-title {
  font-family: 'Montserrat', sans-serif;
  font-size: 0.75rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-muted);
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255,255,255,0.05);
}

.dropdown-links {
  display: grid;
  grid-template-columns: 1fr;
  gap: 4px;
}

.dropdown-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  transition: var(--transition-fast);
}

.dropdown-link:hover {
  background: var(--bg-tertiary);
  color: var(--accent-gold);
  transform: translateX(4px);
}

.dropdown-link-icon { 
  font-size: 1rem;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255,255,255,0.03);
  border-radius: 6px;
}
.dropdown-link-text { font-size: 0.9rem; font-weight: 500; }

/* Header actions */
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.action-btn {
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255,255,255,0.03);
  border: 1px solid var(--border-card);
  border-radius: 50%;
  color: var(--text-primary);
  cursor: pointer;
  transition: var(--transition-base);
}

.action-btn:hover {
  background: var(--bg-tertiary);
  border-color: var(--accent-primary);
  color: var(--accent-gold);
  box-shadow: 0 0 15px rgba(156,167,100,0.3);
  transform: translateY(-2px);
}

.admin-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-dark));
  color: var(--bg-primary);
  font-family: 'Montserrat', sans-serif;
  font-weight: 800;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-radius: 100px;
  transition: var(--transition-base);
  box-shadow: 0 4px 15px rgba(156,167,100,0.2);
}

.admin-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(156,167,100,0.4);
  background: linear-gradient(135deg, var(--accent-hover), var(--accent-primary));
}

/* Search overlay */
.search-overlay {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: rgba(13, 15, 20, 0.98);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--accent-primary);
  padding: 30px 0;
  box-shadow: 0 20px 40px rgba(0,0,0,0.8);
}

.search-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--accent-primary);
  border-radius: 100px;
  padding: 14px 24px;
  box-shadow: 0 0 30px rgba(156,167,100,0.15);
}

.search-icon-input { color: var(--accent-gold); flex-shrink: 0; }

.search-input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-size: 1.1rem;
  font-family: 'Inter', sans-serif;
}

.search-input::placeholder { color: var(--text-muted); }

.search-close {
  background: rgba(255,255,255,0.05);
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 1.2rem;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-fast);
}
.search-close:hover { 
  background: var(--bg-tertiary);
  color: var(--text-primary); 
}

/* Search transitions */
.search-slide-enter-active,
.search-slide-leave-active { transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.search-slide-enter-from,
.search-slide-leave-to { opacity: 0; transform: translateY(-20px); }

/* Hamburger */
.hamburger {
  display: flex;
  flex-direction: column;
  gap: 5px;
  width: 20px;
}
.hamburger span {
  display: block;
  height: 2px;
  background: currentColor;
  border-radius: 2px;
  transition: var(--transition-base);
}
.hamburger.open span:nth-child(1) { transform: translateY(7px) rotate(45deg); }
.hamburger.open span:nth-child(2) { opacity: 0; }
.hamburger.open span:nth-child(3) { transform: translateY(-7px) rotate(-45deg); }

/* Mobile menu */
.mobile-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: rgba(13, 15, 20, 0.98);
  backdrop-filter: blur(20px);
  border-bottom: 2px solid var(--accent-primary);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-height: calc(100vh - 72px);
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0,0,0,0.8);
}

.mobile-menu-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mobile-menu-title {
  font-family: 'Montserrat', sans-serif;
  font-size: 0.8rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--accent-gold);
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255,255,255,0.05);
  margin-bottom: 4px;
}

.mobile-menu-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  font-size: 0.95rem;
  font-weight: 500;
  transition: var(--transition-fast);
  background: rgba(255,255,255,0.02);
}
.mobile-menu-link:hover {
  background: rgba(156, 167, 100, 0.1);
  color: var(--text-primary);
  transform: translateX(4px);
}

.mobile-link-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255,255,255,0.03);
  border-radius: 6px;
  font-size: 1.1rem;
}

/* Slide transitions */
.slide-down-enter-active,
.slide-down-leave-active { transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1); }
.slide-down-enter-from,
.slide-down-leave-to { opacity: 0; transform: translateY(-20px); }

@media (max-width: 1024px) {
  .app-header-wrapper { transform: translateY(0); }
  .top-bar { display: none; }
  .app-header { height: 64px; }
  .header-inner { gap: 20px; }
}
</style>
