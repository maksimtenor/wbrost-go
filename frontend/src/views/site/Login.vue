<template>
  <div class="hold-transition login-page">
    <div class="login-box">
      <div class="login-logo">
        <router-link to="/">
          <img src="/img/wbrost-logo_2.png" alt="logo" height="32">
        </router-link>
      </div>

      <div class="card">
        <div class="card-body login-card-body">
          <p class="login-box-msg">Авторизуйтесь для продолжения</p>

          <!-- Общая ошибка -->
          <div v-if="errors.general" class="alert alert-danger alert-dismissible fade show" role="alert">
            {{ errors.general }}
            <button type="button" class="close" @click="clearError('general')">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>

          <form @submit.prevent="handleLogin" novalidate>
            <div class="form-group has-feedback" :class="{ 'has-error': errors.username }">
              <input
                  type="text"
                  class="form-control"
                  placeholder="Имя пользователя"
                  v-model="form.username"
                  :disabled="loading"
                  @input="clearError('username')"
              />
              <span class="fa fa-user form-control-feedback"></span>
              <span class="help-block-error" v-if="errors.username">{{ errors.username }}</span>
            </div>

            <div class="form-group has-feedback" :class="{ 'has-error': errors.password }">
              <input
                  type="password"
                  class="form-control"
                  placeholder="Пароль"
                  v-model="form.password"
                  :disabled="loading"
                  @input="clearError('password')"
              />
              <span class="fa fa-lock form-control-feedback"></span>
              <span class="help-block-error" v-if="errors.password">{{ errors.password }}</span>
            </div>

            <div class="row">
              <div class="col-8">
                <div class="icheck-primary">
                  <input
                      type="checkbox"
                      id="remember"
                      v-model="form.rememberMe"
                      :disabled="loading"
                  >
                  <label for="remember">Запомнить меня</label>
                </div>
              </div>
              <div class="col-4">
                <button
                    type="submit"
                    class="btn btn-primary btn-block"
                    :disabled="loading"
                >
                  <span v-if="loading">
                    <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                    Вход...
                  </span>
                  <span v-else>Войти</span>
                </button>
              </div>
              <div class="col-12 mt-3">
                <router-link to="/signup">Регистрация</router-link>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { authAPI } from '@/api/auth'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const store = useStore()

    const form = ref({
      username: '',
      password: '',
      rememberMe: localStorage.getItem('rememberMe') === 'true'
    })

    const errors = ref({})
    const loading = ref(false)

    // Восстанавливаем сохраненный username если есть
    onMounted(() => {
      const savedUsername = localStorage.getItem('savedUsername')
      if (savedUsername) {
        form.value.username = savedUsername
      }
    })

    // Функция для очистки ошибки
    const clearError = (field) => {
      if (errors.value[field]) {
        const newErrors = { ...errors.value }
        delete newErrors[field]
        errors.value = newErrors
      }
    }

    const handleLogin = async () => {
      // Сбрасываем ошибки
      errors.value = {}
      loading.value = true

      // Валидация на фронтенде
      let hasError = false
      const newErrors = {}

      if (!form.value.username.trim()) {
        newErrors.username = 'Введите имя пользователя'
        hasError = true
      }

      if (!form.value.password) {
        newErrors.password = 'Введите пароль'
        hasError = true
      }

      if (hasError) {
        errors.value = newErrors
        loading.value = false
        return
      }

      try {
        // Сохраняем имя пользователя для удобства
        if (form.value.rememberMe) {
          localStorage.setItem('savedUsername', form.value.username)
        }

        // Отправляем запрос на бэкенд
        const result = await authAPI.login({
          username: form.value.username,
          password: form.value.password
        })

        console.log('Login result:', result) // Для отладки

        if (result.success) {
          // Сохраняем данные авторизации
          authAPI.setAuthData(result.data.token, result.data.user)

          // Сохраняем в хранилище Vuex
          store.commit('SET_TOKEN', result.data.token)
          store.commit('SET_USER', result.data.user)

          // Сохраняем настройку "запомнить меня"
          if (form.value.rememberMe) {
            localStorage.setItem('rememberMe', 'true')
          } else {
            localStorage.removeItem('rememberMe')
            localStorage.removeItem('savedUsername')
          }
          // Загружаем свежие данные с сервера
          await store.dispatch('loadUserData')
          // Перенаправляем на главную
          router.push('/')

        } else {
          // Обработка ошибок с бэкенда
          const errorMsg = result.error || 'Ошибка авторизации'

          if (errorMsg.toLowerCase().includes('credentials') ||
              errorMsg.toLowerCase().includes('invalid') ||
              errorMsg.toLowerCase().includes('неверный')) {
            errors.value.general = 'Неверное имя пользователя или пароль'
          } else if (errorMsg.toLowerCase().includes('required')) {
            errors.value.general = 'Заполните все обязательные поля'
          } else {
            errors.value.general = errorMsg
          }
        }
      } catch (error) {
        console.error('Login error:', error)

        if (error.response) {
          // Ошибка от сервера
          const status = error.response.status
          const data = error.response.data

          if (status === 401) {
            errors.value.general = 'Неверное имя пользователя или пароль'
          } else if (status === 400) {
            errors.value.general = data?.error || 'Ошибка в данных запроса'
          } else if (status === 500) {
            errors.value.general = 'Ошибка сервера. Попробуйте позже.'
          } else {
            errors.value.general = `Ошибка ${status}: ${data?.error || 'Неизвестная ошибка'}`
          }
        } else if (error.request) {
          // Запрос был отправлен, но ответа нет
          errors.value.general = 'Сервер не отвечает. Проверьте подключение.'
        } else {
          // Ошибка при настройке запроса
          errors.value.general = 'Ошибка отправки запроса: ' + error.message
        }
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      errors,
      loading,
      handleLogin,
      clearError
    }
  }
}
</script>

<style scoped>
.help-block-error {
  color: #dc3545;
  font-size: 0.875rem;
  margin-top: 0.25rem;
  display: block;
}

.has-error .form-control {
  border-color: #dc3545;
}

.has-error .form-control:focus {
  border-color: #dc3545;
  box-shadow: 0 0 0 0.2rem rgba(220, 53, 69, 0.25);
}

.alert {
  margin-bottom: 1rem;
  padding: 0.75rem 1.25rem;
  border-radius: 0.25rem;
}

.alert-danger {
  color: #721c24;
  background-color: #f8d7da;
  border-color: #f5c6cb;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
}

.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f4f6f9;
}

.login-box {
  width: 360px;
}
</style>