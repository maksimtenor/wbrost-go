<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import Navbar from "../../components/layout/Navbar.vue"
import Sidebar from "../../components/layout/Sidebar.vue"

const router = useRouter()

// Данные
const articles = ref([])
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
const fetchArticles = async () => {
  if (!localStorage.getItem('token')) {
    router.push('/login')
    return
  }

  try {
    loading.value = true
    error.value = ''

    const token = localStorage.getItem('token')
    const response = await axios.get('/api/articles', {
      headers: {
        'Authorization': `Bearer ${token}`
      },
      params: {
        page: currentPage.value,
        pageSize: pageSize.value,
        search: searchQuery.value || undefined
      }
    })

    articles.value = response.data.data || []

    // Пагинация
    if (response.data.pagination) {
      currentPage.value = response.data.pagination.current_page
      totalItems.value = response.data.pagination.total_items
      totalPages.value = response.data.pagination.total_pages
    }

    if (articles.value.length === 0) {
      error.value = 'Нет данных'
    }

  } catch (err) {
    console.error('Error fetching articles:', err)
    if (err.response?.data?.error) {
      error.value = `Ошибка: ${err.response.data.error}`
    } else {
      error.value = 'Ошибка загрузки данных'
    }
  } finally {
    loading.value = false
  }
}

// Запрос обновления карточек
const requestArticlesUpdate = async () => {
  if (!localStorage.getItem('token')) {
    router.push('/login')
    return
  }

  try {
    loadingRequest.value = true
    const token = localStorage.getItem('token')

    const response = await axios.post('/api/articles/request',
        {},
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
    )

    if (response.data.success) {
      alert('✅ Запрос на обновление карточек поставлен в очередь')
    }
  } catch (err) {
    console.error('Error requesting articles update:', err)
    if (err.response?.data?.error) {
      alert(`❌ Ошибка: ${err.response.data.error}`)
    } else {
      alert('❌ Ошибка отправки запроса')
    }
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

    await axios.post('/api/articles/cost-price',
        {
          articule: articule,
          cost_price: costPrice
        },
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
    )

    // Обновляем данные на странице
    const article = articles.value.find(a => a.articule === articule)
    if (article) {
      article.cost_price = costPrice
    }

  } catch (err) {
    console.error('Error updating cost price:', err)
    if (err.response?.data?.error) {
      alert(`❌ Ошибка: ${err.response.data.error}`)
    } else {
      alert('❌ Ошибка обновления себестоимости')
    }
  }
}

// Обработчики для редактирования себестоимости
const startEditCostPrice = (articule) => {
  const article = articles.value.find(a => a.articule === articule)
  if (article) {
    article.editingCostPrice = true
    article.tempCostPrice = article.cost_price || ''
  }
}

const saveCostPrice = (articule) => {
  const article = articles.value.find(a => a.articule === articule)
  if (article && article.tempCostPrice !== undefined) {
    updateCostPrice(articule, article.tempCostPrice)
    article.editingCostPrice = false
    delete article.tempCostPrice
  }
}

const cancelEditCostPrice = (articule) => {
  const article = articles.value.find(a => a.articule === articule)
  if (article) {
    article.editingCostPrice = false
    delete article.tempCostPrice
  }
}

// Навигация по страницам
const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchArticles()
}

// Поиск
const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

// Форматирование даты
const formatDate = (dateString) => {
  if (!dateString) return 'Нет данных'
  const date = new Date(dateString)
  return date.toLocaleDateString('ru-RU')
}

// Форматирование цены
const formatPrice = (price) => {
  if (!price || price === 'Не указана' || price === '0') return 'Не указана'
  return `${price} ₽`
}

onMounted(() => {
  fetchArticles()
})

