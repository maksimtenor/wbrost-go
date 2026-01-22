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
/* Общие стили для обеих страниц */
.hold-transition {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f4f6f9;
}

.login-page, .signup-page {
  width: 100%;
}

.login-box, .signup-box {
  width: 100%;
  max-width: 360px;
  margin: 0 auto;
  padding: 20px;
}

.signup-box {
  max-width: 400px;
}

.login-logo, .signup-logo {
  text-align: center;
  margin-bottom: 20px;
}

.login-logo a, .signup-logo a {
  display: inline-block;
}

.card {
  border: 1px solid #d2d6de;
  border-radius: 3px;
  box-shadow: 0 1px 1px rgba(0,0,0,.1);
  background-color: #fff;
}

.login-card-body, .signup-card-body {
  padding: 20px;
}

.login-box-msg, .signup-box-msg {
  margin: 0;
  text-align: center;
  padding: 0 20px 20px;
  font-size: 18px;
  color: #666;
}

/* Стили для полей формы */
.form-group {
  margin-bottom: 1rem;
  position: relative;
}

.has-feedback {
  position: relative;
}

.form-control {
  height: 38px;
  padding: 6px 12px;
  font-size: 14px;
  line-height: 1.42857143;
  color: #555;
  background-color: #fff;
  background-image: none;
  border: 1px solid #d2d6de;
  border-radius: 3px;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075);
  transition: border-color ease-in-out .15s,box-shadow ease-in-out .15s;
  width: 100%;
}

.form-control:focus {
  border-color: #3c8dbc;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(60,141,188,.6);
}

.form-control-feedback {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 2;
  display: block;
  width: 34px;
  height: 34px;
  line-height: 34px;
  text-align: center;
  pointer-events: none;
  color: #777;
}

.has-error .form-control {
  border-color: #dd4b39;
}

.has-error .form-control:focus {
  border-color: #dd4b39;
  box-shadow: 0 0 0 0.2rem rgba(221, 75, 57, 0.25);
}

/* Стили для ошибок */
.help-block-error {
  color: #dd4b39;
  font-size: 12px;
  margin-top: 5px;
  display: block;
}

/* Стили для alert */
.alert {
  margin-bottom: 1rem;
  padding: 0.75rem 1.25rem;
  border-radius: 0.25rem;
  border: 1px solid transparent;
}

.alert-danger {
  color: #721c24;
  background-color: #f8d7da;
  border-color: #f5c6cb;
}

.alert-dismissible {
  padding-right: 3.75rem;
}

.alert-dismissible .close {
  position: absolute;
  top: 0;
  right: 0;
  padding: 0.75rem 1.25rem;
  color: inherit;
  background: transparent;
  border: 0;
  cursor: pointer;
}

/* Стили для кнопок */
.btn {
  display: inline-block;
  font-weight: 400;
  text-align: center;
  white-space: nowrap;
  vertical-align: middle;
  user-select: none;
  border: 1px solid transparent;
  padding: 0.375rem 0.75rem;
  font-size: 1rem;
  line-height: 1.5;
  border-radius: 0.25rem;
  transition: color .15s ease-in-out,background-color .15s ease-in-out,border-color .15s ease-in-out,box-shadow .15s ease-in-out;
  cursor: pointer;
}

.btn-block {
  display: block;
  width: 100%;
}

.btn-primary {
  color: #fff;
  background-color: #3c8dbc;
  border-color: #367fa9;
}

.btn-primary:hover:not(:disabled) {
  background-color: #367fa9;
  border-color: #32769b;
}

.btn-primary:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.btn-secondary {
  color: #fff;
  background-color: #6c757d;
  border-color: #6c757d;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #5a6268;
  border-color: #545b62;
}

/* Стили для чекбоксов */
.icheck-primary {
  display: flex;
  align-items: center;
}

.icheck-primary input[type="checkbox"] {
  margin-right: 5px;
}

.icheck-primary label {
  margin: 0;
  cursor: pointer;
  font-size: 14px;
}

/* Стили для строк */
.row {
  display: flex;
  flex-wrap: wrap;
  margin-right: -15px;
  margin-left: -15px;
}

