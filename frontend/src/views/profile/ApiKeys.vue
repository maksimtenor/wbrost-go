<script setup>
import Sidebar from "../../components/layout/Sidebar.vue";
import Navbar from "../../components/layout/Navbar.vue";
import { onMounted, ref } from 'vue';

const loading = ref(false);
const wbStatus = ref({
  active: null,     // null = еще не проверяли
  message: '',
  lastCheck: null
});

const checkApiStatus = async () => {
  loading.value = true;
  wbStatus.value.message = 'Проверка...';

  try {
    const token = localStorage.getItem('token');
    if (!token) {
      wbStatus.value = {
        active: false,
        message: 'Требуется авторизация',
        lastCheck: new Date()
      };
      return;
    }

    const response = await fetch('http://localhost:8080/api/profile/apikeys/status', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();

    wbStatus.value = {
      active: data.wildberries.active,
      message: data.wildberries.message,
      lastCheck: new Date()
    };

  } catch (error) {
    console.error('Ошибка проверки:', error);
    wbStatus.value = {
      active: false,
      message: `Ошибка: ${error.message}`,
      lastCheck: new Date()
    };
  } finally {
    loading.value = false;
  }
};

// Форматируем время для отображения
const formatTime = (date) => {
  if (!date) return '';
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

onMounted(() => {
  checkApiStatus();
});
</script>

<template>
  <Navbar />
  <Sidebar />

  <div class="content-wrapper">
    <div class="content-header">
      <div class="container-fluid">
        <div class="row mb-2">
          <div class="col-sm-6">
            <h1 class="m-0">Мои API ключи</h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <table class="table table-striped table-bordered detail-view">
        <tbody>
        <tr>
          <th width="30%">Wildberries</th>
          <td>
            <!-- Статус с иконкой -->
            <span v-if="wbStatus.active === null" style="color: #ffc107;">
                <i class="fas fa-spinner fa-spin"></i> {{ wbStatus.message || 'Проверка...' }}
              </span>
            <span v-else-if="wbStatus.active" style="color:#008000;">
                <i class="fas fa-check-circle"></i> Активен
              </span>
            <span v-else style="color:#dc3545;">
                <i class="fas fa-times-circle"></i> Недействителен
              </span>

            <!-- Кнопка перепроверки и время -->
            <div class="mt-2 d-flex align-items-center">
              <button
                  @click="checkApiStatus"
                  :disabled="loading"
                  class="btn btn-xs btn-outline-secondary mr-2"
                  title="Проверить токен еще раз">
                <i v-if="loading" class="fas fa-spinner fa-spin"></i>
                <i v-else class="fas fa-redo"></i>
                {{ loading ? ' Проверка...' : ' Перепроверить' }}
              </button>

              <small v-if="wbStatus.lastCheck" class="text-muted">
                <i class="far fa-clock"></i>
                Последняя проверка: {{ formatTime(wbStatus.lastCheck) }}
              </small>
            </div>

            <!-- Сообщение -->
            <div v-if="wbStatus.message" class="mt-1">
              <small class="text-muted">{{ wbStatus.message }}</small>
            </div>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.btn {
  transition: all 0.2s;
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-xs {
  padding: 0.15rem 0.4rem;
  font-size: 0.75rem;
  line-height: 1;
  border-radius: 0.2rem;
  vertical-align: middle;
}

.mt-1 { margin-top: 0.25rem; }
.mt-2 { margin-top: 0.5rem; }
.mr-2 { margin-right: 0.5rem; }

.d-flex {
  display: flex;
}

.align-items-center {
  align-items: center;
}

/* Анимация для спиннера */
.fa-spin {
  animation: fa-spin 1s infinite linear;
}

@keyframes fa-spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>