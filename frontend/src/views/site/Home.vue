<template>
  <BaseLayout>

    <template #header>
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
    </template>

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
                  <LineChart v-if="lineChartData" :chartData="lineChartData"/>
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
                  <BarChart v-if="barChartData" :chartData="barChartData"/>
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
                            :chartData="monthlyChartData"/>
                  <div v-else class="no-data">
                    <i class="zmdi zmdi-chart"></i>
                    <p>Нет данных по месяцам</p>
                  </div>
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

  </BaseLayout>
</template>

<script>
import {ref, onMounted, computed, onUnmounted} from 'vue'
import {useStore} from 'vuex'
import apiClient from '@/api/client'  // Импортируй твой apiClient
import Navbar from '@/components/layout/Navbar.vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import LineChart from '@/components/charts/LineChart.vue'
import BarChart from '@/components/charts/BarChart.vue'
import BaseLayout from "@/components/layout/BaseLayout.vue";
import BasicForm from "@/components/layout/forms/BasicForm.vue";
const DASHBOARD_AUTO_REFRESH = import.meta.env.DASHBOARD_AUTO_REFRESH;

export default {
  name: 'Home',
  components: {
    BasicForm,
    BaseLayout,
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
      return `₽ ${num.toLocaleString('ru-RU', {minimumFractionDigits: 2, maximumFractionDigits: 2})}`
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
        }, DASHBOARD_AUTO_REFRESH) // 5 секунд
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
@import '@/assets/css/views/site/home.css';
</style>