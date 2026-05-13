<template>
  <div class="book-container">
    <div class="book" :class="[sizeClass]">
      <!-- Front cover -->
      <div class="book-face book-front">
        <img :src="coverUrl || '/images/no-cover.jpg'" :alt="title" class="cover-image" loading="lazy" />
        <div class="cover-lighting" />
      </div>
      
      <!-- Book spine (gáy sách) -->
      <div class="book-face book-spine" :style="{ backgroundColor: spineColor }">
        <span class="spine-text">{{ title }}</span>
        <div class="spine-texture" />
      </div>
      
      <!-- Back cover -->
      <div class="book-face book-back" :style="{ backgroundColor: spineColor }">
        <div class="back-logo">🍵 MatchaComic</div>
      </div>
      
      <!-- Pages edge (cạnh trang sách) -->
      <div class="book-face book-pages-right" />
      <div class="book-face book-pages-top" />
      <div class="book-face book-pages-bottom" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  title: string
  coverUrl?: string
  size?: 'sm' | 'md' | 'lg'
  spineColor?: string
}>(), {
  size: 'md',
  spineColor: '#1e212b', // var(--bg-tertiary)
})

const sizeClass = computed(() => `size-${props.size}`)
</script>

<style scoped>
.book-container {
  perspective: 1200px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px;
}

.book {
  position: relative;
  transform-style: preserve-3d;
  transition: transform 0.6s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.6s ease;
  box-shadow: 0 10px 30px rgba(0,0,0,0.6);
  cursor: pointer;
  background-color: #f5f5f5; /* Core pages color to fill gaps */
}

.book:hover {
  transform: rotateY(-25deg) rotateX(4deg);
  box-shadow: -15px 20px 40px rgba(0,0,0,0.8);
}

/* Sizes */
.size-sm { width: 110px; height: 160px; --thickness: 18px; }
.size-md { width: 140px; height: 210px; --thickness: 24px; }
.size-lg { width: 180px; height: 270px; --thickness: 30px; }

.book-face {
  position: absolute;
  top: 0; left: 0;
  transform-style: preserve-3d;
  backface-visibility: hidden;
}

/* Front Cover - slightly larger to overhang pages */
.book-front {
  width: 103%; height: 102%;
  left: -1.5%; top: -1%;
  transform: translateZ(calc(var(--thickness) / 2));
  border-radius: 2px 4px 4px 2px;
  overflow: hidden;
  background-color: var(--bg-secondary);
}

.cover-image {
  width: 100%; height: 100%;
  object-fit: cover;
}

.cover-lighting {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg, 
    rgba(255,255,255,0.05) 0%, 
    rgba(255,255,255,0) 20%, 
    rgba(0,0,0,0.05) 90%, 
    rgba(0,0,0,0.15) 100%
  );
  pointer-events: none;
}

/* Spine - Left Side */
.book-spine {
  width: var(--thickness); height: 100%;
  left: 0;
  transform: translateX(calc(var(--thickness) / -2)) rotateY(-90deg);
  border-radius: 2px 0 0 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-secondary);
  box-shadow: inset -2px 0 5px rgba(0,0,0,0.4);
}

.spine-texture {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, rgba(255,255,255,0.1), rgba(0,0,0,0.2));
  pointer-events: none;
}

.spine-text {
  writing-mode: vertical-rl;
  text-orientation: mixed;
  color: rgba(255,255,255,0.8);
  font-family: 'Montserrat', sans-serif;
  font-size: 0.65rem;
  font-weight: 700;
  letter-spacing: 1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-height: 80%;
  z-index: 2;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.8);
}

/* Back Cover - slightly larger to overhang pages */
.book-back {
  width: 103%; height: 102%;
  left: -1.5%; top: -1%;
  transform: translateZ(calc(var(--thickness) / -2)) rotateY(180deg);
  border-radius: 4px 2px 2px 4px;
  background-color: var(--bg-secondary);
  box-shadow: inset 4px 0 10px rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-logo {
  font-family: 'Montserrat', sans-serif;
  font-size: 0.8rem;
  color: rgba(255,255,255,0.2);
  font-weight: bold;
}

/* Pages (Right Edge) */
.book-pages-right {
  width: var(--thickness); height: 98%;
  top: 1%;
  right: 0; left: auto;
  transform: translateX(calc(var(--thickness) / 2)) rotateY(90deg);
  background: #e0e0e0;
  background-image: repeating-linear-gradient(0deg, #d0d0d0 0px, #e0e0e0 2px);
  border-radius: 0 4px 4px 0;
}

/* Pages (Top & Bottom Edges) */
.book-pages-top {
  width: 98%; height: var(--thickness);
  top: 0; left: 1%;
  transform: translateY(calc(var(--thickness) / -2)) rotateX(90deg);
  background: #d0d0d0;
  background-image: repeating-linear-gradient(90deg, #c0c0c0 0px, #d0d0d0 2px);
}

.book-pages-bottom {
  width: 98%; height: var(--thickness);
  bottom: 0; top: auto; left: 1%;
  transform: translateY(calc(var(--thickness) / 2)) rotateX(-90deg);
  background: #c0c0c0;
  background-image: repeating-linear-gradient(90deg, #b0b0b0 0px, #c0c0c0 2px);
}
</style>
