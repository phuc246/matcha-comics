<template>
  <div class="admin-create-page">
    <div class="page-header">
      <button class="back-link" @click="router.back()">← Quay lại</button>
      <h2>Thêm Truyện Tranh Mới</h2>
    </div>

    <div class="create-grid">
      <!-- Left: Form Info -->
      <div class="create-form-panel">
        <div class="form-section">
          <h3>Thông tin cơ bản</h3>
          <div class="form-grid">
            <div class="form-group full">
              <label>Tên truyện <span class="required">*</span></label>
              <input type="text" v-model="form.title" maxlength="50" required placeholder="Ví dụ: Solo Leveling (Tối đa 50 ký tự)" />
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

            <div class="form-group" ref="selectContainer">
              <label>Thể loại <span class="required">*</span></label>
              
              <div class="custom-select" :class="{ 'is-open': isGenresDropdownOpen }">
                <div class="select-header" @click="isGenresDropdownOpen = !isGenresDropdownOpen">
                  <div class="selected-genres">
                    <span v-if="form.genres.length === 0">Chọn thể loại...</span>
                    <span v-else v-for="id in form.genres" :key="id" class="genre-tag">
                      {{ availableGenres.find(g => g.id === id)?.name }}
                    </span>
                  </div>
                  <span class="chevron">▼</span>
                </div>
                
                <div class="select-dropdown" v-show="isGenresDropdownOpen">
                  <div class="genres-checkbox-list">
                    <label v-for="genre in availableGenres" :key="genre.id" class="checkbox-item">
                      <input type="checkbox" :value="genre.id" v-model="form.genres" @change="checkGenresLimit" />
                      <span>{{ genre.name }}</span>
                    </label>
                  </div>
                </div>
              </div>

              <small class="help-text">Chọn tối đa 7 thể loại</small>
            </div>
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
                  :toolbar="quillToolbar"
                  placeholder="Tóm tắt nội dung... (Tối đa 700 ký tự)"
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
              <span>Nhấp để tải ảnh bìa</span>
              <input type="file" ref="fileInput" hidden @change="handleFileUpload" />
            </div>
          </div>
        </div>

        <div class="publish-panel">
          <button class="btn-publish" @click="handleSubmit" :disabled="isSubmitting">
            {{ isSubmitting ? 'Đang xử lý...' : 'Đăng truyện' }}
          </button>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
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
  title: 'Thêm Truyện Tranh'
})

const router = useRouter()
const { get, post } = useApi()
const authStore = useAuthStore()
const isSubmitting = ref(false)
const isGenresDropdownOpen = ref(false)
const selectContainer = ref<HTMLElement | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const descTextarea = ref<HTMLTextAreaElement | null>(null)
const form = reactive({
  title: '',
  genres: [],
  author: '',
  publisher: '',
  status: 'ongoing',
  description: '',
  coverUrl: ''
})

const checkGenresLimit = (e: Event) => {
  if (form.genres.length > 7) {
    alert('Chỉ được chọn tối đa 7 thể loại!')
    form.genres.pop() // Revert the last selection
  }
}

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

const availableGenres = ref<any[]>([])

onMounted(async () => {
  document.addEventListener('click', (e) => {
    if (selectContainer.value && !selectContainer.value.contains(e.target as Node)) {
      isGenresDropdownOpen.value = false
    }
  })
  
  try {
    const data = await get<any[]>('/genres')
    availableGenres.value = data || []
  } catch (err) {
    console.error('Error fetching genres:', err)
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
  if (descriptionLength.value > 700) return alert('Mô tả truyện tranh tối đa 700 ký tự!')

  isSubmitting.value = true
  try {
    const payload = {
      title: form.title,
      slug: generateSlug(form.title),
      type: 'comic',
      status: form.status,
      description: form.description,
      author: form.author,
      publisher: form.publisher,
      coverUrl: form.coverUrl,
      genres: form.genres.map(id => ({ id }))
    }
    const headers = { Authorization: `Bearer ${authStore.token}` }
    await post('/admin/stories', payload, { headers })
    alert('Đã thêm truyện tranh thành công!')
    router.push('/admin/comics')
  } catch (error) {
    console.error(error)
    alert('Có lỗi xảy ra khi thêm truyện!')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
/* Resue styles from novels/create.vue */
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

/* Custom Select */
.custom-select { position: relative; width: 100%; }
.select-header { background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; padding: 12px; color: #fff; font-size: 0.9rem; cursor: pointer; display: flex; justify-content: space-between; align-items: center; transition: border-color 0.2s; min-height: 46px; }
.select-header:hover { border-color: #9CA764; }
.custom-select.is-open .select-header { border-color: #9CA764; }
.selected-genres { display: flex; flex-wrap: wrap; gap: 6px; }
.genre-tag { background: rgba(156,167,100,0.2); color: #9CA764; padding: 2px 8px; border-radius: 4px; font-size: 0.8rem; border: 1px solid rgba(156,167,100,0.3); }
.select-dropdown { position: absolute; top: calc(100% + 4px); left: 0; width: 100%; background: #1F2330; border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; padding: 12px; z-index: 10; box-shadow: 0 4px 20px rgba(0,0,0,0.5); max-height: 250px; overflow-y: auto; }
.genres-checkbox-list { display: flex; flex-direction: column; gap: 10px; }
.checkbox-item { display: flex; align-items: center; gap: 8px; color: #fff; cursor: pointer; font-size: 0.9rem; padding: 2px 0; }
.checkbox-item:hover { color: #9CA764; }
.checkbox-item input { accent-color: #9CA764; width: 16px; height: 16px; cursor: pointer; }
.chevron { font-size: 0.7rem; color: #A8A8B3; }

/* Sidebar */
/* Quill Editor Styling */
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
.btn-save { background: #1F2330; color: #A8A8B3; border: 1px solid rgba(255,255,255,0.1); padding: 14px; border-radius: 8px; font-weight: 700; cursor: pointer; }
.btn-publish { background: #9CA764; color: #000; border: none; padding: 14px; border-radius: 8px; font-weight: 700; cursor: pointer; transition: all 0.2s; }
.btn-publish:hover { background: #B8C878; transform: translateY(-2px); box-shadow: 0 4px 15px rgba(156,167,100,0.4); }
</style>
