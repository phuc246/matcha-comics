<template>
  <div class="edit-chapter-page">
    <div class="page-header">
      <button class="back-link" @click="router.back()">← Quay lại truyện</button>
      <h2>Chỉnh sửa Chương: {{ form.title || form.number }}</h2>
    </div>

    <div class="create-grid" v-if="!loading">
      <!-- Left: Form -->
      <div class="form-panel">
        <div class="form-section">
          <h3>📋 Thông tin Chương</h3>
          <div class="form-group">
            <label>Số chương <span class="required">*</span></label>
            <input type="number" v-model="form.number" placeholder="VD: 1" min="0" step="1" disabled />
          </div>
          <div class="form-group">
            <label>Tên chương <span class="required">*</span></label>
            <input type="text" v-model="form.title" placeholder="VD: Chương 1: Sự khởi đầu" />
          </div>
        </div>

        <div class="form-section">
          <h3>✍️ Nội dung chương</h3>
          <div class="quill-chapter-wrapper">
            <QuillEditor
              v-model:content="form.content"
              content-type="html"
              theme="snow"
              placeholder="Nhập nội dung truyện tại đây..."
              :toolbar="fullToolbar"
            />
          </div>
        </div>
      </div>

      <!-- Right: Sidebar -->
      <aside class="sidebar">
        <div class="form-section">
          <h3>📊 Thống kê</h3>
          <div class="stat-row">
            <span>Số từ ước tính:</span>
            <strong>{{ wordCount }} từ</strong>
          </div>
          <div class="stat-row">
            <span>Dung lượng:</span>
            <strong>{{ (form.content.length / 1024).toFixed(1) }} KB</strong>
          </div>
          <div class="stat-row" :class="{ 'warning': charCount > 15000 }">
            <span>Độ dài (chữ thô):</span>
            <strong>{{ charCount }}/15000</strong>
          </div>
        </div>

        <button class="btn-publish" :disabled="isSubmitting" @click="handleSubmit">
          {{ isSubmitting ? '⏳ Đang lưu...' : '💾 Lưu Thay Đổi' }}
        </button>
      </aside>
    </div>
    
    <div v-else class="loading-state">
      <div class="spinner"></div>
      <p>Đang tải dữ liệu chương...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { QuillEditor } from '@vueup/vue-quill'
import Quill from 'quill'
import '@vueup/vue-quill/dist/vue-quill.snow.css'

// Register extra fonts
const Font = Quill.import('formats/font')
Font.whitelist = [false, 'serif', 'monospace', 'arial', 'roboto', 'montserrat', 'inter', 'times-new-roman']
Quill.register(Font, true)

// Register Line Height
const Parchment = Quill.import('parchment')
const LineHeightStyle = new Parchment.Attributor.Style('lineHeight', 'line-height', {
  scope: Parchment.Scope.BLOCK,
  whitelist: ['1', '1.2', '1.5', '1.8', '2', '2.5', '3']
})
Quill.register(LineHeightStyle, true)

definePageMeta({ layout: 'admin' })

const router = useRouter()
const route = useRoute()
const { get, put } = useApi()
const authStore = useAuthStore()

const loading = ref(true)
const isSubmitting = ref(false)

const form = reactive({
  number: 0,
  title: '',
  content: ''
})

