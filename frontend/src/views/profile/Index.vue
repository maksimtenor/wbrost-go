<script setup>
import { ref, computed, onMounted } from 'vue';
import Navbar from "../../components/layout/Navbar.vue";
import Sidebar from "../../components/layout/Sidebar.vue";

// Реактивные данные формы
const formData = ref({
  name: '',
  username: '',
  password: '',
  email: '',
  phone: '',
  taxes: '',
  wb_key: ''
});

const loading = ref(false);
const message = ref({ text: '', type: '' });

// Получаем данные пользователя при загрузке
const loadUserData = async () => {
  try {
    const token = localStorage.getItem('token');
    const response = await fetch('http://localhost:8080/api/auth/me', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const userData = await response.json();
      // Заполняем форму данными пользователя
      formData.value = {
        name: userData.name || '',
        username: userData.username || '',
        password: '', // Пароль всегда пустой
        email: userData.email || '',
        phone: userData.phone || '',
        taxes: userData.taxes || '',
        wb_key: userData.wb_key || ''
      };
    }
  } catch (error) {
    console.error('Ошибка загрузки данных:', error);
  }
};

// Сохранение формы
const saveProfile = async () => {
  loading.value = true;
  message.value = { text: '', type: '' };

  try {
    const token = localStorage.getItem('token');

    const response = await fetch('http://localhost:8080/api/profile/update', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData.value)
    });

    const data = await response.json();

    if (response.ok) {
      message.value = {
        text: 'Данные успешно сохранены!',
        type: 'success'
      };

      // Обновляем данные пользователя
      await loadUserData();

      // Автоочистка сообщения через 3 секунды
      setTimeout(() => {
        message.value = { text: '', type: '' };
      }, 3000);

    } else {
      message.value = {
        text: data.error || 'Ошибка сохранения',
        type: 'error'
      };
    }

  } catch (error) {
    console.error('Ошибка:', error);
    message.value = {
      text: 'Ошибка сети или сервера',
      type: 'error'
    };
  } finally {
    loading.value = false;
  }
};

// PRO статус (из вашего computed)
const proBadgeClass = computed(() => {
  // Здесь нужен доступ к данным пользователя, можно тоже загружать
  return 'color: red'; // Заглушка
});

const proAccountText = computed(() => {
  return 'Неактивирован'; // Заглушка
});

// Загружаем данные при монтировании
onMounted(() => {
  loadUserData();
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
            <h1 class="m-0">Личный кабинет</h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <!-- Сообщение об успехе/ошибке -->
      <div v-if="message.text"
           :class="['alert', message.type === 'success' ? 'alert-success' : 'alert-danger']"
           role="alert">
        {{ message.text }}
      </div>

      <!-- Форма с @submit.prevent -->
      <form id="profile-form" @submit.prevent="saveProfile">
        <div class="form-group field-user-name">
          <label class="col-lg-3 col-form-label mr-lg-3" for="user-name">Имя:</label>
          <input type="text" id="user-name"
                 class="col-lg-3 form-control"
                 v-model="formData.name"
                 placeholder="Введите"
                 style="width: 400px">
        </div>

        <div class="form-group field-user-username">
          <label class="col-lg-3 col-form-label mr-lg-3" for="user-username">Логин:</label>
          <input type="text" id="user-username"
                 class="col-lg-3 form-control"
                 v-model="formData.username"
                 placeholder="Введите"
                 style="width: 400px">
        </div>

        <div class="form-group field-user-password">
          <label class="col-lg-3 col-form-label mr-lg-3" for="user-password">Смена пароля:</label>
          <input type="password" id="user-password"
                 class="col-lg-3 form-control"
                 v-model="formData.password"
                 placeholder="Введите новый (оставьте пустым, если не меняете)"
                 style="width: 400px">
        </div>

        <div class="form-group field-user-email">
          <label class="col-lg-3 col-form-label mr-lg-3" for="user-email">Email:</label>
          <input type="email" id="user-email"
                 class="col-lg-3 form-control"
                 v-model="formData.email"
                 placeholder="Введите"
                 style="width: 400px">
        </div>

        <div class="form-group field-user-phone">
          <label class="col-lg-3 col-form-label mr-lg-3" for="user-phone">Номер телефона:</label>
          <input type="tel" id="user-phone"
                 class="col-lg-3 form-control"
                 v-model="formData.phone"
                 placeholder="Введите"
                 style="width: 400px">
        </div>

        <div class="mb-3">
          <label class="form-label">PRO-аккаунт:</label>
          <span class="form-control" style="width: 400px;" :style="proBadgeClass">
            {{ proAccountText }}
          </span>
        </div>

        <div class="form-group field-user-taxes">
          <label class="col-lg-3 col-form-label mr-lg-3" for="user-taxes">Налоговый %:</label>
          <input type="number" id="user-taxes"
                 class="col-lg-3 form-control"
                 v-model.number="formData.taxes"
                 placeholder="Введите"
                 min="0" max="100"
                 style="width: 400px">
        </div>

        <div class="form-group field-user-wb_key">
          <label class="col-lg-3 col-form-label mr-lg-3" for="user-wb_key">WB API KEY:</label>
          <input type="text" id="user-wb_key"
                 class="col-lg-3 form-control"
                 v-model="formData.wb_key"
                 placeholder="Введите токен Wildberries"
                 style="width: 400px">
        </div>

        <div class="form-group">
          <div>
            <button type="submit"
                    class="btn btn-primary form-button"
                    :disabled="loading">
              {{ loading ? 'Сохранение...' : 'Сохранить' }}
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.form-group {
  margin-bottom: 1rem;
}

.alert {
  padding: 0.75rem 1.25rem;
  margin-bottom: 1rem;
  border: 1px solid transparent;
  border-radius: 0.25rem;
}

.alert-success {
  color: #155724;
  background-color: #d4edda;
  border-color: #c3e6cb;
}

.alert-danger {
  color: #721c24;
  background-color: #f8d7da;
  border-color: #f5c6cb;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>