<template>
  <div class="novel-reader" :class="`theme-${readerSettings.theme}`">
    <div v-if="loading" class="reader-loader">
      <div class="loader-spinner"></div>
    </div>
    
    <div v-else-if="novel" class="reader-container" :style="{ fontSize: `${readerSettings.fontSize}px`, fontFamily: readerSettings.fontFamily }">
      <div class="reader-header-inline">
        <h1>{{ novel.title }}</h1>
        <h2>Chương {{ currentChapterNum }}: {{ mockChapterTitle }}</h2>
      </div>

      <!-- Settings toolbar -->
      <div class="reader-settings-bar">
        <div class="setting-group">
          <button @click="changeFontSize(-2)" class="setting-btn">A-</button>
          <span class="setting-value">{{ readerSettings.fontSize }}px</span>
          <button @click="changeFontSize(2)" class="setting-btn">A+</button>
        </div>
        <div class="setting-group theme-toggles">
          <button @click="setTheme('dark')" class="theme-btn theme-dark" :class="{ active: readerSettings.theme === 'dark' }" title="Tối"></button>
          <button @click="setTheme('light')" class="theme-btn theme-light" :class="{ active: readerSettings.theme === 'light' }" title="Sáng"></button>
          <button @click="setTheme('sepia')" class="theme-btn theme-sepia" :class="{ active: readerSettings.theme === 'sepia' }" title="Vàng nhạt"></button>
        </div>
      </div>

      <!-- Content -->
      <div class="chapter-content" v-html="mockContent"></div>

      <!-- Chapter Navigation -->
      <div class="chapter-navigation">
        <NuxtLink
          v-if="prevChapter"
          :to="`/truyen-chu/${novel.slug}/chapter-${prevChapter}`"
          class="nav-btn prev-btn"
        >
          Chương trước
        </NuxtLink>
        <button v-else class="nav-btn disabled">Chương trước</button>

        <NuxtLink
          v-if="nextChapter"
          :to="`/truyen-chu/${novel.slug}/chapter-${nextChapter}`"
          class="nav-btn next-btn primary"
        >
          Chương sau
        </NuxtLink>
        <button v-else class="nav-btn disabled">Chương sau</button>
      </div>
    </div>
    
    <div v-else class="error-state">
      <p>Không thể tải chương này.</p>
      <NuxtLink to="/truyen-chu" class="back-link">Về danh sách truyện</NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { mockNovels } from '~/data/mock'

definePageMeta({
  layout: 'reader'
})

const route = useRoute()
const loading = ref(true)
const novel = ref<any>(null)
const currentChapterNum = ref(1)

const mockChapterTitle = ref('Khởi đầu mới')
const mockContent = ref(`
  <p>Đây là nội dung văn bản giả lập cho truyện chữ. Trong tương lai, nội dung này sẽ được tải từ API dưới dạng HTML.</p>
  <p>Ánh nắng chói chang chiếu qua khe cửa, gọi cậu thiếu niên thức giấc sau một giấc mộng dài. Cậu vươn vai, cảm nhận luồng linh khí mỏng manh đang lưu chuyển trong không trung.</p>
  <p>"Hôm nay là ngày tuyển chọn đệ tử ngoại môn của Thái Nhất Tông," cậu lẩm bẩm, ánh mắt lộ ra vẻ kiên định.</p>
  <p>Suốt mười năm qua, cậu đã phải chịu vô số ánh mắt khinh miệt chỉ vì đan điền bẩm sinh bị phế. Nhưng không ai biết rằng, sâu trong thức hải của cậu, một viên ngọc giản thần bí đã lặng lẽ hấp thu tinh hoa nhật nguyệt, chờ đợi khoảnh khắc thức tỉnh.</p>
  <p>Cậu bước xuống giường, siết chặt nắm đấm. Một luồng sức mạnh ấm áp bắt đầu chảy dọc theo kinh mạch, phá vỡ từng tầng gông cùm vốn đã trói buộc cậu bấy lâu nay.</p>
  <p>Thế giới tu tiên rộng lớn và tàn khốc, nhưng từ hôm nay, tên của cậu sẽ bắt đầu được viết lên những trang sử rực rỡ nhất.</p>
  <br/>
  <p>...</p>
`)

const readerSettings = reactive({
  fontSize: 20,
  fontFamily: "'Inter', sans-serif",
  theme: 'dark' // dark | light | sepia
})

