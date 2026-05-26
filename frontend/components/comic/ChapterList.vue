<template>
  <div class="chapters-section">
    <div class="chapters-header">
      <h3>{{ title || 'Danh sách chương' }}</h3>
    </div>
    
    <div class="chapters-grid-container" v-if="chapters && chapters.length > 0">
      <NuxtLink
        v-for="chap in [...chapters].reverse()"
        :key="chap.id"
        :to="getChapterLink(chap)"
        class="chapter-item"
      >
        <span class="chapter-name">Chương {{ chap.number }}: {{ chap.title || '' }}</span>
        
        <div class="chapter-right-meta">
          <!-- Like/Heart button -->
          <button 
            class="like-btn" 
            :class="{ 'is-liked': isChapterLiked(chap.id) }" 
            @click.stop.prevent="toggleLikeChapter(chap)"
          >
            <span class="heart-icon">{{ isChapterLiked(chap.id) ? '❤️' : '🤍' }}</span>
            <span class="like-count">{{ chap.likes || 0 }}</span>
          </button>
          
          <span class="chapter-time">{{ timeAgo(chap.createdAt) }}</span>
        </div>
      </NuxtLink>
    </div>
    
    <div v-else class="empty-chapters">
      Chưa có chương nào được đăng.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const props = defineProps<{
  chapters: any[]
  storySlug: string
  type: 'comic' | 'novel'
  title?: string
}>()

const emit = defineEmits<{
  (e: 'like-toggled'): void
}>()

const likedChapters = ref<number[]>([])
const { post } = useApi()

const getChapterLink = (chap: any) => {
  const prefix = props.type === 'novel' ? 'truyen-chu' : 'truyen-tranh'
  return `/${prefix}/${props.storySlug}/chapter-${chap.number}`
}

const isChapterLiked = (chapterId: number) => {
  return likedChapters.value.includes(chapterId)
}

const toggleLikeChapter = async (chap: any) => {
  const isLiked = isChapterLiked(chap.id)
  const endpoint = `/chapters/${chap.id}/${isLiked ? 'unlike' : 'like'}`
  try {
    const res = await post<any>(endpoint, {})
    if (res) {
      // Update local count
      chap.likes = res.likes
      
      // Update state
      if (isLiked) {
        likedChapters.value = likedChapters.value.filter(id => id !== chap.id)
      } else {
        likedChapters.value.push(chap.id)
      }
      localStorage.setItem('liked_chapters', JSON.stringify(likedChapters.value))
      
      // Emit event so parent can recalculate total hearts
      emit('like-toggled')
    }
  } catch (error) {
    console.error('Error toggling chapter like in component:', error)
  }
}

const timeAgo = (dateStr?: string) => {
  if (!dateStr) return 'Vừa cập nhật'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  if (days > 7) return date.toLocaleDateString('vi-VN')
  if (days > 0) return `${days} ngày trước`
  if (hours > 0) return `${hours} giờ trước`
  if (minutes > 0) return `${minutes} phút trước`
  return 'Vừa xong'
}

onMounted(() => {
  const saved = localStorage.getItem('liked_chapters')
  if (saved) {
    try {
      likedChapters.value = JSON.parse(saved)
    } catch {}
  }
})
</script>

<style scoped>
.chapters-section {
  padding: 30px;
}

.chapters-header {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chapters-header h3 {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--text-primary);
}

.chapters-grid-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 8px;
}

.chapters-grid-container::-webkit-scrollbar {
  width: 6px;
}

.chapters-grid-container::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

.chapters-grid-container::-webkit-scrollbar-thumb {
  background: var(--border-active);
  border-radius: 4px;
}

.chapters-grid-container::-webkit-scrollbar-thumb:hover {
  background: var(--accent-primary);
}

.chapter-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  transition: var(--transition-fast);
}

.chapter-item:hover {
  background: rgba(156, 167, 100, 0.08);
  border-color: var(--accent-primary);
  transform: translateX(4px);
}

.chapter-name {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.chapter-item:hover .chapter-name {
  color: var(--accent-primary);
}

.chapter-right-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.like-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  padding: 4px 10px;
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.8rem;
}

.like-btn:hover {
  background: rgba(239, 68, 68, 0.08);
  border-color: rgba(239, 68, 68, 0.3);
  color: #ef4444;
}

.like-btn.is-liked {
  background: rgba(239, 68, 68, 0.12);
  border-color: rgba(239, 68, 68, 0.4);
  color: #ef4444;
  font-weight: 700;
}

.heart-icon {
  font-size: 0.85rem;
  transition: transform 0.2s ease;
}

.like-btn:active .heart-icon {
  transform: scale(1.3);
}

.chapter-time {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.empty-chapters {
  text-align: center;
  color: var(--text-muted);
  padding: 20px;
}

@media (max-width: 600px) {
  .chapters-section {
    padding: 15px;
  }
  .chapter-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    padding: 12px;
  }
  .chapter-right-meta {
    width: 100%;
    justify-content: space-between;
  }
}
</style>
