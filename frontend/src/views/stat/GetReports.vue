<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import Sidebar from "../../components/layout/Sidebar.vue";
import Navbar from "../../components/layout/Navbar.vue";
import apiClient from '@/api/client'

// Данные формы
const formData = ref({
  dateFrom: '2026-01-08',
  dateTo: '2026-01-15'
});

// Состояния
const loading = ref(false);
const reports = ref([]);
const pollingIntervals = ref({});
const error = ref('');
const success = ref('');
const refreshInterval = ref(null);

// Преобразование статуса
const getStatusText = (statusCode) => {
  switch(statusCode) {
    case 0: return 'В обработке';
    case 1: return 'Готово';
    case 2: return 'Ошибка';
    default: return 'Неизвестно';
  }
};

// Получение иконки для статуса
const getStatusIcon = (statusCode) => {
  switch(statusCode) {
    case 0: // В обработке
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"></path>
              </svg>`;
    case 1: // Готово
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>`;
    case 2: // Ошибка
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>`;
    default:
      return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>`;
  }
};

// Получение стиля для статуса
const getStatusStyle = (statusCode) => {
  switch(statusCode) {
    case 0: // В обработке
      return 'status-processing';
    case 1: // Готово
      return 'status-success';
    case 2: // Ошибка
      return 'status-error';
    default:
      return 'status-unknown';
  }
};

// Форматирование даты
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// Загрузка отчетов с бэка
const loadReports = async (silent = false) => {
  if (!silent) {
    loading.value = true;
  }
  error.value = '';
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('Требуется авторизация');
    }

    const response = await apiClient.get('/wb/stats');

    if (response.status === 401) {
      throw new Error('Требуется авторизация');
    }

    const data = response.data;

    // Сохраняем текущие статусы опрашиваемых отчетов
    const oldReportsMap = new Map();
    reports.value.forEach(r => {
      oldReportsMap.set(r.id, {
        statusCode: r.statusCode,
        isProcessing: r.isProcessing
      });
    });

    // Обрабатываем новые данные
    reports.value = data.map(r => ({
      id: r.id,
      statusCode: r.status,
      status: getStatusText(r.status),
      period: `${r.date_from} - ${r.date_to}`,
      createdAt: r.created,
      updatedAt: r.updated,
      comment: r.last_error || '',
      isProcessing: r.status === 0
    }));

    // Запускаем опрос для отчетов в обработке
    startPollingForProcessingReports();

    // if (!silent) {
    //   success.value = 'Список отчетов обновлен';
    //   clearMessages();
    // }

  } catch (err) {
    console.error('Ошибка загрузки отчетов:', err);
    if (!silent) {
      error.value = `Ошибка загрузки: ${err.message}`;
    }
    reports.value = [];
  } finally {
    if (!silent) {
      loading.value = false;
    }
  }
};

// Опрос статуса для конкретного отчета
const pollSingleReport = async (reportId) => {
  try {

    const response = await apiClient.get(`/wb/stats/${reportId}`);
    if (response.status === 401) {
      throw new Error('Требуется авторизация');
    }
    if (response.status === 200) {
      const data = response.data;
      return {
        statusCode: data.status,
        comment: data.last_error || ''
      };
    }
  } catch (error) {
    console.error(`Ошибка опроса отчета ${reportId}:`, error);
  }
  return null;
};

