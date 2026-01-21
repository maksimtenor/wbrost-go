<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import Navbar from "../../components/layout/Navbar.vue"
import Sidebar from "../../components/layout/Sidebar.vue"

const route = useRoute()
const router = useRouter()

const statData = ref([])
const loading = ref(false)
const dateFrom = ref('2025-08-01')
const dateTo = ref('2026-01-15')
const error = ref('')
const summary = ref({})
const taxes = ref(5)

// Пагинация
const currentPage = ref(1)
const pageSize = ref(20)
const totalItems = ref(0)
const totalPages = ref(0)

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

const fetchStatDetails = async () => {
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

    const response = await axios.get('/api/stat/details', {
      headers: {
        'Authorization': `Bearer ${token}`
      },
      params: {
        dateFrom: dateFrom.value,
        dateTo: dateTo.value,
        page: currentPage.value,
        pageSize: pageSize.value
      }
    })

    statData.value = response.data.data || []
    summary.value = response.data.summary || {}
    taxes.value = response.data.taxes || 5

    // Пагинация
    if (response.data.pagination) {
      currentPage.value = response.data.pagination.current_page
      totalItems.value = response.data.pagination.total_items
      totalPages.value = response.data.pagination.total_pages
    }

    if (statData.value.length === 0) {
      error.value = 'Нет данных за выбранный период'
    }

  } catch (err) {
    console.error('Error details:', err)
    if (err.response?.data?.error) {
      error.value = `Ошибка: ${err.response.data.error}`
    } else if (err.message) {
      error.value = `Ошибка сети: ${err.message}`
    } else {
      error.value = 'Ошибка загрузки статистики'
    }
  } finally {
    loading.value = false
  }
}

const handleSubmit = (e) => {
  e.preventDefault()
  currentPage.value = 1 // Сбрасываем на первую страницу при новом поиске
  fetchStatDetails()
}

const formatCurrency = (value) => {
  if (value === null || value === undefined || value === '') return '₽ 0.00'
  const num = typeof value === 'string' ? parseFloat(value) : value
  return `₽ ${num.toFixed(2)}`
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
  fetchStatDetails()
}

onMounted(() => {
  if (route.query.dateFrom) dateFrom.value = route.query.dateFrom
  if (route.query.dateTo) dateTo.value = route.query.dateTo
  if (route.query.page) currentPage.value = parseInt(route.query.page) || 1

  fetchStatDetails()
})
</script>

