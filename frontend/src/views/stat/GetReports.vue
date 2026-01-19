<script setup>
import { ref, onMounted } from 'vue';
import Sidebar from "../../components/layout/Sidebar.vue";
import Navbar from "../../components/layout/Navbar.vue";
import Datepicker from 'vue3-datepicker';
import { format } from 'date-fns';

// Данные формы
const formData = ref({
  dateFrom: new Date('2026-01-08'),
  dateTo: new Date('2026-01-15')
});

// Преобразование даты для API
const formatDateForApi = (date) => format(date, 'yyyy-MM-dd');

// Состояния
const loading = ref(false);
const reports = ref([]);
const pollingIntervals = ref({}); // Храним интервалы опроса

// Преобразование статуса
const getStatusText = (statusCode) => {
  switch(statusCode) {
    case 0: return 'В обработке';
    case 1: return 'Готово';
    case 2: return 'Ошибка';
    default: return 'Неизвестно';
  }
};

// Получение стиля для статуса
const getStatusStyle = (statusCode) => {
  switch(statusCode) {
    case 0: // В обработке
      return 'background-color: #ffc107; color: #212529; padding: 5px 10px; border-radius: 4px; font-weight: 600;';
    case 1: // Готово
      return 'background-color: #28a745; color: white; padding: 5px 10px; border-radius: 4px; font-weight: 600;';
    case 2: // Ошибка
      return 'background-color: #dc3545; color: white; padding: 5px 10px; border-radius: 4px; font-weight: 600;';
    default:
      return 'background-color: #6c757d; color: white; padding: 5px 10px; border-radius: 4px; font-weight: 600;';
  }
};

// Загрузка отчетов с бэка
const loadReports = async () => {
  loading.value = true;
  try {
    const token = localStorage.getItem('token');

    const response = await fetch('http://localhost:8080/api/wb/stats', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.status === 401) {
      alert('Требуется авторизация');
      return;
    }

    if (!response.ok) {
      throw new Error(`Ошибка сервера: ${response.status}`);
    }

    const data = await response.json();

    // Обрабатываем данные - СТРОГО по статусам 0,1,2
    reports.value = data.map(r => ({
      id: r.id,
      statusCode: r.status, // Сохраняем код статуса
      status: getStatusText(r.status), // Текст статуса
      period: `${r.date_from} - ${r.date_to}`,
      createdAt: r.created,
      updatedAt: r.updated,
      comment: r.last_error || '',
      // Добавляем флаг для отслеживания в обработке
      isProcessing: r.status === 0
    }));

    // Запускаем опрос для отчетов в обработке
    startPollingForProcessingReports();

  } catch (error) {
    console.error('Ошибка загрузки отчетов:', error);
    reports.value = [];
  } finally {
    loading.value = false;
  }
};

