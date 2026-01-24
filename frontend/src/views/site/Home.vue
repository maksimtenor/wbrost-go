<template>
  <div class="wrapper hold-transition sidebar-mini">
    <!-- Navbar -->
    <Navbar />
    <!-- Main Sidebar Container -->
    <Sidebar />

    <div class="content-wrapper">
      <!-- Content Header -->
      <div class="content-header">
        <div class="container-fluid">
          <div class="row mb-2">
            <div class="col-sm-6">
              <h1 class="m-0">{{ pageTitle }}</h1>
            </div>
            <div class="col-sm-6 text-right">
              <small v-if="!isGuest" class="text-muted">
                Обновлено: {{ lastUpdated }}
              </small>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="content">
        <div class="container-fluid">
          <!-- Для авторизованных -->
          <div v-if="!isGuest">
            <!-- Индикатор загрузки -->
            <div v-if="loading" class="text-center py-5">
              <div class="spinner-border text-success" role="status">
                <span class="sr-only">Загрузка...</span>
              </div>
              <p class="mt-2">Загрузка данных...</p>
            </div>

            <!-- Сообщение об ошибке -->
            <div v-if="error && !loading" class="alert alert-warning alert-dismissible">
              {{ error }}
              <button type="button" class="close" @click="error = ''">
                <span>&times;</span>
              </button>
            </div>

            <div v-if="!loading && !error">
              <!-- Статистика -->
              <div class="row stats-row">
                <!-- Кол-во продаж -->
                <div class="col-lg-3 col-md-6 col-sm-12 mb-4">
                  <div class="stat-card sales-card">
                    <div class="stat-icon">
                      <i class="zmdi zmdi-shopping-cart"></i>
                    </div>
                    <div class="stat-content">
                      <div class="stat-title">Кол-во продаж</div>
                      <div class="stat-value">{{ stats[0].value }}</div>
                      <div class="stat-period">{{ stats[0].period }}</div>
                    </div>
                    <div class="stat-trend">
                      <i class="zmdi zmdi-trending-up"></i>
                    </div>
                  </div>
                </div>

                <!-- Сумма к перечислению -->
                <div class="col-lg-3 col-md-6 col-sm-12 mb-4">
                  <div class="stat-card revenue-card">
                    <div class="stat-icon">
                      <i class="zmdi zmdi-money"></i>
                    </div>
                    <div class="stat-content">
                      <div class="stat-title">Сумма к перечислению</div>
                      <div class="stat-value">{{ stats[1].value }}</div>
                      <div class="stat-period">{{ stats[1].period }}</div>
                    </div>
                    <div class="stat-trend">
                      <i class="zmdi zmdi-trending-up"></i>
                    </div>
                  </div>
                </div>

                <!-- Возвраты -->
                <div class="col-lg-3 col-md-6 col-sm-12 mb-4">
                  <div class="stat-card returns-card">
                    <div class="stat-icon">
                      <i class="zmdi zmdi-undo"></i>
                    </div>
                    <div class="stat-content">
                      <div class="stat-title">Возвратов</div>
                      <div class="stat-value">{{ stats[2].value }}</div>
                      <div class="stat-period">{{ stats[2].period }}</div>
                    </div>
                    <div class="stat-trend negative">
                      <i class="zmdi zmdi-trending-down"></i>
                    </div>
                  </div>
                </div>

                <!-- Чистая прибыль -->
                <div class="col-lg-3 col-md-6 col-sm-12 mb-4">
                  <div class="stat-card profit-card">
                    <div class="stat-icon">
                      <i class="zmdi zmdi-chart-donut"></i>
                    </div>
                    <div class="stat-content">
                      <div class="stat-title">~ Чистая прибыль</div>
                      <div class="stat-value">{{ stats[3].value }}</div>
                      <div class="stat-period">{{ stats[3].period }}</div>
                    </div>
                    <div class="stat-trend">
                      <i class="zmdi zmdi-trending-up"></i>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Графики -->
              <div class="row charts-row mt-4">
                <!-- Продажи по дням -->
                <div class="col-lg-6 col-md-12 mb-4">
                  <div class="chart-card">
                    <div class="chart-header">
                      <h5 class="chart-title">
                        <i class="zmdi zmdi-chart-line mr-2"></i>
                        Продажи по дням
                      </h5>
                      <div class="chart-actions">
                        <span class="badge badge-success">Актуально</span>
                      </div>
                    </div>
                    <div class="chart-body">
                      <LineChart v-if="lineChartData" :chartData="lineChartData" />
                      <div v-else class="no-data">
                        <i class="zmdi zmdi-chart"></i>
                        <p>Нет данных для графика</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Топ категорий -->
                <div class="col-lg-6 col-md-12 mb-4">
                  <div class="chart-card">
                    <div class="chart-header">
                      <h5 class="chart-title">
                        <i class="zmdi zmdi-chart-bar mr-2"></i>
                        Топ категорий
                      </h5>
                      <div class="chart-actions">
                        <span class="badge badge-info">Топ-5</span>
                      </div>
                    </div>
                    <div class="chart-body">
                      <BarChart v-if="barChartData" :chartData="barChartData" />
                      <div v-else class="no-data">
                        <i class="zmdi zmdi-chart"></i>
                        <p>Нет данных для графика</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Выручка по месяцам -->
                <div class="col-lg-12 col-md-12 mb-4">
                  <div class="chart-card">
                    <div class="chart-header">
                      <h5 class="chart-title">
                        <i class="zmdi zmdi-calendar-alt mr-2"></i>
                        Выручка по месяцам
                        <small v-if="monthlyChartData && monthlyChartData.labels && monthlyChartData.labels.length > 0"
                               class="text-muted ml-2">
                          ({{ monthlyChartData.labels.length }} {{ formatMonthWord(monthlyChartData.labels.length) }})
                        </small>
                      </h5>
                    </div>
                    <div class="chart-body">
                      <!-- Показываем график, если есть хоть какие-то данные -->
                      <BarChart v-if="monthlyChartData && monthlyChartData.datasets &&
                       monthlyChartData.datasets[0] &&
                       monthlyChartData.datasets[0].data &&
                       monthlyChartData.datasets[0].data.length > 0"
                                :chartData="monthlyChartData" />
                      <div v-else class="no-data">
                        <i class="zmdi zmdi-chart"></i>
                        <p>Нет данных по месяцам</p>
                      </div>

                      <!-- Отладочная информация -->
