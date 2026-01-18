<script setup>
import { ref, onMounted } from 'vue';
import Sidebar from "../../components/layout/Sidebar.vue";
import Navbar from "../../components/layout/Navbar.vue";
// Импортируем Datepicker правильно
import Datepicker from 'vue3-datepicker';
import { parseISO, format } from 'date-fns'; // Добавьте date-fns

// Инициализируем с Date объектами
const formData = ref({
  dateFrom: new Date('2026-01-08'),
  dateTo: new Date('2026-01-15')
});

// Функция для преобразования Date в строку для API
const formatDateForApi = (date) => {
  return format(date, 'yyyy-MM-dd');
};

const loading = ref(false);
const reports = ref([]);
const pagination = ref({
  currentPage: 1,
  totalPages: 1,
  perPage: 10
});

// Функция генерации фейковых данных
const generateFakeReports = () => {
  const statuses = ['Готово', 'В обработке', 'Ошибка'];
  const comments = [
    'Отчет успешно сформирован',
    'Обработка данных...',
    '<span class="text-danger">Ошибка: недостаточно данных за указанный период</span>',
    'Загружено 1250 записей',
    '<span class="text-success">✓ Все проверки пройдены</span>',
    '<span class="text-warning">⚠ Частичные данные - отсутствуют данные за выходные</span>',
    '',
    'Экспорт в Excel завершен',
    '<span class="text-danger">Ошибка соединения с базой данных</span>',
    'Агрегация данных выполнена успешно'
  ];

  const fakeData = [];
  for (let i = 1; i <= 35; i++) {
    const status = statuses[Math.floor(Math.random() * statuses.length)];
    const daysAgo = Math.floor(Math.random() * 30);
    const createdDate = new Date();
    createdDate.setDate(createdDate.getDate() - daysAgo);

    const updatedDate = new Date(createdDate);
    updatedDate.setHours(updatedDate.getHours() + Math.floor(Math.random() * 48));

    const periodStart = new Date();
    periodStart.setDate(periodStart.getDate() - Math.floor(Math.random() * 60) - 30);
    const periodEnd = new Date(periodStart);
    periodEnd.setDate(periodEnd.getDate() + Math.floor(Math.random() * 30) + 1);

    fakeData.push({
      id: i,
      status: status,
      period: `${format(periodStart, 'yyyy-MM-dd')} - ${format(periodEnd, 'yyyy-MM-dd')}`,
      createdAt: format(createdDate, 'yyyy-MM-dd HH:mm:ss'),
      updatedAt: format(updatedDate, 'yyyy-MM-dd HH:mm:ss'),
      comment: comments[Math.floor(Math.random() * comments.length)]
    });
  }

  return fakeData.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
};

const allFakeReports = generateFakeReports();

// Функция загрузки отчетов
const loadReports = async () => {
  loading.value = true;

  // Симулируем задержку загрузки
  await new Promise(resolve => setTimeout(resolve, 500));

  try {
    // Используем фейковые данные вместо реального API
    const startIndex = (pagination.value.currentPage - 1) * pagination.value.perPage;
    const endIndex = startIndex + pagination.value.perPage;

    reports.value = allFakeReports.slice(startIndex, endIndex);
    pagination.value.totalPages = Math.ceil(allFakeReports.length / pagination.value.perPage);

    /* ЗАКОММЕНТИРОВАН РЕАЛЬНЫЙ API ЗАПРОС
    const token = localStorage.getItem('token');
    const response = await fetch(`http://localhost:8080/api/reports?page=${pagination.value.currentPage}&perPage=${pagination.value.perPage}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const data = await response.json();
      reports.value = data.reports || [];
      pagination.value = data.pagination || pagination.value;
    }
    */
  } catch (error) {
    console.error('Ошибка загрузки отчетов:', error);
  } finally {
    loading.value = false;
  }
};

