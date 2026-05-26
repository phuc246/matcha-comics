<template>
  <div class="content-protection-system">
    <!-- Floating Security Watermark to deter Screen Recording / Screenshots -->
    <div 
      class="security-watermark-overlay"
      :style="{ 
        transform: `translate3d(${posX}px, ${posY}px, 0)`,
        opacity: watermarkOpacity
      }"
    >
      <div class="watermark-label">
        <span class="lock-icon">🔒</span> {{ watermarkText }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

const props = withDefaults(defineProps<{
  watermarkOpacity?: number
}>(), {
  watermarkOpacity: 0.12
})

const authStore = useAuthStore()

// Generate dynamic secure watermark text based on user session
const watermarkText = computed(() => {
  const identity = authStore.user 
    ? (authStore.user.username || authStore.user.email || 'Thành viên') 
    : 'Khách vãng lai (Guest)'
  return `Bản quyền MatchaComic - Đang xem: ${identity} - Cấm sao chép & quay chụp`
})

// Floating animation logic (DVD-screensaver bouncing bounce box style)
const posX = ref(50)
const posY = ref(100)
let velX = 1.2
let velY = 1.0

let animationFrameId: number | null = null

const updateWatermarkPosition = () => {
  if (typeof window === 'undefined') return

  const width = window.innerWidth - 320 // Subract watermark width estimation
  const height = window.innerHeight - 40 // Subtract watermark height estimation

  posX.value += velX
  posY.value += velY

  // Bounce off horizontal bounds
  if (posX.value >= width) {
    posX.value = width
    velX = -Math.abs(velX)
  } else if (posX.value <= 10) {
    posX.value = 10
    velX = Math.abs(velX)
  }

  // Bounce off vertical bounds
  if (posY.value >= height) {
    posY.value = height
    velY = -Math.abs(velY)
  } else if (posY.value <= 10) {
    posY.value = 10
    velY = Math.abs(velY)
  }

  animationFrameId = requestAnimationFrame(updateWatermarkPosition)
}

// Security event handlers
const blockContext = (e: MouseEvent) => {
  e.preventDefault()
  alert('Nội dung đã được bảo mật bản quyền bởi MatchaComic! Cấm sao chép.')
}

const blockDrag = (e: DragEvent) => {
  if ((e.target as HTMLElement).nodeName === 'IMG') {
    e.preventDefault()
  }
}

const blockKeys = (e: KeyboardEvent) => {
  // Block F12
  if (e.key === 'F12') {
    e.preventDefault()
    return false
  }

  // Block Ctrl+Shift+I, Ctrl+Shift+C, Ctrl+Shift+J (DevTools)
  if (e.ctrlKey && e.shiftKey && (e.key === 'I' || e.key === 'C' || e.key === 'J')) {
    e.preventDefault()
    return false
  }

  // Block Ctrl+U (View Source)
  if (e.ctrlKey && e.key === 'u') {
    e.preventDefault()
    return false
  }

  // Block Ctrl+S (Save Page)
  if (e.ctrlKey && e.key === 's') {
    e.preventDefault()
    return false
  }

  // Block Ctrl+C (Copy)
  if (e.ctrlKey && e.key === 'c') {
    e.preventDefault()
    alert('Bản quyền thuộc về MatchaComic. Tính năng copy đã bị khóa!')
    return false
  }

  // Block Ctrl+P (Print Page)
  if (e.ctrlKey && e.key === 'p') {
    e.preventDefault()
    return false
  }
}

let styleTag: HTMLStyleElement | null = null

onMounted(() => {
  if (typeof window !== 'undefined') {
    // 1. Inject global stylesheet for full block copy and touch styling
    styleTag = document.createElement('style')
    styleTag.innerHTML = `
      body, .chapter-content, .pages-list, img, .comic-reader, .novel-reader {
        user-select: none !important;
        -webkit-user-select: none !important;
        -moz-user-select: none !important;
        -ms-user-select: none !important;
        -webkit-touch-callout: none !important;
      }
      img {
        pointer-events: none !important;
        -webkit-user-drag: none !important;
      }
      @media print {
        body {
          display: none !important;
        }
        .novel-reader, .comic-reader, .reader-container, .reader-main {
          display: none !important;
        }
      }
    `
    document.head.appendChild(styleTag)

    // 2. Add event listeners
    window.addEventListener('contextmenu', blockContext)
    window.addEventListener('dragstart', blockDrag)
    window.addEventListener('keydown', blockKeys)

    // 3. Start watermark animation
    animationFrameId = requestAnimationFrame(updateWatermarkPosition)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    // Clean up style tag
    if (styleTag && styleTag.parentNode) {
      styleTag.parentNode.removeChild(styleTag)
    }

    // Clean up event listeners
    window.removeEventListener('contextmenu', blockContext)
    window.removeEventListener('dragstart', blockDrag)
    window.removeEventListener('keydown', blockKeys)

    // Cancel animation
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId)
    }
  }
})
</script>

<style scoped>
.security-watermark-overlay {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 999999;
  pointer-events: none;
  background: rgba(0, 0, 0, 0.6);
  border: 1px solid rgba(156, 167, 100, 0.4);
  padding: 6px 12px;
  border-radius: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.5), 0 0 10px rgba(156, 167, 100, 0.2);
  white-space: nowrap;
  transition: transform 0.016s linear;
  will-change: transform;
}

.watermark-label {
  font-family: 'Inter', sans-serif;
  font-size: 0.72rem;
  font-weight: 700;
  color: #9CA764;
  letter-spacing: 1px;
  display: flex;
  align-items: center;
  gap: 6px;
  text-shadow: 1px 1px 2px #000;
}

.lock-icon {
  font-size: 0.8rem;
  animation: keyLockPulse 1.5s ease-in-out infinite alternate;
}

@keyframes keyLockPulse {
  from { transform: scale(0.9); opacity: 0.7; }
  to { transform: scale(1.1); opacity: 1; }
}
</style>
