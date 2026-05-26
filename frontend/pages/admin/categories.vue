<template>
  <div class="admin-categories-page">
    <div class="page-header">
      <h2>Quản lý Danh mục (Thanh ngang)</h2>
      <button class="btn-primary" @click="openModal()">+ Thêm Danh mục</button>
    </div>

    <!-- Hướng dẫn sử dụng cho admin -->
    <div class="guide-card">
      <p class="guide-title">💡 Hướng dẫn quản lý cấu trúc Menu</p>
      <p class="guide-text">
        Hệ thống hỗ trợ <strong>Menu 2 cấp</strong> hiển thị trên thanh điều hướng ngang:<br />
        1. <strong>Danh mục cấp 1 (Menu chính)</strong>: Để trống phần "Danh mục cha" (Ví dụ: <em>Truyện Tranh, Truyện Chữ</em>).<br />
        2. <strong>Danh mục cấp 2 (Menu con / Dropdown)</strong>: Chọn "Danh mục cha" tương ứng (Ví dụ: <em>Mới cập nhật, Đã hoàn thành</em> thuộc Truyện Tranh).
      </p>
    </div>

    <!-- Danh sách Danh mục -->
    <div class="table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th>Thứ tự</th>
            <th>ID</th>
            <th>Tên Danh Mục</th>
            <th>Slug</th>
            <th>Đường dẫn (Path)</th>
            <th>Danh mục cha</th>
            <th>Hành động</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="7" class="text-center">Đang tải cấu trúc danh mục...</td>
          </tr>
          <tr v-else-if="flatCategories.length === 0">
            <td colspan="7" class="text-center">Chưa có danh mục nào được khởi tạo</td>
          </tr>
          <tr v-for="cat in flatCategories" :key="cat.id" v-else>
            <td>
              <span class="order-badge">{{ cat.order }}</span>
            </td>
            <td>#{{ cat.id }}</td>
            <td>
              <div class="cat-name-cell">
                <span class="cat-icon" v-if="cat.icon">{{ cat.icon }}</span>
                <span :class="{ 'parent-cat': !cat.parentId, 'child-cat': cat.parentId }">
                  {{ cat.parentId ? '— ' : '' }}{{ cat.name }}
                </span>
              </div>
            </td>
            <td>{{ cat.slug }}</td>
            <td><code class="path-code">{{ cat.path || '—' }}</code></td>
            <td>
              <span v-if="cat.parentId" class="parent-label">
                📂 {{ getParentName(cat.parentId) }}
              </span>
              <span v-else class="root-label">👑 Cấp 1 (Gốc)</span>
            </td>
            <td>
              <div class="actions">
                <button class="btn-edit" @click="openModal(cat)">Sửa</button>
                <button class="btn-delete" @click="deleteCategory(cat.id)">Xóa</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Thêm/Sửa -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content animate-slide-up">
        <h3>{{ editingCategory.id ? 'Sửa Danh mục' : 'Thêm Danh mục' }}</h3>
        
        <form @submit.prevent="saveCategory" class="category-form">
          <div class="form-group">
            <label>Tên Danh mục</label>
            <input type="text" v-model="editingCategory.name" @input="generateSlug" required placeholder="VD: Truyện Tranh, Mới cập nhật..." />
          </div>

          <div class="form-group">
            <label>Slug (Đường dẫn tinh gọn)</label>
            <input type="text" v-model="editingCategory.slug" required placeholder="VD: truyen-tranh, moi-cap-nhat..." />
          </div>

          <div class="form-group">
            <label>Đường dẫn liên kết (Path)</label>
            <input type="text" v-model="editingCategory.path" placeholder="VD: /truyen-tranh hoặc /truyen-tranh?sort=latest" />
            <p class="help-text">Nhập URL của trang khi bấm vào danh mục này.</p>
          </div>

          <div class="form-row">
            <div class="form-group flex-1">
              <label>Icon / Emoji</label>
              <input type="text" v-model="editingCategory.icon" placeholder="VD: 📚, 🎨, 🆕..." />
            </div>

            <div class="form-group flex-1">
              <label>Thứ tự hiển thị (Order)</label>
              <input type="number" v-model.number="editingCategory.order" required />
            </div>
          </div>

          <div class="form-group">
            <label>Danh mục cha (Cấp 1)</label>
            <select v-model="editingCategory.parentId">
              <option :value="null">👑 Không có (Đặt làm Danh mục chính cấp 1)</option>
              <option 
                v-for="parent in availableParents" 
                :key="parent.id" 
                :value="parent.id"
              >
                📂 {{ parent.name }}
              </option>
            </select>
            <p class="help-text">Chọn danh mục cha để hiển thị danh mục này dưới dạng Menu con (Dropdown).</p>
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
import { ref, computed, onMounted } from 'vue'

