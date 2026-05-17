<template>
  <div class="noise-overlay">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
const authStore = useAuthStore()

onMounted(() => {
  authStore.init()
})

useHead({
  htmlAttrs: { lang: 'vi', 'data-theme': 'dark' },
  bodyAttrs: { class: 'antialiased' },
})
</script>
