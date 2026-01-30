<script setup>
import { onMounted, ref } from 'vue';
import apiClient from '@/api/client'
import BaseLayout from "@/components/layout/BaseLayout.vue";

const loading = ref(false);
const wbStatus = ref({
  active: null,
  message: '',
  lastCheck: null
});

const checkApiStatus = async () => {
  loading.value = true;
  wbStatus.value.message = 'Проверка...';

  try {
    const response = await apiClient.get('/profile/apikeys/status');
    const data = response.data;

    wbStatus.value = {
      active: data.wildberries.active,
      message: data.wildberries.message,
      lastCheck: new Date()
    };

  } catch (error) {
    console.error('Ошибка проверки:', error);
    wbStatus.value = {
      active: false,
      message: `Ошибка: ${error.message}`,
      lastCheck: new Date()
    };
  } finally {
    loading.value = false;
  }
};

const formatTime = (date) => {
  if (!date) return '';
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

onMounted(() => {
  checkApiStatus();
});
</script>

<template>
  <BaseLayout>
    <template #title-icon>
      <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
        <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
      </svg>
    </template>
    <template #title>Мои API ключи</template>

    <div class="api-container">
      <div class="api-card">
        <div class="api-header">
          <div class="api-title">
            <svg class="api-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path>
            </svg>
            <h2>Wildberries API</h2>
          </div>

          <transition name="fade">
            <div v-if="wbStatus.active === null" class="status-badge status-checking">
              <svg class="status-icon spinning" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
              </svg>
              <span>{{ wbStatus.message || 'Проверка...' }}</span>
            </div>
            <div v-else-if="wbStatus.active" class="status-badge status-active">
              <svg class="status-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                <polyline points="22 4 12 14.01 9 11.01"></polyline>
              </svg>
              <span>Активен</span>
            </div>
            <div v-else class="status-badge status-inactive">
              <svg class="status-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
              <span>Недействителен</span>
            </div>
          </transition>
        </div>

        <div class="api-body">
          <div v-if="wbStatus.message" class="api-message">
            <svg class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
            <p>{{ wbStatus.message }}</p>
          </div>

          <div class="api-actions">
            <button
                @click="checkApiStatus"
                :disabled="loading"
                class="btn-refresh"
            >
              <svg
                  :class="['btn-icon', { 'spinning': loading }]"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
              >
                <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"></path>
              </svg>
              {{ loading ? 'Проверка...' : 'Перепроверить' }}
            </button>

            <div v-if="wbStatus.lastCheck" class="last-check">
              <svg class="clock-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
              <span>Последняя проверка: {{ formatTime(wbStatus.lastCheck) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Можно добавить другие API в будущем -->
      <div class="api-card api-card-disabled">
        <div class="api-header">
          <div class="api-title">
            <svg class="api-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="14" rx="2" ry="2"></rect>
              <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"></path>
            </svg>
            <h2>Ozon API</h2>
          </div>
          <div class="status-badge status-coming">
            <svg class="status-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
            <span>Скоро</span>
          </div>
        </div>
        <div class="api-body">
          <p class="coming-soon-text">Интеграция с Ozon находится в разработке</p>
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<style scoped>
@import '@/assets/css/views/profile/apikeys.css';
</style>