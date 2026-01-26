<script setup>
import { ref, computed, onMounted } from 'vue';
import Navbar from "../../components/layout/Navbar.vue";
import Sidebar from "../../components/layout/Sidebar.vue";
import apiClient from '@/api/client'

// Реактивные данные формы
const formData = ref({
  name: '',
  username: '',
  password: '',
  email: '',
  phone: '',
  taxes: '',
  wb_key: '',
  pro: ''
});

const loading = ref(false);
const message = ref({ text: '', type: '' });

// Получаем данные пользователя при загрузке
const loadUserData = async () => {
  try {
    const response = await apiClient.get('/auth/me');

    if (response.status === 200) {
      const userData = response.data;
      // Заполняем форму данными пользователя
      formData.value = {
        name: userData.name || '',
        username: userData.username || '',
        password: '', // Пароль всегда пустой
        email: userData.email || '',
        phone: userData.phone || '',
        taxes: userData.taxes || '',
        wb_key: userData.wb_key || '',
        pro: userData.pro || '',
      };
    }
  } catch (error) {
    console.error('Ошибка загрузки данных:', error);
  }
};

// Сохранение формы
const saveProfile = async () => {
  loading.value = true;
  message.value = { text: '', type: '' };

  try {
    const response = await apiClient.post('/profile/update', formData.value);

    const data = response.data;

    if (response.status === 200) {
      message.value = {
        text: 'Данные успешно сохранены!',
        type: 'success'
      };

      // Обновляем данные пользователя
      await loadUserData();

      // Автоочистка сообщения через 3 секунды
      setTimeout(() => {
        message.value = { text: '', type: '' };
      }, 3000);

    } else {
      message.value = {
        text: data.error || 'Ошибка сохранения',
        type: 'error'
      };
    }

  } catch (error) {
    console.error('Ошибка:', error);
    message.value = {
      text: 'Ошибка сети или сервера',
      type: 'error'
    };
  } finally {
    loading.value = false;
  }
};

// PRO статус (основанный на реальных данных пользователя)
const proBadgeClass = computed(() => {
  // Используем данные из formData, которые приходят с сервера
  // Если pro = 1, значит активирован, если 0 - неактивирован
  return formData.value.pro === 1 ? 'pro-badge-active' : 'pro-badge-inactive';
});

const proAccountText = computed(() => {
  return formData.value.pro === 1 ? 'Активирован' : 'Неактивирован';
});