// Запуск опроса для отчетов в обработке
const startPollingForProcessingReports = () => {
  // Останавливаем все предыдущие интервалы
  Object.values(pollingIntervals.value).forEach(interval => {
    clearInterval(interval);
  });
  pollingIntervals.value = {};

  // Находим отчеты в обработке (статус 0)
  const processingReports = reports.value.filter(r => r.statusCode === 0);

  processingReports.forEach(report => {
    const reportId = report.id;

    // Запускаем интервал опроса
    const interval = setInterval(async () => {
      const result = await pollSingleReport(reportId);

      if (result) {
        // Находим отчет в локальном списке
        const reportIndex = reports.value.findIndex(r => r.id === reportId);

        if (reportIndex !== -1) {
          // Обновляем статус только если он изменился
          if (result.statusCode !== reports.value[reportIndex].statusCode) {
            reports.value[reportIndex].statusCode = result.statusCode;
            reports.value[reportIndex].status = getStatusText(result.statusCode);
            reports.value[reportIndex].comment = result.comment || reports.value[reportIndex].comment;
            reports.value[reportIndex].updatedAt = new Date().toISOString().slice(0, 19).replace('T', ' ');
            reports.value[reportIndex].isProcessing = result.statusCode === 0;

            // Если статус больше не "В обработке" - останавливаем опрос
            if (result.statusCode !== 0) {
              clearInterval(pollingIntervals.value[reportId]);
              delete pollingIntervals.value[reportId];
            }

            // Принудительно обновляем реактивность
            reports.value = [...reports.value];
          }
        }
      }
    }, 5000); // Опрос каждые 5 секунд

    pollingIntervals.value[reportId] = interval;

    // Останавливаем опрос через 30 минут на всякий случай
    setTimeout(() => {
      if (pollingIntervals.value[reportId]) {
        clearInterval(pollingIntervals.value[reportId]);
        delete pollingIntervals.value[reportId];
      }
    }, 30 * 60 * 1000);
  });
};
// Запуск автоматического обновления всего списка
const startAutoRefresh = () => {
  // Очищаем предыдущий интервал
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value);
  }

  // Запускаем новое обновление каждые 10 секунд
  refreshInterval.value = setInterval(() => {
    // Обновляем только если есть отчеты или страница активна
    if (reports.value.length > 0) {
      loadReports(true); // silent = true - без показа загрузки
      console.log('Автоматическое обновление списка отчетов...');
    }
  }, 10000); // 10 секунд
};
// Создание нового отчета
const requestReport = async () => {
  if (loading.value) return;

  loading.value = true;
  error.value = '';
  success.value = '';

  try {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('Требуется авторизация');
    }

    const dateFrom = formData.value.dateFrom;
    const dateTo = formData.value.dateTo;

    const response = await apiClient.post('/wb/stats', {
      dateFrom: dateFrom,
      dateTo: dateTo
    });

    const result = response.data;

    // Создаем новый отчет
    const newReport = {
      id: result.id,
      statusCode: 0,
      status: 'В обработке',
      period: `${dateFrom} - ${dateTo}`,
      createdAt: new Date().toISOString().slice(0, 19).replace('T', ' '),
      updatedAt: new Date().toISOString().slice(0, 19).replace('T', ' '),
      comment: 'Отчет поставлен в очередь на обработку...',
      isProcessing: true
    };

    // Добавляем в начало списка
    reports.value.unshift(newReport);

    // Принудительно обновляем реактивность
    reports.value = [...reports.value];

    // Запускаем опрос для нового отчета
    setTimeout(() => {
      startPollingForProcessingReports();
    }, 1000);

    success.value = 'Отчет поставлен в очередь на обработку! Статус будет обновляться автоматически.';

  } catch (err) {
    console.error('Ошибка создания отчета:', err);
    error.value = `Ошибка: ${err.message}`;
  } finally {
    loading.value = false;
  }
};

// Очистка сообщений
const clearMessages = () => {
  setTimeout(() => {
    error.value = '';
    success.value = '';
  }, 5000);
};

// Очистка интервалов
const clearAllPolling = () => {
  Object.values(pollingIntervals.value).forEach(interval => {
    clearInterval(interval);
  });
  pollingIntervals.value = {};
};

// Инициализация
onMounted(() => {
  console.log('API URL from env:', import.meta.env.VITE_API_URL);
  console.log('All env vars:', import.meta.env);
  loadReports();
  startAutoRefresh();
});