const changeFontSize = (delta: number) => {
  const newVal = readerSettings.fontSize + delta
  if (newVal >= 14 && newVal <= 36) {
    readerSettings.fontSize = newVal
  }
}

const setTheme = (theme: string) => {
  readerSettings.theme = theme
}

const prevChapter = computed(() => currentChapterNum.value > 1 ? currentChapterNum.value - 1 : null)
const nextChapter = computed(() => {
  if (novel.value && currentChapterNum.value < novel.value.latestChapter) {
    return currentChapterNum.value + 1
  }
  return null
})

onMounted(async () => {
  const slug = route.params.slug as string
  const chapterParam = route.params.chapter as string
  
  const match = chapterParam.match(/chapter-(\d+)/)
  if (match) {
    currentChapterNum.value = parseInt(match[1])
  }

  await new Promise(r => setTimeout(r, 400))
  novel.value = mockNovels.find(c => c.slug === slug)
  
  loading.value = false
})

useHead(() => ({
  title: novel.value ? `${novel.value.title} - Chương ${currentChapterNum.value}` : 'Đang tải...',
}))
</script>

<style scoped>
.novel-reader {
  width: 100%;
  min-height: 100vh;
  transition: background-color 0.3s, color 0.3s;
}

/* Themes */
.theme-dark {
  background-color: #0d0f14;
  color: #a8a8b3;
}
.theme-light {
  background-color: #f8f9fa;
  color: #333333;
}
.theme-sepia {
  background-color: #f1e7d0;
  color: #5c4b37;
}

.reader-loader { height: 80vh; display: flex; align-items: center; justify-content: center; }
.loader-spinner { width: 40px; height: 40px; border: 3px solid rgba(156,167,100,0.3); border-top-color: #9CA764; border-radius: 50%; animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.reader-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
}

.reader-header-inline {
  text-align: center;
  margin-bottom: 40px;
}
.reader-header-inline h1 {
  font-size: 1.5em;
  font-weight: 800;
  margin-bottom: 8px;
  font-family: 'Montserrat', sans-serif;
}
.reader-header-inline h2 {
  font-size: 1.1em;
  opacity: 0.8;
  font-weight: 500;
}

.reader-settings-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-radius: 12px;
  margin-bottom: 40px;
  font-size: 1rem;
}

.theme-dark .reader-settings-bar { background: rgba(255,255,255,0.05); }
.theme-light .reader-settings-bar { background: rgba(0,0,0,0.05); }
.theme-sepia .reader-settings-bar { background: rgba(0,0,0,0.05); }

.setting-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.setting-btn {
  background: none;
  border: 1px solid currentColor;
  opacity: 0.5;
  border-radius: 4px;
  padding: 4px 12px;
  cursor: pointer;
  color: inherit;
  transition: all 0.2s;
}
.setting-btn:hover { opacity: 1; }

.theme-toggles { gap: 8px; }
.theme-btn {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: transform 0.2s;
}
.theme-btn.active { transform: scale(1.2); border-color: #9CA764; }
.theme-dark { background: #0d0f14; border: 1px solid #333; }
.theme-light { background: #f8f9fa; border: 1px solid #ccc; }
.theme-sepia { background: #f1e7d0; border: 1px solid #d4c5b0; }

.chapter-content {
  line-height: 1.8;
  margin-bottom: 60px;
}

.chapter-content :deep(p) {
  margin-bottom: 1.5em;
}

.chapter-navigation {
  display: flex;
  gap: 16px;
  padding: 40px 0 80px;
  justify-content: center;
}

.nav-btn {
  padding: 14px 28px;
  border-radius: 8px;
  font-weight: 600;
  text-align: center;
  transition: all 0.2s;
  flex: 1;
  max-width: 200px;
  font-size: 1rem;
  border: 1px solid currentColor;
  background: transparent;
  color: inherit;
  opacity: 0.8;
}

.nav-btn:hover:not(.disabled) {
  opacity: 1;
}

.nav-btn.primary {
  background: #9CA764;
  color: #000;
  border: none;
  opacity: 1;
}

.nav-btn.primary:hover:not(.disabled) {
  background: #B8C878;
}

.nav-btn.disabled {
  opacity: 0.3;
  cursor: not-allowed;
}
</style>
