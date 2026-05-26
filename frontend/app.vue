<template>
  <div class="noise-overlay">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <Loader3D :visible="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNuxtApp } from '#app'

const authStore = useAuthStore()
const nuxtApp = useNuxtApp()
const loading = ref(true)

onMounted(() => {
  authStore.init()
  
  // Hide initial landing loader
  setTimeout(() => {
    loading.value = false
  }, 800)

  // Listen to Nuxt route transition events
  nuxtApp.hook('page:start', () => {
    loading.value = true
  })
  
  nuxtApp.hook('page:finish', () => {
    setTimeout(() => {
      loading.value = false
    }, 500)
  })
})

useHead({
  htmlAttrs: { lang: 'vi', 'data-theme': 'dark' },
  bodyAttrs: { class: 'antialiased' },
})
</script>
