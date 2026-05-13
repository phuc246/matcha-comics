<template>
  <div class="admin-catalog-page">
    <div class="page-header-actions">
      <div class="search-filter">
        <input type="text" v-model="search" placeholder="Tìm kiếm truyện tranh..." class="search-input" />
        <select v-model="filterStatus" class="filter-select">
          <option value="">Tất cả trạng thái</option>
          <option value="ongoing">Đang ra</option>
          <option value="completed">Hoàn thành</option>
        </select>
      </div>
      <NuxtLink to="/admin/comics/create" class="btn-primary">
        <span class="icon">+</span> Thêm Truyện Mới
      </NuxtLink>
    </div>

    <div class="table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th width="80">Ảnh</th>
            <th>Tên truyện</th>
            <th>Thể loại</th>
            <th>Trạng thái</th>
            <th>Lượt xem</th>
            <th>Chương</th>
            <th width="120">Thao tác</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in filteredComics" :key="c.id" class="clickable-row" @click="router.push(`/admin/comics/${c.id}`)">
            <td><img :src="c.coverUrl" alt="" class="table-cover" /></td>
            <td>
              <div class="title-cell">
                <span class="title-text">{{ c.title }}</span>
                <span class="slug-text">{{ c.slug }}</span>
              </div>
            </td>
            <td>
              <div class="tags-wrap">
                <span v-for="g in c.genres?.slice(0,2)" :key="g.id || g" class="tag">{{ g.name || g }}</span>
              </div>
            </td>
            <td>
              <span class="status-dot" :class="c.status"></span>
              {{ c.status === 'completed' ? 'Hoàn thành' : 'Đang ra' }}
            </td>
            <td>{{ formatNumber(c.views) }}</td>
            <td>{{ c.chapters?.length || 0 }}</td>
            <td>
              <div class="actions-cell">
                <NuxtLink :to="`/admin/comics/${c.id}/chapters/create`" class="icon-btn chapter" title="Thêm chương" @click.stop><span class="icon">📖</span></NuxtLink>
                <NuxtLink :to="`/admin/comics/edit/${c.id}`" class="icon-btn edit" title="Sửa" @click.stop><span class="icon">✏️</span></NuxtLink>
                <button class="icon-btn delete" title="Xóa" @click.stop="deleteStory(c.id)"><span class="icon">🗑️</span></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination-footer">
      <span class="results-info">Hiển thị {{ filteredComics.length }} trên tổng số {{ comics.length }} truyện</span>
      <div class="pagination-btns">
        <button class="page-btn" disabled>←</button>
        <button class="page-btn active">1</button>
        <button class="page-btn">2</button>
        <button class="page-btn">→</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

definePageMeta({
  layout: 'admin',
  title: 'Quản lý Truyện Tranh'
})

useHead({ title: 'Quản lý Truyện Tranh - Matcha Comic' })

const { get, del } = useApi()
const authStore = useAuthStore()
const router = useRouter()
const comics = ref<any[]>([])

onMounted(async () => {
  try {
    const data = await get<any[]>('/stories?type=comic')
    if (data) {
      comics.value = data.sort((a, b) => b.id - a.id)
    }
  } catch (err) {
    console.error('Error fetching comics:', err)
  }
})

const deleteStory = async (id: number) => {
  if (!confirm('Bạn có chắc chắn muốn xóa truyện này?')) return;
  try {
    await del(`/admin/stories/${id}`, { headers: { Authorization: `Bearer ${authStore.token}` } })
    comics.value = comics.value.filter(c => c.id !== id)
    alert('Đã xóa truyện thành công!')
  } catch (err) {
    alert('Có lỗi xảy ra khi xóa!')
  }
}
const search = ref('')
const filterStatus = ref('')

const filteredComics = computed(() => {
  return comics.value.filter(c => {
    const matchesSearch = c.title.toLowerCase().includes(search.value.toLowerCase())
    const matchesStatus = filterStatus.value ? c.status === filterStatus.value : true
    return matchesSearch && matchesStatus
  })
})

const formatNumber = (n: number) => n.toLocaleString()
</script>

<style scoped>
.admin-catalog-page { display: flex; flex-direction: column; gap: 24px; }

.page-header-actions { display: flex; justify-content: space-between; align-items: center; gap: 20px; flex-wrap: wrap; }
.search-filter { display: flex; gap: 12px; flex: 1; max-width: 600px; }
.search-input { flex: 1; background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 8px; padding: 10px 16px; color: #fff; font-size: 0.9rem; }
.search-input:focus { border-color: #9CA764; outline: none; }
.filter-select { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 8px; padding: 10px 12px; color: #fff; font-size: 0.9rem; outline: none; }

.btn-primary { background: #9CA764; color: #000; padding: 10px 20px; border-radius: 8px; font-weight: 700; text-decoration: none; display: flex; align-items: center; gap: 8px; transition: all 0.2s; }
.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(156,167,100,0.3); }

.table-container { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 12px; overflow: hidden; }
.admin-table { width: 100%; border-collapse: collapse; text-align: left; }
.admin-table th { background: rgba(0,0,0,0.2); padding: 16px; font-size: 0.85rem; color: #5C5C6B; font-weight: 600; }
.admin-table td { padding: 16px; border-bottom: 1px solid rgba(255,255,255,0.03); color: #A8A8B3; font-size: 0.9rem; vertical-align: middle; }
.admin-table tr { cursor: pointer; }
.clickable-row:hover td { background: rgba(156,167,100,0.04) !important; }

.table-cover { width: 48px; height: 64px; object-fit: cover; border-radius: 4px; background: #000; }
.title-cell { display: flex; flex-direction: column; gap: 4px; }
.title-text { color: #fff; font-weight: 600; font-size: 0.95rem; }
.slug-text { font-size: 0.75rem; color: #5C5C6B; font-family: monospace; }

.tags-wrap { display: flex; gap: 6px; flex-wrap: wrap; }
.tag { padding: 2px 8px; background: rgba(255,255,255,0.05); border-radius: 4px; font-size: 0.7rem; color: #F1E8C7; }

.status-dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 6px; }
.status-dot.ongoing { background: #9CA764; box-shadow: 0 0 8px #9CA764; }
.status-dot.completed { background: #4ade80; box-shadow: 0 0 8px #4ade80; }

.actions-cell { display: flex; gap: 8px; }
.icon-btn { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.05); border-radius: 6px; cursor: pointer; transition: all 0.2s; color: inherit; text-decoration: none; }
.icon-btn:hover { background: rgba(255,255,255,0.05); border-color: rgba(255,255,255,0.2); }
.icon-btn.chapter:hover { color: #60a5fa; border-color: #60a5fa; }
.icon-btn.edit:hover { color: #9CA764; border-color: #9CA764; }
.icon-btn.delete:hover { color: #ef4444; border-color: #ef4444; }

.pagination-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 10px; }
.results-info { font-size: 0.85rem; color: #5C5C6B; }
.pagination-btns { display: flex; gap: 8px; }
.page-btn { width: 36px; height: 36px; display: flex; align-items: center; justify-content: center; background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 6px; color: #A8A8B3; cursor: pointer; }
.page-btn.active { background: #9CA764; color: #000; border: none; font-weight: 700; }
.page-btn:disabled { opacity: 0.3; cursor: not-allowed; }
</style>