// Функция запроса нового отчета
const requestReport = async () => {
  loading.value = true;

  // Симулируем задержку
  await new Promise(resolve => setTimeout(resolve, 1000));

  try {
    // Добавляем новый фейковый отчет в начало списка
    const newReport = {
      id: allFakeReports.length + 1,
      status: 'В обработке',
      period: `${formatDateForApi(formData.value.dateFrom)} - ${formatDateForApi(formData.value.dateTo)}`,
      createdAt: format(new Date(), 'yyyy-MM-dd HH:mm:ss'),
      updatedAt: format(new Date(), 'yyyy-MM-dd HH:mm:ss'),
      comment: 'Запрос принят, ожидайте обработки...'
    };

    allFakeReports.unshift(newReport);
    alert('Запрос на отчет отправлен!');
    loadReports(); // Обновляем список

    /* ЗАКОММЕНТИРОВАН РЕАЛЬНЫЙ API ЗАПРОС
    const token = localStorage.getItem('token');

    const response = await fetch('http://localhost:8080/api/reports/request', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        dateFrom: formatDateForApi(formData.value.dateFrom),
        dateTo: formatDateForApi(formData.value.dateTo)
      })
    });

    if (response.ok) {
      alert('Запрос на отчет отправлен!');
      loadReports(); // Обновляем список
    } else {
      const error = await response.json();
      alert(`Ошибка: ${error.message || 'Не удалось запросить отчет'}`);
    }
    */

  } catch (error) {
    console.error('Ошибка:', error);
    alert('Ошибка сети или сервера');
  } finally {
    loading.value = false;
  }
};

// Пагинация
const changePage = (page) => {
  if (page >= 1 && page <= pagination.value.totalPages) {
    pagination.value.currentPage = page;
    loadReports();
  }
};