// Опрос статуса для конкретного отчета
const pollSingleReport = async (reportId) => {
  try {
    const token = localStorage.getItem('token');
    const response = await fetch(`http://localhost:8080/api/wb/stats/${reportId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (response.ok) {
      const data = await response.json();
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
    }, 10000); // Опрос каждые 10 секунд

    pollingIntervals.value[reportId] = interval;

    // Останавливаем опрос через 30 минут на всякий случай
    setTimeout(() => {
      if (pollingIntervals.value[reportId]) {
        clearInterval(pollingIntervals.value[reportId]);
        delete pollingIntervals.value[reportId];
      }
    }, 30 * 60 * 1000); // 30 минут
  });
};

// Создание нового отчета
const requestReport = async () => {
  if (loading.value) return;

  loading.value = true;
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      alert('Требуется авторизация');
      return;
    }

    const dateFrom = formatDateForApi(formData.value.dateFrom);
    const dateTo = formatDateForApi(formData.value.dateTo);

    const response = await fetch('http://localhost:8080/api/wb/stats', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        dateFrom: dateFrom,
        dateTo: dateTo
      })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Ошибка создания отчета');
    }

    const result = await response.json();

    // ВАЖНО: Принудительно устанавливаем статус 0, даже если бэкенд вернул 1
    const newReport = {
      id: result.id,
      statusCode: 0, // ВСЕГДА 0 при создании!
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

    alert('Отчет поставлен в очередь на обработку! Статус будет обновляться автоматически.');

  } catch (error) {
    console.error('Ошибка создания отчета:', error);
    alert(`Ошибка: ${error.message}`);
  } finally {
    loading.value = false;
  }
};

// Очистка интервалов при уничтожении компонента
const clearAllPolling = () => {
  Object.values(pollingIntervals.value).forEach(interval => {
    clearInterval(interval);
  });
  pollingIntervals.value = {};
};

// Инициализация
onMounted(() => {
  loadReports();
});

// Очищаем интервалы при размонтировании
import { onUnmounted } from 'vue';
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
        <h1>Отчет за период</h1>
      </div>
    </div>

    <div class="content container-fluid">
      <!-- Форма запроса -->
      <div class="card card-primary">
        <div class="card-header">Запрос нового отчета</div>
        <div class="card-body">
          <form @submit.prevent="requestReport" class="form-horizontal">
            <div class="row">
              <div class="col-md-3">
                <label>Дата с</label>
                <Datepicker v-model="formData.dateFrom" :format="date => format(date, 'yyyy-MM-dd')" />
              </div>
              <div class="col-md-3">
                <label>Дата по</label>
                <Datepicker v-model="formData.dateTo" :format="date => format(date, 'yyyy-MM-dd')" />
              </div>
              <div class="col-md-2" style="margin-top:32px;">
                <button type="submit" class="btn btn-success btn-block" :disabled="loading">
                  <span v-if="loading">Отправка...</span>
                  <span v-else>Запросить</span>
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>

      <!-- Список отчетов -->
      <div class="card">
        <div class="card-header">История отчетов</div>
        <div class="card-body p-0">
          <div v-if="loading" class="text-center p-3">
            <i class="fas fa-spinner fa-spin"></i> Загрузка...
          </div>
          <div v-else>
            <table class="table table-hover table-striped">
              <thead>
              <tr>
                <th>Статус</th>
                <th>Диапазон</th>
                <th>Дата создания</th>
                <th>Дата изменения</th>
                <th>Комментарий</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="report in reports" :key="report.id">
                <td class="text-center align-middle">
                  <span class="badge" :style="getStatusStyle(report.statusCode)">
                    {{ report.status }}
                    <span v-if="report.statusCode === 0" class="ml-1">
                      <i class="fas fa-spinner fa-spin"></i>
                    </span>
                  </span>
                </td>
                <td>{{ report.period }}</td>
                <td>{{ report.createdAt }}</td>
                <td>{{ report.updatedAt }}</td>
                <td>
                  <div v-if="report.comment" v-html="report.comment"></div>
                  <div v-else class="text-muted">-</div>
                </td>
              </tr>
              <tr v-if="reports.length === 0">
                <td colspan="5" class="text-center py-4">Нет данных об отчетах</td>
              </tr>
              </tbody>
            </table>

            <!-- Кнопка ручного обновления -->
            <div class="card-footer text-center">
              <button @click="loadReports" class="btn btn-sm btn-outline-primary">
                <i class="fas fa-sync-alt"></i> Обновить список
              </button>
              <small class="text-muted ml-2">
                Отчеты в обработке обновляются автоматически каждые 10 секунд
              </small>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Ваши существующие стили остаются */
.card {
  box-shadow: 0 0 1px rgba(0,0,0,.125), 0 1px 3px rgba(0,0,0,.2);
  margin-bottom: 20px;
}

.card-primary {
  border-top: 3px solid #ffa500;
}

.card-header {
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  padding: 0.75rem 1.25rem;
}

.card-primary .card-header {
  background-color: #000000 !important;
  color: white;
}

.table th {
  background-color: #f8f9fa;
  font-weight: 600;
  border-top: 1px solid #dee2e6;
}

.table td {
  vertical-align: middle !important;
}

.badge {
  padding: 5px 10px;
  font-size: 12px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-success {
  padding: 8px 20px;
}

/* Стили для Datepicker */
:deep(.dp__input) {
  height: 38px;
  border-radius: 4px;
  border: 1px solid #ced4da;
  width: 100%;
}

:deep(.dp__input:hover) {
  border-color: #80bdff;
}

:deep(.dp__input:focus) {
  border-color: #80bdff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

/* Анимация для статуса "В обработке" */
@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.7; }
  100% { opacity: 1; }
}

.badge-warning-pulse {
  animation: pulse 2s infinite;
}
</style>