definePageMeta({
  layout: 'admin',
})

useHead({ title: 'Quản lý Danh mục Menu - Admin' })

const { get, post, put, del } = useApi()
const authStore = useAuthStore()

const loading = ref(false)
const saving = ref(false)
const flatCategories = ref<any[]>([])
const showModal = ref(false)

const editingCategory = ref({
  id: null as number | null,
  name: '',
  slug: '',
  path: '',
  icon: '',
  parentId: null as number | null,
  order: 0,
})

const fetchCategories = async () => {
  loading.value = true
  try {
    const data = await get<any[]>('/categories')
    flatCategories.value = data || []
  } catch (error) {
    console.error('Error fetching categories:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCategories()
})

const getParentName = (parentId: number) => {
  const p = flatCategories.value.find(c => c.id === parentId)
  return p ? p.name : 'Không xác định'
}

// Parent options available in modal (Only top level categories, and not itself)
const availableParents = computed(() => {
  return flatCategories.value.filter(c => {
    // Must be a top level category (parentId is null/nil)
    const isTopLevel = !c.parentId
    // Must not be itself (to prevent circular references)
    const isNotSelf = c.id !== editingCategory.value.id
    return isTopLevel && isNotSelf
  })
})

const generateSlug = () => {
  if (!editingCategory.value.id) {
    editingCategory.value.slug = editingCategory.value.name
      .toLowerCase()
      .normalize('NFD').replace(/[\u0300-\u036f]/g, "")
      .replace(/đ/g, "d")
      .replace(/[^a-z0-9\s-]/g, "")
      .trim()
      .replace(/\s+/g, "-")
  }
}

const openModal = (cat: any = null) => {
  if (cat) {
    editingCategory.value = {
      id: cat.id,
      name: cat.name,
      slug: cat.slug,
      path: cat.path || '',
      icon: cat.icon || '',
      parentId: cat.parentId || null,
      order: cat.order || 0,
    }
  } else {
    editingCategory.value = {
      id: null,
      name: '',
      slug: '',
      path: '',
      icon: '',
      parentId: null,
      order: flatCategories.value.length + 1,
    }
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveCategory = async () => {
  saving.value = true
  try {
    const headers = { Authorization: `Bearer ${authStore.token}` }
    
    // Explicitly handle empty string path/icon
    const payload = {
      ...editingCategory.value,
      path: editingCategory.value.path || '',
      icon: editingCategory.value.icon || '',
      parentId: editingCategory.value.parentId || null,
    }

    if (editingCategory.value.id) {
      // Update
      const res = await put<any>(`/admin/categories/${editingCategory.value.id}`, payload, { headers })
      const idx = flatCategories.value.findIndex(c => c.id === editingCategory.value.id)
      if (idx !== -1) flatCategories.value[idx] = res
    } else {
      // Create
      const res = await post<any>('/admin/categories', payload, { headers })
      flatCategories.value.push(res)
    }

    // Sort again by order and id
    flatCategories.value.sort((a, b) => {
      if (a.order !== b.order) return a.order - b.order
      return a.id - b.id
    })

    closeModal()
  } catch (error) {
    alert('Có lỗi xảy ra khi lưu danh mục!')
    console.error(error)
  } finally {
    saving.value = false
  }
}

const deleteCategory = async (id: number) => {
  if (confirm('Bạn có chắc chắn muốn xóa danh mục này? Các danh mục con liên quan sẽ được tự động chuyển về danh mục cấp 1.')) {
    try {
      const headers = { Authorization: `Bearer ${authStore.token}` }
      await del(`/admin/categories/${id}`, { headers })
      
      // Update state
      flatCategories.value = flatCategories.value.filter(c => c.id !== id).map(c => {
        if (c.parentId === id) {
          return { ...c, parentId: null }
        }
        return c
      })
    } catch (error) {
      alert('Có lỗi xảy ra khi xóa danh mục!')
      console.error(error)
    }
  }
}
</script>

<style scoped>
.admin-categories-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
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
  transition: var(--transition-base);
}

.btn-primary:hover {
  background: var(--accent-hover);
  transform: translateY(-2px);
}

/* Guide box */
.guide-card {
  background: rgba(156, 167, 100, 0.08);
  border: 1px dashed var(--accent-primary);
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 24px;
}

.guide-title {
  font-weight: 700;
  color: var(--accent-gold);
  margin-bottom: 6px;
  font-size: 0.95rem;
}

.guide-text {
  font-size: 0.85rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
}

/* Table styling */
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
  padding: 14px 18px;
  text-align: left;
  border-bottom: 1px solid var(--border-subtle);
}

.admin-table th {
  background: rgba(0, 0, 0, 0.2);
  color: var(--text-muted);
  font-weight: 600;
  font-size: 0.85rem;
}

.order-badge {
  background: var(--bg-primary);
  border: 1px solid var(--border-card);
  padding: 3px 8px;
  border-radius: 6px;
  font-family: monospace;
  font-size: 0.8rem;
  color: var(--accent-gold);
}

.cat-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cat-icon {
  font-size: 1.1rem;
}

.parent-cat {
  font-weight: 700;
  color: #fff;
}

.child-cat {
  font-weight: 500;
  color: var(--text-secondary);
}

.path-code {
  font-family: monospace;
  background: rgba(0, 0, 0, 0.2);
  padding: 2px 6px;
  border-radius: 4px;
  color: var(--accent-primary);
  font-size: 0.85rem;
}

.parent-label {
  background: rgba(156, 167, 100, 0.12);
  color: var(--accent-primary);
  padding: 3px 10px;
  border-radius: 100px;
  font-size: 0.78rem;
  font-weight: 600;
}

.root-label {
  background: rgba(255, 215, 0, 0.08);
  color: #FFD700;
  padding: 3px 10px;
  border-radius: 100px;
  font-size: 0.78rem;
  font-weight: 600;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-edit {
  background: rgba(156, 167, 100, 0.2);
  color: var(--accent-primary);
  border: 1px solid rgba(156, 167, 100, 0.3);
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
}

.btn-delete {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
}

/* Modal form */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(5px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: var(--bg-secondary);
  border: 1px solid var(--border-card);
  border-radius: 16px;
  padding: 30px;
  width: 100%;
  max-width: 550px;
}

.modal-content h3 {
  margin-bottom: 20px;
  font-size: 1.3rem;
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: 10px;
}

.category-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.flex-1 {
  flex: 1;
}

.form-group label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  font-weight: 600;
}

.form-group input, .form-group select {
  background: var(--bg-primary);
  border: 1px solid var(--border-card);
  padding: 10px 14px;
  border-radius: 8px;
  color: #fff;
  outline: none;
  font-family: inherit;
}

.form-group input:focus, .form-group select:focus {
  border-color: var(--accent-primary);
}

.help-text {
  font-size: 0.8rem;
  color: var(--text-muted);
  margin: 0;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 10px;
}

.btn-cancel {
  background: transparent;
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
}

.btn-save {
  background: var(--accent-primary);
  color: #000;
  font-weight: 700;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
}
</style>
