<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from "../../components/layout/Navbar.vue"
import Sidebar from "../../components/layout/Sidebar.vue"

const router = useRouter()

// Данные
const products = ref([])
const loading = ref(false)
const loadingRequest = ref(false)
const error = ref('')
const searchQuery = ref('')

// Пагинация
const currentPage = ref(1)
const pageSize = ref(10)
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

// Загрузка карточек товаров
const fetchProducts = async () => {
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

    const response = await fetch('/api/articles', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      params: {
        page: currentPage.value,
        pageSize: pageSize.value
      }
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }

    const data = await response.json()
    products.value = data.data || []

    // Пагинация (замените на реальные данные с сервера)
    if (data.pagination) {
      currentPage.value = data.pagination.current_page
      totalItems.value = data.pagination.total_items
      totalPages.value = data.pagination.total_pages
    } else {
      totalItems.value = products.value.length
      totalPages.value = Math.ceil(products.value.length / pageSize.value)
    }

    if (products.value.length === 0) {
      error.value = 'Нет данных о товарах'
    }

  } catch (err) {
    console.error('Error fetching products:', err)
    error.value = 'Ошибка загрузки данных'
  } finally {
    loading.value = false
  }
}

// Запрос обновления карточек
const requestProductsUpdate = async () => {
  if (!localStorage.getItem('token')) {
    router.push('/login')
    return
  }

  try {
    loadingRequest.value = true
    const token = localStorage.getItem('token')

    const response = await fetch('http://localhost:8080/api/articles/request', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({})
    })

    const data = await response.json()

    if (data.success) {
      // Показать уведомление об успехе
      error.value = ''
      showMessage('✅ Запрос на обновление карточек поставлен в очередь', 'success')
      // Обновляем список
      setTimeout(() => {
        fetchProducts()
      }, 1000)
    }
  } catch (err) {
    console.error('Error requesting products update:', err)
    showMessage('❌ Ошибка отправки запроса', 'error')
  } finally {
    loadingRequest.value = false
  }
}

// Обновление себестоимости
const updateCostPrice = async (articule, costPrice) => {
  if (!localStorage.getItem('token')) {
    router.push('/login')
    return
  }

  try {
    const token = localStorage.getItem('token')

    const response = await fetch('http://localhost:8080/api/articles/cost-price', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        articule: articule,
        cost_price: costPrice
      })
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }

    // Обновляем данные на странице
    const product = products.value.find(p => p.articule === articule)
    if (product) {
      product.cost_price = costPrice
    }

    showMessage('✅ Себестоимость обновлена', 'success')
  } catch (err) {
    console.error('Error updating cost price:', err)
    showMessage('❌ Ошибка обновления себестоимости', 'error')
  }
}

// Вспомогательная функция для сообщений
const showMessage = (text, type) => {
  error.value = text
  setTimeout(() => {
    error.value = ''
  }, 3000)
}

// Обработчики для редактирования себестоимости
const startEditCostPrice = (articule) => {
  const product = products.value.find(p => p.articule === articule)
  if (product) {
    product.editingCostPrice = true
    product.tempCostPrice = product.cost_price || ''
  }
}

const saveCostPrice = (articule) => {
  const product = products.value.find(p => p.articule === articule)
  if (product && product.tempCostPrice !== undefined) {
    updateCostPrice(articule, product.tempCostPrice)
    product.editingCostPrice = false
    delete product.tempCostPrice
  }
}

const cancelEditCostPrice = (articule) => {
  const product = products.value.find(p => p.articule === articule)
  if (product) {
    product.editingCostPrice = false
    delete product.tempCostPrice
  }
}

// Навигация по страницам
const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchProducts()
}

// Фильтрация и пагинация
const filteredProducts = computed(() => {
  let filtered = [...products.value]

  // Поиск
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(product =>
        (product.name && product.name.toLowerCase().includes(query)) ||
        (product.articule && product.articule.toString().toLowerCase().includes(query))
    )
  }

  // Пагинация
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filtered.slice(start, end)
})

