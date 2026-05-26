<template>
  <div class="admin-create-page">
    <div class="page-header">
      <button class="back-link" @click="router.back()">← Quay lại</button>
      <h2>Chỉnh sửa Tiểu Thuyết</h2>
    </div>

    <div class="create-grid" v-if="!loading">
      <!-- Left: Form Info -->
      <div class="create-form-panel">
        <div class="form-section">
          <h3>Thông tin cơ bản</h3>
          <div class="form-grid">
            <div class="form-group full">
              <label>Tên truyện <span class="required">*</span></label>
              <input type="text" v-model="form.title" maxlength="50" required placeholder="Ví dụ: Đấu Phá Thương Khung (Tối đa 50 ký tự)" />
              <small class="char-count">{{ form.title.length }}/50</small>
            </div>

            <div class="form-group">
              <label>Tác giả</label>
              <input type="text" v-model="form.author" placeholder="Tên tác giả..." />
            </div>

            <div class="form-group">
              <label>Nhà xuất bản</label>
              <input type="text" v-model="form.publisher" placeholder="Nhà xuất bản..." />
            </div>
            
            <GenreSelect 
              v-model="form.genres" 
              :required="true"
            />
            
            <div class="form-group">
              <label>Trạng thái</label>
              <div class="radio-group">
                <label class="radio-item">
                  <input type="radio" value="ongoing" v-model="form.status" disabled title="Tự động cập nhật dựa trên số chương" />
                  <span style="opacity: 0.6">Đang ra (Tự động)</span>
                </label>
                <label class="radio-item">
                  <input type="radio" value="hiatus" v-model="form.status" />
                  <span>Tạm ngưng</span>
                </label>
                <label class="radio-item">
                  <input type="radio" value="completed" v-model="form.status" />
                  <span>Hoàn thành</span>
                </label>
              </div>
            </div>
            
            <div class="form-group full">
              <label>Mô tả truyện <span class="required">*</span></label>
              <div class="quill-wrapper">
                <QuillEditor
                  v-model:content="form.description"
                  content-type="html"
                  theme="snow"
                  placeholder="Tóm tắt nội dung... (Tối đa 700 ký tự)"
                  :toolbar="quillToolbar"
                />
              </div>
              <small class="char-count" :class="{ 'error': descriptionLength > 700 }">
                {{ descriptionLength }}/700 (chữ thô)
              </small>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Media/Settings -->
      <aside class="create-sidebar">
        <div class="form-section">
          <h3>Ảnh bìa (Cover)</h3>
          <div class="cover-upload-box">
            <div v-if="form.coverUrl" class="cover-preview">
              <img :src="form.coverUrl" alt="" />
              <button class="remove-cover" @click="form.coverUrl = ''">×</button>
            </div>
            <div v-else class="upload-placeholder" @click="triggerUpload">
              <span class="icon">📸</span>
              <span>Nhấp để tải ảnh</span>
              <input type="file" ref="fileInput" hidden @change="handleFileUpload" />
            </div>
          </div>
        </div>

        <div class="publish-panel">
          <button class="btn-publish" @click="handleSubmit" :disabled="isSubmitting">
            {{ isSubmitting ? 'Đang xử lý...' : 'Lưu Thay Đổi' }}
          </button>
        </div>
      </aside>
    </div>

    <div v-else class="loading-state">
      <div class="spinner"></div>
      <p>Đang tải thông tin truyện...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
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

definePageMeta({
  layout: 'admin',
  title: 'Chỉnh sửa Tiểu Thuyết'
})

const router = useRouter()
const route = useRoute()
const { get, put } = useApi()
const authStore = useAuthStore()

const loading = ref(true)
const isSubmitting = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const form = reactive({
  title: '',
  genres: [] as number[],
  author: '',
  publisher: '',
  status: 'ongoing',
  description: '',
  coverUrl: ''
})

const quillToolbar = [
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

const descriptionLength = computed(() => {
  if (!form.description) return 0
  const temp = document.createElement('div')
  temp.innerHTML = form.description
  return temp.textContent?.length || 0
})

onMounted(async () => {
  try {
    const storyId = route.params.id
    const all = await get<any[]>('/stories?type=novel')
    const story = all?.find(s => s.id === Number(storyId)) || null
    
    if (story) {
      form.title = story.title
      form.author = story.author || ''
      form.publisher = story.publisher || ''
      form.status = story.status
      form.description = story.description || ''
      form.coverUrl = story.coverUrl || ''
      form.genres = story.genres ? story.genres.map((g: any) => g.id) : []
      
      useHead({ title: `Sửa: ${story.title} - Admin` })
    } else {
      alert('Không tìm thấy tiểu thuyết này!')
      router.push('/admin/novels')
    }
  } catch (err) {
    console.error('Error fetching novel details:', err)
  } finally {
    loading.value = false
  }
})

const triggerUpload = () => fileInput.value?.click()

const handleFileUpload = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const reader = new FileReader()
    reader.onload = (re) => {
      form.coverUrl = re.target?.result as string
    }
    reader.readAsDataURL(target.files[0])
  }
}

const generateSlug = (text: string) => {
  return text.toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, "").replace(/đ/g, "d").replace(/[^a-z0-9\s-]/g, "").trim().replace(/\s+/g, "-")
}

