<template>
  <div class="admin-layout" v-if="isAuthenticated">
    <!-- Sidebar -->
    <aside class="admin-sidebar" :class="{ 'collapsed': sidebarCollapsed }">
      <div class="sidebar-header">
        <NuxtLink to="/admin" class="admin-logo">
          <span class="logo-icon">🍵</span>
          <span v-if="!sidebarCollapsed" class="logo-text">Admin<span class="logo-accent">Panel</span></span>
        </NuxtLink>
        <button class="toggle-btn" @click="sidebarCollapsed = !sidebarCollapsed">
          <svg v-if="!sidebarCollapsed" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
          <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </button>
      </div>

      <nav class="sidebar-nav">
        <NuxtLink to="/admin" class="nav-item" active-class="active" exact>
          <span class="nav-icon">📊</span>
          <span v-if="!sidebarCollapsed" class="nav-label">Dashboard</span>
        </NuxtLink>
        <NuxtLink to="/admin/comics" class="nav-item" active-class="active">
          <span class="nav-icon">🎨</span>
          <span v-if="!sidebarCollapsed" class="nav-label">Truyện Tranh</span>
        </NuxtLink>
        <NuxtLink to="/admin/novels" class="nav-item" active-class="active">
          <span class="nav-icon">📖</span>
          <span v-if="!sidebarCollapsed" class="nav-label">Truyện Chữ</span>
        </NuxtLink>
        <NuxtLink to="/admin/genres" class="nav-item" active-class="active">
          <span class="nav-icon">🏷️</span>
          <span v-if="!sidebarCollapsed" class="nav-label">Thể loại</span>
        </NuxtLink>
        <NuxtLink to="/admin/categories" class="nav-item" active-class="active">
          <span class="nav-icon">🗂️</span>
          <span v-if="!sidebarCollapsed" class="nav-label">Danh mục Menu</span>
        </NuxtLink>
        <NuxtLink to="/admin/users" class="nav-item" active-class="active">
          <span class="nav-icon">👥</span>
          <span v-if="!sidebarCollapsed" class="nav-label">Người dùng</span>
        </NuxtLink>
      </nav>

      <div class="sidebar-footer">
        <!-- Storage Usage Widget -->
        <div v-if="!sidebarCollapsed && storageStats" class="storage-widget" :class="{ 'warning': isNearLimit }">
          <div class="storage-info">
            <span class="storage-label">Lưu trữ (R2 Free)</span>
            <span class="storage-value">{{ formatSize(storageStats.totalSize) }} / 10GB</span>
          </div>
          <div class="storage-bar">
            <div class="storage-progress" :style="{ width: storagePercent + '%' }"></div>
          </div>
          <p v-if="isNearLimit" class="storage-warning">⚠️ Sắp chạm hạn mức 10GB!</p>
        </div>

        <NuxtLink to="/" class="nav-item back-site">
          <span class="nav-icon">🏠</span>
          <span v-if="!sidebarCollapsed" class="nav-label">Về trang chủ</span>
        </NuxtLink>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="admin-main">
      <header class="admin-header">
        <div class="header-left">
          <h2 class="page-title">{{ route.meta.title || 'Dashboard' }}</h2>
        </div>
        <div class="header-right">
          <div class="user-profile">
            <div class="avatar">{{ authStore.user?.username?.[0].toUpperCase() || 'A' }}</div>
            <div class="user-info">
              <span class="user-name">{{ authStore.user?.username || 'Super Admin' }}</span>
              <button class="logout-btn" @click="logout">Đăng xuất</button>
            </div>
          </div>
        </div>
      </header>
      
      <div class="admin-content-scroll">
        <div class="admin-content">
          <slot />
        </div>
      </div>
    </div>
  </div>
  
  <!-- Fallback for login page or loading -->
  <div v-else class="admin-unauth">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const sidebarCollapsed = ref(false)

const isAuthenticated = computed(() => authStore.isLoggedIn)

const storageStats = ref<any>(null)
const { get } = useApi()

