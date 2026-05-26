<template>
  <div class="rich-text-editor-container">
    <div class="quill-chapter-wrapper">
      <QuillEditor
        :content="modelValue"
        content-type="html"
        theme="snow"
        :placeholder="placeholder || 'Nhập nội dung truyện tại đây...'"
        :toolbar="fullToolbar"
        @update:content="onContentChange"
      />
    </div>

    <!-- Live Premium Metrics Status Bar -->
    <div class="editor-status-bar">
      <div class="status-item">
        <span class="status-label">Số từ ước tính:</span>
        <strong class="status-value">{{ wordCount }} từ</strong>
      </div>
      <div class="status-item">
        <span class="status-label">Dung lượng:</span>
        <strong class="status-value">{{ kbSize }} KB</strong>
      </div>
      <div class="status-item" :class="{ 'warning': charCount > maxLimit }">
        <span class="status-label">Độ dài (chữ thô):</span>
        <strong class="status-value">{{ charCount }} / {{ maxLimit }}</strong>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
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

const props = withDefaults(defineProps<{
  modelValue: string
  placeholder?: string
  maxLimit?: number
}>(), {
  modelValue: '',
  maxLimit: 15000,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'update:wordCount', count: number): void
  (e: 'update:charCount', count: number): void
}>()

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

// ─── Live Computed Stats ──────────────────────────────────────────────
const wordCount = computed(() => {
  const text = props.modelValue.replace(/<[^>]*>/g, ' ').trim()
  return text ? text.split(/\s+/).length : 0
})

const charCount = computed(() => {
  if (!props.modelValue) return 0
  
  // Safely extract plain text from HTML string in modern environment
  if (typeof document !== 'undefined') {
    const temp = document.createElement('div')
    temp.innerHTML = props.modelValue
    return temp.textContent?.length || 0
  }
  
  // Fallback for SSR
  return props.modelValue.replace(/<[^>]*>/g, '').length
})

const kbSize = computed(() => {
  return (props.modelValue.length / 1024).toFixed(1)
})

const onContentChange = (contentHtml: string) => {
  emit('update:modelValue', contentHtml)
  emit('update:wordCount', wordCount.value)
  emit('update:charCount', charCount.value)
}
</script>

<style scoped>
.rich-text-editor-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

/* Quill White Theme Wrapper */
.quill-chapter-wrapper {
  background: #fff;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px 8px 0 0;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

:deep(.ql-toolbar.ql-snow) {
  border: none;
  border-bottom: 1px solid #ddd;
  background: #f8f9fa;
  padding: 8px 12px;
}

:deep(.ql-container.ql-snow) {
  border: none;
  min-height: 500px;
}

:deep(.ql-editor) {
  font-size: 1.15rem;
  color: #333;
  line-height: 1.8;
  font-family: 'Inter', sans-serif;
  padding: 30px 40px;
}

:deep(.ql-editor.ql-blank::before) {
  color: #999;
  font-style: normal;
  left: 40px;
}

:deep(.ql-snow .ql-stroke) {
  stroke: #444;
}

:deep(.ql-snow .ql-fill) {
  fill: #444;
}

:deep(.ql-snow .ql-picker) {
  color: #444;
}

/* Status Bar styling */
.editor-status-bar {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 24px;
  background: #191d26;
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-top: none;
  border-radius: 0 0 8px 8px;
  padding: 10px 20px;
  font-size: 0.85rem;
  color: #a8a8b3;
  margin-top: -8px; /* Stitch to editor border */
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-label {
  color: #5c5c6b;
}

.status-value {
  color: #f1e8c7;
}

.status-item.warning .status-value {
  color: #ef4444;
  font-weight: 700;
}

/* Font styles overrides inside Quill editor */
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="arial"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="arial"]::before) {
  content: 'Arial';
  font-family: 'Arial';
}
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="roboto"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="roboto"]::before) {
  content: 'Roboto';
  font-family: 'Roboto';
}
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="montserrat"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="montserrat"]::before) {
  content: 'Montserrat';
  font-family: 'Montserrat';
}
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="inter"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="inter"]::before) {
  content: 'Inter';
  font-family: 'Inter';
}
:deep(.ql-snow .ql-picker.ql-font .ql-picker-label[data-value="times-new-roman"]::before),
:deep(.ql-snow .ql-picker.ql-font .ql-picker-item[data-value="times-new-roman"]::before) {
  content: 'Times New Roman';
  font-family: 'Times New Roman';
}

:deep(.ql-font-arial) { font-family: 'Arial', sans-serif; }
:deep(.ql-font-roboto) { font-family: 'Roboto', sans-serif; }
:deep(.ql-font-montserrat) { font-family: 'Montserrat', sans-serif; }
:deep(.ql-font-inter) { font-family: 'Inter', sans-serif; }
:deep(.ql-font-times-new-roman) { font-family: 'Times New Roman', serif; }

/* Custom Line Height Dropdown */
:deep(.ql-snow .ql-picker.ql-lineHeight) {
  width: 100px;
}
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-label[data-value]::before),
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value]::before) {
  content: attr(data-value);
}
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-label::before),
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item::before) {
  content: 'Giãn dòng';
}
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1"]::before) { content: '1.0'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1.2"]::before) { content: '1.2'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1.5"]::before) { content: '1.5'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="1.8"]::before) { content: '1.8'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="2"]::before) { content: '2.0'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="2.5"]::before) { content: '2.5'; }
:deep(.ql-snow .ql-picker.ql-lineHeight .ql-picker-item[data-value="3"]::before) { content: '3.0'; }
</style>