.col-4, .col-6, .col-8, .col-12 {
  position: relative;
  width: 100%;
  padding-right: 15px;
  padding-left: 15px;
}

.col-4 {
  flex: 0 0 33.333333%;
  max-width: 33.333333%;
}

.col-6 {
  flex: 0 0 50%;
  max-width: 50%;
}

.col-8 {
  flex: 0 0 66.666667%;
  max-width: 66.666667%;
}

.col-12 {
  flex: 0 0 100%;
  max-width: 100%;
}

.mt-3 {
  margin-top: 1rem !important;
}

.mt-4 {
  margin-top: 1.5rem !important;
}

/* Стили для ссылок */
a {
  color: #3c8dbc;
  text-decoration: none;
}

a:hover {
  color: #367fa9;
  text-decoration: underline;
}

/* Стили для спиннера */
.spinner-border-sm {
  width: 1rem;
  height: 1rem;
  border-width: 0.2em;
  margin-right: 0.5rem;
  vertical-align: middle;
}

.spinner-border {
  display: inline-block;
  animation: spinner-border .75s linear infinite;
  border: 0.2em solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
}

@keyframes spinner-border {
  to { transform: rotate(360deg); }
}

/* Стили для модального окна */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1040;
  width: 100vw;
  height: 100vh;
  background-color: #000;
  opacity: 0.5;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1050;
  width: 100%;
  height: 100%;
  overflow: hidden;
  outline: 0;
}

.modal.show {
  display: block;
}

.modal-dialog {
  position: relative;
  width: auto;
  margin: 10% auto;
  pointer-events: none;
  max-width: 800px;
}

.modal-content {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  pointer-events: auto;
  background-color: #fff;
  background-clip: padding-box;
  border: 1px solid rgba(0,0,0,.2);
  border-radius: 0.3rem;
  outline: 0;
}

.modal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1rem;
  border-bottom: 1px solid #dee2e6;
  border-top-left-radius: calc(0.3rem - 1px);
  border-top-right-radius: calc(0.3rem - 1px);
}

.modal-title {
  margin-bottom: 0;
  line-height: 1.5;
}

.modal-body {
  position: relative;
  flex: 1 1 auto;
  padding: 1rem;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 1rem;
  border-top: 1px solid #dee2e6;
  border-bottom-right-radius: calc(0.3rem - 1px);
  border-bottom-left-radius: calc(0.3rem - 1px);
}

.modal-footer > * {
  margin-left: 0.25rem;
}

.close {
  float: right;
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1;
  color: #000;
  text-shadow: 0 1px 0 #fff;
  opacity: .5;
  background: transparent;
  border: 0;
  cursor: pointer;
}

.close:hover {
  opacity: .75;
}

/* Дополнительные стили для страницы регистрации */
#privacy_and_terms {
  display: flex;
  align-items: flex-start;
  margin-top: 10px;
}

#privacy_and_terms input[type="checkbox"] {
  margin-top: 3px;
  margin-right: 5px;
}

#privacy_and_terms label {
  font-size: 14px;
  line-height: 1.4;
  margin: 0;
}

#privacy_and_terms.has-error label {
  color: #dd4b39;
}

#privacy_and_terms.has-error a {
  color: #dd4b39;
  text-decoration: underline;
}

.row-cols-1 {
  flex-direction: column;
}

.ml-0 {
  margin-left: 0 !important;
}

.mr-0 {
  margin-right: 0 !important;
}

/* Адаптивность */
@media (max-width: 576px) {
  .login-box, .signup-box {
    padding: 10px;
    max-width: 100%;
  }

  .col-4, .col-6, .col-8 {
    flex: 0 0 100%;
    max-width: 100%;
    margin-bottom: 10px;
  }

  .modal-dialog {
    margin: 5% 10px;
    max-width: calc(100% - 20px);
  }
}

/* Стили для состояний загрузки */
button:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}

/* Улучшенная анимация появления */
.fade {
  transition: opacity 0.15s linear;
}

.fade.show {
  opacity: 1;
}
</style>