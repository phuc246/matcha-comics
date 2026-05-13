<template>
  <div class="admin-users-page">
    <div class="page-header-actions">
      <div class="search-filter">
        <input type="text" v-model="search" placeholder="Tìm kiếm người dùng..." class="search-input" />
      </div>
      <button class="btn-primary">
        <span class="icon">+</span> Thêm Nhân Viên
      </button>
    </div>

    <div class="table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th width="60">#</th>
            <th>Tên người dùng</th>
            <th>Email</th>
            <th>Vai trò</th>
            <th>Trạng thái</th>
            <th>Ngày tạo</th>
            <th width="100">Thao tác</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(u, i) in filteredUsers" :key="u.id">
            <td>{{ i + 1 }}</td>
            <td>
              <div class="user-cell">
                <div class="user-avatar" :style="{ background: u.avatarColor }">{{ u.username[0].toUpperCase() }}</div>
                <span class="username">{{ u.username }}</span>
              </div>
            </td>
            <td>{{ u.email }}</td>
            <td>
              <span class="role-badge" :class="u.role">{{ u.role === 'admin' ? 'Quản trị viên' : 'Biên tập viên' }}</span>
            </td>
            <td>
              <span class="status-indicator" :class="u.status"></span>
              {{ u.status === 'active' ? 'Hoạt động' : 'Bị khóa' }}
            </td>
            <td>{{ u.createdAt }}</td>
            <td>
              <div class="actions-cell">
                <button class="icon-btn edit"><span class="icon">✏️</span></button>
                <button class="icon-btn delete" v-if="u.role !== 'admin'"><span class="icon">🔒</span></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'admin',
  title: 'Quản lý Người dùng'
})

useHead({ title: 'Quản lý Người dùng - Matcha Comic' })

const search = ref('')
const users = ref([
  { id: 1, username: 'superadmin', email: 'admin@matchacomic.com', role: 'admin', status: 'active', createdAt: '2026-01-01', avatarColor: '#9CA764' },
  { id: 2, username: 'editor_manga', email: 'editor1@matchacomic.com', role: 'editor', status: 'active', createdAt: '2026-02-15', avatarColor: '#F1E8C7' },
  { id: 3, username: 'novel_writer', email: 'writer@matchacomic.com', role: 'editor', status: 'active', createdAt: '2026-03-10', avatarColor: '#6B7445' },
  { id: 4, username: 'blocked_user', email: 'spam@gmail.com', role: 'editor', status: 'blocked', createdAt: '2026-04-01', avatarColor: '#5C5C6B' },
])

const filteredUsers = computed(() => {
  return users.value.filter(u => u.username.toLowerCase().includes(search.value.toLowerCase()))
})
</script>

<style scoped>
.admin-users-page { display: flex; flex-direction: column; gap: 24px; }
.page-header-actions { display: flex; justify-content: space-between; align-items: center; gap: 20px; }
.search-filter { flex: 1; max-width: 400px; }
.search-input { width: 100%; background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 8px; padding: 10px 16px; color: #fff; font-size: 0.9rem; }

.btn-primary { background: #9CA764; color: #000; padding: 10px 20px; border-radius: 8px; font-weight: 700; border: none; cursor: pointer; display: flex; align-items: center; gap: 8px; }

.table-container { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 12px; overflow: hidden; }
.admin-table { width: 100%; border-collapse: collapse; text-align: left; }
.admin-table th { background: rgba(0,0,0,0.2); padding: 16px; font-size: 0.85rem; color: #5C5C6B; font-weight: 600; }
.admin-table td { padding: 16px; border-bottom: 1px solid rgba(255,255,255,0.03); color: #A8A8B3; font-size: 0.9rem; }

.user-cell { display: flex; align-items: center; gap: 12px; }
.user-avatar { width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #000; font-weight: 800; font-size: 0.8rem; }
.username { color: #fff; font-weight: 600; }

.role-badge { padding: 4px 10px; border-radius: 100px; font-size: 0.75rem; font-weight: 700; }
.role-badge.admin { background: rgba(156,167,100,0.1); color: #9CA764; }
.role-badge.editor { background: rgba(241,232,199,0.1); color: #F1E8C7; }

.status-indicator { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 6px; }
.status-indicator.active { background: #4ade80; box-shadow: 0 0 8px #4ade80; }
.status-indicator.blocked { background: #ef4444; box-shadow: 0 0 8px #ef4444; }

.actions-cell { display: flex; gap: 8px; }
.icon-btn { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.05); border-radius: 6px; cursor: pointer; color: inherit; }
.icon-btn:hover { background: rgba(255,255,255,0.05); }
</style>
