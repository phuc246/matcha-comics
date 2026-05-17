<template>
  <div class="admin-story-detail" v-if="story">
    <div class="page-header">
      <button class="back-link" @click="router.push('/admin/novels')">← Quay lại danh sách</button>
      <h2>Quản lý Truyện Chữ: {{ story.title }}</h2>
    </div>

    <div class="story-info-card">
      <img :src="story.coverUrl" class="cover-img" alt="Cover" />
      <div class="info-content">
        <h1 class="story-title">{{ story.title }}</h1>
        <div class="meta-row">
          <span class="status-badge" :class="story.status">{{ story.status === 'completed' ? 'Hoàn thành' : 'Đang ra' }}</span>
          <span class="author-text">✍️ {{ story.author || 'Đang cập nhật' }}</span>
          <span class="author-text">👁️ {{ story.views?.toLocaleString() || 0 }} lượt xem</span>
        </div>
        <div class="tags">
          <span v-for="g in story.genres" :key="g.id" class="tag">{{ g.name }}</span>
        </div>
        <div class="desc-box">
          <div class="desc" v-html="story.description"></div>
        </div>
      </div>
      <div class="quick-actions">
        <NuxtLink :to="`/admin/novels/edit/${story.id}`" class="btn-secondary">✏️ Chỉnh sửa thông tin</NuxtLink>
      </div>
    </div>

    <!-- Chapter Management -->
    <div class="chapters-section">
      <div class="section-header">
        <h3>📚 Danh sách Chương ({{ chapters.length }})</h3>
        <NuxtLink :to="`/admin/novels/${story.id}/chapters/create`" class="btn-primary">
          <span>+</span> Thêm Chương Mới
        </NuxtLink>
      </div>

      <div class="table-container">
        <table class="admin-table">
          <thead>
            <tr>
              <th width="120">Số Chương</th>
              <th>Tên Chương</th>
              <th>Server</th>
              <th>Dung lượng</th>
              <th>Ngày đăng</th>
              <th width="120">Thao tác</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="chapters.length === 0">
              <td colspan="5" class="empty-state">
                <span>📭</span>
                <p>Truyện này chưa có chương nào. Hãy thêm chương đầu tiên!</p>
              </td>
            </tr>
            <tr v-for="chap in chapters" :key="chap.id">
              <td><strong class="chap-num">Chương {{ chap.number }}</strong></td>
              <td>{{ chap.title || `Chương ${chap.number}` }}</td>
              <td>
                <div class="server-tags">
                  <span v-for="s in chap.servers" :key="s.id" class="server-badge">{{ s.name }}</span>
                  <span v-if="!chap.servers?.length" class="server-badge empty">Chưa có</span>
                </div>
              </td>
              <td><span class="size-badge">{{ getChapterSize(chap) }}</span></td>
              <td>{{ new Date(chap.createdAt).toLocaleDateString('vi-VN') }}</td>
              <td>
                <div class="actions-cell">
                  <button class="icon-btn view" title="Xem nội dung" @click="previewChapter(chap)">👁️</button>
                  <NuxtLink :to="`/admin/novels/${story.id}/chapters/edit/${chap.id}`" class="icon-btn edit" title="Sửa chương">✏️</NuxtLink>
                  <button class="icon-btn delete" title="Xóa chương" @click="deleteChapter(chap.id)">🗑️</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Preview Chapter Modal -->
    <Teleport to="body">
      <div v-if="showPreview" class="preview-modal" @click.self="showPreview = false">
        <div class="preview-container">
          <div class="preview-header">
            <h3>Xem trước Chương {{ currentPreviewChap?.number }}</h3>
            <button class="close-btn" @click="showPreview = false">✕</button>
          </div>
          <div class="preview-content">
            <div class="novel-content-preview ql-editor" v-html="previewText"></div>
            <div v-if="!previewText" class="empty-preview">Chương này chưa có nội dung</div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>

  <div v-else class="loading-state">
    <div class="spinner"></div>
    <p>Đang tải dữ liệu...</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

definePageMeta({ layout: 'admin' })

const route = useRoute()
const router = useRouter()
const { get, del } = useApi()
const authStore = useAuthStore()

