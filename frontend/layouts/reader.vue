<template>
  <div class="reader-layout">
    <ReaderHeader :title="title" :chapter="chapter" />
    <main class="reader-main">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title?: string; chapter?: string }>()
</script>

<style scoped>
.reader-layout {
  min-height: 100vh;
  background-color: #08090D;
}
.reader-main {
  padding-top: 56px;
}
</style>