const handleSubmit = async () => {
  if (!form.title || !form.description || form.genres.length === 0 || !form.coverUrl) {
    alert('Vui lòng điền đầy đủ các thông tin bắt buộc (Tên truyện, Thể loại, Mô tả, Ảnh bìa)!')
    return
  }
  if (form.title.length > 50) return alert('Tên truyện tối đa 50 ký tự!')
  if (descriptionLength.value > 700) return alert('Mô tả tối đa 700 ký tự (chữ thô)!')

  isSubmitting.value = true
  try {
    const payload = {
      title: form.title,
      slug: generateSlug(form.title),
      type: 'novel',
      status: form.status,
      description: form.description,
      author: form.author,
      publisher: form.publisher,
      coverUrl: form.coverUrl,
      genres: form.genres.map(id => ({ id }))
    }
    const headers = { Authorization: `Bearer ${authStore.token}` }
    await put(`/admin/stories/${route.params.id}`, payload, { headers })
    alert('Đã cập nhật tiểu thuyết thành công!')
    router.push(`/admin/novels/${route.params.id}`)
  } catch (error) {
    console.error(error)
    alert('Có lỗi xảy ra khi cập nhật tiểu thuyết!')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.admin-create-page { display: flex; flex-direction: column; gap: 24px; }
.page-header { display: flex; align-items: center; gap: 16px; }
.back-link { background: none; border: none; color: #9CA764; cursor: pointer; font-weight: 600; }
.page-header h2 { margin: 0; color: #fff; font-family: 'Montserrat', sans-serif; }

.create-grid { display: grid; grid-template-columns: 1fr 320px; gap: 24px; }
@media (max-width: 1024px) { .create-grid { grid-template-columns: 1fr; } }

.form-section { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 12px; padding: 24px; margin-bottom: 24px; }
.form-section h3 { font-size: 1rem; color: #F1E8C7; margin: 0 0 20px 0; border-bottom: 1px solid rgba(255,255,255,0.05); padding-bottom: 12px; }

.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-group.full { grid-column: span 2; }
.form-group label { font-size: 0.85rem; color: #A8A8B3; font-weight: 600; display: flex; justify-content: space-between; }
.form-group label .required { color: #ef4444; }
.form-group input, .form-group textarea, .form-group select { background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; padding: 12px; color: #fff; font-size: 0.9rem; transition: border-color 0.2s; }
.form-group input:focus, .form-group textarea:focus, .form-group select:focus { border-color: #9CA764; outline: none; }
.char-count { font-size: 0.75rem; color: #5C5C6B; text-align: right; margin-top: 4px; }
.help-text { font-size: 0.75rem; color: #9CA764; margin-top: 4px; }

.radio-group { display: flex; gap: 16px; align-items: center; height: 44px; }
.radio-item { display: flex; align-items: center; gap: 8px; color: #fff; font-size: 0.9rem; cursor: pointer; }

.quill-wrapper {
  background: #fff;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 8px;
  overflow: hidden;
}
:deep(.ql-toolbar.ql-snow) {
  border: none;
  border-bottom: 1px solid #ddd;
  background: #f8f9fa;
}
:deep(.ql-container.ql-snow) {
  border: none;
  min-height: 250px;
}
:deep(.ql-editor) {
  font-size: 1rem;
  color: #333;
  line-height: 1.6;
}
:deep(.ql-editor.ql-blank::before) {
  color: #999;
  font-style: normal;
}
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

:deep(.ql-font-arial) { font-family: 'Arial', sans-serif; }
:deep(.ql-font-roboto) { font-family: 'Roboto', sans-serif; }
:deep(.ql-font-montserrat) { font-family: 'Montserrat', sans-serif; }
:deep(.ql-font-inter) { font-family: 'Inter', sans-serif; }
:deep(.ql-font-times-new-roman) { font-family: 'Times New Roman', serif; }

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

:deep(.ql-snow .ql-stroke) { stroke: #444; }
:deep(.ql-snow .ql-fill) { fill: #444; }
:deep(.ql-snow .ql-picker) { color: #444; }

.char-count.error { color: #ef4444; }

.cover-upload-box { width: 100%; aspect-ratio: 2/3; border: 2px dashed rgba(255,255,255,0.1); border-radius: 8px; display: flex; align-items: center; justify-content: center; overflow: hidden; cursor: pointer; transition: all 0.2s; }
.cover-upload-box:hover { border-color: #9CA764; background: rgba(156,167,100,0.05); }
.upload-placeholder { display: flex; flex-direction: column; align-items: center; gap: 10px; color: #5C5C6B; text-align: center; }
.upload-placeholder .icon { font-size: 2rem; }
.cover-preview { width: 100%; height: 100%; position: relative; }
.cover-preview img { width: 100%; height: 100%; object-fit: cover; }
.remove-cover { position: absolute; top: 10px; right: 10px; width: 24px; height: 24px; border-radius: 50%; background: rgba(239,68,68,0.8); color: #fff; border: none; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 1.2rem; }

.publish-panel { display: flex; flex-direction: column; gap: 12px; }
.btn-publish { background: #9CA764; color: #000; border: none; padding: 14px; border-radius: 8px; font-weight: 700; cursor: pointer; transition: all 0.2s; width: 100%; }
.btn-publish:hover { background: #B8C878; transform: translateY(-2px); box-shadow: 0 4px 15px rgba(156,167,100,0.4); }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 80px; color: #A8A8B3; gap: 16px; }
.spinner { width: 40px; height: 40px; border: 3px solid rgba(156,167,100,0.2); border-top-color: #9CA764; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
