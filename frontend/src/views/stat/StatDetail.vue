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
            <h1 class="page-title">
              <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                <circle cx="8.5" cy="7" r="4"></circle>
                <path d="M20 8v6"></path>
                <path d="M23 11h-6"></path>
              </svg>
              Детальная статистика
            </h1>
          </div>
          <div class="col-sm-6 text-right">
            <div v-if="statData.length > 0" class="total-items">
              <svg class="total-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="8" y1="6" x2="21" y2="6"></line>
                <line x1="8" y1="12" x2="21" y2="12"></line>
                <line x1="8" y1="18" x2="21" y2="18"></line>
                <line x1="3" y1="6" x2="3.01" y2="6"></line>
                <line x1="3" y1="12" x2="3.01" y2="12"></line>
                <line x1="3" y1="18" x2="3.01" y2="18"></line>
              </svg>
              <span>Всего позиций: {{ formatNumber(totalItems) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="statistics-container">
        <!-- Сообщение об ошибке -->
        <transition name="fade">
          <div v-if="error" class="message-box message-error">
            <svg class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{{ error }}</span>
          </div>
        </transition>

        <!-- Панель фильтров -->
        <div class="filter-panel">
          <form @submit="handleSubmit" class="filter-form">
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
                    v-model="dateFrom"
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
                    v-model="dateTo"
                    class="filter-input"
                    id="dateTo"
                    required
                >
              </div>

              <div class="filter-group">
                <label class="filter-label" for="pageSize">
                  <svg class="filter-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="8" y1="6" x2="21" y2="6"></line>
                    <line x1="8" y1="12" x2="21" y2="12"></line>
                    <line x1="8" y1="18" x2="21" y2="18"></line>
                    <line x1="3" y1="6" x2="3.01" y2="6"></line>
                    <line x1="3" y1="12" x2="3.01" y2="12"></line>
                    <line x1="3" y1="18" x2="3.01" y2="18"></line>
                  </svg>
                  Показать
                </label>
                <select v-model="pageSize" class="filter-select" id="pageSize">
                  <option value="10">10 строк</option>
                  <option value="20">20 строк</option>
                  <option value="50">50 строк</option>
                  <option value="100">100 строк</option>
                </select>
              </div>

              <div class="filter-group align-bottom">
                <button type="submit" class="filter-btn" :disabled="loading">
                  <svg v-if="!loading" class="filter-btn-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12a9 9 0 0 1-9 9m9-9a9 9 0 0 0-9-9m9 9H3m9 9a9 9 0 0 1-9-9m9 9c1.66 0 3-4 3-9s-1.34-9-3-9m0 18c-1.66 0-3-4-3-9s1.34-9 3-9"></path>
                  </svg>
                  <span class="spinner" v-if="loading"></span>
                  {{ loading ? 'Загрузка...' : 'Показать' }}
                </button>
              </div>
            </div>
          </form>
        </div>

        <!-- Индикатор загрузки -->
        <div v-if="loading" class="loading-container">
          <div class="loading-spinner"></div>
          <p class="loading-text">Загрузка статистики...</p>
        </div>

        <!-- Таблица статистики -->
        <div v-else-if="statData.length > 0" class="statistics-table-container">
          <div class="table-wrapper">
            <table class="statistics-table">
              <thead>
              <tr>
                <th class="table-photo">Фото</th>
                <th class="table-name">Название</th>
                <th class="table-id">Артикул</th>
                <th class="table-logistic">Логистика</th>
                <th class="table-logistic-unit">Лог. на ед.</th>
                <th class="table-deduction">Удержания</th>
                <th class="table-storage">Хранение</th>
                <th class="table-payment">Доплаты</th>
                <th class="table-penalty">Штрафы</th>
                <th class="table-rebill">Возмещение</th>
                <th class="table-revenue">Выручка</th>
                <th class="table-taxes">Налог</th>
                <th class="table-sales">Продажи</th>
                <th class="table-returns">Возвраты</th>
                <th class="table-cost">Себест.</th>
                <th class="table-profit">Прибыль</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(item, index) in statData" :key="item.nm_id" class="stat-row">
                <td class="table-photo">
                  <div class="stat-image-container">
                    <img
                        v-if="item.photo && item.nm_id > 0"
                        :src="item.photo"
                        :alt="item.name"
                        class="stat-image"
                        @error="e => e.target.style.display = 'none'"
                    >
                    <div v-else class="stat-image-placeholder">
                      <svg class="placeholder-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                        <circle cx="8.5" cy="8.5" r="1.5"></circle>
                        <polyline points="21 15 16 10 5 21"></polyline>
                      </svg>
                    </div>
                  </div>
                </td>
                <td class="table-name">
                  <div class="stat-name-container">
                    <h4 class="stat-name">{{ item.name || 'Без названия' }}</h4>
                  </div>
                </td>
                <td class="table-id">
                  <span class="stat-id">{{ item.nm_id }}</span>
                </td>
                <td class="table-logistic">
                  <div class="stat-amount" :class="{ 'amount-negative': parseFloat(item.delivery_rub) < 0 }">
                    {{ formatCurrency(item.delivery_rub) }}
                  </div>
                </td>
                <td class="table-logistic-unit">
                  <div class="stat-amount">
                    {{ formatCurrency(item.delivery_per_unit) }}
                  </div>
                </td>
                <td class="table-deduction">
                  <div class="stat-amount" :class="{ 'amount-negative': parseFloat(item.deduction) < 0 }">
                    {{ formatCurrency(item.deduction) }}
                  </div>
                </td>
                <td class="table-storage">
                  <div class="stat-amount" :class="{ 'amount-negative': parseFloat(item.storage_fee) < 0 }">
                    {{ formatCurrency(item.storage_fee) }}
                  </div>
                </td>
                <td class="table-payment">
                  <div class="stat-amount" :class="{ 'amount-positive': parseFloat(item.additional_payment) > 0 }">
                    {{ formatCurrency(item.additional_payment) }}
                  </div>
                </td>
                <td class="table-penalty">
                  <div class="stat-amount" :class="{ 'amount-negative': parseFloat(item.penalty) < 0 }">
                    {{ formatCurrency(item.penalty) }}
                  </div>
                </td>
                <td class="table-rebill">
                  <div class="stat-amount" :class="{ 'amount-positive': parseFloat(item.rebill_logistic_cost) > 0 }">
                    {{ formatCurrency(item.rebill_logistic_cost) }}
                  </div>
                </td>
                <td class="table-revenue">
                  <div class="stat-amount amount-positive">
                    {{ formatCurrency(item.ppvz_for_pay) }}
                  </div>
                </td>
                <td class="table-taxes">
                  <div class="stat-taxes">
                    <span class="taxes-value">{{ taxes }}%</span>
                  </div>
                </td>
                <td class="table-sales">
                  <div class="stat-sales">
                    <div class="sales-count">{{ formatNumber(item.sales) }} шт.</div>
                    <div class="sales-orders">{{ formatNumber(item.count_sales) }} зак.</div>
                  </div>
                </td>
                <td class="table-returns">
                  <div class="stat-returns">
                    <div class="returns-count">{{ formatNumber(item.returns) }} шт.</div>
                    <div class="returns-orders">{{ formatNumber(item.count_refund) }} зак.</div>
                  </div>
                </td>
                <td class="table-cost">
                  <div class="stat-amount">
                    {{ formatCurrency(0) }}
                  </div>
                </td>
                <td class="table-profit">
                  <div class="stat-profit" :class="{
                    'profit-positive': parseFloat(item.net_profit) > 0,
                    'profit-negative': parseFloat(item.net_profit) < 0
                  }">
                    {{ formatCurrency(item.net_profit) }}
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
              Показано {{ statData.length }} из {{ formatNumber(totalItems) }} позиций
            </div>
          </div>

          <!-- Итоговая сводка -->
          <div class="summary-panel">
            <h3 class="summary-title">
              <svg class="summary-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="20" x2="12" y2="10"></line>
                <line x1="18" y1="20" x2="18" y2="4"></line>
                <line x1="6" y1="20" x2="6" y2="16"></line>
              </svg>
              Итоги за период
            </h3>

            <div class="summary-grid">
              <div class="summary-card">
                <div class="summary-header">
                  <svg class="summary-card-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="1" x2="12" y2="23"></line>
                    <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
                  </svg>
                  <h4 class="summary-card-title">Финансы</h4>
                </div>
                <div class="summary-content">
                  <div class="summary-item">
                    <span class="summary-label">Общая выручка:</span>
                    <span class="summary-value amount-positive">{{ formatCurrency(summary.total_ppvz_for_pay) }}</span>
                  </div>
                  <div class="summary-item">
                    <span class="summary-label">Общая логистика:</span>
                    <span class="summary-value amount-negative">{{ formatCurrency(summary.total_delivery_rub) }}</span>
                  </div>
                  <div class="summary-item">
                    <span class="summary-label">Прибыль:</span>
                    <span class="summary-value" :class="{
                      'profit-positive': parseFloat(summary.total_net_profit) > 0,
                      'profit-negative': parseFloat(summary.total_net_profit) < 0
                    }">
                      {{ formatCurrency(summary.total_net_profit) }}
                    </span>
                  </div>
                </div>
              </div>

              <div class="summary-card">
                <div class="summary-header">
                  <svg class="summary-card-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                    <polyline points="13 2 13 9 20 9"></polyline>
                  </svg>
                  <h4 class="summary-card-title">Продажи</h4>
                </div>
                <div class="summary-content">
                  <div class="summary-item">
                    <span class="summary-label">Заказов:</span>
                    <span class="summary-value">{{ formatNumber(summary.total_count_sales) }}</span>
                  </div>
                  <div class="summary-item">
                    <span class="summary-label">Товаров:</span>
                    <span class="summary-value">{{ formatNumber(summary.total_quantity) }} шт.</span>
                  </div>
                  <div class="summary-item">
                    <span class="summary-label">Уникальных:</span>
                    <span class="summary-value">{{ formatNumber(summary.unique_products) }}</span>
                  </div>
                </div>
              </div>

              <div class="summary-card">
                <div class="summary-header">
                  <svg class="summary-card-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="16"></line>
                    <line x1="8" y1="12" x2="16" y2="12"></line>
                  </svg>
                  <h4 class="summary-card-title">Возвраты</h4>
                </div>
                <div class="summary-content">
                  <div class="summary-item">
                    <span class="summary-label">Заказов:</span>
                    <span class="summary-value">{{ formatNumber(summary.total_count_refund) }}</span>
                  </div>
                  <div class="summary-item">
                    <span class="summary-label">Товаров:</span>
                    <span class="summary-value">{{ formatNumber(summary.total_return_amount) }} шт.</span>
                  </div>
                  <div class="summary-item">
                    <span class="summary-label">Ставка:</span>
                    <span class="summary-value">
                      {{
                        summary.total_quantity > 0
                            ? ((summary.total_return_amount / summary.total_quantity) * 100).toFixed(1)
                            : '0'
                      }}%
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Сообщение, если нет данных -->
        <div v-else-if="!loading && statData.length === 0 && !error" class="empty-state">
          <svg class="empty-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
          <h3 class="empty-title">Нет данных за выбранный период</h3>
          <p class="empty-text">Измените параметры фильтрации или выберите другой период</p>
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

