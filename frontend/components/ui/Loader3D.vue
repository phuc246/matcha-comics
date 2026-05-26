<template>
  <Transition name="fade-out">
    <div v-if="visible" class="loader-3d-overlay">
      <!-- Glow Spotlights -->
      <div class="glow-bg-spot primary"></div>
      <div class="glow-bg-spot secondary"></div>

      <div class="loader-3d-content">
        <!-- Depth HUD Ring behind the book -->
        <div class="hud-ring-wrap">
          <div class="hud-ring inner"></div>
          <div class="hud-ring middle"></div>
          <div class="hud-ring outer"></div>
        </div>

        <!-- 3D Book loader -->
        <div class="book-3d-scene">
          <div class="book-3d">
            <!-- Left Cover -->
            <div class="book-page book-cover cover-left">
              <div class="cover-design">
                <span class="cover-design-logo">🍵</span>
              </div>
            </div>

            <!-- Flipping middle pages -->
            <div class="book-page book-sheet sheet-1">
              <div class="sheet-content left-content"></div>
              <div class="sheet-content right-content"></div>
            </div>
            <div class="book-page book-sheet sheet-2">
              <div class="sheet-content left-content"></div>
              <div class="sheet-content right-content"></div>
            </div>
            <div class="book-page book-sheet sheet-3">
              <div class="sheet-content left-content"></div>
              <div class="sheet-content right-content"></div>
            </div>

            <!-- Right Cover -->
            <div class="book-page book-cover cover-right">
              <div class="cover-design">
                <span class="cover-design-logo">🍵</span>
              </div>
            </div>

            <!-- Spine -->
            <div class="book-spine"></div>
          </div>
        </div>

        <!-- Staggered loading text -->
        <div class="loading-text-container">
          <div class="loader-brand">
            Matcha<span class="gold-text">Comic</span>
          </div>
          <div class="loading-bar-wrap">
            <div class="loading-bar-active"></div>
          </div>
          <div class="loading-subtitle">
            <span v-for="(char, idx) in subtitleChars" :key="idx" :style="{ animationDelay: `${idx * 0.08}s` }">
              {{ char === ' ' ? '\u00A0' : char }}
            </span>
          </div>
        </div>
      </div>

      <!-- Floating particle system -->
      <div class="ambient-particles">
        <div v-for="n in 12" :key="n" class="particle" :style="particleStyles(n)"></div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

const props = withDefaults(defineProps<{
  visible?: boolean
}>(), {
  visible: true
})

const subtitle = 'Đang tải trang...'
const subtitleChars = computed(() => subtitle.split(''))

const particleStyles = (n: number) => {
  const left = Math.random() * 100
  const size = Math.random() * 6 + 3
  const delay = Math.random() * 4
  const duration = Math.random() * 5 + 4
  return {
    left: `${left}%`,
    width: `${size}px`,
    height: `${size}px`,
    animationDelay: `${delay}s`,
    animationDuration: `${duration}s`,
  }
}
</script>

<style scoped>
.loader-3d-overlay {
  position: fixed;
  inset: 0;
  background: #090a0f;
  z-index: 99999;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  perspective: 1500px;
}

/* Glow Spotlights in Background */
.glow-bg-spot {
  position: absolute;
  border-radius: 50%;
  filter: blur(140px);
  pointer-events: none;
  opacity: 0.35;
  mix-blend-mode: screen;
  animation: pulseGlow 6s ease-in-out infinite alternate;
}

.glow-bg-spot.primary {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(156, 167, 100, 0.4) 0%, transparent 70%);
  top: 15%;
  left: 20%;
}

.glow-bg-spot.secondary {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(212, 175, 55, 0.3) 0%, transparent 70%);
  bottom: 15%;
  right: 20%;
  animation-delay: -3s;
}

@keyframes pulseGlow {
  0% { transform: scale(1) translate(0, 0); opacity: 0.25; }
  100% { transform: scale(1.2) translate(30px, -30px); opacity: 0.45; }
}