const filterProducts = () => {
  currentPage.value = 1
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

// Форматирование цены
const formatPrice = (price) => {
  if (!price || price === 'Не указана' || price === '0') return 'Не указана'
  return `${parseFloat(price).toLocaleString('ru-RU')} ₽`
}

onMounted(() => {
  fetchProducts()
})

// Наблюдаем за изменениями pageSize
watch(pageSize, () => {
  currentPage.value = 1
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
                <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path>
                <line x1="7" y1="7" x2="7.01" y2="7"></line>
              </svg>
              Карточки товаров
            </h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="products-container">
        <!-- Сообщение об успехе/ошибке -->
        <transition name="fade">
          <div v-if="error" :class="['message-box', error.includes('✅') ? 'message-success' : 'message-error']">
            <svg v-if="error.includes('✅')" class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
            <svg v-else class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{{ error }}</span>
          </div>
        </transition>

        <!-- Панель управления -->
        <div class="control-panel">
          <div class="search-section">
            <div class="search-box">
              <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"></circle>
                <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
              </svg>
              <input
                  type="text"
                  v-model="searchQuery"
                  class="search-input"
                  placeholder="Поиск по названию или артикулу"
                  @input="filterProducts"
              >
            </div>

            <button
                class="btn-update"
                @click="requestProductsUpdate"
                :disabled="loadingRequest"
            >
              <svg
                  :class="['btn-icon', { 'spinning': loadingRequest }]"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
              >
                <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"></path>
              </svg>
              {{ loadingRequest ? 'Загрузка...' : 'Обновить карточки' }}
            </button>
          </div>

          <div class="pagination-controls">
            <div class="page-size-selector">
              <label class="page-size-label">
                <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
                Показать:
              </label>
              <select v-model="pageSize" class="page-size-select">
                <option value="10">10</option>
                <option value="20">20</option>
                <option value="50">50</option>
                <option value="100">100</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Индикатор загрузки -->
        <div v-if="loading" class="loading-container">
          <div class="loading-spinner"></div>
          <p class="loading-text">Загрузка товаров...</p>
        </div>

        <!-- Таблица товаров -->
        <div v-else-if="filteredProducts.length > 0" class="products-table-container">
          <table class="products-table">
            <thead>
            <tr>
              <th class="table-photo">Фото</th>
              <th class="table-article">Артикул</th>
              <th class="table-name">Название</th>
              <th class="table-created">Создан</th>
              <th class="table-updated">Обновлен</th>
              <th class="table-cost">Себестоимость</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="product in filteredProducts" :key="product.articule" class="product-row">
              <td class="table-photo">
                <div class="product-image-container">
                  <img
                      v-if="product.photo"
                      :src="product.photo"
                      :alt="product.name"
                      class="product-image"
                      @error="e => e.target.style.display = 'none'"
                  >
                  <div v-else class="product-image-placeholder">
                    <svg class="placeholder-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                      <circle cx="8.5" cy="8.5" r="1.5"></circle>
                      <polyline points="21 15 16 10 5 21"></polyline>
                    </svg>
                  </div>
                </div>
              </td>
              <td class="table-article">
                <span class="article-value">{{ product.articule }}</span>
              </td>
              <td class="table-name">
                <div class="product-name-container">
                  <h4 class="product-name">{{ product.name || 'Без названия' }}</h4>
                </div>
              </td>
              <td class="table-created">
                <div class="date-container">
                  <svg class="date-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <polyline points="12 6 12 12 16 14"></polyline>
                  </svg>
                  <span class="date-value">{{ formatDate(product.created) }}</span>
                </div>
              </td>
              <td class="table-updated">
                <div class="date-container">
                  <svg class="date-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <path d="M12 6v6l4 2"></path>
                  </svg>
                  <span class="date-value">{{ formatDate(product.updated) }}</span>
                </div>
              </td>
              <td class="table-cost">
                <div v-if="product.editingCostPrice" class="cost-edit-container">
                  <div class="cost-edit-wrapper">
                    <input
                        type="text"
                        v-model="product.tempCostPrice"
                        class="cost-input"
                        @keyup.enter="saveCostPrice(product.articule)"
                        placeholder="Введите цену"
                        @keyup.esc="cancelEditCostPrice(product.articule)"
                    >
                    <div class="cost-edit-buttons">
                      <button
                          class="cost-btn cost-btn-save"
                          @click="saveCostPrice(product.articule)"
                          title="Сохранить"
                      >
                        <svg class="cost-btn-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                          <polyline points="17 21 17 13 7 13 7 21"></polyline>
                          <polyline points="7 3 7 8 15 8"></polyline>
                        </svg>
                      </button>
                      <button
                          class="cost-btn cost-btn-cancel"
                          @click="cancelEditCostPrice(product.articule)"
                          title="Отмена"
                      >
                        <svg class="cost-btn-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <line x1="18" y1="6" x2="6" y2="18"></line>
                          <line x1="6" y1="6" x2="18" y2="18"></line>
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
                <div v-else class="cost-display">
                  <span class="cost-value">{{ formatPrice(product.cost_price) }}</span>
                  <button
                      class="cost-edit-btn"
                      @click="startEditCostPrice(product.articule)"
                      title="Редактировать"
                  >
                    <svg class="cost-edit-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
            </tbody>
          </table>
        </div>

        <!-- Сообщение, если нет данных -->
        <div v-else-if="!loading && products.length === 0" class="empty-state">
          <svg class="empty-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
          <h3 class="empty-title">Нет данных о товарах</h3>
          <p class="empty-text">Нажмите "Обновить карточки" для загрузки товаров</p>
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
            Показано {{ filteredProducts.length }} из {{ totalItems }} товаров
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

.products-container {
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

/* Панель управления */
.control-panel {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 24px;
  margin-bottom: 24px;
}

.search-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 300px;
  position: relative;
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  width: 20px;
  height: 20px;
  color: #6b7280;
}

.search-input {
  width: 100%;
  padding: 12px 16px 12px 48px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 15px;
  transition: all 0.2s ease;
  background: #fafafa;
}

.search-input:hover {
  border-color: #cbd5e0;
  background: white;
}

.search-input:focus {
  outline: none;
  border-color: #4f46e5;
  background: white;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.btn-update {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
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

.btn-update:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
  background: linear-gradient(135deg, #0da271 0%, #10b981 100%);
}

.btn-update:active:not(:disabled) {
  transform: translateY(0);
}

.btn-update:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-icon {
  width: 18px;
  height: 18px;
}

.pagination-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid #e5e7eb;
  padding-top: 20px;
}

.page-size-selector {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-size-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.label-icon {
  width: 16px;
  height: 16px;
}

.page-size-select {
  padding: 8px 12px;
  border: 2px solid #e5e7eb;
  border-radius: 6px;
  background: white;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s ease;
}

.page-size-select:hover {
  border-color: #cbd5e0;
}

.page-size-select:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
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

/* Таблица товаров */
.products-table-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin-bottom: 24px;
}

.products-table {
  width: 100%;
  border-collapse: collapse;
}

.products-table thead {
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
  border-bottom: 2px solid #e5e7eb;
}

.products-table th {
  padding: 16px 20px;
  text-align: left;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.products-table tbody tr {
  border-bottom: 1px solid #e5e7eb;
  transition: background-color 0.2s ease;
}

.products-table tbody tr:hover {
  background-color: #f9fafb;
}

.products-table tbody tr:last-child {
  border-bottom: none;
}

.products-table td {
  padding: 16px 20px;
  vertical-align: middle;
}

/* Колонки таблицы */
.table-photo {
  width: 80px;
}

.table-article {
  width: 120px;
  min-width: 120px;
}

.table-name {
  min-width: 300px;
}

.table-created,
.table-updated {
  width: 140px;
  min-width: 140px;
}

.table-cost {
  width: 200px;
  min-width: 200px;
}

/* Стили для содержимого таблицы */
.product-image-container {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 4px;
}

.product-image-placeholder {
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

.article-value {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  background: #f3f4f6;
  padding: 4px 10px;
  border-radius: 6px;
  display: inline-block;
}

.product-name-container {
  max-width: 400px;
}

.product-name {
  font-size: 15px;
  font-weight: 500;
  color: #1a202c;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.date-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.date-icon {
  width: 16px;
  height: 16px;
  color: #6b7280;
  flex-shrink: 0;
}

.date-value {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

/* Себестоимость */
.cost-display {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cost-value {
  font-size: 15px;
  color: #10b981;
  font-weight: 600;
  padding: 8px 12px;
  background: #d1fae5;
  border-radius: 8px;
  min-width: 100px;
  text-align: center;
  display: inline-block;
}

.cost-edit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.cost-edit-btn:hover {
  border-color: #4f46e5;
  background: #f5f3ff;
}

.cost-edit-icon {
  width: 16px;
  height: 16px;
  color: #6b7280;
}

/* Редактирование себестоимости */
.cost-edit-container {
  width: 100%;
}

.cost-edit-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f8fafc;
  padding: 4px;
  border-radius: 8px;
  border: 2px solid #e5e7eb;
}

.cost-input {
  flex: 1;
  min-width: 100px;
  padding: 8px 12px;
  border: none;
  background: transparent;
  font-size: 15px;
  font-weight: 500;
  color: #374151;
  outline: none;
}

.cost-input::placeholder {
  color: #9ca3af;
}

.cost-edit-buttons {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.cost-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.cost-btn-save {
  background: #10b981;
}

.cost-btn-save:hover {
  background: #0da271;
  transform: translateY(-2px);
}

.cost-btn-cancel {
  background: #ef4444;
}

.cost-btn-cancel:hover {
  background: #dc2626;
  transform: translateY(-2px);
}

.cost-btn-icon {
  width: 16px;
  height: 16px;
  color: white;
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

/* Пагинация */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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

.spinning {
  animation: spin 1s linear infinite;
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
@media (max-width: 1024px) {
  .products-table-container {
    overflow-x: auto;
  }

  .products-table {
    min-width: 1000px;
  }

  .search-section {
    flex-direction: column;
    align-items: stretch;
  }

  .search-box {
    min-width: 100%;
  }

  .btn-update {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 768px) {
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