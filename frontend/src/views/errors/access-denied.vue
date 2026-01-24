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
.access-denied-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 20px;
}

.error-content {
  text-align: center;
  max-width: 500px;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.error-icon {
  width: 80px;
  height: 80px;
  color: #ef4444;
  margin-bottom: 20px;
}

h1 {
  font-size: 32px;
  color: #1a202c;
  margin-bottom: 16px;
}

p {
  font-size: 18px;
  color: #6b7280;
  margin-bottom: 8px;
}

.home-btn {
  margin-top: 24px;
  padding: 12px 32px;
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.home-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.4);
}
</style>