<template>
  <Navbar />
  <Sidebar />

  <div class="content-wrapper">
    <div class="content-header">
      <div class="container-fluid">
        <div class="row mb-2">
          <div class="col-sm-6">
            <h1 class="m-0">Детальная статистика</h1>
          </div>
          <div class="col-sm-6 text-right">
            <small v-if="statData.length > 0" class="text-muted">
              Всего позиций: {{ formatNumber(totalItems) }}
            </small>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="container-fluid">
        <!-- Форма фильтрации -->
        <div class="card">
          <div class="card-body">
            <form @submit="handleSubmit" id="stat-filter-form">
              <div class="row">
                <div class="col-md-3">
                  <div class="form-group">
                    <label for="dateFrom">Дата с</label>
                    <input
                        type="date"
                        v-model="dateFrom"
                        class="form-control"
                        id="dateFrom"
                        required
                    >
                  </div>
                </div>
                <div class="col-md-3">
                  <div class="form-group">
                    <label for="dateTo">Дата по</label>
                    <input
                        type="date"
                        v-model="dateTo"
                        class="form-control"
                        id="dateTo"
                        required
                    >
                  </div>
                </div>
                <div class="col-md-2">
                  <div class="form-group">
                    <label for="pageSize">Показывать</label>
                    <select v-model="pageSize" class="form-control" id="pageSize">
                      <option value="10">10</option>
                      <option value="20">20</option>
                      <option value="50">50</option>
                      <option value="100">100</option>
                    </select>
                  </div>
                </div>
                <div class="col-md-2 align-self-end">
                  <button type="submit" class="btn btn-success" :disabled="loading">
                    {{ loading ? 'Загрузка...' : 'Показать' }}
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>

        <!-- Сообщение об ошибке -->
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <!-- Таблица статистики -->
        <div class="card" v-if="!loading && statData.length > 0">
          <div class="card-body">
            <div class="table-responsive">
              <table class="table table-striped table-bordered">
                <thead>
                <tr>
                  <th>#</th>
                  <th>Фото</th>
                  <th>Название</th>
                  <th>Артикул</th>
                  <th>Логистика</th>
                  <th>Логистика (средний на единицу товара)</th>
                  <th>Прочие удержания</th>
                  <th>Хранение</th>
                  <th>Доплаты</th>
                  <th>Штрафы</th>
                  <th>Возмещение издержек по перевозке</th>
                  <th>Выручка</th>
                  <th>Налоги %</th>
                  <th>Продажи</th>
                  <th>Возвраты</th>
                  <th>Себестоимость</th>
                  <th>~ Чистая прибыль</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item, index) in statData" :key="item.nm_id">
                  <td class="text-center">{{ formatNumber((currentPage - 1) * pageSize + index + 1) }}</td>
                  <td class="text-center">
                    <img
                        v-if="item.photo && item.nm_id > 0"
                        :src="item.photo"
                        alt="Фото"
                        style="width: 150px; height: 100px;"
                        @error="e => e.target.style.display = 'none'"
                    >
                    <span v-else>Нет фото</span>
                  </td>
                  <td>{{ item.name || 'Нет названия' }}</td>
                  <td>{{ item.nm_id }}</td>
                  <td>{{ formatCurrency(item.delivery_rub) }}</td>
                  <td>{{ formatCurrency(item.delivery_per_unit) }}</td>
                  <td>{{ formatCurrency(item.deduction) }}</td>
                  <td>{{ formatCurrency(item.storage_fee) }}</td>
                  <td>{{ formatCurrency(item.additional_payment) }}</td>
                  <td>{{ formatCurrency(item.penalty) }}</td>
                  <td>{{ formatCurrency(item.rebill_logistic_cost) }}</td>
                  <td>{{ formatCurrency(item.ppvz_for_pay) }}</td>
                  <td>{{ taxes }} %</td>
                  <td>{{ formatNumber(item.sales) }} шт. ({{ formatNumber(item.count_sales) }})</td>
                  <td>{{ formatNumber(item.returns) }} шт. ({{ formatNumber(item.count_refund) }})</td>
                  <td>{{ formatCurrency(0) }}</td>
                  <td>{{ formatCurrency(item.net_profit) }}</td>
                </tr>
                </tbody>
              </table>
            </div>

            <!-- Пагинация -->
            <nav v-if="totalPages > 1" class="mt-3">
              <ul class="pagination justify-content-center">
                <li class="page-item" :class="{ disabled: currentPage === 1 }">
                  <a class="page-link" href="#" @click.prevent="goToPage(currentPage - 1)">
                    &laquo;
                  </a>
                </li>

                <li v-for="page in pagesArray" :key="page"
                    class="page-item" :class="{ active: page === currentPage }">
                  <a class="page-link" href="#" @click.prevent="goToPage(page)">
                    {{ page }}
                  </a>
                </li>

                <li class="page-item" :class="{ disabled: currentPage === totalPages }">
                  <a class="page-link" href="#" @click.prevent="goToPage(currentPage + 1)">
                    &raquo;
                  </a>
                </li>
              </ul>
            </nav>

            <!-- Итоги -->
            <div class="alert alert-info mt-3">
              <h5>Итого за период:</h5>
              <div class="row">
                <div class="col-md-4">
                  <p><strong>Общая выручка:</strong> {{ formatCurrency(summary.total_ppvz_for_pay) }}</p>
                  <p><strong>Общая логистика:</strong> {{ formatCurrency(summary.total_delivery_rub) }}</p>
                </div>
                <div class="col-md-4">
                  <p><strong>Продажи:</strong> {{ formatNumber(summary.total_count_sales) }} заказов,
                    {{ formatNumber(summary.total_quantity) }} шт.</p>
                  <p><strong>Возвраты:</strong> {{ formatNumber(summary.total_count_refund) }} заказов,
                    {{ formatNumber(summary.total_return_amount) }} шт.</p>
                </div>
                <div class="col-md-4">
                  <p><strong>Уникальных товаров:</strong> {{ formatNumber(summary.unique_products) }}</p>
                  <p><strong>~ Чистая прибыль:</strong> {{ formatCurrency(summary.total_net_profit) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Сообщение, если нет данных -->
        <div v-if="!loading && statData.length === 0 && !error" class="alert alert-warning">
          Нет данных за выбранный период
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.table-responsive {
  overflow-x: auto;
}

.table img {
  max-width: 150px;
  max-height: 100px;
  object-fit: contain;
}

.pagination {
  margin-bottom: 0;
}

.page-item.active .page-link {
  background-color: #28a745;
  border-color: #28a745;
}

.page-link {
  color: #28a745;
}

.page-link:hover {
  color: #1e7e34;
}
</style>