// Загружаем данные при монтировании
onMounted(() => {
  loadUserData();
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
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              Личный кабинет
            </h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="profile-container">
        <!-- Сообщение об успехе/ошибке -->
        <transition name="fade">
          <div v-if="message.text" :class="['message-box', `message-${message.type}`]">
            <svg v-if="message.type === 'success'" class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
            <svg v-else class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{{ message.text }}</span>
          </div>
        </transition>

        <!-- Форма -->
        <form class="profile-form" @submit.prevent="saveProfile">
          <div class="form-section">
            <h2 class="section-title">Основная информация</h2>

            <div class="form-grid">
              <div class="form-group">
                <label class="form-label" for="user-name">
                  <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                    <circle cx="12" cy="7" r="4"></circle>
                  </svg>
                  Имя
                </label>
                <input
                    type="text"
                    id="user-name"
                    class="form-input"
                    v-model="formData.name"
                    placeholder="Введите ваше имя"
                >
              </div>

              <div class="form-group">
                <label class="form-label" for="user-username">
                  <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                    <circle cx="8.5" cy="7" r="4"></circle>
                    <polyline points="17 11 19 13 23 9"></polyline>
                  </svg>
                  Логин
                </label>
                <input
                    type="text"
                    id="user-username"
                    class="form-input"
                    v-model="formData.username"
                    placeholder="Введите логин"
                >
              </div>

              <div class="form-group">
                <label class="form-label" for="user-email">
                  <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                    <polyline points="22,6 12,13 2,6"></polyline>
                  </svg>
                  Email
                </label>
                <input
                    type="email"
                    id="user-email"
                    class="form-input"
                    v-model="formData.email"
                    placeholder="example@mail.com"
                >
              </div>

              <div class="form-group">
                <label class="form-label" for="user-phone">
                  <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
                  </svg>
                  Номер телефона
                </label>
                <input
                    type="tel"
                    id="user-phone"
                    class="form-input"
                    v-model="formData.phone"
                    placeholder="+7 (___) ___-__-__"
                >
              </div>
            </div>
          </div>

          <div class="form-section">
            <h2 class="section-title">Безопасность</h2>

            <div class="form-group">
              <label class="form-label" for="user-password">
                <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
                Смена пароля
              </label>
              <input
                  type="password"
                  id="user-password"
                  class="form-input"
                  v-model="formData.password"
                  placeholder="Оставьте пустым, если не хотите менять"
              >
              <p class="form-hint">Введите новый пароль только если хотите его изменить</p>
            </div>
          </div>

          <div class="form-section">
            <h2 class="section-title">PRO статус и настройки</h2>

            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">
                  <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
                  </svg>
                  PRO-аккаунт
                </label>
                <div :class="['pro-badge', proBadgeClass]">
                  <svg class="pro-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                  </svg>
                  {{ proAccountText }}
                </div>
              </div>

              <div class="form-group">
                <label class="form-label" for="user-taxes">
                  <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="1" x2="12" y2="23"></line>
                    <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
                  </svg>
                  Налоговый %
                </label>
                <input
                    type="number"
                    id="user-taxes"
                    class="form-input"
                    v-model.number="formData.taxes"
                    placeholder="0"
                    min="0"
                    max="100"
                >
              </div>
            </div>
          </div>

          <div class="form-section">
            <h2 class="section-title">Интеграции</h2>

            <div class="form-group">
              <label class="form-label" for="user-wb_key">
                <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path>
                </svg>
                Wildberries API KEY
              </label>
              <input
                  type="text"
                  id="user-wb_key"
                  class="form-input"
                  v-model="formData.wb_key"
                  placeholder="Введите ваш API ключ Wildberries"
              >
              <p class="form-hint">Токен для интеграции с Wildberries</p>
            </div>
          </div>

          <div class="form-actions">
            <button
                type="submit"
                class="btn-primary"
                :disabled="loading"
            >
              <svg v-if="!loading" class="btn-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                <polyline points="17 21 17 13 7 13 7 21"></polyline>
                <polyline points="7 3 7 8 15 8"></polyline>
              </svg>
              <span class="spinner" v-if="loading"></span>
              {{ loading ? 'Сохранение...' : 'Сохранить изменения' }}
            </button>
          </div>
        </form>
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

.profile-container {
  //max-width: 900px;
  margin: 0 auto;
  padding: 0 20px;
}

/* Сообщения */
.message-box {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-radius: 12px;
  margin-bottom: 24px;
  font-weight: 500;
  animation: slideDown 0.3s ease-out;
}

.message-success {
  background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
  color: #155724;
  border: 1px solid #b8dcc5;
  box-shadow: 0 4px 12px rgba(21, 87, 36, 0.1);
}

.message-error {
  background: linear-gradient(135deg, #f8d7da 0%, #f5c6cb 100%);
  color: #721c24;
  border: 1px solid #f1b0b7;
  box-shadow: 0 4px 12px rgba(114, 28, 36, 0.1);
}

.message-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

/* Форма */
.profile-form {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.form-section {
  padding: 32px;
  border-bottom: 1px solid #e5e7eb;
}

.form-section:last-child {
  border-bottom: none;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 24px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid #4f46e5;
  display: inline-block;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #374151;
  font-size: 14px;
}

.label-icon {
  width: 18px;
  height: 18px;
  color: #6b7280;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 15px;
  transition: all 0.2s ease;
  background: #fafafa;
}

.form-input:hover {
  border-color: #cbd5e0;
  background: white;
}

.form-input:focus {
  outline: none;
  border-color: #4f46e5;
  background: white;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.form-hint {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}

/* PRO Badge */
.pro-badge {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 14px;
}

.pro-badge-inactive {
  background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%);
  color: #dc2626;
  border: 2px solid #fca5a5;
}
.pro-badge-active {
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  color: #047857;
  border: 2px solid #6ee7b7;
}

.pro-icon {
  width: 20px;
  height: 20px;
}

/* Actions */
.form-actions {
  padding: 24px 32px;
  background: #f9fafb;
  display: flex;
  justify-content: flex-end;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 32px;
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 70, 229, 0.4);
  background: linear-gradient(135deg, #4338ca 0%, #4f46e5 100%);
}

.btn-primary:active:not(:disabled) {
  transform: translateY(0);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-icon {
  width: 20px;
  height: 20px;
}

/* Spinner */
.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* Animations */
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.fade-enter-active, .fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Responsive */
@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .form-section {
    padding: 24px 20px;
  }

  .form-actions {
    padding: 20px;
  }

  .btn-primary {
    width: 100%;
    justify-content: center;
  }
}
</style>