const story = ref<any>(null)
const chapters = ref<any[]>([])
const showPreview = ref(false)
const previewText = ref<string>('')
const currentPreviewChap = ref<any>(null)

const fetchData = async () => {
  try {
    const storyId = route.params.id
    const all = await get<any[]>('/stories?type=novel')
    story.value = all?.find(s => s.id === Number(storyId)) || null
    if (story.value) {
      useHead({ title: `Quản lý: ${story.value.title} - Admin` })
    }

    const chaps = await get<any[]>(`/admin/stories/${storyId}/chapters`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    chapters.value = chaps || []
  } catch (err) {
    console.error(err)
  }
}

onMounted(fetchData)

const getChapterSize = (chap: any) => {
  if (!chap.servers || chap.servers.length === 0) return '0 KB'
  const textServer = chap.servers.find((s: any) => s.type === 'text' || s.name === 'Nội dung chữ')
  if (textServer && textServer.content) {
    const kb = textServer.content.length / 1024
    return kb < 1024 ? `${kb.toFixed(1)} KB` : `${(kb / 1024).toFixed(2)} MB`
  }
  return '0 KB'
}

const previewChapter = (chap: any) => {
  currentPreviewChap.value = chap
  previewText.value = ''
  if (!chap.servers?.length) {
    showPreview.value = true
    return
  }
  
  const firstServer = chap.servers[0]
  try {
    previewText.value = firstServer.content || ''
    showPreview.value = true
  } catch {
    alert('Không thể đọc dữ liệu!')
  }
}

const deleteChapter = async (id: number) => {
  if (!confirm('Xóa chương này? Hành động không thể hoàn tác!')) return
  try {
    await del(`/admin/chapters/${id}`, { headers: { Authorization: `Bearer ${authStore.token}` } })
    chapters.value = chapters.value.filter(c => c.id !== id)
  } catch {
    alert('Xóa chương thất bại!')
  }
}
</script>

<style scoped>
.admin-story-detail { display: flex; flex-direction: column; gap: 24px; padding-bottom: 40px; }
.page-header { display: flex; flex-direction: column; gap: 8px; }
.back-link { background: none; border: none; color: #9CA764; cursor: pointer; font-weight: 600; text-align: left; padding: 0; font-size: 0.95rem; }
.page-header h2 { margin: 0; color: #fff; font-family: 'Montserrat', sans-serif; font-size: 1.4rem; }

.story-info-card { display: flex; gap: 24px; background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 28px; }
@media(max-width:768px){ .story-info-card { flex-direction: column; } }
.cover-img { width: 140px; height: 200px; object-fit: cover; border-radius: 10px; box-shadow: 0 8px 30px rgba(0,0,0,0.6); flex-shrink: 0; }
.info-content { flex: 1; display: flex; flex-direction: column; gap: 14px; }
.story-title { margin: 0; color: #fff; font-size: 1.6rem; font-family: 'Montserrat', sans-serif; }
.meta-row { display: flex; gap: 16px; align-items: center; flex-wrap: wrap; }
.status-badge { padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; font-weight: 700; }
.status-badge.ongoing { background: rgba(156,167,100,0.15); color: #9CA764; border: 1px solid rgba(156,167,100,0.3); }
.status-badge.completed { background: rgba(74,222,128,0.15); color: #4ade80; border: 1px solid rgba(74,222,128,0.3); }
.author-text { color: #A8A8B3; font-size: 0.9rem; }
.tags { display: flex; gap: 8px; flex-wrap: wrap; }
.tag { padding: 4px 12px; background: rgba(241,232,199,0.08); border-radius: 6px; font-size: 0.8rem; color: #F1E8C7; border: 1px solid rgba(241,232,199,0.1); }
.desc-box { background: rgba(0,0,0,0.25); padding: 16px; border-radius: 10px; border: 1px solid rgba(255,255,255,0.04); max-height: 200px; overflow-y: auto; }
.desc { margin: 0; color: #A8A8B3; font-size: 0.9rem; line-height: 1.7; }
.desc :deep(*) { color: #A8A8B3 !important; background-color: transparent !important; font-family: inherit !important; }
.quick-actions { display: flex; flex-direction: column; gap: 10px; flex-shrink: 0; }
.btn-secondary { background: rgba(255,255,255,0.06); color: #fff; padding: 10px 18px; border-radius: 8px; text-decoration: none; font-size: 0.9rem; font-weight: 600; transition: all 0.2s; border: 1px solid rgba(255,255,255,0.1); text-align: center; white-space: nowrap; }
.btn-secondary:hover { border-color: #9CA764; color: #9CA764; }

.chapters-section { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 24px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.section-header h3 { margin: 0; color: #F1E8C7; font-size: 1.1rem; }
.btn-primary { background: #9CA764; color: #000; padding: 10px 20px; border-radius: 8px; font-weight: 700; text-decoration: none; display: flex; align-items: center; gap: 6px; transition: all 0.2s; font-size: 0.9rem; }
.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(156,167,100,0.35); background: #b0c070; }

.table-container { overflow-x: auto; }
.admin-table { width: 100%; border-collapse: collapse; text-align: left; }
.admin-table th { background: rgba(0,0,0,0.25); padding: 14px 16px; font-size: 0.8rem; color: #5C5C6B; font-weight: 700; letter-spacing: 0.05em; text-transform: uppercase; }
.admin-table td { padding: 14px 16px; border-bottom: 1px solid rgba(255,255,255,0.03); color: #A8A8B3; font-size: 0.9rem; vertical-align: middle; }
.admin-table tr:last-child td { border-bottom: none; }
.admin-table tr:hover td { background: rgba(255,255,255,0.02); }
.chap-num { color: #F1E8C7; }
.empty-state { text-align: center; padding: 48px 16px !important; color: #5C5C6B !important; }
.empty-state span { font-size: 2rem; display: block; margin-bottom: 8px; }
.empty-state p { margin: 0; }

.server-tags { display: flex; gap: 6px; flex-wrap: wrap; }
.server-badge { padding: 2px 8px; border-radius: 4px; font-size: 0.75rem; font-weight: 600; background: rgba(96,165,250,0.1); color: #60a5fa; border: 1px solid rgba(96,165,250,0.2); }
.server-badge.empty { background: rgba(255,255,255,0.04); color: #5C5C6B; border-color: rgba(255,255,255,0.05); }

.size-badge { padding: 4px 10px; border-radius: 6px; font-size: 0.8rem; font-weight: 600; background: rgba(241,232,199,0.1); color: #F1E8C7; }

.actions-cell { display: flex; gap: 8px; }
.icon-btn { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.05); border-radius: 6px; cursor: pointer; transition: all 0.2s; font-size: 0.9rem; }
.icon-btn.view:hover { background: rgba(96,165,250,0.1); border-color: #60a5fa; }
.icon-btn.edit:hover { background: rgba(156,167,100,0.1); border-color: #9CA764; }
.icon-btn.delete:hover { background: rgba(239,68,68,0.1); border-color: #ef4444; }

/* Preview Modal */
.preview-modal { position: fixed; inset: 0; z-index: 9999; background: rgba(0,0,0,0.85); display: flex; align-items: center; justify-content: center; padding: 40px; }
.preview-container { background: #141720; width: 100%; max-width: 800px; max-height: 90vh; border-radius: 16px; border: 1px solid rgba(255,255,255,0.1); display: flex; flex-direction: column; overflow: hidden; box-shadow: 0 20px 50px rgba(0,0,0,0.5); }
.preview-header { padding: 16px 24px; border-bottom: 1px solid rgba(255,255,255,0.05); display: flex; justify-content: space-between; align-items: center; }
.preview-header h3 { margin: 0; font-size: 1.1rem; color: #F1E8C7; }
.close-btn { background: none; border: none; color: #5C5C6B; font-size: 1.4rem; cursor: pointer; }
.close-btn:hover { color: #fff; }
.preview-content { flex: 1; overflow-y: auto; padding: 20px; background: #fff; color: #000; border-radius: 0 0 16px 16px; }
.novel-content-preview { font-family: 'Inter', sans-serif; line-height: 1.8; font-size: 1.1rem; }
.empty-preview { padding: 40px; color: #5C5C6B; text-align: center; }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 80px; color: #A8A8B3; gap: 16px; }
.spinner { width: 40px; height: 40px; border: 3px solid rgba(156,167,100,0.2); border-top-color: #9CA764; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
