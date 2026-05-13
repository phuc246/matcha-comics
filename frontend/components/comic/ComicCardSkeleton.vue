<template>
  <!-- Card skeleton for loading state -->
  <div class="skeleton-wrapper">
    <div class="skeleton-cover skeleton" />
    <div class="skeleton-info">
      <div class="skeleton skeleton-title" />
      <div class="skeleton skeleton-chapter" />
    </div>
  </div>
</template>

<style scoped>
.skeleton-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.skeleton-cover {
  width: 100%;
  aspect-ratio: 2/3;
  border-radius: var(--radius-md);
}
.skeleton-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 0 2px;
}
.skeleton-title {
  height: 14px;
  width: 90%;
  border-radius: 4px;
}
.skeleton-chapter {
  height: 11px;
  width: 60%;
  border-radius: 4px;
}
</style>
