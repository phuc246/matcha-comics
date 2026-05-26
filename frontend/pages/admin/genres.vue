<template>
  <div class="admin-genres-page">
    <div class="page-header">
      <h2>Quản lý Thể loại</h2>
      <button class="btn-primary" @click="openModal()">+ Thêm Thể loại</button>
    </div>

    <!-- Danh sách Thể loại -->
    <div class="table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Tên Thể Loại</th>
            <th>Slug</th>
            <th>Màu Badge</th>
            <th>Hành động</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="5" class="text-center">Đang tải...</td>
          </tr>
          <tr v-else-if="genres.length === 0">
            <td colspan="5" class="text-center">Chưa có thể loại nào</td>
          </tr>
          <tr v-for="genre in genres" :key="genre.id" v-else>
            <td>#{{ genre.id }}</td>
            <td>
              <span class="genre-badge" :style="{ backgroundColor: genre.badgeColor || '#333' }">
                {{ genre.name }}
              </span>
            </td>
            <td>{{ genre.slug }}</td>
            <td>
              <div class="color-preview">
                <span class="color-box" :style="{ backgroundColor: genre.badgeColor || '#333' }"></span>
                {{ genre.badgeColor || '#333' }}
              </div>
            </td>
            <td>
              <div class="actions">
                <button class="btn-edit" @click="openModal(genre)">Sửa</button>
                <button class="btn-delete" @click="deleteGenre(genre.id)">Xóa</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Thêm/Sửa -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content animate-slide-up">
        <h3>{{ editingGenre.id ? 'Sửa Thể loại' : 'Thêm Thể loại' }}</h3>
        
        <form @submit.prevent="saveGenre" class="genre-form">
          <div class="form-group">
            <label>Tên Thể loại</label>
            <input type="text" v-model="editingGenre.name" @input="generateSlug" required placeholder="VD: Ngôn Tình, Hành Động..." />
          </div>          <div class="form-group">
            <label>Slug (Đường dẫn)</label>
            <input type="text" v-model="editingGenre.slug" required placeholder="VD: ngon-tinh, hanh-dong..." />
          </div>
          <div class="form-group">
            <label>Màu sắc Badge</label>
            <div class="color-picker-wrap">
              <input type="color" v-model="editingGenre.badgeColor" />
              <input type="text" v-model="editingGenre.badgeColor" placeholder="#HEXCODE" />
            </div>
            <p class="help-text">Chọn màu nền để thẻ thể loại này trông nổi bật hơn.</p>
          </div>

          <div class="form-actions">
            <button type="button" class="btn-cancel" @click="closeModal">Hủy</button>
            <button type="submit" class="btn-save" :disabled="saving">
              {{ saving ? 'Đang lưu...' : 'Lưu lại' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

definePageMeta({
  layout: 'admin',
})

useHead({ title: 'Quản lý Thể loại - Admin' })

const { get, post, put, del } = useApi()
const authStore = useAuthStore()

const loading = ref(false)
const saving = ref(false)
const genres = ref<any[]>([])

const showModal = ref(false)
const editingGenre = ref({ id: null as number | null, name: '', slug: '', badgeColor: '#9CA764' })

const fetchGenres = async () => {
  loading.value = true
  try {
    const data = await get<any[]>('/genres')
    genres.value = data || []
  } catch (error) {
    console.error('Error fetching genres:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchGenres()
})

const generateSlug = () => {
  if (!editingGenre.value.id) {
    editingGenre.value.slug = editingGenre.value.name
      .toLowerCase()
      .normalize('NFD').replace(/[\u0300-\u036f]/g, "")
      .replace(/đ/g, "d")
      .replace(/[^a-z0-9\s-]/g, "")
      .trim()
      .replace(/\s+/g, "-")
  }
}

const openModal = (genre: any = null) => {
  if (genre) {
    editingGenre.value = { ...genre }
  } else {
    editingGenre.value = { id: null, name: '', slug: '', badgeColor: '#9CA764' }
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveGenre = async () => {
  saving.value = true
  try {
    const headers = { Authorization: `Bearer ${authStore.token}` }
    if (editingGenre.value.id) {
      // Update
      const res = await put<any>(`/admin/genres/${editingGenre.value.id}`, editingGenre.value, { headers })
      const idx = genres.value.findIndex(g => g.id === editingGenre.value.id)
      if (idx !== -1) genres.value[idx] = res
    } else {
      // Create
      const res = await post<any>('/admin/genres', editingGenre.value, { headers })
      genres.value.push(res)
    }
    closeModal()
  } catch (error) {
    alert('Có lỗi xảy ra khi lưu thể loại!')
    console.error(error)
  } finally {
    saving.value = false
  }
}

const deleteGenre = async (id: number) => {
  if (confirm('Bạn có chắc chắn muốn xóa thể loại này?')) {
    try {
      const headers = { Authorization: `Bearer ${authStore.token}` }
      await del(`/admin/genres/${id}`, { headers })
      genres.value = genres.value.filter(g => g.id !== id)
    } catch (error) {
      alert('Có lỗi xảy ra khi xóa!')
      console.error(error)
    }
  }
}
</script>

<style scoped>
.admin-genres-page {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.page-header h2 {
  font-size: 1.5rem;
  color: #fff;
}
.btn-primary {
  background: var(--accent-primary);
  color: #000;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 700;
  cursor: pointer;
}

.table-container {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-card);
  border-radius: 12px;
  overflow: hidden;
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
}

.admin-table th, .admin-table td {
  padding: 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-subtle);
}

.admin-table th {
  background: rgba(0,0,0,0.2);
  color: var(--text-muted);
  font-weight: 600;
  font-size: 0.85rem;
}

.genre-badge {
  padding: 4px 10px;
  border-radius: 20px;
  color: #fff;
  font-size: 0.8rem;
  font-weight: 700;
  text-shadow: 0 1px 2px rgba(0,0,0,0.5);
}

.color-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: monospace;
  font-size: 0.9rem;
}
.color-box {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid rgba(255,255,255,0.2);
}

.actions { display: flex; gap: 8px; }
.btn-edit { background: rgba(156,167,100,0.2); color: var(--accent-primary); border: 1px solid rgba(156,167,100,0.3); padding: 6px 12px; border-radius: 6px; cursor: pointer; }
.btn-delete { background: rgba(239,68,68,0.1); color: #ef4444; border: 1px solid rgba(239,68,68,0.2); padding: 6px 12px; border-radius: 6px; cursor: pointer; }

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.modal-content {
  background: var(--bg-secondary);
  border: 1px solid var(--border-card);
  border-radius: 16px;
  padding: 30px;
  width: 100%;
  max-width: 500px;
}
.modal-content h3 {
  margin-bottom: 20px;
  font-size: 1.3rem;
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: 10px;
}

.genre-form { display: flex; flex-direction: column; gap: 16px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-group label { font-size: 0.9rem; color: var(--text-secondary); font-weight: 600; }
.form-group input[type="text"] { background: var(--bg-primary); border: 1px solid var(--border-card); padding: 10px 14px; border-radius: 8px; color: #fff; }

.color-picker-wrap { display: flex; gap: 12px; align-items: center; }
.color-picker-wrap input[type="color"] { width: 40px; height: 40px; border: none; border-radius: 8px; cursor: pointer; background: transparent; }
.color-picker-wrap input[type="text"] { flex: 1; }

.help-text { font-size: 0.8rem; color: var(--text-muted); margin: 0; }

.form-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 10px; }
.btn-cancel { background: transparent; border: 1px solid var(--border-subtle); color: var(--text-secondary); padding: 10px 20px; border-radius: 8px; cursor: pointer; }
.btn-save { background: var(--accent-primary); color: #000; font-weight: 700; border: none; padding: 10px 20px; border-radius: 8px; cursor: pointer; }
</style>