<!--                      <div v-if="monthlyChartData && monthlyChartData.labels" class="debug-info mt-3 small text-muted">-->
<!--                        <p>Месяцы: {{ monthlyChartData.labels.join(', ') }}</p>-->
<!--                        <p>Суммы: {{ monthlyChartData.datasets[0].data.map(d => d.toFixed(2)).join(', ') }}</p>-->
<!--                      </div>-->
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Для гостей -->
          <div v-else class="guest-container">
            <div class="container">
              <div>
                <div class="mx-auto max-w-4xl text-center">
                  <p class="text-5xl text-body-emphasis tracking-tight fw-bold m-0 mt-2">
                    <!--                        Доступные тарифные планы-->
                    Добро пожаловать
                  </p>
                </div>
                <p class="m-0 mt-4 text-lg leading-8 text-body-secondary mx-auto max-w-2xl text-center">
                  <!--                    Выберите план, который соответствует вашим потребностям-->
                  Войдите или зарегистрируйтесь для продолжения
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { useStore } from 'vuex'
import apiClient from '@/api/client'  // Импортируй твой apiClient
import Navbar from '@/components/layout/Navbar.vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import LineChart from '@/components/charts/LineChart.vue'
import BarChart from '@/components/charts/BarChart.vue'

export default {
  name: 'Home',
  components: {
    Navbar,
    Sidebar,
    LineChart,
    BarChart
  },
  setup() {
    const store = useStore()
    const loading = ref(false)
    const error = ref('')
    const lineChartData = ref(null)
    const barChartData = ref(null)
    const monthlyChartData = ref(null)
    const lastUpdated = ref('')
    let refreshInterval = null

    const pageTitle = computed(() => {
      return store.getters.isAuthenticated ? 'Dashboard' : ''
    })

    const isGuest = computed(() => !store.getters.isAuthenticated)

    // Вспомогательная функция для правильного склонения слова "месяц"
    const formatMonthWord = (count) => {
      if (count % 10 === 1 && count % 100 !== 11) return 'месяц'
      if (count % 10 >= 2 && count % 10 <= 4 &&
          (count % 100 < 10 || count % 100 >= 20)) return 'месяца'
      return 'месяцев'
    }

    // Реальные данные
    const stats = ref([
      {
        title: 'Кол-во продаж',
        value: '0',
        period: 'за месяц'
      },
      {
        title: 'Сумма к перечислению',
        value: '₽ 0.00',
        period: 'за месяц'
      },
      {
        title: 'Возвратов',
        value: '0',
        period: 'за месяц'
      },
      {
        title: '~ Чистая прибыль',
        value: '₽ 0.00',
        period: 'за месяц'
      }
    ])

    // Форматирование чисел
    const formatCurrency = (value) => {
      if (value === null || value === undefined || value === '') return '₽ 0.00'
      const num = typeof value === 'string' ? parseFloat(value) : value
      return `₽ ${num.toLocaleString('ru-RU', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
    }

    const formatNumber = (value) => {
      if (value === null || value === undefined || value === '') return '0'
      const num = typeof value === 'string' ? parseFloat(value) : value
      return num.toLocaleString('ru-RU')
    }

    // Форматирование времени
    const formatTime = () => {
      const now = new Date()
      return now.toLocaleTimeString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }
// Анимация обновления времени
    const animateUpdate = () => {
      const timeElement = document.querySelector('.last-updated');
      if (timeElement) {
        timeElement.classList.add('updating');
        setTimeout(() => {
          timeElement.classList.remove('updating');
        }, 500);
      }
    }
    // Подготовка данных для графиков
    const prepareChartData = (chartsData) => {
      // Данные для линейного графика
      if (chartsData.line_chart) {
        const lineData = chartsData.line_chart
        lineChartData.value = {
          labels: lineData.labels || [],
          datasets: [
            {
              label: 'Выручка (руб.)',
              data: lineData.revenue || [],
              borderColor: '#4CAF50',
              backgroundColor: 'rgba(76, 175, 80, 0.1)',
              borderWidth: 2,
              tension: 0.4,
              fill: true
            },
            {
              label: 'Кол-во продаж',
              data: lineData.sales || [],
              borderColor: '#2196F3',
              backgroundColor: 'rgba(33, 150, 243, 0.1)',
              borderWidth: 2,
              tension: 0.4,
              fill: true,
              yAxisID: 'y1'
            }
          ]
        }
      }

      // Данные для столбчатого графика (топ категорий)
      if (chartsData.bar_chart) {
        const barData = chartsData.bar_chart
        barChartData.value = {
          labels: barData.labels || [],
          datasets: [{
            label: 'Выручка по категориям (руб.)',
            data: barData.data || [],
            backgroundColor: [
              'rgba(255, 99, 132, 0.7)',
              'rgba(54, 162, 235, 0.7)',
              'rgba(255, 206, 86, 0.7)',
              'rgba(75, 192, 192, 0.7)',
              'rgba(153, 102, 255, 0.7)'
            ],
            borderColor: [
              'rgba(255, 99, 132, 1)',
              'rgba(54, 162, 235, 1)',
              'rgba(255, 206, 86, 1)',
              'rgba(75, 192, 192, 1)',
              'rgba(153, 102, 255, 1)'
            ],
            borderWidth: 1,
            borderRadius: 4
          }]
        }
      }

      // Данные для графика по месяцам
      if (chartsData.monthly_revenue) {
        const monthlyData = chartsData.monthly_revenue

        // Добавим отладку
        console.log('Monthly revenue data:', monthlyData)
        console.log('Labels:', monthlyData.labels)
        console.log('Data:', monthlyData.data)
        console.log('Has labels?', monthlyData.labels && monthlyData.labels.length > 0)

        if (monthlyData.labels && monthlyData.labels.length > 0) {
          monthlyChartData.value = {
            labels: monthlyData.labels,
            datasets: [{
              label: 'Выручка по месяцам (руб.)',
              data: monthlyData.data || [],
              backgroundColor: 'rgba(63, 81, 181, 0.2)',
              borderColor: '#3F51B5',
              borderWidth: 2,
              borderRadius: 6
            }]
          }
          console.log('Monthly chart data set:', monthlyChartData.value)
        } else {
          console.log('No monthly labels or empty labels array')
          monthlyChartData.value = null
        }
      } else {
        console.log('No monthly_revenue in chartsData')
        monthlyChartData.value = null
      }
    }

// Загрузка данных с сервера
    const fetchDashboardData = async (silent = false) => {
      if (!store.getters.isAuthenticated) return

      try {
        if (!silent) {
          loading.value = true
        }
        error.value = ''

        const token = localStorage.getItem('token')
        const response = await apiClient.get('/dashboard/stats', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })

        const data = response.data

        // Обновляем статистику
        if (data.stats) {
          stats.value = [
            {
              title: 'Кол-во продаж',
              value: formatNumber(data.stats.sales_count),
              period: 'за месяц'
            },
            {
              title: 'Сумма к перечислению',
              value: formatCurrency(data.stats.ppvz_for_pay_total),
              period: 'за месяц'
            },
            {
              title: 'Возвратов',
              value: formatNumber(data.stats.returns_count),
              period: 'за месяц'
            },
            {
              title: '~ Чистая прибыль',
              value: formatCurrency(data.stats.net_profit),
              period: 'за месяц'
            }
          ]
        }

        // Подготавливаем данные для графиков
        if (data.charts) {
          prepareChartData(data.charts)
        } else {
          lineChartData.value = null
          barChartData.value = null
          monthlyChartData.value = null
        }

        // Проверим отдельно monthly_revenue если она на верхнем уровне
        if (data.monthly_revenue) {
          const monthlyData = data.monthly_revenue
          if (monthlyData.labels && monthlyData.labels.length > 0) {
            monthlyChartData.value = {
              labels: monthlyData.labels,
              datasets: [{
                label: 'Выручка по месяцам (руб.)',
                data: monthlyData.data || [],
                backgroundColor: 'rgba(63, 81, 181, 0.2)',
                borderColor: '#3F51B5',
                borderWidth: 2,
                borderRadius: 6
              }]
            }
          }
        }

        // Обновляем время последнего обновления
        lastUpdated.value = formatTime()

      } catch (err) {
        console.error('Error fetching dashboard data:', err)

        if (!silent) {
          if (err.response?.status === 404) {
            error.value = 'Нет данных за выбранный период'
          } else if (err.response?.status === 401) {
            error.value = 'Требуется авторизация'
          } else if (err.message === 'Network Error') {
            error.value = 'Ошибка соединения с сервером'
          } else {
            error.value = 'Ошибка загрузки данных'
          }
        }
      } finally {
        if (!silent) {
          loading.value = false
        }
      }
    }

    onMounted(() => {
      // Проверяем токен в localStorage
      const token = localStorage.getItem('token')
      if (token) {
        store.commit('SET_TOKEN', token)
        store.commit('SET_USER', JSON.parse(localStorage.getItem('user') || '{}'))
      }

      // Загружаем данные сразу
      if (!isGuest.value) {
        fetchDashboardData()

        // Устанавливаем интервал обновления каждые 5 секунд
        refreshInterval = setInterval(() => {
          // Загружаем данные в фоновом режиме
          fetchDashboardData(true)

          // Обновляем время с анимацией
          animateUpdate()
          lastUpdated.value = formatTime()

          console.log('Дашборд обновлен:', new Date().toLocaleTimeString())
        }, 60000) // 5 секунд
      }
    })

    // Очищаем интервал при уничтожении компонента
    onUnmounted(() => {
      if (refreshInterval) {
        clearInterval(refreshInterval)
      }
    })

    return {
      pageTitle,
      isGuest,
      stats,
      loading,
      error,
      lineChartData,
      barChartData,
      monthlyChartData,
      lastUpdated,
      formatMonthWord  // Не забудьте добавить в return
    }
  }
}
</script>

<style scoped>
/* Стили для карточек статистики */
.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.12);
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  border-radius: 12px 12px 0 0;
}

.sales-card::before {
  background: linear-gradient(90deg, #4CAF50, #8BC34A);
}

.revenue-card::before {
  background: linear-gradient(90deg, #2196F3, #03A9F4);
}

.returns-card::before {
  background: linear-gradient(90deg, #F44336, #FF9800);
}

.profit-card::before {
  background: linear-gradient(90deg, #9C27B0, #673AB7);
}

.stat-icon {
  font-size: 2.5rem;
  margin-bottom: 15px;
  opacity: 0.9;
}

.sales-card .stat-icon {
  color: #4CAF50;
}

.revenue-card .stat-icon {
  color: #2196F3;
}

.returns-card .stat-icon {
  color: #F44336;
}

.profit-card .stat-icon {
  color: #9C27B0;
}

.stat-content {
  flex: 1;
}

.stat-title {
  font-size: 0.95rem;
  color: #6c757d;
  font-weight: 500;
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 5px;
  line-height: 1.2;
}

.stat-period {
  font-size: 0.85rem;
  color: #95a5a6;
  font-weight: 400;
}

.stat-trend {
  position: absolute;
  top: 20px;
  right: 20px;
  font-size: 1.5rem;
  color: #4CAF50;
}

.stat-trend.negative {
  color: #F44336;
}

/* Стили для карточек с графиками */
.chart-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.chart-header {
  padding: 20px 25px 10px;
  border-bottom: 1px solid #f1f1f1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
  display: flex;
  align-items: center;
}

.chart-title i {
  font-size: 1.3rem;
}

.chart-actions .badge {
  font-size: 0.75rem;
  padding: 4px 10px;
  border-radius: 20px;
}

.chart-body {
  padding: 20px;
  flex: 1;
  min-height: 400px;
}

.no-data {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 350px;
  color: #95a5a6;
}

.no-data i {
  font-size: 4rem;
  margin-bottom: 20px;
  opacity: 0.5;
}

.no-data p {
  font-size: 1rem;
  margin: 0;
}

/* Стили для гостевого экрана */


.guest-content {
  text-align: center;
  color: white;
  max-width: 800px;
  margin: 0 auto;
}

.guest-title {
  margin-bottom: 30px;
}

.guest-icon {
  font-size: 4rem;
  margin-bottom: 20px;
  opacity: 0.9;
}

.guest-content h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 15px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.guest-subtitle {
  font-size: 1.2rem;
  margin-bottom: 40px;
  opacity: 0.9;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

.guest-features {
  display: flex;
  justify-content: center;
  gap: 40px;
  flex-wrap: wrap;
}

.feature {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.feature i {
  font-size: 2.5rem;
  opacity: 0.9;
}

.feature span {
  font-size: 1rem;
  font-weight: 500;
}

/* Адаптивность */
@media (max-width: 992px) {
  .stat-card {
    padding: 20px;
  }

  .stat-value {
    font-size: 1.6rem;
  }

  .chart-body {
    padding: 15px;
  }
}

@media (max-width: 768px) {
  .guest-content h1 {
    font-size: 2rem;
  }

  .guest-subtitle {
    font-size: 1rem;
  }

  .guest-features {
    gap: 20px;
  }

  .feature i {
    font-size: 2rem;
  }
}

@media (max-width: 576px) {
  .stat-card {
    padding: 15px;
  }

  .stat-value {
    font-size: 1.4rem;
  }

  .chart-header {
    padding: 15px;
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
}

/* Стили для уведомлений */
.alert-warning {
  background-color: #fff3cd;
  border-color: #ffeaa7;
  color: #856404;
}

.alert-dismissible .close {
  position: absolute;
  top: 0;
  right: 0;
  padding: 0.75rem 1.25rem;
  color: inherit;
}

.text-muted {
  color: #6c757d !important;
  font-size: 0.85rem;
}

</style>