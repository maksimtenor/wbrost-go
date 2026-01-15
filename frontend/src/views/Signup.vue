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

          <form @submit.prevent="handleSignup">
            <div class="form-group has-feedback" :class="{ 'has-error': errors.name }">
              <input
                  type="text"
                  class="form-control"
                  placeholder="Имя"
                  v-model="form.name"
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
              />
              <span class="fa fa-lock form-control-feedback"></span>
              <span class="help-block-error" v-if="errors.password">{{ errors.password }}</span>
            </div>

            <div class="row">
              <div class="col-6">
                <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
                  <span v-if="loading">Регистрация...</span>
                  <span v-else>Регистрация</span>
                </button>
              </div>
              <div class="col-6">
                <router-link to="/login" class="btn btn-secondary btn-block">Вход</router-link>
              </div>
            </div>

            <div class="row-cols-1 mt-4">
              <div id="privacy_and_terms" class="row ml-0 mr-0">
                <input
                    id="tems_and_privacy_checkbox"
                    type="checkbox"
                    v-model="form.agreeTerms"
                    required
                >
                <label for="tems_and_privacy_checkbox">
                  Я согласен(на) с
                  <a href="#" @click.prevent="openTerms">условиями</a> и
                  <a href="#" @click.prevent="openPrivacy">политикой конфиденциальности</a>.
                </label>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'Signup',
  setup() {
    const router = useRouter()

    const form = ref({
      name: '',
      username: '',
      email: '',
      password: '',
      agreeTerms: false
    })

    const errors = ref({})
    const loading = ref(false)

    const handleSignup = async () => {
      errors.value = {}
      loading.value = true

      if (!form.value.agreeTerms) {
        errors.value.agreeTerms = 'Необходимо согласиться с условиями'
        loading.value = false
        return
      }

      try {
        // Временная заглушка
        alert('Регистрация прошла успешно! Войдите в систему.')
        router.push('/login')
      } catch (error) {
        errors.value = { general: 'Ошибка регистрации' }
      } finally {
        loading.value = false
      }
    }

    const openTerms = () => {
      alert('Страница условий будет доступна позже')
    }

    const openPrivacy = () => {
      alert('Страница политики конфиденциальности будет доступна позже')
    }

    return {
      form,
      errors,
      loading,
      handleSignup,
      openTerms,
      openPrivacy
    }
  }
}
</script>