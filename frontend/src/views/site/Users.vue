<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from "../../components/layout/Navbar.vue"
import Sidebar from "../../components/layout/Sidebar.vue"
import apiClient from '@/api/client'
import BaseLayout from "../../components/layout/BaseLayout.vue";
import BasicForm from "../../components/layout/forms/BasicForm.vue";

const route = useRoute()
const router = useRouter()

// Данные
const users = ref([])
const params = ref([])
const loading = ref(false)
const error = ref('')
const successMessage = ref('')

// Пагинация
const currentPage = ref(1)
const pageSize = ref(20)
const totalItems = ref(0)
const totalPages = ref(0)

const DEFAULT_TIMEOUT = import.meta.env.DEFAULT_TIMEOUT;

// Вычисляемое свойство для отображения страниц
const pagesArray = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)

  if (end - start + 1 < maxVisible) {
    start = Math.max(1, end - maxVisible + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

const truncateText = (text, maxLength) => {
  if (!text) return '';
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
};

const confirmDelete = (userId) => {
  if (confirm('Вы уверены, что хотите удалить пользователя?')) {
    saveUpdateUser({userId: userId, actionType: 'del', value: 1})
  }
}

const isRecentLogin = (dateString) => {
  if (!dateString) return false;
  const date = new Date(dateString);
  const now = new Date();
  const diffDays = (now - date) / (1000 * 60 * 60 * 24);
  return diffDays <= 30; // Последние 30 дней
};

// Загрузка карточек товаров
const fetchUsers = async () => {
  if (!localStorage.getItem('token')) {
    router.push('/login')
    return
  }

  try {
    loading.value = true
    error.value = ''

    const token = localStorage.getItem('token')
    if (!token) {
      error.value = 'Необходима авторизация'
      return
    }

    const response = await apiClient.get('/site/users', {
      headers: {
        'Authorization': `Bearer ${token}`
      },
      params: {
        page: currentPage.value,
        pageSize: pageSize.value
      }
    })

    users.value = response.data.data || []

    // Пагинация
    if (response.data.pagination) {
      currentPage.value = response.data.pagination.current_page
      totalItems.value = response.data.pagination.total_items
      totalPages.value = response.data.pagination.total_pages
    }

    if (users.value.length === 0) {
      error.value = 'Нет данных'
    }

  } catch (err) {
    console.error('Error details:', err)

    // Обработка ошибки 403 (Forbidden)
    if (err.response?.status === 403) {
      // Перенаправляем на страницу с сообщением об ошибке или на главную
      router.push('/access-denied')
      return
    }

    if (err.response?.data?.error) {
      error.value = `Ошибка: ${err.response.data.error}`
    } else if (err.message) {
      error.value = `Ошибка сети: ${err.message}`
    } else {
      error.value = 'Ошибка загрузки данных'
    }
  } finally {
    loading.value = false
  }
}

const formatNumber = (value) => {
  if (value === null || value === undefined || value === '') return '0'
  const num = typeof value === 'string' ? parseFloat(value) : value
  return new Intl.NumberFormat('ru-RU').format(num)
}

// Навигация по страницам
const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchUsers()
}

// Форматирование даты
const formatDate = (dateString) => {
  if (!dateString) return 'Нет данных'
  const date = new Date(dateString)
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  })
}
const requestUserUpdate = async (params) => {
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      showMessage('❌ Требуется авторизация', 'error');
      return;
    }

    console.log('Sending update request:', params);

    const response = await apiClient.post('/user/update', params);

    console.log('Full response:', response);
    console.log('Response data:', response.data);

    // Проверяем успешность по data.success или статусу
    if (response.status === 200 && response.data?.success !== false) {
      const userIndex = users.value.findIndex(p => p.id === params.userId);

      if (userIndex !== -1) {
        if (params.actionType === 'del' && params.value === 1) {
          users.value.splice(userIndex, 1);
          showMessage(`✅ Пользователь удалён`, 'error');
        } else {
          users.value[userIndex][params.actionType] = params.value;
          // Обновляем реактивность
          users.value = [...users.value];
          if (params.value < 1) {
            showMessage(`✅ Пользователь обновлён`, 'error');
          } else {
            showMessage(`✅ Пользователь обновлён`, 'success');
          }
        }
      }
    } else {
      // Сервер вернул ошибку в теле ответа
      showMessage(`❌ ${response.data?.message || 'Ошибка обновления от сервера'}`, 'error');
    }

  } catch (err) {
    console.error('Error updating user data:', err);

    let errorMsg = '❌ Ошибка обновления';
    if (err.response?.data?.error) {
      errorMsg += `: ${err.response.data.error}`;
    } else if (err.message) {
      errorMsg += `: ${err.message}`;
    }

    showMessage(errorMsg, 'error');
  }
}
const saveUpdateUser = (params) => {
  requestUserUpdate(params)
}
// Вспомогательная функция для сообщений
const showMessage = (text, type) => {
  error.value = text
  setTimeout(() => {
    error.value = ''
  }, DEFAULT_TIMEOUT)
}
onMounted(() => {
  if (route.query.page) currentPage.value = parseInt(route.query.page) || 1

  fetchUsers()
})
</script>
<template>
  <BaseLayout>
    <!-- Заголовок страницы -->
    <template #title-icon>
      <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
        <circle cx="12" cy="7" r="4"></circle>
      </svg>
    </template>

    <template #title>Пользователи</template>

    <template #header-right>
      <div class="total-items">
        <svg class="total-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <line x1="8" y1="6" x2="21" y2="6"></line>
          <line x1="8" y1="12" x2="21" y2="12"></line>
          <line x1="8" y1="18" x2="21" y2="18"></line>
        </svg>
        <span>Всего пользователей: {{ formatNumber(totalItems) }}</span>
      </div>
    </template>

    <!-- Основной контент -->
    <BasicForm>
      <!-- Сообщения -->
      <template #messages>
        <div v-if="error" class="message-box message-error">
          <svg class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
          <span>{{ error }}</span>
        </div>

        <div v-if="successMessage" class="message-box message-success">
          <svg class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
          </svg>
          <span>{{ successMessage }}</span>
        </div>
      </template>

      <!-- Ваш основной контент -->
      <div class="users-container">
        <div class="summary-info">
          <svg class="summary-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          <span>Показаны записи <strong>1-{{ users.length }}</strong> из <strong>{{ formatNumber(totalItems) }}</strong></span>
        </div>

        <!-- Таблица пользователей (ваш код) -->
        <div class="users-table-container">
          <div class="users-table-container">
            <div class="table-header">
              <svg class="table-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                <polyline points="17 21 17 13 7 13 7 21"></polyline>
                <polyline points="7 3 7 8 15 8"></polyline>
              </svg>
              <h3>Список пользователей</h3>
            </div>

            <div class="table-wrapper">
              <table class="users-table">
                <thead>
                <tr>
                  <th class="table-id">
                    <a href="/user?sort=id_user" class="sort-link" data-sort="id_user">
                      ID
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-login">
                    <a href="/user?sort=username" class="sort-link" data-sort="username">
                      Логин
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-email">
                    <a href="/user?sort=email" class="sort-link" data-sort="email">
                      Email
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-name">
                    <a href="/user?sort=name" class="sort-link" data-sort="name">
                      Имя
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-blocked">
                    <a href="/user?sort=block" class="sort-link" data-sort="block">
                      Заблокирован
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-pro">
                    <a href="/user?sort=pro" class="sort-link" data-sort="pro">
                      Pro
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-admin">
                    <a href="/user?sort=admin" class="sort-link" data-sort="admin">
                      Админ
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-phone">Телефон</th>
                  <th class="table-taxes">Налоги %</th>
                  <th class="table-created">
                    <a href="/user?sort=created_at" class="sort-link desc" data-sort="created_at">
                      Создан
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="18 15 12 9 6 15"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-updated">
                    <a href="/user?sort=updated_at" class="sort-link" data-sort="updated_at">
                      Обновлен
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-login-time">
                    <a href="/user?sort=last_login" class="sort-link" data-sort="last_login">
                      Последняя авторизация
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-wb-key">
                    <a href="/user?sort=wb_key" class="sort-link" data-sort="wb_key">
                      Ключ WB
                      <svg class="sort-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="6 9 12 15 18 9"></polyline>
                      </svg>
                    </a>
                  </th>
                  <th class="table-actions">Действия</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="user in users" :key="user.id" class="user-row">
                  <td class="table-id">
                    <span class="id-value">{{ user.id }}</span>
                  </td>
                  <td class="table-login">
                    <span class="login-value">{{ user.username }}</span>
                  </td>
                  <td class="table-email">
                    <span class="email-value">{{ user.email }}</span>
                  </td>
                  <td class="table-name">
                    <span class="name-value">{{ user.name }}</span>
                  </td>
                  <td class="table-blocked">
                      <span :class="['status-badge', user.block ? 'status-error' : 'status-success']">
                        <span class="status-icon">
                          <svg v-if="user.block" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"></circle>
                            <line x1="15" y1="9" x2="9" y2="15"></line>
                            <line x1="9" y1="9" x2="15" y2="15"></line>
                          </svg>
                          <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                            <polyline points="22 4 12 14.01 9 11.01"></polyline>
                          </svg>
                        </span>
                        <span class="status-text">{{ user.block ? 'Да' : 'Нет' }}</span>
                      </span>
                  </td>
                  <td class="table-pro">
                      <span :class="['status-badge', user.pro ? 'status-pro' : 'status-basic']">
                        <span class="status-icon">
                          <svg v-if="user.pro" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M20 6L9 17l-5-5"></path>
                          </svg>
                          <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"></circle>
                            <line x1="15" y1="9" x2="9" y2="15"></line>
                            <line x1="9" y1="9" x2="15" y2="15"></line>
                          </svg>
                        </span>
                        <span class="status-text">{{ user.pro ? 'Да' : 'Нет' }}</span>
                      </span>
                  </td>
                  <td class="table-admin">
                      <span :class="['status-badge', user.admin ? 'status-admin' : 'status-user']">
                        <span class="status-icon">
                          <svg v-if="user.admin" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M12 2L2 7l10 5 10-5-10-5z"></path>
                            <path d="M2 17l10 5 10-5"></path>
                            <path d="M2 12l10 5 10-5"></path>
                          </svg>
                          <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                            <circle cx="12" cy="7" r="4"></circle>
                          </svg>
                        </span>
                        <span class="status-text">{{ user.admin ? 'Да' : 'Нет' }}</span>
                      </span>
                  </td>
                  <td class="table-phone">
                    <span class="phone-value" :title="user.phone">
                          {{ truncateText(user.phone, 15) || '—' }}
                      </span>
                  </td>
                  <td class="table-taxes">
                      <span :class="['tax-value', user.taxes > 0 ? 'tax-active' : '']">
                        {{ user.taxes }}%
                      </span>
                  </td>
                  <td class="table-created">
                    <div class="date-container">
                      <svg class="date-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"></circle>
                        <polyline points="12 6 12 12 16 14"></polyline>
                      </svg>
                      <span class="date-value">{{ formatDate(user.created_at) }}</span>
                    </div>
                  </td>
                  <td class="table-updated">
                    <div class="date-container">
                      <svg class="date-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"></circle>
                        <path d="M12 6v6l4 2"></path>
                      </svg>
                      <span class="date-value">{{ formatDate(user.updated_at) }}</span>
                    </div>
                  </td>
                  <td class="table-login-time">
                    <div class="date-container">
                      <svg class="date-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"></circle>
                        <polyline points="12 6 12 12 15 15"></polyline>
                      </svg>
                      <span :class="['date-value', isRecentLogin(user.last_login) ? 'login-recent' : 'login-old']">
                          {{ formatDate(user.last_login) }}
                        </span>
                    </div>
                  </td>
                  <td class="table-wb-key">
                      <span class="wb-key-value" :title="user.wb_key">
                          {{ truncateText(user.wb_key, 15) || '—' }}
                      </span>
                  </td>
                  <td class="table-actions">
                    <div class="action-buttons">
                      <template v-if="user.pro">
                        <button
                            class="action-btn action-danger"
                            @click="saveUpdateUser({userId: user.id, actionType: 'pro', value: 0})"
                            title="Забрать Pro">
                          <svg class="action-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"></circle>
                            <line x1="15" y1="9" x2="9" y2="15"></line>
                            <line x1="9" y1="9" x2="15" y2="15"></line>
                          </svg>
                          Pro - забрать
                        </button>
                      </template>
                      <template v-if="!user.pro">
                        <button
                            class="action-btn action-success"
                            @click="saveUpdateUser({userId: user.id, actionType: 'pro', value: 1})"
                            title="Выдать Pro">
                          <svg class="action-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"></circle>
                            <line x1="15" y1="9" x2="9" y2="15"></line>
                            <line x1="9" y1="9" x2="15" y2="15"></line>
                          </svg>
                          Pro - выдать
                        </button>
                      </template>

                      <template v-if="!user.admin">
                        <button
                            class="action-btn action-success"
                            @click="saveUpdateUser({userId: user.id, actionType: 'admin', value: 1})"
                            title="Сделать администратором">
                          <svg class="action-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                            <polyline points="22 4 12 14.01 9 11.01"></polyline>
                          </svg>
                          Admin - выдать
                        </button>
                      </template>
                      <template v-if="user.admin">
                        <button
                            class="action-btn action-danger"
                            @click="saveUpdateUser({userId: user.id, actionType: 'admin', value: 0})"
                            title="Забрать админ права">
                          <svg class="action-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                            <polyline points="22 4 12 14.01 9 11.01"></polyline>
                          </svg>
                          Admin - забрать
                        </button>
                      </template>

                      <template v-if="!user.block">
                        <button
                            class="action-btn action-danger"
                            @click="saveUpdateUser({userId: user.id, actionType: 'block', value: 1})"
                            title="Заблокировать пользователя"
                        >
                          <svg class="action-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                          </svg>
                          Заблокировать
                        </button>
                      </template>
                      <template v-if="user.block">
                        <button
                            class="action-btn action-success"
                            @click="saveUpdateUser({userId: user.id, actionType: 'block', value: 0})"
                            title="Разблокиовать пользователя"
                        >
                          <svg class="action-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                          </svg>
                          Разблокиовать
                        </button>
                      </template>

                      <button
                          class="action-btn action-danger"
                          @click="confirmDelete(user.id)"
                          title="Удалить пользователя"
                      >
                        <svg class="action-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <polyline points="3 6 5 6 21 6"></polyline>
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                          <line x1="10" y1="11" x2="10" y2="17"></line>
                          <line x1="14" y1="11" x2="14" y2="17"></line>
                        </svg>
                        Удалить
                      </button>
                    </div>
                  </td>
                </tr>
                </tbody>
              </table>
            </div>
            <!-- Пагинация -->
            <div v-if="totalPages > 1" class="pagination">
              <button
                  class="pagination-btn"
                  @click="goToPage(currentPage - 1)"
                  :disabled="currentPage === 1"
              >
                <svg class="pagination-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="15 18 9 12 15 6"></polyline>
                </svg>
              </button>

              <div class="pagination-pages">
                <button
                    v-for="page in pagesArray"
                    :key="page"
                    class="page-btn"
                    :class="{ active: page === currentPage }"
                    @click="goToPage(page)"
                >
                  {{ page }}
                </button>
              </div>

              <button
                  class="pagination-btn"
                  @click="goToPage(currentPage + 1)"
                  :disabled="currentPage === totalPages"
              >
                <svg class="pagination-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9 18 15 12 9 6"></polyline>
                </svg>
              </button>

              <div class="pagination-info">
                Показано {{ users.length }} из {{ formatNumber(totalItems) }} позиций
              </div>
            </div>
          </div>
        </div>
      </div>
    </BasicForm>
  </BaseLayout>
</template>

<style scoped>
@import '@/assets/css/views/site/users.css';
</style>