const fullToolbar = [
  [{ 'font': [false, 'serif', 'monospace', 'arial', 'roboto', 'montserrat', 'inter', 'times-new-roman'] }, { 'size': ['small', false, 'large', 'huge'] }],
  ['bold', 'italic', 'underline', 'strike'],
  [{ 'color': [] }, { 'background': [] }],
  [{ 'script': 'sub'}, { 'script': 'super' }],
  [{ 'lineHeight': ['1', '1.2', '1.5', '1.8', '2', '2.5', '3'] }],
  [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
  [{ 'list': 'ordered'}, { 'list': 'bullet' }],
  [{ 'indent': '-1'}, { 'indent': '+1' }],
  [{ 'direction': 'rtl' }],
  [{ 'align': [] }],
  ['blockquote', 'code-block'],
  ['link', 'image', 'video', 'formula'],
  ['clean']
]

const wordCount = computed(() => {
  const text = form.content.replace(/<[^>]*>/g, ' ').trim()
  return text ? text.split(/\s+/).length : 0
})

const charCount = computed(() => {
  if (!form.content) return 0
  const temp = document.createElement('div')
  temp.innerHTML = form.content
  return temp.textContent?.length || 0
})

onMounted(async () => {
  try {
    const chap = await get<any>(`/admin/chapters/${route.params.chapterId}`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    
    if (chap) {
      form.number = chap.number
      form.title = chap.title
      
      const textServer = chap.servers?.find((s: any) => s.type === 'text' || s.name === 'Nội dung chữ')
      if (textServer) {
        form.content = textServer.content || ''
      }
    }
  } catch (err) {
    console.error(err)
    alert('Không thể tải dữ liệu chương!')
  } finally {
    loading.value = false
  }
})

const handleSubmit = async () => {
  if (!form.title || !form.content) {
    return alert('Vui lòng điền đầy đủ Tên chương và Nội dung!')
  }
  if (charCount.value > 15000) {
    return alert('Nội dung chương truyện tối đa 15000 ký tự (chữ thô)!')
  }

  isSubmitting.value = true
  try {
    const payload = {
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

    await put(`/admin/chapters/${route.params.chapterId}`, payload, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })

    alert('Cập nhật chương thành công!')
    router.push(`/admin/novels/${route.params.id}`)
  } catch (err: any) {
    console.error(err)
    alert('Lỗi: ' + (err.data?.error || err.message))
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.edit-chapter-page { display: flex; flex-direction: column; gap: 24px; padding-bottom: 40px; }
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
.form-group input:disabled { opacity: 0.6; cursor: not-allowed; }

/* Quill White Theme */
.quill-chapter-wrapper { background: #fff; border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; overflow: hidden; }
:deep(.ql-toolbar.ql-snow) { border: none; border-bottom: 1px solid #ddd; background: #f8f9fa; }
:deep(.ql-container.ql-snow) { border: none; min-height: 500px; }
:deep(.ql-editor) { font-size: 1.15rem; color: #333; line-height: 1.8; font-family: 'Inter', sans-serif; padding: 40px; }
:deep(.ql-editor.ql-blank::before) { color: #999; font-style: normal; left: 40px; }
:deep(.ql-snow .ql-stroke) { stroke: #444; }
:deep(.ql-snow .ql-fill) { fill: #444; }
:deep(.ql-snow .ql-picker) { color: #444; }

.stat-row.warning { color: #ef4444; font-weight: bold; }

.sidebar { position: sticky; top: 24px; display: flex; flex-direction: column; gap: 24px; }
.stat-row { display: flex; justify-content: space-between; font-size: 0.9rem; padding: 8px 0; border-bottom: 1px solid rgba(255,255,255,0.05); }
.stat-row span { color: #A8A8B3; }
.stat-row strong { color: #F1E8C7; }

.btn-publish { background: #9CA764; color: #000; border: none; padding: 16px; border-radius: 12px; font-weight: 800; font-size: 1rem; cursor: pointer; transition: all 0.2s; }
.btn-publish:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(156,167,100,0.4); }

.loading-state { text-align: center; padding: 100px; color: #A8A8B3; }
.spinner { width: 40px; height: 40px; border: 3px solid rgba(156,167,100,0.2); border-top-color: #9CA764; border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 16px; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Quill Font Labels in Toolbar */
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="arial"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="arial"]::before) { content: 'Arial'; font-family: 'Arial'; }
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="roboto"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="roboto"]::before) { content: 'Roboto'; font-family: 'Roboto'; }
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="montserrat"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="montserrat"]::before) { content: 'Montserrat'; font-family: 'Montserrat'; }
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="inter"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="inter"]::before) { content: 'Inter'; font-family: 'Inter'; }
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="times-new-roman"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="times-new-roman"]::before) { content: 'Times New Roman'; font-family: 'Times New Roman'; }

/* Quill Font Classes */
:deep(.ql-font-arial) { font-family: 'Arial', sans-serif; }
:deep(.ql-font-roboto) { font-family: 'Roboto', sans-serif; }
:deep(.ql-font-montserrat) { font-family: 'Montserrat', sans-serif; }
:deep(.ql-font-inter) { font-family: 'Inter', sans-serif; }
:deep(.ql-font-times-new-roman) { font-family: 'Times New Roman', serif; }

/* Custom Line Height Dropdown */
:deep(.ql-snow .ql-picker.ql-lineHeight) { width: 100px; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-label[data-value]::before),
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value]::before) { content: attr(data-value); }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-label::before),
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item::before) { content: 'Giãn dòng'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1"]::before) { content: '1.0'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1.2"]::before) { content: '1.2'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1.5"]::before) { content: '1.5'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1.8"]::before) { content: '1.8'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="2"]::before) { content: '2.0'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="2.5"]::before) { content: '2.5'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="3"]::before) { content: '3.0'; }
</style>