// Очистка при размонтировании
onUnmounted(() => {
  clearAllPolling();
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
                <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                <polyline points="13 2 13 9 20 9"></polyline>
              </svg>
              Отчеты за период
            </h1>
          </div>
          <div class="col-sm-6 text-right">
            <div v-if="reports.length > 0" class="total-items">
              <svg class="total-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="8" y1="6" x2="21" y2="6"></line>
                <line x1="8" y1="12" x2="21" y2="12"></line>
                <line x1="8" y1="18" x2="21" y2="18"></line>
                <line x1="3" y1="6" x2="3.01" y2="6"></line>
                <line x1="3" y1="12" x2="3.01" y2="12"></line>
                <line x1="3" y1="18" x2="3.01" y2="18"></line>
              </svg>
              <span>Всего отчетов: {{ reports.length }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="reports-container">
        <!-- Сообщения -->
        <transition name="fade">
          <div v-if="error" class="message-box message-error" @click="error = ''">
            <svg class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{{ error }}</span>
          </div>
        </transition>

        <transition name="fade">
          <div v-if="success" class="message-box message-success" @click="success = ''">
            <svg class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>{{ success }}</span>
          </div>
        </transition>

        <!-- Панель создания отчета -->
        <div class="filter-panel">
          <div class="panel-header">
            <svg class="panel-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
            <h3>Создать новый отчет</h3>
          </div>
          <form @submit.prevent="requestReport" class="filter-form">
            <div class="filter-grid">
              <div class="filter-group">
                <label class="filter-label" for="dateFrom">
                  <svg class="filter-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                    <line x1="16" y1="2" x2="16" y2="6"></line>
                    <line x1="8" y1="2" x2="8" y2="6"></line>
                    <line x1="3" y1="10" x2="21" y2="10"></line>
                  </svg>
                  Дата с
                </label>
                <input
                    type="date"
                    v-model="formData.dateFrom"
                    class="filter-input"
                    id="dateFrom"
                    required
                >
              </div>

              <div class="filter-group">
                <label class="filter-label" for="dateTo">
                  <svg class="filter-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                    <line x1="16" y1="2" x2="16" y2="6"></line>
                    <line x1="8" y1="2" x2="8" y2="6"></line>
                    <line x1="3" y1="10" x2="21" y2="10"></line>
                    <polyline points="14 14 12 14 12 17"></polyline>
                  </svg>
                  Дата по
                </label>
                <input
                    type="date"
                    v-model="formData.dateTo"
                    class="filter-input"
                    id="dateTo"
                    required
                >
              </div>

              <div class="filter-group align-bottom">
                <button type="submit" class="filter-btn" :disabled="loading">
                  <svg v-if="!loading" class="filter-btn-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                    <polyline points="7 10 12 15 17 10"></polyline>
                    <line x1="12" y1="15" x2="12" y2="3"></line>
                  </svg>
                  <span class="spinner" v-if="loading"></span>
                  {{ loading ? 'Отправка...' : 'Создать отчет' }}
                </button>
              </div>
            </div>
          </form>
        </div>

        <!-- Индикатор загрузки -->
        <div v-if="loading && reports.length === 0" class="loading-container">
          <div class="loading-spinner"></div>
          <p class="loading-text">Загрузка отчетов...</p>
        </div>

        <!-- Таблица отчетов -->
        <div v-else-if="reports.length > 0" class="reports-table-container">
          <div class="table-header">
            <svg class="table-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
              <polyline points="17 21 17 13 7 13 7 21"></polyline>
              <polyline points="7 3 7 8 15 8"></polyline>
            </svg>
            <h3>История отчетов</h3>
          </div>

          <div class="table-wrapper">
            <table class="reports-table">
              <thead>
              <tr>
                <th class="table-status">Статус</th>
                <th class="table-period">Период</th>
                <th class="table-created">Создан</th>
                <th class="table-updated">Обновлен</th>
                <th class="table-comment">Комментарий</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="report in reports" :key="report.id" class="report-row">
                <td class="table-status">
                  <div class="status-badge" :class="getStatusStyle(report.statusCode)">
                    <span class="status-icon" v-html="getStatusIcon(report.statusCode)"></span>
                    <span class="status-text">{{ report.status }}</span>
                    <span v-if="report.statusCode === 0" class="status-spinner">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
                      </svg>
                    </span>
                  </div>
                </td>
                <td class="table-period">
                  <div class="period-text">{{ report.period }}</div>
                </td>
                <td class="table-created">
                  <div class="date-text">{{ formatDate(report.createdAt) }}</div>
                </td>
                <td class="table-updated">
                  <div class="date-text">{{ formatDate(report.updatedAt) }}</div>
                </td>
                <td class="table-comment">
                  <div class="comment-text" v-if="report.comment">{{ report.comment }}</div>
                  <div v-else class="comment-empty">—</div>
                </td>
              </tr>
              </tbody>
            </table>
          </div>

          <!-- Информация об обновлении -->
          <div class="refresh-info">
            <svg class="refresh-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M23 4v6h-6"></path>
              <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
            </svg>
            <span>Отчеты в обработке обновляются автоматически каждые 5 секунд</span>
          </div>
        </div>

        <!-- Пустое состояние -->
        <div v-else-if="!loading && reports.length === 0" class="empty-state">
          <svg class="empty-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
            <polyline points="10 9 9 9 8 9"></polyline>
          </svg>
          <h3 class="empty-title">Отчеты не найдены</h3>
          <p class="empty-text">Создайте первый отчет, указав период выше</p>
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

.total-items {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
}

.total-icon {
  width: 18px;
  height: 18px;
}

.reports-container {
  //max-width: 1400px;
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
  cursor: pointer;
  transition: transform 0.2s ease;
}

.message-box:hover {
  transform: translateY(-1px);
}

.message-error {
  background: linear-gradient(135deg, #f8d7da 0%, #f5c6cb 100%);
  color: #721c24;
  border: 1px solid #f1b0b7;
  box-shadow: 0 4px 12px rgba(114, 28, 36, 0.1);
}

.message-success {
  background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
  color: #155724;
  border: 1px solid #b1dfbb;
  box-shadow: 0 4px 12px rgba(21, 87, 36, 0.1);
}

.message-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

/* Панель создания отчета */
.filter-panel {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 24px;
  margin-bottom: 24px;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 2px solid #e5e7eb;
}

.panel-icon {
  width: 24px;
  height: 24px;
  color: #4f46e5;
}

.panel-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
}

.filter-form {
  width: 100%;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  align-items: end;
}

@media (min-width: 768px) {
  .filter-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group.align-bottom {
  margin-top: auto;
}

.filter-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.filter-icon {
  width: 16px;
  height: 16px;
  color: #6b7280;
}

.filter-input {
  padding: 12px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 14px;
  transition: all 0.2s ease;
  background: #fafafa;
}

.filter-input:hover {
  border-color: #cbd5e0;
  background: white;
}

.filter-input:focus {
  outline: none;
  border-color: #4f46e5;
  background: white;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.filter-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px 16px;
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

.filter-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.4);
  background: linear-gradient(135deg, #4338ca 0%, #4f46e5 100%);
}

.filter-btn:active:not(:disabled) {
  transform: translateY(0);
}

.filter-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.filter-btn-icon {
  width: 18px;
  height: 18px;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* Загрузка */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 3px solid #e5e7eb;
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

.loading-text {
  color: #6b7280;
  font-size: 16px;
  font-weight: 500;
}

/* Таблица отчетов */
.reports-table-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin-bottom: 24px;
}

.table-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 24px 0;
  margin-bottom: 16px;
}

.table-icon {
  width: 24px;
  height: 24px;
  color: #4f46e5;
}

.table-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
}

.table-wrapper {
  overflow-x: auto;
  padding: 0 24px 24px;
}

.reports-table {
  width: 100%;
  border-collapse: collapse;
  min-width: 1000px;
}

.reports-table thead {
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
  border-bottom: 2px solid #e5e7eb;
}

.reports-table th {
  padding: 16px 12px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.reports-table tbody tr {
  border-bottom: 1px solid #e5e7eb;
  transition: background-color 0.2s ease;
}

.reports-table tbody tr:hover {
  background-color: #f9fafb;
}

.reports-table tbody tr:last-child {
  border-bottom: none;
}

.reports-table td {
  padding: 16px 12px;
  vertical-align: middle;
}

/* Колонки таблицы */
.table-status {
  width: 160px;
  min-width: 160px;
}

.table-period {
  width: 200px;
  min-width: 200px;
}

.table-created,
.table-updated {
  width: 180px;
  min-width: 180px;
}

.table-comment {
  min-width: 300px;
}

/* Стили для содержимого таблицы */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
}

.status-processing {
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
  color: #92400e;
  border: 2px solid #fbbf24;
}

.status-success {
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  color: #065f46;
  border: 2px solid #10b981;
}

.status-error {
  background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%);
  color: #991b1b;
  border: 2px solid #ef4444;
}

.status-unknown {
  background: linear-gradient(135deg, #f3f4f6 0%, #e5e7eb 100%);
  color: #374151;
  border: 2px solid #9ca3af;
}

.status-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
}

.status-icon svg {
  width: 100%;
  height: 100%;
}

.status-text {
  flex: 1;
}

.status-spinner {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

.status-spinner svg {
  width: 100%;
  height: 100%;
}

.period-text {
  font-size: 14px;
  font-weight: 500;
  color: #1a202c;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
}

.comment-text {
  font-size: 13px;
  color: #374151;
  line-height: 1.5;
  max-width: 500px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.comment-empty {
  color: #9ca3af;
  font-style: italic;
  font-size: 13px;
}

/* Информация об обновлении */
.refresh-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px 24px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
  color: #6b7280;
  font-size: 13px;
  font-weight: 500;
}

.refresh-icon {
  width: 16px;
  height: 16px;
  color: #9ca3af;
}

/* Пустое состояние */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  width: 64px;
  height: 64px;
  color: #9ca3af;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-text {
  font-size: 15px;
  color: #6b7280;
  margin: 0;
}

/* Анимации */
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

/* Адаптивность */
@media (max-width: 768px) {
  .filter-grid {
    grid-template-columns: 1fr;
  }

  .reports-table {
    min-width: 800px;
  }

  .table-header {
    padding: 20px 20px 0;
  }

  .table-wrapper {
    padding: 0 20px 20px;
  }

  .refresh-info {
    padding: 16px 20px;
    text-align: center;
    font-size: 12px;
  }
}

@media (max-width: 480px) {
  .status-badge {
    flex-direction: column;
    gap: 4px;
    text-align: center;
    padding: 8px;
  }

  .status-text {
    font-size: 12px;
  }
}
</style>