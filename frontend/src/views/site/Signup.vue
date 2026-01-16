<template>
  <div class="hold-transition signup-page">
    <div class="signup-box">
      <div class="signup-logo">
        <router-link to="/">
          <img src="/img/wbrost-logo_2.png" alt="logo" height="32">
        </router-link>
      </div>

      <div class="card">
        <div class="card-body signup-card-body">
          <p class="signup-box-msg">Регистрация</p>

          <!-- Общая ошибка -->
          <div v-if="errors.general" class="alert alert-danger alert-dismissible fade show" role="alert">
            {{ errors.general }}
            <button type="button" class="close" @click="clearError('general')">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>

          <form @submit.prevent="handleSignup" novalidate>
            <div class="form-group has-feedback" :class="{ 'has-error': errors.name }">
              <input
                  type="text"
                  class="form-control"
                  placeholder="Имя"
                  v-model="form.name"
                  :disabled="loading"
                  @input="clearError('name')"
              />
              <span class="fa fa-user form-control-feedback"></span>
              <span class="help-block-error" v-if="errors.name">{{ errors.name }}</span>
            </div>

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

            <div class="form-group has-feedback" :class="{ 'has-error': errors.email }">
              <input
                  type="email"
                  class="form-control"
                  placeholder="Email"
                  v-model="form.email"
                  :disabled="loading"
                  @input="clearError('email')"
              />
              <span class="fa fa-envelope form-control-feedback"></span>
              <span class="help-block-error" v-if="errors.email">{{ errors.email }}</span>
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
              <div class="col-6">
                <button
                    type="submit"
                    class="btn btn-primary btn-block"
                    :disabled="loading"
                >
                  <span v-if="loading">
                    <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                    Регистрация...
                  </span>
                  <span v-else>Регистрация</span>
                </button>
              </div>
              <div class="col-6">
                <router-link to="/login" class="btn btn-secondary btn-block">
                  Вход
                </router-link>
              </div>
            </div>

            <div class="row-cols-1 mt-4">
              <div id="privacy_and_terms" class="row ml-0 mr-0" :class="{ 'has-error': errors.agreeTerms }">
                <input
                    id="tems_and_privacy_checkbox"
                    type="checkbox"
                    v-model="form.agreeTerms"
                    :disabled="loading"
                    @change="clearError('agreeTerms')"
                >
                <label for="tems_and_privacy_checkbox" style="margin-left: 0.5rem;">
                  Я согласен(на) с
                  <a href="#" @click.prevent="openModal('terms')">условиями</a> и
                  <a href="#" @click.prevent="openModal('privacy')">политикой конфиденциальности</a>.
                </label>
                <span class="help-block-error" v-if="errors.agreeTerms">{{ errors.agreeTerms }}</span>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>

  <!-- Модальное окно для условий и политики -->
  <div v-if="showModal" class="modal-backdrop fade show"></div>
  <div v-if="showModal" class="modal fade show d-block" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-lg" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            {{ modalType === 'terms' ? 'Условия использования' : 'Политика конфиденциальности' }}
          </h5>
          <button type="button" class="close" @click="closeModal">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body" style="max-height: 400px; overflow-y: auto;">
          <div v-if="modalType === 'terms'">
            <h5>1. Общие положения</h5>
            <p>1.1. Используя данный сервис, Вы соглашаетесь с настоящими Условиями использования.</p>
            <p>1.2. Сервис предназначен для законного использования.</p>

            <h5>2. Регистрация и учетная запись</h5>
            <p>2.1. Вы должны быть старше 18 лет для регистрации.</p>
            <p>2.2. Вы несете ответственность за сохранность своих учетных данных.</p>

            <h5>3. Обязанности пользователя</h5>
            <p>3.1. Не передавать свои учетные данные третьим лицам.</p>
            <p>3.2. Не использовать сервис для незаконной деятельности.</p>

            <h5>4. Права администрации</h5>
            <p>4.1. Администрация оставляет за собой право изменять условия.</p>
            <p>4.2. Право приостановить доступ при нарушении условий.</p>
          </div>
          <div v-else>
            <h5>1. Сбор информации</h5>
            <p>1.1. Мы собираем информацию, которую Вы предоставляете при регистрации.</p>
            <p>1.2. Данные используются для предоставления услуг сервиса.</p>

            <h5>2. Использование информации</h5>
            <p>2.1. Ваши данные не передаются третьим лицам без согласия.</p>
            <p>2.2. Мы используем данные для улучшения качества сервиса.</p>

            <h5>3. Защита данных</h5>
            <p>3.1. Мы применяем меры для защиты Ваших данных.</p>
            <p>3.2. Пароли хранятся в зашифрованном виде.</p>

            <h5>4. Cookies</h5>
            <p>4.1. Мы используем cookies для работы сервиса.</p>
            <p>4.2. Вы можете отключить cookies в настройках браузера.</p>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal">
            Закрыть
          </button>
          <button v-if="!form.agreeTerms" type="button" class="btn btn-primary" @click="acceptTerms">
            Принимаю условия
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { authAPI } from '@/api/auth'