// Наблюдаем за изменениями pageSize
watch(pageSize, () => {
  currentPage.value = 1
  fetchArticles()
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
            <h1 class="m-0">Товары</h1>
          </div>
          <div class="col-sm-6 text-right">
            <button
                class="btn btn-success"
                @click="requestArticlesUpdate"
                :disabled="loadingRequest"
            >
              <span v-if="loadingRequest" class="spinner-border spinner-border-sm"></span>
              {{ loadingRequest ? 'Загрузка...' : 'Обновить карточки' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="container-fluid">
        <!-- Панель поиска и фильтров -->
        <div class="card">
          <div class="card-body">
            <div class="row">
              <div class="col-md-6">
                <div class="input-group">
                  <input
                      type="text"
                      v-model="searchQuery"
                      class="form-control"
                      placeholder="Поиск по названию или артикулу"
                      @keyup.enter="handleSearch"
                  >
                  <div class="input-group-append">
                    <button class="btn btn-primary" @click="handleSearch">
                      <i class="fas fa-search"></i> Найти
                    </button>
                  </div>
                </div>
              </div>
              <div class="col-md-3">
                <select v-model="pageSize" class="form-control">
                  <option value="10">10 на странице</option>
                  <option value="20">20 на странице</option>
                  <option value="50">50 на странице</option>
                  <option value="100">100 на странице</option>
                </select>
              </div>
            </div>
          </div>
        </div>

        <!-- Сообщение об ошибке -->
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <!-- Таблица карточек -->
        <div class="card" v-if="!loading && articles.length > 0">
          <div class="card-body">
            <div class="table-responsive">
              <table class="table table-striped table-bordered">
                <thead>
                <tr>
                  <th>#</th>
                  <th>Фото</th>
                  <th>Артикул</th>
                  <th>Название</th>
                  <th>Создан</th>
                  <th>Обновлен</th>
                  <th>Себестоимость</th>
                  <th>Действие</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(article, index) in articles" :key="article.articule">
                  <td class="text-center">{{ (currentPage - 1) * pageSize + index + 1 }}</td>
                  <td class="text-center">
                    <img
                        v-if="article.photo"
                        :src="article.photo"
                        alt="Фото"
                        style="width: 80px; height: 80px; object-fit: contain;"
                        @error="e => e.target.style.display = 'none'"
                    >
                    <span v-else class="text-muted">Нет фото</span>
                  </td>
                  <td>{{ article.articule }}</td>
                  <td>{{ article.name || 'Нет названия' }}</td>
                  <td>{{ formatDate(article.created) }}</td>
                  <td>{{ formatDate(article.updated) }}</td>
                  <td>
                    <div v-if="article.editingCostPrice">
                      <div class="input-group input-group-sm">
                        <input
                            type="text"
                            v-model="article.tempCostPrice"
                            class="form-control"
                            @keyup.enter="saveCostPrice(article.articule)"
                        >
                        <div class="input-group-append">
                          <button
                              class="btn btn-success btn-sm"
                              @click="saveCostPrice(article.articule)"
                              title="Сохранить"
                          >
                            <i class="fas fa-check"></i>
                          </button>
                          <button
                              class="btn btn-danger btn-sm"
                              @click="cancelEditCostPrice(article.articule)"
                              title="Отмена"
                          >
                            <i class="fas fa-times"></i>
                          </button>
                        </div>
                      </div>
                    </div>
                    <div v-else>
                      {{ formatPrice(article.cost_price) }}
                      <button
                          class="btn btn-link btn-sm"
                          @click="startEditCostPrice(article.articule)"
                          title="Редактировать"
                      >
                        <i class="fas fa-pencil-alt"></i>
                      </button>
                    </div>
                  </td>
                  <td class="text-center">
                    <button
                        class="btn btn-warning btn-sm"
                        @click="startEditCostPrice(article.articule)"
                        title="Редактировать себестоимость"
                    >
                      <i class="fas fa-edit"></i>
                    </button>
                  </td>
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

              <div class="text-center mt-2 text-muted">
                Показано {{ articles.length }} из {{ totalItems }} товаров
              </div>
            </nav>
          </div>
        </div>

        <!-- Сообщение, если нет данных -->
        <div v-if="!loading && articles.length === 0 && !error" class="alert alert-warning">
          Нет данных. Нажмите "Обновить карточки" для загрузки товаров.
        </div>

        <!-- Индикатор загрузки -->
        <div v-if="loading" class="text-center">
          <div class="spinner-border text-primary" role="status">
            <span class="sr-only">Загрузка...</span>
          </div>
          <p class="mt-2">Загрузка данных...</p>
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
  max-width: 80px;
  max-height: 80px;
  object-fit: contain;
}

.pagination {
  margin-bottom: 0;
}

.page-item.active .page-link {
  background-color: #007bff;
  border-color: #007bff;
}

.page-link {
  color: #007bff;
}

.page-link:hover {
  color: #0056b3;
}

.input-group-append .btn {
  margin-left: 2px;
}
</style>