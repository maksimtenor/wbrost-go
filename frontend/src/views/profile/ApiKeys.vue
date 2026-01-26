<script setup>
import Sidebar from "../../components/layout/Sidebar.vue";
import Navbar from "../../components/layout/Navbar.vue";
import { onMounted, ref } from 'vue';
import apiClient from '@/api/client'

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
  <Navbar />
  <Sidebar />

  <div class="content-wrapper">
    <div class="content-header">
      <div class="container-fluid">
        <div class="row mb-2">
          <div class="col-sm-6">
            <h1 class="page-title">
              <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              Мои API ключи
            </h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
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
    </div>
  </div>
</template>

<style scoped>
.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 700;
  color: #1a202c;
  margin: 0;
}

.title-icon {
  width: 32px;
  height: 32px;
  color: #4f46e5;
}

.api-container {
  margin: 0 auto;
  padding: 0 20px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

/* API Card */
.api-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
}

.api-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.api-card-disabled {
  opacity: 0.7;
}

.api-card-disabled:hover {
  transform: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* Header */
.api-header {
  padding: 24px;
  background: linear-gradient(135deg, #f9fafb 0%, #ffffff 100%);
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.api-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.api-icon {
  width: 28px;
  height: 28px;
  color: #4f46e5;
}

.api-title h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
}

/* Status Badges */
.status-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: 600;
  font-size: 14px;
  animation: fadeIn 0.3s ease;
}

.status-icon {
  width: 18px;
  height: 18px;
}

.status-checking {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  color: #d97706;
  border: 2px solid #fbbf24;
}

.status-active {
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  color: #047857;
  border: 2px solid #6ee7b7;
}

.status-inactive {
  background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%);
  color: #dc2626;
  border: 2px solid #fca5a5;
}

.status-coming {
  background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
  color: #4f46e5;
  border: 2px solid #a5b4fc;
}

/* Body */
.api-body {
  padding: 24px;
}

.api-message {
  display: flex;
  align-items: start;
  gap: 12px;
  padding: 12px 16px;
  background: #f0f9ff;
  border: 1px solid #bfdbfe;
  border-radius: 10px;
  margin-bottom: 20px;
}

.message-icon {
  width: 20px;
  height: 20px;
  color: #3b82f6;
  flex-shrink: 0;
  margin-top: 2px;
}

.api-message p {
  margin: 0;
  color: #1e40af;
  font-size: 14px;
  line-height: 1.5;
}

/* Actions */
.api-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.btn-refresh {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(79, 70, 229, 0.3);
}

.btn-refresh:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.4);
  background: linear-gradient(135deg, #4338ca 0%, #4f46e5 100%);
}

.btn-refresh:active:not(:disabled) {
  transform: translateY(0);
}

.btn-refresh:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-icon {
  width: 18px;
  height: 18px;
}

.last-check {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 13px;
}

.clock-icon {
  width: 16px;
  height: 16px;
}

.coming-soon-text {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
  font-style: italic;
}

/* Animations */
.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.fade-enter-active, .fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

/* Responsive */
@media (max-width: 768px) {
  .api-container {
    grid-template-columns: 1fr;
  }

  .api-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .api-actions {
    flex-direction: column;
    align-items: flex-start;
  }

  .btn-refresh {
    width: 100%;
    justify-content: center;
  }
}
</style>