// Инициализация
onMounted(() => {
  loadReports();
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
            <h1 class="m-0">Отчет за период</h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="container-fluid">
        <!-- Внимание блок -->
        <div class="alert alert-warning" style="margin-bottom: 20px;">
          <strong>Внимание!</strong> Из-за высокой нагрузки, временно добавлено ограничение в 1 месяц!<br>
          Только на загрузку новой... Просмотр уже имеющейся статистики, остается таким же!
        </div>

        <!-- Карточка с формой -->
        <div class="card card-primary">
          <div class="card-header">
            <h3 class="card-title">Запрос нового отчета</h3>
          </div>
          <div class="card-body">
            <form @submit.prevent="requestReport" class="form-horizontal">
              <div class="row">
                <div class="col-md-3">
                  <div class="form-group">
                    <label for="dateFrom">Дата с</label>
                    <Datepicker
                        v-model="formData.dateFrom"
                        :format="date => format(date, 'yyyy-MM-dd')"
                        placeholder="Выберите дату"
                        input-class="form-control"
                        :enable-time-picker="false"
                        auto-apply
                    />
                  </div>
                </div>
                <div class="col-md-3">
                  <div class="form-group">
                    <label for="dateTo">Дата по</label>
                    <Datepicker
                        v-model="formData.dateTo"
                        :format="date => format(date, 'yyyy-MM-dd')"
                        placeholder="Выберите дату"
                        input-class="form-control"
                        :enable-time-picker="false"
                        auto-apply
                    />
                  </div>
                </div>
                <div class="col-md-2">
                  <div class="form-group" style="margin-top: 32px;">
                    <button type="submit" class="btn btn-success btn-block" :disabled="loading">
                      <span v-if="loading">
                        <i class="fas fa-spinner fa-spin"></i> Загрузка...
                      </span>
                      <span v-else>
                        <i class="fas fa-file-download"></i> Запросить
                      </span>
                    </button>
                  </div>
                </div>
              </div>
            </form>
          </div>
        </div>

        <!-- Карточка со списком отчетов -->
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">История отчетов</h3>
          </div>
          <div class="card-body p-0">
            <div v-if="loading" class="text-center p-3">
              <i class="fas fa-spinner fa-spin fa-2x"></i>
              <p>Загрузка...</p>
            </div>

            <div v-else>
              <!-- Пагинация сверху -->
              <div class="card-footer clearfix" v-if="reports.length > 0">
                <div class="float-left">
                  <span class="text-muted">
                    Показаны записи <b>{{ (pagination.currentPage - 1) * pagination.perPage + 1 }}-{{ Math.min(pagination.currentPage * pagination.perPage, allFakeReports.length) }}</b> из <b>{{ allFakeReports.length }}</b>
                  </span>
                </div>

                <nav class="float-right" v-if="pagination.totalPages > 1">
                  <ul class="pagination pagination-sm m-0">
                    <li :class="['page-item', { disabled: pagination.currentPage === 1 }]">
                      <a class="page-link" href="#" @click.prevent="changePage(pagination.currentPage - 1)">
                        <i class="fas fa-chevron-left"></i>
                      </a>
                    </li>

                    <li v-for="page in pagination.totalPages" :key="page"
                        :class="['page-item', { active: pagination.currentPage === page }]">
                      <a class="page-link" href="#" @click.prevent="changePage(page)">{{ page }}</a>
                    </li>

                    <li :class="['page-item', { disabled: pagination.currentPage === pagination.totalPages }]">
                      <a class="page-link" href="#" @click.prevent="changePage(pagination.currentPage + 1)">
                        <i class="fas fa-chevron-right"></i>
                      </a>
                    </li>
                  </ul>
                </nav>
              </div>

              <!-- Таблица -->
              <div class="table-responsive">
                <table class="table table-hover table-striped" style="margin-bottom: 0;">
                  <thead class="thead-dark">
                  <tr>
                    <th style="width: 100px;" class="text-center">Статус</th>
                    <th style="width: 200px;" class="text-center">Диапазон</th>
                    <th style="width: 150px;" class="text-center">Дата создания</th>
                    <th style="width: 150px;" class="text-center">Дата изменения</th>
                    <th class="text-center">Комментарий</th>
                  </tr>
                  </thead>
                  <tbody>
                  <tr v-for="report in reports" :key="report.id">
                    <td class="text-center align-middle">
                        <span class="badge"
                              :class="{
                            'badge-success': report.status === 'Готово',
                            'badge-danger': report.status === 'Ошибка',
                            'badge-warning': report.status === 'В обработке'
                          }">
                          {{ report.status }}
                        </span>
                    </td>
                    <td class="text-center align-middle">{{ report.period }}</td>
                    <td class="text-center align-middle">{{ report.createdAt }}</td>
                    <td class="text-center align-middle">{{ report.updatedAt }}</td>
                    <td>
                      <div v-if="report.comment" style="font-size: 12px; line-height: 1.2;">
                        <span v-html="report.comment"></span>
                      </div>
                      <div v-else class="text-muted">
                        <i class="fas fa-minus"></i>
                      </div>
                    </td>
                  </tr>
                  <tr v-if="reports.length === 0">
                    <td colspan="5" class="text-center py-5">
                      <i class="far fa-folder-open fa-3x text-muted mb-3"></i>
                      <p class="text-muted">Нет данных об отчетах</p>
                    </td>
                  </tr>
                  </tbody>
                </table>
              </div>

              <!-- Пагинация снизу -->
              <div class="card-footer clearfix" v-if="pagination.totalPages > 1 && reports.length > 0">
                <div class="float-left">
                  <div class="btn-group">
                    <button type="button" class="btn btn-default btn-sm"
                            v-for="perPage in [10, 25, 50]"
                            :key="perPage"
                            :class="{ active: pagination.perPage === perPage }"
                            @click="pagination.perPage = perPage; pagination.currentPage = 1; loadReports()">
                      {{ perPage }}
                    </button>
                  </div>
                </div>

                <nav class="float-right">
                  <ul class="pagination pagination-sm m-0">
                    <li :class="['page-item', { disabled: pagination.currentPage === 1 }]">
                      <a class="page-link" href="#" @click.prevent="changePage(pagination.currentPage - 1)">
                        <i class="fas fa-chevron-left"></i>
                      </a>
                    </li>

                    <li v-for="page in pagination.totalPages" :key="page"
                        :class="['page-item', { active: pagination.currentPage === page }]">
                      <a class="page-link" href="#" @click.prevent="changePage(page)">{{ page }}</a>
                    </li>

                    <li :class="['page-item', { disabled: pagination.currentPage === pagination.totalPages }]">
                      <a class="page-link" href="#" @click.prevent="changePage(pagination.currentPage + 1)">
                        <i class="fas fa-chevron-right"></i>
                      </a>
                    </li>
                  </ul>
                </nav>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
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
  font-weight: 600;
}

.badge-success {
  background-color: #28a745;
}

.badge-danger {
  background-color: #dc3545;
}

.badge-warning {
  background-color: #ffc107;
  color: #212529;
}

.btn-success {
  padding: 8px 20px;
}

.form-group {
  margin-bottom: 1rem;
}

label {
  font-weight: 600;
  margin-bottom: 0.5rem;
  display: block;
}

/* Стили для Datepicker */
:deep(.dp__input) {
  height: 38px;
  border-radius: 4px;
  border: 1px solid #ced4da;
}

:deep(.dp__input:hover) {
  border-color: #80bdff;
}

:deep(.dp__input:focus) {
  border-color: #80bdff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

:deep(.dp__calendar) {
  border: 1px solid #dee2e6;
  border-radius: 4px;
}

/* Пагинация */
.pagination .page-item.active .page-link {
  background-color: #007bff;
  border-color: #007bff;
}

/* Иконки */
.fa-spinner {
  margin-right: 5px;
}

/* Адаптивность */
@media (max-width: 768px) {
  .card-body .row {
    flex-direction: column;
  }

  .col-md-2 .form-group {
    margin-top: 10px !important;
  }
}
</style>