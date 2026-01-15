<template>
<!--  <div class="wrapper hold-transition sidebar-mini">-->
    <!-- Navbar -->
<!--    <Navbar />-->

    <!-- Main Sidebar Container -->
<!--    <Sidebar />-->

  <!-- Content Header -->
  <div class="content-header">
    <div class="container-fluid">
      <div class="row mb-2">
        <div class="col-sm-6">
          <h1 class="m-0">{{ pageTitle }}</h1>
        </div>

      </div>
    </div>
  </div>

  <!-- Main Content -->
  <div class="content">
    <div class="container-fluid">
      <!-- Для авторизованных -->
      <div v-if="!isGuest">
        <div class="row">
          <!-- Top Report Start -->
          <div class="col-xlg-3 col-md-3 col-12 mb-30" v-for="stat in stats" :key="stat.title">
            <div class="top-report">
              <!-- Head -->
              <div class="head">
                <h4>{{ stat.title }}</h4>
                <a href="#" class="view"><i class="zmdi zmdi-eye"></i></a>
              </div>

              <!-- Content -->
              <div class="content">
                <h5>{{ stat.period }}</h5>
                <h2>{{ stat.value }}</h2>
              </div>
            </div>
          </div><!-- Top Report End -->
        </div>

        <!-- Графики -->
        <div class="row justify-content-around mt-4">
          <div class="col col-lg-4">
            <div class="card">
              <div class="card-body">
                <h5 class="card-title">Линейный график</h5>
                <div class="chart-container" style="height: 400px; width: 100%;">
                  <canvas ref="lineChart"></canvas>
                </div>
              </div>
            </div>
          </div>
          <div class="col-lg-4 col-md-6 col-sm-6">
            <div class="card">
              <div class="card-body">
                <h5 class="card-title">Столбчатый график</h5>
                <div class="chart-container" style="height: 400px; width: 100%;">
                  <canvas ref="barChart"></canvas>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Для гостей -->
      <div v-else id="pricing" class="overflow-hidden py-7 py-sm-8 py-xl-9">
        <div class="container">
          <div>
            <div class="mx-auto max-w-4xl text-center">
              <p class="text-5xl text-body-emphasis tracking-tight fw-bold m-0 mt-2">
                Добро пожаловать
              </p>
            </div>
            <p class="m-0 mt-4 text-lg leading-8 text-body-secondary mx-auto max-w-2xl text-center">
              Войдите или зарегистрируйтесь для продолжения
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>

    <!-- Main Footer -->
<!--    <Footer />-->
<!--  </div>-->
</template>

<script>
import { ref, onMounted, computed, nextTick } from 'vue'
import { useStore } from 'vuex'
import Chart from 'chart.js/auto'
import Navbar from '@/components/layout/Navbar.vue'
import Sidebar from '@/components/layout/Sidebar.vue'
// import Footer from '@/components/layout/Footer.vue'

export default {
  name: 'Home',
  components: {
    Navbar,
    Sidebar,
    // Footer
  },
  setup() {
    const store = useStore()

    const pageTitle = computed(() => {
      return store.getters.isAuthenticated ? 'Dashboard' : ''
    })

    const isGuest = computed(() => !store.getters.isAuthenticated)

    const lineChart = ref(null)
    const barChart = ref(null)
    let lineChartInstance = null
    let barChartInstance = null

    // Заглушки для данных
    const stats = computed(() => [
      {
        title: 'Кол-во продаж',
        value: '1,250',
        period: 'за месяц'
      },
      {
        title: 'Сумма к перечислению',
        value: '₽ 450,000',
        period: 'за месяц'
      },
      {
        title: 'Возвратов',
        value: '24',
        period: 'за месяц'
      },
      {
        title: '~ Чистая прибыль',
        value: '₽ 120,500',
        period: 'за месяц'
      }
    ])

    // Инициализация графиков
    const initCharts = () => {
      if (!lineChart.value || !barChart.value) return

      // Уничтожаем старые графики если есть
      if (lineChartInstance) {
        lineChartInstance.destroy()
      }
      if (barChartInstance) {
        barChartInstance.destroy()
      }

      // Линейный график
      lineChartInstance = new Chart(lineChart.value, {
        type: 'line',
        data: {
          labels: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл'],
          datasets: [{
            label: 'Продажи',
            data: [65, 59, 80, 81, 56, 55, 40],
            fill: false,
            borderColor: 'rgb(75, 192, 192)',
            tension: 0.1
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'top',
            }
          }
        }
      })

      // Столбчатый график
      barChartInstance = new Chart(barChart.value, {
        type: 'bar',
        data: {
          labels: ['Категория 1', 'Категория 2', 'Категория 3', 'Категория 4', 'Категория 5'],
          datasets: [{
            label: 'Количество',
            data: [12, 19, 3, 5, 2],
            backgroundColor: [
              'rgba(255, 99, 132, 0.2)',
              'rgba(54, 162, 235, 0.2)',
              'rgba(255, 206, 86, 0.2)',
              'rgba(75, 192, 192, 0.2)',
              'rgba(153, 102, 255, 0.2)'
            ],
            borderColor: [
              'rgba(255, 99, 132, 1)',
              'rgba(54, 162, 235, 1)',
              'rgba(255, 206, 86, 1)',
              'rgba(75, 192, 192, 1)',
              'rgba(153, 102, 255, 1)'
            ],
            borderWidth: 1
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'top',
            }
          },
          scales: {
            y: {
              beginAtZero: true
            }
          }
        }
      })
    }

    onMounted(() => {
      // Проверяем токен в localStorage
      const token = localStorage.getItem('token')
      if (token) {
        store.commit('SET_TOKEN', token)
        store.commit('SET_USER', JSON.parse(localStorage.getItem('user') || '{}'))
      }

      // Инициализируем графики после отрисовки DOM
      if (!isGuest.value) {
        nextTick(() => {
          setTimeout(initCharts, 100) // Небольшая задержка для гарантии отрисовки DOM
        })
      }
    })

    return {
      pageTitle,
      isGuest,
      stats,
      lineChart,
      barChart
    }
  }
}
</script>

<style scoped>
/* Пустой - все стили в site.css */
</style>