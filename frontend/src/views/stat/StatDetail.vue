<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import apiClient from '@/api/client'
import axios from 'axios'
import Navbar from "../../components/layout/Navbar.vue"
import Sidebar from "../../components/layout/Sidebar.vue"
import BaseLayout from "../../components/layout/BaseLayout.vue";

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

// Загрузка статистики
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

    const response = await apiClient.get('/stat/details', {
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
const DEFAULT_CURRENCY_TYPE = import.meta.env.VITE_DEFAULT_CURRENCY_TYPE;
const formatCurrency = (value) => {
  if (value === null || value === undefined || value === '') return DEFAULT_CURRENCY_TYPE+' 0.00'
  const num = typeof value === 'string' ? parseFloat(value) : value
  return DEFAULT_CURRENCY_TYPE+` ${num.toFixed(2)}`
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
  <BaseLayout>
    <template #title-icon>
      <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
        <circle cx="8.5" cy="7" r="4"></circle>
        <path d="M20 8v6"></path>
        <path d="M23 11h-6"></path>
      </svg>
    </template>
    <template #title>Детальная статистика</template>
    <template #header-right>
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
    </template>

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
  </BaseLayout>
</template>

<style scoped>
@import '@/assets/css/views/stat/statdetail.css';
</style>