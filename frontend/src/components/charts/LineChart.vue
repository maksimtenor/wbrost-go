<template>
  <div class="chart-container">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import Chart from 'chart.js/auto'

export default {
  name: 'LineChart',
  props: {
    chartData: {
      type: Object,
      default: () => ({})
    }
  },
  setup(props) {
    const chartCanvas = ref(null)
    let chartInstance = null

    const renderChart = () => {
      if (!chartCanvas.value) return

      // Если уже есть график - уничтожаем
      if (chartInstance) {
        chartInstance.destroy()
      }

      // Конфигурация из Yii2 ChartJs виджета
      const config = {
        type: 'line',
        data: props.chartData,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'top',
            }
          }
        }
      }

      chartInstance = new Chart(chartCanvas.value, config)
    }

    onMounted(renderChart)

    // Перерисовываем график при изменении данных
    watch(() => props.chartData, renderChart)

    return {
      chartCanvas
    }
  }
}
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 400px;
  width: 100%;
}
</style>