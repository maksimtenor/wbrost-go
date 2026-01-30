<script setup>
import {ref, onMounted, computed, watch} from 'vue'
import {useRouter} from 'vue-router'
import apiClient from '@/api/client'
import BaseLayout from "@/components/layout/BaseLayout.vue";

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

    const response = await apiClient.get('/articles', {
      params: {
        page: currentPage.value,
        pageSize: pageSize.value
      }
    })

    const data = response.data
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
    const response = await apiClient.post('/articles/request', {
      body: {}
    });

    const data = response.data

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
    const response = await apiClient.post('/articles/cost-price', {
      articule: articule,
      cost_price: costPrice
    })


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
  <BaseLayout>
    <template #title-icon>
      <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
           stroke-width="2">
        <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path>
        <line x1="7" y1="7" x2="7.01" y2="7"></line>
      </svg>
    </template>
    <template #title>Карточки товаров</template>

    <div class="products-container">
      <!-- Сообщение об успехе/ошибке -->
      <transition name="fade">
        <div v-if="error" :class="['message-box', error.includes('✅') ? 'message-success' : 'message-error']">
          <svg v-if="error.includes('✅')" class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
               fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
          </svg>
          <svg v-else class="message-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2">
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
            <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                 stroke="currentColor" stroke-width="2">
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
              <svg class="label-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                   stroke="currentColor" stroke-width="2">
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
                  <svg class="placeholder-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                       stroke="currentColor" stroke-width="2">
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
                <svg class="date-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                     stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
                <span class="date-value">{{ formatDate(product.created) }}</span>
              </div>
            </td>
            <td class="table-updated">
              <div class="date-container">
                <svg class="date-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                     stroke="currentColor" stroke-width="2">
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
                      <svg class="cost-btn-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                           stroke="currentColor" stroke-width="2">
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
                      <svg class="cost-btn-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                           stroke="currentColor" stroke-width="2">
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
                  <svg class="cost-edit-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
                       stroke="currentColor" stroke-width="2">
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
        <svg class="empty-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
             stroke-width="2">
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
          <svg class="pagination-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2">
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
          <svg class="pagination-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>

        <div class="pagination-info">
          Показано {{ filteredProducts.length }} из {{ totalItems }} товаров
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<style scoped>
@import '@/assets/css/views/profile/carts.css';
</style>