.loader-3d-content {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

/* HUD Circle Rings in 3D Perspective behind the book */
.hud-ring-wrap {
  position: absolute;
  width: 400px;
  height: 400px;
  transform: rotateX(72deg) rotateY(0deg) rotateZ(0deg) translateZ(-80px);
  transform-style: preserve-3d;
  pointer-events: none;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.45;
}

.hud-ring {
  position: absolute;
  border: 1px dashed rgba(156, 167, 100, 0.3);
  border-radius: 50%;
  animation: spinRing 25s linear infinite;
  box-shadow: 0 0 20px rgba(156, 167, 100, 0.05);
}

.hud-ring.inner {
  width: 180px;
  height: 180px;
  border-style: double;
  border-width: 3px;
  border-color: rgba(212, 175, 55, 0.25);
  animation-duration: 15s;
  animation-direction: reverse;
}

.hud-ring.middle {
  width: 280px;
  height: 280px;
  border-style: dotted;
  border-width: 2px;
  animation-duration: 35s;
}

.hud-ring.outer {
  width: 380px;
  height: 380px;
  border-color: rgba(156, 167, 100, 0.15);
  animation-duration: 50s;
  animation-direction: reverse;
}

@keyframes spinRing {
  from { transform: rotateZ(0deg); }
  to { transform: rotateZ(360deg); }
}

/* 3D Book Scene */
.book-3d-scene {
  width: 200px;
  height: 150px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 40px;
}

.book-3d {
  width: 140px;
  height: 96px;
  position: relative;
  transform-style: preserve-3d;
  transform: rotateX(28deg) rotateY(-22deg) rotateZ(0deg);
  animation: bookFloat 4s ease-in-out infinite alternate;
}

@keyframes bookFloat {
  0% { transform: rotateX(28deg) rotateY(-22deg) translateY(0px) rotateZ(-2deg); }
  100% { transform: rotateX(32deg) rotateY(-18deg) translateY(-14px) rotateZ(2deg); }
}

.book-page {
  position: absolute;
  width: 70px;
  height: 96px;
  top: 0;
  transform-style: preserve-3d;
  backface-visibility: hidden;
  border-radius: 0 4px 4px 0;
}

/* Covers styling */
.book-cover {
  background: linear-gradient(135deg, #7d874c 0%, #5d6635 100%);
  border: 1px solid rgba(156, 167, 100, 0.5);
  box-shadow: 0 4px 15px rgba(0,0,0,0.5);
}

.cover-left {
  left: 0;
  transform-origin: right center;
  transform: rotateY(-180deg);
  border-radius: 4px 0 0 4px;
}

.cover-right {
  left: 70px;
  transform-origin: left center;
  transform: rotateY(0deg);
  border-radius: 0 4px 4px 0;
}

.cover-design {
  position: absolute;
  inset: 4px;
  border: 1px solid rgba(212, 175, 55, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: inherit;
  background: rgba(0, 0, 0, 0.15);
}

.cover-design-logo {
  font-size: 1.5rem;
  filter: drop-shadow(0 2px 5px rgba(0,0,0,0.4));
  animation: coverLogoPulse 2s ease-in-out infinite alternate;
}

@keyframes coverLogoPulse {
  from { transform: scale(0.95); opacity: 0.8; }
  to { transform: scale(1.05); opacity: 1; }
}

/* Middle flipping sheets */
.book-sheet {
  left: 70px;
  width: 68px;
  height: 92px;
  top: 2px;
  transform-origin: left center;
  background: #fbfbf9;
  border-radius: 0 3px 3px 0;
  box-shadow: inset -2px 0 8px rgba(0,0,0,0.06), 0 2px 5px rgba(0,0,0,0.1);
  animation: flipSheet 2.4s cubic-bezier(0.445, 0.05, 0.55, 0.95) infinite;
}

.sheet-1 { animation-delay: 0s; z-index: 5; }
.sheet-2 { animation-delay: 0.4s; z-index: 4; }
.sheet-3 { animation-delay: 0.8s; z-index: 3; }

.sheet-content {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(rgba(156, 167, 100, 0.08) 1px, transparent 1px);
  background-size: 100% 6px;
  padding: 10px;
}

.sheet-content::before {
  content: '';
  display: block;
  width: 70%;
  height: 4px;
  background: rgba(156, 167, 100, 0.25);
  border-radius: 2px;
  margin-bottom: 8px;
}

.left-content {
  backface-visibility: hidden;
}

.right-content {
  transform: rotateY(180deg);
}

@keyframes flipSheet {
  0% {
    transform: rotateY(0deg);
    z-index: 5;
  }
  50% {
    transform: rotateY(-180deg);
    z-index: 5;
  }
  100% {
    transform: rotateY(-180deg);
    z-index: 1;
  }
}

/* Spine Connection */
.book-spine {
  position: absolute;
  width: 8px;
  height: 96px;
  left: 66px;
  top: 0;
  background: #464d26;
  transform: translateZ(1px);
  box-shadow: 0 0 10px rgba(0,0,0,0.5);
  border-radius: 1px;
}

/* Brand & Loading Info */
.loading-text-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  width: 250px;
}

.loader-brand {
  font-family: 'Montserrat', sans-serif;
  font-size: 1.6rem;
  font-weight: 900;
  letter-spacing: 1px;
  color: #fff;
  text-shadow: 0 0 15px rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  gap: 2px;
}

.gold-text {
  color: transparent;
  background: linear-gradient(135deg, #d4af37, #9CA764);
  -webkit-background-clip: text;
  background-clip: text;
}

.loading-bar-wrap {
  width: 140px;
  height: 2px;
  background: rgba(255,255,255,0.06);
  border-radius: 10px;
  overflow: hidden;
  position: relative;
}

.loading-bar-active {
  position: absolute;
  top: 0;
  left: 0;
  width: 60px;
  height: 100%;
  background: linear-gradient(90deg, transparent, #9CA764, #d4af37, #9CA764, transparent);
  border-radius: inherit;
  animation: shimmerBar 1.8s infinite ease-in-out;
}

@keyframes shimmerBar {
  0% { left: -60px; }
  100% { left: 140px; }
}

.loading-subtitle {
  font-size: 0.8rem;
  color: #a8a8b3;
  font-weight: 500;
  letter-spacing: 3px;
  display: flex;
  margin-left: 3px;
}

.loading-subtitle span {
  display: inline-block;
  animation: letterFloat 1.6s ease-in-out infinite;
}

@keyframes letterFloat {
  0%, 100% { transform: translateY(0); filter: drop-shadow(0 0 0 transparent); color: #a8a8b3; }
  50% { transform: translateY(-4px); filter: drop-shadow(0 0 8px rgba(156,167,100,0.6)); color: #9CA764; }
}

/* Background Ambient Particles */
.ambient-particles {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 1;
}

.particle {
  position: absolute;
  bottom: -20px;
  background: rgba(156, 167, 100, 0.45);
  box-shadow: 0 0 10px rgba(156, 167, 100, 0.6);
  border-radius: 50%;
  animation: riseParticle 6s linear infinite;
  opacity: 0;
}

@keyframes riseParticle {
  0% {
    transform: translateY(0) rotate(0deg) scale(0.8);
    opacity: 0;
  }
  15% {
    opacity: 0.6;
  }
  85% {
    opacity: 0.6;
  }
  100% {
    transform: translateY(-110vh) rotate(360deg) scale(1.2);
    opacity: 0;
  }
}

/* Vue Transitions */
.fade-out-leave-active {
  transition: opacity 0.8s cubic-bezier(0.25, 1, 0.5, 1), transform 0.8s cubic-bezier(0.25, 1, 0.5, 1);
}

.fade-out-leave-to {
  opacity: 0;
  transform: scale(1.04);
}
</style>