.statistics-container {
  max-width: 1600px;
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

/* Панель фильтров */
.filter-panel {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 24px;
  margin-bottom: 24px;
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

@media (min-width: 1024px) {
  .filter-grid {
    grid-template-columns: repeat(4, 1fr);
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

.filter-select {
  padding: 12px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 14px;
  color: #374151;
  background: #fafafa;
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-select:hover {
  border-color: #cbd5e0;
  background: white;
}

.filter-select:focus {
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
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.3);
}

.filter-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
  background: linear-gradient(135deg, #0da271 0%, #10b981 100%);
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

/* Таблица */
.statistics-table-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin-bottom: 24px;
}

.table-wrapper {
  overflow-x: auto;
  padding: 20px;
}

.statistics-table {
  width: 100%;
  border-collapse: collapse;
  min-width: 1500px;
}

.statistics-table thead {
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
  border-bottom: 2px solid #e5e7eb;
}

.statistics-table th {
  padding: 16px 12px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.statistics-table tbody tr {
  border-bottom: 1px solid #e5e7eb;
  transition: background-color 0.2s ease;
}

.statistics-table tbody tr:hover {
  background-color: #f9fafb;
}

.statistics-table tbody tr:last-child {
  border-bottom: none;
}

.statistics-table td {
  padding: 12px;
  vertical-align: middle;
  font-size: 13px;
}

/* Колонки таблицы */
.table-photo {
  width: 80px;
}

.table-name {
  width: 200px;
  min-width: 200px;
}

.table-id {
  width: 100px;
  min-width: 100px;
}

.table-logistic,
.table-logistic-unit,
.table-deduction,
.table-storage,
.table-payment,
.table-penalty,
.table-rebill,
.table-revenue,
.table-cost,
.table-profit {
  width: 120px;
  min-width: 120px;
}

.table-taxes {
  width: 80px;
  min-width: 80px;
}

.table-sales,
.table-returns {
  width: 140px;
  min-width: 140px;
}

/* Стили для содержимого таблицы */
.stat-image-container {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 4px;
}

.stat-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-icon {
  width: 24px;
  height: 24px;
  color: #9ca3af;
}

.stat-name-container {
  max-width: 200px;
}

.stat-name {
  font-size: 14px;
  font-weight: 500;
  color: #1a202c;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.stat-id {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 6px;
  display: inline-block;
}

/* Стили для сумм */
.stat-amount {
  font-size: 13px;
  font-weight: 600;
  padding: 6px 10px;
  border-radius: 6px;
  text-align: right;
  background: #f8fafc;
}

.amount-positive {
  color: #10b981;
  background: #d1fae5;
}

.amount-negative {
  color: #ef4444;
  background: #fee2e2;
}

.stat-taxes {
  text-align: center;
}

.taxes-value {
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  padding: 6px 10px;
  background: #f3f4f6;
  border-radius: 6px;
  display: inline-block;
}

.stat-sales,
.stat-returns {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.sales-count,
.returns-count {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
}

.sales-orders,
.returns-orders {
  font-size: 12px;
  color: #6b7280;
}

.stat-profit {
  font-size: 14px;
  font-weight: 700;
  padding: 8px 12px;
  border-radius: 8px;
  text-align: right;
}

.profit-positive {
  color: #10b981;
  background: #d1fae5;
  border: 2px solid #a7f3d0;
}

.profit-negative {
  color: #ef4444;
  background: #fee2e2;
  border: 2px solid #fecaca;
}

/* Итоговая панель */
.summary-panel {
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
  border-top: 1px solid #e5e7eb;
  padding: 24px;
}

.summary-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 20px 0;
}

.summary-icon {
  width: 24px;
  height: 24px;
  color: #4f46e5;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.summary-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.summary-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid #4f46e5;
}

.summary-card-icon {
  width: 20px;
  height: 20px;
  color: #4f46e5;
}

.summary-card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
}

.summary-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.summary-value {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

/* Пагинация */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 20px;
  border-top: 1px solid #e5e7eb;
}

.pagination-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
}

.pagination-btn:hover:not(:disabled) {
  border-color: #4f46e5;
  background: #f5f3ff;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-icon {
  width: 20px;
  height: 20px;
  color: #374151;
}

.pagination-pages {
  display: flex;
  gap: 8px;
}

.page-btn {
  min-width: 40px;
  height: 40px;
  padding: 0 8px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  background: white;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s ease;
}

.page-btn:hover:not(.active) {
  border-color: #4f46e5;
  background: #f5f3ff;
}

.page-btn.active {
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  color: white;
  border-color: #4f46e5;
  box-shadow: 0 2px 8px rgba(79, 70, 229, 0.3);
}

.pagination-info {
  margin-left: 16px;
  font-size: 14px;
  color: #6b7280;
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

  .summary-grid {
    grid-template-columns: 1fr;
  }

  .pagination {
    flex-wrap: wrap;
    gap: 12px;
  }

  .pagination-info {
    width: 100%;
    text-align: center;
    margin-left: 0;
    margin-top: 8px;
  }
}
</style>