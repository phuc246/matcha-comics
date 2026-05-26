<template>
  <div class="create-chapter-page" v-if="story">
    <div class="page-header">
      <button class="back-link" @click="router.back()">← Quay lại truyện</button>
      <h2>Thêm Chương Mới: {{ story.title }}</h2>
    </div>

    <div class="create-grid">
      <!-- Left: Form -->
      <div class="form-panel">
        <div class="form-section">
          <h3>📋 Thông tin Chương</h3>
          <div class="form-group">
            <label>Số chương <span class="required">*</span></label>
            <input type="number" v-model="form.number" placeholder="VD: 1" min="0" step="1" />
          </div>
          <div class="form-group">
            <label>Tên chương <span class="required">*</span></label>
            <input type="text" v-model="form.title" placeholder="VD: Chương 1: Sự khởi đầu" />
          </div>
        </div>

        <div class="form-section">
          <h3>✍️ Nội dung chương</h3>
          <RichTextEditor
            v-model="form.content"
            placeholder="Nhập nội dung truyện tại đây..."
          />
        </div>
      </div>

      <!-- Right: Sidebar -->
      <aside class="sidebar">
        <button class="btn-publish" :disabled="isSubmitting" @click="handleSubmit">
          {{ isSubmitting ? '⏳ Đang lưu...' : '🚀 Đăng Chương' }}
        </button>
      </aside>
    </div>
  </div>
  <div v-else class="loading-state">
    <div class="spinner"></div>
    <p>Đang tải thông tin truyện...</p>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

definePageMeta({ layout: 'admin' })

const router = useRouter()
const route = useRoute()
const { get, post } = useApi()
const authStore = useAuthStore()

const story = ref<any>(null)
const isSubmitting = ref(false)

const form = reactive({
  number: 1,
  title: '',
  content: ''
})

onMounted(async () => {
  try {
    const storyId = route.params.id
    const stories = await get<any[]>('/stories?type=novel')
    story.value = stories?.find(s => s.id === Number(storyId)) || null
    
    if (story.value) {
      // Auto-calculate next chapter number
      const chaps = await get<any[]>(`/admin/stories/${storyId}/chapters`, {
        headers: { Authorization: `Bearer ${authStore.token}` }
      })
      if (chaps && chaps.length > 0) {
        const maxNum = Math.max(...chaps.map(c => c.number || 0))
        form.number = maxNum + 1
        form.title = `Chương ${maxNum + 1}`
      }
    }
  } catch (err) {
    console.error(err)
  }
})

const handleSubmit = async () => {
  if (!form.number || !form.title || !form.content) {
    return alert('Vui lòng điền đầy đủ Số chương, Tên chương và Nội dung!')
  }
  
  // Calculate raw character count
  const rawText = form.content.replace(/<[^>]*>/g, '').trim()
  if (rawText.length > 15000) {
    return alert('Nội dung chương truyện tối đa 15000 ký tự (chữ thô)!')
  }

  isSubmitting.value = true
  try {
    const payload = {
      story_id: Number(route.params.id),
      number: Number(form.number),
      title: form.title,
      servers: [
        {
          name: 'Nội dung chữ',
          type: 'text',
          content: form.content
        }
      ]
    }

    await post('/admin/chapters', payload, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })

    alert(`✅ Đã đăng ${form.title} thành công!`)
    router.push(`/admin/novels`)
  } catch (err: any) {
    console.error(err)
    alert('Lỗi: ' + (err.data?.error || err.message))
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.create-chapter-page { display: flex; flex-direction: column; gap: 24px; padding-bottom: 40px; }
.page-header { display: flex; align-items: center; gap: 16px; }
.back-link { background: none; border: none; color: #9CA764; cursor: pointer; font-weight: 600; font-size: 0.95rem; }
.page-header h2 { margin: 0; color: #fff; font-family: 'Montserrat', sans-serif; }

.create-grid { display: grid; grid-template-columns: 1fr 300px; gap: 24px; align-items: start; }
@media(max-width:1024px) { .create-grid { grid-template-columns: 1fr; } }

.form-panel { display: flex; flex-direction: column; gap: 24px; }
.form-section { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 24px; display: flex; flex-direction: column; gap: 16px; }
.form-section h3 { margin: 0; color: #F1E8C7; font-size: 1rem; border-bottom: 1px solid rgba(255,255,255,0.05); padding-bottom: 12px; }

.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-group label { font-size: 0.85rem; color: #A8A8B3; font-weight: 600; }
.required { color: #ef4444; }
.form-group input { background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; padding: 12px; color: #fff; font-size: 0.9rem; }
.form-group input:focus { border-color: #9CA764; outline: none; }

.sidebar { position: sticky; top: 24px; display: flex; flex-direction: column; gap: 24px; }

.btn-publish { background: #9CA764; color: #000; border: none; padding: 16px; border-radius: 12px; font-weight: 800; font-size: 1rem; cursor: pointer; transition: all 0.2s; }
.btn-publish:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(156,167,100,0.4); }

.loading-state { text-align: center; padding: 100px; color: #A8A8B3; }
.spinner { width: 40px; height: 40px; border: 3px solid rgba(156,167,100,0.2); border-top-color: #9CA764; border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 16px; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
