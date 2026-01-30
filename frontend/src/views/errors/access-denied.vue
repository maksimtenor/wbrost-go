<!-- AccessDenied.vue -->
<template>
  <div class="access-denied-container">
    <div class="error-content">
      <svg class="error-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <h1>Доступ запрещен</h1>
      <p>У вас недостаточно прав для просмотра этой страницы.</p>
      <p>Требуются права администратора.</p>
      <button @click="goHome" class="home-btn">
        Вернуться на главную
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/dashboard')
}
</script>

<style scoped>
@import '@/assets/css/views/errors/access-denied.css';
</style>