const fetchStorageStats = async () => {
  try {
    const data = await get<any>('/admin/storage-stats', {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    storageStats.value = data
  } catch (err) {
    console.error('Error fetching storage stats:', err)
  }
}

onMounted(async () => {
  // Ensure store is initialized from localStorage
  if (!authStore.initialized) {
    authStore.init()
  }

  if (route.path !== '/admin/login') {
    if (!authStore.isLoggedIn) {
      router.push('/admin/login')
    } else {
      fetchStorageStats()
    }
  }
})

const storagePercent = computed(() => {
  if (!storageStats.value) return 0
  return Math.min(100, (storageStats.value.totalSize / storageStats.value.limitBytes) * 100)
})

const isNearLimit = computed(() => {
  if (!storageStats.value) return false
  // 9GB limit = 9 * 1024 * 1024 * 1024
  return storageStats.value.totalSize > (9 * 1024 * 1024 * 1024)
})

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const logout = () => {
  authStore.logout()
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background-color: #0d0f14;
  color: #fff;
  font-family: 'Inter', sans-serif;
}

/* Sidebar */
.admin-sidebar {
  width: 260px;
  background: #141720;
  border-right: 1px solid rgba(255,255,255,0.05);
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
}

.admin-sidebar.collapsed { width: 80px; }

.sidebar-header {
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  border-bottom: 1px solid rgba(255,255,255,0.05);
}

.admin-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-family: 'Montserrat', sans-serif;
  font-weight: 800;
  font-size: 1.2rem;
  color: #F1E8C7;
  text-decoration: none;
}

.logo-accent { color: #9CA764; }

.toggle-btn {
  background: none;
  border: none;
  color: #A8A8B3;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
}
.toggle-btn:hover { background: rgba(255,255,255,0.05); color: #fff; }

.sidebar-nav {
  flex: 1;
  padding: 20px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  color: #A8A8B3;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
}

.nav-item:hover {
  background: rgba(156,167,100,0.1);
  color: #F1E8C7;
}

.nav-item.active {
  background: linear-gradient(135deg, rgba(156,167,100,0.2), rgba(156,167,100,0.05));
  color: #9CA764;
  border-left: 3px solid #9CA764;
}

.nav-icon { font-size: 1.2rem; display: flex; align-items: center; justify-content: center; width: 24px; }
.nav-label { white-space: nowrap; }

.sidebar-footer { padding: 20px 12px; border-top: 1px solid rgba(255,255,255,0.05); }

/* Storage Widget */
.storage-widget {
  padding: 12px;
  background: rgba(0,0,0,0.2);
  border-radius: 10px;
  margin-bottom: 16px;
  border: 1px solid rgba(255,255,255,0.03);
}
.storage-info {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  margin-bottom: 8px;
}
.storage-label { color: #5C5C6B; }
.storage-value { color: #A8A8B3; font-weight: 600; }
.storage-bar {
  height: 6px;
  background: rgba(255,255,255,0.05);
  border-radius: 3px;
  overflow: hidden;
}
.storage-progress {
  height: 100%;
  background: #9CA764;
  border-radius: 3px;
  transition: width 0.5s ease;
}
.storage-widget.warning .storage-progress { background: #ef4444; }
.storage-warning {
  color: #ef4444;
  font-size: 0.7rem;
  margin-top: 8px;
  font-weight: 600;
  text-align: center;
}

.back-site {
  color: #5C5C6B;
}
.back-site:hover { color: #A8A8B3; background: rgba(255,255,255,0.05); }

/* Main */
.admin-main { flex: 1; display: flex; flex-direction: column; min-width: 0; }

.admin-header {
  height: 70px;
  background: #141720;
  border-bottom: 1px solid rgba(255,255,255,0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 30px;
  flex-shrink: 0;
}

.page-title { margin: 0; font-size: 1.2rem; font-weight: 600; font-family: 'Montserrat', sans-serif; color: #fff; }

.user-profile { display: flex; align-items: center; gap: 12px; }
.avatar { width: 40px; height: 40px; border-radius: 50%; background: linear-gradient(135deg, #9CA764, #6B7445); display: flex; align-items: center; justify-content: center; font-weight: bold; color: #000; }
.user-info { display: flex; flex-direction: column; }
.user-name { font-size: 0.9rem; font-weight: 600; }
.logout-btn { background: none; border: none; color: #ef4444; font-size: 0.8rem; cursor: pointer; padding: 0; text-align: left; }
.logout-btn:hover { text-decoration: underline; }

.admin-content-scroll { flex: 1; overflow-y: auto; background: #0d0f14; }
.admin-content { padding: 30px; max-width: 1400px; margin: 0 auto; min-height: 100%; }

.admin-unauth {
  min-height: 100vh;
  background: #0d0f14;
}
</style>