export default {
  name: 'Signup',
  setup() {
    const router = useRouter()
    const store = useStore()

    const form = ref({
      name: '',
      username: '',
      email: '',
      password: '',
      agreeTerms: false
    })

    const errors = ref({})
    const loading = ref(false)
    const showModal = ref(false)
    const modalType = ref('terms') // 'terms' или 'privacy'

    // Функция для очистки ошибки
    const clearError = (field) => {
      if (errors.value[field]) {
        const newErrors = { ...errors.value }
        delete newErrors[field]
        errors.value = newErrors
      }
    }

    const openModal = (type) => {
      modalType.value = type
      showModal.value = true
    }

    const closeModal = () => {
      showModal.value = false
    }

    const acceptTerms = () => {
      form.value.agreeTerms = true
      closeModal()
      clearError('agreeTerms')
    }

    const validateForm = () => {
      const newErrors = {}
      let isValid = true

      if (!form.value.name.trim()) {
        newErrors.name = 'Введите имя'
        isValid = false
      } else if (form.value.name.length < 2) {
        newErrors.name = 'Имя должно быть не менее 2 символов'
        isValid = false
      }

      if (!form.value.username.trim()) {
        newErrors.username = 'Введите имя пользователя'
        isValid = false
      } else if (form.value.username.length < 3) {
        newErrors.username = 'Имя пользователя должно быть не менее 3 символов'
        isValid = false
      } else if (!/^[a-zA-Z0-9_]+$/.test(form.value.username)) {
        newErrors.username = 'Можно использовать только латинские буквы, цифры и подчеркивание'
        isValid = false
      }

      if (!form.value.email.trim()) {
        newErrors.email = 'Введите email'
        isValid = false
      } else if (!/^\S+@\S+\.\S+$/.test(form.value.email)) {
        newErrors.email = 'Введите корректный email'
        isValid = false
      }

      if (!form.value.password) {
        newErrors.password = 'Введите пароль'
        isValid = false
      } else if (form.value.password.length < 6) {
        newErrors.password = 'Пароль должен быть не менее 6 символов'
        isValid = false
      }

      if (!form.value.agreeTerms) {
        newErrors.agreeTerms = 'Необходимо согласиться с условиями'
        isValid = false
      }

      if (!isValid) {
        errors.value = newErrors
      }

      return isValid
    }

    const handleSignup = async () => {
      // Сбрасываем ошибки
      errors.value = {}
      loading.value = true

      // Фронтенд валидация
      if (!validateForm()) {
        loading.value = false
        return
      }

      try {
        // Отправляем запрос на бэкенд
        const result = await authAPI.signup({
          name: form.value.name,
          username: form.value.username,
          email: form.value.email,
          password: form.value.password
        })

        console.log('Signup result:', result) // Для отладки

        if (result.success) {
          // Сохраняем данные авторизации
          authAPI.setAuthData(result.data.token, result.data.user)

          // Сохраняем в хранилище Vuex
          store.commit('SET_TOKEN', result.data.token)
          store.commit('SET_USER', result.data.user)

          // Перенаправляем на главную
          router.push('/')

        } else {
          // Обработка ошибок с бэкенда
          if (result.validationErrors) {
            // Ошибки валидации от сервера
            errors.value = result.validationErrors
          } else {
            const errorMsg = result.error || 'Ошибка регистрации'

            if (errorMsg.toLowerCase().includes('already exists') ||
                errorMsg.toLowerCase().includes('уже существует')) {
              if (errorMsg.toLowerCase().includes('username')) {
                errors.value.username = 'Это имя пользователя уже занято'
              } else if (errorMsg.toLowerCase().includes('email')) {
                errors.value.email = 'Этот email уже зарегистрирован'
              } else {
                errors.value.general = errorMsg
              }
            } else {
              errors.value.general = errorMsg
            }
          }
        }
      } catch (error) {
        console.error('Signup error:', error)

        if (error.response) {
          // Ошибка от сервера
          const status = error.response.status
          const data = error.response.data

          if (status === 400 && data?.errors) {
            // Ошибки валидации от сервера
            errors.value = data.errors
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
      showModal,
      modalType,
      handleSignup,
      openModal,
      closeModal,
      acceptTerms,
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

.signup-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f4f6f9;
}

.signup-box {
  width: 400px;
}

#privacy_and_terms.has-error label {
  color: #dc3545;
}

#privacy_and_terms.has-error a {
  color: #dc3545;
  text-decoration: underline;
}

.modal-backdrop {
  opacity: 0.5;
}
</style>