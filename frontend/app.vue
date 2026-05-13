<template>
  <div class="noise-overlay">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
useHead({
  htmlAttrs: { lang: 'vi', 'data-theme': 'dark' },
  bodyAttrs: { class: 'antialiased' },
})
</script>
